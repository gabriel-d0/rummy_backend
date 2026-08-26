/**
 * Font System — ensures Inter and JetBrains Mono are loaded before Phaser Text creation.
 * Solves blocky text: Phaser was creating text before Google Fonts finished loading,
 * so it fell back to monospace bitmap and looked distorted with stroke.
 */

export const FontFamily = {
  inter: "'Inter', system-ui, -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif",
  mono: "'JetBrains Mono', ui-monospace, monospace",
  display: "'Inter', system-ui, sans-serif",
} as const;

export const FontSize = {
  xs: "10px",
  sm: "11px",
  base: "13px",
  lg: "16px",
  xl: "22px",
  title: "26px",
} as const;

// Preload link is in index.html, but we also ensure via document.fonts
let fontsReady: Promise<void> | null = null;

export function loadFonts(): Promise<void> {
  if (fontsReady) return fontsReady;
  fontsReady = (async () => {
    if (typeof document === "undefined" || !("fonts" in document)) return;
    const fontFaces = [
      '400 16px Inter',
      '500 16px Inter',
      '600 16px Inter',
      '700 16px Inter',
      '400 16px "JetBrains Mono"',
      '600 16px "JetBrains Mono"',
    ];
    try {
      await Promise.all(fontFaces.map((f) => (document as unknown as { fonts: { load: (s: string) => Promise<unknown> } }).fonts.load(f)));
      await (document as unknown as { fonts: { ready: Promise<void> } }).fonts.ready;
    } catch {
      // ignore, fallback will still render
    }
  })();
  return fontsReady;
}

// Helper to create Phaser Text with correct font — ensures family and style are consistent
export type TextStylePreset = "title" | "subtitle" | "label" | "mono" | "debug";

export function textStyle(preset: TextStylePreset): Phaser.Types.GameObjects.Text.TextStyle {
  switch (preset) {
    case "title":
      return {
        fontFamily: FontFamily.display,
        fontSize: FontSize.title,
        color: "#ffffff",
        fontStyle: "bold",
        align: "center",
        // Use shadow instead of stroke for crisp Inter
        shadow: { offsetX: 0, offsetY: 2, color: "#0a2e1a", blur: 8, fill: true },
      };
    case "subtitle":
      return {
        fontFamily: FontFamily.inter,
        fontSize: FontSize.sm,
        color: "#e0f2e0",
        fontStyle: "bold",
        backgroundColor: "#0a2e1aee",
        padding: { x: 12, y: 7 },
        align: "center",
      };
    case "label":
      return {
        fontFamily: FontFamily.inter,
        fontSize: FontSize.base,
        color: "#e8f5e9",
        fontStyle: "bold",
        backgroundColor: "#0a2e1a",
        padding: { x: 14, y: 8 },
        align: "center",
      };
    case "mono":
      return {
        fontFamily: FontFamily.mono,
        fontSize: FontSize.xs,
        color: "#a5d6a7",
        fontStyle: "bold",
        backgroundColor: "#0a2e1a",
        padding: { x: 10, y: 6 },
        align: "center",
      };
    case "debug":
      return {
        fontFamily: FontFamily.inter,
        fontSize: "9px",
        color: "#ffffff",
        fontStyle: "600",
        backgroundColor: "#0a2e1a99",
        padding: { x: 4, y: 2 },
      };
    default:
      return { fontFamily: FontFamily.inter, fontSize: FontSize.base, color: "#ffffff" };
  }
}

// Convenience: create text via preset
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
