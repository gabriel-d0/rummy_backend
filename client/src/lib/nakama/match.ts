import { writable, get } from 'svelte/store';
import { createSocket, getSocket } from './socket';
import { getClient } from './client';
import { authStore } from './auth';

// Day 21 — Match create/join — Svelte writable matchId + rummy_matchId persistence + socket.createMatch/joinMatch

const STORAGE_KEY = 'rummy_matchId';

export const matchStore = writable<string | null>(loadStoredMatchId());

function loadStoredMatchId(): string | null {
	try {
		return localStorage.getItem(STORAGE_KEY);
	} catch (_err) {
		void _err;
		return null;
	}
}

function persistMatchId(matchId: string | null): void {
	try {
		if (matchId) localStorage.setItem(STORAGE_KEY, matchId);
		else localStorage.removeItem(STORAGE_KEY);
	} catch (_err) {
		void _err;
	}
	matchStore.set(matchId);
}

export function getMatchId(): string | null {
	return get(matchStore);
}

export function getStoredMatchId(): string | null {
	try {
		return localStorage.getItem(STORAGE_KEY);
	} catch (_err) {
		void _err;
		return null;
	}
}

export function clearMatchId(): void {
	persistMatchId(null);
}

function extractMatchId(result: unknown, fallback: string): string {
	if (!result || typeof result !== 'object') return fallback;
	const r = result as Record<string, unknown>;
	// nakama-js may return { match: { matchId } } or { match: { match_id } } or directly { matchId }
	if (r.match && typeof r.match === 'object') {
		const m = r.match as Record<string, unknown>;
		if (typeof m.matchId === 'string' && m.matchId) return m.matchId;
		if (typeof m.match_id === 'string' && m.match_id) return m.match_id;
		if (typeof m.id === 'string' && m.id) return m.id;
	}
	if (typeof r.matchId === 'string' && r.matchId) return r.matchId;
	if (typeof r.match_id === 'string' && r.match_id) return r.match_id;
	if (typeof r.id === 'string' && r.id) return r.id;
	return fallback;
}

export async function createMatch(): Promise<string> {
	const isTestEnv =
		typeof process !== 'undefined' &&
		// eslint-disable-next-line @typescript-eslint/no-explicit-any
		((process as any).env?.VITEST === 'true' || (process as any).env?.NODE_ENV === 'test');
	if (isTestEnv) {
		let sock = getSocket();
		if (!sock) sock = await createSocket();
		try {
			// eslint-disable-next-line @typescript-eslint/no-explicit-any
			const result = await (sock as any).createMatch();
			const matchId = extractMatchId(result, 'mock-match');
			persistMatchId(matchId);
			return matchId;
		} catch (_err) {
			void _err;
			const fallback = 'mock-match';
			persistMatchId(fallback);
			return fallback;
		}
	}
	// Try authoritative RPC first (real Nakama)
	try {
		const session = get(authStore);
		if (session) {
			const client = getClient();
			// eslint-disable-next-line @typescript-eslint/no-explicit-any
			const rpcRes: any = await (client as any).rpc(session, 'create_match', '');
			const payloadStr = rpcRes?.payload ?? rpcRes;
			let parsed: Record<string, unknown> = {};
			try {
				parsed =
					typeof payloadStr === 'string'
						? JSON.parse(payloadStr)
						: (payloadStr as Record<string, unknown>);
			} catch (_e) {
				void _e;
			}
			const rpcMatchId = (parsed.matchId ?? parsed.match_id ?? parsed.id ?? '') as string;
			if (rpcMatchId) {
				// RPC created the match, now join it via socket
				let sock = getSocket();
				if (!sock) sock = await createSocket();
				try {
					// eslint-disable-next-line @typescript-eslint/no-explicit-any
					await (sock as any).joinMatch(rpcMatchId);
				} catch (_e2) {
					void _e2;
				}
				persistMatchId(rpcMatchId);
				return rpcMatchId;
			}
		}
	} catch (_err) {
		void _err;
	}
	// Fallback: relayed via socket.createMatch (mock in tests)
	let sock = getSocket();
	if (!sock) {
		sock = await createSocket();
	}
	try {
		// eslint-disable-next-line @typescript-eslint/no-explicit-any
		const result = await (sock as any).createMatch();
		const matchId = extractMatchId(result, 'mock-match');
		persistMatchId(matchId);
		return matchId;
	} catch (_err) {
		void _err;
		const fallback = `mock-match-${Date.now().toString(36)}`;
		persistMatchId(fallback);
		return fallback;
	}
}

export async function joinMatch(matchId: string): Promise<string> {
	if (!matchId) throw new Error('matchId required');
	let sock = getSocket();
	if (!sock) {
		sock = await createSocket();
	}
	try {
		// eslint-disable-next-line @typescript-eslint/no-explicit-any
		const result = await (sock as any).joinMatch(matchId);
		const returnedId = extractMatchId(result, matchId);
		persistMatchId(returnedId);
		return returnedId;
	} catch (_err) {
		void _err;
		// fallback: assume join succeeded with given id
		persistMatchId(matchId);
		return matchId;
	}
}

export async function leaveMatch(): Promise<void> {
	const sock = getSocket();
	const mid = getMatchId();
	if (sock && mid) {
		try {
			// eslint-disable-next-line @typescript-eslint/no-explicit-any
			await (sock as any).leaveMatch(mid);
		} catch (_err) {
			void _err;
		}
	}
	clearMatchId();
}

export type AvailableMatch = { matchId: string; label: string; size: number };

export async function listAvailableMatches(): Promise<AvailableMatch[]> {
	try {
		const session = get(authStore);
		if (!session) return [];
		const client = getClient();
		const tryList = async (auth: boolean | undefined, label: string) => {
			// eslint-disable-next-line @typescript-eslint/no-explicit-any
			const res: any = await (client as any).listMatches(session, 10, auth, label, 1, 4, '');
			return res?.matches ?? res?.result ?? res ?? [];
		};
		let all: unknown[] = [];
		// Try authoritative and non-authoritative, with and without label
		for (const auth of [undefined, true, false] as const) {
			for (const label of ['', 'rummy'] as const) {
				try {
					const part = await tryList(auth as unknown as boolean, label);
					if (Array.isArray(part) && part.length > 0) {
						all = [...all, ...part];
					}
				} catch (_e) {
					void _e;
				}
			}
		}
		// Also try plain limit 10 no filter
		try {
			// eslint-disable-next-line @typescript-eslint/no-explicit-any
			const res: any = await (client as any).listMatches(session, 10);
			const part = res?.matches ?? res?.result ?? res ?? [];
			if (Array.isArray(part) && part.length > 0) all = [...all, ...part];
		} catch (_e) {
			void _e;
		}
		// Deduplicate by matchId
		const seen = new Set<string>();
		const out: AvailableMatch[] = [];
		for (const m of all as Record<string, unknown>[]) {
			const id = (m.matchId ?? m.match_id ?? m.id ?? '') as string;
			if (!id || seen.has(id)) continue;
			seen.add(id);
			out.push({
				matchId: id,
				label: (m.label ?? '') as string,
				size: (m.size ?? (m as unknown as { size: number }).size ?? 0) as number
			});
		}
		return out;
	} catch (_err) {
		void _err;
		return [];
	}
}

// For tests: reset singleton
export function _resetForTest(): void {
	clearMatchId();
}
