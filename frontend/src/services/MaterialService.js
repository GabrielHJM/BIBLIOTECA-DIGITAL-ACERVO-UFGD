import api from './api';

export default {
	// Listar materiais com paginação
	listar(limit = 10, offset = 0) {
		return api.get(`/materiais?limit=${limit}&offset=${offset}`);
	},

	// Pesquisar materiais por termo ou categoria
	pesquisar(termo = '', categoria = '', limit = 10, offset = 0) {
		let url = `/materiais?limit=${limit}&offset=${offset}`;
		if (termo) url += `&q=${encodeURIComponent(termo)}`;
		if (categoria) url += `&categoria=${encodeURIComponent(categoria)}`;
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
	}
};
