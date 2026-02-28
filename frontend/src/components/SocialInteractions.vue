<template>
  <v-card class="social-card overflow-hidden" elevation="3" rounded="lg">
    <v-toolbar flat color="transparent" dense>
      <v-toolbar-title class="text-subtitle-1 font-weight-bold grey--text text--darken-2">
        Interações Sociais
      </v-toolbar-title>
      <v-spacer></v-spacer>
      <v-chip outlined small color="primary" class="font-weight-medium">
        {{ likeCount }} {{ likeCount === 1 ? 'Curtida' : 'Curtidas' }}
      </v-chip>
    </v-toolbar>

    <v-divider></v-divider>

    <v-card-text>
      <div class="d-flex align-center mb-4">
        <v-btn
          icon
          large
          @click="toggleLike"
          :color="liked ? 'error' : 'grey darken-1'"
          :loading="loadingLike"
          class="pulse-hover"
        >
          <v-icon>{{ liked ? 'mdi-heart' : 'mdi-heart-outline' }}</v-icon>
        </v-btn>

        <v-btn icon large color="grey darken-1" @click="shareDialog = true" class="ml-2">
          <v-icon>mdi-share-variant</v-icon>
        </v-btn>

        <v-btn icon large color="grey darken-1" @click="messageDialog = true" class="ml-2">
          <v-icon>mdi-send-outline</v-icon>
        </v-btn>
      </div>

      <div class="comments-section mt-4">
        <h3 class="text-subtitle-2 font-weight-bold mb-3 d-flex align-center">
          <v-icon small left>mdi-comment-text-multiple-outline</v-icon>
          Comentários ({{ comments.length }})
        </h3>

        <v-slide-y-transition group>
          <v-list v-if="comments.length > 0" key="list" class="transparent pa-0">
            <v-list-item v-for="(c, index) in comments" :key="c.id" class="px-0 py-2">
              <v-list-item-avatar color="primary lighten-4" size="36">
                <span class="primary--text text-caption font-weight-bold">{{ index + 1 }}</span>
              </v-list-item-avatar>
              <v-list-item-content>
                <div class="d-flex justify-space-between align-center mb-1">
                  <span class="text-caption font-weight-bold grey--text text--darken-3">Usuário {{ c.usuario_id }}</span>
                  <span class="text-caption grey--text">{{ formatDate(c.created_at) }}</span>
                </div>
                <div class="text-body-2 grey--text text--darken-1 comment-text pa-2 rounded-lg grey lighten-4">
                  {{ c.texto }}
                </div>
              </v-list-item-content>
            </v-list-item>
          </v-list>
          <div v-else key="empty" class="text-center py-6 grey--text text-body-2 font-italic">
            Seja o primeiro a comentar!
          </div>
        </v-slide-y-transition>

        <div class="mt-4">
          <v-textarea
            v-model="newComment"
            label="Compartilhe sua opinião..."
            rows="2"
            auto-grow
            outlined
            hide-details
            dense
            rounded
            class="comment-area"
            background-color="grey lighten-5"
          ></v-textarea>
          <div class="d-flex justify-end mt-2">
            <v-btn
              color="primary"
              depressed
              rounded
              @click="postComment"
              :disabled="!newComment.trim()"
              :loading="loadingComment"
              class="text-none"
            >
              Publicar Comentário
            </v-btn>
          </div>
        </div>
      </div>
    </v-card-text>

    <!-- Share Dialog -->
    <v-dialog v-model="shareDialog" max-width="400" transition="dialog-bottom-transition">
      <v-card rounded="xl">
        <v-card-title class="headline flex-column align-start">
          <div class="text-h6 font-weight-bold primary--text">Compartilhar</div>
          <div class="text-caption grey--text">Espalhe o conhecimento com seus amigos</div>
        </v-card-title>
        <v-card-text class="pt-4">
          <v-row>
            <v-col v-for="p in platforms" :key="p.name" cols="4" class="text-center">
              <v-btn icon x-large :color="p.color" @click="selectPlatform(p.name)" class="mb-1">
                <v-icon>{{ p.icon }}</v-icon>
              </v-btn>
              <div class="text-caption font-weight-medium">{{ p.name }}</div>
            </v-col>
          </v-row>
        </v-card-text>
        <v-card-actions class="pa-4 justify-end">
          <v-btn text rounded @click="shareDialog = false">Fechar</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Message Dialog -->
    <v-dialog v-model="messageDialog" max-width="500" transition="dialog-bottom-transition">
      <v-card rounded="xl">
        <v-card-title class="headline primary white--text">
          <v-icon left color="white">mdi-email-send-outline</v-icon>
          Enviar para Amigo
        </v-card-title>
        <v-card-text class="pa-6">
          <v-text-field
            v-model="friendId"
            label="ID do Amigo"
            prepend-inner-icon="mdi-account-search"
            outlined
            rounded
            dense
          ></v-text-field>
          <v-textarea
            v-model="messageText"
            label="Mensagem (opcional)"
            placeholder="Olha que livro legal que encontrei na Biblioteca Digital!"
            rows="3"
            outlined
            rounded
            dense
          ></v-textarea>
        </v-card-text>
        <v-card-actions class="pa-4 pt-0">
          <v-spacer></v-spacer>
          <v-btn text rounded @click="messageDialog = false">Cancelar</v-btn>
          <v-btn
            color="primary"
            depressed
            rounded
            @click="sendMessage"
            :loading="loadingMessage"
            :disabled="!friendId"
          >
            Enviar Agora
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <v-snackbar v-model="snackbar" :color="snackbarColor" timeout="3000" rounded="pill">
      {{ snackbarText }}
      <template v-slot:action="{ attrs }">
        <v-btn text v-bind="attrs" @click="snackbar = false">Fechar</v-btn>
      </template>
    </v-snackbar>
  </v-card>
