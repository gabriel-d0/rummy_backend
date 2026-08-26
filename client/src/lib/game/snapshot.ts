// Day 17 — Snapshot types — mirrors Go internal/match/visibility.go:36 Version 1

export const SnapshotVersion = 1;

export type TileInstance = {
	ID: string;
	Colour: number;
	Rank: number;
	IsJoker: boolean;
};

export type DiscardEntry = {
	Tile: TileInstance;
	IsOpeningDiscard: boolean;
	Index: number;
};

export type TableMeld = {
	ID: string;
	Kind: string;
	Tiles: TileInstance[];
	JokerReps: Record<string, TileInstance>;
	OwnerSeat: number;
};

export type PublicPlayer = {
	id: string;
	seat: number;
	hasOpened: boolean;
	rackCount: number;
};

export type PublicSnapshot = {
	v: number;
	gamePhase: string;
	turnPhase: string;
	currentSeat: number;
	players: PublicPlayer[];
	stockCount: number;
	discardRow: DiscardEntry[];
	tableMelds: TableMeld[];
	winner: number;
};

export type PrivateSnapshot = PublicSnapshot & {
	ownRack: TileInstance[];
	ownSeat: number;
};

export function isValidPublicSnapshot(snap: unknown): snap is PublicSnapshot {
	if (typeof snap !== 'object' || snap === null) return false;
	const s = snap as unknown as Record<string, unknown>;
	if (s.v !== SnapshotVersion) return false;
	if (typeof s.gamePhase !== 'string' || typeof s.turnPhase !== 'string') return false;
	if (
		typeof s.currentSeat !== 'number' ||
		typeof s.stockCount !== 'number' ||
		typeof s.winner !== 'number'
	)
		return false;
	if (!Array.isArray(s.players) || !Array.isArray(s.discardRow) || !Array.isArray(s.tableMelds))
		return false;
	return true;
}

export function isValidPrivateSnapshot(snap: unknown): snap is PrivateSnapshot {
	if (!isValidPublicSnapshot(snap)) return false;
	const s = snap as unknown as Record<string, unknown>;
	if (!Array.isArray(s.ownRack) || typeof s.ownSeat !== 'number') return false;
	return true;
}

export function checkNoLeak(publicJson: string, privateIds: string[]): boolean {
	for (const id of privateIds) {
		if (publicJson.includes(id)) return false;
	}
	return true;
}
