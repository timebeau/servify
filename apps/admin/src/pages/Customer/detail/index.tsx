import React, { useState, useEffect } from 'react';
import { ProDescriptions } from '@ant-design/pro-components';
import { ProCard } from '@ant-design/pro-components';
import { Tag, Button, Space, Spin, List, message } from 'antd';
import { useParams, history } from '@umijs/max';
import { CUSTOMER_SOURCE_MAP } from '@/utils/constants';
import { getCustomer, getCustomerActivity } from '@/services/customer';

const CustomerDetailPage: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const [loading, setLoading] = useState(true);
  const [customer, setCustomer] = useState<API.Customer | null>(null);
  const [activities, setActivities] = useState<any[]>([]);

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
        <Button style={{ marginTop: 16 }} onClick={() => history.back()}>
          返回
        </Button>
      </div>
    );
  }

  const sourceText = CUSTOMER_SOURCE_MAP[customer.source];

  return (
    <div>
      <ProCard
        title="客户详情"
        extra={
          <Space>
            <Button onClick={() => history.back()}>返回</Button>
            <Button>编辑</Button>
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
    </div>
  );
};

export default CustomerDetailPage;
