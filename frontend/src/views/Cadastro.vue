<template>
	<v-container class="d-flex justify-center align-center min-vh-100">
		<v-card class="ios-card-premium" elevation="12">
			<div class="d-flex justify-start mb-2">
				<v-btn icon variant="text" color="white" @click="$router.push('/')" size="small">
					<v-icon>mdi-arrow-left</v-icon>
				</v-btn>
			</div>

			<v-row justify="center" class="mb-4">
				<v-img
					:src="caminho"
					alt="Avatar"
					max-width="100"
					class="drop-shadow"
					contain
				></v-img>
			</v-row>

			<v-row>
				<v-col cols="12" class="text-center mb-6">
					<h1 class="text-white text-h4 font-weight-bold">CADASTRE-SE</h1>
				</v-col>
			</v-row>

			<v-row>
				<v-col cols="12" md="6" class="py-1">
					<label class="ios-label">Nome:</label>
					<input v-model="nome" type="text" placeholder="Seu nome completo" class="ios-input-field w-100" />
				</v-col>
				<v-col cols="12" md="6" class="py-1">
					<label class="ios-label">Instituição:</label>
					<input v-model="instituicao" type="text" placeholder="Escola ou Faculdade" class="ios-input-field w-100" />
				</v-col>
			</v-row>

			<v-row class="mt-4">
				<v-col cols="12" md="6" class="py-1">
					<label class="ios-label">Email:</label>
					<input v-model="email" type="email" placeholder="seu@email.com" class="ios-input-field w-100" />
				</v-col>
				<v-col cols="12" md="6" class="py-1">
					<label class="ios-label">Senha:</label>
					<input v-model="senha" type="password" placeholder="Mínimo 6 caracteres" class="ios-input-field w-100" />
				</v-col>
			</v-row>

			<v-row class="mt-4 align-center">
				<v-col cols="12" md="6" class="py-1">
					<label class="ios-label">Eu sou:</label>
					<v-select
						class="ios-select-field"
						:items="['Aluno', 'Professor', 'Curador Especializado', 'Suporte']"
						v-model="tipoUsuario"
						variant="solo"
						flat
						hide-details
						density="comfortable"
						bg-color="white"
						placeholder="Selecione seu perfil"
					></v-select>
				</v-col>
				<v-col cols="12" md="6" class="text-center mt-6">
					<v-btn size="x-large" class="ios-action-btn w-100" elevation="4" @click="cadastrar" :loading="loading" :disabled="loading">
						Cadastrar-me
					</v-btn>
				</v-col>
			</v-row>

			<v-row class="mt-8">
				<v-col cols="12" class="text-center">
					<router-link to="/login" class="ios-secondary-link">
						<v-icon color="white" size="24" class="mr-2">mdi-arrow-left</v-icon>
						<span class="back-text">Já tenho cadastro? Ir para o Login</span>
					</router-link>
				</v-col>
			</v-row>

			<v-snackbar v-model="snackbar" :color="snackbarColor" location="top" rounded="pill">
				{{ snackbarText }}
			</v-snackbar>
		</v-card>
	</v-container>
</template>

<script>
import UsuarioService from '@/services/UsuarioService'

export default {
	name: 'CadastroPage',
	data: () => ({
		caminho: require("@/assets/images/site-images/login/person-img.png"),
		tipoUsuario: null,
		nome: '',
		instituicao: '',
		email: '',
		senha: '',
		loading: false,
		snackbar: false,
		snackbarText: '',
		snackbarColor: 'success'
	}),
	methods: {
		async cadastrar() {
			if (!this.nome || !this.email || !this.senha || !this.tipoUsuario) {
				this.snackbarText = "Por favor, preencha todos os campos obrigatórios."
				this.snackbarColor = "warning"
				this.snackbar = true
				return
			}

			// Validação simples de email no frontend
			const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
			if (!emailRegex.test(this.email)) {
				this.snackbarText = "Por favor, insira um email válido."
				this.snackbarColor = "warning"
				this.snackbar = true
				return
			}

			if (this.senha.length < 6) {
				this.snackbarText = "A senha deve ter pelo menos 6 caracteres."
				this.snackbarColor = "warning"
				this.snackbar = true
				return
			}

			this.loading = true
			try {
				let tipo = 0
				switch(this.tipoUsuario) {
					case 'Aluno': tipo = 1; break;
					case 'Professor': tipo = 2; break;
					case 'Curador Especializado': tipo = 3; break;
					case 'Suporte': tipo = 4; break;
					default: tipo = 1;
				}

				await UsuarioService.cadastrar({
					nome: this.nome,
					email: this.email,
					senha: this.senha,
					tipo: tipo
				})

				this.snackbarText = "Cadastro realizado com sucesso! Bem-vindo(a)!"
				this.snackbarColor = "success"
				this.snackbar = true

				setTimeout(() => {
					this.$router.push('/login')
				}, 1500)

			} catch (error) {
				console.error(error)
				// Agora o interceptor retorna o erro como string
				this.snackbarText = typeof error === 'string' ? error : "Erro ao realizar cadastro. Tente outro email."
				this.snackbarColor = "error"
				this.snackbar = true
			} finally {
				this.loading = false
			}
		}
	}
}
</script>

<style scoped>
	.min-vh-100 {
		min-height: 100vh;
	}

	.ios-card-premium {
		width: 95%;
		max-width: 700px;
		background: rgba(45, 78, 115, 0.85) !important;
		backdrop-filter: blur(20px);
		border-radius: 24px !important;
		border: 1px solid rgba(255, 255, 255, 0.1);
		transition: transform 0.3s var(--spring-easing);
		padding: 24px !important;
	}

	.ios-card-premium:hover {
		transform: translateY(-5px);
	}

	.ios-label {
		color: rgba(255, 255, 255, 0.9);
		font-weight: 600;
		font-size: 14px;
		margin-left: 12px;
		display: block;
		margin-bottom: 4px;
	}

	.ios-input-field {
		background: rgba(255, 255, 255, 0.9) !important;
		border: none !important;
		border-radius: 16px;
		padding: 12px 20px;
		font-size: 15px;
		color: #333;
		transition: all 0.2s ease;
	}

	.ios-input-field:focus {
		outline: none;
		background: white !important;
		box-shadow: 0 0 0 3px rgba(0, 184, 212, 0.4);
	}

	.ios-select-field :deep(.v-field) {
		border-radius: 16px !important;
		background: rgba(255, 255, 255, 0.9) !important;
	}

	.ios-action-btn {
		background: linear-gradient(135deg, #00B8D4 0%, #0097A7 100%) !important;
		color: white !important;
		border-radius: 16px !important;
		height: 54px !important;
		text-transform: none !important;
		font-size: 18px !important;
		font-weight: 700 !important;
		letter-spacing: 0.5px !important;
	}

	.ios-secondary-link {
		color: white;
		text-decoration: none;
		font-weight: 500;
		display: inline-flex;
		align-items: center;
		opacity: 0.8;
		transition: all 0.2s;
	}

	.ios-secondary-link:hover {
		opacity: 1;
		transform: scale(1.05);
	}

	.back-text {
		font-size: 16px;
	}

	.drop-shadow {
		filter: drop-shadow(0 10px 15px rgba(0,0,0,0.3));
	}
</style>
