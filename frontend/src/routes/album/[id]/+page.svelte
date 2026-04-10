<script lang="ts">
	import { invalidateAll } from '$app/navigation';
	import AssetThumbnail from '$lib/components/AssetThumbnail.svelte';
	import type { PageProps } from './$types';

	let { params, data }: PageProps = $props();
	let files: FileList | null = $state(null);

	$inspect(data)

	function uploadFile() {
		if (files) {
			const file = files[0];

			const formData = new FormData();
			formData.append('file', file);
			formData.append('albumId', params.id);

			fetch(`/api/asset`, {
				method: 'POST',
				body: formData
			}).then((res) => {
				if (res.ok) {
					invalidateAll();
				} else {
					alert('Failed to upload files');
				}
			});
		}
	}
</script>

<div class="text-xl font-bold">Album {data.name}</div>

<div class="form-group mb-4">
	<label for="file" class="mb-2 block font-bold">File</label>
	<input id="file" type="file" class="file-input" accept="image/*" bind:files />
	<button class="mt-2 rounded bg-green-500 px-4 py-2 text-white" onclick={uploadFile}>
		Upload
	</button>
</div>

<div class="flex flex-wrap">
	{#each data.assets as asset (asset)}
		<div class="m-1 block h-[200px] grow object-cover">
			<AssetThumbnail id={asset}></AssetThumbnail>
		</div>
	{/each}
</div>
