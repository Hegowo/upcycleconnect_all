<template>
  <div class="space-y-5">
    <div class="flex items-center justify-between gap-3 flex-wrap">
      <div>
        <h2 class="text-xl sm:text-2xl font-bold text-[#001d32]">Employés &amp; plannings</h2>
        <p class="text-sm text-[#64748b] mt-0.5">Gérez les comptes employés et leur emploi du temps. Sélectionnez un employé pour voir son planning.</p>
      </div>
      <div class="flex gap-2 flex-wrap">
        <button @click="openEventForm()" class="px-3 py-2 rounded-lg text-white text-sm font-semibold hover:opacity-90" style="background-color:#1B8848;">+ Événement (plusieurs employés)</button>
        <button @click="openEmpForm()" class="px-3 py-2 rounded-lg text-white text-sm font-semibold hover:opacity-90 inline-flex items-center gap-1.5" style="background-color:#006d35;"><PlusIcon class="w-4 h-4" /> Nouvel employé</button>
      </div>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-3 gap-5">

      <div class="lg:col-span-1 card p-3 flex flex-col" style="max-height: calc(100vh - 180px);">
        <div class="relative mb-2">
          <MagnifyingGlassIcon class="w-4 h-4 absolute left-3 top-1/2 -translate-y-1/2 text-gray-500" />
          <input v-model="search" type="text" placeholder="Rechercher un employé…" class="w-full pl-9 pr-3 py-2 bg-[#f8fafc] border border-[#e5e7eb] rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-[#006d35]/30" />
        </div>
        <div class="overflow-y-auto flex-1 -mx-1 px-1">
          <div v-if="loading" class="py-8 text-center text-[#64748b] text-sm">Chargement…</div>
          <div v-else-if="!filteredEmployees.length" class="py-8 text-center text-gray-500 text-sm">Aucun employé.</div>
          <button v-for="e in filteredEmployees" :key="e.id" @click="selectEmployee(e)"
            class="w-full text-left px-3 py-2.5 rounded-lg mb-1 flex items-center gap-3 transition border"
            :class="selectedId === e.id ? 'bg-[#f0fdf4] border-[#006d35]' : 'border-transparent hover:bg-[#f8fafc]'">
            <div class="w-9 h-9 rounded-full flex items-center justify-center text-white text-xs font-bold shrink-0" style="background:linear-gradient(135deg,#006d35,#1b8848);">{{ initials(e) }}</div>
            <div class="min-w-0 flex-1">
              <p class="font-semibold text-[#001d32] text-sm truncate">{{ e.first_name }} {{ e.last_name }}</p>
              <p class="text-xs text-gray-500 truncate">{{ e.email }}</p>
            </div>
            <span v-if="e.status !== 'active'" class="text-[10px] font-bold text-gray-500 bg-gray-200 px-1.5 py-0.5 rounded-full shrink-0">susp.</span>
          </button>
        </div>
      </div>

      <div class="lg:col-span-2 space-y-4">
        <div v-if="!selected" class="card p-12 text-center text-[#64748b]">
          <CalendarDaysIcon class="w-12 h-12 text-[#cbd5e1] mx-auto mb-3" />
          <p class="font-semibold text-[#001d32]">Sélectionnez un employé</p>
          <p class="text-sm mt-1">Choisissez un employé à gauche pour gérer son emploi du temps.</p>
        </div>

        <template v-else>
          <div class="card p-5 flex items-center justify-between gap-3 flex-wrap">
            <div class="flex items-center gap-3">
              <div class="w-11 h-11 rounded-full flex items-center justify-center text-white text-sm font-bold" style="background:linear-gradient(135deg,#006d35,#1b8848);">{{ initials(selected) }}</div>
              <div>
                <p class="font-bold text-[#001d32]">{{ selected.first_name }} {{ selected.last_name }}</p>
                <p class="text-xs text-[#64748b]">{{ selected.email }} · {{ selected.status === 'active' ? 'Actif' : 'Suspendu' }}</p>
              </div>
            </div>
            <div class="flex gap-2">
              <button @click="openEmpForm(selected)" class="px-3 py-1.5 rounded-lg bg-[#edf4ff] text-[#1a73e8] text-xs font-semibold hover:bg-[#d0e8ff]">Modifier</button>
              <button @click="removeEmployee(selected)" class="px-3 py-1.5 rounded-lg bg-red-50 text-red-600 text-xs font-semibold hover:bg-red-100">Supprimer</button>
            </div>
          </div>

          <div class="card p-5">
            <div class="flex items-center justify-between mb-3">
              <h3 class="font-bold text-[#001d32]">Emploi du temps récurrent</h3>
              <span class="text-xs text-gray-500">Répété chaque semaine</span>
            </div>
            <div class="grid grid-cols-2 sm:grid-cols-4 lg:grid-cols-7 gap-2">
              <div v-for="(day, di) in dayNames" :key="di" class="border border-[#f1f5f9] rounded-xl p-2 min-h-[90px] flex flex-col">
                <p class="text-xs font-bold text-[#40617f] mb-1.5">{{ day }}</p>
                <div class="space-y-1 flex-1">
                  <div v-for="s in shiftsByDay(di + 1)" :key="s.id" class="group bg-[#dcfce7] border border-[#22c55e]/40 rounded-md px-1.5 py-1 text-[11px] text-[#166534] relative">
                    <p class="font-semibold">{{ s.start_time }}–{{ s.end_time }}</p>
                    <p v-if="s.label" class="truncate opacity-80">{{ s.label }}</p>
                    <button @click="deleteShift(s)" class="absolute top-0.5 right-0.5 w-4 h-4 rounded-full bg-black/50 text-white text-[9px] flex items-center justify-center opacity-0 group-hover:opacity-100 transition">✕</button>
                  </div>
                </div>
                <button @click="openShiftForm(di + 1)" class="mt-1 text-[11px] text-[#006d35] font-semibold hover:underline">+ créneau</button>
              </div>
            </div>
          </div>

          <div class="card p-5">
            <h3 class="font-bold text-[#001d32] mb-3">Événements à venir</h3>
            <div v-if="!events.length" class="text-sm text-gray-500 py-2">Aucun événement à venir pour cet employé.</div>
            <div v-else class="space-y-2">
              <div v-for="ev in events" :key="ev.id" class="flex items-start justify-between gap-3 border border-[#f1f5f9] rounded-xl p-3">
                <div class="min-w-0">
                  <p class="font-semibold text-[#001d32] text-sm">{{ ev.title }}</p>
                  <p class="text-xs text-[#40617f]">{{ formatDT(ev.start_at) }} → {{ formatTime(ev.end_at) }}<span v-if="ev.location"> · {{ ev.location }}</span></p>
                  <p v-if="ev.members.length > 1" class="text-[11px] text-gray-500 mt-0.5">Avec : {{ ev.members.map(m => m.name).join(', ') }}</p>
                </div>
                <button @click="deleteEvent(ev)" class="text-xs font-semibold text-red-600 hover:underline shrink-0">Suppr.</button>
              </div>
            </div>
          </div>
        </template>
      </div>
    </div>

    <Teleport to="body">

      <div v-if="empForm.show" class="fixed inset-0 z-50 flex items-center justify-center p-4">
        <div class="absolute inset-0 bg-black/60 backdrop-blur-sm" @click="empForm.show = false" />
        <div class="relative bg-white rounded-2xl shadow-2xl w-full max-w-md p-6">
          <h3 class="font-bold text-[#001d32] text-lg mb-4">{{ empForm.id ? 'Modifier l\'employé' : 'Nouvel employé' }}</h3>
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
            <div><label class="label">Prénom</label><input v-model="empForm.first_name" type="text" class="input" /></div>
            <div><label class="label">Nom</label><input v-model="empForm.last_name" type="text" class="input" /></div>
            <div v-if="!empForm.id" class="sm:col-span-2"><label class="label">Email</label><input v-model="empForm.email" type="email" class="input" /></div>
            <div><label class="label">{{ empForm.id ? 'Nouveau mot de passe' : 'Mot de passe' }}</label><input v-model="empForm.password" type="password" class="input" placeholder="Min. 8 caractères" /></div>
            <div v-if="empForm.id"><label class="label">Statut</label><select v-model="empForm.status" class="input"><option value="active">Actif</option><option value="suspended">Suspendu</option></select></div>
          </div>
          <p v-if="empForm.error" class="text-red-600 text-sm mt-3">{{ empForm.error }}</p>
          <div class="flex gap-2 mt-5">
            <button @click="saveEmployee" :disabled="empForm.saving" class="px-5 py-2 rounded-lg text-white text-sm font-semibold disabled:opacity-50" style="background-color:#006d35;">Enregistrer</button>
            <button @click="empForm.show = false" class="px-5 py-2 rounded-lg bg-gray-100 text-[#40617f] text-sm font-semibold">Annuler</button>
          </div>
        </div>
      </div>

      <div v-if="shiftForm.show" class="fixed inset-0 z-50 flex items-center justify-center p-4">
        <div class="absolute inset-0 bg-black/60 backdrop-blur-sm" @click="shiftForm.show = false" />
        <div class="relative bg-white rounded-2xl shadow-2xl w-full max-w-sm p-6">
          <h3 class="font-bold text-[#001d32] text-lg mb-1">Créneau récurrent</h3>
          <p class="text-xs text-[#64748b] mb-4">{{ selected?.first_name }} {{ selected?.last_name }} · {{ dayNames[shiftForm.weekday - 1] }}</p>
          <div class="space-y-3">
            <div>
              <label class="label">Jour</label>
              <select v-model.number="shiftForm.weekday" class="input"><option v-for="(n, i) in dayNames" :key="i" :value="i + 1">{{ n }}</option></select>
            </div>
            <div class="grid grid-cols-2 gap-3">
              <div><label class="label">Début</label><input v-model="shiftForm.start_time" type="time" class="input" /></div>
              <div><label class="label">Fin</label><input v-model="shiftForm.end_time" type="time" class="input" /></div>
            </div>
            <div><label class="label">Libellé (optionnel)</label><input v-model="shiftForm.label" type="text" class="input" placeholder="Ex: Atelier" /></div>
          </div>
          <p v-if="shiftForm.error" class="text-red-600 text-sm mt-3">{{ shiftForm.error }}</p>
          <div class="flex gap-2 mt-5">
            <button @click="saveShift" :disabled="shiftForm.saving" class="px-5 py-2 rounded-lg text-white text-sm font-semibold disabled:opacity-50" style="background-color:#006d35;">Ajouter</button>
            <button @click="shiftForm.show = false" class="px-5 py-2 rounded-lg bg-gray-100 text-[#40617f] text-sm font-semibold">Annuler</button>
          </div>
        </div>
      </div>

      <div v-if="eventForm.show" class="fixed inset-0 z-50 flex items-center justify-center p-4">
        <div class="absolute inset-0 bg-black/60 backdrop-blur-sm" @click="eventForm.show = false" />
        <div class="relative bg-white rounded-2xl shadow-2xl w-full max-w-md p-6 max-h-[92vh] overflow-y-auto">
          <h3 class="font-bold text-[#001d32] text-lg mb-4">Événement (plusieurs employés)</h3>
          <div class="space-y-3">
            <div><label class="label">Titre</label><input v-model="eventForm.title" type="text" class="input" /></div>
            <div><label class="label">Lieu (optionnel)</label><input v-model="eventForm.location" type="text" class="input" /></div>
            <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
              <div><label class="label">Début</label><input v-model="eventForm.start_at" type="datetime-local" class="input" /></div>
              <div><label class="label">Fin</label><input v-model="eventForm.end_at" type="datetime-local" class="input" /></div>
            </div>
            <div>
              <div class="flex items-center justify-between mb-1">
                <label class="label mb-0">Employés concernés ({{ eventForm.employee_ids.length }})</label>
                <button type="button" @click="toggleAll" class="text-xs text-[#006d35] font-semibold hover:underline">{{ allSelected ? 'Tout décocher' : 'Tout cocher' }}</button>
              </div>
              <input v-model="eventEmpSearch" type="text" placeholder="Filtrer…" class="input mb-2 py-1.5 text-sm" />
              <div class="max-h-44 overflow-y-auto border border-[#e5e7eb] rounded-lg p-2 space-y-1">
                <label v-for="e in eventFilteredEmployees" :key="e.id" class="flex items-center gap-2 text-sm text-[#001d32]">
                  <input type="checkbox" :value="e.id" v-model="eventForm.employee_ids" /> {{ e.first_name }} {{ e.last_name }}
                </label>
                <p v-if="!employees.length" class="text-xs text-gray-500">Aucun employé.</p>
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
import { PlusIcon, MagnifyingGlassIcon, CalendarDaysIcon } from '@heroicons/vue/24/outline'
import api from '@/services/api'

