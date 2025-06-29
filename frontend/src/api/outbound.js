import apiClient from './client'

export const outboundAPI = {
  // Get all outbound orders
  getAll: async (params = {}) => {
    const response = await apiClient.get('/outbound/orders', { params })
    return response
  },

  // Get outbound order by ID
  getById: async (id) => {
    const response = await apiClient.get(`/outbound/orders/${id}`)
    return response
  },

  // Create outbound order
  create: async (orderData) => {
    const response = await apiClient.post('/outbound/orders', orderData)
    return response
  },

  // Update outbound order
  update: async (id, orderData) => {
    const response = await apiClient.put(`/outbound/orders/${id}`, orderData)
    return response
  },

  // Delete outbound order
  delete: async (id) => {
    const response = await apiClient.delete(`/outbound/orders/${id}`)
    return response
  }
}