<script lang="ts">
	import { goto } from '$app/navigation';
	import { authClient } from '$lib/auth-client';
	import { mdiImageSearch, mdiMagnify, mdiImageAlbum, mdiLogout } from '@mdi/js';
	import Icon from 'mdi-svelte';
	import { onMount } from 'svelte';
	import { createHash } from '@better-auth/utils/hash';

	let searchInput: HTMLInputElement;

	let { title } = $props();
	let avatarSrc = $state('');

	function doSearch() {
		const search = searchInput.value;
		const url = new URL('/search', location.origin);
		url.searchParams.append('search', search);

		goto(url);
	}

	onMount(async () => {
		const session = await authClient.getSession();
		if (!session.data) {
			console.log('session not found?');
			return;
		}

		const email = session.data?.user.email;
		const hashVal = await createHash('SHA-256', 'hex').digest(email);

		avatarSrc = `https://gravatar.com/avatar/${hashVal} `;

		console.log(avatarSrc);
	});
</script>

<div class="navbar bg-base-100 shadow-sm">
	<div class="flex-1">
		{@render title()}
	</div>
	<div class="flex-none">
		<ul class="menu menu-horizontal px-1">
			<li><a href="/album"><Icon path={mdiImageAlbum} />Albums</a></li>
		</ul>
	</div>

	<div class="flex gap-2 me-4">
		<div class="join">
			<label class="input join-item px-1">
				<input type="text" placeholder="search" bind:this={searchInput} />
			</label>
			<button class="btn join-item btn-neutral" onclick={() => doSearch()}>
				<Icon path={mdiImageSearch} />
			</button>
		</div>
		<div class="dropdown dropdown-end">
			<div tabindex="0" role="button" class="btn avatar btn-circle btn-ghost">
				<div class="w-10 rounded-full">
					<img alt="User" src={avatarSrc} />
				</div>
			</div>
			<ul
				tabindex="-1"
				class="dropdown-content menu z-1 mt-3 w-52 menu-sm rounded-box bg-base-300 p-2 shadow"
			>
				<li>
					<a href="/logout" data-sveltekit-preload-data="off"><Icon path={mdiLogout}/> Logout</a>
				</li>
			</ul>
		</div>
	</div>
</div>
