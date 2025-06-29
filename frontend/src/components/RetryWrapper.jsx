import React, { useState, useEffect } from 'react'
import { Result, Button, Spin, Alert } from 'antd'
import { ReloadOutlined, ExclamationCircleOutlined } from '@ant-design/icons'

const RetryWrapper = ({
  children,
  onRetry,
  maxRetries = 3,
  initialLoading = true,
  loadingText = '加载中...',
  errorTitle = '加载失败',
  errorSubTitle = '数据加载失败，请重试',
  showRetryCount = false,
  minHeight = 300,
  size = 'default' // 'small', 'default', 'large'
}) => {
  const [loading, setLoading] = useState(initialLoading)
  const [error, setError] = useState(null)
  const [retryCount, setRetryCount] = useState(0)

  // Execute the retry function
  const executeRetry = async () => {
    if (!onRetry) return

    setLoading(true)
    setError(null)

    try {
      await onRetry()
      setRetryCount(0) // Reset retry count on success
    } catch (err) {
      setError(err)
      console.error('RetryWrapper error:', err)
    } finally {
      setLoading(false)
    }
  }

  // Handle retry button click
  const handleRetry = () => {
    if (retryCount < maxRetries) {
      setRetryCount(prev => prev + 1)
      executeRetry()
    }
  }

  // Initial load
  useEffect(() => {
    if (initialLoading && onRetry) {
      executeRetry()
    }
  }, []) // eslint-disable-line react-hooks/exhaustive-deps

  // Provide context to children
  const contextValue = {
    loading,
    error,
    retryCount,
    setLoading,
    setError,
    retry: executeRetry
  }

  // Loading state
  if (loading) {
    const spinSize = size === 'large' ? 'large' : size === 'small' ? 'small' : 'default'
    
    return (
      <div style={{
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        minHeight,
        flexDirection: 'column'
      }}>
        <Spin size={spinSize} tip={loadingText} />
      </div>
    )
  }

  // Error state
  if (error) {
    const canRetry = retryCount < maxRetries
    
    return (
      <div style={{ 
        padding: '20px',
        textAlign: 'center',
        minHeight 
      }}>
        <Result
          status="error"
          icon={<ExclamationCircleOutlined style={{ color: '#ff4d4f' }} />}
          title={errorTitle}
          subTitle={
            <div>
              <div>{errorSubTitle}</div>
              {showRetryCount && retryCount > 0 && (
                <div style={{ marginTop: '8px', color: '#666' }}>
                  已重试 {retryCount} 次
                </div>
              )}
              <div style={{ marginTop: '8px', color: '#ff4d4f', fontSize: '14px' }}>
                {error.message || '未知错误'}
              </div>
            </div>
          }
          extra={[
            canRetry && (
              <Button 
                key="retry" 
                type="primary" 
                icon={<ReloadOutlined />}
                onClick={handleRetry}
              >
                重试 {retryCount > 0 ? `(${retryCount}/${maxRetries})` : ''}
              </Button>
            ),
            <Button 
              key="details" 
              onClick={() => {
                console.error('Error details:', error)
              }}
            >
              查看控制台详情
            </Button>
          ].filter(Boolean)}
        />
        
        {!canRetry && (
          <Alert
            message="重试次数已达上限"
            description="请检查网络连接或联系技术支持"
            type="warning"
            showIcon
            style={{ marginTop: '16px' }}
          />
        )}
      </div>
    )
  }

  // Success state - render children with context
  return typeof children === 'function' 
    ? children(contextValue)
    : React.cloneElement(children, { retryContext: contextValue })
}

export default RetryWrapper