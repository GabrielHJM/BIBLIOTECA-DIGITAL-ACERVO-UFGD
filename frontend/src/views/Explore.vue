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
								<label class="filter-label">Categoria Principal</label>
								<v-select
									v-model="filters.categoria"
									:items="categoriasMock"
									placeholder="Todas as Categorias"
									variant="solo-inverted"
									density="comfortable"
									hide-details
									class="mt-2"
									rounded="lg"
									clearable
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
				<v-col v-for="(livro, index) in livros" :key="index" cols="12" md="6" lg="4" class="pa-4">
					<v-card class="ios-item-card" elevation="8" @click="$router.push('/estudo/' + livro.id)">
						<div class="source-badge" v-if="livro.fonte">
							{{ livro.fonte }}
						</div>
						<v-row no-gutters>
							<v-col cols="5" class="pa-3">
								<div class="book-cover-wrapper premium-shadow">
									<img :src="livro.capa_url || 'https://images.unsplash.com/photo-1544947950-fa07a98d237f?q=80&w=200'" :alt="livro.titulo" class="book-cover" />
								</div>
							</v-col>
							<v-col cols="7" class="pa-4 text-left">
								<h3 class="item-title mb-2">{{ livro.titulo }}</h3>
								<div class="item-details mt-2">
									<div class="d-flex align-center mb-1">
										<v-icon size="14" class="mr-1">mdi-account-edit</v-icon>
										<span class="text-truncate">{{ livro.autor }}</span>
									</div>
									<div class="d-flex align-center mb-1">
										<v-icon size="14" class="mr-1">mdi-tag-outline</v-icon>
										<span>{{ livro.categoria }}</span>
									</div>
									<div class="d-flex align-center" v-if="livro.ano_publicacao">
										<v-icon size="14" class="mr-1">mdi-calendar-range</v-icon>
										<span>{{ livro.ano_publicacao }}</span>
									</div>
								</div>
								<div class="mt-4 d-flex align-center">
									<v-rating :model-value="livro.media_nota || 0" density="compact" color="amber" active-color="amber" size="x-small" readonly></v-rating>
									<span class="ml-2 text-caption opacity-60">({{ livro.total_avaliacoes || 0 }})</span>
								</div>
							</v-col>
						</v-row>
						<v-card-actions class="pa-3 justify-end">
							<v-btn class="ios-btn-open" variant="tonal" color="cyan" size="small" rounded="pill">
								Ler Material
								<v-icon class="ml-1" size="14">mdi-arrow-right</v-icon>
							</v-btn>
						</v-card-actions>
					</v-card>
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

export default {
	name: 'ExplorePageExtended',
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
					this.filters.ano_fim
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
		background: rgba(255, 255, 255, 0.06) !important;
		backdrop-filter: blur(20px);
		border-radius: 24px !important;
		border: 1px solid rgba(255, 255, 255, 0.08);
		transition: all 0.4s cubic-bezier(0.16, 1, 0.3, 1);
		position: relative;
		overflow: hidden;
	}
	.ios-item-card:hover {
		transform: translateY(-8px);
		background: rgba(255, 255, 255, 0.1) !important;
		border-color: rgba(0, 184, 212, 0.4);
		box-shadow: 0 20px 40px rgba(0,0,0,0.3) !important;
	}

	.source-badge {
		position: absolute;
		top: 12px;
		right: 12px;
		background: rgba(0, 184, 212, 0.2);
		color: #00D4E8;
		font-size: 9px;
		font-weight: 800;
		padding: 3px 8px;
		border-radius: 20px;
		border: 1px solid rgba(0, 184, 212, 0.3);
		text-transform: uppercase;
		z-index: 2;
	}

	.book-cover-wrapper {
		width: 100%;
		height: 180px;
		border-radius: 12px;
		overflow: hidden;
		background: rgba(0,0,0,0.2);
	}
	.book-cover { width: 100%; height: 100%; object-fit: contain; padding: 4px; }
	.premium-shadow { box-shadow: 0 10px 25px rgba(0,0,0,0.5); }

	.item-title {
		font-size: 1rem;
		font-weight: 700;
		color: white;
		line-height: 1.3;
		height: 2.6em;
		overflow: hidden;
		display: -webkit-box;
		-webkit-line-clamp: 2;
		-webkit-box-orient: vertical;
	}

	.item-details { color: rgba(255,255,255,0.6); font-size: 12px; }

	.ios-btn-open { font-size: 11px !important; font-weight: 700 !important; }

	.ios-item-card-skeleton {
		background: rgba(255, 255, 255, 0.04) !important;
		border-radius: 24px !important;
	}
</style>
