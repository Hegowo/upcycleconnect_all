<template>
  <div class="bg-[#f7f9ff] min-h-screen pb-16">
    <div class="max-w-[960px] mx-auto px-4 sm:px-6 py-8">

      <RouterLink to="/projets" class="inline-flex items-center gap-1.5 text-sm text-[#40617f] hover:text-[#006d35] mb-6 transition">
        <ArrowLeftIcon class="w-4 h-4" /> Tous les projets
      </RouterLink>

      <div v-if="loading" class="py-24 flex justify-center">
        <div class="w-8 h-8 border-4 border-[#006d35] border-t-transparent rounded-full animate-spin" />
      </div>

      <div v-else-if="!project" class="bg-white rounded-[24px] border border-[#edf4ff] p-16 text-center">
        <p class="font-jakarta font-bold text-[#001d32] text-lg">Projet introuvable</p>
        <RouterLink to="/projets" class="text-[#006d35] text-sm font-semibold hover:underline mt-2 inline-block">Retour à la galerie</RouterLink>
      </div>

      <template v-else>
        <div class="flex items-center gap-2 mb-3 flex-wrap">
          <span v-if="project.category" class="text-xs font-semibold text-[#006d35] bg-[#f0fdf4] px-2.5 py-1 rounded-full">{{ project.category }}</span>
          <span class="text-xs font-bold px-2.5 py-1 rounded-full" :class="statusBadge(project.status)">{{ statusLabel(project.status) }}</span>
        </div>
        <h1 class="font-jakarta font-extrabold text-[#001d32] text-3xl sm:text-4xl tracking-tight">{{ project.title }}</h1>
        <p class="text-[#40617f] text-sm mt-2">Par {{ project.provider_name || 'Artisan' }}</p>
        <p v-if="project.impact_label" class="inline-flex items-center gap-1.5 text-[#006d35] text-sm font-semibold mt-3 bg-[#f0fdf4] px-3 py-1.5 rounded-xl">
          🌱 {{ project.impact_label }}
        </p>

        <div v-if="images.length" class="mt-6">
          <div class="rounded-[24px] overflow-hidden bg-[#edf4ff] aspect-[16/10]">
            <img :src="images[activeImg]" class="w-full h-full object-cover" />
          </div>
          <div v-if="images.length > 1" class="flex gap-2 mt-3 overflow-x-auto pb-1">
            <button v-for="(img, i) in images" :key="i" @click="activeImg = i"
              class="w-20 h-20 rounded-xl overflow-hidden shrink-0 border-2 transition"
              :class="i === activeImg ? 'border-[#006d35]' : 'border-transparent opacity-70 hover:opacity-100'">
              <img :src="img" class="w-full h-full object-cover" />
            </button>
          </div>
        </div>

        <div v-if="project.description" class="bg-white rounded-[24px] border border-[#edf4ff] p-6 mt-6">
          <h2 class="font-jakarta font-bold text-[#001d32] text-lg mb-2">À propos du projet</h2>
          <p class="text-[#40617f] text-sm leading-relaxed whitespace-pre-line">{{ project.description }}</p>
        </div>

        <div v-if="project.steps && project.steps.length" class="bg-white rounded-[24px] border border-[#edf4ff] p-6 mt-6">
          <h2 class="font-jakarta font-bold text-[#001d32] text-lg mb-4">Étapes de transformation</h2>
          <div class="space-y-4">
            <div v-for="(step, idx) in project.steps" :key="step.id" class="flex gap-4">
              <div class="flex flex-col items-center shrink-0">
                <div class="w-8 h-8 rounded-full flex items-center justify-center text-sm font-bold shrink-0"
                  :class="step.completed_at ? 'bg-[#006d35] text-white' : 'bg-[#e5e7eb] text-[#40617f]'">{{ idx + 1 }}</div>
                <div v-if="idx < project.steps.length - 1" class="w-0.5 flex-1 bg-[#e5e7eb] mt-1" />
              </div>
              <div class="flex-1 min-w-0 pb-2">
                <p class="font-semibold text-[#001d32] text-sm">{{ step.title }}</p>
                <p v-if="step.description" class="text-[#40617f] text-sm mt-0.5 whitespace-pre-line">{{ step.description }}</p>
                <div v-if="stepImages(step).length" class="flex gap-2 mt-2 flex-wrap">
                  <img v-for="(src, i) in stepImages(step)" :key="i" :src="src" class="w-32 h-32 rounded-xl object-cover border border-[#edf4ff]" />
                </div>
              </div>
            </div>
          </div>
        </div>
      </template>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { RouterLink, useRoute } from 'vue-router'
import { ArrowLeftIcon } from '@heroicons/vue/24/outline'
import { publicGet } from '@/services/publicApi'

const route = useRoute()
const project = ref(null)
const loading = ref(true)
const activeImg = ref(0)

const images = computed(() => {
  if (!project.value) return []
  const list = (project.value.images || []).map(i => i.url)
  if (!list.length && project.value.cover_image) return [project.value.cover_image]
  return list
})

function stepImages(step) {
  const list = (step.images || []).map(i => i.url)
  if (!list.length && step.image_url) return [step.image_url]
  return list
}

function statusLabel(s) { return { in_progress: 'En cours', completed: 'Terminé', showcased: 'Mis en avant' }[s] || s }
function statusBadge(s) {
  if (s === 'showcased') return 'bg-[#dcfce7] text-[#166534]'
  if (s === 'completed') return 'bg-blue-100 text-blue-700'
  return 'bg-yellow-100 text-yellow-800'
}

onMounted(async () => {
  try {
    project.value = await publicGet(`/projects/${route.params.id}`)
  } catch {
    project.value = null
  } finally {
    loading.value = false
  }
})
</script>
