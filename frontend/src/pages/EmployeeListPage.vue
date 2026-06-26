<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between gap-3 flex-wrap">
      <div>
        <h2 class="text-xl sm:text-2xl font-bold text-[#001d32]">Employés</h2>
        <p class="text-sm text-[#64748b] mt-0.5">Comptes du personnel à droits limités (formations, conseils, modération du forum).</p>
      </div>
      <button @click="openForm()" class="px-4 py-2 text-sm font-semibold rounded-lg text-white transition hover:opacity-90 flex items-center gap-2" style="background-color:#006d35;">
        <PlusIcon class="w-4 h-4" /> Nouvel employé
      </button>
    </div>

    <div v-if="showForm" class="card p-6">
      <h3 class="font-bold text-[#001d32] mb-4">{{ form.id ? 'Modifier l\'employé' : 'Nouvel employé' }}</h3>
      <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
        <div>
          <label class="label">Prénom</label>
          <input v-model="form.first_name" type="text" class="input" />
        </div>
        <div>
          <label class="label">Nom</label>
          <input v-model="form.last_name" type="text" class="input" />
        </div>
        <div v-if="!form.id">
          <label class="label">Email</label>
          <input v-model="form.email" type="email" class="input" />
        </div>
        <div>
          <label class="label">{{ form.id ? 'Nouveau mot de passe (optionnel)' : 'Mot de passe' }}</label>
          <input v-model="form.password" type="password" class="input" placeholder="Min. 8 caractères" />
        </div>
        <div v-if="form.id">
          <label class="label">Statut</label>
          <select v-model="form.status" class="input">
            <option value="active">Actif</option>
            <option value="suspended">Suspendu</option>
          </select>
        </div>
      </div>
      <p v-if="error" class="text-red-600 text-sm mt-3">{{ error }}</p>
      <div class="flex gap-3 mt-5">
        <button @click="save" :disabled="saving" class="px-5 py-2 rounded-lg text-white text-sm font-semibold hover:opacity-90 disabled:opacity-50" style="background-color:#1B8848;">
          {{ saving ? 'Enregistrement…' : 'Enregistrer' }}
        </button>
        <button @click="showForm = false" class="px-5 py-2 rounded-lg bg-gray-100 text-[#40617f] text-sm font-semibold hover:bg-gray-200">Annuler</button>
      </div>
    </div>

    <div v-if="loading" class="card p-10 text-center text-[#64748b] text-sm">Chargement…</div>
    <div v-else-if="!employees.length" class="card p-10 text-center text-[#64748b] text-sm">Aucun employé pour le moment.</div>
    <div v-else class="card overflow-hidden">
      <table class="w-full text-sm">
        <thead class="bg-[#f8fafc] text-[#64748b] text-xs uppercase">
          <tr>
            <th class="text-left px-5 py-3">Nom</th>
            <th class="text-left px-5 py-3">Email</th>
            <th class="text-left px-5 py-3">Statut</th>
            <th class="text-right px-5 py-3">Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="e in employees" :key="e.id" class="border-t border-[#f1f5f9]">
            <td class="px-5 py-3 font-semibold text-[#001d32]">{{ e.first_name }} {{ e.last_name }}</td>
            <td class="px-5 py-3 text-[#40617f]">{{ e.email }}</td>
            <td class="px-5 py-3">
              <span :class="e.status === 'active' ? 'bg-green-100 text-green-800' : 'bg-gray-200 text-gray-600'" class="text-xs font-bold px-2 py-0.5 rounded-full">
                {{ e.status === 'active' ? 'Actif' : 'Suspendu' }}
              </span>
            </td>
            <td class="px-5 py-3 text-right">
              <button @click="openForm(e)" class="text-xs font-semibold text-[#1a73e8] hover:underline mr-3">Modifier</button>
              <button @click="remove(e)" class="text-xs font-semibold text-red-600 hover:underline">Supprimer</button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { PlusIcon } from '@heroicons/vue/24/outline'
import api from '@/services/api'

const employees = ref([])
const loading = ref(true)
const showForm = ref(false)
const saving = ref(false)
const error = ref('')
const form = ref({ id: null, first_name: '', last_name: '', email: '', password: '', status: 'active' })

async function fetchEmployees() {
  loading.value = true
  try {
    const { data } = await api.get('/employees')
    employees.value = data.data || []
  } catch { employees.value = [] } finally { loading.value = false }
}

function openForm(e) {
  error.value = ''
  form.value = e
    ? { id: e.id, first_name: e.first_name, last_name: e.last_name, email: e.email, password: '', status: e.status }
    : { id: null, first_name: '', last_name: '', email: '', password: '', status: 'active' }
  showForm.value = true
}

async function save() {
  error.value = ''
  saving.value = true
  try {
    if (form.value.id) {
      const payload = { first_name: form.value.first_name, last_name: form.value.last_name, status: form.value.status }
      if (form.value.password) payload.password = form.value.password
      await api.put(`/employees/${form.value.id}`, payload)
    } else {
      await api.post('/employees', {
        email: form.value.email,
        password: form.value.password,
        first_name: form.value.first_name,
        last_name: form.value.last_name,
      })
    }
    showForm.value = false
    await fetchEmployees()
  } catch (e) {
    error.value = e.response?.data?.message || 'Erreur lors de l\'enregistrement.'
  } finally { saving.value = false }
}

async function remove(e) {
  if (!confirm(`Supprimer l'employé ${e.first_name} ${e.last_name} ?`)) return
  try {
    await api.delete(`/employees/${e.id}`)
    await fetchEmployees()
  } catch (err) { alert(err.response?.data?.message || 'Suppression impossible.') }
}

onMounted(fetchEmployees)
</script>
