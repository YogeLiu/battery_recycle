import React from 'react'
import { Card, Statistic } from 'antd'

const StatCard = ({
  title,
  value,
  suffix,
  prefix,
  precision,
  icon,
  color = '#1890ff',
  loading = false,
  onClick,
  bordered = true,
  ...cardProps
}) => {
  const cardStyle = onClick ? { cursor: 'pointer' } : {}

  return (
    <Card 
      bordered={bordered}
      loading={loading}
      onClick={onClick}
      style={cardStyle}
      {...cardProps}
    >
      <div style={{ display: 'flex', alignItems: 'center' }}>
        {icon && (
          <div style={{
            fontSize: 32,
            color,
            marginRight: 16,
            display: 'flex',
            alignItems: 'center'
          }}>
            {icon}
          </div>
        )}
        
        <div style={{ flex: 1 }}>
          <Statistic
            title={title}
            value={value}
            suffix={suffix}
            prefix={prefix}
            precision={precision}
            valueStyle={{ color, fontSize: icon ? 24 : 32 }}
          />
        </div>
      </div>
    </Card>
  )
}

export default StatCard