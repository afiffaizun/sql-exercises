<script setup>
import { ref, watch } from 'vue'

const props = defineProps({
  show: { type: Boolean, default: false },
  produk: { type: Object, default: null },
  categories: { type: Array, default: () => [] },
  saving: { type: Boolean, default: false }
})

const emit = defineEmits(['close', 'submit'])

const isEdit = ref(false)
const errors = ref({})
const form = ref({
  nama: '',
  kategori: '',
  harga: '',
  stok: '',
  deskripsi: ''
})

const kategoriOptions = ['Elektronik', 'Aksesoris', 'Makanan', 'Furniture', 'Pakaian']

watch(
  () => props.show,
  (val) => {
    if (!val) return
    errors.value = {}
    if (props.produk) {
      isEdit.value = true
      form.value = {
        nama: props.produk.nama,
        kategori: props.produk.kategori,
        harga: props.produk.harga,
        stok: props.produk.stok,
        deskripsi: props.produk.deskripsi
      }
    } else {
      isEdit.value = false
      form.value = { nama: '', kategori: '', harga: '', stok: '', deskripsi: '' }
    }
  }
)

function validate() {
  const e = {}
  if (!form.value.nama || !form.value.nama.trim()) e.nama = 'Nama produk wajib diisi'
  if (form.value.harga === '' || Number(form.value.harga) <= 0) e.harga = 'Harga harus lebih dari 0'
  if (form.value.stok === '' || Number(form.value.stok) < 0) e.stok = 'Stok tidak boleh negatif'
  errors.value = e
  return Object.keys(e).length === 0
}

function submit() {
  if (!validate()) return
  emit('submit', {
    nama: form.value.nama.trim(),
    kategori: form.value.kategori.trim(),
    harga: Number(form.value.harga),
    stok: Number(form.value.stok),
    deskripsi: form.value.deskripsi.trim()
  })
}
</script>

<template>
  <transition name="fade">
    <div
      v-if="show"
      class="fixed inset-0 z-50 flex items-center justify-center bg-slate-900/50 p-4 backdrop-blur-sm"
      @click.self="emit('close')"
    >
      <div class="w-full max-w-lg animate-pop-in rounded-2xl bg-white shadow-2xl">
        <div class="flex items-center justify-between border-b border-slate-100 px-6 py-4">
          <div>
            <h2 class="text-lg font-bold text-slate-800">
              {{ isEdit ? 'Edit Produk' : 'Tambah Produk' }}
            </h2>
            <p class="text-xs text-slate-400">
              {{ isEdit ? 'Perbarui informasi produk' : 'Lengkapi data produk baru' }}
            </p>
          </div>
          <button
            type="button"
            class="rounded-lg p-2 text-slate-400 transition hover:bg-slate-100 hover:text-slate-600"
            @click="emit('close')"
          >
            <svg class="h-5 w-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path stroke-linecap="round" d="M6 6l12 12M18 6L6 18" />
            </svg>
          </button>
        </div>

        <form class="space-y-4 px-6 py-5" @submit.prevent="submit">
          <div>
            <label class="label">Nama Produk <span class="text-red-500">*</span></label>
            <input v-model="form.nama" class="input" placeholder="cth. Laptop ASUS" />
            <p v-if="errors.nama" class="mt-1 text-xs text-red-500">{{ errors.nama }}</p>
          </div>

          <div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
            <div>
              <label class="label">Kategori</label>
              <input
                v-model="form.kategori"
                class="input"
                list="kategori-list"
                placeholder="cth. Elektronik"
              />
              <datalist id="kategori-list">
                <option v-for="k in [...new Set([...kategoriOptions, ...categories])]" :key="k" :value="k" />
              </datalist>
            </div>
            <div>
              <label class="label">Stok <span class="text-red-500">*</span></label>
              <input v-model="form.stok" type="number" min="0" class="input" placeholder="0" />
              <p v-if="errors.stok" class="mt-1 text-xs text-red-500">{{ errors.stok }}</p>
            </div>
          </div>

          <div>
            <label class="label">Harga (Rp) <span class="text-red-500">*</span></label>
            <input v-model="form.harga" type="number" min="0" step="any" class="input" placeholder="0" />
            <p v-if="errors.harga" class="mt-1 text-xs text-red-500">{{ errors.harga }}</p>
          </div>

          <div>
            <label class="label">Deskripsi</label>
            <textarea v-model="form.deskripsi" rows="3" class="input resize-none" placeholder="Deskripsi singkat produk"></textarea>
          </div>

          <div class="flex justify-end gap-2 border-t border-slate-100 pt-4">
            <button type="button" class="btn-ghost" @click="emit('close')">Batal</button>
            <button type="submit" class="btn-primary" :disabled="saving">
              <svg v-if="saving" class="h-4 w-4 animate-spin" viewBox="0 0 24 24" fill="none">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.4 0 0 5.4 0 12h4z" />
              </svg>
              {{ saving ? 'Menyimpan...' : isEdit ? 'Simpan Perubahan' : 'Simpan Produk' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </transition>
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
