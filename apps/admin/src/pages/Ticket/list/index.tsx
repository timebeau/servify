import React, { useCallback, useEffect, useRef, useState } from 'react';
import { ProTable, ModalForm, ProFormText, ProFormTextArea, ProFormSelect, StatisticCard } from '@ant-design/pro-components';
import type { ProColumns, ActionType } from '@ant-design/pro-components';
import { Tag, Button, Space, message, Row, Col } from 'antd';
import { PlusOutlined } from '@ant-design/icons';
import { navigateTo } from '@/lib/navigation';
import { TICKET_STATUS_MAP, TICKET_PRIORITY_MAP } from '@/utils/constants';
import { listTickets, createTicket } from '@/services/ticket';
import { getRemoteAssistTicketStats } from '@/services/statistics';

const TicketListPage: React.FC = () => {
  const actionRef = useRef<ActionType>();
  const [quickView, setQuickView] = useState<'all' | 'remote_assist'>('all');
  const [remoteAssistStats, setRemoteAssistStats] = useState({
    total: 0,
    open: 0,
    resolved: 0,
    closed: 0,
    resolvedRate: 0,
    closedRate: 0,
    avgCloseHours: 0,
  });

  const fetchRemoteAssistStats = useCallback(async () => {
    try {
      const result = await getRemoteAssistTicketStats();
      setRemoteAssistStats({
        total: result?.total || 0,
        open: result?.open || 0,
        resolved: result?.resolved || 0,
        closed: result?.closed || 0,
        resolvedRate: result?.resolved_rate || 0,
        closedRate: result?.closed_rate || 0,
        avgCloseHours: result?.avg_close_hours || 0,
      });
    } catch (error) {
      console.error('获取远程协助工单统计失败:', error);
    }
  }, []);

  useEffect(() => {
    void fetchRemoteAssistStats();
  }, [fetchRemoteAssistStats]);

  const renderSourceTag = (record: API.Ticket) => {
    const customFields = record.custom_fields || {};
    const tags = Array.isArray(record.tag_list)
      ? record.tag_list
      : Array.isArray(record.tags)
        ? record.tags
        : typeof record.tags === 'string'
          ? record.tags.split(',').map((item) => item.trim()).filter(Boolean)
          : [];
    const isRemoteAssist =
      record.source === 'remote_assist' ||
      record.category === 'remote-assist' ||
      customFields.source === 'remote_assist' ||
      tags.includes('remote_assist');

    if (isRemoteAssist) {
      return <Tag color="blue">远程协助</Tag>;
    }

    if (record.source) {
      return <Tag>{record.source}</Tag>;
    }

    return '-';
  };

  const columns: ProColumns<API.Ticket>[] = [
    {
      title: 'ID',
      dataIndex: 'id',
      width: 80,
      copyable: true,
    },
    {
      title: '标题',
      dataIndex: 'title',
      ellipsis: true,
      search: true,
    },
    {
      title: '状态',
      dataIndex: 'status',
      width: 100,
      valueType: 'select',
      valueEnum: Object.fromEntries(
        Object.entries(TICKET_STATUS_MAP).map(([k, v]) => [k, { text: v.text }]),
      ),
      render: (_, record) => {
        const item = TICKET_STATUS_MAP[record.status];
        return item ? <Tag color={item.color}>{item.text}</Tag> : record.status;
      },
    },
    {
      title: '优先级',
      dataIndex: 'priority',
      width: 100,
      valueType: 'select',
      valueEnum: Object.fromEntries(
        Object.entries(TICKET_PRIORITY_MAP).map(([k, v]) => [k, { text: v.text }]),
      ),
      render: (_, record) => {
        const item = TICKET_PRIORITY_MAP[record.priority];
        return item ? <Tag color={item.color}>{item.text}</Tag> : record.priority;
      },
    },
    {
      title: '来源',
      dataIndex: 'source',
      width: 120,
      valueType: 'select',
      valueEnum: {
        remote_assist: { text: '远程协助' },
        web: { text: 'Web' },
        chat: { text: 'Chat' },
        email: { text: 'Email' },
        phone: { text: 'Phone' },
      },
      render: (_, record) => renderSourceTag(record),
    },
    {
      title: '标签',
      dataIndex: 'tag',
      hideInTable: true,
    },
    {
      title: '客户',
      dataIndex: 'customer_name',
      width: 120,
    },
    {
      title: '客服',
      dataIndex: 'agent_name',
      width: 120,
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      valueType: 'dateTime',
      width: 180,
      sorter: true,
    },
    {
      title: '操作',
      valueType: 'option',
      width: 120,
      render: (_, record) => (
        <Space>
          <a onClick={() => navigateTo(`/ticket/detail/${record.id}`)}>查看</a>
          <a>编辑</a>
        </Space>
      ),
    },
  ];

  const applyQuickView = (view: 'all' | 'remote_assist') => {
    setQuickView(view);
    actionRef.current?.reload();
  };

  return (
    <>
      <Row gutter={[16, 16]} style={{ marginBottom: 16 }}>
        <Col xs={24} sm={8}>
          <StatisticCard statistic={{ title: '远程协助工单', value: remoteAssistStats.total }} />
        </Col>
        <Col xs={24} sm={8}>
          <StatisticCard statistic={{ title: '待处理', value: remoteAssistStats.open }} />
        </Col>
        <Col xs={24} sm={8}>
          <StatisticCard
            statistic={{
              title: '已解决',
              value: remoteAssistStats.resolved,
              description: <span>解决率 {(remoteAssistStats.resolvedRate * 100).toFixed(1)}%</span>,
            }}
          />
        </Col>
        <Col xs={24} sm={8}>
          <StatisticCard
            statistic={{
              title: '已关闭',
              value: remoteAssistStats.closed,
              description: <span>关闭率 {(remoteAssistStats.closedRate * 100).toFixed(1)}%</span>,
            }}
          />
        </Col>
        <Col xs={24} sm={8}>
          <StatisticCard
            statistic={{
              title: '平均关闭时长',
              value: remoteAssistStats.avgCloseHours.toFixed(1),
              suffix: '小时',
            }}
          />
        </Col>
      </Row>
      <ProTable<API.Ticket>
        headerTitle={quickView === 'remote_assist' ? '工单列表 · 远程协助视图' : '工单列表'}
        rowKey="id"
        actionRef={actionRef}
        columns={columns}
        toolBarRender={() => [
          <Space key="views">
            <Button
              type={quickView === 'all' ? 'primary' : 'default'}
              onClick={() => applyQuickView('all')}
            >
              全部工单
            </Button>
            <Button
              type={quickView === 'remote_assist' ? 'primary' : 'default'}
              onClick={() => applyQuickView('remote_assist')}
            >
              远程协助
            </Button>
          </Space>,
          <ModalForm
            key="create"
            title="新建工单"
            trigger={
              <Button type="primary" icon={<PlusOutlined />}>
                新建工单
              </Button>
            }
            onFinish={async (values) => {
              try {
                await createTicket(values);
                message.success('工单创建成功');
                void fetchRemoteAssistStats();
                actionRef.current?.reload();
                return true;
              } catch (error) {
                message.error('工单创建失败');
                return false;
              }
            }}
          >
            <ProFormText
              name="title"
              label="标题"
              rules={[{ required: true, message: '请输入工单标题' }]}
            />
            <ProFormTextArea
              name="description"
              label="描述"
            />
            <ProFormSelect
              name="priority"
              label="优先级"
              valueEnum={Object.fromEntries(
                Object.entries(TICKET_PRIORITY_MAP).map(([k, v]) => [k, { text: v.text }]),
              )}
            />
          </ModalForm>,
        ]}
        request={async (params) => {
          try {
            const result = await listTickets({
              page: params.current,
              page_size: params.pageSize,
              status: params.status,
              priority: params.priority,
              source: quickView === 'remote_assist' ? 'remote_assist' : params.source,
              tag: quickView === 'remote_assist' ? 'remote_assist' : params.tag,
              search: params.title,
            });
            return {
              data: result?.data || [],
              total: result?.total || 0,
              success: true,
            };
          } catch (error) {
            console.error('获取工单列表失败:', error);
            message.error('获取工单列表失败');
            return { data: [], total: 0, success: false };
          }
        }}
        pagination={{ defaultPageSize: 20 }}
      />
    </>
  );
};

export default TicketListPage;
