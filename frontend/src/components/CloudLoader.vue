<template>
  <div class="cloud-loader-wrapper" :style="{ minHeight: height }">
    <div class="cloud-pulse-container">
      <v-icon class="cloud-icon" size="64" :color="color">mdi-cloud</v-icon>
      <div class="rain-drops">
        <div class="drop drop-1"></div>
        <div class="drop drop-2"></div>
        <div class="drop drop-3"></div>
      </div>
      <div class="cloud-shadow"></div>
    </div>
    <div v-if="text" class="cloud-text mt-4">{{ text }}</div>
  </div>
</template>

<script>
export default {
  name: 'CloudLoader',
  props: {
    text: {
      type: String,
      default: 'Carregando...'
    },
    color: {
      type: String,
      default: '#00E5FF' // Padrão azul neon da Acervus Core
    },
    height: {
      type: String,
      default: '300px'
    }
  }
}
</script>

<style scoped>
.cloud-loader-wrapper {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  width: 100%;
}

.cloud-pulse-container {
  position: relative;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}

.cloud-icon {
  animation: pulseCloud 1.5s ease-in-out infinite alternate, floatCloud 3s ease-in-out infinite;
  filter: drop-shadow(0 0 15px rgba(0, 229, 255, 0.4));
  z-index: 2;
}

.cloud-shadow {
  width: 40px;
  height: 8px;
  background: rgba(0, 0, 0, 0.3);
  border-radius: 50%;
  margin-top: 10px;
  animation: shadowScale 3s ease-in-out infinite;
  filter: blur(2px);
}

.cloud-text {
  font-family: 'Inter', sans-serif;
  font-size: 1rem;
  font-weight: 600;
  color: #a1a1a6;
  letter-spacing: 1px;
  animation: pulseText 1.5s ease-in-out infinite alternate;
}

@keyframes pulseCloud {
  0% {
    transform: scale(0.95);
    filter: drop-shadow(0 0 10px rgba(0, 229, 255, 0.2));
  }
  100% {
    transform: scale(1.05);
    filter: drop-shadow(0 0 25px rgba(0, 229, 255, 0.6));
  }
}

@keyframes floatCloud {
  0%, 100% {
    transform: translateY(0);
  }
  50% {
    transform: translateY(-10px);
  }
}

@keyframes shadowScale {
  0%, 100% {
    transform: scale(1);
    opacity: 0.3;
  }
  50% {
    transform: scale(0.7);
    opacity: 0.1;
  }
}

@keyframes pulseText {
  0% { opacity: 0.6; }
  100% { opacity: 1; }
}

/* Rain Effect */
.rain-drops {
  position: absolute;
  top: 40px;
  display: flex;
  justify-content: center;
  gap: 12px;
  z-index: 1;
}

.drop {
  width: 2px;
  height: 12px;
  background: linear-gradient(to bottom, transparent, rgba(0, 229, 255, 0.8));
  border-radius: 4px;
  opacity: 0;
  animation: rainFall 1.2s infinite ease-in;
}

.drop-1 { animation-delay: 0.1s; height: 10px; margin-top: 5px; }
.drop-2 { animation-delay: 0.5s; height: 14px; }
.drop-3 { animation-delay: 0.3s; height: 12px; margin-top: 2px; }

@keyframes rainFall {
  0% { transform: translateY(0); opacity: 0; }
  20% { opacity: 1; }
  80% { transform: translateY(35px); opacity: 1; height: 4px; }
  100% { transform: translateY(40px) scaleY(0); opacity: 0; }
}
</style>
