import React, { useRef, useState } from 'react';
import { ProTable } from '@ant-design/pro-components';
import type { ActionType, ProColumns } from '@ant-design/pro-components';
import { Button, Input, Tag, message } from 'antd';
import { getTransferHistory, getWaitingQueue, processQueue } from '@/services/sessionTransfer';

const RoutingPage: React.FC = () => {
  const queueActionRef = useRef<ActionType>();
  const historyActionRef = useRef<ActionType>();
  const [sessionId, setSessionId] = useState('');

  const queueColumns: ProColumns<API.TransferQueueRecord>[] = [
    { title: '会话ID', dataIndex: 'session_id', width: 180 },
    { title: '原因', dataIndex: 'reason', ellipsis: true },
    {
      title: '目标技能',
      dataIndex: 'target_skills',
      ellipsis: true,
      render: (_, record) => record.target_skills || '-',
    },
    {
      title: '优先级',
      dataIndex: 'priority',
      width: 120,
      render: (_, record) => <Tag>{record.priority || 'normal'}</Tag>,
    },
    {
      title: '状态',
      dataIndex: 'status',
      width: 120,
      render: (_, record) => <Tag>{record.status || 'waiting'}</Tag>,
    },
    {
      title: '入队时间',
      dataIndex: 'queued_at',
      valueType: 'dateTime',
      width: 180,
    },
    {
      title: '分配时间',
      dataIndex: 'assigned_at',
      valueType: 'dateTime',
      width: 180,
    },
  ];

  const transferColumns: ProColumns<API.SessionTransferRecord>[] = [
    { title: 'ID', dataIndex: 'id', width: 80 },
    { title: '会话ID', dataIndex: 'session_id', width: 180 },
    { title: '原客服ID', dataIndex: 'from_agent_id', width: 120 },
    { title: '目标客服ID', dataIndex: 'to_agent_id', width: 120 },
    { title: '原因', dataIndex: 'reason', ellipsis: true },
    { title: '备注', dataIndex: 'notes', ellipsis: true },
    {
      title: '转接时间',
      dataIndex: 'transferred_at',
      valueType: 'dateTime',
      width: 180,
    },
  ];

  return (
    <div>
      <ProTable<API.TransferQueueRecord>
        headerTitle="等待队列"
        rowKey="session_id"
        columns={queueColumns}
        actionRef={queueActionRef}
        toolBarRender={() => [
          <Button
            key="process"
            type="primary"
            onClick={async () => {
              try {
                await processQueue();
                message.success('队列处理完成');
                queueActionRef.current?.reload();
                historyActionRef.current?.reload();
              } catch (error) {
                message.error('处理队列失败');
              }
            }}
          >
            处理队列
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
        actionRef={historyActionRef}
        style={{ marginTop: 16 }}
        toolBarRender={() => [
          <Input
            key="session-id"
            allowClear
            placeholder="按会话ID筛选"
            style={{ width: 240 }}
            value={sessionId}
            onChange={(event) => {
              setSessionId(event.target.value);
            }}
            onPressEnter={() => historyActionRef.current?.reload()}
          />,
          <Button key="search" type="primary" onClick={() => historyActionRef.current?.reload()}>
            查询
          </Button>,
        ]}
        request={async () => {
          try {
            const result = await getTransferHistory(sessionId, 50);
            return { data: result.items, total: result.count, success: true };
          } catch (error) {
            console.error('获取转接历史失败:', error);
            return { data: [], total: 0, success: true };
          }
        }}
        search={false}
        pagination={{ defaultPageSize: 10 }}
      />
    </div>
  );
};

export default RoutingPage;
