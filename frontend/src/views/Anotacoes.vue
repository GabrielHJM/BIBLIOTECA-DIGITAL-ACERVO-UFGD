<template>
  <div class="anotacoes-page">
    <!-- Header compacto e arredondado -->
    <div class="mb-6">
      <div class="page-header">
        <div class="d-flex align-center mb-3">
          <v-btn
            variant="text"
            color="white"
            prepend-icon="mdi-chevron-left"
            class="back-btn"
            @click="$router.push('/dashboard')"
          >
            <span class="back-label">Painel</span>
          </v-btn>
        </div>
        <div class="d-flex align-center justify-space-between flex-wrap gap-4">
          <div>
            <h1 class="page-title text-white">Minhas Anotações</h1>
            <p class="page-subtitle text-white opacity-60">Seu caderno digital de estudos</p>
          </div>
          <v-btn
            class="nova-btn font-weight-bold"
            prepend-icon="mdi-plus"
            size="large"
            @click="abrirModal(null)"
          >
            Nova Anotação
          </v-btn>
        </div>
      </div>
    </div>

    <!-- Content -->
    <v-container class="content-area mt-4">
      <!-- Loading -->
      <div v-if="loading" class="d-flex justify-center align-center" style="min-height:300px">
        <v-progress-circular indeterminate color="cyan-accent-3" size="64"></v-progress-circular>
      </div>

      <!-- Empty state -->
      <div v-else-if="anotacoes.length === 0" class="empty-state text-center py-16">
        <v-icon size="72" color="white" class="opacity-30 mb-4">mdi-notebook-outline</v-icon>
        <h3 class="text-h5 text-white opacity-60 mb-2">Nenhuma anotação ainda</h3>
        <p class="text-body-2 text-white opacity-40 mb-6">Comece a resumir seus estudos agora mesmo.</p>
        <v-btn class="nova-btn font-weight-bold" prepend-icon="mdi-plus" @click="abrirModal(null)">
          Nova Anotação
        </v-btn>
      </div>

      <!-- Grid de notas -->
      <v-row v-else class="notes-grid">
        <v-col
          v-for="nota in anotacoes"
          :key="nota.id"
          cols="12" sm="6" md="4" lg="3"
        >
          <div class="note-card" @click="abrirModal(nota)">
            <div class="note-card-top">
              <div class="d-flex justify-space-between align-start">
                <h3 class="note-title text-white">{{ nota.titulo || 'Sem Título' }}</h3>
                <v-btn
                  icon="mdi-delete-outline"
                  variant="text"
                  size="x-small"
                  color="white"
                  class="delete-btn opacity-50"
                  @click.stop="deletarAnotacao(nota.id)"
                ></v-btn>
              </div>
              <p class="note-content text-white opacity-80 mt-2">{{ nota.conteudo }}</p>
            </div>
            <div class="note-footer d-flex justify-space-between align-center">
              <div v-if="nota.material_id" class="d-flex align-center">
                <v-icon size="12" color="white" class="mr-1 opacity-60">mdi-book-open-variant</v-icon>
                <span class="text-caption text-white opacity-50">Vinculado</span>
              </div>
              <div v-else></div>
              <span class="text-caption text-white opacity-40">{{ formatData(nota.data_atualizacao) }}</span>
            </div>
          </div>
        </v-col>
      </v-row>
    </v-container>

    <!-- Modal de criação / edição -->
    <v-dialog v-model="modalAberto" max-width="580px" persistent>
      <v-card
        class="modal-card"
        :style="{ color: 'white' }"
        rounded="xl"
      >
        <v-card-title class="d-flex justify-space-between align-center pa-5 pb-3">
          <div>
            <span class="modal-title">{{ notaAtual.id ? 'Editar Anotação' : 'Nova Anotação' }}</span>
            <div class="modal-subtitle">{{ notaAtual.id ? 'Atualize o conteúdo da sua nota' : 'Capture seus pensamentos e ideias' }}</div>
          </div>
          <v-btn icon="mdi-close" variant="text" color="white" size="small" class="opacity-60" @click="fecharModal"></v-btn>
        </v-card-title>

        <v-divider class="opacity-10"></v-divider>

        <v-card-text class="pa-5">
          <v-text-field
            v-model="notaAtual.titulo"
            label="Título da anotação"
            placeholder="Ex: Resumo de Capítulo 3..."
            variant="outlined"
            color="cyan-accent-3"
            class="mb-4"
            maxlength="80"
            counter
            hide-details="auto"
            prepend-inner-icon="mdi-format-title"
          ></v-text-field>

          <v-textarea
            v-model="notaAtual.conteudo"
            label="Conteúdo"
            placeholder="Escreva aqui seus resumos, ideias ou pontos importantes..."
            variant="outlined"
            color="cyan-accent-3"
            rows="5"
            auto-grow
            class="mb-5"
            maxlength="2000"
            counter
            hide-details="auto"
            prepend-inner-icon="mdi-text"
          ></v-textarea>

        </v-card-text>

        <v-divider class="opacity-10"></v-divider>

        <v-card-actions class="pa-4">
          <v-spacer></v-spacer>
          <v-btn variant="text" color="white" class="opacity-60 mr-2" @click="fecharModal">Cancelar</v-btn>
          <v-btn
            class="nova-btn font-weight-bold px-6"
            :loading="salvando"
            @click="salvarAnotacao"
          >
            Salvar
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import AnotacaoService from '@/services/AnotacaoService';

