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
						<div class="premium-book-cover-ios-fav">
							<div class="mesh-gradient"></div>
							<div class="glass-overlay-ios-fav">
								<div class="cover-content-ios-fav">
									<v-icon class="cover-icon-ios mb-1" size="32" color="rgba(255,255,255,0.8)">{{ getBookIcon(livro.categoria, livro.titulo) }}</v-icon>
									<div class="cover-title-ios-fav">{{ livro.titulo }}</div>
									<div class="cover-summary-ios-fav">{{ livro.resumo }}</div>
								</div>
							</div>
							<div class="glass-shine"></div>
							
							<div class="fav-badge">
								<v-icon size="14" color="white">mdi-heart</v-icon>
							</div>
						</div>
						<div class="pa-4">
							<div class="text-subtitle-1 font-weight-bold truncate-2-lines text-white mb-1">{{ livro.titulo }}</div>
							
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
				<p class="text-body-2 text-white opacity-20 mt-2">Explore o Acervus Core para adicionar novos livros aqui.</p>
				<v-btn color="primary" variant="outlined" class="mt-8 rounded-pill" to="/explorar">Explorar Acervo</v-btn>
			</div>

			<v-row v-else justify="center" class="mt-12">
				<CloudLoader text="Carregando favoritos..." height="200px" />
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
	backdrop-filter: blur(10px);
}

.fav-card-premium {
	background: rgba(255, 255, 255, 0.03) !important;
	border: 1px solid rgba(255, 255, 255, 0.1);
	border-radius: 20px;
	overflow: hidden;
	transition: all 0.4s cubic-bezier(0.16, 1, 0.3, 1);
	cursor: pointer;
	height: 100%;
	backdrop-filter: blur(8px);
}

.fav-card-premium:hover {
	transform: translateY(-8px);
	background: rgba(255, 255, 255, 0.07);
	border-color: rgba(0, 122, 255, 0.3);
	box-shadow: 0 20px 40px rgba(0, 0, 0, 0.4);
}

.premium-book-cover-ios-fav {
	width: 100%;
	aspect-ratio: 1 / 1.4;
	max-height: 240px;
	margin: 0 auto;
	position: relative;
	overflow: hidden;
	border-radius: 20px 20px 0 0;
	background: #000;
	display: flex;
	align-items: center;
	justify-content: center;
	box-shadow: 0 10px 30px rgba(0, 0, 0, 0.4);
	transition: all 0.6s cubic-bezier(0.16, 1, 0.3, 1);
}

.mesh-gradient {
	position: absolute;
	top: -50%;
	left: -50%;
	width: 200%;
	height: 200%;
	background: 
		radial-gradient(circle at 20% 20%, #007AFF 0%, transparent 40%),
		radial-gradient(circle at 80% 25%, #5AC8FA 0%, transparent 40%),
		radial-gradient(circle at 30% 80%, #0051FF 0%, transparent 40%),
		radial-gradient(circle at 75% 75%, #00D4E8 0%, transparent 40%),
		radial-gradient(circle at 50% 50%, #0033FF 0%, transparent 50%);
	filter: blur(40px);
	animation: mesh-rotate 15s linear infinite;
	opacity: 0.7;
}

@keyframes mesh-rotate {
	from { transform: rotate(0deg); }
	to { transform: rotate(360deg); }
}

.glass-overlay-ios-fav {
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
	padding: 16px;
	text-align: center;
}

.cover-content-ios-fav {
	position: relative;
	z-index: 3;
}

.cover-title-ios-fav {
	color: #ffffff;
	font-family: 'Outfit', 'Inter', sans-serif;
	font-size: 0.85rem;
	font-weight: 800;
	line-height: 1.2;
	margin-bottom: 6px;
	display: -webkit-box;
	-webkit-line-clamp: 3;
	-webkit-box-orient: vertical;
	overflow: hidden;
	text-shadow: 0 2px 10px rgba(0, 0, 0, 0.5);
}

.cover-summary-ios-fav {
	color: rgba(255, 255, 255, 0.6);
	font-size: 0.65rem;
	font-weight: 500;
	line-height: 1.3;
	display: -webkit-box;
	-webkit-line-clamp: 5;
	-webkit-box-orient: vertical;
	overflow: hidden;
	text-align: center;
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

.fav-card-premium:hover .premium-book-cover-ios-fav {
	transform: scale(1.02) translateY(-2px);
}

.fav-card-premium:hover .mesh-gradient {
	opacity: 1;
	animation-duration: 8s;
}

.fav-badge {
	position: absolute;
	top: 12px;
	right: 12px;
	background: rgba(233, 30, 99, 0.9);
	width: 28px;
	height: 28px;
	border-radius: 50%;
	display: flex;
	align-items: center;
	justify-content: center;
	backdrop-filter: blur(4px);
	z-index: 10;
}

.truncate-2-lines {
	display: -webkit-box;
	-webkit-line-clamp: 2;
	-webkit-box-orient: vertical;
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
