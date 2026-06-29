<template>
  <div class="space-y-5">
    <div class="flex items-center justify-between gap-3 flex-wrap">
      <div>
        <h2 class="text-xl sm:text-2xl font-bold text-[#001d32]">Planning du personnel</h2>
        <p class="text-sm text-[#64748b] mt-0.5">
          {{ isAdmin ? 'Emplois du temps récurrents et événements exceptionnels des employés.' : 'Votre emploi du temps et vos événements.' }}
        </p>
      </div>
      <div v-if="isAdmin" class="flex gap-2 flex-wrap">
        <button @click="openShiftForm" class="px-3 py-2 rounded-lg text-white text-sm font-semibold hover:opacity-90" style="background-color:#006d35;">+ Créneau récurrent</button>
        <button @click="openEventForm" class="px-3 py-2 rounded-lg text-white text-sm font-semibold hover:opacity-90" style="background-color:#1B8848;">+ Événement</button>
      </div>
    </div>

    <div class="flex items-center justify-between gap-3 flex-wrap">
      <div v-if="isAdmin" class="flex items-center gap-2">
        <label class="text-sm text-[#64748b]">Employé :</label>
        <select v-model="selectedEmployee" @change="reload" class="px-3 py-1.5 bg-white border border-[#e5e7eb] rounded-lg text-sm">
          <option value="">Tous les employés</option>
          <option v-for="e in employees" :key="e.id" :value="e.id">{{ e.first_name }} {{ e.last_name }}</option>
        </select>
      </div>
      <div class="flex items-center gap-2 ml-auto">
        <button @click="shiftWeek(-1)" class="px-2.5 py-1.5 rounded-lg border border-[#e5e7eb] text-sm hover:bg-gray-50">‹</button>
        <button @click="goToday" class="px-3 py-1.5 rounded-lg border border-[#e5e7eb] text-sm hover:bg-gray-50">Aujourd'hui</button>
        <button @click="shiftWeek(1)" class="px-2.5 py-1.5 rounded-lg border border-[#e5e7eb] text-sm hover:bg-gray-50">›</button>
        <span class="text-sm font-semibold text-[#001d32] ml-2">{{ weekLabel }}</span>
      </div>
    </div>

    <div class="card overflow-x-auto">
      <div class="min-w-[760px]">
        <div class="grid" style="grid-template-columns: 56px repeat(7, 1fr);">
          <div class="border-b border-[#f1f5f9]"></div>
          <div v-for="(d, i) in weekDays" :key="i" class="border-b border-l border-[#f1f5f9] px-2 py-2 text-center">
            <p class="text-xs font-bold text-[#001d32]">{{ dayNames[i] }}</p>
            <p class="text-xs text-[#94a3b8]">{{ d.getDate() }}/{{ d.getMonth() + 1 }}</p>
          </div>
        </div>
        <div class="grid relative" style="grid-template-columns: 56px repeat(7, 1fr);">
          <div>
            <div v-for="h in hours" :key="h" class="text-[10px] text-[#94a3b8] text-right pr-1.5 border-b border-[#f8fafc]" :style="{ height: hourHeight + 'px' }">{{ h }}h</div>
          </div>
          <div v-for="(d, i) in weekDays" :key="i" class="relative border-l border-[#f1f5f9]">
            <div v-for="h in hours" :key="h" class="border-b border-[#f8fafc]" :style="{ height: hourHeight + 'px' }"></div>
            <div v-for="b in blocksForDay(i, d)" :key="b.key"
              class="absolute left-0.5 right-0.5 rounded-md px-1.5 py-0.5 text-[10px] leading-tight overflow-hidden cursor-pointer"
              :class="b.type === 'event' ? 'bg-[#fef3c7] border border-[#f59e0b] text-[#92400e]' : 'bg-[#dcfce7] border border-[#22c55e] text-[#166534]'"
              :style="{ top: b.top + 'px', height: b.height + 'px' }"
              @click="onBlockClick(b)"
              :title="b.tooltip">
              <p class="font-bold truncate">{{ b.title }}</p>
              <p class="truncate">{{ b.time }}</p>
              <p v-if="b.subtitle" class="truncate opacity-80">{{ b.subtitle }}</p>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="flex gap-4 text-xs text-[#64748b]">
      <span class="inline-flex items-center gap-1.5"><span class="w-3 h-3 rounded bg-[#dcfce7] border border-[#22c55e]"></span> Créneau récurrent</span>
      <span class="inline-flex items-center gap-1.5"><span class="w-3 h-3 rounded bg-[#fef3c7] border border-[#f59e0b]"></span> Événement exceptionnel</span>
    </div>

    <Teleport to="body">
      <div v-if="shiftForm.show" class="fixed inset-0 z-50 flex items-center justify-center p-4">
        <div class="absolute inset-0 bg-black/60 backdrop-blur-sm" @click="shiftForm.show = false" />
        <div class="relative bg-white rounded-2xl shadow-2xl w-full max-w-md p-6">
          <h3 class="font-bold text-[#001d32] text-lg mb-4">Créneau récurrent</h3>
          <p class="text-xs text-[#64748b] mb-4">Ex : tous les lundis de 09:00 à 12:00. Répété chaque semaine.</p>
          <div class="space-y-3">
            <div>
              <label class="label">Employé</label>
              <select v-model="shiftForm.employee_id" class="input">
                <option :value="null">— Choisir —</option>
                <option v-for="e in employees" :key="e.id" :value="e.id">{{ e.first_name }} {{ e.last_name }}</option>
              </select>
            </div>
            <div>
              <label class="label">Jour</label>
              <select v-model.number="shiftForm.weekday" class="input">
                <option v-for="(n, i) in dayNames" :key="i" :value="i + 1">{{ n }}</option>
              </select>
            </div>
            <div class="grid grid-cols-2 gap-3">
              <div><label class="label">Début</label><input v-model="shiftForm.start_time" type="time" class="input" /></div>
              <div><label class="label">Fin</label><input v-model="shiftForm.end_time" type="time" class="input" /></div>
            </div>
            <div><label class="label">Libellé (optionnel)</label><input v-model="shiftForm.label" type="text" class="input" placeholder="Ex: Atelier couture" /></div>
          </div>
          <p v-if="shiftForm.error" class="text-red-600 text-sm mt-3">{{ shiftForm.error }}</p>
          <div class="flex gap-2 mt-5">
            <button @click="saveShift" :disabled="shiftForm.saving" class="px-5 py-2 rounded-lg text-white text-sm font-semibold disabled:opacity-50" style="background-color:#006d35;">Enregistrer</button>
            <button @click="shiftForm.show = false" class="px-5 py-2 rounded-lg bg-gray-100 text-[#40617f] text-sm font-semibold">Annuler</button>
          </div>
        </div>
      </div>

      <div v-if="eventForm.show" class="fixed inset-0 z-50 flex items-center justify-center p-4">
        <div class="absolute inset-0 bg-black/60 backdrop-blur-sm" @click="eventForm.show = false" />
        <div class="relative bg-white rounded-2xl shadow-2xl w-full max-w-md p-6 max-h-[92vh] overflow-y-auto">
          <h3 class="font-bold text-[#001d32] text-lg mb-4">Événement exceptionnel</h3>
          <div class="space-y-3">
            <div><label class="label">Titre</label><input v-model="eventForm.title" type="text" class="input" /></div>
            <div><label class="label">Lieu (optionnel)</label><input v-model="eventForm.location" type="text" class="input" /></div>
            <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
              <div><label class="label">Début</label><input v-model="eventForm.start_at" type="datetime-local" class="input" /></div>
              <div><label class="label">Fin</label><input v-model="eventForm.end_at" type="datetime-local" class="input" /></div>
            </div>
            <div><label class="label">Description (optionnel)</label><textarea v-model="eventForm.description" rows="2" class="input resize-none" /></div>
            <div>
              <label class="label">Employés concernés</label>
              <div class="max-h-40 overflow-y-auto border border-[#e5e7eb] rounded-lg p-2 space-y-1">
                <label v-for="e in employees" :key="e.id" class="flex items-center gap-2 text-sm text-[#001d32]">
                  <input type="checkbox" :value="e.id" v-model="eventForm.employee_ids" /> {{ e.first_name }} {{ e.last_name }}
                </label>
                <p v-if="!employees.length" class="text-xs text-[#94a3b8]">Aucun employé. Créez-en d'abord dans « Employés ».</p>
              </div>
            </div>
          </div>
          <p v-if="eventForm.error" class="text-red-600 text-sm mt-3">{{ eventForm.error }}</p>
          <div class="flex gap-2 mt-5">
            <button @click="saveEvent" :disabled="eventForm.saving" class="px-5 py-2 rounded-lg text-white text-sm font-semibold disabled:opacity-50" style="background-color:#1B8848;">Enregistrer</button>
            <button @click="eventForm.show = false" class="px-5 py-2 rounded-lg bg-gray-100 text-[#40617f] text-sm font-semibold">Annuler</button>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import api from '@/services/api'
