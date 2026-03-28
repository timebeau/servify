import React, { useState, useEffect } from 'react';
import { ProCard, StatisticCard } from '@ant-design/pro-components';
import { Row, Col, Spin } from 'antd';
import { getSatisfactionStats } from '@/services/satisfaction';

const SatisfactionOverviewPage: React.FC = () => {
  const [loading, setLoading] = useState(true);
  const [stats, setStats] = useState<any>(null);

  useEffect(() => {
    const fetchStats = async () => {
      setLoading(true);
      try {
        const result = await getSatisfactionStats();
        if (result) {
          setStats(result);
        }
      } catch (error) {
        console.error('获取满意度统计失败:', error);
      } finally {
        setLoading(false);
      }
    };
    fetchStats();
  }, []);

  if (loading) {
    return (
      <div style={{ textAlign: 'center', padding: 80 }}>
        <Spin size="large" tip="加载中..." />
      </div>
    );
  }

  return (
    <div>
      <Row gutter={[16, 16]}>
        <Col xs={24} sm={12} lg={6}>
          <StatisticCard
            statistic={{
              title: '平均满意度',
              value: stats?.avg_score || 0,
              suffix: '分',
            }}
          />
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <StatisticCard
            statistic={{
              title: '评价总数',
              value: stats?.total_count || 0,
            }}
          />
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <StatisticCard
            statistic={{
              title: '好评率',
              value: stats?.positive_rate ? `${(stats.positive_rate * 100).toFixed(1)}%` : '0%',
            }}
          />
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <StatisticCard
            statistic={{
              title: '差评率',
              value: stats?.negative_rate ? `${(stats.negative_rate * 100).toFixed(1)}%` : '0%',
            }}
          />
        </Col>
      </Row>

      <ProCard title="满意度趋势" style={{ marginTop: 16 }} collapsible>
        <div style={{ textAlign: 'center', padding: 40, color: '#999' }}>
          满意度趋势图表占位 - 待接入数据
        </div>
      </ProCard>

      <ProCard title="评分分布" style={{ marginTop: 16 }} collapsible>
        <div style={{ textAlign: 'center', padding: 40, color: '#999' }}>
          评分分布图表占位 - 待接入数据
        </div>
      </ProCard>
    </div>
  );
};

export default SatisfactionOverviewPage;
