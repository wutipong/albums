<script lang="ts">
	import AssetThumbnail from '$lib/components/AssetThumbnail.svelte';
	import ItemViewer from '$lib/components/ItemViewer.svelte';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();
</script>

<div class="relative h-screen w-screen flex flex-col">
	<div class="text-xl font-bold">Album {data.name}</div>

	<div class="overflow-auto">
		<div class="flex flex-wrap bg-base-300">
			{#each data.assets as asset (asset)}
				<AssetThumbnail id={asset}></AssetThumbnail>
			{/each}
		</div>
	</div>
	<ItemViewer id={data.assets[0]}></ItemViewer>
</div>
