<template>
	<div class="home-container mt-8">
		<v-container fluid class="pa-4">
			<v-row>
				<!-- Main Content: Suggested Materials -->
				<v-col cols="12" lg="8" class="pa-4">
					<div class="section-header mb-6">
						<h1 class="text-white text-h3 font-weight-bold tracking-tight">Comece seus Estudos</h1>
						<p class="text-white opacity-70 mt-2">Materiais selecionados para sua jornada acadêmica</p>
					</div>

					<v-row>
						<v-col
							v-for="(livro, index) in livros.slice(0, 4)"
							:key="index"
							cols="12"
							md="6"
							class="pa-3"
						>
							<v-card class="ios-item-card" elevation="8" :style="{ animationDelay: `${index * 100}ms` }">
								<v-row no-gutters>
									<v-col cols="5" class="pa-3">
										<div class="book-cover-wrapper">
											<img :src="livro.link" :alt="livro.nome" class="book-cover" />
										</div>
									</v-col>
									<v-col cols="7" class="pa-4 text-left">
										<h4 class="item-category mb-1">{{ livro.categoria || 'Sistemas de Inf' }}</h4>
										<h3 class="item-title mb-2">{{ livro.nome }}</h3>
										<div class="item-details">
											<p><strong>Publicado:</strong> {{ livro.publicado || '05/10/2022' }}</p>
											<p><strong>Páginas:</strong> {{ livro.paginas || 220 }}</p>
										</div>
										<v-btn class="ios-btn-ghost mt-4 px-0" variant="text" size="small" @click="$router.push('/login')">
											Ver detalhes <v-icon size="small">mdi-chevron-right</v-icon>
										</v-btn>
									</v-col>
								</v-row>
							</v-card>
						</v-col>
					</v-row>
				</v-col>

				<!-- Sidebar: Why us? -->
				<v-col cols="12" lg="4" class="pa-4">
					<v-card class="ios-card-premium pa-6 benefits-card" elevation="12">
						<h2 class="text-white text-h4 font-weight-bold mb-8 leading-tight">
							O que torna a <span class="text-cyan">Biblioteca Digital</span> excelente?
						</h2>

						<div class="benefits-list">
							<div class="benefit-item d-flex align-center mb-6">
								<div class="benefit-icon-wrapper mr-4">
									<img src="../assets/images/site-images/home/infinite.png" alt="Infinite" />
								</div>
								<div>
									<h4 class="text-white font-weight-bold">Sem Limites</h4>
									<p class="text-white opacity-60 text-caption">Acesso ilimitado a todo o acervo</p>
								</div>
							</div>

							<div class="benefit-item d-flex align-center mb-6">
								<div class="benefit-icon-wrapper mr-4">
									<img src="../assets/images/site-images/home/4.png" alt="Instant" />
								</div>
								<div>
									<h4 class="text-white font-weight-bold">Acesso Instantâneo</h4>
									<p class="text-white opacity-60 text-caption">Estude de qualquer lugar, a qualquer hora</p>
								</div>
							</div>

							<div class="benefit-item d-flex align-center mb-6">
								<div class="benefit-icon-wrapper mr-4">
									<img src="../assets/images/site-images/home/2.png" alt="Custom" />
								</div>
								<div>
									<h4 class="text-white font-weight-bold">Personalização</h4>
									<p class="text-white opacity-60 text-caption">Sua própria estante virtual organizada</p>
								</div>
							</div>

							<div class="benefit-item d-flex align-center mb-6">
								<div class="benefit-icon-wrapper mr-4">
									<img src="../assets/images/site-images/home/5.png" alt="Free" />
								</div>
								<div>
									<h4 class="text-white font-weight-bold">Gratuito</h4>
									<p class="text-white opacity-60 text-caption">Conhecimento livre para todos</p>
								</div>
							</div>
						</div>

						<v-divider class="my-8 opacity-20"></v-divider>

						<div class="cta-section">
							<p class="text-white font-italic mb-6">"O mundo acadêmico está ao seu alcance."</p>
							<v-btn class="ios-action-btn w-100" elevation="6" @click="$router.push('/cadastro')">
								Cadastre-se Agora
							</v-btn>
						</div>
					</v-card>
				</v-col>
			</v-row>
		</v-container>
	</div>
</template>

<script>
import livros from '../../../livros.json'

export default {
	name: 'HomePage',
	data: () => ({
		livros: livros
	})
}
</script>

<style scoped>
	.home-container {
		min-height: 100vh;
		padding-bottom: 80px;
	}

	.tracking-tight {
		letter-spacing: -1.5px !important;
	}

	.opacity-70 { opacity: 0.7; }
	.opacity-60 { opacity: 0.6; }
	.opacity-20 { opacity: 0.2; }
	.text-cyan { color: var(--ios-cyan); }

	.ios-card-premium {
		background: rgba(58, 99, 145, 0.6) !important;
		backdrop-filter: blur(20px);
		border-radius: 32px !important;
		border: 1px solid rgba(255, 255, 255, 0.1);
	}

	.benefits-card {
		height: 100%;
		display: flex;
		flex-direction: column;
		justify-content: center;
	}

	.benefit-icon-wrapper {
		width: 54px;
		height: 54px;
		background: rgba(255, 255, 255, 0.1);
		border-radius: 14px;
		display: flex;
		align-items: center;
		justify-content: center;
		padding: 10px;
	}

	.benefit-icon-wrapper img {
		width: 100%;
		height: 100%;
		object-fit: contain;
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
	}

	.book-cover {
		max-width: 100%;
		max-height: 100%;
		object-fit: contain;
		transition: transform 0.5s ease;
	}

	.ios-item-card:hover .book-cover {
		transform: scale(1.1);
	}

	.item-category {
		color: var(--ios-cyan);
		font-size: 12px;
		text-transform: uppercase;
		font-weight: 800;
		letter-spacing: 1px;
	}

	.item-title {
		color: white;
		font-size: 1.2rem;
		font-weight: 700;
		line-height: 1.2;
	}

	.item-details {
		color: rgba(255, 255, 255, 0.7);
		font-size: 13px;
	}

	.ios-btn-ghost {
		color: white !important;
		text-transform: none !important;
		font-weight: 600 !important;
	}

	.ios-action-btn {
		background: linear-gradient(135deg, #00B8D4 0%, #0097A7 100%) !important;
		color: white !important;
		border-radius: 16px !important;
		height: 54px !important;
		text-transform: none !important;
		font-size: 18px !important;
		font-weight: 700 !important;
	}

	@keyframes ios-reveal {
		from { opacity: 0; transform: scale(0.9) translateY(30px); }
		to { opacity: 1; transform: scale(1) translateY(0); }
	}
</style>
