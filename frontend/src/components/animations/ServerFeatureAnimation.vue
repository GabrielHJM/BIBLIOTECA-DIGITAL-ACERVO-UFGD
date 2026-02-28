<template>
  <div class="server-animation-container" :class="type">
    <div class="server-core">
      <v-icon size="24" class="core-icon">{{ serverIcon }}</v-icon>
      <div class="glass-orb"></div>
      <div class="core-glow"></div>
    </div>

    <!-- Infinity Data Flow (iOS Fluid style) -->
    <div v-if="type === 'infinity'" class="infinity-flow">
      <svg viewBox="0 0 100 60" class="flow-svg">
        <path id="infPath" d="M 25,30 C 25,10 45,10 50,30 C 55,50 75,50 75,30 C 75,10 55,10 50,30 C 45,50 25,50 25,30" fill="none" />
        <circle r="2" fill="currentColor">
          <animateMotion dur="3s" repeatCount="indefinite" path="M 25,30 C 25,10 45,10 50,30 C 55,50 75,50 75,30 C 75,10 55,10 50,30 C 45,50 25,50 25,30" />
        </circle>
        <circle r="1.5" fill="currentColor" opacity="0.6">
          <animateMotion dur="3s" begin="1s" repeatCount="indefinite" path="M 25,30 C 25,10 45,10 50,30 C 55,50 75,50 75,30 C 75,10 55,10 50,30 C 45,50 25,50 25,30" />
        </circle>
        <circle r="1" fill="currentColor" opacity="0.3">
          <animateMotion dur="3s" begin="2s" repeatCount="indefinite" path="M 25,30 C 25,10 45,10 50,30 C 55,50 75,50 75,30 C 75,10 55,10 50,30 C 45,50 25,50 25,30" />
        </circle>
      </svg>
    </div>

    <!-- Soft Sonar Ping (iOS style) -->
    <div v-if="type === 'timer'" class="sonar-pings">
      <div v-for="i in 3" :key="i" class="sonar-ring" :style="`--delay: ${i * 0.6}s`"></div >
    </div>

    <!-- Floating Glass Cards (iOS style) -->
    <div v-if="type === 'personalization'" class="glass-cards">
      <div v-for="i in 3" :key="i" class="glass-card" :style="`--delay: ${i * 0.4}s; --offset: ${i * 10}px`"></div >
    </div>

    <!-- Living Shield (iOS Shield Fix) -->
    <div v-if="type === 'security'" class="living-shield">
      <div class="shield-sphere"></div>
      <div class="shield-sweep"></div>
      <div class="shield-nodes">
        <div v-for="i in 4" :key="i" class="shield-node" :style="`--rot: ${i * 90}deg`"></div>
      </div>
    </div>

    <!-- Organic Network (iOS Connection style) -->
    <div v-if="type === 'global'" class="organic-network">
      <div v-for="i in 6" :key="i" class="network-branch" :style="`--rot: ${i * 60}deg; --len: ${30 + (i % 2) * 10}px`"></div >
      <div v-for="i in 3" :key="'p'+i" class="network-pulse" :style="`--delay: ${i * 0.8}s`"></div >
    </div>
  </div>
</template>

<script>
export default {
  name: 'ServerFeatureAnimation',
  props: {
    type: {
      type: String,
      required: true
    },
    icon: {
      type: String,
      default: 'mdi-server'
    }
  },
  computed: {
    serverIcon() {
      return this.icon;
    }
  }
}
</script>

<style scoped>
.server-animation-container {
  position: relative;
  width: 50px;
  height: 50px;
  display: flex;
  align-items: center;
  justify-content: center;
  perspective: 1000px;
  will-change: transform;
}

.server-core {
  position: relative;
  z-index: 10;
  display: flex;
  align-items: center;
  justify-content: center;
}

.core-icon {
  position: relative;
  z-index: 2;
  filter: drop-shadow(0 0 8px currentColor);
}

.glass-orb {
  position: absolute;
  width: 32px;
  height: 32px;
  background: rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(4px);
  border-radius: 50%;
  border: 1px solid rgba(255, 255, 255, 0.2);
  box-shadow: inset 0 0 10px rgba(255, 255, 255, 0.2);
  will-change: transform, opacity;
}

