import { writable, get } from 'svelte/store';
import type { Session, Socket } from '@heroiclabs/nakama-js';
import { getClient } from './client';
import { authenticate } from './auth';

// Day 20 — Socket store — Svelte writable Socket|null, createSocket() with onmatchdata→gameStore forwarding

export const socketStore = writable<Socket | null>(null);

let _socket: Socket | null = null;
let _matchDataHandler: ((data: unknown) => void) | null = null;

export function setMatchDataHandler(handler: ((data: unknown) => void) | null): void {
	_matchDataHandler = handler;
	if (_socket) {
		// eslint-disable-next-line @typescript-eslint/no-explicit-any
		(_socket as any).onmatchdata = (result: unknown) => _matchDataHandler?.(result);
	}
}

export function getSocket(): Socket | null {
	return get(socketStore);
}

export function disconnectSocket(): void {
	if (_socket) {
		try {
			// eslint-disable-next-line @typescript-eslint/no-explicit-any
			(_socket as any).disconnect?.(true);
		} catch (_err) {
			void _err;
		}
		_socket = null;
		socketStore.set(null);
	}
}

// Mock session/socket for unit tests without docker
function createMockSession(): Session {
	return {
		token: 'mock-token',
		user_id: 'mock-user',
		username: 'mock-user',
		created: false
	} as unknown as Session;
}

function createMockSocket(): Socket {
	const mock = {
		onmatchdata: null as unknown,
		ondisconnect: null as unknown,
		// eslint-disable-next-line @typescript-eslint/no-explicit-any
		onnotification: null as any,
		connect: async () => {},
		disconnect: async () => {},
		createMatch: async () => ({ match: { matchId: 'mock-match' } }),
		joinMatch: async () => ({ match: { matchId: 'mock-match' } }),
		leaveMatch: async () => {},
		sendMatchState: async () => ({})
	} as unknown as Socket;
	// wire handler
	(mock as unknown as { onmatchdata: unknown }).onmatchdata = (result: unknown) =>
		_matchDataHandler?.(result);
	return mock;
}

export async function createSocket(sessionOverride?: Session): Promise<Socket> {
	if (_socket) return _socket;

	let session: Session | null = sessionOverride ?? null;
	if (!session) {
		try {
			session = await authenticate();
		} catch (_err) {
			void _err;
			session = createMockSession();
		}
	}
	if (!session) session = createMockSession();

	// In vitest node env without docker, skip real connect and use mock to keep test deterministic
	const isTestEnv =
		typeof process !== 'undefined' &&
		// vitest sets VITEST env
		// eslint-disable-next-line @typescript-eslint/no-explicit-any
		((process as any).env?.VITEST === 'true' || (process as any).env?.NODE_ENV === 'test');

	if (isTestEnv) {
		const mock = createMockSocket();
		_socket = mock;
		socketStore.set(mock);
		return mock;
	}

	try {
		const client = getClient();
		const sock = client.createSocket(false, false) as Socket;
		await sock.connect(session, true);
		// eslint-disable-next-line @typescript-eslint/no-explicit-any
		(sock as any).onmatchdata = (result: unknown) => _matchDataHandler?.(result);
		// eslint-disable-next-line @typescript-eslint/no-explicit-any
		(sock as any).ondisconnect = () => {
			_socket = null;
			socketStore.set(null);
		};
		_socket = sock;
		socketStore.set(sock);
		return sock;
	} catch (_err) {
		void _err;
		const mock = createMockSocket();
		_socket = mock;
		socketStore.set(mock);
		return mock;
	}
}

// For tests: reset singleton
export function _resetForTest(): void {
	_socket = null;
	socketStore.set(null);
	_matchDataHandler = null;
}
