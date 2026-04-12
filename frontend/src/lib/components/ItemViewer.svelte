<script lang="ts">
	let { id, show = $bindable(false) } = $props();
</script>

{#if show}
	<div class="absolute top-0 right-0 bottom-0 left-0 backdrop-blur-lg backdrop-brightness-50">
		<div class="h-full w-full">
			<img src={`/api/asset/${id}/view`} alt={id} class="m-auto h-full w-full object-contain" />
		</div>

		<button 
			class="btn btn-ghost absolute top-4 right-4 rounded-full" 
			onclick={() => show = false}>Close
		</button>
	</div>
{/if}
