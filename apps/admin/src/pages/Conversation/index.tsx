import React, { useState, useEffect, useRef, useCallback } from 'react';
import { ProTable, ProCard } from '@ant-design/pro-components';
import type { ProColumns } from '@ant-design/pro-components';
import {
  Button,
  Empty,
  Input,
  Modal,
  Select,
  Space,
  Spin,
  Tag,
  message,
  Tooltip,
} from 'antd';
import {
  getConversationMessages,
  sendConversationMessage,
  assignAgent,
  transferConversation,
  closeConversation,
} from '@/services/conversation';
import { getWorkspaceOverview } from '@/services/workspace';

const STATUS_MAP: Record<string, { color: string; label: string }> = {
  active: { color: 'green', label: '进行中' },
  waiting_human: { color: 'orange', label: '等待客服' },
  transferred: { color: 'blue', label: '已转派' },
  closed: { color: 'default', label: '已结束' },
};

const SENDER_MAP: Record<string, { label: string; color: string }> = {
  agent: { label: '客服', color: '#1677ff' },
  customer: { label: '客户', color: '#52c41a' },
  ai: { label: 'AI', color: '#722ed1' },
  system: { label: '系统', color: '#999' },
};

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
  const [hasMore, setHasMore] = useState(false);
  const [transferModalOpen, setTransferModalOpen] = useState(false);
  const [transferTarget, setTransferTarget] = useState<number | null>(null);
  const [operating, setOperating] = useState(false);
  const chatEndRef = useRef<HTMLDivElement>(null);

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

  const scrollToBottom = useCallback(() => {
    setTimeout(() => chatEndRef.current?.scrollIntoView({ behavior: 'smooth' }), 50);
  }, []);

  useEffect(() => {
    const fetchMessages = async () => {
      if (!selectedId) {
        setMessages([]);
        setHasMore(false);
        return;
      }
      setMessageLoading(true);
      try {
        const result = await getConversationMessages(selectedId, { limit: 50 });
        setMessages(result?.data || []);
        setHasMore((result?.data || []).length >= 50);
        scrollToBottom();
      } catch (error) {
        console.error('获取会话消息失败:', error);
        message.error('获取会话消息失败');
      } finally {
        setMessageLoading(false);
      }
    };
    fetchMessages();
  }, [selectedId, scrollToBottom]);

  const handleLoadMore = async () => {
    if (!selectedId || messages.length === 0) return;
    const oldestId = messages[0]?.id;
    if (!oldestId) return;
    try {
      const result = await getConversationMessages(selectedId, { before: oldestId, limit: 50 });
      const older = result?.data || [];
      setMessages((prev) => [...older, ...prev]);
      setHasMore(older.length >= 50);
    } catch {
      message.error('加载更多消息失败');
    }
  };

  const handleSend = async () => {
    if (!selectedId || !draft.trim()) return;
    setSending(true);
    try {
      const payload = await sendConversationMessage(selectedId, { content: draft.trim() });
      const item = payload?.data as API.ConversationMessage | undefined;
      if (item) {
        setMessages((prev) => [...prev, item]);
      }
      setDraft('');
      message.success('消息已发送');
      scrollToBottom();
    } catch (error) {
      console.error('发送消息失败:', error);
      message.error('发送消息失败');
    } finally {
      setSending(false);
    }
  };

  const handleAssign = async () => {
    if (!selectedId) return;
    setOperating(true);
    try {
      // Auto-assign: pick first available agent from overview
      const agents = overview?.agent_stats?.available_agents || [];
      if (agents.length === 0) {
        message.warning('当前没有可用客服');
        return;
      }
      await assignAgent(selectedId, { agent_id: agents[0].id });
      message.success('已接管会话');
      // Refresh messages to get system event
      const result = await getConversationMessages(selectedId, { limit: 50 });
      setMessages(result?.data || []);
    } catch (error) {
      message.error('接管失败: ' + (error as Error).message);
    } finally {
      setOperating(false);
    }
  };

  const handleTransfer = async () => {
    if (!selectedId || !transferTarget) return;
    setOperating(true);
    try {
      await transferConversation(selectedId, { to_agent_id: transferTarget });
      message.success('已转派会话');
      setTransferModalOpen(false);
      setTransferTarget(null);
      const result = await getConversationMessages(selectedId, { limit: 50 });
      setMessages(result?.data || []);
    } catch (error) {
      message.error('转派失败: ' + (error as Error).message);
    } finally {
      setOperating(false);
    }
  };

  const handleClose = async () => {
    if (!selectedId) return;
    setOperating(true);
    try {
      await closeConversation(selectedId);
      message.success('会话已结束');
      const result = await getConversationMessages(selectedId, { limit: 50 });
      setMessages(result?.data || []);
    } catch (error) {
      message.error('关闭失败: ' + (error as Error).message);
    } finally {
      setOperating(false);
    }
  };

  const columns: ProColumns<ConversationRecord>[] = [
    { title: 'ID', dataIndex: 'id', width: 80, render: (_, r) => <Tooltip title={r.id}><span style={{ fontFamily: 'monospace', fontSize: 12 }}>{r.id.length > 8 ? `${r.id.slice(0, 8)}...` : r.id}</span></Tooltip> },
    { title: '客户', dataIndex: 'customer_name', width: 120, search: true },
    { title: '客服', dataIndex: 'agent_name', width: 120 },
    { title: '渠道', dataIndex: 'platform', width: 80 },
    { title: '状态', dataIndex: 'status', width: 100, render: (_, r) => { const c = STATUS_MAP[r.status] || { color: 'default', label: r.status }; return <Tag color={c.color}>{c.label}</Tag>; } },
    { title: '开始时间', dataIndex: 'started_at', valueType: 'dateTime', width: 160 },
  ];

  const sessions = overview?.recent_sessions || [];
  const selectedSession = sessions.find((s) => s.id === selectedId) || null;
  const sessionStatus = selectedSession?.status || '';
  const statusCfg = STATUS_MAP[sessionStatus] || {};
  const isClosed = sessionStatus === 'closed';
  const isWaiting = sessionStatus === 'waiting_human';
  const isActive = sessionStatus === 'active';

  const onlineAgents = (overview?.agent_stats?.available_agents || []).filter(
    (a: any) => a.id !== selectedSession?.agent_id,
  );

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
            style: { cursor: 'pointer', background: selectedId === record.id ? '#e6f7ff' : undefined },
          })}
          request={async () => ({ data: sessions, total: sessions.length, success: true })}
        />
      </div>
      <div style={{ flex: 1 }}>
        <ProCard
          title={
            selectedId ? (
              <Space>
                <span>会话 {selectedId.slice(0, 8)}</span>
                {sessionStatus && <Tag color={statusCfg.color}>{statusCfg.label || sessionStatus}</Tag>}
              </Space>
            ) : '聊天区域'
          }
          extra={
            selectedSession ? (
              <Space size="middle">
                <Tag>{selectedSession.platform || 'unknown'}</Tag>
                <span>{selectedSession.customer_name || '未识别客户'}</span>
                <span>{selectedSession.agent_name || '待分配客服'}</span>
                {!isClosed && (
                  <>
                    {isWaiting && (
                      <Button size="small" type="primary" onClick={handleAssign} loading={operating}>
                        接管
                      </Button>
                    )}
                    {isActive && (
                      <Button size="small" onClick={() => setTransferModalOpen(true)} loading={operating}>
                        转派
                      </Button>
                    )}
                    {!isClosed && (
                      <Button size="small" danger onClick={handleClose} loading={operating}>
                        结束
                      </Button>
                    )}
                  </>
                )}
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
                <div style={{ textAlign: 'center', padding: 80 }}><Spin /></div>
              ) : messages.length === 0 ? (
                <Empty description="当前会话暂无消息。" style={{ marginTop: 96 }} />
              ) : (
                <>
                  {hasMore && (
                    <div style={{ textAlign: 'center', marginBottom: 12 }}>
                      <Button size="small" type="link" onClick={handleLoadMore}>加载更多消息</Button>
                    </div>
                  )}
                  {messages.map((item) => {
                    const isAgent = item.sender === 'agent';
                    const align = isAgent ? 'flex-end' : 'flex-start';
                    const background = isAgent ? '#1677ff' : '#fff';
                    const color = isAgent ? '#fff' : '#000';
                    const senderCfg = SENDER_MAP[item.sender] || { label: item.sender, color: '#999' };
                    return (
                      <div key={item.id} style={{ display: 'flex', justifyContent: align, marginBottom: 12 }}>
                        <div style={{
                          maxWidth: '72%', background, color,
                          border: isAgent ? 'none' : '1px solid #f0f0f0',
                          borderRadius: 12, padding: '10px 12px',
                          boxShadow: '0 1px 2px rgba(0,0,0,0.04)',
                        }}>
                          <div style={{ fontSize: 12, opacity: 0.75, marginBottom: 4, display: 'flex', gap: 8, alignItems: 'center' }}>
                            <span style={{ color: senderCfg.color, fontWeight: 500 }}>{senderCfg.label}</span>
                            <span>{new Date(item.created_at).toLocaleString()}</span>
                          </div>
                          <div style={{ whiteSpace: 'pre-wrap', wordBreak: 'break-word' }}>{item.content}</div>
                        </div>
                      </div>
                    );
                  })}
                  <div ref={chatEndRef} />
                </>
              )}
            </div>
            {!isClosed && selectedId && (
              <>
                <Input.TextArea
                  rows={3}
                  placeholder="输入消息..."
                  value={draft}
                  onChange={(e) => setDraft(e.target.value)}
                  style={{ marginTop: 8 }}
                  disabled={sending}
                  onPressEnter={(e) => { if (!e.shiftKey) { e.preventDefault(); handleSend(); } }}
                />
                <div style={{ marginTop: 8, display: 'flex', justifyContent: 'flex-end' }}>
                  <Button type="primary" onClick={handleSend} loading={sending} disabled={!draft.trim()}>
                    发送消息
                  </Button>
                </div>
              </>
            )}
          </div>
        </ProCard>
      </div>

      <Modal
        title="转派会话"
        open={transferModalOpen}
        onCancel={() => { setTransferModalOpen(false); setTransferTarget(null); }}
        onOk={handleTransfer}
        confirmLoading={operating}
        okButtonProps={{ disabled: !transferTarget }}
      >
        <p style={{ marginBottom: 12 }}>选择目标客服：</p>
        <Select
          style={{ width: '100%' }}
          placeholder="选择客服"
          value={transferTarget}
          onChange={(v) => setTransferTarget(v)}
          options={onlineAgents.map((a: any) => ({ label: a.name || `客服 #${a.id}`, value: a.id }))}
        />
      </Modal>
    </div>
  );
};

export default ConversationPage;
