// Package phoenix :: tx.go
// Minimal Solana transaction builder, signer, and RPC submitter.
//
// Flow:
//   1. Load ed25519 keypair from Solana CLI JSON file
//   2. Fetch recent blockhash via Solana RPC
//   3. Build legacy transaction message from Phoenix instructions
//   4. Sign with ed25519 private key
//   5. Submit base64-encoded transaction via sendTransaction RPC
package phoenix

import (
	"bytes"
	"context"
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// ── Base58 ────────────────────────────────────────────────────────────

const base58Alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

func base58Decode(s string) ([]byte, error) {
	alphabet := []byte(base58Alphabet)
	lookup := [256]int{}
	for i := range lookup {
		lookup[i] = -1
	}
	for i, c := range alphabet {
		lookup[c] = i
	}

	result := make([]byte, 0, len(s))
	for _, c := range []byte(s) {
		carry := lookup[c]
		if carry < 0 {
			return nil, fmt.Errorf("invalid base58 char %q", c)
		}
		for j := len(result) - 1; j >= 0; j-- {
			carry += 58 * int(result[j])
			result[j] = byte(carry & 0xff)
			carry >>= 8
		}
		for carry > 0 {
			result = append([]byte{byte(carry & 0xff)}, result...)
			carry >>= 8
		}
	}
	// Leading zeros
	for _, c := range []byte(s) {
		if c != alphabet[0] {
			break
		}
		result = append([]byte{0}, result...)
	}
	return result, nil
}

func base58Encode(b []byte) string {
	alphabet := []byte(base58Alphabet)
	result := []byte{}
	for _, byt := range b {
		carry := int(byt)
		for j := len(result) - 1; j >= 0; j-- {
			carry += 256 * int(result[j])
			result[j] = byte(carry % 58)
			carry /= 58
		}
		for carry > 0 {
			result = append([]byte{byte(carry % 58)}, result...)
			carry /= 58
		}
	}
	// Leading zero bytes → '1'
	for _, byt := range b {
		if byt != 0 {
			break
		}
		result = append([]byte{0}, result...)
	}
	encoded := make([]byte, len(result))
	for i, idx := range result {
		encoded[i] = alphabet[idx]
	}
	return string(encoded)
}

// ── Keypair loading ───────────────────────────────────────────────────

// Keypair holds an ed25519 private key and its corresponding public key.
type Keypair struct {
	PrivateKey ed25519.PrivateKey
	PublicKey  [32]byte
}

// LoadKeypair reads a Solana CLI keypair file (JSON array of 64 bytes).
func LoadKeypair(path string) (*Keypair, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read keypair %s: %w", path, err)
	}
	var bytes64 []byte
	if err := json.Unmarshal(data, &bytes64); err != nil {
		// Try as array of ints (Solana CLI format)
		var ints []int
		if err2 := json.Unmarshal(data, &ints); err2 != nil {
			return nil, fmt.Errorf("parse keypair: %w", err)
		}
		bytes64 = make([]byte, len(ints))
		for i, v := range ints {
			bytes64[i] = byte(v)
		}
	}
	if len(bytes64) != 64 {
		return nil, fmt.Errorf("keypair must be 64 bytes, got %d", len(bytes64))
	}
	priv := ed25519.NewKeyFromSeed(bytes64[:32])
	var pub [32]byte
	copy(pub[:], bytes64[32:])
	return &Keypair{PrivateKey: priv, PublicKey: pub}, nil
}

// Pubkey returns the base58-encoded public key string.
func (kp *Keypair) Pubkey() string {
	return base58Encode(kp.PublicKey[:])
}

// ── compact-u16 encoding (Solana wire format) ─────────────────────────

func compactU16(n int) []byte {
	switch {
	case n < 0x80:
		return []byte{byte(n)}
	case n < 0x4000:
		return []byte{byte(n&0x7f | 0x80), byte(n >> 7)}
	default:
		return []byte{byte(n&0x7f | 0x80), byte((n>>7)&0x7f | 0x80), byte(n >> 14)}
	}
}

