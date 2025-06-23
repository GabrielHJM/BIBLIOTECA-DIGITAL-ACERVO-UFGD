<template>
	Hello Login
</template>

<script>
import livros from '../../../livros.json'

export default {
	name: 'LoginPage',
	data: () => ({
		pack: require('../../package.json'),
		livros: livros
	}),
	mounted() {
		console.log('livro: ', this.livros[0])
	}
}
</script>
