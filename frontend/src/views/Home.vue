<template>
	<div class="mx-auto mt-5">
		<v-container class="pa-0" fluid="">
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
																<li class="custom-item"><strong>Título:</strong> Tempos Digitais</li>
																<li class="custom-item"><strong>Publicado em:</strong> 05/10/2022</li>
																<li class="custom-item">
																	<strong>Páginas:</strong> 220 |
																	<strong>Capítulos:</strong> 5
																</li>
																<li class="custom-item"><strong>Volume:</strong> 1</li>
																<li class="custom-item"><strong>Edição:</strong> 1</li>
																<li class="custom-item"><strong>Organizadores:</strong> Hélio Lemes Costa Jr</li>
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
															:src="livros[0].link"
															style="position: absolute; top: 0; left: 0; width: 100%; height: 100%; object-fit: contain;"
														/>
													</div>
												</v-col>
												<v-col cols="6">
													<div class="pa-2" style="height: 300px; overflow: auto;">
														<ul style="list-style: none; padding: 0; margin: 0;">
															<li class="custom-item"><strong>Título:</strong> Tempos Digitais</li>
															<li class="custom-item"><strong>Publicado em:</strong> 05/10/2022</li>
															<li class="custom-item">
																<strong>Páginas:</strong> 220 |
																<strong>Capítulos:</strong> 5
															</li>
															<li class="custom-item"><strong>Volume:</strong> 1</li>
															<li class="custom-item"><strong>Edição:</strong> 1</li>
															<li class="custom-item"><strong>Organizadores:</strong> Hélio Lemes Costa Jr</li>
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
										<v-col cols="12" class="custom-title">História</v-col>
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
															<li class="custom-item"><strong>Título:</strong> Tempos Digitais</li>
															<li class="custom-item"><strong>Publicado em:</strong> 05/10/2022</li>
															<li class="custom-item">
																<strong>Páginas:</strong> 220 |
																<strong>Capítulos:</strong> 5
															</li>
															<li class="custom-item"><strong>Volume:</strong> 1</li>
															<li class="custom-item"><strong>Edição:</strong> 1</li>
															<li class="custom-item"><strong>Organizadores:</strong> Hélio Lemes Costa Jr</li>
														</ul>
													</div>
												</v-col>
											</v-row>
										</v-col>
									</v-row>
								</v-col>
								<v-col cols="6">
									<v-row no-gutters="">
										<v-col cols="12" class="custom-title">Engenharia de Energia</v-col>
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
															<li class="custom-item"><strong>Título:</strong> Tempos Digitais</li>
															<li class="custom-item"><strong>Publicado em:</strong> 05/10/2022</li>
															<li class="custom-item">
																<strong>Páginas:</strong> 220 |
																<strong>Capítulos:</strong> 5
															</li>
															<li class="custom-item"><strong>Volume:</strong> 1</li>
															<li class="custom-item"><strong>Edição:</strong> 1</li>
															<li class="custom-item"><strong>Organizadores:</strong> Hélio Lemes Costa Jr</li>
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
						color="black"
						vertical
					></v-divider>

					<!-- Tela á direita, mostrando as vantegens da Biblioteca Digital -->
					<v-col cols="5">
						<v-sheet class="pa-4 text-white text-h5 text-left" elevation="0" style="background-color: transparent;">
							O QUE TORNA A <strong>BIBLIOTECA DIGITAL ACESSÍVEL</strong> UMA EXCELENTE OPÇÃO?
						</v-sheet>

						<div class="border space">
							<v-row>
								<v-col cls="6">
									<img
										:src="imgsSide[0]"
										style="width: 50px; height: 50px; object-fit: contain;"
									/>
									Sem Limites
								</v-col>
								<v-col cls="6">Acesso Instantâneo</v-col>
							</v-row>

							<v-row>
								<v-col cls="6">Personalização</v-col>
								<v-col cls="6">Gratuito</v-col>
							</v-row>

							<v-row>
								<v-col cls="12">Sem Restrição Geográfica</v-col>
							</v-row>

							<v-row>
								<v-col cls="12">
									<p>Junte-se a nós nesta emocionante jornada de descoberta literária. Com a Biblioteca Digital Livre, o mundo acadêmico está ao seu alcançe. Cadastre-se agora e começe a explorar!</p>
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
		livros: livros,
		imgsSide: [
			'../assets/images/site-images/home/1.png',
			'../assets/images/site-images/home/2.jpg',
			'../assets/images/site-images/home/3.jpg',
			'../assets/images/site-images/home/4.jpg',
			'../assets/images/site-images/home/5.jpg'
		]
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
</style>
