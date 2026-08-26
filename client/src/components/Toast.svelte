<script lang="ts">
	import { errorStore, clearError } from '$lib/game/errorStore';
	import { fade, fly } from 'svelte/transition';

	const err = $derived($errorStore);
</script>

{#if err}
	<div
		data-testid="toast"
		data-error-code={err.code}
		data-request-id={err.requestId ?? ''}
		data-op={err.op ?? ''}
		role="alert"
		aria-live="assertive"
		in:fly={{ y: 20, duration: 200 }}
		out:fade={{ duration: 200 }}
		class="pointer-events-auto fixed bottom-4 left-1/2 z-50 flex max-w-[90vw] -translate-x-1/2 flex-col gap-1 rounded-xl bg-[#dc2626] px-4 py-3 text-sm font-bold text-white shadow-xl"
	>
		<div class="flex items-center gap-2">
			<span class="rounded bg-white px-1.5 py-0.5 text-[10px] font-black text-[#dc2626]"
				>{err.code}</span
			>
			<span class="flex-1">{err.message}</span>
			<button
				onclick={clearError}
				aria-label="Close"
				class="rounded bg-white/20 px-2 py-0.5 text-xs text-white hover:bg-white/30">✕</button
			>
		</div>
		{#if err.details && Object.keys(err.details).length > 0}
			<div class="text-xs font-normal opacity-90">
				{Object.entries(err.details)
					.map(([k, v]) => `${k}: ${v}`)
					.join(' • ')}
			</div>
		{/if}
		{#if err.requestId}
			<div class="font-mono text-[10px] opacity-70">req {err.requestId} • op {err.op ?? ''}</div>
		{/if}
	</div>
{/if}
