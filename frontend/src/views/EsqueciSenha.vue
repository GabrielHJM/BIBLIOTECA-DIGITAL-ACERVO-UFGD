<template>
	<v-container class="d-flex justify-center align-center min-vh-100">
		<v-card class="ios-card-premium pa-8" elevation="12">
			<v-row justify="center" class="mb-6">
				<v-img
					:src="caminho"
					alt="Redefinir"
					max-width="120"
					class="drop-shadow"
					contain
				></v-img>
			</v-row>

            <v-row justify="center" class="mb-6">
                <h2 class="text-white font-weight-bold">Redefinir Senha</h2>
            </v-row>

			<v-row>
				<v-col cols="12" class="text-left py-1 pb-2"><label class="ios-label">Email:</label></v-col>
				<v-col cols="12" class="py-0">
					<input v-model="email" type="text" placeholder="seu@email.com" class="ios-input-field w-100" />
				</v-col>
			</v-row>

			<v-row class="mt-4">
				<v-col cols="12" class="text-left py-1 pb-2"><label class="ios-label">Nova Senha:</label></v-col>
				<v-col cols="12" class="py-0">
					<input v-model="senha" type="password" placeholder="Mínimo 6 caracteres" class="ios-input-field w-100" />
				</v-col>
			</v-row>

            <v-row class="mt-4">
				<v-col cols="12" class="text-left py-1 pb-2"><label class="ios-label">Confirmar Nova Senha:</label></v-col>
				<v-col cols="12" class="py-0">
					<input v-model="confirmarSenha" type="password" placeholder="Repita a nova senha" class="ios-input-field w-100" />
				</v-col>
			</v-row>

			<v-row justify="center" class="mt-8">
				<v-btn class="ios-action-btn w-100" elevation="4" @click="redefinirSenha" :loading="loading" :disabled="loading">
					Redefinir
				</v-btn>
			</v-row>

			<v-row class="mt-10">
				<v-col cols="12" class="text-center">
					<router-link to="/login" class="ios-secondary-link">
						<v-icon color="white" size="24" class="mr-2">mdi-chevron-left</v-icon>
						<span class="back-text">Voltar para o Login</span>
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
	name: 'EsqueciSenhaPage',
	data: () => ({
		caminho: require("@/assets/images/site-images/login/person-img.png"),
		loading: false,
		email: '',
		senha: '',
        confirmarSenha: '',
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
				this.snackbarText = error.response?.data || "Falha ao redefinir senha. Verifique o email."
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
		width: 100%;
		max-width: 420px;
		background: rgba(58, 99, 145, 0.8) !important;
		backdrop-filter: blur(20px);
		border-radius: 32px !important;
		border: 1px solid rgba(255, 255, 255, 0.1);
        transition: transform 0.3s var(--spring-easing);
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
