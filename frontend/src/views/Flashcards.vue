<template>
	<div class="flashcards-container mt-16 pt-8 pb-16">
		<v-container>
			<div class="d-flex align-center flex-wrap gap-4 mb-10">
				<div>
					<h1 class="text-h3 font-weight-black text-white mb-2">Flashcards</h1>
					<p class="text-subtitle-1 text-white opacity-60">Memorização Prática com Estilo</p>
				</div>
				<v-spacer></v-spacer>
				<v-btn
					color="primary"
					prepend-icon="mdi-plus"
					class="text-none font-weight-bold rounded-pill"
					size="large"
					elevation="4"
					@click="showCreateDialog = true"
				>
					Novo Cartão
				</v-btn>
			</div>

			<div v-if="loading" class="text-center pa-16">
				<CloudLoader text="Carregando cartões..." height="200px" />
			</div>

			<div v-else-if="cards.length === 0" class="text-center pa-16 ios-glass-card reader-rounding border-cyan">
				<v-icon size="80" color="primary" class="mb-4">mdi-cards-outline</v-icon>
				<h2 class="text-h5 text-white">Pronto para começar?</h2>
				<p class="text-body-1 text-white opacity-40 mt-2 mb-8">
					Crie seus próprios cartões ou explore materiais para estudar.
				</p>
				<div class="d-flex justify-center gap-4">
					<v-btn color="primary" variant="tonal" rounded="pill" size="large" @click="showCreateDialog = true">Criar Manualmente</v-btn>
					<v-btn color="primary" rounded="pill" size="large" @click="$router.push('/explorar')">Explorar Livros</v-btn>
				</div>
			</div>

			<div v-else class="flashcard-flow mx-auto" style="max-width: 600px">
				<!-- Progress bar -->
				<div class="d-flex align-center justify-space-between mb-4">
					<span class="text-caption text-white opacity-60">Questão {{ currentIndex + 1 }} de {{ cards.length }}</span>
					<v-progress-linear :model-value="((currentIndex + 1) / cards.length) * 100" color="primary" height="8" rounded></v-progress-linear>
				</div>

				<!-- The Card -->
				<div class="position-relative">
					<v-card
						class="study-card mb-8 reader-rounding"
						:class="{ flipped: isFlipped }"
						@click="flipCard"
						elevation="20"
					>
						<div class="card-inner">
							<!-- Front (Question) -->
							<div class="card-front pa-10 d-flex flex-column align-center justify-center">
								<div class="text-overline text-primary-lighten-3 mb-4">PERGUNTA</div>
								<div class="text-h5 text-center font-weight-bold line-height-relaxed">
									{{ currentCard.pergunta }}
								</div>
								<div class="mt-8 text-caption opacity-40">Clique para revelar a resposta</div>
							</div>

							<!-- Back (Answer) -->
							<div class="card-back pa-10 d-flex flex-column align-center justify-center">
								<div class="text-overline text-white mb-4">RESPOSTA</div>
								<div class="text-h5 text-center font-weight-medium line-height-relaxed">
									{{ currentCard.resposta }}
								</div>
								<v-btn
									icon="mdi-trash-can-outline"
									variant="text"
									color="white"
									class="position-absolute delete-btn"
									@click.stop="confirmDelete(currentCard.id)"
								></v-btn>
							</div>
						</div>
					</v-card>
				</div>

				<!-- Controls -->
				<div v-if="isFlipped" class="d-flex flex-column align-center animate-up">
					<p class="text-subtitle-2 text-white opacity-60 mb-5">Qual foi a sua facilidade?</p>
					<div class="d-flex gap-4">
						<v-btn color="red-darken-2" class="px-8 flex-1" rounded="pill" size="large" @click="revisar(0)">Difícil</v-btn>
						<v-btn color="#007AFF" class="px-8 flex-1 text-white font-weight-bold" rounded="pill" size="large" @click="revisar(1)">Fácil</v-btn>
					</div>
				</div>
			</div>
		</v-container>

		<!-- Create Dialog -->
		<v-dialog v-model="showCreateDialog" max-width="500">
			<v-card class="ios-glass-card border-cyan pa-6" rounded="xl">
				<h2 class="text-h5 font-weight-bold text-white mb-6">Novo Flashcard</h2>
				<v-textarea
					v-model="newCard.pergunta"
					label="Pergunta"
					variant="outlined"
					color="primary"
					rows="3"
					class="mb-4"
					rounded="lg"
				></v-textarea>
				<v-textarea
					v-model="newCard.resposta"
					label="Resposta"
					variant="outlined"
					color="primary"
					rows="3"
					class="mb-6"
					rounded="lg"
				></v-textarea>
				<div class="d-flex gap-3">
					<v-btn variant="text" color="white" block class="flex-1" @click="showCreateDialog = false">Cancelar</v-btn>
					<v-btn color="primary" block class="flex-1 font-weight-bold" @click="salvarNovoCard" :loading="saving">Salvar</v-btn>
				</div>
			</v-card>
		</v-dialog>

		<v-snackbar v-model="snackbar" :timeout="3000" color="primary" rounded="pill">
			{{ snackbarText }}
		</v-snackbar>
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
		loading: true,
		showCreateDialog: false,
		saving: false,
		newCard: {
			pergunta: '',
			resposta: ''
		},
		snackbar: false,
		snackbarText: ''
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
					this.snackbarText = "Sessão finalizada! Bom trabalho.";
					this.snackbar = true;
					setTimeout(() => this.$router.push('/dashboard'), 1500);
				}
			} catch (error) {
				console.error('Erro ao revisar flashcard:', error);
			}
		},
		async salvarNovoCard() {
			if (!this.newCard.pergunta || !this.newCard.resposta) return;
			this.saving = true;
			try {
				const cardData = {
					...this.newCard,
					usuario_id: auth.getUser().id,
					material_id: 0 // Manual card
				};
				await EstudoService.criarFlashcard(cardData);
				this.showCreateDialog = false;
				this.newCard = { pergunta: '', resposta: '' };
				this.snackbarText = "Flashcard criado com sucesso!";
				this.snackbar = true;
				this.carregarCards();
			} catch (error) {
				console.error('Erro ao criar flashcard:', error);
			} finally {
				this.saving = false;
			}
		},
		async confirmDelete(id) {
			if (!confirm("Deseja realmente excluir este flashcard?")) return;
			try {
				await EstudoService.deletarFlashcard(id, auth.getUser().id);
				this.snackbarText = "Flashcard removido.";
				this.snackbar = true;
				this.carregarCards();
			} catch (error) {
				console.error('Erro ao deletar:', error);
			}
		},
		async carregarCards() {
			this.loading = true;
			try {
				const response = await EstudoService.listarFlashcards(auth.getUser().id);
				// Filter for review or show all if empty
				this.cards = response.data;
				this.currentIndex = 0;
				this.isFlipped = false;
			} catch (error) {
				console.error('Erro ao carregar flashcards:', error);
			} finally {
				this.loading = false;
			}
		}
	},
	async mounted() {
		this.carregarCards();
	}
}
</script>

