import {
	PUBLIC_NAKAMA_HOST,
	PUBLIC_NAKAMA_PORT,
	PUBLIC_NAKAMA_KEY,
	PUBLIC_NAKAMA_USE_SSL
} from '$env/static/public';

// Day 6 — Env wiring — typed VITE_NAKAMA_* via $env/static/public with fallbacks
// PUBLIC_ prefix is SvelteKit's public env (exposed to client), VITE_ kept for compat.

export const NAKAMA_HOST = PUBLIC_NAKAMA_HOST ?? '127.0.0.1';
export const NAKAMA_PORT = PUBLIC_NAKAMA_PORT ?? '7350';
export const NAKAMA_KEY = PUBLIC_NAKAMA_KEY ?? 'defaultkey';
export const NAKAMA_USE_SSL = PUBLIC_NAKAMA_USE_SSL === 'true';

export const config = {
	host: NAKAMA_HOST,
	port: NAKAMA_PORT,
	key: NAKAMA_KEY,
	useSsl: NAKAMA_USE_SSL
} as const;
