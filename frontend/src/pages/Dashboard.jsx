import React, { useState, useEffect } from 'react'
import { 
  Card, 
  Row, 
  Col, 
  Statistic, 
  Button, 
  Space, 
  Table,
  Tag,
  message,
  Spin,
  Alert
} from 'antd'
import { 
  DashboardOutlined, 
  InboxOutlined, 
  SendOutlined, 
  DatabaseOutlined,
  WarningOutlined,
  ReloadOutlined,
  TrendingUpOutlined,
  DollarOutlined,
  ShoppingCartOutlined,
  StockOutlined
} from '@ant-design/icons'
import { useNavigate } from 'react-router-dom'
import { inboundAPI } from '../api/inbound'
import { outboundAPI } from '../api/outbound'
import { inventoryAPI } from '../api/inventory'
import { categoriesAPI } from '../api/categories'
import { reportsAPI } from '../api/reports'

const Dashboard = () => {
  const [loading, setLoading] = useState(true)
  const [stats, setStats] = useState({
    todayInbound: 0,
    todayOutbound: 0,
    totalCategories: 0,
    totalStockWeight: 0,
    totalStockValue: 0,
    lowStockCount: 0,
    outOfStockCount: 0,
    monthlyRevenue: 0
  })
  const [recentInbound, setRecentInbound] = useState([])
  const [recentOutbound, setRecentOutbound] = useState([])
  const [lowStockItems, setLowStockItems] = useState([])
  const navigate = useNavigate()

  // Low stock threshold
  const LOW_STOCK_THRESHOLD = 10

  // Load dashboard data
  const loadDashboardData = async () => {
    setLoading(true)
    try {
      // Get today's date range
      const today = new Date()
      const startOfDay = new Date(today.setHours(0, 0, 0, 0))
      const endOfDay = new Date(today.setHours(23, 59, 59, 999))
      
      // Get current month date range
      const startOfMonth = new Date(today.getFullYear(), today.getMonth(), 1)
      const endOfMonth = new Date(today.getFullYear(), today.getMonth() + 1, 0)

      // Load all data in parallel
      const [
        inboundData,
        outboundData,
        inventoryData,
        categoriesData,
        todayInboundData,
        todayOutboundData,
        monthlyOutboundData
      ] = await Promise.all([
        inboundAPI.getAll({ limit: 5 }), // Recent 5 orders
        outboundAPI.getAll({ limit: 5 }), // Recent 5 orders
        inventoryAPI.getAll(),
        categoriesAPI.getAll(),
        inboundAPI.getAll({ 
          start_date: startOfDay.toISOString().split('T')[0],
          end_date: endOfDay.toISOString().split('T')[0]
        }),
        outboundAPI.getAll({ 
          start_date: startOfDay.toISOString().split('T')[0],
          end_date: endOfDay.toISOString().split('T')[0]
        }),
        outboundAPI.getAll({ 
          start_date: startOfMonth.toISOString().split('T')[0],
          end_date: endOfMonth.toISOString().split('T')[0]
        })
      ])

      // Calculate statistics
      const totalStockWeight = (inventoryData || []).reduce(
        (sum, item) => sum + (item.current_weight_kg || 0), 0
      )

      const totalStockValue = (inventoryData || []).reduce((sum, item) => {
        const category = categoriesData?.find(cat => cat.id === item.category_id)
        return sum + ((item.current_weight_kg || 0) * (category?.unit_price || 0))
      }, 0)

      const lowStockCount = (inventoryData || []).filter(item => {
        const weight = item.current_weight_kg || 0
        return weight > 0 && weight <= LOW_STOCK_THRESHOLD
      }).length

      const outOfStockCount = (inventoryData || []).filter(item => 
        (item.current_weight_kg || 0) === 0
      ).length

      const monthlyRevenue = (monthlyOutboundData || []).reduce(
        (sum, order) => sum + (order.total_amount || 0), 0
      )

      // Prepare low stock items for display
      const lowStockItemsData = (inventoryData || [])
        .filter(item => {
          const weight = item.current_weight_kg || 0
          return weight <= LOW_STOCK_THRESHOLD
        })
        .map(item => {
          const category = categoriesData?.find(cat => cat.id === item.category_id)
          return {
            ...item,
            category_name: category?.name || '未知品类',
            unit_price: category?.unit_price || 0,
            stock_status: (item.current_weight_kg || 0) === 0 ? 'out_of_stock' : 'low_stock'
          }
        })
        .slice(0, 10) // Show top 10 items

      setStats({
        todayInbound: (todayInboundData || []).length,
        todayOutbound: (todayOutboundData || []).length,
        totalCategories: (categoriesData || []).length,
        totalStockWeight,
        totalStockValue,
        lowStockCount,
        outOfStockCount,
        monthlyRevenue
      })

      setRecentInbound((inboundData || []).slice(0, 5))
      setRecentOutbound((outboundData || []).slice(0, 5))
      setLowStockItems(lowStockItemsData)

    } catch (error) {
      message.error('加载仪表板数据失败：' + error.message)
    } finally {
      setLoading(false)
    }
  }

  // Load data on mount
  useEffect(() => {
    loadDashboardData()
  }, [])

  // Recent orders columns
  const recentOrderColumns = [
    {
      title: '单号',
      dataIndex: 'order_number',
      key: 'order_number',
      width: 140,
      render: (text) => <span style={{ fontFamily: 'monospace' }}>{text}</span>
    },
    {
      title: '客户/供应商',
      key: 'partner',
      width: 120,
      render: (_, record) => record.customer_name || record.supplier_name || '-'
    },
    {
      title: '金额',
      dataIndex: 'total_amount',
      key: 'total_amount',
      width: 100,
      render: (amount) => `¥${(amount || 0).toFixed(2)}`
    },
    {
      title: '时间',
      dataIndex: 'created_at',
      key: 'created_at',
      width: 100,
      render: (date) => new Date(date).toLocaleDateString('zh-CN')
    }
  ]

  // Low stock columns
  const lowStockColumns = [
    {
      title: '品类名称',
      dataIndex: 'category_name',
      key: 'category_name'
    },
    {
      title: '当前库存',
      dataIndex: 'current_weight_kg',
      key: 'current_weight_kg',
      render: (weight) => `${(weight || 0).toFixed(2)} kg`
    },
    {
      title: '状态',
      dataIndex: 'stock_status',
      key: 'stock_status',
      render: (status) => (
        <Tag color={status === 'out_of_stock' ? 'red' : 'orange'}>
          {status === 'out_of_stock' ? '缺货' : '库存不足'}
        </Tag>
      )
    }
  ]

  if (loading) {
    return (
      <div style={{ 
        display: 'flex', 
        justifyContent: 'center', 
        alignItems: 'center', 
        height: '400px' 
      }}>
        <Spin size="large" tip="加载仪表板数据..." />
      </div>
    )
  }

  return (
    <div>
      <div className="page-header">
        <h1 className="page-title">
          <DashboardOutlined /> 仪表板
        </h1>
        <Space>
          <Button 
            icon={<ReloadOutlined />} 
            onClick={loadDashboardData}
            loading={loading}
          >
            刷新数据
          </Button>
        </Space>
      </div>

      {/* Alerts for critical issues */}
      {(stats.outOfStockCount > 0 || stats.lowStockCount > 0) && (
        <Alert
          type="warning"
          showIcon
          style={{ marginBottom: 16 }}
          message={`库存警告: ${stats.outOfStockCount} 个品类缺货，${stats.lowStockCount} 个品类库存不足`}
          description="请及时补充库存以避免业务中断"
          action={
            <Button size="small" onClick={() => navigate('/inventory')}>
              查看库存
            </Button>
          }
        />
      )}
      
      {/* Key Performance Indicators */}
      <Row gutter={[16, 16]} style={{ marginBottom: 24 }}>
        <Col xs={24} sm={12} lg={6}>
          <Card hoverable onClick={() => navigate('/inbound')}>
            <Statistic
              title="今日入库"
              value={stats.todayInbound}
              suffix="单"
              valueStyle={{ color: '#52c41a' }}
              prefix={<InboxOutlined />}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card hoverable onClick={() => navigate('/outbound')}>
            <Statistic
              title="今日出库"
              value={stats.todayOutbound}
              suffix="单"
              valueStyle={{ color: '#1890ff' }}
              prefix={<SendOutlined />}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card hoverable onClick={() => navigate('/inventory')}>
            <Statistic
              title="库存类目"
              value={stats.totalCategories}
              suffix="种"
              valueStyle={{ color: '#722ed1' }}
              prefix={<DatabaseOutlined />}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card hoverable onClick={() => navigate('/inventory')}>
            <Statistic
              title="总库存重量"
              value={stats.totalStockWeight}
              suffix="kg"
              precision={2}
              valueStyle={{ color: '#fa8c16' }}
              prefix={<StockOutlined />}
            />
          </Card>
        </Col>
      </Row>

      {/* Secondary Metrics */}
      <Row gutter={[16, 16]} style={{ marginBottom: 24 }}>
        <Col xs={24} sm={12} lg={6}>
          <Card>
            <Statistic
              title="库存总价值"
              value={stats.totalStockValue}
              precision={2}
              valueStyle={{ color: '#f5222d' }}
              prefix="¥"
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card>
            <Statistic
              title="本月营收"
              value={stats.monthlyRevenue}
              precision={2}
              valueStyle={{ color: '#13c2c2' }}
              prefix="¥"
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card>
            <Statistic
              title="库存不足"
              value={stats.lowStockCount}
              suffix="种"
              valueStyle={{ color: '#faad14' }}
              prefix={<WarningOutlined />}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card>
            <Statistic
              title="缺货品类"
              value={stats.outOfStockCount}
              suffix="种"
              valueStyle={{ color: '#ff4d4f' }}
              prefix={<WarningOutlined />}
            />
          </Card>
        </Col>
      </Row>

      {/* Data Tables */}
      <Row gutter={[16, 16]}>
        <Col xs={24} lg={12}>
          <Card 
            title="最近入库单" 
            extra={
              <Button 
                type="link" 
                onClick={() => navigate('/inbound')}
              >
                查看全部
              </Button>
            }
          >
            <Table
              dataSource={recentInbound}
              columns={recentOrderColumns}
              pagination={false}
              size="small"
              rowKey="id"
              locale={{ emptyText: '暂无入库记录' }}
            />
          </Card>
        </Col>
        <Col xs={24} lg={12}>
          <Card 
            title="最近出库单" 
            extra={
              <Button 
                type="link" 
                onClick={() => navigate('/outbound')}
              >
                查看全部
              </Button>
            }
          >
            <Table
              dataSource={recentOutbound}
              columns={recentOrderColumns}
              pagination={false}
              size="small"
              rowKey="id"
              locale={{ emptyText: '暂无出库记录' }}
            />
          </Card>
        </Col>
      </Row>

      {/* Low Stock Alert */}
      {lowStockItems.length > 0 && (
        <Row style={{ marginTop: 16 }}>
          <Col span={24}>
            <Card 
              title={
                <span style={{ color: '#faad14' }}>
                  <WarningOutlined /> 库存预警
                </span>
              }
              extra={
                <Button 
                  type="link" 
                  onClick={() => navigate('/inventory')}
                >
                  查看库存
                </Button>
              }
            >
              <Table
                dataSource={lowStockItems}
                columns={lowStockColumns}
                pagination={false}
                size="small"
                rowKey="category_id"
              />
            </Card>
          </Col>
        </Row>
      )}

      {/* System Functions */}
      <Row style={{ marginTop: 16 }}>
        <Col span={24}>
          <Card title="系统功能" extra={
            <Button type="link" onClick={() => navigate('/reports')}>
              查看报表
            </Button>
          }>
            <Row gutter={[16, 16]}>
              <Col xs={24} sm={12} md={6}>
                <Card 
                  size="small" 
                  className="stats-card"
                  hoverable
                  onClick={() => navigate('/inbound')}
                >
                  <InboxOutlined style={{ fontSize: 24, color: '#52c41a' }} />
                  <div style={{ marginTop: 8 }}>采购入库</div>
                </Card>
              </Col>
              <Col xs={24} sm={12} md={6}>
                <Card 
                  size="small" 
                  className="stats-card"
                  hoverable
                  onClick={() => navigate('/outbound')}
                >
                  <SendOutlined style={{ fontSize: 24, color: '#1890ff' }} />
                  <div style={{ marginTop: 8 }}>销售出库</div>
                </Card>
              </Col>
              <Col xs={24} sm={12} md={6}>
                <Card 
                  size="small" 
                  className="stats-card"
                  hoverable
                  onClick={() => navigate('/inventory')}
                >
                  <DatabaseOutlined style={{ fontSize: 24, color: '#722ed1' }} />
                  <div style={{ marginTop: 8 }}>库存管理</div>
                </Card>
              </Col>
              <Col xs={24} sm={12} md={6}>
                <Card 
                  size="small" 
                  className="stats-card"
                  hoverable
                  onClick={() => navigate('/reports')}
                >
                  <DashboardOutlined style={{ fontSize: 24, color: '#fa8c16' }} />
                  <div style={{ marginTop: 8 }}>报表分析</div>
                </Card>
              </Col>
            </Row>
          </Card>
        </Col>
      </Row>
    </div>
  )
}

export default Dashboard