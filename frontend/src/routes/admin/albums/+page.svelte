<script lang="ts">
	import { mdiPlayCircle } from '@mdi/js';
	import Icon from 'mdi-svelte';
	import type { PageProps } from './$types';
	import Toast from '$lib/components/Toast.svelte';

	let { data }: PageProps = $props();
	const formatter = new Intl.NumberFormat('en-US', {
		notation: 'compact',
		compactDisplay: 'short',
		maximumFractionDigits: 1 // Controls decimal points (e.g., 1.5K vs 1.52K)
	});

	$inspect(data);

	let toast: Toast;

	async function notifyPopulateMissingCover() {
		try {
			const resp = await fetch('/api/album/cover');
			if (resp.ok) {
				toast.add('Album cover update request has been made.', 'info');
			} else {
				toast.add('Album cover request fails.', 'error');
			}
		} catch (e) {
			toast.add('Album cover update request fails.', 'error');
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
			<div class="stat-title">Missing Cover</div>
			<div class="stat-value text-warning">{formatter.format(data.missingCover)}</div>
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
				<td>Populate Missing Album Cover</td>
				<td>
					<button class="btn-small btn btn-primary" onclick={() => notifyPopulateMissingCover()}>
						<Icon path={mdiPlayCircle} />
						Start
					</button>
				</td>
			</tr>
		</tbody>
	</table>
</div>

<Toast bind:this={toast} />
