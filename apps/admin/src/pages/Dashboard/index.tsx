import React, { useState, useEffect } from 'react';
import { ProCard, StatisticCard } from '@ant-design/pro-components';
import { Col, Row, Spin } from 'antd';
import { getDashboardStats } from '@/services/statistics';

const DashboardPage: React.FC = () => {
  const [loading, setLoading] = useState(true);
  const [stats, setStats] = useState<API.DashboardStats>({
    total_sessions: 0,
    total_tickets: 0,
    total_customers: 0,
    total_agents: 0,
    today_tickets: 0,
    today_sessions: 0,
    today_messages: 0,
    open_tickets: 0,
    assigned_tickets: 0,
    resolved_tickets: 0,
    closed_tickets: 0,
    online_agents: 0,
    busy_agents: 0,
    active_sessions: 0,
    avg_response_time: 0,
    avg_resolution_time: 0,
    customer_satisfaction: 0,
    ai_usage_today: 0,
    weknora_usage_today: 0,
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
              value: stats.total_sessions,
              description: (
                <StatisticCard.Statistic
                  title="今日新增"
                  value={stats.today_sessions}
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
                  title="今日新增"
                  value={stats.today_tickets}
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
                  title="在线客服"
                  value={stats.online_agents}
                />
              ),
            }}
          />
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <StatisticCard
            statistic={{
              title: '满意度',
              value: stats.customer_satisfaction
                ? `${(stats.customer_satisfaction * 20).toFixed(1)}%`
                : '0%',
              description: (
                <StatisticCard.Statistic
                  title="活跃会话"
                  value={stats.active_sessions}
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
        <Row gutter={[16, 16]}>
          <Col xs={24} sm={8}>
            <StatisticCard
              statistic={{ title: '今日消息量', value: stats.today_messages }}
            />
          </Col>
          <Col xs={24} sm={8}>
            <StatisticCard
              statistic={{ title: '已解决工单', value: stats.resolved_tickets }}
            />
          </Col>
          <Col xs={24} sm={8}>
            <StatisticCard
              statistic={{
                title: '平均响应时长',
                value: `${stats.avg_response_time.toFixed(1)}s`,
              }}
            />
          </Col>
        </Row>
      </ProCard>
    </div>
  );
};

export default DashboardPage;
