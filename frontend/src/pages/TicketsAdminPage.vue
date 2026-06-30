<template>
  <div class="space-y-5">
    <div class="flex items-center justify-between flex-wrap gap-3">
      <div>
        <h2 class="text-xl sm:text-2xl font-bold text-[#001d32]">Tickets de contact</h2>
        <p class="text-sm text-[#40617f] mt-0.5">Demandes envoyées par les utilisateurs.</p>
      </div>
      <div class="flex gap-1.5 flex-wrap">
        <button
          v-for="f in filters"
          :key="f.value"
          @click="setFilter(f.value)"
          :class="filter === f.value ? 'bg-[#006d35] text-white border-[#006d35]' : 'bg-white text-[#40617f] border-[#e5e7eb] hover:border-[#006d35]'"
          class="px-3 py-1.5 rounded-full text-xs font-semibold border transition"
        >
          {{ f.label }}<span v-if="f.value==='open' && meta.open"> ({{ meta.open }})</span><span v-if="f.value==='in_progress' && meta.in_progress"> ({{ meta.in_progress }})</span>
        </button>
      </div>
    </div>

    <div class="grid lg:grid-cols-[340px_1fr] gap-5">

      <div class="card overflow-hidden flex flex-col lg:max-h-[72vh]">
        <div v-if="loadingList" class="p-4 space-y-3">
          <div v-for="n in 5" :key="n" class="h-14 bg-gray-50 rounded-lg animate-pulse" />
        </div>
        <div v-else-if="!tickets.length" class="p-8 text-center text-sm text-gray-500">Aucun ticket.</div>
        <div v-else class="divide-y divide-[#f1f5f9] overflow-y-auto">
          <button
            v-for="t in tickets"
            :key="t.id"
            @click="openTicket(t.id)"
            :class="['w-full text-left px-4 py-3 transition', selected?.id === t.id ? 'bg-[#f0fdf4]' : 'hover:bg-[#f8fafc]']"
          >
            <div class="flex items-center justify-between gap-2 mb-1">
              <span class="font-semibold text-sm text-[#001d32] truncate">{{ t.subject }}</span>
              <span class="shrink-0 text-[10px] font-bold uppercase px-2 py-0.5 rounded-full" :class="statusClass(t.status)">{{ statusLabel(t.status) }}</span>
            </div>
            <p class="text-xs text-[#64748b] truncate">{{ t.user_name }} · {{ t.user_email }}</p>
            <p class="text-[11px] text-gray-500 mt-0.5">{{ formatDate(t.last_message_at) }}</p>
          </button>
        </div>
      </div>

      <div class="card overflow-hidden flex flex-col lg:max-h-[72vh]">
        <template v-if="selected">
          <div class="px-5 py-3.5 border-b border-[#f1f5f9] shrink-0">
            <div class="flex items-center justify-between gap-3 flex-wrap">
              <div class="min-w-0">
                <h3 class="text-base font-bold text-[#001d32] truncate">{{ selected.subject }}</h3>
                <p class="text-xs text-[#64748b]">{{ selected.user_name }} · {{ selected.user_email }}</p>
                <RouterLink v-if="selected.ref" :to="selected.ref.path" target="_blank" class="inline-flex items-center gap-1 mt-1 text-xs text-[#006d35] hover:underline">
                  <LinkIcon class="w-3.5 h-3.5" /> {{ refTypeLabel(selected.ref.type) }} : {{ selected.ref.label }}
                </RouterLink>
              </div>
              <select :value="selected.status" aria-label="Statut du ticket" @change="changeStatus($event.target.value)" class="shrink-0 px-3 py-1.5 bg-[#f8fafc] border border-[#e5e7eb] rounded-lg text-xs font-semibold focus:outline-none focus:ring-2 focus:ring-[#006d35]/30">
                <option value="open">Ouvert</option>
                <option value="in_progress">En cours</option>
                <option value="resolved">Résolu</option>
                <option value="closed">Clôturé</option>
              </select>
            </div>
          </div>

          <div ref="scrollArea" class="flex-1 overflow-y-auto p-4 sm:p-5 space-y-4 bg-[#f8fafc]">
            <div v-for="m in selected.messages" :key="m.id" :class="['flex gap-2.5', m.is_staff ? 'flex-row-reverse' : '']">
              <div class="w-8 h-8 rounded-full flex items-center justify-center text-white text-xs font-bold shrink-0" :style="{ backgroundColor: m.is_staff ? '#006d35' : '#40617f' }">
                {{ (m.sender_name || '?').charAt(0).toUpperCase() }}
              </div>
              <div class="max-w-[75%]">
                <div :class="['rounded-2xl px-3.5 py-2.5', m.is_staff ? 'bg-[#006d35] text-white' : 'bg-white border border-[#e2e8f0]']">
                  <p v-if="!m.is_staff" class="text-[10px] font-bold uppercase text-[#40617f] mb-0.5">{{ m.sender_name }}</p>
                  <p v-if="m.content" class="text-sm whitespace-pre-wrap break-words">{{ m.content }}</p>
                  <div v-if="m.images && m.images.length" class="flex flex-wrap gap-2 mt-2">
                    <a v-for="img in m.images" :key="img" :href="img" target="_blank" class="block w-24 h-24 rounded-lg overflow-hidden border" :class="m.is_staff ? 'border-white/30' : 'border-[#e5e7eb]'">
                      <img :src="img" class="w-full h-full object-cover" />
                    </a>
                  </div>
                </div>
                <p :class="['text-[10px] text-gray-500 mt-1', m.is_staff ? 'text-right' : '']">{{ formatTime(m.created_at) }}</p>
              </div>
            </div>
          </div>

          <div v-if="selected.status !== 'closed'" class="border-t border-[#f1f5f9] p-3 shrink-0">
            <div v-if="reply.images.length" class="flex flex-wrap gap-2 mb-2">
              <div v-for="(img, i) in reply.images" :key="img" class="relative w-14 h-14 rounded-lg overflow-hidden border border-[#e5e7eb]">
                <img :src="img" class="w-full h-full object-cover" />
                <button @click="reply.images.splice(i,1)" class="absolute top-0.5 right-0.5 w-5 h-5 rounded-full bg-black/60 text-white text-xs flex items-center justify-center">×</button>
              </div>
            </div>
            <div class="flex items-end gap-2">
              <label class="shrink-0 p-2 rounded-lg text-[#40617f] hover:bg-[#f1f5f9] cursor-pointer" for="lfticket-content">
                <PhotoIcon class="w-5 h-5" />
                <input type="file" accept="image/*" class="hidden" @change="onFilePick" :disabled="uploading || reply.images.length >= 6" />
              </label>
              <textarea v-model="reply.content" rows="1" placeholder="Répondre..." @keydown.enter.exact.prevent="sendReply" class="flex-1 px-3 py-2 bg-[#f8fafc] border border-[#e5e7eb] rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-[#006d35]/30 resize-none max-h-32"  id="lfticket-content"/>
              <button @click="sendReply" :disabled="sending" class="shrink-0 px-4 py-2 rounded-xl text-white font-bold text-sm hover:opacity-90 disabled:opacity-50 transition" style="background:#006d35;">Envoyer</button>
            </div>
          </div>
          <div v-else class="border-t border-[#f1f5f9] p-4 text-center text-sm text-gray-500 shrink-0">Ticket clôturé.</div>
        </template>

        <div v-else class="flex-1 flex items-center justify-center p-8 text-center text-sm text-gray-500">
          Sélectionnez un ticket pour afficher la conversation.
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, nextTick, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { PhotoIcon, LinkIcon } from '@heroicons/vue/24/outline'
import api from '@/services/api'

