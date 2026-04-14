<script lang="ts">
	import { mdiChevronLeft, mdiChevronRight, mdiClose, mdiDotsVertical } from '@mdi/js';
	import Icon from 'mdi-svelte';
	import 'vidstack/bundle';

	let {
		asset = $bindable({ id: '<placeholder>' }),
		show = $bindable(false),
		next,
		previous,
		hasNext = false,
		hasPrevious = false,
		menu
	} = $props();
</script>

{#if show}
	<div
		role="presentation"
		class="absolute top-0 right-0 bottom-0 left-0 backdrop-blur-lg backdrop-brightness-50"
	>
		{#if asset.type === 'image'}
			<div class="h-full w-full">
				<img
					src={`/api/asset/${asset.id}/view`}
					alt={asset.id}
					class="m-auto h-full w-full object-contain"
				/>
			</div>
		{/if}
		{#if asset.type === 'video'}
			<media-player
				class="h-full w-full"
				title={asset.original}
				src={`/api/asset/${asset.id}/view`}
			>
				<media-provider></media-provider>
				<media-video-layout></media-video-layout>
			</media-player>
		{/if}
		<div
			class="absolute top-1/2 left-4 -translate-y-1/2 bg-transparent"
			data-theme={asset.type === 'video' ? 'dark' : null}
		>
			<button
				class="btn btn-circle btn-ghost btn-lg"
				class:btn-disabled={!hasPrevious}
				onclick={() => {
					previous();
				}}
			>
				<Icon path={mdiChevronLeft} />
			</button>
		</div>
		<div
			class="absolute top-1/2 right-4 -translate-y-1/2 bg-transparent"
			data-theme={asset.type === 'video' ? 'dark' : null}
		>
			<button
				class="btn btn-circle btn-ghost btn-lg"
				class:btn-disabled={!hasNext}
				onclick={() => {
					next();
				}}
			>
				<Icon path={mdiChevronRight} />
			</button>
		</div>
		<div
			class="absolute top-4 right-4 flex flex-row-reverse gap-4 rounded-full bg-transparent"
			data-theme={asset.type === 'video' ? 'dark' : null}
		>
			<button class="btn btn-circle btn-ghost btn-lg" onclick={() => (show = false)}>
				<Icon path={mdiClose} />
			</button>
			{#if menu}
				<button
					class="btn btn-circle btn-ghost btn-lg"
					popovertarget="popover-1"
					style="anchor-name:--anchor-1"
				>
					<Icon path={mdiDotsVertical} />
				</button>
				<ul
					class="menu dropdown w-52 rounded-box bg-base-300 shadow-sm"
					popover
					id="popover-1"
					style="position-anchor:--anchor-1"
				>
					{@render menu()}
				</ul>
			{/if}
		</div>
	</div>
{/if}
