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
					alt="Login"
					max-width="100"
					class="drop-shadow"
					contain
				></v-img>
			</v-row>

			<v-row>
				<v-col cols="12" class="text-center mb-2">
					<h2 class="text-white text-h5 font-weight-bold">BEM-VINDO</h2>
				</v-col>
			</v-row>

			<v-row>
				<v-col cols="12" class="text-left py-1 pb-2"><label class="ios-label">Email:</label></v-col>
				<v-col cols="12" class="py-0">
					<input v-model="email" type="text" placeholder="seu@email.com" class="ios-input-field w-100" />
				</v-col>
			</v-row>

			<v-row class="mt-4">
				<v-col cols="12" class="text-left py-1 pb-2"><label class="ios-label">Senha:</label></v-col>
				<v-col cols="12" class="py-0">
					<input v-model="senha" type="password" placeholder="Sua senha" class="ios-input-field w-100" />
				</v-col>
				<v-col cols="12" class="text-right mt-1">
					<router-link to="/esqueci-senha" class="ios-link-small">Esqueci minha senha</router-link>
				</v-col>
			</v-row>

			<v-row justify="center" class="mt-8">
				<v-btn class="ios-action-btn w-100" elevation="4" @click="login" :loading="loading" :disabled="loading">
					Entrar
				</v-btn>
			</v-row>

			<v-row class="mt-10">
				<v-col cols="12" class="text-center">
					<router-link to="/cadastro" class="ios-secondary-link">
						<v-icon size="20" class="mr-2">mdi-account-plus-outline</v-icon>
						<span>Não tenho cadastro</span>
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
import auth from '@/auth'
import UsuarioService from '@/services/UsuarioService'

export default {
	name: 'LoginPage',
	data: () => ({
		caminho: require("@/assets/images/site-images/login/person-img.png"),
		loading: false,
		email: '',
		senha: '',
		snackbar: false,
		snackbarText: '',
		snackbarColor: 'error'
	}),
	methods: {
		async login() {
			if (!this.email || !this.senha) {
				this.snackbarText = "Preencha todos os campos."
				this.snackbarColor = "warning"
				this.snackbar = true
				return
			}

			this.loading = true
			try {
				const response = await UsuarioService.login({
					email: this.email,
					senha: this.senha
				})

				if (response.data && response.data.token) {
					auth.login(response.data)
					this.$router.push('/dashboard')
				} else {
					throw new Error("Token não recebido")
				}
			} catch (error) {
				console.error(error)
				this.snackbarText = typeof error === 'string' ? error : "Falha no login. Verifique suas credenciais."
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
		width: 90%;
		max-width: 380px;
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
	}

	.ios-input-field {
		background: rgba(255, 255, 255, 0.9) !important;
		border: none !important;
		border-radius: 16px;
		padding: 14px 20px;
		font-size: 16px;
		color: #333;
		transition: all 0.2s ease;
	}

	.ios-input-field:focus {
		outline: none;
		background: white !important;
		box-shadow: 0 0 0 3px rgba(0, 184, 212, 0.4);
	}

	.ios-link-small {
		color: rgba(255, 255, 255, 0.7);
		font-size: 13px;
		text-decoration: none;
		transition: color 0.2s;
	}

	.ios-link-small:hover {
		color: #00B8D4;
		text-decoration: underline;
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

	.drop-shadow {
		filter: drop-shadow(0 10px 15px rgba(0,0,0,0.3));
	}
</style>
