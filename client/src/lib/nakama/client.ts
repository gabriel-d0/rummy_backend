import { Client, Session } from '@heroiclabs/nakama-js';

// Day 4: Nakama JS client — hardcoded 127.0.0.1:7350 defaultkey (env wiring Day 6)
const HOST = '127.0.0.1';
const PORT = '7350';
const KEY = 'defaultkey';
const USE_SSL = false;

function getOrCreateDeviceId(): string {
	const key = 'rummy_device_id';
	let id: string | null = null;
	try {
		id = localStorage.getItem(key);
	} catch (_err) {
		void _err;
	}
	if (!id) {
		if (typeof crypto !== 'undefined' && 'randomUUID' in crypto) {
			id = (crypto as unknown as { randomUUID: () => string }).randomUUID();
		} else {
			id = 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, (c) => {
				const r = (Math.random() * 16) | 0;
				const v = c === 'x' ? r : (r & 0x3) | 0x8;
				return v.toString(16);
			});
		}
		try {
			localStorage.setItem(key, id!);
		} catch (_err2) {
			void _err2;
		}
	}
	return id!;
}

let client: Client | null = null;
let session: Session | null = null;
let socket: unknown = null;

export function getClient(): Client {
	if (!client) {
		client = new Client(KEY, HOST, PORT, USE_SSL);
	}
	return client;
}

export async function authenticate(username?: string): Promise<Session> {
	const c = getClient();
	const deviceId = getOrCreateDeviceId();
	const uname = username ?? `rummy-${deviceId.slice(0, 8)}`;
	session = await c.authenticateDevice(deviceId, true, uname);
	try {
		localStorage.setItem('rummy_token', session!.token as string);
		localStorage.setItem('rummy_userId', session!.user_id as string);
	} catch (_err) {
		void _err;
	}
	return session!;
}

export function getSession(): Session | null {
	if (session) return session;
	try {
		const token = localStorage.getItem('rummy_token');
		if (!token) return null;
	} catch (_err) {
		void _err;
	}
	return null;
}

export async function createSocket(): Promise<unknown> {
	const c = getClient();
	const s = (getSession() ?? (await authenticate())) as Session;
	if (!socket) {
		socket = c.createSocket(false, false);
		await (socket as { connect: (sess: Session, b: boolean) => Promise<void> }).connect(s, true);
	}
	return socket;
}
