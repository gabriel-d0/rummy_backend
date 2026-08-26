import { Client, Session } from "@heroiclabs/nakama-js";

// Day 3: Nakama JS client — device auth via defaultkey at 127.0.0.1:7350, lazy socket

const HOST = (import.meta as unknown as { env?: Record<string, string> })?.env?.VITE_NAKAMA_HOST ?? "127.0.0.1";
const PORT = (import.meta as unknown as { env?: Record<string, string> })?.env?.VITE_NAKAMA_PORT ?? "7350";
const KEY = (import.meta as unknown as { env?: Record<string, string> })?.env?.VITE_NAKAMA_KEY ?? "defaultkey";
const USE_SSL = (import.meta as unknown as { env?: Record<string, string> })?.env?.VITE_NAKAMA_USE_SSL === "true";

function getOrCreateDeviceId(): string {
  const key = "rummy_device_id";
  let id = localStorage.getItem(key);
  if (!id) {
    if (typeof crypto !== "undefined" && "randomUUID" in crypto) {
      id = (crypto as unknown as { randomUUID: () => string }).randomUUID();
    } else {
      id = "xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx".replace(/[xy]/g, (c) => {
        const r = (Math.random() * 16) | 0;
        const v = c === "x" ? r : (r & 0x3) | 0x8;
        return v.toString(16);
      });
    }
    localStorage.setItem(key, id!);
  }
  return id!;
}

let client: Client | null = null;
let session: Session | null = null;
let socket: unknown = null;

export function getClient(): Client {
  if (!client) {
    client = new Client(KEY, HOST, PORT, USE_SSL);
    console.log(`Nakama Client ${HOST}:${PORT} key=${KEY} ssl=${USE_SSL} — Day 3`);
  }
  return client;
}

export async function authenticate(username?: string): Promise<Session> {
  const c = getClient();
  const deviceId = getOrCreateDeviceId();
  const uname = username ?? `rummy-${deviceId.slice(0, 8)}`;
  session = await c.authenticateDevice(deviceId, true, uname);
  localStorage.setItem("rummy_token", session!.token as string);
  localStorage.setItem("rummy_userId", session!.user_id as string);
  console.log(`Authenticated ${uname} deviceId=${deviceId} userId=${session!.user_id} — Day 3`);
  return session!;
}

export function getSession(): Session | null {
  if (session) return session;
  const token = localStorage.getItem("rummy_token");
  if (!token) return null;
  return null;
}

export async function createSocket(): Promise<unknown> {
  const c = getClient();
  const s = getSession() ?? (await authenticate());
  if (!socket) {
    socket = c.createSocket(false, false);
    await (socket as { connect: (s: Session, b: boolean) => Promise<void> }).connect(s, true);
    console.log("Socket connected — Day 3");
  }
  return socket;
}
