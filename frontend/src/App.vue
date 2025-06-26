<template>
	<v-app>
		<v-app-bar :elevation="2" color="#242424" height="80">
			<template v-slot:prepend>
				<img src="@/assets/images/site-images/login/img-logo-menu-bar.png" alt="Logo" style="height: 70px; margin-right: 10px; margin-left: 20px;" />
				<v-app-bar-title>BIBLIOTECA DIGITAL</v-app-bar-title>
			</template>

			<div class="input-wrapper">
				<v-icon>mdi-magnify</v-icon>
				<input type="text" class="custom-input w-75" />
			</div>

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
		background-color: #3a6391 !important;
	}

	.v-app-bar-title {
		font-size: 35px !important;
		font-weight: bolder;
	}

	.input-wrapper {
		display: flex;
		justify-content: center;
		align-items: center;
		flex: 1;
	}

	.custom-input {
		width: 300px;
		border: 2px solid black;
		border-radius: 25px;
		padding: 8px 12px;
		font-size: 16px;
		color: #000;
		background-color: white;
	}
</style>
