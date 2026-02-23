<template>
  <!-- FAB Trigger -->
  <div
    class="a11y-fab"
    :class="{ 'a11y-fab--active': state.panelOpen }"
    @click="state.panelOpen = !state.panelOpen"
    role="button"
    aria-label="Abrir painel de acessibilidade"
    :aria-expanded="state.panelOpen"
  >
    <!-- Universal Accessibility Symbol (W3C/WAI) – person with arms outstretched -->
    <svg v-if="!state.panelOpen" xmlns="http://www.w3.org/2000/svg" width="28" height="28" viewBox="0 0 24 24" fill="white" aria-hidden="true">
      <!-- Head -->
      <circle cx="12" cy="3.5" r="2"/>
      <!-- Body -->
      <line x1="12" y1="5.5" x2="12" y2="13" stroke="white" stroke-width="2" stroke-linecap="round"/>
      <!-- Arms outstretched -->
      <line x1="5" y1="9" x2="19" y2="9" stroke="white" stroke-width="2" stroke-linecap="round"/>
      <!-- Left leg -->
      <line x1="12" y1="13" x2="8" y2="21" stroke="white" stroke-width="2" stroke-linecap="round"/>
      <!-- Right leg -->
      <line x1="12" y1="13" x2="16" y2="21" stroke="white" stroke-width="2" stroke-linecap="round"/>
    </svg>
    <v-icon v-else color="white" size="26">mdi-close</v-icon>
    <div class="fab-pulse" v-if="!state.panelOpen"></div>
  </div>

  <!-- Backdrop -->
  <Transition name="a11y-backdrop">
    <div v-if="state.panelOpen" class="a11y-backdrop" @click="state.panelOpen = false"></div>
  </Transition>

  <!-- Panel -->
  <Transition name="a11y-panel">
    <aside v-if="state.panelOpen" class="a11y-panel" role="dialog" aria-label="Painel de Acessibilidade" aria-modal="true">

      <!-- Header -->
      <div class="a11y-header">
        <div class="a11y-header-icon">
          <v-icon color="#00B8D4" size="22">mdi-tune</v-icon>
        </div>
        <div>
          <div class="a11y-header-title">Acessibilidade</div>
          <div class="a11y-header-sub">Personalize sua experiência</div>
        </div>
        <button class="a11y-reset-btn" @click="resetAll" title="Redefinir tudo">
          <v-icon size="16">mdi-refresh</v-icon>
          Redefinir
        </button>
      </div>

      <!-- Divider -->
      <div class="a11y-divider"></div>

      <!-- Section: Visão -->
      <div class="a11y-section-label">Visão</div>
      <div class="a11y-grid">

        <button
          class="a11y-btn"
          :class="{ 'a11y-btn--on': state.highContrast }"
          @click="toggle('highContrast')"
          aria-pressed="state.highContrast"
        >
          <div class="a11y-btn-icon">
            <v-icon size="22" :color="state.highContrast ? '#fff' : '#00B8D4'">mdi-contrast-circle</v-icon>
          </div>
          <span class="a11y-btn-label">Alto Contraste</span>
          <div class="a11y-toggle" :class="{ 'a11y-toggle--on': state.highContrast }"></div>
        </button>

        <button
          class="a11y-btn"
          :class="{ 'a11y-btn--on': state.largeText }"
          @click="toggle('largeText')"
          aria-pressed="state.largeText"
        >
          <div class="a11y-btn-icon">
            <v-icon size="22" :color="state.largeText ? '#fff' : '#00B8D4'">mdi-format-size</v-icon>
          </div>
          <span class="a11y-btn-label">Texto Grande</span>
          <div class="a11y-toggle" :class="{ 'a11y-toggle--on': state.largeText }"></div>
        </button>

        <button
          class="a11y-btn"
          :class="{ 'a11y-btn--on': state.extraLargeText }"
          @click="toggle('extraLargeText')"
          aria-pressed="state.extraLargeText"
        >
          <div class="a11y-btn-icon">
            <v-icon size="22" :color="state.extraLargeText ? '#fff' : '#00B8D4'">mdi-format-header-1</v-icon>
          </div>
          <span class="a11y-btn-label">Texto Extra Grande</span>
          <div class="a11y-toggle" :class="{ 'a11y-toggle--on': state.extraLargeText }"></div>
        </button>

        <button
          class="a11y-btn"
          :class="{ 'a11y-btn--on': state.dyslexiaFont }"
          @click="toggle('dyslexiaFont')"
          aria-pressed="state.dyslexiaFont"
        >
          <div class="a11y-btn-icon">
            <v-icon size="22" :color="state.dyslexiaFont ? '#fff' : '#00B8D4'">mdi-format-font</v-icon>
          </div>
          <span class="a11y-btn-label">Fonte Dislexia</span>
          <div class="a11y-toggle" :class="{ 'a11y-toggle--on': state.dyslexiaFont }"></div>
        </button>

        <button
          class="a11y-btn"
          :class="{ 'a11y-btn--on': state.letterSpacing }"
          @click="toggle('letterSpacing')"
          aria-pressed="state.letterSpacing"
        >
          <div class="a11y-btn-icon">
            <v-icon size="22" :color="state.letterSpacing ? '#fff' : '#00B8D4'">mdi-format-pilcrow</v-icon>
          </div>
          <span class="a11y-btn-label">Espaço Entre Letras</span>
          <div class="a11y-toggle" :class="{ 'a11y-toggle--on': state.letterSpacing }"></div>
        </button>

        <button
          class="a11y-btn"
          :class="{ 'a11y-btn--on': state.lineSpacing }"
          @click="toggle('lineSpacing')"
          aria-pressed="state.lineSpacing"
        >
          <div class="a11y-btn-icon">
            <v-icon size="22" :color="state.lineSpacing ? '#fff' : '#00B8D4'">mdi-format-line-spacing</v-icon>
          </div>
          <span class="a11y-btn-label">Espaço Entre Linhas</span>
          <div class="a11y-toggle" :class="{ 'a11y-toggle--on': state.lineSpacing }"></div>
        </button>

        <button
          class="a11y-btn"
          :class="{ 'a11y-btn--on': state.highlightLinks }"
          @click="toggle('highlightLinks')"
          aria-pressed="state.highlightLinks"
        >
          <div class="a11y-btn-icon">
            <v-icon size="22" :color="state.highlightLinks ? '#fff' : '#00B8D4'">mdi-link-variant</v-icon>
          </div>
          <span class="a11y-btn-label">Destacar Links</span>
          <div class="a11y-toggle" :class="{ 'a11y-toggle--on': state.highlightLinks }"></div>
        </button>

        <button
          class="a11y-btn"
          :class="{ 'a11y-btn--on': state.largeCursor }"
          @click="toggle('largeCursor')"
          aria-pressed="state.largeCursor"
        >
          <div class="a11y-btn-icon">
            <v-icon size="22" :color="state.largeCursor ? '#fff' : '#00B8D4'">mdi-cursor-default-outline</v-icon>
          </div>
          <span class="a11y-btn-label">Cursor Grande</span>
          <div class="a11y-toggle" :class="{ 'a11y-toggle--on': state.largeCursor }"></div>
        </button>

      </div>

      <!-- Divider -->
      <div class="a11y-divider"></div>

      <!-- Section: Movimento & Leitura -->
      <div class="a11y-section-label">Movimento & Leitura</div>
      <div class="a11y-grid">

        <button
          class="a11y-btn"
          :class="{ 'a11y-btn--on': state.reduceMotion }"
          @click="toggle('reduceMotion')"
          aria-pressed="state.reduceMotion"
        >
          <div class="a11y-btn-icon">
            <v-icon size="22" :color="state.reduceMotion ? '#fff' : '#00B8D4'">mdi-motion-pause</v-icon>
          </div>
          <span class="a11y-btn-label">Reduzir Animações</span>
          <div class="a11y-toggle" :class="{ 'a11y-toggle--on': state.reduceMotion }"></div>
        </button>

        <button
          class="a11y-btn"
          :class="{ 'a11y-btn--on': state.readingGuide }"
          @click="toggle('readingGuide')"
          aria-pressed="state.readingGuide"
        >
          <div class="a11y-btn-icon">
            <v-icon size="22" :color="state.readingGuide ? '#fff' : '#00B8D4'">mdi-ray-start-arrow</v-icon>
          </div>
          <span class="a11y-btn-label">Guia de Leitura</span>
          <div class="a11y-toggle" :class="{ 'a11y-toggle--on': state.readingGuide }"></div>
        </button>

      </div>

      <!-- Divider -->
      <div class="a11y-divider"></div>

      <!-- Section: Áudio -->
      <div class="a11y-section-label">Áudio</div>
      <div class="a11y-grid">

        <button
          class="a11y-btn a11y-btn--wide"
          :class="{ 'a11y-btn--on': state.isSpeaking, 'a11y-btn--speaking': state.isSpeaking }"
          @click="speakPage"
          :aria-label="state.isSpeaking ? 'Parar leitura' : 'Ler página em voz alta'"
        >
          <div class="a11y-btn-icon">
            <v-icon size="22" :color="state.isSpeaking ? '#fff' : '#00B8D4'">
              {{ state.isSpeaking ? 'mdi-stop-circle' : 'mdi-text-to-speech' }}
            </v-icon>
          </div>
          <span class="a11y-btn-label">
            {{ state.isSpeaking ? 'Parar Leitura em Voz Alta' : 'Ler Página em Voz Alta' }}
          </span>
          <span v-if="state.isSpeaking" class="a11y-speaking-dot"></span>
        </button>

      </div>

      <!-- Footer -->
      <div class="a11y-footer">
        <v-icon size="12" color="rgba(255,255,255,0.3)" class="mr-1">mdi-shield-check</v-icon>
        As preferências são salvas automaticamente
      </div>
    </aside>
  </Transition>
