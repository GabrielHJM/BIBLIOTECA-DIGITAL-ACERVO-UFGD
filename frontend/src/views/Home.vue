<template>
	<div class="mx-auto border">
		<v-container class="pa-0" fluid="">
			<v-card variant="tonal" min-width="100%">
				<v-row no-gutters="">
					<!-- Tela à esquerda, mostrando os matériais sugeridos-->
					<v-col>
						<v-chip
							variant="text"
							size="large"
							label
						>
							COMEÇE SEUS ESTUDOS
						</v-chip>
						<div>
							<v-row no-gutters="" class="border">
								<v-col cols="6">
									<div>
										<v-row no-gutters="">
											<v-col cols="12">Sistemas de Inf</v-col>
											<v-col cols="12">
												<!-- Div que mostra oconteúdo relacionado a {Sistemas de Inf} -->
												<div>
													<v-row no-gutters="">
														<!-- Imagem -->
														<v-col cols="6" class="border pa-0">
															<v-img
																width="250"
																aspect-ratio="1/1"
																cover
																:src="livros[0].link"
															></v-img>
														</v-col>
														<!-- Descrição -->
														<v-col cols="6" class="border pa-0">
															<v-card
																class="mx-auto"
																max-width="300"
															>
																<v-toolbar color="cyan-lighten-1">
																	<v-icon icon="mdi-book" variant="text"></v-icon>
																	<v-toolbar-title>{{livros[0].nome}}</v-toolbar-title>
																</v-toolbar>

																<v-list density="compact">
																	<v-list-item color="primary">
																		<template v-slot:prepend>
																			<v-icon icon="mdi-account"></v-icon>
																		</template>

																		<v-list-item-title v-text="livros[0].autor"></v-list-item-title>
																	</v-list-item>

																	<v-list-item color="primary">
																		<template v-slot:prepend>
																			<v-icon icon="mdi-format-page-break"></v-icon>
																		</template>

																		<v-list-item-title>Páginas: {{ livros[0].paginas }}</v-list-item-title>
																	</v-list-item>

																	<v-list-item color="primary">
																		<template v-slot:prepend>
																			<v-icon icon="mdi-counter"></v-icon>
																		</template>

																		<v-list-item-title>Capítulos: {{ livros[0].capitulos }}</v-list-item-title>
																	</v-list-item>

																	<v-list-item color="primary">
																		<template v-slot:prepend>
																			<v-icon icon="mdi-numeric"></v-icon>
																		</template>

																		<v-list-item-title>Volume: {{ livros[0].volume }}</v-list-item-title>
																	</v-list-item>

																	<v-list-item color="primary">
																		<template v-slot:prepend>
																			<v-icon icon="mdi-numeric-0-box-multiple-outline"></v-icon>
																		</template>

																		<v-list-item-title>Edição: {{ livros[0].edicao }}</v-list-item-title>
																	</v-list-item>

																	<v-list-item color="primary">
																		<template v-slot:prepend>
																			<v-icon icon="mdi-numeric-0-box-multiple-outline"></v-icon>
																		</template>

																		<v-list-item-title>Edição: {{ livros[0].edicao }}</v-list-item-title>
																	</v-list-item>
																</v-list>
															</v-card>
														</v-col>
													</v-row>
												</div>
											</v-col>
										</v-row>
									</div>
								</v-col>
								<v-col cols="6">
									<v-row no-gutters="">
										<v-col cols="12">Psicologia</v-col>
										<v-col cols="12">Div's</v-col>
									</v-row>
								</v-col>
							</v-row>
							<v-row no-gutters="" class="border">
								<v-col cols="6">
									<v-row no-gutters="">
										<v-col cols="12">História</v-col>
										<v-col cols="12">Div's</v-col>
									</v-row>
								</v-col>
								<v-col cols="6">
									<v-row no-gutters="">
										<v-col cols="12">Engenharia de Energia</v-col>
										<v-col cols="12">Div's</v-col>
									</v-row>
								</v-col>
							</v-row>
							<v-row no-gutters="" class="border">
								<v-col cols="6">
									<v-row no-gutters="">
										<v-col cols="12">Eng. de Computação</v-col>
										<v-col cols="12">Div's</v-col>
									</v-row>
								</v-col>
								<v-col cols="6">
									<v-row no-gutters="">
										<v-col cols="12">Eng. Civil</v-col>
										<v-col cols="12">Div's</v-col>
									</v-row>
								</v-col>
							</v-row>
						</div>
					</v-col>
					<v-divider
						class="border-opacity-50"
						color="black"
						vertical
					></v-divider>
					<!-- Tela á direita, mostrando as vantegens da Biblioteca Digital -->
					<v-col cols="5">
						<v-chip
							variant="text"
							size="large"
							label
						>
							<p>O QUE TORNA A <strong>BIBLIOTECA DIGITAL ACESSÍVEL</strong> UMA EXCELENTE OPÇÃO?</p>
						</v-chip>
					</v-col>
				</v-row>
			</v-card>
		</v-container>
	</div>
</template>

<script>
import livros from '../../../livros.json'

export default {
	name: 'HomePage',
	data: () => ({
		pack: require('../../package.json'),
		livros: livros
	})
}
</script>

<style>
	.border {
		border: 1px solid black;
	}
</style>
