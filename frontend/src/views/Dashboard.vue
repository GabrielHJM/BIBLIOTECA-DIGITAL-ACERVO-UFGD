<template>
	<div class="dashboard-container position-relative">
		<v-container fluid class="pa-0 pa-md-2">
			<!-- Expert Statistics Dashboard Section -->
			<v-row v-if="stats && !loading" class="mb-4 px-4 fade-in">
				<v-col cols="12" md="12">
					<div class="premium-card-blur pa-5 h-100">
						<div class="d-flex align-center justify-space-between mb-4">
							<h3 class="text-overline font-weight-bold opacity-60 d-flex align-center">
								<v-icon color="amber" class="mr-2" size="18">mdi-medal</v-icon>
								Minhas Conquistas
							</h3>
						</div>
						<div class="badges-scroll-wrap">
							<div v-if="stats.badges && stats.badges.length > 0" class="d-flex flex-wrap gap-2">
								<v-chip
									v-for="badge in stats.badges"
									:key="badge"
									color="amber"
									variant="tonal"
									class="badge-chip-premium px-4 py-2"
									prepend-icon="mdi-shield-check"
								>
									{{ badge }}
								</v-chip>
							</div>
							<div v-else class="text-caption opacity-30 d-flex flex-column align-center justify-center py-4">
								<v-icon size="40" color="white" class="opacity-20 mb-2">mdi-lock-outline</v-icon>
								<span>Leia para desbloquear medalhas!</span>
							</div>
						</div>
					</div>
				</v-col>
			</v-row>

			<!-- Study Hub Section -->
			<v-row class="mb-6 px-4">
				<v-col cols="12">
					<h2 class="text-h5 font-weight-bold mb-4 d-flex align-center" :class="isDarkTheme ? 'text-white' : 'text-primary-darken-4'">
						<v-icon color="amber" class="mr-2">mdi-school-outline</v-icon>
						Hub de Estudos
					</h2>
				</v-col>
				<v-col cols="12">
					<v-card class="study-hub-card notes premium-glow reader-rounding cursor-pointer" @click="$router.push('/anotacoes')">
						<div class="d-flex align-center pa-6">
							<div class="hub-icon-wrap cyan">
								<v-icon size="40" color="cyan">mdi-notebook-edit-outline</v-icon>
							</div>
							<div class="ml-6">
								<h3 class="text-h5 font-weight-black" :class="isDarkTheme ? 'text-white' : 'text-primary-darken-4'">Minhas Anotações</h3>
								<p class="text-body-2 opacity-60" :class="isDarkTheme ? 'text-white' : 'text-primary-darken-2'">Acesse todas as notas que você fez durante a leitura.</p>
							</div>
							<v-spacer></v-spacer>
							<v-icon color="white" class="opacity-20">mdi-chevron-right</v-icon>
						</div>
					</v-card>
				</v-col>
			</v-row>

			<!-- Insights Section (New) -->
			<v-row v-if="stats && !loading" class="mb-6 px-4 fade-in">
				<v-col cols="12">
					<h2 class="text-h5 font-weight-bold mb-4 d-flex align-center" :class="isDarkTheme ? 'text-white' : 'text-primary-darken-4'">
						<v-icon color="primary" class="mr-2">mdi-chart-areaspline</v-icon>
						Insights de Estudo
					</h2>
				</v-col>

				<!-- Top Categorias -->
				<v-col cols="12" md="6">
					<PremiumCard title="Áreas Mais Estudadas" icon="mdi-shape-outline" icon-color="primary">
						<div v-if="Object.keys(stats.categorias).length > 0" class="mt-4">
							<div v-for="(count, cat) in stats.categorias" :key="cat" class="mb-4">
								<div class="d-flex justify-space-between mb-1">
									<span class="text-body-2 font-weight-bold" :class="isDarkTheme ? 'text-white' : 'text-primary-darken-4'">{{ cat }}</span>
									<span class="text-caption opacity-60">{{ count }} sessões</span>
								</div>
								<div class="stat-progress-premium">
									<div class="progress-fill" :style="{ width: Math.min((count / maxCategoryCount) * 100, 100) + '%', background: 'linear-gradient(90deg, #007AFF, #00C1FF)' }"></div>
								</div>
							</div>
						</div>
						<div v-else class="text-center py-6 opacity-40">
							<v-icon size="40">mdi-chart-bubble</v-icon>
							<p>Nenhum dado de categoria ainda</p>
						</div>
					</PremiumCard>
				</v-col>

				<!-- Top Livros -->
				<v-col cols="12" md="6">
					<PremiumCard title="Livros Mais Vistos" icon="mdi-trophy-outline" icon-color="amber">
						<v-list class="pa-0 bg-transparent">
							<v-list-item
								v-for="(livro, idx) in stats.top_livros"
								:key="livro.id"
								class="px-0 mb-3 top-book-item"
								@click="$router.push('/estudo/' + livro.id)"
							>
								<template v-slot:prepend>
									<div class="top-book-rank mr-4" :class="'rank-' + (idx + 1)">
										{{ idx + 1 }}
									</div>
									<v-avatar size="48" rounded="lg" class="mr-4">
									<div class="mini-premium-cover-ios-stats">
										<div class="mesh-gradient-mini"></div>
										<div class="glass-overlay-mini">
											<v-icon size="14" color="rgba(255,255,255,0.7)">mdi-trophy-variant-outline</v-icon>
										</div>
									</div>
									</v-avatar>
								</template>
								<v-list-item-title class="text-body-2 font-weight-bold" :class="isDarkTheme ? 'text-white' : 'text-primary-darken-4'">{{ livro.titulo }}</v-list-item-title>
								<v-list-item-subtitle class="text-caption opacity-60">{{ livro.qtd }} acessos</v-list-item-subtitle>
								<template v-slot:append>
									<v-icon size="18" class="opacity-20">mdi-chevron-right</v-icon>
								</template>
							</v-list-item>
							<div v-if="!stats.top_livros || stats.top_livros.length === 0" class="text-center py-6 opacity-40">
								<v-icon size="40">mdi-book-off-outline</v-icon>
								<p>Inicie uma leitura para ver métricas</p>
							</div>
						</v-list>
					</PremiumCard>
				</v-col>
			</v-row>

			<!-- Header Section -->
			<v-row align="center" class="mb-6 px-4">
				<v-col cols="12" sm="4" md="4" class="py-2">
					<h1 class="text-h4 text-sm-h3 font-weight-bold tracking-tight" :class="isDarkTheme ? 'text-white' : 'text-primary'">Recentes</h1>
				</v-col>
				<v-col cols="12" sm="5" md="5" class="py-2 position-relative">
					<v-text-field
						v-model="searchQuery"
						prepend-inner-icon="mdi-magnify"
						placeholder="Pesquisar..."
						variant="solo"
						class="ios-search-bar"
						hide-details
						@keyup.enter="handleEnter"
						@keydown.down="navigateSuggestions(1)"
						@keydown.up="navigateSuggestions(-1)"
						@keydown.esc="showSuggestions = false"
						@input="debouncedGetSuggestions"
					></v-text-field>

					<!-- Apple-style Suggestions Dropdown -->
					<div v-if="showSuggestions && suggestions.length > 0" class="suggestions-dropdown-premium fade-in-fast">
						<div
							v-for="(suggestion, idx) in suggestions"
							:key="suggestion.id"
							class="suggestion-item"
							:class="{ 'active': suggestionIndex === idx }"
							@click="selectSuggestion(suggestion)"
							@mouseover="suggestionIndex = idx"
						>
							<v-icon size="18" class="mr-3 opacity-40">{{ getBookIcon(suggestion.categoria, suggestion.titulo) }}</v-icon>
							<div class="suggestion-content">
								<div class="suggestion-title">{{ suggestion.titulo }}</div>
								<div class="suggestion-meta">{{ suggestion.autor }} • {{ suggestion.categoria }}</div>
							</div>
							<v-spacer></v-spacer>
							<v-icon size="14" class="opacity-20">mdi-arrow-top-left</v-icon>
						</div>
					</div>
				</v-col>
				<v-col cols="12" sm="3" md="3" class="text-right py-2">
					<v-btn class="ios-filter-btn w-100 w-sm-auto" elevation="2" @click="buscar(true)">
						<span>Buscar</span>
						<v-icon class="ml-2">mdi-magnify</v-icon>
					</v-btn>
				</v-col>
			</v-row>

			<!-- Histórico de Leitura Section -->
			<div v-if="historico.length > 0 && !loading" class="mb-6 mt-2">
				<h2 class="text-h5 font-weight-bold mb-4 px-4 px-sm-2 d-flex align-center" :class="isDarkTheme ? 'text-white' : 'text-grey-darken-3'">
					<v-icon color="primary" class="mr-2">mdi-history</v-icon>
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
							<div class="mini-premium-cover-ios-dashboard">
								<div class="mesh-gradient-mini"></div>
								<div class="glass-overlay-mini">
									<div class="cover-content-ios-mini">
										<v-icon class="cover-icon-ios mb-1" size="24" color="rgba(255,255,255,0.7)">{{ getBookIcon(livro.categoria, livro.titulo) }}</v-icon>
										<div class="cover-title-ios-mini">{{ livro.titulo }}</div>
										<div class="cover-summary-ios-mini">{{ livro.resumo }}</div>
									</div>
								</div>
								<div class="glass-shine"></div>
							</div>
							<div class="pa-3">
								<div class="text-caption font-weight-bold history-title-premium" :class="isDarkTheme ? 'text-white' : 'text-grey-darken-3'">{{ livro.titulo }}</div>
							</div>
						</v-card>
					</v-slide-group-item>
				</v-slide-group>
			</div>



			<!-- Meus Favoritos Section -->
			<div v-if="favoritos.length > 0 && !loading" class="mb-6 mt-2">
				<h2 class="text-h5 font-weight-bold mb-4 px-4 px-sm-2 d-flex align-center" :class="isDarkTheme ? 'text-white' : 'text-grey-darken-3'">
					<v-icon color="pink" class="mr-2">mdi-heart</v-icon>
					Meus Favoritos
					<v-spacer></v-spacer>
					<v-btn variant="text" color="pink" size="small" class="text-none" @click="$router.push('/favoritos')">Ver Tudo</v-btn>
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
							<div class="mini-premium-cover-ios-dashboard">
								<div class="mesh-gradient-mini"></div>
								<div class="glass-overlay-mini">
									<div class="cover-content-ios-mini">
										<v-icon class="cover-icon-ios mb-1" size="24" color="rgba(255,255,255,0.7)">{{ getBookIcon(livro.categoria, livro.titulo) }}</v-icon>
										<div class="cover-title-ios-mini">{{ livro.titulo }}</div>
										<div class="cover-summary-ios-mini">{{ livro.resumo }}</div>
									</div>
								</div>
								<div class="glass-shine"></div>
							</div>
							<div class="pa-2">
								<div class="text-caption font-weight-bold history-title" :class="isDarkTheme ? 'text-white' : 'text-grey-darken-3'">{{ livro.titulo }}</div>
							</div>
						</v-card>
					</v-slide-group-item>
				</v-slide-group>
			</div>

			<!-- Content Grid -->
			<div>
				<v-row class="px-2">
					<v-col
						v-for="(livro, index) in livros"
						:key="livro.id"
						cols="12"
						sm="12"
						md="6"
						lg="6"
						class="pa-4"
					>
						<BookCard
							:book="livro"
							:is-favorited="isFavorited(livro.id)"
							:animation-delay="index * 50"
							@toggle-favorite="toggleFavorite"
							@share="shareBook(livro.id)"
						/>
					</v-col>
				</v-row>

				<v-row v-if="loading" justify="center" class="mt-6 mb-8">
					<v-progress-circular indeterminate color="primary" size="48"></v-progress-circular>
				</v-row>

				<!-- Intersection Observer Sentinel -->
				<div ref="loadMoreSentinel" class="load-more-sentinel-modern"></div>
			</div>
		</v-container>

	</div>
