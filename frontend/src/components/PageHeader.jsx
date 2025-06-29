import React from 'react'
import { Button, Space } from 'antd'

const PageHeader = ({ 
  title, 
  subTitle, 
  icon, 
  extra,
  actions,
  children 
}) => {
  return (
    <div className="page-header">
      <div style={{ 
        display: 'flex', 
        justifyContent: 'space-between', 
        alignItems: 'center',
        marginBottom: subTitle || children ? 16 : 0
      }}>
        <div>
          <h1 className="page-title">
            {icon && <span style={{ marginRight: 8 }}>{icon}</span>}
            {title}
          </h1>
          {subTitle && (
            <p style={{ margin: 0, color: '#8c8c8c', fontSize: 14 }}>
              {subTitle}
            </p>
          )}
        </div>
        
        {(extra || actions) && (
          <div>
            {actions && (
              <Space>
                {actions.map((action, index) => (
                  <Button key={index} {...action.props}>
                    {action.label}
                  </Button>
                ))}
              </Space>
            )}
            {extra}
          </div>
        )}
      </div>
      
      {children}
    </div>
  )
}

export default PageHeader