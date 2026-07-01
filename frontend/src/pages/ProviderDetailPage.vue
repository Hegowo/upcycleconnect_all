<template>
  <div class="max-w-2xl space-y-6">
    <RouterLink to="/admin/providers" class="inline-flex items-center gap-1 text-sm text-gray-500 hover:text-gray-700">
      {{ t('common.back') }}
    </RouterLink>

    <div v-if="loading" class="card p-8 space-y-4">
      <div class="flex items-center gap-4">
        <div class="w-16 h-16 rounded-full bg-gray-100 animate-pulse"></div>
        <div class="space-y-2">
          <div class="h-5 bg-gray-100 rounded animate-pulse w-44"></div>
          <div class="h-4 bg-gray-100 rounded animate-pulse w-60"></div>
        </div>
      </div>
      <div class="grid grid-cols-2 gap-4 pt-4">
        <div v-for="n in 4" :key="n" class="h-10 bg-gray-100 rounded animate-pulse"></div>
      </div>
    </div>

    <template v-else-if="provider">
      <div class="card p-8">
        <div class="flex items-start justify-between mb-6">
          <div class="flex items-center gap-4">
            <div class="w-16 h-16 rounded-full flex items-center justify-center text-white text-xl font-bold shrink-0" style="background-color:#1B5B88;">
              {{ (provider.profile?.company_name?.[0] || provider.first_name?.[0] || '').toUpperCase() }}{{ (provider.profile?.company_name?.split(' ')?.[1]?.[0] || provider.last_name?.[0] || '').toUpperCase() }}
            </div>
            <div>
              <h2 class="text-xl font-bold text-gray-900">{{ provider.profile?.company_name || `${provider.first_name} ${provider.last_name}` }}</h2>
              <p class="text-sm text-gray-500 mt-0.5">{{ provider.first_name }} {{ provider.last_name }} — {{ provider.email }}</p>
            </div>
          </div>
          <AppBadge :label="provider.profile?.status || '—'" />
        </div>

        <dl class="grid grid-cols-1 sm:grid-cols-2 gap-x-6 gap-y-4 text-sm border-t border-gray-100 pt-6">
          <div>
            <dt class="text-xs text-gray-500 uppercase font-medium mb-1">{{ t('providers.fieldSiret') }}</dt>
            <dd class="text-gray-900 font-medium font-mono">{{ provider.profile?.siret || '—' }}</dd>
          </div>
          <div>
            <dt class="text-xs text-gray-500 uppercase font-medium mb-1">{{ t('providers.fieldWebsite') }}</dt>
            <dd class="font-medium">
              <a v-if="provider.profile?.website" :href="provider.profile.website" target="_blank" rel="noopener" class="text-blue-600 hover:underline truncate block">
                {{ provider.profile.website }}
              </a>
              <span v-else class="text-gray-500">—</span>
            </dd>
          </div>
          <div class="col-span-2">
            <dt class="text-xs text-gray-500 uppercase font-medium mb-1">{{ t('providers.fieldDescription') }}</dt>
            <dd class="text-gray-700 leading-relaxed">{{ provider.profile?.description || t('providers.noDescription') }}</dd>
          </div>

          <div class="col-span-2">
            <dt class="text-xs text-gray-500 uppercase font-medium mb-2">Document officiel (Kbis)</dt>
            <dd>
              <button
                v-if="provider.profile?.has_kbis"
                @click="consultKbis"
                class="inline-flex items-center gap-2 px-4 py-2 rounded-xl border border-[#006d35] text-[#006d35] text-sm font-semibold hover:bg-[#f0fdf4] transition"
              >
                <DocumentArrowDownIcon class="w-4 h-4" />
                Consulter le Kbis
              </button>
              <span v-else class="inline-flex items-center gap-1.5 text-amber-600 text-sm">
                <ExclamationTriangleIcon class="w-4 h-4" />
                Aucun Kbis déposé
              </span>
            </dd>
          </div>
        </dl>

        <div class="flex gap-3 pt-6 mt-2 border-t border-gray-100">
          <button
            v-if="provider.profile?.status !== 'approved'"
            @click="changeStatus('approved')"
            class="inline-flex items-center gap-2 px-4 py-2 rounded-xl text-white text-sm font-semibold transition hover:opacity-90"
            style="background-color:#1B8848;"
          >
            {{ t('providers.actionApprove') }}
          </button>
          <button
            v-if="provider.profile?.status !== 'suspended'"
            @click="changeStatus('suspended')"
            class="inline-flex items-center gap-2 px-4 py-2 rounded-xl text-white text-sm font-semibold transition hover:opacity-90"
            style="background-color:#dc2626;"
          >
            {{ t('providers.actionSuspend') }}
          </button>
        </div>

        <div class="mt-6 pt-4 border-t border-red-100">
          <div class="flex items-start gap-3 rounded-xl bg-red-50 border border-red-200 p-4">
            <TrashIcon class="w-5 h-5 text-red-600 shrink-0 mt-0.5" />
            <div class="flex-1 min-w-0">
              <p class="text-sm font-semibold text-red-800">Supprimer toutes les données (RGPD)</p>
              <p class="text-xs text-red-700/80 mt-0.5">Suppression définitive et irréversible du prestataire et de toutes ses données. Un e-mail de confirmation RGPD lui sera envoyé.</p>
            </div>
            <button
              @click="purgeConfirm = true"
              class="shrink-0 inline-flex items-center gap-1.5 px-3 py-2 rounded-lg text-white text-xs font-bold transition hover:opacity-90"
              style="background-color:#b91c1c;"
            >
              <TrashIcon class="w-3.5 h-3.5" /> Supprimer les données
            </button>
          </div>
        </div>
      </div>
    </template>

    <AppConfirmDialog
      :show="purgeConfirm"
      title="Supprimer toutes les données ?"
      :message="purgeMessage"
      confirm-label="Supprimer définitivement"
      confirm-variant="danger"
      :loading="purgeLoading"
      @confirm="executePurge"
      @cancel="purgeConfirm = false"
    />
  </div>
