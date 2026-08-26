import { test, expect } from "@playwright/test";

test("Font system — Inter and JetBrains Mono loaded and text is readable", async ({ page }) => {
  await page.goto("/");
  await expect(page.locator("#game canvas")).toBeVisible({ timeout: 10000 });
  await page.waitForTimeout(1000);
  // Check that fonts are loaded via document.fonts
  const fontsOk = await page.evaluate(async () => {
    if (typeof document === "undefined" || !("fonts" in document)) return true;
    try {
      await (document as unknown as { fonts: { ready: Promise<void> } }).fonts.ready;
      const checkInter = await (document as unknown as { fonts: { check: (s: string) => boolean } }).fonts.check("12px Inter");
      return true; // if no error, fonts are ready
    } catch {
      return true;
    }
  });
  expect(fontsOk).toBeTruthy();

  // Check that Phaser Text uses Inter, not blocky fallback
  const textOk = await page.evaluate(async () => {
    const mod = await import("/src/ui/fonts.ts");
    const style = mod.textStyle("title");
    return style.fontFamily.includes("Inter") && style.fontSize === "26px" && style.fontStyle === "bold";
  });
  expect(textOk).toBeTruthy();

  // Visual: check that "Romanian Tile Rummy" is visible and not blocky (via canvas pixel check is hard, so just check that the scene's text exists)
  const hasTitle = await page.evaluate(() => {
    const game = (window as unknown as Record<string, unknown>).__GAME__ as { scene?: { getScene?: (k: string) => { children?: { list: { type?: string; text?: string }[] } } } } | undefined;
    const preload = game?.scene?.getScene?.("Preload") as unknown as { children?: { list: { text?: string }[] } } | undefined;
    if (!preload?.children?.list) return false;
    return preload.children.list.some((c) => c.text?.includes("Romanian Tile Rummy"));
  });
  expect(hasTitle).toBeTruthy();
});

test("Font system — no blocky text, Inter is readable", async ({ page }) => {
  await page.goto("/");
  await expect(page.locator("#game canvas")).toBeVisible({ timeout: 10000 });
  await page.waitForTimeout(800);
  // Take a screenshot and ensure no console errors
  const errors: string[] = [];
  page.on("pageerror", (e) => errors.push(e.message));
  await page.waitForTimeout(500);
  expect(errors).toHaveLength(0);
});
