<template>
	<v-app>
		<v-main>
			<v-app-bar :elevation="2" color="#242424">
				<template v-slot:prepend>
					<v-app-bar-nav-icon></v-app-bar-nav-icon>
				</template>

				<v-app-bar-title>Application Bar</v-app-bar-title>
			</v-app-bar>

			<v-btn>
				<v-icon color="green-darken-2" icon="mdi-domain" size="large" @click="$router.push({ path: '/' })"></v-icon>
				Botão Vuetify 09
			</v-btn>
		</v-main>

		<Footer />
	</v-app>
</template>

<script>
import Footer from './components/Footer.vue';

export default {
	name: 'App',
	components: {Footer}
}
</script>

<style>
#app {
	font-family: Avenir, Helvetica, Arial, sans-serif;
	-webkit-font-smoothing: antialiased;
	-moz-osx-font-smoothing: grayscale;
	text-align: center;
	color: #2c3e50;
	margin-top: 60px;
}
</style>
