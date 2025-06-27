<template>
	<div class="mx-auto mt-8">
		<v-container class="pa-0" fluid>
			<v-card min-width="100%" rounded="xl" style="background-color: rgba(36, 36, 36, 0.50);">
				<v-row dense>
					<!-- div a esquerda com as informações do conteúdo renderizado -->
					<v-col cols="4">
						<div>
							<!-- Titulo -->
							<div>
								<v-row dense class="mt-4">
									<v-col cols="2"><v-icon color="white" class="pt-4">mdi-file-document-multiple</v-icon></v-col>
									<v-col cols="10" class="text-left">
										<h3 class="titulo">Estrutura de Dados e Algoritmos</h3>
									</v-col>
								</v-row>
							</div>

							<!-- div com as informações do livro -->
							<div class="mt-8">
								<v-row dense>
									<v-col cols="4">
										<div class="ml-4" style="height: 300px; overflow: hidden; position: relative;">
											<img
												:src="livros[0].link"
												style="position: absolute; top: 0; left: 0; width: 100%; height: 100%; object-fit: contain;"
											/>
										</div>
									</v-col>

									<v-col cols="8">
										<div class="pa-2" style="height: 300px; overflow: auto;">
											<ul style="list-style: none; padding: 0; margin: 0;">
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

								<!-- Avaliar -->
								<v-row no-gutters class="mt-4">
									<v-col cols="12">
										<div class="ml-8">
											<v-row>
												<v-col cols="4" align-self="center">
													<p class="rating-text text-left">Avaliar: </p>
												</v-col>
												<v-col cols="8" align-self="center">
													<v-rating
														hover
														:length="5"
														:size="50"
														:model-value="3"
														active-color="yellow"
													/>
												</v-col>
											</v-row>
										</div>
									</v-col>
								</v-row>

								<!-- Compartilhar -->
								<v-row no-gutters class="mb-8 mt-4">
									<v-col cols="12">
										<div class="ml-8">
											<v-row>
												<v-col cols="4">
													<p class="rating-text text-left" style="padding-right: 5px;">Compartilhar: </p>
												</v-col>
												<v-col cols="8">
													<v-row>
														<v-col cols="3">
															<v-hover>
																<template v-slot:default="{ isHovering, props }">
																	<a href="">
																		<v-icon
																			size="x-large"
																			v-bind="props"
																			:color="isHovering ? 'purple' : 'white'"
																		>
																			mdi-instagram
																		</v-icon>
																	</a>
																</template>
															</v-hover>
														</v-col>

														<v-col cols="3">
															<v-hover>
																<template v-slot:default="{ isHovering, props }">
																	<a href="">
																		<v-icon
																			size="x-large"
																			v-bind="props"
																			:color="isHovering ? 'purple' : 'white'"
																		>
																			mdi-facebook
																		</v-icon>
																	</a>
																</template>
															</v-hover>
														</v-col>

														<v-col cols="3">
															<v-hover>
																<template v-slot:default="{ isHovering, props }">
																	<a href="">
																		<v-icon
																			size="x-large"
																			v-bind="props"
																			:color="isHovering ? 'purple' : 'white'"
																		>
																			mdi-reddit
																		</v-icon>
																	</a>
																</template>
															</v-hover>
														</v-col>

														<v-col cols="3">
															<v-hover>
																<template v-slot:default="{ isHovering, props }">
																	<a href="">
																		<v-icon
																			size="x-large"
																			v-bind="props"
																			:color="isHovering ? 'purple' : 'white'"
																		>
																			mdi-whatsapp
																		</v-icon>
																	</a>
																</template>
															</v-hover>
														</v-col>
													</v-row>
												</v-col>
											</v-row>
										</div>
									</v-col>
								</v-row>
							</div>
						</div>
					</v-col>

					<v-divider
						class="border-opacity-50"
						color="white"
						vertical
					></v-divider>

					<!-- div a direita com o conteúdo renderizado -->
					<v-col cols="8">
						<div class="modelo-conteudo">
							<div id="embed-doc" v-html="iframeHTML" class="embed-doc"></div>
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
	name: 'EstudoPage',
	data: () => ({
		pack: require('../../package.json'),
		livros: livros
	}),
	methods: {
		embedDocumento(pdfUrl) {
			const container = document.getElementById('embed-doc');

			container.innerHTML = '';

			const iframe = document.createElement('iframe');
			iframe.setAttribute('src', pdfUrl);
			iframe.setAttribute('width', '100%');
			iframe.setAttribute('height', '100%');
			iframe.style.border = 'none';

			container.appendChild(iframe);
		}
	},
	mounted() {
		console.log('livro: ', this.livros[0]);
		// const exemploPDF = '/pdfs/exemplo.pdf';
		const exemploPDF = 'https://dn790006.ca.archive.org/0/items/estruturas-de-dados-e-algoritmos-em-java-pdfdrive/Estruturas%20de%20dados%20e%20algoritmos%20em%20JAVA%20(%20PDFDrive%20).pdf'
		this.embedDocumento(exemploPDF);
	}
}
</script>

<style>
	.border {
		border: 1px solid black;
	}

	#embed-doc {
		width: 100%;
		height: 700px;
	}

	#embed-doc iframe {
		width: 100%;
		overflow-y: initial;
	}

	.titulo {
		font-weight: bolder;
		font-size: 1.5rem;
		color: white;
	}

	.custom-item {
		color: white;
		font-size: 15px;
		line-height: 2.5;
	}

	li {
		text-align: left;
		padding-top: 16px;
	}

	.rating-text {
		color: white;
		font-size: 1.3rem;
		font-family: Avenir, Helvetica, Arial, sans-serif;
		font-weight: bolder;
	}
</style>
