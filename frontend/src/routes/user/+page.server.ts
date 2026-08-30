import type { PageServerLoad } from './$types';
import { SHA256 } from 'bun';
export const load: PageServerLoad = async ({ locals }) => {
	const hashVal = SHA256.hash(locals.user.email, 'hex');
	const avatarSrc = `https://gravatar.com/avatar/${hashVal} `;

	return {
		user: locals.user,
		session: locals.session,
		avatarSrc
	};
};
