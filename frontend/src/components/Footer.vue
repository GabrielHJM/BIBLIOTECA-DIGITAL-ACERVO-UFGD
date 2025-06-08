<template>

  <v-footer class="d-flex flex-column" color="#242424" rounded="lg">
		<div class="d-flex w-100 align-center px-4 py-2">
			<strong>Get connected with us on social networks!</strong>

			<div class="d-flex ga-2 ms-auto">
				<v-btn
				v-for="icon in icons"
				:key="icon"
				:icon="icon"
				size="small"
				variant="plain"
				></v-btn>
			</div>
		</div>

		<div class="px-4 py-2 bg-surface-variant text-center w-100 rounded-lg">
			<span class="text-footer"> Biblioteca Digital 2025 - Versão: <a class="link">{{ pack.version }} </a> | Projeto de Sistemas Computacionais II - UFGD | Copyright 2024 Todos os direitos reservados &copy;</span>
		</div>
  </v-footer>

</template>

<script>

export default {
	name: 'Footer-Vue',
	data: () => ({
		icons: [
			'mdi-facebook',
			'mdi-twitter',
			'mdi-linkedin',
			'mdi-instagram',
		],
		pack: require('../../package.json')
	})
}
</script>
