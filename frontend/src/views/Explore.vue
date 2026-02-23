<template>
	<div class="explore-container mt-8">
		<v-container fluid class="pa-4">
			<!-- Header Section -->
			<v-row align="center" class="mb-6 px-4">
				<v-col cols="12" sm="4" md="4" class="py-2">
					<h1 class="text-white text-h4 text-sm-h3 font-weight-bold tracking-tight">Explorar</h1>
				</v-col>
				<v-col cols="12" sm="5" md="5" class="py-2">
					<v-text-field
						v-model="searchQuery"
						prepend-inner-icon="mdi-magnify"
						placeholder="Busque por título, autor ou categoria..."
						variant="solo"
						class="ios-search-bar"
						hide-details
						@keyup.enter="buscar"
					></v-text-field>
				</v-col>
				<v-col cols="12" sm="3" md="3" class="text-right py-2">
					<v-btn class="ios-filter-btn w-100 w-sm-auto" elevation="2" @click="buscar" :loading="loading">
						<span>Buscar</span>
						<v-icon class="ml-2">mdi-magnify</v-icon>
					</v-btn>
				</v-col>
			</v-row>

			<!-- Content Grid -->
			<v-row class="px-2" v-if="!loading">
				<v-col
					v-for="(livro, index) in livros"
					:key="index"
					cols="12"
					md="6"
					lg="4"
					class="pa-4"
				>
					<v-card class="ios-item-card" elevation="8" :style="{ animationDelay: `${index * 100}ms` }" @click="$router.push('/estudo/' + livro.id)">
						<v-row no-gutters>
							<!-- Image Section -->
							<v-col cols="5" class="pa-3">
								<div class="book-cover-wrapper">
									<img :src="livro.capa_url || 'https://images.unsplash.com/photo-1544947950-fa07a98d237f?q=80&w=200'" :alt="livro.titulo" class="book-cover" />
								</div>
							</v-col>

							<!-- Info Section -->
							<v-col cols="7" class="pa-4 text-left">
								<h3 class="item-title mb-2">{{ livro.titulo }}</h3>
								<div class="item-details">
									<p><strong>Autor:</strong> {{ livro.autor }}</p>
									<p><strong>Categoria:</strong> {{ livro.categoria }}</p>
								</div>

								<div class="mt-4">
									<v-rating
										:model-value="livro.media_nota || 0"
										density="compact"
										color="amber"
										active-color="amber"
										size="small"
										readonly
									></v-rating>
								</div>
							</v-col>
						</v-row>

						<v-divider class="mx-4 opacity-50"></v-divider>
						<v-card-actions class="pa-3 justify-end">
							<v-btn class="ios-btn-open" variant="flat" size="small">
								Ver Detalhes
							</v-btn>
						</v-card-actions>
					</v-card>
				</v-col>

				<v-col cols="12" v-if="livros.length === 0" class="text-center mt-12">
					<v-icon size="64" color="white" class="opacity-20 mb-4">mdi-book-search-outline</v-icon>
					<h2 class="text-white opacity-60">Nenhum material encontrado</h2>
					<p class="text-white opacity-40">Tente buscar por termos diferentes ou navegue por categorias.</p>
				</v-col>
			</v-row>

			<v-row v-else justify="center" class="mt-12">
				<v-progress-circular indeterminate color="cyan" size="64"></v-progress-circular>
			</v-row>
		</v-container>
	</div>
</template>

<script>
import MaterialService from '@/services/MaterialService'

export default {
	name: 'ExplorePage',
	data: () => ({
		livros: [],
		searchQuery: '',
		loading: false
	}),
	watch: {
		'$route.query.q': {
			immediate: true,
			handler(val) {
				if (val) {
					this.searchQuery = val
					this.buscar()
				}
			}
		},
		'$route.query.categoria': {
			immediate: true,
			handler(val) {
				if (val) {
					this.buscarPorCategoria(val)
				}
			}
		}
	},
	created() {
		if (!this.$route.query.q && !this.$route.query.categoria) {
			this.buscar()
		}
	},
	methods: {
		async buscar() {
			this.loading = true
			try {
				const response = await MaterialService.pesquisar(this.searchQuery)
				this.livros = response.data || []
			} catch (error) {
				console.error('Erro ao buscar materiais:', error)
			} finally {
				this.loading = false
			}
		},
		async buscarPorCategoria(cat) {
			this.loading = true
			try {
				const response = await MaterialService.pesquisar('', cat)
				this.livros = response.data || []
			} catch (error) {
				console.error('Erro ao buscar por categoria:', error)
			} finally {
				this.loading = false
			}
		}
	}
}
</script>

<style scoped>
	.explore-container {
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

	.ios-item-card {
		background: rgba(45, 78, 115, 0.5) !important;
		backdrop-filter: blur(20px);
		border-radius: 24px !important;
		border: 1px solid rgba(255, 255, 255, 0.1);
		overflow: hidden;
		transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
		opacity: 0;
		animation: ios-reveal 0.6s cubic-bezier(0.4, 0, 0.2, 1) forwards;
		cursor: pointer;
	}

	.ios-item-card:hover {
		transform: translateY(-8px);
		background: rgba(45, 78, 115, 0.7) !important;
		box-shadow: 0 20px 40px rgba(0,0,0,0.3) !important;
	}

	.book-cover-wrapper {
		width: 100%;
		height: 160px;
		border-radius: 12px;
		overflow: hidden;
		background: rgba(0,0,0,0.2);
		display: flex;
		align-items: center;
		justify-content: center;
	}

	.book-cover {
		max-width: 100%;
		max-height: 100%;
		object-fit: contain;
	}

	.item-title {
		color: white;
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
		color: rgba(255, 255, 255, 0.7);
		font-size: 13px;
	}

	.ios-btn-open {
		background: #00B8D4 !important;
		color: white !important;
		border-radius: 10px !important;
		text-transform: none !important;
		font-weight: 600 !important;
	}

	.opacity-20 { opacity: 0.2; }
	.opacity-40 { opacity: 0.4; }
	.opacity-60 { opacity: 0.6; }

	@keyframes ios-reveal {
		from { opacity: 0; transform: scale(0.95) translateY(20px); }
		to { opacity: 1; transform: scale(1) translateY(0); }
	}
</style>
