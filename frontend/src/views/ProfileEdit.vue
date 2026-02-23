<template>
	<v-container class="profile-container d-flex justify-center align-center">
		<v-card class="ios-profile-card" elevation="12">
			<!-- Header -->
			<div class="d-flex align-center mb-6">
				<v-btn icon variant="text" color="white" @click="$router.back()" size="small" class="mr-2">
					<v-icon>mdi-arrow-left</v-icon>
				</v-btn>
				<h2 class="text-white text-h5 font-weight-bold">Editar Perfil</h2>
			</div>

			<!-- Avatar Section -->
			<div class="avatar-section d-flex flex-column align-center mb-8">
				<div class="avatar-wrapper" @click="openPhotoSelector">
					<v-avatar size="120" class="profile-avatar-glass">
						<v-img v-if="user.foto_url || previewUrl" :src="previewUrl || user.foto_url" cover></v-img>
						<v-icon v-else size="64" color="white">mdi-account</v-icon>
						<div class="avatar-overlay">
							<v-icon color="white">mdi-camera</v-icon>
						</div>
					</v-avatar>
				</div>
				<span class="mt-3 text-caption text-white opacity-70">Toque para alterar a foto</span>
			</div>

			<!-- Form Section -->
			<div class="form-section">
				<div class="mb-4">
					<label class="ios-label">Nome de Usuário</label>
					<input v-model="user.nome" type="text" placeholder="Seu nome" class="ios-input-field w-100 mt-1" />
				</div>

				<div class="mb-6">
					<label class="ios-label">E-mail</label>
					<input v-model="user.email" type="email" placeholder="Seu e-mail" class="ios-input-field w-100 mt-1" />
				</div>

				<v-btn
					class="ios-action-btn w-100"
					elevation="4"
					:loading="loading"
					@click="saveProfile"
				>
					Salvar Alterações
				</v-btn>
			</div>
		</v-card>

		<!-- Photo Selector Modal -->
		<v-dialog v-model="photoModal" max-width="400" transition="dialog-bottom-transition">
			<v-card class="ios-modal-card">
				<v-card-title class="text-white text-center pt-4">Escolha uma opção</v-card-title>
				<v-list class="ios-modal-list bg-transparent">
					<v-list-item @click="selectLocalFile" class="modal-item">
						<template v-slot:prepend>
							<v-icon color="#00B8D4">mdi-folder-outline</v-icon>
						</template>
						<v-list-item-title class="text-white">Arquivos Locais</v-list-item-title>
					</v-list-item>

					<v-list-item @click="simulateCloud('Google Drive')" class="modal-item">
						<template v-slot:prepend>
							<v-icon color="#4285F4">mdi-google-drive</v-icon>
						</template>
						<v-list-item-title class="text-white">Google Drive</v-list-item-title>
					</v-list-item>

					<v-list-item @click="simulateCloud('Dropbox')" class="modal-item">
						<template v-slot:prepend>
							<v-icon color="#0061FF">mdi-dropbox</v-icon>
						</template>
						<v-list-item-title class="text-white">Dropbox</v-list-item-title>
					</v-list-item>

					<v-list-item @click="simulateCloud('OneDrive')" class="modal-item">
						<template v-slot:prepend>
							<v-icon color="#0078D4">mdi-microsoft-onedrive</v-icon>
						</template>
						<v-list-item-title class="text-white">OneDrive</v-list-item-title>
					</v-list-item>
				</v-list>
				<v-card-actions class="justify-center pb-4">
					<v-btn variant="text" color="rgba(255,255,255,0.6)" @click="photoModal = false">Cancelar</v-btn>
				</v-card-actions>
			</v-card>
		</v-dialog>

		<!-- Hidden file input -->
		<input
			type="file"
			ref="fileInput"
			style="display: none"
			accept="image/*"
			@change="onFileSelected"
		/>

		<v-snackbar v-model="snackbar" :color="snackbarColor" location="top" rounded="pill">
			{{ snackbarText }}
		</v-snackbar>
	</v-container>
</template>

<script>
import auth from '@/auth'
import UsuarioService from '@/services/UsuarioService'

