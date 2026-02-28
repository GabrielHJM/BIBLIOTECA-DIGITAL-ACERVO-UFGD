import { createApp } from 'vue'
import router from './router'

// Vuetify
import vuetify from './plugins/vuetify'

// Components
import App from './App.vue'

createApp(App).use(router).use(vuetify).mount('#app')

// PWA Service Worker Registration
if ('serviceWorker' in navigator) {
	window.addEventListener('load', () => {
		navigator.serviceWorker.register('/service-worker.js')
			.then(reg => console.log('Expert PWA Service Worker registrado!', reg))
			.catch(err => console.log('Falha ao registrar PWA Service Worker:', err));
	});
}
