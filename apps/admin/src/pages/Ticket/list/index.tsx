import React, { useRef } from 'react';
import { ProTable, ModalForm, ProFormText, ProFormTextArea, ProFormSelect } from '@ant-design/pro-components';
import type { ProColumns, ActionType } from '@ant-design/pro-components';
import { Tag, Button, Space, message } from 'antd';
import { PlusOutlined } from '@ant-design/icons';
import { history } from '@umijs/max';
import { TICKET_STATUS_MAP, TICKET_PRIORITY_MAP } from '@/utils/constants';
import { listTickets, createTicket } from '@/services/ticket';

const TicketListPage: React.FC = () => {
  const actionRef = useRef<ActionType>();

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
          <a onClick={() => history.push(`/ticket/detail/${record.id}`)}>查看</a>
          <a>编辑</a>
        </Space>
      ),
    },
  ];

  return (
    <>
      <ProTable<API.Ticket>
        headerTitle="工单列表"
        rowKey="id"
        actionRef={actionRef}
        columns={columns}
        toolBarRender={() => [
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
                Object.entries(TICKET_PRIORITY_MAP).map(([k, v]) => [k, { label: v.text }]),
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
              search: params.title,
            });
            return {
              data: result?.data || [],
              total: result?.total || 0,
              success: true,
            };
          } catch (error) {
            console.error('获取工单列表失败:', error);
            return { data: [], total: 0, success: true };
          }
        }}
        pagination={{ defaultPageSize: 20 }}
      />
    </>
  );
};

export default TicketListPage;
