import api from './api';

export default {
	// Cadastro de usuário
	cadastrar(dados) {
		return api.post('/usuarios', dados);
	},

	// Login
	login(dados) {
		return api.post('/login', dados);
	},

	// Redefinir senha
	redefinirSenha(dados) {
		return api.post('/redefinir-senha', dados);
	}
};
