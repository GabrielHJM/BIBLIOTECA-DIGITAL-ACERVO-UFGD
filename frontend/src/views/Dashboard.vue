<template>
	<div class="dashboard-container mt-8 position-relative">
		<v-snackbar v-model="snackbar" :timeout="3000" color="cyan darken-2" location="top right">
			{{ snackbarMsg }}
		</v-snackbar>
		<v-container fluid class="pa-4">
			<!-- Expert Statistics Dashboard Section -->
			<v-row v-if="stats && !loading" class="mb-8 px-4 fade-in">
				<!-- Meta Semanal Progress -->
				<v-col cols="12" md="4">
					<PremiumCard
						title="Meta Semanal"
						icon="mdi-flag-checkered"
						icon-color="amber"
						padding="20px"
						@click="dialogMeta = true"
						class="cursor-pointer"
					>
						<template #action>
							<v-btn icon="mdi-pencil" size="x-small" variant="text" color="amber"></v-btn>
						</template>

						<div class="stat-content">
							<div class="d-flex align-end">
								<span class="stat-value counter-anim">{{ stats.paginas_lidas_semana || 0 }}</span>
								<span class="text-caption ml-2 mb-2 opacity-60">/ {{ stats.meta_paginas_semana || 100 }} pág.</span>
							</div>
						</div>

						<div class="stat-progress-premium mt-4">
							<div class="progress-fill" :style="{ width: Math.min((stats.paginas_lidas_semana / stats.meta_paginas_semana) * 100, 100) + '%', background: 'linear-gradient(90deg, #FFAB00, #FFD600)' }"></div>
						</div>

						<div class="text-right mt-1 text-caption font-weight-bold" color="amber">
							{{ Math.round((stats.paginas_lidas_semana / stats.meta_paginas_semana) * 100) }}% da meta
						</div>
					</PremiumCard>
				</v-col>

				<v-col cols="12" md="3">
					<PremiumCard
						title="Páginas Lidas"
						icon="mdi-book-open-page-variant"
						icon-color="cyan"
						padding="24px"
					>
						<div class="stat-content">
							<span class="stat-value counter-anim">{{ stats.total_paginas || 0 }}</span>
							<p class="text-caption opacity-60 mb-0">Total acumulado</p>
						</div>
					</PremiumCard>
				</v-col>

				<v-col cols="12" md="5">
					<div class="stat-card badges premium-glow">
						<div class="d-flex align-center justify-space-between mb-4">
							<h3 class="stat-card-title d-flex align-center">
								<v-icon color="amber" class="mr-2">mdi-medal</v-icon>
								Minhas Conquistas
							</h3>
						</div>
						<div class="badges-scroll-wrap">
							<div v-if="stats.badges && stats.badges.length > 0" class="d-flex flex-wrap gap-2">
								<v-tooltip v-for="badge in stats.badges" :key="badge" location="top">
									<template v-slot:activator="{ props }">
										<v-chip
											v-bind="props"
											color="amber"
											variant="tonal"
											class="badge-chip px-4 py-2"
											prepend-icon="mdi-shield-check"
										>
											{{ badge }}
										</v-chip>
									</template>
									<span>Conquista Desbloqueada!</span>
								</v-tooltip>
							</div>
							<div v-else class="text-caption opacity-50 d-flex flex-column align-center justify-center py-4">
								<v-icon size="48" color="grey lighten-1">mdi-lock-outline</v-icon>
								<span>Continue lendo para ganhar medalhas!</span>
							</div>
						</div>
					</div>
				</v-col>
			</v-row>

			<!-- Study Hub Section (New) -->
			<v-row class="mb-8 px-4">
				<v-col cols="12">
					<h2 class="text-h5 font-weight-bold mb-4 d-flex align-center text-white">
						<v-icon color="amber" class="mr-2">mdi-school-outline</v-icon>
						Hub de Estudos
					</h2>
				</v-col>
				<v-col cols="12" md="6">
					<v-card class="study-hub-card flashcards premium-glow" @click="$router.push('/flashcards')" rounded="xl">
						<div class="d-flex align-center pa-6">
							<div class="hub-icon-wrap amber">
								<v-icon size="40" color="amber">mdi-cards-playing-outline</v-icon>
							</div>
							<div class="ml-6">
								<h3 class="text-h5 font-weight-black text-white">Praticar Flashcards</h3>
								<p class="text-body-2 text-white opacity-60">Memorize conceitos importantes com repetição espaçada.</p>
							</div>
							<v-spacer></v-spacer>
							<v-icon color="white" class="opacity-20">mdi-chevron-right</v-icon>
						</div>
					</v-card>
				</v-col>
				<v-col cols="12" md="6">
					<v-card class="study-hub-card notes premium-glow" rounded="xl">
						<div class="d-flex align-center pa-6">
							<div class="hub-icon-wrap cyan">
								<v-icon size="40" color="cyan">mdi-notebook-edit-outline</v-icon>
							</div>
							<div class="ml-6">
								<h3 class="text-h5 font-weight-black text-white">Minhas Anotações</h3>
								<p class="text-body-2 text-white opacity-60">Acesse todas as notas que você fez durante a leitura.</p>
							</div>
							<v-spacer></v-spacer>
							<v-chip color="cyan" size="small" class="font-weight-bold">Em breve</v-chip>
						</div>
					</v-card>
				</v-col>
			</v-row>

			<!-- Header Section -->
			<v-row align="center" class="mb-6 px-4">
				<v-col cols="12" sm="4" md="4" class="py-2">
					<h1 class="text-h4 text-sm-h3 font-weight-bold tracking-tight" :class="$vuetify.theme.current.dark ? 'text-white' : 'text-primary'">Recentes</h1>
				</v-col>
				<v-col cols="12" sm="5" md="5" class="py-2">
					<v-text-field
						v-model="searchQuery"
						prepend-inner-icon="mdi-magnify"
						placeholder="Pesquisar..."
						variant="solo"
						class="ios-search-bar"
						hide-details
						@keyup.enter="buscar"
					></v-text-field>
				</v-col>
				<v-col cols="12" sm="3" md="3" class="text-right py-2">
					<v-btn class="ios-filter-btn w-100 w-sm-auto" elevation="2" @click="buscar">
						<span>Buscar</span>
						<v-icon class="ml-2">mdi-magnify</v-icon>
					</v-btn>
				</v-col>
			</v-row>

			<!-- Histórico de Leitura Section -->
			<div v-if="historico.length > 0 && !loading" class="mb-10 mt-4">
				<h2 class="text-h5 font-weight-bold mb-4 px-4 px-sm-2 d-flex align-center" :class="$vuetify.theme.current.dark ? 'text-white' : 'text-grey-darken-3'">
					<v-icon color="cyan" class="mr-2">mdi-history</v-icon>
					Continuar Lendo
				</h2>
				<v-slide-group show-arrows class="pa-2 premium-slider">
					<v-slide-group-item
						v-for="(livro, index) in historico"
						:key="'hist-'+index"
					>
						<v-card
							class="ios-history-card ma-2 premium-shadow-hover"
							elevation="0"
							@click="$router.push('/estudo/' + livro.id)"
						>
							<div class="history-cover-wrap">
								<v-img
									:src="livro.capa_url || 'https://images.unsplash.com/photo-1544947950-fa07a98d237f?q=80&w=200'"
									cover
									height="160"
									class="history-img"
								></v-img>
								<div class="history-overlay">
									<v-icon color="white" size="32">mdi-play-circle-outline</v-icon>
								</div>
							</div>
							<div class="pa-3">
								<div class="text-caption font-weight-bold history-title-premium" :class="$vuetify.theme.current.dark ? 'text-white' : 'text-grey-darken-3'">{{ livro.titulo }}</div>
								<div class="opacity-40 ultra-caption mt-1" :class="$vuetify.theme.current.dark ? 'text-white' : 'text-grey-darken-2'">{{ livro.autor }}</div>
							</div>
						</v-card>
					</v-slide-group-item>
				</v-slide-group>
			</div>

			<!-- Recomendações Inteligentes (Topic 2 IA) -->
			<div v-if="recomendacoes.length > 0 && !loading" class="mb-10 mt-4">
				<h2 class="text-h5 font-weight-bold mb-4 px-4 px-sm-2 d-flex align-center" :class="$vuetify.theme.current.dark ? 'text-white' : 'text-grey-darken-3'">
					<v-icon color="amber" class="mr-2">mdi-auto-fix</v-icon>
					Recomendado para Você (IA)
				</h2>
				<v-slide-group show-arrows class="pa-2">
					<v-slide-group-item
						v-for="(livro, index) in recomendacoes"
						:key="'rec-'+index"
					>
						<v-card
							class="ios-history-card ma-2"
							elevation="4"
							@click="$router.push('/estudo/' + livro.id)"
						>
							<v-img
								:src="livro.capa_url || 'https://images.unsplash.com/photo-1544947950-fa07a98d237f?q=80&w=200'"
								cover
								height="150"
							></v-img>
							<div class="pa-2">
								<div class="text-caption font-weight-bold history-title" :class="$vuetify.theme.current.dark ? 'text-white' : 'text-grey-darken-3'">{{ livro.titulo }}</div>
								<div class="opacity-60 ultra-caption" :class="$vuetify.theme.current.dark ? 'text-white' : 'text-grey-darken-2'">{{ livro.autor }}</div>
							</div>
						</v-card>
					</v-slide-group-item>
				</v-slide-group>
			</div>

			<!-- Meus Favoritos Section -->
			<div v-if="favoritos.length > 0 && !loading" class="mb-10 mt-4">
				<h2 class="text-h5 font-weight-bold mb-4 px-4 px-sm-2 d-flex align-center" :class="$vuetify.theme.current.dark ? 'text-white' : 'text-grey-darken-3'">
					<v-icon color="pink" class="mr-2">mdi-heart</v-icon>
					Meus Favoritos
				</h2>
				<v-slide-group show-arrows class="pa-2">
					<v-slide-group-item
						v-for="(livro, index) in favoritos"
						:key="'fav-'+index"
					>
						<v-card
							class="ios-history-card ma-2"
							elevation="4"
							@click="$router.push('/estudo/' + livro.id)"
						>
							<v-img
								:src="livro.capa_url || 'https://images.unsplash.com/photo-1544947950-fa07a98d237f?q=80&w=200'"
								cover
								height="150"
							></v-img>
							<div class="pa-2">
								<div class="text-caption font-weight-bold history-title" :class="$vuetify.theme.current.dark ? 'text-white' : 'text-grey-darken-3'">{{ livro.titulo }}</div>
								<div class="opacity-60 ultra-caption" :class="$vuetify.theme.current.dark ? 'text-white' : 'text-grey-darken-2'">{{ livro.autor }}</div>
							</div>
						</v-card>
					</v-slide-group-item>
				</v-slide-group>
			</div>

			<!-- Content Grid -->
			<v-row class="px-2" v-if="!loading">
				<v-col
					v-for="(livro, index) in livros"
					:key="livro.id"
					cols="12"
					md="6"
					lg="4"
					class="pa-4"
				>
					<BookCard
						:book="livro"
						:is-favorited="isFavorited(livro.id)"
						:animation-delay="index * 100"
						@toggle-favorite="toggleFavorite"
						@share="shareBook(livro.id)"
					/>
				</v-col>
			</v-row>

			<v-row v-else justify="center" class="mt-12">
				<v-progress-circular indeterminate color="cyan" size="64"></v-progress-circular>
			</v-row>
		</v-container>

		<!-- Meta Semanal Dialog -->
		<v-dialog v-model="dialogMeta" max-width="400">
			<v-card class="ios-dialog-premium pa-4">
				<v-card-title class="text-h5 font-weight-bold d-flex align-center">
					<v-icon color="amber" class="mr-3">mdi-flag-checkered</v-icon>
					Meta de Leitura
				</v-card-title>
				<v-card-text>
					<p class="mb-6 opacity-70">Defina quantas páginas você pretende ler por semana para manter seu ritmo de estudos.</p>
					<v-slider
						v-model="novaMeta"
						min="10"
						max="1000"
						step="10"
						color="amber"
						thumb-label="always"
						class="mt-8"
					></v-slider>
					<div class="text-center text-h6 font-weight-bold text-amber">
						{{ novaMeta }} páginas / semana
					</div>
				</v-card-text>
				<v-card-actions class="mt-4">
					<v-spacer></v-spacer>
					<v-btn variant="text" @click="dialogMeta = false">Cancelar</v-btn>
					<v-btn color="amber" variant="flat" class="rounded-lg px-6" @click="salvarMeta" :loading="salvandoMeta">
						Salvar Meta
					</v-btn>
				</v-card-actions>
			</v-card>
		</v-dialog>
	</div>
