const CACHE_NAME = 'bib-digital-v1';
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

// Estratégia Stale-While-Revalidate
self.addEventListener('fetch', (event) => {
	// Ignora requisições de API para não cachear dados sensíveis de forma agressiva aqui
	// (Pode ser ajustado para cachear prefetched data se necessário)
	if (event.request.url.includes('/api/')) {
		return;
	}

	event.respondWith(
		caches.open(CACHE_NAME).then((cache) => {
			return cache.match(event.request).then((cachedResponse) => {
				const fetchedResponse = fetch(event.request).then((networkResponse) => {
					cache.put(event.request, networkResponse.clone());
					return networkResponse;
				});

				return cachedResponse || fetchedResponse;
			});
		})
	);
});
