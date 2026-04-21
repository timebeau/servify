import React from 'react';
import { ProTable } from '@ant-design/pro-components';
import type { ProColumns } from '@ant-design/pro-components';
import { Button, Space, Tag, message } from 'antd';
import { PlusOutlined } from '@ant-design/icons';
import { listShifts, deleteShift } from '@/services/shift';
import { getErrorMessage } from '@/utils/error';

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
      width: 140,
      search: false,
      render: (_, record) => record.agent?.name || record.agent?.username || record.agent_name || `#${record.agent_id}`,
    },
    {
      title: '班次类型',
      dataIndex: 'shift_type',
      width: 120,
      search: false,
    },
    {
      title: '开始时间',
      dataIndex: 'start_time',
      width: 180,
      valueType: 'dateTime',
    },
    {
      title: '结束时间',
      dataIndex: 'end_time',
      width: 180,
      valueType: 'dateTime',
    },
    {
      title: '状态',
      dataIndex: 'status',
      width: 100,
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
              } catch (error: unknown) {
                message.error(getErrorMessage(error, '删除失败'));
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
          });
          return {
            data: result.data,
            total: result.total,
            success: true,
          };
        } catch (error: unknown) {
          console.error('获取班次列表失败:', error);
          return { data: [], total: 0, success: true };
        }
      }}
      pagination={{ defaultPageSize: 20 }}
    />
  );
};

export default ShiftPage;
