<template>
	<div class="explore-container">
		<v-container fluid class="px-4 py-0">
			<!-- Advanced Header Section -->
			<v-row dense align="center" class="mt-2 mb-2">
				<v-col cols="12">
					<div class="filter-bar-premium d-flex flex-wrap align-center pa-2">
						<!-- Categoria Filter -->
						<v-col cols="12" sm="4" class="pa-2">
							<v-select
								v-model="filters.categoria"
								:items="['TODOS', ...categoriasMock]"
								label="Categoria"
								variant="solo-filled"
								density="comfortable"
								rounded="lg"
								hide-details
								clearable
								prepend-inner-icon="mdi-shape-outline"
								class="filter-select-premium"
								@update:modelValue="buscar"
							></v-select>
						</v-col>

						<!-- Ano Filter -->
						<v-col cols="12" sm="4" class="pa-2">
							<v-select
								v-model="filters.ano_inicio"
								:items="yearsList"
								label="Ano de Publicação"
								variant="solo-filled"
								density="comfortable"
								rounded="lg"
								hide-details
								clearable
								prepend-inner-icon="mdi-calendar-range"
								class="filter-select-premium"
								@update:modelValue="buscar"
							></v-select>
						</v-col>

						<!-- Ordenação Filter -->
						<v-col cols="12" sm="4" class="pa-2">
							<v-select
								v-model="filters.sort"
								:items="sortOptions"
								item-title="label"
								item-value="value"
								label="Ordenar por"
								variant="solo-filled"
								density="comfortable"
								rounded="lg"
								hide-details
								prepend-inner-icon="mdi-sort-variant"
								class="filter-select-premium"
								@update:modelValue="buscar"
							></v-select>
						</v-col>
					</div>


				</v-col>
			</v-row>



			<!-- Content Grid -->
			<div class="results-container">
				<v-row dense>
					<v-col v-for="(livro, index) in livros" :key="livro.id" cols="12" sm="12" md="6" lg="6" class="pa-2">
						<BookCard
							:book="livro"
							:animation-delay="index * 30"
							:is-favorited="isBookFavorited(livro.id)"
							@toggle-favorite="onToggleFavorite"
							@share="shareBook(livro.id)"
						/>
					</v-col>

					<v-col cols="12" v-if="(!livros || livros.length === 0) && !loading" class="text-center py-12">
						<v-icon size="64" color="primary" class="mb-4 opacity-50">mdi-text-search-variant</v-icon>
						<h3 :class="isDarkTheme ? 'text-white' : 'text-primary-darken-4'">Nenhum resultado</h3>
						<p :class="isDarkTheme ? 'text-white' : 'text-primary-darken-2'" class="opacity-60">Tente outros termos ou filtros.</p>
					</v-col>
				</v-row>

				<!-- Timeout Warning -->
				<v-row dense v-if="loading && showTimeoutWarning" class="mb-4">
					<v-col cols="12">
						<v-alert
							type="info"
							variant="tonal"
							class="rounded-xl font-weight-medium"
							icon="mdi-rocket-launch-outline"
							elevation="0"
						>
							A busca está demorando um pouco mais. O Acervus Core está conectando às bibliotecas e bases acadêmicas externas para expandir seus resultados!
						</v-alert>
					</v-col>
				</v-row>

				<!-- Loading More Shimmer -->
				<v-row dense v-if="loading" class="mt-8 mb-8" justify="center">
					<CloudLoader text="Buscando no acervo..." height="200px" />
				</v-row>

				<!-- Intersection Observer Sentinel -->
				<div ref="loadMoreSentinel" class="load-more-sentinel-modern"></div>
			</div>
		</v-container>
	</div>
</template>

<script>
import MaterialService from '@/services/MaterialService'
import { computed } from 'vue'
import { useTheme } from 'vuetify'

