<template>
	<footer class="premium-footer">
		<div class="footer-container">
			<div class="footer-grid">
				<!-- Brand Bio -->
				<div class="footer-col brand-col">
					<div class="d-flex align-center mb-4">
						<img src="@/assets/images/site-images/login/img-logo-menu-bar.png" alt="Logo" class="footer-logo" />
						<div class="logo-text-stack ml-3">
							<h2 class="footer-title">ACERVUS</h2>
							<div class="footer-subtitle">CORE</div>
						</div>
					</div>
					<p class="footer-bio">
						A sua biblioteca digital inteligente. Democrática, rápida e sem limites. Acesse milhares de recursos educacionais em uma plataforma feita para a sua evolução.
					</p>
				</div>

				<!-- Quick Links -->
				<div class="footer-col">
					<h3 class="col-title">Navegação</h3>
					<ul class="footer-links">
						<li><router-link to="/">Início</router-link></li>
						<li><router-link to="/explorar">Explorar Acervo</router-link></li>
						<li><router-link to="/sobre-nos">Sobre o Projeto</router-link></li>
					</ul>
				</div>

				<!-- Legal / More -->
				<div class="footer-col">
					<h3 class="col-title">Legal & Suporte</h3>
					<ul class="footer-links">
						<li><a href="#">Termos de Uso</a></li>
						<li><a href="#">Política de Privacidade</a></li>
						<li><a href="mailto:contato@acervus.core">Suporte Técnico</a></li>
					</ul>
				</div>

				<!-- Social & Contact -->
				<div class="footer-col">
					<h3 class="col-title">Conecte-se</h3>
					<div class="social-icons">
						<a href="#" class="social-btn"><v-icon size="20">mdi-github</v-icon></a>
						<a href="#" class="social-btn"><v-icon size="20">mdi-linkedin</v-icon></a>
						<a href="#" class="social-btn"><v-icon size="20">mdi-instagram</v-icon></a>
					</div>
					<div class="mt-4 text-caption opacity-60">
						Universidade Federal da Grande Dourados (UFGD)
					</div>
				</div>
			</div>

			<!-- Divider -->
			<div class="footer-divider"></div>

			<!-- Copyright -->
			<div class="footer-bottom">
				<p class="copyright-text">
					&copy; {{ new Date().getFullYear() }} Acervus Core. Todos os direitos reservados.
				</p>
				<p class="author-text">
					Desenvolvido com <v-icon color="pink" size="14" class="mx-1">mdi-heart</v-icon> por Gabriel.
				</p>
			</div>
		</div>
	</footer>
</template>

<script>
export default {
	name: 'AppFooter'
}
</script>

<style scoped>
.premium-footer {
	background: rgba(5, 10, 20, 0.8);
	backdrop-filter: blur(20px);
	-webkit-backdrop-filter: blur(20px);
	border-top: 1px solid rgba(255, 255, 255, 0.05);
	padding: 80px 0 30px;
	color: #ffffff;
	margin-top: 40px;
}

.footer-container {
	max-width: 1200px;
	margin: 0 auto;
	padding: 0 24px;
}

.footer-grid {
	display: grid;
	grid-template-columns: 2fr 1fr 1fr 1fr;
	gap: 40px;
	margin-bottom: 60px;
}

.footer-logo {
	height: 32px;
	filter: drop-shadow(0 2px 10px rgba(0,0,0,0.5));
}

.logo-text-stack {
	display: flex;
	flex-direction: column;
	line-height: 1;
}

.footer-title {
	font-size: 18px;
	font-weight: 900;
	margin: 0;
	letter-spacing: -0.5px;
}

.footer-subtitle {
	font-size: 10px;
	font-weight: 800;
	color: #00E5FF;
	letter-spacing: 2px;
}

.footer-bio {
	font-size: 0.9rem;
	line-height: 1.6;
	color: rgba(255, 255, 255, 0.6);
	max-width: 350px;
}

.col-title {
	font-size: 1.1rem;
	font-weight: 700;
	margin-bottom: 24px;
	letter-spacing: -0.2px;
}

.footer-links {
	list-style: none;
	padding: 0;
	margin: 0;
	display: flex;
	flex-direction: column;
	gap: 12px;
}

.footer-links a {
	color: rgba(255, 255, 255, 0.6);
	text-decoration: none;
	font-size: 0.9rem;
	transition: color 0.3s ease;
	display: inline-block;
}

.footer-links a:hover {
	color: #00E5FF;
}

.social-icons {
	display: flex;
	gap: 12px;
}

.social-btn {
	width: 40px;
	height: 40px;
	border-radius: 50%;
	background: rgba(255, 255, 255, 0.05);
	color: white;
	display: flex;
	align-items: center;
	justify-content: center;
	text-decoration: none;
	transition: all 0.3s var(--spring-easing);
	border: 1px solid rgba(255, 255, 255, 0.1);
}

.social-btn:hover {
	transform: translateY(-4px);
	background: rgba(255, 255, 255, 0.15);
	border-color: #00E5FF;
	color: #00E5FF;
}

.footer-divider {
	height: 1px;
	background: linear-gradient(90deg, transparent, rgba(255,255,255,0.1), transparent);
	margin-bottom: 30px;
}

.footer-bottom {
	display: flex;
	align-items: center;
	justify-content: space-between;
	flex-wrap: wrap;
	gap: 16px;
}

.copyright-text, .author-text {
	font-size: 0.85rem;
	color: rgba(255, 255, 255, 0.5);
	margin: 0;
}

@media (max-width: 900px) {
	.footer-grid {
		grid-template-columns: 1fr 1fr;
	}
	.brand-col {
		grid-column: 1 / -1;
	}
}

@media (max-width: 600px) {
	.footer-grid {
		grid-template-columns: 1fr;
	}
	.footer-bottom {
		flex-direction: column;
		text-align: center;
		justify-content: center;
	}
}
</style>
