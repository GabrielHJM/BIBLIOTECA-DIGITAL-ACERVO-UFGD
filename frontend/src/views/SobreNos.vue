<template>
	<div class="sobre-container">
		<v-container class="py-12 px-4 px-md-8">
			<!-- Header Section -->
			<v-row justify="center" class="mb-12">
				<v-col cols="12" md="10" class="text-center">
					<div class="header-reveal">
						<h1 class="text-premium-gradient text-h2 font-weight-black mb-4">Conheça Nossa História</h1>
						<p class="text-subtitle-1 text-white opacity-70 max-width-600 mx-auto">
							Unindo a tradição do conhecimento com o futuro da tecnologia digital.
						</p>
					</div>
				</v-col>
			</v-row>

			<!-- Main Content Card -->
			<v-row justify="center">
				<v-col cols="12" md="10">
					<v-card class="ios-glass-card pa-8 pa-md-12" flat color="transparent">
						<!-- Introdução -->
						<section class="mb-12 content-reveal-1">
							<h2 class="text-accent-cyan text-h4 font-weight-bold mb-6 d-flex align-center">
								<v-icon color="var(--ios-cyan)" class="mr-3">mdi-book-open-variant</v-icon>
								Preservando o Legado
							</h2>
							<p class="premium-body-text">
								A preservação do conhecimento e o acesso irrestrito à leitura sempre foram pilares fundamentais para o desenvolvimento da sociedade. Esta biblioteca digital foi concebida para honrar essa tradição, unindo o rigor e a organização dos acervos clássicos à eficiência da tecnologia atual. Nosso objetivo é facilitar o acesso à informação de maneira estruturada, confiável e sem complicações.
							</p>
						</section>

						<v-divider class="glass-divider mb-12"></v-divider>

						<!-- Missão, Visão e Valores Grid -->
						<v-row class="mb-12 content-reveal-2">
							<v-col cols="12" sm="4">
								<div class="value-card">
									<v-icon size="40" class="mb-4">mdi-target</v-icon>
									<h3 class="text-h5 font-weight-bold mb-2">Missão</h3>
									<p class="text-body-2 opacity-80">Garantir acesso universal e estável ao conhecimento através de uma plataforma de alta performance.</p>
								</div>
							</v-col>
							<v-col cols="12" sm="4">
								<div class="value-card">
									<v-icon size="40" class="mb-4">mdi-eye-outline</v-icon>
									<h3 class="text-h5 font-weight-bold mb-2">Visão</h3>
									<p class="text-body-2 opacity-80">Tornar-se a referência digital em gestão de acervos acadêmicos e literários.</p>
								</div>
							</v-col>
							<v-col cols="12" sm="4">
								<div class="value-card border-none">
									<v-icon size="40" class="mb-4">mdi-shield-check-outline</v-icon>
									<h3 class="text-h5 font-weight-bold mb-2">Valores</h3>
									<p class="text-body-2 opacity-80">Excelência técnica, segurança de dados e compromisso com a educação gratuita.</p>
								</div>
							</v-col>
						</v-row>

						<!-- O Projeto -->
						<v-row class="mb-12 content-reveal-3">
							<v-col cols="12">
								<h2 class="text-accent-cyan text-h4 font-weight-bold mb-6">O Projeto</h2>
								<p class="premium-body-text">
									Este sistema é fruto de dedicação acadêmica e técnica, projetado e desenvolvido por <strong>Gabriel Henrique Jesus de Melo</strong>. A plataforma representa a culminação do Trabalho de Conclusão de Curso (TCC) em Sistemas de Informação pela <strong>Universidade Federal da Grande Dourados (UFGD)</strong>.
								</p>
							</v-col>
						</v-row>

						<!-- Rodapé do Card -->
						<v-divider class="glass-divider mb-12"></v-divider>

						<v-row justify="center" align="center" class="content-reveal-4">
							<v-col cols="12" sm="auto">
								<div class="d-flex align-center justify-center">
									<div class="text-center">
										<div class="text-subtitle-1 font-weight-bold text-white">Gabriel Henrique Jesus de Melo</div>
										<div class="text-caption text-white opacity-60">Desenvolvedor & Criador</div>
									</div>
								</div>
							</v-col>
						</v-row>

						<v-row justify="center" class="mt-12 content-reveal-4">
							<v-btn class="ios-premium-btn" elevation="0" @click="$router.push('/')">
								<v-icon left class="mr-2">mdi-arrow-left</v-icon>
								Voltar ao Início
							</v-btn>
						</v-row>
					</v-card>
				</v-col>
			</v-row>
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
		background: radial-gradient(circle at top right, rgba(0, 184, 212, 0.1), transparent 40%),
					radial-gradient(circle at bottom left, rgba(0, 122, 255, 0.1), transparent 40%);
		background-color: transparent;
		padding-top: 2rem;
	}

	.text-premium-gradient {
		background: linear-gradient(135deg, #FFFFFF 0%, rgba(255,255,255,0.7) 100%);
		-webkit-background-clip: text;
		background-clip: text;
		-webkit-text-fill-color: transparent;
		letter-spacing: -1.5px !important;
	}

	.max-width-600 {
		max-width: 600px;
	}

	.text-accent-cyan {
		color: var(--ios-cyan) !important;
		letter-spacing: -0.5px;
	}

	.ios-glass-card {
		background: rgba(var(--v-theme-surface), 0.4) !important;
		backdrop-filter: blur(30px) saturate(160%) !important;
		-webkit-backdrop-filter: blur(30px) saturate(160%) !important;
		border-radius: 48px !important;
		border: 1px solid rgba(255, 255, 255, 0.1) !important;
		box-shadow: 0 32px 64px rgba(0, 0, 0, 0.3) !important;
		overflow: hidden;
	}

	.premium-body-text {
		color: white;
		opacity: 0.85;
		font-size: 1.15rem;
		line-height: 1.8;
		text-align: justify;
		font-weight: 400;
	}

	.glass-divider {
		border-color: rgba(255, 255, 255, 0.1) !important;
	}

	.value-card {
		padding: 24px;
		text-align: center;
		color: white;
		border-right: 1px solid rgba(255, 255, 255, 0.1);
		transition: transform 0.3s ease;
	}

	.value-card:hover {
		transform: translateY(-5px);
	}

	@media (max-width: 600px) {
		.value-card {
			border-right: none;
			border-bottom: 1px solid rgba(255, 255, 255, 0.1);
		}
	}

	.border-none {
		border: none !important;
	}

	.tech-chip {
		background: rgba(0, 184, 212, 0.1) !important;
		border: 1px solid rgba(0, 184, 212, 0.2) !important;
		font-weight: 600;
	}

	.border-premium {
		border: 2px solid rgba(0, 184, 212, 0.4) !important;
	}

	.ios-premium-btn {
		background: rgba(255, 255, 255, 0.1) !important;
		backdrop-filter: blur(10px);
		color: white !important;
		border: 1px solid rgba(255, 255, 255, 0.2) !important;
		border-radius: 20px !important;
		padding: 0 32px !important;
		height: 56px !important;
		text-transform: none !important;
		font-size: 1.1rem !important;
		font-weight: 600 !important;
		transition: all 0.3s cubic-bezier(0.16, 1, 0.3, 1) !important;
	}

	.ios-premium-btn:hover {
		background: white !important;
		color: black !important;
		transform: scale(1.05);
	}

	.gap-4 {
		gap: 1rem;
	}

	/* Animations */
	.header-reveal {
		animation: reveal 1s var(--spring-easing) forwards;
	}

	.content-reveal-1 { animation: reveal 1s var(--spring-easing) 0.2s forwards; opacity: 0; }
	.content-reveal-2 { animation: reveal 1s var(--spring-easing) 0.4s forwards; opacity: 0; }
	.content-reveal-3 { animation: reveal 1s var(--spring-easing) 0.6s forwards; opacity: 0; }
	.content-reveal-4 { animation: reveal 1s var(--spring-easing) 0.8s forwards; opacity: 0; }

	.float-animation {
		animation: float 6s ease-in-out infinite;
	}

	@keyframes float {
		0%, 100% { transform: translateY(0); }
		50% { transform: translateY(-15px); }
	}

	@keyframes reveal {
		from {
			opacity: 0;
			transform: translateY(30px);
		}
		to {
			opacity: 1;
			transform: translateY(0);
		}
	}

	.drop-shadow {
		filter: drop-shadow(0 20px 30px rgba(0,0,0,0.5));
	}

	.profile-circle {
		border-radius: 50% !important;
		aspect-ratio: 1/1;
		object-fit: cover;
		border: 4px solid rgba(0, 184, 212, 0.4);
	}
</style>
