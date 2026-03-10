<template>
  <div
    class="premium-card"
    :class="[
      { 'premium-card-dark': isDarkTheme },
      { 'premium-card-blur': blur },
      { 'premium-card-hover': hover }
    ]"
    :style="containerStyle"
  >
    <div v-if="$slots.header || title" class="premium-card-header d-flex align-center justify-space-between mb-4">
      <slot name="header">
        <h2 class="text-h5 font-weight-bold d-flex align-center">
          <v-icon v-if="icon" :color="iconColor" class="mr-2">{{ icon }}</v-icon>
          {{ title }}
        </h2>
      </slot>
      <div v-if="$slots.action" class="premium-card-action">
        <slot name="action"></slot>
      </div>
    </div>

    <div class="premium-card-content">
      <slot></slot>
    </div>

    <div v-if="$slots.footer" class="premium-card-footer mt-4">
      <slot name="footer"></slot>
    </div>
  </div>
</template>

<script setup>
/* global defineProps */
import { computed } from 'vue';
import { useTheme } from 'vuetify';

const props = defineProps({
  title: String,
  icon: String,
  iconColor: {
    type: String,
    default: 'primary'
  },
  padding: {
    type: String,
    default: '18px'
  },
  blur: {
    type: Boolean,
    default: true
  },
  hover: {
    type: Boolean,
    default: true
  }
});

const theme = useTheme();
const isDarkTheme = computed(() => theme.global.current.value.dark);

const containerStyle = computed(() => ({
  padding: props.padding
}));
</script>

<style scoped>
.premium-card {
  background: var(--glass-bg);
  backdrop-filter: var(--glass-blur);
  -webkit-backdrop-filter: var(--glass-blur);
  border-radius: 24px;
  border: 1px solid var(--glass-border);
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.1);
  transition: all 0.4s cubic-bezier(0.165, 0.84, 0.44, 1);
  position: relative;
  overflow: hidden;
  color: #ffffff;
}

.premium-card-dark {
  background: var(--v-theme-surface);
  border-color: rgba(255, 255, 255, 0.05);
  color: white;
}

.premium-card-blur {
  background: var(--glass-bg) !important;
  backdrop-filter: var(--glass-blur) !important;
}

.premium-card-hover:hover {
  transform: translateY(-5px);
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.08);
}

.premium-card-header h2 {
  font-family: 'Outfit', sans-serif;
  letter-spacing: -0.5px;
}
</style>
