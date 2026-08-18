<script lang="ts">
	import AlbumItem from '$lib/components/AlbumItem.svelte';
	import type { PageProps } from './$types';
	import NavBar from '$lib/components/NavBar.svelte';
	import { mdiFilter, mdiOrderBoolAscending } from '@mdi/js';
	import Icon from 'mdi-svelte';
	import { authClient } from '$lib/auth-client';

	let { data }: PageProps = $props();

	const session = authClient.useSession();

	let filter = $state('');
	let order = $state('name asc');

	$inspect(data);

	let albums = $derived(
		data.albums
			.filter((album) => album.name.toLowerCase().includes(filter.toLowerCase()))
			.sort((a, b) => {
				switch (order) {
					case 'name asc':
						return a.name.localeCompare(b.name);
					case 'name desc':
						return b.name.localeCompare(a.name);

					case 'created_at asc':
						return new Date(a.created_at).getTime() - new Date(b.created_at).getTime();
					case 'created_at desc':
						return new Date(b.created_at).getTime() - new Date(a.created_at).getTime();

					default:
						return 0;
				}
			})
	);
</script>

<svelte:head>
	<title>Albums</title>
</svelte:head>

<div class="relative flex h-screen w-screen flex-col">
	<NavBar />

	<div class="flex w-full justify-end gap-2 border-1 border-base-300 bg-base-300 p-2 shadow">
		<label class="select w-full max-w-xs">
			<span class="label">
				<Icon path={mdiOrderBoolAscending} />
			</span>
			<select class="select" bind:value={order}>
				<option value="name asc">Name Ascending</option>
				<option value="name desc">Name Descending</option>
				<option value="created_at asc">Creation Time Ascending</option>
				<option value="created_at desc">Creation Time Ascending</option>
			</select>
		</label>

		<label class="input w-full max-w-xs">
			<span class="label">
				<Icon path={mdiFilter} />
			</span>
			<input type="text" placeholder="Type here to filter albums" bind:value={filter} />
		</label>
	</div>

	<div class="mx-4 overflow-auto pt-4 pb-20">
		<div class="flex flex-wrap justify-evenly gap-2">
			{#each albums as album (album.id)}
				<AlbumItem {album} aspect={data.aspect}/>
			{/each}
		</div>
	</div>
</div>
