// Day 3 — Design tokens — single source for colours, fonts, spacing, radius
// No hard-coded hex in components; import from here.

export const colors = {
  felt: "#0a2e1a",
  feltLight: "#0f5d2e",
  wood: "#5d4037",
  woodLight: "#6d4c41",
  tile: {
    red: "#e53935",
    yellow: "#f9a825",
    blue: "#1e88e5",
    black: "#212121",
  },
} as const;

export const fonts = {
  sans: "Inter, system-ui, -apple-system, sans-serif",
  mono: "JetBrains Mono, ui-monospace, monospace",
} as const;

export const spacing = {
  xs: "4px",
  sm: "8px",
  md: "12px",
  lg: "16px",
  xl: "24px",
} as const;

export const radius = {
  sm: "8px",
  md: "12px",
  lg: "16px",
  xl: "18px",
} as const;

export const tile = {
  rack: { w: 64, h: 90 },
  table: { w: 52, h: 72 },
} as const;
