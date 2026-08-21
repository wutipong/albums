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
			const resp = await fetch('/api/album/embedding');
			if (resp.ok) {
				toast.add('Embedding update request has been made.', 'info');
			} else {
				toast.add('Embedding update request fails.', 'error');
			}
		} catch (e) {
			toast.add('Embedding update request fails.', 'error');
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
				<td>Populate Missing Image Embedding</td>
				<td>
					<button class="btn-small btn btn-primary" onclick={() => notifyUpdateAllImageEmbedding()}>
						<Icon path={mdiPlayCircle} />
						Start
					</button>
				</td>
			</tr>
		</tbody>
	</table>
</div>

<Toast bind:this={toast} />
