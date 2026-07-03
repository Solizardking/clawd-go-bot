export type ParsedAgentSessionKey = {
  raw: string;
  agentId: string;
  channel?: string;
  peer?: string;
};

export function parseAgentSessionKey(value: string | null | undefined): ParsedAgentSessionKey | null {
  const raw = value?.trim();
  if (!raw) return null;

  if (raw === "main") {
    return { raw, agentId: "main", channel: "main" };
  }

  const parts = raw.split(":").filter((part) => part.length > 0);
  if (parts.length === 0) return null;

  if (parts[0] === "agent") {
    const agentId = parts[1]?.trim();
    if (!agentId) return null;
    return {
      raw,
      agentId,
      channel: parts[2],
      peer: parts.length > 3 ? parts.slice(3).join(":") : undefined,
    };
  }

  const agentId = parts[0]?.trim();
  if (!agentId) return null;
  return {
    raw,
    agentId,
    channel: parts[1],
    peer: parts.length > 2 ? parts.slice(2).join(":") : undefined,
  };
}
