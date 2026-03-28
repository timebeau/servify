import React, { useState, useEffect } from 'react';
import { ProTable, ProCard } from '@ant-design/pro-components';
import type { ProColumns } from '@ant-design/pro-components';
import { Tag, Input, Spin } from 'antd';
import { getWorkspaceOverview } from '@/services/workspace';

interface ConversationRecord {
  id: string;
  customer_name?: string;
  agent_name?: string;
  channel?: string;
  status: string;
  created_at: string;
  updated_at: string;
}

const ConversationPage: React.FC = () => {
  const [selectedId, setSelectedId] = useState<string | null>(null);
  const [overview, setOverview] = useState<any>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchOverview = async () => {
      try {
        const result = await getWorkspaceOverview();
        if (result) {
          setOverview(result);
        }
      } catch (error) {
        console.error('获取工作区概览失败:', error);
      } finally {
        setLoading(false);
      }
    };
    fetchOverview();
  }, []);

  const columns: ProColumns<ConversationRecord>[] = [
    {
      title: 'ID',
      dataIndex: 'id',
      width: 80,
    },
    {
      title: '客户',
      dataIndex: 'customer_name',
      width: 120,
      search: true,
    },
    {
      title: '客服',
      dataIndex: 'agent_name',
      width: 120,
    },
    {
      title: '渠道',
      dataIndex: 'channel',
      width: 100,
    },
    {
      title: '状态',
      dataIndex: 'status',
      width: 100,
      render: (_, record) => <Tag>{record.status}</Tag>,
    },
    {
      title: '开始时间',
      dataIndex: 'created_at',
      valueType: 'dateTime',
      width: 180,
    },
  ];

  return (
    <div style={{ display: 'flex', gap: 16, height: 'calc(100vh - 120px)' }}>
      <div style={{ width: 480, flexShrink: 0 }}>
        <ProTable<ConversationRecord>
          headerTitle="会话列表"
          rowKey="id"
          columns={columns}
          search={{ filterType: 'light' }}
          tableAlertRender={false}
          scroll={{ y: 'calc(100vh - 220px)' }}
          pagination={{ defaultPageSize: 20 }}
          onRow={(record) => ({
            onClick: () => setSelectedId(record.id),
            style: {
              cursor: 'pointer',
              background: selectedId === record.id ? '#e6f7ff' : undefined,
            },
          })}
          request={async () => {
            try {
              const sessions = overview?.sessions || [];
              return { data: sessions, total: sessions.length, success: true };
            } catch (error) {
              console.error('获取会话列表失败:', error);
              return { data: [], total: 0, success: true };
            }
          }}
        />
      </div>
      <div style={{ flex: 1 }}>
        <ProCard
          title={selectedId ? `会话 ${selectedId}` : '聊天区域'}
          style={{ height: '100%' }}
        >
          <div style={{ display: 'flex', flexDirection: 'column', height: '100%' }}>
            <div style={{ flex: 1, textAlign: 'center', padding: 40, color: '#999' }}>
              聊天消息区域占位 - 待接入数据
            </div>
            <Input.TextArea
              rows={3}
              placeholder="输入消息..."
              style={{ marginTop: 8 }}
              disabled
            />
          </div>
        </ProCard>
      </div>
    </div>
  );
};

export default ConversationPage;
