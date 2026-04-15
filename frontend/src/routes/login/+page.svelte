<script lang="ts">
	import { authClient } from '$lib/auth-client';
	import Icon from 'mdi-svelte';
	import { mdiLogin } from '@mdi/js';

	const { data } = $props();

	async function performLogin() {
		await authClient.signIn.oauth2(data);
	}
</script>

<div class="relative h-screen w-screen items-center justify-center">
	<button
		class="btn absolute top-1/2 left-1/2 -translate-1/2 gap-4 btn-outline btn-xl btn-primary"
		onclick={() => performLogin()}
	>
		<Icon path={mdiLogin} />
		Login
	</button>

	<div class="absolute right-0 bottom-0 left-0 m-4">Albums &COPY; 20027 Wutipong Wongsakuldej</div>
</div>
