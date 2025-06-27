<template>
	<div class="mx-auto mt-8">
		<v-container class="pa-0 border" fluid>
			<v-card min-width="100%" rounded="xl" style="background-color: rgba(36, 36, 36, 0.50);">

			</v-card>
		</v-container>
	</div>
</template>

<script>
import livros from '../../../livros.json'

export default {
	name: 'EstudoPage',
	data: () => ({
		pack: require('../../package.json'),
		livros: livros
	}),
	mounted() {
		console.log('livro: ', this.livros[0])
	}
}
</script>
