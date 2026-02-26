<template>
	<div class="dashboard-container mt-8 position-relative">
		<v-snackbar v-model="snackbar" :timeout="3000" color="cyan darken-2" location="top right">
			{{ snackbarMsg }}
		</v-snackbar>
		<v-container fluid class="pa-4">
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
				<v-slide-group show-arrows class="pa-0">
					<v-slide-group-item
						v-for="(livro, index) in historico"
						:key="'hist-'+index"
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
				<v-slide-group show-arrows class="pa-0">
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
					:key="index"
					cols="12"
					md="6"
					lg="4"
					class="pa-4"
				>
					<v-card class="ios-item-card" elevation="8" :style="{ animationDelay: `${index * 100}ms` }">
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
									<p><strong>Publicado:</strong> {{ livro.ano_publicacao }}</p>
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

						<!-- Actions Footer -->
						<v-divider class="mx-4 opacity-50"></v-divider>
						<v-card-actions class="pa-3 justify-space-between">
							<v-btn icon="" variant="text" :color="isFavorited(livro.id) ? 'pink' : ($vuetify.theme.current.dark ? 'white' : 'grey-darken-1')" size="small" @click="toggleFavorite(livro)">
								<v-icon>{{ isFavorited(livro.id) ? 'mdi-heart' : 'mdi-heart-outline' }}</v-icon>
							</v-btn>
							<div class="share-actions">
								<v-btn icon="mdi-share-variant" variant="text" :color="$vuetify.theme.current.dark ? 'white' : 'grey-darken-1'" size="small" @click="shareBook(livro.id)"></v-btn>
								<v-btn class="ios-btn-open" variant="flat" size="small" @click="$router.push('/estudo/' + livro.id)">
									Ler Agora
								</v-btn>
							</div>
						</v-card-actions>
					</v-card>
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
import auth from '@/auth'

export default {
	name: 'DashboardPage',
	data: () => ({
		livros: [],
		historico: [],
		favoritos: [],
		searchQuery: '',
		loading: false,
		user: {},
		snackbar: false,
		snackbarMsg: ''
	}),
	created() {
		this.user = auth.getUser()
		this.buscar()
		if (this.user && this.user.id) {
			this.buscarHistorico()
			this.buscarFavoritos()
		}
	},
	methods: {
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
		logout() {
			auth.logout()
			this.$router.push('/login')
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
		background: rgba(var(--v-theme-surface), 0.8) !important;
		backdrop-filter: blur(30px) saturate(120%);
		border-radius: 28px !important;
		border: 1px solid rgba(128, 128, 128, 0.2);
		overflow: hidden;
		transition: all 0.5s cubic-bezier(0.16, 1, 0.3, 1);
		opacity: 0;
		animation: ios-reveal 0.7s cubic-bezier(0.16, 1, 0.3, 1) forwards;
	}

	.ios-item-card:hover {
		transform: translateY(-10px) scale(1.02);
		background: rgba(var(--v-theme-surface), 0.95) !important;
		box-shadow: 0 30px 60px rgba(0,0,0,0.2) !important;
		border-color: rgba(128, 128, 128, 0.4) !important;
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
</style>