import { useAuthStore } from '@/stores/auth'

const auth = useAuthStore()
const isAdmin = computed(() => auth.isAdmin)

const dayNames = ['Lundi', 'Mardi', 'Mercredi', 'Jeudi', 'Vendredi', 'Samedi', 'Dimanche']
const startHour = 7
const endHour = 21
const hourHeight = 44
const hours = Array.from({ length: endHour - startHour }, (_, i) => startHour + i)

const employees = ref([])
const shifts = ref([])
const events = ref([])
const selectedEmployee = ref('')
const weekStart = ref(mondayOf(new Date()))

function mondayOf(date) {
  const d = new Date(date)
  d.setHours(0, 0, 0, 0)
  const day = (d.getDay() + 6) % 7
  d.setDate(d.getDate() - day)
  return d
}
const weekDays = computed(() => Array.from({ length: 7 }, (_, i) => {
  const d = new Date(weekStart.value)
  d.setDate(d.getDate() + i)
  return d
}))
const weekLabel = computed(() => {
  const a = weekDays.value[0], b = weekDays.value[6]
  const f = (d) => `${d.getDate()}/${d.getMonth() + 1}`
  return `${f(a)} – ${f(b)}`
})

function timeToFloat(t) {
  if (!t) return 0
  const [h, m] = t.split(':').map(Number)
  return h + (m || 0) / 60
}
function sameDate(a, b) {
  return a.getFullYear() === b.getFullYear() && a.getMonth() === b.getMonth() && a.getDate() === b.getDate()
}
function pos(startFloat, endFloat) {
  const top = (startFloat - startHour) * hourHeight
  const height = Math.max(16, (endFloat - startFloat) * hourHeight)
  return { top, height }
}

