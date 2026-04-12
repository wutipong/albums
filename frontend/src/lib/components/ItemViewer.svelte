<script lang="ts">
	let { id, show = $bindable(false) } = $props();
</script>

{#if show}
	<div class="absolute top-0 right-0 bottom-0 left-0 backdrop-blur-lg backdrop-brightness-50">
		<div class="h-full w-full">
			<img src={`/api/asset/${id}/view`} alt={id} class="m-auto h-full w-full object-contain" />
		</div>

		<button class="btn absolute top-4 right-4 btn-lg btn-circle btn-ghost" onclick={() => (show = false)}>
			Close
		</button>

		<div class="fab">
			<!-- a focusable div with tabindex is necessary to work on all browsers. role="button" is necessary for accessibility -->
			<div tabindex="0" role="button" class="btn btn-circle btn-ghost btn-lg">:</div>

			<!-- buttons that show up when FAB is open -->
			<button class="btn btn-circle btn-lg">A</button>
			<button class="btn btn-circle btn-lg">B</button>
			<button class="btn btn-circle btn-lg">C</button>
		</div>
	</div>
{/if}
