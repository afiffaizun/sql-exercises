<script setup>
import { computed, onMounted, ref, watch } from 'vue'
import api from './api'
import ProdukForm from './components/ProdukForm.vue'
import Toast from './components/Toast.vue'

const produkList = ref([])
const stats = ref({ total_produk: 0, total_stok: 0, stok_rendah: 0, nilai_stok: 0 })
const loading = ref(false)
const saving = ref(false)

const search = ref('')
const kategoriFilter = ref('')

const showForm = ref(false)
const editing = ref(null)

const deleteTarget = ref(null)
const deleting = ref(false)

const toasts = ref([])
let toastId = 0
let searchTimer = null

const categories = computed(() => {
  const set = new Set(produkList.value.map((p) => p.kategori).filter(Boolean))
  return [...set]
})

const badgeColors = [
  'bg-indigo-50 text-indigo-600 ring-indigo-200',
  'bg-emerald-50 text-emerald-600 ring-emerald-200',
  'bg-amber-50 text-amber-600 ring-amber-200',
  'bg-sky-50 text-sky-600 ring-sky-200',
  'bg-rose-50 text-rose-600 ring-rose-200',
  'bg-violet-50 text-violet-600 ring-violet-200'
]

function badgeClass(kategori) {
  if (!kategori) return 'bg-slate-100 text-slate-500 ring-slate-200'
  let hash = 0
  for (let i = 0; i < kategori.length; i++) hash = kategori.charCodeAt(i) + ((hash << 5) - hash)
  return badgeColors[Math.abs(hash) % badgeColors.length]
}

function formatRupiah(value) {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    minimumFractionDigits: 0
  }).format(value || 0)
}

function addToast(type, message) {
  const id = ++toastId
  toasts.value.push({ id, type, message })
  setTimeout(() => {
    toasts.value = toasts.value.filter((t) => t.id !== id)
  }, 3500)
}

async function fetchProduk() {
  loading.value = true
  try {
    const params = {}
    if (search.value) params.q = search.value
    if (kategoriFilter.value) params.kategori = kategoriFilter.value
    const { data } = await api.getProduk(params)
    produkList.value = data.data || []
  } catch (err) {
    addToast('error', 'Gagal memuat data produk')
  } finally {
    loading.value = false
  }
}

async function fetchStats() {
  try {
    const { data } = await api.getStats()
    stats.value = data
  } catch (err) {
    /* diamkan, statistik tidak kritis */
  }
}

function openCreate() {
  editing.value = null
  showForm.value = true
}

function openEdit(item) {
  editing.value = item
  showForm.value = true
}

async function handleSubmit(payload) {
  saving.value = true
  try {
    if (editing.value) {
      await api.updateProduk(editing.value.id, payload)
      addToast('success', 'Produk berhasil diperbarui')
    } else {
      await api.createProduk(payload)
      addToast('success', 'Produk berhasil ditambahkan')
    }
    showForm.value = false
    await Promise.all([fetchProduk(), fetchStats()])
  } catch (err) {
    const msg = err?.response?.data?.error || 'Terjadi kesalahan saat menyimpan'
    addToast('error', msg)
  } finally {
    saving.value = false
  }
}

function askDelete(item) {
  deleteTarget.value = item
}

async function doDelete() {
  if (!deleteTarget.value) return
  deleting.value = true
  try {
    await api.deleteProduk(deleteTarget.value.id)
    addToast('success', 'Produk berhasil dihapus')
    deleteTarget.value = null
    await Promise.all([fetchProduk(), fetchStats()])
  } catch (err) {
    addToast('error', 'Gagal menghapus produk')
  } finally {
    deleting.value = false
  }
}

watch(search, () => {
  clearTimeout(searchTimer)
  searchTimer = setTimeout(fetchProduk, 350)
})

watch(kategoriFilter, fetchProduk)

onMounted(() => {
  fetchProduk()
  fetchStats()
})
</script>

