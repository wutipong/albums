<script lang="ts">
	import { mdiChevronLeft, mdiChevronRight, mdiClose, mdiDotsVertical } from '@mdi/js';
	import Hammer from 'hammerjs';
	import Icon from 'mdi-svelte';
	import 'vidstack/bundle';
	import { MediaPlayerElement } from 'vidstack/elements';

	let {
		viewURL = 'http://example.com',
		assetType = 'image',
		filename = '',
		show = $bindable(false),
		next,
		previous,
		hasNext = false,
		hasPrevious = false,
		menu
	} = $props();

	let loading = $state(true);
	let mediaPlayer: MediaPlayerElement | null = $state(null);

	$effect(() => {
		loading = true;
	});

	function doNext() {
		if (mediaPlayer) {
			mediaPlayer.pause();
		}
		next();
	}

	function doPrevious() {
		if (mediaPlayer) {
			mediaPlayer.pause();
		}
		previous();
	}

	function doClose() {
		if (mediaPlayer) {
			mediaPlayer.pause();
		}
		show = false;
	}

	function handleKeyDown(event: KeyboardEvent) {
		if (!show) return;
		if (event.key === 'Escape') {
			doClose();
		} else if (event.key === 'ArrowRight' && hasNext) {
			if (hasNext) {
				doNext();
			}
		} else if (event.key === 'ArrowLeft' && hasPrevious) {
			if (hasPrevious) {
				doPrevious();
			}
		}
	}

	function hammerJsAttachment(element: HTMLElement) {
		let manager = new Hammer.Manager(element);
		let swipe = new Hammer.Swipe();
		manager.add(swipe);
		manager.on('swipeleft', () => {
			if (hasNext) {
				doNext();
			}
		});

		manager.on('swiperight', () => {
			if (hasPrevious) {
				doPrevious();
			}
		});
	}
</script>

<svelte:window onkeydown={handleKeyDown} />

<div
	role="presentation"
	class="absolute top-0 right-0 bottom-0 left-0 backdrop-blur-lg backdrop-brightness-50"
	class:hidden={!show}
	{@attach hammerJsAttachment}
>
	{#if assetType === 'image'}
		<div class="h-full w-full">
			<img
				class:hidden={loading}
				onload={() => (loading = false)}
				onerror={() => (loading = false)}
				src={viewURL}
				alt={filename}
				class="m-auto h-full w-full object-contain"
			/>
		</div>
	{/if}
	{#if assetType === 'video'}
		{#key viewURL}
			<media-player title={filename} src={viewURL} bind:this={mediaPlayer} autoplay controls>
				<media-provider>
					<source src={viewURL} type="video/mp4" />
				</media-provider>
				<media-video-layout></media-video-layout>
			</media-player>
		{/key}
	{/if}
	<div class="absolute top-1/2 left-4 -translate-y-1/2 bg-transparent">
		<button
			class="btn btn-circle btn-lg btn-soft"
			class:btn-disabled={!hasPrevious}
			onclick={() => {
				doPrevious();
			}}
		>
			<Icon path={mdiChevronLeft} />
		</button>
	</div>
	<div class="absolute top-1/2 right-4 -translate-y-1/2 bg-transparent">
		<button
			class="btn btn-circle btn-lg btn-soft"
			class:btn-disabled={!hasNext}
			onclick={() => {
				doNext();
			}}
		>
			<Icon path={mdiChevronRight} />
		</button>
	</div>
	<div
		class="absolute top-4 right-4 flex flex-row-reverse gap-4 rounded-full bg-transparent"
		data-theme={assetType === 'video' ? 'dark' : null}
	>
		<button
			class="btn btn-circle btn-lg btn-soft"
			onclick={() => {
				doClose();
			}}
		>
			<Icon path={mdiClose} />
		</button>
		{#if menu}
			<button
				class="btn btn-circle btn-lg btn-soft"
				popovertarget="popover-1"
				style="anchor-name:--anchor-1"
			>
				<Icon path={mdiDotsVertical} />
			</button>
			<ul
				class="menu dropdown w-52 rounded-box bg-base-100 shadow-sm"
				popover
				id="popover-1"
				style="position-anchor:--anchor-1"
			>
				{@render menu()}
			</ul>
		{/if}
	</div>
</div>

<style lang="scss">
	media-player {
		width: 100%;
		height: 100%;
		aspect-ratio: unset;
	}

	:global(media-provider video) {
		width: 95%;
		height: 95%;
	}
</style>