const dayNames = ['Lundi', 'Mardi', 'Mercredi', 'Jeudi', 'Vendredi', 'Samedi', 'Dimanche']

const employees = ref([])
const loading = ref(true)
const search = ref('')
const selectedId = ref(null)
const shifts = ref([])
const events = ref([])

const filteredEmployees = computed(() => {
  const q = search.value.trim().toLowerCase()
  if (!q) return employees.value
  return employees.value.filter(e => `${e.first_name} ${e.last_name} ${e.email}`.toLowerCase().includes(q))
})
const selected = computed(() => employees.value.find(e => e.id === selectedId.value) || null)

function initials(e) {
  return ((e.first_name?.[0] || '') + (e.last_name?.[0] || '')).toUpperCase() || '?'
}
function shiftsByDay(weekday) {
  return shifts.value.filter(s => s.weekday === weekday).sort((a, b) => a.start_time.localeCompare(b.start_time))
}
function pad(n) { return String(n).padStart(2, '0') }
function formatDT(iso) {
  const d = new Date(iso)
  return `${pad(d.getDate())}/${pad(d.getMonth() + 1)} ${pad(d.getHours())}:${pad(d.getMinutes())}`
}
function formatTime(iso) {
  const d = new Date(iso)
  return `${pad(d.getHours())}:${pad(d.getMinutes())}`
}

