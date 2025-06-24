<template>
	<v-app>
		<v-app-bar :elevation="2" color="#242424" height="80">
			<template v-slot:prepend>
				<v-app-bar-nav-icon></v-app-bar-nav-icon>
			</template>

			<v-app-bar-title>Application Bar</v-app-bar-title>
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
	margin: 0;
	padding: 0;
	background-color: red;
}

.v-application {
  background-color: #3a6391 !important; /* ou #3a6391 se quiser outra */
}
</style>
