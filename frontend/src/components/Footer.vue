<template>
	<v-footer
		class="d-flex align-center justify-center ga-2 flex-wrap flex-grow-1 py-3"
		color="#242424"
	>
		<v-btn
			v-for="link in links"
			:key="link"
			:text="link"
			variant="text"
			rounded
		></v-btn>

		<div class="flex-1-0-100 text-center mt-2">
			Biblioteca Digital (Projeto de Sistemas Computacionais II - UFGD) | <strong>Versão: 1.0</strong> - Todos os Direitos Reservados
		</div>
	</v-footer>
</template>

<script>

export default {
	name: 'Footer-Vue',
	data: () => ({
		links: [
			'Home',
			'About Us',
			'Team',
			'Services',
			'Blog',
			'Contact Us',
		],
		pack: require('../../package.json')
	})
}
</script>