export default {
	name: 'ExplorePageExtended',
	inject: ['notify', 'fetchGlobalFavorites', 'getGlobalFavorites', 'globalFavorites'],
	components: {
	},
	data: () => ({
		livros: [],
		loading: false,
		showFilters: false,
		filters: {
			q: '',
			categoria: 'TODOS',
			ano_inicio: null,
			ano_fim: null,
			sort: 'recent'
		},
		isFetching: false,
		hasInitialFetchDone: false,
		offset: 0,
		limit: 16,
		hasMore: true,

		categoriasMock: ['TECNOLOGIA', 'SAÚDE', 'MATEMÁTICA', 'CIÊNCIAS', 'HISTÓRIA', 'EDUCAÇÃO', 'JURÍDICO', 'LITERATURA BRASILEIRA', 'CONTABILIDADE', 'DIREITO BRASILEIRO', 'SAÚDE PÚBLICA BRASIL', 'TECNOLOGIA BRASIL'],
		sortOptions: [
			{ label: 'Ordem Alfabética (A-Z)', value: 'az' },
			{ label: 'Ordem Alfabética (Z-A)', value: 'za' },
			{ label: 'Data de Publicação (Mais novos)', value: 'recent' },
			{ label: 'Data de Publicação (Antigos)', value: 'oldest' },
			{ label: 'Aleatório', value: 'random' }
		],
		yearsList: Array.from({length: 40}, (_, i) => new Date().getFullYear() - i),
		searchTimeout: null,
		longRequestTimeout: null,
		showTimeoutWarning: false,
		observer: null
	}),
	setup() {
		const theme = useTheme();
		const isDarkTheme = computed(() => theme.global.current.value.dark);
		return { isDarkTheme };
	},
	computed: {
		favoritos() {
			return this.globalFavorites?.list || [];
		},
		activeFiltersCount() {
			let count = 0;
			if (this.filters.ano_inicio) count++;
			if (this.filters.categoria && this.filters.categoria !== 'TODOS') count++;
			return count;
		}
	},
	mounted() {
		this.initObserver();
		document.addEventListener('scroll', this.handleScroll, { passive: true, capture: true });
	},
	watch: {
		'$route.query': {
			immediate: true,
			handler(query) {
				// Prevent double fetching if already fetching
				if (this.loading) return;
				let changed = false;
				
				const newQ = query.q || '';
				if (newQ !== this.filters.q) {
					this.filters.q = newQ;
					changed = true;
				}
				
				const newCategoria = query.categoria || 'TODOS';
				if (newCategoria !== this.filters.categoria) {
					this.filters.categoria = newCategoria;
					changed = true;
				}

				// Fetch once correctly
				if (changed || !this.hasInitialFetchDone) {
					// Atraso intencional para permitir fluidez na transição ios-page
					setTimeout(() => {
						this.buscar();
					}, 300);
					this.hasInitialFetchDone = true;
				}
			}
		}
	},
	unmounted() {
		this.isFetching = false;
		this.loading = false;
		if (this.longRequestTimeout) clearTimeout(this.longRequestTimeout);
		if (this.searchTimeout) clearTimeout(this.searchTimeout);
		if (this.observer) {
			this.observer.disconnect();
		}
	},
	methods: {
		async buscar(reset = true) {
			if (this.isFetching) return;

			if (reset) {
				this.offset = 0;
				this.livros = [];
				this.hasMore = true;
			}

			if (!this.hasMore && !reset) return;

			if (this.searchTimeout) clearTimeout(this.searchTimeout);

			this.isFetching = true;
			this.loading = true;
			this.showTimeoutWarning = false;

			if (this.longRequestTimeout) clearTimeout(this.longRequestTimeout);
			this.longRequestTimeout = setTimeout(() => {
				if (this.loading && reset) {
					this.showTimeoutWarning = true;
				}
			}, 3000);

			try {
				// Sincronizar favoritos se necessário
				const userStr = localStorage.getItem('user');
				if (userStr && this.favoritos.length === 0) {
					await this.fetchGlobalFavorites();
				}

				// Trata o valor "TODOS" como string vazia para a API
				const categoriaParaBusca = (this.filters.categoria === 'TODOS' || !this.filters.categoria) ? '' : this.filters.categoria;

				const response = await MaterialService.pesquisar(
					this.filters.q,
					categoriaParaBusca,
					'', // fonte
					this.filters.ano_inicio,
					this.filters.ano_fim,
					this.limit,
					this.offset,
					this.filters.sort
				)

				const novosLivrosRaw = response.data || [];

				if (reset) {
					this.livros = novosLivrosRaw;
				} else {
					// Deduplicação por ID para evitar repetições na rolagem infinita
					const idsExistentes = new Set(this.livros.map(l => l.id));
					const novosLivrosFiltrados = novosLivrosRaw.filter(l => !idsExistentes.has(l.id));

					if (novosLivrosFiltrados.length > 0) {
						this.livros = [...this.livros, ...novosLivrosFiltrados];
					}
				}

				this.hasMore = novosLivrosRaw.length === this.limit;
				if (novosLivrosRaw.length > 0) {
					this.offset += this.limit;
				}

				if (reset && this.livros.length === 0) {
					this.notify('Nenhum material encontrado para sua busca.', 'info')
				}
			} catch (error) {
				console.error('Erro na pesquisa avançada:', error)
				this.notify('Erro ao realizar busca. Tente novamente.', 'error')
			} finally {
				if (this.longRequestTimeout) clearTimeout(this.longRequestTimeout);
				this.loading = false;
				this.isFetching = false;
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
		handleScroll() {
			if (this.loading || !this.hasMore || this.isFetching) return;
			
			if (this.$refs.loadMoreSentinel) {
				const rect = this.$refs.loadMoreSentinel.getBoundingClientRect();
				// Se a sentinela estiver a 800px ou menos de entrar na tela (ou já na tela)
				if (rect.top > 0 && rect.top <= window.innerHeight + 800) {
					this.loadMore();
				}
			}
		},
		onSearchInput() {
			if (this.searchTimeout) clearTimeout(this.searchTimeout);
			this.searchTimeout = setTimeout(() => {
				this.buscar();
			}, 400);
		},
		limparFiltros() {
			this.filters = { q: '', categoria: 'TODOS', ano_inicio: null, ano_fim: null };
			this.$router.replace({ path: '/explorar', query: {} });
			this.buscar();
		},
		clearSearch() {
			this.filters.q = '';
			const query = { ...this.$route.query };
			delete query.q;
			this.$router.replace({ path: '/explorar', query });
			this.buscar();
		},
		async onToggleFavorite(livro) {
			try {
				const userStr = localStorage.getItem('user')
				if (!userStr) {
					this.notify('Faça login para favoritar materiais!', 'warning')
					return
				}
				const user = JSON.parse(userStr)
				const currentlyFavorited = this.isBookFavorited(livro.id)

				await MaterialService.favoritar(user.id, livro.id, !currentlyFavorited)

				// Atualiza lista global
				await this.fetchGlobalFavorites()

				if (currentlyFavorited) {
					this.notify('Removido dos favoritos', 'info')
				} else {
					this.notify('Adicionado aos favoritos!', 'success')
				}
			} catch (err) {
				console.error(err)
				this.notify('Erro ao atualizar favorito', 'error')
			}
		},

		isBookFavorited(bookId) {
			return this.favoritos.some(f => f.id === bookId)
		},
		shareBook(id) {
			const link = `${window.location.origin}/estudo/${id}`
			navigator.clipboard.writeText(link).then(() => {
				this.notify('Link copiado para a área de transferência!', 'success')
			})
		}
	},
	beforeUnmount() {
		if (this.searchTimeout) clearTimeout(this.searchTimeout);
		if (this.observer) this.observer.disconnect();
		document.removeEventListener('scroll', this.handleScroll, { capture: true });
	}
}
</script>

<style scoped>
	.explore-container {
		min-height: calc(100vh - 80px);
		padding-bottom: 40px;
		margin-top: 0;
		background: transparent;
	}
	.results-container {
		background: var(--glass-bg);
		backdrop-filter: var(--glass-blur);
		-webkit-backdrop-filter: var(--glass-blur);
		border-radius: 32px;
		padding: 24px;
		border: 1px solid var(--glass-border);
		margin-top: 16px;
	}

	.filter-bar-premium {
		background: rgba(255, 255, 255, 0.03);
		backdrop-filter: blur(10px);
		border-radius: 24px;
		border: 1px solid rgba(255, 255, 255, 0.08);
	}

	.filter-select-premium :deep(.v-field) {
		background: rgba(255, 255, 255, 0.05) !important;
		border-radius: 16px !important;
		border: 1px solid rgba(255, 255, 255, 0.1);
		transition: all 0.3s ease;
	}

	.filter-select-premium :deep(.v-field--focused) {
		border-color: var(--ios-blue) !important;
		background: rgba(0, 122, 255, 0.05) !important;
	}

	:deep(.v-label) {
		font-size: 13px;
		font-weight: 600;
		opacity: 0.7;
	}
	.load-more-sentinel-modern {
		height: 20px;
		width: 100%;
		margin-top: 10px;
		visibility: hidden;
	}

	.category-scroll-container {
		overflow-x: auto;
		-ms-overflow-style: none;
		scrollbar-width: none;
	}
	.category-scroll-container::-webkit-scrollbar { display: none; }

	.premium-chip {
		background: var(--glass-bg) !important;
		backdrop-filter: var(--glass-blur);
		border: 1px solid var(--glass-border) !important;
		color: #ffffff !important;
		transition: all 0.3s ease;
	}
	.premium-chip.v-chip--active {
		background: linear-gradient(135deg, rgba(0, 122, 255, 0.3), rgba(0, 122, 153, 0.5)) !important;
		border-color: rgba(0, 122, 255, 0.5) !important;
		box-shadow: 0 4px 12px rgba(0, 122, 255, 0.2);
	}

	.ios-search-bar :deep(.v-field) {
		background: rgba(255, 255, 255, 0.05) !important;
		border-radius: 20px !important;
		border: 1px solid rgba(255, 255, 255, 0.1);
		color: #ffffff !important;
	}

	.search-btn {
		background: linear-gradient(135deg, #007AFF 0%, #0056B3 100%) !important;
		color: white !important;
		font-weight: 800 !important;
	}

	.ios-filter-toggle {
		background: rgba(255, 255, 255, 0.05) !important;
		color: white !important;
		border-color: rgba(255, 255, 255, 0.2) !important;
	}

	.ios-glass-card {
		background: var(--glass-bg) !important;
		backdrop-filter: var(--glass-blur);
		border: 1px solid var(--glass-border);
	}

	.filter-label {
		display: block;
		font-size: 11px;
		font-weight: 700;
		color: var(--v-theme-on-surface);
		opacity: 0.5;
		text-transform: uppercase;
		letter-spacing: 1px;
		margin-left: 4px;
	}

	:deep(.ios-item-card) {
		background: var(--glass-bg) !important;
		backdrop-filter: var(--glass-blur);
		border: 1px solid var(--glass-border) !important;
		border-radius: 20px !important;
		transition: all 0.3s ease;
		box-shadow: 0 8px 32px rgba(0, 0, 0, 0.2) !important;
	}
	:deep(.ios-item-card:hover) {
		transform: translateY(-5px);
		background: rgba(255, 255, 255, 0.1) !important;
		border-color: var(--ios-blue) !important;
	}

	.item-title {
		font-size: 1.1rem;
		font-weight: 800;
		letter-spacing: -0.5px;
		color: #ffffff;
		line-height: 1.2;
		height: 2.4em;
		overflow: hidden;
		display: -webkit-box;
		-webkit-line-clamp: 2;
		line-clamp: 2;
		-webkit-box-orient: vertical;
	}

	.item-details { color: rgba(255, 255, 255, 0.7); font-size: 13px; font-weight: 500; }

	.ios-btn-open { font-size: 11px !important; font-weight: 700 !important; }

	.ios-item-card-skeleton {
		background: rgba(255, 255, 255, 0.04) !important;
		border-radius: 24px !important;
	}
</style>
