<template>
	<div class="apple-home-root" :class="{ 'light-theme': !isDarkTheme }">
		
		<!-- Hero Section Apple Style -->
		<section class="apple-hero" ref="heroSection">
			<div class="hero-content-center">
				<div class="pill-badge fade-up-delay-1">A sua biblioteca inteligente</div>
				<h1 class="apple-huge-title fade-up-delay-2">
					Acervus Core.
				</h1>
				<h2 class="apple-subtitle fade-up-delay-3">
					O conhecimento ao seu alcance.<br>
					Sem limites. Sem barreiras.
				</h2>
				
				<div class="hero-actions fade-up-delay-4">
					<template v-if="!isLoggedIn">
						<button class="apple-btn-primary" @click="$router.push('/cadastro')">Começar Agora</button>
						<button class="apple-btn-link" @click="$router.push('/login')">Entrar <v-icon size="16">mdi-chevron-right</v-icon></button>
					</template>
					<template v-else>
						<button class="apple-btn-primary" @click="$router.push('/dashboard')">Ver Meu Dashboard</button>
					</template>
				</div>
			</div>
			
			<!-- Subtle hero glow -->
			<div class="hero-glow"></div>
		</section>

		<!-- Categories (Shelves) Apple Style -->
		<section class="apple-categories-section">
			<div class="section-container">
				<div class="apple-section-header fade-on-scroll">
					<h3 class="apple-section-title">Explore as estantes.</h3>
					<p class="apple-section-desc">Milhares de recursos organizados para você.</p>
					<button class="apple-btn-link mt-2" @click="$router.push({ path: '/explorar', query: { categoria: 'Todos' } })">
						Explorar Tudo <v-icon size="16">mdi-chevron-right</v-icon>
					</button>
				</div>

				<div class="apple-shelves-grid">
					<div 
						v-for="(cat, idx) in categoriasMock" 
						:key="cat.nome"
						class="apple-shelf-card glass-panel fade-on-scroll"
						:style="{ transitionDelay: `${idx * 0.05}s` }"
					>
						<div class="shelf-header" @click="$router.push({ name: 'explorar', query: { categoria: cat.nome } })">
							<div class="shelf-icon-box" :style="{ color: cat.iconColor }">
								<v-icon size="24">{{ cat.icon }}</v-icon>
							</div>
							<h4 class="shelf-title">{{ cat.nome }}</h4>
							<v-spacer></v-spacer>
							<v-icon color="rgba(255,255,255,0.3)">mdi-chevron-right</v-icon>
						</div>

						<div class="shelf-books">
							<template v-if="loading">
								<div v-for="s in 3" :key="s" class="apple-skeleton-book"></div>
							</template>
							<template v-else-if="cat.livros && cat.livros.length > 0">
								<div 
									v-for="livro in cat.livros.slice(0,3)" 
									:key="livro.id"
									class="apple-book-row"
									@click="$router.push('/estudo/' + livro.id)"
								>
									<div class="book-cover-mini">
										<v-icon size="20" color="rgba(255,255,255,0.7)">{{ getBookIcon(livro.categoria, livro.titulo) }}</v-icon>
									</div>
									<div class="book-info">
										<div class="book-title">{{ livro.titulo }}</div>
										<div class="book-author">{{ livro.autor || 'Autor Desconhecido' }}</div>
									</div>
								</div>
							</template>
							<div v-else class="empty-shelf">
								Nenhum livro disponível
							</div>
						</div>
					</div>
				</div>
			</div>
		</section>

		<!-- Features Innovation Apple Style -->
		<section class="apple-features-section">
			<div class="section-container">
				<div class="apple-section-header center-text fade-on-scroll">
					<h3 class="apple-section-title">Inovação e Acessibilidade.</h3>
					<p class="apple-section-desc">Um ecossistema feito para a sua evolução.</p>
				</div>

				<div class="apple-features-grid">
					<div 
						v-for="(feature, fIdx) in features" 
						:key="feature.title"
						class="apple-feature-box glass-panel fade-on-scroll"
						:style="{ transitionDelay: `${fIdx * 0.1}s` }"
					>
						<div class="feature-icon" :style="{ color: feature.iconColor }">
							<v-icon size="32">{{ feature.icon }}</v-icon>
						</div>
						<h4 class="feature-title">{{ feature.title }}</h4>
						<p class="feature-desc">{{ feature.desc }}</p>
					</div>
				</div>
			</div>
		</section>

	</div>
