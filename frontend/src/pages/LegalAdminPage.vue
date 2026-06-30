<template>
  <div class="max-w-4xl space-y-6">
    <div>
      <h2 class="text-xl sm:text-2xl font-bold text-[#001d32]">Documents légaux</h2>
      <p class="text-sm text-[#40617f] mt-0.5">Modifiez à tout moment les mentions légales, les CGU/CGV et la politique de confidentialité.</p>
    </div>

    <div v-if="loading" class="py-16 text-center">
      <div class="w-8 h-8 border-4 border-[#006d35] border-t-transparent rounded-full animate-spin mx-auto" />
    </div>

    <div v-else class="card overflow-hidden">
      <div class="flex gap-1 p-3 bg-[#f8fafc] border-b border-[#f1f5f9] overflow-x-auto">
        <button
          v-for="d in docs"
          :key="d.slug"
          @click="select(d.slug)"
          :class="active === d.slug ? 'bg-white shadow text-[#001d32] font-semibold' : 'text-gray-500 hover:text-gray-700'"
          class="px-4 py-2 rounded-lg text-sm transition whitespace-nowrap"
        >
          {{ tabLabel(d.slug) }}
        </button>
      </div>

      <div v-if="current" class="p-5 space-y-4">
        <div>
          <label class="block text-xs font-semibold text-[#40617f] uppercase mb-1">Titre</label>
          <input v-model="current.title" type="text" class="w-full px-3 py-2.5 bg-[#f8fafc] border border-[#e5e7eb] rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-[#006d35]/30" />
        </div>
        <div>
          <label class="block text-xs font-semibold text-[#40617f] uppercase mb-1">Contenu</label>
          <TiptapEditor v-model="current.content" min-height="320px" placeholder="Rédigez le document..." />
        </div>
        <p v-if="current.updated_at" class="text-xs text-[#94a3b8]">Dernière mise à jour : {{ formatDate(current.updated_at) }}</p>

        <div class="flex items-center justify-end gap-3">
          <span v-if="savedMsg" class="text-sm text-[#16a34a] font-medium">{{ savedMsg }}</span>
          <button @click="save" :disabled="saving" class="px-6 py-2.5 rounded-xl text-white font-bold text-sm hover:opacity-90 disabled:opacity-50 transition" style="background:linear-gradient(135deg,#006d35,#1b8848);">
            {{ saving ? 'Enregistrement...' : 'Enregistrer' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import api from '@/services/api'
import TiptapEditor from '@/components/TiptapEditor.vue'

const docs = ref([])
const active = ref('')
const loading = ref(true)
const saving = ref(false)
const savedMsg = ref('')

const LABELS = {
  'mentions-legales': 'Mentions légales',
  'cgu-cgv': 'CGU / CGV',
  'confidentialite': 'Confidentialité',
}
function tabLabel(slug) { return LABELS[slug] || slug }

const current = computed(() => docs.value.find(d => d.slug === active.value))

function formatDate(iso) {
  return new Date(iso).toLocaleDateString('fr-FR', { day: 'numeric', month: 'long', year: 'numeric', hour: '2-digit', minute: '2-digit' })
}

function select(slug) { active.value = slug; savedMsg.value = '' }

async function load() {
  loading.value = true
  try {
    const { data } = await api.get('/legal')
    docs.value = data.data || []
    if (docs.value.length) active.value = docs.value[0].slug
  } finally {
    loading.value = false
  }
}

async function save() {
  if (!current.value) return
  saving.value = true
  savedMsg.value = ''
  try {
    const { data } = await api.put(`/legal/${current.value.slug}`, {
      title: current.value.title,
      content: current.value.content,
    })
    const idx = docs.value.findIndex(d => d.slug === current.value.slug)
    if (idx !== -1) docs.value[idx] = data.data
    savedMsg.value = 'Enregistré ✓'
    setTimeout(() => { savedMsg.value = '' }, 3000)
  } finally {
    saving.value = false
  }
}

onMounted(load)
</script>
