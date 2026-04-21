import React from 'react';
import { ProTable } from '@ant-design/pro-components';
import type { ProColumns } from '@ant-design/pro-components';
import { Space, Tag } from 'antd';
import { listIntegrations } from '@/services/appMarket';

const AppMarketPage: React.FC = () => {
  const columns: ProColumns<API.Integration>[] = [
    {
      title: 'ID',
      dataIndex: 'id',
      width: 80,
    },
    {
      title: '应用名称',
      dataIndex: 'name',
      search: true,
    },
    {
      title: '分类',
      dataIndex: 'category',
      width: 140,
    },
    {
      title: '厂商',
      dataIndex: 'vendor',
      width: 160,
      search: false,
    },
    {
      title: '状态',
      dataIndex: 'enabled',
      width: 100,
      render: (_, record) => (
        <Tag color={record.enabled ? 'green' : 'default'}>
          {record.enabled ? '已启用' : '未启用'}
        </Tag>
      ),
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      valueType: 'dateTime',
      width: 180,
    },
    {
      title: '操作',
      valueType: 'option',
      width: 100,
      render: (_, record) => (
        <Space>
          {record.enabled ? <a>配置</a> : <a>启用</a>}
        </Space>
      ),
    },
  ];

  return (
    <ProTable<API.Integration>
      headerTitle="应用市场"
      rowKey="id"
      columns={columns}
      request={async (params) => {
        try {
          const result = await listIntegrations({
            page: params.current,
            page_size: params.pageSize,
            search: typeof params.name === 'string' ? params.name : undefined,
            category: typeof params.category === 'string' ? params.category : undefined,
          });
          return {
            data: result.data || [],
            total: result.total || 0,
            success: true,
          };
        } catch (error) {
          console.error('获取集成列表失败:', error);
          return { data: [], total: 0, success: true };
        }
      }}
      pagination={{ defaultPageSize: 20 }}
    />
  );
};

export default AppMarketPage;
