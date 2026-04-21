import React from 'react';
import { ProTable } from '@ant-design/pro-components';
import type { ProColumns } from '@ant-design/pro-components';
import { Button, Space, Switch, message } from 'antd';
import { PlusOutlined } from '@ant-design/icons';
import { listAutomations, deleteAutomation, runAutomation } from '@/services/automation';
import { getErrorMessage } from '@/utils/error';

const AutomationPage: React.FC = () => {
  const columns: ProColumns<API.Automation>[] = [
    {
      title: 'ID',
      dataIndex: 'id',
      width: 80,
    },
    {
      title: '规则名称',
      dataIndex: 'name',
      search: true,
    },
    {
      title: '触发事件',
      dataIndex: 'event',
      width: 180,
      search: false,
    },
    {
      title: '执行动作',
      dataIndex: 'actions',
      width: 240,
      search: false,
      render: (_, record) => {
        if (typeof record.actions === 'string') {
          return record.actions;
        }
        if (record.actions) {
          return JSON.stringify(record.actions);
        }
        return '-';
      },
    },
    {
      title: '状态',
      dataIndex: 'active',
      width: 100,
      search: false,
      render: (_, record) => (
        <Switch
          checked={Boolean(record.active ?? record.enabled)}
          checkedChildren="启用"
          unCheckedChildren="停用"
          disabled
          onChange={() => {}}
        />
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
      width: 160,
      render: (_, record) => (
        <Space>
          <a>编辑</a>
          <a
            onClick={async () => {
              try {
                await runAutomation(record.id);
                message.success('规则已执行');
              } catch (error: unknown) {
                message.error(getErrorMessage(error, '执行失败'));
              }
            }}
          >
            执行
          </a>
          <a
            onClick={async () => {
              try {
                await deleteAutomation(record.id);
                message.success('规则已删除');
              } catch (error: unknown) {
                message.error(getErrorMessage(error, '删除失败'));
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
    <ProTable<API.Automation>
      headerTitle="自动化规则"
      rowKey="id"
      columns={columns}
      toolBarRender={() => [
        <Button key="create" type="primary" icon={<PlusOutlined />}>
          新建规则
        </Button>,
      ]}
      request={async (params) => {
        try {
          const result = await listAutomations();
          const keyword = typeof params.name === 'string' ? params.name.trim().toLowerCase() : '';
          let data = result.data;
          if (keyword) {
            data = data.filter((item) => item.name.toLowerCase().includes(keyword));
          }
          const total = data.length;
          const current = params.current || 1;
          const pageSize = params.pageSize || 20;
          const pageData = data.slice((current - 1) * pageSize, current * pageSize);
          return {
            data: pageData,
            total,
            success: true,
          };
        } catch (error: unknown) {
          console.error('获取自动化规则失败:', error);
          return { data: [], total: 0, success: true };
        }
      }}
      pagination={{ defaultPageSize: 20 }}
    />
  );
};

export default AutomationPage;
