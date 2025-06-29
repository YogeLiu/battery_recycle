import React from 'react'
import { Result, Button } from 'antd'
import { ExclamationCircleOutlined, ReloadOutlined } from '@ant-design/icons'

class ErrorBoundary extends React.Component {
  constructor(props) {
    super(props)
    this.state = { 
      hasError: false, 
      error: null,
      errorInfo: null 
    }
  }

  static getDerivedStateFromError(error) {
    // Update state so the next render will show the fallback UI
    return { hasError: true }
  }

  componentDidCatch(error, errorInfo) {
    // Log error details
    console.error('ErrorBoundary caught an error:', error, errorInfo)
    
    this.setState({
      error,
      errorInfo
    })

    // You can also log the error to an error reporting service here
    // logErrorToService(error, errorInfo)
  }

  handleReload = () => {
    // Reset error state and reload
    this.setState({ 
      hasError: false, 
      error: null, 
      errorInfo: null 
    })
    
    // If onRetry is provided, call it; otherwise reload the page
    if (this.props.onRetry) {
      this.props.onRetry()
    } else {
      window.location.reload()
    }
  }

  render() {
    if (this.state.hasError) {
      // Fallback UI
      return (
        <div style={{ 
          padding: '50px 20px',
          textAlign: 'center',
          minHeight: '400px',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center'
        }}>
          <Result
            status="error"
            icon={<ExclamationCircleOutlined style={{ color: '#ff4d4f' }} />}
            title="页面出现错误"
            subTitle={
              this.props.showErrorDetails 
                ? `错误信息: ${this.state.error?.message || '未知错误'}`
                : "抱歉，页面遇到了一些问题。请尝试刷新页面或联系技术支持。"
            }
            extra={[
              <Button 
                key="reload" 
                type="primary" 
                icon={<ReloadOutlined />}
                onClick={this.handleReload}
              >
                重新加载
              </Button>,
              ...(this.props.showErrorDetails ? [
                <Button 
                  key="details" 
                  onClick={() => {
                    console.log('Error details:', {
                      error: this.state.error,
                      errorInfo: this.state.errorInfo
                    })
                    alert(`详细错误信息已输出到控制台。\n\n${this.state.error?.stack || '无堆栈信息'}`)
                  }}
                >
                  查看详情
                </Button>
              ] : [])
            ]}
          />
        </div>
      )
    }

    return this.props.children
  }
}

export default ErrorBoundary