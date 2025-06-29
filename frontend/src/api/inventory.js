import apiClient from './client'

export const inventoryAPI = {
  // Get all inventory
  getAll: async () => {
    const response = await apiClient.get('/inventory')
    return response
  },

  // Get inventory by category ID
  getByCategoryId: async (categoryId) => {
    const response = await apiClient.get(`/inventory/${categoryId}`)
    return response
  }
}