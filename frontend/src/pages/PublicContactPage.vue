<template>
  <div class="max-w-6xl mx-auto px-4 sm:px-6 py-8 sm:py-12">

    <div v-if="!userAuth.isLoggedIn" class="max-w-md mx-auto bg-white rounded-2xl border border-[#e2e8f0] shadow-sm p-8 text-center">
      <ChatBubbleLeftRightIcon class="w-12 h-12 mx-auto text-[#006d35] mb-4" />
      <h1 class="text-xl font-bold text-[#001d32] mb-2">Contactez-nous</h1>
      <p class="text-[#64748b] text-sm mb-6">Connectez-vous pour ouvrir une demande auprès de notre équipe et suivre vos échanges.</p>
      <RouterLink to="/connexion" class="inline-block px-6 py-2.5 rounded-full text-white font-bold text-sm" style="background:#006d35;">Se connecter</RouterLink>
    </div>

    <div v-else>
      <div class="flex items-center justify-between flex-wrap gap-3 mb-6">
        <div>
          <h1 class="text-2xl font-bold text-[#001d32]">Mes demandes</h1>
          <p class="text-sm text-[#64748b] mt-0.5">Échangez avec notre équipe support.</p>
        </div>
        <button @click="startNew" class="inline-flex items-center gap-2 px-4 py-2.5 rounded-xl text-white font-bold text-sm hover:opacity-90 transition" style="background:linear-gradient(135deg,#006d35,#1b8848);">
          <PlusIcon class="w-4 h-4" /> Nouvelle demande
        </button>
      </div>

      <div class="grid md:grid-cols-[320px_1fr] gap-5">

        <div class="bg-white rounded-2xl border border-[#e2e8f0] shadow-sm overflow-hidden md:max-h-[70vh] flex flex-col">
          <div v-if="loadingList" class="p-4 space-y-3">
            <div v-for="n in 4" :key="n" class="h-14 bg-gray-50 rounded-lg animate-pulse" />
          </div>
          <div v-else-if="!tickets.length" class="p-8 text-center text-sm text-[#64748b]">
            Aucune demande pour l'instant.
          </div>
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
              <p class="text-xs text-[#64748b]">{{ formatDate(t.last_message_at) }}</p>
            </button>
          </div>
        </div>

        <div class="bg-white rounded-2xl border border-[#e2e8f0] shadow-sm overflow-hidden flex flex-col md:max-h-[70vh]">

          <div v-if="mode === 'create'" class="p-5 sm:p-6 overflow-y-auto">
            <h2 class="text-lg font-bold text-[#001d32] mb-4">Nouvelle demande</h2>
            <div class="space-y-4">
              <div>
                <label class="block text-xs font-semibold text-[#40617f] uppercase mb-1">Sujet</label>
                <input v-model="form.subject" type="text" maxlength="200" placeholder="Résumé de votre demande" class="w-full px-3 py-2.5 bg-[#f8fafc] border border-[#e5e7eb] rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-[#006d35]/30" />
              </div>

              <div class="grid sm:grid-cols-2 gap-3">
                <div>
                  <label class="block text-xs font-semibold text-[#40617f] uppercase mb-1">Élément concerné (optionnel)</label>
                  <select v-model="form.refType" @change="onRefTypeChange" class="w-full px-3 py-2.5 bg-[#f8fafc] border border-[#e5e7eb] rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-[#006d35]/30">
                    <option value="">Aucun</option>
                    <option value="prestation">Prestation</option>
                    <option value="event">Événement</option>
                    <option value="project">Projet</option>
                  </select>
                </div>
                <div v-if="form.refType">
                  <label class="block text-xs font-semibold text-[#40617f] uppercase mb-1">Sélection</label>
                  <select v-model="form.refId" class="w-full px-3 py-2.5 bg-[#f8fafc] border border-[#e5e7eb] rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-[#006d35]/30">
                    <option :value="null">— Choisir —</option>
                    <option v-for="o in refOptions" :key="o.id" :value="o.id">{{ o.title }}</option>
                  </select>
                </div>
              </div>

              <div>
                <label class="block text-xs font-semibold text-[#40617f] uppercase mb-1">Message</label>
                <textarea v-model="form.content" rows="5" placeholder="Décrivez votre demande..." class="w-full px-3 py-2.5 bg-[#f8fafc] border border-[#e5e7eb] rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-[#006d35]/30 resize-y" />
              </div>

              <div>
                <div class="flex flex-wrap gap-2 mb-2">
                  <div v-for="(img, i) in form.images" :key="img" class="relative w-16 h-16 rounded-lg overflow-hidden border border-[#e5e7eb]">
                    <img :src="img" class="w-full h-full object-cover" />
                    <button @click="form.images.splice(i,1)" class="absolute top-0.5 right-0.5 w-5 h-5 rounded-full bg-black/60 text-white text-xs flex items-center justify-center">×</button>
                  </div>
                </div>
                <label class="inline-flex items-center gap-2 text-sm text-[#006d35] font-semibold cursor-pointer">
                  <PhotoIcon class="w-5 h-5" /> {{ uploading ? 'Envoi...' : 'Ajouter une image' }}
                  <input type="file" accept="image/*" class="hidden" @change="onFilePick($event, form.images)" :disabled="uploading || form.images.length >= 6" />
                </label>
              </div>

              <p v-if="formError" class="text-red-600 text-sm">{{ formError }}</p>

              <div class="flex justify-end gap-3 pt-2">
                <button @click="mode = selected ? 'view' : 'empty'" class="px-4 py-2 text-sm text-[#40617f] font-medium">Annuler</button>
                <button @click="submitTicket" :disabled="sending" class="px-6 py-2.5 rounded-xl text-white font-bold text-sm hover:opacity-90 disabled:opacity-50 transition" style="background:linear-gradient(135deg,#006d35,#1b8848);">
                  {{ sending ? 'Envoi...' : 'Envoyer' }}
                </button>
              </div>
            </div>
          </div>

          <template v-else-if="mode === 'view' && selected">
            <div class="px-5 py-3.5 border-b border-[#f1f5f9] shrink-0">
              <div class="flex items-center justify-between gap-3">
                <h2 class="text-base font-bold text-[#001d32] truncate">{{ selected.subject }}</h2>
                <span class="shrink-0 text-[10px] font-bold uppercase px-2 py-0.5 rounded-full" :class="statusClass(selected.status)">{{ statusLabel(selected.status) }}</span>
              </div>
              <RouterLink v-if="selected.ref" :to="selected.ref.path" class="inline-flex items-center gap-1 mt-1 text-xs text-[#006d35] hover:underline">
                <LinkIcon class="w-3.5 h-3.5" /> {{ refTypeLabel(selected.ref.type) }} : {{ selected.ref.label }}
              </RouterLink>
            </div>

            <div ref="scrollArea" class="flex-1 overflow-y-auto p-4 sm:p-5 space-y-4 bg-[#f8fafc]">
              <div v-for="m in selected.messages" :key="m.id" :class="['flex gap-2.5', m.is_staff ? '' : 'flex-row-reverse']">
                <div class="w-8 h-8 rounded-full flex items-center justify-center text-white text-xs font-bold shrink-0" :style="{ backgroundColor: m.is_staff ? '#006d35' : '#40617f' }">
                  {{ (m.sender_name || '?').charAt(0).toUpperCase() }}
                </div>
                <div :class="['max-w-[75%]']">
                  <div :class="['rounded-2xl px-3.5 py-2.5', m.is_staff ? 'bg-white border border-[#e2e8f0]' : 'bg-[#006d35] text-white']">
                    <p v-if="m.is_staff" class="text-[10px] font-bold uppercase text-[#006d35] mb-0.5">{{ m.sender_name }} · Équipe</p>
                    <p v-if="m.content" class="text-sm whitespace-pre-wrap break-words">{{ m.content }}</p>
                    <div v-if="m.images && m.images.length" class="flex flex-wrap gap-2 mt-2">
                      <a v-for="img in m.images" :key="img" :href="img" target="_blank" class="block w-24 h-24 rounded-lg overflow-hidden border" :class="m.is_staff ? 'border-[#e5e7eb]' : 'border-white/30'">
                        <img :src="img" class="w-full h-full object-cover" />
                      </a>
                    </div>
                  </div>
                  <p :class="['text-[10px] text-[#64748b] mt-1', m.is_staff ? '' : 'text-right']">{{ formatTime(m.created_at) }}</p>
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
                <label class="shrink-0 p-2 rounded-lg text-[#40617f] hover:bg-[#f1f5f9] cursor-pointer">
                  <PhotoIcon class="w-5 h-5" />
                  <input type="file" accept="image/*" class="hidden" @change="onFilePick($event, reply.images)" :disabled="uploading || reply.images.length >= 6" />
                </label>
                <textarea v-model="reply.content" rows="1" placeholder="Votre message..." @keydown.enter.exact.prevent="sendReply" class="flex-1 px-3 py-2 bg-[#f8fafc] border border-[#e5e7eb] rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-[#006d35]/30 resize-none max-h-32" />
                <button @click="sendReply" :disabled="sending" class="shrink-0 px-4 py-2 rounded-xl text-white font-bold text-sm hover:opacity-90 disabled:opacity-50 transition" style="background:#006d35;">Envoyer</button>
              </div>
            </div>
            <div v-else class="border-t border-[#f1f5f9] p-4 text-center text-sm text-[#64748b] shrink-0">
              Cette demande est clôturée.
            </div>
          </template>

          <div v-else class="flex-1 flex items-center justify-center p-8 text-center text-sm text-[#64748b]">
            Sélectionnez une demande ou créez-en une nouvelle.
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, nextTick, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { ChatBubbleLeftRightIcon, PlusIcon, PhotoIcon, LinkIcon } from '@heroicons/vue/24/outline'
import { useUserAuthStore } from '@/stores/userAuth'
import { userApi, publicGet } from '@/services/publicApi'

