<template>
	<v-app class="ios-app">
		<v-app-bar v-if="showBar" :elevation="0" class="glass-app-bar" height="90">

			<!-- Logo → clica e vai para home -->
			<template v-slot:prepend>
				<div class="header-logo-container d-flex align-center ml-4 logo-clickable" @click="$router.push('/')">
					<img src="@/assets/images/site-images/login/img-logo-menu-bar.png" alt="Logo" class="logo-img-original" />
					<div class="logo-text-stack ml-2">
						<h1 class="original-title">BIBLIOTECA DIGITAL</h1>
					</div>
				</div>
			</template>

			<v-spacer />

			<!-- Barra de Pesquisa funcional -->
			<div class="search-container-center">
				<v-text-field
					v-model="searchInput"
					placeholder="Pesquise livros, autores, categorias..."
					variant="solo"
					rounded="pill"
					flat
					density="compact"
					hide-details
					prepend-inner-icon="mdi-magnify"
					append-inner-icon="mdi-arrow-right-circle"
					:loading="loading"
					class="original-search-field"
					@focus="isSearchFocused = true"
					@blur="isSearchFocused = false"
					@keyup.enter="doSearch"
					@click:append-inner="doSearch"
				></v-text-field>
			</div>

			<v-spacer />

			<!-- Barra de progresso global -->
			<v-progress-linear
				v-show="loading"
				indeterminate
				color="#00B8D4"
				absolute
				bottom
				height="2"
			></v-progress-linear>

			<!-- Ações do usuário -->
			<div class="nav-actions-original d-flex align-center mr-6" style="gap: 12px;">

				<!-- Dropdown de Usuário -->
				<v-menu
					v-model="userMenuOpen"
					:close-on-content-click="true"
					location="bottom end"
					offset="12"
					transition="slide-y-transition"
				>
					<template v-slot:activator="{ props }">
						<div class="user-trigger d-flex align-center" v-bind="props" style="cursor:pointer; gap: 10px;">
							<div class="text-right hidden-sm-and-down">
								<div class="welcome-text">{{ isLoggedIn ? 'Bem-vindo,' : 'Olá,' }}</div>
								<div class="username-text">{{ isLoggedIn ? userEmail : 'Visitante' }}</div>
							</div>
							<v-avatar size="44" class="header-avatar-glass">
								<v-icon color="white" size="26">{{ isLoggedIn ? 'mdi-account-circle' : 'mdi-account-circle-outline' }}</v-icon>
							</v-avatar>
						</div>
					</template>

					<!-- Dropdown Menu -->
					<v-card class="user-dropdown-card" min-width="220" elevation="16">
						<!-- Header do dropdown -->
						<div class="dropdown-header pa-4 pb-3">
							<div class="dropdown-avatar-row d-flex align-center mb-2" style="gap: 10px;">
								<v-avatar size="38" color="rgba(0,184,212,0.15)">
									<v-icon color="#00B8D4" size="22">{{ isLoggedIn ? 'mdi-account' : 'mdi-account-outline' }}</v-icon>
								</v-avatar>
								<div>
									<div class="dropdown-user-name">{{ isLoggedIn ? userEmail : 'Visitante' }}</div>
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
			<v-container fluid class="pa-0">
				<router-view v-slot="{ Component }">
					<transition name="ios-page" mode="out-in">
						<component :is="Component" />
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
import { ref } from 'vue'
import auth from '@/auth'
import { useAccessibility } from '@/composables/useAccessibility'

