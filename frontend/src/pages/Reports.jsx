import React, { useState, useEffect } from 'react'
import { 
  Form, 
  DatePicker, 
  Button, 
  message, 
  Space, 
  Row,
  Col,
  Card,
  Statistic,
  Table,
  Divider,
  Select,
  Tabs
} from 'antd'
import { 
  BarChartOutlined,
  InboxOutlined,
  SendOutlined,
  DatabaseOutlined,
  DollarOutlined,
  TrendingUpOutlined,
  CalendarOutlined,
  FileTextOutlined
} from '@ant-design/icons'
import { 
  PageHeader, 
  ExportButton
} from '../components'
import { reportsAPI } from '../api/reports'
import { inboundAPI } from '../api/inbound'
import { outboundAPI } from '../api/outbound'
import { inventoryAPI } from '../api/inventory'
import { categoriesAPI } from '../api/categories'

const { RangePicker } = DatePicker
const { TabPane } = Tabs

const Reports = () => {
  const [loading, setLoading] = useState(false)
  const [dateRange, setDateRange] = useState(null)
  const [reportForm] = Form.useForm()
  
  // Report data states
  const [summaryData, setSummaryData] = useState({})
  const [inboundData, setInboundData] = useState([])
  const [outboundData, setOutboundData] = useState([])
  const [inventoryData, setInventoryData] = useState([])
  const [categoryStats, setCategoryStats] = useState([])

  // Load summary data
  const loadSummaryData = async (params = {}) => {
    setLoading(true)
    try {
      const result = await reportsAPI.getSummary(params)
      setSummaryData(result || {})
    } catch (error) {
      message.error('加载汇总数据失败：' + error.message)
    } finally {
      setLoading(false)
    }
  }

  // Load detailed data
  const loadDetailedData = async (params = {}) => {
    try {
      // Load parallel data
      const [inbound, outbound, inventory, categories] = await Promise.all([
        inboundAPI.getAll(params),
        outboundAPI.getAll(params),
        inventoryAPI.getAll(),
        categoriesAPI.getAll()
      ])

      setInboundData(inbound || [])
      setOutboundData(outbound || [])
      setInventoryData(inventory || [])
      
      // Calculate category statistics
      calculateCategoryStats(inbound || [], outbound || [], inventory || [], categories || [])
    } catch (error) {
      message.error('加载详细数据失败：' + error.message)
    }
  }

  // Calculate category statistics
  const calculateCategoryStats = (inbound, outbound, inventory, categories) => {
    const categoryMap = new Map()
    
    // Initialize with categories
    categories.forEach(cat => {
      categoryMap.set(cat.id, {
        id: cat.id,
        name: cat.name,
        unit_price: cat.unit_price,
        inbound_weight: 0,
        inbound_amount: 0,
        outbound_weight: 0,
        outbound_amount: 0,
        current_stock: 0,
        stock_value: 0
      })
    })
    
    // Calculate inbound stats
    inbound.forEach(order => {
      if (order.items) {
        order.items.forEach(item => {
          const stat = categoryMap.get(item.category_id)
          if (stat) {
            stat.inbound_weight += item.net_weight_kg || 0
            stat.inbound_amount += item.subtotal_amount || 0
          }
        })
      }
    })
    
    // Calculate outbound stats
    outbound.forEach(order => {
      if (order.items) {
        order.items.forEach(item => {
          const stat = categoryMap.get(item.category_id)
          if (stat) {
            stat.outbound_weight += item.weight_kg || 0
            stat.outbound_amount += item.subtotal_amount || 0
          }
        })
      }
    })
    
    // Add current inventory
    inventory.forEach(inv => {
      const stat = categoryMap.get(inv.category_id)
      if (stat) {
        stat.current_stock = inv.current_weight_kg || 0
        stat.stock_value = stat.current_stock * stat.unit_price
      }
    })
    
    setCategoryStats(Array.from(categoryMap.values()))
  }

  // Handle date range change
  const handleDateRangeChange = (dates) => {
    setDateRange(dates)
    const params = {}
    if (dates) {
      params.start_date = dates[0].format('YYYY-MM-DD')
      params.end_date = dates[1].format('YYYY-MM-DD')
    }
    loadSummaryData(params)
    loadDetailedData(params)
  }

  // Load data on mount
  useEffect(() => {
    loadSummaryData()
    loadDetailedData()
  }, [])

  // Category performance columns
  const categoryColumns = [
    {
      title: '品类名称',
      dataIndex: 'name',
      key: 'name',
      render: (text) => <strong>{text}</strong>
    },
    {
      title: '入库重量 (kg)',
      dataIndex: 'inbound_weight',
      key: 'inbound_weight',
      render: (val) => val?.toFixed(2) || '0.00',
      sorter: (a, b) => (a.inbound_weight || 0) - (b.inbound_weight || 0)
    },
    {
      title: '入库金额 (元)',
      dataIndex: 'inbound_amount',
      key: 'inbound_amount',
      render: (val) => `¥${val?.toFixed(2) || '0.00'}`,
      sorter: (a, b) => (a.inbound_amount || 0) - (b.inbound_amount || 0)
    },
    {
      title: '出库重量 (kg)',
      dataIndex: 'outbound_weight',
      key: 'outbound_weight',
      render: (val) => val?.toFixed(2) || '0.00',
      sorter: (a, b) => (a.outbound_weight || 0) - (b.outbound_weight || 0)
    },
    {
      title: '出库金额 (元)',
      dataIndex: 'outbound_amount',
      key: 'outbound_amount',
      render: (val) => `¥${val?.toFixed(2) || '0.00'}`,
      sorter: (a, b) => (a.outbound_amount || 0) - (b.outbound_amount || 0)
    },
    {
      title: '当前库存 (kg)',
      dataIndex: 'current_stock',
      key: 'current_stock',
      render: (val) => val?.toFixed(2) || '0.00',
      sorter: (a, b) => (a.current_stock || 0) - (b.current_stock || 0)
    },
    {
      title: '库存价值 (元)',
      dataIndex: 'stock_value',
      key: 'stock_value',
      render: (val) => `¥${val?.toFixed(2) || '0.00'}`,
      sorter: (a, b) => (a.stock_value || 0) - (b.stock_value || 0)
    }
  ]

  // Inbound summary columns
  const inboundColumns = [
    {
      title: '入库单号',
      dataIndex: 'order_number',
      key: 'order_number'
    },
    {
      title: '供应商',
      dataIndex: 'supplier_name',
      key: 'supplier_name'
    },
    {
      title: '总重量 (kg)',
      dataIndex: 'total_weight_kg',
      key: 'total_weight_kg',
      render: (val) => val?.toFixed(2) || '0.00'
    },
    {
      title: '总金额 (元)',
      dataIndex: 'total_amount',
      key: 'total_amount',
      render: (val) => `¥${val?.toFixed(2) || '0.00'}`
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      key: 'created_at',
      render: (date) => new Date(date).toLocaleDateString('zh-CN')
    }
  ]

  // Outbound summary columns
  const outboundColumns = [
    {
      title: '出库单号',
      dataIndex: 'order_number',
      key: 'order_number'
    },
    {
      title: '客户名称',
      dataIndex: 'customer_name',
      key: 'customer_name'
    },
    {
      title: '总重量 (kg)',
      dataIndex: 'total_weight_kg',
      key: 'total_weight_kg',
      render: (val) => val?.toFixed(2) || '0.00'
    },
    {
      title: '总金额 (元)',
      dataIndex: 'total_amount',
      key: 'total_amount',
      render: (val) => `¥${val?.toFixed(2) || '0.00'}`
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      key: 'created_at',
      render: (date) => new Date(date).toLocaleDateString('zh-CN')
    }
  ]

  return (
    <div>
      <PageHeader
        title="报表分析"
        subTitle="业务数据统计与分析报告"
        icon={<BarChartOutlined />}
        extra={
          <Space>
            <RangePicker 
              onChange={handleDateRangeChange}
              placeholder={['开始日期', '结束日期']}
              style={{ width: 240 }}
            />
            <Button 
              type="primary"
              onClick={() => {
                loadSummaryData()
                loadDetailedData()
              }}
            >
              刷新数据
            </Button>
            <ExportButton 
              data={categoryStats}
              filename="business_report"
              title="导出报表"
            />
          </Space>
        }
      >
        {/* Summary Statistics */}
        <Row gutter={16} style={{ marginBottom: 24 }}>
          <Col xs={24} sm={12} lg={6}>
            <Card>
              <Statistic
                title="总入库重量"
                value={summaryData.total_inbound_weight || inboundData.reduce((sum, item) => sum + (item.total_weight_kg || 0), 0)}
                suffix="kg"
                precision={2}
                valueStyle={{ color: '#52c41a' }}
                prefix={<InboxOutlined />}
              />
            </Card>
          </Col>
          <Col xs={24} sm={12} lg={6}>
            <Card>
              <Statistic
                title="总出库重量"
                value={summaryData.total_outbound_weight || outboundData.reduce((sum, item) => sum + (item.total_weight_kg || 0), 0)}
                suffix="kg"
                precision={2}
                valueStyle={{ color: '#1890ff' }}
                prefix={<SendOutlined />}
              />
            </Card>
          </Col>
          <Col xs={24} sm={12} lg={6}>
            <Card>
              <Statistic
                title="总入库金额"
                value={summaryData.total_inbound_amount || inboundData.reduce((sum, item) => sum + (item.total_amount || 0), 0)}
                precision={2}
                valueStyle={{ color: '#722ed1' }}
                prefix="¥"
              />
            </Card>
          </Col>
          <Col xs={24} sm={12} lg={6}>
            <Card>
              <Statistic
                title="总出库金额"
                value={summaryData.total_outbound_amount || outboundData.reduce((sum, item) => sum + (item.total_amount || 0), 0)}
                precision={2}
                valueStyle={{ color: '#fa8c16' }}
                prefix="¥"
              />
            </Card>
          </Col>
        </Row>

        {/* Additional Metrics */}
        <Row gutter={16} style={{ marginBottom: 24 }}>
          <Col xs={24} sm={12} lg={6}>
            <Card>
              <Statistic
                title="入库单数量"
                value={inboundData.length}
                suffix="单"
                valueStyle={{ color: '#13c2c2' }}
              />
            </Card>
          </Col>
          <Col xs={24} sm={12} lg={6}>
            <Card>
              <Statistic
                title="出库单数量"
                value={outboundData.length}
                suffix="单"
                valueStyle={{ color: '#eb2f96' }}
              />
            </Card>
          </Col>
          <Col xs={24} sm={12} lg={6}>
            <Card>
              <Statistic
                title="当前库存价值"
                value={categoryStats.reduce((sum, cat) => sum + (cat.stock_value || 0), 0)}
                precision={2}
                valueStyle={{ color: '#f5222d' }}
                prefix="¥"
              />
            </Card>
          </Col>
          <Col xs={24} sm={12} lg={6}>
            <Card>
              <Statistic
                title="净收益"
                value={(summaryData.total_outbound_amount || outboundData.reduce((sum, item) => sum + (item.total_amount || 0), 0)) - 
                      (summaryData.total_inbound_amount || inboundData.reduce((sum, item) => sum + (item.total_amount || 0), 0))}
                precision={2}
                valueStyle={{ 
                  color: ((summaryData.total_outbound_amount || 0) - (summaryData.total_inbound_amount || 0)) >= 0 
                    ? '#52c41a' : '#ff4d4f' 
                }}
                prefix="¥"
              />
            </Card>
          </Col>
        </Row>
      </PageHeader>

      {/* Detailed Reports Tabs */}
      <Tabs defaultActiveKey="category" size="large">
        <TabPane 
          tab={
            <span>
              <BarChartOutlined />
              品类分析
            </span>
          } 
          key="category"
        >
          <Card title="品类业务分析" extra={<FileTextOutlined />}>
            <Table
              dataSource={categoryStats}
              columns={categoryColumns}
              loading={loading}
              rowKey="id"
              pagination={{
                showSizeChanger: true,
                showTotal: (total) => `共 ${total} 个品类`
              }}
              scroll={{ x: 800 }}
            />
          </Card>
        </TabPane>

        <TabPane 
          tab={
            <span>
              <InboxOutlined />
              入库明细
            </span>
          } 
          key="inbound"
        >
          <Card title="入库业务明细" extra={<FileTextOutlined />}>
            <Table
              dataSource={inboundData}
              columns={inboundColumns}
              loading={loading}
              rowKey="id"
              pagination={{
                showSizeChanger: true,
                showTotal: (total) => `共 ${total} 条记录`
              }}
            />
          </Card>
        </TabPane>

        <TabPane 
          tab={
            <span>
              <SendOutlined />
              出库明细
            </span>
          } 
          key="outbound"
        >
          <Card title="出库业务明细" extra={<FileTextOutlined />}>
            <Table
              dataSource={outboundData}
              columns={outboundColumns}
              loading={loading}
              rowKey="id"
              pagination={{
                showSizeChanger: true,
                showTotal: (total) => `共 ${total} 条记录`
              }}
            />
          </Card>
        </TabPane>

        <TabPane 
          tab={
            <span>
              <DatabaseOutlined />
              库存概览
            </span>
          } 
          key="inventory"
        >
          <Card title="当前库存状况" extra={<FileTextOutlined />}>
            <Table
              dataSource={inventoryData.map(inv => {
                const cat = categoryStats.find(c => c.id === inv.category_id)
                return {
                  ...inv,
                  category_name: cat?.name || '未知品类',
                  unit_price: cat?.unit_price || 0,
                  stock_value: (inv.current_weight_kg || 0) * (cat?.unit_price || 0)
                }
              })}
              columns={[
                { title: '品类名称', dataIndex: 'category_name', key: 'category_name' },
                { title: '当前库存 (kg)', dataIndex: 'current_weight_kg', key: 'current_weight_kg', render: val => val?.toFixed(2) },
                { title: '单价 (元/kg)', dataIndex: 'unit_price', key: 'unit_price', render: val => `¥${val?.toFixed(2)}` },
                { title: '库存价值 (元)', dataIndex: 'stock_value', key: 'stock_value', render: val => `¥${val?.toFixed(2)}` },
                { title: '最后更新', dataIndex: 'updated_at', key: 'updated_at', render: date => date ? new Date(date).toLocaleString('zh-CN') : '-' }
              ]}
              loading={loading}
              rowKey="category_id"
              pagination={{
                showSizeChanger: true,
                showTotal: (total) => `共 ${total} 条记录`
              }}
            />
          </Card>
        </TabPane>
      </Tabs>
    </div>
  )
}

export default Reports