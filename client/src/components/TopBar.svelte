<script lang="ts">
  let { players = 4, masa = 1, seconds = 0 } = $props<{ players?: number; masa?: number; seconds?: number }>();
</script>

<header class="h-12 bg-black/90 backdrop-blur border-b border-white/10 flex items-center justify-between px-3 sm:px-4">
  <div class="flex items-center gap-3">
    <div class="w-8 h-8 rounded-lg bg-amber-400 grid place-items-center font-black text-black text-sm">R</div>
    <div>
      <div class="text-sm font-black tracking-tight text-white leading-none flex items-center gap-2">
        REMI ETALAT <span class="text-[10px] bg-amber-400 text-black px-1.5 py-0.5 rounded font-bold">ETALAT</span>
      </div>
      <div class="text-[10px] leading-none text-white/60 tracking-widest">PREMIUM • ONLINE</div>
    </div>
  </div>
  <div class="flex items-center gap-2">
    <div class="hidden sm:flex items-center gap-1.5 text-xs bg-white/10 border border-white/10 px-3 py-1.5 rounded-full">
      <span class="w-2 h-2 rounded-full bg-emerald-500 animate-pulse"></span> MASA {masa} • {players} JUCĂTORI
    </div>
    <div class="text-xs bg-white/10 px-2.5 py-1.5 rounded-full">🕒 {seconds}s</div>
    <button class="hidden sm:block text-xs bg-white text-black px-3 py-1.5 rounded-full font-bold">REGULI</button>
    <button class="text-xs bg-amber-400 text-black px-3 py-1.5 rounded-full font-bold">JOC NOU</button>
  </div>
</header>
