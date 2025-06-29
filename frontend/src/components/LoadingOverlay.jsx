import React from 'react'
import { Spin } from 'antd'

const LoadingOverlay = ({
  loading = false,
  children,
  tip = '加载中...',
  size = 'large',
  delay = 200,
  overlay = true,
  minHeight,
  style = {}
}) => {
  if (!loading) {
    return children
  }

  // Full overlay loading
  if (overlay) {
    return (
      <div style={{ position: 'relative', ...style }}>
        {children}
        <div style={{
          position: 'absolute',
          top: 0,
          left: 0,
          right: 0,
          bottom: 0,
          background: 'rgba(255, 255, 255, 0.8)',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          zIndex: 1000,
          backdropFilter: 'blur(2px)',
          borderRadius: '8px'
        }}>
          <Spin size={size} tip={tip} delay={delay} />
        </div>
      </div>
    )
  }

  // Centered loading without overlay
  return (
    <div style={{
      display: 'flex',
      alignItems: 'center',
      justifyContent: 'center',
      minHeight: minHeight || '200px',
      flexDirection: 'column',
      ...style
    }}>
      <Spin size={size} tip={tip} delay={delay} />
    </div>
  )
}

export default LoadingOverlay