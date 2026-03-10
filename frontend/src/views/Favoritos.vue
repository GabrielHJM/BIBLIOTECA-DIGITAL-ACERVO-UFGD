<template>
	<div class="favoritos-container">
		<v-container class="max-width-content py-8">
			<div class="d-flex align-center justify-space-between mb-8">
				<div class="d-flex align-center">
					<v-btn icon="mdi-arrow-left" variant="text" @click="$router.go(-1)" class="mr-4 text-white opacity-60"></v-btn>
					<h1 class="text-h4 font-weight-bold text-white d-flex align-center">
						<v-icon color="pink" class="mr-3">mdi-heart</v-icon>
						Meus Favoritos
					</h1>
				</div>
				<v-text-field
					v-model="searchQuery"
					prepend-inner-icon="mdi-magnify"
					placeholder="Buscar nos favoritos..."
					variant="solo"
					class="ios-search-bar max-width-search"
					hide-details
					density="compact"
				></v-text-field>
			</div>

			<v-row v-if="filteredFavoritos.length > 0">
				<v-col
					v-for="(livro, index) in filteredFavoritos"
					:key="livro.id"
					cols="12"
					sm="6"
					md="4"
					lg="3"
					class="pa-3"
				>
					<div
						class="fav-card-premium animate-fade-in"
						:style="{ animationDelay: (index * 0.05) + 's' }"
						@click="$router.push('/estudo/' + livro.id)"
					>
						<v-img
							:src="livro.capa_url || 'https://images.unsplash.com/photo-1544947950-fa07a98d237f?q=80&w=200'"
							cover
							height="220"
							class="fav-cover"
						>
							<div class="fav-badge">
								<v-icon size="14" color="white">mdi-heart</v-icon>
							</div>
						</v-img>
						<div class="pa-4">
							<div class="text-subtitle-1 font-weight-bold truncate-2-lines text-white mb-1">{{ livro.titulo }}</div>
							<div class="text-caption opacity-60 text-white truncate-1-line">{{ livro.autor }}</div>

							<div class="d-flex align-center justify-space-between mt-3">
								<v-chip size="x-small" variant="tonal" color="primary" class="font-weight-bold">{{ livro.categoria }}</v-chip>
								<v-btn
									icon="mdi-heart-broken"
									variant="text"
									size="small"
									color="pink"
									@click.stop="removerFavorito(livro)"
									title="Remover dos favoritos"
								></v-btn>
							</div>
						</div>
					</div>
				</v-col>
			</v-row>

			<div v-else-if="!loading" class="empty-favorites-state py-16 text-center">
				<v-icon size="80" color="rgba(255,255,255,0.1)" class="mb-4">mdi-heart-off-outline</v-icon>
				<h3 class="text-h5 text-white opacity-40">Nenhum favorito encontrado</h3>
				<p class="text-body-2 text-white opacity-20 mt-2">Explore a biblioteca para adicionar novos livros aqui.</p>
				<v-btn color="primary" variant="outlined" class="mt-8 rounded-pill" to="/explorar">Explorar Biblioteca</v-btn>
			</div>

			<v-row v-else justify="center" class="mt-12">
				<v-progress-circular indeterminate color="primary" size="64"></v-progress-circular>
			</v-row>
		</v-container>
	</div>
</template>

<script>
import MaterialService from '@/services/MaterialService'
import auth from '@/auth'

export default {
	name: 'FavoritosPage',
	inject: ['notify', 'fetchGlobalFavorites', 'getGlobalFavorites', 'globalFavorites'],
	data: () => ({
		loading: true,
		searchQuery: '',
		user: null
	}),
	computed: {
		favoritos() {
			return this.globalFavorites?.list || [];
		},
		filteredFavoritos() {
			if (!this.searchQuery) return this.favoritos;
			const q = this.searchQuery.toLowerCase();
			return this.favoritos.filter(f =>
				f.titulo.toLowerCase().includes(q) ||
				f.autor.toLowerCase().includes(q)
			);
		}
	},
	created() {
		this.user = auth.getUser();
		if (this.user) {
			this.buscarFavoritos();
		} else {
			this.$router.push('/login');
		}
	},
	methods: {
		async buscarFavoritos() {
			this.loading = true;
			try {
				await this.fetchGlobalFavorites();
			} catch (error) {
				console.error('Erro ao buscar favoritos:', error);
				this.notify('Não foi possível carregar seus favoritos.', 'error');
			} finally {
				this.loading = false;
			}
		},
		async removerFavorito(livro) {
			try {
				await MaterialService.favoritar(this.user.id, livro.id, false);
				await this.fetchGlobalFavorites();
				this.notify('Material removido dos favoritos.', 'info');
			} catch (error) {
				console.error('Erro ao remover favorito:', error);
				this.notify('Erro ao remover favorito.', 'error');
			}
		}
	}
}
</script>

<style scoped>
.favoritos-container {
	min-height: 100vh;
	background: transparent;
	padding-top: 60px;
}

.max-width-content {
	max-width: 1400px !important;
}

.max-width-search {
	max-width: 300px;
}

.ios-search-bar {
	background: rgba(255, 255, 255, 0.05) !important;
	border-radius: 12px !important;
}

.fav-card-premium {
	background: rgba(255, 255, 255, 0.03);
	border: 1px solid rgba(255, 255, 255, 0.05);
	border-radius: 24px;
	overflow: hidden;
	transition: all 0.4s cubic-bezier(0.16, 1, 0.3, 1);
	cursor: pointer;
	height: 100%;
}

.fav-card-premium:hover {
	transform: translateY(-8px);
	background: rgba(255, 255, 255, 0.07);
	border-color: rgba(0, 122, 255, 0.3);
	box-shadow: 0 20px 40px rgba(0, 0, 0, 0.3);
}

.fav-cover {
	border-radius: 0 0 16px 16px;
	transition: transform 0.6s ease;
}

.fav-card-premium:hover .fav-cover {
	transform: scale(1.05);
}

.fav-badge {
	position: absolute;
	top: 12px;
	right: 12px;
	background: rgba(233, 30, 99, 0.8);
	width: 28px;
	height: 28px;
	border-radius: 50%;
	display: flex;
	align-items: center;
	justify-content: center;
	backdrop-filter: blur(4px);
}

.truncate-2-lines {
	display: -webkit-box;
	-webkit-line-clamp: 2;
	-webkit-box-orient: vertical;
	line-clamp: 2;
	overflow: hidden;
	height: 3em;
	line-height: 1.5;
}

.truncate-1-line {
	white-space: nowrap;
	overflow: hidden;
	text-overflow: ellipsis;
}

.animate-fade-in {
	opacity: 0;
	animation: fadeIn 0.8s ease forwards;
}

@keyframes fadeIn {
	from { opacity: 0; transform: translateY(10px); }
	to { opacity: 1; transform: translateY(0); }
}

@media (max-width: 600px) {
	.max-width-search {
		max-width: 100%;
		margin-top: 16px;
	}
	.d-flex.align-center.justify-space-between {
		flex-direction: column;
		align-items: flex-start !important;
	}
}
</style>
