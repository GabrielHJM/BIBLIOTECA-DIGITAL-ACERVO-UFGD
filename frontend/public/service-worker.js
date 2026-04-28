const CACHE_NAME = 'bib-digital-v2';
const ASSETS = [
	'/',
	'/index.html',
	'/manifest.json',
	'/favicon.ico',
	'/logo-biblioteca.png'
];

// Instalação: Cacheia assets estáticos
self.addEventListener('install', (event) => {
	event.waitUntil(
		caches.open(CACHE_NAME).then((cache) => {
			return cache.addAll(ASSETS);
		})
	);
	self.skipWaiting();
});

// Ativação: Limpa caches antigos
self.addEventListener('activate', (event) => {
	event.waitUntil(
		caches.keys().then((cacheNames) => {
			return Promise.all(
				cacheNames.map((cacheName) => {
					if (cacheName !== CACHE_NAME) {
						return caches.delete(cacheName);
					}
				})
			);
		})
	);
});

self.addEventListener('fetch', (event) => {
	// Estratégia Network First para garantir que os usuários sempre recebam
	// as atualizações mais recentes do GitHub e Backend automaticamente sem CTRL+F5.
	event.respondWith(
		fetch(event.request)
			.then((networkResponse) => {
				// Atualiza o cache com a versão mais nova
				return caches.open(CACHE_NAME).then((cache) => {
					cache.put(event.request, networkResponse.clone());
					return networkResponse;
				});
			})
			.catch(() => {
				// Fallback offline: se falhar (sem internet), tenta pegar do cache
				return caches.match(event.request);
			})
	);
});
