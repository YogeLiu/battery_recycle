import React, { useEffect } from 'react'
import { Modal, Form } from 'antd'

const FormModal = ({
  title,
  visible,
  onCancel,
  onSubmit,
  form,
  initialValues,
  loading = false,
  width = 600,
  children,
  okText = '确定',
  cancelText = '取消'
}) => {
  // Set initial values when modal opens
  useEffect(() => {
    if (visible && initialValues) {
      form.setFieldsValue(initialValues)
    }
  }, [visible, initialValues, form])

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields()
      await onSubmit(values)
    } catch (error) {
      console.error('Form validation failed:', error)
    }
  }

  const handleCancel = () => {
    form.resetFields()
    onCancel()
  }

  return (
    <Modal
      title={title}
      open={visible}
      onOk={handleSubmit}
      onCancel={handleCancel}
      width={width}
      confirmLoading={loading}
      okText={okText}
      cancelText={cancelText}
      destroyOnClose
    >
      <Form
        form={form}
        layout="vertical"
        preserve={false}
      >
        {children}
      </Form>
    </Modal>
  )
}

export default FormModal