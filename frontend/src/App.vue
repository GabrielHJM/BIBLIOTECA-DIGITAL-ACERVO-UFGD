<template>
	<v-app class="ios-app">
		<!-- Conditionally render App Bar -->
		<v-app-bar app v-if="showBar" :elevation="0" class="glass-app-bar" height="80">
			<!-- Logo Section -->
			<div class="header-logo-container d-flex align-center ml-2 ml-sm-8 logo-clickable" @click="$router.push('/')">
				<img src="@/assets/images/site-images/login/img-logo-menu-bar.png" alt="Logo" class="logo-img-original" />
				<div class="logo-text-stack ml-3 hidden-sm-and-down">
					<h1 class="original-title">BIBLIOTECA</h1>
					<div class="subtitle-accent">DIGITAL</div>
				</div>
			</div>

			<v-spacer></v-spacer>

			<!-- Central Search Bar -->
			<div class="search-wrapper-global mx-4">
				<v-text-field
					v-model="searchInput"
					ref="globalSearch"
					placeholder="Buscar conhecimento... (Ctrl+K)"
					variant="solo"
					rounded="pill"
					flat
					density="compact"
					hide-details
					prepend-inner-icon="mdi-magnify"
					:loading="loading"
					class="ios-search-field"
					:class="{ 'search-active': isSearchFocused }"
					@focus="isSearchFocused = true"
					@blur="isSearchFocused = false"
					@keyup.enter="doSearch"
				>
					<template v-slot:append-inner>
						<kbd class="search-kbd hidden-sm-and-down">⌘ K</kbd>
					</template>
				</v-text-field>
			</div>

			<v-spacer></v-spacer>

			<!-- User Actions -->
			<div class="nav-actions-wrapper d-flex align-center mr-2 mr-sm-8" style="gap: 8px;">
				<v-btn icon variant="text" class="theme-toggle-btn" @click="toggleTheme">
					<v-icon size="22">{{ isDarkTheme ? 'mdi-sun-side' : 'mdi-moon-waning-crescent' }}</v-icon>
				</v-btn>

				<!-- Dropdown de Usuário -->
				<v-menu
					v-model="userMenuOpen"
					:close-on-content-click="true"
					location="bottom end"
					offset="12"
					transition="slide-y-transition"
				>
					<template v-slot:activator="{ props }">
						<div class="user-trigger d-flex align-center ml-2" v-bind="props" style="cursor:pointer; gap: 10px;">
							<div class="text-right hidden-sm-and-down">
								<div class="welcome-text">{{ isLoggedIn ? 'Bem-vindo,' : 'Olá,' }}</div>
								<div class="username-text">{{ isLoggedIn ? userDisplayName : 'Visitante' }}</div>
							</div>
							<v-avatar size="44" class="header-avatar-glass">
								<v-img v-if="isLoggedIn && userAvatar" :src="userAvatar" cover></v-img>
								<v-icon v-else color="white" size="26">{{ isLoggedIn ? 'mdi-account-circle' : 'mdi-account-circle-outline' }}</v-icon>
							</v-avatar>
						</div>
					</template>

					<!-- Dropdown Menu -->
					<v-card class="user-dropdown-card" min-width="220" elevation="16">
						<!-- Header do dropdown -->
						<div class="dropdown-header pa-4 pb-3">
							<div class="dropdown-avatar-row d-flex align-center mb-2" style="gap: 10px;">
								<v-avatar size="38" color="rgba(0,184,212,0.15)">
									<v-img v-if="isLoggedIn && userAvatar" :src="userAvatar" cover></v-img>
									<v-icon v-else color="#00B8D4" size="22">{{ isLoggedIn ? 'mdi-account' : 'mdi-account-outline' }}</v-icon>
								</v-avatar>
								<div>
									<div class="dropdown-user-name">{{ isLoggedIn ? userDisplayName : 'Visitante' }}</div>
									<div class="dropdown-user-role">{{ isLoggedIn ? 'Usuário autenticado' : 'Não autenticado' }}</div>
								</div>
							</div>
						</div>

						<v-divider class="opacity-10"></v-divider>

						<!-- Opções para visitante (não logado) -->
						<template v-if="!isLoggedIn">
							<v-list class="dropdown-list py-2">
								<v-list-item
									prepend-icon="mdi-login"
									title="Entrar"
									subtitle="Acesse sua conta"
									class="dropdown-item"
									@click="$router.push('/login')"
								>
									<template v-slot:append>
										<v-icon size="16" color="rgba(255,255,255,0.3)">mdi-chevron-right</v-icon>
									</template>
								</v-list-item>
								<v-list-item
									prepend-icon="mdi-account-plus"
									title="Criar Conta"
									subtitle="Cadastre-se gratuitamente"
									class="dropdown-item dropdown-item--highlight"
									@click="$router.push('/cadastro')"
								>
									<template v-slot:append>
										<v-icon size="16" color="rgba(255,255,255,0.3)">mdi-chevron-right</v-icon>
									</template>
								</v-list-item>
								<v-divider class="opacity-10 my-1"></v-divider>
								<v-list-item
									prepend-icon="mdi-help-circle-outline"
									title="Sobre Nós"
									class="dropdown-item"
									@click="$router.push('/sobre-nos')"
								>
									<template v-slot:append>
										<v-icon size="16" color="rgba(255,255,255,0.3)">mdi-chevron-right</v-icon>
									</template>
								</v-list-item>
							</v-list>
						</template>

						<!-- Opções para usuário logado -->
						<template v-else>
							<v-list class="dropdown-list py-2">
								<v-list-item
									prepend-icon="mdi-account-edit-outline"
									title="Editar Perfil"
									subtitle="Mude seu nome e foto"
									class="dropdown-item"
									@click="$router.push('/perfil')"
								>
									<template v-slot:append>
										<v-icon size="16" color="rgba(255,255,255,0.3)">mdi-chevron-right</v-icon>
									</template>
								</v-list-item>
								<v-list-item
									prepend-icon="mdi-view-dashboard"
									title="Meu Painel"
									subtitle="Acesse seu dashboard"
									class="dropdown-item"
									@click="$router.push('/dashboard')"
								>
									<template v-slot:append>
										<v-icon size="16" color="rgba(255,255,255,0.3)">mdi-chevron-right</v-icon>
									</template>
								</v-list-item>
								<v-list-item
									prepend-icon="mdi-heart-outline"
									title="Favoritos"
									class="dropdown-item"
									@click="$router.push('/dashboard')"
								>
									<template v-slot:append>
										<v-icon size="16" color="rgba(255,255,255,0.3)">mdi-chevron-right</v-icon>
									</template>
								</v-list-item>
								<v-divider class="opacity-10 my-1"></v-divider>
								<v-list-item
									prepend-icon="mdi-logout"
									title="Sair"
									class="dropdown-item dropdown-item--danger"
									@click="logout"
								></v-list-item>
							</v-list>
						</template>
					</v-card>
				</v-menu>
			</div>
		</v-app-bar>

		<v-main>
			<v-container fluid class="pa-0" :class="{ 'public-page-container': !showBar }">
				<router-view v-slot="{ Component }">
						<transition name="ios-page" mode="out-in">
							<component :is="Component" :key="$route.fullPath" />
						</transition>
				</router-view>
			</v-container>
		</v-main>

		<Footer />

		<!-- Global Accessibility Panel -->
		<AccessibilityPanel />
	</v-app>
