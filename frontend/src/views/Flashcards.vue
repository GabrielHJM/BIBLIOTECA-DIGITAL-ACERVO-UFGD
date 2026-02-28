<template>
	<div class="flashcards-container mt-16 pt-8 pb-16">
		<v-container>
			<div class="text-center mb-10">
				<h1 class="text-h3 font-weight-black text-white mb-2">Modo Flashcards</h1>
				<p class="text-subtitle-1 text-white opacity-60">Memorização Prática com Repetição Espaçada</p>
			</div>

			<div v-if="loading" class="text-center pa-16">
				<v-progress-circular indeterminate color="amber" size="64"></v-progress-circular>
			</div>

			<div v-else-if="cards.length === 0" class="text-center pa-16 ios-glass-card rounded-xl">
				<v-icon size="80" color="amber" class="mb-4">mdi-cards-outline</v-icon>
				<h2 class="text-h5 text-white">Nenhum flashcard disponível</h2>
				<p class="text-body-1 text-white opacity-40 mt-2 mb-8">
					Use o botão "Gerar Flashcards IA" na tela de qualquer livro para começar.
				</p>
				<v-btn color="amber" rounded="pill" size="large" @click="$router.push('/explorar')">Explorar Livros</v-btn>
			</div>

			<div v-else class="flashcard-flow mx-auto" style="max-width: 600px">
				<!-- Progress bar -->
				<div class="d-flex align-center justify-space-between mb-4">
					<span class="text-caption text-white opacity-60">Progresso: {{ currentIndex + 1 }} de {{ cards.length }}</span>
					<v-progress-linear :model-value="((currentIndex + 1) / cards.length) * 100" color="amber" height="6" rounded></v-progress-linear>
				</div>

				<!-- The Card -->
				<v-card
					class="study-card mb-8"
					:class="{ flipped: isFlipped }"
					@click="flipCard"
					elevation="12"
					rounded="xl"
				>
					<div class="card-inner">
						<!-- Front (Question) -->
						<div class="card-front pa-10 d-flex flex-column align-center justify-center">
							<div class="text-overline text-amber mb-4">PERGUNTA</div>
							<div class="text-h5 text-center font-weight-bold line-height-relaxed">
								{{ currentCard.pergunta }}
							</div>
							<div class="mt-8 text-caption opacity-40">Clique para revelar a resposta</div>
						</div>

						<!-- Back (Answer) -->
						<div class="card-back pa-10 d-flex flex-column align-center justify-center">
							<div class="text-overline text-cyan mb-4">RESPOSTA</div>
							<div class="text-h5 text-center font-weight-medium line-height-relaxed">
								{{ currentCard.resposta }}
							</div>
						</div>
					</div>
				</v-card>

				<!-- Controls -->
				<div v-if="isFlipped" class="d-flex flex-column align-center fade-in">
					<p class="text-subtitle-2 text-white opacity-60 mb-4">Como foi sua performance?</p>
					<div class="d-flex gap-4">
						<v-btn color="red-darken-2" class="px-6" rounded="lg" @click="revisar(0)">Difícil (Hoje)</v-btn>
						<v-btn color="cyan-darken-2" class="px-6" rounded="lg" @click="revisar(1)">Fácil (7d)</v-btn>
					</div>
				</div>
			</div>
		</v-container>
	</div>
</template>

<script>
import EstudoService from '@/services/EstudoService';
import auth from '@/auth';

export default {
	name: 'FlashcardsView',
	data: () => ({
		cards: [],
		currentIndex: 0,
		isFlipped: false,
		loading: true
	}),
	computed: {
		currentCard() {
			return this.cards[this.currentIndex] || {};
		}
	},
	methods: {
		flipCard() {
			this.isFlipped = !this.isFlipped;
		},
		async revisar(dificuldade) {
			try {
				await EstudoService.revisarFlashcard(this.currentCard.id, dificuldade);
				if (this.currentIndex < this.cards.length - 1) {
					this.currentIndex++;
					this.isFlipped = false;
				} else {
					// End of session
					alert("Sessão finalizada! Bom trabalho.");
					this.$router.push('/dashboard');
				}
			} catch (error) {
				console.error('Erro ao revisar flashcard:', error);
			}
		}
	},
	async mounted() {
		try {
			const response = await EstudoService.listarFlashcards(auth.getUser().id);
			this.cards = response.data.filter(c => new Date(c.proxima_revisao) <= new Date());
			if (this.cards.length === 0 && response.data.length > 0) {
				// If no cards due today, show all for practice
				this.cards = response.data;
			}
		} catch (error) {
			console.error('Erro ao carregar flashcards:', error);
		} finally {
			this.loading = false;
		}
	}
}
</script>

<style scoped>
	.flashcards-container { min-height: 100vh; background: #0f172a; }

	.ios-glass-card {
		background: rgba(255, 255, 255, 0.03);
		backdrop-filter: blur(20px);
		border: 1px solid rgba(255, 255, 255, 0.05);
	}

	.study-card {
		height: 400px;
		perspective: 1000px;
		background: transparent !important;
		cursor: pointer;
	}

	.card-inner {
		position: relative;
		width: 100%;
		height: 100%;
		text-align: center;
		transition: transform 0.6s cubic-bezier(0.4, 0, 0.2, 1);
		transform-style: preserve-3d;
	}

	.study-card.flipped .card-inner {
		transform: rotateY(180deg);
	}

	.card-front, .card-back {
		position: absolute;
		width: 100%;
		height: 100%;
		-webkit-backface-visibility: hidden;
		backface-visibility: hidden;
		border-radius: 24px;
		background: rgba(255, 255, 255, 0.05);
		border: 1px solid rgba(255, 255, 255, 0.1);
		color: white;
	}

	.card-back {
		transform: rotateY(180deg);
		background: linear-gradient(135deg, rgba(0, 184, 212, 0.1), rgba(0, 0, 0, 0.8));
	}

	.line-height-relaxed { line-height: 1.6; }
	.gap-4 { gap: 16px; }

	.fade-in {
		animation: fadeIn 0.4s ease-out;
	}

	@keyframes fadeIn {
		from { opacity: 0; transform: translateY(5px); }
		to { opacity: 1; transform: translateY(0); }
	}
</style>
