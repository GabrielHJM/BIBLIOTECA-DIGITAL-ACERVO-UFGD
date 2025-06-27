<template>
	<div class="mx-auto mt-5">
		<v-container class="pa-0 border" fluid>
			<v-card min-width="100%" rounded="xl" style="background-color: rgba(36, 36, 36, 0.50);">
				<v-row no-gutters="" class="ma-4">
					<v-col cols="6">
						<v-chip
							variant="text"
							size="large"
							label
							class="text-white text-h4 ma-2"
						>
							Recentes
						</v-chip>
					</v-col>

					<v-col cols="6">
						<v-chip
							variant="text"
							label
							class="text-white text-h6 ma-2"
						>
							<p>Filtrar</p>
							<v-icon class="ml-2">mdi-filter</v-icon>
						</v-chip>
					</v-col>
				</v-row>

				<v-row no-gutters="" class="ma-4">

				</v-row>

				<v-row no-gutters="" class="ma-4">

				</v-row>
			</v-card>
		</v-container>
	</div>
</template>

<script>
import livros from '../../../livros.json'

export default {
	name: 'DashboardPage',
	data: () => ({
		pack: require('../../package.json'),
		livros: livros
	}),
	mounted() {
		console.log('livro: ', this.livros[0])
	}
}
</script>

<style>
	.border {
		border: 1px solid black;
	}
</style>
