<script lang="ts">
	import { mdiImage, mdiPlayCircle, mdiVideo } from '@mdi/js';
	import type { PageProps } from './$types';
	import Icon from 'mdi-svelte';
	import Toast from '$lib/components/Toast.svelte';

	let { data }: PageProps = $props();
	const formatter = new Intl.NumberFormat('en-US', {
		notation: 'compact',
		compactDisplay: 'short',
		maximumFractionDigits: 1 // Controls decimal points (e.g., 1.5K vs 1.52K)
	});

	let toast: Toast;
	$inspect(data);

	async function notifyUpdateAllImageEmbedding() {
		try {
			const resp = await fetch('/api/asset/embedding');
			if (resp.ok) {
				toast.add('Embedding update request has been made.', 'info');
			} else {
				toast.add('Embedding update request fails.', 'error');
			}
		} catch (e) {
			toast.add('Embedding update request fails.', 'error');
		}
	}

	async function notifyProcessAllAssets(missingOnly: boolean) {
		try {
			const resp = await fetch(`/api/asset/process?missingOnly=${missingOnly}`);
			if (resp.ok) {
				toast.add('Asset processing request has been made.', 'info');
			} else {
				toast.add('Asset processing  request fails.', 'error');
			}
		} catch (e) {
			toast.add('Asset processing  request fails.', 'error');
		}
	}
</script>

<svelte:head>
	<title>Albums: Administration</title>
</svelte:head>

<div class="w-full p-4 pt-8">
	<h2 class="p-4 text-2xl font-semibold">Statistics</h2>
	<div class="stats shadow">
		<div class="stat">
			<div class="stat-title">Total</div>
			<div class="stat-value text-primary">{formatter.format(data.total)}</div>
		</div>

		<div class="stat">
			<div class="stat-title">Pending</div>
			<div class="stat-value text-primary">{formatter.format(data.pending)}</div>
		</div>

		<div class="stat">
			<div class="stat-title">Failed Assets</div>
			<div class="stat-value text-error">{formatter.format(data.failed)}</div>
		</div>

		<div class="stat">
			<div class="stat-title">Uploading</div>
			<div class="stat-value text-warning">{formatter.format(data.uploading)}</div>
		</div>
	</div>

	<h2 class="p-4 text-2xl font-semibold">Assets by types</h2>
	<div class="stats shadow">
		<div class="stat">
			<div class="stat-figure text-primary">
				<Icon path={mdiImage} size={2} />
			</div>
			<div class="stat-title">Images</div>
			<div class="stat-value text-primary">{formatter.format(data.images)}</div>
			<div class="stat-desc">
				{formatter.format((data.embeddings * 100n) / data.images)}% searchable.
			</div>
		</div>

		<div class="stat">
			<div class="stat-figure text-primary">
				<Icon path={mdiVideo} size={2} />
			</div>
			<div class="stat-title">Video</div>
			<div class="stat-value text-primary">{formatter.format(data.video)}</div>
			<div class="stat-desc">Transcoded</div>
		</div>
	</div>

	<h2 class="p-4 text-2xl font-semibold">Operations</h2>
	<table class="table">
		<thead>
			<tr>
				<td>Operation</td>
				<td></td>
			</tr>
		</thead>
		<tbody>
			<tr>
				<td>Populate missing image embedding</td>
				<td>
					<button class="btn-small btn btn-primary" onclick={() => notifyUpdateAllImageEmbedding()}>
						<Icon path={mdiPlayCircle} />
						Start
					</button>
				</td>
			</tr>
			<tr>
				<td colspan="2" class="italic">
					This generates embedding for items that doesn't have it generated during processing, so
					the items are searchable.
				</td>
			</tr>
			<tr>
				<td>Re-process all items</td>
				<td>
					<button class="btn-small btn btn-primary" onclick={() => notifyProcessAllAssets(false)}>
						<Icon path={mdiPlayCircle} />
						Start
					</button>
				</td>
			</tr>
			<tr>
				<td colspan="2" class="italic">
					Reprocess all items to update the assets information. Avoid repettively triggering this
					action as it can causes massive jobs being queued.
				</td>
			</tr>
			<tr>
				<td>Re-process failed and pending Items</td>
				<td>
					<button class="btn-small btn btn-primary" onclick={() => notifyProcessAllAssets(true)}>
						<Icon path={mdiPlayCircle} />
						Start
					</button>
				</td>
			</tr>
			<tr>
				<td colspan="2" class="italic">
					Reprocess items that marked as 'failed' and 'new' items to update the assets information.
					Avoid repettively triggering this action as it can causes massive jobs being queued.
				</td>
			</tr>
		</tbody>
	</table>
</div>

<Toast bind:this={toast} />
