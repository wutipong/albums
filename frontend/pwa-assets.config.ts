import {
	createAppleSplashScreens,
	defaultAssetName,
	defaultSplashScreenName,
	defineConfig,
	minimal2023Preset
} from '@vite-pwa/assets-generator/config';

const preset = {
	...minimal2023Preset,
	transparent: {
		...minimal2023Preset.transparent,
		sizes: [192, 512],
		resizeOptions: {
			fit: 'contain',
			background: '#222222'
		},
		favicons: []
	},
	maskable: {
		...minimal2023Preset.maskable,
		resizeOptions: {
			fit: 'contain',
			background: '#222222'
		}
	},
	apple: {
		...minimal2023Preset.apple,
		resizeOptions: {
			fit: 'contain',
			background: '#222222'
		}
	},
	assetName: (type, size) => `pwa/${defaultAssetName(type, size)}`
};

preset.appleSplashScreens = createAppleSplashScreens(
	{
		// padding: 0.3,
		resizeOptions: { fit: 'contain', background: '#222222' },
		linkMediaOptions: {
			log: true,
			addMediaScreen: true,
			xhtml: true
		},
		name: (landscape, size) => `pwa/${defaultSplashScreenName(landscape, size)}`
	},
	['iPad Air 9.7"']
);

export default defineConfig({
	headLinkOptions: {
		preset: '2023',
		resolveSvgName: () => 'pwa/favicon.svg'
	},
	preset: preset,
	images: ['static/logo.svg']
});
