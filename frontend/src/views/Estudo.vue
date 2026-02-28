<template>
	<div class="estudo-container mt-8" v-if="!loading && material">
		<v-container fluid class="pa-0 pa-sm-4">
			<v-row no-gutters>
				<!-- Sidebar: Book Info & Citation Tools -->
				<v-col cols="12" lg="3" class="pa-4" v-if="!zenMode">
					<v-card class="ios-glass-card pa-6 sticky-sidebar" rounded="xl">
						<v-btn class="ios-btn-back mb-6" variant="text" @click="goBack" density="comfortable">
							<v-icon class="mr-2">mdi-arrow-left</v-icon> Voltar
						</v-btn>
						<!-- ... rest of sidebar ... -->

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
								<span class="meta-label">Páginas</span>
								<span class="meta-value">{{ material.paginas || 'N/A' }} fls.</span>
							</div>
							<div class="meta-item">
								<span class="meta-label">Fonte</span>
								<span class="meta-value">{{ material.fonte || 'Digital' }}</span>
							</div>
							<div class="meta-item">
								<span class="meta-label">Categoria</span>
								<span class="meta-value">{{ material.categoria }}</span>
							</div>
							<div class="meta-item">
								<span class="meta-label">Identifier (ISBN/DOI)</span>
								<span class="meta-value text-truncate" :title="material.isbn || 'Não disponível'">{{ material.isbn || 'N/A' }}</span>
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

						<!-- Reading Options -->
						<div class="mt-6">
							<div class="text-caption text-uppercase font-weight-bold opacity-40 mb-3 tracking-widest text-center">Opções de Leitura</div>

							<v-btn
								v-if="material.pdf_url"
								block
								color="cyan"
								theme="dark"
								class="mb-3 text-none font-weight-black action-btn-premium"
								prepend-icon="mdi-book-open-page-variant"
								rounded="lg"
								@click="handleReadOnline"
								elevation="4"
							>
								Ler no Navegador
							</v-btn>

							<v-btn
								v-if="material.pdf_url"
								block
								variant="tonal"
								color="cyan"
								class="mb-3 text-none font-weight-bold"
								prepend-icon="mdi-cloud-download-outline"
								rounded="lg"
								:href="material.pdf_url"
								target="_blank"
								download
							>
								Baixar para Ler Offline
							</v-btn>

							<v-btn
								v-if="material.pdf_url"
								block
								variant="outlined"
								color="white"
								class="text-none"
								prepend-icon="mdi-open-in-new"
								rounded="lg"
								:href="material.pdf_url"
								target="_blank"
							>
								Ver no Site Original
							</v-btn>

							<!-- AI Summary Toggle -->
							<v-btn
								block
								color="amber-darken-2"
								class="mt-4 text-none font-weight-black ai-btn-premium"
								prepend-icon="mdi-auto-fix"
								rounded="lg"
								@click="gerarResumoIA"
								:loading="loadingResumo"
							>
								Gerar Resumo IA
							</v-btn>
						</div>
					</v-card>
				</v-col>

				<!-- Main Content: Reading View -->
				<v-col cols="12" lg="9" class="pa-4">
					<v-card class="reader-card" rounded="xl">
						<!-- Custom Reader Toolbar -->
						<div class="reader-toolbar d-flex align-center px-4 py-2">
							<div class="d-flex align-center">
								<v-btn
									:icon="zenMode ? 'mdi-close' : 'mdi-auto-fix'"
									size="small"
									variant="tonal"
									color="cyan"
									class="mr-3"
									@click="zenMode = !zenMode"
									:title="zenMode ? 'Sair do Modo Zen' : 'Ativar Modo Zen'"
								></v-btn>
								<span class="font-weight-bold text-caption text-uppercase tracking-widest text-white">{{ zenMode ? 'MODO ZEN ATIVO' : 'Visualizador Acadêmico' }}</span>
							</div>
							<v-spacer></v-spacer>
							<div class="d-flex align-center gap-2">
								<v-btn
									icon="mdi-note-edit-outline"
									size="small"
									variant="elevated"
									color="cyan"
									class="mr-2"
									@click="notesDrawer = !notesDrawer"
									title="Minhas Anotações"
								></v-btn>
								<v-btn
									icon="mdi-robot-happy"
									size="small"
									variant="elevated"
									color="amber-darken-2"
									class="mr-2 pulse-amber"
									@click="chatDrawer = !chatDrawer"
									title="Perguntar ao Livro (IA)"
								></v-btn>
								<v-btn icon="mdi-magnify-plus-outline" size="x-small" variant="text" color="white" @click="zoomIn"></v-btn>
								<v-btn icon="mdi-magnify-minus-outline" size="x-small" variant="text" color="white" @click="zoomOut"></v-btn>
								<v-divider vertical class="mx-2 opacity-20"></v-divider>
								<v-btn :icon="isFullscreen ? 'mdi-fullscreen-exit' : 'mdi-fullscreen'" size="small" variant="text" color="cyan" @click="toggleFullscreen" title="Foco total"></v-btn>
							</div>
						</div>

						<div class="viewer-area fill-height" ref="viewerContainer">
							<!-- AI Chat Overlay/Drawer -->
							<v-navigation-drawer
								v-model="chatDrawer"
								location="right"
								temporary
								width="350"
								class="ios-chat-drawer"
								overlay-opacity="0.3"
							>
								<div class="chat-header pa-4 d-flex align-center">
									<v-icon color="amber" class="mr-2">mdi-robot-happy</v-icon>
									<span class="font-weight-black text-white">Ask the Book (IA)</span>
									<v-spacer></v-spacer>
									<v-btn icon="mdi-close" size="small" variant="text" @click="chatDrawer = false"></v-btn>
								</div>

								<div class="chat-messages pa-4" ref="chatScroll">
									<div v-if="chatMessages.length === 0" class="empty-chat d-flex flex-column align-center justify-center opacity-40 py-10">
										<v-icon size="48">mdi-chat-question-outline</v-icon>
										<p class="text-caption mt-2 text-center">Tire suas dúvidas sobre o conteúdo deste livro agora!</p>
									</div>
									<div v-for="(msg, i) in chatMessages" :key="i" :class="['chat-bubble-wrap', msg.role]">
										<div class="chat-bubble">
											{{ msg.text }}
										</div>
									</div>
									<div v-if="sendingChat" class="chat-bubble-wrap ai">
										<div class="chat-bubble typing">
											<v-progress-circular indeterminate size="16" width="2" color="amber"></v-progress-circular>
											<span class="ml-2">Pensando...</span>
										</div>
									</div>
								</div>

								<div class="chat-input-wrap pa-4 bg-black">
									<v-text-field
										v-model="chatInput"
										placeholder="Faça uma pergunta..."
										variant="solo"
										density="comfortable"
										rounded="pill"
										hide-details
										@keyup.enter="enviarPerguntaIA"
										class="chat-field"
									>
										<template v-slot:append-inner>
											<v-btn icon="mdi-send" variant="text" size="small" color="amber" @click="enviarPerguntaIA" :disabled="!chatInput.trim() || sendingChat"></v-btn>
										</template>
									</v-text-field>
								</div>
							</v-navigation-drawer>

							<!-- Study Notes Drawer -->
							<v-navigation-drawer
								v-model="notesDrawer"
								location="right"
								temporary
								width="350"
								class="ios-notes-drawer"
							>
								<div class="notes-header pa-4 d-flex align-center">
									<v-icon color="cyan" class="mr-2">mdi-note-edit-outline</v-icon>
									<span class="font-weight-black text-white">Minhas Anotações</span>
									<v-spacer></v-spacer>
									<v-btn icon="mdi-close" size="small" variant="text" @click="notesDrawer = false"></v-btn>
								</div>

								<div class="pa-4">
									<v-textarea
										v-model="novaAnotacao"
										placeholder="Escreva uma nota sobre este livro..."
										variant="solo"
										bg-color="rgba(255,255,255,0.05)"
										rows="3"
										hide-details
										class="mb-3 rounded-lg"
									></v-textarea>
									<v-btn block color="cyan" class="text-none font-weight-bold" @click="salvarAnotacao" :loading="salvandoNota" :disabled="!novaAnotacao.trim()">
										Salvar Anotação
									</v-btn>
								</div>

								<v-divider class="opacity-10"></v-divider>

								<div class="notes-list pa-4">
									<div v-if="anotacoes.length === 0" class="empty-notes text-center opacity-30 py-10">
										<v-icon size="48">mdi-text-box-plus-outline</v-icon>
										<p class="text-caption mt-2">Nenhuma anotação ainda.</p>
									</div>
									<v-card v-for="nota in anotacoes" :key="nota.id" class="nota-card mb-4 pa-3" variant="tonal" color="cyan">
										<div class="d-flex justify-space-between align-start mb-2">
											<v-chip size="x-small" color="cyan" variant="flat">Pág. {{ nota.pagina || '?' }}</v-chip>
											<v-btn icon="mdi-delete" size="x-small" variant="text" color="red" @click="deletarAnotacao(nota.id)"></v-btn>
										</div>
										<div class="text-body-2 text-white">{{ nota.conteudo }}</div>
										<div class="text-right text-overline opacity-40 mt-1">{{ formatarData(nota.data_criacao) }}</div>
									</v-card>
								</div>
							</v-navigation-drawer>

							<template v-if="isAuthenticated">
								<!-- PDF Proxy Viewer -->
								<div v-if="isEmbeddablePDF(material.pdf_url)" class="pdf-container">
									<iframe
										:src="pdfProxyURL"
										class="safe-iframe"
										width="100%"
										height="100%"
										frameborder="0"
										allow="fullscreen"
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
					<v-card class="ios-glass-card mt-6 pa-8" rounded="xl" v-if="material.descricao || resumoIA">
						<div class="d-flex align-center justify-space-between mb-6">
							<h3 class="text-h5 font-weight-bold text-white d-flex align-center">
								<v-icon :color="resumoIA ? 'amber' : 'cyan'" class="mr-3">{{ resumoIA ? 'mdi-auto-fix' : 'mdi-text-long' }}</v-icon>
								{{ resumoIA ? 'Resumo Inteligente (IA)' : 'Resumo Acadêmico' }}
							</h3>
							<v-btn v-if="resumoIA" variant="text" size="small" color="white" class="opacity-40" @click="resumoIA = ''">Ver Original</v-btn>
						</div>

						<div v-if="resumoIA" class="resumo-ia-content fade-in">
							<div class="text-body-1 text-white opacity-90 leading-relaxed abstract-text mb-8 whitespace-pre-wrap">
								{{ resumoIA }}
							</div>
						</div>
						<div v-else class="text-body-1 text-white opacity-80 leading-relaxed abstract-text mb-8">
							{{ cleanText(material.descricao) }}
						</div>

						<v-divider class="opacity-10 mb-6"></v-divider>

						<!-- Detailed Specs Row -->
						<div class="d-flex flex-wrap gap-6">
							<div class="spec-card">
								<div class="spec-icon"><v-icon>mdi-file-document-outline</v-icon></div>
								<div class="spec-info">
									<div class="spec-label">Extensão</div>
									<div class="spec-value">{{ material.paginas || 'N/A' }} páginas</div>
								</div>
							</div>
							<div class="spec-card" v-if="material.isbn">
								<div class="spec-icon"><v-icon>mdi-barcode</v-icon></div>
								<div class="spec-info">
									<div class="spec-label">Identificação</div>
									<div class="spec-value font-weight-medium">{{ material.isbn }}</div>
								</div>
							</div>
							<div class="spec-card">
								<div class="spec-icon"><v-icon>mdi-earth</v-icon></div>
								<div class="spec-info">
									<div class="spec-label">Origem</div>
									<div class="spec-value">{{ material.fonte || 'Repositório Público' }}</div>
								</div>
							</div>
						</div>

						<!-- Tags Cloud -->
						<div class="mt-8" v-if="material.tags && material.tags.length">
							<div class="text-caption text-uppercase font-weight-bold opacity-40 mb-3 tracking-widest">Palavras-chave</div>
							<div class="d-flex flex-wrap gap-2">
								<v-chip
									v-for="tag in material.tags"
									:key="tag"
									size="small"
									variant="tonal"
									color="cyan"
									class="px-4"
								>
									{{ tag }}
								</v-chip>
							</div>
						</div>
					</v-card>
				</v-col>
			</v-row>

			<v-snackbar v-model="copySnack" :timeout="2000" color="cyan darken-2" location="bottom">
				Citação copiada com sucesso!
			</v-snackbar>
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
import EstudoService from '@/services/EstudoService'
import auth from '@/auth'

