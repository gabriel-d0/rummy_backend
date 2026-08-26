import { test, expect } from "@playwright/test";

// Day 28: Error display — OpServerError 102 red toast for 3s
// Tests what we implemented until now: ErrorToast, sync onServerError, envelope merge, requestId

test.describe("Error handling — Day 28", () => {
  test("ErrorToast module exists with showErrorToast and ServerError", async ({ page }) => {
    await page.goto("/");
    await expect(page.locator("#game canvas")).toBeVisible({ timeout: 10000 });
    const hasModule = await page.evaluate(async () => {
      const mod = await import("/src/ui/ErrorToast.ts");
      return (
        typeof mod.showErrorToast === "function" &&
        typeof mod.clearErrorToasts === "function" &&
        typeof mod.getLastError === "function" &&
        typeof mod.clearLastError === "function"
      );
    });
    expect(hasModule).toBeTruthy();
  });

  test("showErrorToast stores lastError and creates DOM toast", async ({ page }) => {
    await page.goto("/");
    await expect(page.locator("#game canvas")).toBeVisible({ timeout: 10000 });
    await page.waitForTimeout(500);
    const toastOk = await page.evaluate(async () => {
      const toast = await import("/src/ui/ErrorToast.ts");
      toast.clearErrorToasts();
      toast.showErrorToast({ code: "not_your_turn", message: "not your turn", op: 2, requestId: "req-001" });
      const err = toast.getLastError();
      const hasError = err !== null && err.code === "not_your_turn" && err.requestId === "req-001" && err.op === 2;
      // Check DOM container was created
      const container = document.getElementById("error-toasts");
      const hasContainer = !!container;
      const hasToast = hasContainer && container!.children.length > 0;
      let toastTextOk = false;
      if (hasToast) {
        const toastEl = container!.lastChild as HTMLElement;
        toastTextOk = toastEl.textContent !== null && toastEl.textContent.includes("not_your_turn");
      }
      // Cleanup
      toast.clearErrorToasts();
      return hasError && hasContainer && hasToast && toastTextOk;
    });
    expect(toastOk).toBeTruthy();
  });

  test("OpServerError toast shows red background and auto-removes after 3s", async ({ page }) => {
    await page.goto("/");
    await expect(page.locator("#game canvas")).toBeVisible({ timeout: 10000 });
    await page.waitForTimeout(300);
    const beforeCount = await page.evaluate(() => document.getElementById("error-toasts")?.children.length ?? 0);
    await page.evaluate(async () => {
      const toast = await import("/src/ui/ErrorToast.ts");
      toast.showErrorToast({ code: "wrong_phase", message: "wrong phase", op: 3 });
    });
    await page.waitForTimeout(200);
    const hasToast = await page.evaluate((before) => {
      const c = document.getElementById("error-toasts");
      return !!c && c.children.length > before;
    }, beforeCount);
    expect(hasToast).toBeTruthy();
    // Check style is red
    const isRed = await page.evaluate(() => {
      const c = document.getElementById("error-toasts");
      if (!c || c.children.length === 0) return false;
      const el = c.lastChild as HTMLElement;
      const style = window.getComputedStyle(el);
      // background should be #dc2626 or rgb(220, 38, 38)
      return style.backgroundColor.includes("220") || el.style.cssText.includes("dc2626") || el.style.background.includes("dc2626");
    });
    expect(isRed).toBeTruthy();
    // Wait for auto-remove (3s + 300ms fade)
    await page.waitForTimeout(3600);
    const afterCount = await page.evaluate(() => document.getElementById("error-toasts")?.children.length ?? 0);
    // After 3.6s, toast should be removed (or at least not growing indefinitely)
    expect(afterCount).toBeLessThanOrEqual(beforeCount);
  });

  test("handleMatchData merges envelope requestId/op for OpServerError", async ({ page }) => {
    await page.goto("/");
    await expect(page.locator("#game canvas")).toBeVisible({ timeout: 10000 });
    const mergeOk = await page.evaluate(async () => {
      const sync = await import("/src/state/sync.ts");
      const toast = await import("/src/ui/ErrorToast.ts");
      toast.clearErrorToasts();
      const envelope = JSON.stringify({
        v: 1,
        op: 102,
        requestId: "req-merge-test",
        payload: { code: "bad_payload", message: "tileId required" },
      });
      sync.handleMatchData(102, envelope);
      const err = toast.getLastError();
      return err !== null && err.code === "bad_payload" && err.requestId === "req-merge-test" && err.op === 102;
    });
    expect(mergeOk).toBeTruthy();
  });

  test("multiple error codes are handled: bad_version, not_opened, already_opened", async ({ page }) => {
    await page.goto("/");
    await expect(page.locator("#game canvas")).toBeVisible({ timeout: 10000 });
    const codesOk = await page.evaluate(async () => {
      const toast = await import("/src/ui/ErrorToast.ts");
      const sync = await import("/src/state/sync.ts");
      const codes = ["bad_version", "bad_payload", "not_your_turn", "wrong_phase", "not_opened", "already_opened"];
      for (const code of codes) {
        toast.clearErrorToasts();
        sync.handleMatchData(102, JSON.stringify({ v: 1, op: 102, payload: { code, message: "test" } }));
        const err = toast.getLastError();
        if (!err || err.code !== code) return false;
      }
      toast.clearErrorToasts();
      return true;
    });
    expect(codesOk).toBeTruthy();
  });

  test("clearErrorToasts clears lastError and DOM", async ({ page }) => {
    await page.goto("/");
    await expect(page.locator("#game canvas")).toBeVisible({ timeout: 10000 });
    const clearOk = await page.evaluate(async () => {
      const toast = await import("/src/ui/ErrorToast.ts");
      toast.showErrorToast({ code: "test_clear", message: "clear me" });
      const hasErrorBefore = toast.getLastError() !== null;
      toast.clearErrorToasts();
      const hasErrorAfter = toast.getLastError() === null;
      const container = document.getElementById("error-toasts");
      const isEmpty = !container || container.children.length === 0 || container.innerHTML === "";
      return hasErrorBefore && hasErrorAfter && isEmpty;
    });
    expect(clearOk).toBeTruthy();
  });
});
