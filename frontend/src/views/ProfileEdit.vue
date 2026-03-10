<template>
	<v-container class="profile-container d-flex justify-center align-center">
		<v-card class="ios-profile-card" elevation="0">
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
				<v-chip size="small" :color="userRoleName === 'Professor' ? 'primary' : 'white'" :variant="userRoleName === 'Professor' ? 'flat' : 'outlined'" class="mt-4 font-weight-bold" style="color: #0b192c !important;">
					<v-icon start size="16">{{ userRoleName === 'Professor' ? 'mdi-school' : 'mdi-account-school' }}</v-icon>
					{{ userRoleName }}
				</v-chip>
				<span class="mt-2 text-caption text-white opacity-70">Toque para alterar a foto</span>
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

				<!-- Delete Account Button -->
				<v-btn
					class="w-100 mt-4 text-none d-flex align-center justify-center font-weight-medium"
					color="#ff3b30"
					variant="text"
					@click="confirmDeleteModal = true"
				>
					<v-icon left size="20" class="mr-2">mdi-delete-outline</v-icon>
					Excluir Minha Conta
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
							<v-icon color="primary">mdi-folder-outline</v-icon>
						</template>
						<v-list-item-title class="text-white">Arquivos Locais</v-list-item-title>
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

		<!-- Delete Account Confirmation Modal -->
		<v-dialog v-model="confirmDeleteModal" max-width="400" transition="dialog-bottom-transition">
			<v-card class="ios-modal-card">
				<v-card-title class="text-white text-center pt-6 text-h6 font-weight-bold" style="color: #ff3b30 !important;">
					<v-icon color="#ff3b30" class="mb-2" size="36">mdi-alert-circle-outline</v-icon>
					<br>Excluir Conta
				</v-card-title>
				<v-card-text class="text-white text-center opacity-80 pb-6 px-6">
					Tem certeza de que deseja excluir sua conta permanentemente? Esta ação não pode ser desfeita e todos os seus dados serão apagados.
				</v-card-text>
				<v-card-actions class="d-flex flex-column px-6 pb-6" style="gap: 12px;">
					<v-btn
						class="w-100 rounded-lg font-weight-bold ma-0"
						color="#ff3b30"
						variant="flat"
						height="48"
						:loading="deleting"
						@click="deleteAccount"
					>
						Sim, Excluir Conta
					</v-btn>
					<v-btn
						class="w-100 rounded-lg font-weight-medium ma-0 text-white"
						style="background: rgba(255, 255, 255, 0.1);"
						variant="flat"
						height="48"
						@click="confirmDeleteModal = false"
					>
						Cancelar
					</v-btn>
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
			deleting: false,
			photoModal: false,
			confirmDeleteModal: false,
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
	computed: {
		userRoleName() {
			if (this.user.tipo_usuario_id === 1) return 'Aluno';
			if (this.user.tipo_usuario_id === 2) return 'Professor';
			return 'Membro';
		}
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
		},
		async deleteAccount() {
			this.deleting = true;
			try {
				await UsuarioService.apagar(this.user.id);

				this.snackbarText = "Conta excluída com sucesso.";
				this.snackbarColor = "info";
				this.snackbar = true;

				auth.logout();

				setTimeout(() => {
					this.$router.push('/');
				}, 1500);
			} catch (error) {
				console.error("Erro ao excluir conta:", error);
				this.snackbarText = "Erro ao excluir conta. Tente novamente mais tarde.";
				this.snackbarColor = "error";
				this.snackbar = true;
			} finally {
				this.deleting = false;
				this.confirmDeleteModal = false;
			}
		}
	}
}
</script>

<style scoped>
.profile-container {
	min-height: calc(100vh - 90px);
	padding: 0;
}

.ios-profile-card {
	width: 100%;
	max-width: 450px;
	background: var(--glass-bg) !important;
	backdrop-filter: var(--glass-blur);
	-webkit-backdrop-filter: var(--glass-blur);
	border-radius: 28px !important;
	border: 1px solid var(--glass-border);
	padding: 32px !important;
	box-shadow: 0 32px 64px rgba(0, 0, 0, 0.2) !important;
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
	background: rgba(255, 255, 255, 0.05) !important;
	border: 1px solid rgba(255, 255, 255, 0.1) !important;
	border-radius: 16px;
	padding: 14px 20px;
	font-size: 16px;
	color: #ffffff;
	transition: all 0.2s ease;
}

.ios-input-field:focus {
	outline: none;
	background: rgba(255, 255, 255, 0.15) !important;
	border-color: rgba(0, 122, 255, 0.5) !important;
}

.ios-input--disabled {
	opacity: 0.5;
	cursor: not-allowed;
}

.ios-action-btn {
	background: linear-gradient(135deg, #007AFF 0%, #0056B3 100%) !important;
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
	background: var(--glass-bg) !important;
	backdrop-filter: var(--glass-blur);
	-webkit-backdrop-filter: var(--glass-blur);
	border-radius: 24px !important;
	border: 1px solid var(--glass-border);
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