</template>

<script>
import MaterialService from '@/services/MaterialService'
import auth from '@/auth'
import BookCard from '@/components/BookCard.vue'
import PremiumCard from '@/components/PremiumCard.vue'

export default {
	name: 'DashboardPage',
	components: {
		BookCard,
		PremiumCard
	},
	data: () => ({
		livros: [],
		historico: [],
		favoritos: [],
		recomendacoes: [],
		searchQuery: '',
		loading: false,
		user: {},
		stats: null,
		snackbar: false,
		snackbarMsg: '',
		dialogMeta: false,
		novaMeta: 100,
		salvandoMeta: false
	}),
	computed: {
		topCategorias() {
			if (!this.stats || !this.stats.categorias) return {};
			// Retornar as top 3
			return Object.fromEntries(
				Object.entries(this.stats.categorias).slice(0, 3)
			);
		}
	},
	created() {
		this.user = auth.getUser()
		this.buscar()
		if (this.user && this.user.id) {
			this.buscarHistorico()
			this.buscarFavoritos()
			this.buscarEstatisticas()
			this.buscarRecomendacoes()
		}
	},
	methods: {
		async buscarEstatisticas() {
			try {
				const response = await MaterialService.obterEstatisticas(this.user.id)
				this.stats = response.data
			} catch (error) {
				console.error('Erro ao buscar estatísticas:', error)
			}
		},
		async buscar() {
			this.loading = true
			try {
				let response;
				if (this.searchQuery) {
					response = await MaterialService.pesquisar(this.searchQuery)
				} else {
					response = await MaterialService.listar(9, 0)
				}
				this.livros = response.data || []
			} catch (error) {
				console.error('Erro ao buscar materiais:', error)
			} finally {
				this.loading = false
			}
		},
		async buscarHistorico() {
			try {
				const response = await MaterialService.listarHistorico(this.user.id)
				this.historico = response.data || []
			} catch (error) {
				console.error('Erro ao buscar histórico:', error)
			}
		},
		async buscarFavoritos() {
			try {
				const response = await MaterialService.listarFavoritos(this.user.id)
				this.favoritos = response.data || []
			} catch (error) {
				console.error('Erro ao buscar favoritos:', error)
			}
		},
		async buscarRecomendacoes() {
			try {
				const response = await MaterialService.obterRecomendacoes(this.user.id, 6)
				this.recomendacoes = response.data || []
			} catch (error) {
				console.error('Erro ao buscar recomendações inteligentes:', error)
			}
		},
		isFavorited(materialId) {
			return this.favoritos.some(f => f.id === materialId)
		},
		async toggleFavorite(livro) {
			if (!this.user || !this.user.id) return

			const currentlyFavorited = this.isFavorited(livro.id)
			const novoStatus = !currentlyFavorited

			try {
				await MaterialService.favoritar(this.user.id, livro.id, novoStatus)
				if (novoStatus) {
					this.favoritos.unshift(livro)
					this.showSnackbar('Adicionado aos favoritos!')
				} else {
					this.favoritos = this.favoritos.filter(f => f.id !== livro.id)
					this.showSnackbar('Removido dos favoritos.')
				}
			} catch (err) {
				console.error('Erro ao favoritar/desfavoritar:', err)
				this.showSnackbar('Erro ao atualizar favorito.')
			}
		},
		shareBook(id) {
			const link = `${window.location.origin}/estudo/${id}`
			navigator.clipboard.writeText(link).then(() => {
				this.showSnackbar('Link copiado para a área de transferência!')
			}).catch(err => {
				console.error('Erro ao copiar link: ', err)
				this.showSnackbar('Não foi possível copiar o link.')
			})
		},
		showSnackbar(msg) {
			this.snackbarMsg = msg
			this.snackbar = true
		},
		async salvarMeta() {
			this.salvandoMeta = true
			try {
				const UsuarioService = (await import('@/services/UsuarioService')).default
				await UsuarioService.atualizarMeta(this.user.id, this.novaMeta)
				this.stats.meta_paginas_semana = this.novaMeta
				this.dialogMeta = false
				this.showSnackbar('Meta atualizada com sucesso!')
			} catch (error) {
				console.error('Erro ao salvar meta:', error)
				this.showSnackbar('Erro ao salvar meta.')
			} finally {
				this.salvandoMeta = false
			}
		},
		logout() {
			auth.logout()
			this.$router.push('/login')
		}
	},
	watch: {
		stats(newVal) {
			if (newVal && newVal.meta_paginas_semana) {
				this.novaMeta = newVal.meta_paginas_semana
			}
		}
	}
}
</script>

