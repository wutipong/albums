<script lang="ts">
	import NavBar from '$lib/components/NavBar.svelte';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();
	const formatter = new Intl.NumberFormat('en-US', {
		notation: 'compact',
		compactDisplay: 'short',
		maximumFractionDigits: 1 // Controls decimal points (e.g., 1.5K vs 1.52K)
	});
</script>

<div class="relative flex h-screen w-screen flex-col">
	<NavBar />
	<div class="overflow-auto p-4 pt-8">
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
				<div class="stat-title">Total Albums</div>
				<div class="stat-value text-primary">{formatter.format(data.albumCount)}</div>
			</div>
		</div>

		<div class="overflow-x-auto">
			<table class="table">
				<thead>
					<tr>
						<th>#</th>
						<th>ID</th>
						<th>Albums</th>
						<th>Created At</th>
					</tr>
				</thead>
				<tbody>
					{#if data.pendingCount === 0n}
						<tr>
							<td colspan="4" class="text-center">No pending assets found.</td>
						</tr>
					{:else}
						{#each data.pendings as asset, index}
							<tr>
								<th>{index + 1}</th>
								<td>{asset.asset_id}</td>
								<td>{asset.album_name}</td>
								<td>{asset.created_at.toLocaleString()}</td>
							</tr>
						{/each}

						<tr>
							<td colspan="4" class="text-center">Some pending assets are not displayed.</td>
						</tr>
					{/if}
				</tbody>
				<tfoot>
					<tr>
						<th>#</th>
						<th>ID</th>
						<th>Albums</th>
						<th>Created At</th>
					</tr>
				</tfoot>
			</table>
		</div>
	</div>
</div>
