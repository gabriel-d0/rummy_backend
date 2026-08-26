import { createSocket, getSocket, disconnectSocket } from './socket';
import { joinMatch, getStoredMatchId, getMatchId } from './match';
import { handleMatchData } from '../game/store';
import { OpServerState } from './protocol';
import type { PrivateSnapshot } from '../game/snapshot';

// Day 26 — Reconnection — rejoin — socket.connect + joinMatch(matchId) expects OpServerState 100 PrivateSnapshot for that Seat only, Rack rehydrates from OwnRack not old

export async function reconnect(): Promise<string | null> {
	const storedId = getStoredMatchId() ?? getMatchId();
	if (!storedId) return null;

	// Ensure socket (re)connected — createSocket handles mock in vitest
	let sock = getSocket();
	if (!sock) {
		try {
			sock = await createSocket();
		} catch (_err) {
			void _err;
			return null;
		}
	} else {
		// In real Nakama, socket may be disconnected; ensure connected
		// For mock, socket is already connected
		try {
			// If socket has disconnect flag, recreate
			// eslint-disable-next-line @typescript-eslint/no-explicit-any
			if ((sock as any).isConnected === false) {
				disconnectSocket();
				sock = await createSocket();
			}
		} catch (_err) {
			void _err;
		}
	}

	try {
		const rejoinedId = await joinMatch(storedId);
		return rejoinedId;
	} catch (_err) {
		void _err;
		return null;
	}
}

// Helper for tests: simulate server sending new PrivateSnapshot after rejoin
// Rack must rehydrate from OwnRack new, not old
export function simulatePrivateAfterReconnect(
	newOwnRackIds: string[],
	seat = 0,
	currentSeat = seat
): PrivateSnapshot {
	const snap: PrivateSnapshot = {
		v: 1,
		gamePhase: 'Playing',
		turnPhase: 'MeldOrDiscard',
		currentSeat,
		players: [
			{ id: 'alice', seat: 0, hasOpened: true, rackCount: seat === 0 ? newOwnRackIds.length : 7 },
			{ id: 'bob', seat: 1, hasOpened: false, rackCount: seat === 1 ? newOwnRackIds.length : 14 }
		],
		stockCount: 70,
		discardRow: [],
		tableMelds: [],
		winner: -1,
		ownRack: newOwnRackIds.map((id) => ({ ID: id, Colour: 1, Rank: 5, IsJoker: false })),
		ownSeat: seat
	};
	// Feed via handleMatchData as server would (OpServerState 100)
	handleMatchData(OpServerState, JSON.stringify(snap));
	return snap;
}

export function hasStoredMatch(): boolean {
	try {
		return !!localStorage.getItem('rummy_matchId');
	} catch (_err) {
		void _err;
		return false;
	}
}
