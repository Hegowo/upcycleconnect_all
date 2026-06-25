<template>
  <div class="bg-[#f7f9ff] min-h-screen pb-16">
    <div class="max-w-[960px] mx-auto px-4 sm:px-6 py-10">

      <div class="mb-8">
        <button @click="router.push('/profil/pro')"
          class="text-sm text-[#40617f] hover:text-[#006d35] flex items-center gap-1.5 mb-4 transition">
          <ArrowLeftIcon class="w-4 h-4" />
          Retour à l'espace pro
        </button>
        <div class="flex items-center justify-between gap-4 flex-wrap">
          <div>
            <h1 class="font-jakarta font-extrabold text-[#001d32] text-3xl tracking-tight">Mes prestations</h1>
            <p class="text-[#40617f] text-sm mt-1">Créez et gérez les prestations que vous proposez.</p>
          </div>
          <button @click="openCreateForm"
            class="flex items-center gap-2 bg-[#006d35] text-white font-bold text-sm px-5 py-3 rounded-xl hover:bg-[#1b8848] transition">
            <PlusIcon class="w-4 h-4" />
            Nouvelle prestation
          </button>
        </div>
      </div>

      <div v-if="loading" class="flex items-center justify-center py-24">
        <div class="w-8 h-8 border-4 border-[#006d35] border-t-transparent rounded-full animate-spin" />
      </div>

      <template v-else>
        <div v-if="showForm" class="bg-white rounded-[32px] p-8 mb-8 shadow-sm">
          <h2 class="font-jakarta font-bold text-[#001d32] text-xl mb-6">
            {{ editingId ? 'Modifier la prestation' : 'Nouvelle prestation' }}
          </h2>
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-5">
            <div class="sm:col-span-2">
              <label class="block text-xs font-semibold text-[#40617f] uppercase tracking-wide mb-1">Titre</label>
              <input v-model="form.title" type="text"
                class="w-full border border-gray-200 rounded-xl px-4 py-3 text-sm text-[#001d32] outline-none focus:ring-2 focus:ring-[#006d35]/20" />
            </div>
            <div>
              <label class="block text-xs font-semibold text-[#40617f] uppercase tracking-wide mb-1">Catégorie</label>
              <select v-model="form.category_id"
                class="w-full border border-gray-200 rounded-xl px-4 py-3 text-sm text-[#001d32] outline-none focus:ring-2 focus:ring-[#006d35]/20 bg-white">
                <option :value="null">Aucune</option>
                <option v-for="cat in categories" :key="cat.id" :value="cat.id">{{ cat.name }}</option>
              </select>
            </div>
            <div>
              <label class="block text-xs font-semibold text-[#40617f] uppercase tracking-wide mb-1">Type de prix</label>
              <select v-model="form.price_type"
                class="w-full border border-gray-200 rounded-xl px-4 py-3 text-sm text-[#001d32] outline-none focus:ring-2 focus:ring-[#006d35]/20 bg-white">
                <option value="fixed">Prix fixe</option>
                <option value="hourly">Tarif horaire</option>
                <option value="quote">Sur devis</option>
              </select>
            </div>
            <div v-if="form.price_type !== 'quote'">
              <label class="block text-xs font-semibold text-[#40617f] uppercase tracking-wide mb-1">Prix (€)</label>
              <input v-model="form.price" type="number" step="0.01" min="0" placeholder="0 pour gratuit"
                class="w-full border border-gray-200 rounded-xl px-4 py-3 text-sm text-[#001d32] outline-none focus:ring-2 focus:ring-[#006d35]/20" />
            </div>
            <div class="sm:col-span-2">
              <label class="block text-xs font-semibold text-[#40617f] uppercase tracking-wide mb-1">Description</label>
              <textarea v-model="form.description" rows="4"
                class="w-full border border-gray-200 rounded-xl px-4 py-3 text-sm text-[#001d32] outline-none focus:ring-2 focus:ring-[#006d35]/20 resize-none" />
            </div>
          </div>

          <p v-if="formError" class="mt-4 text-sm text-red-600 font-medium">{{ formError }}</p>

          <div class="flex gap-3 mt-6">
            <button @click="saveForm" :disabled="formLoading"
              class="bg-[#006d35] text-white font-bold text-sm px-6 py-3 rounded-xl hover:bg-[#1b8848] transition disabled:opacity-60">
              {{ formLoading ? 'Enregistrement…' : 'Enregistrer' }}
            </button>
            <button @click="closeForm"
              class="bg-gray-100 text-[#40617f] font-bold text-sm px-6 py-3 rounded-xl hover:bg-gray-200 transition">
              Annuler
            </button>
          </div>
        </div>

        <div v-if="!prestations.length && !showForm" class="flex flex-col items-center justify-center py-20 text-center">
          <BriefcaseIcon class="w-14 h-14 text-[#40617f]/30 mb-4" />
          <p class="font-jakarta font-bold text-[#001d32] text-lg mb-1">Aucune prestation</p>
          <p class="text-[#40617f] text-sm">Créez votre première prestation pour commencer.</p>
        </div>

        <div v-if="prestations.length" class="flex flex-col gap-4">
          <div v-for="s in prestations" :key="s.id"
            class="bg-white rounded-2xl p-6 shadow-sm flex flex-col sm:flex-row sm:items-center gap-4">
            <div class="flex-1 min-w-0">
              <div class="flex items-center gap-2 flex-wrap mb-1">
                <span :class="['text-xs font-bold px-2.5 py-1 rounded-full', statusClass(s.status)]">{{ statusLabel(s.status) }}</span>
                <span class="text-xs text-[#40617f]">{{ priceLabel(s) }}</span>
              </div>
              <p class="font-jakarta font-bold text-[#001d32] text-base truncate">{{ s.title }}</p>
              <p v-if="s.category" class="text-[#40617f] text-xs mt-0.5">{{ s.category.name }}</p>
            </div>
            <div class="flex items-center gap-2 shrink-0 flex-wrap">
              <button v-if="s.status === 'draft' || s.status === 'rejected'" @click="editPrestation(s)"
                class="text-xs bg-[#edf4ff] text-[#1a73e8] font-semibold px-3 py-2 rounded-xl hover:bg-[#d0e8ff] transition">
                Modifier
              </button>
              <button v-if="s.status === 'draft' || s.status === 'rejected'" @click="submitPrestation(s)"
                class="text-xs bg-[#d1fae5] text-[#065f46] font-semibold px-3 py-2 rounded-xl hover:bg-[#a7f3d0] transition">
                Soumettre
              </button>
              <button v-if="s.status !== 'published'" @click="deletePrestation(s)"
                class="text-xs bg-red-50 text-red-600 font-semibold px-3 py-2 rounded-xl hover:bg-red-100 transition">
                Supprimer
              </button>
            </div>
          </div>
        </div>
      </template>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { userApi, publicGet } from '@/services/publicApi'
