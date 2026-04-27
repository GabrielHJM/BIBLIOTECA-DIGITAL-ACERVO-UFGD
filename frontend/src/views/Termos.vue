<template>
	<div class="termos-root apple-home-root">
		<section class="apple-hero" style="min-height: 40vh;">
			<div class="hero-content-center">
				<h1 class="apple-huge-title apple-gradient-text fade-up-delay-1">Legal e Privacidade</h1>
				<h2 class="apple-subtitle fade-up-delay-2 mt-4">
					Transparência e segurança para todos os usuários.
				</h2>
			</div>
			<div class="hero-glow"></div>
		</section>

		<section class="section-container" style="max-width: 900px; padding-bottom: 100px;">
			<div class="glass-panel terms-panel fade-up-delay-3">
				<h3 class="terms-title mb-4">Termos de Uso</h3>
				<p class="terms-text">
					Bem-vindo ao Acervus Core. Ao acessar e utilizar nossa plataforma, você concorda com os seguintes termos.
					O Acervus Core é um projeto educacional voltado à democratização do conhecimento. Os materiais disponibilizados
					são de natureza acadêmica e de domínio público ou indexados por fontes abertas.
				</p>
				<p class="terms-text mb-8">
					O uso indevido da plataforma, tentativas de engenharia reversa, uso de automação abusiva (crawlers) ou 
					qualquer atividade que prejudique a estabilidade do sistema resultará no banimento permanente da conta.
				</p>

				<v-divider class="my-8" style="border-color: rgba(255,255,255,0.1);"></v-divider>

				<h3 class="terms-title mb-4">Política de Privacidade</h3>
				<p class="terms-text">
					Sua privacidade é nossa prioridade. Nós coletamos apenas as informações mínimas e necessárias para o funcionamento
					adequado de sua conta (Nome, Email, e Preferências de Leitura). 
				</p>
				<p class="terms-text">
					<strong>Armazenamento de Dados:</strong> Suas senhas são rigorosamente criptografadas. Nós nunca venderemos, 
					alugaremos ou distribuiremos seus dados pessoais para terceiros para fins publicitários.
				</p>
				<p class="terms-text mb-8">
					<strong>Uso de Cookies:</strong> Utilizamos apenas os cookies estritamente necessários (tokens de sessão) para 
					manter você conectado de forma segura.
				</p>
				
				<div class="text-center mt-12">
					<p class="text-caption opacity-60">Última atualização: Abril de 2026</p>
				</div>
			</div>
		</section>
	</div>
</template>

<script>
export default {
	name: 'TermosView'
}
</script>

<style scoped>
.termos-root { padding-top: 80px; }
.terms-panel { padding: 60px 40px; border-radius: 24px; text-align: left; }
.terms-title { font-size: 1.8rem; font-weight: 700; color: #00E5FF; }
.terms-text { font-size: 1rem; line-height: 1.7; color: rgba(255,255,255,0.8); margin-bottom: 16px; }

/* Global overrides das classes Apple para este arquivo local */
.apple-home-root { background-color: transparent; min-height: 100vh; font-family: 'Inter', sans-serif; color: #f5f5f7; }
.apple-hero { position: relative; display: flex; align-items: center; justify-content: center; text-align: center; padding: 40px 20px; overflow: hidden; }
.hero-content-center { position: relative; z-index: 2; max-width: 800px; }
.apple-huge-title { font-size: clamp(2.5rem, 8vw, 4.5rem); font-weight: 800; letter-spacing: -0.04em; line-height: 1.05; }
.apple-subtitle { font-size: clamp(1.1rem, 2vw, 1.5rem); font-weight: 500; color: #a1a1a6; line-height: 1.4; }
.apple-gradient-text { background: linear-gradient(135deg, #00E5FF 0%, #007AFF 100%); -webkit-background-clip: text; -webkit-text-fill-color: transparent; background-clip: text; color: transparent; }
.hero-glow { position: absolute; top: 50%; left: 50%; width: 60vw; height: 60vw; max-width: 800px; max-height: 800px; background: radial-gradient(circle, rgba(10, 132, 255, 0.15) 0%, transparent 60%); transform: translate(-50%, -50%); pointer-events: none; z-index: 1; filter: blur(60px); }
.glass-panel { background: rgba(255, 255, 255, 0.04); backdrop-filter: blur(40px) saturate(200%); -webkit-backdrop-filter: blur(40px) saturate(200%); border: 1px solid rgba(255, 255, 255, 0.08); box-shadow: 0 20px 40px rgba(0, 0, 0, 0.2); }
.fade-up-delay-1 { animation: appleFadeUp 1s forwards; opacity: 0; }
.fade-up-delay-2 { animation: appleFadeUp 1s forwards 0.1s; opacity: 0; }
.fade-up-delay-3 { animation: appleFadeUp 1s forwards 0.2s; opacity: 0; }
@keyframes appleFadeUp { from { opacity: 0; transform: translateY(30px); } to { opacity: 1; transform: translateY(0); } }
.section-container { max-width: 1200px; margin: 0 auto; padding: 0 24px; position: relative; z-index: 2; }
</style>
