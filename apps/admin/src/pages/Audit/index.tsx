import React from 'react';
import { ProTable } from '@ant-design/pro-components';
import type { ProColumns } from '@ant-design/pro-components';
import { Tag } from 'antd';
import { listAuditLogs } from '@/services/audit';

const AuditPage: React.FC = () => {
  const columns: ProColumns<API.AuditLog>[] = [
    {
      title: 'ID',
      dataIndex: 'id',
      width: 80,
    },
    {
      title: '操作',
      dataIndex: 'action',
      width: 100,
      search: true,
      render: (_, record) => <Tag>{record.action}</Tag>,
    },
    {
      title: '资源类型',
      dataIndex: 'resource_type',
      width: 120,
      search: true,
    },
    {
      title: '资源ID',
      dataIndex: 'resource_id',
      width: 120,
      copyable: true,
    },
    {
      title: '操作者类型',
      dataIndex: 'principal_kind',
      width: 120,
    },
    {
      title: '操作者ID',
      dataIndex: 'actor_user_id',
      width: 120,
    },
    {
      title: '成功',
      dataIndex: 'success',
      width: 80,
      render: (_, record) => (
        <Tag color={record.success ? 'green' : 'red'}>
          {record.success ? '成功' : '失败'}
        </Tag>
      ),
    },
    {
      title: '详情',
      dataIndex: 'request_json',
      ellipsis: true,
      search: false,
      render: (_, record) => record.request_json || '-',
    },
    {
      title: '时间',
      dataIndex: 'created_at',
      valueType: 'dateTime',
      width: 180,
      sorter: true,
    },
  ];

  return (
    <ProTable<API.AuditLog>
      headerTitle="审计日志"
      rowKey="id"
      columns={columns}
      request={async (params) => {
        try {
          const result = await listAuditLogs({
            page: params.current,
            page_size: params.pageSize,
            action: params.action,
            resource_type: params.resource_type,
          });
          return {
            data: result?.data || [],
            total: result?.total || 0,
            success: true,
          };
        } catch (error) {
          console.error('获取审计日志失败:', error);
          return { data: [], total: 0, success: true };
        }
      }}
      search={{ filterType: 'light' }}
      pagination={{ defaultPageSize: 20 }}
    />
  );
};

export default AuditPage;
