<template>
	<div class="estudo-view-root">
		<div class="estudo-container" v-if="!loading && material">
			<v-container class="max-width-content position-relative py-8">
				<!-- Header -->
				<div class="d-flex align-center mb-8">
					<v-btn variant="text" @click="goBack" prepend-icon="mdi-arrow-left" class="text-none font-weight-bold" :class="$vuetify.theme.current.dark ? 'text-white' : 'text-primary'" rounded="pill">
						Voltar
					</v-btn>
				</div>

				<v-row class="content-row">
					<!-- Left Column: Visual & Actions -->
					<v-col cols="12" md="4" lg="3">
						<div class="book-visual-card mb-6">
							<v-img :src="material.capa_url || 'https://images.unsplash.com/photo-1544947950-fa07a98d237f?q=80&w=400'" cover class="book-cover rounded-lg shadow-elevation-2">
								<template v-slot:placeholder>
									<div class="d-flex align-center justify-center fill-height bg-grey-darken-4">
										<v-icon color="rgba(255,255,255,0.2)" size="64">mdi-book-multiple</v-icon>
									</div>
								</template>
							</v-img>
						</div>

						<div class="action-buttons">
							<v-btn
								block
								color="primary"
								class="mb-3 text-none font-weight-bold ios-btn"
								@click="abrirPDF"
								prepend-icon="mdi-book-open-page-variant"
								elevation="0"
							>
								Acessar Material
							</v-btn>
						</div>
					</v-col>

					<!-- Right Column: Info -->
					<v-col cols="12" md="8" lg="9" class="info-col">
						<div class="mb-6">
							<div class="d-flex align-center mb-3 flex-wrap gap-2">
								<v-chip size="small" variant="flat" color="primary" class="font-weight-bold text-uppercase">{{ material.categoria }}</v-chip>
								<span class="text-caption opacity-60" :class="$vuetify.theme.current.dark ? 'text-white' : 'text-primary-darken-2'">{{ material.fonte || 'Repositório Público' }}</span>
							</div>
							<div class="d-flex align-center justify-space-between mb-2">
								<h1 class="text-h4 font-weight-bold mb-0" :class="$vuetify.theme.current.dark ? 'text-white' : 'text-primary-darken-4'">{{ material.titulo }}</h1>
								<v-btn
									icon
									variant="text"
									:color="isFavorited ? 'pink' : ($vuetify.theme.current.dark ? 'white' : 'primary')"
									@click="toggleFavorite"
									class="favorite-btn"
								>
									<v-icon ref="heartIcon" size="32">{{ isFavorited ? 'mdi-heart' : 'mdi-heart-outline' }}</v-icon>
								</v-btn>
							</div>
							<div class="d-flex align-center">
								<v-icon color="primary" size="20" class="mr-2">mdi-account</v-icon>
								<span class="text-h6 text-primary font-weight-medium mr-6">{{ material.autor }}</span>

								<v-rating
									:model-value="material.media_nota || 0"
									density="compact"
									color="amber"
									active-color="amber"
									size="small"
									@update:model-value="onRateBook($event)"
									hover
								></v-rating>
								<span class="ml-2 text-caption opacity-60" style="color: #ffffff !important;">({{ material.total_avaliacoes || 0 }})</span>
							</div>
						</div>

						<div class="metadata-grid mb-8">
							<div class="metadata-item">
								<span class="label">Publicado</span>
								<span class="value">{{ material.ano_publicacao || 'N/A' }}</span>
							</div>
							<div class="metadata-item">
								<span class="label">Páginas</span>
								<span class="value">{{ material.paginas || 'N/A' }}</span>
							</div>
							<div class="metadata-item">
								<span class="label">Identificação</span>
								<span class="value text-truncate" :title="material.isbn">{{ material.isbn || 'Digital ID' }}</span>
							</div>
						</div>

						<div class="description-section mb-8">
							<h3 class="text-subtitle-1 font-weight-bold text-primary mb-3 d-flex align-center">
								<v-icon size="18" class="mr-2">mdi-text-box-outline</v-icon> Resumo
							</h3>
							<div class="description-card">
								<p class="text-body-1 opacity-80 mb-0 line-height-relaxed" style="color: #ffffff !important;" v-text="material.descricao || 'Nenhum resumo disponível para este material.'">
								</p>
							</div>
						</div>

						<v-btn
							variant="tonal"
							color="primary"
							prepend-icon="mdi-format-quote-close"
							class="text-none font-weight-bold rounded-pill ios-btn"
							@click="showCitationDialog = true"
						>
							Gerar Citação Acadêmica
						</v-btn>
					</v-col>
				</v-row>

				<v-snackbar v-model="snackbar" :timeout="3000" :color="snackbarColor" location="bottom" rounded="pill">
					{{ snackbarText }}
				</v-snackbar>
			</v-container>

			<!-- Citation Dialog -->
			<v-dialog v-model="showCitationDialog" max-width="500">
				<v-card class="citation-card pa-6" rounded="xl">
					<div class="d-flex align-center mb-4">
						<h3 class="text-h6 font-weight-bold" style="color: #ffffff !important;">Citação</h3>
						<v-spacer></v-spacer>
						<v-btn icon="mdi-close" variant="text" size="small" @click="showCitationDialog = false"></v-btn>
					</div>

					<v-tabs v-model="citationTab" color="primary" density="compact" class="mb-4">
						<v-tab value="ABNT" class="text-none font-weight-bold">ABNT</v-tab>
						<v-tab value="APA" class="text-none font-weight-bold">APA</v-tab>
						<v-tab value="BibTeX" class="text-none font-weight-bold">BibTeX</v-tab>
					</v-tabs>

					<v-window v-model="citationTab">
						<v-window-item v-for="style in ['ABNT', 'APA', 'BibTeX']" :key="style" :value="style">
							<div class="citation-copy-box pa-4 rounded-lg" @click="copyCitation(style)">
								<p class="text-body-2 font-mono opacity-80 mb-0" style="color: #ffffff !important;">
									{{ getCitation(style) }}
								</p>
								<div class="mt-3 d-flex align-center">
									<v-icon size="14" color="primary" class="mr-1">mdi-content-copy</v-icon>
									<span class="text-caption color-primary font-weight-bold" style="color: var(--ios-blue) !important;">Copiar</span>
								</div>
							</div>
						</v-window-item>
					</v-window>
				</v-card>
			</v-dialog>
		</div>

		<!-- Loading State -->
		<div v-else-if="loading" class="d-flex justify-center align-center fill-height loading-wrapper">
			<v-progress-circular indeterminate color="primary" size="64" width="6"></v-progress-circular>
		</div>
	</div>
