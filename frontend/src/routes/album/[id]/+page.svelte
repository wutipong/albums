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
	let hasNext = $state(true);
	let hasPrevious = $state(true);

	let toast: Toast;

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
	<div class="flex text-xl md:ms-4">
		<Icon path={mdiImageAlbum}></Icon>
		{data.name}
	</div>
{/snippet}

<div class="relative flex h-screen w-screen flex-col bg-base-100">
	<NavBar {title}></NavBar>
	<div class="mx-4 overflow-auto pt-8">
		<div class="flex flex-wrap justify-evenly gap-4">
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
		{hasNext}
		{hasPrevious}
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