export default {
	name: 'EstudoPremium',
	data: () => ({
		material: null,
		loading: true,
		isAuthenticated: false,
		showCitationDialog: false,
		zoomLevel: 100,
		zenMode: false,
		isFullscreen: false,
		copySnack: false,
		chatDrawer: false,
		chatInput: '',
		chatMessages: [],
		sendingChat: false,
		resumoIA: '',
		loadingResumo: false,
		notesDrawer: false,
		novaAnotacao: '',
		anotacoes: [],
		salvandoNota: false
	}),
	computed: {
		pdfProxyURL() {
			if (!this.material || !this.material.pdf_url) return '';
			return MaterialService.getProxyPdfUrl(this.material.pdf_url);
		}
	},
	methods: {
		cleanText(text) {
			if (!text) return '';
			// Remove XML/HTML tags (like <jats:p>)
			let cleaned = text.replace(/<[^>]*>/g, '');
			// Remove extra spaces
			return cleaned.trim();
		},
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
			this.copySnack = true;
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
				lowerUrl.includes('.pdf?') ||
				lowerUrl.includes('sci_pdf') ||
				lowerUrl.includes('pdf_file');
		},
		goBack() {
			window.history.length > 1 ? this.$router.go(-1) : this.$router.push('/explorar');
		},
		async enviarPerguntaIA() {
			if (!this.chatInput.trim() || this.sendingChat) return;

			const msg = this.chatInput;
			this.chatMessages.push({ role: 'user', text: msg });
			this.chatInput = '';
			this.sendingChat = true;

			this.$nextTick(() => { this.scrollToBottom(); });

			try {
				const response = await MaterialService.perguntar(this.material.id, msg);
				this.chatMessages.push({ role: 'ai', text: response.data.resposta });
			} catch (error) {
				console.error('Erro ao perguntar IA:', error);
				this.chatMessages.push({ role: 'ai', text: 'Desculpe, tive um problema ao processar sua pergunta. Verifique sua chave de API.' });
			} finally {
				this.sendingChat = false;
				this.$nextTick(() => { this.scrollToBottom(); });
			}
		},
		async gerarResumoIA() {
			this.loadingResumo = true;
			try {
				const response = await MaterialService.obterResumo(this.material.id);
				this.resumoIA = response.data.resumo;
				// Scroll to abstract section
				document.querySelector('.abstract-text').scrollIntoView({ behavior: 'smooth' });
			} catch (error) {
				console.error('Erro ao gerar resumo IA:', error);
			} finally {
				this.loadingResumo = false;
			}
		},
		scrollToBottom() {
			const container = this.$refs.chatScroll;
			if (container) {
				container.scrollTop = container.scrollHeight;
			}
		},
		async carregarAnotacoes() {
			try {
				const response = await EstudoService.listarAnotacoes(auth.getUser().id, this.material.id);
				this.anotacoes = response.data;
			} catch (error) {
				console.error('Erro ao carregar anotações:', error);
			}
		},
		async salvarAnotacao() {
			this.salvandoNota = true;
			try {
				await EstudoService.criarAnotacao({
					usuario_id: auth.getUser().id,
					material_id: this.material.id,
					conteudo: this.novaAnotacao,
					pagina: 0, // Placeholder, can be improved to detect current page
					cor: 'cyan'
				});
				this.novaAnotacao = '';
				await this.carregarAnotacoes();
			} catch (error) {
				console.error('Erro ao salvar anotação:', error);
			} finally {
				this.salvandoNota = false;
			}
		},
		async deletarAnotacao(id) {
			try {
				await EstudoService.deletarAnotacao(id, auth.getUser().id);
				await this.carregarAnotacoes();
			} catch (error) {
				console.error('Erro ao deletar anotação:', error);
			}
		},
		formatarData(dateStr) {
			if (!dateStr) return '';
			const d = new Date(dateStr);
			return d.toLocaleDateString();
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
				this.carregarAnotacoes();
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
		height: 80vh;
		min-height: 500px;
		max-height: 800px;
		background: #000 !important;
		display: flex;
		flex-direction: column;
		border: 1px solid rgba(255, 255, 255, 0.08);
		overflow: hidden;
	}

	.reader-toolbar { background: #1a1a1a; border-bottom: 1px solid rgba(255,255,255,0.05); transition: background 0.3s ease; }
	.zen-mode .reader-toolbar { background: #000; }

	.viewer-area { flex: 1; position: relative; transition: all 0.5s ease; }
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

	.spec-card {
		display: flex;
		align-items: center;
		gap: 12px;
		background: rgba(255,255,255,0.03);
		padding: 12px 20px;
		border-radius: 12px;
		border: 1px solid rgba(255,255,255,0.05);
	}
	.spec-icon {
		width: 36px; height: 36px;
		border-radius: 10px;
		background: rgba(0,184,212,0.1);
		color: #00B8D4;
		display: flex; align-items: center; justify-content: center;
	}
	.spec-label { font-size: 10px; color: rgba(255,255,255,0.4); text-transform: uppercase; font-weight: 800; }
	.spec-value { font-size: 13px; color: white; font-weight: 600; }
	.tracking-widest { letter-spacing: 2px !important; }
	.gap-2 { gap: 8px; }
	.gap-4 { gap: 16px; }
	.gap-6 { gap: 24px; }

	@media (max-width: 1264px) { .sticky-sidebar { position: relative; top: 0; } }

	/* AI Styles */
	.ai-btn-premium {
		background: linear-gradient(135deg, #FFAB00, #FF6F00) !important;
		box-shadow: 0 4px 15px rgba(255, 171, 0, 0.3) !important;
		transition: all 0.3s ease;
	}
	.ai-btn-premium:hover {
		transform: translateY(-2px);
		box-shadow: 0 8px 25px rgba(255, 171, 0, 0.4) !important;
	}

	.ios-chat-drawer {
		background: rgba(18, 24, 38, 0.95) !important;
		backdrop-filter: blur(20px);
		border-left: 1px solid rgba(255,255,255,0.1);
		display: flex;
		flex-direction: column;
	}

	.chat-header {
		border-bottom: 1px solid rgba(255,255,255,0.1);
	}

	.chat-messages {
		flex: 1;
		overflow-y: auto;
		display: flex;
		flex-direction: column;
		gap: 12px;
	}

	.chat-bubble-wrap {
		display: flex;
		width: 100%;
	}

	.chat-bubble-wrap.user { justify-content: flex-end; }
	.chat-bubble-wrap.ai { justify-content: flex-start; }

	.chat-bubble {
		max-width: 85%;
		padding: 10px 14px;
		border-radius: 18px;
		font-size: 14px;
		line-height: 1.4;
	}

	.user .chat-bubble {
		background: var(--ios-cyan);
		color: white;
		border-bottom-right-radius: 4px;
	}

	.ai .chat-bubble {
		background: rgba(255,255,255,0.1);
		color: white;
		border-bottom-left-radius: 4px;
		border: 1px solid rgba(255,255,255,0.05);
	}

	.chat-bubble.typing {
		display: flex;
		align-items: center;
		color: rgba(255,255,255,0.5);
	}

	.chat-input-wrap {
		border-top: 1px solid rgba(255,255,255,0.1);
	}

	:deep(.chat-field .v-field) {
		background: rgba(255,255,255,0.05) !important;
		border: 1px solid rgba(255,255,255,0.1) !important;
	}

	.pulse-amber {
		animation: pulse-amber 2s infinite;
	}

	@keyframes pulse-amber {
		0% { box-shadow: 0 0 0 0 rgba(255, 171, 0, 0.4); }
		70% { box-shadow: 0 0 0 10px rgba(255, 171, 0, 0); }
		100% { box-shadow: 0 0 0 0 rgba(255, 171, 0, 0); }
	}

	.fade-in {
		animation: fadeIn 0.5s ease-out;
	}

	@keyframes fadeIn {
		from { opacity: 0; transform: translateY(10px); }
		to { opacity: 1; transform: translateY(0); }
	}

	.whitespace-pre-wrap { white-space: pre-wrap !important; }

	.ios-notes-drawer {
		background: rgba(18, 24, 38, 0.98) !important;
		backdrop-filter: blur(20px);
		border-left: 1px solid rgba(0, 184, 212, 0.2);
	}
	.notes-header { border-bottom: 1px solid rgba(255,255,255,0.05); }
	.nota-card {
		background: rgba(0, 184, 212, 0.05) !important;
		border: 1px solid rgba(0, 184, 212, 0.1);
		transition: all 0.3s ease;
	}
	.nota-card:hover { background: rgba(0, 184, 212, 0.1) !important; transform: scale(1.02); }
</style>
