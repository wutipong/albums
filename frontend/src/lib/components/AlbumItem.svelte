<script lang="ts">
	import notAvailableSvg from '$lib/assets/not-available-small.svg?raw';

	interface Props {
		album: any;
		aspect: 'portrait' | 'landscape';
	}

	let { album, aspect = 'landscape' }: Props = $props();
</script>

<div class="block rounded-xl p-4 hover:bg-base-100 hover:shadow-xl" id={album.id}>
	<div class="w-[300px] overflow-clip rounded-xl border-1 border-base-300">
		<figure
			class="block w-[300px]"
			class:h-[450px]={aspect === 'portrait'}
			class:h-[200px]={aspect === 'landscape'}
		>
			<a href={`/album/${album.id}/`} class="block h-full w-full">
				<div
					class="relative w-[300px]"
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
							loading="lazy"
						/>
					{/if}
					<div
						class={`absolute bottom-0 left-0 h-[100px] w-full 
                            bg-gradient-to-b from-base-100/0 via-base-100/70 via-30% to-base-100/100 
                            p-4 break-all`}
						class:hidden={aspect != 'portrait'}
					>
						{album.name}
					</div>
				</div>
			</a>
		</figure>
		<div class="h-[100px] w-full overflow-clip p-4 break-all" class:hidden={aspect != 'landscape'}>
			<a href={`/album/${album.id}/`}>
				{album.name}
			</a>
		</div>
	</div>
</div>
