import React, { useState, useEffect } from 'react';
import { ProCard, StatisticCard } from '@ant-design/pro-components';
import { Col, Row, Spin } from 'antd';
import { getDashboardStats } from '@/services/statistics';

const DashboardPage: React.FC = () => {
  const [loading, setLoading] = useState(true);
  const [stats, setStats] = useState<API.DashboardStats>({
    total_conversations: 0,
    total_tickets: 0,
    total_customers: 0,
    avg_satisfaction: 0,
    open_tickets: 0,
    online_agents: 0,
  });

  useEffect(() => {
    const fetchStats = async () => {
      setLoading(true);
      try {
        const result = await getDashboardStats();
        if (result) {
          setStats(result);
        }
      } catch (error) {
        console.error('获取仪表板统计失败:', error);
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
              title: '总会话数',
              value: stats.total_conversations,
              description: (
                <StatisticCard.Statistic
                  title="较昨日"
                  value="暂无数据"
                />
              ),
            }}
          />
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <StatisticCard
            statistic={{
              title: '工单数',
              value: stats.total_tickets,
              description: (
                <StatisticCard.Statistic
                  title="较昨日"
                  value="暂无数据"
                />
              ),
            }}
          />
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <StatisticCard
            statistic={{
              title: '客户数',
              value: stats.total_customers,
              description: (
                <StatisticCard.Statistic
                  title="较昨日"
                  value="暂无数据"
                />
              ),
            }}
          />
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <StatisticCard
            statistic={{
              title: '满意度',
              value: stats.avg_satisfaction
                ? `${(stats.avg_satisfaction * 20).toFixed(1)}%`
                : '0%',
              description: (
                <StatisticCard.Statistic
                  title="较昨日"
                  value="暂无数据"
                />
              ),
            }}
          />
        </Col>
      </Row>

      <ProCard
        title="趋势图表"
        style={{ marginTop: 16 }}
        collapsible
      >
        <div style={{ textAlign: 'center', padding: 40, color: '#999' }}>
          图表区域占位 - 待接入数据
        </div>
      </ProCard>
    </div>
  );
};

export default DashboardPage;