const route = useRoute()
const tickets = ref([])
const meta = ref({ open: 0, in_progress: 0 })
const selected = ref(null)
const loadingList = ref(false)
const sending = ref(false)
const uploading = ref(false)
const filter = ref('')
const scrollArea = ref(null)
const reply = reactive({ content: '', images: [] })

const filters = [
  { value: '', label: 'Tous' },
  { value: 'open', label: 'Ouverts' },
  { value: 'in_progress', label: 'En cours' },
  { value: 'resolved', label: 'Résolus' },
  { value: 'closed', label: 'Clôturés' },
]

function statusLabel(s) {
  return { open: 'Ouvert', in_progress: 'En cours', resolved: 'Résolu', closed: 'Clôturé' }[s] || s
}
function statusClass(s) {
  return {
    open: 'bg-[#dbeafe] text-[#1e40af]',
    in_progress: 'bg-[#fef9c3] text-[#854d0e]',
    resolved: 'bg-[#dcfce7] text-[#166534]',
    closed: 'bg-[#f1f5f9] text-[#475569]',
  }[s] || 'bg-[#f1f5f9] text-[#475569]'
}
function refTypeLabel(t) {
  return { prestation: 'Prestation', event: 'Événement', project: 'Projet' }[t] || t
}
function formatDate(iso) {
  return new Date(iso).toLocaleDateString('fr-FR', { day: 'numeric', month: 'short', hour: '2-digit', minute: '2-digit' })
}
function formatTime(iso) {
  return new Date(iso).toLocaleString('fr-FR', { day: 'numeric', month: 'short', hour: '2-digit', minute: '2-digit' })
}

async function loadList() {
  loadingList.value = true
  try {
    const { data } = await api.get('/tickets', { params: filter.value ? { status: filter.value } : {} })
    tickets.value = data.data || []
    if (data.meta) meta.value = data.meta
  } finally {
    loadingList.value = false
  }
}

function setFilter(v) {
  filter.value = v
  loadList()
}

async function openTicket(id) {
  const { data } = await api.get(`/tickets/${id}`)
  selected.value = data.data
  await nextTick()
  scrollBottom()
}

function scrollBottom() {
  if (scrollArea.value) scrollArea.value.scrollTop = scrollArea.value.scrollHeight
}

async function onFilePick(e) {
  const file = e.target.files?.[0]
  if (!file) return
  uploading.value = true
  try {
    const fd = new FormData()
    fd.append('image', file)
    const res = await fetch('/api/admin/v1/ticket-attachments', {
      method: 'POST',
      headers: { Authorization: `Bearer ${localStorage.getItem('admin_token')}` },
      body: fd,
    })
    if (res.ok) {
      const data = await res.json()
      reply.images.push(data.url)
    }
  } finally {
    uploading.value = false
    if (e.target) e.target.value = ''
  }
}

async function sendReply() {
  if (!reply.content.trim() && !reply.images.length) return
  sending.value = true
  try {
    await api.post(`/tickets/${selected.value.id}/messages`, { content: reply.content.trim(), images: reply.images })
    reply.content = ''
    reply.images = []
    await openTicket(selected.value.id)
    loadList()
  } finally {
    sending.value = false
  }
}

async function changeStatus(status) {
  await api.patch(`/tickets/${selected.value.id}/status`, { status })
  selected.value.status = status
  loadList()
}

onMounted(async () => {
  await loadList()
  const id = Number(route.query.id)
  if (id) openTicket(id)
})
</script>
