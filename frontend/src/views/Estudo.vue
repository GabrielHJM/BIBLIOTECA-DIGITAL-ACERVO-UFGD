<template>
	<div class="estudo-container mt-8" v-if="!loading && material">
		<v-container fluid class="pa-4">
			<v-row>
				<!-- Sidebar: Book Info -->
				<v-col cols="12" lg="4" class="pa-4">
					<v-card class="ios-card-premium pa-6 sticky-sidebar" elevation="12">
						<v-btn class="ios-btn-back mb-6" variant="text" @click="$router.push('/dashboard')">
							<v-icon class="mr-2">mdi-chevron-left</v-icon> Voltar
						</v-btn>

						<div class="text-center mb-6">
							<div class="book-cover-wrapper mx-auto mb-4">
								<img :src="material.capa_url" :alt="material.titulo" class="book-cover" />
							</div>
							<h1 class="text-white text-h5 font-weight-bold leading-tight">{{ material.titulo }}</h1>
						</div>

						<v-divider class="opacity-20 mb-6"></v-divider>

						<div class="item-details text-white">
							<div class="detail-row mb-3 d-flex justify-space-between">
								<span class="opacity-60">Publicado</span>
								<span class="font-weight-bold">{{ material.ano_publicacao }}</span>
							</div>
							<div class="detail-row mb-3 d-flex justify-space-between">
								<span class="opacity-60">ISBN</span>
								<span class="font-weight-bold">{{ material.isbn || 'N/A' }}</span>
							</div>
							<div class="detail-row mb-3 d-flex justify-space-between">
								<span class="opacity-60">Categoria</span>
								<span class="font-weight-bold">{{ material.categoria }}</span>
							</div>
							<div class="detail-row mb-3 d-flex justify-space-between">
								<span class="opacity-60">Autor</span>
								<span class="font-weight-bold">{{ material.autor }}</span>
							</div>
						</div>

						<v-divider class="opacity-20 my-6"></v-divider>

						<div class="rating-section mb-6">
							<label class="text-white opacity-60 text-caption d-block mb-2">Avalie este conteúdo</label>
							<v-rating
								hover
								:length="5"
								:model-value="material.media_nota || 0"
								active-color="amber"
								color="grey-lighten-1"
								density="comfortable"
							></v-rating>
						</div>

						<div class="share-section">
							<label class="text-white opacity-60 text-caption d-block mb-3">Compartilhar</label>
							<div class="d-flex justify-space-between">
								<v-btn icon="mdi-instagram" variant="tonal" color="white" density="comfortable"></v-btn>
								<v-btn icon="mdi-facebook" variant="tonal" color="white" density="comfortable"></v-btn>
								<v-btn icon="mdi-reddit" variant="tonal" color="white" density="comfortable"></v-btn>
								<v-btn icon="mdi-whatsapp" variant="tonal" color="white" density="comfortable"></v-btn>
							</div>
						</div>
					</v-card>
				</v-col>

				<!-- Main Content: PDF Viewer -->
				<v-col cols="12" lg="8" class="pa-4">
					<v-card class="ios-viewer-card" elevation="12">
						<div class="viewer-header pa-4 d-flex justify-space-between align-center">
							<div class="d-flex align-center">
								<v-icon color="cyan" class="mr-2">mdi-file-pdf-box</v-icon>
								<span class="text-white font-weight-medium">Modo de Leitura</span>
							</div>
							<div class="viewer-actions" v-if="material.pdf_url">
								<v-btn icon="mdi-fullscreen" variant="text" color="white" size="small" @click="toggleFullscreen"></v-btn>
								<v-btn icon="mdi-download" variant="text" color="white" size="small" :href="material.pdf_url" target="_blank"></v-btn>
							</div>
						</div>
						<div class="content-body">
							<div id="embed-doc" class="embed-doc" v-if="isAuthenticated"></div>
							<div v-else class="login-cta-container d-flex flex-column align-center justify-center pa-8 text-center fill-height">
								<v-icon size="80" color="cyan" class="mb-4 opacity-20">mdi-lock-outline</v-icon>
								<h2 class="text-h5 font-weight-bold text-white mb-2">Conteúdo Restrito</h2>
								<p class="text-white opacity-60 mb-6">
									Crie uma conta ou faça login para acessar este material completo e utilizar as ferramentas de estudo.
								</p>
								<div class="d-flex flex-column w-100" style="gap: 12px; max-width: 280px;">
									<v-btn color="cyan" size="large" rounded="pill" block @click="$router.push('/login')">
										Login
									</v-btn>
									<v-btn variant="outlined" color="white" size="large" rounded="pill" block @click="$router.push('/cadastro')">
										Cadastre-se Grátis
									</v-btn>
								</div>
							</div>
							<div v-if="isAuthenticated && !material.pdf_url" class="text-center mt-12 text-white opacity-60">
								PDF não disponível para este material.
							</div>
						</div>
					</v-card>
				</v-col>
			</v-row>
		</v-container>
	</div>

	<div v-else-if="loading" class="d-flex justify-center align-center" style="height: 80vh">
		<v-progress-circular indeterminate color="cyan" size="64"></v-progress-circular>
	</div>

	<div v-else class="text-center mt-12 text-white">
		<h2>Material não encontrado</h2>
		<v-btn class="mt-4" @click="goBack">Voltar</v-btn>
	</div>