import { useUserAuthStore } from '@/stores/userAuth'
import { ArrowLeftIcon, PlusIcon, BriefcaseIcon } from '@heroicons/vue/24/outline'

const router = useRouter()
const userAuth = useUserAuthStore()

const loading = ref(true)
const prestations = ref([])
const categories = ref([])
const showForm = ref(false)
const editingId = ref(null)
const formLoading = ref(false)
const formError = ref('')

const form = ref({ title: '', description: '', category_id: null, price_type: 'fixed', price: '' })

function openCreateForm() {
  editingId.value = null
  form.value = { title: '', description: '', category_id: null, price_type: 'fixed', price: '' }
  formError.value = ''
  showForm.value = true
}

function editPrestation(s) {
  editingId.value = s.id
  form.value = {
    title: s.title || '',
    description: s.description || '',
    category_id: s.category?.id || null,
    price_type: s.price_type || 'fixed',
    price: s.price != null ? String(s.price) : '',
  }
  formError.value = ''
  showForm.value = true
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

function closeForm() {
  showForm.value = false
  editingId.value = null
  formError.value = ''
}

async function saveForm() {
  formError.value = ''
  if (!form.value.title.trim()) {
    formError.value = 'Le titre est requis.'
    return
  }
  const payload = {
    title: form.value.title.trim(),
    description: form.value.description || null,
    category_id: form.value.category_id || null,
    price_type: form.value.price_type,
    price: form.value.price_type === 'quote' ? null : (form.value.price === '' ? '0' : String(form.value.price)),
  }
  formLoading.value = true
  try {
    if (editingId.value) {
      const saved = await userApi(`/provider/prestations/${editingId.value}`, { method: 'PUT', body: JSON.stringify(payload) })
      const idx = prestations.value.findIndex(p => p.id === editingId.value)
      if (idx !== -1) prestations.value[idx] = saved
    } else {
      const saved = await userApi('/provider/prestations', { method: 'POST', body: JSON.stringify(payload) })
      prestations.value.unshift(saved)
    }
    closeForm()
  } catch (e) {
    formError.value = e.message || 'Erreur lors de la sauvegarde.'
  } finally {
    formLoading.value = false
  }
}

async function submitPrestation(s) {
  if (!confirm('Soumettre cette prestation pour validation par un administrateur ?')) return
  try {
    const saved = await userApi(`/provider/prestations/${s.id}/submit`, { method: 'PUT' })
    const idx = prestations.value.findIndex(p => p.id === s.id)
    if (idx !== -1) prestations.value[idx] = saved
  } catch {}
}

async function deletePrestation(s) {
  if (!confirm('Supprimer définitivement cette prestation ?')) return
  try {
    await userApi(`/provider/prestations/${s.id}`, { method: 'DELETE' })
    prestations.value = prestations.value.filter(p => p.id !== s.id)
  } catch {}
}

function statusLabel(s) {
  return {
    draft: 'Brouillon',
    pending: 'En attente de validation',
    published: 'Publiée',
    suspended: 'Suspendue',
    rejected: 'Refusée',
  }[s] || s
}

function statusClass(s) {
  return {
    draft: 'bg-gray-100 text-gray-600',
    pending: 'bg-yellow-50 text-yellow-700',
    published: 'bg-green-50 text-green-700',
    suspended: 'bg-orange-50 text-orange-700',
    rejected: 'bg-red-50 text-red-600',
  }[s] || 'bg-gray-100 text-gray-600'
}

function priceLabel(s) {
  if (s.price_type === 'quote') return 'Sur devis'
  const n = parseFloat(s.price)
  if (!n || isNaN(n)) return 'Gratuit'
  const formatted = new Intl.NumberFormat('fr-FR', { style: 'currency', currency: 'EUR' }).format(n)
  return s.price_type === 'hourly' ? `${formatted}/h` : formatted
}

onMounted(async () => {
  if (!userAuth.isLoggedIn) {
    router.push('/connexion?redirect=/profil/pro/prestations')
    return
  }
  try {
    const [pres, cats] = await Promise.all([
      userApi('/provider/prestations?per_page=100').catch(() => ({ data: [] })),
      publicGet('/categories').catch(() => ({ data: [] })),
    ])
    prestations.value = pres.data || []
    categories.value = cats.data || cats || []
  } finally {
    loading.value = false
  }
})
</script>
