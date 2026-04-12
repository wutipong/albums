<script lang="ts">
	import AssetThumbnail from '$lib/components/AssetThumbnail.svelte';
	import ItemViewer from '$lib/components/ItemViewer.svelte';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();
	let currentId = $derived(data.assets[0]);
	let viewer: ItemViewer;
	let showViewer = $state(false)
</script>

<div class="relative flex h-screen w-screen flex-col">
	<div class="text-xl font-bold">Album {data.name}</div>

	<div class="overflow-auto">
		<div class="flex flex-wrap bg-base-300">
			{#each data.assets as asset (asset)}
				<AssetThumbnail
					id={asset}
					onclick={(id: string) => {
						currentId = id;
						showViewer = true;
					}}
				/>
			{/each}
		</div>
	</div>
	<ItemViewer id={currentId} bind:show={showViewer} ></ItemViewer>
</div>
