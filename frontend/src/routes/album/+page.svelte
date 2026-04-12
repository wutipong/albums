<script lang="ts">
	import AlbumItem from '$lib/components/AlbumItem.svelte';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();

	async function notifyScanCache() {
		await fetch('/api/asset/scan');
	}
</script>

<div class="relative flex flex-col h-screen w-screen">
	<div class="navbar bg-base-100 shadow-sm">
		<div class="flex-1">
			<div class="btn text-xl btn-ghost">Albums</div>
		</div>
		<div class="flex-none">
			<ul class="menu menu-horizontal px-1">
				<li><a href="/album">Albums</a></li>
			</ul>
		</div>
	</div>

	<div class="overflow-auto">
		<div class="grid grid-cols-1 gap-2 bg-base-300 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 flex-wrap">
			{#each data.albums as album (album.id)}
				<AlbumItem {album} />
			{/each}
		</div>
	</div>
</div>
