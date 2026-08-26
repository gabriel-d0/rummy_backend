<script lang="ts">
	// Day 5 placeholder — Tile will render 1..13 + Joly, red/yellow/blue/black, rack vs table size
	let { colour = 1, rank = 1 } = $props<{ colour?: number; rank?: number }>();
</script>

<div class="grid h-[90px] w-[64px] place-items-center rounded-lg border bg-white text-xs">
	{rank}-{colour}
</div>
