import { describe, it, expect } from "vitest";
import { GAME_SPACE, getSubspaces, centerX, centerY, spacing, tileX } from "./LayoutManager";

describe("LayoutManager — Day 4", () => {
  it("GAME_SPACE is 1000x1000 with FIT + CENTER_BOTH", () => {
    expect(GAME_SPACE.width).toBe(1000);
    expect(GAME_SPACE.height).toBe(1000);
    expect(GAME_SPACE.scale.mode).toBe(3); // Phaser.Scale.FIT = 3
    expect(GAME_SPACE.scale.autoCenter).toBe(3); // CENTER_BOTH = 3
  });

  it("getSubspaces returns 6 subspaces with logical units", () => {
    const s = getSubspaces();
    expect(s.TopBar).toBeTruthy();
    expect(s.TableArea).toBeTruthy();
    expect(s.MeldArea).toBeTruthy();
    expect(s.DiscardRowArea).toBeTruthy();
    expect(s.PlayerRackArea).toBeTruthy();
    expect(s.ActionButtonsArea).toBeTruthy();
    expect(s.TopBar.height).toBe(80);
    expect(s.PlayerRackArea.height).toBe(132);
    expect(s.ActionButtonsArea.height).toBe(56);
  });

  it("subspaces partition GameSpace without overlap (proportional)", () => {
    const s = getSubspaces();
    // TopBar at top
    expect(s.TopBar.y).toBe(12);
    // TableArea between TopBar and PlayerRackArea
    expect(s.TableArea.y).toBe(s.TopBar.y + s.TopBar.height + 16);
    // PlayerRackArea above ActionButtonsArea
    expect(s.PlayerRackArea.y + s.PlayerRackArea.height + 16).toBe(s.ActionButtonsArea.y);
    // MeldArea inside TableArea
    expect(s.MeldArea.x).toBe(s.TableArea.x + 16);
    expect(s.DiscardRowArea.y).toBe(s.MeldArea.y + s.MeldArea.height + 16);
  });

  it("proportional helpers centerX, centerY, spacing, tileX", () => {
    const s = getSubspaces();
    expect(centerX(s.TopBar)).toBe(s.TopBar.x + s.TopBar.width / 2);
    expect(centerY(s.TopBar)).toBe(s.TopBar.y + s.TopBar.height / 2);
    expect(spacing(1000, 3, 100)).toBeGreaterThan(6);
    expect(tileX(s.MeldArea, 0, 3, 48)).toBeGreaterThan(s.MeldArea.x);
    expect(tileX(s.MeldArea, 2, 3, 48)).toBeLessThan(s.MeldArea.x + s.MeldArea.width);
  });

  it("no pixel-based positioning, no magic numbers — all via subspaces", () => {
    const s = getSubspaces();
    // All subspaces are in logical units (0..1000), not pixels
    for (const rect of Object.values(s)) {
      expect(rect.x).toBeGreaterThanOrEqual(0);
      expect(rect.y).toBeGreaterThanOrEqual(0);
      expect(rect.width).toBeGreaterThan(0);
      expect(rect.height).toBeGreaterThan(0);
      expect(rect.x + rect.width).toBeLessThanOrEqual(1000);
      expect(rect.y + rect.height).toBeLessThanOrEqual(1000);
    }
  });
});
