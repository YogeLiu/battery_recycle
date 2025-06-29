import React, { useState, useEffect } from 'react'
import { Form, Input, Select, Switch, Button, message, Space, Tag } from 'antd'
import { UserOutlined, PlusOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons'
import { 
  PageHeader, 
  SearchForm, 
  DataTable, 
  FormModal, 
  StatusTag,
  EmptyState,
  ExportButton
} from '../components'
import { usersAPI } from '../api/users'
import { useAuth } from '../contexts/AuthContext'

const Users = () => {
  const [data, setData] = useState([])
  const [loading, setLoading] = useState(false)
  const [modalVisible, setModalVisible] = useState(false)
  const [editingRecord, setEditingRecord] = useState(null)
  const [searchForm] = Form.useForm()
  const [modalForm] = Form.useForm()
  const { isAdmin, user: currentUser } = useAuth()

  // User roles
  const userRoles = [
    { value: 'super_admin', label: '超级管理员', color: 'red' },
    { value: 'normal_user', label: '普通用户', color: 'blue' }
  ]

  // Table columns
  const columns = [
    {
      title: 'ID',
      dataIndex: 'id',
      key: 'id',
      width: 80,
    },
    {
      title: '用户名',
      dataIndex: 'username',
      key: 'username',
      width: 120,
      render: (text) => <strong>{text}</strong>
    },
    {
      title: '真实姓名',
      dataIndex: 'real_name',
      key: 'real_name',
      width: 120,
    },
    {
      title: '角色',
      dataIndex: 'role',
      key: 'role',
      width: 120,
      render: (role) => {
        const roleInfo = userRoles.find(r => r.value === role)
        return roleInfo ? (
          <Tag color={roleInfo.color}>{roleInfo.label}</Tag>
        ) : role
      }
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
    },
    {
      title: '最后登录',
      dataIndex: 'last_login_at',
      key: 'last_login_at',
      width: 180,
      render: (date) => date ? new Date(date).toLocaleString('zh-CN') : '从未登录'
    }
  ]

  // Load data
  const loadData = async (searchParams = {}) => {
    setLoading(true)
    try {
      const result = await usersAPI.getAll()
      let filteredData = result || []
      
      // Apply search filters
      if (searchParams.username) {
        filteredData = filteredData.filter(item => 
          item.username?.toLowerCase().includes(searchParams.username.toLowerCase())
        )
      }
      if (searchParams.real_name) {
        filteredData = filteredData.filter(item => 
          item.real_name?.toLowerCase().includes(searchParams.real_name.toLowerCase())
        )
      }
      if (searchParams.role) {
        filteredData = filteredData.filter(item => item.role === searchParams.role)
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
      modalForm.setFieldsValue({
        ...record,
        password: '' // Don't show existing password
      })
    }
  }

  // Save
  const handleSave = async (values) => {
    try {
      // Remove empty password for updates
      if (editingRecord && !values.password) {
        delete values.password
      }
      
      if (editingRecord) {
        await usersAPI.update(editingRecord.id, values)
        message.success('更新成功')
      } else {
        await usersAPI.create(values)
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
    // Prevent deleting self
    if (record.id === currentUser?.id) {
      message.error('不能删除自己的账号')
      return
    }
    
    try {
      await usersAPI.delete(record.id)
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

  // Only allow admin access
  if (!isAdmin()) {
    return (
      <div className="empty-container">
        <h3>访问被拒绝</h3>
        <p>您需要管理员权限才能访问此页面</p>
      </div>
    )
  }

  return (
    <div>
      <PageHeader
        title="用户管理"
        subTitle="管理系统用户账号和权限"
        icon={<UserOutlined />}
        actions={[
          {
            label: '新增用户',
            props: {
              type: 'primary',
              icon: <PlusOutlined />,
              onClick: () => handleEdit()
            }
          }
        ]}
        extra={
          <Space>
            <ExportButton 
              data={data.map(item => ({
                ...item,
                role_name: userRoles.find(r => r.value === item.role)?.label || item.role
              }))}
              filename="system_users"
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
            name="username" 
            label="用户名"
          >
            <Input placeholder="请输入用户名" />
          </Form.Item>
          
          <Form.Item 
            name="real_name" 
            label="真实姓名"
          >
            <Input placeholder="请输入真实姓名" />
          </Form.Item>
          
          <Form.Item 
            name="role" 
            label="角色"
          >
            <Select 
              placeholder="请选择角色"
              allowClear
              options={userRoles}
            />
          </Form.Item>
          
          <Form.Item 
            name="is_active" 
            label="状态"
          >
            <Select 
              placeholder="请选择状态"
              allowClear
              options={[
                { value: true, label: '启用' },
                { value: false, label: '禁用' }
              ]}
            />
          </Form.Item>
        </SearchForm>
      </PageHeader>

      <DataTable
        dataSource={data}
        columns={columns}
        loading={loading}
        onEdit={handleEdit}
        onDelete={handleDelete}
        deleteConfirmTitle="确定要删除此用户吗？删除后无法恢复。"
        rowKey="id"
        canDelete={(record) => record.id !== currentUser?.id}
        pagination={{
          showSizeChanger: true,
          showQuickJumper: true,
          showTotal: (total) => `共 ${total} 条记录`
        }}
      />

      {data.length === 0 && !loading && (
        <EmptyState
          description="暂无用户数据"
          actionText="新增用户"
          onAction={() => handleEdit()}
        />
      )}

      <FormModal
        title={editingRecord ? '编辑用户' : '新增用户'}
        visible={modalVisible}
        onCancel={handleCancel}
        onSubmit={handleSave}
        form={modalForm}
        width={600}
      >
        <Form.Item
          name="username"
          label="用户名"
          rules={[
            { required: true, message: '请输入用户名' },
            { min: 3, message: '用户名至少3个字符' },
            { max: 50, message: '用户名不能超过50字符' },
            { pattern: /^[a-zA-Z0-9_]+$/, message: '用户名只能包含字母、数字和下划线' }
          ]}
        >
          <Input placeholder="请输入用户名" disabled={!!editingRecord} />
        </Form.Item>

        <Form.Item
          name="real_name"
          label="真实姓名"
          rules={[
            { required: true, message: '请输入真实姓名' },
            { max: 50, message: '姓名不能超过50字符' }
          ]}
        >
          <Input placeholder="请输入真实姓名" />
        </Form.Item>

        <Form.Item
          name="password"
          label={editingRecord ? '新密码' : '密码'}
          rules={[
            { required: !editingRecord, message: '请输入密码' },
            { min: 6, message: '密码至少6个字符' }
          ]}
        >
          <Input.Password 
            placeholder={editingRecord ? '留空则不修改密码' : '请输入密码'} 
          />
        </Form.Item>

        <Form.Item
          name="role"
          label="角色"
          rules={[
            { required: true, message: '请选择角色' }
          ]}
          initialValue="normal_user"
        >
          <Select 
            placeholder="请选择角色"
            options={userRoles}
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

export default Users