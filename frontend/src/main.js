import { createApp } from 'vue'
import router from './router'

// Vuetify
import 'vuetify/styles'
import { createVuetify } from 'vuetify'
import * as components from 'vuetify/components'
import * as directives from 'vuetify/directives'

// Components
import App from './App.vue'

const vuetify = createVuetify({
	components,
	directives,
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

createApp(App).use(router).use(vuetify).mount('#app')
