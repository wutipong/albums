<script lang="ts">
	import AssetThumbnail from '$lib/components/AssetThumbnail.svelte';
	import AssetViewer from '$lib/components/AssetViewer.svelte';
	import Icon from 'mdi-svelte';
	import type { PageProps } from './$types';
	import { mdiDownload, mdiImageAlbum } from '@mdi/js';
	import NavBar from '$lib/components/NavBar.svelte';
	import Toast from '$lib/components/Toast.svelte';

	let { data, params }: PageProps = $props();
	let asset = $state({
		id: '<placeholder>',
		type: 'image',
		view_url: 'http://example.com',
		filename: ''
	});
	let showViewer = $state(false);
	let currentIndex = $state(0);

	let nextIndex = $state(-1)
	let prevIndex = $state(-1)

	let toast: Toast;

	function findPrevious(assets: any[], index: number): number {
		return assets.slice(0, index).findLastIndex((asset: any, index, arr) => {
			if (asset == undefined) return false;
			return asset.process_status == 'processed';
		});
	}

	function findNext(assets: any[], index: number): number {
		return assets.findIndex((asset: any, i, arr) => {
			if (i <= index) return false
			if (asset == undefined) return false;
			return asset.process_status == 'processed';
		});
	}

	async function setAlbumCover(albumId: string, assetId: string) {
		await fetch(`/api/album/${albumId}/cover`, {
			method: 'POST',
			body: JSON.stringify({ asset_id: assetId })
		});

		toast.add(
			'Album cover change has been queued. It will take some time before the change is applied.',
			'info'
		);
	}

	function next() {
		if (nextIndex === -1) {
			return
		} 

		currentIndex = nextIndex;
		onIndexUpdated(currentIndex);
	}

	function previous() {
		if (prevIndex === -1) {
			return
		} 
		currentIndex = prevIndex
		onIndexUpdated(currentIndex);
	}

	function onIndexUpdated(index: number) {
		nextIndex = findNext(data.assets, index);
		prevIndex = findPrevious(data.assets, index);

		currentIndex = index;
		asset = data.assets[index];
	}
</script>

{#snippet title()}
	<div class="flex text-xl md:ms-4">
		<Icon path={mdiImageAlbum}></Icon>
		{data.name}
	</div>
{/snippet}

<div class="relative flex h-screen w-screen flex-col">
	<NavBar {title} album={data}></NavBar>
	<div class="mx-4 overflow-auto pt-8">
		<div class="flex flex-wrap justify-evenly gap-1">
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
	<AssetViewer
		bind:assetType={asset.type}
		bind:viewURL={asset.view_url}
		bind:filename={asset.filename}
		bind:show={showViewer}
		{next}
		{previous}
		hasNext={nextIndex != -1}
		hasPrevious={prevIndex != -1}
		menu={viewMenu}
	/>
</div>

<Toast bind:this={toast} />

{#snippet viewMenu()}
	<li>
		<a href={`/api/asset/${asset.id}/original/`} target="_blank">
			<Icon path={mdiDownload} /> Download.
		</a>
	</li>
	<li>
		<button
			onclick={() => {
				setAlbumCover(params.id, asset.id);
			}}
		>
			<Icon path={mdiImageAlbum} /> Set as album cover.
		</button>
	</li>
{/snippet}
