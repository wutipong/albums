<script lang="ts">
	import NavBar from '$lib/components/NavBar.svelte';
	import { mdiLogout } from '@mdi/js';
	import Icon from 'mdi-svelte';

	let { data } = $props();
</script>

<svelte:head>
	<title>Albums: User</title>
</svelte:head>

<div class="flex h-screen w-screen flex-col">
	<NavBar user={data.user} />

	<div class="overflow-auto p-4 pt-8">
		<article class="mx-auto prose h-full w-full md:w-200">
			<div class="flex flex-row gap-8">
				<div class="avatar">
					<div class="h-24 w-24 rounded-full">
						<img src={data.avatarSrc} alt="avatar" class="my-0!" />
					</div>
				</div>
				<div>
					<h2 class="mt-0">{data.user.name}</h2>
					{#if data.user.role && data.user.role == 'admin'}
						<h3>Administrator</h3>
					{/if}
					<p><a href={`mailto:${data.user.email}`}>{data.user.email}</a></p>
					<a class="btn btn-soft" href="/logout" data-sveltekit-preload-data="off">
						<Icon path={mdiLogout} />Logout
					</a>
				</div>
			</div>
		</article>
	</div>
</div>
