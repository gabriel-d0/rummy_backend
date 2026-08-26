import { describe, it, expect } from "vitest";
import { FontFamily, FontSize, textStyle, loadFonts } from "./fonts";

describe("Font System", () => {
  it("FontFamily has Inter and JetBrains Mono", () => {
    expect(FontFamily.inter).toContain("Inter");
    expect(FontFamily.mono).toContain("JetBrains Mono");
    expect(FontFamily.display).toContain("Inter");
  });

  it("FontSize presets are defined", () => {
    expect(FontSize.title).toBe("26px");
    expect(FontSize.base).toBe("13px");
    expect(FontSize.xs).toBe("10px");
  });

  it("textStyle returns correct presets", () => {
    const title = textStyle("title");
    expect(title.fontFamily).toContain("Inter");
    expect(title.fontSize).toBe("26px");
    expect(title.fontStyle).toBe("bold");
    expect(title.color).toBe("#ffffff");

    const subtitle = textStyle("subtitle");
    expect(subtitle.fontSize).toBe("11px");
    expect(subtitle.backgroundColor).toBe("#0a2e1aee");

    const label = textStyle("label");
    expect(label.fontSize).toBe("13px");

    const mono = textStyle("mono");
    expect(mono.fontFamily).toContain("JetBrains Mono");

    const debug = textStyle("debug");
    expect(debug.fontSize).toBe("9px");
  });

  it("loadFonts returns a promise", async () => {
    const p = loadFonts();
    expect(p).toBeInstanceOf(Promise);
    await expect(p).resolves.toBeUndefined();
    // Second call should return same promise
    const p2 = loadFonts();
    expect(p2).toBe(p);
  });
});
