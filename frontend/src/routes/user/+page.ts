import type { PageLoad } from './$types';
import { authClient } from '$lib/auth-client';
import { createHash } from '@better-auth/utils/hash';

export const load: PageLoad = async () => {
	const session = await authClient.getSession();
	if (!session.data) {
		console.log('session not found?');
		return;
	}

	const name = session.data.user.name;
	const email = session.data.user.email;
	const hashVal = await createHash('SHA-256', 'hex').digest(email);
	const avatarSrc = `https://gravatar.com/avatar/${hashVal}`;
    const role = session.data.user.role?? "user"

	return {
		name,
		email,
		avatarSrc,
        role,
	};
};
