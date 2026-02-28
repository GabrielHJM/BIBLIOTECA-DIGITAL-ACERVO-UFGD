<template>
	<div class="explore-container mt-8">
		<v-container fluid class="pa-4">
			<!-- Advanced Header Section -->
			<v-row align="center" class="mb-6 px-4">
				<v-col cols="12" md="8" lg="9" class="py-2">
					<div class="search-premium-box d-flex align-center">
						<v-text-field
							v-model="filters.q"
							prepend-inner-icon="mdi-magnify"
							placeholder="Busque por título, autor, DOI ou assunto..."
							variant="solo"
							class="ios-search-bar flex-grow-1"
							hide-details
							@update:modelValue="onSearchInput"
							@keyup.enter="buscar"
						></v-text-field>
						<v-btn
							class="search-btn ml-4"
							elevation="4"
							@click="buscar"
							:loading="loading"
							height="56"
							rounded="xl"
						>
							<v-icon class="mr-2">mdi-magnify</v-icon>
							EXPLORAR
						</v-btn>
					</div>
				</v-col>
				<v-col cols="12" md="4" lg="3" class="py-2">
					<v-btn
						class="ios-filter-toggle w-100"
						variant="outlined"
						prepend-icon="mdi-tune-variant"
						@click="showFilters = !showFilters"
						height="56"
						rounded="xl"
					>
						Filtros Avançados
						<v-chip size="x-small" color="cyan" class="ml-2" v-if="activeFiltersCount > 0">
							{{ activeFiltersCount }}
						</v-chip>
					</v-btn>
				</v-col>
			</v-row>

			<!-- Expandable Filter Panel -->
			<v-expand-transition>
				<div v-show="showFilters" class="filter-panel px-4 mb-6">
					<v-card class="ios-glass-card pa-6" rounded="xl">
						<v-row>
							<!-- Source Filter -->
							<v-col cols="12" sm="4">
								<label class="filter-label">Fonte do Material</label>
								<v-select
									v-model="filters.fonte"
									:items="fontesList"
									placeholder="Todas as Fontes"
									variant="solo-inverted"
									density="comfortable"
									hide-details
									class="mt-2"
									rounded="lg"
									clearable
								></v-select>
							</v-col>
							<!-- Year Range -->
							<v-col cols="12" sm="4">
								<label class="filter-label">Ano de Publicação (A partir de)</label>
								<v-select
									v-model="filters.ano_inicio"
									:items="yearsList"
									placeholder="Qualquer ano"
									variant="solo-inverted"
									density="comfortable"
									hide-details
									class="mt-2"
									rounded="lg"
									clearable
								></v-select>
							</v-col>
							<v-col cols="12" sm="4">
								<label class="filter-label">Ordenar Por</label>
								<v-select
									v-model="filters.sort"
									:items="sortOptions"
									item-title="label"
									item-value="value"
									placeholder="Relevância"
									variant="solo-inverted"
									density="comfortable"
									hide-details
									class="mt-2"
									rounded="lg"
									@update:modelValue="buscar"
								></v-select>
							</v-col>
						</v-row>
						<v-row dense class="mt-4">
							<v-col cols="12" class="d-flex justify-end gap-2">
								<v-btn variant="text" size="small" @click="limparFiltros">Limpar</v-btn>
								<v-btn color="cyan" size="small" rounded="pill" @click="buscar">Aplicar Filtros</v-btn>
							</v-col>
						</v-row>
					</v-card>
				</div>
			</v-expand-transition>

			<!-- Quick Categories -->
			<v-row class="px-4 mb-4" v-if="!showFilters">
				<v-col cols="12">
					<v-slide-group v-model="filters.categoria" show-arrows @update:modelValue="buscar">
						<v-slide-group-item v-slot="{ isSelected, toggle }" value="">
							<v-chip :color="isSelected ? 'cyan' : 'white'" :variant="isSelected ? 'flat' : 'outlined'" class="ma-2 px-4" @click="toggle" size="large">
								Todos
							</v-chip>
						</v-slide-group-item>
						<v-slide-group-item v-for="cat in categoriasMock" :key="cat" :value="cat" v-slot="{ isSelected, toggle }">
							<v-chip :color="isSelected ? 'cyan' : 'white'" :variant="isSelected ? 'flat' : 'outlined'" class="ma-2 px-4" @click="toggle" size="large">
								{{ cat }}
							</v-chip>
						</v-slide-group-item>
					</v-slide-group>
				</v-col>
			</v-row>

			<!-- Content Grid -->
			<v-row class="px-2" v-if="!loading">
				<v-col v-for="(livro, index) in livros" :key="livro.id" cols="12" md="6" lg="4" class="pa-4">
					<BookCard
						:book="livro"
						:animation-delay="index * 50"
						@toggle-favorite="onToggleFavorite"
						@share="shareBook(livro.id)"
					/>
				</v-col>

				<v-col cols="12" v-if="livros.length === 0" class="text-center mt-12">
					<v-icon size="80" color="cyan" class="mb-4 opacity-20">mdi-text-search-variant</v-icon>
					<h2 class="opacity-60 text-white">Nenhum resultado encontrado</h2>
					<p class="opacity-40 text-white px-4">Refine seus filtros ou tente termos mais abrangentes.</p>
				</v-col>
			</v-row>

			<!-- Skeleton Loading State -->
			<v-row class="px-2" v-else>
				<v-col v-for="n in 6" :key="n" cols="12" md="6" lg="4" class="pa-4">
					<v-skeleton-loader
						class="ios-item-card-skeleton"
						type="image, article, actions"
						theme="dark"
					></v-skeleton-loader>
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
			fonte: '',
			ano_inicio: null,
			ano_fim: null
		},
		fontesList: ['SciELO', 'CAPES', 'IEEE'],
		categoriasMock: ['Ciência', 'Tecnologia', 'Educação', 'Medicina', 'Engenharia', 'Filosofia'],
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
			if (this.filters.fonte) count++;
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
					this.filters.fonte,
					this.filters.ano_inicio,
					this.filters.ano_fim,
					null, // tags
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
			this.filters = { q: this.filters.q, categoria: '', fonte: '', ano_inicio: null, ano_fim: null };
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
	.explore-container { min-height: 100vh; padding-bottom: 80px; }
	.search-premium-box { width: 100%; }

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

	.ios-item-card {
		background: rgba(var(--v-theme-surface), 0.7) !important;
		backdrop-filter: blur(20px) saturate(180%);
		-webkit-backdrop-filter: blur(20px) saturate(180%);
		border-radius: 24px !important;
		border: 1px solid rgba(var(--v-border-color), 0.1);
		transition: all 0.4s var(--spring-easing);
		position: relative;
		overflow: hidden;
		opacity: 0;
		animation: revealCard 0.6s var(--spring-easing) forwards;
	}
	.ios-item-card:hover {
		transform: translateY(-8px) scale(1.02);
		background: rgba(var(--v-theme-surface), 0.9) !important;
		border-color: var(--ios-cyan);
		box-shadow: 0 20px 40px rgba(0,0,0,0.15) !important;
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
		-webkit-box-orient: vertical;
	}

	.item-details { color: rgba(var(--v-theme-on-surface), 0.6); font-size: 13px; font-weight: 500; }

	.ios-btn-open { font-size: 11px !important; font-weight: 700 !important; }

	.ios-item-card-skeleton {
		background: rgba(255, 255, 255, 0.04) !important;
		border-radius: 24px !important;
	}
</style>