<template>
  <div class="min-h-screen">
    <!-- Header -->
    <header class="relative overflow-hidden bg-gradient-to-br from-indigo-600 via-indigo-500 to-violet-600 pb-24 pt-10 text-white">
      <div class="pointer-events-none absolute -right-16 -top-16 h-64 w-64 rounded-full bg-white/10 blur-2xl"></div>
      <div class="pointer-events-none absolute -bottom-24 left-1/3 h-72 w-72 rounded-full bg-violet-400/20 blur-3xl"></div>
      <div class="mx-auto flex max-w-6xl items-center justify-between px-4 sm:px-6">
        <div class="flex items-center gap-3">
          <div class="flex h-11 w-11 items-center justify-center rounded-2xl bg-white/15 text-2xl ring-1 ring-white/25 backdrop-blur">
            📦
          </div>
          <div>
            <h1 class="text-xl font-bold tracking-tight sm:text-2xl">Manajemen Produk</h1>
            <p class="text-sm text-indigo-100">Kelola data barang dengan mudah</p>
          </div>
        </div>
        <button class="btn bg-white/15 text-white ring-1 ring-white/25 backdrop-blur hover:bg-white/25" @click="openCreate">
          <svg class="h-4 w-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
            <path stroke-linecap="round" d="M12 5v14M5 12h14" />
          </svg>
          <span class="hidden sm:inline">Tambah Produk</span>
        </button>
      </div>
    </header>

    <main class="mx-auto -mt-16 max-w-6xl px-4 pb-16 sm:px-6">
      <!-- Stat cards -->
      <section class="grid grid-cols-2 gap-4 lg:grid-cols-4">
        <div class="rounded-2xl border border-slate-100 bg-white p-5 shadow-sm">
          <div class="flex items-center justify-between">
            <p class="text-xs font-medium uppercase tracking-wide text-slate-400">Total Produk</p>
            <span class="rounded-lg bg-indigo-50 p-1.5 text-indigo-500">
              <svg class="h-4 w-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
              </svg>
            </span>
          </div>
          <p class="mt-3 text-2xl font-bold text-slate-800">{{ stats.total_produk }}</p>
        </div>

        <div class="rounded-2xl border border-slate-100 bg-white p-5 shadow-sm">
          <div class="flex items-center justify-between">
            <p class="text-xs font-medium uppercase tracking-wide text-slate-400">Total Stok</p>
            <span class="rounded-lg bg-emerald-50 p-1.5 text-emerald-500">
              <svg class="h-4 w-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M3 3v18h18M7 15l4-4 3 3 5-6" />
              </svg>
            </span>
          </div>
          <p class="mt-3 text-2xl font-bold text-slate-800">{{ stats.total_stok }}</p>
        </div>

        <div class="rounded-2xl border border-slate-100 bg-white p-5 shadow-sm">
          <div class="flex items-center justify-between">
            <p class="text-xs font-medium uppercase tracking-wide text-slate-400">Stok Rendah</p>
            <span class="rounded-lg bg-rose-50 p-1.5 text-rose-500">
              <svg class="h-4 w-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M12 9v4m0 4h.01M10.29 3.86L1.82 18a2 2 0 001.71 3h16.94a2 2 0 001.71-3L13.71 3.86a2 2 0 00-3.42 0z" />
              </svg>
            </span>
          </div>
          <p class="mt-3 text-2xl font-bold text-slate-800">{{ stats.stok_rendah }}</p>
        </div>

        <div class="rounded-2xl border border-slate-100 bg-white p-5 shadow-sm">
          <div class="flex items-center justify-between">
            <p class="text-xs font-medium uppercase tracking-wide text-slate-400">Nilai Stok</p>
            <span class="rounded-lg bg-amber-50 p-1.5 text-amber-500">
              <svg class="h-4 w-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M12 1v22M17 5H9.5a3.5 3.5 0 000 7h5a3.5 3.5 0 010 7H6" />
              </svg>
            </span>
          </div>
          <p class="mt-3 truncate text-2xl font-bold text-slate-800">{{ formatRupiah(stats.nilai_stok) }}</p>
        </div>
      </section>

      <!-- Table card -->
      <section class="mt-6 rounded-2xl border border-slate-100 bg-white shadow-sm">
        <div class="flex flex-col gap-3 border-b border-slate-100 p-4 sm:flex-row sm:items-center sm:justify-between">
          <div class="relative w-full sm:max-w-xs">
            <svg class="pointer-events-none absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-slate-400" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="11" cy="11" r="7" /><path stroke-linecap="round" d="M21 21l-4.3-4.3" />
            </svg>
            <input v-model="search" class="input pl-9" placeholder="Cari produk atau kategori..." />
          </div>
          <select v-model="kategoriFilter" class="input w-full sm:w-48">
            <option value="">Semua Kategori</option>
            <option v-for="k in categories" :key="k" :value="k">{{ k }}</option>
          </select>
        </div>

        <!-- Loading skeleton -->
        <div v-if="loading" class="divide-y divide-slate-50">
          <div v-for="n in 5" :key="n" class="flex items-center gap-4 p-4">
            <div class="h-10 w-10 animate-pulse rounded-xl bg-slate-100"></div>
            <div class="flex-1 space-y-2">
              <div class="h-3 w-1/3 animate-pulse rounded bg-slate-100"></div>
              <div class="h-3 w-1/5 animate-pulse rounded bg-slate-100"></div>
            </div>
            <div class="h-3 w-20 animate-pulse rounded bg-slate-100"></div>
          </div>
        </div>

        <!-- Empty state -->
        <div v-else-if="produkList.length === 0" class="flex flex-col items-center justify-center px-4 py-16 text-center">
          <div class="mb-3 text-4xl">🗂️</div>
          <p class="font-semibold text-slate-700">Belum ada produk</p>
          <p class="mt-1 text-sm text-slate-400">
            {{ search || kategoriFilter ? 'Coba ubah kata kunci pencarian.' : 'Tambahkan produk pertama Anda sekarang.' }}
          </p>
          <button v-if="!search && !kategoriFilter" class="btn-primary mt-4" @click="openCreate">Tambah Produk</button>
        </div>

        <!-- Table -->
        <div v-else class="overflow-x-auto">
          <table class="w-full min-w-[640px] text-left text-sm">
            <thead>
              <tr class="border-b border-slate-100 text-xs uppercase tracking-wide text-slate-400">
                <th class="px-4 py-3 font-semibold">Produk</th>
                <th class="px-4 py-3 font-semibold">Kategori</th>
                <th class="px-4 py-3 text-right font-semibold">Harga</th>
                <th class="px-4 py-3 text-center font-semibold">Stok</th>
                <th class="px-4 py-3 text-right font-semibold">Aksi</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-50">
              <tr v-for="item in produkList" :key="item.id" class="group transition hover:bg-slate-50/70">
                <td class="px-4 py-3">
                  <div class="flex items-center gap-3">
                    <div class="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl bg-gradient-to-br from-indigo-500 to-violet-500 text-sm font-bold text-white">
                      {{ (item.nama || '?').charAt(0).toUpperCase() }}
                    </div>
                    <div class="min-w-0">
                      <p class="truncate font-semibold text-slate-700">{{ item.nama }}</p>
                      <p class="truncate text-xs text-slate-400">{{ item.deskripsi || 'Tanpa deskripsi' }}</p>
                    </div>
                  </div>
                </td>
                <td class="px-4 py-3">
                  <span class="inline-flex rounded-full px-2.5 py-1 text-xs font-medium ring-1 ring-inset" :class="badgeClass(item.kategori)">
                    {{ item.kategori || 'Tanpa kategori' }}
                  </span>
                </td>
                <td class="px-4 py-3 text-right font-medium text-slate-700">{{ formatRupiah(item.harga) }}</td>
                <td class="px-4 py-3 text-center">
                  <span
                    class="inline-flex min-w-[2.5rem] justify-center rounded-lg px-2 py-1 text-xs font-semibold"
                    :class="item.stok < 10 ? 'bg-rose-50 text-rose-600' : 'bg-emerald-50 text-emerald-600'"
                  >
                    {{ item.stok }}
                  </span>
                </td>
                <td class="px-4 py-3">
                  <div class="flex justify-end gap-1">
                    <button
                      class="rounded-lg p-2 text-slate-400 transition hover:bg-indigo-50 hover:text-indigo-600"
                      title="Edit"
                      @click="openEdit(item)"
                    >
                      <svg class="h-4 w-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <path d="M11 4H4a2 2 0 00-2 2v14a2 2 0 002 2h14a2 2 0 002-2v-7" />
                        <path d="M18.5 2.5a2.12 2.12 0 013 3L12 15l-4 1 1-4 9.5-9.5z" />
                      </svg>
                    </button>
                    <button
                      class="rounded-lg p-2 text-slate-400 transition hover:bg-rose-50 hover:text-rose-600"
                      title="Hapus"
                      @click="askDelete(item)"
                    >
                      <svg class="h-4 w-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <path d="M3 6h18M8 6V4a2 2 0 012-2h4a2 2 0 012 2v2m2 0v14a2 2 0 01-2 2H7a2 2 0 01-2-2V6" />
                        <path stroke-linecap="round" d="M10 11v6M14 11v6" />
                      </svg>
                    </button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <div v-if="!loading && produkList.length > 0" class="border-t border-slate-100 px-4 py-3 text-xs text-slate-400">
          Menampilkan {{ produkList.length }} produk
        </div>
      </section>
    </main>

    <!-- Form modal -->
    <ProdukForm
      :show="showForm"
      :produk="editing"
      :categories="categories"
      :saving="saving"
      @close="showForm = false"
      @submit="handleSubmit"
    />

    <!-- Delete confirm -->
    <transition name="fade">
      <div
        v-if="deleteTarget"
        class="fixed inset-0 z-50 flex items-center justify-center bg-slate-900/50 p-4 backdrop-blur-sm"
        @click.self="deleteTarget = null"
      >
        <div class="w-full max-w-sm animate-pop-in rounded-2xl bg-white p-6 text-center shadow-2xl">
          <div class="mx-auto mb-3 flex h-12 w-12 items-center justify-center rounded-full bg-rose-50 text-rose-500">
            <svg class="h-6 w-6" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M12 9v4m0 4h.01M10.29 3.86L1.82 18a2 2 0 001.71 3h16.94a2 2 0 001.71-3L13.71 3.86a2 2 0 00-3.42 0z" />
            </svg>
          </div>
          <h3 class="text-lg font-bold text-slate-800">Hapus produk?</h3>
          <p class="mt-1 text-sm text-slate-500">
            Produk <span class="font-semibold text-slate-700">{{ deleteTarget.nama }}</span> akan dihapus permanen.
          </p>
          <div class="mt-5 flex justify-center gap-2">
            <button class="btn-ghost" @click="deleteTarget = null">Batal</button>
            <button
              class="btn bg-rose-600 text-white shadow-sm hover:bg-rose-700 focus:ring-rose-400"
              :disabled="deleting"
              @click="doDelete"
            >
              {{ deleting ? 'Menghapus...' : 'Ya, Hapus' }}
            </button>
          </div>
        </div>
      </div>
    </transition>

    <Toast :toasts="toasts" />
  </div>
</template>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
