<template>
  <div class="relative" ref="root">
    <button
      @click="open = !open"
      class="flex items-center gap-1 pl-2 pr-1.5 py-2 rounded-full hover:bg-[#edf4ff] transition"
      aria-label="Changer de langue"
      aria-haspopup="true"
      :aria-expanded="open"
    >
      <span class="text-xl leading-none">{{ flag(current) }}</span>
      <ChevronDownIcon class="w-3.5 h-3.5 text-[#40617f]" />
    </button>

    <Transition name="lang-pop">
      <div v-if="open" class="absolute right-0 top-12 w-48 bg-white rounded-2xl shadow-xl border border-[#e2e8f0] py-2 z-50">
        <p class="px-4 pb-1.5 text-[10px] font-bold uppercase tracking-wider text-[#94a3b8]">Langue</p>
        <button
          v-for="loc in locales"
          :key="loc.code"
          @click="choose(loc.code)"
          class="w-full flex items-center gap-3 px-4 py-2 text-sm hover:bg-[#f8fafc] transition"
          :class="loc.code === current ? 'text-[#006d35] font-semibold' : 'text-[#334155]'"
        >
          <span class="text-lg leading-none">{{ flag(loc.code) }}</span>
          <span class="flex-1 text-left">{{ loc.name }}</span>
          <CheckIcon v-if="loc.code === current" class="w-4 h-4 text-[#006d35] shrink-0" />
        </button>
      </div>
    </Transition>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ChevronDownIcon, CheckIcon } from '@heroicons/vue/24/outline'
import { availableLocales, setLocale } from '@/utils/i18n.js'

const { locale } = useI18n()
const open = ref(false)
const root = ref(null)
const locales = availableLocales
const current = computed(() => locale.value)

const FLAGS = {
  fr: '🇫🇷', en: '🇬🇧', es: '🇪🇸', de: '🇩🇪', it: '🇮🇹', pt: '🇵🇹', nl: '🇳🇱',
  ar: '🇸🇦', zh: '🇨🇳', ja: '🇯🇵', ru: '🇷🇺', pl: '🇵🇱', tr: '🇹🇷',
}
function flag(code) {
  return FLAGS[(code || '').slice(0, 2).toLowerCase()] || '🌐'
}

async function choose(code) {
  await setLocale(code)
  open.value = false
}

function onDocClick(e) {
  if (open.value && root.value && !root.value.contains(e.target)) open.value = false
}
onMounted(() => document.addEventListener('mousedown', onDocClick))
onUnmounted(() => document.removeEventListener('mousedown', onDocClick))
</script>

<style scoped>
.lang-pop-enter-active, .lang-pop-leave-active { transition: opacity 0.15s ease, transform 0.15s ease; }
.lang-pop-enter-from, .lang-pop-leave-to { opacity: 0; transform: translateY(-4px); }
</style>
