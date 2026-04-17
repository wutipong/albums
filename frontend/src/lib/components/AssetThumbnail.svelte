<script lang="ts">
	let { asset, onclick = (id: string) => {} } = $props();

	const TARGET_HEIGHT = 200;

	let preview = $state(false);
	let thumbnailWidth = $derived((TARGET_HEIGHT * asset.thumbnail_width) / asset.thumbnail_height);
	let thumbnailHeight = TARGET_HEIGHT;
	let assetType = $derived(asset.type);
	let imageFrames = $derived(asset.image_frames);
</script>

<button
	tabindex="0"
	class={`block h-[${thumbnailHeight}px] m-1 overflow-hidden rounded-xl cursor-pointer hover:shadow-xl`}
	style={`width: ${thumbnailWidth}px;`}
	onmouseenter={() => (preview = true)}
	onmouseleave={() => (preview = false)}
	onclick={() => {
		onclick(asset);
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
				src={asset.thumbnail_url}
				alt={asset.id}
				class:hidden={preview}
			/>
		</div>

		<div
			class:hidden={!preview}
			class="box-border h-full w-full overflow-hidden rounded-xl"
			style={`width: ${thumbnailWidth}px; height: ${thumbnailHeight}px;`}
		>
			<img
				width={thumbnailWidth}
				height={thumbnailHeight}
				src={asset.preview_url}
				alt={asset.id}
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
</button>
