<template>
  <v-card
    class="ios-item-card premium-shadow-hover"
    elevation="12"
    :style="cardStyle"
  >
    <v-row no-gutters>
      <!-- Image Section -->
      <v-col cols="5" class="pa-3">
        <div class="book-cover-wrapper position-relative">
          <img
            :src="book.capa_url || defaultCover"
            :alt="book.titulo"
            class="book-cover"
          />
          <v-chip class="source-badge" size="x-small" color="cyan-darken-3" variant="flat">
            {{ book.fonte || 'Repositório' }}
          </v-chip>
        </div>
      </v-col>

      <!-- Info Section -->
      <v-col cols="7" class="pa-3 text-left">
        <h3 class="item-title mb-1 text-truncate">{{ book.titulo }}</h3>
        <div class="item-details">
          <p class="text-truncate"><strong>Autor:</strong> {{ book.autor }}</p>
          <p><strong>Categoria:</strong> {{ book.categoria }}</p>
          <p v-if="book.ano_publicacao"><strong>Ano:</strong> {{ book.ano_publicacao }}</p>
        </div>

        <div class="mt-2 d-flex align-center justify-space-between">
          <v-rating
            :model-value="book.media_nota || 0"
            density="compact"
            color="amber"
            active-color="amber"
            size="small"
            readonly
          ></v-rating>

          <div class="xp-badge ml-2">+{{ book.xp || 10 }} XP</div>
        </div>

        <div class="mt-2 text-caption d-flex align-center">
          <v-icon size="small" color="blue-lighten-3" class="mr-1">mdi-brain</v-icon>
          <span class="difficulty-text">Dificuldade: </span>
          <div class="difficulty-dots ml-1">
            <span
              v-for="i in 5"
              :key="i"
              class="dot"
              :class="{ 'active': i <= (book.dificuldade || 1) }"
            ></span>
          </div>
        </div>
      </v-col>
    </v-row>

    <!-- Actions Footer -->
    <v-divider class="mx-4 opacity-30"></v-divider>
    <v-card-actions class="pa-3 justify-space-between">
      <v-btn
        icon
        variant="text"
        :color="isFavorited ? 'pink' : ($vuetify.theme.current.dark ? 'white' : 'grey-darken-1')"
        size="small"
        @click="$emit('toggle-favorite', book)"
      >
        <v-icon>{{ isFavorited ? 'mdi-heart' : 'mdi-heart-outline' }}</v-icon>
      </v-btn>

      <div class="d-flex align-center gap-2">
        <v-btn
          icon="mdi-share-variant"
          variant="text"
          :color="$vuetify.theme.current.dark ? 'white' : 'grey-darken-1'"
          size="small"
          @click="$emit('share', book)"
        ></v-btn>

        <v-btn
          class="ios-btn-open"
          variant="flat"
          size="small"
          @click="$router.push('/estudo/' + book.id)"
        >
          {{ actionLabel }}
        </v-btn>
      </div>
    </v-card-actions>
  </v-card>
</template>

<script setup>
/* global defineProps, defineEmits */
import { computed } from 'vue';

const props = defineProps({
  book: {
    type: Object,
    required: true
  },
  isFavorited: {
    type: Boolean,
    default: false
  },
  animationDelay: {
    type: Number,
    default: 0
  },
  actionLabel: {
    type: String,
    default: 'Ler Agora'
  }
});

defineEmits(['toggle-favorite', 'share']);

const defaultCover = 'https://images.unsplash.com/photo-1507842217343-583bb7270b66?q=80&w=400';

const cardStyle = computed(() => ({
  animationDelay: `${props.animationDelay}ms`
}));
</script>

<style scoped>
.book-cover-wrapper {
  perspective: 1200px;
  height: 150px;
}

.book-cover {
  width: 100%;
  height: 100%;
  object-fit: cover;
  border-radius: 12px;
  box-shadow: 0 10px 20px rgba(0, 0, 0, 0.2);
  transition: transform 0.6s cubic-bezier(0.16, 1, 0.3, 1);
}

.premium-glass-card:hover .book-cover {
  transform: rotateY(-12deg) translateY(-5px) scale(1.08);
  box-shadow: 0 15px 30px rgba(0, 0, 0, 0.3);
}

.source-badge {
  position: absolute;
  top: 8px;
  left: 8px;
  font-weight: 800;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  box-shadow: 0 4px 10px rgba(0,0,0,0.5);
}

.item-title {
  font-family: 'Outfit', sans-serif;
  font-size: 1.1rem;
  font-weight: 700;
  line-height: 1.2;
  color: white !important;
  height: 2.6em;
  overflow: hidden;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}

.item-details p {
  font-size: 0.85rem;
  margin-bottom: 2px;
  opacity: 0.9;
  color: rgba(255, 255, 255, 0.8) !important;
}

.ios-btn-open {
  background: linear-gradient(135deg, #00BCD4 0%, #2196F3 100%) !important;
  color: white !important;
  text-transform: none;
  font-weight: 600;
  border-radius: 16px;
  padding: 0 16px;
}

.xp-badge {
  background: rgba(0, 255, 157, 0.1);
  color: #00ff9d !important;
  border: 1px solid rgba(0, 255, 157, 0.3);
  padding: 2px 8px;
  border-radius: 8px;
  font-size: 0.75rem;
  font-weight: 800;
  font-family: 'Outfit', sans-serif;
}

.difficulty-text {
  font-weight: 600;
  opacity: 0.8;
}

.difficulty-dots {
  display: flex;
  gap: 3px;
}

.dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.2);
}

.dot.active {
  background: #00BCD4;
  box-shadow: 0 0 5px #00BCD4;
}

.gap-2 {
  gap: 8px;
}
</style>
