// import Vue from 'vue';
// import Vuetify from 'vuetify/lib';

// Vue.use(Vuetify);

// export default new Vuetify({});

// import '@mdi/font/css/materialdesignicons.css'
import 'vuetify/styles'

// const vuetify = createVuetify({
//   ssr: true,
// })

import { createVuetify } from 'vuetify'
import * as components from 'vuetify/components'
import * as directives from 'vuetify/directives'

export default createVuetify({
	components,
	directives,
	icons: {
		defaultSet: 'mdi',
	},
	theme: {
		defaultTheme: 'dark',
		themes: {
			dark: {
				dark: true,
				colors: {
					primary: '#007AFF',
					secondary: '#5AC8FA',
					background: '#0B192C',
					surface: '#1E2A38',
				}
			}
		}
	}
})
