<template>
	<div class="sobre-container mt-8">
		<v-container class="d-flex justify-center align-center py-12">
			<v-card class="ios-card-premium pa-10" elevation="12">
				<v-row justify="center" class="mb-8">
					<v-img
						src="@/assets/images/site-images/login/person-img.png"
						alt="Sobre Nós"
						max-width="120"
						class="drop-shadow"
						contain
					></v-img>
				</v-row>

				<v-row>
					<v-col cols="12" class="text-center mb-6">
						<h1 class="text-white text-h3 font-weight-bold tracking-tight">Sobre Nós</h1>
					</v-col>
				</v-row>

				<v-row class="mb-10">
					<v-col cols="12">
						<p class="text-body-1 text-white opacity-90 leading-relaxed text-justify">
							A preservação do conhecimento e o acesso irrestrito à leitura sempre foram pilares fundamentais para o desenvolvimento da sociedade. Esta biblioteca digital foi concebida para honrar essa tradição, unindo o rigor e a organização dos acervos clássicos à eficiência da tecnologia atual. Nosso objetivo é facilitar o acesso à informação de maneira estruturada, confiável e sem complicações.
						</p>
					</v-col>
				</v-row>

				<v-row class="mb-10">
					<v-col cols="12">
						<h2 class="text-cyan text-h5 font-weight-bold mb-4">O Projeto</h2>
						<p class="text-body-1 text-white opacity-90 leading-relaxed text-justify">
							Este sistema é fruto de dedicação acadêmica e técnica, projetado e desenvolvido por Gabriel. A plataforma representa a culminação do Trabalho de Conclusão de Curso (TCC) em Sistemas de Informação pela Universidade Federal da Grande Dourados (UFGD). A iniciativa reflete o compromisso genuíno em aplicar os fundamentos consolidados da computação para resolver problemas reais de gestão de acervos.
						</p>
					</v-col>
				</v-row>

				<v-row class="mb-10">
					<v-col cols="12">
						<h2 class="text-cyan text-h5 font-weight-bold mb-4">Nossa Missão</h2>
						<p class="text-body-1 text-white opacity-90 leading-relaxed text-justify">
							Fornecer uma plataforma estável, segura e direta ao ponto, garantindo que o acervo esteja sempre à disposição. Valorizamos a durabilidade do conhecimento e a importância de um sistema que cumpra o seu papel com excelência, permitindo que os usuários encontrem exatamente o que precisam.
						</p>
					</v-col>
				</v-row>

				<v-divider class="opacity-20 mb-8"></v-divider>

				<v-row justify="center">
					<v-btn class="ios-action-btn px-10" elevation="6" @click="$router.push('/')">
						Voltar ao Início
					</v-btn>
				</v-row>
			</v-card>
		</v-container>
	</div>
</template>

<script>
export default {
	name: 'SobreNosPage'
}
</script>

<style scoped>
	.sobre-container {
		min-height: 100vh;
		background-color: var(--ios-bg);
		padding-bottom: 80px;
	}

	.ios-card-premium {
		width: 100%;
		max-width: 900px;
		background: rgba(58, 99, 145, 0.7) !important;
		backdrop-filter: blur(20px);
		border-radius: 40px !important;
		border: 1px solid rgba(255, 255, 255, 0.1);
		animation: ios-reveal 0.8s var(--spring-easing) forwards;
	}

	.tracking-tight {
		letter-spacing: -2px !important;
	}

	.text-cyan {
		color: var(--ios-cyan) !important;
	}

	.opacity-90 { opacity: 0.9; }
	.opacity-20 { opacity: 0.2; }

	.leading-relaxed {
		line-height: 1.8 !important;
		font-size: 1.1rem !important;
	}

	.ios-action-btn {
		background: linear-gradient(135deg, #00B8D4 0%, #0097A7 100%) !important;
		color: white !important;
		border-radius: 16px !important;
		height: 54px !important;
		text-transform: none !important;
		font-size: 18px !important;
		font-weight: 700 !important;
	}

	.drop-shadow {
		filter: drop-shadow(0 10px 20px rgba(0,0,0,0.4));
	}

	@keyframes ios-reveal {
		from {
			opacity: 0;
			transform: scale(0.95) translateY(40px);
		}
		to {
			opacity: 1;
			transform: scale(1) translateY(0);
		}
	}
</style>
