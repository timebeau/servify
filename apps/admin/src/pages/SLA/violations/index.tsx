import React from 'react';
import { ProTable } from '@ant-design/pro-components';
import type { ProColumns } from '@ant-design/pro-components';
import { Tag, message } from 'antd';
import { listSLAViolations, resolveSLAViolation } from '@/services/sla';
import { getErrorMessage } from '@/utils/error';

const VIOLATION_TYPE_LABELS: Record<string, string> = {
  first_response: '首次响应超时',
  resolution: '解决超时',
  escalation: '升级超时',
};

const SLAViolationsPage: React.FC = () => {
  const columns: ProColumns<API.SLAViolation>[] = [
    {
      title: 'ID',
      dataIndex: 'id',
      width: 80,
    },
    {
      title: '工单ID',
      dataIndex: 'ticket_id',
      width: 100,
    },
    {
      title: '违规类型',
      dataIndex: 'violation_type',
      width: 160,
      search: false,
      render: (_, record) => VIOLATION_TYPE_LABELS[record.violation_type] || record.violation_type,
    },
    {
      title: '违规时间',
      dataIndex: 'violated_at',
      valueType: 'dateTime',
      width: 180,
      search: false,
    },
    {
      title: '状态',
      dataIndex: 'status',
      width: 100,
      valueType: 'select',
      valueEnum: {
        pending: { text: '待处理' },
        resolved: { text: '已解决' },
      },
      render: (_, record) => (
        <Tag color={record.resolved ? 'green' : 'red'}>
          {record.resolved ? '已解决' : '待处理'}
        </Tag>
      ),
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
        record.resolved ? null : (
          <a
            onClick={async () => {
              try {
                await resolveSLAViolation(record.id);
                message.success('违规已处理');
              } catch (error: unknown) {
                message.error(getErrorMessage(error, '处理失败'));
              }
            }}
          >
            处理
          </a>
        ),
    },
  ];

  return (
    <ProTable<API.SLAViolation>
      headerTitle="SLA 违规记录"
      rowKey="id"
      columns={columns}
      request={async (params) => {
        try {
          const result = await listSLAViolations({
            page: params.current,
            page_size: params.pageSize,
            status: typeof params.status === 'string' ? params.status : undefined,
          });
          return {
            data: result.data,
            total: result.total,
            success: true,
          };
        } catch (error: unknown) {
          console.error('获取 SLA 违规记录失败:', error);
          return { data: [], total: 0, success: true };
        }
      }}
      search={{ filterType: 'light' }}
      pagination={{ defaultPageSize: 20 }}
    />
  );
};

export default SLAViolationsPage;
