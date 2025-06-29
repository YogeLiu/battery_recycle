import React from 'react'
import { Button, Dropdown, message } from 'antd'
import { DownloadOutlined, ExportOutlined, FileExcelOutlined, FileTextOutlined } from '@ant-design/icons'

const ExportButton = ({
  data,
  filename = 'export',
  formats = ['csv', 'excel'],
  title = '导出',
  icon = <ExportOutlined />,
  size = 'default',
  type = 'default',
  onExport,
  columns, // Column definitions for better headers
  ...buttonProps
}) => {
  // Chinese header mapping for common fields
  const fieldNameMap = {
    // Common fields
    'id': 'ID',
    'name': '名称',
    'username': '用户名',
    'real_name': '真实姓名',
    'role': '角色',
    'role_name': '角色',
    'is_active': '状态',
    'created_at': '创建时间',
    'updated_at': '更新时间',
    'last_login_at': '最后登录时间',
    
    // Categories
    'category_name': '品类名称',
    'description': '描述',
    'unit_price': '单价(元/kg)',
    
    // Orders
    'order_number': '单号',
    'supplier_name': '供应商',
    'customer_name': '客户名称',
    'contact_person': '联系人',
    'contact_phone': '联系电话',
    'vehicle_number': '车牌号',
    'delivery_address': '配送地址',
    'total_weight_kg': '总重量(kg)',
    'total_amount': '总金额(元)',
    'status': '状态',
    'operator_name': '操作员',
    'remarks': '备注',
    
    // Inventory
    'current_weight_kg': '当前库存(kg)',
    'stock_value': '库存价值(元)',
    'inventory_value': '库存价值(元)',
    
    // Items
    'gross_weight_kg': '毛重(kg)',
    'tare_weight_kg': '皮重(kg)',
    'net_weight_kg': '净重(kg)',
    'weight_kg': '重量(kg)',
    'subtotal_amount': '小计(元)',
    
    // Reports
    'inbound_weight': '入库重量(kg)',
    'inbound_amount': '入库金额(元)',
    'outbound_weight': '出库重量(kg)',
    'outbound_amount': '出库金额(元)',
    'current_stock': '当前库存(kg)'
  }

  const formatValue = (value, key) => {
    if (value === null || value === undefined) return ''
    
    // Format different data types
    if (typeof value === 'boolean') {
      return value ? '是' : '否'
    }
    
    if (typeof value === 'number') {
      // Format currency fields
      if (key.includes('amount') || key.includes('price') || key.includes('value')) {
        return `¥${value.toFixed(2)}`
      }
      // Format weight fields
      if (key.includes('weight')) {
        return `${value.toFixed(2)}kg`
      }
      return value.toString()
    }
    
    if (typeof value === 'string' && value.includes('T') && value.includes('Z')) {
      // Format dates
      try {
        return new Date(value).toLocaleString('zh-CN')
      } catch {
        return value
      }
    }
    
    return value.toString()
  }

  const getHeaders = (data) => {
    if (!data || data.length === 0) return []
    
    // Use columns if provided, otherwise extract from data
    if (columns && columns.length > 0) {
      return columns
        .filter(col => col.dataIndex && col.title)
        .map(col => ({
          key: col.dataIndex,
          label: col.title
        }))
    }
    
    // Auto-generate headers from data keys
    const keys = Object.keys(data[0])
    return keys.map(key => ({
      key,
      label: fieldNameMap[key] || key
    }))
  }

  const exportToCSV = (data, filename, forExcel = false) => {
    if (!data || data.length === 0) {
      message.warning('没有数据可导出')
      return
    }

    const headers = getHeaders(data)
    const csvContent = [
      // Header row
      headers.map(h => `"${h.label}"`).join(','),
      // Data rows
      ...data.map(row => 
        headers.map(h => {
          const value = formatValue(row[h.key], h.key)
          return `"${value.toString().replace(/"/g, '""')}"`
        }).join(',')
      )
    ].join('\n')

    const BOM = '\ufeff' // UTF-8 BOM for proper Chinese character display
    const mimeType = forExcel ? 'application/vnd.ms-excel' : 'text/csv'
    const extension = forExcel ? 'xls' : 'csv'
    
    const blob = new Blob([BOM + csvContent], { 
      type: `${mimeType};charset=utf-8;` 
    })
    
    downloadFile(blob, `${filename}.${extension}`)
    
    message.success(`数据已导出为 ${extension.toUpperCase()} 格式`)
  }

  const exportToExcel = (data, filename) => {
    exportToCSV(data, filename, true)
  }

  const exportToJSON = (data, filename) => {
    if (!data || data.length === 0) {
      message.warning('没有数据可导出')
      return
    }

    const jsonContent = JSON.stringify(data, null, 2)
    const blob = new Blob([jsonContent], { type: 'application/json' })
    
    downloadFile(blob, `${filename}.json`)
    message.success('数据已导出为 JSON 格式')
  }

  const downloadFile = (blob, filename) => {
    const link = document.createElement('a')
    const url = URL.createObjectURL(blob)
    link.setAttribute('href', url)
    link.setAttribute('download', filename)
    link.style.visibility = 'hidden'
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    URL.revokeObjectURL(url)
  }

  const handleExport = (format) => {
    if (onExport) {
      onExport(format)
      return
    }

    if (!data || data.length === 0) {
      message.warning('没有数据可导出')
      return
    }

    const timestamp = new Date().toISOString().slice(0, 16).replace(/[:-]/g, '')
    const filenameWithTime = `${filename}_${timestamp}`

    try {
      switch (format) {
        case 'csv':
          exportToCSV(data, filenameWithTime)
          break
        case 'excel':
        case 'xls':
          exportToExcel(data, filenameWithTime)
          break
        case 'json':
          exportToJSON(data, filenameWithTime)
          break
        default:
          message.error(`不支持的导出格式: ${format}`)
      }
    } catch (error) {
      message.error('导出失败: ' + error.message)
    }
  }

  // Format options with icons
  const formatOptions = {
    csv: { label: 'CSV 格式', icon: <FileTextOutlined /> },
    excel: { label: 'Excel 格式', icon: <FileExcelOutlined /> },
    xls: { label: 'Excel 格式', icon: <FileExcelOutlined /> },
    json: { label: 'JSON 格式', icon: <DownloadOutlined /> }
  }

  // If only one format, show simple button
  if (formats.length === 1) {
    const format = formats[0]
    const option = formatOptions[format] || { label: format.toUpperCase(), icon: <DownloadOutlined /> }
    
    return (
      <Button
        icon={icon}
        size={size}
        type={type}
        onClick={() => handleExport(format)}
        {...buttonProps}
      >
        {title}
      </Button>
    )
  }

  // Multiple formats, show dropdown
  const menuItems = formats.map(format => {
    const option = formatOptions[format] || { label: format.toUpperCase(), icon: <DownloadOutlined /> }
    return {
      key: format,
      label: `导出为 ${option.label}`,
      icon: option.icon,
      onClick: () => handleExport(format)
    }
  })

  return (
    <Dropdown 
      menu={{ items: menuItems }} 
      placement="bottomRight"
      trigger={['click']}
    >
      <Button
        icon={icon}
        size={size}
        type={type}
        {...buttonProps}
      >
        {title}
      </Button>
    </Dropdown>
  )
}

export default ExportButton