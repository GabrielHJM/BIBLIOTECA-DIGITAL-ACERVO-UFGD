<template>
	<div class="home-wrapper" :class="{ 'light-theme': !isDarkTheme }">
		<!-- Atmospheric background orbs -->
		<div class="bg-orb bg-orb-1"></div>
		<div class="bg-orb bg-orb-2"></div>

		<!-- Hero Section -->
		<section class="hero-section">
			<div class="hero-content">
				<div class="hero-badge-pill mb-6">
					<v-icon size="14" color="var(--ios-cyan)" class="mr-2">mdi-sparkles</v-icon>
					O Conhecimento ao seu alcance
				</div>
				<h1 class="hero-title-modern">A Sua Biblioteca<br><span class="accent-gradient">Digital Inteligente</span></h1>
				<p class="hero-description">Acesse milhares de recursos educacionais, livros e artigos científicos em uma plataforma feita para a sua evolução.</p>

				<div class="hero-actions mt-10">
					<template v-if="!isLoggedIn">
						<button class="ios-btn-main mr-4" @click="$router.push('/cadastro')">
							Começar Agora — É Grátis
						</button>
						<button class="ios-btn-secondary" @click="$router.push('/login')">
							Entrar
						</button>
					</template>
					<template v-else>
						<button class="ios-btn-main" @click="$router.push('/dashboard')">
							Ver Meu Dashboard
						</button>
					</template>
				</div>
			</div>

			<!-- Stats Floating Bar -->
			<div class="stats-glass-bar">
				<div class="stat-unit" v-for="stat in stats" :key="stat.label">
					<span class="stat-val">{{ stat.value }}</span>
					<span class="stat-lbl">{{ stat.label }}</span>
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
					class="category-card animate-fade-in premium-hover"
					:style="{
						animationDelay: (idx * 0.1) + 's',
						borderTop: `4px solid ${cat.iconColor}`
					}"
				>
					<div class="category-card-header" @click="$router.push({ name: 'explorar', query: { categoria: cat.nome } })">
						<div class="category-icon-wrap" :style="{ background: `${cat.iconColor}20` }">
							<v-icon :color="cat.iconColor" size="24" class="icon-bounce">{{ cat.icon }}</v-icon>
						</div>
						<h3 class="category-name">{{ cat.nome }}</h3>
						<v-btn icon="mdi-chevron-right" variant="text" size="small" color="cyan" class="chevron-btn"></v-btn>
					</div>

					<!-- Books list -->
					<div class="books-list">
						<template v-if="loading">
							<v-skeleton-loader
								v-for="s in 3"
								:key="'sk'+s"
								type="list-item-avatar-two-line"
								class="skeleton-book-premium"
							></v-skeleton-loader>
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
						:class="{ 'pill-dimmed': hoveredFeature !== null && hoveredFeature !== fIdx }"
						:style="[
							{ animationDelay: (fIdx * 0.1) + 's' },
							pillTransforms[fIdx] || {}
						]"
						@mouseenter="hoveredFeature = fIdx"
						@mousemove="handleMouseMove($event, fIdx)"
						@mouseleave="resetPill(fIdx)"
					>
						<!-- Specular Light Reflection -->
						<div class="pill-specular" :style="pillTransforms[fIdx] ? { left: pillTransforms[fIdx]['--mouse-x'], top: pillTransforms[fIdx]['--mouse-y'] } : {}"></div>
						<div class="feature-pill-icon" :style="{ background: hoveredFeature === fIdx ? 'rgba(255,255,255,0.15)' : feature.bg }">
							<transition name="fade-scale">
								<server-feature-animation
									v-if="hoveredFeature === fIdx"
									:type="feature.animType"
									:icon="feature.icon"
									:style="{ color: feature.iconColor }"
								/>
								<v-icon v-else :color="feature.iconColor" size="20">{{ feature.icon }}</v-icon>
							</transition>
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
import ServerFeatureAnimation from '@/components/animations/ServerFeatureAnimation.vue';

