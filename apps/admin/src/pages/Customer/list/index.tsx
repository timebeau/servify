import React, { useRef } from 'react';
import { ProTable } from '@ant-design/pro-components';
import type { ProColumns, ActionType } from '@ant-design/pro-components';
import { Tag, Button, Space } from 'antd';
import { PlusOutlined } from '@ant-design/icons';
import { history } from '@umijs/max';
import { CUSTOMER_SOURCE_MAP } from '@/utils/constants';
import { listCustomers } from '@/services/customer';

const CustomerListPage: React.FC = () => {
  const actionRef = useRef<ActionType>();

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
        const text = CUSTOMER_SOURCE_MAP[record.source];
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
          <a onClick={() => history.push(`/customer/detail/${record.id}`)}>查看</a>
          <a>编辑</a>
        </Space>
      ),
    },
  ];

  return (
    <ProTable<API.Customer>
      headerTitle="客户列表"
      rowKey="id"
      actionRef={actionRef}
      columns={columns}
      toolBarRender={() => [
        <Button key="create" type="primary" icon={<PlusOutlined />}>
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
  );
};

export default CustomerListPage;
