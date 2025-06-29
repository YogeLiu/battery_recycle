import React from 'react'
import { Empty, Button } from 'antd'
import { PlusOutlined } from '@ant-design/icons'

const EmptyState = ({
  description = '暂无数据',
  image,
  actionText,
  onAction,
  actionIcon = <PlusOutlined />,
  actionType = 'primary',
  children
}) => {
  return (
    <div className="empty-container">
      <Empty
        image={image}
        description={description}
      >
        {actionText && onAction && (
          <Button 
            type={actionType}
            icon={actionIcon}
            onClick={onAction}
          >
            {actionText}
          </Button>
        )}
        {children}
      </Empty>
    </div>
  )
}

export default EmptyState