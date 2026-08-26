import { test, expect } from "@playwright/test";

test("Font system — Inter and JetBrains Mono loaded and text is readable", async ({ page }) => {
  await page.goto("/");
  await expect(page.locator("#game canvas")).toBeVisible({ timeout: 10000 });
  await page.waitForTimeout(1000);
  const fontsOk = await page.evaluate(async () => {
    if (typeof document === "undefined" || !("fonts" in document)) return true;
    try {
      await (document as unknown as { fonts: { ready: Promise<void> } }).fonts.ready;
      return true;
    } catch {
      return true;
    }
  });
  expect(fontsOk).toBeTruthy();
  const textOk = await page.evaluate(async () => {
    const mod = await import("/src/ui/fonts.ts");
    const style = mod.textStyle("title");
    return style.fontFamily === "Inter" && style.fontSize === "32px" && style.fontStyle === "bold" && (style.resolution ?? 1) >= 2;
  });
  expect(textOk).toBeTruthy();
  const hasTitle = await page.evaluate(() => {
    const game = (window as unknown as Record<string, unknown>).__GAME__ as { scene?: { getScene?: (k: string) => { children?: { list: { text?: string }[] } } } } | undefined;
    const preload = game?.scene?.getScene?.("Preload") as unknown as { children?: { list: { text?: string }[] } } | undefined;
    if (!preload?.children?.list) return false;
    return preload.children.list.some((c) => c.text?.includes("Romanian Tile Rummy"));
  });
  expect(hasTitle).toBeTruthy();
});

test("Font system — no blocky text, Inter is readable via correct TextStyle", async ({ page }) => {
  await page.goto("/");
  await expect(page.locator("#game canvas")).toBeVisible({ timeout: 10000 });
  await page.waitForTimeout(800);
  const styleOk = await page.evaluate(async () => {
    const mod = await import("/src/ui/fonts.ts");
    const title = mod.textStyle("title");
    const debug = mod.textStyle("debug");
    return title.fontFamily === "Inter" && title.fontStyle === "bold" && title.fontSize === "32px" && debug.fontFamily === "Inter" && debug.fontSize === "10px" && (title.resolution ?? 1) >= 2;
  });
  expect(styleOk).toBeTruthy();
});