</template>

<script>
import Footer from './components/Footer.vue'
import AccessibilityPanel from './components/AccessibilityPanel.vue'
import { ref, computed } from 'vue'
import auth, { state as authState } from '@/auth'
import { useAccessibility } from '@/composables/useAccessibility'
import { useTheme } from 'vuetify'

export default {
	name: 'App',
	components: { Footer, AccessibilityPanel },
	data() {
		return {
			publicRoutes: ['/login', '/cadastro', '/esqueci-senha'],
			searchInput: '',
			isSearchFocused: false,
			userMenuOpen: false,
		}
	},
	methods: {
		logout() {
			auth.logout()
			this.$router.push('/')
		},
		async doSearch() {
			const q = this.searchInput.trim()
			if (!q) return

			this.loading = true
			try {
				// Pequeno delay para feedback visual se for muito rápido
				await new Promise(resolve => setTimeout(resolve, 500))

				if (auth.isAuthenticated()) {
					this.$router.push({ path: '/dashboard', query: { q } })
				} else {
					this.$router.push({ path: '/explorar', query: { q } })
				}
			} finally {
				this.loading = false
			}
		},
		handleGlobalKeydown(e) {
			if ((e.ctrlKey || e.metaKey) && e.key === 'k') {
				e.preventDefault();
				this.$refs.globalSearch?.focus();
			}
		}
	},
	computed: {
		showBar() {
			return !this.publicRoutes.includes(this.$route.path)
		},
		isLoggedIn() {
			return authState.isAuthenticated
		},
		userDisplayName() {
			const user = authState.user
			return user?.nome || user?.email || 'Usuário'
		},
		userAvatar() {
			const user = authState.user
			return user?.foto_url || null
		}
	},
	watch: {
		'$route.path'() {
			this.showBar
			// Limpa busca ao navegar
			this.searchInput = ''
		},
		'$route.query.q'(val) {
			if (val) this.searchInput = val
		}
	},
	setup() {
		const loading = ref(false)
		const { init } = useAccessibility()
		const theme = useTheme()

		const toggleTheme = () => {
			theme.global.name.value = theme.global.current.value.dark ? 'light' : 'dark'
		}

		const isDarkTheme = computed(() => theme.global.current.value.dark)

		return {
			loading,
			init,
			cleanupAccessibility: useAccessibility().cleanup,
			theme,
			toggleTheme,
			isDarkTheme
		}
	},
	mounted() {
		this.init()
		// Restaura query de busca se existir na URL
		if (this.$route.query.q) {
			this.searchInput = this.$route.query.q
		}

		// Keyboard shortcut Ctrl+K
		window.addEventListener('keydown', this.handleGlobalKeydown);

		// Router hooks for loading
		this.$router.beforeEach((to, from, next) => {
			this.loading = true;
			next();
		});
		this.$router.afterEach(() => {
			setTimeout(() => { this.loading = false; }, 300);
		});
	},
	beforeUnmount() {
		window.removeEventListener('keydown', this.handleGlobalKeydown);
		if (this.cleanupAccessibility) {
			this.cleanupAccessibility()
		}
	}
}
</script>

