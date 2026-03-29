import React, { useState, useEffect } from 'react';
import { ProTable, ProCard } from '@ant-design/pro-components';
import type { ProColumns } from '@ant-design/pro-components';
import { Button, Empty, Input, Space, Spin, Tag, message } from 'antd';
import { getConversationMessages, sendConversationMessage } from '@/services/conversation';
import { getWorkspaceOverview } from '@/services/workspace';

interface ConversationRecord {
  id: string;
  customer_name?: string;
  agent_name?: string;
  platform?: string;
  status: string;
  started_at: string;
}

const ConversationPage: React.FC = () => {
  const [selectedId, setSelectedId] = useState<string | null>(null);
  const [overview, setOverview] = useState<API.WorkspaceOverview | null>(null);
  const [loading, setLoading] = useState(true);
  const [messageLoading, setMessageLoading] = useState(false);
  const [sending, setSending] = useState(false);
  const [messages, setMessages] = useState<API.ConversationMessage[]>([]);
  const [draft, setDraft] = useState('');

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

  useEffect(() => {
    const fetchMessages = async () => {
      if (!selectedId) {
        setMessages([]);
        return;
      }
      setMessageLoading(true);
      try {
        const result = await getConversationMessages(selectedId);
        setMessages(result?.data || []);
      } catch (error) {
        console.error('获取会话消息失败:', error);
        message.error('获取会话消息失败');
      } finally {
        setMessageLoading(false);
      }
    };
    fetchMessages();
  }, [selectedId]);

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
      dataIndex: 'platform',
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
      dataIndex: 'started_at',
      valueType: 'dateTime',
      width: 180,
    },
  ];

  const sessions = overview?.recent_sessions || [];
  const selectedSession = sessions.find((item) => item.id === selectedId) || null;

  const handleSend = async () => {
    if (!selectedId || !draft.trim()) {
      return;
    }
    setSending(true);
    try {
      const payload = await sendConversationMessage(selectedId, { content: draft.trim() });
      const item = payload?.data as API.ConversationMessage | undefined;
      if (item) {
        setMessages((prev) => [...prev, item]);
      }
      setDraft('');
      message.success('消息已发送');
    } catch (error) {
      console.error('发送消息失败:', error);
      message.error('发送消息失败');
    } finally {
      setSending(false);
    }
  };

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
          extra={
            selectedSession ? (
              <Space size="middle">
                <Tag>{selectedSession.platform || 'unknown'}</Tag>
                <span>{selectedSession.customer_name || '未识别客户'}</span>
                <span>{selectedSession.agent_name || '待分配客服'}</span>
              </Space>
            ) : null
          }
          style={{ height: '100%' }}
        >
          <div style={{ display: 'flex', flexDirection: 'column', height: '100%' }}>
            <div style={{ flex: 1, overflowY: 'auto', padding: 12, background: '#fafafa', borderRadius: 8 }}>
              {!selectedId ? (
                <Empty description="请先从左侧选择一个会话。" style={{ marginTop: 96 }} />
              ) : messageLoading ? (
                <div style={{ textAlign: 'center', padding: 80 }}>
                  <Spin />
                </div>
              ) : messages.length === 0 ? (
                <Empty description="当前会话暂无消息。" style={{ marginTop: 96 }} />
              ) : (
                messages.map((item) => {
                  const isAgent = item.sender === 'agent';
                  const align = isAgent ? 'flex-end' : 'flex-start';
                  const background = isAgent ? '#1677ff' : '#fff';
                  const color = isAgent ? '#fff' : '#000';
                  return (
                    <div
                      key={item.id}
                      style={{ display: 'flex', justifyContent: align, marginBottom: 12 }}
                    >
                      <div
                        style={{
                          maxWidth: '72%',
                          background,
                          color,
                          border: isAgent ? 'none' : '1px solid #f0f0f0',
                          borderRadius: 12,
                          padding: '10px 12px',
                          boxShadow: '0 1px 2px rgba(0,0,0,0.04)',
                        }}
                      >
                        <div style={{ fontSize: 12, opacity: 0.75, marginBottom: 4 }}>
                          {item.sender} · {new Date(item.created_at).toLocaleString()}
                        </div>
                        <div style={{ whiteSpace: 'pre-wrap', wordBreak: 'break-word' }}>
                          {item.content}
                        </div>
                      </div>
                    </div>
                  );
                })
              )}
            </div>
            <Input.TextArea
              rows={3}
              placeholder="输入消息..."
              value={draft}
              onChange={(e) => setDraft(e.target.value)}
              style={{ marginTop: 8 }}
              disabled={!selectedId || sending}
            />
            <div style={{ marginTop: 8, display: 'flex', justifyContent: 'flex-end' }}>
              <Button type="primary" onClick={handleSend} loading={sending} disabled={!selectedId || !draft.trim()}>
                发送消息
              </Button>
            </div>
          </div>
        </ProCard>
      </div>
    </div>
  );
};

export default ConversationPage;
