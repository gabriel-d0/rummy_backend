/**
 * Font System — crisp Inter via correct Phaser TextStyle + high DPI + font-display swap.
 * Fixes blocky: was using backgroundColor + shadow + late font load → fallback bitmap.
 */

export const FontFamily = {
  inter: "Inter",
  mono: "JetBrains Mono",
} as const;

let fontsReady: Promise<void> | null = null;

export function loadFonts(): Promise<void> {
  if (fontsReady) return fontsReady;
  fontsReady = (async () => {
    if (typeof document === "undefined" || !("fonts" in document)) return;
    try {
      // Load critical weights via CSS Font Loading API
      await Promise.all([
        (document as unknown as { fonts: { load: (s: string) => Promise<unknown> } }).fonts.load('700 32px "Inter"'),
        (document as unknown as { fonts: { load: (s: string) => Promise<unknown> } }).fonts.load('600 14px "Inter"'),
        (document as unknown as { fonts: { load: (s: string) => Promise<unknown> } }).fonts.load('600 11px "Inter"'),
        (document as unknown as { fonts: { load: (s: string) => Promise<unknown> } }).fonts.load('400 12px "Inter"'),
      ]);
      await (document as unknown as { fonts: { ready: Promise<void> } }).fonts.ready;
      // Small delay to ensure Phaser's canvas picks up the loaded font
      await new Promise((r) => setTimeout(r, 50));
    } catch {
      // ignore
    }
  })();
  return fontsReady;
}

export type TextStylePreset = "title" | "subtitle" | "label" | "mono" | "debug";

export function textStyle(preset: TextStylePreset): Phaser.Types.GameObjects.Text.TextStyle {
  const resolution = typeof window !== "undefined" ? Math.max(2, window.devicePixelRatio || 1) : 2;
  switch (preset) {
    case "title":
      return {
        fontFamily: FontFamily.inter,
        fontSize: "32px",
        fontStyle: "bold",
        color: "#ffffff",
        align: "center",
        resolution,
        // No backgroundColor, no stroke, no shadow — just crisp Inter 700
      };
    case "subtitle":
      return {
        fontFamily: FontFamily.inter,
        fontSize: "12px",
        fontStyle: "bold",
        color: "#d7e8d7",
        align: "center",
        resolution,
        backgroundColor: "#0a2e1aee",
        padding: { x: 12, y: 6 },
      };
    case "label":
      return {
        fontFamily: FontFamily.inter,
        fontSize: "13px",
        fontStyle: "bold",
        color: "#e8f5e9",
        backgroundColor: "#0a2e1a",
        padding: { x: 14, y: 8 },
        align: "center",
        resolution,
      };
    case "mono":
      return {
        fontFamily: FontFamily.mono,
        fontSize: "10px",
        fontStyle: "bold",
        color: "#a5d6a7",
        backgroundColor: "#0a2e1a",
        padding: { x: 10, y: 6 },
        align: "center",
        resolution,
      };
    case "debug":
      return {
        fontFamily: FontFamily.inter,
        fontSize: "10px",
        color: "#ffffff",
        fontStyle: "bold",
        backgroundColor: "#0a2e1acc",
        padding: { x: 5, y: 3 },
        resolution,
      };
    default:
      return { fontFamily: FontFamily.inter, fontSize: "13px", color: "#ffffff", fontStyle: "bold", resolution };
  }
}

export function createText(
  scene: Phaser.Scene,
  x: number,
  y: number,
  text: string,
  preset: TextStylePreset,
  origin = 0.5
): Phaser.GameObjects.Text {
  const t = scene.add.text(x, y, text, textStyle(preset));
  t.setOrigin(origin, origin);
  return t;
}
