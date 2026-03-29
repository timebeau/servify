import React, { useState, useEffect } from 'react';
import { ProDescriptions } from '@ant-design/pro-components';
import { ProCard } from '@ant-design/pro-components';
import { Tag, Steps, Button, Space, Divider, Input, List, Spin, message } from 'antd';
import { goBack, useDetailParams } from '@/lib/navigation';
import { TICKET_STATUS_MAP, TICKET_PRIORITY_MAP } from '@/utils/constants';
import { getTicket, getComments, addComment } from '@/services/ticket';

const TicketDetailPage: React.FC = () => {
  const { id } = useDetailParams();
  const [loading, setLoading] = useState(true);
  const [ticket, setTicket] = useState<API.Ticket | null>(null);
  const [comments, setComments] = useState<any[]>([]);
  const [commentContent, setCommentContent] = useState('');

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
      try {
        const result = await getComments(Number(id));
        setComments(Array.isArray(result) ? result : result?.data || []);
      } catch (error) {
        console.error('获取评论失败:', error);
      }
    };
    fetchComments();
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
    try {
      await addComment(Number(id), { content: commentContent });
      message.success('评论已添加');
      setCommentContent('');
      const result = await getComments(Number(id));
      setComments(Array.isArray(result) ? result : result?.data || []);
    } catch (error) {
      message.error('添加评论失败');
    }
  };

  return (
    <div>
      <ProCard
        title="工单详情"
        extra={
          <Space>
            <Button onClick={goBack}>返回</Button>
            <Button type="primary">回复</Button>
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

      <ProCard title="沟通记录" style={{ marginTop: 16 }}>
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
          <Button type="primary" onClick={handleAddComment} disabled={!commentContent.trim()}>
            提交评论
          </Button>
        </div>
      </ProCard>
    </div>
  );
};

export default TicketDetailPage;
