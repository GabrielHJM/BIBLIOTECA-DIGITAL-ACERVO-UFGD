<template>
	<v-container class="d-flex justify-center">
		<v-card class="ma-7 pa-5" style="width: 75%; height: 700px; background-color: #3a6391;" elevation="0">
			<v-row>
				<v-sheet class="pa-4 text-white text-h4 text-left" elevation="0" style="background-color: transparent; font-weight: bolder;">
					CADASTRE-SE
				</v-sheet>
			</v-row>

			<v-row class="mt-8">
				<v-col cols="6">
					<v-row>
						<v-col cols="12" class="text-left"><label>Nome:</label></v-col>
						<v-col cols="12"><input type="text" class="custom-input w-100" /></v-col>
					</v-row>
				</v-col>
				<v-col cols="6">
					<v-row>
						<v-col cols="12" class="text-left"><label>Instituição de Ensino:</label></v-col>
						<v-col cols="12"><input type="text" class="custom-input w-100" /></v-col>
					</v-row>
				</v-col>
			</v-row>

			<v-row class="mt-8">
				<v-col cols="6">
					<v-row>
						<v-col cols="12" class="text-left"><label>Email:</label></v-col>
						<v-col cols="12"><input type="text" class="custom-input w-100" /></v-col>
					</v-row>
				</v-col>
				<v-col cols="6">
					<v-row>
						<v-col cols="12" class="text-left"><label>Senha:</label></v-col>
						<v-col cols="12"><input type="password" class="custom-input w-100" /></v-col>
					</v-row>
				</v-col>
			</v-row>

			<v-row class="mt-8">
				<v-col cols="6">
					<v-row>
						<v-col cols="12" class="text-left"><label>Eu sou:</label></v-col>
						<v-col cols="12">
							<v-select
								class="custom-input custom-select w-100"
								:items="['Aluno', 'Professor', 'Curador Especializado', 'Suporte']"
								variant="plain"
								hide-details
								density="comfortable"
								v-model="tipoUsuario"
								menu-icon=""
								chips
							></v-select>
						</v-col>
					</v-row>
				</v-col>
				<v-col cols="6">
					<v-row>
						<v-col cols="12"><label></label></v-col>
						<v-col cols="12" class="text-left">
							<v-btn size="large" class="login-button mx-auto" elevation="0">Cadastrar-me</v-btn>
						</v-col>
					</v-row>
				</v-col>
			</v-row>
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
		caminho: require("@/assets/images/site-images/login/person-img.png"),
		tipoUsuario: ''
	}),
	mounted() {
		console.log('livro: ', this.livros[0])
	}
}
</script>

<style>
	label {
		color: white;
		font-weight: bolder;
		padding-left: 10px;
	}

	.border {
		border: 1px solid white;
	}

	.custom-input {
		border: 2px solid black;
		border-radius: 25px;
		padding: 8px 12px;
		font-size: 16px;
		color: #000;
		background-color: white;
		height: 40px;
		box-sizing: border-box;
	}

	.login-button {
		background-color: #00B8D4;
		color: white;
		padding: 10px 24px;
		border-radius: 25px;
		text-transform: none;
		font-size: 16px;
	}

	.login-button:hover {
		background-color: #00B8D4;
	}

	.space-v-row {
		margin-top: 25px;
	}

	.custom-select.v-input {
		padding: 0 !important;
		min-height: 0 !important;
		height: 40px;
	}

	.custom-select.v-input .v-input__control {
		padding: 0 !important;
		height: 100%;
		align-items: center;
	}

	.custom-select.v-input .v-field__input {
		padding: 0 8px !important;
		height: 100%;
		display: flex;
		align-items: center;
	}

	.custom-select.v-input {
		padding: 0 !important;
		min-height: 0 !important;
		height: 40px !important;
	}

	.custom-select.v-input .v-input__control {
		padding: 0 !important;
		height: 100% !important;
		align-items: center;
	}

	.custom-select.v-input .v-field__field {
		height: 100% !important;
		display: flex;
		align-items: center;
	}

	.custom-select.v-input .v-field__input {
		padding: 0 12px !important;
		height: 100% !important;
		display: flex;
		align-items: center;
	}

	.custom-select.v-input .v-field {
		border: none !important;
		box-shadow: none !important;
	}
</style>
