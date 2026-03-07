<template>
	<div class="login-absolute-center">
		<!-- Ambient Background Blobs for Glassmorphism -->
		<div class="ambient-blob blob-1"></div>
		<div class="ambient-blob blob-2"></div>

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
						</v-col>
					</v-row>

					<v-row dense class="mb-2">
						<v-col cols="12" class="px-2 mb-2">
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
		mostrarSenha: false,
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
		background: radial-gradient(circle at top left, rgba(0, 184, 212, 0.08), transparent 40%),
					radial-gradient(circle at bottom right, rgba(0, 184, 212, 0.05), transparent 40%);
		overflow: hidden;
		height: 100vh;
		width: 100vw;
	}

	.ambient-blob {
		position: absolute;
		border-radius: 50%;
		filter: blur(80px);
		z-index: 0;
		opacity: 0.6;
	}
	.blob-1 {
		width: 300px; height: 300px;
		background: rgba(0, 184, 212, 0.4);
		top: -50px; left: -100px;
	}
	.blob-2 {
		width: 400px; height: 400px;
		background: rgba(144, 202, 249, 0.2);
		bottom: -100px; right: -100px;
	}

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
		background: rgba(255, 255, 255, 0.03) !important;
		backdrop-filter: blur(40px) saturate(180%);
		-webkit-backdrop-filter: blur(40px) saturate(180%);
		border-radius: clamp(24px, 4vw, 32px) !important;
		border: 1px solid rgba(255, 255, 255, 0.1);
		padding: clamp(16px, 4vw, 32px) clamp(16px, 4vw, 30px) !important;
		box-shadow: 0 32px 64px rgba(0, 0, 0, 0.4) !important;
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
		color: #ffffff;
	}

	.login-subtitle {
		font-size: 16px;
		opacity: 0.7;
		font-weight: 500;
		color: #ffffff;
	}

	.ios-label {
		display: block;
		font-size: 11px;
		font-weight: 700;
		letter-spacing: 1px;
		margin-bottom: 8px;
		margin-left: 8px;
		opacity: 0.6;
		color: #ffffff;
	}

	.ios-input-modern {
		width: 100%;
		box-sizing: border-box;
		background: rgba(0, 0, 0, 0.25);
		border: 1px solid rgba(255, 255, 255, 0.08);
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
		background: rgba(0, 0, 0, 0.4);
		border-color: rgba(0, 184, 212, 0.6);
		box-shadow: 0 0 0 4px rgba(0, 184, 212, 0.15);
	}

	.ios-select-modern :deep(.v-field) {
		background: rgba(0, 0, 0, 0.25) !important;
		border: 1px solid rgba(255, 255, 255, 0.08) !important;
		border-radius: 18px !important;
		color: #ffffff !important;
		transition: all 0.3s cubic-bezier(0.16, 1, 0.3, 1);
		box-shadow: none !important;
		padding: 6px 0;
	}
	.ios-select-modern :deep(.v-field--focused) {
		background: rgba(0, 0, 0, 0.4) !important;
		border-color: rgba(0, 184, 212, 0.6) !important;
		box-shadow: 0 0 0 4px rgba(0, 184, 212, 0.15) !important;
	}
	.ios-select-modern :deep(.v-select__selection-text) {
		color: #ffffff !important;
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
		color: #00B8D4 !important;
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
		background: linear-gradient(135deg, #00B8D4 0%, #007A99 100%) !important;
		color: white !important;
		border-radius: 20px !important;
		text-transform: none !important;
		font-size: 18px !important;
		font-weight: 700 !important;
		letter-spacing: -0.2px;
		transition: transform 0.2s cubic-bezier(0.16, 1, 0.3, 1), box-shadow 0.2s !important;
		box-shadow: 0 12px 24px rgba(0, 184, 212, 0.3) !important;
	}

	.ios-primary-btn:active {
		transform: scale(0.97);
	}

	.signup-link {
		font-size: 15px;
		text-decoration: none;
		color: rgba(255, 255, 255, 0.6);
		font-weight: 500;
	}

	.accent-text {
		color: #00B8D4;
		font-weight: 700;
	}

	.drop-shadow {
		filter: drop-shadow(0 16px 24px rgba(0, 0, 0, 0.4));
	}

	.back-btn {
		opacity: 0.6;
		transition: opacity 0.2s;
		color: white;
	}

	.back-btn:hover {
		opacity: 1;
	}
</style>
