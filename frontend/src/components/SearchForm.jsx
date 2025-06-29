import React from 'react'
import { Form, Row, Col, Button, Space } from 'antd'
import { SearchOutlined, ReloadOutlined } from '@ant-design/icons'

const SearchForm = ({ 
  form,
  onFinish,
  onReset,
  children,
  loading = false,
  span = 8
}) => {
  const handleReset = () => {
    form.resetFields()
    if (onReset) {
      onReset()
    }
  }

  return (
    <div className="search-form">
      <Form
        form={form}
        layout="vertical"
        onFinish={onFinish}
        autoComplete="off"
      >
        <Row gutter={16}>
          {React.Children.map(children, (child, index) => (
            <Col key={index} span={span}>
              {child}
            </Col>
          ))}
          
          <Col span={span}>
            <Form.Item label=" ">
              <Space>
                <Button 
                  type="primary" 
                  htmlType="submit"
                  icon={<SearchOutlined />}
                  loading={loading}
                >
                  查询
                </Button>
                <Button 
                  icon={<ReloadOutlined />}
                  onClick={handleReset}
                >
                  重置
                </Button>
              </Space>
            </Form.Item>
          </Col>
        </Row>
      </Form>
    </div>
  )
}

export default SearchForm