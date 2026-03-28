<template>
	<div class="login-absolute-center">
		<!-- Ambient mesh gradient is now handled globally in App.vue -->

		<div class="content-wrapper">
			<v-card class="ios-login-card position-relative" elevation="0">
				<!-- Header / Back Button -->
				<div class="card-header-actions">
					<v-btn icon="mdi-chevron-left" variant="text" @click="$router.push('/')" size="large" class="back-btn"></v-btn>
				</div>

				<div class="text-center mt-2 mb-4">
					<div class="login-icon-wrapper mb-4">
						<v-img :src="caminho" width="64" class="mx-auto drop-shadow"></v-img>
					</div>
					<h2 class="login-title">Cadastre-se</h2>
					<p class="login-subtitle">Crie sua biblioteca pessoal</p>
				</div>

				<div class="login-form">
					<v-row dense>
						<v-col cols="12" sm="6" class="px-2 mb-1">
							<label class="ios-label">NOME</label>
							<input v-model="nome" type="text" placeholder="Seu nome completo" class="ios-input-modern" />
						</v-col>
						<v-col cols="12" sm="6" class="px-2 mb-1">
							<label class="ios-label">INSTITUIÇÃO</label>
							<input v-model="instituicao" type="text" placeholder="Escola ou Faculdade" class="ios-input-modern" />
						</v-col>
					</v-row>

					<v-row dense>
						<v-col cols="12" sm="6" class="px-2 mb-1">
							<label class="ios-label">E-MAIL</label>
							<input v-model="email" type="email" placeholder="seu@email.com" class="ios-input-modern" />
						</v-col>
						<v-col cols="12" sm="6" class="px-2 mb-1">
							<label class="ios-label">SENHA</label>
							<div class="password-wrapper">
								<input
									v-model="senha"
									:type="mostrarSenha ? 'text' : 'password'"
									placeholder="Mínimo 6 caracteres"
									class="ios-input-modern pr-12"
								/>
								<v-btn
									icon
									variant="text"
									class="password-toggle-btn"
									@click="mostrarSenha = !mostrarSenha"
									tabindex="-1"
								>
									<v-icon color="white" size="20">
										{{ mostrarSenha ? 'mdi-eye-off' : 'mdi-eye' }}
									</v-icon>
								</v-btn>
							</div>
							<!-- Password Strength Indicator -->
							<div v-if="senha" class="mt-2 px-2">
								<div class="strength-bar-container">
									<div :class="['strength-bar', passwordStrength.colorClass]" :style="{ width: (passwordStrength.score / 5) * 100 + '%' }"></div>
								</div>
								<div class="d-flex justify-space-between mt-1">
									<span class="text-caption text-white opacity-60">Segurança da senha</span>
									<span :class="['text-caption font-weight-bold', passwordStrength.textClass]">{{ passwordStrength.label }}</span>
								</div>
							</div>
						</v-col>
					</v-row>

					<v-row dense class="mb-2">
						<v-col cols="12" sm="6" class="px-2 mb-1">
							<label class="ios-label">CONFIRMAR SENHA</label>
							<div class="password-wrapper">
								<input
									v-model="confirmarSenha"
									:type="mostrarSenha ? 'text' : 'password'"
									placeholder="Repita sua senha"
									class="ios-input-modern pr-12"
								/>
							</div>
						</v-col>
						<v-col cols="12" sm="6" class="px-2 mb-2">
							<label class="ios-label">PERFIL</label>
							<v-select
								class="ios-select-modern"
								:items="['Aluno', 'Professor']"
								v-model="tipoUsuario"
								variant="solo"
								flat
								hide-details
								placeholder="Selecione seu perfil"
							></v-select>
						</v-col>
					</v-row>

					<v-btn class="ios-primary-btn w-100 mb-6 mt-1" @click="cadastrar" :loading="loading" height="56" elevation="0">
						Criar Conta
					</v-btn>

					<div class="text-center mb-1">
						<router-link to="/login" class="signup-link">
							Já tem conta? <span class="accent-text">Fazer login</span>
						</router-link>
					</div>
				</div>

				<v-snackbar v-model="snackbar" :color="snackbarColor" location="top" rounded="pill">
					{{ snackbarText }}
				</v-snackbar>
			</v-card>
		</div>
	</div>
