import React from 'react'
import { Card, Spin } from 'antd'

const LoadingCard = ({ 
  loading = true, 
  children, 
  tip = '加载中...',
  minHeight = 200,
  ...cardProps 
}) => {
  if (!loading) {
    return (
      <Card {...cardProps}>
        {children}
      </Card>
    )
  }

  return (
    <Card {...cardProps}>
      <div style={{
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        minHeight,
        flexDirection: 'column'
      }}>
        <Spin size="large" tip={tip} />
      </div>
    </Card>
  )
}

export default LoadingCard