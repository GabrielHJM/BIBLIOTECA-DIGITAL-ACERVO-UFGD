<template>
  <v-card
    class="premium-card-blur premium-shadow-hover clickable-card"
    elevation="0"
    :style="{ animationDelay: `${animationDelay}ms` }"
    @click="$router.push('/estudo/' + book.id)"
  >
    <v-row no-gutters>
      <!-- Image Section -->
      <v-col cols="5" class="pa-3">
        <div class="book-cover-wrapper position-relative">
          <!-- Premium Mandatory Glass Cover -->
          <div class="premium-book-cover-ios">
            <div class="mesh-gradient"></div>
            <div class="glass-overlay-ios">
              <div class="cover-content-ios">
                <v-icon class="cover-icon-ios mb-1" size="28" color="rgba(255,255,255,0.8)">{{ getBookIcon(book.categoria, book.titulo) }}</v-icon>
                <div class="cover-title-ios">{{ book.titulo }}</div>
                <div class="cover-summary-ios">{{ book.resumo }}</div>
              </div>
            </div>
            <div class="glass-shine"></div>
          </div>

          <v-chip class="source-badge" size="x-small" color="primary" variant="flat">
            {{ book.fonte || 'Repositório' }}
          </v-chip>
        </div>
      </v-col>

      <!-- Info Section -->
      <v-col cols="7" class="pa-3 text-left">
        <h3 class="item-title mb-2" :title="book.titulo">{{ book.titulo }}</h3>
        <div class="item-details">
          <p class="text-truncate"><strong>Autor:</strong> {{ book.autor }}</p>
          <p><strong>Categoria:</strong> {{ book.categoria }}</p>
          <p v-if="book.ano_publicacao"><strong>Ano:</strong> {{ book.ano_publicacao }}</p>
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
        @click.stop="handleToggleFavorite"
      >
        <v-icon ref="heartIcon">{{ isFavorited ? 'mdi-heart' : 'mdi-heart-outline' }}</v-icon>
      </v-btn>

      <div class="d-flex align-center gap-2">
        <v-btn
          icon="mdi-share-variant"
          variant="text"
          :color="$vuetify.theme.current.dark ? 'white' : 'grey-darken-1'"
          size="small"
          @click.stop="$emit('share', book)"
        ></v-btn>

        <v-btn
          class="ios-btn-open"
          variant="flat"
          size="small"
          size="small"
          @click.stop="$router.push('/estudo/' + book.id)"
        >
          {{ actionLabel }}
        </v-btn>
      </div>
    </v-card-actions>
  </v-card>
</template>

<script setup>
/* eslint-disable no-undef */
import { ref } from 'vue';
import { gsap } from 'gsap';

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
    default: 'Mais Detalhes'
  }
});

const emit = defineEmits(['toggle-favorite', 'share']);

const heartIcon = ref(null);

const handleToggleFavorite = () => {
  if (heartIcon.value) {
    const el = heartIcon.value.$el || heartIcon.value;
    const tl = gsap.timeline();
    tl.to(el, { scale: 1.5, rotation: 15, duration: 0.15, ease: "back.out(2)" })
      .to(el, { rotation: -15, duration: 0.1, ease: "none" })
      .to(el, { scale: 1, rotation: 0, duration: 0.3, ease: "elastic.out(1, 0.3)" });
  }
  emit('toggle-favorite', props.book);
};

const getBookIcon = (category, title) => {
  const text = ((category || '') + ' ' + (title || '')).toLowerCase();
  if (text.includes('tecnologia') || text.includes('comput') || text.includes('software') || text.includes('program') || text.includes('digital')) return 'mdi-laptop';
  if (text.includes('saúde') || text.includes('medicina') || text.includes('biolog') || text.includes('enferm') || text.includes('médic')) return 'mdi-heart-pulse';
  if (text.includes('direito') || text.includes('lei') || text.includes('juríd') || text.includes('advog')) return 'mdi-gavel';
  if (text.includes('matemát') || text.includes('física') || text.includes('cálculo') || text.includes('engenh')) return 'mdi-calculator';
  if (text.includes('história') || text.includes('socio') || text.includes('psico') || text.includes('filo')) return 'mdi-bank';
  if (text.includes('literat') || text.includes('poesia') || text.includes('romance')) return 'mdi-feather';
  if (text.includes('educação') || text.includes('ensino') || text.includes('pedagog')) return 'mdi-school';
  return 'mdi-book-open-page-variant';
};
</script>

