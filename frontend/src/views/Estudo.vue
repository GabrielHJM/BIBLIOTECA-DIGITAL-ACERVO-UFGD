<template>
	<div class="estudo-container mt-8" v-if="!loading && material">
		<v-container fluid class="pa-0 pa-sm-4">
			<v-row no-gutters>
				<!-- Sidebar: Book Info & Citation Tools -->
				<v-col cols="12" lg="3" class="pa-4">
					<v-card class="ios-glass-card pa-6 sticky-sidebar" rounded="xl">
						<v-btn class="ios-btn-back mb-6" variant="text" @click="goBack" density="comfortable">
							<v-icon class="mr-2">mdi-arrow-left</v-icon> Voltar
						</v-btn>

						<div class="text-center mb-6">
							<div class="book-cover-wrapper-premium mx-auto mb-4">
								<v-img :src="material.capa_url || 'https://images.unsplash.com/photo-1544947950-fa07a98d237f?q=80&w=200'" cover class="book-cover-img">
									<template v-slot:placeholder>
										<div class="d-flex align-center justify-center fill-height bg-grey-darken-4">
											<v-icon color="rgba(255,255,255,0.2)" size="48">mdi-book-open-blank-variant</v-icon>
										</div>
									</template>
								</v-img>
							</div>
							<h1 class="text-h6 font-weight-black lh-tight text-white">{{ material.titulo }}</h1>
							<p class="text-caption text-cyan mt-1 font-weight-bold">{{ material.autor }}</p>
						</div>

						<v-divider class="opacity-10 mb-6"></v-divider>

						<!-- Metadata grid -->
						<div class="meta-grid">
							<div class="meta-item">
								<span class="meta-label">Publicação</span>
								<span class="meta-value">{{ material.ano_publicacao || 'N/A' }}</span>
							</div>
							<div class="meta-item">
								<span class="meta-label">Fonte</span>
								<span class="meta-value">{{ material.fonte || 'Digital' }}</span>
							</div>
							<div class="meta-item">
								<span class="meta-label">Categoria</span>
								<span class="meta-value">{{ material.categoria }}</span>
							</div>
						</div>

						<!-- Citation Button -->
						<v-btn
							block
							variant="outlined"
							color="cyan"
							class="mt-6 text-none font-weight-bold"
							prepend-icon="mdi-format-quote-close"
							rounded="lg"
							@click="showCitationDialog = true"
						>
							Citar este Material
						</v-btn>

						<v-btn
							v-if="material.pdf_url"
							block
							color="cyan"
							class="mt-3 text-none font-weight-bold"
							prepend-icon="mdi-book-open-page-variant"
							rounded="lg"
							@click="handleReadOnline"
						>
							Ler agora
						</v-btn>

						<v-btn
							v-if="material.pdf_url"
							block
							variant="text"
							color="white"
							class="mt-2 text-none"
							prepend-icon="mdi-cloud-download-outline"
							size="small"
							:href="material.pdf_url"
							target="_blank"
							download
						>
							Baixar PDF
						</v-btn>
					</v-card>
				</v-col>

				<!-- Main Content: Reading View -->
				<v-col cols="12" lg="9" class="pa-4">
					<v-card class="reader-card" rounded="xl">
						<!-- Custom Reader Toolbar -->
						<div class="reader-toolbar d-flex align-center px-4 py-2">
							<div class="d-flex align-center">
								<v-icon color="cyan" class="mr-2">mdi-auto-fix</v-icon>
								<span class="font-weight-bold text-caption text-uppercase tracking-widest text-white">Visualizador Acadêmico</span>
							</div>
							<v-spacer></v-spacer>
							<div class="d-flex align-center gap-2">
								<v-btn icon="mdi-magnify-plus-outline" size="x-small" variant="text" color="white" @click="zoomIn"></v-btn>
								<v-btn icon="mdi-magnify-minus-outline" size="x-small" variant="text" color="white" @click="zoomOut"></v-btn>
								<v-divider vertical class="mx-2 opacity-20"></v-divider>
								<v-btn icon="mdi-fullscreen" size="small" variant="text" color="cyan" @click="toggleFullscreen" title="Foco total"></v-btn>
							</div>
						</div>

						<div class="viewer-area fill-height" ref="viewerContainer">
							<template v-if="isAuthenticated">
								<!-- PDF Embed Viewer -->
								<div v-if="isEmbeddablePDF(material.pdf_url)" class="pdf-container">
									<iframe
										:src="`https://docs.google.com/viewer?url=${encodeURIComponent(material.pdf_url)}&embedded=true`"
										class="safe-iframe"
										width="100%"
										height="100%"
										frameborder="0"
									></iframe>
								</div>

								<!-- External Material Fallback -->
								<div v-else-if="material.pdf_url" class="external-reader-overlay d-flex flex-column align-center justify-center text-center pa-10">
									<div class="external-icon-box mb-6">
										<v-icon size="80" color="cyan">mdi-earth-arrow-right</v-icon>
									</div>
									<h2 class="text-h4 font-weight-black text-white mb-3">Portal Acadêmico Externo</h2>
									<p class="text-body-1 text-white opacity-60 mb-8 max-w-lg">
										Este material faz parte de uma base de dados externa ({{ material.fonte || 'Repositório Digital' }}).
										Para garantir o acesso à versão mais recente e ferramentas interativas, utilize o portal original.
									</p>
									<v-btn
										variant="elevated"
										color="cyan"
										size="x-large"
										rounded="pill"
										:href="material.pdf_url"
										target="_blank"
										class="px-10 font-weight-bold"
										elevation="8"
									>
										Abrir no Portal de Origem <v-icon right size="20" class="ml-3">mdi-open-in-new</v-icon>
									</v-btn>

									<v-alert
										v-if="material.pdf_url.includes('doi.org')"
										variant="tonal"
										color="cyan"
										icon="mdi-information-outline"
										class="mt-10 text-caption max-w-md"
										rounded="lg"
									>
										Este material possui um DOI (Digital Object Identifier), garantindo sua procedência acadêmica oficial.
									</v-alert>
								</div>

								<!-- Fully Unavailable State -->
								<div v-else class="unavailable-state d-flex flex-column align-center justify-center text-center pa-6">
									<div class="empty-state-icon mb-4">
										<v-icon size="64" color="white" class="opacity-20">mdi-file-edit-outline</v-icon>
									</div>
									<h3 class="text-h6 text-white opacity-40 font-weight-medium">Conteúdo Restrito ou em Manutenção</h3>
									<p class="text-caption text-white opacity-20 mt-1">
										Estamos trabalhando para disponibilizar este material digitalmente em breve.
									</p>
									<v-btn variant="text" color="white" class="mt-4 opacity-40" size="small" prepend-icon="mdi-alert-circle-outline">
										Notificar erro de link
									</v-btn>
								</div>
							</template>
							<div v-else class="lock-screen d-flex flex-column align-center justify-center pa-10 text-center">
								<div class="lock-icon-wrap mb-6">
									<v-icon size="64" color="white">mdi-lock-open-outline</v-icon>
								</div>
								<h2 class="text-h4 font-weight-black text-white mb-2">Conhecimento Protegido</h2>
								<p class="text-white opacity-60 mb-8 max-w-lg">
									Este material acadêmico requer identificação para acesso completo. Entre em sua conta para utilizar o leitor e as ferramentas de citação.
								</p>
								<div class="d-flex gap-4">
									<v-btn color="cyan" rounded="pill" size="large" @click="$router.push('/login')" class="px-8">Login</v-btn>
									<v-btn variant="outlined" color="white" rounded="pill" size="large" @click="$router.push('/cadastro')">Criar Conta</v-btn>
								</div>
							</div>
						</div>
					</v-card>

					<!-- Description & Abstract -->
					<v-card class="ios-glass-card mt-6 pa-8" rounded="xl" v-if="material.descricao">
						<h3 class="text-h5 font-weight-bold text-white mb-4 d-flex align-center">
							<v-icon color="cyan" class="mr-3">mdi-text-long</v-icon> Resumo / Descrição
						</h3>
						<div class="text-body-1 text-white opacity-80 leading-relaxed abstract-text">
							{{ material.descricao }}
						</div>
					</v-card>
				</v-col>
			</v-row>
		</v-container>

		<!-- Citation Dialog -->
		<v-dialog v-model="showCitationDialog" max-width="600">
			<v-card class="ios-glass-card pa-6" rounded="xl">
				<div class="d-flex align-center mb-6">
					<v-icon color="cyan" size="24" class="mr-2">mdi-format-quote-close</v-icon>
					<h3 class="text-h6 font-weight-bold text-white">Gerador de Citação</h3>
					<v-spacer></v-spacer>
					<v-btn icon="mdi-close" variant="text" size="small" @click="showCitationDialog = false"></v-btn>
				</div>

				<div class="citation-box pa-4 mb-4 rounded-lg bg-black-opacity-40">
					<div class="text-overline text-cyan mb-1">ABNT (NBR 6023)</div>
					<p class="text-caption text-white font-mono leading-tight citation-content" @click="copyCitation($event)">
						{{ getCitation('ABNT') }}
					</p>
				</div>

				<div class="citation-box pa-4 mb-4 rounded-lg bg-black-opacity-40">
					<div class="text-overline text-cyan mb-1">APA (7ª Edition)</div>
					<p class="text-caption text-white font-mono leading-tight citation-content" @click="copyCitation($event)">
						{{ getCitation('APA') }}
					</p>
				</div>

				<div class="citation-box pa-4 rounded-lg bg-black-opacity-40">
					<div class="text-overline text-cyan mb-1">BibTeX</div>
					<p class="text-caption text-white font-mono leading-tight citation-content whitespace-pre" @click="copyCitation($event)">
						{{ getCitation('BibTeX') }}
					</p>
				</div>

				<v-alert type="info" variant="tonal" class="mt-4 text-caption" density="compact">
					Clique na citação para copiar automaticamente para a área de transferência.
				</v-alert>
			</v-card>
		</v-dialog>
	</div>

	<div v-else-if="loading" class="d-flex justify-center align-center fill-height" style="height: 100vh">
		<v-progress-circular indeterminate color="cyan" size="80" width="8"></v-progress-circular>
	</div>
