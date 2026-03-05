<template>
	<div class="explore-container">
		<v-container fluid class="px-4 py-0">
			<!-- Advanced Header Section -->
			<v-row dense align="center" class="mt-2 mb-2">
				<v-col cols="12" md="8" lg="9">
					<div class="search-premium-box d-flex align-center">
						<v-text-field
							v-model="filters.q"
							prepend-inner-icon="mdi-magnify"
							placeholder="Busque por título, autor, DOI ou assunto..."
							variant="solo"
							class="ios-search-bar"
							hide-details
							@update:modelValue="onSearchInput"
							@keyup.enter="buscar"
							density="comfortable"
						></v-text-field>
						<v-btn
							class="search-btn ml-3"
							@click="buscar"
							:loading="loading"
							height="48"
							rounded="lg"
						>
							EXPLORAR
						</v-btn>
					</div>
				</v-col>
				<v-col cols="12" md="4" lg="3">
					<v-btn
						class="ios-filter-toggle w-100"
						variant="tonal"
						prepend-icon="mdi-tune-variant"
						@click="showFilters = !showFilters"
						height="48"
						rounded="lg"
					>
						Filtros
						<v-chip size="x-small" color="cyan" class="ml-2" v-if="activeFiltersCount > 0">
							{{ activeFiltersCount }}
						</v-chip>
					</v-btn>
				</v-col>
			</v-row>

			<!-- Expandable Filter Panel -->
			<v-expand-transition>
				<div v-show="showFilters" class="filter-panel mb-4">
					<v-card class="ios-glass-card pa-4" rounded="lg">
						<v-row dense>
							<v-col cols="12" sm="6">
								<v-select
									v-model="filters.ano_inicio"
									:items="yearsList"
									label="Ano de Publicação"
									variant="filled"
									density="compact"
									hide-details
									clearable
								></v-select>
							</v-col>
							<v-col cols="12" sm="6">
								<v-select
									v-model="filters.sort"
									:items="sortOptions"
									item-title="label"
									item-value="value"
									label="Ordenar Por"
									variant="filled"
									density="compact"
									hide-details
									@update:modelValue="buscar"
								></v-select>
							</v-col>
						</v-row>
						<div class="d-flex justify-end mt-3 gap-2">
							<v-btn variant="text" size="small" @click="limparFiltros">Limpar</v-btn>
							<v-btn color="cyan" size="small" @click="buscar">Aplicar</v-btn>
						</div>
					</v-card>
				</div>
			</v-expand-transition>

			<!-- Quick Categories -->
			<div class="category-row mb-4" v-if="!showFilters">
				<v-chip-group v-model="filters.categoria" @update:modelValue="buscar" mandatory color="cyan">
					<v-chip value="" class="premium-chip" variant="tonal" filter>Todos</v-chip>
					<v-chip v-for="cat in categoriasMock" :key="cat" :value="cat" class="premium-chip" variant="tonal" filter>
						{{ cat }}
					</v-chip>
				</v-chip-group>
			</div>

			<!-- Content Grid -->
			<v-row dense v-if="!loading">
				<v-col v-for="(livro, index) in livros" :key="livro.id" cols="12" sm="6" md="4" lg="3" class="pa-2">
					<BookCard
						:book="livro"
						:animation-delay="index * 30"
						@toggle-favorite="onToggleFavorite"
						@share="shareBook(livro.id)"
					/>
				</v-col>

				<v-col cols="12" v-if="livros.length === 0" class="text-center py-12">
					<v-icon size="64" color="cyan" class="mb-4 opacity-50">mdi-text-search-variant</v-icon>
					<h3 class="text-white">Nenhum resultado</h3>
					<p class="text-white opacity-60">Tente outros termos ou filtros.</p>
				</v-col>
			</v-row>

			<!-- Skeleton -->
			<v-row dense v-else>
				<v-col v-for="n in 8" :key="n" cols="12" sm="6" md="4" lg="3" class="pa-2">
					<v-skeleton-loader type="card" theme="dark" class="rounded-xl"></v-skeleton-loader>
				</v-col>
			</v-row>
		</v-container>
	</div>
</template>

<script>
import MaterialService from '@/services/MaterialService'
import BookCard from '@/components/BookCard.vue'