const userAuth = useUserAuthStore()
const route = useRoute()

const tickets = ref([])
const loadingList = ref(false)
const selected = ref(null)
const mode = ref('empty')
const sending = ref(false)
const uploading = ref(false)
const formError = ref('')
const scrollArea = ref(null)

const form = reactive({ subject: '', content: '', images: [], refType: '', refId: null })
const reply = reactive({ content: '', images: [] })
const refOptions = ref([])

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
    const { data } = await userApi('/tickets')
    tickets.value = data || []
  } finally {
    loadingList.value = false
  }
}

async function openTicket(id) {
  const { data } = await userApi(`/tickets/${id}`)
  selected.value = data
  mode.value = 'view'
  await nextTick()
  scrollBottom()
  const idx = tickets.value.findIndex(t => t.id === id)
  if (idx !== -1) tickets.value[idx] = { ...tickets.value[idx], status: data.status, last_message_at: data.last_message_at }
}

function scrollBottom() {
  if (scrollArea.value) scrollArea.value.scrollTop = scrollArea.value.scrollHeight
}

function startNew() {
  form.subject = ''
  form.content = ''
  form.images = []
  form.refType = ''
  form.refId = null
  formError.value = ''
  refOptions.value = []
  mode.value = 'create'
}

async function onRefTypeChange() {
  form.refId = null
  refOptions.value = []
  if (!form.refType) return
  const path = { prestation: '/prestations', event: '/events', project: '/projects' }[form.refType]
  try {
    const res = await publicGet(path)
    const list = res.data || res || []
    refOptions.value = list.map(x => ({ id: x.id, title: x.title || x.name || `#${x.id}` }))
  } catch {
    refOptions.value = []
  }
}

