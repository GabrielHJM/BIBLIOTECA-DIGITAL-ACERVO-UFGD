<template>
	<div class="suporte-root apple-home-root">
		<section class="apple-hero" style="min-height: 50vh;">
			<div class="hero-content-center">
				<h1 class="apple-huge-title apple-gradient-text fade-up-delay-1">Suporte Técnico</h1>
				<h2 class="apple-subtitle fade-up-delay-2 mt-4">
					Estamos aqui para ajudar com qualquer dúvida ou problema no Acervus Core.
				</h2>
			</div>
			<div class="hero-glow"></div>
		</section>

		<section class="section-container" style="max-width: 800px; padding-bottom: 100px;">
			<div class="glass-panel support-panel fade-up-delay-3">
				<div class="d-flex align-center mb-6">
					<v-icon size="40" color="#00E5FF" class="mr-4">mdi-email-fast-outline</v-icon>
					<div>
						<h3 class="support-title">Envie-nos um e-mail</h3>
						<p class="support-desc">Nossa equipe de especialistas responderá o mais rápido possível.</p>
					</div>
				</div>

				<div class="support-details mb-8">
					<div class="detail-item">
						<strong>E-mail Direto:</strong> <span>acervuscore@gmail.com</span>
					</div>
					<div class="detail-item">
						<strong>Horário de Atendimento:</strong> <span>Segunda a Sexta, das 09h às 18h (Horário de Brasília)</span>
					</div>
					<div class="detail-item">
						<strong>Tempo Médio de Resposta:</strong> <span>Até 24 horas úteis</span>
					</div>
				</div>

				<div class="text-center">
					<a href="mailto:acervuscore@gmail.com" class="apple-btn-accent text-decoration-none">
						<v-icon size="20" class="mr-2">mdi-email-outline</v-icon> Abrir E-mail
					</a>
				</div>
			</div>
		</section>
	</div>
</template>

<script>
export default {
	name: 'SuporteView'
}
</script>

<style scoped>
.suporte-root {
	padding-top: 80px;
}
.support-panel {
	padding: 40px;
	border-radius: 24px;
}
.support-title {
	font-size: 1.5rem;
	font-weight: 700;
	margin: 0;
}
.support-desc {
	color: #a1a1a6;
	margin: 0;
	font-size: 1rem;
}
.support-details {
	background: rgba(0, 0, 0, 0.2);
	border-radius: 16px;
	padding: 24px;
	border: 1px solid rgba(255, 255, 255, 0.05);
}
.detail-item {
	display: flex;
	flex-direction: column;
	margin-bottom: 16px;
}
.detail-item:last-child {
	margin-bottom: 0;
}
.detail-item strong {
	color: #00E5FF;
	font-size: 0.9rem;
	margin-bottom: 4px;
}
.detail-item span {
	font-size: 1.05rem;
	color: #f5f5f7;
}

/* Copiado da Home.vue */
.apple-home-root { background-color: transparent; min-height: 100vh; font-family: 'Inter', sans-serif; color: #f5f5f7; }
.apple-hero { position: relative; display: flex; align-items: center; justify-content: center; text-align: center; padding: 40px 20px; overflow: hidden; }
.hero-content-center { position: relative; z-index: 2; max-width: 800px; }
.apple-huge-title { font-size: clamp(2.5rem, 8vw, 4.5rem); font-weight: 800; letter-spacing: -0.04em; line-height: 1.05; }
.apple-subtitle { font-size: clamp(1.1rem, 2vw, 1.5rem); font-weight: 500; color: #a1a1a6; line-height: 1.4; }
.apple-gradient-text { background: linear-gradient(135deg, #00E5FF 0%, #007AFF 100%); -webkit-background-clip: text; -webkit-text-fill-color: transparent; background-clip: text; color: transparent; }
.hero-glow { position: absolute; top: 50%; left: 50%; width: 60vw; height: 60vw; max-width: 800px; max-height: 800px; background: radial-gradient(circle, rgba(10, 132, 255, 0.15) 0%, transparent 60%); transform: translate(-50%, -50%); pointer-events: none; z-index: 1; filter: blur(60px); }
.glass-panel { background: rgba(255, 255, 255, 0.04); backdrop-filter: blur(40px) saturate(200%); -webkit-backdrop-filter: blur(40px) saturate(200%); border: 1px solid rgba(255, 255, 255, 0.08); box-shadow: 0 20px 40px rgba(0, 0, 0, 0.2); }
.apple-btn-accent { background: linear-gradient(135deg, #007AFF 0%, #00E5FF 100%); color: #ffffff; padding: 14px 32px; border-radius: 100px; font-size: 16px; font-weight: 700; border: none; cursor: pointer; transition: all 0.4s cubic-bezier(0.16, 1, 0.3, 1); display: inline-flex; align-items: center; box-shadow: 0 10px 20px rgba(0, 122, 255, 0.3); }
.apple-btn-accent:hover { transform: scale(1.05) translateY(-2px); box-shadow: 0 15px 30px rgba(0, 122, 255, 0.5); }
.fade-up-delay-1 { animation: appleFadeUp 1s forwards; opacity: 0; }
.fade-up-delay-2 { animation: appleFadeUp 1s forwards 0.1s; opacity: 0; }
.fade-up-delay-3 { animation: appleFadeUp 1s forwards 0.2s; opacity: 0; }
@keyframes appleFadeUp { from { opacity: 0; transform: translateY(30px); } to { opacity: 1; transform: translateY(0); } }
.section-container { max-width: 1200px; margin: 0 auto; padding: 0 24px; position: relative; z-index: 2; }
</style>
