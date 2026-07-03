import { afterEach, beforeEach, describe, expect, it } from "vitest";

import { ClawdBotApp } from "./app";

const originalConnect = ClawdBotApp.prototype.connect;

function mountApp(pathname: string) {
  window.history.replaceState({}, "", pathname);
  const app = document.createElement("clawdbot-app") as ClawdBotApp;
  document.body.append(app);
  return app;
}

beforeEach(() => {
  ClawdBotApp.prototype.connect = () => {
    // no-op: avoid real gateway WS connections in browser tests
  };
  window.__CLAWDBOT_CONTROL_UI_BASE_PATH__ = undefined;
  localStorage.clear();
  document.body.innerHTML = "";
});

afterEach(() => {
  ClawdBotApp.prototype.connect = originalConnect;
  window.__CLAWDBOT_CONTROL_UI_BASE_PATH__ = undefined;
  localStorage.clear();
  document.body.innerHTML = "";
});

describe("chat markdown rendering", () => {
  it("renders markdown inside tool output sidebar", async () => {
    const app = mountApp("/chat");
    await app.updateComplete;

    const timestamp = Date.now();
    app.chatMessages = [
      {
        role: "assistant",
        content: [
          { type: "toolcall", name: "noop", arguments: {} },
          { type: "toolresult", name: "noop", text: "Hello **world**" },
        ],
        timestamp,
      },
    ];

    await app.updateComplete;

    const toolCard = Array.from(app.querySelectorAll(".chat-tool-card")).find(
      (card) => card.textContent?.includes("Hello **world**"),
    ) as HTMLElement | undefined;
    expect(toolCard).toBeTruthy();
    toolCard?.click();

    await app.updateComplete;

    const strong = app.querySelector(".sidebar-markdown strong");
    expect(strong?.textContent).toBe("world");
  });
});
