import { getRequestEvent } from '$app/server';
import { env } from '$env/dynamic/private';
import { betterAuth } from 'better-auth';
import { admin, genericOAuth } from 'better-auth/plugins';
import { sveltekitCookies } from 'better-auth/svelte-kit';
import { apiKey } from '@better-auth/api-key';
import { Pool } from 'pg';
import { building } from '$app/environment';

export const auth = (building)? undefined:  betterAuth({
	secret: env.BETTER_AUTH_SECRET || 'placeholder-secret',
	database: new Pool({
		connectionString: env.DATABASE_URL
	}),

	plugins: [
		admin(),
		sveltekitCookies(getRequestEvent),
		apiKey({
			rateLimit: { enabled: false },
			enableSessionForAPIKeys: true
		}),
		genericOAuth({
			config: [
				{
					providerId: env.OIDC_PROVIDER_ID || 'placeholder-provider',
					clientId: env.OIDC_CLIENT_ID || 'placeholder-client',
					clientSecret: env.OIDC_SECRET || 'placeholder-secret',
					discoveryUrl: env.OIDC_DISCOVERY_URL || 'placeholder-discovery',
					scopes: ['openid', 'email', 'profile', 'groups'],
					mapProfileToUser: async (profile) => {
						const groups = profile.groups;
						const isAdmin = groups?.includes('admin');

						return {
							name: profile.name,
							email: profile.email,

							role: isAdmin ? 'admin' : 'user' // Maps directly to user.additionalFields.role
						};
					}
				}
			]
		})
	]
});
