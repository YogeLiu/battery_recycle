import React, { useState, useEffect } from 'react'
import { 
  Form, 
  Input, 
  Select, 
  Button, 
  message, 
  Space, 
  Tag,
  Progress,
  Alert,
  Row,
  Col,
  Card,
  Statistic
} from 'antd'
import { 
  DatabaseOutlined, 
  WarningOutlined,
  CheckCircleOutlined,
  ExclamationCircleOutlined,
  SearchOutlined
} from '@ant-design/icons'
import { 
  PageHeader, 
  SearchForm, 
  DataTable, 
  EmptyState,
  ExportButton
} from '../components'
import { inventoryAPI } from '../api/inventory'
import { categoriesAPI } from '../api/categories'

const Inventory = () => {
  const [data, setData] = useState([])
  const [categories, setCategories] = useState([])
  const [loading, setLoading] = useState(false)
  const [searchForm] = Form.useForm()
  const [stats, setStats] = useState({
    totalCategories: 0,
    totalWeight: 0,
    lowStockCount: 0,
    outOfStockCount: 0
  })

  // Low stock threshold (can be configurable)
  const LOW_STOCK_THRESHOLD = 10

  // Table columns
  const columns = [
    {
      title: '品类名称',
      dataIndex: 'category_name',
      key: 'category_name',
      width: 200,
      render: (text) => <strong>{text}</strong>
    },
    {
      title: '当前库存 (kg)',
      dataIndex: 'current_weight_kg',
      key: 'current_weight_kg',
      width: 150,
      render: (weight) => {
        const value = weight || 0
        let color = '#52c41a' // green
        if (value === 0) color = '#ff4d4f' // red
        else if (value <= LOW_STOCK_THRESHOLD) color = '#faad14' // orange
        
        return (
          <span style={{ color, fontWeight: 'bold', fontSize: '16px' }}>
            {value.toFixed(2)}
          </span>
        )
      }
    },
    {
      title: '库存状态',
      dataIndex: 'current_weight_kg',
      key: 'status',
      width: 120,
      render: (weight) => {
        const value = weight || 0
        if (value === 0) {
          return <Tag color="red" icon={<ExclamationCircleOutlined />}>缺货</Tag>
        } else if (value <= LOW_STOCK_THRESHOLD) {
          return <Tag color="orange" icon={<WarningOutlined />}>库存不足</Tag>
        } else {
          return <Tag color="green" icon={<CheckCircleOutlined />}>库存充足</Tag>
        }
      }
    },
    {
      title: '库存健康度',
      dataIndex: 'current_weight_kg',
      key: 'health',
      width: 150,
      render: (weight) => {
        const value = weight || 0
        let percent = 0
        let status = 'exception'
        
        if (value === 0) {
          percent = 0
          status = 'exception'
        } else if (value <= LOW_STOCK_THRESHOLD) {
          percent = Math.min((value / LOW_STOCK_THRESHOLD) * 50, 50)
          status = 'active'
        } else {
          percent = Math.min(50 + ((value - LOW_STOCK_THRESHOLD) / 50) * 50, 100)
          status = 'success'
        }
        
        return (
          <Progress
            percent={percent}
            size="small"
            status={status}
            format={(percent) => `${Math.round(percent)}%`}
            strokeWidth={8}
          />
        )
      }
    },
    {
      title: '单价 (元/kg)',
      dataIndex: 'unit_price',
      key: 'unit_price',
      width: 120,
      render: (price) => price ? `¥${price.toFixed(2)}` : '-'
    },
    {
      title: '库存价值 (元)',
      key: 'inventory_value',
      width: 150,
      render: (_, record) => {
        const value = (record.current_weight_kg || 0) * (record.unit_price || 0)
        return (
          <span style={{ fontWeight: 'bold', color: '#1890ff' }}>
            ¥{value.toFixed(2)}
          </span>
        )
      }
    },
    {
      title: '最后更新',
      dataIndex: 'updated_at',
      key: 'updated_at',
      width: 180,
      render: (date) => date ? new Date(date).toLocaleString('zh-CN') : '-'
    }
  ]

  // Load categories
  const loadCategories = async () => {
    try {
      const result = await categoriesAPI.getAll()
      setCategories(result || [])
    } catch (error) {
      message.error('加载品类失败：' + error.message)
    }
  }

  // Load inventory data
  const loadData = async (searchParams = {}) => {
    setLoading(true)
    try {
      const result = await inventoryAPI.getAll()
      let inventoryData = result || []
      
      // Enhance data with category information
      inventoryData = inventoryData.map(item => {
        const category = categories.find(cat => cat.id === item.category_id)
        return {
          ...item,
          category_name: category?.name || '未知品类',
          unit_price: category?.unit_price || 0
        }
      })
      
      // Apply search filters
      if (searchParams.category_name) {
        inventoryData = inventoryData.filter(item => 
          item.category_name?.toLowerCase().includes(searchParams.category_name.toLowerCase())
        )
      }
      
      if (searchParams.stock_status) {
        inventoryData = inventoryData.filter(item => {
          const weight = item.current_weight_kg || 0
          switch (searchParams.stock_status) {
            case 'out_of_stock':
              return weight === 0
            case 'low_stock':
              return weight > 0 && weight <= LOW_STOCK_THRESHOLD
            case 'in_stock':
              return weight > LOW_STOCK_THRESHOLD
            default:
              return true
          }
        })
      }
      
      setData(inventoryData)
      calculateStats(inventoryData)
    } catch (error) {
      message.error('加载库存数据失败：' + error.message)
    } finally {
      setLoading(false)
    }
  }

  // Calculate statistics
  const calculateStats = (inventoryData) => {
    const stats = {
      totalCategories: inventoryData.length,
      totalWeight: inventoryData.reduce((sum, item) => sum + (item.current_weight_kg || 0), 0),
      lowStockCount: inventoryData.filter(item => {
        const weight = item.current_weight_kg || 0
        return weight > 0 && weight <= LOW_STOCK_THRESHOLD
      }).length,
      outOfStockCount: inventoryData.filter(item => (item.current_weight_kg || 0) === 0).length
    }
    setStats(stats)
  }

  // Search
  const handleSearch = (values) => {
    loadData(values)
  }

  // Reset search
  const handleReset = () => {
    searchForm.resetFields()
    loadData()
  }

  // Load data on mount
  useEffect(() => {
    loadCategories()
  }, [])

  useEffect(() => {
    if (categories.length > 0) {
      loadData()
    }
  }, [categories])

  // Get alerts
  const getAlerts = () => {
    const alerts = []
    
    if (stats.outOfStockCount > 0) {
      alerts.push({
        type: 'error',
        message: `有 ${stats.outOfStockCount} 个品类缺货`,
        description: '请及时补充库存以避免业务中断'
      })
    }
    
    if (stats.lowStockCount > 0) {
      alerts.push({
        type: 'warning',
        message: `有 ${stats.lowStockCount} 个品类库存不足`,
        description: `库存低于 ${LOW_STOCK_THRESHOLD}kg 的品类需要关注`
      })
    }
    
    return alerts
  }

  return (
    <div>
      <PageHeader
        title="库存管理"
        subTitle="查看和监控废旧电池的库存状态"
        icon={<DatabaseOutlined />}
        extra={
          <Space>
            <Button 
              type="primary" 
              icon={<SearchOutlined />}
              onClick={() => loadData()}
            >
              刷新数据
            </Button>
            <ExportButton 
              data={data}
              filename="inventory_report"
              title="导出库存报告"
            />
          </Space>
        }
      >
        {/* Statistics Cards */}
        <Row gutter={16} style={{ marginBottom: 24 }}>
          <Col xs={24} sm={12} lg={6}>
            <Card>
              <Statistic
                title="品类总数"
                value={stats.totalCategories}
                suffix="种"
                valueStyle={{ color: '#1890ff' }}
              />
            </Card>
          </Col>
          <Col xs={24} sm={12} lg={6}>
            <Card>
              <Statistic
                title="总库存重量"
                value={stats.totalWeight}
                suffix="kg"
                precision={2}
                valueStyle={{ color: '#52c41a' }}
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
                prefix={<ExclamationCircleOutlined />}
              />
            </Card>
          </Col>
        </Row>

        {/* Alerts */}
        {getAlerts().map((alert, index) => (
          <Alert
            key={index}
            type={alert.type}
            message={alert.message}
            description={alert.description}
            showIcon
            style={{ marginBottom: 16 }}
          />
        ))}

        <SearchForm
          form={searchForm}
          onFinish={handleSearch}
          onReset={handleReset}
          loading={loading}
        >
          <Form.Item name="category_name" label="品类名称">
            <Input placeholder="请输入品类名称" />
          </Form.Item>
          
          <Form.Item name="stock_status" label="库存状态">
            <Select 
              placeholder="请选择库存状态"
              allowClear
              options={[
                { value: 'in_stock', label: '库存充足' },
                { value: 'low_stock', label: '库存不足' },
                { value: 'out_of_stock', label: '缺货' }
              ]}
            />
          </Form.Item>
        </SearchForm>
      </PageHeader>

      <DataTable
        dataSource={data}
        columns={columns}
        loading={loading}
        rowKey="category_id"
        pagination={{
          showSizeChanger: true,
          showQuickJumper: true,
          showTotal: (total) => `共 ${total} 条记录`,
          defaultPageSize: 20
        }}
        rowClassName={(record) => {
          const weight = record.current_weight_kg || 0
          if (weight === 0) return 'row-out-of-stock'
          if (weight <= LOW_STOCK_THRESHOLD) return 'row-low-stock'
          return ''
        }}
      />

      {data.length === 0 && !loading && (
        <EmptyState
          description="暂无库存数据"
          actionText="刷新数据"
          onAction={() => loadData()}
        />
      )}

      <style jsx>{`
        .row-out-of-stock {
          background-color: #fff2f0 !important;
        }
        .row-low-stock {
          background-color: #fffbe6 !important;
        }
      `}</style>
    </div>
  )
}

export default Inventory