</template>

<script>
import MaterialService from '@/services/MaterialService'
import auth from '@/auth'

export default {
	name: 'EstudoPage',
	data: () => ({
		material: null,
		loading: true,
		isAuthenticated: false
	}),
	methods: {
		goBack() {
			if (this.isAuthenticated) {
				this.$router.push('/dashboard')
			} else {
				this.$router.push('/explorar')
			}
		},
		embedDocumento(pdfUrl) {
			this.$nextTick(() => {
				const container = document.getElementById('embed-doc');
				if (!container) return;

				container.innerHTML = '';
				const iframe = document.createElement('iframe');
				iframe.setAttribute('src', pdfUrl + '#toolbar=0&navpanes=0&scrollbar=0');
				iframe.setAttribute('width', '100%');
				iframe.setAttribute('height', '100%');
				iframe.style.border = 'none';
				container.appendChild(iframe);
			});
		},
		toggleFullscreen() {
			if (!this.isAuthenticated) return
			const elem = document.getElementById('embed-doc');
			if (elem && elem.requestFullscreen) {
				elem.requestFullscreen();
			}
		}
	},
	async mounted() {
		this.isAuthenticated = auth.isAuthenticated();
		const id = this.$route.params.id;
		if (!id) {
			this.loading = false;
			return;
		}

		try {
			const response = await MaterialService.obterDetalhes(id);
			this.material = response.data;
			if (this.isAuthenticated && this.material && this.material.pdf_url) {
				this.embedDocumento(this.material.pdf_url);
			}
		} catch (error) {
			console.error('Erro ao buscar detalhes do material:', error);
		} finally {
			this.loading = false;
		}
	}
}
</script>

<style scoped>
	.login-cta-container {
		background: rgba(0,0,0,0.2);
		border-radius: 20px;
	}
	.estudo-container {
		min-height: 100vh;
		padding-bottom: 40px;
	}

	.ios-card-premium {
		background: rgba(58, 99, 145, 0.6) !important;
		backdrop-filter: blur(20px);
		border-radius: 32px !important;
		border: 1px solid rgba(255, 255, 255, 0.1);
	}

	.sticky-sidebar {
		position: sticky;
		top: 100px;
	}

	@media (max-width: 1264px) {
		.sticky-sidebar {
			position: relative;
			top: 0;
			margin-bottom: 20px;
		}
	}

	.ios-btn-back {
		color: rgba(255, 255, 255, 0.8) !important;
		text-transform: none !important;
		font-weight: 600 !important;
	}

	.book-cover-wrapper {
		width: 160px;
		height: 220px;
		border-radius: 16px;
		overflow: hidden;
		background: rgba(0,0,0,0.2);
		box-shadow: 0 12px 24px rgba(0,0,0,0.4);
	}

	.book-cover {
		width: 100%;
		height: 100%;
		object-fit: contain;
	}

	.opacity-60 { opacity: 0.6; }
	.opacity-20 { opacity: 0.2; }

	.ios-viewer-card {
		background: #2a2a2a !important;
		border-radius: 32px !important;
		overflow: hidden;
		height: 850px;
		border: 1px solid rgba(255, 255, 255, 0.05);
		display: flex;
		flex-direction: column;
	}

	.viewer-header {
		background: rgba(0,0,0,0.3);
		border-bottom: 1px solid rgba(255, 255, 255, 0.05);
	}

	.content-body {
		flex: 1;
		background: #1a1a1a;
	}

	.embed-doc {
		width: 100%;
		height: 100%;
	}

	@media (max-width: 600px) {
		.ios-card-premium {
			border-radius: 24px !important;
			padding: 20px !important;
		}

		.book-cover-wrapper {
			width: 120px;
			height: 160px;
		}

		.ios-viewer-card {
			height: 500px;
			border-radius: 24px !important;
		}

		.viewer-header {
			padding: 12px !important;
		}
	}

	:deep(.v-rating__item) {
		padding: 0 4px !important;
	}
</style>
