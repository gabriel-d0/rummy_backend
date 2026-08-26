import { describe, it, expect } from "vitest";
import { FontFamily, textStyle, loadFonts } from "./fonts";

describe("Font System", () => {
  it("FontFamily has Inter and JetBrains Mono", () => {
    expect(FontFamily.inter).toBe("Inter");
    expect(FontFamily.mono).toBe("JetBrains Mono");
  });

  it("textStyle returns correct presets with Inter bold", () => {
    const title = textStyle("title");
    expect(title.fontFamily).toBe("Inter");
    expect(title.fontSize).toBe("26px");
    expect(title.fontStyle).toBe("bold");
    expect(title.color).toBe("#ffffff");
    expect(title.shadow).toBeTruthy();

    const subtitle = textStyle("subtitle");
    expect(subtitle.fontFamily).toBe("Inter");
    expect(subtitle.fontSize).toBe("11px");
    expect(subtitle.fontStyle).toBe("bold");
    expect(subtitle.backgroundColor).toBe("#0a2e1aee");

    const label = textStyle("label");
    expect(label.fontFamily).toBe("Inter");
    expect(label.fontSize).toBe("13px");

    const mono = textStyle("mono");
    expect(mono.fontFamily).toBe("JetBrains Mono");
    expect(mono.fontSize).toBe("10px");

    const debug = textStyle("debug");
    expect(debug.fontFamily).toBe("Inter");
    expect(debug.fontSize).toBe("10px");
  });

  it("loadFonts returns a promise and caches", async () => {
    const p = loadFonts();
    expect(p).toBeInstanceOf(Promise);
    await expect(p).resolves.toBeUndefined();
    const p2 = loadFonts();
    expect(p2).toBe(p);
  });
});
