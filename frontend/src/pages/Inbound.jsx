import React, { useState, useEffect } from 'react'
import {
  Form,
  Input,
  InputNumber,
  Select,
  Button,
  message,
  Space,
  Card,
  Table,
  Divider,
  Row,
  Col,
  Modal,
  DatePicker,
  Tag
} from 'antd'
import {
  InboxOutlined,
  PlusOutlined,
  DeleteOutlined,
  EyeOutlined,
  PrinterOutlined,
  SearchOutlined,
  ClearOutlined
} from '@ant-design/icons'
import {
  PageHeader,
  SearchForm,
  DataTable,
  EmptyState,
  ExportButton,
  PrintButton
} from '../components'
import { inboundAPI } from '../api/inbound'
import { categoriesAPI } from '../api/categories'
import { useAuth } from '../contexts/AuthContext'

const { RangePicker } = DatePicker

const Inbound = () => {
  const [data, setData] = useState([])
  const [categories, setCategories] = useState([])
  const [loading, setLoading] = useState(false)
  const [modalVisible, setModalVisible] = useState(false)
  const [viewModalVisible, setViewModalVisible] = useState(false)
  const [viewingRecord, setViewingRecord] = useState(null)
  const [editingRecord, setEditingRecord] = useState(null)
  const [searchForm] = Form.useForm()
  const [orderForm] = Form.useForm()
  const { user } = useAuth()

  // Order form state
  const [orderItems, setOrderItems] = useState([])
  const [totalAmount, setTotalAmount] = useState(0)

  // Table columns for order list
  const columns = [
    {
      title: '入库单号',
      dataIndex: 'order_number',
      key: 'order_number',
      width: 160,
      render: (text) => <strong>{text}</strong>
    },
    {
      title: '供应商',
      dataIndex: 'supplier_name',
      key: 'supplier_name',
      width: 150,
    },
    {
      title: '总重量 (kg)',
      dataIndex: 'total_weight_kg',
      key: 'total_weight_kg',
      width: 120,
      render: (weight) => weight?.toFixed(2) || '0.00'
    },
    {
      title: '总金额 (元)',
      dataIndex: 'total_amount',
      key: 'total_amount',
      width: 120,
      render: (amount) => `¥${amount?.toFixed(2) || '0.00'}`
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      width: 100,
      render: (status) => {
        const statusMap = {
          'pending': { color: 'orange', text: '待入库' },
          'completed': { color: 'green', text: '已完成' },
          'cancelled': { color: 'red', text: '已取消' }
        }
        const statusInfo = statusMap[status] || { color: 'default', text: status }
        return <Tag color={statusInfo.color}>{statusInfo.text}</Tag>
      }
    },
    {
      title: '操作员',
      dataIndex: 'operator_name',
      key: 'operator_name',
      width: 100,
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      key: 'created_at',
      width: 180,
      render: (date) => new Date(date).toLocaleString('zh-CN')
    }
  ]

  // Item columns for order form
  const itemColumns = [
    {
      title: '品类',
      dataIndex: 'category_id',
      key: 'category_id',
      width: 150,
      render: (_, record, index) => (
        <Select
          value={record.category_id}
          onChange={(value) => handleItemChange(index, 'category_id', value)}
          placeholder="选择品类"
          style={{ width: '100%' }}
          options={categories.map(cat => ({
            value: cat.id,
            label: cat.name
          }))}
        />
      )
    },
    {
      title: '毛重 (kg)',
      dataIndex: 'gross_weight_kg',
      key: 'gross_weight_kg',
      width: 120,
      render: (_, record, index) => (
        <InputNumber
          value={record.gross_weight_kg}
          onChange={(value) => handleItemChange(index, 'gross_weight_kg', value)}
          placeholder="毛重"
          precision={2}
          min={0}
          style={{ width: '100%' }}
        />
      )
    },
    {
      title: '皮重 (kg)',
      dataIndex: 'tare_weight_kg',
      key: 'tare_weight_kg',
      width: 120,
      render: (_, record, index) => (
        <InputNumber
          value={record.tare_weight_kg}
          onChange={(value) => handleItemChange(index, 'tare_weight_kg', value)}
          placeholder="皮重"
          precision={2}
          min={0}
          style={{ width: '100%' }}
        />
      )
    },
    {
      title: '净重 (kg)',
      dataIndex: 'net_weight_kg',
      key: 'net_weight_kg',
      width: 120,
      render: (_, record) => (
        <strong>{((record.gross_weight_kg || 0) - (record.tare_weight_kg || 0)).toFixed(2)}</strong>
      )
    },
    {
      title: '单价 (元/kg)',
      dataIndex: 'unit_price',
      key: 'unit_price',
      width: 120,
      render: (_, record, index) => (
        <InputNumber
          value={record.unit_price}
          onChange={(value) => handleItemChange(index, 'unit_price', value)}
          placeholder="单价"
          precision={2}
          min={0}
          style={{ width: '100%' }}
        />
      )
    },
    {
      title: '小计 (元)',
      dataIndex: 'subtotal_amount',
      key: 'subtotal_amount',
      width: 120,
      render: (_, record) => {
        const netWeight = (record.gross_weight_kg || 0) - (record.tare_weight_kg || 0)
        const subtotal = netWeight * (record.unit_price || 0)
        return <strong>¥{subtotal.toFixed(2)}</strong>
      }
    },
    {
      title: '操作',
      key: 'action',
      width: 80,
      render: (_, record, index) => (
        <Button
          type="text"
          danger
          icon={<DeleteOutlined />}
          onClick={() => removeItem(index)}
          size="small"
        />
      )
    }
  ]

  // Load categories
  const loadCategories = async () => {
    try {
      const result = await categoriesAPI.getAll()
      setCategories(result?.filter(cat => cat.is_active) || [])
    } catch (error) {
      message.error('加载品类失败：' + error.message)
    }
  }

  // Load data
  const loadData = async (searchParams = {}) => {
    setLoading(true)
    try {
      const result = await inboundAPI.getAll(searchParams)
      setData(result || [])
    } catch (error) {
      message.error('加载数据失败：' + error.message)
    } finally {
      setLoading(false)
    }
  }

  // Search
  const handleSearch = (values) => {
    const searchParams = {}
    if (values.order_number) searchParams.order_number = values.order_number
    if (values.supplier_name) searchParams.supplier_name = values.supplier_name
    if (values.status) searchParams.status = values.status
    if (values.date_range) {
      searchParams.start_date = values.date_range[0].format('YYYY-MM-DD')
      searchParams.end_date = values.date_range[1].format('YYYY-MM-DD')
    }
    loadData(searchParams)
  }

  // Reset search
  const handleReset = () => {
    searchForm.resetFields()
    loadData()
  }

  // Create/Edit order
  const handleCreateOrder = () => {
    setEditingRecord(null)
    setModalVisible(true)
    setOrderItems([])
    setTotalAmount(0)
    orderForm.resetFields()
  }

  // Add item
  const addItem = () => {
    const newItem = {
      id: Date.now(),
      category_id: null,
      gross_weight_kg: 0,
      tare_weight_kg: 0,
      unit_price: 0
    }
    setOrderItems([...orderItems, newItem])
  }

  // Remove item
  const removeItem = (index) => {
    const newItems = [...orderItems]
    newItems.splice(index, 1)
    setOrderItems(newItems)
    calculateTotal(newItems)
  }

  // Handle item change
  const handleItemChange = (index, field, value) => {
    const newItems = [...orderItems]
    newItems[index][field] = value

    // Auto-fill unit price from category
    if (field === 'category_id') {
      const category = categories.find(cat => cat.id === value)
      if (category) {
        newItems[index].unit_price = category.unit_price
      }
    }

    setOrderItems(newItems)
    calculateTotal(newItems)
  }

  // Calculate total
  const calculateTotal = (items) => {
    const total = items.reduce((sum, item) => {
      const netWeight = (item.gross_weight_kg || 0) - (item.tare_weight_kg || 0)
      const subtotal = netWeight * (item.unit_price || 0)
      return sum + subtotal
    }, 0)
    setTotalAmount(total)
  }

  // Save order
  const handleSaveOrder = async (values) => {
    if (orderItems.length === 0) {
      message.error('请至少添加一个项目')
      return
    }

    try {
      const orderData = {
        ...values,
        items: orderItems.map(item => ({
          category_id: item.category_id,
          gross_weight_kg: item.gross_weight_kg,
          tare_weight_kg: item.tare_weight_kg,
          net_weight_kg: (item.gross_weight_kg || 0) - (item.tare_weight_kg || 0),
          unit_price: item.unit_price,
          subtotal_amount: ((item.gross_weight_kg || 0) - (item.tare_weight_kg || 0)) * (item.unit_price || 0)
        })),
        total_amount: totalAmount,
        total_weight_kg: orderItems.reduce((sum, item) =>
          sum + ((item.gross_weight_kg || 0) - (item.tare_weight_kg || 0)), 0
        )
      }

      await inboundAPI.create(orderData)
      message.success('入库单创建成功')
      setModalVisible(false)
      loadData()
    } catch (error) {
      message.error('保存失败：' + error.message)
    }
  }

  // View order details
  const handleViewOrder = async (record) => {
    try {
      const result = await inboundAPI.getById(record.id)
      setViewingRecord(result)
      setViewModalVisible(true)
    } catch (error) {
      message.error('加载详情失败：' + error.message)
    }
  }

  // Cancel modal
  const handleCancel = () => {
    setModalVisible(false)
    orderForm.resetFields()
    setOrderItems([])
    setTotalAmount(0)
  }

  // Load data on mount
  useEffect(() => {
    loadData()
    loadCategories()
  }, [])

  return (
    <div>
      <PageHeader
        title="采购入库管理"
        subTitle="管理废旧电池的采购入库业务"
        icon={<InboxOutlined />}
        actions={[
          {
            label: '新建入库单',
            props: {
              type: 'primary',
              icon: <PlusOutlined />,
              onClick: handleCreateOrder
            }
          }
        ]}
        extra={
          <Space>
            <ExportButton
              data={data}
              filename="inbound_orders"
              title="导出数据"
            />
          </Space>
        }
      >
        <SearchForm
          form={searchForm}
          onFinish={handleSearch}
          onReset={handleReset}
          loading={loading}
        >
          <Form.Item name="order_number" label="入库单号">
            <Input placeholder="请输入入库单号" />
          </Form.Item>

          <Form.Item name="supplier_name" label="供应商">
            <Input placeholder="请输入供应商名称" />
          </Form.Item>

          <Form.Item name="status" label="状态">
            <Select
              placeholder="请选择状态"
              allowClear
              options={[
                { value: 'pending', label: '待入库' },
                { value: 'completed', label: '已完成' },
                { value: 'cancelled', label: '已取消' }
              ]}
            />
          </Form.Item>

          <Form.Item name="date_range" label="创建时间">
            <RangePicker />
          </Form.Item>
        </SearchForm>
      </PageHeader>

      <DataTable
        dataSource={data}
        columns={columns}
        loading={loading}
        onView={handleViewOrder}
        rowKey="id"
        pagination={{
          showSizeChanger: true,
          showQuickJumper: true,
          showTotal: (total) => `共 ${total} 条记录`
        }}
      />

      {data.length === 0 && !loading && (
        <EmptyState
          description="暂无入库数据"
          actionText="新建入库单"
          onAction={handleCreateOrder}
        />
      )}

      {/* Create/Edit Order Modal */}
      <Modal
        title="新建入库单"
        visible={modalVisible}
        onCancel={handleCancel}
        footer={null}
        width={1200}
        destroyOnClose
      >
        <Form
          form={orderForm}
          layout="vertical"
          onFinish={handleSaveOrder}
        >
          <Row gutter={16}>
            <Col span={12}>
              <Form.Item
                name="supplier_name"
                label="供应商"
                rules={[{ required: true, message: '请输入供应商名称' }]}
              >
                <Input placeholder="请输入供应商名称" />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item
                name="contact_person"
                label="联系人"
              >
                <Input placeholder="请输入联系人" />
              </Form.Item>
            </Col>
          </Row>

          <Row gutter={16}>
            <Col span={12}>
              <Form.Item
                name="contact_phone"
                label="联系电话"
              >
                <Input placeholder="请输入联系电话" />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item
                name="vehicle_number"
                label="车牌号"
              >
                <Input placeholder="请输入车牌号" />
              </Form.Item>
            </Col>
          </Row>

          <Form.Item
            name="remarks"
            label="备注"
          >
            <Input.TextArea placeholder="请输入备注信息" rows={3} />
          </Form.Item>

          <Divider>入库项目</Divider>

          <div style={{ marginBottom: 16 }}>
            <Button
              type="dashed"
              onClick={addItem}
              icon={<PlusOutlined />}
              style={{ width: '100%' }}
            >
              添加项目
            </Button>
          </div>

          <Table
            dataSource={orderItems}
            columns={itemColumns}
            pagination={false}
            rowKey="id"
            scroll={{ x: 800 }}
            locale={{ emptyText: '请添加入库项目' }}
          />

          <div style={{ marginTop: 16, textAlign: 'right' }}>
            <Space size="large">
              <span>
                总重量: <strong>{orderItems.reduce((sum, item) =>
                  sum + ((item.gross_weight_kg || 0) - (item.tare_weight_kg || 0)), 0
                ).toFixed(2)} kg</strong>
              </span>
              <span>
                总金额: <strong style={{ color: '#f5222d', fontSize: '18px' }}>
                  ¥{totalAmount.toFixed(2)}
                </strong>
              </span>
            </Space>
          </div>

          <div className="form-actions">
            <Button onClick={handleCancel}>
              取消
            </Button>
            <Button
              type="primary"
              htmlType="submit"
              disabled={orderItems.length === 0}
            >
              创建入库单
            </Button>
          </div>
        </Form>
      </Modal>

      {/* View Order Modal */}
      <Modal
        title="入库单详情"
        visible={viewModalVisible}
        onCancel={() => setViewModalVisible(false)}
        footer={[
          <PrintButton
            key="print"
            data={viewingRecord}
            orderType="inbound"
            title="打印入库单"
          />,
          <Button key="close" onClick={() => setViewModalVisible(false)}>
            关闭
          </Button>
        ]}
        width={1000}
      >
        {viewingRecord && (
          <div>
            <Row gutter={16}>
              <Col span={12}>
                <p><strong>入库单号:</strong> {viewingRecord.order_number}</p>
                <p><strong>供应商:</strong> {viewingRecord.supplier_name}</p>
                <p><strong>联系人:</strong> {viewingRecord.contact_person}</p>
              </Col>
              <Col span={12}>
                <p><strong>联系电话:</strong> {viewingRecord.contact_phone}</p>
                <p><strong>车牌号:</strong> {viewingRecord.vehicle_number}</p>
                <p><strong>操作员:</strong> {viewingRecord.operator_name}</p>
              </Col>
            </Row>

            <Divider>项目明细</Divider>

            <Table
              dataSource={viewingRecord.items || []}
              columns={[
                { title: '品类', dataIndex: 'category_name', key: 'category_name' },
                { title: '毛重 (kg)', dataIndex: 'gross_weight_kg', key: 'gross_weight_kg', render: val => val?.toFixed(2) },
                { title: '皮重 (kg)', dataIndex: 'tare_weight_kg', key: 'tare_weight_kg', render: val => val?.toFixed(2) },
                { title: '净重 (kg)', dataIndex: 'net_weight_kg', key: 'net_weight_kg', render: val => val?.toFixed(2) },
                { title: '单价 (元/kg)', dataIndex: 'unit_price', key: 'unit_price', render: val => `¥${val?.toFixed(2)}` },
                { title: '小计 (元)', dataIndex: 'subtotal_amount', key: 'subtotal_amount', render: val => `¥${val?.toFixed(2)}` }
              ]}
              pagination={false}
              size="small"
            />

            <div style={{ marginTop: 16, textAlign: 'right' }}>
              <Space size="large">
                <span>总重量: <strong>{viewingRecord.total_weight_kg?.toFixed(2)} kg</strong></span>
                <span>总金额: <strong style={{ color: '#f5222d', fontSize: '16px' }}>¥{viewingRecord.total_amount?.toFixed(2)}</strong></span>
              </Space>
            </div>

            {viewingRecord.remarks && (
              <>
                <Divider>备注</Divider>
                <p>{viewingRecord.remarks}</p>
              </>
            )}
          </div>
        )}
      </Modal>
    </div>
  )
}

export default Inbound