import React from 'react'
import { Table, Button, Space, Popconfirm } from 'antd'
import { EditOutlined, DeleteOutlined, EyeOutlined } from '@ant-design/icons'

const DataTable = ({
  dataSource,
  columns,
  loading = false,
  pagination = true,
  onView,
  onEdit,
  onDelete,
  showActions = true,
  actionWidth = 120,
  deleteConfirmTitle = '确定要删除这条记录吗？',
  rowKey = 'id',
  ...tableProps
}) => {
  // Build action column
  const actionColumn = showActions ? {
    title: '操作',
    key: 'actions',
    width: actionWidth,
    fixed: 'right',
    render: (_, record) => (
      <Space size="small">
        {onView && (
          <Button
            type="link"
            size="small"
            icon={<EyeOutlined />}
            onClick={() => onView(record)}
          >
            查看
          </Button>
        )}
        {onEdit && (
          <Button
            type="link"
            size="small"
            icon={<EditOutlined />}
            onClick={() => onEdit(record)}
          >
            编辑
          </Button>
        )}
        {onDelete && (
          <Popconfirm
            title={deleteConfirmTitle}
            onConfirm={() => onDelete(record)}
            okText="确定"
            cancelText="取消"
          >
            <Button
              type="link"
              size="small"
              danger
              icon={<DeleteOutlined />}
            >
              删除
            </Button>
          </Popconfirm>
        )}
      </Space>
    ),
  } : null

  // Combine columns with action column
  const finalColumns = actionColumn 
    ? [...columns, actionColumn]
    : columns

  return (
    <Table
      rowKey={rowKey}
      dataSource={dataSource}
      columns={finalColumns}
      loading={loading}
      pagination={pagination}
      scroll={{ x: 'max-content' }}
      {...tableProps}
    />
  )
}

export default DataTable