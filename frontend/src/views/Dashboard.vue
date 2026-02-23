<template>
	<div class="dashboard-container mt-8">
		<v-container fluid class="pa-4">
			<!-- Header Section -->
			<v-row align="center" class="mb-8 px-4">
				<v-col cols="12" md="6">
					<h1 class="text-white text-h3 font-weight-bold tracking-tight">Recentes</h1>
				</v-col>
				<v-col cols="12" md="6" class="text-right">
					<v-btn class="ios-filter-btn" elevation="2">
						<span>Filtrar</span>
						<v-icon class="ml-2">mdi-filter-variant</v-icon>
					</v-btn>
				</v-col>
			</v-row>

			<!-- Content Grid -->
			<v-row class="px-2">
				<v-col
					v-for="(livro, index) in livros.slice(0, 6)"
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
									<img :src="livro.link" :alt="livro.nome" class="book-cover" />
								</div>
							</v-col>

							<!-- Info Section -->
							<v-col cols="7" class="pa-4 text-left">
								<h3 class="item-title mb-2">{{ livro.nome }}</h3>
								<div class="item-details">
									<p><strong>Publicado:</strong> {{ livro.publicado || '05/10/2022' }}</p>
									<p><strong>Páginas:</strong> {{ livro.paginas || 220 }}</p>
									<p><strong>Autor:</strong> {{ livro.autor || 'Michael T. Goodrich' }}</p>
								</div>

								<div class="mt-4">
									<v-rating
										:model-value="4"
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
							<v-btn icon="mdi-heart-outline" variant="text" color="white" size="small"></v-btn>
							<div class="share-actions">
								<v-btn icon="mdi-share-variant" variant="text" color="white" size="small"></v-btn>
								<v-btn class="ios-btn-open" variant="flat" size="small" @click="$router.push('/estudo')">
									Ler Agora
								</v-btn>
							</div>
						</v-card-actions>
					</v-card>
				</v-col>
			</v-row>
		</v-container>
	</div>
</template>

<script>
import livros from '../../../livros.json'
import auth from '@/auth'

export default {
	name: 'DashboardPage',
	data: () => ({
		livros: livros,
		user: {}
	}),
	created() {
		this.user = auth.getUser()
	},
	methods: {
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

	/* Item Card Style */
	.ios-item-card {
		background: rgba(45, 45, 45, 0.6) !important;
		backdrop-filter: blur(20px);
		border-radius: 24px !important;
		border: 1px solid rgba(255, 255, 255, 0.1);
		overflow: hidden;
		transition: all 0.4s var(--spring-easing);
		opacity: 0;
		animation: ios-reveal 0.6s var(--spring-easing) forwards;
	}

	.ios-item-card:hover {
		transform: translateY(-8px) scale(1.02);
		background: rgba(55, 55, 55, 0.8) !important;
		box-shadow: 0 20px 40px rgba(0,0,0,0.3) !important;
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
		color: white;
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
		color: rgba(255, 255, 255, 0.7);
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

	:deep(.v-rating__item) {
		padding: 0 !important;
	}
</style>
