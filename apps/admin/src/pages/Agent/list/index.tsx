import React from 'react';
import { ProTable } from '@ant-design/pro-components';
import type { ProColumns } from '@ant-design/pro-components';
import { Tag, Button, Space, Progress } from 'antd';
import { PlusOutlined } from '@ant-design/icons';
import { navigateTo } from '@/lib/navigation';
import { AGENT_STATUS_MAP } from '@/utils/constants';
import { listAgents } from '@/services/agent';

const AgentListPage: React.FC = () => {
  const columns: ProColumns<API.Agent>[] = [
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
      title: '状态',
      dataIndex: 'status',
      width: 100,
      valueType: 'select',
      valueEnum: Object.fromEntries(
        Object.entries(AGENT_STATUS_MAP).map(([k, v]) => [k, { text: v.text }]),
      ),
      render: (_, record) => {
        const item = AGENT_STATUS_MAP[record.status];
        return item ? <Tag color={item.color}>{item.text}</Tag> : record.status;
      },
    },
    {
      title: '技能',
      dataIndex: 'skills',
      search: false,
      render: (_, record) =>
        record.skills?.map((skill) => <Tag key={skill}>{skill}</Tag>),
    },
    {
      title: '会话数',
      dataIndex: 'current_sessions',
      width: 160,
      search: false,
      render: (_, record) => {
        const percent = record.max_sessions
          ? Math.round(((record.current_sessions || 0) / record.max_sessions) * 100)
          : 0;
        return (
          <Space>
            <span>
              {record.current_sessions || 0}/{record.max_sessions}
            </span>
            <Progress percent={percent} size="small" style={{ width: 60 }} />
          </Space>
        );
      },
    },
    {
      title: '操作',
      valueType: 'option',
      width: 120,
      render: (_, record) => (
        <Space>
          <a onClick={() => navigateTo(`/agent/detail/${record.id}`)}>查看</a>
          <a>编辑</a>
        </Space>
      ),
    },
  ];

  return (
    <ProTable<API.Agent>
      headerTitle="客服列表"
      rowKey="id"
      columns={columns}
      toolBarRender={() => [
        <Button key="create" type="primary" icon={<PlusOutlined />}>
          添加客服
        </Button>,
      ]}
      request={async (params) => {
        try {
          const result = await listAgents({
            page: params.current,
            page_size: params.pageSize,
            status: params.status,
            search: params.name,
          });
          return {
            data: result?.data || [],
            total: result?.total || 0,
            success: true,
          };
        } catch (error) {
          console.error('获取客服列表失败:', error);
          return { data: [], total: 0, success: true };
        }
      }}
      pagination={{ defaultPageSize: 20 }}
    />
  );
};

export default AgentListPage;
