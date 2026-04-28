<template>
	<v-app class="ios-app" :class="['platform-' + platform]">
		<!-- Mesh Gradient Background -->
		<div class="premium-bg">
			<div class="mesh-gradient"></div>
		</div>

		<!-- Conditionally render App Bar -->
		<v-app-bar app v-if="showBar" :elevation="0" class="glass-app-bar" height="80">
			<!-- Logo Section -->
			<div class="header-logo-container d-flex align-center ml-2 ml-sm-8 logo-clickable" @click="$router.push('/')">
				<img src="@/assets/images/site-images/login/img-logo-menu-bar.png" alt="Logo" class="logo-img-original" />
				<div class="logo-text-stack ml-3 hidden-sm-and-down">
					<h1 class="original-title">ACERVUS</h1>
					<div class="subtitle-accent">CORE</div>
				</div>
			</div>

			<v-spacer></v-spacer>

			<!-- Central Search Bar -->
			<div class="search-wrapper-global mx-4">
				<v-text-field
					v-model="searchInput"
					ref="globalSearch"
					placeholder="Buscar conhecimento..."
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
				</v-text-field>
			</div>

			<v-spacer></v-spacer>

			<!-- Header Desktop Links -->
			<div class="header-links hidden-md-and-down mr-6">
				<router-link to="/explorar" class="header-btn"><v-icon size="18" class="mr-2">mdi-compass-outline</v-icon> Explorar</router-link>
				<router-link v-if="isLoggedIn" to="/dashboard" class="header-btn"><v-icon size="18" class="mr-2">mdi-view-dashboard-outline</v-icon> Dashboard</router-link>
				<router-link to="/sobre-nos" class="header-btn"><v-icon size="18" class="mr-2">mdi-information-outline</v-icon> Sobre Nós</router-link>
			</div>

			<!-- User Actions -->
			<div class="nav-actions-wrapper d-flex align-center mr-2 mr-sm-8" style="gap: 8px;">

				<!-- Notification Menu -->
				<v-menu
					v-model="notificationMenuOpen"
					:close-on-content-click="false"
					location="bottom end"
					offset="12"
					transition="slide-y-transition"
				>
					<template v-slot:activator="{ props }">
						<v-btn icon variant="text" v-bind="props" class="notification-btn" @click="fetchNotifications">
							<v-badge
								v-if="unreadCount > 0"
								:content="unreadCount"
								color="primary"
								offset-x="2"
								offset-y="2"
							>
								<v-icon size="22">mdi-bell-outline</v-icon>
							</v-badge>
							<v-icon v-else size="22">mdi-bell-outline</v-icon>
						</v-btn>
					</template>

					<v-card class="notification-dropdown-card" width="320" elevation="16">
						<div class="pa-4 d-flex align-center justify-space-between">
							<span class="text-subtitle-1 font-weight-bold">Notificações</span>
							<div class="d-flex align-center" style="gap: 8px;">
								<v-btn v-if="unreadCount > 0" variant="text" size="x-small" color="primary" class="text-none" @click="markAllAsRead">Lidas</v-btn>
								<v-btn v-if="notifications.length > 0" variant="text" size="x-small" color="error" class="text-none" @click="clearAllNotifications">Limpar</v-btn>
							</div>
						</div>
						<v-divider></v-divider>
						<v-list v-if="notifications.length > 0" class="notification-list pa-0" max-height="400">
							<v-list-item
								v-for="n in notifications"
								:key="n.id"
								:class="{ 'unread-item': !n.lida }"
								class="notification-item py-3"
								@click="markAsRead(n)"
							>
								<template v-slot:prepend>
									<v-avatar size="32" :color="getNotificationColor(n.tipo).bg" class="mr-3">
										<v-icon size="18" :color="getNotificationColor(n.tipo).icon">{{ getNotificationIcon(n.tipo) }}</v-icon>
									</v-avatar>
								</template>
								<v-list-item-title class="text-subtitle-2 font-weight-bold">{{ n.titulo }}</v-list-item-title>
								<v-list-item-subtitle class="text-caption mt-1">{{ n.mensagem }}</v-list-item-subtitle>
								<template v-slot:append>
									<div class="d-flex flex-column align-end">
										<div class="text-xxs opacity-40 mb-1">{{ formatTime(n.data_criacao) }}</div>
										<v-btn
											v-if="!n.lida"
											icon="mdi-check"
											size="x-small"
											variant="text"
											color="primary"
											@click.stop="markAsRead(n)"
											title="Marcar como lida"
										></v-btn>
									</div>
								</template>
							</v-list-item>
						</v-list>
						<div v-else class="pa-10 text-center opacity-40">
							<v-icon size="48" class="mb-2">mdi-bell-off-outline</v-icon>
							<div class="text-body-2">Nenhuma notificação por aqui</div>
						</div>
					</v-card>
				</v-menu>

				<!-- Unified User Toggle -->
				<v-btn
					icon
					variant="text"
					@click="drawer = !drawer"
					class="ml-2 avatar-toggle-btn"
					title="Menu do Usuário"
				>
					<v-avatar size="40" class="header-avatar-glass">
						<v-img v-if="isLoggedIn && userAvatar" :src="userAvatar" cover></v-img>
						<v-icon v-else color="white" size="24">{{ isLoggedIn ? 'mdi-account-circle' : 'mdi-account-circle-outline' }}</v-icon>
					</v-avatar>
				</v-btn>
			</div>
		</v-app-bar>

		<!-- Menu Lateral (Sidebar) -->
		<v-navigation-drawer
			v-model="drawer"
			location="right"
			temporary
			touchless
			:permanent="false"
			scrim="rgba(0, 10, 20, 0.7)"
			:width="isMobile ? 280 : 350"
			class="premium-drawer"
		>
			<div class="drawer-header pa-6 pt-10">
				<div class="d-flex align-center mb-6">
					<v-avatar size="64" class="header-avatar-glass mr-4">
						<v-img v-if="isLoggedIn && userAvatar" :src="userAvatar" cover></v-img>
						<v-icon v-else color="white" size="32">{{ isLoggedIn ? 'mdi-account' : 'mdi-account-outline' }}</v-icon>
					</v-avatar>
					<div>
						<div class="font-weight-bold text-truncate" style="max-width: 180px;">{{ isLoggedIn ? userDisplayName : 'Visitante' }}</div>
						<div class="text-caption opacity-60">{{ isLoggedIn ? userRoleName : 'Acesse sua conta' }}</div>
					</div>
				</div>
				<v-divider class="opacity-10"></v-divider>
			</div>

			<v-list class="px-4 nav-list-premium">
				<v-list-subheader class="text-xxs font-weight-bold opacity-40 mb-2 uppercase">NAVEGAÇÃO</v-list-subheader>
				
				<v-list-item
					prepend-icon="mdi-home-outline"
					title="Início"
					@click="$router.push('/'); drawer = false"
					class="drawer-item"
				></v-list-item>

				<v-list-item
					prepend-icon="mdi-compass-outline"
					title="Explorar"
					@click="$router.push('/explorar'); drawer = false"
					class="drawer-item"
				></v-list-item>

				<v-list-item
					v-if="isLoggedIn"
					prepend-icon="mdi-view-dashboard-outline"
					title="Meu Painel"
					@click="$router.push('/dashboard'); drawer = false"
					class="drawer-item"
				></v-list-item>

				<v-list-item
					v-if="isLoggedIn"
					prepend-icon="mdi-heart-outline"
					title="Favoritos"
					@click="$router.push('/favoritos'); drawer = false"
					class="drawer-item"
				></v-list-item>

				<v-divider class="my-4 opacity-10"></v-divider>
				<v-list-subheader class="text-xxs font-weight-bold opacity-40 mb-2 uppercase">CONTA</v-list-subheader>

				<template v-if="!isLoggedIn">
					<v-list-item
						prepend-icon="mdi-login"
						title="Entrar"
						@click="$router.push('/login'); drawer = false"
						class="drawer-item"
					></v-list-item>
					<v-list-item
						prepend-icon="mdi-account-plus-outline"
						title="Criar Conta"
						@click="$router.push('/cadastro'); drawer = false"
						class="drawer-item"
					></v-list-item>
				</template>

				<template v-else>
					<v-list-item
						prepend-icon="mdi-account-circle-outline"
						title="Meu Perfil"
						@click="$router.push('/perfil'); drawer = false"
						class="drawer-item"
					></v-list-item>
					<v-list-item
						prepend-icon="mdi-logout"
						title="Sair"
						color="error"
						@click="logout(); drawer = false"
						class="drawer-item danger-item"
					></v-list-item>
				</template>

				<v-divider class="my-4 opacity-10"></v-divider>
				<v-list-item
					prepend-icon="mdi-information-outline"
					title="Sobre Nós"
					@click="$router.push('/sobre-nos'); drawer = false"
					class="drawer-item"
				></v-list-item>
			</v-list>

			<template v-slot:append>
				<div class="pa-6 text-center">
					<div class="text-xxs opacity-30 font-weight-bold">ACERVUS CORE v1.0</div>
				</div>
			</template>
		</v-navigation-drawer>

		<v-main>
			<!-- Private/Main App Layout Wrapper -->
			<div v-if="showBar" class="responsive-container">
				<router-view v-slot="{ Component }">
					<transition name="ios-page" mode="out-in">
						<keep-alive include="ExplorePageExtended,HomeView">
							<component :is="Component" :key="$route.fullPath" />
						</keep-alive>
					</transition>
				</router-view>
			</div>
			<!-- Public Layout Wrapper (No app bar, full screen auth views) -->
			<div v-else class="public-page-container">
				<router-view v-slot="{ Component }">
					<transition name="ios-page" mode="out-in">
						<keep-alive>
							<component :is="Component" :key="$route.fullPath" />
						</keep-alive>
					</transition>
				</router-view>
			</div>
		</v-main>

		<!-- Footer -->
		<AppFooter v-if="showBar" />

		<!-- Global Accessibility Panel -->

		<!-- Global Accessibility Panel -->
		<AccessibilityPanel />

		<!-- Global Notifications (iOS Style) -->
		<IOSNotification />
	</v-app>
