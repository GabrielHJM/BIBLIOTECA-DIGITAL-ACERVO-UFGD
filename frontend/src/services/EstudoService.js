import api from './api';

export default {


	// Flashcards
	listarFlashcards(usuarioId, materialId = 0) {
		let url = `/estudo/flashcards?usuario_id=${usuarioId}`;
		if (materialId > 0) url += `&material_id=${materialId}`;
		return api.get(url);
	},



	revisarFlashcard(id, dificuldade) {
		return api.put('/estudo/flashcards/revisar', { id, dificuldade });
	}
};
