<template>
	<div class="mx-auto mt-8">
		<v-container class="pa-0 border" fluid>
			<v-card min-width="100%" rounded="xl" style="background-color: rgba(36, 36, 36, 0.50);">
				<v-row no-gutters="" class="ma-4">
					<v-col cols="6" class="text-left">
						<v-chip
							variant="text"
							size="large"
							label
							class="text-white text-h4 ma-2"
						>
							Recentes
						</v-chip>
					</v-col>

					<v-col cols="6" class="text-right">
						<v-chip
							variant="text"
							label
							class="text-white text-h6 ma-2"
						>
							<p>Filtrar</p>
							<v-icon class="ml-2">mdi-filter</v-icon>
						</v-chip>
					</v-col>
				</v-row>

				<!-- v-row com o conteúdo já aceito pelo usuário -->
				<v-row no-gutters="" class="ma-4">
					<v-col cols="12" class="mt-6">
						<!-- Conteúdos -->
						<v-row no-gutters>
							<v-col cols="4">
								<!-- Conteúdo 01 -->
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

								<!-- Avaliar -->
								<v-row no-gutters class="mt-16">
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
								<v-row no-gutters class="mb-16">
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
							</v-col>

							<v-divider
								class="border-opacity-50"
								color="white"
								vertical
							></v-divider>

							<v-col cols="4">
								<!-- Conteúdo 02 -->
								<v-row no-gutter>
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

								<!-- Avaliar -->
								<v-row no-gutters class="mt-16">
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
								<v-row no-gutters class="mb-16">
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
							</v-col>

							<v-divider
								class="border-opacity-50"
								color="white"
								vertical
							></v-divider>

							<v-col cols="4">
								<!-- Conteúdo 03 -->
								<v-row no-gutter>
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

								<!-- Avaliar -->
								<v-row no-gutters class="mt-16">
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
								<v-row no-gutters class="mb-16">
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
							</v-col>
						</v-row>
					</v-col>
				</v-row>
			</v-card>
		</v-container>
	</div>
</template>

<script>
import livros from '../../../livros.json'

export default {
	name: 'DashboardPage',
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

	.custom-item {
		color: white;
		font-size: 15px;
		line-height: 2.5;
	}

	li {
		text-align: left;
	}

	.rating-text {
		color: white;
		font-size: 1.3rem;
		font-family: Avenir, Helvetica, Arial, sans-serif;
		font-weight: bolder;
	}
</style>