<style>
	:root {
		--ios-blue: #007AFF;
		--ios-cyan: #00B8D4;
		--ios-bg: #F2F2F7;
		--ios-bg-dark: #000000;
		--ios-card: rgba(255, 255, 255, 0.7);
		--ios-card-dark: rgba(28, 28, 30, 0.7);
		--spring-easing: cubic-bezier(0.16, 1, 0.3, 1);
		--apple-font: -apple-system, BlinkMacSystemFont, "SF Pro Text", "SF Pro Display", "Helvetica Neue", Arial, sans-serif;
	}

	@import url('https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700;800;900&display=swap');

	body {
		margin: 0;
		padding: 0;
		font-family: 'Inter', var(--apple-font);
		overflow-x: hidden;
		-webkit-font-smoothing: antialiased;
		background-color: rgb(var(--v-theme-background));
	}

	.ios-app {
		background: transparent !important;
	}

	/* Glassmorphism AppBar */
	.glass-app-bar {
		background: rgba(var(--v-theme-surface), 0.7) !important;
		backdrop-filter: blur(20px) saturate(180%) !important;
		-webkit-backdrop-filter: blur(20px) saturate(180%) !important;
		border-bottom: 0.5px solid rgba(var(--v-border-color), 0.1) !important;
		transition: all 0.4s var(--spring-easing) !important;
		z-index: 1000 !important;
	}

	.header-logo-container {
		transition: transform 0.3s var(--spring-easing);
	}

	.header-logo-container:hover {
		transform: scale(1.02);
	}

	.logo-img-original {
		height: 36px;
		filter: drop-shadow(0 4px 12px rgba(0,0,0,0.1));
	}

	.logo-text-stack {
		display: flex;
		flex-direction: column;
		line-height: 1;
	}

	.original-title {
		font-size: 16px !important;
		font-weight: 800 !important;
		letter-spacing: -0.5px;
		margin: 0;
	}

	.subtitle-accent {
		font-size: 10px;
		font-weight: 900;
		color: var(--ios-cyan);
		letter-spacing: 2px;
	}

	/* Enhanced Search Bar */
	.search-wrapper-global {
		flex: 1 1 auto;
		max-width: 600px;
		width: 100%;
	}

	.ios-search-field :deep(.v-field) {
		background: rgba(var(--v-theme-on-surface), 0.05) !important;
		border-radius: 14px !important;
		height: 42px !important;
		transition: all 0.3s var(--spring-easing);
		border: 1px solid transparent !important;
	}

	.search-active :deep(.v-field) {
		background: rgba(var(--v-theme-surface), 0.9) !important;
		border-color: var(--ios-cyan) !important;
		box-shadow: 0 8px 24px rgba(0,0,0,0.12) !important;
		transform: scale(1.02);
	}

	.search-kbd {
		background: rgba(var(--v-theme-on-surface), 0.1);
		padding: 2px 6px;
		border-radius: 6px;
		font-size: 10px;
		font-weight: 600;
		color: rgba(var(--v-theme-on-surface), 0.5);
		margin-right: 8px;
	}

	.theme-toggle-btn {
		opacity: 0.7;
		transition: opacity 0.2s;
	}

	.theme-toggle-btn:hover {
		opacity: 1;
	}

	.header-avatar-glass {
		border: 1.5px solid rgba(var(--v-theme-on-surface), 0.1);
		transition: transform 0.2s;
	}

	.header-avatar-glass:hover {
		transform: scale(1.05);
	}

	.welcome-text {
		font-size: 10px;
		opacity: 0.6;
		font-weight: 500;
		text-transform: uppercase;
		letter-spacing: 0.5px;
	}

	.username-text {
		font-size: 13px;
		font-weight: 700;
	}

	/* Transitions Style iOS 17 */
	.ios-page-enter-active,
	.ios-page-leave-active {
		transition: opacity 0.4s var(--spring-easing),
					transform 0.4s var(--spring-easing);
	}

	.ios-page-enter-from {
		opacity: 0;
		transform: scale(0.98) translateY(10px);
	}

	.ios-page-leave-to {
		opacity: 0;
		transform: scale(1.02) translateY(-10px);
	}

	.public-page-container {
		display: flex !important;
		flex-direction: column !important;
		align-items: center !important;
		justify-content: center !important;
		min-height: 100vh !important;
		padding: 40px 20px !important;
		overflow-y: auto !important;
	}

	.public-page {
		padding-top: 0 !important;
		padding-left: 0 !important;
		min-height: 100vh;
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		overflow-y: auto;
	}

	@media (max-width: 600px) {
		.search-wrapper-global {
			margin: 0 10px;
		}
		.nav-actions-wrapper .text-right {
			display: none;
		}
	}
</style>