async function fetchEmployees() {
  loading.value = true
  try { employees.value = (await api.get('/employees')).data.data || [] }
  catch { employees.value = [] } finally { loading.value = false }
}
async function selectEmployee(e) {
  selectedId.value = e.id
  await loadPlanning()
}
async function loadPlanning() {
  if (!selectedId.value) return
  const from = new Date(); from.setHours(0, 0, 0, 0)
  try {
    const [sh, ev] = await Promise.all([
      api.get('/staff/shifts', { params: { employee_id: selectedId.value } }),
      api.get('/staff/events', { params: { employee_id: selectedId.value, from: from.toISOString() } }),
    ])
    shifts.value = sh.data.data || []
    events.value = ev.data.data || []
  } catch { shifts.value = []; events.value = [] }
}

const empForm = ref({ show: false })
function openEmpForm(e) {
  empForm.value = e
    ? { show: true, id: e.id, first_name: e.first_name, last_name: e.last_name, email: e.email, password: '', status: e.status, saving: false, error: '' }
    : { show: true, id: null, first_name: '', last_name: '', email: '', password: '', status: 'active', saving: false, error: '' }
}
async function saveEmployee() {
  const f = empForm.value
  f.saving = true; f.error = ''
  try {
    if (f.id) {
      const payload = { first_name: f.first_name, last_name: f.last_name, status: f.status }
      if (f.password) payload.password = f.password
      await api.put(`/employees/${f.id}`, payload)
    } else {
      await api.post('/employees', { email: f.email, password: f.password, first_name: f.first_name, last_name: f.last_name })
    }
    f.show = false
    await fetchEmployees()
  } catch (e) { f.error = e.response?.data?.message || 'Erreur.' } finally { f.saving = false }
}
async function removeEmployee(e) {
  if (!confirm(`Supprimer l'employé ${e.first_name} ${e.last_name} ?`)) return
  try {
    await api.delete(`/employees/${e.id}`)
    if (selectedId.value === e.id) selectedId.value = null
    await fetchEmployees()
  } catch (err) { alert(err.response?.data?.message || 'Suppression impossible.') }
}

