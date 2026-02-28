import api from './api';

export default {
	// Listar materiais com paginação
	listar(limit = 10, offset = 0) {
		return api.get(`/materiais?limit=${limit}&offset=${offset}`);
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
		return api.get(url);
	},

	// Obter detalhes de um material específico
	obterDetalhes(id) {
		return api.get(`/materiais/detalhes?id=${id}`);
	},

	// Obter recomendações para um usuário
	obterRecomendacoes(usuarioId, limit = 5) {
		return api.get(`/materiais/recomendacoes?usuario_id=${usuarioId}&limit=${limit}`);
	},

	// Favoritar ou desfavoritar material
	favoritar(usuarioId, materialId, status) {
		return api.post('/materiais/favoritar', {
			usuario_id: usuarioId,
			material_id: materialId,
			favoritar: status
		});
	},

	// Listar favoritos do usuário
	listarFavoritos(usuarioId) {
		return api.get(`/materiais/favoritos?usuario_id=${usuarioId}`);
	},

	// Avaliar um material
	avaliar(usuarioId, materialId, nota, comentario = '') {
		return api.post('/materiais/avaliar', {
			usuario_id: usuarioId,
			material_id: materialId,
			nota: nota,
			comentario: comentario
		});
	},

	// Realizar empréstimo de um material
	emprestar(usuarioId, materialId) {
		return api.post('/materiais/emprestar', {
			usuario_id: usuarioId,
			material_id: materialId
		});
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
		const baseURL = api.defaults.baseURL || 'http://localhost:8080/api';
		return `${baseURL}/materiais/pdf/proxy?url=${encodeURIComponent(pdfUrl)}`;
	},

	// Obter estatísticas do usuário
	obterEstatisticas(usuarioId) {
		return api.get(`/usuario/estatisticas?usuario_id=${usuarioId}`);
	},

	// Chatbot IA com o livro
	perguntar(id, pergunta) {
		return api.post(`/materiais/${id}/chat`, { pergunta });
	},

	// Gerar resumo IA do livro
	obterResumo(id) {
		return api.get(`/materiais/${id}/resumo`);
	}
};
