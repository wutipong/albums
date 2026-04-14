<script lang="ts">
	import { goto } from '$app/navigation';
	import { mdiImageSearch, mdiMagnify, mdiImageAlbum } from '@mdi/js';
	import Icon from 'mdi-svelte';

	let searchInput: HTMLInputElement;

	let { title } = $props();

	function doSearch() {
		const search = searchInput.value;
		const url = new URL('/search', location.origin);
		url.searchParams.append('search', search);

		goto(url);
	}
</script>

<div class="navbar bg-base-100 shadow-sm">
	<div class="flex-1">
		{@render title()}
	</div>
	<div class="flex-none">
		<ul class="menu menu-horizontal px-1">
			<li><a href="/album"><Icon path={mdiImageAlbum} />Albums</a></li>
		</ul>
    </div>
    <div class="flex-none">
		<div class="join">
			<label class="input join-item">
				<input type="text" placeholder="search" bind:this={searchInput} />
			</label>
			<button class="btn join-item btn-neutral" onclick={() => doSearch()}>
				<Icon path={mdiImageSearch} />
			</button>
		</div>
	</div>
</div>
