<template>
	<div class="home-wrapper" :class="{ 'light-theme': !isDarkTheme }">
		<!-- Atmospheric background orbs -->
		<div class="bg-orb bg-orb-1"></div>
		<div class="bg-orb bg-orb-2"></div>

		<!-- Hero Section -->
		<section class="hero-section">
			<div class="hero-content">
				<div class="hero-badge">
					<v-icon size="14" color="#00B8D4" class="mr-1">mdi-library</v-icon>
					Biblioteca Digital
				</div>
				<h1 class="hero-title">Explore o Conhecimento<br><span class="hero-title-accent">Sem Fronteiras</span></h1>
				<p class="hero-subtitle">Acesse milhares de livros, artigos e materiais acadêmicos de forma gratuita e instantânea.</p>
				<div class="hero-cta-group">
					<template v-if="!isLoggedIn">
						<button class="btn-primary" @click="$router.push('/cadastro')">
							<v-icon size="18" class="mr-2">mdi-account-plus</v-icon>
							Criar Conta Grátis
						</button>
						<button class="btn-ghost" @click="$router.push('/login')">
							<v-icon size="18" class="mr-2">mdi-login</v-icon>
							Entrar
						</button>
					</template>
					<template v-else>
						<button class="btn-primary" @click="$router.push('/dashboard')">
							<v-icon size="18" class="mr-2">mdi-view-dashboard</v-icon>
							Acessar Meu Dashboard
						</button>
					</template>
				</div>
			</div>
			<!-- Stats row -->
			<div class="stats-row">
				<div class="stat-item" v-for="stat in stats" :key="stat.label">
					<div class="stat-number">{{ stat.value }}</div>
					<div class="stat-label">{{ stat.label }}</div>
				</div>
			</div>
		</section>

		<!-- Categories Section -->
		<section class="categories-section">
			<div class="section-header">
				<h2 class="section-title">COMECE SEUS ESTUDOS</h2>
				<div class="title-underline"></div>

				<div class="acervos-filter mt-8 d-flex justify-center">
					<v-select
						v-model="selectedAcervo"
						:items="acervosList"
						item-title="title"
						item-value="value"
						variant="outlined"
						rounded="pill"
						density="compact"
						hide-details
						class="acervo-select"
						@update:modelValue="fetchMateriais"
					>
						<template v-slot:prepend-inner>
							<v-icon color="#00B8D4" class="mr-2">mdi-bookshelf</v-icon>
						</template>
					</v-select>
				</div>
			</div>

			<div class="categories-grid">
				<div
					v-for="(cat, idx) in categoriasMock"
					:key="cat.nome"
					class="category-card animate-fade-in"
					:style="{ animationDelay: (idx * 0.08) + 's' }"
				>
					<div class="category-card-header" @click="$router.push({ name: 'explorar', query: { categoria: cat.nome } })">
						<div class="category-icon-wrap">
							<v-icon :color="cat.iconColor" size="20">{{ cat.icon }}</v-icon>
						</div>
						<h3 class="category-name">{{ cat.nome }}</h3>
						<button class="ver-tudo-btn">
							Ver Tudo <v-icon size="14">mdi-chevron-right</v-icon>
						</button>
					</div>

					<!-- Books list -->
					<div class="books-list">
						<template v-if="loading">
							<div v-for="s in 2" :key="'sk'+s" class="skeleton-book"></div>
						</template>
						<template v-else-if="cat.livros && cat.livros.length > 0">
							<div
								v-for="livro in cat.livros"
								:key="livro.id"
								class="book-row"
								@click="$router.push('/estudo/' + livro.id)"
							>
								<div class="book-cover">
									<v-img
										:src="livro.capa_url || 'https://images.unsplash.com/photo-1544947950-fa07a98d237f?q=80&w=100'"
										width="48"
										height="68"
										cover
										class="rounded-lg shadow-sm"
									>
										<template v-slot:placeholder>
											<div class="d-flex align-center justify-center fill-height bg-glass-dark">
												<v-icon color="rgba(255,255,255,0.2)">mdi-book-open-variant</v-icon>
											</div>
										</template>
									</v-img>
								</div>
								<div class="book-info">
									<div class="book-title">{{ livro.titulo }}</div>
									<div class="book-author">
										{{ livro.autor }}
									</div>
									<div class="book-meta-row">
										<span class="book-tag">{{ livro.ano_publicacao }}</span>
										<span v-if="livro.media_nota" class="book-rating">
											<v-icon size="10" color="amber" class="mr-1">mdi-star</v-icon>
											{{ livro.media_nota.toFixed(1) }}
										</span>
									</div>
								</div>
								<v-icon size="18" class="book-arrow">mdi-chevron-right</v-icon>
							</div>
						</template>
						<div v-else class="empty-state">
							<v-icon size="24" style="color:rgba(255,255,255,0.2)" class="mb-1">mdi-book-off-outline</v-icon>
							<span>Nenhum livro disponível</span>
						</div>
					</div>
				</div>
			</div>
		</section>

		<!-- Features / Innovation Section -->
		<section class="features-section">
			<div class="features-inner">
				<div class="features-text-col">
					<h2 class="features-headline">Por que a Biblioteca Digital Acessível é inovadora?</h2>
					<p class="features-desc">Um ecossistema de leitura pensado para democratizar o acesso ao conhecimento, sem barreiras e sem custos.</p>
					<button class="btn-primary mt-6" @click="$router.push('/cadastro')">
						<v-icon size="18" class="mr-2">mdi-rocket-launch</v-icon>
						Comece Agora
					</button>
				</div>

				<div class="features-cards-col">
					<div
						v-for="(feature, fIdx) in features"
						:key="feature.title"
						class="feature-pill animate-slide-up"
						:style="{ animationDelay: (fIdx * 0.1) + 's' }"
					>
						<div class="feature-pill-icon" :style="{ background: feature.bg }">
							<v-icon :color="feature.iconColor" size="20">{{ feature.icon }}</v-icon>
						</div>
						<div class="feature-pill-text">
							<div class="feature-pill-title">{{ feature.title }}</div>
							<div class="feature-pill-desc">{{ feature.desc }}</div>
						</div>
					</div>
				</div>
			</div>
		</section>

		<!-- CTA Banner -->
		<section class="cta-section">
			<div class="cta-inner">
				<v-icon size="40" color="white" class="mb-4">mdi-book-open-page-variant</v-icon>
				<h2 class="cta-title">Pronto para explorar?</h2>
				<p class="cta-subtitle">Junte-se a milhares de estudantes e profissionais que já utilizam nossa plataforma de forma gratuita.</p>
				<div class="cta-actions">
					<template v-if="!isLoggedIn">
						<button class="btn-white" @click="$router.push('/cadastro')">Criar Conta Grátis</button>
						<button class="btn-outline-white" @click="$router.push('/login')">Já tenho conta</button>
					</template>
					<template v-else>
						<button class="btn-white" @click="$router.push('/dashboard')">Acessar Meu Dashboard</button>
					</template>
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
	name: 'HomePagePremium',
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
			selectedAcervo: '',
			acervosList: [
				{ title: 'Todos os Acervos', value: '' },
				{ title: 'SciELO', value: 'SciELO' },
				{ title: 'CAPES', value: 'CAPES' },
				{ title: 'IEEE', value: 'IEEE' },
				{ title: 'OpenAlex', value: 'OpenAlex' }
			],
			stats: [
				{ value: '5K+', label: 'Materiais Livres' },
				{ value: '50+', label: 'Categorias' },
				{ value: 'Grátis', label: 'Acesso' },
				{ value: '24/7', label: 'Disponível' }
			],
			categoriasMock: [
				{ nome: 'TECNOLOGIA', livros: [], icon: 'mdi-laptop', iconColor: '#00B8D4' },
				{ nome: 'SAÚDE', livros: [], icon: 'mdi-heart-pulse', iconColor: '#FF6B9D' },
				{ nome: 'MATEMÁTICA', livros: [], icon: 'mdi-calculator-variant', iconColor: '#FFD60A' },
				{ nome: 'CIÊNCIAS', livros: [], icon: 'mdi-flask', iconColor: '#39FF14' },
				{ nome: 'HISTÓRIA', livros: [], icon: 'mdi-castle', iconColor: '#BF5AF2' },
				{ nome: 'CONTABILIDADE', livros: [], icon: 'mdi-currency-usd', iconColor: '#FF9F0A' }
			],
			features: [
				{ title: 'Sem Limites', desc: 'Acesse quantos livros quiser, sem restrições.', icon: 'mdi-infinity', iconColor: '#00F2FE', bg: 'rgba(0,242,254,0.12)' },
				{ title: 'Acesso Instantâneo', desc: 'Leitura online sem necessidade de download.', icon: 'mdi-timer-outline', iconColor: '#39FF14', bg: 'rgba(57,255,20,0.12)' },
				{ title: 'Personalização', desc: 'Favoritos, histórico e recomendações.', icon: 'mdi-palette-outline', iconColor: '#BF5AF2', bg: 'rgba(191,90,242,0.12)' },
				{ title: 'Sem Taxas Ocultas', desc: '100% gratuito, sempre e para todos.', icon: 'mdi-currency-usd-off', iconColor: '#FF6B6B', bg: 'rgba(255,107,107,0.12)' },
				{ title: 'Global', desc: 'Acesse de qualquer lugar do mundo.', icon: 'mdi-earth', iconColor: '#FFD60A', bg: 'rgba(255,214,10,0.12)' }
			]
		}
	},
	async mounted() {
		await this.fetchMateriais();
	},
	methods: {
		async fetchMateriais() {
			this.loading = true;
			try {
				// Simula um pequeno delay para mostrar os esqueletos (UX mais "smooth")
				await new Promise(resolve => setTimeout(resolve, 800));

				const promises = this.categoriasMock.map(async (cat) => {
					// Usa o nome da categoria usando filtro correto 'categoria' em vez de busca livre, exigindo 3 aleatórios locais
					const response = await MaterialService.pesquisar('', cat.nome, this.selectedAcervo, 0, 0, 3, 0, 'random');
					// Com o novo interceptor, response já contém o array de materiais
					cat.livros = Array.isArray(response.data) ? response.data : [];
				});
				await Promise.all(promises);
			} catch (err) {
				console.error("Erro ao carregar dados da Home:", err);
				// Aqui poderíamos emitir um alerta global ou local
			} finally {
				this.loading = false;
			}
		}
	}
}
</script>

