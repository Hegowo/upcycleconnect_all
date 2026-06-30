<template>
  <div class="max-w-3xl mx-auto px-6 py-12">
    <div v-if="loading" class="py-20 text-center">
      <div class="w-8 h-8 border-4 border-[#006d35] border-t-transparent rounded-full animate-spin mx-auto" />
    </div>
    <div v-else-if="error" class="py-20 text-center text-[#64748b]">
      {{ error }}
    </div>
    <div v-else>
      <h1 class="text-3xl font-bold text-[#001d32] mb-2">{{ doc.title }}</h1>
      <p v-if="doc.updated_at" class="text-sm text-[#64748b] mb-8">
        Dernière mise à jour : {{ formatDate(doc.updated_at) }}
      </p>
      <article class="legal-content" v-html="doc.content" />
    </div>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { publicGet } from '@/services/publicApi'

const route = useRoute()
const doc = ref({ title: '', content: '', updated_at: '' })
const loading = ref(true)
const error = ref('')

function formatDate(iso) {
  try {
    return new Date(iso).toLocaleDateString('fr-FR', { day: 'numeric', month: 'long', year: 'numeric' })
  } catch {
    return ''
  }
}

async function load(slug) {
  loading.value = true
  error.value = ''
  try {
    const { data } = await publicGet(`/legal/${slug}`)
    doc.value = data
  } catch {
    error.value = 'Ce document n\'est pas disponible pour le moment.'
  } finally {
    loading.value = false
  }
}

watch(() => route.meta.slug, (slug) => { if (slug) load(slug) }, { immediate: true })
</script>

<style scoped>
.legal-content :deep(h1),
.legal-content :deep(h2) { font-size: 1.4rem; font-weight: 700; color: #001d32; margin: 1.5rem 0 0.75rem; }
.legal-content :deep(h3) { font-size: 1.15rem; font-weight: 700; color: #001d32; margin: 1.25rem 0 0.5rem; }
.legal-content :deep(p) { color: #334155; line-height: 1.7; margin-bottom: 0.85rem; }
.legal-content :deep(ul),
.legal-content :deep(ol) { color: #334155; line-height: 1.7; margin: 0 0 0.85rem 1.5rem; }
.legal-content :deep(ul) { list-style: disc; }
.legal-content :deep(ol) { list-style: decimal; }
.legal-content :deep(a) { color: #006d35; text-decoration: underline; }
.legal-content :deep(strong) { font-weight: 700; color: #001d32; }
</style>
