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
	},

	// Atualizar perfil
	atualizar(id, dados) {
		return api.put(`/usuarios/${id}`, dados);
	},

	// Atualizar meta de leitura
	atualizarMeta(id, meta) {
		return api.put(`/usuarios/${id}/meta`, { meta_paginas_semana: meta });
	}
};
