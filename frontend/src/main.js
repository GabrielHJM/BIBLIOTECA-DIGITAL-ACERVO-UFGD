import { createApp } from 'vue'
import router from './router'

// Vuetify
import '@mdi/font/css/materialdesignicons.css'
import vuetify from './plugins/vuetify'

// Design System
import './assets/styles/design-system.css'

// Components
import App from './App.vue'
import BookCard from './components/BookCard.vue'
import CloudLoader from './components/CloudLoader.vue'

const app = createApp(App)
app.component('BookCard', BookCard)
app.component('CloudLoader', CloudLoader)
app.use(router).use(vuetify).mount('#app')

// PWA Service Worker Registration
if ('serviceWorker' in navigator) {
	window.addEventListener('load', () => {
		navigator.serviceWorker.register('/service-worker.js')
			.then(reg => console.log('Expert PWA Service Worker registrado!', reg))
			.catch(err => console.log('Falha ao registrar PWA Service Worker:', err));
	});
}
