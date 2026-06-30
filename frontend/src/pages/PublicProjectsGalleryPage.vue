<template>
  <div class="bg-[#f7f9ff] min-h-screen pb-16">
    <section class="px-4 sm:px-6 pt-10 sm:pt-14 pb-6 max-w-[1280px] mx-auto">
      <h1 class="font-jakarta font-extrabold text-[#001d32] text-4xl sm:text-5xl tracking-tight">
        Projets d'upcycling
      </h1>
      <p class="text-[#40617f] text-lg mt-3 max-w-2xl leading-relaxed">
        Découvrez les transformations réalisées par nos artisans : objets sauvés, matières détournées et savoir-faire en action.
      </p>
    </section>

    <section class="px-4 sm:px-6 max-w-[1280px] mx-auto">
      <div v-if="categories.length" class="flex flex-wrap gap-2 mb-8">
        <button @click="activeCategory = 'all'"
          class="px-4 py-2 rounded-full text-sm font-semibold border-2 transition-all"
          :class="activeCategory === 'all' ? 'bg-[#006d35] text-white border-[#006d35]' : 'bg-white text-[#40617f] border-[#e5e7eb] hover:border-[#006d35]/30'">
          Toutes
        </button>
        <button v-for="cat in categories" :key="cat" @click="activeCategory = cat"
          class="px-4 py-2 rounded-full text-sm font-semibold border-2 transition-all"
          :class="activeCategory === cat ? 'bg-[#006d35] text-white border-[#006d35]' : 'bg-white text-[#40617f] border-[#e5e7eb] hover:border-[#006d35]/30'">
          {{ cat }}
        </button>
      </div>

      <div v-if="loading" class="py-24 flex justify-center">
        <div class="w-8 h-8 border-4 border-[#006d35] border-t-transparent rounded-full animate-spin" />
      </div>

      <div v-else-if="!filtered.length" class="bg-white rounded-[24px] border border-[#edf4ff] p-16 text-center">
        <SparklesIcon class="w-12 h-12 text-[#64748b]/40 mx-auto mb-3" />
        <p class="font-jakarta font-bold text-[#001d32] text-lg">Aucun projet à découvrir pour le moment</p>
        <p class="text-[#40617f] text-sm mt-1">Revenez bientôt, nos artisans préparent de belles transformations.</p>
      </div>

      <div v-else class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
        <RouterLink v-for="p in filtered" :key="p.id" :to="`/projets/${p.id}`"
          class="group bg-white rounded-[24px] border border-[#edf4ff] overflow-hidden hover:shadow-[0_12px_40px_0_rgba(0,29,50,0.10)] transition">
          <div class="h-48 bg-[#edf4ff] overflow-hidden">
            <img v-if="cover(p)" :src="cover(p)" :alt="p.title" class="w-full h-full object-cover group-hover:scale-105 transition duration-500" />
            <div v-else class="w-full h-full flex items-center justify-center">
              <PhotoIcon class="w-12 h-12 text-[#64748b]/40" />
            </div>
          </div>
          <div class="p-5">
            <div class="flex items-center gap-2 mb-2 flex-wrap">
              <span v-if="p.category" class="text-xs font-semibold text-[#006d35] bg-[#f0fdf4] px-2 py-0.5 rounded-full">{{ p.category }}</span>
              <span v-if="p.status === 'showcased'" class="text-xs font-bold text-[#166534] bg-[#dcfce7] px-2 py-0.5 rounded-full">Mis en avant</span>
            </div>
            <h2 class="font-jakarta font-bold text-[#001d32] text-lg leading-snug line-clamp-2">{{ p.title }}</h2>
            <p class="text-[#40617f] text-sm mt-1 line-clamp-2">{{ p.description }}</p>
            <p v-if="p.impact_label" class="text-[#006d35] text-xs font-semibold mt-2">🌱 {{ p.impact_label }}</p>
            <div class="flex items-center justify-between mt-3 pt-3 border-t border-[#f1f5f9]">
              <p class="text-xs text-[#40617f]">Par {{ p.provider_name || 'Artisan' }}</p>
              <span v-if="p.step_count" class="text-xs text-[#64748b]">{{ p.step_count }} étape{{ p.step_count > 1 ? 's' : '' }}</span>
            </div>
          </div>
        </RouterLink>
      </div>
    </section>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { RouterLink } from 'vue-router'
import { PhotoIcon, SparklesIcon } from '@heroicons/vue/24/outline'
import { publicGet } from '@/services/publicApi'

const projects = ref([])
const loading = ref(true)
const activeCategory = ref('all')

function cover(p) {
  if (p.images && p.images.length) return p.images[0].url
  return p.cover_image || null
}

const categories = computed(() => {
  const set = new Set()
  projects.value.forEach(p => { if (p.category) set.add(p.category) })
  return [...set].sort()
})

const filtered = computed(() =>
  activeCategory.value === 'all'
    ? projects.value
    : projects.value.filter(p => p.category === activeCategory.value)
)

onMounted(async () => {
  try {
    const res = await publicGet('/projects')
    projects.value = res.data || []
  } catch {
    projects.value = []
  } finally {
    loading.value = false
  }
})
</script>
