<script lang="ts">
	import type { PageProps } from './$types';

	let { params }: PageProps = $props();
</script>

<div>Album {params.id}</div>