<style scoped>
/* ===========================
   BASE
=========================== */
.home-wrapper {
	background-color: #2d5178;
	min-height: 100vh;
	position: relative;
	overflow-x: hidden;
	font-family: -apple-system, BlinkMacSystemFont, "SF Pro Text", "Helvetica Neue", Arial, sans-serif;
	transition: background-color 0.3s ease, color 0.3s ease;
}

/* Light Theme Overrides */
.home-wrapper.light-theme {
	background-color: #f5f7fa;
	color: #1a3a5c;
}

.home-wrapper.light-theme .hero-title,
.home-wrapper.light-theme .hero-subtitle,
.home-wrapper.light-theme .section-title,
.home-wrapper.light-theme .category-name,
.home-wrapper.light-theme .book-title,
.home-wrapper.light-theme .features-headline,
.home-wrapper.light-theme .features-desc,
.home-wrapper.light-theme .feature-pill-title,
.home-wrapper.light-theme .cta-title,
.home-wrapper.light-theme .cta-subtitle {
	color: #1a3a5c !important;
}

.home-wrapper.light-theme .stat-number {
	color: #00B8D4 !important;
}

.home-wrapper.light-theme .stat-label,
.home-wrapper.light-theme .book-author,
.home-wrapper.light-theme .book-rating,
.home-wrapper.light-theme .book-tag,
.home-wrapper.light-theme .feature-pill-desc {
	color: #4a6a8c !important;
}

