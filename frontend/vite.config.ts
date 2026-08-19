import tailwindcss from '@tailwindcss/vite';
import { sveltekit } from '@sveltejs/kit/vite';
import { SvelteKitPWA } from '@vite-pwa/sveltekit';
import { defineConfig } from 'vite';
import { vite as vidstack } from 'vidstack/plugins';

export default defineConfig({
	plugins: [
		vidstack(),
		tailwindcss(),
		sveltekit(),
		SvelteKitPWA({
			registerType: 'autoUpdate',
			injectRegister: false,
			manifest: {
				name: 'Albums',
				short_name: 'Albums',
				start_url: '/',
				scope: '/',
				display: 'standalone',
				background_color: '#ffffff',
				theme_color: '#ff3e00',
				icons: [
					{
						src: '/pwa/pwa-64x64.png',
						sizes: '64x64',
						type: 'image/png'
					},
					{
						src: '/pwa/pwa-192x192.png',
						sizes: '192x192',
						type: 'image/png'
					},
					{
						src: '/pwa/pwa-512x512.png',
						sizes: '512x512',
						type: 'image/png'
					},
					{
						src: '/pwa/maskable-icon-512x512.png',
						sizes: '512x512',
						type: 'image/png',
						purpose: 'maskable'
					},
					{
						src: '/pwa/apple-touch-icon-180x180.png',
						sizes: '180x180',
						type: 'image/png'
					}
				]
			},
			workbox: {
				globPatterns: [
					'client/**/*.{js,css,ico,png,svg,webp,webmanifest}',
					'prerendered/**/*.{html,json}'
				]
			}
		})
	]
});