export default {
	name: 'ExplorePageExtended',
	components: {
		BookCard
	},
	data: () => ({
		livros: [],
		loading: false,
		showFilters: false,
		filters: {
			q: '',
			categoria: '',
			ano_inicio: null,
			ano_fim: null
		},

		categoriasMock: ['TECNOLOGIA', 'SAÚDE', 'MATEMÁTICA', 'CIÊNCIAS', 'HISTÓRIA', 'CONTABILIDADE', 'Artigo Periódico'],
		sortOptions: [
			{ label: 'Relevância', value: '' },
			{ label: 'Melhor Avaliados', value: 'rating' },
			{ label: 'Recentes', value: 'id' },
			{ label: 'Aleatório', value: 'random' }
		],
		yearsList: Array.from({length: 30}, (_, i) => 2025 - i),
		searchTimeout: null
	}),
	computed: {
		activeFiltersCount() {
			let count = 0;
			if (this.filters.ano_inicio) count++;
			if (this.filters.categoria) count++;
			return count;
		}
	},
	watch: {
		'$route.query': {
			immediate: true,
			handler(query) {
				if (query.q) this.filters.q = query.q;
				if (query.categoria) this.filters.categoria = query.categoria;
				this.buscar();
			}
		}
	},
	methods: {
		async buscar() {
			if (this.searchTimeout) clearTimeout(this.searchTimeout);
			this.loading = true
			try {
				const response = await MaterialService.pesquisar(
					this.filters.q,
					this.filters.categoria,
					'', // fonte
					this.filters.ano_inicio,
					this.filters.ano_fim,
					20,   // limit
					0,    // offset
					this.filters.sort
				)
				this.livros = response.data || []
			} catch (error) {
				console.error('Erro na pesquisa avançada:', error)
			} finally {
				this.loading = false
			}
		},
		onSearchInput() {
			if (this.searchTimeout) clearTimeout(this.searchTimeout);
			this.searchTimeout = setTimeout(() => {
				this.buscar();
			}, 600);
		},
		limparFiltros() {
			this.filters = { q: this.filters.q, categoria: '', ano_inicio: null, ano_fim: null };
			this.buscar();
		},
		async onToggleFavorite(livro) {
			// Simples implementação de feedback, já que favoritos exigem login
			// Esta view pode ser acessada por visitantes
			try {
				const user = JSON.parse(localStorage.getItem('user'))
				if (!user) {
					alert('Faça login para favoritar materiais!')
					return
				}
				await MaterialService.favoritar(user.id, livro.id, true)
			} catch (err) {
				console.error(err)
			}
		},
		shareBook(id) {
			const link = `${window.location.origin}/estudo/${id}`
			navigator.clipboard.writeText(link).then(() => {
				alert('Link copiado!')
			})
		}
	}
}
</script>

<style scoped>
	.explore-container { min-height: calc(100vh - 80px); padding-bottom: 40px; margin-top: -12px; }
	.search-premium-box { width: 100%; }

	.category-scroll-container {
		overflow-x: auto;
		-ms-overflow-style: none;
		scrollbar-width: none;
	}
	.category-scroll-container::-webkit-scrollbar { display: none; }

	.premium-chip {
		background: rgba(255, 255, 255, 0.05) !important;
		backdrop-filter: blur(10px);
		border: 1px solid rgba(255, 255, 255, 0.1) !important;
		color: white !important;
		transition: all 0.3s ease;
	}
	.premium-chip.v-chip--active {
		background: linear-gradient(135deg, rgba(0, 184, 212, 0.3), rgba(0, 122, 153, 0.5)) !important;
		border-color: rgba(0, 184, 212, 0.5) !important;
		box-shadow: 0 4px 12px rgba(0, 184, 212, 0.2);
	}

	.ios-search-bar :deep(.v-field) {
		background: rgba(255, 255, 255, 0.08) !important;
		border-radius: 20px !important;
		border: 1px solid rgba(255, 255, 255, 0.1);
		color: white !important;
	}

	.search-btn {
		background: linear-gradient(135deg, #00B8D4 0%, #007A99 100%) !important;
		color: white !important;
		font-weight: 800 !important;
	}

	.ios-filter-toggle {
		background: rgba(255, 255, 255, 0.05) !important;
		color: white !important;
		border-color: rgba(255, 255, 255, 0.2) !important;
	}

	.ios-glass-card {
		background: rgba(255, 255, 255, 0.05) !important;
		backdrop-filter: blur(20px);
		border: 1px solid rgba(255, 255, 255, 0.1);
	}

	.filter-label {
		display: block;
		font-size: 11px;
		font-weight: 700;
		color: rgba(255, 255, 255, 0.5);
		text-transform: uppercase;
		letter-spacing: 1px;
		margin-left: 4px;
	}

	:deep(.ios-item-card) {
		background: rgba(30, 41, 59, 0.7) !important;
		backdrop-filter: blur(16px);
		border: 1px solid rgba(255, 255, 255, 0.15) !important;
		border-radius: 20px !important;
		transition: all 0.3s ease;
		box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3) !important;
	}
	:deep(.ios-item-card:hover) {
		transform: translateY(-5px);
		background: rgba(30, 41, 59, 0.9) !important;
		border-color: #00B8D4 !important;
	}

	.item-title {
		font-size: 1.1rem;
		font-weight: 800;
		letter-spacing: -0.5px;
		color: rgb(var(--v-theme-on-surface));
		line-height: 1.2;
		height: 2.4em;
		overflow: hidden;
		display: -webkit-box;
		-webkit-line-clamp: 2;
		line-clamp: 2;
		-webkit-box-orient: vertical;
	}

	.item-details { color: rgba(var(--v-theme-on-surface), 0.6); font-size: 13px; font-weight: 500; }

	.ios-btn-open { font-size: 11px !important; font-weight: 700 !important; }

	.ios-item-card-skeleton {
		background: rgba(255, 255, 255, 0.04) !important;
		border-radius: 24px !important;
	}
</style>
