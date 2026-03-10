import api from './api';

export default {
	listar(usuarioId) {
		return api.get(`/notificacoes?usuario_id=${usuarioId}`);
	},

	marcarComoLida(id) {
		return api.put(`/notificacoes/ler?id=${id}`);
	},

	criar(data) {
		return api.post('/notificacoes', data);
	},

	limparTudo(usuarioId) {
		return api.delete(`/notificacoes?usuario_id=${usuarioId}`);
	}
};