</template>

<script>
import axios from 'axios'

export default {
  name: 'SocialInteractions',
  props: {
    materialId: { type: Number, required: true },
    userId: { type: Number, required: true }
  },
  data() {
    return {
      liked: false,
      likeCount: 0,
      comments: [],
      newComment: '',
      shareDialog: false,
      messageDialog: false,
      friendId: null,
      messageText: '',

      loadingLike: false,
      loadingComment: false,
      loadingMessage: false,

      snackbar: false,
      snackbarText: '',
      snackbarColor: 'success',

      platforms: [
        { name: 'Facebook', icon: 'mdi-facebook', color: '#1877F2' },
        { name: 'Twitter', icon: 'mdi-twitter', color: '#1DA1F2' },
        { name: 'WhatsApp', icon: 'mdi-whatsapp', color: '#25D366' },
        { name: 'LinkedIn', icon: 'mdi-linkedin', color: '#0A66C2' },
        { name: 'Email', icon: 'mdi-email', color: '#EA4335' },
        { name: 'Link', icon: 'mdi-link-variant', color: '#757575' }
      ]
    }
  },
  created() {
    this.fetchLikes()
    this.fetchComments()
  },
  methods: {
    async fetchLikes() {
      try {
        const res = await axios.get(`/material/likes?material_id=${this.materialId}&usuario_id=${this.userId}`)
        this.likeCount = res.data.count
        this.liked = res.data.likedByUser
      } catch (e) { console.error(e) }
    },
    async toggleLike() {
      if (this.loadingLike) return
      this.loadingLike = true
      try {
        const res = await axios.post('/material/like', { usuario_id: this.userId, material_id: this.materialId })
        this.liked = res.data.liked
        this.fetchLikes() // Count refresh
      } catch (e) {
        this.showFeedback('Erro ao processar curtida', 'error')
      } finally {
        this.loadingLike = false
      }
    },
    async fetchComments() {
      try {
        const res = await axios.get(`/material/comments?material_id=${this.materialId}`)
        this.comments = res.data || []
      } catch (e) { console.error(e) }
    },
    async postComment() {
      if (!this.newComment.trim() || this.loadingComment) return
      this.loadingComment = true
      try {
        await axios.post('/material/comment', {
          usuario_id: this.userId,
          material_id: this.materialId,
          texto: this.newComment,
          created_at: new Date().toISOString()
        })
        this.newComment = ''
        await this.fetchComments()
        this.showFeedback('Comentário publicado!')
      } catch (e) {
        this.showFeedback('Erro ao publicar comentário', 'error')
      } finally {
        this.loadingComment = false
      }
    },
    async selectPlatform(platform) {
      try {
        await axios.post('/material/share', {
          usuario_id: this.userId,
          material_id: this.materialId,
          plataforma: platform,
          shared_at: new Date().toISOString()
        })
        this.shareDialog = false
        this.showFeedback(`Compartilhado no ${platform}!`)
      } catch (e) {
        this.showFeedback('Erro ao registrar compartilhamento', 'error')
      }
    },
    async sendMessage() {
      if (!this.friendId || this.loadingMessage) return
      this.loadingMessage = true
      try {
        await axios.post('/material/message', {
          remetente_id: this.userId,
          destinatario_id: parseInt(this.friendId),
          material_id: this.materialId,
          texto: this.messageText || 'Olha que livro legal que encontrei na Biblioteca Digital!',
          sent_at: new Date().toISOString()
        })
        this.messageDialog = false
        this.friendId = null
        this.messageText = ''
        this.showFeedback('Livro enviado com sucesso!')
      } catch (e) {
        this.showFeedback('Erro ao enviar mensagem', 'error')
      } finally {
        this.loadingMessage = false
      }
    },
    showFeedback(text, color = 'success') {
      this.snackbarText = text
      this.snackbarColor = color
      this.snackbar = true
    },
    formatDate(dateStr) {
      if (!dateStr) return ''
      const date = new Date(dateStr)
      return date.toLocaleDateString('pt-BR', { day: '2-digit', month: 'short' }) + ' ' +
             date.toLocaleTimeString('pt-BR', { hour: '2-digit', minute: '2-digit' })
    }
  }
}
</script>

<style scoped>
.social-card {
  max-width: 600px;
  margin: auto;
  border: 1px solid rgba(0, 0, 0, 0.05);
}

.comment-text {
  line-height: 1.4;
  font-size: 0.9rem;
}

.comment-area ::v-deep fieldset {
  border-color: rgba(0, 0, 0, 0.1) !important;
}

.pulse-hover:hover {
  transform: scale(1.1);
  transition: transform 0.2s ease-in-out;
}

.comments-section {
  border-top: 1px dashed rgba(0, 0, 0, 0.1);
  padding-top: 16px;
}
</style>
