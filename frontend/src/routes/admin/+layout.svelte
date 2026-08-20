<script lang="ts">
	import NavBar from '$lib/components/NavBar.svelte';
	import { mdiImageAlbum, mdiTrayFull } from '@mdi/js';
	import Icon from 'mdi-svelte';

	let { children } = $props();
</script>

<NavBar />
<div class="drawer lg:drawer-open">
	<input id="my-drawer-3" type="checkbox" class="drawer-toggle" />
	<div class="drawer-content">
		{@render children()}
	</div>
	<div class="drawer-side">
		<label for="my-drawer-3" aria-label="close sidebar" class="drawer-overlay"></label>
		<ul class="menu min-h-full w-80 bg-base-200 p-4">
			<li><a href="/admin/albums"><Icon path={mdiImageAlbum} />Albums</a></li>
			<li><a href="/admin/queue"><Icon path={mdiTrayFull} /> Processing Queue</a></li>
		</ul>
	</div>
</div>
