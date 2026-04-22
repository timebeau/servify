import React from 'react';
import { ProTable } from '@ant-design/pro-components';
import type { ProColumns, ActionType } from '@ant-design/pro-components';
import { Button, Space, Tag, message } from 'antd';
import { PlusOutlined } from '@ant-design/icons';
import { useRef } from 'react';
import { listShifts, deleteShift } from '@/services/shift';
import { getErrorMessage } from '@/utils/error';

const ShiftPage: React.FC = () => {
  const actionRef = useRef<ActionType>();
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
      valueType: 'select',
      valueEnum: {
        morning: { text: '早班' },
        afternoon: { text: '午班' },
        evening: { text: '晚班' },
        night: { text: '夜班' },
      },
    },
    {
      title: '开始时间',
      dataIndex: 'start_time',
      width: 180,
      valueType: 'dateTime',
      search: false,
    },
    {
      title: '结束时间',
      dataIndex: 'end_time',
      width: 180,
      valueType: 'dateTime',
      search: false,
    },
    {
      title: '状态',
      dataIndex: 'status',
      width: 100,
      valueType: 'select',
      valueEnum: {
        scheduled: { text: '已排班' },
        active: { text: '进行中' },
        completed: { text: '已完成' },
        cancelled: { text: '已取消' },
      },
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
                actionRef.current?.reload();
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
      actionRef={actionRef}
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
            shift_type: typeof params.shift_type === 'string' ? [params.shift_type] : undefined,
            status: typeof params.status === 'string' ? [params.status] : undefined,
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
      search={{ filterType: 'light' }}
      pagination={{ defaultPageSize: 20 }}
    />
  );
};

export default ShiftPage;