.home-wrapper.light-theme .stats-row,
.home-wrapper.light-theme .category-card,
.home-wrapper.light-theme .book-row,
.home-wrapper.light-theme .feature-pill {
	background: rgba(255, 255, 255, 0.8) !important;
	border-color: rgba(0, 0, 0, 0.08) !important;
	box-shadow: 0 4px 15px rgba(0,0,0,0.05) !important;
}

.home-wrapper.light-theme .features-section {
	background: rgba(255, 255, 255, 0.6) !important;
	border-color: rgba(0, 0, 0, 0.05) !important;
}

.home-wrapper.light-theme .category-icon-wrap,
.home-wrapper.light-theme .book-tag {
	background: rgba(0, 184, 212, 0.1) !important;
}

.home-wrapper.light-theme .book-arrow,
.home-wrapper.light-theme .empty-state i {
	color: rgba(0, 0, 0, 0.3) !important;
}

.home-wrapper.light-theme .cta-inner {
	background: rgba(255, 255, 255, 0.9) !important;
	border-color: rgba(0, 184, 212, 0.2) !important;
	box-shadow: 0 10px 40px rgba(0,184,212,0.15) !important;
}

.home-wrapper.light-theme .btn-ghost {
	color: #1a3a5c !important;
	border-color: rgba(0, 0, 0, 0.1) !important;
	background: rgba(0, 0, 0, 0.02) !important;
}

