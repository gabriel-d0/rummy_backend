// Day 29-30: PublicSnapshot/PrivateSnapshot types mirroring Go internal/match/visibility.go:36
export interface TileInstance { ID: string; Colour: number; Rank: number; IsJoker: boolean; }
export interface DiscardEntry { Tile: TileInstance; IsOpeningDiscard: boolean; Index: number; }
export interface TableMeld { ID: string; Kind: string; Tiles: TileInstance[]; JokerReps: Record<string, TileInstance>; OwnerSeat: number; }
export interface PublicSnapshot {
  v: number;
  gamePhase: string;
  turnPhase: string;
  currentSeat: number;
  players: { id: string; seat: number; hasOpened: boolean; rackCount: number }[];
  stockCount: number;
  discardRow: DiscardEntry[];
  tableMelds: TableMeld[];
  winner: number;
}
export interface PrivateSnapshot extends PublicSnapshot {
  ownRack: TileInstance[];
  ownSeat: number;
}