export default {
	name: 'App',
	components: { Footer, AccessibilityPanel },
	data() {
		return {
			publicRoutes: ['/login', '/cadastro', '/esqueci-senha'],
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
		}
	},
	computed: {
		showBar() {
			return !this.publicRoutes.includes(this.$route.path)
		},
		isLoggedIn() {
			return auth.isAuthenticated()
		},
		userEmail() {
			const user = auth.getUser()
			return user?.email || 'Usuário'
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
		const searchInput = ref('')
		const loading = ref(false)
		const { init } = useAccessibility()

		return {
			searchInput,
			loading,
			init
		}
	},
	mounted() {
		this.init()
		// Restaura query de busca se existir na URL
		if (this.$route.query.q) {
			this.searchInput = this.$route.query.q
		}
	}
}
</script>

<style>
	:root {
		--ios-blue: #007AFF;
		--ios-cyan: #5AC8FA;
		--ios-bg: #325178; /* Tom de azul do Login */
		--ios-card: rgba(45, 78, 115, 0.85);
		--spring-easing: cubic-bezier(0.4, 0, 0.2, 1);
		--apple-font: -apple-system, BlinkMacSystemFont, "SF Pro Text", "Helvetica Neue", Arial, sans-serif;
	}

	body {
		margin: 0;
		padding: 0;
		font-family: var(--apple-font);
		background-color: var(--ios-bg);
		overflow-x: hidden;
	}

	.ios-app {
		background-color: var(--ios-bg) !important;
	}

	/* Glassmorphism AppBar */
	.glass-app-bar {
		background: rgba(45, 78, 115, 0.6) !important;
		backdrop-filter: blur(25px) saturate(160%) !important;
		-webkit-backdrop-filter: blur(25px) saturate(160%) !important;
		border-bottom: 1px solid rgba(255, 255, 255, 0.1) !important;
	}

	.logo-img-original {
		height: 44px;
		filter: drop-shadow(0 4px 8px rgba(0,0,0,0.3));
		transition: transform 0.3s ease;
	}

	.logo-img-original:hover {
		transform: scale(1.05);
	}

	.original-title {
		font-size: 22px !important;
		font-weight: 900 !important;
		color: white !important;
		letter-spacing: 0.5px;
		text-shadow: 0 2px 10px rgba(0,0,0,0.2);
	}

	.search-container-center {
		flex: 0 1 500px;
		margin: 0 40px;
	}

	.original-search-field :deep(.v-field) {
		background: rgba(255, 255, 255, 0.1) !important;
		backdrop-filter: blur(10px);
		border-radius: 14px !important;
		height: 42px !important;
		border: 1px solid rgba(255, 255, 255, 0.1) !important;
		box-shadow: none !important;
		transition: all 0.3s ease;
	}

	.original-search-field :deep(.v-field--focused) {
		background: rgba(255, 255, 255, 0.2) !important;
		border-color: rgba(255, 255, 255, 0.3) !important;
	}

	.original-search-field :deep(input) {
		color: white !important;
		font-weight: 400;
	}

	.original-search-field :deep(.v-field__prepend-inner) {
		color: rgba(255, 255, 255, 0.6) !important;
	}

	.header-avatar-glass {
		background: rgba(255, 255, 255, 0.1) !important;
		backdrop-filter: blur(10px);
	}

	.welcome-text {
		color: white;
		font-size: 12px;
		opacity: 0.9;
		line-height: 1;
	}

	.username-text {
		color: white;
		font-size: 14px;
		font-weight: 800;
	}


	/* iOS Suggestions */
	.ios-suggestions {
		background: rgba(255, 255, 255, 0.9);
		backdrop-filter: blur(20px);
		border-radius: 16px;
		box-shadow: 0 10px 40px rgba(0,0,0,0.3);
		overflow: hidden;
		z-index: 9999;
	}

	.ios-suggestion-item {
		padding: 12px 20px;
		cursor: pointer;
		font-size: 15px;
		color: #333;
		transition: background 0.2s;
	}

	.ios-suggestion-item:hover {
		background: rgba(0, 122, 255, 0.1);
	}

	/* iOS Buttons */
	.ios-btn-primary {
		background: var(--ios-cyan) !important;
		color: white !important;
		border-radius: 12px !important;
		text-transform: none !important;
		font-weight: 600 !important;
		padding: 0 24px !important;
		transition: transform 0.2s var(--spring-easing) !important;
	}

	.ios-btn-primary:active {
		transform: scale(0.95);
	}

	.ios-btn-ghost {
		color: white !important;
		transition: opacity 0.2s !important;
	}

	.ios-btn-ghost:hover {
		opacity: 0.7;
	}

	/* Transitions */
	.ios-page-enter-active,
	.ios-page-leave-active {
		transition: opacity 0.4s var(--spring-easing), transform 0.4s var(--spring-easing);
	}

	.ios-page-enter-from {
		opacity: 0;
		transform: translateY(20px) scale(0.98);
	}

	.ios-page-leave-to {
		opacity: 0;
		transform: translateY(-20px) scale(1.02);
	}

	.fade-enter-active, .fade-leave-active {
		transition: opacity 0.3s;
	}
	.fade-enter-from, .fade-leave-to {
		opacity: 0;
	}

	/* Responsive Media Queries */
	@media (max-width: 960px) {
		.logo-img {
			height: 48px;
			margin: 0 10px;
		}

		.app-title {
			font-size: 20px !important;
		}
	}

	@media (max-width: 600px) {
		.glass-app-bar {
			height: 64px !important;
		}

		.logo-img {
			height: 40px;
			margin: 0 5px;
		}

		.ios-search-input {
			font-size: 14px;
		}

		.ios-btn-primary {
			padding: 0 12px !important;
			font-size: 13px !important;
		}

		.search-container {
			margin: 0 5px;
		}

		.search-wrapper {
			padding: 2px 12px;
		}
	}
</style>
