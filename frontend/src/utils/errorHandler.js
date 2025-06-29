import { message, notification } from 'antd'

// Error types
export const ERROR_TYPES = {
  NETWORK: 'NETWORK',
  VALIDATION: 'VALIDATION',
  AUTHENTICATION: 'AUTHENTICATION',
  AUTHORIZATION: 'AUTHORIZATION',
  SERVER: 'SERVER',
  UNKNOWN: 'UNKNOWN'
}

// Determine error type from error object
export const getErrorType = (error) => {
  if (!error) return ERROR_TYPES.UNKNOWN

  // Network errors
  if (error.code === 'NETWORK_ERROR' || error.message?.includes('Network Error')) {
    return ERROR_TYPES.NETWORK
  }

  // HTTP status based errors
  if (error.response?.status) {
    const status = error.response.status
    if (status === 401) return ERROR_TYPES.AUTHENTICATION
    if (status === 403) return ERROR_TYPES.AUTHORIZATION
    if (status >= 400 && status < 500) return ERROR_TYPES.VALIDATION
    if (status >= 500) return ERROR_TYPES.SERVER
  }

  return ERROR_TYPES.UNKNOWN
}

// Get user-friendly error message
export const getErrorMessage = (error, context = '') => {
  const errorType = getErrorType(error)
  const contextPrefix = context ? `${context}: ` : ''

  // Extract server error message if available
  const serverMessage = error.response?.data?.msg || 
                       error.response?.data?.message || 
                       error.message

  switch (errorType) {
    case ERROR_TYPES.NETWORK:
      return `${contextPrefix}网络连接失败，请检查网络连接后重试`
    
    case ERROR_TYPES.AUTHENTICATION:
      return `${contextPrefix}身份验证失败，请重新登录`
    
    case ERROR_TYPES.AUTHORIZATION:
      return `${contextPrefix}权限不足，无法执行此操作`
    
    case ERROR_TYPES.VALIDATION:
      return `${contextPrefix}${serverMessage || '请求参数错误'}`
    
    case ERROR_TYPES.SERVER:
      return `${contextPrefix}服务器错误，请稍后重试`
    
    default:
      return `${contextPrefix}${serverMessage || '操作失败，请重试'}`
  }
}

// Show error message with appropriate method
export const showErrorMessage = (error, context = '', options = {}) => {
  const {
    useNotification = false,
    duration = 4,
    key,
    placement = 'topRight'
  } = options

  const errorMessage = getErrorMessage(error, context)
  const errorType = getErrorType(error)

  if (useNotification) {
    const notificationType = errorType === ERROR_TYPES.NETWORK ? 'warning' : 'error'
    
    notification[notificationType]({
      message: context || '操作失败',
      description: errorMessage,
      duration,
      key,
      placement
    })
  } else {
    message.error(errorMessage, duration, key)
  }

  // Log error for debugging
  console.error(`Error in ${context}:`, error)
}

// Retry wrapper with exponential backoff
export const withRetry = async (fn, maxRetries = 3, baseDelay = 1000) => {
  let lastError

  for (let attempt = 0; attempt <= maxRetries; attempt++) {
    try {
      return await fn()
    } catch (error) {
      lastError = error
      
      // Don't retry on client errors (4xx) except 408, 429
      const status = error.response?.status
      if (status >= 400 && status < 500 && status !== 408 && status !== 429) {
        throw error
      }

      // Don't retry on last attempt
      if (attempt === maxRetries) {
        throw error
      }

      // Exponential backoff with jitter
      const delay = baseDelay * Math.pow(2, attempt) + Math.random() * 1000
      await new Promise(resolve => setTimeout(resolve, delay))
    }
  }

  throw lastError
}

// Create error handler for specific contexts
export const createErrorHandler = (context, options = {}) => {
  return (error) => {
    showErrorMessage(error, context, options)
    return Promise.reject(error)
  }
}

// Global error handlers for common scenarios
export const errorHandlers = {
  // Data loading errors
  dataLoad: createErrorHandler('数据加载失败'),
  
  // Form submission errors
  formSubmit: createErrorHandler('提交失败'),
  
  // Delete operation errors
  delete: createErrorHandler('删除失败'),
  
  // Export errors
  export: createErrorHandler('导出失败'),
  
  // Print errors
  print: createErrorHandler('打印失败'),
  
  // Authentication errors
  auth: createErrorHandler('认证失败', { useNotification: true }),
  
  // Network errors with notification
  network: createErrorHandler('网络错误', { 
    useNotification: true,
    duration: 6
  })
}

// Error boundary error handler
export const handleErrorBoundaryError = (error, errorInfo) => {
  console.error('ErrorBoundary caught error:', error, errorInfo)
  
  // Could send to error reporting service here
  // logErrorToService(error, errorInfo)
  
  notification.error({
    message: '应用程序错误',
    description: '页面遇到意外错误，已自动记录。如果问题持续出现，请联系技术支持。',
    duration: 8,
    placement: 'topRight'
  })
}

// API response interceptor error handler
export const handleAPIError = (error) => {
  const errorType = getErrorType(error)
  
  // Special handling for authentication errors
  if (errorType === ERROR_TYPES.AUTHENTICATION) {
    // Handled by axios interceptor - redirect to login
    return Promise.reject(error)
  }
  
  // For other errors, just log and pass through
  console.error('API Error:', error)
  return Promise.reject(error)
}

export default {
  getErrorType,
  getErrorMessage,
  showErrorMessage,
  withRetry,
  createErrorHandler,
  errorHandlers,
  handleErrorBoundaryError,
  handleAPIError,
  ERROR_TYPES
}