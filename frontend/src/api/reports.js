import apiClient from './client'

export const reportsAPI = {
  // Get summary report
  getSummary: async (params = {}) => {
    const response = await apiClient.get('/reports/summary', { params })
    return response
  }
}