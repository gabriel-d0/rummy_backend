import { describe, it, expect, vi, beforeEach } from "vitest";
import { getClient } from "./client";

describe("Nakama client — Day 4 SvelteKit", () => {
  beforeEach(() => {
    if (typeof globalThis.localStorage === "undefined") {
      (globalThis as unknown as Record<string, unknown>).localStorage = {
        getItem: vi.fn(() => null),
        setItem: vi.fn(),
        removeItem: vi.fn(),
        key: () => null,
        length: 0,
        clear: vi.fn(),
      } as unknown as Storage;
    }
  });

  it("getClient creates Client with defaultkey", () => {
    const client = getClient();
    expect(client).toBeTruthy();
    expect(typeof client.authenticateDevice).toBe("function");
  });

  it("getClient is singleton", () => {
    const c1 = getClient();
    const c2 = getClient();
    expect(c1).toBe(c2);
  });
});