</template>

<script>
import auth from '@/auth'
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
		confirmarSenha: '',
		mostrarSenha: false,
		loading: false,
		snackbar: false,
		snackbarText: '',
		snackbarColor: 'success'
	}),
	computed: {
		passwordStrength() {
			const s = this.senha;
			if (!s) return { score: 0, label: '', colorClass: '', textClass: '' };
			
			let score = 0;
			if (s.length >= 8) score++;
			if (/[A-Z]/.test(s)) score++;
			if (/[a-z]/.test(s)) score++;
			if (/[0-9]/.test(s)) score++;
			if (/[^A-Za-z0-9]/.test(s)) score++;

			if (score <= 2) return { score, label: 'Fraca', colorClass: 'bg-error', textClass: 'text-error' };
			if (score <= 4) return { score, label: 'Forte', colorClass: 'bg-warning', textClass: 'text-warning' };
			return { score, label: 'Complexa', colorClass: 'bg-success', textClass: 'text-success' };
		}
	},
	methods: {
		async cadastrar() {
			if (!this.nome || !this.email || !this.senha || !this.tipoUsuario || !this.confirmarSenha) {
				this.snackbarText = "Por favor, preencha todos os campos obrigatórios."
				this.snackbarColor = "warning"
				this.snackbar = true
				return
			}

			if (this.senha !== this.confirmarSenha) {
				this.snackbarText = "As senhas não coincidem."
				this.snackbarColor = "error"
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
					default: tipo = 1;
				}

				const response = await UsuarioService.cadastrar({
					nome: this.nome,
					email: this.email,
					senha: this.senha,
					tipo: tipo
				})

				this.snackbarText = "Cadastro realizado com sucesso! Bem-vindo(a)!"
				this.snackbarColor = "success"
				this.snackbar = true

				// Auto-login se o backend retornar token
				if (response.data && response.data.token) {
					auth.login(response.data)
					this.snackbarText = "Cadastro realizado com sucesso! Redirecionando..."

					setTimeout(() => {
						this.$router.push('/dashboard').catch(err => {
							console.error("Erro no redirecionamento:", err)
							window.location.href = '/dashboard'
						})
					}, 1000)
				} else {
					this.snackbarText = "Cadastro realizado! Por favor, faça login."
					setTimeout(() => {
						this.$router.push('/login')
					}, 1500)
				}

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
	.login-absolute-center {
		position: absolute;
		top: 0; left: 0; right: 0; bottom: 0;
		display: flex;
		justify-content: center;
		align-items: center;
		background: transparent;
		overflow: hidden;
		height: 100vh;
		width: 100vw;
	}

	/* Blobs removed to avoid background clutter over the new vibrant mesh */

	.content-wrapper {
		width: 100%;
		display: flex;
		justify-content: center;
		align-items: center;
		z-index: 1;
		padding: clamp(24px, 4vw, 64px) clamp(16px, 4vw, 64px);
		height: 100%;
		box-sizing: border-box;
	}

	.ios-login-card {
		width: 100%;
		max-width: 580px;
		background: var(--glass-bg) !important;
		backdrop-filter: var(--glass-blur);
		-webkit-backdrop-filter: var(--glass-blur);
		border-radius: clamp(24px, 4vw, 32px) !important;
		border: 1px solid var(--glass-border);
		padding: clamp(16px, 4vw, 32px) clamp(16px, 4vw, 30px) !important;
		box-shadow: 0 32px 64px rgba(0, 0, 0, 0.2) !important;
		overflow-y: auto;
		max-height: calc(100vh - 32px);
	}

	/* Mobile Specific Adjustments */
	@media (max-width: 600px) {
		.login-absolute-center {
			padding: 12px;
			align-items: flex-start;
			overflow-y: auto;
		}
		.content-wrapper {
			padding: 0;
			height: auto;
		}
		.ios-login-card {
			margin: 16px 0;
			border-radius: 20px !important;
			max-height: none;
			padding: 24px 20px !important;
		}
		.login-title {
			font-size: 26px;
		}
		.ios-input-modern {
			padding: 12px 16px;
		}
	}

	.card-header-actions {
		position: absolute;
		top: 24px;
		left: 24px;
	}

	.login-title {
		font-size: 32px;
		font-weight: 800;
		letter-spacing: -1px;
		margin-bottom: 8px;
		color: var(--v-theme-on-surface);
	}

	.login-subtitle {
		font-size: 16px;
		opacity: 0.7;
		font-weight: 500;
		color: var(--v-theme-on-surface);
	}

	.ios-label {
		display: block;
		font-size: 11px;
		font-weight: 700;
		letter-spacing: 1px;
		margin-bottom: 8px;
		margin-left: 8px;
		opacity: 0.6;
		color: var(--v-theme-on-surface);
	}

	.ios-input-modern {
		width: 100%;
		box-sizing: border-box;
		background: rgba(255, 255, 255, 0.05);
		border: 1px solid rgba(255, 255, 255, 0.1);
		border-radius: 14px;
		padding: 14px 20px;
		font-size: 15px;
		transition: all 0.3s cubic-bezier(0.16, 1, 0.3, 1);
		color: #ffffff;
	}

	.ios-input-modern::placeholder {
		color: rgba(255,255,255,0.4);
	}

	.ios-input-modern:focus {
		outline: none;
		background: rgba(255, 255, 255, 0.1);
		border-color: var(--ios-blue);
		box-shadow: 0 0 0 4px rgba(var(--ios-blue-rgb), 0.2);
	}

	.ios-select-modern :deep(.v-field) {
		background: var(--glass-bg) !important;
		border: 1px solid var(--glass-border) !important;
		border-radius: 18px !important;
		color: var(--v-theme-on-surface) !important;
		transition: all 0.3s cubic-bezier(0.16, 1, 0.3, 1);
		box-shadow: none !important;
		padding: 6px 0;
	}
	.ios-select-modern :deep(.v-field--focused) {
		background: rgba(0, 0, 0, 0.4) !important;
		border-color: rgba(0, 122, 255, 0.6) !important;
		box-shadow: 0 0 0 4px rgba(0, 122, 255, 0.15) !important;
	}
	.ios-select-modern :deep(.v-select__selection-text) {
		color: var(--v-theme-on-surface) !important;
	}

	.password-wrapper {
		position: relative;
		display: flex;
		align-items: center;
	}

	.password-toggle-btn {
		position: absolute;
		right: 6px;
		opacity: 0.7;
		color: #007AFF !important;
		transition: all 0.2s ease;
	}
	.password-toggle-btn:hover {
		opacity: 1;
		transform: scale(1.1);
	}

	.pr-12 {
		padding-right: 48px !important;
	}

	.ios-primary-btn {
		background: linear-gradient(135deg, #007AFF 0%, #0056B3 100%) !important;
		color: white !important;
		border-radius: 20px !important;
		text-transform: none !important;
		font-size: 18px !important;
		font-weight: 700 !important;
		letter-spacing: -0.2px;
		transition: transform 0.2s cubic-bezier(0.16, 1, 0.3, 1), box-shadow 0.2s !important;
		box-shadow: 0 12px 24px rgba(0, 122, 255, 0.3) !important;
	}

	.ios-primary-btn:active {
		transform: scale(0.97);
	}

	.signup-link {
		font-size: 15px;
		text-decoration: none;
		color: var(--v-theme-on-surface);
		opacity: 0.6;
		font-weight: 500;
	}

	.accent-text {
		color: #007AFF;
		font-weight: 700;
	}

	.drop-shadow {
		filter: drop-shadow(0 16px 24px rgba(0, 0, 0, 0.4));
	}

	.back-btn {
		opacity: 0.6;
		transition: opacity 0.2s;
		color: var(--v-theme-on-surface);
	}

	.back-btn:hover {
		opacity: 1;
	}

	.strength-bar-container {
		width: 100%;
		height: 4px;
		background: rgba(255, 255, 255, 0.1);
		border-radius: 2px;
		overflow: hidden;
	}

	.strength-bar {
		height: 100%;
		transition: all 0.4s cubic-bezier(0.16, 1, 0.3, 1);
	}
</style>
