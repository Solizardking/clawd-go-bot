import { describe, expect, it } from "vitest";

import { parseAgentSessionKey } from "./session-key";

describe("parseAgentSessionKey", () => {
  it("parses agent-prefixed main session keys", () => {
    expect(parseAgentSessionKey("agent:main:main")).toEqual({
      raw: "agent:main:main",
      agentId: "main",
      channel: "main",
    });
  });

  it("parses agent-prefixed channel session keys", () => {
    expect(parseAgentSessionKey("agent:main:whatsapp:dm:+15555550123")).toEqual({
      raw: "agent:main:whatsapp:dm:+15555550123",
      agentId: "main",
      channel: "whatsapp",
      peer: "dm:+15555550123",
    });
  });

  it("parses legacy routing session keys", () => {
    expect(parseAgentSessionKey("clawd:telegram:1234")).toEqual({
      raw: "clawd:telegram:1234",
      agentId: "clawd",
      channel: "telegram",
      peer: "1234",
    });
  });

  it("normalizes the main alias", () => {
    expect(parseAgentSessionKey("main")).toEqual({
      raw: "main",
      agentId: "main",
      channel: "main",
    });
  });
});
