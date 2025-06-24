<template>
	<v-container class="d-flex justify-center border">
		<v-card class="ma-7 pa-5" style="width: 80%; height: 700px; background-color: #3a6391;" elevation="0">
			
		</v-card>
	</v-container>
</template>

<script>
import livros from '../../../livros.json'

export default {
	name: 'CadastroPage',
	data: () => ({
		pack: require('../../package.json'),
		livros: livros,
		caminho: require("@/assets/images/site-images/login/person-img.png")
	}),
	mounted() {
		console.log('livro: ', this.livros[0])
	}
}
</script>

<style>
	.border {
		border: 1px solid white;
	}
</style>
