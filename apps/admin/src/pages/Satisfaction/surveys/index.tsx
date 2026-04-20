import React from 'react';
import { ProTable } from '@ant-design/pro-components';
import type { ProColumns } from '@ant-design/pro-components';
import { Tag } from 'antd';
import { listSurveys } from '@/services/satisfaction';

const SatisfactionSurveysPage: React.FC = () => {
  const columns: ProColumns<API.SatisfactionSurvey>[] = [
    {
      title: 'ID',
      dataIndex: 'id',
      width: 80,
    },
    {
      title: '客户',
      dataIndex: 'customer',
      width: 120,
    },
    {
      title: '客服',
      dataIndex: 'agent',
      width: 120,
    },
    {
      title: '工单ID',
      dataIndex: 'ticket_id',
      width: 100,
    },
    {
      title: '评分',
      dataIndex: 'rating',
      width: 80,
      sorter: true,
      render: (_, record) => (
        <Tag color={record.rating >= 4 ? 'green' : record.rating >= 3 ? 'orange' : 'red'}>
          {record.rating}分
        </Tag>
      ),
    },
    {
      title: '评论',
      dataIndex: 'comment',
      ellipsis: true,
    },
    {
      title: '评价时间',
      dataIndex: 'created_at',
      valueType: 'dateTime',
      width: 180,
      sorter: true,
    },
  ];

  return (
    <ProTable<API.SatisfactionSurvey>
      headerTitle="评价列表"
      rowKey="id"
      columns={columns}
      request={async (params) => {
        try {
          const result = await listSurveys({
            page: params.current,
            page_size: params.pageSize,
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