.home-wrapper.light-theme .btn-outline-white {
	color: #1a3a5c !important;
	border-color: #00B8D4 !important;
}

.home-wrapper.light-theme .cta-inner .v-icon.mb-4 {
	color: #00B8D4 !important;
}

/* Acervos Select styles */
.acervos-filter {
	width: 100%;
}
.acervo-select {
	max-width: 320px;
}
.acervo-select :deep(.v-field) {
	background: rgba(255, 255, 255, 0.1) !important;
	border-color: rgba(255, 255, 255, 0.2) !important;
	color: white;
	border-radius: 20px !important;
}
.acervo-select :deep(.v-field__input) {
	color: white;
	font-weight: 600;
}
.light-theme .acervo-select :deep(.v-field) {
	background: rgba(255, 255, 255, 0.9) !important;
	border-color: rgba(0, 184, 212, 0.3) !important;
	box-shadow: 0 4px 12px rgba(0,0,0,0.05) !important;
}
.light-theme .acervo-select :deep(.v-field__input) {
	color: #1a3a5c;
}
.light-theme .acervo-select :deep(.v-icon) {
	color: #00B8D4 !important;
}

/* Background orbs */
.bg-orb {
	position: absolute;
	border-radius: 50%;
	filter: blur(80px);
	pointer-events: none;
	z-index: 0;
}
.bg-orb-1 {
	width: 600px; height: 600px;
	top: -100px; right: -100px;
	background: radial-gradient(circle, rgba(0,184,212,0.12) 0%, transparent 70%);
}
.bg-orb-2 {
	width: 500px; height: 500px;
	bottom: 20%; left: -120px;
	background: radial-gradient(circle, rgba(0,122,255,0.10) 0%, transparent 70%);
}

