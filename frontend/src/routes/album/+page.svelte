<script lang="ts">
	import AlbumItem from '$lib/components/AlbumItem.svelte';
	import type { PageProps } from './$types';
	import NavBar from '$lib/components/NavBar.svelte';
	import { mdiAlbum, mdiImageAlbum } from '@mdi/js';
	import Icon from 'mdi-svelte';
	import { authClient } from '$lib/auth-client';
	import { onMount } from 'svelte';

	let { data }: PageProps = $props();

	const session = authClient.useSession();
</script>

<div>
	{#if $session.data}
		<div>
			<p>
				{$session.data.user.name}
			</p>
		</div>
		{:else}
		<div>session not found</div>
		<div> {$session.error}</div>
		<div> {$session.isPending}</div>
		<div>{$session.isRefetching}</div>
	{/if}
</div>

{#snippet title()}
	<div class="flex text-xl md:ms-4">
		<Icon path={mdiImageAlbum}></Icon>
		Albums
	</div>
{/snippet}

<div class="relative flex h-screen w-screen flex-col bg-base-100">
	<NavBar {title} />

	<div class="m-4 overflow-auto pt-8">
		<div class="flex flex-wrap justify-evenly gap-4">
			{#each data.albums as album (album.id)}
				<AlbumItem {album} />
			{/each}
		</div>
	</div>
</div>
