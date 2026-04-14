<script lang="ts">
	import AlbumItem from '$lib/components/AlbumItem.svelte';
	import type { PageProps } from './$types';
	import NavBar from '$lib/components/NavBar.svelte';

	let { data }: PageProps = $props();
</script>

{#snippet title()}
	<div class="text-xl">Albums </div>
{/snippet}

<div class="relative flex h-screen w-screen flex-col bg-base-300">
	<NavBar title={title}/>

	<div class="overflow-auto mt-8">
		<div class="grid grid-cols-1 flex-wrap gap-2 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
			{#each data.albums as album (album.id)}
				<AlbumItem {album} />
			{/each}
		</div>
	</div>
</div>
