import { reactive } from 'vue'

// ─── State global singleton ───────────────────────────────────────────────────
const state = reactive({
	highContrast: false,
	largeText: false,   // 120% font scale
	extraLargeText: false,   // 150% font scale
	dyslexiaFont: false,   // OpenDyslexic
	reduceMotion: false,
	highlightLinks: false,
	largeCursor: false,
	letterSpacing: false,   // extra letter spacing
	lineSpacing: false,   // extra line height
	readingGuide: false,   // horizontal bar tracking mouse
	isSpeaking: false,
	panelOpen: false,
})

const STORAGE_KEY = 'a11y-prefs'

// ─── Apply all active classes to <html> ──────────────────────────────────────
function applyClasses() {
	const el = document.documentElement
	el.classList.toggle('a11y-high-contrast', state.highContrast)
	el.classList.toggle('a11y-large-text', state.largeText && !state.extraLargeText)
	el.classList.toggle('a11y-extra-large-text', state.extraLargeText)
	el.classList.toggle('a11y-dyslexia', state.dyslexiaFont)
	el.classList.toggle('a11y-reduce-motion', state.reduceMotion)
	el.classList.toggle('a11y-highlight-links', state.highlightLinks)
	el.classList.toggle('a11y-large-cursor', state.largeCursor)
	el.classList.toggle('a11y-letter-spacing', state.letterSpacing)
	el.classList.toggle('a11y-line-spacing', state.lineSpacing)
	el.classList.toggle('a11y-reading-guide', state.readingGuide)
}

// ─── Persist prefs ────────────────────────────────────────────────────────────
function save() {
	// eslint-disable-next-line no-unused-vars
	const { panelOpen: _p, isSpeaking: _s, ...prefs } = state
	localStorage.setItem(STORAGE_KEY, JSON.stringify(prefs))
}

function load() {
	try {
		const saved = JSON.parse(localStorage.getItem(STORAGE_KEY) || '{}')
		Object.assign(state, saved)
	} catch (_) { /* ignore parse errors */ }
}

// ─── Text-to-Speech ──────────────────────────────────────────────────────────
let utterance = null

function speakPage() {
	if (!window.speechSynthesis) return
	if (state.isSpeaking) {
		stopSpeaking()
		return
	}
	const text = document.body.innerText
		.replace(/\s+/g, ' ')
		.trim()
		.substring(0, 8000)

	utterance = new SpeechSynthesisUtterance(text)
	utterance.lang = 'pt-BR'
	utterance.rate = 0.95
	utterance.pitch = 1
	utterance.onend = () => { state.isSpeaking = false }
	utterance.onerror = () => { state.isSpeaking = false }

	// Try to use a Portuguese voice
	const voices = window.speechSynthesis.getVoices()
	const ptVoice = voices.find(v => v.lang.startsWith('pt'))
	if (ptVoice) utterance.voice = ptVoice

	state.isSpeaking = true
	window.speechSynthesis.speak(utterance)
}

function stopSpeaking() {
	if (window.speechSynthesis) window.speechSynthesis.cancel()
	state.isSpeaking = false
}

// ─── Reading guide (follows mouse vertically) ─────────────────────────────────
let guideEl = null

function updateGuidePosition(e) {
	if (guideEl) guideEl.style.top = e.clientY + 'px'
}

function initReadingGuide() {
	if (guideEl) return
	guideEl = document.createElement('div')
	guideEl.id = 'a11y-reading-guide-bar'
	document.body.appendChild(guideEl)
	document.addEventListener('mousemove', updateGuidePosition)
}

function cleanupReadingGuide() {
	if (guideEl) {
		document.removeEventListener('mousemove', updateGuidePosition)
		guideEl.remove()
		guideEl = null
	}
}

// ─── Toggle helpers ───────────────────────────────────────────────────────────
function toggle(key) {
	state[key] = !state[key]
	// Mutual exclusion for text sizes
	if (key === 'largeText' && state.largeText) state.extraLargeText = false
	if (key === 'extraLargeText' && state.extraLargeText) state.largeText = false
	applyClasses()
	save()
}

function resetAll() {
	stopSpeaking()
	Object.assign(state, {
		highContrast: false,
		largeText: false,
		extraLargeText: false,
		dyslexiaFont: false,
		reduceMotion: false,
		highlightLinks: false,
		largeCursor: false,
		letterSpacing: false,
		lineSpacing: false,
		readingGuide: false,
		isSpeaking: false,
	})
	applyClasses()
	save()
}

// ─── Composable export ────────────────────────────────────────────────────────
export function useAccessibility() {
	// Called once from App.vue mounted
	function init() {
		load()
		applyClasses()
		initReadingGuide()
		// Ensure voices are loaded (async in some browsers)
		if (window.speechSynthesis) {
			window.speechSynthesis.onvoiceschanged = () => { /* trigger voices load */ }
		}
	}

	function cleanup() {
		cleanupReadingGuide()
		stopSpeaking()
	}

	return { state, toggle, resetAll, speakPage, stopSpeaking, init, cleanup }
}
