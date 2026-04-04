<script lang="ts">
	import type { PageProps } from './$types';

	let { params, data }: PageProps = $props();
</script>

<div class="text-xl font-bold">Album {params.id}</div>

<div class="flex flex-wrap">
	{#each data.assets as asset (asset)}
		<div class="h-40 grow object-cover block m-1">
			<img src={`/api/thumbnail/${asset}`} alt={asset} class="max-w-full h-full"/>
		</div>
	{/each}
</div>