</template>

<script setup>
import { useAccessibility } from '@/composables/useAccessibility'

const { state, toggle, resetAll, speakPage } = useAccessibility()
</script>

<style>
/* ===========================
   ACCESSIBILITY CSS CLASSES
   Applied to <html> element
=========================== */

/* High Contrast */
html.a11y-high-contrast {
  filter: contrast(170%) brightness(1.05) !important;
}
html.a11y-high-contrast img {
  filter: contrast(130%);
}

/* Large Text */
html.a11y-large-text {
  font-size: 120% !important;
}
html.a11y-extra-large-text {
  font-size: 150% !important;
}

/* Dyslexia Font */
@import url('https://fonts.cdnfonts.com/css/open-dyslexic');
html.a11y-dyslexia,
html.a11y-dyslexia * {
  font-family: 'OpenDyslexic', Arial, sans-serif !important;
}

/* Reduce Motion */
html.a11y-reduce-motion *,
html.a11y-reduce-motion *::before,
html.a11y-reduce-motion *::after {
  animation-duration: 0.001ms !important;
  animation-iteration-count: 1 !important;
  transition-duration: 0.001ms !important;
}

/* Highlight Links */
html.a11y-highlight-links a,
html.a11y-highlight-links [role="link"],
html.a11y-highlight-links .router-link-active {
  outline: 2px solid #FFD600 !important;
  background: rgba(255, 214, 0, 0.15) !important;
  border-radius: 4px !important;
  padding: 0 2px !important;
  text-decoration: underline !important;
}

