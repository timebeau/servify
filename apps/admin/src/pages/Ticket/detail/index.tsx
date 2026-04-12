import React, { useState, useEffect } from 'react';
import { ProDescriptions } from '@ant-design/pro-components';
import { ProCard } from '@ant-design/pro-components';
import { Tag, Steps, Button, Space, Divider, Input, List, Spin, message, Empty, Alert } from 'antd';
import { goBack, useDetailParams } from '@/lib/navigation';
import { navigateTo } from '@/lib/navigation';
import { TICKET_STATUS_MAP, TICKET_PRIORITY_MAP } from '@/utils/constants';
import { getTicket, getComments, addComment, getTicketConversations } from '@/services/ticket';

function extractRemoteAssistInfo(ticket: API.Ticket) {
  const customFields = ticket.custom_fields || {};
  const structuredRemoteAssist = customFields.remote_assist as Record<string, unknown> | undefined;

  if (structuredRemoteAssist && typeof structuredRemoteAssist === 'object') {
    return {
      sessionId:
        typeof structuredRemoteAssist.session_id === 'string'
          ? structuredRemoteAssist.session_id
          : undefined,
      customer:
        typeof structuredRemoteAssist.customer_name === 'string'
          ? structuredRemoteAssist.customer_name
          : undefined,
      channel:
        typeof structuredRemoteAssist.platform === 'string'
          ? structuredRemoteAssist.platform
          : undefined,
      assistState:
        typeof structuredRemoteAssist.assist_state_label === 'string'
          ? structuredRemoteAssist.assist_state_label
          : typeof structuredRemoteAssist.assist_state === 'string'
            ? structuredRemoteAssist.assist_state
            : undefined,
      result:
        typeof structuredRemoteAssist.result_preset === 'string' && structuredRemoteAssist.result_preset
          ? structuredRemoteAssist.result_preset
          : undefined,
      followup:
        typeof structuredRemoteAssist.summary === 'string'
          ? structuredRemoteAssist.summary
          : undefined,
      structured: true,
    };
  }

  const title = ticket.title || '';
  const description = ticket.description || '';
  const lines = description
    .split('\n')
    .map((line) => line.trim())
    .filter(Boolean);

  const isRemoteAssistTicket =
    title.includes('远程协助') ||
    lines.some((line) => line.startsWith('会话ID:')) ||
    lines.some((line) => line.startsWith('协助结果:')) ||
    lines.some((line) => line.startsWith('远程协助状态:'));

  if (!isRemoteAssistTicket) {
    return null;
  }

  const sessionLine = lines.find((line) => line.startsWith('会话ID:'));
  const customerLine = lines.find((line) => line.startsWith('客户:'));
  const channelLine = lines.find((line) => line.startsWith('渠道:'));
  const assistStateLine = lines.find((line) => line.startsWith('远程协助状态:'));
  const resultLine = lines.find((line) => line.startsWith('协助结果:'));
  const followupLine = lines.find((line) => line.startsWith('后续动作:'));

  return {
    sessionId: sessionLine?.replace('会话ID:', '').trim(),
    customer: customerLine?.replace('客户:', '').trim(),
    channel: channelLine?.replace('渠道:', '').trim(),
    assistState: assistStateLine?.replace('远程协助状态:', '').trim(),
    result: resultLine?.replace('协助结果:', '').trim(),
    followup: followupLine?.replace('后续动作:', '').trim(),
    structured: false,
  };
}

