<script lang="ts">
	import notAvailableSvg from '$lib/assets/not-available-small.svg?raw';

	interface Props {
		album: any;
		aspect: 'portrait' | 'landscape';
	}

	let { album, aspect = 'landscape' }: Props = $props();
	let imgHeight = $derived(aspect == 'landscape' ? 200 : 450);

	let coverLoading = $state(true);
</script>

<div
	class="card m-4 w-[300px] border-1 border-base-300 shadow hover:bg-base-100 hover:shadow-xl hover:my-1"
	class:h-[450px]={aspect === 'portrait'}
	class:h-[300px]={aspect === 'landscape'}
	id={album.id}
>
	<figure
		class="block"
		class:h-[450px]={aspect === 'portrait'}
		class:h-[200px]={aspect === 'landscape'}
	>
		<a href={`/album/${album.id}/`} class="block h-full w-full">
			<div
				class="relative"
				class:h-[450px]={aspect === 'portrait'}
				class:h-[200px]={aspect === 'landscape'}
			>
				{#if album.cover_url == ''}
					{@html notAvailableSvg}
				{:else}
					<img
						class="w-[300px] object-cover"
						src={album.cover_url}
						alt={album.id}
						class:h-[450px]={aspect === 'portrait'}
						class:h-[200px]={aspect === 'landscape'}
						onload={() => (coverLoading = false)}
						loading="lazy"
					/>
				{/if}
				<div
					class="absolute bottom-0 left-0 h-[100px] w-full p-4 bg-gradient-to-b from-base-100/0 via-base-100/70 via-30% to-base-100/100"
					class:hidden={aspect != 'portrait'}
				>
					{album.name}
				</div>
			</div>
		</a>
	</figure>
	<div class="card-body h-4 overflow-clip" class:hidden={aspect != 'landscape'}>
		<a href={`/album/${album.id}/`}>
			{album.name}
		</a>
	</div>
</div>
