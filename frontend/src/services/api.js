import axios from 'axios';

const api = axios.create({
	baseURL: 'http://localhost:8080', // Ajuste a porta conforme .env do backend
	headers: {
		'Content-Type': 'application/json',
	},
});

// Interceptor para adicionar token se existir (para futuro)
api.interceptors.request.use(
	(config) => {
		const userStr = localStorage.getItem('user');
		if (userStr) {
			const user = JSON.parse(userStr);
			if (user.token) {
				config.headers.Authorization = `Bearer ${user.token}`;
			}
		}
		return config;
	},
	(error) => {
		return Promise.reject(error);
	}
);

// Interceptor de resposta para padronizar o retorno e o erro
api.interceptors.response.use(
	(response) => {
		const res = response.data;
		// Se o backend retornou success: true, retornamos apenas o data
		if (res && res.success) {
			return { ...response, data: res.data };
		}
		// Caso contrário, tratamos como erro
		return Promise.reject(res?.error || 'Erro desconhecido no servidor');
	},
	(error) => {
		const message = error.response?.data?.error || error.message || 'Erro na comunicação com a API';
		return Promise.reject(message);
	}
);

export default api;
