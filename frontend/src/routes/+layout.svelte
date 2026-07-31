<script lang="ts">
	import './layout.css';
	import favicon from '$lib/assets/favicon.svg';
	// This injects Temporal onto the global window/globalThis object
	import 'temporal-polyfill/global';

	let { children } = $props();


</script>

<svelte:head><link rel="icon" href={favicon} /></svelte:head>
<div class="vh-full vw-full bg-base-200">
	{@render children()}
</div>
