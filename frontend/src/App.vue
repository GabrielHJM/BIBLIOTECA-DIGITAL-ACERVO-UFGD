<template>
	<v-app>
		<v-app-bar :elevation="2" color="#242424" height="80">
			<template v-slot:prepend>
				<img src="@/assets/images/site-images/login/img-logo-menu-bar.png" alt="Logo" style="height: 70px; margin-right: 10px; margin-left: 20px;" />
				<v-app-bar-title>BIBLIOTECA DIGITAL</v-app-bar-title>
			</template>

			<v-spacer />

			<div class="input-wrapper" ref="inputWrapperRef">
				<i class="mdi mdi-magnify input-icon"></i>
				<input type="text" class="custom-input w-100" placeholder="Pesquisar por material..." v-model="searchInput" />

				<Teleport to="body">
					<div
						class="suggestions-dropdown"
						v-if="suggestions.length && searchInput"
						:style="dropdownStyle"
					>
						<div
							class="suggestion-item"
							v-for="(item, index) in suggestions"
							:key="index"
							@click="selectSuggestion(item)"
						>
							{{ item }}
						</div>
					</div>
				</Teleport>
			</div>

			<v-spacer />

			<v-btn icon>
				<v-icon>mdi-magnify</v-icon>
			</v-btn>

			<v-btn icon>
				<v-icon>mdi-heart</v-icon>
			</v-btn>

			<v-btn icon>
				<v-icon>mdi-dots-vertical</v-icon>
			</v-btn>
		</v-app-bar>

		<v-main>
			<v-container fluid="false">
				<router-view />
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

export default {
	name: 'App',
	components: { Footer },
	setup() {
		const search = ref('')
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
				top: `${rect.bottom + window.scrollY}px`,
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
			console.log('input:', val, 'sugestões:', suggestions.value)
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
			search,
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
	#app {
		font-family: Avenir, Helvetica, Arial, sans-serif;
		-webkit-font-smoothing: antialiased;
		-moz-osx-font-smoothing: grayscale;
		text-align: center;
		color: #2c3e50;
		margin: 0;
		padding: 0;
		background-color: red;
	}

	.v-application {
		background-color: #3a6391 !important;
	}

	.v-app-bar-title {
		font-size: 35px !important;
		font-weight: bolder;
	}

	.input-wrapper {
		position: relative;
		width: 300px;
		display: flex;
		align-items: center;
		flex: 15;
		margin-left: 20px;
	}

	.input-icon {
		position: absolute;
		left: 12px;
		font-size: 20px;
		color: #555;
	}

	.custom-input {
		width: 100%;
		padding: 8px 12px 8px 36px;
		border: 2px solid black;
		border-radius: 25px;
		font-size: 16px;
		color: #000;
		background-color: white;
	}

	.suggestions-dropdown {
		position: absolute;
		top: 100%;
		left: 0;
		right: 0;
		background-color: white;
		border: 1px solid #ccc;
		border-radius: 0 0 12px 12px;
		box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
		z-index: 9999;
		max-height: 200px;
		overflow-y: auto;
	}

	.suggestion-item {
		padding: 8px 12px;
		cursor: pointer;
		transition: background-color 0.2s ease;
	}

	.suggestion-item:hover {
		background-color: #f0f0f0;
	}
</style>