</template>

<script>
import MaterialService from '@/services/MaterialService'
import auth from '@/auth'
import PremiumCard from '@/components/PremiumCard.vue'
import debounce from 'lodash.debounce'
import { computed } from 'vue'
import { useTheme } from 'vuetify'

export default {
	name: 'DashboardPage',
	inject: ['notify', 'fetchGlobalFavorites', 'getGlobalFavorites', 'globalFavorites'],
	components: {
		PremiumCard
	},
	setup() {
		const theme = useTheme();
		const isDarkTheme = computed(() => theme.global.current.value.dark);
		return { isDarkTheme };
	},
	data: () => ({
		livros: [],
		historico: [],
		searchQuery: '',
		loading: false,
		user: {},
		stats: null,
		limit: 16,
		offset: 0,
		observer: null,
		suggestions: [],
		showSuggestions: false,
		suggestionIndex: -1
	}),
	computed: {
		favoritos() {
			return this.globalFavorites?.list || [];
		},
		topCategorias() {
			if (!this.stats || !this.stats.categorias) return {};
			// Retornar as top 3
			return Object.fromEntries(
				Object.entries(this.stats.categorias).slice(0, 3)
			);
		},
		maxCategoryCount() {
			if (!this.stats || !this.stats.categorias) return 1;
			const counts = Object.values(this.stats.categorias);
			return counts.length > 0 ? Math.max(...counts) : 1;
		}
	},
	created() {
		this.user = auth.getUser()
		this.buscar(true)
		if (this.user && this.user.id) {
			this.buscarHistorico()
			this.buscarFavoritos()
			this.buscarEstatisticas()
		}	},
	mounted() {
		this.initObserver();
	},
	beforeUnmount() {
		if (this.observer) this.observer.disconnect();
	},
	watch: {
		searchQuery(newVal) {
			if (!newVal) {
				this.suggestions = [];
				this.showSuggestions = false;
			}
		}
	},
	methods: {
		debouncedGetSuggestions: debounce(function() {
			this.getSuggestions();
		}, 300),
		async getSuggestions() {
			if (!this.searchQuery || this.searchQuery.length < 2) {
				this.suggestions = [];
				this.showSuggestions = false;
				return;
			}
			try {
				const response = await MaterialService.obterSugestoes(this.searchQuery);
				this.suggestions = response.data || [];
				this.showSuggestions = this.suggestions.length > 0;
				this.suggestionIndex = -1;
			} catch (err) {
				console.error('Erro ao buscar sugestões:', err);
			}
		},
		navigateSuggestions(dir) {
			if (!this.showSuggestions) return;
			this.suggestionIndex += dir;
			if (this.suggestionIndex < -1) this.suggestionIndex = this.suggestions.length - 1;
			if (this.suggestionIndex >= this.suggestions.length) this.suggestionIndex = -1;
		},
		handleEnter() {
			if (this.suggestionIndex >= 0 && this.suggestions[this.suggestionIndex]) {
				this.selectSuggestion(this.suggestions[this.suggestionIndex]);
			} else {
				this.showSuggestions = false;
				this.buscar(true);
			}
		},
		selectSuggestion(suggestion) {
			this.searchQuery = suggestion.titulo;
			this.showSuggestions = false;
			this.$router.push('/estudo/' + suggestion.id);
		},
		async buscarEstatisticas() {
			try {
				const response = await MaterialService.obterEstatisticas(this.user.id);
				this.stats = response.data;

				// Notificar novas conquistas através do App.vue (parent)
				if (this.$root.checkNewAchievements) {
					this.$root.checkNewAchievements(this.stats);
				}
			} catch (error) {
				console.error('Erro ao buscar estatísticas:', error);
			}
		},
		async buscar(reset = true) {
			if (reset) {
				this.offset = 0;
				this.livros = [];
				this.hasMore = true;
			}
			if (!this.hasMore && !reset) return;

			this.loading = true;
			try {
				const response = await MaterialService.pesquisar(
					this.searchQuery, '', '', 0, 0, this.limit, this.offset, 'relevancia'
				);
				const novosLivros = response.data || [];
				
				if (reset) {
					this.livros = novosLivros;
				} else {
					this.livros = [...this.livros, ...novosLivros];
				}
				
				this.hasMore = novosLivros.length === this.limit;
				this.offset += this.limit;
			} catch (error) {
				console.error('Erro ao buscar materiais:', error);
			} finally {
				this.loading = false;
			}
		},
		async loadMore() {
			if (this.loading || !this.hasMore) return;
			await this.buscar(false);
		},
		initObserver() {
			const options = {
				root: null,
				rootMargin: '400px', // Trigger earlier for smoother "YouTube" feel
				threshold: 0.1
			};

			this.observer = new IntersectionObserver((entries) => {
				if (entries[0].isIntersecting && !this.loading && this.hasMore) {
					this.loadMore();
				}
			}, options);

			if (this.$refs.loadMoreSentinel) {
				this.observer.observe(this.$refs.loadMoreSentinel);
			}
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
				await this.fetchGlobalFavorites()
			} catch (error) {
				console.error('Erro ao buscar favoritos:', error)
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
				await this.fetchGlobalFavorites()

				if (novoStatus) {
					this.showSnackbar('Adicionado aos favoritos!')
				} else {
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
		showSnackbar(msg, type = 'success') {
			this.notify(msg, type)
		},
		logout() {
			auth.logout()
			this.$router.push('/login')
		}
	},
}
</script>

<style scoped>
	.dashboard-container {
		min-height: 100vh;
		padding-bottom: 40px;
	}
	.load-more-sentinel-modern {
		height: 20px;
		width: 100%;
		visibility: hidden;
	}

	/* Suggestions Styles */
	.suggestions-dropdown-premium {
		position: absolute;
		top: 100%;
		left: 16px;
		right: 16px;
		background: rgba(var(--v-theme-surface), 0.95);
		backdrop-filter: blur(25px) saturate(180%);
		-webkit-backdrop-filter: blur(25px) saturate(180%);
		border-radius: 20px;
		margin-top: 8px;
		z-index: 1000;
		box-shadow: 0 20px 50px rgba(0,0,0,0.4);
		border: 1px solid rgba(255, 255, 255, 0.1);
		overflow: hidden;
		padding: 8px;
	}

	.suggestion-item {
		display: flex;
		align-items: center;
		padding: 12px 16px;
		border-radius: 12px;
		cursor: pointer;
		transition: all 0.2s ease;
	}

	.suggestion-item:hover, .suggestion-item.active {
		background: rgba(var(--v-theme-primary), 0.15);
		transform: translateX(5px);
	}

	.suggestion-title {
		font-weight: 700;
		font-size: 0.95rem;
		color: var(--v-theme-on-surface);
	}

	.suggestion-meta {
		font-size: 0.75rem;
		opacity: 0.5;
	}

	.fade-in-fast {
		animation: fadeInFast 0.2s ease-out;
	}

	@keyframes fadeInFast {
		from { opacity: 0; transform: translateY(-10px); }
		to { opacity: 1; transform: translateY(0); }
	}

	.tracking-tight {
		letter-spacing: -1.5px !important;
	}

	.ios-filter-btn {
		background: rgba(var(--ios-blue-rgb), 0.1) !important;
		backdrop-filter: blur(10px);
		color: var(--ios-blue) !important;
		border-radius: 12px !important;
		text-transform: none !important;
		border: 1px solid rgba(var(--ios-blue-rgb), 0.2);
	}

	.ios-search-bar {
		background: rgba(var(--ios-blue-rgb), 0.05) !important;
		border-radius: 16px !important;
	}

	:deep(.ios-search-bar .v-field) {
		background: transparent !important;
		border-radius: 16px !important;
		color: var(--v-theme-on-surface) !important;
	}

	/* Item Card Style */
	.ios-item-card {
		background: var(--glass-bg) !important;
		backdrop-filter: var(--glass-blur);
		-webkit-backdrop-filter: var(--glass-blur);
		border-radius: 32px !important;
		border: 1px solid var(--glass-border);
		overflow: hidden;
		transition: all 0.5s var(--spring-easing);
		opacity: 0;
		animation: ios-reveal 0.7s var(--spring-easing) forwards;
	}

	.reader-rounding { border-radius: 32px !important; }

	.ios-item-card:hover {
		transform: translateY(-10px) scale(1.02);
		background: rgba(255, 255, 255, 0.1) !important;
		box-shadow: 0 30px 60px rgba(0,0,0,0.3) !important;
		border-color: var(--ios-blue) !important;
	}

	/* Expert Stats Cards */
	.stat-card {
		background: var(--glass-bg);
		backdrop-filter: var(--glass-blur);
		border-radius: 32px;
		padding: 24px;
		height: 100%;
		border: 1px solid var(--glass-border);
		transition: all 0.5s cubic-bezier(0.16, 1, 0.3, 1);
		position: relative;
		overflow: hidden;
	}

	.stat-card::after {
		content: '';
		position: absolute;
		top: 0; right: 0; width: 100px; height: 100px;
		background: radial-gradient(circle, rgba(0, 122, 255, 0.05) 0%, transparent 70%);
		pointer-events: none;
	}

	.premium-glow:hover {
		background: rgba(var(--v-theme-surface), 0.6);
		border-color: rgba(0, 122, 255, 0.3);
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
		background: linear-gradient(90deg, #007AFF, #00C1FF);
		border-radius: 10px;
		box-shadow: 0 0 15px rgba(0, 122, 255, 0.3);
	}

	.premium-bar :deep(.v-progress-linear__determinate) {
		background: linear-gradient(90deg, #007AFF, #00C1FF) !important;
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
		line-clamp: 2;
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
		background: var(--glass-bg) !important;
		backdrop-filter: var(--glass-blur);
		border-radius: 24px !important;
		border: 1px solid rgba(128, 128, 128, 0.2);
		overflow: hidden;
		cursor: pointer;
		transition: transform 0.3s ease;
	}

	.ios-history-card:hover {
		transform: scale(1.05);
		background: rgba(255, 255, 255, 0.1) !important;
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

	.mini-premium-cover {
		width: 100%;
		height: 100%;
		border-radius: 8px;
		position: relative;
		overflow: hidden;
		display: flex;
		align-items: center;
		justify-content: center;
		background: linear-gradient(145deg, #1e293b 0%, #0f172a 100%);
		border: 1px solid rgba(255, 255, 255, 0.1);
	}

	.mini-premium-cover.has-image {
		background: transparent;
	}

	.cover-image-inner {
		position: absolute;
		top: 0;
		left: 0;
		width: 100%;
		height: 100%;
		z-index: 1;
	}

	.cover-overlay {
		position: absolute;
		top: 0;
		left: 0;
		width: 100%;
		height: 100%;
		background: linear-gradient(to top, rgba(15, 23, 42, 0.6) 0%, transparent 100%);
		z-index: 2;
	}

	.cover-content {
		position: relative;
		z-index: 3;
	}

	.mini-premium-cover-ios-stats {
		width: 48px;
		height: 48px;
		border-radius: 12px;
		position: relative;
		overflow: hidden;
		background: #000;
		display: flex;
		align-items: center;
		justify-content: center;
		flex-shrink: 0;
	}

	.mini-premium-cover-ios-dashboard {
		width: 100%;
		height: 160px;
		position: relative;
		overflow: hidden;
		background: #000;
		display: flex;
		align-items: center;
		justify-content: center;
		transition: all 0.6s cubic-bezier(0.16, 1, 0.3, 1);
	}

	.mesh-gradient-mini {
		position: absolute;
		inset: -50%;
		background: 
			radial-gradient(circle at 30% 30%, #007AFF 0%, transparent 60%),
			radial-gradient(circle at 70% 70%, #5AC8FA 0%, transparent 60%);
		filter: blur(20px);
		opacity: 0.8;
	}

	.glass-overlay-mini {
		position: absolute;
		inset: 8px;
		background: rgba(255, 255, 255, 0.05);
		backdrop-filter: blur(20px) saturate(180%);
		border: 1px solid rgba(255, 255, 255, 0.15);
		border-radius: 12px;
		z-index: 2;
		display: flex;
		flex-direction: column;
		justify-content: center;
		align-items: center;
		padding: 12px;
		text-align: center;
	}

	.cover-content-ios-mini {
		position: relative;
		z-index: 3;
	}

	.cover-title-ios-mini {
		color: #ffffff;
		font-family: 'Outfit', 'Inter', sans-serif;
		font-size: 0.55rem;
		font-weight: 800;
		line-height: 1.1;
		display: -webkit-box;
		-webkit-line-clamp: 2;
		-webkit-box-orient: vertical;
		overflow: hidden;
		margin-bottom: 2px;
		text-shadow: 0 1px 4px rgba(0, 0, 0, 0.5);
	}

	.cover-summary-ios-mini {
		color: rgba(255, 255, 255, 0.5);
		font-size: 0.45rem;
		font-weight: 400;
		line-height: 1.2;
		display: -webkit-box;
		-webkit-line-clamp: 2;
		-webkit-box-orient: vertical;
		overflow: hidden;
	}

	.glass-shine {
		position: absolute;
		top: 0;
		left: 0;
		width: 100%;
		height: 100%;
		background: linear-gradient(135deg, rgba(255,255,255,0.15) 0%, transparent 50%, rgba(255,255,255,0.05) 100%);
		pointer-events: none;
		z-index: 4;
	}

	.book-cover {
		width: 100%;
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
		line-clamp: 2;
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
		background: var(--ios-blue) !important;
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
		background: var(--glass-bg) !important;
		backdrop-filter: var(--glass-blur);
		border: 1px solid var(--glass-border);
		border-radius: 32px !important;
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
	.study-hub-card.notes:hover { border-color: rgba(0, 122, 255, 0.4); }

	.top-book-item {
		transition: all 0.3s var(--spring-easing);
		cursor: pointer;
		border-radius: 16px !important;
	}
	.top-book-item:hover {
		background: rgba(255, 255, 255, 0.05) !important;
		transform: translateX(8px);
	}

	.top-book-rank {
		width: 28px;
		height: 28px;
		display: flex;
		align-items: center;
		justify-content: center;
		border-radius: 50%;
		font-weight: 800;
		font-size: 14px;
		background: rgba(255, 255, 255, 0.1);
		color: white;
	}
	.rank-1 { background: linear-gradient(135deg, #FFD700, #FFA000); color: black; }
	.rank-2 { background: linear-gradient(135deg, #C0C0C0, #9E9E9E); color: black; }
	.rank-3 { background: linear-gradient(135deg, #CD7F32, #8D6E63); color: black; }

	.hub-icon-wrap {
		width: 70px; height: 70px;
		border-radius: 20px;
		display: flex;
		align-items: center;
		justify-content: center;
		flex-shrink: 0;
	}
	.hub-icon-wrap.amber { background: rgba(255, 171, 0, 0.1); }
	.hub-icon-wrap.cyan { background: rgba(0, 122, 255, 0.1); }
</style>
