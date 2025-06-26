<template>
	<v-app>
		<v-app-bar :elevation="2" color="#242424" height="80">
			<template v-slot:prepend>
				<img src="@/assets/images/site-images/login/img-logo-menu-bar.png" alt="Logo" style="height: 70px; margin-right: 10px; margin-left: 20px;" />
				<v-app-bar-title>BIBLIOTECA DIGITAL</v-app-bar-title>
			</template>

			<v-spacer />

			<div class="input-wrapper">
				<i class="mdi mdi-magnify input-icon"></i>
				<input type="text" class="custom-input w-100" />
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
import Footer from './components/Footer.vue';
import { ref, watch } from 'vue'
import debounce from 'lodash.debounce' // npm install lodash.debounce

export default {
	name: 'App',
	components: {Footer},
	setup() {
		const search = ref('')
		const searchInput = ref('')
		const suggestions = ref([])
		const loading = ref(false)

		const fetchSuggestions = debounce(async (query) => {
			if (!query || query.length < 2) {
				suggestions.value = []
				return
			}

			loading.value = true
			try {
				// Substitua pela sua API real
				const response = await fetch(`https://api.exemplo.com/livros?q=${query}`)
				const data = await response.json()

				suggestions.value = data.map((item) => item.titulo)
			} catch (err) {
				console.error(err)
			} finally {
				loading.value = false
			}
		}, 400)

		watch(searchInput, (val) => {
			fetchSuggestions(val)
		})

		return {
			search,
			searchInput,
			suggestions,
			loading,
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
</style>
