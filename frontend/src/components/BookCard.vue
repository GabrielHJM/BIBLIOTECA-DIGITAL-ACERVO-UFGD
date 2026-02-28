<template>
  <v-card
    class="ios-item-card premium-shadow-hover"
    elevation="8"
    :style="cardStyle"
  >
    <v-row no-gutters>
      <!-- Image Section -->
      <v-col cols="5" class="pa-3">
        <div class="book-cover-wrapper">
          <img
            :src="book.capa_url || defaultCover"
            :alt="book.titulo"
            class="book-cover"
          />
        </div>
      </v-col>

      <!-- Info Section -->
      <v-col cols="7" class="pa-4 text-left">
        <h3 class="item-title mb-1 text-truncate">{{ book.titulo }}</h3>
        <div class="item-details">
          <p class="text-truncate"><strong>Autor:</strong> {{ book.autor }}</p>
          <p><strong>Categoria:</strong> {{ book.categoria }}</p>
          <p v-if="book.ano_publicacao"><strong>Ano:</strong> {{ book.ano_publicacao }}</p>
        </div>

        <div class="mt-2">
          <v-rating
            :model-value="book.media_nota || 0"
            density="compact"
            color="amber"
            active-color="amber"
            size="small"
            readonly
          ></v-rating>
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

const defaultCover = 'https://images.unsplash.com/photo-1544947950-fa07a98d237f?q=80&w=200';

const cardStyle = computed(() => ({
  animationDelay: `${props.animationDelay}ms`
}));
</script>

<style scoped>
.ios-item-card {
  border-radius: 20px;
  background: rgba(255, 255, 255, 0.7);
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.3);
  transition: all 0.4s cubic-bezier(0.165, 0.84, 0.44, 1);
  overflow: hidden;
  height: 100%;
}

.dark .ios-item-card {
  background: rgba(30, 30, 30, 0.7);
  border-color: rgba(255, 255, 255, 0.1);
}

.book-cover-wrapper {
  perspective: 1000px;
  height: 140px;
}

.book-cover {
  width: 100%;
  height: 100%;
  object-fit: cover;
  border-radius: 12px;
  box-shadow: 5px 5px 15px rgba(0, 0, 0, 0.2);
  transition: transform 0.5s ease;
}

.ios-item-card:hover .book-cover {
  transform: rotateY(-10deg) scale(1.05);
}

.item-title {
  font-family: 'Outfit', sans-serif;
  font-size: 1.1rem;
  font-weight: 700;
  line-height: 1.2;
}

.item-details p {
  font-size: 0.85rem;
  margin-bottom: 2px;
  opacity: 0.8;
}

.ios-btn-open {
  background: linear-gradient(135deg, #00BCD4 0%, #2196F3 100%) !important;
  color: white !important;
  text-transform: none;
  font-weight: 600;
  border-radius: 12px;
  padding: 0 16px;
}

.gap-2 {
  gap: 8px;
}
</style>
