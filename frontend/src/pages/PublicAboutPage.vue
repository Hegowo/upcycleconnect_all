<template>
  <div class="max-w-5xl mx-auto px-6 py-12 sm:py-16">

    <header class="text-center max-w-3xl mx-auto mb-14">
      <span class="inline-block text-xs font-bold uppercase tracking-wider text-[#006d35] bg-[#f0fdf4] px-3 py-1 rounded-full mb-4">Depuis 2021</span>
      <h1 class="text-3xl sm:text-4xl font-extrabold text-[#001d32] tracking-tight">Qui sommes-nous ?</h1>
      <p class="mt-4 text-lg text-[#40617f] leading-relaxed">
        UpcycleConnect, la plateforme de l'upcycling intelligent : donner une seconde vie aux objets, ensemble.
      </p>
    </header>

    <section class="prose-section space-y-5 text-[#334155] leading-relaxed">
      <p>
        Fondée en 2021, <strong>UpcycleConnect</strong> est une entreprise innovante dédiée à la réduction des déchets
        et à la valorisation des matériaux grâce à l'upcycling. Notre idée de départ : créer une communauté active où
        particuliers, artisans et entreprises peuvent donner une seconde vie à des objets et matériaux inutilisés ou obsolètes.
      </p>
      <p>
        Grâce à une plateforme web pensée pour faciliter la mise en relation, le suivi des projets et la mesure de l'impact
        environnemental, nous avons connu une croissance rapide. UpcycleConnect est aujourd'hui devenue une plateforme de
        référence pour l'upcycling en France et en Europe, au service d'une économie plus circulaire et durable.
      </p>
    </section>

    <section class="mt-14">
      <h2 class="text-2xl font-bold text-[#001d32] mb-6">Ce que nous faisons</h2>
      <div class="grid sm:grid-cols-2 gap-4">
        <div v-for="item in actions" :key="item.title" class="flex gap-3 bg-white rounded-2xl border border-[#e2e8f0] p-5">
          <div class="w-10 h-10 rounded-xl bg-[#f0fdf4] flex items-center justify-center shrink-0">
            <component :is="item.icon" class="w-5 h-5 text-[#006d35]" />
          </div>
          <div>
            <h3 class="font-semibold text-[#001d32] mb-1">{{ item.title }}</h3>
            <p class="text-sm text-[#64748b] leading-relaxed">{{ item.desc }}</p>
          </div>
        </div>
      </div>
    </section>

    <section class="mt-14">
      <h2 class="text-2xl font-bold text-[#001d32] mb-6">En quelques chiffres</h2>
      <div class="grid grid-cols-2 sm:grid-cols-4 gap-4">
        <div v-for="stat in stats" :key="stat.label" class="bg-white rounded-2xl border border-[#e2e8f0] p-5 text-center">
          <p class="text-2xl sm:text-3xl font-extrabold text-[#006d35]">{{ stat.value }}</p>
          <p class="text-xs text-[#64748b] mt-1 leading-snug">{{ stat.label }}</p>
        </div>
      </div>
    </section>

    <section class="mt-14">
      <h2 class="text-2xl font-bold text-[#001d32] mb-2">Nos implantations</h2>
      <p class="text-[#334155] leading-relaxed">
        Notre siège se situe au <strong>174, rue La Fayette, dans le 10<sup>e</sup> arrondissement de Paris</strong>, complété
        par des locaux dans les 11<sup>e</sup>, 13<sup>e</sup> et 16<sup>e</sup> arrondissements, ainsi qu'à Bourg-la-Reine,
        Ivry et Montreuil. Notre premier relais hors de France se trouve en Suisse, et de nouvelles annexes en province sont prévues.
      </p>
    </section>

    <section class="mt-14">
      <h2 class="text-2xl font-bold text-[#001d32] mb-2">La direction</h2>
      <p class="text-[#64748b] text-sm mb-6">L'équipe de direction du siège qui porte le projet UpcycleConnect au quotidien.</p>
      <div class="grid sm:grid-cols-2 lg:grid-cols-3 gap-4">
        <div v-for="m in leadership" :key="m.name" class="bg-white rounded-2xl border border-[#e2e8f0] p-5 flex items-center gap-4">
          <img v-if="m.photo" :src="m.photo" :alt="m.name" class="w-14 h-14 rounded-full object-cover shrink-0" />
          <div v-else class="w-14 h-14 rounded-full flex items-center justify-center text-white text-base font-bold shrink-0" :style="{ background: m.color }">
            {{ initials(m.name) }}
          </div>
          <div class="min-w-0">
            <p class="font-bold text-[#001d32] truncate">{{ m.name }}</p>
            <p class="text-sm text-[#006d35] font-medium">{{ m.role }}</p>
          </div>
        </div>
      </div>

      <h3 class="text-lg font-bold text-[#001d32] mt-8 mb-4">Sur nos autres sites</h3>
      <div class="grid sm:grid-cols-2 gap-4">
        <div v-for="m in sites" :key="m.name" class="bg-white rounded-2xl border border-[#e2e8f0] p-5 flex items-center gap-4">
          <img v-if="m.photo" :src="m.photo" :alt="m.name" class="w-14 h-14 rounded-full object-cover shrink-0" />
          <div v-else class="w-14 h-14 rounded-full flex items-center justify-center text-white text-base font-bold shrink-0" :style="{ background: m.color }">
            {{ initials(m.name) }}
          </div>
          <div class="min-w-0">
            <p class="font-bold text-[#001d32] truncate">{{ m.name }}</p>
            <p class="text-sm text-[#006d35] font-medium">{{ m.role }}</p>
          </div>
        </div>
      </div>
    </section>

    <section class="mt-14 bg-[#f0fdf4] rounded-2xl border border-[#bbf7d0] p-8 text-center">
      <h2 class="text-xl font-bold text-[#001d32] mb-2">Rejoignez le mouvement</h2>
      <p class="text-[#40617f] mb-5 max-w-xl mx-auto">Donnez une seconde vie à vos objets ou proposez vos talents d'artisan : la communauté UpcycleConnect vous attend.</p>
      <div class="flex flex-wrap items-center justify-center gap-3">
        <RouterLink to="/inscription" class="px-6 py-2.5 rounded-full text-white font-bold text-sm hover:opacity-90 transition" style="background:#006d35;">S'inscrire</RouterLink>
        <RouterLink to="/prestations" class="px-6 py-2.5 rounded-full font-bold text-sm text-[#006d35] border-2 border-[#006d35]/20 hover:border-[#006d35]/50 transition">Découvrir les prestations</RouterLink>
      </div>
    </section>
  </div>
