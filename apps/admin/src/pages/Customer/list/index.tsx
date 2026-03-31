import React, { useRef, useState } from 'react';
import { ProTable } from '@ant-design/pro-components';
import type { ProColumns, ActionType } from '@ant-design/pro-components';
import { Tag, Button, Space, Modal, Form, Input, message } from 'antd';
import { PlusOutlined } from '@ant-design/icons';
import { navigateTo } from '@/lib/navigation';
import { CUSTOMER_SOURCE_MAP } from '@/utils/constants';
import { listCustomers, createCustomer } from '@/services/customer';

const CustomerListPage: React.FC = () => {
  const actionRef = useRef<ActionType>();
  const [createOpen, setCreateOpen] = useState(false);
  const [creating, setCreating] = useState(false);
  const [form] = Form.useForm();

  const handleCreate = async () => {
    try {
      const values = await form.validateFields();
      setCreating(true);
      const tags = values.tags
        ? values.tags.split(',').map((s: string) => s.trim()).filter(Boolean)
        : [];
      await createCustomer({
        name: values.name,
        email: values.email,
        phone: values.phone,
        company: values.company,
        tags,
      });
      message.success('客户创建成功');
      setCreateOpen(false);
      form.resetFields();
      actionRef.current?.reload();
    } catch (error: any) {
      if (error?.errorFields) return;
      message.error('创建失败: ' + (error?.message || '未知错误'));
    } finally {
      setCreating(false);
    }
  };

  const columns: ProColumns<API.Customer>[] = [
    {
      title: '姓名',
      dataIndex: 'name',
      search: true,
    },
    {
      title: '邮箱',
      dataIndex: 'email',
      copyable: true,
    },
    {
      title: '电话',
      dataIndex: 'phone',
      copyable: true,
    },
    {
      title: '公司',
      dataIndex: 'company',
    },
    {
      title: '来源',
      dataIndex: 'source',
      width: 100,
      valueType: 'select',
      valueEnum: Object.fromEntries(
        Object.entries(CUSTOMER_SOURCE_MAP).map(([k, v]) => [k, { text: v }]),
      ),
      render: (_, record) => {
        const text = record.source ? CUSTOMER_SOURCE_MAP[record.source] : undefined;
        return text ? <Tag>{text}</Tag> : record.source;
      },
    },
    {
      title: '标签',
      dataIndex: 'tags',
      search: false,
      render: (_, record) =>
        record.tags?.map((tag) => <Tag key={tag}>{tag}</Tag>),
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      valueType: 'dateTime',
      width: 180,
      sorter: true,
    },
    {
      title: '操作',
      valueType: 'option',
      width: 120,
      render: (_, record) => (
        <Space>
          <a onClick={() => navigateTo(`/customer/detail/${record.id}`)}>查看</a>
          <a onClick={() => navigateTo(`/customer/detail/${record.id}`)}>编辑</a>
        </Space>
      ),
    },
  ];

  return (
    <>
      <ProTable<API.Customer>
        headerTitle="客户列表"
        rowKey="id"
        actionRef={actionRef}
        columns={columns}
        toolBarRender={() => [
          <Button
            key="create"
            type="primary"
            icon={<PlusOutlined />}
            onClick={() => setCreateOpen(true)}
          >
            新建客户
          </Button>,
        ]}
        request={async (params) => {
          try {
            const result = await listCustomers({
              page: params.current,
              page_size: params.pageSize,
              search: params.name,
              source: params.source,
            });
            return {
              data: result?.data || [],
              total: result?.total || 0,
              success: true,
            };
          } catch (error) {
            console.error('获取客户列表失败:', error);
            return { data: [], total: 0, success: true };
          }
        }}
        pagination={{ defaultPageSize: 20 }}
      />

      <Modal
        title="新建客户"
        open={createOpen}
        onCancel={() => {
          setCreateOpen(false);
          form.resetFields();
        }}
        onOk={handleCreate}
        confirmLoading={creating}
        okText="创建"
      >
        <Form form={form} layout="vertical" style={{ marginTop: 16 }}>
          <Form.Item
            name="name"
            label="姓名"
            rules={[{ required: true, message: '请输入姓名' }]}
          >
            <Input placeholder="客户姓名" />
          </Form.Item>
          <Form.Item name="email" label="邮箱" rules={[{ type: 'email', message: '请输入有效邮箱' }]}>
            <Input placeholder="邮箱" />
          </Form.Item>
          <Form.Item name="phone" label="电话">
            <Input placeholder="电话号码" />
          </Form.Item>
          <Form.Item name="company" label="公司">
            <Input placeholder="公司名称" />
          </Form.Item>
          <Form.Item name="tags" label="标签（逗号分隔）">
            <Input placeholder="如: VIP, 企业, 渠道A" />
          </Form.Item>
        </Form>
      </Modal>
    </>
  );
};

export default CustomerListPage;