/* Large Cursor */
html.a11y-large-cursor,
html.a11y-large-cursor * {
  cursor: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='40' height='40' viewBox='0 0 40 40'%3E%3Cpath d='M5 5 L5 30 L12 23 L18 35 L22 33 L16 21 L25 21 Z' fill='black' stroke='white' stroke-width='2'/%3E%3C/svg%3E") 5 5, auto !important;
}

/* Letter Spacing */
html.a11y-letter-spacing,
html.a11y-letter-spacing p,
html.a11y-letter-spacing span,
html.a11y-letter-spacing div,
html.a11y-letter-spacing li {
  letter-spacing: 0.12em !important;
  word-spacing: 0.16em !important;
}

/* Line Spacing */
html.a11y-line-spacing,
html.a11y-line-spacing p,
html.a11y-line-spacing div,
html.a11y-line-spacing li,
html.a11y-line-spacing span {
  line-height: 2 !important;
}

/* Reading Guide Bar */
#a11y-reading-guide-bar {
  position: fixed;
  left: 0;
  width: 100%;
  height: 60px;
  background: rgba(255, 214, 0, 0.08);
  border-top: 2px solid rgba(255, 214, 0, 0.5);
  border-bottom: 2px solid rgba(255, 214, 0, 0.5);
  pointer-events: none;
  z-index: 9990;
  transform: translateY(-50%);
  display: none;
}
html.a11y-reading-guide #a11y-reading-guide-bar {
  display: block;
}
</style>

