import api from './api';

export default {
	// Materiais Pendentes
	listarPendentes() {
		return api.get('/admin/materiais/pendentes');
	},

	// Moderar (Aprovar/Rejeitar)
	moderarMaterial(id, status, curadorId) {
		return api.post('/admin/materiais/moderar', {
			id: id,
			status: status,
			curador_id: curadorId
		});
	},

	// Métricas Globais
	obterMetricas() {
		return api.get('/admin/metricas');
	},

	// Deletar material
	deletarMaterial(id) {
		return api.delete(`/admin/materiais/${id}`);
	}
};
