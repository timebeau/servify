import React from 'react';
import { ProTable } from '@ant-design/pro-components';
import type { ProColumns } from '@ant-design/pro-components';
import { Tag, Space, message } from 'antd';
import { SLA_VIOLATION_STATUS_MAP } from '@/utils/constants';
import { listSLAViolations, resolveSLAViolation } from '@/services/sla';

interface SLAViolationRecord {
  id: number;
  ticket_id: string;
  policy_name: string;
  metric: string;
  target: number;
  actual: number;
  status: string;
  created_at: string;
}

const SLAViolationsPage: React.FC = () => {
  const columns: ProColumns<SLAViolationRecord>[] = [
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
      title: '策略名称',
      dataIndex: 'policy_name',
    },
    {
      title: '指标',
      dataIndex: 'metric',
      width: 120,
    },
    {
      title: '目标值(分)',
      dataIndex: 'target',
      width: 120,
    },
    {
      title: '实际值(分)',
      dataIndex: 'actual',
      width: 120,
    },
    {
      title: '状态',
      dataIndex: 'status',
      width: 100,
      valueType: 'select',
      valueEnum: Object.fromEntries(
        Object.entries(SLA_VIOLATION_STATUS_MAP).map(([k, v]) => [k, { text: v.text }]),
      ),
      render: (_, record) => {
        const item = SLA_VIOLATION_STATUS_MAP[record.status];
        return item ? <Tag color={item.color}>{item.text}</Tag> : record.status;
      },
    },
    {
      title: '发生时间',
      dataIndex: 'created_at',
      valueType: 'dateTime',
      width: 180,
      sorter: true,
    },
    {
      title: '操作',
      valueType: 'option',
      width: 80,
      render: (_, record) => (
        <a
          onClick={async () => {
            try {
              await resolveSLAViolation(record.id);
              message.success('违规已处理');
            } catch (error) {
              message.error('处理失败');
            }
          }}
        >
          处理
        </a>
      ),
    },
  ];

  return (
    <ProTable<SLAViolationRecord>
      headerTitle="SLA 违规记录"
      rowKey="id"
      columns={columns}
      request={async (params) => {
        try {
          const result = await listSLAViolations({
            page: params.current,
            page_size: params.pageSize,
            status: params.status,
          });
          return {
            data: result?.data || [],
            total: result?.total || 0,
            success: true,
          };
        } catch (error) {
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
