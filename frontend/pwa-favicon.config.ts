import { defaultAssetName, defineConfig } from '@vite-pwa/assets-generator/config';

export default defineConfig({
	headLinkOptions: {
		preset: '2023',
		resolveSvgName: () => 'pwa/favicon.svg'
	},
	preset: {
		assetName: (type, size) => `pwa/${defaultAssetName(type, size)}`,
		transparent: {
			sizes: [64],
			resizeOptions: {
				fit: 'contain',
				background: '#222222'
			},
			favicons: [[48, 'pwa/favicon.ico']]
		},
		maskable: {
			sizes: []
		},
		apple: {
			sizes: []
		}
	},
	images: ['static/favicon.svg']
});