const anotacoes = ref([]);
const loading = ref(true);
const modalAberto = ref(false);
const salvando = ref(false);

const notaVazia = { id: null, titulo: '', conteudo: '', cor: 'rgba(0,188,212,0.25)', material_id: null };
const notaAtual = ref({ ...notaVazia });

const carregarAnotacoes = async () => {
  loading.value = true;
  try {
    // api.js interceptor já desempacota {success:true, data:[...]} → retorna o array diretamente
    const res = await AnotacaoService.listar();
    anotacoes.value = Array.isArray(res) ? res : (res || []);
  } catch (error) {
    console.error('Erro ao buscar anotações', error);
    anotacoes.value = [];
  } finally {
    loading.value = false;
  }
};

const abrirModal = (nota) => {
  notaAtual.value = nota ? { ...nota } : { ...notaVazia };
  modalAberto.value = true;
};

const fecharModal = () => {
  modalAberto.value = false;
  setTimeout(() => { notaAtual.value = { ...notaVazia }; }, 200);
};

const salvarAnotacao = async () => {
  if (!notaAtual.value.conteudo.trim()) return;
  salvando.value = true;
  try {
    if (notaAtual.value.id) {
      await AnotacaoService.atualizar(notaAtual.value.id, notaAtual.value);
    } else {
      await AnotacaoService.criar(notaAtual.value);
    }
    await carregarAnotacoes();
    fecharModal();
  } catch (error) {
    console.error('Erro ao salvar', error);
  } finally {
    salvando.value = false;
  }
};

const deletarAnotacao = async (id) => {
  if (confirm('Tem certeza que deseja excluir esta anotação?')) {
    try {
      await AnotacaoService.excluir(id);
      await carregarAnotacoes();
    } catch (error) {
      console.error('Erro ao excluir', error);
    }
  }
};

const formatData = (isoString) => {
  if (!isoString) return '';
  return new Intl.DateTimeFormat('pt-BR', { day: '2-digit', month: 'short' }).format(new Date(isoString));
};

onMounted(carregarAnotacoes);
</script>

<style scoped>
/* ── Container principal (sem background próprio, herda o tema) ── */
.anotacoes-page {
  min-height: 100vh;
  padding-bottom: 60px;
}

/* ── Header compacto e arredondado ─────────────────────────── */
.page-header {
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 20px;
  padding: 16px 20px 20px;
}