// ── Solana transaction builder ────────────────────────────────────────

// BuildAndSign constructs a signed legacy Solana transaction from Phoenix instructions.
// blockhash is the recent blockhash (base58).
func BuildAndSign(kp *Keypair, instructions []Instruction, blockhash string) ([]byte, error) {
	blockhashBytes, err := base58Decode(blockhash)
	if err != nil {
		return nil, fmt.Errorf("decode blockhash: %w", err)
	}
	if len(blockhashBytes) != 32 {
		return nil, fmt.Errorf("blockhash must be 32 bytes, got %d", len(blockhashBytes))
	}

	// Collect all unique account keys; fee payer (signer) must be first.
	type accountEntry struct {
		pubkey     [32]byte
		isSigner   bool
		isWritable bool
	}
	keyIndex := map[[32]byte]int{}
	var accounts []accountEntry

	addAccount := func(pubkeyStr string, isSigner, isWritable bool) (int, error) {
		bs, err := base58Decode(pubkeyStr)
		if err != nil {
			return 0, fmt.Errorf("decode pubkey %s: %w", pubkeyStr, err)
		}
		var pk [32]byte
		if len(bs) != 32 {
			return 0, fmt.Errorf("pubkey must be 32 bytes, got %d for %s", len(bs), pubkeyStr)
		}
		copy(pk[:], bs)
		if idx, ok := keyIndex[pk]; ok {
			// Upgrade writable/signer flags
			if isSigner {
				accounts[idx].isSigner = true
			}
			if isWritable {
				accounts[idx].isWritable = true
			}
			return idx, nil
		}
		idx := len(accounts)
		keyIndex[pk] = idx
		accounts = append(accounts, accountEntry{pk, isSigner, isWritable})
		return idx, nil
	}

	// Insert fee payer (our wallet) first as signer + writable
	feePayerIdx, err := addAccount(kp.Pubkey(), true, true)
	if err != nil {
		return nil, err
	}
	_ = feePayerIdx

	// Collect all accounts from instructions
	type builtIx struct {
		programIdx int
		accountIdxs []byte
		data        []byte
	}
	builtIxs := make([]builtIx, 0, len(instructions))
	for _, ix := range instructions {
		progIdx, err := addAccount(ix.ProgramID, false, false)
		if err != nil {
			return nil, err
		}
		acctIdxs := make([]byte, len(ix.Keys))
		for j, k := range ix.Keys {
			idx, err := addAccount(k.Pubkey, k.IsSigner, k.IsWritable)
			if err != nil {
				return nil, err
			}
			acctIdxs[j] = byte(idx)
		}
		builtIxs = append(builtIxs, builtIx{progIdx, acctIdxs, ix.DataBytes()})
	}

	// Sort accounts: writable signers, readonly signers, writable non-signers, readonly non-signers
	// For simplicity: signers first, then non-signers. Fee payer is always index 0.
	numRequiredSigs := 0
	numReadonlySigned := 0
	numReadonlyUnsigned := 0
	for _, a := range accounts {
		if a.isSigner {
			numRequiredSigs++
			if !a.isWritable {
				numReadonlySigned++
			}
		} else if !a.isWritable {
			numReadonlyUnsigned++
		}
	}

	// ── Build message ─────────────────────────────────────────────────
	var msg bytes.Buffer

	// Header
	msg.WriteByte(byte(numRequiredSigs))
	msg.WriteByte(byte(numReadonlySigned))
	msg.WriteByte(byte(numReadonlyUnsigned))

	// Account keys
	msg.Write(compactU16(len(accounts)))
	for _, a := range accounts {
		msg.Write(a.pubkey[:])
	}

	// Recent blockhash
	msg.Write(blockhashBytes)

	// Instructions
	msg.Write(compactU16(len(builtIxs)))
	for _, ix := range builtIxs {
		msg.WriteByte(byte(ix.programIdx))
		msg.Write(compactU16(len(ix.accountIdxs)))
		msg.Write(ix.accountIdxs)
		msg.Write(compactU16(len(ix.data)))
		msg.Write(ix.data)
	}

	msgBytes := msg.Bytes()

	// ── Sign ──────────────────────────────────────────────────────────
	sig := ed25519.Sign(kp.PrivateKey, msgBytes)

	// ── Serialize transaction ─────────────────────────────────────────
	var tx bytes.Buffer
	tx.Write(compactU16(numRequiredSigs)) // signature count
	tx.Write(sig)                         // first (fee payer) signature
	// Remaining signers get zero signatures (Phoenix ix typically only needs 1 signer)
	for i := 1; i < numRequiredSigs; i++ {
		tx.Write(make([]byte, 64))
	}
	tx.Write(msgBytes)

	return tx.Bytes(), nil
}

