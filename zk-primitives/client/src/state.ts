/**
 * Light Protocol state tree helpers.
 *
 * Wraps the Photon indexer queries (via `@lightprotocol/stateless.js`)
 * and the address derivation / proof-fetching boilerplate for our
 * compressed accounts.
 */

import type { PublicKey } from "@solana/web3.js";
import type { Bytes32 } from "./types.js";

export interface ValidityProof {
  /** Compressed proof, 128 bytes (Groth16 over the Merkle proof). */
  compressedProof: Uint8Array;
  /** Root indices for each leaf. */
  rootIndices: number[];
  /** Addresses (optional, for create-with-address flows). */
  addresses?: Uint8Array[];
}

export interface PackedAddressTreeInfo {
  addressMerkleTreePubkeyIndex: number;
  addressQueuePubkeyIndex: number;
  rootIndex: number;
}

export interface PackedStateTreeInfo {
  stateTree: PublicKey;
  outputQueue: PublicKey;
  rootIndex: number;
}

/** Fetch a non-inclusion proof for a list of nullifier addresses. */
export async function fetchValidityProofV2(args: {
  rpc: any;
  hashes?: Uint8Array[];
  addressesWithTrees?: { address: Uint8Array; tree: PublicKey }[];
}): Promise<ValidityProof> {
  const rpc = args.rpc as {
    getValidityProofV0?: (
      hashes?: Uint8Array[],
      newAddresses?: { address: Uint8Array; tree: PublicKey }[],
    ) => Promise<{ compressedProof: Uint8Array; rootIndices: number[]; addresses?: Uint8Array[] }>;
  };
  if (!rpc.getValidityProofV0) {
    throw new Error("RPC client does not expose getValidityProofV0; use a Light-compatible RPC wrapper.");
  }
  const result = await rpc.getValidityProofV0(args.hashes ?? [], args.addressesWithTrees ?? []);
  return {
    compressedProof: result.compressedProof,
    rootIndices: result.rootIndices,
    addresses: result.addresses,
  };
}

/** Fetch the current v2 address tree info. */
export async function fetchAddressTreeV2(rpc: any): Promise<PublicKey> {
  if (typeof rpc.getAddressTreeV2 !== "function") {
    throw new Error("RPC client does not expose getAddressTreeV2; use a Light-compatible RPC wrapper.");
  }
  return rpc.getAddressTreeV2();
}

/** Fetch a random v2 state tree for new outputs. */
export async function fetchRandomStateTreeV2(rpc: any): Promise<{
  tree: PublicKey;
  queue: PublicKey;
}> {
  if (typeof rpc.getRandomStateTreeInfo !== "function") {
    throw new Error("RPC client does not expose getRandomStateTreeInfo; use a Light-compatible RPC wrapper.");
  }
  const info = await rpc.getRandomStateTreeInfo();
  return { tree: info.tree, queue: info.queue };
}

/** Pack tree info into a `PackedAccounts` struct for the on-chain CPI. */
export async function packAccounts(args: {
  rpc: any;
  programId: PublicKey;
  treeInfo: { addressTree: PublicKey; stateTree: PublicKey; outputQueue: PublicKey };
  proof: ValidityProof;
}) {
  const sdk = await import("@lightprotocol/stateless.js");
  const { PackedAccounts, SystemAccountMetaConfig } = sdk;
  const packed = new PackedAccounts();
  packed.addSystemAccountsV2(SystemAccountMetaConfig.new(args.programId));

  const addressMerkleTreeIndex = packed.insertOrGet(args.treeInfo.addressTree);
  const outputStateTreeIndex = packed.insertOrGet(args.treeInfo.stateTree);

  return {
    packed,
    addressMerkleTreeIndex,
    outputStateTreeIndex,
  };
}
