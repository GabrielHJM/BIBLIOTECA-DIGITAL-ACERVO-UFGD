<template>
	<v-app class="ios-app">
		<v-app-bar v-if="showBar" :elevation="0" class="glass-app-bar" height="80">
			<template v-slot:prepend>
				<img src="@/assets/images/site-images/login/img-logo-menu-bar.png" alt="Logo" class="logo-img" />
				<v-app-bar-title class="app-title">BIBLIOTECA DIGITAL</v-app-bar-title>
			</template>

			<v-spacer />

			<div class="search-container" ref="inputWrapperRef">
				<div class="search-wrapper" :class="{ 'search-active': isSearchFocused }">
					<i class="mdi mdi-magnify search-icon"></i>
					<input
						type="text"
						class="ios-search-input"
						placeholder="Buscar materiais..."
						v-model="searchInput"
						@focus="isSearchFocused = true"
						@blur="isSearchFocused = false"
					/>
				</div>

				<Teleport to="body">
					<transition name="fade">
						<div
							class="ios-suggestions"
							v-if="suggestions.length && searchInput"
							:style="dropdownStyle"
						>
							<div
								class="ios-suggestion-item"
								v-for="(item, index) in suggestions"
								:key="index"
								@click="selectSuggestion(item)"
							>
								{{ item }}
							</div>
						</div>
					</transition>
				</Teleport>
			</div>

			<v-spacer />

			<div class="nav-actions">
				<router-link v-if="$route.path === '/'" to="/login">
					<v-btn class="ios-btn-primary" elevation="0">
						<v-icon class="mr-2">mdi-login</v-icon>
						Entrar
					</v-btn>
				</router-link>

				<template v-else>
					<v-btn class="ios-btn-ghost mr-2" elevation="0" @click="logout" icon="mdi-logout"></v-btn>
					<v-btn class="ios-btn-ghost mr-2" elevation="0" icon="mdi-heart-outline"></v-btn>
					<v-btn class="ios-btn-ghost" elevation="0" icon="mdi-dots-vertical"></v-btn>
				</template>
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
	</v-app>
</template>

<script>
import Footer from './components/Footer.vue'
import { ref, watch, onMounted, onBeforeUnmount, nextTick } from 'vue'
import debounce from 'lodash.debounce'
import livros from '../../livros.json'
import auth from '@/auth'

export default {
	name: 'App',
	components: { Footer },
	data() {
		return {
			publicRoutes: ['/login', '/cadastro', '/esqueci-senha'],
			isSearchFocused: false
		}
	},
	methods: {
		logout() {
			auth.logout()
			this.$router.push('/')
		}
	},
	computed: {
		showBar() {
			return !this.publicRoutes.includes(this.$route.path)
		}
	},
	watch: {
		'$route.path'() {
			// força reatividade
			this.showBar
		}
	},
	setup() {
		const searchInput = ref('')
		const suggestions = ref([])
		const loading = ref(false)
		const inputWrapperRef = ref(null)
		const dropdownStyle = ref({ top: '0px', left: '0px', width: '300px' })

		const fetchSuggestions = debounce(async (query) => {
			if (!query || query.length < 2) {
				suggestions.value = []
				return
			}

			loading.value = true
			try {
				// Simulação com JSON local
				const data = livros
				suggestions.value = data
				.filter((item) => item.nome.toLowerCase().includes(query.toLowerCase()))
				.map((item) => item.nome)
			} catch (err) {
				console.error(err)
			} finally {
				loading.value = false
			}
		}, 400)

		const updateDropdownPosition = () => {
			if (inputWrapperRef.value) {
				const rect = inputWrapperRef.value.getBoundingClientRect()
				dropdownStyle.value = {
					top: `${rect.bottom + window.scrollY + 8}px`,
					left: `${rect.left + window.scrollX}px`,
					width: `${rect.width}px`,
					position: 'absolute',
				}
			}
		}

		const selectSuggestion = (item) => {
			searchInput.value = item
			suggestions.value = []
		}

		watch(searchInput, async (val) => {
			if (!val) {
				suggestions.value = []
				return
			}
			await nextTick()
			updateDropdownPosition()
			fetchSuggestions(val)
		})

		const handleClickOutside = (event) => {
			if (
				inputWrapperRef.value &&
				!inputWrapperRef.value.contains(event.target)
			) {
				suggestions.value = []
			}
		}

		onMounted(() => {
			document.addEventListener('click', handleClickOutside)
			window.addEventListener('resize', updateDropdownPosition)
			window.addEventListener('scroll', updateDropdownPosition, true)
		})

		onBeforeUnmount(() => {
			document.removeEventListener('click', handleClickOutside)
			window.removeEventListener('resize', updateDropdownPosition)
			window.removeEventListener('scroll', updateDropdownPosition, true)
		})

		return {
			searchInput,
			suggestions,
			loading,
			selectSuggestion,
			inputWrapperRef,
			dropdownStyle,
		}
	},
}
</script>

<style>
	:root {
		--ios-blue: #007AFF;
		--ios-cyan: #5AC8FA;
		--ios-bg: #3a6391;
		--ios-card: rgba(255, 255, 255, 0.95);
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
		background: rgba(36, 36, 36, 0.7) !important;
		backdrop-filter: blur(20px) saturate(180%) !important;
		-webkit-backdrop-filter: blur(20px) saturate(180%) !important;
		transition: background 0.3s ease;
	}

	.app-title {
		font-size: 28px !important;
		font-weight: 800 !important;
		letter-spacing: -0.5px;
		color: white !important;
	}

	.logo-img {
		height: 60px;
		margin: 0 15px 0 20px;
		filter: drop-shadow(0 2px 4px rgba(0,0,0,0.2));
	}

	/* iOS Search Input */
	.search-container {
		flex: 1;
		max-width: 500px;
		margin: 0 20px;
	}

	.search-wrapper {
		position: relative;
		background: rgba(255, 255, 255, 0.15);
		border-radius: 20px;
		padding: 4px 16px;
		display: flex;
		align-items: center;
		transition: all 0.3s var(--spring-easing);
		border: 1px solid rgba(255, 255, 255, 0.1);
	}

	.search-active {
		background: rgba(255, 255, 255, 0.25) !important;
		transform: scale(1.02);
		box-shadow: 0 4px 20px rgba(0,0,0,0.2);
	}

	.search-icon {
		font-size: 20px;
		color: rgba(255, 255, 255, 0.7);
		margin-right: 12px;
	}

	.ios-search-input {
		width: 100%;
		border: none !important;
		background: transparent !important;
		color: white !important;
		font-size: 16px;
		padding: 8px 0;
	}

	.ios-search-input::placeholder {
		color: rgba(255, 255, 255, 0.5);
	}

	.ios-search-input:focus {
		outline: none;
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
</style>