</template>

<script>
import MaterialService from '@/services/MaterialService'
import auth from '@/auth'

export default {
	name: 'EstudoPremium',
	data: () => ({
		material: null,
		loading: true,
		isAuthenticated: false,
		showCitationDialog: false,
		zoomLevel: 100
	}),
	methods: {
		getCitation(style) {
			if (!this.material) return '';
			const autor = this.material.autor || 'AUTOR DESCONHECIDO';
			const titulo = this.material.titulo || 'Título não informado';
			const ano = this.material.ano_publicacao || new Date().getFullYear();

			if (style === 'ABNT') {
				// SOBRENOME, Nome. Titulo: subtitulo. Local: Editora, ano.
				const autorUpper = autor.toUpperCase();
				return `${autorUpper}. ${titulo}. Biblioteca Digital, ${ano}. Disponível em: Online.`;
			}
			if (style === 'APA') {
				// Surname, N. (Year). Title. Publisher.
				return `${autor} (${ano}). ${titulo}. Digital Library Archive.`;
			}
			if (style === 'BibTeX') {
				const citeKey = autor.split(' ')[0].toLowerCase() + ano;
				return `@article{${citeKey},
  author = {${autor}},
  title = {${titulo}},
  year = {${ano}},
  journal = {Biblioteca Digital PSCII}
}`;
			}
			return '';
		},
		copyCitation(event) {
			const text = event.target.innerText;
			navigator.clipboard.writeText(text);
			// Aqui poderíamos disparar um toast/snack
		},
		handleReadOnline() {
			if (this.$refs.viewerContainer) {
				this.$refs.viewerContainer.scrollIntoView({ behavior: 'smooth' });
			}
		},
		toggleFullscreen() {
			const elem = this.$refs.viewerContainer;
			if (elem && elem.requestFullscreen) {
				elem.requestFullscreen();
			}
		},
		zoomIn() { this.zoomLevel = Math.min(this.zoomLevel + 10, 200); },
		zoomOut() { this.zoomLevel = Math.max(this.zoomLevel - 10, 50); },
		isEmbeddablePDF(url) {
			if (!url) return false;
			const lowerUrl = url.toLowerCase();
			return lowerUrl.endsWith('.pdf') ||
				lowerUrl.includes('format=pdf') ||
				lowerUrl.includes('/pdf/') ||
				lowerUrl.includes('.pdf?');
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
			console.error('Erro ao carregar material premium:', error);
		} finally {
			this.loading = false;
		}
	}
}
</script>

