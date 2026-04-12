<script lang="ts">
	let { id } = $props();
</script>

<div class="absolute top-0 bottom-0 left-0 right-0 backdrop-brightness-50 backdrop-blur-lg">
	<div class="w-full h-full">
		<img src={`/api/asset/${id}/view`} alt={id} class="max-w-full max-h-full m-auto"/>
	</div>
</div>