export default {
	name: 'HomeView',
	components: {
		ServerFeatureAnimation
	},
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
			hoveredFeature: null,
			pillTransforms: {},
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
				{ title: 'Sem Limites', desc: 'Acesse quantos livros quiser, sem restrições.', icon: 'mdi-infinity', iconColor: '#00F2FE', bg: 'rgba(0,242,254,0.12)', animType: 'infinity' },
				{ title: 'Acesso Instantâneo', desc: 'Leitura online sem necessidade de download.', icon: 'mdi-timer-outline', iconColor: '#39FF14', bg: 'rgba(57,255,20,0.12)', animType: 'timer' },
				{ title: 'Personalização', desc: 'Favoritos, histórico e recomendações.', icon: 'mdi-palette-outline', iconColor: '#BF5AF2', bg: 'rgba(191,90,242,0.12)', animType: 'personalization' },
				{ title: 'Sem Taxas Ocultas', desc: '100% gratuito, sempre e para todos.', icon: 'mdi-currency-usd-off', iconColor: '#FF6B6B', bg: 'rgba(255,107,107,0.12)', animType: 'security' },
				{ title: 'Global', desc: 'Acesse de qualquer lugar do mundo.', icon: 'mdi-earth', iconColor: '#FFD60A', bg: 'rgba(255,214,10,0.12)', animType: 'global' }
			]
		}
	},
	async mounted() {
		await this.fetchMateriais();
		window.addEventListener('scroll', this.handleParallax);
	},
	beforeUnmount() {
		window.removeEventListener('scroll', this.handleParallax);
	},
	methods: {
		handleParallax() {
			const scrolled = window.scrollY;
			const hero = document.querySelector('.hero-content');
			if (hero) {
				hero.style.transform = `translateY(${scrolled * 0.3}px)`;
				hero.style.opacity = `${1 - scrolled / 500}`;
			}
		},
		async fetchMateriais() {
			this.loading = true;
			try {
				// Simula um pequeno delay para mostrar os esqueletos (UX mais "smooth")
				await new Promise(resolve => setTimeout(resolve, 800));

				const promises = this.categoriasMock.map(async (cat) => {
					try {
						// Usa o nome da categoria usando filtro correto 'categoria' em vez de busca livre, exigindo 3 aleatórios locais
						const response = await MaterialService.pesquisar('', cat.nome, this.selectedAcervo, 0, 0, 3, 0, 'random');

						// Com o novo interceptor, response já contém o array de materiais ou response.data
						let livros = [];
						if (response && Array.isArray(response)) {
							livros = response;
						} else if (response && response.data && Array.isArray(response.data)) {
							livros = response.data;
						}

						cat.livros = livros;
					} catch (err) {
						console.error(`Erro ao carregar livros para a categoria ${cat.nome}:`, err);
						cat.livros = [];
					}
				});
				await Promise.all(promises);
			} catch (err) {
				console.error("Erro ao carregar dados da Home:", err);
				// Aqui poderíamos emitir um alerta global ou local
			} finally {
				this.loading = false;
			}
		},
		handleMouseMove(e, index) {
			const el = e.currentTarget;
			const rect = el.getBoundingClientRect();
			const x = e.clientX - rect.left;
			const y = e.clientY - rect.top;

			// Tilt calculations (-8 to 8 degrees for a subtle premium feel)
			const centerX = rect.width / 2;
			const centerY = rect.height / 2;
			const rotateX = ((y - centerY) / centerY) * -8;
			const rotateY = ((x - centerX) / centerX) * 8;

			// Magnet displacement (up to 4px)
			const moveX = ((x - centerX) / centerX) * 4;
			const moveY = ((y - centerY) / centerY) * 4;

			this.pillTransforms[index] = {
				transform: `perspective(1000px) rotateX(${rotateX}deg) rotateY(${rotateY}deg) translate3d(${moveX}px, ${moveY}px, 0)`,
				'--mouse-x': `${x}px`,
				'--mouse-y': `${y}px`
			};
		},
		resetPill(index) {
			this.hoveredFeature = null;
			this.pillTransforms[index] = {
				transform: 'perspective(1000px) rotateX(0deg) rotateY(0deg) translate3d(0, 0, 0)',
				'--mouse-x': '50%',
				'--mouse-y': '50%'
			};
		}
	}
}
</script>

