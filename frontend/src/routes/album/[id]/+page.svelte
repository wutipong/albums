<script lang="ts">
	import AssetThumbnail from '$lib/components/AssetThumbnail.svelte';
	import ItemViewer from '$lib/components/ItemViewer.svelte';
	import Icon from 'mdi-svelte';
	import type { PageProps } from './$types';
	import { mdiDownload, mdiImageAlbum } from '@mdi/js';

	let { data, params }: PageProps = $props();
	let currentId = $state('');
	let showViewer = $state(false);
	let currentIndex = $state(0);
	let hasNext = $state(true);
	let hasPrevious = $state(true);

	async function setAlbumCover(albumId: string, assetId: string) {
		await fetch(`/api/album/${albumId}/cover`, {
			method: 'POST',
			body: JSON.stringify({ asset_id: assetId })
		});
	}

	function next() {
		if (hasNext) {
			currentIndex++;
		} else {
			return;
		}
		onIndexUpdated(currentIndex);
	}

	function previous() {
		if (hasPrevious) {
			currentIndex--;
		} else {
			return;
		}

		onIndexUpdated(currentIndex);
	}

	function onIndexUpdated(index: number) {
		if (index == data.assets.length - 1) hasNext = false;
		else hasNext = true;
		if (index == 0) hasPrevious = false;
		else hasPrevious = true;

		currentIndex = index;
		currentId = data.assets[index];
	}
</script>

<div class="relative flex h-screen w-screen flex-col">
	<div class="text-xl font-bold">Album {data.name}</div>

	<div class="overflow-auto">
		<div class="flex flex-wrap bg-base-300">
			{#each data.assets as asset, index (asset)}
				<AssetThumbnail
					id={asset}
					onclick={(id: string) => {
						onIndexUpdated(index);
						showViewer = true;
					}}
				/>
			{/each}
		</div>
	</div>
	<ItemViewer
		bind:id={currentId}
		bind:show={showViewer}
		{next}
		{previous}
		{hasNext}
		{hasPrevious}
		menu={viewMenu}
	/>
</div>

{#snippet viewMenu()}
	<a href={`/api/asset/${currentId}/original/`} target="_blank" class="btn btn-soft">
		<Icon path={mdiDownload} /> Download.
	</a>

	<button
		class="btn btn-soft"
		onclick={() => {
			setAlbumCover(params.id, currentId);
		}}
	>
		<Icon path={mdiImageAlbum} /> Set as album cover.
	</button>
{/snippet}
