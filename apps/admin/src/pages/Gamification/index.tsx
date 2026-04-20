import React from 'react';
import { ProTable } from '@ant-design/pro-components';
import type { ProColumns } from '@ant-design/pro-components';
import { Tag } from 'antd';
import { getLeaderboard } from '@/services/gamification';

const GamificationPage: React.FC = () => {
  const columns: ProColumns<API.LeaderboardRecord>[] = [
    {
      title: '排名',
      dataIndex: 'rank',
      width: 80,
      sorter: true,
      render: (_, record) => {
        const color =
          record.rank === 1
            ? '#FFD700'
            : record.rank === 2
            ? '#C0C0C0'
            : record.rank === 3
            ? '#CD7F32'
            : undefined;
        return (
          <span style={{ fontWeight: 'bold', color }}>{record.rank}</span>
        );
      },
    },
    {
      title: '客服',
      dataIndex: 'agent',
      width: 120,
    },
    {
      title: '积分',
      dataIndex: 'score',
      width: 100,
      sorter: true,
    },
    {
      title: '解决工单数',
      dataIndex: 'resolved_tickets',
      width: 120,
      sorter: true,
    },
    {
      title: '平均评分',
      dataIndex: 'avg_rating',
      width: 100,
      sorter: true,
      render: (_, record) => (
        <Tag
          color={
            record.avg_rating >= 4
              ? 'green'
              : record.avg_rating >= 3
              ? 'orange'
              : 'red'
          }
        >
          {record.avg_rating}
        </Tag>
      ),
    },
    {
      title: '平均响应时间(秒)',
      dataIndex: 'avg_response_time',
      width: 160,
      sorter: true,
    },
  ];

  return (
    <ProTable<API.LeaderboardRecord>
      headerTitle="客服排行榜"
      rowKey="id"
      columns={columns}
      request={async (params) => {
        try {
          const result = await getLeaderboard({
            page: params.current,
            page_size: params.pageSize,
          });
          const data = Array.isArray(result) ? result : result.data || [];
          return {
            data,
            total: Array.isArray(result) ? data.length : result.total || data.length,
            success: true,
          };
        } catch (error) {
          console.error('获取排行榜失败:', error);
          return { data: [], total: 0, success: true };
        }
      }}
      search={false}
      pagination={{ defaultPageSize: 20 }}
    />
  );
};

export default GamificationPage;