export default {
	name: 'ProfileEdit',
	data() {
		return {
			user: {
				id: null,
				nome: '',
				email: '',
				foto_url: ''
			},
			loading: false,
			photoModal: false,
			previewUrl: null,
			snackbar: false,
			snackbarText: '',
			snackbarColor: 'success'
		}
	},
	created() {
		const cachedUser = auth.getUser()
		this.user = { ...cachedUser }
	},
	methods: {
		openPhotoSelector() {
			this.photoModal = true
		},
		selectLocalFile() {
			this.$refs.fileInput.click()
			this.photoModal = false
		},
		onFileSelected(event) {
			const file = event.target.files[0]
			if (!file) return

			// Em um app real, faríamos upload para um S3/Firebase aqui
			// Para demonstração, usaremos um FileReader para preview local
			const reader = new FileReader()
			reader.onload = (e) => {
				this.previewUrl = e.target.result
				this.user.foto_url = this.previewUrl // Mocking upload
			}
			reader.readAsDataURL(file)
		},
		simulateCloud(provider) {
			this.photoModal = false
			this.snackbarText = `Simulando integração com ${provider}...`
			this.snackbarColor = "info"
			this.snackbar = true

			// Simulação de delay de rede e seleção
			setTimeout(() => {
				this.previewUrl = "https://i.pravatar.cc/300?u=" + this.user.email
				this.user.foto_url = this.previewUrl
				this.snackbarText = `Foto selecionada do ${provider}`
				this.snackbarColor = "success"
				this.snackbar = true
			}, 1500)
		},
		async saveProfile() {
			this.loading = true
			try {
				await UsuarioService.atualizar(this.user.id, this.user)

				// Atualiza o estado global
				auth.login({ ...this.user }) // auth.login atualiza o state reativo

				this.snackbarText = "Perfil atualizado com sucesso!"
				this.snackbarColor = "success"
				this.snackbar = true

				setTimeout(() => {
					this.$router.push('/dashboard')
				}, 1500)
			} catch (error) {
				console.error(error)
				this.snackbarText = "Erro ao salvar perfil."
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
.profile-container {
	min-height: calc(100vh - 90px);
	padding: 40px 20px;
}

.ios-profile-card {
	width: 100%;
	max-width: 450px;
	background: rgba(45, 78, 115, 0.85) !important;
	backdrop-filter: blur(20px);
	border-radius: 28px !important;
	border: 1px solid rgba(255, 255, 255, 0.1);
	padding: 32px !important;
}

.avatar-wrapper {
	position: relative;
	cursor: pointer;
	transition: transform 0.3s ease;
}

.avatar-wrapper:hover {
	transform: scale(1.05);
}

.profile-avatar-glass {
	background: rgba(255, 255, 255, 0.1);
	border: 4px solid rgba(255, 255, 255, 0.2);
	overflow: hidden;
}

.avatar-overlay {
	position: absolute;
	top: 0;
	left: 0;
	right: 0;
	bottom: 0;
	background: rgba(0, 0, 0, 0.3);
	display: flex;
	align-items: center !important;
	justify-content: center;
	opacity: 0;
	transition: opacity 0.3s;
}

.avatar-wrapper:hover .avatar-overlay {
	opacity: 1;
}

.ios-label {
	color: rgba(255, 255, 255, 0.9);
	font-weight: 600;
	font-size: 14px;
	margin-left: 12px;
}

.ios-input-field {
	background: rgba(255, 255, 255, 0.1) !important;
	border: 1px solid rgba(255, 255, 255, 0.1) !important;
	border-radius: 16px;
	padding: 14px 20px;
	font-size: 16px;
	color: white;
	transition: all 0.2s ease;
}

.ios-input-field:focus {
	outline: none;
	background: rgba(255, 255, 255, 0.15) !important;
	border-color: rgba(0, 184, 212, 0.5) !important;
}

.ios-input--disabled {
	opacity: 0.5;
	cursor: not-allowed;
}

.ios-action-btn {
	background: linear-gradient(135deg, #00B8D4 0%, #0097A7 100%) !important;
	color: white !important;
	border-radius: 18px !important;
	height: 56px !important;
	text-transform: none !important;
	font-size: 17px !important;
	font-weight: 700 !important;
	margin-top: 10px;
}

/* Modals */
.ios-modal-card {
	background: rgba(30, 60, 90, 0.95) !important;
	backdrop-filter: blur(30px);
	border-radius: 24px !important;
	border: 1px solid rgba(255, 255, 255, 0.1);
}

.modal-item {
	margin: 8px 16px;
	border-radius: 12px;
	background: rgba(255, 255, 255, 0.05);
	transition: background 0.2s;
}

.modal-item:hover {
	background: rgba(255, 255, 255, 0.1);
}

.ios-modal-list :deep(.v-list-item-title) {
	font-weight: 600;
}
</style>