</template>

<script setup>
import { RouterLink } from 'vue-router'
import {
  MegaphoneIcon, ArchiveBoxIcon, MagnifyingGlassIcon, ChartBarIcon,
  UsersIcon, AcademicCapIcon,
} from '@heroicons/vue/24/outline'
import photos from './aboutPhotos'

const actions = [
  { title: 'Annonces & dépôts', desc: 'Déposez vos objets dans nos points de collecte après validation, ou publiez une annonce de don ou de vente.', icon: MegaphoneIcon },
  { title: 'Inventaire intelligent', desc: 'Un suivi de la disponibilité et de l\'état des objets mis en ligne, avec des filtres de recherche avancés.', icon: ArchiveBoxIcon },
  { title: 'Suivi de projets', desc: 'Documentez chaque étape de transformation de vos objets et partagez vos créations avec la communauté.', icon: MagnifyingGlassIcon },
  { title: 'Upcycling Score', desc: 'Mesurez votre impact : quantité de déchets évités et économies de ressources réalisées.', icon: ChartBarIcon },
  { title: 'Communauté', desc: 'Forums, ateliers et événements pour échanger les bonnes pratiques et encourager l\'innovation.', icon: UsersIcon },
  { title: 'Sensibilisation', desc: 'Des ressources pédagogiques et des ateliers pour sensibiliser aux bonnes pratiques du développement durable.', icon: AcademicCapIcon },
]

const stats = [
  { value: '2021', label: 'Année de création' },
  { value: '260', label: 'Collaborateurs au siège' },
  { value: 'France & Europe', label: 'Présence' },
  { value: '100 %', label: 'Économie circulaire' },
]

const leadership = [
  { name: 'Sylvain Levy', role: 'Président-Directeur Général', color: '#006d35', photo: photos.sylvain_levy },
  { name: 'Pierre Chabrier', role: 'Directeur des Ressources Humaines', color: '#1b8848', photo: photos.pierre_chabrier },
  { name: 'Ronnand Peuplus', role: 'Directeur Commercial', color: '#40617f', photo: photos.ronnand_peuplus },
  { name: 'Laink Terracid', role: 'Directrice Marketing', color: '#7c3aed', photo: photos.laink_terracid },
  { name: 'Norman Thavaud', role: 'Directeur des Systèmes d\'Information', color: '#0e7490', photo: photos.norman_thavaud },
]

const sites = [
  { name: 'Antoine Maclair', role: 'Happiness Manager — Montreuil', color: '#b45309', photo: photos.antoine_maclair },
  { name: 'Frédéric Molas', role: 'Directeur régional — Suisse', color: '#001d32', photo: photos.frederic_molas },
]

function initials(name) {
  return name.split(' ').map(p => p[0]).join('').slice(0, 2).toUpperCase()
}
</script>
