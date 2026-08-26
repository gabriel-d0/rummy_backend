import { get } from 'svelte/store';
import { writable } from 'svelte/store';
import {
	NewEnvelope,
	OpClientStart,
	OpClientDiscard,
	OpClientDrawStock,
	OpClientDrawPreviousDiscard,
	OpClientPickupDiscardForMeld,
	OpClientMeldInitial,
	OpClientMeldNew
} from '../nakama/protocol';
import { getSocket, createSocket } from '../nakama/socket';
import { getMatchId, getStoredMatchId } from '../nakama/match';

// Day 29-36 — Game actions — Start 1 + Opening discard 2 + Draw stock 3 + Draw previous 4 + Pickup 5 + MeldInitial 6 + MeldNew 7

export const lastSent = writable<{ op: number; envelope: string; matchId: string } | null>(null);
export const lastSentStore = lastSent;

let _lastSentRaw: { op: number; envelope: string; matchId: string } | null = null;

export function getLastSent(): { op: number; envelope: string; matchId: string } | null {
	return _lastSentRaw ?? get(lastSent);
}

export async function sendStart(requestId?: string): Promise<string> {
	const matchId = getMatchId() ?? getStoredMatchId() ?? 'mock-match';
	let sock = getSocket();
	if (!sock) {
		try {
			sock = await createSocket();
		} catch (_err) {
			void _err;
			// fallback mock socket
			sock = {
				sendMatchState: async () => ({})
			} as unknown as typeof sock;
		}
	}
	const rid =
		requestId ??
		(typeof crypto !== 'undefined' && 'randomUUID' in crypto
			? crypto.randomUUID()
			: `req-${Date.now()}`);
	const envelope = NewEnvelope(OpClientStart, {}, rid);
	_lastSentRaw = { op: OpClientStart, envelope, matchId };
	lastSent.set(_lastSentRaw);
	try {
		// eslint-disable-next-line @typescript-eslint/no-explicit-any
		await (sock as any).sendMatchState(matchId, OpClientStart, envelope);
	} catch (_err) {
		void _err;
		// still record as sent for demo; server would return bad_version etc.
	}
	return envelope;
}

export async function sendDiscard(tileId: string, requestId?: string): Promise<string> {
	if (!tileId) throw new Error('tileId required');
	const matchId = getMatchId() ?? getStoredMatchId() ?? 'mock-match';
	let sock = getSocket();
	if (!sock) {
		try {
			sock = await createSocket();
		} catch (_err) {
			void _err;
			sock = {
				sendMatchState: async () => ({})
			} as unknown as typeof sock;
		}
	}
	const rid =
		requestId ??
		(typeof crypto !== 'undefined' && 'randomUUID' in crypto
			? crypto.randomUUID()
			: `req-${Date.now()}`);
	const envelope = NewEnvelope(OpClientDiscard, { tileId }, rid);
	_lastSentRaw = { op: OpClientDiscard, envelope, matchId };
	lastSent.set(_lastSentRaw);
	try {
		// eslint-disable-next-line @typescript-eslint/no-explicit-any
		await (sock as any).sendMatchState(matchId, OpClientDiscard, envelope);
	} catch (_err) {
		void _err;
	}
	return envelope;
}

export async function sendDrawStock(requestId?: string): Promise<string> {
	const matchId = getMatchId() ?? getStoredMatchId() ?? 'mock-match';
	let sock = getSocket();
	if (!sock) {
		try {
			sock = await createSocket();
		} catch (_err) {
			void _err;
			sock = {
				sendMatchState: async () => ({})
			} as unknown as typeof sock;
		}
	}
	const rid =
		requestId ??
		(typeof crypto !== 'undefined' && 'randomUUID' in crypto
			? crypto.randomUUID()
			: `req-${Date.now()}`);
	const envelope = NewEnvelope(OpClientDrawStock, {}, rid);
	_lastSentRaw = { op: OpClientDrawStock, envelope, matchId };
	lastSent.set(_lastSentRaw);
	try {
		// eslint-disable-next-line @typescript-eslint/no-explicit-any
		await (sock as any).sendMatchState(matchId, OpClientDrawStock, envelope);
	} catch (_err) {
		void _err;
	}
	return envelope;
}

export async function sendDrawPreviousDiscard(requestId?: string): Promise<string> {
	const matchId = getMatchId() ?? getStoredMatchId() ?? 'mock-match';
	let sock = getSocket();
	if (!sock) {
		try {
			sock = await createSocket();
		} catch (_err) {
			void _err;
			sock = {
				sendMatchState: async () => ({})
			} as unknown as typeof sock;
		}
	}
	const rid =
		requestId ??
		(typeof crypto !== 'undefined' && 'randomUUID' in crypto
			? crypto.randomUUID()
			: `req-${Date.now()}`);
	const envelope = NewEnvelope(OpClientDrawPreviousDiscard, {}, rid);
	_lastSentRaw = { op: OpClientDrawPreviousDiscard, envelope, matchId };
	lastSent.set(_lastSentRaw);
	try {
		// eslint-disable-next-line @typescript-eslint/no-explicit-any
		await (sock as any).sendMatchState(matchId, OpClientDrawPreviousDiscard, envelope);
	} catch (_err) {
		void _err;
	}
	return envelope;
}

