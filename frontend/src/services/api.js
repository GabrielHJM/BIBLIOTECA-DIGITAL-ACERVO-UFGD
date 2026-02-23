import axios from 'axios';

const api = axios.create({
	baseURL: 'http://localhost:8081', // Ajuste a porta conforme .env do backend
	headers: {
		'Content-Type': 'application/json',
	},
});

// Interceptor para adicionar token se existir (para futuro)
api.interceptors.request.use(
	(config) => {
		const token = localStorage.getItem('token');
		if (token) {
			config.headers.Authorization = `Bearer ${token}`;
		}
		return config;
	},
	(error) => {
		return Promise.reject(error);
	}
);

export default api;
