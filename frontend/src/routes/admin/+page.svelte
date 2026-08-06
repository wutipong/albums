<script lang="ts">
	import NavBar from '$lib/components/NavBar.svelte';
	import type { PageProps } from './$types';
    import Icon from 'mdi-svelte';
    import {mdiInformationOutline} from '@mdi/js';

	let { data }: PageProps = $props();
	const formatter = new Intl.NumberFormat('en-US', {
		notation: 'compact',
		compactDisplay: 'short',
		maximumFractionDigits: 1 // Controls decimal points (e.g., 1.5K vs 1.52K)
	});

	$inspect(data);
</script>

<div class="relative flex h-screen w-screen flex-col">
	<NavBar />
	<div class="overflow-auto p-4 pt-8">
        <h2 class="text-2xl font-semibold">Statistics</h2>
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

		<h2 class="text-2xl font-semibold">Processing Queue</h2>
		<div role="alert" class="alert alert-info">
			<Icon path={mdiInformationOutline} />
			<span>This table shows up to 20 items being wait in the queue.</span>
		</div>

		<div class="overflow-x-auto">
			<table class="table">
				<thead>
					<tr>
						<th>ID</th>
						<th>Asset ID</th>
						<th>Type</th>
						<th>Album</th>
						<th>Created At</th>
					</tr>
				</thead>
				<tbody>
					{#if data.queueItems.length === 0}
						<tr>
							<td colspan="5" class="text-center">No processing assets found.</td>
						</tr>
					{:else}
						{#each data.queueItems as asset, index}
							<tr>
								<th>{asset.id}</th>
								<td>{asset.asset_id}</td>
								<td>{asset.type}</td>
								<td
									><a href={`album/${asset.album_id}`} class="link link-primary"
										>{asset.album_name}</a
									></td
								>
								<td>{asset.created_at.toLocaleString()}</td>
							</tr>
						{/each}
					{/if}
				</tbody>
				<tfoot>
					<tr>
						<th>ID</th>
						<th>Asset ID</th>
						<th>Type</th>
						<th>Album</th>
						<th>Created At</th>
					</tr>
				</tfoot>
			</table>
		</div>
	</div>
</div>
