import React, { useRef } from 'react';
import { ProTable } from '@ant-design/pro-components';
import type { ProColumns, ActionType } from '@ant-design/pro-components';
import { Tag, Button, Space, message } from 'antd';
import { getWaitingQueue, getTransferHistory, processQueue } from '@/services/sessionTransfer';

const RoutingPage: React.FC = () => {
  const queueActionRef = useRef<ActionType>();

  const queueColumns: ProColumns<API.TransferQueueRecord>[] = [
    { title: 'ID', dataIndex: 'id', width: 80 },
    { title: '客户', dataIndex: 'customer', width: 120 },
    { title: '渠道', dataIndex: 'channel', width: 100 },
    {
      title: '等待时间(秒)',
      dataIndex: 'wait_time',
      width: 120,
      sorter: true,
    },
    {
      title: '优先级',
      dataIndex: 'priority',
      width: 100,
      render: (_, record) => <Tag>{record.priority}</Tag>,
    },
    {
      title: '进入队列时间',
      dataIndex: 'created_at',
      valueType: 'dateTime',
      width: 180,
    },
    {
      title: '操作',
      valueType: 'option',
      width: 80,
      render: () => <a>分配</a>,
    },
  ];

  const transferColumns: ProColumns<API.SessionTransferRecord>[] = [
    { title: 'ID', dataIndex: 'id', width: 80 },
    { title: '工单ID', dataIndex: 'ticket_id', width: 100 },
    { title: '原客服', dataIndex: 'from_agent', width: 120 },
    { title: '目标客服', dataIndex: 'to_agent', width: 120 },
    { title: '原因', dataIndex: 'reason', ellipsis: true },
    {
      title: '转接时间',
      dataIndex: 'created_at',
      valueType: 'dateTime',
      width: 180,
    },
  ];

  return (
    <div>
      <ProTable<API.TransferQueueRecord>
        headerTitle="等待队列"
        rowKey="id"
        columns={queueColumns}
        actionRef={queueActionRef}
        toolBarRender={() => [
          <Button
            key="refresh"
            type="primary"
            onClick={async () => {
              try {
                await processQueue();
                message.success('队列已处理');
                queueActionRef.current?.reload();
              } catch (error) {
                message.error('处理队列失败');
              }
            }}
          >
            刷新队列
          </Button>,
        ]}
        request={async () => {
          try {
            const result = await getWaitingQueue();
            const data = result.data || [];
            return { data, total: result.count || data.length, success: true };
          } catch (error) {
            console.error('获取等待队列失败:', error);
            return { data: [], total: 0, success: true };
          }
        }}
        search={false}
        pagination={{ defaultPageSize: 10 }}
      />

      <ProTable<API.SessionTransferRecord>
        headerTitle="转接历史"
        rowKey="id"
        columns={transferColumns}
        style={{ marginTop: 16 }}
        request={async () => {
          try {
            const result = await getTransferHistory('');
            return { data: result, total: result.length, success: true };
          } catch (error) {
            console.error('获取转接历史失败:', error);
            return { data: [], total: 0, success: true };
          }
        }}
        search={{ filterType: 'light' }}
        pagination={{ defaultPageSize: 10 }}
      />
    </div>
  );
};

export default RoutingPage;
