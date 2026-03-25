import api from './api';

const cache = new Map();

export default {
	// Helper for caching GET requests
	async _cachedGet(url, ttl = 300000) { // Default 5 mins TTL
		const now = Date.now();
		if (cache.has(url)) {
			const { data, timestamp } = cache.get(url);
			if (now - timestamp < ttl) {
				return Promise.resolve(data);
			}
			cache.delete(url); // expired
		}

		const response = await api.get(url);

		// Retornamos um objeto que simula a resposta do Axios (com .data)
		const data = response.data !== undefined ? response.data : response;
		const result = { ...response, data };

		// Limit cache size to 50 entries to avoid memory bloat
		if (cache.size > 50) {
			const firstKey = cache.keys().next().value;
			cache.delete(firstKey);
		}

		cache.set(url, { data: result, timestamp: now });
		return result;
	},

	// Listar materiais com paginação
	listar(limit = 10, offset = 0) {
		return this._cachedGet(`/materiais?limit=${limit}&offset=${offset}`);
	},

	// Pesquisar materiais por termo ou categoria com filtros avançados
	pesquisar(termo = '', categoria = '', fonte = '', anoInicio = 0, anoFim = 0, limit = 10, offset = 0, sort = '') {
		let url = `/materiais?limit=${limit}&offset=${offset}`;
		if (termo) url += `&q=${encodeURIComponent(termo)}`;
		if (categoria) url += `&categoria=${encodeURIComponent(categoria)}`;
		if (fonte) url += `&fonte=${encodeURIComponent(fonte)}`;
		if (anoInicio) url += `&ano_inicio=${anoInicio}`;
		if (anoFim) url += `&ano_fim=${anoFim}`;
		if (sort) url += `&sort=${encodeURIComponent(sort)}`;
		return this._cachedGet(url, 60000); // 1 minute TTL for search queries
	},

	// Obter dados consolidados para o dashboard (múltiplas categorias em um request)
	dashboard() {
		return this._cachedGet('/materiais/dashboard', 300000); // 5 mins TTL
	},

	// Obter detalhes de um material específico
	obterDetalhes(id) {
		return api.get(`/materiais/detalhes?id=${id}`);
	},



	// Favoritar ou desfavoritar material
	favoritar(usuarioId, materialId, status) {
		return api.post('/materiais/favoritar', {
			usuario_id: usuarioId,
			material_id: materialId,
			favoritar: status
		});
	},

	// Avaliar material com nota (1-5)
	avaliar(usuarioId, materialId, nota) {
		return api.post('/materiais/avaliar', {
			usuario_id: usuarioId,
			material_id: materialId,
			nota: nota
		});
	},

	// Listar favoritos do usuário
	listarFavoritos(usuarioId) {
		return api.get(`/materiais/favoritos?usuario_id=${usuarioId}`);
	},





	// Registrar leitura no histórico
	registrarLeitura(usuarioId, materialId) {
		return api.post('/materiais/historico', {
			usuario_id: usuarioId,
			material_id: materialId
		});
	},

	// Listar histórico de leitura do usuário
	listarHistorico(usuarioId) {
		return api.get(`/materiais/historico?usuario_id=${usuarioId}`);
	},

	// Listar avaliações de um material
	listarAvaliacoes(materialId) {
		return api.get(`/materiais/avaliacoes?id=${materialId}`);
	},

	// Obter URL do proxy de PDF
	getProxyPdfUrl(pdfUrl) {
		return `/api/materiais/pdf/proxy?url=${encodeURIComponent(pdfUrl)}`;
	},

	// Obter estatísticas do usuário
	obterEstatisticas(usuarioId) {
		return api.get(`/usuario/estatisticas?usuario_id=${usuarioId}`);
	},


};
