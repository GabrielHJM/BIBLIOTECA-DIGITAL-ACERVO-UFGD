import api from './api';

class AnotacaoService {
	async criar(anotacao) {
		const response = await api.post('/anotacoes', anotacao);
		return response.data;
	}

	async listar() {
		const response = await api.get('/anotacoes');
		return response.data;
	}

	async atualizar(id, anotacao) {
		const response = await api.put(`/anotacoes/${id}`, anotacao);
		return response.data;
	}

	async excluir(id) {
		const response = await api.delete(`/anotacoes/${id}`);
		return response.data;
	}
}

export default new AnotacaoService();
