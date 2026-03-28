import React from 'react';
import { ProTable } from '@ant-design/pro-components';
import type { ProColumns } from '@ant-design/pro-components';
import { Button, Space, Tag, message } from 'antd';
import { PlusOutlined } from '@ant-design/icons';
import { listShifts, deleteShift } from '@/services/shift';

const ShiftPage: React.FC = () => {
  const columns: ProColumns<API.Shift>[] = [
    {
      title: 'ID',
      dataIndex: 'id',
      width: 80,
    },
    {
      title: '客服',
      dataIndex: 'agent_name',
      width: 120,
      search: true,
    },
    {
      title: '开始时间',
      dataIndex: 'start_time',
      width: 180,
    },
    {
      title: '结束时间',
      dataIndex: 'end_time',
      width: 180,
    },
    {
      title: '状态',
      dataIndex: 'status',
      width: 80,
      render: (_, record) => (
        <Tag color={record.status === 'active' ? 'green' : 'default'}>
          {record.status === 'active' ? '进行中' : record.status}
        </Tag>
      ),
    },
    {
      title: '操作',
      valueType: 'option',
      width: 120,
      render: (_, record) => (
        <Space>
          <a>编辑</a>
          <a
            onClick={async () => {
              try {
                await deleteShift(record.id);
                message.success('班次已删除');
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
    <ProTable<API.Shift>
      headerTitle="班次管理"
      rowKey="id"
      columns={columns}
      toolBarRender={() => [
        <Button key="create" type="primary" icon={<PlusOutlined />}>
          新建班次
        </Button>,
      ]}
      request={async (params) => {
        try {
          const result = await listShifts({
            page: params.current,
            page_size: params.pageSize,
            agent_id: params.agent_name ? Number(params.agent_name) : undefined,
          });
          return {
            data: result?.data || [],
            total: result?.total || 0,
            success: true,
          };
        } catch (error) {
          console.error('获取班次列表失败:', error);
          return { data: [], total: 0, success: true };
        }
      }}
      pagination={{ defaultPageSize: 20 }}
    />
  );
};

export default ShiftPage;
