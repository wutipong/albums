<script lang="ts">
	import AssetThumbnail from '$lib/components/AssetThumbnail.svelte';
	import ItemViewer from '$lib/components/ItemViewer.svelte';
	import Icon from 'mdi-svelte';
	import type { PageProps } from './$types';
	import { mdiDownload, mdiImageAlbum, mdiImageSearch, mdiImageSearchOutline } from '@mdi/js';
	import NavBar from '$lib/components/NavBar.svelte';

	let { data, params }: PageProps = $props();
	let asset = $state({ id: '<placeholder>', album_id: '' });
	let showViewer = $state(false);
	let currentIndex = $state(0);
	let hasNext = $state(true);
	let hasPrevious = $state(true);

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
		asset = data.assets[index];
	}
</script>

{#snippet title()}
	<div class="flex text-xl">
		<Icon path={mdiImageSearchOutline}></Icon>
		{data.search}
	</div>
{/snippet}

<div class="relative flex h-screen w-screen flex-col bg-base-300">
	<NavBar {title} />
	<div class="overflow-auto mt-8">
		<div class="flex flex-wrap">
			{#each data.assets as asset, index (asset)}
				<AssetThumbnail
					{asset}
					onclick={(asset: any) => {
						onIndexUpdated(index);
						showViewer = true;
					}}
				/>
			{/each}
		</div>
	</div>
	<ItemViewer
		bind:asset
		bind:show={showViewer}
		{next}
		{previous}
		{hasNext}
		{hasPrevious}
		menu={viewMenu}
	/>
</div>

{#snippet viewMenu()}
	<a href={`/album/${asset.album_id}/`} class="btn btn-soft">
		<Icon path={mdiImageAlbum} /> View album.
	</a>

	<a href={`/api/asset/${asset.id}/original/`} target="_blank" class="btn btn-soft">
		<Icon path={mdiDownload} /> Download.
	</a>
{/snippet}
