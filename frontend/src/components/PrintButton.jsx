import React from 'react'
import { Button } from 'antd'
import { PrinterOutlined } from '@ant-design/icons'

const PrintButton = ({
  data,
  elementId,
  title = '打印',
  icon = <PrinterOutlined />,
  size = 'default',
  type = 'default',
  orderType = 'order', // 'inbound', 'outbound', or 'order'
  ...buttonProps
}) => {
  // Generate HTML content for order data
  const generateOrderHTML = (orderData, type) => {
    if (!orderData) return ''

    const isInbound = type === 'inbound' || orderData.supplier_name
    const orderTypeText = isInbound ? '入库单' : '出库单'
    const partnerLabel = isInbound ? '供应商' : '客户名称'
    const partnerName = orderData.supplier_name || orderData.customer_name
    const addressLabel = isInbound ? '车牌号' : '配送地址'
    const addressValue = orderData.vehicle_number || orderData.delivery_address

    return `
      <div style="max-width: 800px; margin: 0 auto; padding: 20px; font-family: Arial, sans-serif;">
        <div style="text-align: center; border-bottom: 2px solid #333; padding-bottom: 20px; margin-bottom: 30px;">
          <h1 style="margin: 0; font-size: 28px; color: #333;">废旧电池回收ERP系统</h1>
          <h2 style="margin: 10px 0 0 0; font-size: 24px; color: #666;">${orderTypeText}</h2>
        </div>
        
        <div style="margin-bottom: 30px;">
          <table style="width: 100%; border-collapse: collapse;">
            <tr>
              <td style="padding: 8px 0; font-weight: bold; width: 20%;">${orderTypeText}号:</td>
              <td style="padding: 8px 0; width: 30%;">${orderData.order_number || '-'}</td>
              <td style="padding: 8px 0; font-weight: bold; width: 20%;">${partnerLabel}:</td>
              <td style="padding: 8px 0; width: 30%;">${partnerName || '-'}</td>
            </tr>
            <tr>
              <td style="padding: 8px 0; font-weight: bold;">联系人:</td>
              <td style="padding: 8px 0;">${orderData.contact_person || '-'}</td>
              <td style="padding: 8px 0; font-weight: bold;">联系电话:</td>
              <td style="padding: 8px 0;">${orderData.contact_phone || '-'}</td>
            </tr>
            <tr>
              <td style="padding: 8px 0; font-weight: bold;">${addressLabel}:</td>
              <td style="padding: 8px 0;">${addressValue || '-'}</td>
              <td style="padding: 8px 0; font-weight: bold;">操作员:</td>
              <td style="padding: 8px 0;">${orderData.operator_name || '-'}</td>
            </tr>
            <tr>
              <td style="padding: 8px 0; font-weight: bold;">创建时间:</td>
              <td style="padding: 8px 0;" colspan="3">${orderData.created_at ? new Date(orderData.created_at).toLocaleString('zh-CN') : '-'}</td>
            </tr>
          </table>
        </div>

        <div style="margin-bottom: 30px;">
          <h3 style="margin-bottom: 15px; font-size: 18px; border-bottom: 1px solid #ddd; padding-bottom: 10px;">项目明细</h3>
          <table style="width: 100%; border-collapse: collapse; border: 1px solid #ddd;">
            <thead>
              <tr style="background-color: #f5f5f5;">
                <th style="border: 1px solid #ddd; padding: 12px; text-align: left;">品类名称</th>
                ${isInbound ? `
                  <th style="border: 1px solid #ddd; padding: 12px; text-align: center;">毛重 (kg)</th>
                  <th style="border: 1px solid #ddd; padding: 12px; text-align: center;">皮重 (kg)</th>
                  <th style="border: 1px solid #ddd; padding: 12px; text-align: center;">净重 (kg)</th>
                ` : `
                  <th style="border: 1px solid #ddd; padding: 12px; text-align: center;">重量 (kg)</th>
                `}
                <th style="border: 1px solid #ddd; padding: 12px; text-align: center;">单价 (元/kg)</th>
                <th style="border: 1px solid #ddd; padding: 12px; text-align: center;">小计 (元)</th>
              </tr>
            </thead>
            <tbody>
              ${(orderData.items || []).map(item => `
                <tr>
                  <td style="border: 1px solid #ddd; padding: 12px;">${item.category_name || '-'}</td>
                  ${isInbound ? `
                    <td style="border: 1px solid #ddd; padding: 12px; text-align: center;">${(item.gross_weight_kg || 0).toFixed(2)}</td>
                    <td style="border: 1px solid #ddd; padding: 12px; text-align: center;">${(item.tare_weight_kg || 0).toFixed(2)}</td>
                    <td style="border: 1px solid #ddd; padding: 12px; text-align: center; font-weight: bold;">${(item.net_weight_kg || 0).toFixed(2)}</td>
                  ` : `
                    <td style="border: 1px solid #ddd; padding: 12px; text-align: center; font-weight: bold;">${(item.weight_kg || 0).toFixed(2)}</td>
                  `}
                  <td style="border: 1px solid #ddd; padding: 12px; text-align: center;">¥${(item.unit_price || 0).toFixed(2)}</td>
                  <td style="border: 1px solid #ddd; padding: 12px; text-align: center; font-weight: bold;">¥${(item.subtotal_amount || 0).toFixed(2)}</td>
                </tr>
              `).join('')}
            </tbody>
          </table>
        </div>

        <div style="margin-bottom: 30px;">
          <table style="width: 100%; border-collapse: collapse;">
            <tr>
              <td style="padding: 8px 0; text-align: right; font-size: 16px;">
                <strong>总重量: ${(orderData.total_weight_kg || 0).toFixed(2)} kg</strong>
              </td>
            </tr>
            <tr>
              <td style="padding: 8px 0; text-align: right; font-size: 18px;">
                <strong style="color: #f5222d;">总金额: ¥${(orderData.total_amount || 0).toFixed(2)}</strong>
              </td>
            </tr>
          </table>
        </div>

        ${orderData.remarks ? `
          <div style="margin-bottom: 30px;">
            <h3 style="margin-bottom: 10px; font-size: 16px;">备注</h3>
            <p style="margin: 0; padding: 10px; background-color: #f9f9f9; border-left: 4px solid #1890ff;">${orderData.remarks}</p>
          </div>
        ` : ''}

        <div style="margin-top: 50px; text-align: center; border-top: 1px solid #ddd; padding-top: 20px;">
          <p style="margin: 0; color: #666; font-size: 14px;">打印时间: ${new Date().toLocaleString('zh-CN')}</p>
          <p style="margin: 5px 0 0 0; color: #666; font-size: 12px;">废旧电池回收ERP系统 - 专业的库存管理解决方案</p>
        </div>
      </div>
    `
  }

  const handlePrint = () => {
    let printContent = ''
    let printTitle = title

    if (data) {
      // Print structured data (order details)
      printContent = generateOrderHTML(data, orderType)
      printTitle = `${data.order_number || ''} - ${title}`
    } else if (elementId) {
      // Print specific element
      const element = document.getElementById(elementId)
      if (element) {
        printContent = element.innerHTML
      } else {
        console.error(`Element with ID '${elementId}' not found`)
        return
      }
    } else {
      // Print current page
      window.print()
      return
    }

    // Open new window for printing
    const printWindow = window.open('', '_blank')
    printWindow.document.write(`
      <!DOCTYPE html>
      <html>
        <head>
          <title>${printTitle}</title>
          <meta charset="utf-8">
          <style>
            * {
              box-sizing: border-box;
            }
            body { 
              font-family: 'Microsoft YaHei', Arial, sans-serif; 
              margin: 0; 
              padding: 20px;
              font-size: 14px;
              line-height: 1.4;
            }
            h1, h2, h3 {
              margin: 0;
              padding: 0;
            }
            table {
              border-collapse: collapse;
              width: 100%;
            }
            @media print {
              body { 
                margin: 0; 
                padding: 15px;
              }
              .no-print { 
                display: none !important; 
              }
              @page {
                margin: 1cm;
                size: A4;
              }
            }
          </style>
        </head>
        <body>
          ${printContent}
        </body>
      </html>
    `)
    
    printWindow.document.close()
    
    // Wait for content to load then print
    printWindow.onload = () => {
      setTimeout(() => {
        printWindow.print()
        printWindow.close()
      }, 250)
    }
  }

  return (
    <Button
      icon={icon}
      size={size}
      type={type}
      onClick={handlePrint}
      {...buttonProps}
    >
      {title}
    </Button>
  )
}

export default PrintButton