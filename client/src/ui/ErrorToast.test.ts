import { describe, it, expect, beforeEach, afterEach, vi } from "vitest";
import { showErrorToast, clearErrorToasts, getLastError, clearLastError } from "./ErrorToast";
import { handleMatchData, onServerError } from "../state/sync";

// Day 28: ErrorToast — shows OpServerError 102 code/message/requestId/op as red toast for 3s

describe("ErrorToast — Day 28", () => {
  beforeEach(() => {
    clearLastError();
    clearErrorToasts();
    // Mock localStorage for handleMatchData tests
    if (typeof globalThis.localStorage === "undefined") {
      (globalThis as any).localStorage = {
        getItem: () => null,
        setItem: () => {},
        removeItem: () => {},
      };
    }
  });

  afterEach(() => {
    vi.restoreAllMocks();
    clearErrorToasts();
  });

  it("showErrorToast stores lastError without throwing", () => {
    expect(() =>
      showErrorToast({ code: "not_your_turn", message: "not your turn", op: 2 })
    ).not.toThrow();
    expect(getLastError()?.code).toBe("not_your_turn");
    expect(getLastError()?.op).toBe(2);
  });

  it("showErrorToast handles bad_payload with requestId", () => {
    showErrorToast({
      code: "bad_payload",
      message: "tileId required",
      requestId: "req-123",
      op: 2,
    });
    const err = getLastError();
    expect(err?.code).toBe("bad_payload");
    expect(err?.requestId).toBe("req-123");
    expect(err?.op).toBe(2);
  });

  it("onServerError via sync.ts stores and toasts", () => {
    onServerError({ code: "wrong_phase", message: "wrong phase", op: 3 });
    expect(getLastError()?.code).toBe("wrong_phase");
  });

  it("handleMatchData routes op 102 envelope to ErrorToast with merge", () => {
    const envelope = JSON.stringify({
      v: 1,
      op: 102,
      requestId: "req-xyz",
      payload: { code: "bad_version", message: "bad version", op: 102 },
    });
    handleMatchData(102, envelope);
    const err = getLastError();
    expect(err?.code).toBe("bad_version");
    expect(err?.requestId).toBe("req-xyz");
    expect(err?.op).toBe(102);
  });

  it("handleMatchData merges envelope requestId when payload missing", () => {
    const envelope = JSON.stringify({
      v: 1,
      op: 102,
      requestId: "req-merge",
      payload: { code: "not_opened", message: "not opened" },
    });
    handleMatchData(102, envelope);
    const err = getLastError();
    expect(err?.code).toBe("not_opened");
    expect(err?.requestId).toBe("req-merge");
  });

  it("showErrorToast does not throw when document is undefined (node)", () => {
    const origDoc = (globalThis as any).document;
    (globalThis as any).document = undefined;
    expect(() =>
      showErrorToast({ code: "already_opened", message: "already opened", op: 6 })
    ).not.toThrow();
    expect(getLastError()?.code).toBe("already_opened");
    (globalThis as any).document = origDoc;
  });

  it("clearErrorToasts clears lastError", () => {
    showErrorToast({ code: "bad_json", message: "bad json" });
    expect(getLastError()).not.toBeNull();
    clearErrorToasts();
    expect(getLastError()).toBeNull();
  });
});

// DOM toast rendering — minimal mock verifies that showErrorToast does not throw when document exists
describe("ErrorToast DOM rendering (jsdom)", () => {
  it("creates toast element when document exists (mock DOM)", () => {
    const origDoc = (globalThis as any).document;
    const containerEl: any = {
      id: "error-toasts",
      style: {},
      children: [],
      firstChild: null,
      appendChild(child: any) {
        this.children.push(child);
        this.firstChild = this.children[0];
      },
      removeChild: vi.fn(),
      innerHTML: "",
    };
    let createdContainer: any = null;
    (globalThis as any).document = {
      getElementById: vi.fn().mockImplementation((id: string) => {
        if (id === "error-toasts") return createdContainer ?? containerEl;
        return null;
      }),
      createElement: vi.fn().mockImplementation((tag: string) => {
        if (tag === "div") {
          const el: any = {
            tagName: "DIV",
            textContent: "",
            style: {} as any,
            setAttribute: vi.fn(),
            remove: vi.fn(),
            getAttribute: vi.fn(),
          };
          // If this is the container being created
          if (!createdContainer) {
            // First div created is the container
            createdContainer = el;
            el.id = "error-toasts";
            el.children = [];
            el.firstChild = null;
            el.appendChild = function (child: any) {
              this.children.push(child);
              this.firstChild = this.children[0];
            };
            el.innerHTML = "";
            el.style = {} as any;
          }
          return el;
        }
        return {} as any;
      }),
      body: {
        appendChild: vi.fn().mockImplementation((el: any) => {
          if (el.id === "error-toasts") createdContainer = el;
          return el;
        }),
      },
    } as any;

    expect(() =>
      showErrorToast({ code: "test_code", message: "test message", requestId: "r1", op: 2 })
    ).not.toThrow();
    expect(getLastError()?.code).toBe("test_code");

    (globalThis as any).document = origDoc;
  });
});
