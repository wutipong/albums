<script lang="ts">
	let { id } = $props();
</script>

<div class="absolute top-0 right-0 bottom-0 left-0 backdrop-blur-xl backdrop-brightness-50">
	<div class="absolute m-auto">
		<img src={`/api/asset/${id}/view`} alt={id} />
	</div>
</div>