.back-btn {
  font-size: 15px !important;
  text-transform: none !important;
  letter-spacing: 0 !important;
  opacity: 0.8;
  transition: opacity 0.2s;
}
.back-btn:hover { opacity: 1; }
.back-label { font-size: 15px; }

.page-title {
  font-size: 2rem;
  font-weight: 700;
  letter-spacing: -0.5px;
  line-height: 1.2;
}
.page-subtitle { font-size: 0.95rem; margin-top: 2px; }

/* ── Botão padrão Nova Anotação ────────────────────────────── */
.nova-btn {
  background: linear-gradient(135deg, #00BCD4, #00838F) !important;
  color: white !important;
  border-radius: 12px !important;
  text-transform: none !important;
  letter-spacing: 0.3px !important;
  box-shadow: 0 4px 15px rgba(0, 188, 212, 0.3) !important;
  transition: all 0.3s ease !important;
}
.nova-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(0, 188, 212, 0.4) !important;
}

/* ── Content ────────────────────────────────────────────────── */
.content-area { max-width: 1200px; }
.notes-grid { padding: 8px 0; }

/* ── Note Cards ─────────────────────────────────────────────── */
.note-card {
  background: rgba(0, 188, 212, 0.1);
  border: 1px solid rgba(0, 188, 212, 0.2);
  border-radius: 20px;
  padding: 20px;
  min-height: 180px;
  max-height: 240px;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.25, 0.8, 0.25, 1);
  overflow: hidden;
  position: relative;
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
}

.note-card::before {
  content: '';
  position: absolute;
  inset: 0;
  border-radius: inherit;
  background: linear-gradient(135deg, rgba(0,188,212,0.05) 0%, transparent 60%);
  pointer-events: none;
}

.note-card:hover {
  transform: translateY(-6px) scale(1.01);
  border-color: rgba(0, 188, 212, 0.5);
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.4), 0 0 0 1px rgba(0,188,212,0.2), 0 0 30px rgba(0,188,212,0.1) !important;
}

.note-card-top { flex: 1; overflow: hidden; }

.note-title {
  font-size: 1rem;
  font-weight: 700;
  line-height: 1.3;
  display: -webkit-box;
  -webkit-line-clamp: 1;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.note-content {
  font-size: 0.875rem;
  line-height: 1.55;
  display: -webkit-box;
  -webkit-line-clamp: 4;
  -webkit-box-orient: vertical;
  overflow: hidden;
  white-space: pre-wrap;
}

.note-footer {
  border-top: 1px solid rgba(255, 255, 255, 0.07);
  padding-top: 10px;
  margin-top: 10px;
}

.delete-btn { opacity: 0 !important; transition: opacity 0.2s; }
.note-card:hover .delete-btn { opacity: 0.6 !important; }
.delete-btn:hover { opacity: 1 !important; }

/* ── Modal ──────────────────────────────────────────────────── */
.modal-card {
  background: rgb(var(--v-theme-surface)) !important;
  border: 1px solid rgba(255, 255, 255, 0.06);
}

.modal-title {
  font-size: 1.15rem;
  font-weight: 700;
  letter-spacing: -0.3px;
  color: white;
}

.modal-subtitle {
  font-size: 0.8rem;
  opacity: 0.5;
  margin-top: 2px;
  color: white;
}
:deep(.modal-field input),
:deep(.modal-field textarea) {
  color: white !important;
}

/* ── Color Dots ─────────────────────────────────────────────── */
.color-dot {
  width: 22px;
  height: 22px;
  border-radius: 50%;
  border: 2px solid transparent;
  cursor: pointer;
  transition: all 0.2s;
  outline: none;
}
.color-dot:hover { transform: scale(1.15); }
.color-dot--active {
  border-color: white;
  transform: scale(1.2);
  box-shadow: 0 0 0 2px rgba(255,255,255,0.3);
}

/* ── Gaps ────────────────────────────────────────────────────── */
.gap-2 { gap: 8px; }
.gap-4 { gap: 16px; }
</style>