.core-glow {
  position: absolute;
  width: 40px;
  height: 40px;
  background: radial-gradient(circle, currentColor 0%, transparent 70%);
  opacity: 0.2;
  animation: glowPulse 3s infinite ease-in-out;
}

@keyframes glowPulse {
  0%, 100% { transform: scale(1); opacity: 0.15; }
  50% { transform: scale(1.3); opacity: 0.3; }
}

/* iOS Infinity Flow */
.infinity-flow {
  position: absolute;
  width: 80px;
  height: 50px;
  color: #00F2FE;
}
.flow-svg {
  width: 100%;
  height: 100%;
  filter: drop-shadow(0 0 4px currentColor);
}

/* iOS Sonar Pings */
.sonar-pings {
  position: absolute;
  width: 100%;
  height: 100%;
}
.sonar-ring {
  position: absolute;
  width: 100%;
  height: 100%;
  border: 1.5px solid currentColor;
  border-radius: 50%;
  top: 0; left: 0;
  animation: sonarRipple 2s infinite cubic-bezier(0.25, 0.46, 0.45, 0.94);
  animation-delay: var(--delay);
  opacity: 0;
  will-change: transform, opacity;
}
@keyframes sonarRipple {
  0% { transform: scale(0.5); opacity: 0; }
  20% { opacity: 0.5; }
  100% { transform: scale(2.2); opacity: 0; }
}

/* iOS Glass Cards */
.glass-cards {
  position: absolute;
  width: 40px;
  height: 40px;
}
.glass-card {
  position: absolute;
  width: 20px;
  height: 14px;
  background: rgba(255, 255, 255, 0.15);
  backdrop-filter: blur(2px);
  border: 1px solid rgba(255, 255, 255, 0.3);
  border-radius: 4px;
  left: 50%; top: 50%;
  margin-left: -10px; margin-top: -7px;
  animation: cardFloat 3s infinite ease-in-out;
  animation-delay: var(--delay);
  will-change: transform, opacity;
}
@keyframes cardFloat {
  0%, 100% { transform: translate(var(--offset), -15px) rotate(5deg); opacity: 0.4; }
  50% { transform: translate(calc(var(--offset) * -1), 15px) rotate(-5deg); opacity: 0.8; }
}

/* iOS Living Shield (Fix) */
.living-shield {
  position: absolute;
  width: 54px;
  height: 54px;
  color: #FF6B6B;
}
.shield-sphere {
  position: absolute;
  width: 100%;
  height: 100%;
  border: 1px solid rgba(255, 107, 107, 0.3);
  border-radius: 50%;
  background: radial-gradient(circle at 30% 30%, rgba(255, 255, 255, 0.1), transparent);
}
.shield-sweep {
  position: absolute;
  width: 100%;
  height: 100%;
  border-radius: 50%;
  border-top: 2px solid currentColor;
  animation: rotateShield 2s infinite cubic-bezier(0.4, 0, 0.2, 1);
}
.shield-node {
  position: absolute;
  width: 4px;
  height: 4px;
  background: white;
  border-radius: 50%;
  top: 50%; left: 50%;
  transform: rotate(var(--rot)) translateY(-27px);
  box-shadow: 0 0 8px currentColor;
}
@keyframes rotateShield {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

/* iOS Organic Network */
.organic-network {
  position: absolute;
  width: 100%;
  height: 100%;
  color: #FFD60A;
}
.network-branch {
  position: absolute;
  width: var(--len);
  height: 1px;
  background: linear-gradient(90deg, currentColor, transparent);
  top: 50%; left: 50%;
  transform-origin: left center;
  transform: rotate(var(--rot));
  opacity: 0.4;
}
.network-pulse {
  position: absolute;
  width: 6px;
  height: 6px;
  background: white;
  border-radius: 50%;
  top: 50%; left: 50%;
  margin-left: -3px; margin-top: -3px;
  box-shadow: 0 0 10px currentColor;
  animation: networkPulse 2.4s infinite ease-out;
  animation-delay: var(--delay);
}
@keyframes networkPulse {
  0% { transform: scale(0) translate(0); opacity: 1; }
  20% { opacity: 1; }
  100% { transform: scale(1.2) translate(40px, 25px); opacity: 0; }
}
</style>
