import React from 'react';
import { ProTable } from '@ant-design/pro-components';
import type { ProColumns } from '@ant-design/pro-components';
import { Button, Space, Tag, message } from 'antd';
import { PlusOutlined } from '@ant-design/icons';
import { TICKET_PRIORITY_MAP } from '@/utils/constants';
import { listSLAConfigs, deleteSLAConfig } from '@/services/sla';

const SLAConfigsPage: React.FC = () => {
  const columns: ProColumns<API.SLAConfig>[] = [
    {
      title: 'ID',
      dataIndex: 'id',
      width: 80,
    },
    {
      title: '策略名称',
      dataIndex: 'name',
      search: true,
    },
    {
      title: '优先级',
      dataIndex: 'priority',
      width: 100,
      render: (_, record) => {
        const item = TICKET_PRIORITY_MAP[record.priority];
        return item ? <Tag color={item.color}>{item.text}</Tag> : record.priority;
      },
    },
    {
      title: '首次响应时间(分)',
      dataIndex: 'first_response_time',
      width: 150,
    },
    {
      title: '解决时间(分)',
      dataIndex: 'resolution_time',
      width: 140,
    },
    {
      title: '状态',
      dataIndex: 'enabled',
      width: 80,
      render: (_, record) => (
        <Tag color={record.enabled ? 'green' : 'default'}>
          {record.enabled ? '启用' : '停用'}
        </Tag>
      ),
    },
    {
      title: '操作',
      valueType: 'option',
      width: 120,
      render: (_, record) => (
        <Space>
          <a>编辑</a>
          <a
            onClick={async () => {
              try {
                await deleteSLAConfig(record.id);
                message.success('策略已删除');
              } catch (error) {
                message.error('删除失败');
              }
            }}
          >
            删除
          </a>
        </Space>
      ),
    },
  ];

  return (
    <ProTable<API.SLAConfig>
      headerTitle="SLA 策略配置"
      rowKey="id"
      columns={columns}
      toolBarRender={() => [
        <Button key="create" type="primary" icon={<PlusOutlined />}>
          新建策略
        </Button>,
      ]}
      request={async (params) => {
        try {
          const result = await listSLAConfigs({
            page: params.current,
            page_size: params.pageSize,
          });
          return {
            data: result?.data || [],
            total: result?.total || 0,
            success: true,
          };
        } catch (error) {
          console.error('获取 SLA 配置失败:', error);
          return { data: [], total: 0, success: true };
        }
      }}
      pagination={{ defaultPageSize: 20 }}
    />
  );
};

export default SLAConfigsPage;
