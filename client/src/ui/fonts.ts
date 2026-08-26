/**
 * Font System — ensures Inter and JetBrains Mono are loaded before Phaser Text creation.
 * Uses correct Phaser TextStyle per https://docs.phaser.io/api-documentation/typedef/types-gameobjects-text#TextStyle
 * - fontFamily: "'Inter', system-ui, sans-serif" (quoted if needed)
 * - fontSize: "26px" (string with px)
 * - fontStyle: "bold" (not "700" — Phaser only understands bold/normal/italic)
 * - For weight, use font: "700 26px Inter" shorthand which overrides fontFamily/fontSize/fontStyle
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
    const families = [
      '700 26px "Inter"',
      '600 13px "Inter"',
      '600 11px "Inter"',
      '700 11px "Inter"',
      '600 10px "JetBrains Mono"',
      '400 9px "Inter"',
    ];
    try {
      await Promise.all(
        families.map((f) => (document as unknown as { fonts: { load: (s: string) => Promise<unknown> } }).fonts.load(f))
      );
      await (document as unknown as { fonts: { ready: Promise<void> } }).fonts.ready;
    } catch {
      // ignore
    }
  })();
  return fontsReady;
}

export type TextStylePreset = "title" | "subtitle" | "label" | "mono" | "debug";

export function textStyle(preset: TextStylePreset): Phaser.Types.GameObjects.Text.TextStyle {
  switch (preset) {
    case "title":
      return {
        fontFamily: FontFamily.inter,
        fontSize: "26px",
        fontStyle: "bold",
        color: "#ffffff",
        align: "center",
        shadow: { offsetX: 0, offsetY: 2, color: "#0a2e1a", blur: 6, fill: true, stroke: false },
      };
    case "subtitle":
      return {
        fontFamily: FontFamily.inter,
        fontSize: "11px",
        fontStyle: "bold",
        color: "#e0f2e0",
        backgroundColor: "#0a2e1aee",
        padding: { x: 12, y: 7 },
        align: "center",
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
      };
    case "debug":
      return {
        fontFamily: FontFamily.inter,
        fontSize: "10px",
        color: "#ffffff",
        fontStyle: "bold",
        backgroundColor: "#0a2e1acc",
        padding: { x: 5, y: 3 },
      };
    default:
      return { fontFamily: FontFamily.inter, fontSize: "13px", color: "#ffffff", fontStyle: "bold" };
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
  // Ensure crisp rendering
  t.setResolution(window.devicePixelRatio || 1);
  return t;
}
