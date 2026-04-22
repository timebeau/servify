import React, { useRef } from 'react';
import { ProTable } from '@ant-design/pro-components';
import type { ProColumns, ActionType } from '@ant-design/pro-components';
import { Tag, message } from 'antd';
import { listSurveys, resendSurvey } from '@/services/satisfaction';
import { getErrorMessage } from '@/utils/error';

const SatisfactionSurveysPage: React.FC = () => {
  const actionRef = useRef<ActionType>();
  const columns: ProColumns<API.SatisfactionSurvey>[] = [
    {
      title: 'ID',
      dataIndex: 'id',
      width: 80,
    },
    {
      title: '客户ID',
      dataIndex: 'customer_id',
      width: 120,
      search: false,
    },
    {
      title: '客服ID',
      dataIndex: 'agent_id',
      width: 120,
      search: false,
    },
    {
      title: '工单ID',
      dataIndex: 'ticket_id',
      width: 100,
    },
    {
      title: '渠道',
      dataIndex: 'channel',
      width: 100,
      valueType: 'select',
      valueEnum: {
        email: { text: 'Email' },
        sms: { text: 'SMS' },
        web: { text: 'Web' },
      },
      render: (_, record) => <Tag>{record.channel || '-'}</Tag>,
    },
    {
      title: '状态',
      dataIndex: 'status',
      width: 100,
      valueType: 'select',
      valueEnum: {
        queued: { text: '待发送' },
        sent: { text: '已发送' },
        completed: { text: '已完成' },
        expired: { text: '已过期' },
      },
      render: (_, record) => (
        <Tag color={record.status === 'completed' ? 'green' : record.status === 'expired' ? 'red' : 'blue'}>
          {record.status || '-'}
        </Tag>
      ),
    },
    {
      title: '发送时间',
      dataIndex: 'sent_at',
      valueType: 'dateTime',
      width: 180,
      search: false,
    },
    {
      title: '过期时间',
      dataIndex: 'expires_at',
      valueType: 'dateTime',
      width: 180,
      search: false,
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      valueType: 'dateTime',
      width: 180,
      search: false,
    },
    {
      title: '操作',
      valueType: 'option',
      width: 100,
      render: (_, record) =>
        record.status === 'completed' ? null : (
          <a
            onClick={async () => {
              try {
                await resendSurvey(record.id);
                message.success('问卷已重发');
                actionRef.current?.reload();
              } catch (error: unknown) {
                message.error(getErrorMessage(error, '重发失败'));
              }
            }}
          >
            重发
          </a>
        ),
    },
  ];

  return (
    <ProTable<API.SatisfactionSurvey>
      headerTitle="评价列表"
      rowKey="id"
      actionRef={actionRef}
      columns={columns}
      request={async (params) => {
        try {
          const result = await listSurveys({
            page: params.current,
            page_size: params.pageSize,
            ticket_id: typeof params.ticket_id === 'number' ? params.ticket_id : undefined,
            status: typeof params.status === 'string' ? [params.status] : undefined,
            channel: typeof params.channel === 'string' ? [params.channel] : undefined,
          });
          return {
            data: result?.data || [],
            total: result?.total || 0,
            success: true,
          };
        } catch (error) {
          console.error('获取评价列表失败:', error);
          return { data: [], total: 0, success: true };
        }
      }}
      search={{ filterType: 'light' }}
      pagination={{ defaultPageSize: 20 }}
    />
  );
};

export default SatisfactionSurveysPage;
