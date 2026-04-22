import React, { useRef, useState } from 'react';
import { ProTable } from '@ant-design/pro-components';
import type { ProColumns, ActionType } from '@ant-design/pro-components';
import { Button, Space, Tag, message, Modal, Form, Input, Switch } from 'antd';
import { PlusOutlined } from '@ant-design/icons';
import { listMacros, deleteMacro, createMacro, updateMacro } from '@/services/macro';
import { getErrorMessage, isFormValidationError } from '@/utils/error';

const MacroPage: React.FC = () => {
  const actionRef = useRef<ActionType>();
  const [modalOpen, setModalOpen] = useState(false);
  const [editingMacro, setEditingMacro] = useState<API.Macro | null>(null);
  const [submitting, setSubmitting] = useState(false);
  const [form] = Form.useForm();

  const openCreate = () => {
    setEditingMacro(null);
    form.resetFields();
    form.setFieldsValue({ language: 'zh' });
    setModalOpen(true);
  };

  const openEdit = (record: API.Macro) => {
    setEditingMacro(record);
    form.setFieldsValue({
      name: record.name,
      description: record.description,
      language: record.language || 'zh',
      content: record.content,
      active: record.active ?? true,
    });
    setModalOpen(true);
  };

  const openView = (record: API.Macro) => {
    Modal.info({
      title: record.name,
      width: 720,
      content: (
        <div style={{ marginTop: 16 }}>
          <p>语言：{record.language || '-'}</p>
          <p>状态：{record.active ? '启用' : '停用'}</p>
          <p>描述：{record.description || '-'}</p>
          <div style={{ whiteSpace: 'pre-wrap', lineHeight: 1.8 }}>{record.content}</div>
        </div>
      ),
    });
  };

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields();
      setSubmitting(true);
      if (editingMacro) {
        await updateMacro(editingMacro.id, {
          description: values.description,
          language: values.language,
          content: values.content,
          active: values.active,
        });
        message.success('宏已更新');
      } else {
        await createMacro({
          name: values.name,
          description: values.description,
          language: values.language,
          content: values.content,
        });
        message.success('宏已创建');
      }
      setModalOpen(false);
      form.resetFields();
      actionRef.current?.reload();
    } catch (error: unknown) {
      if (isFormValidationError(error)) return;
      message.error(getErrorMessage(error, editingMacro ? '更新失败' : '创建失败'));
    } finally {
      setSubmitting(false);
    }
  };
  const columns: ProColumns<API.Macro>[] = [
    {
      title: 'ID',
      dataIndex: 'id',
      width: 80,
    },
    {
      title: '名称',
      dataIndex: 'name',
      search: true,
    },
    {
      title: '描述',
      dataIndex: 'description',
      ellipsis: true,
      search: false,
    },
    {
      title: '语言',
      dataIndex: 'language',
      width: 100,
      search: false,
    },
    {
      title: '状态',
      dataIndex: 'active',
      width: 100,
      search: false,
      render: (_, record) => (
        <Tag color={record.active ? 'green' : 'default'}>
          {record.active ? '启用' : '停用'}
        </Tag>
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
      width: 160,
      render: (_, record) => (
        <Space>
          <a onClick={() => openView(record)}>查看</a>
          <a onClick={() => openEdit(record)}>编辑</a>
          <a
            onClick={async () => {
              try {
                await deleteMacro(record.id);
                message.success('宏已删除');
                actionRef.current?.reload();
              } catch (error) {
                message.error('删除失败');
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
    <ProTable<API.Macro>
      headerTitle="宏模板"
      rowKey="id"
      actionRef={actionRef}
      columns={columns}
      toolBarRender={() => [
        <Button key="create" type="primary" icon={<PlusOutlined />} onClick={openCreate}>
          新建宏
        </Button>,
      ]}
      request={async (params) => {
        try {
          const result = await listMacros();
          const keyword = typeof params.name === 'string' ? params.name.trim().toLowerCase() : '';
          let data = result.data;
          if (keyword) {
            data = data.filter((item) => item.name.toLowerCase().includes(keyword));
          }
          const total = data.length;
          const current = params.current || 1;
          const pageSize = params.pageSize || 20;
          const pageData = data.slice((current - 1) * pageSize, current * pageSize);
          return {
            data: pageData,
            total,
            success: true,
          };
        } catch (error) {
          console.error('获取宏列表失败:', error);
          return { data: [], total: 0, success: true };
        }
      }}
      pagination={{ defaultPageSize: 20 }}
    >
      <Modal
        title={editingMacro ? '编辑宏' : '新建宏'}
        open={modalOpen}
        onCancel={() => {
          setModalOpen(false);
          form.resetFields();
        }}
        onOk={handleSubmit}
        confirmLoading={submitting}
        okText={editingMacro ? '保存' : '创建'}
        width={720}
      >
        <Form form={form} layout="vertical" style={{ marginTop: 16 }}>
          <Form.Item name="name" label="名称" rules={[{ required: true, message: '请输入名称' }]}>
            <Input disabled={Boolean(editingMacro)} placeholder="例如：标准回复宏" />
          </Form.Item>
          <Form.Item name="description" label="描述">
            <Input placeholder="用途说明" />
          </Form.Item>
          <Form.Item name="language" label="语言">
            <Input placeholder="例如：zh" />
          </Form.Item>
          {editingMacro ? (
            <Form.Item name="active" label="启用" valuePropName="checked">
              <Switch checkedChildren="启用" unCheckedChildren="停用" />
            </Form.Item>
          ) : null}
          <Form.Item name="content" label="内容" rules={[{ required: true, message: '请输入宏内容' }]}>
            <Input.TextArea rows={10} placeholder="宏文本内容" />
          </Form.Item>
        </Form>
      </Modal>
    </ProTable>
  );
};

export default MacroPage;
