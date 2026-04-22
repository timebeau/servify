import React, { useRef, useState } from 'react';
import { ProTable } from '@ant-design/pro-components';
import type { ActionType, ProColumns } from '@ant-design/pro-components';
import { PlusOutlined } from '@ant-design/icons';
import { Button, Form, Input, Modal, Space, Switch, message } from 'antd';
import { createAutomation, deleteAutomation, listAutomations, runAutomation } from '@/services/automation';
import { getErrorMessage, isFormValidationError } from '@/utils/error';

function prettyJson(value?: string | Record<string, unknown>) {
  if (!value) {
    return '[]';
  }
  if (typeof value === 'string') {
    try {
      return JSON.stringify(JSON.parse(value), null, 2);
    } catch {
      return value;
    }
  }
  return JSON.stringify(value, null, 2);
}

const AutomationPage: React.FC = () => {
  const actionRef = useRef<ActionType>();
  const [form] = Form.useForm();
  const [modalOpen, setModalOpen] = useState(false);
  const [submitting, setSubmitting] = useState(false);

  const openView = (record: API.Automation) => {
    Modal.info({
      title: record.name,
      width: 760,
      content: (
        <div style={{ marginTop: 16, display: 'grid', gap: 12 }}>
          <div>事件：{record.event || record.trigger_type || '-'}</div>
          <div>状态：{record.active ?? record.enabled ? '启用' : '停用'}</div>
          <div>
            条件：
            <pre style={{ whiteSpace: 'pre-wrap', marginTop: 8 }}>{prettyJson(record.conditions)}</pre>
          </div>
          <div>
            动作：
            <pre style={{ whiteSpace: 'pre-wrap', marginTop: 8 }}>{prettyJson(record.actions)}</pre>
          </div>
        </div>
      ),
    });
  };

  const handleCreate = async () => {
    try {
      const values = await form.validateFields();
      setSubmitting(true);
      await createAutomation({
        name: values.name?.trim(),
        event: values.event?.trim(),
        conditions: JSON.parse(values.conditions),
        actions: JSON.parse(values.actions),
        active: values.active ?? true,
      });
      message.success('规则已创建');
      setModalOpen(false);
      form.resetFields();
      actionRef.current?.reload();
    } catch (error: unknown) {
      if (isFormValidationError(error)) {
        return;
      }
      if (error instanceof SyntaxError) {
        message.error('条件或动作不是合法 JSON');
        return;
      }
      message.error(getErrorMessage(error, '创建失败'));
    } finally {
      setSubmitting(false);
    }
  };

  const columns: ProColumns<API.Automation>[] = [
    {
      title: 'ID',
      dataIndex: 'id',
      width: 80,
    },
    {
      title: '规则名称',
      dataIndex: 'name',
      search: true,
    },
    {
      title: '触发事件',
      dataIndex: 'event',
      width: 180,
      search: false,
    },
    {
      title: '执行动作',
      dataIndex: 'actions',
      width: 240,
      search: false,
      render: (_, record) => {
        if (typeof record.actions === 'string') {
          return record.actions;
        }
        if (record.actions) {
          return JSON.stringify(record.actions);
        }
        return '-';
      },
    },
    {
      title: '状态',
      dataIndex: 'active',
      width: 100,
      search: false,
      render: (_, record) => (
        <Switch
          checked={Boolean(record.active ?? record.enabled)}
          checkedChildren="启用"
          unCheckedChildren="停用"
          disabled
          onChange={() => {}}
        />
      ),
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      valueType: 'dateTime',
      width: 180,
      search: false,
    },
    {
      title: '操作',
      valueType: 'option',
      width: 180,
      render: (_, record) => (
        <Space>
          <a onClick={() => openView(record)}>查看</a>
          <a
            onClick={async () => {
              try {
                await runAutomation(record.id);
                message.success('规则已执行');
              } catch (error: unknown) {
                message.error(getErrorMessage(error, '执行失败'));
              }
            }}
          >
            执行
          </a>
          <a
            onClick={async () => {
              try {
                await deleteAutomation(record.id);
                message.success('规则已删除');
                actionRef.current?.reload();
              } catch (error: unknown) {
                message.error(getErrorMessage(error, '删除失败'));
              }
            }}
          >
            删除
          </a>
        </Space>
      ),
    },
  ];

  return (
    <>
      <ProTable<API.Automation>
        headerTitle="自动化规则"
        rowKey="id"
        actionRef={actionRef}
        columns={columns}
        toolBarRender={() => [
          <Button
            key="create"
            type="primary"
            icon={<PlusOutlined />}
            onClick={() => {
              form.resetFields();
              form.setFieldsValue({
                active: true,
                conditions: '[]',
                actions: '[]',
              });
              setModalOpen(true);
            }}
          >
            新建规则
          </Button>,
        ]}
        request={async (params) => {
          try {
            const result = await listAutomations();
            const keyword = typeof params.name === 'string' ? params.name.trim().toLowerCase() : '';
            let data = result.data;
            if (keyword) {
              data = data.filter((item) => item.name.toLowerCase().includes(keyword));
            }
            const total = data.length;
            const current = params.current || 1;
            const pageSize = params.pageSize || 20;
            return {
              data: data.slice((current - 1) * pageSize, current * pageSize),
              total,
              success: true,
            };
          } catch (error: unknown) {
            console.error('获取自动化规则失败:', error);
            return { data: [], total: 0, success: true };
          }
        }}
        pagination={{ defaultPageSize: 20 }}
      />

      <Modal
        title="新建规则"
        open={modalOpen}
        onCancel={() => {
          setModalOpen(false);
          form.resetFields();
        }}
        onOk={handleCreate}
        confirmLoading={submitting}
        okText="创建"
        width={760}
      >
        <Form form={form} layout="vertical" style={{ marginTop: 16 }}>
          <Form.Item name="name" label="规则名称" rules={[{ required: true, message: '请输入规则名称' }]}>
            <Input placeholder="例如：高优先级自动加急" />
          </Form.Item>
          <Form.Item name="event" label="触发事件" rules={[{ required: true, message: '请输入触发事件' }]}>
            <Input placeholder="例如：ticket_created" />
          </Form.Item>
          <Form.Item
            name="conditions"
            label="条件 JSON"
            rules={[{ required: true, message: '请输入条件 JSON' }]}
          >
            <Input.TextArea rows={6} placeholder='例如：[{"field":"priority","op":"eq","value":"urgent"}]' />
          </Form.Item>
          <Form.Item
            name="actions"
            label="动作 JSON"
            rules={[{ required: true, message: '请输入动作 JSON' }]}
          >
            <Input.TextArea rows={6} placeholder='例如：[{"type":"set_priority","params":{"priority":"urgent"}}]' />
          </Form.Item>
          <Form.Item name="active" label="启用" valuePropName="checked">
            <Switch checkedChildren="启用" unCheckedChildren="停用" />
          </Form.Item>
        </Form>
      </Modal>
    </>
  );
};

export default AutomationPage;