<style scoped>
	.dashboard-container {
		min-height: 100vh;
		padding-bottom: 80px;
	}

	.tracking-tight {
		letter-spacing: -1.5px !important;
	}

	.ios-filter-btn {
		background: rgba(255, 255, 255, 0.15) !important;
		backdrop-filter: blur(10px);
		color: white !important;
		border-radius: 12px !important;
		text-transform: none !important;
		border: 1px solid rgba(255, 255, 255, 0.1);
	}

	.ios-search-bar {
		background: rgba(255, 255, 255, 0.1) !important;
		border-radius: 16px !important;
	}

	:deep(.ios-search-bar .v-field) {
		background: transparent !important;
		border-radius: 16px !important;
		color: white !important;
	}

	/* Item Card Style */
	.ios-item-card {
		background: rgba(var(--v-theme-surface), 0.7) !important;
		backdrop-filter: blur(20px) saturate(180%);
		-webkit-backdrop-filter: blur(20px) saturate(180%);
		border-radius: 28px !important;
		border: 1px solid rgba(var(--v-border-color), 0.1);
		overflow: hidden;
		transition: all 0.5s var(--spring-easing);
		opacity: 0;
		animation: ios-reveal 0.7s var(--spring-easing) forwards;
	}

	.ios-item-card:hover {
		transform: translateY(-10px) scale(1.02);
		background: rgba(var(--v-theme-surface), 0.9) !important;
		box-shadow: 0 30px 60px rgba(0,0,0,0.15) !important;
		border-color: var(--ios-cyan) !important;
	}

	/* Expert Stats Cards */
	.stat-card {
		background: rgba(var(--v-theme-surface), 0.4);
		backdrop-filter: blur(40px) saturate(150%);
		border-radius: 28px;
		padding: 24px;
		height: 100%;
		border: 1px solid rgba(255, 255, 255, 0.08);
		transition: all 0.5s cubic-bezier(0.16, 1, 0.3, 1);
		position: relative;
		overflow: hidden;
	}

	.stat-card::after {
		content: '';
		position: absolute;
		top: 0; right: 0; width: 100px; height: 100px;
		background: radial-gradient(circle, rgba(0, 184, 212, 0.05) 0%, transparent 70%);
		pointer-events: none;
	}

	.premium-glow:hover {
		background: rgba(var(--v-theme-surface), 0.6);
		border-color: rgba(0, 184, 212, 0.3);
		transform: translateY(-8px) scale(1.01);
		box-shadow: 0 20px 40px rgba(0,0,0,0.2) !important;
	}

	.stat-icon-wrap {
		width: 52px;
		height: 52px;
		border-radius: 14px;
		display: flex;
		align-items: center;
		justify-content: center;
		margin-bottom: 24px;
		transition: transform 0.4s ease;
	}
	.stat-card:hover .stat-icon-wrap { transform: scale(1.1) rotate(-5deg); }

	.stat-progress-premium {
		height: 6px;
		background: rgba(255,255,255,0.05);
		border-radius: 10px;
		overflow: hidden;
	}
	.progress-fill {
		height: 100%;
		background: linear-gradient(90deg, #00B8D4, #007A99);
		border-radius: 10px;
		box-shadow: 0 0 15px rgba(0, 184, 212, 0.3);
	}

	.premium-bar :deep(.v-progress-linear__determinate) {
		background: linear-gradient(90deg, #007A99, #00B8D4) !important;
		border-radius: 10px !important;
	}

	.history-cover-wrap {
		position: relative;
		overflow: hidden;
		border-radius: 16px 16px 0 0;
	}
	.history-overlay {
		position: absolute;
		top: 0; left: 0; width: 100%; height: 100%;
		background: rgba(0,0,0,0.4);
		display: flex; align-items: center; justify-content: center;
		opacity: 0;
		transition: opacity 0.3s ease;
	}
	.ios-history-card:hover .history-overlay { opacity: 1; }
	.history-img { transition: transform 0.6s cubic-bezier(0.16, 1, 0.3, 1); }
	.ios-history-card:hover .history-img { transform: scale(1.15); }

	.history-title-premium {
		font-size: 13px !important;
		line-height: 1.3;
		height: 2.6em;
		overflow: hidden;
		display: -webkit-box;
		-webkit-line-clamp: 2;
		-webkit-box-orient: vertical;
		font-weight: 700 !important;
	}

	.fade-in {
		animation: fadeIn 1s ease-out;
	}

	@keyframes fadeIn {
		from { opacity: 0; transform: translateY(20px); }
		to { opacity: 1; transform: translateY(0); }
	}

	.ios-history-card {
		width: 130px;
		background: rgba(var(--v-theme-surface), 0.8) !important;
		backdrop-filter: blur(10px);
		border-radius: 16px !important;
		border: 1px solid rgba(128, 128, 128, 0.2);
		overflow: hidden;
		cursor: pointer;
		transition: transform 0.3s ease;
	}

	.ios-history-card:hover {
		transform: scale(1.05);
		background: rgba(var(--v-theme-surface), 0.95) !important;
	}

	.history-title {
		white-space: nowrap;
		overflow: hidden;
		text-overflow: ellipsis;
	}

	.ultra-caption {
		font-size: 10px;
	}

	.book-cover-wrapper {
		width: 100%;
		height: 180px;
		border-radius: 12px;
		overflow: hidden;
		background: rgba(0,0,0,0.2);
		display: flex;
		align-items: center;
		justify-content: center;
		box-shadow: 0 4px 12px rgba(0,0,0,0.2);
	}

	.book-cover {
		max-width: 100%;
		max-height: 100%;
		object-fit: contain;
		transition: transform 0.5s ease;
	}

	.ios-item-card:hover .book-cover {
		transform: scale(1.1) rotate(2deg);
	}

	.item-title {
		color: var(--v-theme-on-surface);
		font-size: 1.1rem;
		font-weight: 700;
		line-height: 1.2;
		height: 2.4em;
		overflow: hidden;
		display: -webkit-box;
		-webkit-line-clamp: 2;
		-webkit-box-orient: vertical;
	}

	.item-details {
		color: var(--v-theme-on-surface);
		opacity: 0.7;
		font-size: 13px;
	}

	.item-details p {
		margin: 2px 0;
	}

	.ios-btn-open {
		background: var(--ios-cyan) !important;
		color: white !important;
		border-radius: 10px !important;
		text-transform: none !important;
		font-weight: 600 !important;
		margin-left: 8px;
	}

	.share-actions {
		display: flex;
		align-items: center;
	}

	/* Animations */
	@keyframes ios-reveal {
		from {
			opacity: 0;
			transform: scale(0.9) translateY(30px);
		}
		to {
			opacity: 1;
			transform: scale(1) translateY(0);
		}
	}

	@media (max-width: 600px) {
		.dashboard-container {
			padding-bottom: 40px;
		}

		.tracking-tight {
			letter-spacing: -1px !important;
		}

		.ios-item-card {
			border-radius: 20px !important;
		}

		.item-title {
			font-size: 1rem;
		}

		.book-cover-wrapper {
			height: 140px;
		}

		.item-details {
			font-size: 12px;
		}
	}

	:deep(.v-rating__item) {
		padding: 0 !important;
	}

	/* Gamification Styles */
	.badges-scroll-wrap {
		max-height: 120px;
		overflow-y: auto;
		scrollbar-width: none; /* Firefox */
	}
	.badges-scroll-wrap::-webkit-scrollbar {
		display: none; /* Safari and Chrome */
	}

	.badge-chip {
		font-weight: 700 !important;
		font-size: 11px !important;
		letter-spacing: 0.5px;
		transition: all 0.3s ease;
		background: rgba(255, 171, 0, 0.15) !important;
		border: 1px solid rgba(255, 171, 0, 0.2);
	}

	.badge-chip:hover {
		transform: scale(1.05) translateY(-2px);
		background: rgba(255, 171, 0, 0.25) !important;
		box-shadow: 0 4px 12px rgba(255, 171, 0, 0.2);
	}

	.ios-dialog-premium {
		background: rgba(var(--v-theme-surface), 0.8) !important;
		backdrop-filter: blur(25px) saturate(180%);
		border-radius: 32px !important;
		border: 1px solid rgba(255, 255, 255, 0.1);
	}

	.stat-card.goal {
		cursor: pointer;
	}

	.stat-card.badges {
		background: rgba(var(--v-theme-surface), 0.2);
		border-color: rgba(255, 171, 0, 0.2);
	}

	.gap-2 { gap: 8px; }

	/* Study Hub */
	.study-hub-card {
		background: rgba(var(--v-theme-surface), 0.5) !important;
		backdrop-filter: blur(20px);
		border: 1px solid rgba(255, 255, 255, 0.1);
		cursor: pointer;
		transition: all 0.4s var(--spring-easing);
	}
	.study-hub-card.flashcards:hover { border-color: rgba(255, 171, 0, 0.4); }
	.study-hub-card.notes:hover { border-color: rgba(0, 184, 212, 0.4); }

	.hub-icon-wrap {
		width: 70px; height: 70px;
		border-radius: 20px;
		display: flex; align-items: center; justify-content: center;
	}
	.hub-icon-wrap.amber { background: rgba(255, 171, 0, 0.1); }
	.hub-icon-wrap.cyan { background: rgba(0, 184, 212, 0.1); }
</style>
