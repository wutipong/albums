<script lang="ts">
	import AssetThumbnail from '$lib/components/AssetThumbnail.svelte';
	import ItemViewer from '$lib/components/ItemViewer.svelte';
	import Icon from 'mdi-svelte';
	import type { PageProps } from './$types';
	import { mdiDownload, mdiImageAlbum } from '@mdi/js';

	let { data, params }: PageProps = $props();
	let currentId = $derived(data.assets[0]);
	let showViewer = $state(false);

	async function setAlbumCover(albumId :string, assetId :string){
		await fetch(`/api/album/${albumId}/cover`, {
			method: "POST",
			body: JSON.stringify({asset_id: assetId})
		})
	}
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
	<ItemViewer id={currentId} bind:show={showViewer} menu={viewMenu}></ItemViewer>
</div>

{#snippet viewMenu()}
	<a href={`/api/asset/${currentId}/original/`} target="_blank" class="btn btn-soft">
		<Icon path={mdiDownload} /> Download.
	</a>

	<button class="btn btn-soft"
		onclick={()=>{setAlbumCover(params.id, currentId)}}
	>
		<Icon path={mdiImageAlbum} /> Set as album cover.
	</button>
{/snippet}
