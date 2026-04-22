import React, { useRef, useState } from 'react';
import { ProTable } from '@ant-design/pro-components';
import type { ActionType, ProColumns } from '@ant-design/pro-components';
import { PlusOutlined } from '@ant-design/icons';
import { Button, Form, Input, Modal, Select, Space, Switch, Tag, message } from 'antd';
import {
  createCustomField,
  deleteCustomField,
  listCustomFields,
  updateCustomField,
} from '@/services/customField';
import { getErrorMessage, isFormValidationError } from '@/utils/error';

const FIELD_TYPE_OPTIONS = [
  { label: '文本', value: 'string' },
  { label: '数字', value: 'number' },
  { label: '布尔', value: 'boolean' },
  { label: '日期', value: 'date' },
  { label: '单选', value: 'select' },
  { label: '多选', value: 'multiselect' },
];

const RESOURCE_OPTIONS = [{ label: '工单', value: 'ticket' }];

function parseOptions(value?: string) {
  return value
    ? value
        .split(',')
        .map((item) => item.trim())
        .filter(Boolean)
    : [];
}

function parseJsonInput(value?: string) {
  const text = value?.trim();
  if (!text) {
    return undefined;
  }
  return JSON.parse(text);
}

const CustomFieldPage: React.FC = () => {
  const actionRef = useRef<ActionType>();
  const [form] = Form.useForm();
  const [modalOpen, setModalOpen] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  const [editingField, setEditingField] = useState<API.CustomField | null>(null);

  const closeModal = () => {
    setModalOpen(false);
    setEditingField(null);
    form.resetFields();
  };

  const openCreate = () => {
    setEditingField(null);
    form.resetFields();
    form.setFieldsValue({
      resource: 'ticket',
      type: 'string',
      required: false,
      active: true,
    });
    setModalOpen(true);
  };

  const openEdit = (record: API.CustomField) => {
    setEditingField(record);
    form.setFieldsValue({
      name: record.name,
      key: record.key,
      type: record.type || record.field_type || 'string',
      resource: record.resource || record.entity_type || 'ticket',
      required: record.required,
      active: record.active ?? true,
      options:
        Array.isArray(record.options) && record.options.length > 0 ? record.options.join(', ') : '',
      validation: '',
      show_when: '',
    });
    setModalOpen(true);
  };

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields();
      setSubmitting(true);

      const payload = {
        name: values.name?.trim(),
        type: values.type,
        required: values.required ?? false,
        active: values.active ?? true,
        options: ['select', 'multiselect'].includes(values.type)
          ? parseOptions(values.options)
          : undefined,
        validation: parseJsonInput(values.validation),
        show_when: parseJsonInput(values.show_when),
      };

      if (editingField) {
        await updateCustomField(editingField.id, payload);
        message.success('自定义字段已更新');
      } else {
        await createCustomField({
          ...payload,
          key: values.key?.trim(),
          resource: values.resource,
        });
        message.success('自定义字段已创建');
      }

      closeModal();
      actionRef.current?.reload();
    } catch (error: unknown) {
      if (isFormValidationError(error)) {
        return;
      }
      if (error instanceof SyntaxError) {
        message.error('校验规则或显示条件不是合法 JSON');
        return;
      }
      message.error(getErrorMessage(error, editingField ? '更新失败' : '创建失败'));
    } finally {
      setSubmitting(false);
    }
  };

  const columns: ProColumns<API.CustomField>[] = [
    {
      title: 'ID',
      dataIndex: 'id',
      width: 80,
    },
    {
      title: '字段名称',
      dataIndex: 'name',
      search: true,
    },
    {
      title: '字段标识',
      dataIndex: 'key',
    },
    {
      title: '字段类型',
      dataIndex: 'type',
      width: 120,
      search: false,
      render: (_, record) => <Tag>{record.type || record.field_type}</Tag>,
    },
    {
      title: '适用资源',
      dataIndex: 'resource',
      width: 120,
      render: (_, record) => <Tag color="blue">{record.resource || record.entity_type}</Tag>,
    },
    {
      title: '必填',
      dataIndex: 'required',
      width: 80,
      render: (_, record) => (
        <Tag color={record.required ? 'red' : 'default'}>{record.required ? '是' : '否'}</Tag>
      ),
    },
    {
      title: '启用',
      dataIndex: 'active',
      width: 80,
      search: false,
      render: (_, record) => (
        <Tag color={record.active ? 'green' : 'default'}>{record.active ? '启用' : '停用'}</Tag>
      ),
    },
    {
      title: '操作',
      valueType: 'option',
      width: 140,
      render: (_, record) => (
        <Space>
          <a onClick={() => openEdit(record)}>编辑</a>
          <a
            onClick={async () => {
              try {
                await deleteCustomField(record.id);
                message.success('字段已删除');
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
      <ProTable<API.CustomField>
        headerTitle="自定义字段"
        rowKey="id"
        actionRef={actionRef}
        columns={columns}
        toolBarRender={() => [
          <Button key="create" type="primary" icon={<PlusOutlined />} onClick={openCreate}>
            新建字段
          </Button>,
        ]}
        request={async (params) => {
          try {
            const result = await listCustomFields({
              resource: typeof params.resource === 'string' ? params.resource : undefined,
              active: false,
            });
            const keyword = typeof params.name === 'string' ? params.name.trim().toLowerCase() : '';
            const resource = typeof params.resource === 'string' ? params.resource.trim() : '';
            let data = result.data;
            if (keyword) {
              data = data.filter((item) => item.name.toLowerCase().includes(keyword));
            }
            if (resource) {
              data = data.filter((item) => (item.resource || item.entity_type) === resource);
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
            console.error('获取自定义字段失败:', error);
            return { data: [], total: 0, success: true };
          }
        }}
        pagination={{ defaultPageSize: 20 }}
      />

      <Modal
        title={editingField ? '编辑字段' : '新建字段'}
        open={modalOpen}
        onCancel={closeModal}
        onOk={handleSubmit}
        confirmLoading={submitting}
        okText={editingField ? '保存' : '创建'}
        width={720}
      >
        <Form form={form} layout="vertical" style={{ marginTop: 16 }}>
          <Form.Item name="name" label="字段名称" rules={[{ required: true, message: '请输入字段名称' }]}>
            <Input placeholder="例如：处理优先级备注" />
          </Form.Item>
          <Form.Item
            name="key"
            label="字段标识"
            rules={[
              { required: true, message: '请输入字段标识' },
              { pattern: /^[a-z][a-z0-9_]*$/, message: '仅支持小写字母、数字和下划线，且需字母开头' },
            ]}
          >
            <Input disabled={Boolean(editingField)} placeholder="例如：priority_note" />
          </Form.Item>
          <Form.Item name="type" label="字段类型" rules={[{ required: true, message: '请选择字段类型' }]}>
            <Select options={FIELD_TYPE_OPTIONS} />
          </Form.Item>
          <Form.Item
            name="resource"
            label="适用资源"
            rules={[{ required: true, message: '请选择适用资源' }]}
          >
            <Select disabled={Boolean(editingField)} options={RESOURCE_OPTIONS} />
          </Form.Item>
          <Form.Item name="required" label="必填" valuePropName="checked">
            <Switch checkedChildren="是" unCheckedChildren="否" />
          </Form.Item>
          <Form.Item name="active" label="启用" valuePropName="checked">
            <Switch checkedChildren="启用" unCheckedChildren="停用" />
          </Form.Item>
          <Form.Item
            noStyle
            shouldUpdate={(prev, next) => prev.type !== next.type}
          >
            {({ getFieldValue }) =>
              ['select', 'multiselect'].includes(getFieldValue('type')) ? (
                <Form.Item
                  name="options"
                  label="选项"
                  rules={[{ required: true, message: '请选择输入选项' }]}
                >
                  <Input placeholder="使用逗号分隔，例如：高, 中, 低" />
                </Form.Item>
              ) : null
            }
          </Form.Item>
          <Form.Item name="validation" label="校验规则 JSON">
            <Input.TextArea rows={4} placeholder='例如：{"min":1,"max":10}' />
          </Form.Item>
          <Form.Item name="show_when" label="显示条件 JSON">
            <Input.TextArea rows={4} placeholder='例如：{"status":"open"}' />
          </Form.Item>
        </Form>
      </Modal>
    </>
  );
};

export default CustomFieldPage;
