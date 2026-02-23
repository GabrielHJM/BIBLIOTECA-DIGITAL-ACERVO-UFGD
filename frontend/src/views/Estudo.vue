<template>
	<div class="estudo-container mt-8">
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
								<img :src="livros[0].link" :alt="livros[0].nome" class="book-cover" />
							</div>
							<h1 class="text-white text-h5 font-weight-bold leading-tight">{{ livros[0].nome }}</h1>
						</div>

						<v-divider class="opacity-20 mb-6"></v-divider>

						<div class="item-details text-white">
							<div class="detail-row mb-3 d-flex justify-space-between">
								<span class="opacity-60">Publicado</span>
								<span class="font-weight-bold">05/10/2022</span>
							</div>
							<div class="detail-row mb-3 d-flex justify-space-between">
								<span class="opacity-60">Páginas</span>
								<span class="font-weight-bold">220</span>
							</div>
							<div class="detail-row mb-3 d-flex justify-space-between">
								<span class="opacity-60">Volume</span>
								<span class="font-weight-bold">1</span>
							</div>
							<div class="detail-row mb-3 d-flex justify-space-between">
								<span class="opacity-60">Autor</span>
								<span class="font-weight-bold">Michael T. Goodrich</span>
							</div>
						</div>

						<v-divider class="opacity-20 my-6"></v-divider>

						<div class="rating-section mb-6">
							<label class="text-white opacity-60 text-caption d-block mb-2">Avalie este conteúdo</label>
							<v-rating
								hover
								:length="5"
								:model-value="4"
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
							<div class="viewer-actions">
								<v-btn icon="mdi-fullscreen" variant="text" color="white" size="small"></v-btn>
								<v-btn icon="mdi-download" variant="text" color="white" size="small"></v-btn>
							</div>
						</div>
						<div class="content-body">
							<div id="embed-doc" class="embed-doc"></div>
						</div>
					</v-card>
				</v-col>
			</v-row>
		</v-container>
	</div>
</template>

<script>
import livros from '../../../livros.json'

export default {
	name: 'EstudoPage',
	data: () => ({
		livros: livros
	}),
	methods: {
		embedDocumento(pdfUrl) {
			const container = document.getElementById('embed-doc');
			if (!container) return;

			container.innerHTML = '';
			const iframe = document.createElement('iframe');
			iframe.setAttribute('src', pdfUrl + '#toolbar=0&navpanes=0&scrollbar=0');
			iframe.setAttribute('width', '100%');
			iframe.setAttribute('height', '100%');
			iframe.style.border = 'none';
			container.appendChild(iframe);
		}
	},
	mounted() {
		const exemploPDF = 'https://dn790006.ca.archive.org/0/items/estruturas-de-dados-e-algoritmos-em-java-pdfdrive/Estruturas%20de%20dados%20e%20algoritmos%20em%20JAVA%20(%20PDFDrive%20).pdf'
		this.embedDocumento(exemploPDF);
	}
}
</script>

<style scoped>
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

	:deep(.v-rating__item) {
		padding: 0 4px !important;
	}
</style>
