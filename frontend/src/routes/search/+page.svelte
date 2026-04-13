<script lang="ts">
	import AssetThumbnail from '$lib/components/AssetThumbnail.svelte';
	import ItemViewer from '$lib/components/ItemViewer.svelte';
	import Icon from 'mdi-svelte';
	import type { PageProps } from './$types';
	import { mdiDownload } from '@mdi/js';

	let { data, params }: PageProps = $props();
	let asset = $state({id:'<placeholder>'});
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

<div class="relative flex h-screen w-screen flex-col bg-base-300">
	<div class="navbar bg-base-100 shadow-sm">
		<div class="flex-1">
			<div class="btn text-xl btn-ghost">{data.search}</div>
		</div>
		<div class="flex-none">
			<ul class="menu menu-horizontal px-1">
				<li><a href="/album">Albums</a></li>
			</ul>
		</div>
	</div>
	<div class="overflow-auto">
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
		bind:asset={asset}
		bind:show={showViewer}
		{next}
		{previous}
		{hasNext}
		{hasPrevious}
		menu={viewMenu}
	/>
</div>

{#snippet viewMenu()}
	<a href={`/api/asset/${asset.id}/original/`} target="_blank" class="btn btn-soft">
		<Icon path={mdiDownload} /> Download.
	</a>
{/snippet}
