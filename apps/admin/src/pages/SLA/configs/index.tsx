import React, { useRef, useState } from 'react';
import { ProTable } from '@ant-design/pro-components';
import type { ProColumns, ActionType } from '@ant-design/pro-components';
import { Button, Space, Tag, message, Modal, Form, Input, InputNumber, Switch } from 'antd';
import { PlusOutlined } from '@ant-design/icons';
import { TICKET_PRIORITY_MAP } from '@/utils/constants';
import { listSLAConfigs, deleteSLAConfig, createSLAConfig, updateSLAConfig } from '@/services/sla';

const SLAConfigsPage: React.FC = () => {
  const actionRef = useRef<ActionType>();
  const [modalOpen, setModalOpen] = useState(false);
  const [editingId, setEditingId] = useState<number | null>(null);
  const [submitting, setSubmitting] = useState(false);
  const [form] = Form.useForm();

  const openCreate = () => {
    setEditingId(null);
    form.resetFields();
    form.setFieldsValue({ enabled: true });
    setModalOpen(true);
  };

  const openEdit = (record: API.SLAConfig) => {
    setEditingId(record.id);
    form.setFieldsValue({
      name: record.name,
      priority: record.priority,
      first_response_time: record.first_response_time,
      resolution_time: record.resolution_time,
      enabled: record.enabled,
    });
    setModalOpen(true);
  };

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields();
      setSubmitting(true);
      if (editingId) {
        await updateSLAConfig(editingId, values);
        message.success('策略已更新');
      } else {
        await createSLAConfig(values);
        message.success('策略已创建');
      }
      setModalOpen(false);
      form.resetFields();
      actionRef.current?.reload();
    } catch (error: any) {
      if (error?.errorFields) return;
      message.error('操作失败: ' + (error?.message || '未知错误'));
    } finally {
      setSubmitting(false);
    }
  };

  const handleDelete = async (id: number) => {
    Modal.confirm({
      title: '确认删除',
      content: '确定要删除此 SLA 策略吗？',
      okText: '确认删除',
      okButtonProps: { danger: true },
      onOk: async () => {
        try {
          await deleteSLAConfig(id);
          message.success('策略已删除');
          actionRef.current?.reload();
        } catch (error) {
          message.error('删除失败');
        }
      },
    });
  };

  const columns: ProColumns<API.SLAConfig>[] = [
    { title: 'ID', dataIndex: 'id', width: 80 },
    { title: '策略名称', dataIndex: 'name', search: true },
    {
      title: '优先级',
      dataIndex: 'priority',
      width: 100,
      render: (_, record) => {
        const item = TICKET_PRIORITY_MAP[record.priority];
        return item ? <Tag color={item.color}>{item.text}</Tag> : record.priority;
      },
    },
    { title: '首次响应时间(分)', dataIndex: 'first_response_time', width: 150 },
    { title: '解决时间(分)', dataIndex: 'resolution_time', width: 140 },
    {
      title: '状态',
      dataIndex: 'enabled',
      width: 80,
      render: (_, record) => (
        <Tag color={record.enabled ? 'green' : 'default'}>
          {record.enabled ? '启用' : '停用'}
        </Tag>
      ),
    },
    {
      title: '操作',
      valueType: 'option',
      width: 120,
      render: (_, record) => (
        <Space>
          <a onClick={() => openEdit(record)}>编辑</a>
          <a onClick={() => handleDelete(record.id)} style={{ color: '#ff4d4f' }}>删除</a>
        </Space>
      ),
    },
  ];

  return (
    <>
      <ProTable<API.SLAConfig>
        headerTitle="SLA 策略配置"
        rowKey="id"
        actionRef={actionRef}
        columns={columns}
        toolBarRender={() => [
          <Button key="create" type="primary" icon={<PlusOutlined />} onClick={openCreate}>
            新建策略
          </Button>,
        ]}
        request={async (params) => {
          try {
            const result = await listSLAConfigs({
              page: params.current,
              page_size: params.pageSize,
            });
            return {
              data: result?.data || [],
              total: result?.total || 0,
              success: true,
            };
          } catch (error) {
            console.error('获取 SLA 配置失败:', error);
            return { data: [], total: 0, success: true };
          }
        }}
        pagination={{ defaultPageSize: 20 }}
      />

      <Modal
        title={editingId ? '编辑 SLA 策略' : '新建 SLA 策略'}
        open={modalOpen}
        onCancel={() => { setModalOpen(false); form.resetFields(); }}
        onOk={handleSubmit}
        confirmLoading={submitting}
        okText={editingId ? '保存' : '创建'}
      >
        <Form form={form} layout="vertical" style={{ marginTop: 16 }}>
          <Form.Item name="name" label="策略名称" rules={[{ required: true, message: '请输入策略名称' }]}>
            <Input placeholder="如: 标准SLA" />
          </Form.Item>
          <Form.Item name="priority" label="适用优先级" rules={[{ required: true, message: '请输入优先级' }]}>
            <Input placeholder="如: high, medium, low" />
          </Form.Item>
          <Form.Item name="first_response_time" label="首次响应时间(分钟)" rules={[{ required: true, message: '请输入时间' }]}>
            <InputNumber min={1} style={{ width: '100%' }} placeholder="如: 30" />
          </Form.Item>
          <Form.Item name="resolution_time" label="解决时间(分钟)" rules={[{ required: true, message: '请输入时间' }]}>
            <InputNumber min={1} style={{ width: '100%' }} placeholder="如: 480" />
          </Form.Item>
          <Form.Item name="enabled" label="启用" valuePropName="checked">
            <Switch checkedChildren="启用" unCheckedChildren="停用" />
          </Form.Item>
        </Form>
      </Modal>
    </>
  );
};

export default SLAConfigsPage;
