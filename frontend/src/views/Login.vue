<template>
	<v-container class="d-flex justify-center">
		<v-card class="ma-7 pa-5" style="width: 30%; height: 700px; background-color: #3a6391;" elevation="0">
			<v-row justify="center">
				<v-img
					:src="caminho"
					alt="Logo ou imagem de login"
					max-width="150"
					contain
				></v-img>
			</v-row>

			<v-row>
				<v-col cols="12" class="text-left"><label>Nome:</label></v-col>
				<v-col cols="12"><input type="text" class="custom-input w-100" /></v-col>
			</v-row>

			<v-row justify="center" class="text-left">
				<v-col cols="12"><label>Senha:</label></v-col>
				<v-col cols="12"><input type="password" class="custom-input w-100" /></v-col>
				<v-col cols="12"><a href="#" class="forgot-password-link">Esqueci minha senha</a></v-col>
			</v-row>

			<v-row justify="center" class="mt-4">
				<v-btn class="login-button" elevation="0">Login</v-btn>
			</v-row>

			<v-row class="mt-10 text-left">
				<v-col>
					<a href="#" class="no-account-link">
						<v-icon color="white" size="30">mdi-link-variant</v-icon>
						<span class="no-account"> - Não tenho cadastro</span>
					</a>
				</v-col>
			</v-row>
		</v-card>
	</v-container>
</template>

<script>
import livros from '../../../livros.json'

export default {
	name: 'LoginPage',
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
	.custom-input {
		border: 2px solid black;
		border-radius: 25px;
		padding: 8px 12px;
		font-size: 16px;
		color: #000;
		background-color: white;
	}

	label {
		color: white;
	}

	.forgot-password-link {
		font-size: 12px;
		color: white;
		margin: 0;
		padding-left: 20px;
		text-decoration: none;
	}

	.forgot-password-link:hover {
		text-decoration: underline;
	}

	.login-button {
		background-color: #00B8D4; /* azul-esmeralda */
		color: white;
		padding: 10px 24px;
		border-radius: 25px;
		text-transform: none; /* evita letras maiúsculas automáticas */
		font-size: 16px;
	}

	.login-button:hover {
		background-color: #00B8D4; /* tom mais escuro no hover */
	}

	.no-account {
		color: white;
		font-size: 20px;
	}

	.no-account-link {
		text-decoration: none;
	}

	.no-account-link:hover {
		text-decoration: underline;
		color: white;
	}
</style>
