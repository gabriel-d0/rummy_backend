<script lang="ts">
	import { errorStore, onServerError, clearError } from '$lib/game/errorStore';
	import { handleMatchData } from '$lib/game/store';
	import { OpServerError } from '$lib/nakama/protocol';

	function triggerBadPayload() {
		handleMatchData(
			OpServerError,
			JSON.stringify({
				code: 'bad_payload',
				message: 'invalid tileId format',
				details: { tileId: 'bad-id', field: 'tileId' },
				requestId: 'req-123',
				op: 2
			})
		);
	}

	function triggerNotYourTurn() {
		onServerError({
			code: 'not_your_turn',
			message: 'not your turn',
			details: { currentSeat: '1', yourSeat: '0' },
			requestId: 'req-124',
			op: 2
		});
	}

	function triggerBadRequest() {
		onServerError({
			code: 'bad_request',
			message: 'invalid discard tileId not in rack',
			requestId: 'req-125',
			op: 2
		});
	}

	function triggerLeaked() {
		onServerError({
			code: 'LEAKED',
			message: 'private data leak detected',
			details: { leaked: 'OwnRack' },
			requestId: 'req-leak',
			op: 102
		});
	}

	function triggerWithDetails() {
		onServerError({
			code: 'bad_payload',
			message: 'meld invalid',
			details: { reason: 'duplicate tile', meldId: 'm1' },
			requestId: 'req-999',
			op: 6
		});
	}

	function doClear() {
		clearError();
	}

	let err = $derived($errorStore);
</script>

<div class="min-h-screen bg-[#0a2e1a] p-4">
	<h1 class="font-bold text-white">Error Toast Demo — Day 45</h1>
	<p class="text-sm text-white/60">
		OpServerError 102 code/message/details/requestId/op 3s bg #dc2626 • Toast via errorStore
	</p>
	<div class="mt-3 flex flex-wrap gap-2">
		<button onclick={triggerBadPayload} class="rounded bg-red-600 px-3 py-1 text-sm text-white"
			>Trigger bad_payload</button
		>
		<button onclick={triggerNotYourTurn} class="rounded bg-amber-600 px-3 py-1 text-sm text-white"
			>Trigger not_your_turn</button
		>
		<button onclick={triggerBadRequest} class="rounded bg-orange-600 px-3 py-1 text-sm text-white"
			>Trigger bad_request</button
		>
		<button onclick={triggerLeaked} class="rounded bg-violet-600 px-3 py-1 text-sm text-white"
			>Trigger LEAKED</button
		>
		<button onclick={triggerWithDetails} class="rounded bg-sky-600 px-3 py-1 text-sm text-white"
			>Trigger with details</button
		>
		<button onclick={doClear} class="rounded bg-white/10 px-3 py-1 text-sm text-white">Clear</button
		>
	</div>
	<div class="mt-3 text-xs text-white/80" data-testid="error-info">
		code:{err?.code ?? 'none'} message:{err?.message ?? 'none'} requestId:{err?.requestId ?? 'none'} op:{String(
			err?.op ?? 'none'
		)}
	</div>
	<div class="mt-2 text-[11px] text-white/40" data-testid="error-details">
		{err?.details ? JSON.stringify(err.details) : 'no details'}
	</div>
	<div class="mt-8 text-xs text-white/40">
		Toast appears fixed bottom center, bg #dc2626, shows code • message • details • req/op,
		auto-dismiss 3s, manual ✕. Tested via OpServerError 102 handleMatchData.
	</div>
</div>
