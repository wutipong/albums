<script lang="ts">
	import AssetThumbnail from '$lib/components/AssetThumbnail.svelte';
	import AssetViewer from '$lib/components/AssetViewer.svelte';
	import Icon from 'mdi-svelte';
	import type { PageProps } from './$types';
	import {
		mdiClipboardOutline,
		mdiDownload,
		mdiImageAlbum,
		mdiImageSearch,
		mdiInformationOutline
	} from '@mdi/js';
	import NavBar from '$lib/components/NavBar.svelte';
	import AssetInfoDialog from '$lib/components/AssetInfoDialog.svelte';
	import Toast from '$lib/components/Toast.svelte';
	import { copyImageToClipboard } from '$lib/clipboard';

	let { data, params }: PageProps = $props();
	let asset: any = $state({ id: '<placeholder>', album_id: '' });
	let showViewer = $state(false);
	let currentIndex = $state(0);
	let hasNext = $state(true);
	let hasPrevious = $state(true);
	let assetInfoDialog: AssetInfoDialog;
	let toast: Toast;

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

<svelte:head>
	<title>Albums: Search -- {data.search}</title>
</svelte:head>

<div class="relative flex h-screen w-screen flex-col">
	<NavBar user={data.user} />
	<div class="flex w-full gap-2 border-1 border-base-300 bg-base-300 p-2 shadow">
		<span class="label">
			<Icon path={mdiImageSearch} />
			Search: {data.search}
		</span>
	</div>
	<div class="overflow-auto p-4 pt-8">
		<div class="flex flex-wrap justify-evenly">
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
		assetType={data.assets[currentIndex].type}
		viewURL={data.assets[currentIndex].view_url}
		filename={data.assets[currentIndex].filename}
		bind:show={showViewer}
		{next}
		{previous}
		{hasNext}
		{hasPrevious}
		menu={viewMenu}
	/>
</div>

<Toast bind:this={toast} />

<AssetInfoDialog bind:this={assetInfoDialog} />

{#snippet viewMenu()}
	<li>
		<button
			onclick={() => {
				assetInfoDialog.show(asset);
			}}
		>
			<Icon path={mdiInformationOutline} /> Asset information
		</button>
	</li>
	<li>
		<a href={`/album/${asset.album_id}/`}>
			<Icon path={mdiImageAlbum} /> View album
		</a>
	</li>
	<li>
		<button
			disabled={asset.copy_url == undefined || asset.copy_url == ''}
			onclick={() => {
				copyImageToClipboard(asset.copy_url);
				toast.add('Image copied to clipboard', 'success');
			}}
		>
			<Icon path={mdiClipboardOutline} /> Copy to clipboard
		</button>
	</li>
	<li>
		<a href={asset.original_url} target="_blank">
			<Icon path={mdiDownload} /> Download
		</a>
	</li>
{/snippet}
