<script lang="ts">
	import { privateStore, publicStore } from '$lib/game/store';

	const gamePhase = $derived($privateStore?.gamePhase ?? $publicStore?.gamePhase ?? '');
	const winner = $derived($privateStore?.winner ?? $publicStore?.winner ?? -1);
	const isWinner = $derived(gamePhase === 'RoundComplete' && winner >= 0);

	function restart() {
		// Day 39 — RESTART MASA placeholder: clear winner overlay (demo only)
		// Real would navigate or reset match via Nakama; for now just reload
		try {
			window.location.reload();
		} catch (_err) {
			void _err;
		}
	}
</script>

{#if isWinner}
	<div
		data-testid="winner-overlay"
		class="fixed inset-0 z-50 flex items-center justify-center bg-black/70 backdrop-blur-sm pointer-events-none"
	>
		<div class="pointer-events-auto rounded-2xl bg-white p-6 text-center shadow-2xl">
			<div class="text-sm font-bold tracking-widest text-emerald-600">CÂȘTIGĂTOR</div>
			<div class="mt-2 text-2xl font-black text-slate-900" data-testid="winner-text">
				Winner {winner}
			</div>
			<div class="mt-1 text-xs text-slate-500">Seat {winner} a câștigat masa • RoundComplete</div>
			<button
				onclick={restart}
				data-testid="restart-btn"
				class="mt-4 rounded-xl bg-amber-400 px-6 py-2.5 text-sm font-bold text-black hover:bg-amber-300"
				>RESTART MASA</button
			>
		</div>
	</div>
{/if}
