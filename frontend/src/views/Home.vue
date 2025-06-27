<template>
	<div class="mx-auto mt-5">
		<v-container class="pa-0" fluid>
			<v-card min-width="100%" rounded="xl" style="background-color: rgba(36, 36, 36, 0.50);">
				<v-row no-gutters="" class="ma-4">
					<!-- Tela à esquerda, mostrando os matériais sugeridos-->
					<v-col>
						<v-chip
							variant="text"
							size="large"
							label
							class="text-white text-h4 ma-2"
						>
							COMEÇE SEUS ESTUDOS
						</v-chip>
						<div>
							<!-- Session 01 -->
							<v-row no-gutters>
								<!-- Material 01 -->
								<v-col cols="6">
									<div>
										<v-row no-gutters>
											<v-col cols="12" class="custom-title">Sistemas de Inf</v-col>
											<v-col cols="12">
												<!-- Linha que contém duas divs lado a lado -->
												<v-row no-gutters>
													<v-col cols="6">
														<div class="pa-2" style="height: 300px; overflow: hidden; position: relative;">
															<img
																:src="livros[0].link"
																style="position: absolute; top: 0; left: 0; width: 100%; height: 100%; object-fit: contain;"
															/>
														</div>
													</v-col>
													<v-col cols="6">
														<div class="pa-2" style="height: 300px; overflow: auto;">
															<ul style="list-style: none; padding: 0; margin: 0;">
																<li class="custom-item"><strong>Título:</strong> Estrutura de Dados</li>
																<li class="custom-item"><strong>Publicado em:</strong> 05/10/2022</li>
																<li class="custom-item">
																	<strong>Páginas:</strong> 220 |
																	<strong>Capítulos:</strong> 5
																</li>
																<li class="custom-item"><strong>Volume:</strong> 1</li>
																<li class="custom-item"><strong>Edição:</strong> 4</li>
																<li class="custom-item"><strong>Organizadores:</strong> Michael T. Goodrich</li>
															</ul>
														</div>
													</v-col>
												</v-row>
											</v-col>
										</v-row>
									</div>
								</v-col>

								<!-- Material 02 -->
								<v-col cols="6">
									<v-row no-gutters="">
										<v-col cols="12" class="custom-title">Psicologia</v-col>
										<v-col cols="12">
											<!-- Linha que contém duas divs lado a lado -->
											<v-row no-gutters>
												<v-col cols="6">
													<div class="pa-2" style="height: 300px; overflow: hidden; position: relative;">
														<img
															:src="livros[6].link"
															style="position: absolute; top: 0; left: 0; width: 100%; height: 100%; object-fit: contain;"
														/>
													</div>
												</v-col>
												<v-col cols="6">
													<div class="pa-2" style="height: 300px; overflow: auto;">
														<ul style="list-style: none; padding: 0; margin: 0;">
															<li class="custom-item"><strong>Título:</strong> Psicopatologia</li>
															<li class="custom-item"><strong>Publicado em:</strong> 03/05/2018</li>
															<li class="custom-item">
																<strong>Páginas:</strong> 335 |
																<strong>Capítulos:</strong> 12
															</li>
															<li class="custom-item"><strong>Volume:</strong> 1</li>
															<li class="custom-item"><strong>Edição:</strong> 3</li>
															<li class="custom-item"><strong>Organizadores:</strong> Paulo Dalgalarrondo</li>
														</ul>
													</div>
												</v-col>
											</v-row>
										</v-col>
									</v-row>
								</v-col>
							</v-row>

							<!-- Session 02 -->
							<v-row no-gutters="" class="mt-6">
								<v-col cols="6">
									<v-row no-gutters="">
										<v-col cols="12" class="custom-title">Sistemas de Inf</v-col>
										<v-col cols="12">
											<!-- Linha que contém duas divs lado a lado -->
											<v-row no-gutters>
												<v-col cols="6">
													<div class="pa-2" style="height: 300px; overflow: hidden; position: relative;">
														<img
															:src="livros[2].link"
															style="position: absolute; top: 0; left: 0; width: 100%; height: 100%; object-fit: contain;"
														/>
													</div>
												</v-col>
												<v-col cols="6">
													<div class="pa-2" style="height: 300px; overflow: auto;">
														<ul style="list-style: none; padding: 0; margin: 0;">
															<li class="custom-item"><strong>Título:</strong> Meu primeiro livro de programação</li>
															<li class="custom-item"><strong>Publicado em:</strong> 08/06/2024</li>
															<li class="custom-item">
																<strong>Páginas:</strong> 144 |
																<strong>Capítulos:</strong> 1
															</li>
															<li class="custom-item"><strong>Volume:</strong> 1</li>
															<li class="custom-item"><strong>Edição:</strong> 1</li>
															<li class="custom-item"><strong>Organizadores:</strong> Nestor Burlamaqui</li>
														</ul>
													</div>
												</v-col>
											</v-row>
										</v-col>
									</v-row>
								</v-col>
								<v-col cols="6">
									<v-row no-gutters="">
										<v-col cols="12" class="custom-title">Sistemas de Inf</v-col>
										<v-col cols="12">
											<!-- Linha que contém duas divs lado a lado -->
											<v-row no-gutters>
												<v-col cols="6">
													<div class="pa-2" style="height: 300px; overflow: hidden; position: relative;">
														<img
															:src="livros[3].link"
															style="position: absolute; top: 0; left: 0; width: 100%; height: 100%; object-fit: contain;"
														/>
													</div>
												</v-col>
												<v-col cols="6">
													<div class="pa-2" style="height: 300px; overflow: auto;">
														<ul style="list-style: none; padding: 0; margin: 0;">
															<li class="custom-item"><strong>Título:</strong> Introdução á teoria da computação</li>
															<li class="custom-item"><strong>Publicado em:</strong> 24/07/2005</li>
															<li class="custom-item">
																<strong>Páginas:</strong> 448 |
																<strong>Capítulos:</strong> 10
															</li>
															<li class="custom-item"><strong>Volume:</strong> 1</li>
															<li class="custom-item"><strong>Edição:</strong> 2</li>
															<li class="custom-item"><strong>Organizadores:</strong> Michael Sipser</li>
														</ul>
													</div>
												</v-col>
											</v-row>
										</v-col>
									</v-row>
								</v-col>
							</v-row>
						</div>
					</v-col>

					<v-divider
						class="border-opacity-50"
						color="white"
						vertical
					></v-divider>

					<!-- Tela á direita, mostrando as vantegens da Biblioteca Digital -->
					<v-col cols="5">
						<v-sheet class="pa-4 text-white text-h5 text-left" elevation="0" style="background-color: transparent;">
							O QUE TORNA A <strong>BIBLIOTECA DIGITAL ACESSÍVEL</strong> UMA EXCELENTE OPÇÃO?
						</v-sheet>

						<div class="space pt-16">
							<v-row>
								<v-col cls="6">
									<v-row class="text-center">
										<v-col>
											<img
												src="../assets/images/site-images/home/infinite.png"
												style="width: 75px; height: 75px; object-fit: contain;"
											/>
										</v-col>

										<v-col class="text-left mt-6">
											<v-chip
												color="blue"
												variant="outlined"
												label
											>
												Sem Limites
											</v-chip>
										</v-col>
									</v-row>
								</v-col>
								<v-col cls="6">
									<v-row class="text-center">
										<v-col>
											<img
												src="../assets/images/site-images/home/4.png"
												style="width: 75px; height: 75px; object-fit: contain;"
											/>
										</v-col>

										<v-col class="text-left mt-6">
											<v-chip
												color="blue"
												variant="outlined"
												label
											>
												Acesso Instantâneo
											</v-chip>
										</v-col>
									</v-row>
								</v-col>
							</v-row>

							<v-row>
								<v-col cls="6">
									<v-row class="text-center">
										<v-col>
											<img
												src="../assets/images/site-images/home/2.png"
												style="width: 75px; height: 75px; object-fit: contain;"
											/>
										</v-col>

										<v-col class="text-left mt-6">
											<v-chip
												color="blue"
												variant="outlined"
												label
											>
												Personalização
											</v-chip>
										</v-col>
									</v-row>
								</v-col>
								<v-col cls="6">
									<v-row class="text-center">
										<v-col>
											<img
												src="../assets/images/site-images/home/5.png"
												style="width: 75px; height: 75px; object-fit: contain;"
											/>
										</v-col>

										<v-col class="text-left mt-6">
											<v-chip
												color="blue"
												variant="outlined"
												label
											>
												Gratuito
											</v-chip>
										</v-col>
									</v-row>
								</v-col>
							</v-row>

							<v-row>
								<v-col cls="12">
									<v-row class="text-center">
										<v-col>
											<img
												src="../assets/images/site-images/home/3.png"
												style="width: 75px; height: 75px; object-fit: contain;"
											/>
										</v-col>

										<v-col class="text-left mt-6">
											<v-chip
												color="blue"
												variant="outlined"
												label
											>
												Sem Restrição Geográfica
											</v-chip>
										</v-col>
									</v-row>
								</v-col>
							</v-row>

							<v-row>
								<v-col cls="12">
									<p class="welcome-text">Junte-se a nós nesta emocionante jornada de descoberta literária. Com a Biblioteca Digital Livre, o mundo acadêmico está ao seu alcance. Cadastre-se agora e começe a explorar!</p>
								</v-col>
							</v-row>
						</div>
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
	}),
	mounted() {
		console.log('livro: ', this.livros[0])
	}
}
</script>

<style>
	.border {
		border: 1px solid black;
	}

	.space {
		margin-left: 15px;
	}

	ul {
		list-style: none;
		padding: 0;
		margin: 0;
		display: flex;
		flex-direction: column;
		justify-content: space-between;
		height: 100%;
	}

	li {
		text-align: left;
	}

	.custom-title {
		color: white;
		font-size: 20px;
		font-weight: bold;
	}

	.custom-item {
		color: white;
		font-size: 15px;
		line-height: 2.5;
	}

	.welcome-text {
		padding: 50px 20px;
		font-weight: bolder;
		font-size: 1.5rem;
		color: white;
		text-align: justify;
	}
</style>
