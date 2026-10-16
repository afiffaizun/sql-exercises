import axios from 'axios'

const api = axios.create({
  baseURL: '/api',
  headers: { 'Content-Type': 'application/json' }
})

export default {
  getProduk(params) {
    return api.get('/produk', { params })
  },
  getStats() {
    return api.get('/stats')
  },
  createProduk(data) {
    return api.post('/produk', data)
  },
  updateProduk(id, data) {
    return api.put(`/produk/${id}`, data)
  },
  deleteProduk(id) {
    return api.delete(`/produk/${id}`)
  }
}