function blocksForDay(dayIndex, date) {
  const out = []
  for (const s of shifts.value) {
    if (s.weekday !== dayIndex + 1) continue
    const { top, height } = pos(timeToFloat(s.start_time), timeToFloat(s.end_time))
    out.push({
      key: 'shift-' + s.id + '-' + dayIndex, type: 'shift', raw: s, top, height,
      title: s.label || (isAdmin.value && !selectedEmployee.value ? s.employee_name : 'Travail'),
      time: `${s.start_time}–${s.end_time}`,
      subtitle: isAdmin.value && !selectedEmployee.value ? s.employee_name : '',
      tooltip: `${s.employee_name} · ${s.start_time}–${s.end_time}`,
    })
  }
  for (const e of events.value) {
    const start = new Date(e.start_at), end = new Date(e.end_at)
    if (!sameDate(start, date)) continue
    const { top, height } = pos(start.getHours() + start.getMinutes() / 60, end.getHours() + end.getMinutes() / 60)
    out.push({
      key: 'event-' + e.id, type: 'event', raw: e, top, height,
      title: e.title,
      time: `${pad(start.getHours())}:${pad(start.getMinutes())}–${pad(end.getHours())}:${pad(end.getMinutes())}`,
      subtitle: e.members.map(m => m.name).join(', '),
      tooltip: `${e.title} · ${e.members.map(m => m.name).join(', ')}`,
    })
  }
  return out
}
function pad(n) { return String(n).padStart(2, '0') }

