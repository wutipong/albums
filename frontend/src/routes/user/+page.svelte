<script lang="ts">
	import { authClient } from '$lib/auth-client';
	import NavBar from '$lib/components/NavBar.svelte';
	import { mdiAlert, mdiKeyMinus, mdiKeyPlus, mdiLogout } from '@mdi/js';
	import Icon from 'mdi-svelte';
	import { onMount } from 'svelte';
	import type { ApiKey } from '@better-auth/api-key/types';
	import { createHash } from '@better-auth/utils/hash';

	let apiKeys: Omit<ApiKey, 'key'>[] = $state([]);
	let apiKeyModal: HTMLDialogElement;
	let apiNewKey = $state('');

	let { data } = $props();

	let avatarSrc = $state('')

	onMount(async () => {
		const keys = await authClient.apiKey.list();
		if (keys.data) {
			apiKeys = keys.data.apiKeys;
		}

		const hashVal = await createHash('SHA-256', 'hex').digest(data.user.email);

		avatarSrc = `https://gravatar.com/avatar/${hashVal} `;
	});

	async function addNewApiKey() {
		const { data, error } = await authClient.apiKey.create({});
		if (data) {
			apiNewKey = data.key;
			apiKeyModal.showModal();
		}

		const keys = await authClient.apiKey.list();
		if (keys.data) {
			apiKeys = keys.data.apiKeys;
		}
	}

	async function deleteApiKey(id: string) {
		const { data, error } = await authClient.apiKey.delete({ keyId: id });

		const keys = await authClient.apiKey.list();
		if (keys.data) {
			apiKeys = keys.data.apiKeys;
		}
	}
</script>

<svelte:head>
	<title>Albums: User</title>
</svelte:head>

<div class="flex h-screen w-screen flex-col">
	<NavBar user={data.user}/>

	<div class="overflow-auto p-4 pt-8">
		<article class="mx-auto prose h-full w-full md:w-200">
			<div class="flex flex-row gap-8">
				<div class="avatar">
					<div class="h-24 w-24 rounded-full">
						<img src={avatarSrc} alt="avatar" class="my-0!" />
					</div>
				</div>
				<div>
					<h2 class="mt-0">{data.user.name}</h2>
					<p><a href={`mailto:${data.user.email}`}>{data.user.email}</a></p>
					<a class="btn btn-soft" href="/logout" data-sveltekit-preload-data="off">
						<Icon path={mdiLogout} />Logout
					</a>
				</div>
			</div>
			<hr />
			{#if data.user.role == 'admin'}
				<h2>API Keys</h2>
				<table>
					<thead>
						<tr>
							<th>Key</th>
							<th>Created At</th>
							<th>Action</th>
						</tr>
					</thead>
					<tbody>
						{#each apiKeys as apiKey}
							<tr>
								<td>{apiKey.start}...</td>
								<td>{apiKey.createdAt.toLocaleString()}</td>
								<td>
									<button class="btn btn-sm" onclick={() => deleteApiKey(apiKey.id)}>
										<Icon path={mdiKeyMinus} /> Delete
									</button>
								</td>
							</tr>
						{/each}
						<tr>
							<td></td>
							<td></td>
							<td>
								<button class="btn btn-sm btn-primary" onclick={() => addNewApiKey()}>
									<Icon path={mdiKeyPlus} /> Add
								</button>
							</td>
						</tr>
					</tbody>
				</table>
			{/if}
		</article>
	</div>
</div>

<dialog class="modal" bind:this={apiKeyModal}>
	<div class="modal-box">
		<h3 class="text-lg font-bold">New key added</h3>
		<div role="alert" class="alert py-4 alert-warning">
			<Icon path={mdiAlert} />
			<span>This API key will not be visible again! </span>
		</div>
		<p class="py-4 font-mono text-wrap break-all">{apiNewKey}</p>
		<div class="modal-action">
			<form method="dialog">
				<!-- if there is a button in form, it will close the modal -->
				<button class="btn">Close</button>
			</form>
		</div>
	</div>
</dialog>
