<template>
	<div class="explore-container">
		<v-container fluid class="px-4 py-0">
			<!-- Advanced Header Section -->
			<v-row dense align="center" class="mt-2 mb-2">
				<v-col cols="12">
					<div class="search-premium-box d-flex align-center">
						<v-text-field
							v-model="filters.q"
							prepend-inner-icon="mdi-magnify"
							placeholder="Busque por título, autor, DOI ou assunto..."
							variant="solo"
							class="ios-search-bar flex-grow-1"
							hide-details
							@update:modelValue="onSearchInput"
							@keyup.enter="buscar(true)"
							density="comfortable"
						></v-text-field>

						<v-btn
							class="ios-filter-toggle-premium ml-3"
							variant="tonal"
							icon="mdi-tune-variant"
							@click="showFilters = true"
							height="48"
							width="48"
							rounded="lg"
						></v-btn>

						<v-btn
							class="search-btn ml-3 hidden-sm-and-down"
							@click="buscar(true)"
							:loading="loading && (livros?.length === 0)"
							height="48"
							rounded="lg"
							min-width="120"
						>
							EXPLORAR
						</v-btn>
					</div>

					<div v-if="activeFiltersCount > 0" class="d-flex flex-wrap gap-2 mt-3 mb-1">
						<v-chip
							v-if="filters.categoria"
							closable
							size="small"
							color="primary"
							variant="tonal"
							@click:close="filters.categoria = ''; buscar(true)"
						>
							Categoria: {{ filters.categoria }}
						</v-chip>
						<v-chip
							v-if="filters.ano_inicio"
							closable
							size="small"
							color="primary"
							variant="tonal"
							@click:close="filters.ano_inicio = null; buscar(true)"
						>
							Desde: {{ filters.ano_inicio }}
						</v-chip>
					</div>
				</v-col>
			</v-row>

			<!-- Modern Filter Drawer -->
			<v-navigation-drawer
				v-model="showFilters"
				location="right"
				temporary
				width="320"
				class="filter-drawer-premium"
			>
				<div class="pa-6">
					<div class="d-flex align-center justify-space-between mb-8">
						<h2 class="text-h6 font-weight-black" :class="isDarkTheme ? 'text-white' : 'text-primary-darken-4'">Refinar Busca</h2>
						<v-btn icon="mdi-close" variant="text" size="small" @click="showFilters = false"></v-btn>
					</div>

					<div class="filter-section mb-6">
						<label class="filter-label mb-3 d-block">Categoria</label>
						<v-select
							v-model="filters.categoria"
							:items="categoriasMock"
							placeholder="Todas as categorias"
							variant="solo-filled"
							density="comfortable"
							rounded="lg"
							hide-details
							clearable
							class="filter-select-premium"
						></v-select>
					</div>

					<div class="filter-section mb-6">
						<label class="filter-label mb-3 d-block">Ano de Publicação</label>
						<v-select
							v-model="filters.ano_inicio"
							:items="yearsList"
							placeholder="Qualquer ano"
							variant="solo-filled"
							density="comfortable"
							rounded="lg"
							hide-details
							clearable
							class="filter-select-premium"
						></v-select>
					</div>

					<div class="filter-section mb-10">
						<label class="filter-label mb-3 d-block">Ordenar Resultados</label>
						<v-select
							v-model="filters.sort"
							:items="sortOptions"
							item-title="label"
							item-value="value"
							variant="solo-filled"
							density="comfortable"
							rounded="lg"
							hide-details
							class="filter-select-premium"
						></v-select>
					</div>

					<v-btn color="primary" block @click="buscar(true); showFilters = false" height="48" class="font-weight-bold">
						APLICAR FILTROS
					</v-btn>

					<v-btn variant="text" block @click="limparFiltros" class="mt-2 text-none opacity-60">
						Limpar tudo
					</v-btn>
				</div>
			</v-navigation-drawer>

			<!-- Quick Categories -->
			<div class="category-row mb-4" v-if="!showFilters">
				<v-chip-group v-model="filters.categoria" @update:modelValue="buscar" mandatory color="primary">
					<v-chip value="" class="premium-chip" variant="tonal" filter>Todos</v-chip>
					<v-chip v-for="cat in categoriasMock" :key="cat" :value="cat" class="premium-chip" variant="tonal" filter>
						{{ cat }}
					</v-chip>
				</v-chip-group>
			</div>

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

				<!-- Loading More Shimmer -->
				<v-row dense v-if="loading">
					<v-col v-for="n in 4" :key="'skeleton-more-'+n" cols="12" sm="12" md="6" lg="6" class="pa-2">
						<v-skeleton-loader type="card" :theme="isDarkTheme ? 'dark' : 'light'" class="rounded-xl"></v-skeleton-loader>
					</v-col>
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
			categoria: '',
			ano_inicio: null,
			ano_fim: null,
			sort: 'az'
		},
		hasInitialFetchDone: false,
		offset: 0,
		limit: 16,
		hasMore: true,

		categoriasMock: ['TECNOLOGIA', 'SAÚDE', 'MATEMÁTICA', 'CIÊNCIAS', 'HISTÓRIA DO BRASIL', 'EDUCAÇÃO', 'JURÍDICO', 'LITERATURA BRASILEIRA', 'CONTABILIDADE', 'DIREITO BRASILEIRO', 'SAÚDE PÚBLICA BRASIL', 'TECNOLOGIA BRASIL'],
		sortOptions: [
			{ label: 'Ordem Alfabética (A-Z)', value: 'az' },
			{ label: 'Ordem Alfabética (Z-A)', value: 'za' },
			{ label: 'Data de Publicação (Mais novos)', value: 'recent' },
			{ label: 'Data de Publicação (Antigos)', value: 'oldest' },
			{ label: 'Aleatório', value: 'random' }
		],
		yearsList: Array.from({length: 30}, (_, i) => 2025 - i),
		searchTimeout: null,
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
			if (this.filters.categoria) count++;
			return count;
		}
	},
	mounted() {
		this.initObserver();
	},
	watch: {
		'$route.query': {
			immediate: true,
			handler(query) {
				// Prevent double fetching if already fetching
				if (this.loading) return;
				let changed = false;
				if (query.q && query.q !== this.filters.q) {
					this.filters.q = query.q;
					changed = true;
				}
				if (query.categoria && query.categoria !== this.filters.categoria) {
					this.filters.categoria = query.categoria;
					changed = true;
				}

				// Fetch once correctly
				if (changed || !this.hasInitialFetchDone) {
					this.buscar();
					this.hasInitialFetchDone = true;
				}
			}
		}
	},
	unmounted() {
		if (this.observer) {
			this.observer.disconnect();
		}
	},
	methods: {
		async buscar(reset = true) {
			if (reset) {
				this.offset = 0;
				this.livros = [];
				this.hasMore = true;
			}

			if (!this.hasMore && !reset) return;

			if (this.searchTimeout) clearTimeout(this.searchTimeout);

			this.loading = true
			try {
				// Sincronizar favoritos se necessário
				const userStr = localStorage.getItem('user');
				if (userStr && this.favoritos.length === 0) {
					await this.fetchGlobalFavorites();
				}

				const response = await MaterialService.pesquisar(
					this.filters.q,
					this.filters.categoria,
					'', // fonte
					this.filters.ano_inicio,
					this.filters.ano_fim,
					this.limit,
					this.offset,
					this.filters.sort
				)

				const novosLivros = response.data || [];
				if (reset) {
					this.livros = novosLivros;
				} else {
					this.livros = [...this.livros, ...novosLivros];
				}

				this.hasMore = novosLivros.length === this.limit;
				this.offset += this.limit;

				if (reset && this.livros.length === 0) {
					this.notify('Nenhum material encontrado para sua busca.', 'info')
				}
			} catch (error) {
				console.error('Erro na pesquisa avançada:', error)
				this.notify('Erro ao realizar busca. Tente novamente.', 'error')
			} finally {
				this.loading = false
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
		onSearchInput() {
			if (this.searchTimeout) clearTimeout(this.searchTimeout);
			this.searchTimeout = setTimeout(() => {
				this.buscar();
			}, 400);
		},
		limparFiltros() {
			this.filters = { q: this.filters.q, categoria: '', ano_inicio: null, ano_fim: null };
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
		padding: 32px;
		border: 1px solid var(--glass-border);
		margin-top: 24px;
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