</template>

<script>
import MaterialService from '@/services/MaterialService'
import auth from '@/auth'
import { gsap } from 'gsap'

export default {
	name: 'EstudoPremiumRedesign',
	inject: ['notify', 'fetchGlobalFavorites', 'globalFavorites'],
	data: () => ({
		material: null,
		loading: true,
		isAuthenticated: false,
		showCitationDialog: false,
		citationTab: 'ABNT',
		snackbar: false,
		snackbarText: '',
		snackbarColor: 'success'
	}),
	computed: {
		isFavorited() {
			if (!this.material) return false;
			const favs = this.globalFavorites?.value || [];
			return favs.some(f => f.id === this.material.id);
		}
	},
	methods: {
		async toggleFavorite() {
			if (!this.isAuthenticated) {
				this.snackbarText = 'Faça login para favoritar!';
				this.snackbarColor = 'warning';
				this.snackbar = true;
				return;
			}
			try {
				const novoStatus = !this.isFavorited;
				await MaterialService.favoritar(auth.getUser().id, this.material.id, novoStatus);

				// Animation
				if (this.$refs.heartIcon) {
					const el = this.$refs.heartIcon.$el || this.$refs.heartIcon;
					gsap.fromTo(el, { scale: 1 }, { scale: 1.5, duration: 0.2, yoyo: true, repeat: 1, ease: "back.out(3)" });
				}

				this.snackbarText = novoStatus ? 'Adicionado aos favoritos!' : 'Removido dos favoritos.';
				this.snackbarColor = 'success';
				this.snackbar = true;
			} catch (e) {
				console.error('Erro ao favoritar:', e);
				this.snackbarText = 'Erro ao atualizar favorito.';
				this.snackbarColor = 'error';
				this.snackbar = true;
			}
		},
		async onRateBook(nota) {
			if (!this.isAuthenticated) {
				this.snackbarText = 'Faça login para avaliar!';
				this.snackbarColor = 'warning';
				this.snackbar = true;
				return;
			}
			try {
				await MaterialService.avaliar(auth.getUser().id, this.material.id, nota);
				this.material.media_nota = nota; // Feedback imediato visual
				this.snackbarText = 'Avaliação enviada com sucesso!';
				this.snackbarColor = 'success';
				this.snackbar = true;
			} catch (e) {
				console.error('Erro ao avaliar:', e);
				this.snackbarText = 'Erro ao enviar avaliação.';
				this.snackbarColor = 'error';
				this.snackbar = true;
			}
		},
		getCitation(style) {
			if (!this.material) return '';
			const autor = this.material.autor || 'AUTOR DESCONHECIDO';
			const titulo = this.material.titulo || 'Título não informado';
			const ano = this.material.ano_publicacao || new Date().getFullYear();

			if (style === 'ABNT') {
				return `${autor.toUpperCase()}. ${titulo}. Biblioteca Digital, ${ano}.`;
			}
			if (style === 'APA') {
				return `${autor} (${ano}). ${titulo}. Digital Library.`;
			}
			if (style === 'BibTeX') {
				return `@article{citekey, author={${autor}}, title={${titulo}}, year={${ano}}, journal={BD-PSCII}}`;
			}
			return '';
		},
		async copyCitation(style) {
			const text = this.getCitation(style);
			await navigator.clipboard.writeText(text);
			this.snackbarText = 'Citação copiada para transferência!';
			this.snackbarColor = 'success';
			this.snackbar = true;
		},
		abrirPDF() {
			if (!this.material || !this.material.pdf_url) return;
			// Abre diretamente assumindo que é um link de PDF válido
			window.open(this.material.pdf_url, '_blank');
		},
		baixarPDF() {
			if (!this.material || !this.material.pdf_url) return;
			// Clean programmatic interaction that doesn't leak DOM elements or event listeners
			const isIos = /iPad|iPhone|iPod/.test(navigator.userAgent) && !window.MSStream;
			if(isIos) {
				window.location.assign(this.material.pdf_url);
			} else {
				Object.assign(document.createElement('a'), {
					href: this.material.pdf_url,
					download: `${this.material.titulo || 'document'}.pdf`,
					target: '_blank'
				}).click();
			}
		},
		goBack() {
			window.history.length > 1 ? this.$router.go(-1) : this.$router.push('/explorar');
		}
	},
	async mounted() {
		this.isAuthenticated = auth.isAuthenticated();
		const id = this.$route.params.id;
		try {
			const response = await MaterialService.obterDetalhes(id);
			this.material = response.data;
			if (this.isAuthenticated && this.material) {
				MaterialService.registrarLeitura(auth.getUser().id, this.material.id);
			}
		} catch (error) {
			console.error('Erro ao carregar material:', error);
		} finally {
			this.loading = false;
		}
	}
}
</script>

