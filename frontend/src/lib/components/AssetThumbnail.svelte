<script lang="ts">
	import notAvailableSvg from '$lib/assets/not-available-small.svg?raw';
	let { asset, onclick = (id: string) => {} } = $props();

	const TARGET_HEIGHT = 200;

	let preview = $state(false);
	let thumbnailWidth = $derived((TARGET_HEIGHT * asset.thumbnail_width) / asset.thumbnail_height);
	let thumbnailHeight = TARGET_HEIGHT;
	let assetType = $derived(asset.type);
	let imageFrames = $derived(asset.image_frames);

	let thumbnailLoading = $state(true);
	let previewLoading = $state(true);

	let rootElement: HTMLElement | undefined = $state();

	export function scrollIntoView(arg?: boolean | ScrollIntoViewOptions) {
		rootElement?.scrollIntoView(arg);
	}
</script>

<button
	bind:this={rootElement}
	tabindex="0"
	class={`block cursor-pointer overflow-hidden rounded-xl p-4 
		hover:bg-base-content/20 hover:shadow-xl`}
	onmouseenter={() => (preview = true)}
	onmouseleave={() => (preview = false)}
	onclick={() => {
		onclick(asset);
	}}
	disabled={asset.process_status != 'processed'}
>
	<div class="relative h-full w-full">
		<div
			class:hidden={preview}
			class="box-border h-full w-full overflow-hidden rounded-xl border-1 border-base-content/20"
			style={`width: ${thumbnailWidth}px; height: ${thumbnailHeight}px;`}
		>
			{#if asset.thumbnail_url === ''}
				{@html notAvailableSvg}
			{:else}
				<img
					width={thumbnailWidth}
					height={thumbnailHeight}
					src={asset.thumbnail_url}
					alt={asset.id}
					class="h-full w-full"
					onload={() => (thumbnailLoading = false)}
					loading="lazy"
				/>
				<div class="h-full w-full skeleton bg-base-100" class:hidden={!thumbnailLoading}></div>
			{/if}
		</div>

		<div
			class:hidden={!preview}
			class="box-border h-full w-full overflow-hidden rounded-xl border-1 border-base-content/20"
			style={`width: ${thumbnailWidth}px; height: ${thumbnailHeight}px;`}
		>
			{#if asset.preview_url == ''}
				{@html notAvailableSvg}
			{:else}
				<img
					width={thumbnailWidth}
					height={thumbnailHeight}
					src={asset.preview_url}
					alt={asset.id}
					class="h-full w-full"
					onload={() => (previewLoading = false)}
					loading="lazy"
				/>
				<div
					class="h-full w-full skeleton animate-pulse bg-base-200"
					class:hidden={!previewLoading}
				></div>
			{/if}
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
