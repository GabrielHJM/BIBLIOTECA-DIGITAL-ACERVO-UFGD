<template>
	<div class="admin-container mt-16 pt-8 pb-16">
		<v-container>
			<div class="d-flex align-center justify-space-between mb-10">
				<div>
					<h1 class="text-h3 font-weight-black text-white mb-2">Painel do Curador</h1>
					<p class="text-subtitle-1 text-white opacity-60">Gestão de Conteúdo e Métricas Globais</p>
				</div>
				<v-btn color="cyan" variant="tonal" prepend-icon="mdi-plus" rounded="lg" @click="dialogNovo = true">
					Novo Material
				</v-btn>
			</div>

			<!-- Metrics Row -->
			<v-row v-if="metricas" class="mb-12">
				<v-col cols="12" md="3" v-for="(val, key) in statsBrief" :key="key">
					<v-card class="metric-card premium-glow pa-6 rounded-xl" elevation="4">
						<div class="text-overline opacity-60">{{ key }}</div>
						<div class="text-h4 font-weight-black mt-2">{{ val }}</div>
					</v-card>
				</v-col>
			</v-row>

			<!-- Tabs for Moderation and Extended Metrics -->
			<v-tabs v-model="tab" color="cyan" align-tabs="start" class="mb-6">
				<v-tab value="pendentes">Moderação ({{ pendentes.length }})</v-tab>
				<v-tab value="metricas">Estatísticas Detalhadas</v-tab>
			</v-tabs>

			<v-window v-model="tab">
				<!-- Moderation Tab -->
				<v-window-item value="pendentes">
					<v-card class="moderation-card ios-glass-card rounded-xl overflow-hidden">
						<v-table class="admin-table">
							<thead>
								<tr>
									<th>Título</th>
									<th>Autor</th>
									<th>Fonte</th>
									<th class="text-right">Ações</th>
								</tr>
							</thead>
							<tbody v-if="pendentes.length > 0">
								<tr v-for="item in pendentes" :key="item.id">
									<td class="font-weight-bold">{{ item.titulo }}</td>
									<td>{{ item.autor }}</td>
									<td>
										<v-chip size="x-small" :color="item.fonte === 'Local' ? 'cyan' : 'amber'">
											{{ item.fonte || 'Local' }}
										</v-chip>
									</td>
									<td class="text-right">
										<v-btn icon="mdi-check" color="success" variant="text" size="small" @click="moderar(item.id, 'aprovado')"></v-btn>
										<v-btn icon="mdi-close" color="red" variant="text" size="small" @click="moderar(item.id, 'rejeitado')"></v-btn>
										<v-btn icon="mdi-eye" color="white" variant="text" size="small" @click="verDetalhes(item)"></v-btn>
									</td>
								</tr>
							</tbody>
							<tbody v-else>
								<tr>
									<td colspan="4" class="text-center py-10 opacity-30">Nenhum material pendente para aprovação.</td>
								</tr>
							</tbody>
						</v-table>
					</v-card>
				</v-window-item>

				<!-- Detailed Metrics Tab -->
				<v-window-item value="metricas">
					<v-row>
						<v-col cols="12" md="6">
							<v-card class="ios-glass-card pa-6 rounded-xl">
								<h3 class="mb-4">Distribuição por Fonte</h3>
								<div v-if="metricas" class="source-list">
									<div v-for="(count, fonte) in metricas.fontes" :key="fonte" class="d-flex justify-space-between mb-2">
										<span>{{ fonte || 'Manual/Local' }}</span>
										<span class="font-weight-black text-cyan">{{ count }} obras</span>
									</div>
								</div>
							</v-card>
						</v-col>
						<v-col cols="12" md="6">
							<v-card class="ios-glass-card pa-6 rounded-xl d-flex flex-column justify-center align-center h-100">
								<v-icon size="64" color="cyan" class="mb-4 opacity-20">mdi-chart-timeline-variant</v-icon>
								<p class="text-center opacity-40">Gráficos de tendência em tempo real estarão disponíveis na v2.</p>
							</v-card>
						</v-col>
					</v-row>
				</v-window-item>
			</v-window>
		</v-container>
	</div>
</template>

<script>
import AdminService from '@/services/AdminService';
import auth from '@/auth';

export default {
	name: 'AdminDashboard',
	data: () => ({
		tab: 'pendentes',
		pendentes: [],
		metricas: null,
		loading: true,
		dialogNovo: false
	}),
	computed: {
		statsBrief() {
			if (!this.metricas) return {};
			return {
				"Alunos": this.metricas.total_usuarios,
				"Livros": this.metricas.total_materiais,
				"Leituras": this.metricas.total_leituras,
				"Pendente": this.pendentes.length
			};
		}
	},
	methods: {
		async carregarTudo() {
			try {
				const [pResp, mResp] = await Promise.all([
					AdminService.listarPendentes(),
					AdminService.obterMetricas()
				]);
				this.pendentes = pResp.data || [];
				this.metricas = mResp.data;
			} catch (err) {
				console.error('Erro ao carregar admin:', err);
			} finally {
				this.loading = false;
			}
		},
		async moderar(id, status) {
			try {
				await AdminService.moderarMaterial(id, status, auth.getUser().id);
				await this.carregarTudo();
			} catch (err) {
				console.error('Erro ao moderar:', err);
			}
		},
		verDetalhes(item) {
			if (item.pdf_url) window.open(item.pdf_url, '_blank');
		}
	},
	mounted() {
		this.carregarTudo();
	}
}
</script>

<style scoped>
	.admin-container { min-height: 100vh; background: #0f172a; color: white; }

	.metric-card {
		background: rgba(255, 255, 255, 0.03) !important;
		border: 1px solid rgba(0, 184, 212, 0.1);
	}

	.ios-glass-card {
		background: rgba(255, 255, 255, 0.02);
		backdrop-filter: blur(20px);
		border: 1px solid rgba(255, 255, 255, 0.05);
	}

	.admin-table { background: transparent !important; color: white; }
	.admin-table th { color: rgba(255,255,255,0.4) !important; text-transform: uppercase; font-size: 11px; letter-spacing: 1px; }

	.premium-glow:hover {
		box-shadow: 0 0 30px rgba(0, 184, 212, 0.15) !important;
		border-color: rgba(0, 184, 212, 0.4);
		transform: translateY(-5px);
	}
</style>
