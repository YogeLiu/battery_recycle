import axios from 'axios'

// Create axios instance
const apiClient = axios.create({
  baseURL: import.meta.env.VITE_API_URL || (import.meta.env.PROD ? 'http://backend:8036/jxc/v1' : '/jxc/v1'),
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// Request interceptor to add auth token
apiClient.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// Response interceptor to handle errors
apiClient.interceptors.response.use(
  (response) => {
    // Return data directly for successful responses
    return response.data
  },
  (error) => {
    if (error.response) {
      // Server responded with error status
      const { status, data } = error.response
      
      if (status === 401) {
        // Unauthorized - clear token and redirect to login
        localStorage.removeItem('token')
        localStorage.removeItem('user')
        window.location.href = '/login'
      }
      
      // Reject with the error message from server
      return Promise.reject(new Error(data.msg || 'Request failed'))
    } else if (error.request) {
      // Network error
      return Promise.reject(new Error('Network error - please check your connection'))
    } else {
      // Something else happened
      return Promise.reject(error)
    }
  }
)

export default apiClient