const TicketDetailPage: React.FC = () => {
  const { id } = useDetailParams();
  const [loading, setLoading] = useState(true);
  const [ticket, setTicket] = useState<API.Ticket | null>(null);
  const [comments, setComments] = useState<any[]>([]);
  const [commentContent, setCommentContent] = useState('');
  const [commentsError, setCommentsError] = useState<string | null>(null);
  const [submittingComment, setSubmittingComment] = useState(false);
  const [conversations, setConversations] = useState<any[]>([]);
  const [conversationsLoading, setConversationsLoading] = useState(false);
  const [conversationsError, setConversationsError] = useState<string | null>(null);

  useEffect(() => {
    const fetchData = async () => {
      if (!id) return;
      setLoading(true);
      try {
        const ticketData = await getTicket(Number(id));
        if (ticketData) {
          setTicket(ticketData);
        }
      } catch (error) {
        console.error('获取工单详情失败:', error);
      } finally {
        setLoading(false);
      }
    };
    fetchData();
  }, [id]);

  useEffect(() => {
    const fetchComments = async () => {
      if (!id) return;
      setCommentsError(null);
      try {
        const result = await getComments(Number(id));
        setComments(Array.isArray(result) ? result : result?.data || []);
      } catch (error) {
        console.error('获取评论失败:', error);
        setCommentsError('获取评论失败，请刷新重试');
      }
    };
    fetchComments();
  }, [id]);

  useEffect(() => {
    const fetchConversations = async () => {
      if (!id) return;
      setConversationsLoading(true);
      setConversationsError(null);
      try {
        const result = await getTicketConversations(Number(id));
        setConversations(result?.data || []);
      } catch (error) {
        console.error('获取关联会话失败:', error);
        setConversationsError('获取关联会话失败，请刷新重试');
      } finally {
        setConversationsLoading(false);
      }
    };
    fetchConversations();
  }, [id]);

  if (loading) {
    return (
      <div style={{ textAlign: 'center', padding: 80 }}>
        <Spin size="large" tip="加载中..." />
      </div>
    );
  }

  if (!ticket) {
    return (
      <div style={{ textAlign: 'center', padding: 80, color: '#999' }}>
        工单不存在或加载失败
        <br />
        <Button style={{ marginTop: 16 }} onClick={goBack}>
          返回
        </Button>
      </div>
    );
  }

  const statusItem = TICKET_STATUS_MAP[ticket.status];
  const priorityItem = TICKET_PRIORITY_MAP[ticket.priority];
  const remoteAssistInfo = extractRemoteAssistInfo(ticket);

  const statusStepMap: Record<string, number> = {
    open: 0,
    assigned: 1,
    in_progress: 2,
    pending_customer: 3,
    resolved: 4,
    closed: 5,
  };

  const statusSteps = [
    { title: '待处理', description: '工单已创建' },
    { title: '已分配', description: '已分配客服' },
    { title: '处理中', description: '客服正在处理' },
    { title: '待客户回复', description: '等待客户确认' },
    { title: '已解决', description: '问题已解决' },
    { title: '已关闭', description: '工单已关闭' },
  ];

  const handleAddComment = async () => {
    if (!commentContent.trim() || !id) return;
    setSubmittingComment(true);
    try {
      await addComment(Number(id), { content: commentContent });
      message.success('评论已添加');
      setCommentContent('');
      const result = await getComments(Number(id));
      setComments(Array.isArray(result) ? result : result?.data || []);
    } catch (error) {
      message.error('添加评论失败');
    } finally {
      setSubmittingComment(false);
    }
  };

  const platformColorMap: Record<string, string> = {
    web: 'blue',
    wechat: 'green',
    whatsapp: 'cyan',
    voice: 'purple',
  };

  return (
    <div>
      <ProCard
        title="工单详情"
        extra={
          <Space>
            <Button onClick={goBack}>返回</Button>
          </Space>
        }
      >
        <ProDescriptions
          column={2}
          dataSource={ticket}
          columns={[
            { title: '工单ID', dataIndex: 'id' },
            { title: '标题', dataIndex: 'title' },
            {
              title: '状态',
              dataIndex: 'status',
              render: () =>
                statusItem ? (
                  <Tag color={statusItem.color}>{statusItem.text}</Tag>
                ) : (
                  ticket.status
                ),
            },
            {
              title: '优先级',
              dataIndex: 'priority',
              render: () =>
                priorityItem ? (
                  <Tag color={priorityItem.color}>{priorityItem.text}</Tag>
                ) : (
                  ticket.priority
                ),
            },
            { title: '客户', dataIndex: 'customer_name' },
            { title: '客服', dataIndex: 'agent_name' },
            { title: '创建时间', dataIndex: 'created_at' },
            { title: '更新时间', dataIndex: 'updated_at' },
          ]}
        />
      </ProCard>

      <ProCard title="状态流转" style={{ marginTop: 16 }}>
        <Steps current={statusStepMap[ticket.status] ?? 0} items={statusSteps} />
      </ProCard>

      {remoteAssistInfo && (
        <ProCard title="远程协助来源" style={{ marginTop: 16 }}>
          <Alert
            type="info"
            showIcon
            message={(
              <Space wrap>
                <Tag color="blue">来自远程协助</Tag>
                {remoteAssistInfo.assistState && <Tag>{remoteAssistInfo.assistState}</Tag>}
                {remoteAssistInfo.result && <Tag color="green">{remoteAssistInfo.result}</Tag>}
              </Space>
            )}
            description={(
              <div style={{ display: 'grid', gap: 8 }}>
                {remoteAssistInfo.sessionId && (
                  <div>
                    会话 ID:
                    {' '}
                    <Button
                      type="link"
                      size="small"
                      style={{ paddingInline: 4 }}
                      onClick={() => navigateTo(`/conversation?id=${remoteAssistInfo.sessionId}`)}
                    >
                      {remoteAssistInfo.sessionId}
                    </Button>
                  </div>
                )}
                {remoteAssistInfo.customer && <div>客户: {remoteAssistInfo.customer}</div>}
                {remoteAssistInfo.channel && <div>渠道: {remoteAssistInfo.channel}</div>}
                <div>数据来源: {remoteAssistInfo.structured ? '结构化字段' : '描述文本解析'}</div>
                {remoteAssistInfo.followup && <div>后续动作: {remoteAssistInfo.followup}</div>}
              </div>
            )}
          />
        </ProCard>
      )}

      <ProCard title="沟通记录" style={{ marginTop: 16 }}>
        {commentsError && (
          <Alert type="error" message={commentsError} showIcon style={{ marginBottom: 16 }} />
        )}
        <List
          dataSource={comments}
          locale={{ emptyText: '暂无评论' }}
          renderItem={(item: any) => (
            <List.Item>
              <List.Item.Meta
                title={
                  <Space>
                    <span>{item.author || '未知'}</span>
                    <span style={{ color: '#999', fontSize: 12 }}>
                      {item.created_at}
                    </span>
                    {item.internal && <Tag>内部备注</Tag>}
                  </Space>
                }
                description={item.content}
              />
            </List.Item>
          )}
        />
        <Divider />
        <Input.TextArea
          rows={3}
          placeholder="输入评论内容..."
          value={commentContent}
          onChange={(e) => setCommentContent(e.target.value)}
          style={{ marginTop: 8 }}
        />
        <div style={{ marginTop: 8, textAlign: 'right' }}>
          <Button type="primary" onClick={handleAddComment} disabled={!commentContent.trim()} loading={submittingComment}>
            提交评论
          </Button>
        </div>
      </ProCard>

      {/* 关联会话 */}
      <ProCard title="关联会话" style={{ marginTop: 16 }}>
        {conversationsError && (
          <Alert type="error" message={conversationsError} showIcon style={{ marginBottom: 16 }} />
        )}
        <Spin spinning={conversationsLoading}>
          {conversations.length === 0 ? (
            <Empty description="该工单没有关联的会话记录" />
          ) : (
            <List
              dataSource={conversations}
              locale={{ emptyText: '暂无关联会话' }}
              renderItem={(session: any) => (
                <List.Item
                  actions={[
                    <Button
                      key="view"
                      size="small"
                      type="link"
                      onClick={() => navigateTo(`/conversation?id=${session.id}`)}
                    >
                      查看消息
                    </Button>,
                  ]}
                >
                  <List.Item.Meta
                    title={
                      <Space>
                        <Tag color={platformColorMap[session.platform] || 'default'}>
                          {session.platform || '未知'}
                        </Tag>
                        <span>{session.customer_name || '未知客户'}</span>
                        {session.started_at && (
                          <span style={{ color: '#999', fontSize: 12 }}>
                            {session.started_at}
                          </span>
                        )}
                      </Space>
                    }
                  />
                </List.Item>
              )}
            />
          )}
        </Spin>
      </ProCard>
    </div>
  );
};

export default TicketDetailPage;