const shiftForm = ref({ show: false })
function openShiftForm(weekday) {
  shiftForm.value = { show: true, weekday, start_time: '09:00', end_time: '12:00', label: '', saving: false, error: '' }
}
async function saveShift() {
  const f = shiftForm.value
  if (f.end_time <= f.start_time) { f.error = 'La fin doit être après le début.'; return }
  f.saving = true; f.error = ''
  try {
    await api.post('/staff/shifts', { employee_id: selectedId.value, weekday: f.weekday, start_time: f.start_time, end_time: f.end_time, label: f.label || null })
    f.show = false
    await loadPlanning()
  } catch (e) { f.error = e.response?.data?.message || 'Erreur.' } finally { f.saving = false }
}
async function deleteShift(s) {
  await api.delete(`/staff/shifts/${s.id}`)
  await loadPlanning()
}

const eventForm = ref({ show: false, employee_ids: [] })
const eventEmpSearch = ref('')
const eventFilteredEmployees = computed(() => {
  const q = eventEmpSearch.value.trim().toLowerCase()
  if (!q) return employees.value
  return employees.value.filter(e => `${e.first_name} ${e.last_name} ${e.email}`.toLowerCase().includes(q))
})
const allSelected = computed(() => employees.value.length > 0 && eventForm.value.employee_ids.length === employees.value.length)
function toggleAll() {
  eventForm.value.employee_ids = allSelected.value ? [] : employees.value.map(e => e.id)
}
function openEventForm() {
  const d = new Date()
  const day = `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}`
  eventForm.value = { show: true, title: '', location: '', start_at: `${day}T09:00`, end_at: `${day}T11:00`, employee_ids: selectedId.value ? [selectedId.value] : [], saving: false, error: '' }
  eventEmpSearch.value = ''
}
async function saveEvent() {
  const f = eventForm.value
  if (!f.title.trim()) { f.error = 'Titre requis.'; return }
  if (!f.employee_ids.length) { f.error = 'Sélectionnez au moins un employé.'; return }
  f.saving = true; f.error = ''
  try {
    await api.post('/staff/events', { title: f.title.trim(), location: f.location || null, start_at: f.start_at, end_at: f.end_at, employee_ids: f.employee_ids })
    f.show = false
    await loadPlanning()
  } catch (e) { f.error = e.response?.data?.message || 'Erreur.' } finally { f.saving = false }
}
async function deleteEvent(ev) {
  if (!confirm(`Supprimer l'événement « ${ev.title} » ?`)) return
  await api.delete(`/staff/events/${ev.id}`)
  await loadPlanning()
}

onMounted(fetchEmployees)
</script>
