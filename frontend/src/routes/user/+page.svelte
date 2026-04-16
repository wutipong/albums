<script lang="ts">
	import { authClient } from '$lib/auth-client';
	import NavBar from '$lib/components/NavBar.svelte';
	import { mdiAccount, mdiAlert, mdiLogout } from '@mdi/js';
	import Icon from 'mdi-svelte';
	import { onMount } from 'svelte';
	import { createHash } from '@better-auth/utils/hash';
	import type { ApiKey } from '@better-auth/api-key/types';

	let name = $state('');
	let email = $state('');
	let avatarSrc = $state('');

	let apiKeys: Omit<ApiKey, 'key'>[] = $state([]);
	let apiKeyModal: HTMLDialogElement;
	let apiNewKey = $state('');

	onMount(async () => {
		const session = await authClient.getSession();
		if (!session.data) {
			console.log('session not found?');
			return;
		}

		name = session.data.user.name;
		email = session.data.user.email;
		const hashVal = await createHash('SHA-256', 'hex').digest(email);

		avatarSrc = `https://gravatar.com/avatar/${hashVal}`;

		const keys = await authClient.apiKey.list();
		if (keys.data) {
			apiKeys = keys.data.apiKeys;
		}
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

{#snippet title()}
	<div class="ms-4 flex flex-row gap-4 text-xl">
		<Icon path={mdiAccount} />
		User
	</div>
{/snippet}

<div class="flex h-screen w-screen flex-col">
	<NavBar {title}></NavBar>
	<div class="overflow-auto p-4 pt-8">
		<article class="mx-auto prose h-full w-full md:w-200">
			<div class="flex flex-row gap-8">
				<div class="avatar">
					<div class="w-24 h-24 rounded-full">
						<img src={avatarSrc} alt="avatar" class="my-0!" />
					</div>
				</div>
				<div>
					<h2 class="mt-0">{name}</h2>
					<p><a href={`mailto:${email}`}>{email}</a></p>
					<a class="btn btn-soft" href="/logout" data-sveltekit-preload-data="off">
						<Icon path={mdiLogout} />Logout
					</a>
				</div>
			</div>
			<hr />

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
								<button class="btn btn-sm" onclick={() => deleteApiKey(apiKey.id)}> Delete </button>
							</td>
						</tr>
					{/each}
				</tbody>
			</table>
			<button class="btn w-full btn-primary" onclick={() => addNewApiKey()}>Add new API key</button>
		</article>
	</div>
</div>

<dialog class="modal" bind:this={apiKeyModal}>
	<div class="modal-box">
		<h3 class="text-lg font-bold">New key added</h3>
		<div role="alert" class="alert py-4 alert-warning">
			<Icon path={mdiAlert} />
			<span>Warning: This API key will not be visible again! </span>
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
