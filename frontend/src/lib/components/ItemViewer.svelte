<script lang="ts">
	import { mdiClose, mdiDotsVertical, mdiDownload } from '@mdi/js';
	import Icon from 'mdi-svelte';
	let { id, show = $bindable(false) } = $props();
</script>

{#if show}
	<div
		role="presentation"
		class="absolute top-0 right-0 bottom-0 left-0 backdrop-blur-lg backdrop-brightness-50"
	>
		<div class="h-full w-full">
			<img src={`/api/asset/${id}/view`} alt={id} class="m-auto h-full w-full object-contain" />
		</div>

		<button
			class="btn absolute top-4 right-4 btn-circle btn-ghost btn-lg"
			onclick={() => (show = false)}
		>
			<Icon path={mdiClose} />
		</button>

		<div class="fab">
			<!-- a focusable div with tabindex is necessary to work on all browsers. role="button" is necessary for accessibility -->
			<div tabindex="0" role="button" class="btn btn-circle btn-ghost btn-lg">
				<Icon path={mdiDotsVertical} />
			</div>

			<a href={`/api/asset/${id}/original/`} target="_blank" class="btn btn-circle btn-lg">
				<Icon path={mdiDownload} />
			</a>
			<button class="btn btn-circle btn-lg">B</button>
			<button class="btn btn-circle btn-lg">C</button>
		</div>
	</div>
{/if}