</template>

<script>
import AccessibilityPanel from './components/AccessibilityPanel.vue'
import IOSNotification from './components/IOSNotification.vue'
import AppFooter from './components/Footer.vue'
import { iosNotificationStore } from '@/services/IOSNotificationStore'
import '@/assets/styles/premium.css'
import { ref, computed, reactive } from 'vue'
import auth, { state as authState } from '@/auth'
import NotificationService from '@/services/NotificationService'
import MaterialService from '@/services/MaterialService'
import { useAccessibility } from '@/composables/useAccessibility'
import { useDevice } from '@/composables/useDevice'

export default {
	name: 'App',
	components: { AccessibilityPanel, IOSNotification, AppFooter },
	data() {
		return {
			publicRoutes: ['/login', '/cadastro', '/esqueci-senha'],
			searchInput: '',
			isSearchFocused: false,
			notificationMenuOpen: false,
			notifications: [],
			lastWelcomeId: null,
			drawer: false
		}
	},
	provide() {
		return {
			notify: this.notify.bind(this),
			fetchGlobalFavorites: this.fetchFavorites.bind(this),
			getGlobalFavorites: this.fetchFavorites.bind(this),
			globalFavorites: this.favoritesStore
		}
	},
	methods: {
		notify(text, type = 'info') {
			iosNotificationStore.addNotification(text, type);
		},
		logout() {
			auth.logout()
			this.favoritesStore.list = [] // Clear favorites to prevent ghost books for next user
			this.$router.push('/')
		},
		async doSearch() {
			const q = this.searchInput.trim()
			if (!q) return

			this.loading = true
			try {
				// Pequeno delay para feedback visual se for muito rápido
				await new Promise(resolve => setTimeout(resolve, 500))

				// ALWAYS redirect to /explorar to show search results properly
				if (this.$route.path !== '/explorar') {
					this.$router.push({ path: '/explorar', query: { q, categoria: 'TODOS' } })
				} else {
					// If already on /explorar, update the query which will trigger the watch and fetch results
					this.$router.replace({ path: '/explorar', query: { q, categoria: 'TODOS' } })
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
		},
		async fetchNotifications() {
			if (!this.isLoggedIn) return;
			try {
				const res = await NotificationService.listar(authState.user.id);
				this.notifications = res.data || [];
			} catch (e) {
				console.error(e);
				this.notifications = [];
			}
		},
		async markAsRead(n) {
			if (n.lida) return;
			try {
				await NotificationService.marcarComoLida(n.id);
				n.lida = true;
			} catch (e) { console.error(e); }
		},
		async clearAllNotifications() {
			try {
				if (this.isLoggedIn) {
					await NotificationService.limparTudo(authState.user.id);
					this.notifications = [];
				}
			} catch (e) {
				console.error('Erro ao limpar notificações:', e);
			}
		},
		async markAllAsRead() {
			const unread = this.notifications.filter(n => !n.lida);
			for (const n of unread) {
				await this.markAsRead(n);
			}
		},
		getNotificationIcon(tipo) {
			const icons = {
				'boas-vindas': 'mdi-hand-wave',
				'conquista': 'mdi-trophy',
				'meta': 'mdi-target',
				'info': 'mdi-information'
			};
			return icons[tipo] || 'mdi-bell';
		},
		getNotificationColor(tipo) {
			const colors = {
				'boas-vindas': { bg: 'rgba(0, 122, 255, 0.1)', icon: 'primary' },
				'conquista': { bg: 'rgba(255, 193, 7, 0.1)', icon: 'amber' },
				'meta': { bg: 'rgba(76, 175, 80, 0.1)', icon: 'green' },
				'info': { bg: 'rgba(255, 255, 255, 0.05)', icon: 'white' }
			};
			return colors[tipo] || { bg: 'rgba(255, 255, 255, 0.05)', icon: 'white' };
		},
		formatTime(dateStr) {
			if (!dateStr) return '';
			// Garantir que a data do backend seja interpretada como UTC se vier sem timezone explicito
			const date = new Date(dateStr.endsWith('Z') || dateStr.includes('+') ? dateStr : dateStr.replace(' ', 'T') + 'Z');
			return date.toLocaleTimeString('pt-BR', { hour: '2-digit', minute: '2-digit' });
		},
		handleLoginEffect() {
			if (this.isLoggedIn) {
				const sessionKey = `welcome_notified_${authState.user.id}`;
				if (sessionStorage.getItem(sessionKey)) return;

				this.lastWelcomeId = authState.user.id;
				sessionStorage.setItem(sessionKey, 'true');

				const now = new Date();
				const dateTime = now.toLocaleString('pt-BR', {
					weekday: 'long',
					day: 'numeric',
					month: 'long',
					hour: '2-digit',
					minute: '2-digit'
				});

				// Dynamic notify
				this.notify(`Bem-vindo de volta! Acesso em: ${dateTime}`, 'info');

				// Create persistent welcome notification if it doesn't exist
				NotificationService.criar({
					usuario_id: authState.user.id,
					titulo: 'Boas-vindas!',
					mensagem: `Você acessou às ${now.toLocaleTimeString('pt-BR')}`,
					tipo: 'boas-vindas'
				}).then(() => this.fetchNotifications());
			}
		},
		async checkNewAchievements(newStats) {
			if (!newStats || !newStats.badges) return;

			const storageKey = `badges_seen_${authState.user.id}`;
			const seenBadges = JSON.parse(localStorage.getItem(storageKey) || '[]');

			const newBadges = newStats.badges.filter(b => !seenBadges.includes(b));

			if (newBadges.length > 0) {
				newBadges.forEach(badge => {
					this.notify(`Nova Conquista Desbloqueada: ${badge}!`, 'success');

					// Persistir no banco como notificação também
					NotificationService.criar({
						usuario_id: authState.user.id,
						titulo: 'Nova Conquista!',
						mensagem: `Você desbloqueou a medalha: ${badge}`,
						tipo: 'conquista'
					});
				});

				localStorage.setItem(storageKey, JSON.stringify([...seenBadges, ...newBadges]));
				this.fetchNotifications();
			}
		},
		async fetchFavorites() {
			if (!this.isLoggedIn) return;
			try {
				const res = await MaterialService.listarFavoritos(authState.user.id);
				this.favoritesStore.list = res.data || [];
			} catch (e) { console.error(e); }
		},
		startPolling() {
			this.stopPolling(); // Safely clear existing
			this.refreshInterval = setInterval(this.fetchNotifications, 60000);
		},
		stopPolling() {
			if (this.refreshInterval) {
				clearInterval(this.refreshInterval);
				this.refreshInterval = null;
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
		userRoleName() {
			const user = authState.user
			// 1: Aluno, 2: Professor (Baseado no fluxo de cadastro/banco)
			if (user?.tipo === 1 || user?.tipo_usuario_id === 1) return 'Aluno';
			if (user?.tipo === 2 || user?.tipo_usuario_id === 2) return 'Professor';
			return 'Membro';
		},
		userAvatar() {
			const user = authState.user
			return user?.foto_url || null
		},
		unreadCount() {
			if (!this.notifications || !Array.isArray(this.notifications)) return 0;
			return this.notifications.filter(n => !n.lida).length;
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
		},
		isLoggedIn(val) {
			if (val) {
				this.handleLoginEffect();
				this.fetchNotifications();
				this.fetchFavorites();
				// Reinicia o polling se logar
				this.startPolling();
			} else {
				this.notifications = [];
				this.lastWelcomeId = null;
				this.stopPolling();
			}
		}
	},
	setup() {
		const loading = ref(false)
		const { init } = useAccessibility()
		const { platform, isMobile } = useDevice()

		const favoritesStore = reactive({
			list: []
		})

		return {
			loading,
			init,
			favoritesStore,
			platform,
			isMobile,
			favoritos: computed(() => favoritesStore.list),
			cleanupAccessibility: useAccessibility().cleanup
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

		if (this.isLoggedIn) {
			this.handleLoginEffect();
			this.fetchNotifications();
			this.fetchFavorites();
			this.startPolling();
		}
	},
	beforeUnmount() {
		window.removeEventListener('keydown', this.handleGlobalKeydown);
		this.stopPolling();
		if (this.cleanupAccessibility) {
			this.cleanupAccessibility()
		}
	}
}
</script>

<style>
	/* Global Defaults */

	@import url('https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700;800;900&display=swap');

	body {
		margin: 0;
		padding: 0;
		font-family: 'Inter', var(--apple-font);
		overflow-x: hidden;
		background: radial-gradient(circle at 30% 30%, #102136 0%, #0B192C 100%) fixed !important;
		color: #ffffff !important;
		-webkit-font-smoothing: antialiased;
	}

	.v-application {
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
		transition: all 0.3s cubic-bezier(0.16, 1, 0.3, 1);
		cursor: pointer;
	}

	.header-logo-container:hover {
		transform: scale(1.03) translateY(-1px);
		filter: drop-shadow(0 4px 12px rgba(0, 122, 255, 0.15));
	}

	.header-logo-container:active {
		transform: scale(0.97);
	}

	.logo-img-original {
		height: 36px;
		filter: drop-shadow(0 4px 12px rgba(0,0,0,0.1));
		transition: transform 0.4s cubic-bezier(0.16, 1, 0.3, 1);
	}

	.header-logo-container:hover .logo-img-original {
		transform: rotate(5deg) scale(1.05);
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
		transition: color 0.3s ease;
	}

	.header-logo-container:hover .original-title {
		color: #007AFF;
	}

	.subtitle-accent {
		font-size: 10px;
		font-weight: 900;
		color: var(--ios-cyan);
		letter-spacing: 2px;
		transition: all 0.3s ease;
	}

	.header-links {
		display: flex;
		gap: 16px;
		align-items: center;
	}
	
	.header-btn {
		color: rgba(255,255,255,0.85);
		text-decoration: none;
		font-size: 0.95rem;
		font-weight: 600;
		transition: all 0.4s var(--spring-easing);
		padding: 8px 18px;
		border-radius: 100px;
		display: flex;
		align-items: center;
		background: rgba(255, 255, 255, 0.05);
		border: 1px solid rgba(255, 255, 255, 0.08);
	}
	
	.header-btn:hover {
		color: #ffffff;
		background: rgba(255, 255, 255, 0.12);
		border-color: rgba(255, 255, 255, 0.2);
		transform: translateY(-2px);
		box-shadow: 0 8px 16px rgba(0,0,0,0.2);
	}

	.header-logo-container:hover .subtitle-accent {
		color: #00E5FF;
		text-shadow: 0 0 10px rgba(0, 229, 255, 0.4);
	}

	/* Enhanced Search Bar */
	.search-wrapper-global {
		flex: 1 1 auto;
		max-width: 600px;
		width: 100%;
	}

	.ios-search-field :deep(.v-field) {
		background: rgba(var(--v-theme-on-surface), 0.05) !important;
		backdrop-filter: blur(10px);
		border-radius: 16px !important;
		height: 44px !important;
		transition: all 0.4s var(--spring-easing);
		border: 1px solid rgba(255, 255, 255, 0.08) !important;
		box-shadow: inset 0 1px 1px rgba(255, 255, 255, 0.05);
	}

	.search-active :deep(.v-field) {
		background: rgba(var(--v-theme-surface), 0.4) !important;
		backdrop-filter: blur(25px) saturate(200%) !important;
		border-color: rgba(0, 122, 255, 0.5) !important;
		box-shadow: 0 10px 30px rgba(0, 0, 0, 0.2), 0 0 0 4px rgba(0, 122, 255, 0.1) !important;
		transform: translateY(-1px);
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

	/* Transitions Style iOS 17 - Platform Optimized */
	.ios-page-enter-active,
	.ios-page-leave-active {
		transition: opacity 0.4s var(--spring-easing),
					transform 0.4s var(--spring-easing);
	}

	/* Mobile Optimization: Snappier transitions */
	.platform-mobile .ios-page-enter-active,
	.platform-mobile .ios-page-leave-active {
		transition: opacity 0.25s ease-out, transform 0.25s ease-out;
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
		padding: clamp(24px, 5vw, 64px) clamp(16px, 5vw, 64px) !important;
		overflow-y: auto !important;
		width: 100%;
		background: transparent !important;
	}

	.responsive-container {
		width: 100%;
		max-width: 1600px;
		margin: 0 auto;
		padding-left: clamp(16px, 3vw, 40px) !important;
		padding-right: clamp(16px, 3vw, 40px) !important;
		padding-top: clamp(16px, 2vw, 24px) !important;
		padding-bottom: clamp(24px, 4vw, 60px) !important;
		background: transparent !important;
	}

	/* Mobile-Specific UX Improvements */
	.platform-mobile .search-wrapper-global {
		max-width: 100%;
		margin: 0 8px;
	}

	.platform-mobile .logo-text-stack {
		display: none; /* Already handled by hidden-sm-and-down, but enforcing for UX */
	}

	.platform-mobile .glass-app-bar {
		height: 70px !important;
	}

	.platform-mobile .ios-search-field :deep(.v-field) {
		height: 40px !important;
		border-radius: 12px !important;
	}

	/* Overall Fluidity: GPU acceleration for core elements */
	.platform-mobile .responsive-container {
		padding-left: 12px !important;
		padding-right: 12px !important;
		padding-top: 10px !important;
	}

	/* Overall Fluidity: GPU acceleration for core elements */
	.ios-app * {
		backface-visibility: hidden;
		-webkit-backface-visibility: hidden;
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
		.glass-app-bar {
			height: 64px !important;
		}
		.logo-img-original {
			height: 28px;
		}
		.search-wrapper-global {
			margin: 0 8px;
			max-width: 180px;
		}
		.ios-search-field :deep(.v-field) {
			height: 36px !important;
			font-size: 13px !important;
		}
		.responsive-container {
			padding-left: 16px !important;
			padding-right: 16px !important;
			padding-top: 16px !important;
		}
		.nav-actions-wrapper .text-right {
			display: none;
		}
		.header-avatar-glass {
			size: 36px !important;
		}
	}

	/* Premium Snackbar Styling */
	.premium-snackbar :deep(.v-snackbar__wrapper) {
		background: rgba(var(--v-theme-surface), 0.8) !important;
		backdrop-filter: blur(15px) saturate(150%) !important;
		-webkit-backdrop-filter: blur(15px) saturate(150%) !important;
		border: 1px solid rgba(255, 255, 255, 0.1) !important;
		box-shadow: 0 12px 32px rgba(0, 0, 0, 0.2) !important;
		color: rgb(var(--v-theme-on-surface)) !important;
	}

	.notification-dropdown-card,
	.user-dropdown-card {
		background: rgba(20, 30, 45, 0.7) !important;
		backdrop-filter: blur(30px) saturate(210%) !important;
		-webkit-backdrop-filter: blur(30px) saturate(210%) !important;
		border: 1px solid rgba(255, 255, 255, 0.12) !important;
		border-radius: 24px !important;
		box-shadow: 0 30px 60px rgba(0, 0, 0, 0.4), 0 0 0 1px rgba(255, 255, 255, 0.05) inset !important;
		overflow: hidden;
		animation: menu-pop 0.4s var(--spring-easing);
	}

	@keyframes menu-pop {
		from { opacity: 0; transform: translateY(10px) scale(0.95); }
		to { opacity: 1; transform: translateY(0) scale(1); }
	}

	.dropdown-item, .notification-item {
		transition: all 0.2s ease !important;
		margin: 4px 8px !important;
		border-radius: 12px !important;
	}

	.dropdown-item:hover, .notification-item:hover {
		background: rgba(255, 255, 255, 0.08) !important;
		transform: translateX(4px);
	}

	.unread-item {
		background: rgba(0, 122, 255, 0.05) !important;
	}

	.fav-mini-item {
		min-height: 48px !important;
		border-radius: 12px !important;
		margin: 4px 8px !important;
		transition: all 0.2s ease;
	}
	.fav-mini-item:hover {
		background: rgba(255, 255, 255, 0.08) !important;
		transform: translateX(4px);
	}
	.fav-mini-item :deep(.v-list-item-title) {
		font-size: 0.9rem !important;
		font-weight: 600 !important;
		color: white !important;
	}
	.fav-mini-item :deep(.v-list-item-subtitle) {
		font-size: 0.75rem !important;
		opacity: 0.7;
	}

	/* Premium Drawer Styling */
	.premium-drawer {
		background: rgba(15, 25, 40, 0.85) !important;
		backdrop-filter: blur(40px) saturate(200%) !important;
		-webkit-backdrop-filter: blur(40px) saturate(200%) !important;
		border-left: 1px solid rgba(255, 255, 255, 0.1) !important;
		box-shadow: -20px 0 50px rgba(0, 0, 0, 0.5) !important;
	}

	.drawer-item {
		border-radius: 16px !important;
		margin-bottom: 4px !important;
		transition: all 0.3s var(--spring-easing) !important;
		color: rgba(255, 255, 255, 0.7) !important;
	}

	.drawer-item:hover {
		background: rgba(255, 255, 255, 0.08) !important;
		transform: translateX(-4px);
		color: white !important;
	}

	.drawer-item :deep(.v-icon) {
		opacity: 0.6;
		transition: transform 0.3s ease;
	}

	.drawer-item:hover :deep(.v-icon) {
		opacity: 1;
		transform: scale(1.1);
		color: var(--ios-blue);
	}

	.danger-item:hover {
		background: rgba(255, 59, 48, 0.1) !important;
		color: #FF6B6B !important;
	}

	.text-xxs {
		font-size: 10px;
		letter-spacing: 1px;
	}
</style>
