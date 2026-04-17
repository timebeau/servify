import React, { useRef, useState } from 'react';
import { ProTable } from '@ant-design/pro-components';
import type { ProColumns, ActionType } from '@ant-design/pro-components';
import { Tag, Button, Space, message, Modal, Form, Input, Select, Switch } from 'antd';
import { PlusOutlined } from '@ant-design/icons';
import { navigateTo } from '@/lib/navigation';
import { listDocs, deleteDoc, createDoc, updateDoc } from '@/services/knowledge';

const CATEGORIES = ['产品文档', '常见问题', '操作指南', 'API文档', '其他'];

const KnowledgeListPage: React.FC = () => {
  const actionRef = useRef<ActionType>();
  const [modalOpen, setModalOpen] = useState(false);
  const [modalType, setModalType] = useState<'create' | 'edit'>('create');
  const [editingDoc, setEditingDoc] = useState<API.KnowledgeDoc | null>(null);
  const [submitting, setSubmitting] = useState(false);
  const [form] = Form.useForm();

  const openCreate = () => {
    setModalType('create');
    setEditingDoc(null);
    form.resetFields();
    setModalOpen(true);
  };

  const openEdit = (record: API.KnowledgeDoc) => {
    setModalType('edit');
    setEditingDoc(record);
    form.setFieldsValue({
      title: record.title,
      category: record.category,
      content: record.content,
      tags: record.tags?.join(', '),
      is_public: record.is_public ?? false,
    });
    setModalOpen(true);
  };

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields();
      setSubmitting(true);
      const tags = values.tags
        ? values.tags.split(',').map((s: string) => s.trim()).filter(Boolean)
        : [];

      if (modalType === 'create') {
        await createDoc({
          title: values.title,
          content: values.content || '',
          category: values.category,
          tags,
          is_public: values.is_public ?? false,
        });
        message.success('文档创建成功');
      } else if (editingDoc) {
        await updateDoc(editingDoc.id, {
          title: values.title,
          content: values.content,
          category: values.category,
          tags,
          is_public: values.is_public ?? false,
        });
        message.success('文档已更新');
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
      content: '删除后不可恢复，确定要删除此文档吗？',
      okText: '确认删除',
      okButtonProps: { danger: true },
      onOk: async () => {
        try {
          await deleteDoc(id);
          message.success('文档已删除');
          actionRef.current?.reload();
        } catch (error) {
          message.error('删除失败');
        }
      },
    });
  };

  const columns: ProColumns<API.KnowledgeDoc>[] = [
    {
      title: 'ID',
      dataIndex: 'id',
      width: 80,
    },
    {
      title: '标题',
      dataIndex: 'title',
      search: true,
      ellipsis: true,
    },
    {
      title: '分类',
      dataIndex: 'category',
      width: 120,
      valueType: 'select',
      valueEnum: Object.fromEntries(CATEGORIES.map((c) => [c, { text: c }])),
    },
    {
      title: '状态',
      dataIndex: 'status',
      width: 100,
      render: (_, record) => <Tag>{record.status}</Tag>,
    },
    {
      title: '公开',
      dataIndex: 'is_public',
      width: 100,
      search: false,
      render: (_, record) => (
        <Tag color={record.is_public ? 'green' : 'default'}>
          {record.is_public ? '公开' : '内部'}
        </Tag>
      ),
    },
    {
      title: '标签',
      dataIndex: 'tags',
      search: false,
      render: (_, record) =>
        record.tags?.map((tag) => <Tag key={tag}>{tag}</Tag>),
    },
    {
      title: '更新时间',
      dataIndex: 'updated_at',
      valueType: 'dateTime',
      width: 180,
      sorter: true,
    },
    {
      title: '操作',
      valueType: 'option',
      width: 180,
      render: (_, record) => (
        <Space>
          <a onClick={() => navigateTo(`/knowledge/detail/${record.id}`)}>查看</a>
          <a onClick={() => openEdit(record)}>编辑</a>
          <a onClick={() => handleDelete(record.id)} style={{ color: '#ff4d4f' }}>删除</a>
        </Space>
      ),
    },
  ];

  return (
    <>
      <ProTable<API.KnowledgeDoc>
        headerTitle="知识库文档"
        rowKey="id"
        actionRef={actionRef}
        columns={columns}
        toolBarRender={() => [
          <Button key="create" type="primary" icon={<PlusOutlined />} onClick={openCreate}>
            新建文档
          </Button>,
        ]}
        request={async (params) => {
          try {
            const result = await listDocs({
              page: params.current,
              page_size: params.pageSize,
              search: params.title,
              category: params.category,
              status: params.status,
            });
            return {
              data: result?.data || [],
              total: result?.total || 0,
              success: true,
            };
          } catch (error) {
            console.error('获取知识库列表失败:', error);
            return { data: [], total: 0, success: true };
          }
        }}
        pagination={{ defaultPageSize: 20 }}
      />

      <Modal
        title={modalType === 'create' ? '新建文档' : '编辑文档'}
        open={modalOpen}
        onCancel={() => { setModalOpen(false); form.resetFields(); }}
        onOk={handleSubmit}
        confirmLoading={submitting}
        okText={modalType === 'create' ? '创建' : '保存'}
        width={640}
      >
        <Form form={form} layout="vertical" style={{ marginTop: 16 }}>
          <Form.Item name="title" label="标题" rules={[{ required: true, message: '请输入标题' }]}>
            <Input placeholder="文档标题" />
          </Form.Item>
          <Form.Item name="category" label="分类">
            <Select placeholder="选择分类" allowClear options={CATEGORIES.map((c) => ({ label: c, value: c }))} />
          </Form.Item>
          <Form.Item name="tags" label="标签（逗号分隔）">
            <Input placeholder="如: 入门, API, 常见问题" />
          </Form.Item>
          <Form.Item
            name="is_public"
            label="公开到公共知识库"
            valuePropName="checked"
            initialValue={false}
          >
            <Switch checkedChildren="公开" unCheckedChildren="内部" />
          </Form.Item>
          <Form.Item name="content" label="内容" rules={[{ required: true, message: '请输入内容' }]}>
            <Input.TextArea rows={8} placeholder="文档内容（支持 Markdown）" />
          </Form.Item>
        </Form>
      </Modal>
    </>
  );
};

export default KnowledgeListPage;
