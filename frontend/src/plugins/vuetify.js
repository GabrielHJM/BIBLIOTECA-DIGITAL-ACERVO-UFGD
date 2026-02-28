// import Vue from 'vue';
// import Vuetify from 'vuetify/lib';

// Vue.use(Vuetify);

// export default new Vuetify({});

// import '@mdi/font/css/materialdesignicons.css'
// import 'vuetify/styles'

// const vuetify = createVuetify({
//   ssr: true,
// })

import { createVuetify } from 'vuetify'
import { aliases, mdi } from 'vuetify/iconsets/mdi'
import * as components from 'vuetify/components'
import * as directives from 'vuetify/directives'

export default createVuetify({
	components,
	directives,
	icons: {
		defaultSet: 'mdi',
		aliases,
		sets: {
			mdi,
		},
	},
	theme: {
		defaultTheme: 'dark',
		themes: {
			light: {
				dark: false,
				colors: {
					primary: '#00B8D4',
					secondary: '#5AC8FA',
					background: '#f0f2f5',
					surface: '#ffffff',
				}
			},
			dark: {
				dark: true,
				colors: {
					primary: '#00B8D4',
					secondary: '#5AC8FA',
					background: '#2d5178',
					surface: '#1e3857',
				}
			}
		}
	}
})
