import React, { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import { ProTable, StatisticCard } from '@ant-design/pro-components';
import type { ActionType, ProColumns } from '@ant-design/pro-components';
import { PlusOutlined } from '@ant-design/icons';
import { Button, Col, Form, Input, Modal, Row, Select, Space, Tag, message } from 'antd';
import { navigateTo, useQueryParam } from '@/lib/navigation';
import { listCustomers } from '@/services/customer';
import { getRemoteAssistTicketStats } from '@/services/statistics';
import { createTicket, listTickets, updateTicket } from '@/services/ticket';
import { TICKET_PRIORITY_MAP, TICKET_STATUS_MAP } from '@/utils/constants';
import { getErrorMessage, isFormValidationError } from '@/utils/error';

const TICKET_SOURCE_OPTIONS = [
  { label: '远程协助', value: 'remote_assist' },
  { label: 'Web', value: 'web' },
  { label: 'Chat', value: 'chat' },
  { label: 'Email', value: 'email' },
  { label: 'Phone', value: 'phone' },
];

const TICKET_STATUS_OPTIONS = Object.entries(TICKET_STATUS_MAP).map(([value, item]) => ({
  label: item.text,
  value,
}));

const TICKET_PRIORITY_OPTIONS = Object.entries(TICKET_PRIORITY_MAP).map(([value, item]) => ({
  label: item.text,
  value,
}));

function parseTags(value?: string) {
  return value
    ? value
        .split(',')
        .map((item) => item.trim())
        .filter(Boolean)
    : [];
}

const TicketListPage: React.FC = () => {
  const actionRef = useRef<ActionType>();
  const autoOpenedRef = useRef(false);
  const [form] = Form.useForm();
  const [quickView, setQuickView] = useState<'all' | 'remote_assist'>('all');
  const [modalOpen, setModalOpen] = useState(false);
  const [editingTicket, setEditingTicket] = useState<API.Ticket | null>(null);
  const [submitting, setSubmitting] = useState(false);
  const [customerOptions, setCustomerOptions] = useState<Array<{ label: string; value: number }>>([]);
  const [remoteAssistStats, setRemoteAssistStats] = useState({
    total: 0,
    open: 0,
    resolved: 0,
    closed: 0,
    resolvedRate: 0,
    closedRate: 0,
    avgCloseHours: 0,
  });
  const createQuery = useQueryParam('create');
  const customerIdQuery = useQueryParam('customer_id');

  const prefilledCustomerId = useMemo(() => {
    const parsed = Number(customerIdQuery);
    return Number.isInteger(parsed) && parsed > 0 ? parsed : undefined;
  }, [customerIdQuery]);

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

  useEffect(() => {
    const fetchCustomers = async () => {
      try {
        const result = await listCustomers({ page: 1, page_size: 200 });
        const options = result.data.map((customer) => ({
          label: customer.name ? `${customer.name} (#${customer.id})` : `客户 #${customer.id}`,
          value: customer.id,
        }));

        if (prefilledCustomerId && !options.some((item) => item.value === prefilledCustomerId)) {
          options.unshift({
            label: `客户 #${prefilledCustomerId}`,
            value: prefilledCustomerId,
          });
        }

        setCustomerOptions(options);
      } catch (error) {
        console.error('获取客户列表失败:', error);
      }
    };

    void fetchCustomers();
  }, [prefilledCustomerId]);

  const renderSourceTag = (record: API.Ticket) => {
    const customFields = record.custom_fields || {};
    const tags = Array.isArray(record.tag_list)
      ? record.tag_list
      : Array.isArray(record.tags)
        ? record.tags
        : typeof record.tags === 'string'
          ? record.tags
              .split(',')
              .map((item) => item.trim())
              .filter(Boolean)
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

  const closeModal = () => {
    setModalOpen(false);
    setEditingTicket(null);
    form.resetFields();
  };

  const openCreate = (customerId?: number) => {
    setEditingTicket(null);
    form.resetFields();
    form.setFieldsValue({
      customer_id: customerId,
      priority: 'medium',
      source: quickView === 'remote_assist' ? 'remote_assist' : 'web',
    });
    setModalOpen(true);
  };

  const openEdit = (record: API.Ticket) => {
    setEditingTicket(record);
    form.setFieldsValue({
      title: record.title,
      description: record.description,
      category: record.category,
      priority: record.priority,
      status: record.status,
      tags: Array.isArray(record.tag_list)
        ? record.tag_list.join(', ')
        : typeof record.tags === 'string'
          ? record.tags
          : '',
    });
    setModalOpen(true);
  };

  useEffect(() => {
    if (createQuery !== '1' || autoOpenedRef.current) {
      return;
    }

    autoOpenedRef.current = true;
    openCreate(prefilledCustomerId);

    const url = new URL(window.location.href);
    url.searchParams.delete('create');
    window.history.replaceState({}, '', `${url.pathname}${url.search}`);
  }, [createQuery, prefilledCustomerId]);

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields();
      setSubmitting(true);

      if (editingTicket) {
        await updateTicket(editingTicket.id, {
          title: values.title?.trim(),
          description: values.description?.trim() || undefined,
          category: values.category?.trim() || undefined,
          priority: values.priority,
          status: values.status,
          tags: values.tags ? parseTags(values.tags).join(',') : undefined,
        });
        message.success('工单已更新');
      } else {
        await createTicket({
          title: values.title?.trim(),
          description: values.description?.trim() || undefined,
          customer_id: values.customer_id,
          category: values.category?.trim() || undefined,
          priority: values.priority,
          source: values.source,
          tags: parseTags(values.tags),
        });
        message.success('工单已创建');
      }

      closeModal();
      void fetchRemoteAssistStats();
      actionRef.current?.reload();
    } catch (error: unknown) {
      if (isFormValidationError(error)) {
        return;
      }
      message.error(getErrorMessage(error, editingTicket ? '更新失败' : '创建失败'));
    } finally {
      setSubmitting(false);
    }
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
        Object.entries(TICKET_STATUS_MAP).map(([key, item]) => [key, { text: item.text }]),
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
        Object.entries(TICKET_PRIORITY_MAP).map(([key, item]) => [key, { text: item.text }]),
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
      valueEnum: Object.fromEntries(
        TICKET_SOURCE_OPTIONS.map((item) => [item.value, { text: item.label }]),
      ),
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
      width: 140,
    },
    {
      title: '客服',
      dataIndex: 'agent_name',
      width: 140,
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
      width: 140,
      render: (_, record) => (
        <Space>
          <a onClick={() => navigateTo(`/ticket/detail/${record.id}`)}>查看</a>
          <a onClick={() => openEdit(record)}>编辑</a>
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
        headerTitle={quickView === 'remote_assist' ? '工单列表 / 远程协助视图' : '工单列表'}
        rowKey="id"
        actionRef={actionRef}
        columns={columns}
        toolBarRender={() => [
          <Space key="views">
            <Button type={quickView === 'all' ? 'primary' : 'default'} onClick={() => applyQuickView('all')}>
              全部工单
            </Button>
            <Button
              type={quickView === 'remote_assist' ? 'primary' : 'default'}
              onClick={() => applyQuickView('remote_assist')}
            >
              远程协助
            </Button>
          </Space>,
          <Button
            key="create"
            type="primary"
            icon={<PlusOutlined />}
            onClick={() => openCreate(prefilledCustomerId)}
          >
            新建工单
          </Button>,
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

      <Modal
        title={editingTicket ? '编辑工单' : '新建工单'}
        open={modalOpen}
        onCancel={closeModal}
        onOk={handleSubmit}
        confirmLoading={submitting}
        okText={editingTicket ? '保存' : '创建'}
        width={720}
      >
        <Form form={form} layout="vertical" style={{ marginTop: 16 }}>
          {!editingTicket ? (
            <Form.Item
              name="customer_id"
              label="客户"
              rules={[{ required: true, message: '请选择客户' }]}
            >
              <Select
                showSearch
                options={customerOptions}
                optionFilterProp="label"
                placeholder="选择客户"
              />
            </Form.Item>
          ) : null}
          <Form.Item name="title" label="标题" rules={[{ required: true, message: '请输入工单标题' }]}>
            <Input placeholder="工单标题" />
          </Form.Item>
          <Form.Item name="description" label="描述">
            <Input.TextArea rows={4} placeholder="问题描述" />
          </Form.Item>
          <Form.Item name="category" label="分类">
            <Input placeholder="例如：billing / remote-assist" />
          </Form.Item>
          <Form.Item
            name="priority"
            label="优先级"
            rules={[{ required: true, message: '请选择优先级' }]}
          >
            <Select options={TICKET_PRIORITY_OPTIONS} />
          </Form.Item>
          {!editingTicket ? (
            <Form.Item name="source" label="来源">
              <Select options={TICKET_SOURCE_OPTIONS} />
            </Form.Item>
          ) : (
            <Form.Item
              name="status"
              label="状态"
              rules={[{ required: true, message: '请选择状态' }]}
            >
              <Select options={TICKET_STATUS_OPTIONS} />
            </Form.Item>
          )}
          <Form.Item name="tags" label="标签">
            <Input placeholder="使用逗号分隔，例如：vip, remote_assist" />
          </Form.Item>
        </Form>
      </Modal>
    </>
  );
};

export default TicketListPage;
