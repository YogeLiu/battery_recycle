import React, { useState, useEffect } from 'react'
import { Form, Input, InputNumber, Switch, Button, message, Space } from 'antd'
import { AppstoreOutlined, PlusOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons'
import { 
  PageHeader, 
  SearchForm, 
  DataTable, 
  FormModal, 
  StatusTag,
  EmptyState,
  ExportButton
} from '../components'
import { categoriesAPI } from '../api/categories'
import { useAuth } from '../contexts/AuthContext'

const Categories = () => {
  const [data, setData] = useState([])
  const [loading, setLoading] = useState(false)
  const [modalVisible, setModalVisible] = useState(false)
  const [editingRecord, setEditingRecord] = useState(null)
  const [searchForm] = Form.useForm()
  const [modalForm] = Form.useForm()
  const { isAdmin } = useAuth()

  // Table columns
  const columns = [
    {
      title: 'ID',
      dataIndex: 'id',
      key: 'id',
      width: 80,
    },
    {
      title: '品类名称',
      dataIndex: 'name',
      key: 'name',
      width: 150,
      render: (text) => <strong>{text}</strong>
    },
    {
      title: '描述',
      dataIndex: 'description',
      key: 'description',
      ellipsis: true,
    },
    {
      title: '单价 (元/kg)',
      dataIndex: 'unit_price',
      key: 'unit_price',
      width: 120,
      render: (price) => `¥${price?.toFixed(2) || '0.00'}`
    },
    {
      title: '状态',
      dataIndex: 'is_active',
      key: 'is_active',
      width: 100,
      render: (isActive) => <StatusTag status={isActive} />
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      key: 'created_at',
      width: 180,
      render: (date) => new Date(date).toLocaleString('zh-CN')
    }
  ]

  // Load data
  const loadData = async (searchParams = {}) => {
    setLoading(true)
    try {
      const result = await categoriesAPI.getAll()
      let filteredData = result || []
      
      // Apply search filters
      if (searchParams.name) {
        filteredData = filteredData.filter(item => 
          item.name?.toLowerCase().includes(searchParams.name.toLowerCase())
        )
      }
      if (searchParams.is_active !== undefined) {
        filteredData = filteredData.filter(item => item.is_active === searchParams.is_active)
      }
      
      setData(filteredData)
    } catch (error) {
      message.error('加载数据失败：' + error.message)
    } finally {
      setLoading(false)
    }
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

  // Create/Edit
  const handleEdit = (record = null) => {
    setEditingRecord(record)
    setModalVisible(true)
    if (record) {
      modalForm.setFieldsValue(record)
    }
  }

  // Save
  const handleSave = async (values) => {
    try {
      if (editingRecord) {
        await categoriesAPI.update(editingRecord.id, values)
        message.success('更新成功')
      } else {
        await categoriesAPI.create(values)
        message.success('创建成功')
      }
      setModalVisible(false)
      modalForm.resetFields()
      setEditingRecord(null)
      loadData()
    } catch (error) {
      message.error('保存失败：' + error.message)
    }
  }

  // Delete
  const handleDelete = async (record) => {
    try {
      await categoriesAPI.delete(record.id)
      message.success('删除成功')
      loadData()
    } catch (error) {
      message.error('删除失败：' + error.message)
    }
  }

  // Cancel modal
  const handleCancel = () => {
    setModalVisible(false)
    modalForm.resetFields()
    setEditingRecord(null)
  }

  // Load data on mount
  useEffect(() => {
    loadData()
  }, [])

  return (
    <div>
      <PageHeader
        title="电池品类管理"
        subTitle="管理废旧电池的品类信息和价格"
        icon={<AppstoreOutlined />}
        actions={isAdmin() ? [
          {
            label: '新增品类',
            props: {
              type: 'primary',
              icon: <PlusOutlined />,
              onClick: () => handleEdit()
            }
          }
        ] : []}
        extra={
          <Space>
            <ExportButton 
              data={data}
              filename="battery_categories"
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
          <Form.Item 
            name="name" 
            label="品类名称"
          >
            <Input placeholder="请输入品类名称" />
          </Form.Item>
          
          <Form.Item 
            name="is_active" 
            label="状态"
          >
            <Switch 
              checkedChildren="启用" 
              unCheckedChildren="禁用"
              onChange={(checked) => searchForm.setFieldsValue({ is_active: checked })}
            />
          </Form.Item>
        </SearchForm>
      </PageHeader>

      <DataTable
        dataSource={data}
        columns={columns}
        loading={loading}
        onEdit={isAdmin() ? handleEdit : null}
        onDelete={isAdmin() ? handleDelete : null}
        deleteConfirmTitle="确定要删除此品类吗？删除后无法恢复。"
        pagination={{
          showSizeChanger: true,
          showQuickJumper: true,
          showTotal: (total) => `共 ${total} 条记录`
        }}
      />

      {data.length === 0 && !loading && (
        <EmptyState
          description="暂无品类数据"
          actionText={isAdmin() ? "新增品类" : null}
          onAction={isAdmin() ? () => handleEdit() : null}
        />
      )}

      <FormModal
        title={editingRecord ? '编辑品类' : '新增品类'}
        visible={modalVisible}
        onCancel={handleCancel}
        onSubmit={handleSave}
        form={modalForm}
        width={600}
      >
        <Form.Item
          name="name"
          label="品类名称"
          rules={[
            { required: true, message: '请输入品类名称' },
            { max: 100, message: '名称长度不能超过100字符' }
          ]}
        >
          <Input placeholder="请输入品类名称" />
        </Form.Item>

        <Form.Item
          name="description"
          label="描述"
          rules={[
            { max: 255, message: '描述长度不能超过255字符' }
          ]}
        >
          <Input.TextArea 
            placeholder="请输入品类描述" 
            rows={3}
          />
        </Form.Item>

        <Form.Item
          name="unit_price"
          label="单价 (元/kg)"
          rules={[
            { required: true, message: '请输入单价' },
            { type: 'number', min: 0, message: '单价不能为负数' }
          ]}
        >
          <InputNumber
            placeholder="请输入单价"
            precision={2}
            min={0}
            style={{ width: '100%' }}
            addonAfter="元/kg"
          />
        </Form.Item>

        <Form.Item
          name="is_active"
          label="状态"
          valuePropName="checked"
          initialValue={true}
        >
          <Switch 
            checkedChildren="启用" 
            unCheckedChildren="禁用"
          />
        </Form.Item>
      </FormModal>
    </div>
  )
}

export default Categories