<style scoped>
.book-cover-wrapper {
  perspective: 1200px;
  width: 100%;
}

.premium-book-cover-ios {
  width: 100%;
  aspect-ratio: 1 / 1.4;
  max-height: 150px;
  margin: 0 auto;
  position: relative;
  overflow: hidden;
  border-radius: 12px;
  background: #000;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.5);
  transition: all 0.6s cubic-bezier(0.16, 1, 0.3, 1);
}

.mesh-gradient {
  position: absolute;
  top: -50%;
  left: -50%;
  width: 200%;
  height: 200%;
  background: 
    radial-gradient(circle at 20% 20%, #007AFF 0%, transparent 40%),
    radial-gradient(circle at 80% 25%, #5AC8FA 0%, transparent 40%),
    radial-gradient(circle at 30% 80%, #0051FF 0%, transparent 40%),
    radial-gradient(circle at 75% 75%, #00D4E8 0%, transparent 40%),
    radial-gradient(circle at 50% 50%, #0033FF 0%, transparent 50%);
  filter: blur(40px);
  animation: mesh-rotate 15s linear infinite;
  opacity: 0.7;
}

@keyframes mesh-rotate {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

.glass-overlay-ios {
  position: absolute;
  inset: 6px;
  background: rgba(255, 255, 255, 0.05);
  backdrop-filter: blur(20px) saturate(180%);
  border: 1px solid rgba(255, 255, 255, 0.15);
  border-radius: 10px;
  z-index: 2;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  padding: 12px;
  text-align: center;
}

.cover-content-ios {
  position: relative;
  z-index: 3;
}

.cover-title-ios {
  color: #ffffff;
  font-family: 'Outfit', 'Inter', sans-serif;
  font-size: 0.75rem;
  font-weight: 800;
  line-height: 1.1;
  max-width: 100%;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  text-shadow: 0 2px 10px rgba(0, 0, 0, 0.5);
  margin-bottom: 4px;
}

.cover-summary-ios {
  color: rgba(255, 255, 255, 0.6);
  font-size: 0.55rem;
  font-weight: 500;
  line-height: 1.3;
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
  text-align: center;
}

.glass-shine {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: linear-gradient(135deg, rgba(255,255,255,0.15) 0%, transparent 50%, rgba(255,255,255,0.05) 100%);
  pointer-events: none;
  z-index: 4;
}

.clickable-card {
  cursor: pointer;
  transition: all 0.4s var(--spring-easing) !important;
}

.clickable-card:active {
  transform: scale(0.98);
}

.premium-card-blur:hover .premium-book-cover-ios {
  transform: scale(1.05) translateY(-5px);
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.6);
}

.premium-card-blur:hover .mesh-gradient {
  opacity: 1;
  animation-duration: 8s;
}

.source-badge {
  position: absolute;
  top: 8px;
  left: 8px;
  font-weight: 800;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  box-shadow: 0 4px 10px rgba(0,0,0,0.5);
  z-index: 20 !important; /* Sempre acima das camadas de vidro */
  pointer-events: none;
}

@media (max-width: 600px) {
  .source-badge {
    top: 4px;
    left: 4px;
    font-size: 0.6rem !important;
    padding: 0 4px !important;
  }
}

.item-title {
  font-family: 'Outfit', sans-serif;
  font-size: 1rem;
  font-weight: 700;
  line-height: 1.3;
  color: var(--v-theme-on-surface) !important;
  display: -webkit-box;
  -webkit-line-clamp: 4;
  -webkit-box-orient: vertical;
  overflow: hidden;
  text-overflow: ellipsis;
}

.item-details p {
  font-size: 0.85rem;
  margin-bottom: 2px;
  opacity: 0.9;
  color: var(--v-theme-on-surface) !important;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.ios-btn-open {
  background: linear-gradient(135deg, #007AFF 0%, #0056B3 100%) !important;
  color: white !important;
  text-transform: none;
  font-weight: 600;
  border-radius: 16px;
  padding: 0 16px;
}





.gap-2 {
  gap: 8px;
}
</style>
