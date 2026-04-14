<script lang="ts">
	import AlbumItem from '$lib/components/AlbumItem.svelte';
	import type { PageProps } from './$types';
	import NavBar from '$lib/components/NavBar.svelte';
	import { mdiAlbum, mdiImageAlbum } from '@mdi/js';
	import Icon from 'mdi-svelte';

	let { data }: PageProps = $props();
</script>

{#snippet title()}
	<div class="flex text-xl md:ms-4">
		<Icon path={mdiImageAlbum}></Icon>
		Albums
	</div>
{/snippet}

<div class="relative flex h-screen w-screen flex-col bg-base-300">
	<NavBar {title}/>

	<div class="overflow-auto pt-8 m-4">
		<div class="flex flex-wrap gap-4 justify-evenly">
			{#each data.albums as album (album.id)}
				<AlbumItem {album} />
			{/each}
		</div>
	</div>
</div>
