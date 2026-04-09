<script lang="ts">
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

<ul>
	{#each data.albums as album (album.id)}
		<li>
			<a href={`/album/${album.id}`} class="mb-2 block rounded border p-2">
				{album.name}
			</a>
		</li>
	{/each}
</ul>