/* ===========================
   BUTTONS (shared)
=========================== */
.btn-primary {
	display: inline-flex;
	align-items: center;
	background: linear-gradient(135deg, #00B8D4 0%, #0097A7 100%);
	color: white;
	border: none;
	border-radius: 14px;
	padding: 14px 28px;
	font-size: 15px;
	font-weight: 700;
	cursor: pointer;
	transition: transform 0.2s ease, box-shadow 0.2s ease;
	box-shadow: 0 6px 20px rgba(0,184,212,0.4);
}
.btn-primary:hover {
	transform: translateY(-2px);
	box-shadow: 0 10px 30px rgba(0,184,212,0.5);
}
.btn-primary:active { transform: scale(0.97); }

.btn-ghost {
	display: inline-flex;
	align-items: center;
	background: rgba(255,255,255,0.1);
	color: white;
	border: 1px solid rgba(255,255,255,0.2);
	border-radius: 14px;
	padding: 14px 28px;
	font-size: 15px;
	font-weight: 600;
	cursor: pointer;
	backdrop-filter: blur(10px);
	transition: all 0.2s ease;
}
.btn-ghost:hover {
	background: rgba(255,255,255,0.18);
	border-color: rgba(255,255,255,0.35);
}

.btn-white {
	display: inline-flex;
	align-items: center;
	background: white;
	color: #1a3a5c;
	border: none;
	border-radius: 14px;
	padding: 14px 30px;
	font-size: 15px;
	font-weight: 700;
	cursor: pointer;
	transition: all 0.2s ease;
	box-shadow: 0 6px 20px rgba(0,0,0,0.15);
}
.btn-white:hover { transform: translateY(-2px); }

.btn-outline-white {
	display: inline-flex;
	align-items: center;
	background: transparent;
	color: white;
	border: 2px solid rgba(255,255,255,0.5);
	border-radius: 14px;
	padding: 13px 30px;
	font-size: 15px;
	font-weight: 600;
	cursor: pointer;
	transition: all 0.2s ease;
}
.btn-outline-white:hover {
	border-color: white;
	background: rgba(255,255,255,0.1);
}

/* ===========================
   HERO SECTION
=========================== */
.hero-section {
	position: relative;
	z-index: 1;
	display: flex;
	flex-direction: column;
	align-items: center;
	text-align: center;
	padding: 80px 24px 60px;
}

.hero-badge {
	display: inline-flex;
	align-items: center;
	background: rgba(0,184,212,0.15);
	border: 1px solid rgba(0,184,212,0.3);
	border-radius: 20px;
	padding: 6px 16px;
	font-size: 12px;
	font-weight: 600;
	color: #00D4E8;
	letter-spacing: 0.5px;
	text-transform: uppercase;
	margin-bottom: 24px;
}

.hero-title {
	font-size: clamp(2rem, 5vw, 3.2rem);
	font-weight: 900;
	color: white;
	line-height: 1.15;
	letter-spacing: -0.5px;
	margin-bottom: 20px;
}
.hero-title-accent {
	background: linear-gradient(135deg, #00B8D4, #00E5FF);
	-webkit-background-clip: text;
	-webkit-text-fill-color: transparent;
	background-clip: text;
}

.hero-subtitle {
	font-size: 1.1rem;
	color: rgba(255,255,255,0.72);
	max-width: 520px;
	line-height: 1.6;
	margin-bottom: 36px;
}

.hero-cta-group {
	display: flex;
	gap: 14px;
	flex-wrap: wrap;
	justify-content: center;
}

/* Stats */
.stats-row {
	display: flex;
	gap: 0;
	margin-top: 56px;
	background: rgba(255,255,255,0.07);
	border: 1px solid rgba(255,255,255,0.12);
	border-radius: 20px;
	backdrop-filter: blur(20px);
	overflow: hidden;
}
.stat-item {
	flex: 1;
	padding: 20px 30px;
	text-align: center;
	border-right: 1px solid rgba(255,255,255,0.1);
}
.stat-item:last-child { border-right: none; }
.stat-number {
	font-size: 1.6rem;
	font-weight: 900;
	color: white;
	line-height: 1;
}
.stat-label {
	font-size: 0.75rem;
	color: rgba(255,255,255,0.55);
	margin-top: 4px;
	text-transform: uppercase;
	letter-spacing: 0.5px;
}

/* ===========================
   CATEGORIES SECTION
=========================== */
.categories-section {
	position: relative;
	z-index: 1;
	padding: 20px 32px 60px;
	max-width: 1280px;
	margin: 0 auto;
}

.section-header {
	text-align: center;
	margin-bottom: 36px;
}
.section-title {
	font-size: 1.8rem;
	font-weight: 900;
	color: white;
	letter-spacing: 2px;
}
.title-underline {
	width: 60px;
	height: 3px;
	background: linear-gradient(90deg, #007AFF, #00D4E8);
	margin: 10px auto 0;
	border-radius: 2px;
}

.categories-grid {
	display: grid;
	grid-template-columns: repeat(auto-fill, minmax(340px, 1fr));
	gap: 20px;
}

/* Category Cards – white card style matching the mockup */
.category-card {
	background: rgba(255,255,255,0.08);
	border: 1px solid rgba(255,255,255,0.12);
	border-radius: 22px;
	padding: 22px;
	backdrop-filter: blur(20px);
	transition: transform 0.3s ease, box-shadow 0.3s ease, border-color 0.3s ease;
	opacity: 0;
}
.category-card:hover {
	transform: translateY(-6px);
	box-shadow: 0 16px 40px rgba(0,0,0,0.3);
	border-color: rgba(0,184,212,0.3);
}

.category-card-header {
	display: flex;
	align-items: center;
	margin-bottom: 16px;
	gap: 10px;
	cursor: pointer;
}
.category-icon-wrap {
	width: 36px; height: 36px;
	background: rgba(255,255,255,0.1);
	border-radius: 10px;
	display: flex;
	align-items: center;
	justify-content: center;
	flex-shrink: 0;
}
.category-name {
	font-size: 0.9rem;
	font-weight: 800;
	color: white;
	letter-spacing: 1px;
	flex: 1;
}
.ver-tudo-btn {
	display: inline-flex;
	align-items: center;
	background: none;
	border: none;
	color: #00D4E8;
	font-size: 12px;
	font-weight: 600;
	cursor: pointer;
	padding: 0;
	transition: opacity 0.2s;
}
.ver-tudo-btn:hover { opacity: 0.7; }

/* Book rows */
.books-list { display: flex; flex-direction: column; gap: 10px; }

.book-row {
	display: flex;
	align-items: center;
	gap: 16px;
	background: rgba(255,255,255,0.03);
	border: 1px solid rgba(255,255,255,0.04);
	border-radius: 16px;
	padding: 12px 16px;
	cursor: pointer;
	transition: all 0.4s cubic-bezier(0.16, 1, 0.3, 1);
}
.book-row:hover {
	background: rgba(255,255,255,0.08);
	border-color: rgba(0,184,212,0.3);
	transform: translateX(6px);
}
.book-cover {
	border-radius: 10px;
	overflow: hidden;
	box-shadow: 0 10px 20px rgba(0,0,0,0.4);
	flex-shrink: 0;
	transition: transform 0.3s ease;
}
.book-row:hover .book-cover {
	transform: scale(1.05) rotate(1deg);
}
.bg-glass-dark {
	background: rgba(0, 0, 0, 0.3);
	backdrop-filter: blur(5px);
}
.book-info { flex: 1; overflow: hidden; }
.book-title {
	font-size: 0.9rem;
	font-weight: 700;
	color: white;
	white-space: nowrap;
	overflow: hidden;
	text-overflow: ellipsis;
}
.book-author {
	font-size: 0.72rem;
	color: rgba(255,255,255,0.55);
	display: flex;
	align-items: center;
	margin-top: 3px;
}
.book-meta-row {
	display: flex;
	align-items: center;
	gap: 8px;
	margin-top: 4px;
}
.book-tag {
	font-size: 0.65rem;
	background: rgba(255,255,255,0.08);
	color: rgba(255,255,255,0.6);
	border-radius: 6px;
	padding: 2px 7px;
}
.book-rating {
	font-size: 0.65rem;
	color: rgba(255,255,255,0.65);
	display: flex;
	align-items: center;
}
.book-arrow { color: rgba(255,255,255,0.25); flex-shrink: 0; }

/* Skeleton */
.skeleton-book {
	height: 82px;
	background: rgba(255,255,255,0.05);
	border-radius: 14px;
	animation: skeleton-pulse 1.5s ease-in-out infinite;
}
@keyframes skeleton-pulse {
	0%, 100% { opacity: 0.5; }
	50% { opacity: 1; }
}

.empty-state {
	display: flex;
	flex-direction: column;
	align-items: center;
	padding: 24px;
	color: rgba(255,255,255,0.35);
	font-size: 0.8rem;
}

/* ===========================
   FEATURES SECTION
=========================== */
.features-section {
	position: relative;
	z-index: 1;
	padding: 60px 32px;
	background: rgba(0,0,0,0.15);
	border-top: 1px solid rgba(255,255,255,0.06);
	border-bottom: 1px solid rgba(255,255,255,0.06);
}
.features-inner {
	max-width: 1100px;
	margin: 0 auto;
	display: grid;
	grid-template-columns: 1fr 1fr;
	gap: 60px;
	align-items: center;
}
.features-headline {
	font-size: 1.9rem;
	font-weight: 900;
	color: white;
	line-height: 1.25;
	margin-bottom: 14px;
}
.features-desc {
	font-size: 1rem;
	color: rgba(255,255,255,0.65);
	line-height: 1.6;
}

.features-cards-col { display: flex; flex-direction: column; gap: 12px; }

.feature-pill {
	display: flex;
	align-items: center;
	gap: 14px;
	background: rgba(255,255,255,0.06);
	border: 1px solid rgba(255,255,255,0.1);
	border-radius: 16px;
	padding: 14px 18px;
	backdrop-filter: blur(10px);
	transition: all 0.25s ease;
	opacity: 0;
}
.feature-pill:hover {
	background: rgba(255,255,255,0.1);
	border-color: rgba(255,255,255,0.2);
	transform: translateX(-4px);
}
.feature-pill-icon {
	width: 42px; height: 42px;
	border-radius: 12px;
	display: flex;
	align-items: center;
	justify-content: center;
	flex-shrink: 0;
}
.feature-pill-title {
	font-size: 0.95rem;
	font-weight: 700;
	color: white;
}
.feature-pill-desc {
	font-size: 0.78rem;
	color: rgba(255,255,255,0.55);
	margin-top: 2px;
}

/* ===========================
   CTA SECTION
=========================== */
.cta-section {
	position: relative;
	z-index: 1;
	padding: 80px 24px;
	text-align: center;
}
.cta-inner {
	display: flex;
	flex-direction: column;
	align-items: center;
	background: rgba(0,184,212,0.1);
	border: 1px solid rgba(0,184,212,0.25);
	border-radius: 32px;
	padding: 60px 40px;
	max-width: 640px;
	margin: 0 auto;
	backdrop-filter: blur(20px);
}
.cta-title {
	font-size: 2rem;
	font-weight: 900;
	color: white;
	margin-bottom: 10px;
}
.cta-subtitle {
	color: rgba(255,255,255,0.65);
	font-size: 1rem;
	margin-bottom: 32px;
}
.cta-btn-group {
	display: flex;
	gap: 14px;
	flex-wrap: wrap;
	justify-content: center;
}

/* ===========================
   ACCESSIBILITY FAB
=========================== */


/* ===========================
   ANIMATIONS
=========================== */
.animate-fade-in {
	animation: fadeInUp 0.7s ease forwards;
}
@keyframes fadeInUp {
	from { opacity: 0; transform: translateY(20px); }
	to { opacity: 1; transform: translateY(0); }
}

.animate-slide-up {
	animation: slideUp 0.6s cubic-bezier(0.23, 1, 0.32, 1) forwards;
}
@keyframes slideUp {
	from { opacity: 0; transform: translateY(16px); }
	to { opacity: 1; transform: translateY(0); }
}

/* ===========================
   RESPONSIVE
=========================== */
@media (max-width: 900px) {
	.features-inner {
		grid-template-columns: 1fr;
		gap: 36px;
	}
	.features-headline { font-size: 1.5rem; }
}

@media (max-width: 700px) {
	.hero-section { padding: 60px 16px 40px; }
	.categories-section { padding: 16px 16px 40px; }
	.categories-grid { grid-template-columns: 1fr; }
	.features-section { padding: 40px 16px; }
	.stats-row { flex-wrap: wrap; }
	.stat-item { border-right: none; border-bottom: 1px solid rgba(255,255,255,0.08); }
	.stat-item:last-child { border-bottom: none; }
	.cta-inner { padding: 40px 20px; }
}
</style>
