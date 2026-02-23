<template>
	<v-footer class="ios-floating-footer-wrapper">
		<div class="footer-hover-trigger"></div>
		<div class="ios-floating-footer pa-4">
			<div class="footer-links-container">
				<v-btn
					v-for="link in links"
					:key="link.title"
					:to="link.link"
					variant="text"
					class="ios-footer-link"
					rounded="xl"
				>
					{{ link.title }}
				</v-btn>
			</div>

			<div class="footer-bottom-info mt-2">
				<span class="copyright">Biblioteca Digital (Projeto de Sistemas de Informação SI - UFGD) | Copyright {{ new Date().getFullYear() }} - Todos os Direitos Reservados ©</span>
			</div>
		</div>
	</v-footer>
</template>

<script>
export default {
	name: 'Footer-Vue',
	data: () => ({
		links: [
			{title: "Início", link: "/"},
			{title: "Sobre Nós", link: "/sobre-nos"},
			{title: "Login", link: "/login"}
		],
		pack: require('../../package.json')
	})
}
</script>

<style scoped>
	.ios-floating-footer-wrapper {
		background: transparent !important;
		position: fixed;
		bottom: 0;
		left: 0;
		right: 0;
		width: 100%;
		height: 40px; /* Área de detecção reduzida para não bloquear links */
		display: flex;
		justify-content: center;
		align-items: flex-end;
		pointer-events: none; /* Não bloqueia cliques nos elementos de baixo */
		z-index: 1000;
		padding-bottom: 20px;
	}

	.ios-floating-footer-wrapper * {
		pointer-events: auto; /* Reativa cliques para os botões internos */
	}

	.ios-floating-footer {
		background: rgba(30, 30, 30, 0.6) !important;
		backdrop-filter: blur(20px) saturate(180%);
		-webkit-backdrop-filter: blur(20px) saturate(180%);
		border: 1px solid rgba(255, 255, 255, 0.1);
		border-radius: 24px;
		width: auto;
		min-width: 400px;
		max-width: 90%;
		box-shadow: 0 10px 30px rgba(0, 0, 0, 0.5);
		display: flex;
		flex-direction: column;
		align-items: center;

		/* Estado Inicial: Escondido */
		opacity: 0;
		transform: translateY(100px) scale(0.95);
		transition: all 0.5s var(--spring-easing);
		pointer-events: none;
	}

	.ios-floating-footer-wrapper:hover .ios-floating-footer,
	.ios-floating-footer-wrapper:active .ios-floating-footer,
	.ios-floating-footer:hover {
		opacity: 1 !important;
		transform: translateY(0) scale(1) !important;
		pointer-events: auto !important;
	}

	.footer-hover-trigger {
		position: absolute;
		bottom: 0;
		left: 0;
		right: 0;
		height: 60px; /* Sensibilidade do gatilho */
		pointer-events: auto;
		z-index: -1;
	}

	.footer-links-container {
		display: flex;
		gap: 8px;
		flex-wrap: wrap;
		justify-content: center;
	}

	.ios-footer-link {
		color: rgba(255, 255, 255, 0.8) !important;
		text-transform: none !important;
		font-weight: 600 !important;
		letter-spacing: 0.5px;
		transition: all 0.3s ease;
	}

	.ios-footer-link:hover {
		color: var(--ios-cyan) !important;
		background: rgba(255, 255, 255, 0.05) !important;
		transform: translateY(-2px);
	}

	.footer-bottom-info {
		display: flex;
		align-items: center;
		font-size: 11px;
		color: rgba(255, 255, 255, 0.4);
		text-transform: uppercase;
		letter-spacing: 1px;
	}

	.opacity-20 { opacity: 0.2; }

	@keyframes footer-slide-up {
		from {
			opacity: 0;
			transform: translateY(40px) scale(0.95);
		}
		to {
			opacity: 1;
			transform: translateY(0) scale(1);
		}
	}

	/* Responsive tweaks */
	@media (max-width: 600px) {
		.ios-floating-footer {
			min-width: 300px;
			padding: 12px !important;
			opacity: 1 !important; /* Force visible on mobile */
			transform: translateY(0) scale(1) !important;
			pointer-events: auto !important;
			bottom: 10px;
		}
		.ios-footer-link {
			font-size: 13px;
		}
		.ios-floating-footer-wrapper {
			height: auto;
			padding-bottom: 10px;
		}
	}
</style>
