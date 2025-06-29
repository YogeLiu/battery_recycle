import React from 'react'
import { Tag } from 'antd'

const StatusTag = ({ status, statusMap }) => {
  // Default status configurations
  const defaultStatusMap = {
    // User status
    true: { color: 'green', text: '启用' },
    false: { color: 'red', text: '禁用' },
    
    // Order status
    completed: { color: 'green', text: '已完成' },
    cancelled: { color: 'red', text: '已取消' },
    pending: { color: 'orange', text: '待处理' },
    
    // User roles
    super_admin: { color: 'purple', text: '超级管理员' },
    normal: { color: 'blue', text: '普通用户' },
    
    // General status
    active: { color: 'green', text: '活跃' },
    inactive: { color: 'red', text: '非活跃' }
  }

  const finalStatusMap = { ...defaultStatusMap, ...statusMap }
  const config = finalStatusMap[status]

  if (!config) {
    return <Tag>{status}</Tag>
  }

  return (
    <Tag color={config.color}>
      {config.text}
    </Tag>
  )
}

export default StatusTag