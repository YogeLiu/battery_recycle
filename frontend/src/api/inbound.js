import apiClient from './client'

export const inboundAPI = {
  // Get all inbound orders
  getAll: async (params = {}) => {
    const response = await apiClient.get('/inbound/orders', { params })
    return response
  },

  // Get inbound order by ID
  getById: async (id) => {
    const response = await apiClient.get(`/inbound/orders/${id}`)
    return response
  },

  // Create inbound order
  create: async (orderData) => {
    const response = await apiClient.post('/inbound/orders', orderData)
    return response
  },

  // Update inbound order
  update: async (id, orderData) => {
    const response = await apiClient.put(`/inbound/orders/${id}`, orderData)
    return response
  },

  // Delete inbound order
  delete: async (id) => {
    const response = await apiClient.delete(`/inbound/orders/${id}`)
    return response
  }
}