<script lang="ts">
	let { id, show = $bindable(false) } = $props();
</script>

{#if show}
	<div class="absolute top-0 right-0 bottom-0 left-0 backdrop-blur-lg backdrop-brightness-50">
		<div class="h-full w-full">
			<img src={`/api/asset/${id}/view`} alt={id} class="m-auto max-h-full max-w-full" />
		</div>
	</div>
{/if}
