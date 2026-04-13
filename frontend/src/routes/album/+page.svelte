<script lang="ts">
	import AlbumItem from '$lib/components/AlbumItem.svelte';
	import { mdiImageSearch } from '@mdi/js';
	import type { PageProps } from './$types';
	import Icon from 'mdi-svelte';
	import { goto } from '$app/navigation';

	let { data }: PageProps = $props();

	let searchInput: HTMLInputElement;

	function doSearch() {
		const search = searchInput.value;
		const url = new URL('/search', location.origin);
		url.searchParams.append('search', search);

		goto(url);
	}
</script>

<div class="relative flex h-screen w-screen flex-col">
	<div class="navbar bg-base-100 shadow-sm">
		<div class="flex-1">
			<div class="btn text-xl btn-ghost">Albums</div>
		</div>
		<div class="flex-none">
			<ul class="menu menu-horizontal px-1">
				<li><a href="/album">Albums</a></li>
			</ul>
			<div class="join">
				<div>
					<label class="input join-item">
						<Icon path={mdiImageSearch} />
						<input type="text" placeholder="search" bind:this={searchInput} />
					</label>
				</div>
				<button class="btn join-item btn-neutral" onclick={() => doSearch()}> Search </button>
			</div>
		</div>
	</div>

	<div class="overflow-auto bg-base-300">
		<div class="grid grid-cols-1 flex-wrap gap-2 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
			{#each data.albums as album (album.id)}
				<AlbumItem {album} />
			{/each}
		</div>
	</div>
</div>
