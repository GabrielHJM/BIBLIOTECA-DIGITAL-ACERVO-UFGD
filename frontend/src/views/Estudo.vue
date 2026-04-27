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
							<!-- Premium Mandatory Glass Cover (Large) -->
							<div class="premium-book-cover-ios-large">
								<div class="mesh-gradient-large"></div>
								<div class="glass-overlay-ios-large">
									<div class="cover-content-ios-large">
										<v-icon class="cover-icon-ios mb-4" size="80" color="rgba(255,255,255,0.7)">{{ getBookIcon(material.categoria, material.titulo) }}</v-icon>
										<div class="cover-title-ios-large">{{ material.titulo }}</div>
										<div class="cover-summary-ios-large mt-3">{{ material.resumo }}</div>
									</div>
								</div>
								<div class="glass-shine"></div>
							</div>
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
								<div class="d-flex align-center">
									<v-btn
										icon="mdi-share-variant"
										variant="text"
										:color="$vuetify.theme.current.dark ? 'white' : 'grey-darken-1'"
										@click="shareBook"
										class="mr-2"
										title="Copiar Link do Livro"
									></v-btn>
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
							</div>
							<div class="d-flex align-center">
								<v-icon color="primary" size="20" class="mr-2">mdi-account</v-icon>
								<span class="text-h6 text-primary font-weight-medium mr-6">{{ material.autor }}</span>


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

			<!-- Citation Premium Dialog -->
			<v-dialog v-model="showCitationDialog" max-width="600" transition="dialog-bottom-transition" scrim="rgba(0, 10, 20, 0.85)">
				<div class="premium-citation-modal">
					<div class="glass-modal-bg"></div>
					<div class="modal-content-container pa-6">
						<div class="d-flex align-center mb-6">
							<v-icon color="primary" class="mr-3" size="28">mdi-format-quote-open</v-icon>
							<h3 class="text-h5 font-weight-black text-white">Gerar Citação</h3>
							<v-spacer></v-spacer>
							<v-btn icon="mdi-close" variant="tonal" color="white" class="ios-btn-close" size="small" @click="showCitationDialog = false"></v-btn>
						</div>

						<v-tabs v-model="citationTab" color="primary" bg-color="rgba(255,255,255,0.05)" class="premium-tabs rounded-xl mb-6" hide-slider>
							<v-tab value="ABNT" class="text-none font-weight-bold flex-1-1">ABNT</v-tab>
							<v-tab value="APA" class="text-none font-weight-bold flex-1-1">APA</v-tab>
							<v-tab value="BibTeX" class="text-none font-weight-bold flex-1-1">BibTeX</v-tab>
						</v-tabs>

						<v-window v-model="citationTab" class="rounded-xl overflow-visible">
							<v-window-item v-for="style in ['ABNT', 'APA', 'BibTeX']" :key="style" :value="style">
								<div class="citation-glass-box" @click="copyCitation(style)">
									<div class="citation-text">
										{{ getCitation(style) }}
									</div>
									<div class="copy-action-bar">
										<v-icon size="16" color="#00E5FF">mdi-content-copy</v-icon>
										<span>CLIQUE PARA COPIAR</span>
									</div>
								</div>
							</v-window-item>
						</v-window>
					</div>
				</div>
			</v-dialog>
		</div>

		<!-- Loading State -->
		<div v-else-if="loading" class="d-flex justify-center align-center fill-height loading-wrapper mt-16">
			<CloudLoader text="Abrindo material..." height="200px" />
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
			if (!this.material || !this.globalFavorites || !this.globalFavorites.list) return false;
			return this.globalFavorites.list.some(f => f.id == this.material.id);
		},
	},
	methods: {
		async toggleFavorite() {
			if (!this.isAuthenticated) {
				this.snackbarText = 'Faça login para favoritar!';
				this.snackbarColor = 'warning';
				this.snackbar = true;
				return;
			}

			// Feedback instantâneo (Animação GSAP primeiro)
			if (this.$refs.heartIcon) {
				const el = this.$refs.heartIcon.$el || this.$refs.heartIcon;
				const tl = gsap.timeline();
				tl.to(el, { scale: 1.5, rotation: 15, duration: 0.15, ease: "back.out(2)" })
					.to(el, { rotation: -15, duration: 0.1, ease: "none" })
					.to(el, { scale: 1, rotation: 0, duration: 0.3, ease: "elastic.out(1, 0.3)" });
			}

			try {
				const novoStatus = !this.isFavorited;
				await MaterialService.favoritar(auth.getUser().id, this.material.id, novoStatus);

				// Buscar favoritos globalmente para atualizar o status isFavorited na DOM
				await this.fetchGlobalFavorites();

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
		async shareBook() {
			try {
				const link = `${window.location.origin}/estudo/${this.material.id}`;
				await navigator.clipboard.writeText(link);
				this.snackbarText = 'Link copiado com sucesso para a área de transferência!';
				this.snackbarColor = 'success';
				this.snackbar = true;
			} catch (e) {
				console.error('Erro ao copiar link:', e);
				this.snackbarText = 'Não foi possível copiar o link.';
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
				return `${autor.toUpperCase()}. ${titulo}. Acervus Core, ${ano}.`;
			}
			if (style === 'APA') {
				return `${autor} (${ano}). ${titulo}. Acervus Core.`;
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
			// Força abertura direta do PDF no visualizador do navegador
			const url = this.material.pdf_url;
			const win = window.open(url, '_blank');
			if (win) {
				win.focus();
			} else {
				// Fallback if popup blocked
				window.location.assign(url);
			}
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
		},
		getBookIcon(category, title) {
			const text = ((category || '') + ' ' + (title || '')).toLowerCase();
			if (text.includes('tecnologia') || text.includes('comput') || text.includes('software') || text.includes('program') || text.includes('digital')) return 'mdi-laptop';
			if (text.includes('saúde') || text.includes('medicina') || text.includes('biolog') || text.includes('enferm') || text.includes('médic')) return 'mdi-heart-pulse';
			if (text.includes('direito') || text.includes('lei') || text.includes('juríd') || text.includes('advog')) return 'mdi-gavel';
			if (text.includes('matemát') || text.includes('física') || text.includes('cálculo') || text.includes('engenh')) return 'mdi-calculator';
			if (text.includes('história') || text.includes('socio') || text.includes('psico') || text.includes('filo')) return 'mdi-bank';
			if (text.includes('literat') || text.includes('poesia') || text.includes('romance')) return 'mdi-feather';
			if (text.includes('educação') || text.includes('ensino') || text.includes('pedagog')) return 'mdi-school';
			return 'mdi-book-open-page-variant';
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
			this.notify('Material indisponível ou não encontrado.', 'error');
			this.$router.replace('/explorar');
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

	.premium-book-cover-ios-large {
		width: 100%;
		aspect-ratio: 1 / 1.4; /* Aspecto mais fiel a livro real */
		max-height: 600px;
		margin: 0 auto;
		position: relative;
		overflow: hidden;
		border-radius: 32px;
		background: #000;
		display: flex;
		align-items: center;
		justify-content: center;
		box-shadow: 0 30px 60px rgba(0, 0, 0, 0.6);
		transition: all 0.8s cubic-bezier(0.16, 1, 0.3, 1);
	}

	.mesh-gradient-large {
		position: absolute;
		top: -50%;
		left: -50%;
		width: 200%;
		height: 200%;
		background: 
			radial-gradient(circle at 20% 20%, #007AFF 0%, transparent 40%),
			radial-gradient(circle at 80% 25%, #5AC8FA 0%, transparent 40%),
			radial-gradient(circle at 30% 80%, #0051FF 0%, transparent 40%),
			radial-gradient(circle at 75% 75%, #00D4E8 0%, transparent 40%),
			radial-gradient(circle at 50% 50%, #0033FF 0%, transparent 50%);
		filter: blur(80px);
		animation: mesh-rotate 25s linear infinite;
		opacity: 0.8;
	}

	@keyframes mesh-rotate {
		from { transform: rotate(0deg); }
		to { transform: rotate(360deg); }
	}

	.glass-overlay-ios-large {
		position: absolute;
		inset: 24px;
		background: rgba(255, 255, 255, 0.05);
		backdrop-filter: blur(40px) saturate(200%);
		border: 1px solid rgba(255, 255, 255, 0.2);
		border-radius: 20px;
		z-index: 2;
		display: flex;
		flex-direction: column;
		justify-content: center;
		align-items: center;
		padding: 48px;
		text-align: center;
		overflow: hidden;
	}

	.cover-content-ios-large {
		position: relative;
		z-index: 3;
		width: 100%;
		display: flex;
		flex-direction: column;
		align-items: center;
	}

	.cover-title-ios-large {
		color: #ffffff;
		font-family: 'Outfit', 'Inter', sans-serif;
		font-size: 2rem;
		font-weight: 800;
		line-height: 1.1;
		width: 100%;
		max-width: 100%;
		display: -webkit-box;
		-webkit-line-clamp: 3;
		-webkit-box-orient: vertical;
		overflow: hidden;
		text-shadow: 0 4px 20px rgba(0, 0, 0, 0.6);
		word-break: break-word;
	}

	.cover-summary-ios-large {
		color: rgba(255, 255, 255, 0.7);
		font-size: 1rem;
		font-weight: 500;
		line-height: 1.5;
		max-width: 80%;
		display: -webkit-box;
		-webkit-line-clamp: 4;
		-webkit-box-orient: vertical;
		overflow: hidden;
		text-align: center;
		word-break: break-word;
	}

	.glass-shine {
		position: absolute;
		top: 0;
		left: 0;
		width: 100%;
		height: 100%;
		background: linear-gradient(135deg, rgba(255,255,255,0.2) 0%, transparent 50%, rgba(255,255,255,0.1) 100%);
		pointer-events: none;
		z-index: 4;
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

	.premium-citation-modal {
		position: relative;
		border-radius: 32px;
		overflow: hidden;
		border: 1px solid rgba(255,255,255,0.15);
		box-shadow: 0 40px 100px rgba(0,0,0,0.8);
		transform: translateY(0);
	}

	.glass-modal-bg {
		position: absolute;
		top: 0; left: 0; right: 0; bottom: 0;
		background: rgba(15, 23, 42, 0.7);
		backdrop-filter: blur(40px) saturate(200%);
		-webkit-backdrop-filter: blur(40px) saturate(200%);
		z-index: 1;
	}

	.modal-content-container {
		position: relative;
		z-index: 2;
	}

	.ios-btn-close {
		border-radius: 50% !important;
		background: rgba(255,255,255,0.1) !important;
	}

	.premium-tabs .v-tab {
		border-radius: 100px !important;
		margin: 4px;
		transition: all 0.3s ease;
	}
	.premium-tabs .v-tab--selected {
		background: linear-gradient(135deg, #007AFF, #00C1FF) !important;
		color: white !important;
		box-shadow: 0 4px 15px rgba(0,122,255,0.4);
	}

	.citation-glass-box {
		background: rgba(255, 255, 255, 0.03);
		border: 1px solid rgba(255, 255, 255, 0.1);
		border-radius: 20px;
		padding: 24px;
		cursor: pointer;
		transition: all 0.4s var(--spring-easing);
		position: relative;
		overflow: hidden;
	}

	.citation-glass-box:hover {
		background: rgba(255, 255, 255, 0.08);
		border-color: rgba(0, 229, 255, 0.4);
		transform: translateY(-4px);
		box-shadow: 0 10px 30px rgba(0,0,0,0.3);
	}

	.citation-text {
		color: #ffffff;
		font-family: monospace;
		font-size: 0.95rem;
		line-height: 1.6;
		opacity: 0.9;
		margin-bottom: 20px;
	}

	.copy-action-bar {
		display: flex;
		align-items: center;
		justify-content: center;
		gap: 8px;
		background: rgba(0, 229, 255, 0.1);
		padding: 10px;
		border-radius: 12px;
		color: #00E5FF;
		font-weight: 800;
		font-size: 0.75rem;
		letter-spacing: 1px;
	}

	.loading-wrapper {
		min-height: 80vh;
	}

	@media (max-width: 960px) {
		.estudo-container { padding-top: 40px; }
		.text-h4 { font-size: 2rem !important; }
	}

</style>