<style scoped>
	.estudo-container {
		min-height: 100vh;
		background: transparent;
		padding-top: 40px;
		padding-bottom: 60px;
	}

	.max-width-content {
		max-width: 1000px !important;
	}

	.book-cover {
		width: 100%;
		aspect-ratio: 1/1.4;
		border: 1px solid rgba(255, 255, 255, 0.1);
	}

	.shadow-elevation-2 {
		box-shadow: 0 8px 24px rgba(0, 0, 0, 0.4) !important;
	}

	.ios-btn {
		border-radius: 12px !important;
		letter-spacing: 0.2px;
	}

	.gap-2 {
		gap: 8px;
	}

	.metadata-grid {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(120px, 1fr));
		gap: 16px;
		background: var(--glass-bg);
		border-radius: 16px;
		padding: 20px;
		border: 1px solid var(--glass-border);
	}

	.metadata-item {
		display: flex;
		flex-direction: column;
	}

	.metadata-item .label {
		font-size: 12px;
		font-weight: 600;
		color: var(--v-theme-on-surface);
		opacity: 0.5;
		margin-bottom: 4px;
		text-transform: uppercase;
	}

	.metadata-item .value {
		font-size: 16px;
		font-weight: 500;
		color: var(--v-theme-on-surface);
	}

	.description-card {
		background: var(--glass-bg);
		backdrop-filter: var(--glass-blur);
		-webkit-backdrop-filter: var(--glass-blur);
		border-radius: 16px;
		padding: 24px;
		border: 1px solid var(--glass-border);
	}

	.line-height-relaxed {
		line-height: 1.6;
	}

	.citation-card {
		background: var(--glass-bg) !important;
		backdrop-filter: var(--glass-blur);
		-webkit-backdrop-filter: var(--glass-blur);
		border: 1px solid var(--glass-border) !important;
		color: #ffffff !important;
	}

	.citation-copy-box {
		background: var(--glass-bg);
		border: 1px solid var(--glass-border);
		cursor: pointer;
		transition: background 0.2s;
	}

	.citation-copy-box:hover {
		background: rgba(0, 0, 0, 0.5);
	}

	.loading-wrapper {
		min-height: 80vh;
	}

	@media (max-width: 960px) {
		.estudo-container { padding-top: 40px; }
		.text-h4 { font-size: 2rem !important; }
	}
</style>