function onBlockClick(b) {
  if (!isAdmin.value) return
  if (b.type === 'shift') {
    if (confirm(`Supprimer le créneau récurrent de ${b.raw.employee_name} (${b.raw.start_time}–${b.raw.end_time}) ?`)) deleteShift(b.raw.id)
  } else {
    if (confirm(`Supprimer l'événement « ${b.raw.title} » ?`)) deleteEvent(b.raw.id)
  }
}

function shiftWeek(n) { const d = new Date(weekStart.value); d.setDate(d.getDate() + n * 7); weekStart.value = d; reload() }
function goToday() { weekStart.value = mondayOf(new Date()); reload() }

async function reload() {
  const params = {}
  if (isAdmin.value && selectedEmployee.value) params.employee_id = selectedEmployee.value
  const from = weekDays.value[0].toISOString()
  const to = new Date(weekDays.value[6]); to.setHours(23, 59, 59)
  try {
    const [sh, ev] = await Promise.all([
      api.get('/staff/shifts', { params }),
      api.get('/staff/events', { params: { ...params, from, to: to.toISOString() } }),
    ])
    shifts.value = sh.data.data || []
    events.value = ev.data.data || []
  } catch { shifts.value = []; events.value = [] }
}

const shiftForm = ref({ show: false })
function openShiftForm() {
  shiftForm.value = { show: true, employee_id: selectedEmployee.value || null, weekday: 1, start_time: '09:00', end_time: '12:00', label: '', saving: false, error: '' }
}
async function saveShift() {
  const f = shiftForm.value
  if (!f.employee_id) { f.error = 'Choisissez un employé.'; return }
  f.saving = true; f.error = ''
  try {
    await api.post('/staff/shifts', { employee_id: f.employee_id, weekday: f.weekday, start_time: f.start_time, end_time: f.end_time, label: f.label || null })
    f.show = false
    await reload()
  } catch (e) { f.error = e.response?.data?.message || 'Erreur.' } finally { f.saving = false }
}
async function deleteShift(id) { await api.delete(`/staff/shifts/${id}`); await reload() }

const eventForm = ref({ show: false })
function openEventForm() {
  const d = weekDays.value[0]
  const day = `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}`
  eventForm.value = { show: true, title: '', location: '', description: '', start_at: `${day}T09:00`, end_at: `${day}T11:00`, employee_ids: [], saving: false, error: '' }
}
async function saveEvent() {
  const f = eventForm.value
  if (!f.title.trim()) { f.error = 'Titre requis.'; return }
  if (!f.start_at || !f.end_at) { f.error = 'Dates requises.'; return }
  f.saving = true; f.error = ''
  try {
    await api.post('/staff/events', {
      title: f.title.trim(), location: f.location || null, description: f.description || null,
      start_at: f.start_at, end_at: f.end_at, employee_ids: f.employee_ids,
    })
    f.show = false
    await reload()
  } catch (e) { f.error = e.response?.data?.message || 'Erreur.' } finally { f.saving = false }
}
async function deleteEvent(id) { await api.delete(`/staff/events/${id}`); await reload() }

onMounted(async () => {
  if (isAdmin.value) {
    try { employees.value = (await api.get('/employees')).data.data || [] } catch { employees.value = [] }
  }
  await reload()
})
</script>
