<script lang="ts">
	import { invalidateAll } from '$app/navigation';
	import type { PageProps } from './$types';

	let { params, data }: PageProps = $props();

	async function createAlbum() {
		const name = 'New Album ' + new Date().toLocaleTimeString();
		if (name) {
			const res = await fetch('/api/album', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ name })
			});
			if (res.ok) {
				invalidateAll();
			} else {
				alert('Failed to create album');
			}
		}
	}
</script>

<div>Album List</div>

<button onclick={createAlbum} class="mb-4 rounded bg-blue-500 px-4 py-2 text-white"
	>Create New Album</button
>

<ul>
	{#each data.albums as album (album.id)}
		<li>
			<a href={`/album/${album.id}`} class="mb-2 block rounded border p-2">
				{album.name}
			</a>
		</li>
	{/each}
</ul>
