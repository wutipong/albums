<script lang="ts">
	import AlbumItem from '$lib/components/AlbumItem.svelte';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();

	async function notifyScanCache() {
		await fetch('/api/asset/scan');
	}
</script>

<div>Album List</div>

<button onclick={notifyScanCache} class="mb-4 rounded bg-blue-500 px-4 py-2 text-white">
	Scan Cache
</button>

<div class="grid grid-cols-5 xs:grid-cols-1 gap-2 bg-base-200">
  	{#each data.albums as album (album.id)}
		<AlbumItem {album}/>
	{/each}
</div>