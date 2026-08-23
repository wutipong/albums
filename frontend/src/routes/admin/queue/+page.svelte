<script lang="ts">
	import { mdiInformationOutline } from '@mdi/js';
	import type { PageProps } from './$types';
	import Icon from 'mdi-svelte';

	let { data }: PageProps = $props();
	let progress = $derived(Number((100n * data.queueProcessed) / data.queueTotal));
	const formatter = new Intl.NumberFormat('en-US', {
		notation: 'compact',
		compactDisplay: 'short',
		maximumFractionDigits: 1 // Controls decimal points (e.g., 1.5K vs 1.52K)
	});
</script>

<div class="w-full p-4">
	<h2 class="p-4 text-2xl font-semibold">Processing Queue</h2>
	<div class="m-4">
		<progress class="progress w-full" value={progress} max={100}></progress>
		{data.queueProcessed}/{data.queueTotal}
	</div>

	<div class="stats shadow">
		<div class="stat">
			<div class="stat-title">#Jobs</div>
			<div class="stat-value text-primary">{formatter.format(data.queueTotal)}</div>
		</div>

		<div class="stat">
			<div class="stat-title">Pending Jobs</div>
			<div class="stat-value text-primary">{formatter.format(data.queuePending)}</div>
		</div>

		<div class="stat">
			<div class="stat-title">Failed Jobs</div>
			<div class="stat-value text-error">{formatter.format(data.queueFailed)}</div>
		</div>
	</div>

	<hr />
	<h3 class="p-4 text-xl font-semibold">Queue Details</h3>

	<div role="alert" class="alert alert-info">
		<Icon path={mdiInformationOutline} />
		<span>This table shows up to 20 items waiting in the queue.</span>
	</div>

	<div class="overflow-x-auto">
		<table class="table">
			<thead>
				<tr>
					<th>ID</th>
					<th>Command</th>
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
							<th>{asset.command}</th>
							<td>{asset.asset_id}</td>
							<td>{asset.type}</td>
							<td>
								<a href={`/album/${asset.album_id}`} class="link link-primary">
									{asset.album_name}
								</a>
							</td>
							<td>{asset.created_at.toLocaleString()}</td>
						</tr>
					{/each}
				{/if}
			</tbody>
			<tfoot>
				<tr>
					<th>ID</th>
					<th>Command</th>
					<th>Asset ID</th>
					<th>Type</th>
					<th>Album</th>
					<th>Created At</th>
				</tr>
			</tfoot>
		</table>
	</div>
</div>