</template>

<script setup>import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { providerService } from '@/services/userService'
import api from '@/services/api'
import { useToast } from '@/utils/useToast'
import AppBadge from '@/components/AppBadge.vue'
import AppConfirmDialog from '@/components/AppConfirmDialog.vue'
import { DocumentArrowDownIcon, ExclamationTriangleIcon, TrashIcon } from '@heroicons/vue/24/outline'

const { t } = useI18n()
const route    = useRoute()
const router   = useRouter()
const toast    = useToast()
const provider = ref(null)
const loading  = ref(true)
const purgeConfirm = ref(false)
const purgeLoading = ref(false)
const purgeMessage = computed(() => {
  const name = provider.value ? `${provider.value.first_name} ${provider.value.last_name}` : ''
  return `Cette action supprimera DÉFINITIVEMENT ${name} et l'intégralité de ses données (prestations, réservations, campagnes, messages, etc.). Elle est irréversible et un e-mail de confirmation RGPD sera envoyé. Confirmez-vous ?`
})

async function executePurge() {
  purgeLoading.value = true
  try {
    await api.post(`/users/${provider.value.id}/purge`)
    toast.showSuccess('Toutes les données du prestataire ont été supprimées.')
    purgeConfirm.value = false
    router.push('/admin/providers')
  } catch (e) {
    toast.showError(e?.response?.data?.message || 'Erreur lors de la suppression.')
  } finally {
    purgeLoading.value = false
  }
}

async function fetchProvider() {
  try {
    provider.value = await providerService.get(route.params.id)
  } catch {
    toast.showError(t('providers.toastLoadError'))
    router.push('/admin/providers')
  } finally {
    loading.value = false
  }
}

async function consultKbis() {
  try {
    const blob = await providerService.downloadKbis(provider.value.id)
    const url = URL.createObjectURL(blob)
    window.open(url, '_blank', 'noopener')
    setTimeout(() => URL.revokeObjectURL(url), 60000)
  } catch {
    toast.showError('Impossible de charger le Kbis.')
  }
}

async function changeStatus(status) {
  try {
    await providerService.updateStatus(provider.value.id, status)
    provider.value.profile.status = status
    toast.showSuccess(t('providers.toastStatusUpdated', { status }))
  } catch {
    toast.showError(t('providers.toastStatusError'))
  }
}

onMounted(fetchProvider)
</script>