async function onFilePick(e, target) {
  const file = e.target.files?.[0]
  if (!file) return
  uploading.value = true
  try {
    const fd = new FormData()
    fd.append('image', file)
    const res = await fetch('/api/v1/ticket-attachments', {
      method: 'POST',
      headers: { Authorization: `Bearer ${userAuth.token}` },
      body: fd,
    })
    if (res.ok) {
      const data = await res.json()
      target.push(data.url)
    }
  } finally {
    uploading.value = false
    if (e.target) e.target.value = ''
  }
}

async function submitTicket() {
  formError.value = ''
  if (!form.subject.trim()) { formError.value = 'Le sujet est requis.'; return }
  if (!form.content.trim() && !form.images.length) { formError.value = 'Un message ou une image est requis.'; return }
  sending.value = true
  try {
    const body = {
      subject: form.subject.trim(),
      content: form.content.trim(),
      images: form.images,
      ref_type: form.refType || null,
      ref_id: form.refType && form.refId ? form.refId : null,
    }
    const { data } = await userApi('/tickets', { method: 'POST', body: JSON.stringify(body) })
    await loadList()
    await openTicket(data.id)
  } catch (err) {
    formError.value = err.message || 'Erreur lors de l\'envoi.'
  } finally {
    sending.value = false
  }
}

async function sendReply() {
  if (!reply.content.trim() && !reply.images.length) return
  sending.value = true
  try {
    const body = { content: reply.content.trim(), images: reply.images }
    await userApi(`/tickets/${selected.value.id}/messages`, { method: 'POST', body: JSON.stringify(body) })
    reply.content = ''
    reply.images = []
    await openTicket(selected.value.id)
  } finally {
    sending.value = false
  }
}

onMounted(async () => {
  if (!userAuth.isLoggedIn) return
  await loadList()
  const id = Number(route.query.id)
  if (id) openTicket(id)
})
</script>