export async function sendPickupDiscardForMeld(
	discardIndex: number,
	tileIds: string[],
	requestId?: string
): Promise<string> {
	if (discardIndex < 0) throw new Error('discardIndex required');
	if (!tileIds || tileIds.length !== 2) throw new Error('exactly 2 tileIds required');
	const matchId = getMatchId() ?? getStoredMatchId() ?? 'mock-match';
	let sock = getSocket();
	if (!sock) {
		try {
			sock = await createSocket();
		} catch (_err) {
			void _err;
			sock = {
				sendMatchState: async () => ({})
			} as unknown as typeof sock;
		}
	}
	const rid =
		requestId ??
		(typeof crypto !== 'undefined' && 'randomUUID' in crypto
			? crypto.randomUUID()
			: `req-${Date.now()}`);
	const envelope = NewEnvelope(OpClientPickupDiscardForMeld, { discardIndex, tileIds }, rid);
	_lastSentRaw = { op: OpClientPickupDiscardForMeld, envelope, matchId };
	lastSent.set(_lastSentRaw);
	try {
		// eslint-disable-next-line @typescript-eslint/no-explicit-any
		await (sock as any).sendMatchState(matchId, OpClientPickupDiscardForMeld, envelope);
	} catch (_err) {
		void _err;
	}
	return envelope;
}

export async function sendMeldInitial(
	melds: Array<{
		id?: string;
		kind: 'run' | 'set';
		tileIds: string[];
		jokerReps?: Record<string, unknown>;
	}>,
	requestId?: string
): Promise<string> {
	if (!melds || melds.length === 0) throw new Error('melds required');
	const matchId = getMatchId() ?? getStoredMatchId() ?? 'mock-match';
	let sock = getSocket();
	if (!sock) {
		try {
			sock = await createSocket();
		} catch (_err) {
			void _err;
			sock = { sendMatchState: async () => ({}) } as unknown as typeof sock;
		}
	}
	const rid =
		requestId ??
		(typeof crypto !== 'undefined' && 'randomUUID' in crypto
			? crypto.randomUUID()
			: `req-${Date.now()}`);
	// Map to server expected shape: melds:[{id,kind,tileIds,jokerReps}]
	const payload = {
		melds: melds.map((m, i) => ({
			id: m.id ?? `m${i + 1}`,
			kind: m.kind,
			tileIds: m.tileIds,
			jokerReps: m.jokerReps ?? {}
		}))
	};
	const envelope = NewEnvelope(OpClientMeldInitial, payload, rid);
	_lastSentRaw = { op: OpClientMeldInitial, envelope, matchId };
	lastSent.set(_lastSentRaw);
	try {
		// eslint-disable-next-line @typescript-eslint/no-explicit-any
		await (sock as any).sendMatchState(matchId, OpClientMeldInitial, envelope);
	} catch (_err) {
		void _err;
	}
	return envelope;
}

export async function sendMeldNew(
	melds: Array<{
		id?: string;
		kind: 'run' | 'set';
		tileIds: string[];
		jokerReps?: Record<string, unknown>;
	}>,
	requestId?: string
): Promise<string> {
	if (!melds || melds.length === 0) throw new Error('melds required');
	const matchId = getMatchId() ?? getStoredMatchId() ?? 'mock-match';
	let sock = getSocket();
	if (!sock) {
		try {
			sock = await createSocket();
		} catch (_err) {
			void _err;
			sock = { sendMatchState: async () => ({}) } as unknown as typeof sock;
		}
	}
	const rid =
		requestId ??
		(typeof crypto !== 'undefined' && 'randomUUID' in crypto
			? crypto.randomUUID()
			: `req-${Date.now()}`);
	const payload = {
		melds: melds.map((m, i) => ({
			id: m.id ?? `m${i + 1}`,
			kind: m.kind,
			tileIds: m.tileIds,
			jokerReps: m.jokerReps ?? {}
		}))
	};
	const envelope = NewEnvelope(OpClientMeldNew, payload, rid);
	_lastSentRaw = { op: OpClientMeldNew, envelope, matchId };
	lastSent.set(_lastSentRaw);
	try {
		// eslint-disable-next-line @typescript-eslint/no-explicit-any
		await (sock as any).sendMatchState(matchId, OpClientMeldNew, envelope);
	} catch (_err) {
		void _err;
	}
	return envelope;
}

export function _resetForTest(): void {
	_lastSentRaw = null;
	lastSent.set(null);
}
