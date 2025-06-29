import apiClient from './client'

export const categoriesAPI = {
  // Get all categories
  getAll: async () => {
    const response = await apiClient.get('/categories')
    return response
  },

  // Get category by ID
  getById: async (id) => {
    const response = await apiClient.get(`/categories/${id}`)
    return response
  },

  // Create category
  create: async (categoryData) => {
    const response = await apiClient.post('/categories', categoryData)
    return response
  },

  // Update category
  update: async (id, categoryData) => {
    const response = await apiClient.put(`/categories/${id}`, categoryData)
    return response
  },

  // Delete category
  delete: async (id) => {
    const response = await apiClient.delete(`/categories/${id}`)
    return response
  }
}