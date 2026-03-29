import React, { useState, useEffect } from 'react';
import { ProDescriptions } from '@ant-design/pro-components';
import { ProCard } from '@ant-design/pro-components';
import { Tag, Button, Switch, Space, Spin, message } from 'antd';
import { goBack, useDetailParams } from '@/lib/navigation';
import { AGENT_STATUS_MAP } from '@/utils/constants';
import { getAgent, updateAgentStatus } from '@/services/agent';

const AgentDetailPage: React.FC = () => {
  const { id } = useDetailParams();
  const [loading, setLoading] = useState(true);
  const [agent, setAgent] = useState<API.Agent | null>(null);

  useEffect(() => {
    const fetchAgent = async () => {
      if (!id) return;
      setLoading(true);
      try {
        const result = await getAgent(Number(id));
        if (result) {
          setAgent(result);
        }
      } catch (error) {
        console.error('获取客服详情失败:', error);
      } finally {
        setLoading(false);
      }
    };
    fetchAgent();
  }, [id]);

  if (loading) {
    return (
      <div style={{ textAlign: 'center', padding: 80 }}>
        <Spin size="large" tip="加载中..." />
      </div>
    );
  }

  if (!agent) {
    return (
      <div style={{ textAlign: 'center', padding: 80, color: '#999' }}>
        客服不存在或加载失败
        <br />
        <Button style={{ marginTop: 16 }} onClick={goBack}>
          返回
        </Button>
      </div>
    );
  }

  const statusItem = AGENT_STATUS_MAP[agent.status];

  const handleStatusToggle = async (checked: boolean) => {
    if (!id) return;
    try {
      const newStatus = checked ? 'online' : 'offline';
      await updateAgentStatus(Number(id), newStatus);
      setAgent({ ...agent, status: newStatus });
      message.success(`状态已切换为${checked ? '在线' : '离线'}`);
    } catch (error) {
      message.error('状态切换失败');
    }
  };

  return (
    <div>
      <ProCard
        title="客服详情"
        extra={
          <Space>
            <Button onClick={goBack}>返回</Button>
            <Button>编辑</Button>
            <Button type="primary">分配工单</Button>
          </Space>
        }
      >
        <ProDescriptions
          column={2}
          dataSource={agent}
          columns={[
            { title: '客服ID', dataIndex: 'id' },
            { title: '姓名', dataIndex: 'name' },
            { title: '邮箱', dataIndex: 'email', copyable: true },
            {
              title: '状态',
              dataIndex: 'status',
              render: () =>
                statusItem ? (
                  <Tag color={statusItem.color}>{statusItem.text}</Tag>
                ) : (
                  agent.status
                ),
            },
            {
              title: '技能',
              dataIndex: 'skills',
              render: () =>
                agent.skills && agent.skills.length > 0
                  ? agent.skills.map((s) => <Tag key={s}>{s}</Tag>)
                  : '-',
            },
            {
              title: '会话数',
              render: () => `${agent.current_sessions || 0}/${agent.max_sessions || '-'}`,
            },
            { title: '创建时间', dataIndex: 'created_at' },
          ]}
        />
      </ProCard>

      <ProCard title="在线状态切换" style={{ marginTop: 16 }}>
        <Space>
          <span>在线状态：</span>
          <Switch
            checkedChildren="在线"
            unCheckedChildren="离线"
            checked={agent.status === 'online' || agent.status === 'busy'}
            onChange={handleStatusToggle}
          />
        </Space>
      </ProCard>
    </div>
  );
};

export default AgentDetailPage;
