// Day 28: ErrorToast — shows OpServerError 102 code/message/requestId/op as red toast for 3s
// Mirrors Go internal/protocol/errors.go:12 OpServerError 102 with codes bad_payload/not_your_turn/wrong_phase etc.
// Logs not_your_turn and other error codes, never leaks private rack.

export interface ServerError {
  code: string;
  message: string;
  requestId?: string;
  op?: number;
  details?: unknown;
}

let lastError: ServerError | null = null;

export function getLastError(): ServerError | null {
  return lastError;
}

export function clearLastError(): void {
  lastError = null;
}

function getOrCreateContainer(): HTMLElement | null {
  if (typeof document === "undefined") return null;
  let container = document.getElementById("error-toasts");
  if (!container) {
    container = document.createElement("div");
    container.id = "error-toasts";
    container.style.cssText =
      "position:fixed;top:12px;right:12px;z-index:9999;display:flex;flex-direction:column;align-items:flex-end;pointer-events:none;max-width:min(420px,90vw);";
    document.body.appendChild(container);
  }
  return container;
}

export function showErrorToast(error: ServerError): void {
  lastError = error;
  const code = error.code ?? "unknown";
  const msg = error.message ?? "";
  const req = error.requestId ? ` req:${error.requestId}` : "";
  const op = error.op !== undefined ? ` op:${error.op}` : "";
  console.log(
    `received op 102 OpServerError code=${code} message=${msg} requestId=${error.requestId} op=${error.op} — Day 28`
  );

  const container = getOrCreateContainer();
  if (!container) return;

  const toast = document.createElement("div");
  toast.setAttribute("role", "alert");
  toast.setAttribute("data-error-code", code);
  if (error.op !== undefined) toast.setAttribute("data-error-op", String(error.op));
  if (error.requestId) toast.setAttribute("data-request-id", error.requestId);
  toast.textContent = `${code}: ${msg}${op}${req}`;
  toast.style.cssText =
    "background:#dc2626;color:#ffffff;padding:10px 14px;margin:4px;border-radius:6px;font-family:monospace;font-size:12px;line-height:1.4;opacity:0.97;box-shadow:0 2px 8px rgba(0,0,0,0.35);pointer-events:auto;max-width:100%;word-break:break-word;transition:opacity 0.3s ease;";
  // Keep at most 5 toasts visible
  const children = (container as unknown as { children?: { length: number } }).children;
  if (children && children.length >= 5) {
    while ((container as unknown as { children: { length: number } }).children.length >= 5) {
      container.firstChild?.remove();
    }
  }
  container.appendChild(toast);

  // Auto-remove after 3s with fade
  setTimeout(() => {
    toast.style.opacity = "0";
    setTimeout(() => toast.remove(), 300);
  }, 3000);
}

export function clearErrorToasts(): void {
  lastError = null;
  if (typeof document === "undefined") return;
  const container = document.getElementById("error-toasts");
  if (container) container.innerHTML = "";
}

// Phaser integration helper — if a Scene is available, also show a Phaser text toast
// For Day 28 we primarily use DOM; this helper is for future canvas-only mode.
export function showErrorToastInScene(scene: Phaser.Scene, error: ServerError): void {
  showErrorToast(error);
  try {
    const w = scene.scale.width;
    const text = scene.add.text(w / 2, 24, `${error.code}: ${error.message}`, {
      fontFamily: "monospace",
      fontSize: "12px",
      color: "#ffffff",
      backgroundColor: "#dc2626",
      padding: { x: 10, y: 6 },
    });
    text.setOrigin(0.5, 0);
    text.setDepth(10000);
    text.setData("isErrorToast", true);
    scene.tweens.add({
      targets: text,
      alpha: 0,
      duration: 300,
      delay: 3000,
      onComplete: () => text.destroy(),
    });
  } catch {
    // ignore if scene not ready or headless test
  }
}
