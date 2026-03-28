import React from 'react';
import { ProTable } from '@ant-design/pro-components';
import type { ProColumns } from '@ant-design/pro-components';
import { Tag, Button, Space, message } from 'antd';
import { revokeUserTokens } from '@/services/security';

interface UserSecurityRecord {
  id: number;
  username: string;
  email: string;
  role: string;
  last_login: string;
  is_active: boolean;
}

interface TokenRecord {
  id: string;
  user: string;
  token_name: string;
  created_at: string;
  expires_at: string;
  last_used: string;
}

const SecurityPage: React.FC = () => {
  const userColumns: ProColumns<UserSecurityRecord>[] = [
    {
      title: 'ID',
      dataIndex: 'id',
      width: 80,
    },
    {
      title: '用户名',
      dataIndex: 'username',
      search: true,
    },
    {
      title: '邮箱',
      dataIndex: 'email',
    },
    {
      title: '角色',
      dataIndex: 'role',
      width: 100,
      render: (_, record) => <Tag color="blue">{record.role}</Tag>,
    },
    {
      title: '最后登录',
      dataIndex: 'last_login',
      valueType: 'dateTime',
      width: 180,
    },
    {
      title: '状态',
      dataIndex: 'is_active',
      width: 80,
      render: (_, record) => (
        <Tag color={record.is_active ? 'green' : 'red'}>
          {record.is_active ? '活跃' : '禁用'}
        </Tag>
      ),
    },
    {
      title: '操作',
      valueType: 'option',
      width: 160,
      render: (_, record) => (
        <Space>
          <a>重置密码</a>
          <a
            onClick={async () => {
              try {
                await revokeUserTokens(record.id);
                message.success('用户 Token 已撤销');
              } catch (error) {
                message.error('操作失败');
              }
            }}
          >
            撤销Token
          </a>
          <a>禁用</a>
        </Space>
      ),
    },
  ];

  const tokenColumns: ProColumns<TokenRecord>[] = [
    { title: 'ID', dataIndex: 'id', width: 80 },
    { title: '用户', dataIndex: 'user', width: 120 },
    { title: 'Token 名称', dataIndex: 'token_name' },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      valueType: 'dateTime',
      width: 180,
    },
    {
      title: '过期时间',
      dataIndex: 'expires_at',
      valueType: 'dateTime',
      width: 180,
    },
    {
      title: '最后使用',
      dataIndex: 'last_used',
      valueType: 'dateTime',
      width: 180,
    },
    {
      title: '操作',
      valueType: 'option',
      width: 80,
      render: () => (
        <Button
          type="link"
          danger
          onClick={() => message.info('Token 撤销功能待对接')}
        >
          撤销
        </Button>
      ),
    },
  ];

  return (
    <div>
      <ProTable<UserSecurityRecord>
        headerTitle="用户安全管理"
        rowKey="id"
        columns={userColumns}
        request={async () => {
          return { data: [], total: 0, success: true };
        }}
        pagination={{ defaultPageSize: 20 }}
      />

      <ProTable<TokenRecord>
        headerTitle="Token 管理"
        rowKey="id"
        columns={tokenColumns}
        style={{ marginTop: 16 }}
        toolBarRender={() => [
          <Button key="revoke-all" danger>
            批量撤销
          </Button>,
        ]}
        request={async () => {
          return { data: [], total: 0, success: true };
        }}
        search={{ filterType: 'light' }}
        pagination={{ defaultPageSize: 20 }}
      />
    </div>
  );
};

export default SecurityPage;
