<script lang="ts">
	import AlbumItem from '$lib/components/AlbumItem.svelte';
	import type { PageProps } from './$types';
	import NavBar from '$lib/components/NavBar.svelte';
	import { mdiClose, mdiFilter, mdiImageAlbum } from '@mdi/js';
	import Icon from 'mdi-svelte';
	import { authClient } from '$lib/auth-client';

	let { data }: PageProps = $props();

	const session = authClient.useSession();

	let filter = $state('');
</script>

{#snippet title()}
	<div class="flex text-xl md:ms-4">
		<Icon path={mdiImageAlbum}></Icon>
		Albums
	</div>
{/snippet}

<div class="relative flex h-screen w-screen flex-col">
	<NavBar {title} />

	<div class="bg-base-300 shadow w-full p-2 justify-end gap-2 border-1 border-base-300 flex">
			<label class="input w-full max-w-xs">
				<span class="label">
					<Icon path={mdiFilter} />
				</span>
				<input
					type="text"
					placeholder="Type here to filter albums"
					bind:value={filter}
				/>
			</label>
		</div>

	<div class="mx-4 overflow-auto pt-4 pb-20">
		<div class="flex flex-wrap justify-evenly gap-2">
			{#each data.albums.filter((album) => album.name
					.toLowerCase()
					.includes(filter)) as album (album.id)}
				<AlbumItem {album} />
			{/each}
		</div>
	</div>
</div>
