import { Client, Session } from "@heroiclabs/nakama-js";

// Day 3, 22-24: Nakama JS client — device auth via defaultkey at 127.0.0.1:7350, socket, create/join match
// Stores token/matchId in localStorage for reconnect, lazy inits socket.

const HOST = (import.meta as any).env?.VITE_NAKAMA_HOST ?? "127.0.0.1";
const PORT = (import.meta as any).env?.VITE_NAKAMA_PORT ?? "7350";
const KEY = (import.meta as any).env?.VITE_NAKAMA_KEY ?? "defaultkey";
const USE_SSL = (import.meta as any).env?.VITE_NAKAMA_USE_SSL === "true";

function getOrCreateDeviceId(): string {
  const key = "rummy_device_id";
  let id = localStorage.getItem(key);
  if (!id) {
    // Use crypto.randomUUID if available, else fallback
    if (typeof crypto !== "undefined" && "randomUUID" in crypto) {
      id = (crypto as any).randomUUID();
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
let socket: any | null = null;

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
  // create: true will create account if not exists
  session = await c.authenticateDevice(deviceId, true, uname);
  localStorage.setItem("rummy_token", session!.token as string);
  localStorage.setItem("rummy_userId", session!.user_id as string);
  console.log(`Authenticated ${uname} deviceId=${deviceId} userId=${session!.user_id} (Day 3)`);
  return session!;
}

export function getSession(): Session | null {
  if (session) return session;
  const token = localStorage.getItem("rummy_token");
  if (!token) return null;
  // Rehydrate minimal session from token (nakama-js will validate on use)
  // For Day 3 we just return null and require re-auth if no in-memory session
  return null;
}

export async function createSocket(): Promise<any> {
  const c = getClient();
  const s = getSession() ?? (await authenticate());
  if (!socket) {
    socket = c.createSocket(false, false);
    // Day 27: Receive match state — parse Envelope and route op 100/101/102/103
    socket.onmatchdata = (matchData: any) => {
      const opCode = matchData.op_code ?? matchData.opCode ?? 0;
      const data = matchData.data;
      // Dynamically import to avoid circular deps
      import("../state/sync").then(({ handleMatchData }) => {
        handleMatchData(opCode, data);
      });
    };
    // Day 33: keep matchId and userId on disconnect for reconnection
    const keepForReconnect = () => {
      try {
        const matchId = localStorage.getItem("rummy_matchId");
        const userId = localStorage.getItem("rummy_userId");
        const token = localStorage.getItem("rummy_token");
        console.log(
          `Socket disconnected — keeping matchId=${matchId} userId=${userId} for reconnection — Day 33`
        );
        // Ensure they remain stored (they already are, but re-set to be safe)
        if (matchId) localStorage.setItem("rummy_matchId", matchId);
        if (userId) localStorage.setItem("rummy_userId", userId);
        if (token) localStorage.setItem("rummy_token", token);
      } catch {
        // ignore in test env
      }
    };
    socket.ondisconnect = keepForReconnect;
    // nakama-js may also emit onDisconnect (camelCase) or via socket.on('disconnect')
    try {
      socket.onDisconnect = keepForReconnect;
    } catch {
      // ignore
    }
    await socket.connect(s, true);
    console.log("Socket connected — Day 3 (Day 22) with onmatchdata handler Day 27");
  }
  return socket;
}

export async function createMatch(): Promise<string> {
  const sock = await createSocket();
  const match = await sock.createMatch();
  console.log(`Match created ${match.matchId} — Day 23`);
  localStorage.setItem("rummy_matchId", match.matchId);
  return match.matchId;
}

export async function joinMatch(matchId?: string): Promise<string> {
  const sock = await createSocket();
  const id = matchId ?? localStorage.getItem("rummy_matchId") ?? (await createMatch());
  const match = await sock.joinMatch(id);
  console.log(`Joined match ${match.matchId} — Day 24`);
  localStorage.setItem("rummy_matchId", match.matchId);
  return match.matchId;
}

export async function ensureAuthenticated(): Promise<Session> {
  if (session) return session;
  return authenticate();
}

// Day 34: reconnect — socket.connect + joinMatch(matchId) and expect OpServerState 100 PrivateSnapshot for that Seat only
// Mirrors Go rummy_match.go:79 — MatchLeave keeps Players/Racks, MatchJoin re-sends PrivateView to that Seat only
export async function reconnect(): Promise<string | null> {
  const storedMatchId = localStorage.getItem("rummy_matchId");
  const storedToken = localStorage.getItem("rummy_token");
  if (!storedMatchId) {
    console.log("reconnect: no stored matchId — cannot rejoin");
    return null;
  }
  if (!storedToken) {
    console.log("reconnect: no stored token — re-authenticating");
  }
  try {
    const sock = await createSocket();
    // If socket is already created but disconnected, try to re-connect
    try {
      const sess = getSession() ?? (await authenticate());
      // Only reconnect if not already connected — check for isConnected or similar
      // For Day 34 we just attempt to connect again (idempotent)
      if (typeof sock.connect === "function") {
        try {
          await sock.connect(sess, true);
        } catch {
          // ignore if already connected
        }
      }
    } catch {
      // ignore connect errors, proceed to join
    }
    const match = await sock.joinMatch(storedMatchId);
    const joinedId = match.matchId ?? storedMatchId;
    localStorage.setItem("rummy_matchId", joinedId);
    console.log(
      `Reconnected — Day 34 rejoin ${joinedId} expecting OpServerState 100 PrivateSnapshot for ownSeat`
    );
    // The server will send PrivateSnapshot via onmatchdata, which will be handled by sync.ts
    // and rehydrated in RackScene via subscribePrivateSnapshot immediate replay from localStorage per-Seat
    return joinedId;
  } catch (e) {
    console.log("reconnect failed", e);
    return null;
  }
}

export function getStoredMatchId(): string | null {
  try {
    return localStorage.getItem("rummy_matchId");
  } catch {
    return null;
  }
}

export function getStoredUserId(): string | null {
  try {
    return localStorage.getItem("rummy_userId");
  } catch {
    return null;
  }
}

export function isReconnectionAvailable(): boolean {
  try {
    return !!localStorage.getItem("rummy_matchId") && !!localStorage.getItem("rummy_userId");
  } catch {
    return false;
  }
}
