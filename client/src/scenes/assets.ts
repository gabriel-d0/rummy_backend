// Day 10: Asset manifest — single source of truth for all 4 assets
export const ASSET_MANIFEST = {
  tile: "assets/tile.png",
  joker: "assets/joker.png",
  table: "assets/table.png",
  rack: "assets/rack.png",
} as const;

export type AssetKey = keyof typeof ASSET_MANIFEST;
