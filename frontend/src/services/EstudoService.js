import api from './api';

export default {
	// Anotações
	listarAnotacoes(usuarioId, materialId) {
		return api.get(`/estudo/anotacoes?usuario_id=${usuarioId}&material_id=${materialId}`);
	},

	criarAnotacao(anotacao) {
		return api.post('/estudo/anotacoes', anotacao);
	},

	deletarAnotacao(id, usuarioId) {
		return api.delete(`/estudo/anotacoes?id=${id}&usuario_id=${usuarioId}`);
	},

	// Flashcards
	listarFlashcards(usuarioId, materialId = 0) {
		let url = `/estudo/flashcards?usuario_id=${usuarioId}`;
		if (materialId > 0) url += `&material_id=${materialId}`;
		return api.get(url);
	},

	gerarFlashcardsIA(usuarioId, materialId) {
		return api.post('/estudo/flashcards/gerar', {
			usuario_id: usuarioId,
			material_id: materialId
		});
	},

	revisarFlashcard(id, dificuldade) {
		return api.put('/estudo/flashcards/revisar', { id, dificuldade });
	}
};