<style scoped>
/* ===========================
   FAB BUTTON
=========================== */
.a11y-fab {
  position: fixed;
  right: 36px;
  top: 50%;
  transform: translateY(-50%);
  background: linear-gradient(135deg, #007AFF, #00B8D4);
  width: 56px;
  height: 56px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  z-index: 3000;
  transition: all 0.3s cubic-bezier(0.175, 0.885, 0.32, 1.275);
  box-shadow: 0 8px 24px rgba(0,122,255,0.45);
  user-select: none;
}
.a11y-fab:hover {
  transform: translateY(-50%) scale(1.12);
  box-shadow: 0 12px 36px rgba(0,122,255,0.6);
}
.a11y-fab--active {
  background: linear-gradient(135deg, #FF3B30, #FF6B6B);
  transform: translateY(-50%) rotate(15deg);
  box-shadow: 0 8px 24px rgba(255,59,48,0.4);
}
.a11y-fab--active:hover {
  transform: translateY(-50%) rotate(15deg) scale(1.1);
}

.fab-pulse {
  position: absolute;
  width: 100%;
  height: 100%;
  border-radius: 50%;
  background: #007AFF;
  opacity: 0.3;
  animation: pulse-ring 2.2s infinite;
}
@keyframes pulse-ring {
  0%   { transform: scale(1);   opacity: 0.3; }
  100% { transform: scale(1.8); opacity: 0; }
}

/* ===========================
   BACKDROP
=========================== */
.a11y-backdrop {
  position: fixed;
  inset: 0;
  background: rgba(0,0,0,0.5);
  z-index: 2900;
  backdrop-filter: blur(4px);
}
.a11y-backdrop-enter-active,
.a11y-backdrop-leave-active { transition: opacity 0.25s ease; }
.a11y-backdrop-enter-from,
.a11y-backdrop-leave-to { opacity: 0; }

/* ===========================
   PANEL
=========================== */
.a11y-panel {
  position: fixed;
  right: 104px;
  top: 50%;
  transform: translateY(-50%);
  width: 320px;
  max-height: 85vh;
  overflow-y: auto;
  background: rgba(18, 38, 64, 0.98);
  border: 1px solid rgba(255,255,255,0.12);
  border-radius: 28px;
  padding: 24px 20px;
  z-index: 2950;
  box-shadow: 0 24px 80px rgba(0,0,0,0.6), 0 0 0 1px rgba(0,184,212,0.1);
  scrollbar-width: thin;
  scrollbar-color: rgba(255,255,255,0.15) transparent;
}
.a11y-panel::-webkit-scrollbar { width: 4px; }
.a11y-panel::-webkit-scrollbar-track { background: transparent; }
.a11y-panel::-webkit-scrollbar-thumb { background: rgba(255,255,255,0.15); border-radius: 2px; }

.a11y-panel-enter-active,
.a11y-panel-leave-active { transition: opacity 0.25s ease, transform 0.25s cubic-bezier(0.23,1,0.32,1); }
.a11y-panel-enter-from,
.a11y-panel-leave-to { opacity: 0; transform: translateY(-50%) translateX(20px); }

/* Header */
.a11y-header {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 16px;
}
.a11y-header-icon {
  width: 40px; height: 40px;
  background: rgba(0,184,212,0.12);
  border-radius: 12px;
  display: flex; align-items: center; justify-content: center;
  flex-shrink: 0;
}
.a11y-header-title {
  font-size: 1rem;
  font-weight: 800;
  color: white;
  line-height: 1;
}
.a11y-header-sub {
  font-size: 0.72rem;
  color: rgba(255,255,255,0.45);
  margin-top: 3px;
}
.a11y-reset-btn {
  margin-left: auto;
  display: inline-flex;
  align-items: center;
  gap: 4px;
  background: rgba(255,255,255,0.07);
  border: 1px solid rgba(255,255,255,0.1);
  border-radius: 8px;
  padding: 5px 10px;
  font-size: 11px;
  color: rgba(255,255,255,0.6);
  cursor: pointer;
  transition: all 0.2s;
  flex-shrink: 0;
}
.a11y-reset-btn:hover {
  background: rgba(255,59,48,0.15);
  border-color: rgba(255,59,48,0.3);
  color: #FF6B6B;
}

/* Divider */
.a11y-divider {
  height: 1px;
  background: rgba(255,255,255,0.07);
  margin: 14px 0 12px;
}

/* Section label */
.a11y-section-label {
  font-size: 0.65rem;
  font-weight: 700;
  color: rgba(255,255,255,0.35);
  text-transform: uppercase;
  letter-spacing: 1px;
  margin-bottom: 10px;
}

/* Grid */
.a11y-grid {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

/* Feature buttons */
.a11y-btn {
  display: flex;
  align-items: center;
  gap: 10px;
  width: 100%;
  background: rgba(255,255,255,0.05);
  border: 1px solid rgba(255,255,255,0.08);
  border-radius: 14px;
  padding: 11px 14px;
  cursor: pointer;
  transition: all 0.2s ease;
  text-align: left;
  color: rgba(255,255,255,0.8);
}
.a11y-btn:hover {
  background: rgba(255,255,255,0.09);
  border-color: rgba(0,184,212,0.25);
}
.a11y-btn--on {
  background: linear-gradient(135deg, rgba(0,184,212,0.25), rgba(0,122,255,0.2)) !important;
  border-color: rgba(0,184,212,0.4) !important;
  color: white !important;
}
.a11y-btn--speaking {
  background: linear-gradient(135deg, rgba(39,174,96,0.25), rgba(0,200,100,0.15)) !important;
  border-color: rgba(39,210,96,0.4) !important;
}
.a11y-btn--wide { flex: 1; }

.a11y-btn-icon {
  width: 34px; height: 34px;
  background: rgba(255,255,255,0.06);
  border-radius: 10px;
  display: flex; align-items: center; justify-content: center;
  flex-shrink: 0;
}
.a11y-btn--on .a11y-btn-icon {
  background: rgba(0,184,212,0.2);
}

.a11y-btn-label {
  font-size: 0.83rem;
  font-weight: 600;
  flex: 1;
}

/* Toggle switch */
.a11y-toggle {
  width: 34px; height: 20px;
  background: rgba(255,255,255,0.15);
  border-radius: 10px;
  position: relative;
  transition: background 0.2s;
  flex-shrink: 0;
}
.a11y-toggle::after {
  content: '';
  position: absolute;
  top: 3px; left: 3px;
  width: 14px; height: 14px;
  background: white;
  border-radius: 50%;
  transition: transform 0.2s cubic-bezier(0.34,1.56,0.64,1);
  box-shadow: 0 1px 4px rgba(0,0,0,0.3);
}
.a11y-toggle--on {
  background: #00B8D4;
}
.a11y-toggle--on::after {
  transform: translateX(14px);
}

/* Speaking animation */
.a11y-speaking-dot {
  width: 8px; height: 8px;
  background: #27AE60;
  border-radius: 50%;
  animation: blink 1s ease infinite;
  flex-shrink: 0;
}
@keyframes blink {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.2; }
}

/* Footer */
.a11y-footer {
  margin-top: 16px;
  font-size: 0.68rem;
  color: rgba(255,255,255,0.25);
  text-align: center;
  display: flex;
  align-items: center;
  justify-content: center;
}

/* Responsive */
@media (max-width: 600px) {
  .a11y-panel {
    right: 8px;
    left: 8px;
    width: auto;
    top: auto;
    bottom: 100px;
    transform: none;
    max-height: 70vh;
  }
  .a11y-panel-enter-from,
  .a11y-panel-leave-to { transform: translateY(20px); }
  .a11y-fab { right: 16px; top: auto; bottom: 28px; transform: none; }
  .a11y-fab:hover { transform: scale(1.1); }
  .a11y-fab--active { transform: rotate(15deg); }
}
</style>