<style scoped>
	.flashcards-container {
		min-height: 100vh;
		background: var(--ios-bg);
	}

	.ios-glass-card {
		background: var(--glass-bg);
		backdrop-filter: var(--glass-blur);
		border: 1px solid var(--glass-border);
	}

	.border-cyan {
		border: 1px solid rgba(0, 122, 255, 0.2) !important;
	}

	.study-card {
		height: 400px;
		perspective: 1200px;
		background: transparent !important;
		cursor: pointer;
	}

	.card-inner {
		position: relative;
		width: 100%;
		height: 100%;
		text-align: center;
		transition: transform 0.7s cubic-bezier(0.4, 0, 0.2, 1);
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
		border-radius: 40px;
		background: rgba(255, 255, 255, 0.05);
		border: 1px solid rgba(255, 255, 255, 0.1);
		color: white;
		box-shadow: 0 10px 30px rgba(0,0,0,0.5);
	}

	.card-back {
		transform: rotateY(180deg);
		background: var(--premium-gradient);
		border: 1px solid rgba(0, 122, 255, 0.3);
	}

	.delete-btn {
		top: 20px;
		right: 20px;
		opacity: 0.6;
	}

	.delete-btn:hover {
		opacity: 1;
		background: rgba(244, 67, 54, 0.2);
	}

	.reader-rounding { border-radius: 40px !important; }
	.line-height-relaxed { line-height: 1.6; }
	.gap-3 { gap: 12px; }
	.gap-4 { gap: 16px; }

	.animate-up {
		animation: slideUp 0.5s cubic-bezier(0.4, 0, 0.2, 1);
	}

	@keyframes slideUp {
		from { opacity: 0; transform: translateY(20px); }
		to { opacity: 1; transform: translateY(0); }
	}

	@media (max-width: 600px) {
		.study-card { height: 350px; }
		.text-h3 { font-size: 2.2rem !important; }
	}
</style>
