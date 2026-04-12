<script lang="ts">
	import { onMount } from 'svelte';
	let { id, onclick = (id: string) => {} } = $props();

	let thumbnailWidth = $state(0);
	let thumbnailHeight = $state(0);
	let available = $state(false);
	let videoDuration = $state('');
	let imageFrames = $state(1);
	let assetType = $state('image');

	let preview = $state(false);

	onMount(async () => {
		const resp = await fetch(`/api/asset/${id}/`);
		const obj = await resp.json();

		const TARGET_HEIGHT = 200;
		const ratio = TARGET_HEIGHT / obj.thumbnail_height;

		thumbnailWidth = obj.thumbnail_width * ratio;
		thumbnailHeight = TARGET_HEIGHT;
		available = obj.available;

		assetType = obj.type;
		videoDuration = obj.video_duration.toString();
		imageFrames = obj.image_frames;
	});
</script>

<a
	role="button"
	tabindex="0"
    href="#"
	class={`block h-[${thumbnailHeight}px] m-1 overflow-hidden rounded-xl`}
	style={`width: ${thumbnailWidth}px;`}
	onmouseenter={() => (preview = true)}
	onmouseleave={() => (preview = false)}
	onclick={() => {
		onclick(id);
	}}
>
	<div class="relative h-full w-full">
		<div
			class:hidden={preview}
			class="box-border h-full w-full overflow-hidden rounded-xl"
			style={`width: ${thumbnailWidth}px; height: ${thumbnailHeight}px;`}
		>
			<img
				width={thumbnailWidth}
				height={thumbnailHeight}
				src={`/api/asset/${id}/thumbnail`}
				alt={id}
				class:hidden={preview}
			/>
		</div>

		<div
			class:hidden={!preview}
			class="box-border h-full w-full overflow-hidden rounded-xl border-4"
			style={`width: ${thumbnailWidth}px; height: ${thumbnailHeight}px;`}
		>
			<img
				width={thumbnailWidth}
				height={thumbnailHeight}
				src={`/api/asset/${id}/preview`}
				alt={id}
				class="h-full w-full"
			/>
		</div>

		<div class="absolute top-1 right-2 place-items-end">
			{#if assetType === 'video'}
				<div class="badge">Video</div>
			{/if}

			{#if imageFrames > 1}
				<div class="badge">Animation</div>
			{/if}
		</div>
	</div>
</a>
