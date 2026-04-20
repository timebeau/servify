import React, { useState, useEffect } from 'react';
import { ProDescriptions } from '@ant-design/pro-components';
import { ProCard } from '@ant-design/pro-components';
import { Tag, Button, Switch, Space, Spin, message, Modal, Form, Input, InputNumber } from 'antd';
import { goBack, useDetailParams } from '@/lib/navigation';
import { AGENT_STATUS_MAP } from '@/utils/constants';
import { getAgent, updateAgentStatus, updateAgent } from '@/services/agent';
import { getErrorMessage, isFormValidationError } from '@/utils/error';

const AgentDetailPage: React.FC = () => {
  const { id } = useDetailParams();
  const [loading, setLoading] = useState(true);
  const [agent, setAgent] = useState<API.Agent | null>(null);
  const [editOpen, setEditOpen] = useState(false);
  const [saving, setSaving] = useState(false);
  const [form] = Form.useForm();

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

  const handleStatusToggle = async (checked: boolean) => {
    if (!id) return;
    try {
      const newStatus = checked ? 'online' : 'offline';
      await updateAgentStatus(Number(id), newStatus);
      setAgent({ ...agent!, status: newStatus });
      message.success(`状态已切换为${checked ? '在线' : '离线'}`);
    } catch (error: unknown) {
      message.error(getErrorMessage(error, '状态切换失败'));
    }
  };

  const handleEdit = async () => {
    try {
      const values = await form.validateFields();
      setSaving(true);
      const skills = values.skills
        ? values.skills.split(',').map((s: string) => s.trim()).filter(Boolean)
        : [];
      await updateAgent(Number(id), {
        name: values.name,
        email: values.email,
        skills,
        max_concurrent: values.max_concurrent,
      });
      message.success('客服信息已更新');
      setEditOpen(false);
      // 刷新数据
      const result = await getAgent(Number(id));
      if (result) setAgent(result);
    } catch (error: unknown) {
      if (isFormValidationError(error)) return;
      message.error(getErrorMessage(error, '更新失败'));
    } finally {
      setSaving(false);
    }
  };

  const openEditModal = () => {
    if (!agent) return;
    form.setFieldsValue({
      name: agent.name,
      email: agent.email,
      skills: agent.skills?.join(', ') || '',
      max_concurrent: agent.max_sessions || 5,
    });
    setEditOpen(true);
  };

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

  return (
    <div>
      <ProCard
        title="客服详情"
        extra={
          <Space>
            <Button onClick={goBack}>返回</Button>
            <Button onClick={openEditModal}>编辑</Button>
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

      <Modal
        title="编辑客服"
        open={editOpen}
        onCancel={() => setEditOpen(false)}
        onOk={handleEdit}
        confirmLoading={saving}
        okText="保存"
      >
        <Form form={form} layout="vertical" style={{ marginTop: 16 }}>
          <Form.Item name="name" label="姓名" rules={[{ required: true, message: '请输入姓名' }]}>
            <Input placeholder="客服姓名" />
          </Form.Item>
          <Form.Item name="email" label="邮箱" rules={[
            { required: true, message: '请输入邮箱' },
            { type: 'email', message: '请输入有效邮箱' },
          ]}>
            <Input placeholder="邮箱" />
          </Form.Item>
          <Form.Item name="skills" label="技能标签（逗号分隔）">
            <Input placeholder="如: 售前, 售后, 技术" />
          </Form.Item>
          <Form.Item name="max_concurrent" label="最大并发会话数">
            <InputNumber min={1} max={50} style={{ width: '100%' }} />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};

export default AgentDetailPage;
