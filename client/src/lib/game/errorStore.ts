import { writable } from 'svelte/store';

export type ServerError = {
	code: string;
	message: string;
	details?: Record<string, string>;
	requestId?: string;
	op?: number;
};

export const errorStore = writable<ServerError | null>(null);

let clearTimer: ReturnType<typeof setTimeout> | null = null;

export function onServerError(raw: unknown): boolean {
	try {
		let obj: Record<string, unknown> | null = null;
		if (typeof raw === 'string') obj = JSON.parse(raw) as Record<string, unknown>;
		else if (raw instanceof Uint8Array)
			obj = JSON.parse(new TextDecoder().decode(raw)) as Record<string, unknown>;
		else if (raw && typeof raw === 'object') obj = raw as Record<string, unknown>;
		else return false;
		if (!obj || typeof obj.code !== 'string' || typeof obj.message !== 'string') return false;
		const err: ServerError = {
			code: obj.code as string,
			message: obj.message as string,
			details: (obj.details as Record<string, string>) ?? undefined,
			requestId: (obj.requestId as string) ?? undefined,
			op: (obj.op as number) ?? (obj.OpCode as number) ?? undefined
		};
		if (clearTimer) clearTimeout(clearTimer);
		errorStore.set(err);
		clearTimer = setTimeout(() => errorStore.set(null), 3000);
		return true;
	} catch (_err) {
		void _err;
		return false;
	}
}

export function clearError(): void {
	if (clearTimer) {
		clearTimeout(clearTimer);
		clearTimer = null;
	}
	errorStore.set(null);
}

export function _resetForTest(): void {
	if (clearTimer) {
		clearTimeout(clearTimer);
		clearTimer = null;
	}
	errorStore.set(null);
}
