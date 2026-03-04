<template>
	<div class="login-absolute-center">
		<!-- Ambient Background Blobs for Glassmorphism -->
		<div class="ambient-blob blob-1"></div>
		<div class="ambient-blob blob-2"></div>

		<div class="content-wrapper">
			<v-card class="ios-login-card position-relative" elevation="0">
				<!-- Header / Back Button -->
				<div class="card-header-actions">
					<v-btn icon="mdi-chevron-left" variant="text" @click="$router.push('/login')" size="large" class="back-btn"></v-btn>
				</div>

				<div class="text-center mt-2 mb-6">
					<div class="login-icon-wrapper mb-4">
						<v-img :src="caminho" width="64" class="mx-auto drop-shadow"></v-img>
					</div>
					<h2 class="login-title">Recuperação</h2>
					<p class="login-subtitle">Redefina sua senha de acesso</p>
				</div>

				<div class="login-form">
					<div class="input-group mb-4">
						<label class="ios-label">E-MAIL</label>
						<input v-model="email" type="email" placeholder="seu@email.com" class="ios-input-modern" />
					</div>

					<div class="input-group mb-4">
						<label class="ios-label">NOVA SENHA</label>
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
					</div>

					<div class="input-group mb-8">
						<label class="ios-label">CONFIRMAR NOVA SENHA</label>
						<div class="password-wrapper">
							<input
								v-model="confirmarSenha"
								:type="mostrarConfirmarSenha ? 'text' : 'password'"
								placeholder="Repita a nova senha"
								class="ios-input-modern pr-12"
							/>
							<v-btn
								icon
								variant="text"
								class="password-toggle-btn"
								@click="mostrarConfirmarSenha = !mostrarConfirmarSenha"
								tabindex="-1"
							>
								<v-icon color="white" size="20">
									{{ mostrarConfirmarSenha ? 'mdi-eye-off' : 'mdi-eye' }}
								</v-icon>
							</v-btn>
						</div>
					</div>

					<v-btn class="ios-primary-btn w-100 mb-6" @click="redefinirSenha" :loading="loading" height="56" elevation="0" :disabled="loading">
						Redefinir Senha
					</v-btn>

					<div class="text-center">
						<router-link to="/login" class="signup-link">
							Lembrou da senha? <span class="accent-text">Fazer login</span>
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
import UsuarioService from '@/services/UsuarioService'

export default {
	name: 'EsqueciSenhaPage',
	data: () => ({
		caminho: require("@/assets/images/site-images/login/person-img.png"),
		loading: false,
		email: '',
		senha: '',
        confirmarSenha: '',
		mostrarSenha: false,
		mostrarConfirmarSenha: false,
		snackbar: false,
		snackbarText: '',
		snackbarColor: 'error'
	}),
	methods: {
		async redefinirSenha() {
            if (!this.email || !this.senha || !this.confirmarSenha) {
                this.snackbarText = "Preencha todos os campos."
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

			this.loading = true
			try {
				await UsuarioService.redefinirSenha({
					email: this.email,
					senha: this.senha
				})

                this.snackbarText = "Senha redefinida com sucesso!"
                this.snackbarColor = "success"
                this.snackbar = true

                setTimeout(() => {
                    this.$router.push('/login')
                }, 2000)
			} catch (error) {
				console.error(error)
				this.snackbarText = typeof error === 'string' ? error : "Falha ao redefinir senha. Verifique o email."
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
		max-width: 440px;
		background: rgba(255, 255, 255, 0.03) !important;
		backdrop-filter: blur(40px) saturate(180%);
		-webkit-backdrop-filter: blur(40px) saturate(180%);
		border-radius: clamp(24px, 4vw, 32px) !important;
		border: 1px solid rgba(255, 255, 255, 0.1);
		padding: clamp(24px, 5vw, 40px) clamp(20px, 5vw, 36px) !important;
		box-shadow: 0 32px 64px rgba(0, 0, 0, 0.4) !important;
		overflow-y: auto;
		max-height: calc(100vh - 48px);
	}

	.card-header-actions {
		position: absolute;
		top: 24px;
		left: 24px;
	}

	.login-title {
		font-size: 28px;
		font-weight: 800;
		letter-spacing: -1px;
		margin-bottom: 4px;
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
		border-radius: 16px;
		padding: 16px 20px;
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
