<template>
  <div class="bg-[#f7f9ff] min-h-screen pb-16">
    <div class="max-w-[960px] mx-auto px-4 sm:px-6 py-10">

      <div class="mb-8">
        <button @click="router.push('/profil/pro')"
          class="text-sm text-[#40617f] hover:text-[#006d35] flex items-center gap-1.5 mb-4 transition">
          <ArrowLeftIcon class="w-4 h-4" />
          Retour à l'espace pro
        </button>
        <h1 class="font-jakarta font-extrabold text-[#001d32] text-3xl tracking-tight">Mes réservations</h1>
        <p class="text-[#40617f] text-sm mt-1">Les réservations reçues sur vos prestations.</p>
      </div>

      <div class="flex flex-wrap gap-2 mb-6">
        <button v-for="tab in tabs" :key="tab.key"
          @click="activeTab = tab.key"
          class="px-4 py-2 rounded-full text-sm font-semibold border-2 transition-all"
          :class="activeTab === tab.key
            ? 'bg-[#006d35] text-white border-[#006d35]'
            : 'bg-white text-[#40617f] border-[#e5e7eb] hover:border-[#006d35]/30'">
          {{ tab.label }}
          <span class="ml-1.5 text-xs opacity-80">{{ countFor(tab.key) }}</span>
        </button>
      </div>

      <div v-if="loading" class="flex items-center justify-center py-24">
        <div class="w-8 h-8 border-4 border-[#006d35] border-t-transparent rounded-full animate-spin" />
      </div>

      <div v-else-if="filtered.length === 0" class="flex flex-col items-center justify-center py-20 text-center">
        <InboxIcon class="w-14 h-14 text-[#40617f]/30 mb-4" />
        <p class="font-jakarta font-bold text-[#001d32] text-lg mb-1">Aucune réservation</p>
        <p class="text-[#40617f] text-sm">Vous n'avez pas encore de réservation dans cette catégorie.</p>
      </div>

      <div v-else class="flex flex-col gap-3">
        <button v-for="r in filtered" :key="r.id"
          @click="router.push(`/profil/reservations/${r.id}`)"
          class="w-full text-left bg-white rounded-2xl p-5 shadow-sm flex items-center gap-4 hover:shadow-md transition">

          <div class="shrink-0 w-11 h-11 rounded-full flex items-center justify-center text-white font-jakarta font-bold text-sm overflow-hidden"
            style="background: linear-gradient(135deg, #006d35, #1b8848);">
            <img v-if="r.client?.avatar_url" :src="r.client.avatar_url" class="w-full h-full object-cover" />
            <span v-else>{{ initials(r.client) }}</span>
          </div>

          <div class="flex-1 min-w-0">
            <div class="flex items-center gap-2 flex-wrap mb-0.5">
              <span :class="['text-xs font-bold px-2.5 py-1 rounded-full', statusClass(r.status)]">
                {{ statusLabel(r.status) }}
              </span>
              <span class="text-xs text-[#94a3b8]">{{ formatDate(r.created_at) }}</span>
            </div>
            <p class="font-jakarta font-bold text-[#001d32] text-base truncate">{{ r.prestation?.title || 'Prestation supprimée' }}</p>
            <p class="text-[#40617f] text-xs mt-0.5 truncate">{{ clientName(r.client) }}</p>
          </div>

          <div class="shrink-0 text-right">
            <p class="font-jakarta font-bold text-[#001d32] text-sm">{{ amountLabel(r) }}</p>
            <ChevronRightIcon class="w-4 h-4 text-[#94a3b8] inline-block mt-1" />
          </div>
        </button>
      </div>

    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { userApi } from '@/services/publicApi'
import { useUserAuthStore } from '@/stores/userAuth'
import { ArrowLeftIcon, InboxIcon, ChevronRightIcon } from '@heroicons/vue/24/outline'

const router = useRouter()
const { locale } = useI18n()
const userAuth = useUserAuthStore()

const loading = ref(true)
const reservations = ref([])
const activeTab = ref('all')

const tabs = [
  { key: 'all', label: 'Toutes' },
  { key: 'pending', label: 'En attente' },
  { key: 'confirmed', label: 'Confirmées' },
  { key: 'cancelled', label: 'Annulées' },
]

const groupOf = (s) => {
  if (s === 'paid') return 'confirmed'
  if (s === 'cancelled' || s === 'failed') return 'cancelled'
  return 'pending'
}

const filtered = computed(() =>
  activeTab.value === 'all'
    ? reservations.value
    : reservations.value.filter(r => groupOf(r.status) === activeTab.value)
)

function countFor(key) {
  if (key === 'all') return reservations.value.length
  return reservations.value.filter(r => groupOf(r.status) === key).length
}

function clientName(c) {
  if (!c) return 'Client'
  const n = `${c.first_name || ''} ${c.last_name || ''}`.trim()
  return n || c.email || 'Client'
}

function initials(c) {
  if (!c) return '?'
  const a = (c.first_name || '').charAt(0)
  const b = (c.last_name || '').charAt(0)
  return (a + b).toUpperCase() || (c.email || '?').charAt(0).toUpperCase()
}

function statusLabel(s) {
  return {
    pending: 'En attente de paiement',
    paid: 'Confirmée',
    quote_requested: 'Devis demandé',
    quote_issued: 'Devis envoyé',
    cancelled: 'Annulée',
    failed: 'Échec du paiement',
  }[s] || s
}

function statusClass(s) {
  return {
    pending: 'bg-amber-50 text-amber-700',
    paid: 'bg-green-50 text-green-700',
    quote_requested: 'bg-blue-50 text-blue-700',
    quote_issued: 'bg-blue-50 text-blue-700',
    cancelled: 'bg-gray-100 text-gray-500',
    failed: 'bg-red-50 text-red-600',
  }[s] || 'bg-gray-100 text-gray-600'
}

function amountLabel(r) {
  if (r.status === 'quote_requested') return 'Devis'
  if (!r.amount_cents) return 'Gratuit'
  return new Intl.NumberFormat('fr-FR', { style: 'currency', currency: 'EUR' }).format(r.amount_cents / 100)
}

function formatDate(iso) {
  if (!iso) return ''
  return new Date(iso).toLocaleDateString(locale.value === 'en' ? 'en-US' : 'fr-FR', {
    day: 'numeric', month: 'short', year: 'numeric',
  })
}

onMounted(async () => {
  if (!userAuth.isLoggedIn) {
    router.push('/connexion?redirect=/profil/pro/reservations')
    return
  }
  try {
    const res = await userApi('/provider/reservations?per_page=100')
    reservations.value = res.data || []
  } catch {
    reservations.value = []
  } finally {
    loading.value = false
  }
})
</script>
