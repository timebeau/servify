import React, { useState, useEffect } from 'react';
import { ProCard, StatisticCard } from '@ant-design/pro-components';
import { Row, Col, Spin, Empty, Result, Button } from 'antd';
import { Column, Pie } from '@ant-design/charts';
import { getSatisfactionStats } from '@/services/satisfaction';

type TrendChartDatum = {
  date: string;
  value: number;
  count: number;
};

type DistributionChartDatum = {
  rating: string;
  count: number;
};

const SatisfactionOverviewPage: React.FC = () => {
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [stats, setStats] = useState<API.SatisfactionStats | null>(null);

  const fetchStats = async () => {
    setLoading(true);
    setError(null);
    try {
      const result = await getSatisfactionStats();
      if (result) {
        setStats(result);
      }
    } catch (err) {
      console.error('获取满意度统计失败:', err);
      setError('获取满意度统计失败，请重试');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchStats();
  }, []);

  if (loading) {
    return (
      <div style={{ textAlign: 'center', padding: 80 }}>
        <Spin size="large" tip="加载中..." />
      </div>
    );
  }

  if (error) {
    return (
      <Result
        status="error"
        title="加载失败"
        subTitle={error}
        extra={<Button type="primary" onClick={fetchStats}>重试</Button>}
      />
    );
  }

  // 满意度趋势图数据
  const trendData = (stats?.trend_data || []).map((item) => ({
    date: item.date,
    value: Number(item.average_rating.toFixed(2)),
    count: item.count,
  }));

  // 评分分布数据
  const distEntries = Object.entries(stats?.rating_distribution || {}) as [string, number][];
  const distributionData = distEntries.map(([rating, count]) => ({
    rating: `${rating} 星`,
    count: Number(count),
  }));

  return (
    <div>
      <Row gutter={[16, 16]}>
        <Col xs={24} sm={12} lg={6}>
          <StatisticCard
            statistic={{
              title: '平均满意度',
              value: stats?.average_rating ? Number(stats.average_rating).toFixed(1) : 0,
              suffix: '分',
            }}
          />
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <StatisticCard
            statistic={{
              title: '评价总数',
              value: stats?.total_ratings || 0,
            }}
          />
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <StatisticCard
            statistic={{
              title: '好评率',
              value: stats?.category_stats?.positive
                ? `${(stats.category_stats.positive.average_rating * 20).toFixed(1)}%`
                : '0%',
            }}
          />
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <StatisticCard
            statistic={{
              title: '差评率',
              value: stats?.category_stats?.negative
                ? `${(stats.category_stats.negative.average_rating * 20).toFixed(1)}%`
                : '0%',
            }}
          />
        </Col>
      </Row>

      <ProCard title="满意度趋势" style={{ marginTop: 16 }} collapsible>
        {trendData.length === 0 ? (
          <Empty description="暂无趋势数据" image={Empty.PRESENTED_IMAGE_SIMPLE} />
        ) : (
          <Column
            data={trendData}
            xField="date"
            yField="value"
            color="#1677ff"
            label={{
              text: (origin: TrendChartDatum) => origin.value?.toFixed(1),
              position: 'top' as const,
            }}
            meta={{
              count: { alias: '评价数' },
            }}
            tooltip={{
              fields: ['date', 'value', 'count'],
              formatter: (datum: TrendChartDatum) => ({
                name: datum.date,
                value: `评分: ${datum.value} / 评价数: ${datum.count || 0}`,
              }),
            }}
            style={{ height: 300 }}
          />
        )}
      </ProCard>

      <ProCard title="评分分布" style={{ marginTop: 16 }} collapsible>
        {distributionData.length === 0 ? (
          <Empty description="暂无评分数据" image={Empty.PRESENTED_IMAGE_SIMPLE} />
        ) : (
          <Pie
            data={distributionData}
            angleField="count"
            colorField="rating"
            label={{
              type: 'inner' as const,
              offset: '-30%',
              content: (obj: DistributionChartDatum) => `${obj.rating}: ${obj.count}`,
            }}
            style={{ height: 300 }}
          />
        )}
      </ProCard>
    </div>
  );
};

export default SatisfactionOverviewPage;