<style scoped>
	.estudo-container { min-height: 100vh; background: #121826; }

	.ios-glass-card {
		background: rgba(255, 255, 255, 0.04) !important;
		backdrop-filter: blur(30px) saturate(180%);
		border: 1px solid rgba(255, 255, 255, 0.1);
	}

	.sticky-sidebar { position: sticky; top: 100px; }

	.book-cover-wrapper-premium {
		width: 140px;
		height: 200px;
		border-radius: 16px;
		overflow: hidden;
		box-shadow: 0 15px 35px rgba(0,0,0,0.6);
		background: #000;
	}

	.meta-grid { display: grid; grid-template-columns: 1fr; gap: 12px; }
	.meta-item { display: flex; flex-direction: column; }
	.meta-label { font-size: 10px; color: rgba(255, 255, 255, 0.4); text-transform: uppercase; letter-spacing: 1px; font-weight: 700; }
	.meta-value { color: white; font-weight: 600; font-size: 13px; }

	.reader-card {
		height: 800px;
		background: #000 !important;
		display: flex;
		flex-direction: column;
		border: 1px solid rgba(255, 255, 255, 0.08);
		overflow: hidden;
	}

	.reader-toolbar { background: #1a1a1a; border-bottom: 1px solid rgba(255,255,255,0.05); }

	.viewer-area { flex: 1; position: relative; }
	.safe-iframe { border: none; }

	.lock-screen {
		background: linear-gradient(135deg, rgba(0,0,0,0.95), rgba(0,184,212,0.1));
		height: 100%;
	}
	.lock-icon-wrap {
		width: 120px; height: 120px;
		background: rgba(0, 184, 212, 0.1);
		border-radius: 50%;
		display: flex; align-items: center; justify-content: center;
		border: 1px solid rgba(255,255,255,0.1);
	}

	.bg-black-opacity-40 { background: rgba(0,0,0,0.4); }
	.font-mono { font-family: 'DM Mono', monospace !important; }
	.citation-content { cursor: pointer; transition: color 0.2s; }
	.citation-content:hover { color: #00B8D4 !important; }

	.abstract-text { line-height: 1.8 !important; }
	.whitespace-pre { white-space: pre !important; }

	@media (max-width: 1264px) { .sticky-sidebar { position: relative; top: 0; } }
</style>
