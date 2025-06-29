import apiClient from './client'

export const usersAPI = {
  // Get all users
  getAll: async () => {
    const response = await apiClient.get('/users')
    return response
  },

  // Get user by ID
  getById: async (id) => {
    const response = await apiClient.get(`/users/${id}`)
    return response
  },

  // Create user
  create: async (userData) => {
    const response = await apiClient.post('/users', userData)
    return response
  },

  // Update user
  update: async (id, userData) => {
    const response = await apiClient.put(`/users/${id}`, userData)
    return response
  },

  // Delete user
  delete: async (id) => {
    const response = await apiClient.delete(`/users/${id}`)
    return response
  }
}