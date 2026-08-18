<script lang="ts">
	import { mdiClose } from '@mdi/js';
	import Icon from 'mdi-svelte';
	import { Temporal } from 'temporal-polyfill';

	interface Asset {
		created_at: Date;
		id: string;
		modified_at: Date;
		album_id: string;
		filename: string;
		image_frames: number;
		original: string;
		preview: string;
		process_status: 'pending' | 'processed' | 'processing' | 'uploading';
		thumbnail: string;
		type: 'animated' | 'audio' | 'image' | 'video';
		video_duration: Temporal.Duration;
		view: string;
	}

	let asset: Asset = $state({
		created_at: new Date(),
		id: '',
		modified_at: new Date(),
		album_id: '',
		filename: '',
		image_frames: 0,
		original: '',
		preview: '',
		process_status: 'pending',
		thumbnail: '',
		type: 'image',
		video_duration: Temporal.Duration.from({ seconds: 0 }),
		view: ''
	});

	export function show(a: Asset) {
		asset = a;
		console;
		dialog.showModal();
	}

	let dialog: HTMLDialogElement;
</script>

<dialog class="modal" bind:this={dialog}>
	<div class="modal-box">
		<form method="dialog">
			<button class="btn absolute top-2 right-2 btn-circle btn-ghost btn-sm">
				<Icon path={mdiClose} />
				<span class="sr-only">Close</span>
			</button>
		</form>
		<h3 class="text-lg font-bold">Infomation</h3>
		<div class="py-4">
			<table class="table w-full">
				<tbody>
					<tr>
						<th>Asset ID</th>
						<td>{asset.id}</td>
					</tr>
					<tr>
						<th>Filename</th>
						<td>{asset.filename}</td>
					</tr>
					<tr>
						<th>Created At</th>
						<td>{asset.created_at.toLocaleString()}</td>
					</tr>
					<tr>
						<th>Modified At</th>
						<td>{asset.modified_at.toLocaleString()}</td>
					</tr>
					<tr>
						<th>Type</th>
						<td>{asset.type}</td>
					</tr>
					{#if asset.type === 'video'}
						<tr>
							<th>Video Duration</th>
							<td>{asset.video_duration.toLocaleString('en-US', { style: 'digital' })}</td>
						</tr>
					{/if}
					{#if asset.image_frames > 1}
						<tr>
							<th>Image Frames</th>
							<td>{asset.image_frames}</td>
						</tr>
					{/if}
					<tr>
						<th>Original </th>
						<td>{asset.original}</td>
					</tr>
					<tr>
						<th>Preview </th>
						<td>{asset.preview}</td>
					</tr>
					<tr>
						<th>Thumbnail </th>
						<td>{asset.thumbnail}</td>
					</tr>
				</tbody>
			</table>
		</div>
	</div>
</dialog>
