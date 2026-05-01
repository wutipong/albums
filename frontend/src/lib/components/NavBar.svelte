<script lang="ts">
	import { goto } from '$app/navigation';
	import { authClient } from '$lib/auth-client';
	import { mdiImageSearch, mdiImageAlbum, mdiLogout, mdiAccount } from '@mdi/js';
	import Icon from 'mdi-svelte';
	import { onMount } from 'svelte';
	import { createHash } from '@better-auth/utils/hash';

	let searchInput: HTMLInputElement;

	let { title, album = null } = $props();
	let avatarSrc = $state('');

	let albumsUrl = $derived.by(() => {
		if (album) {
			return `/album#${album.id}`;
		} else {
			return '/album';
		}
	});

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
	});

	let searchDialog: HTMLDialogElement;
</script>

<div class="navbar shadow-sm">
	<div class="flex-1">
		{@render title()}
	</div>
	<div class="flex-none">
		<ul class="menu menu-horizontal px-1">
			<li>
				<a href={albumsUrl}>
					<Icon path={mdiImageAlbum} />
					<div class="hidden md:block">Albums</div></a
				>
			</li>
		</ul>
		<ul class="menu menu-horizontal px-1 md:hidden">
			<li>
				<button onclick={() => searchDialog.showModal()}>
					<Icon path={mdiImageSearch} />
					<div class="hidden md:block">Search</div>
				</button>
			</li>
		</ul>
	</div>

	<div class="me-4 hidden gap-2 md:flex">
		<div class="join">
			<div>
				<input class="input join-item" type="text" placeholder="search" bind:this={searchInput} />
			</div>
			<button class="btn join-item" onclick={() => doSearch()}>
				<Icon path={mdiImageSearch} />
			</button>
		</div>
	</div>
	<div class="dropdown dropdown-end">
		<div tabindex="0" role="button" class="btn avatar btn-circle btn-ghost">
			<div class="w-10 rounded-full">
				<img alt="User" src={avatarSrc} />
			</div>
		</div>
		<ul
			tabindex="-1"
			class="dropdown-content menu z-1 mt-3 w-52 menu-sm rounded-box bg-base-100 p-2 shadow-xl"
		>
			<li>
				<a href="/user"><Icon path={mdiAccount} /> User</a>
			</li>
			<li></li>
			<li>
				<a href="/logout" data-sveltekit-preload-data="off"><Icon path={mdiLogout} /> Logout</a>
			</li>
		</ul>
	</div>
</div>

<dialog class="modal" bind:this={searchDialog}>
	<div class="modal-box">
		<h3 class="text-lg font-bold">Search</h3>
		<div class="py-4">
			<input
				class="input-bordered input w-full"
				type="text"
				placeholder="search"
				bind:this={searchInput}
				onkeydown={(e) => {
					if (e.key === 'Enter') {
						doSearch();
						searchDialog.close();
					}
				}}
			/>
		</div>
	</div>
</dialog>
