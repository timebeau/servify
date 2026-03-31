import React, { useState, useEffect } from 'react';
import { ProDescriptions } from '@ant-design/pro-components';
import { ProCard } from '@ant-design/pro-components';
import { Tag, Button, Space, Spin, List, message, Modal, Form, Input } from 'antd';
import { goBack, useDetailParams } from '@/lib/navigation';
import { CUSTOMER_SOURCE_MAP } from '@/utils/constants';
import { getCustomer, getCustomerActivity, updateCustomer } from '@/services/customer';

const CustomerDetailPage: React.FC = () => {
  const { id } = useDetailParams();
  const [loading, setLoading] = useState(true);
  const [customer, setCustomer] = useState<API.Customer | null>(null);
  const [activities, setActivities] = useState<any[]>([]);
  const [editOpen, setEditOpen] = useState(false);
  const [saving, setSaving] = useState(false);
  const [form] = Form.useForm();

  useEffect(() => {
    const fetchCustomer = async () => {
      if (!id) return;
      setLoading(true);
      try {
        const result = await getCustomer(Number(id));
        if (result) {
          setCustomer(result);
        }
      } catch (error) {
        console.error('获取客户详情失败:', error);
      } finally {
        setLoading(false);
      }
    };
    fetchCustomer();
  }, [id]);

  useEffect(() => {
    const fetchActivity = async () => {
      if (!id) return;
      try {
        const result = await getCustomerActivity(Number(id));
        setActivities(Array.isArray(result) ? result : result?.data || []);
      } catch (error) {
        console.error('获取活动记录失败:', error);
      }
    };
    fetchActivity();
  }, [id]);

  if (loading) {
    return (
      <div style={{ textAlign: 'center', padding: 80 }}>
        <Spin size="large" tip="加载中..." />
      </div>
    );
  }

  if (!customer) {
    return (
      <div style={{ textAlign: 'center', padding: 80, color: '#999' }}>
        客户不存在或加载失败
        <br />
        <Button style={{ marginTop: 16 }} onClick={goBack}>
          返回
        </Button>
      </div>
    );
  }

  const sourceText = customer.source ? CUSTOMER_SOURCE_MAP[customer.source] : undefined;

  const handleEdit = async () => {
    try {
      const values = await form.validateFields();
      setSaving(true);
      await updateCustomer(Number(id), values);
      message.success('客户信息已更新');
      setEditOpen(false);
      const result = await getCustomer(Number(id));
      if (result) setCustomer(result);
    } catch (error: any) {
      if (error?.errorFields) return;
      message.error('更新失败');
    } finally {
      setSaving(false);
    }
  };

  const openEditModal = () => {
    form.setFieldsValue({
      name: customer?.name,
      email: customer?.email,
      phone: customer?.phone,
      company: customer?.company,
    });
    setEditOpen(true);
  };

  return (
    <div>
      <ProCard
        title="客户详情"
        extra={
          <Space>
            <Button onClick={goBack}>返回</Button>
            <Button onClick={openEditModal}>编辑</Button>
            <Button type="primary">新建工单</Button>
          </Space>
        }
      >
        <ProDescriptions
          column={2}
          dataSource={customer}
          columns={[
            { title: '客户ID', dataIndex: 'id' },
            { title: '姓名', dataIndex: 'name' },
            { title: '邮箱', dataIndex: 'email', copyable: true },
            { title: '电话', dataIndex: 'phone', copyable: true },
            { title: '公司', dataIndex: 'company' },
            {
              title: '来源',
              dataIndex: 'source',
              render: () =>
                sourceText ? <Tag>{sourceText}</Tag> : customer.source,
            },
            { title: '创建时间', dataIndex: 'created_at' },
          ]}
        />
      </ProCard>

      <ProCard title="标签" style={{ marginTop: 16 }}>
        {customer.tags && customer.tags.length > 0 ? (
          <Space wrap>
            {customer.tags.map((tag) => (
              <Tag key={tag} color="blue">
                {tag}
              </Tag>
            ))}
          </Space>
        ) : (
          <div style={{ textAlign: 'center', padding: 20, color: '#999' }}>
            暂无标签
          </div>
        )}
      </ProCard>

      <ProCard title="活动记录" style={{ marginTop: 16 }}>
        <List
          dataSource={activities}
          locale={{ emptyText: '暂无活动记录' }}
          renderItem={(item: any) => (
            <List.Item>
              <List.Item.Meta
                title={item.action || item.description || '-'}
                description={item.created_at || '-'}
              />
            </List.Item>
          )}
        />
      </ProCard>

      <Modal
        title="编辑客户"
        open={editOpen}
        onCancel={() => setEditOpen(false)}
        onOk={handleEdit}
        confirmLoading={saving}
        okText="保存"
      >
        <Form form={form} layout="vertical" style={{ marginTop: 16 }}>
          <Form.Item name="name" label="姓名" rules={[{ required: true, message: '请输入姓名' }]}>
            <Input placeholder="客户姓名" />
          </Form.Item>
          <Form.Item name="email" label="邮箱">
            <Input placeholder="邮箱" />
          </Form.Item>
          <Form.Item name="phone" label="电话">
            <Input placeholder="电话号码" />
          </Form.Item>
          <Form.Item name="company" label="公司">
            <Input placeholder="公司名称" />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};

export default CustomerDetailPage;
