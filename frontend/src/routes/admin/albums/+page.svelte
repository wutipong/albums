<script lang="ts">
	import type { PageProps } from './$types';
	import Icon from 'mdi-svelte';
	import { mdiInformationOutline } from '@mdi/js';
	import { authClient } from '$lib/auth-client';


	let { data }: PageProps = $props();
	const formatter = new Intl.NumberFormat('en-US', {
		notation: 'compact',
		compactDisplay: 'short',
		maximumFractionDigits: 1 // Controls decimal points (e.g., 1.5K vs 1.52K)
	});

	$inspect(data);
</script>

<svelte:head>
	<title>Albums: Administration</title>
</svelte:head>

<div class="w-full p-4 pt-8">
	<h2 class="p-4 text-2xl font-semibold">Statistics</h2>
	<div class="stats shadow">
		<div class="stat">
			<div class="stat-title">Total Assets</div>
			<div class="stat-value text-primary">{formatter.format(data.count)}</div>
		</div>

		<div class="stat">
			<div class="stat-title">Pending Assets</div>
			<div class="stat-value text-primary">{formatter.format(data.pendingCount)}</div>
		</div>

		<div class="stat">
			<div class="stat-title">Failed Assets</div>
			<div class="stat-value text-error">{formatter.format(data.failedCount)}</div>
		</div>

		<div class="stat">
			<div class="stat-title">Total Albums</div>
			<div class="stat-value text-primary">{formatter.format(data.albumCount)}</div>
		</div>
	</div>
</div>