</template>

<script>
import MaterialService from '../services/MaterialService';
import { state as authState } from '@/auth';
import { useTheme } from 'vuetify';
import { computed } from 'vue';

export default {
	name: 'HomeView',
	setup() {
		const theme = useTheme();
		const isDarkTheme = computed(() => theme.global.current.value.dark);
		return { isDarkTheme };
	},
	computed: {
		isLoggedIn() {
			return authState.isAuthenticated;
		}
	},
	data() {
		return {
			loading: true,
			categoriasMock: [
				{ nome: 'TECNOLOGIA', livros: [], icon: 'mdi-laptop', iconColor: '#0A84FF' },
				{ nome: 'SAÚDE', livros: [], icon: 'mdi-heart-pulse', iconColor: '#FF375F' },
				{ nome: 'MATEMÁTICA', livros: [], icon: 'mdi-calculator-variant', iconColor: '#FF9F0A' },
				{ nome: 'CIÊNCIAS', livros: [], icon: 'mdi-flask', iconColor: '#32D74B' },
				{ nome: 'HISTÓRIA', livros: [], icon: 'mdi-castle', iconColor: '#BF5AF2' },
				{ nome: 'CONTABILIDADE', livros: [], icon: 'mdi-currency-usd', iconColor: '#FFD60A' }
			],
			features: [
				{ title: 'Sem Limites', desc: 'Acesse quantos livros quiser, sem restrições e onde estiver.', icon: 'mdi-infinity', iconColor: '#0A84FF' },
				{ title: 'Instantâneo', desc: 'Leitura online sem necessidade de download ou espera.', icon: 'mdi-lightning-bolt', iconColor: '#FFD60A' },
				{ title: 'Personalizado', desc: 'Favoritos, histórico e recomendações sob medida.', icon: 'mdi-auto-fix', iconColor: '#BF5AF2' },
				{ title: 'Gratuito', desc: 'Sempre gratuito e democrático. O conhecimento é de todos.', icon: 'mdi-currency-usd-off', iconColor: '#FF375F' }
			]
		}
	},
	async mounted() {
		await this.fetchMateriais();
		this.setupIntersectionObserver();
	},
	methods: {
		setupIntersectionObserver() {
			const options = {
				root: null,
				rootMargin: '0px',
				threshold: 0.15
			};
			
			const observer = new IntersectionObserver((entries, observer) => {
				entries.forEach(entry => {
					if (entry.isIntersecting) {
						entry.target.classList.add('is-visible');
						observer.unobserve(entry.target);
					}
				});
			}, options);

			const elements = document.querySelectorAll('.fade-on-scroll');
			elements.forEach(el => observer.observe(el));
		},
		async fetchMateriais() {
			this.loading = true;
			try {
				const response = await MaterialService.dashboard();
				const dashboardData = response.data || response;
				
				this.categoriasMock.forEach(cat => {
					if (dashboardData[cat.nome]) {
						cat.livros = dashboardData[cat.nome];
					}
				});
			} catch (err) {
				console.error("Erro ao carregar dados da Home:", err);
			} finally {
				this.loading = false;
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
/* Typography & Base */
.apple-home-root {
	background-color: transparent;
	min-height: 100vh;
	font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif;
	color: #f5f5f7;
	overflow-x: hidden;
}

.section-container {
	max-width: 1200px;
	margin: 0 auto;
	padding: 0 24px;
}

/* Glass Panel Utility */
.glass-panel {
	background: rgba(255, 255, 255, 0.04);
	backdrop-filter: blur(40px) saturate(200%);
	-webkit-backdrop-filter: blur(40px) saturate(200%);
	border: 1px solid rgba(255, 255, 255, 0.08);
	border-radius: 24px;
	box-shadow: 0 20px 40px rgba(0, 0, 0, 0.2);
}

/* Hero Section */
.apple-hero {
	position: relative;
	min-height: 90vh;
	display: flex;
	align-items: center;
	justify-content: center;
	text-align: center;
	padding: 100px 20px;
	overflow: hidden;
}

.hero-content-center {
	position: relative;
	z-index: 2;
	max-width: 800px;
}

.pill-badge {
	display: inline-block;
	padding: 6px 16px;
	background: rgba(255, 255, 255, 0.1);
	border: 1px solid rgba(255, 255, 255, 0.2);
	border-radius: 100px;
	font-size: 13px;
	font-weight: 600;
	color: #f5f5f7;
	margin-bottom: 24px;
	backdrop-filter: blur(10px);
}

.apple-huge-title {
	font-size: clamp(3rem, 10vw, 6rem);
	font-weight: 800;
	letter-spacing: -0.04em;
	line-height: 1.05;
	margin-bottom: 24px;
	background: linear-gradient(180deg, #ffffff 0%, #a1a1a6 100%);
	-webkit-background-clip: text;
	-webkit-text-fill-color: transparent;
	background-clip: text;
}

.apple-subtitle {
	font-size: clamp(1.2rem, 3vw, 1.8rem);
	font-weight: 500;
	letter-spacing: -0.01em;
	color: #a1a1a6;
	line-height: 1.4;
	margin-bottom: 40px;
}

.hero-actions {
	display: flex;
	align-items: center;
	justify-content: center;
	gap: 20px;
}

.apple-btn-primary {
	background: #f5f5f7;
	color: #1d1d1f;
	padding: 14px 28px;
	border-radius: 100px;
	font-size: 15px;
	font-weight: 600;
	border: none;
	cursor: pointer;
	transition: all 0.3s ease;
}

.apple-btn-primary:hover {
	transform: scale(1.02);
	background: #ffffff;
	box-shadow: 0 0 20px rgba(255,255,255,0.2);
}

.apple-btn-link {
	background: transparent;
	color: #2997ff;
	font-size: 17px;
	font-weight: 500;
	border: none;
	cursor: pointer;
	display: flex;
	align-items: center;
	transition: opacity 0.3s ease;
}

.apple-btn-link:hover {
	opacity: 0.8;
}

.hero-glow {
	position: absolute;
	top: 50%;
	left: 50%;
	width: 60vw;
	height: 60vw;
	max-width: 800px;
	max-height: 800px;
	background: radial-gradient(circle, rgba(10, 132, 255, 0.15) 0%, transparent 60%);
	transform: translate(-50%, -50%);
	pointer-events: none;
	z-index: 1;
	filter: blur(60px);
}

/* Animations */
.fade-up-delay-1 { animation: appleFadeUp 1s cubic-bezier(0.16, 1, 0.3, 1) 0.1s forwards; opacity: 0; }
.fade-up-delay-2 { animation: appleFadeUp 1s cubic-bezier(0.16, 1, 0.3, 1) 0.2s forwards; opacity: 0; }
.fade-up-delay-3 { animation: appleFadeUp 1s cubic-bezier(0.16, 1, 0.3, 1) 0.3s forwards; opacity: 0; }
.fade-up-delay-4 { animation: appleFadeUp 1s cubic-bezier(0.16, 1, 0.3, 1) 0.4s forwards; opacity: 0; }

@keyframes appleFadeUp {
	from { opacity: 0; transform: translateY(30px); }
	to { opacity: 1; transform: translateY(0); }
}

.fade-on-scroll {
	opacity: 0;
	transform: translateY(40px);
	transition: opacity 0.8s cubic-bezier(0.16, 1, 0.3, 1), transform 0.8s cubic-bezier(0.16, 1, 0.3, 1);
	will-change: opacity, transform;
}

.fade-on-scroll.is-visible {
	opacity: 1;
	transform: translateY(0);
}

/* Sections */
.apple-categories-section, .apple-features-section {
	padding: 60px 0;
	position: relative;
	z-index: 5;
}

.apple-section-header {
	margin-bottom: 60px;
}

.center-text {
	text-align: center;
}

.apple-section-title {
	font-size: clamp(2.5rem, 5vw, 4rem);
	font-weight: 700;
	letter-spacing: -0.03em;
	color: #f5f5f7;
	line-height: 1.1;
	margin-bottom: 12px;
}

.apple-section-desc {
	font-size: clamp(1.2rem, 2vw, 1.5rem);
	color: #a1a1a6;
	font-weight: 500;
}

/* Shelves Grid */
.apple-shelves-grid {
	display: grid;
	grid-template-columns: repeat(auto-fit, minmax(320px, 1fr));
	gap: 24px;
}

.apple-shelf-card {
	padding: 24px;
	transition: transform 0.4s ease, border-color 0.4s ease;
}

.apple-shelf-card:hover {
	transform: scale(1.02);
	border-color: rgba(255,255,255,0.2);
}

.shelf-header {
	display: flex;
	align-items: center;
	gap: 12px;
	margin-bottom: 24px;
	cursor: pointer;
}

.shelf-icon-box {
	background: rgba(255,255,255,0.1);
	width: 40px;
	height: 40px;
	border-radius: 10px;
	display: flex;
	align-items: center;
	justify-content: center;
}

.shelf-title {
	font-size: 17px;
	font-weight: 600;
	letter-spacing: -0.01em;
}

.shelf-books {
	display: flex;
	flex-direction: column;
	gap: 12px;
}

.apple-book-row {
	display: flex;
	align-items: center;
	gap: 16px;
	padding: 12px;
	border-radius: 12px;
	background: rgba(255,255,255,0.02);
	transition: background 0.2s ease;
	cursor: pointer;
}

.apple-book-row:hover {
	background: rgba(255,255,255,0.08);
}

.book-cover-mini {
	width: 48px;
	height: 64px;
	background: #000;
	position: relative;
	overflow: hidden;
	border-radius: 8px;
	display: flex;
	align-items: center;
	justify-content: center;
	box-shadow: 0 4px 10px rgba(0,0,0,0.3);
}

.book-cover-mini::before {
	content: '';
	position: absolute;
	top: -50%; left: -50%; width: 200%; height: 200%;
	background: radial-gradient(circle at 30% 30%, #007AFF 0%, transparent 60%), radial-gradient(circle at 70% 70%, #5AC8FA 0%, transparent 60%);
	filter: blur(15px);
	opacity: 0.8;
}

.book-cover-mini::after {
	content: '';
	position: absolute;
	inset: 2px;
	background: rgba(255, 255, 255, 0.05);
	backdrop-filter: blur(10px) saturate(180%);
	border: 1px solid rgba(255, 255, 255, 0.1);
	border-radius: 6px;
}

.book-cover-mini .v-icon {
	position: relative;
	z-index: 2;
}

.book-info {
	flex: 1;
	overflow: hidden;
}

.book-title {
	font-size: 15px;
	font-weight: 600;
	white-space: nowrap;
	overflow: hidden;
	text-overflow: ellipsis;
	margin-bottom: 4px;
}

.book-author {
	font-size: 13px;
	color: #a1a1a6;
}

.empty-shelf {
	font-size: 14px;
	color: #a1a1a6;
	text-align: center;
	padding: 20px 0;
}

.apple-skeleton-book {
	height: 88px;
	border-radius: 12px;
	background: rgba(255,255,255,0.05);
	animation: pulse 1.5s infinite;
}

@keyframes pulse {
	0% { opacity: 0.6; }
	50% { opacity: 1; }
	100% { opacity: 0.6; }
}

/* Features Grid */
.apple-features-grid {
	display: grid;
	grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
	gap: 24px;
	margin-top: 60px;
}

.apple-feature-box {
	padding: 32px;
	text-align: left;
	display: flex;
	flex-direction: column;
	gap: 16px;
}

.feature-icon {
	margin-bottom: 8px;
}

.feature-title {
	font-size: 20px;
	font-weight: 600;
	letter-spacing: -0.01em;
}

.feature-desc {
	font-size: 15px;
	color: #a1a1a6;
	line-height: 1.5;
}
</style>