// ── Solana RPC ────────────────────────────────────────────────────────

type rpcClient struct {
	url        string
	httpClient *http.Client
}

func newRPCClient(rpcURL string) *rpcClient {
	return &rpcClient{
		url:        rpcURL,
		httpClient: &http.Client{Timeout: 20 * time.Second},
	}
}

func (r *rpcClient) call(ctx context.Context, method string, params any) (json.RawMessage, error) {
	payload, err := json.Marshal(map[string]any{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  method,
		"params":  params,
	})
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, r.url, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := r.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("rpc request: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var rpcResp struct {
		Result json.RawMessage `json:"result"`
		Error  *struct {
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.Unmarshal(body, &rpcResp); err != nil {
		return nil, fmt.Errorf("parse rpc response: %w", err)
	}
	if rpcResp.Error != nil {
		return nil, fmt.Errorf("rpc error: %s", rpcResp.Error.Message)
	}
	return rpcResp.Result, nil
}

// GetLatestBlockhash fetches the latest blockhash from the Solana RPC.
func GetLatestBlockhash(ctx context.Context, rpcURL string) (string, error) {
	r := newRPCClient(rpcURL)
	result, err := r.call(ctx, "getLatestBlockhash", []map[string]string{{"commitment": "confirmed"}})
	if err != nil {
		return "", err
	}
	var resp struct {
		Value struct {
			Blockhash string `json:"blockhash"`
		} `json:"value"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		return "", fmt.Errorf("parse blockhash: %w", err)
	}
	return resp.Value.Blockhash, nil
}

// SendTransaction submits a signed transaction to Solana RPC.
// Returns the transaction signature (base58) on success.
func SendTransaction(ctx context.Context, rpcURL string, txBytes []byte) (string, error) {
	r := newRPCClient(rpcURL)
	encoded := base64.StdEncoding.EncodeToString(txBytes)
	result, err := r.call(ctx, "sendTransaction", []any{
		encoded,
		map[string]any{
			"encoding":            "base64",
			"preflightCommitment": "confirmed",
		},
	})
	if err != nil {
		return "", err
	}
	var sig string
	if err := json.Unmarshal(result, &sig); err != nil {
		return "", fmt.Errorf("parse signature: %w", err)
	}
	return sig, nil
}

// SignAndSend is a convenience wrapper: fetches blockhash, signs instructions, submits tx.
func SignAndSend(ctx context.Context, kp *Keypair, instructions []Instruction, rpcURL string) (string, error) {
	blockhash, err := GetLatestBlockhash(ctx, rpcURL)
	if err != nil {
		return "", fmt.Errorf("get blockhash: %w", err)
	}
	txBytes, err := BuildAndSign(kp, instructions, blockhash)
	if err != nil {
		return "", fmt.Errorf("build tx: %w", err)
	}
	return SendTransaction(ctx, rpcURL, txBytes)
}
