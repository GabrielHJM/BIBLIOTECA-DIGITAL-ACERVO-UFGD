<template>
	<div class="login-absolute-center">
		<v-card class="ios-login-card" elevation="0">
			<!-- Header / Back Button -->
			<div class="card-header-actions mb-4">
				<v-btn icon="mdi-chevron-left" variant="text" @click="$router.push('/')" size="small" class="back-btn"></v-btn>
			</div>

				<div class="text-center mb-8">
					<div class="login-icon-wrapper mb-4">
						<v-img :src="caminho" width="80" class="mx-auto drop-shadow"></v-img>
					</div>
					<h2 class="login-title">Boas-vindas</h2>
					<p class="login-subtitle">Acesse sua biblioteca pessoal</p>
				</div>

				<div class="login-form">
					<div class="input-group mb-4">
						<label class="ios-label">E-MAIL</label>
						<input v-model="email" type="email" placeholder="exemplo@email.com" class="ios-input-modern" />
					</div>

					<div class="input-group mb-2">
						<label class="ios-label">SENHA</label>
						<input v-model="senha" type="password" placeholder="••••••••" class="ios-input-modern" />
					</div>

					<div class="text-right mb-8">
						<router-link to="/esqueci-senha" class="forgot-link">Esqueceu a senha?</router-link>
					</div>

					<v-btn class="ios-primary-btn w-100 mb-6" @click="login" :loading="loading" height="52" elevation="0">
						Entrar
					</v-btn>

					<div class="text-center">
						<router-link to="/cadastro" class="signup-link">
							Não tem conta? <span class="accent-text">Criar agora</span>
						</router-link>
					</div>
				</div>

				<v-snackbar v-model="snackbar" :color="snackbarColor" location="top" rounded="pill">
					{{ snackbarText }}
				</v-snackbar>
			</v-card>
	</div>
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
	.login-absolute-center {
		width: 100%;
		display: flex;
		justify-content: center;
		margin: auto;
	}

	.ios-login-card {
		width: 100%;
		max-width: 400px;
		background: rgba(var(--v-theme-surface), 0.8) !important;
		backdrop-filter: blur(20px) saturate(180%);
		-webkit-backdrop-filter: blur(20px) saturate(180%);
		border-radius: 32px !important;
		border: 1px solid rgba(var(--v-border-color), 0.1);
		padding: 32px !important;
		box-shadow: 0 20px 50px rgba(0, 0, 0, 0.15) !important;
	}

	.login-title {
		font-size: 28px;
		font-weight: 800;
		letter-spacing: -1px;
		margin-bottom: 4px;
	}

	.login-subtitle {
		font-size: 15px;
		opacity: 0.6;
		font-weight: 500;
	}

	.ios-label {
		display: block;
		font-size: 11px;
		font-weight: 800;
		letter-spacing: 1px;
		margin-bottom: 8px;
		margin-left: 4px;
		opacity: 0.5;
	}

	.ios-input-modern {
		width: 100%;
		background: rgba(var(--v-theme-on-surface), 0.05);
		border: 1px solid transparent;
		border-radius: 14px;
		padding: 14px 18px;
		font-size: 16px;
		transition: all 0.3s var(--spring-easing);
		color: rgb(var(--v-theme-on-surface));
	}

	.ios-input-modern:focus {
		outline: none;
		background: rgba(var(--v-theme-surface), 1);
		border-color: var(--ios-cyan);
		box-shadow: 0 0 0 4px rgba(0, 184, 212, 0.1);
	}

	.forgot-link {
		font-size: 13px;
		color: var(--ios-cyan);
		text-decoration: none;
		font-weight: 600;
	}

	.ios-primary-btn {
		background: var(--ios-cyan) !important;
		color: white !important;
		border-radius: 14px !important;
		text-transform: none !important;
		font-size: 16px !important;
		font-weight: 700 !important;
		letter-spacing: -0.2px;
		transition: transform 0.2s var(--spring-easing) !important;
	}

	.ios-primary-btn:active {
		transform: scale(0.98);
	}

	.signup-link {
		font-size: 14px;
		text-decoration: none;
		color: rgba(var(--v-theme-on-surface), 0.6);
		font-weight: 500;
	}

	.accent-text {
		color: var(--ios-cyan);
		font-weight: 700;
	}

	.drop-shadow {
		filter: drop-shadow(0 12px 24px rgba(0, 184, 212, 0.25));
	}

	.back-btn {
		margin-left: -12px;
		opacity: 0.5;
		transition: opacity 0.2s;
	}

	.back-btn:hover {
		opacity: 1;
	}
</style>