<style scoped>
/* ===========================
   BASE
=========================== */
.home-wrapper {
	background-color: rgb(var(--v-theme-background));
	min-height: 100vh;
	position: relative;
	overflow-x: hidden;
	transition: background-color 0.4s var(--spring-easing);
}

/* Hero Styles */
.hero-section {
	padding: 120px 24px 100px;
	display: flex;
	flex-direction: column;
	align-items: center;
	text-align: center;
}

.hero-badge-pill {
	display: inline-flex;
	align-items: center;
	padding: 8px 16px;
	background: rgba(var(--v-theme-on-surface), 0.05);
	border-radius: 100px;
	font-size: 13px;
	font-weight: 700;
	color: var(--ios-cyan);
}

.hero-title-modern {
	font-size: clamp(2.5rem, 8vw, 4.5rem);
	font-weight: 900;
	letter-spacing: -2px;
	line-height: 0.95;
	margin-bottom: 24px;
}

.accent-gradient {
	background: linear-gradient(135deg, var(--ios-cyan), #007AFF);
	-webkit-background-clip: text;
	-webkit-text-fill-color: transparent;
}

.hero-description {
	font-size: 1.1rem;
	max-width: 600px;
	opacity: 0.7;
	line-height: 1.6;
	margin: 0 auto;
}

.ios-btn-main {
	background: var(--ios-cyan);
	color: white;
	padding: 16px 32px;
	border-radius: 16px;
	font-weight: 700;
	font-size: 16px;
	box-shadow: 0 10px 30px rgba(0, 184, 212, 0.3);
	transition: all 0.3s var(--spring-easing);
}

.ios-btn-main:hover {
	transform: translateY(-4px);
	box-shadow: 0 15px 40px rgba(0, 184, 212, 0.4);
}

.ios-btn-secondary {
	padding: 16px 32px;
	border-radius: 16px;
	font-weight: 600;
	background: rgba(var(--v-theme-on-surface), 0.05);
	transition: all 0.3s var(--spring-easing);
}

.ios-btn-secondary:hover {
	background: rgba(var(--v-theme-on-surface), 0.1);
}

/* Stats Glass Bar */
.stats-glass-bar {
	margin: 80px auto 0;
	display: flex;
	justify-content: center;
	flex-wrap: wrap;
	max-width: fit-content;
	background: rgba(var(--v-theme-surface), 0.5);
	backdrop-filter: blur(20px);
	border-radius: 24px;
	border: 1px solid rgba(var(--v-border-color), 0.1);
	padding: 8px;
}

.stat-unit {
	padding: 20px 40px;
	display: flex;
	flex-direction: column;
	border-right: 1px solid rgba(var(--v-border-color), 0.1);
}

.stat-unit:last-child {
	border-right: none;
}

.stat-val {
	font-size: 24px;
	font-weight: 900;
	letter-spacing: -1px;
}

.stat-lbl {
	font-size: 12px;
	opacity: 0.5;
	font-weight: 700;
	text-transform: uppercase;
}

.categories-grid {
	display: grid;
	grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
	gap: 24px;
	margin-top: 24px;
}

/* Category Cards – white card style matching the mockup */
.category-card {
	background: rgba(255,255,255,0.05);
	border: 1px solid rgba(255,255,255,0.08);
	border-radius: 28px; /* Mais arredondado */
	padding: 24px;
	backdrop-filter: blur(30px);
	transition: all 0.5s cubic-bezier(0.16, 1, 0.3, 1);
	opacity: 0;
	position: relative;
	overflow: hidden;
}
.category-card::before {
	content: '';
	position: absolute;
	top: 0; left: 0; width: 100%; height: 100%;
	background: linear-gradient(180deg, rgba(255,255,255,0.05) 0%, transparent 100%);
	pointer-events: none;
}
.category-card:hover {
	transform: translateY(-8px) scale(1.02);
	background: rgba(255,255,255,0.12);
	box-shadow: 0 30px 60px rgba(0,0,0,0.4);
	border-color: rgba(255,255,255,0.2);
}

.category-card-header {
	display: flex;
	align-items: center;
	margin-bottom: 20px;
	gap: 12px;
	cursor: pointer;
}
.category-icon-wrap {
	width: 44px; height: 44px;
	border-radius: 12px;
	display: flex;
	align-items: center;
	justify-content: center;
	flex-shrink: 0;
	transition: transform 0.3s ease;
}
.category-card:hover .category-icon-wrap {
	transform: scale(1.1) rotate(5deg);
}
.icon-bounce { transition: transform 0.3s cubic-bezier(0.175, 0.885, 0.32, 1.275); }
.category-card:hover .icon-bounce { transform: scale(1.2); }
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
.skeleton-book-premium {
	background: rgba(255, 255, 255, 0.05) !important;
	border-radius: 16px !important;
	margin-bottom: 8px;
}
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
	transition: all 0.4s cubic-bezier(0.34, 1.56, 0.64, 1);
	opacity: 0;
	cursor: default;
}
.feature-pill:hover {
	background: rgba(255,255,255,0.08);
	border-color: rgba(255,255,255,0.2);
	transform: scale(1.04) translateX(-4px);
	box-shadow: 0 15px 45px rgba(0,0,0,0.25);
	z-index: 10;
	border-color: rgba(255, 255, 255, 0.4);
}
.pill-specular {
	position: absolute;
	width: 150px;
	height: 150px;
	background: radial-gradient(circle, rgba(255,255,255,0.15) 0%, transparent 70%);
	pointer-events: none;
	border-radius: 50%;
	transform: translate(-50%, -50%);
	z-index: 1;
	opacity: 0;
	transition: opacity 0.3s ease;
}
.feature-pill:hover .pill-specular {
	opacity: 1;
}
.pill-dimmed {
	opacity: 0.35 !important;
	filter: blur(1.5px) grayscale(0.4);
	transform: scale(0.975);
	transition: all 0.6s cubic-bezier(0.34, 1, 0.64, 1) !important;
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
   BTN CLASSES FOR CTA
=========================== */
.btn-white {
	background: white;
	color: #00B8D4;
	border: none;
	padding: 16px 32px;
	border-radius: 100px;
	font-weight: 800;
	font-size: 15px;
	box-shadow: 0 10px 20px rgba(0,0,0,0.1);
	transition: all 0.3s var(--spring-easing);
	margin: 0 8px;
}
.btn-white:hover {
	transform: translateY(-4px);
	box-shadow: 0 15px 30px rgba(0,0,0,0.15);
}

.btn-outline-white {
	background: rgba(255,255,255,0.1);
	color: white;
	border: 1px solid rgba(255,255,255,0.3);
	padding: 16px 32px;
	border-radius: 100px;
	font-weight: 700;
	font-size: 15px;
	backdrop-filter: blur(10px);
	transition: all 0.3s var(--spring-easing);
	margin: 0 8px;
}
.btn-outline-white:hover {
	background: rgba(255,255,255,0.2);
	border-color: white;
}

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

/* Transition for icon animation swap */
.fade-scale-enter-active,
.fade-scale-leave-active {
	transition: all 0.4s cubic-bezier(0.34, 1.56, 0.64, 1);
	position: absolute;
}

.fade-scale-enter-from {
	opacity: 0;
	transform: scale(0.4) rotate(-10deg);
}
.fade-scale-leave-to {
	opacity: 0;
	transform: scale(1.4) rotate(10deg);
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
