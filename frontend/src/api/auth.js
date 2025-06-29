import apiClient from './client'

export const authAPI = {
  // Login
  login: async (credentials) => {
    const response = await apiClient.post('/auth/login', credentials)
    return response
  },

  // Logout (client-side only)
  logout: () => {
    localStorage.removeItem('token')
    localStorage.removeItem('user')
  }
}