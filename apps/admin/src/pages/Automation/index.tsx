import React from 'react';
import { ProTable } from '@ant-design/pro-components';
import type { ProColumns } from '@ant-design/pro-components';
import { Alert, Button, Space, Switch, Tag, message } from 'antd';
import { PlusOutlined } from '@ant-design/icons';
import { listAutomations, deleteAutomation, runAutomation } from '@/services/automation';

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
      title: '触发条件',
      dataIndex: 'trigger_type',
      width: 200,
    },
    {
      title: '执行动作',
      dataIndex: 'actions',
      width: 200,
      render: (_, record) =>
        record.actions ? JSON.stringify(record.actions) : '-',
    },
    {
      title: '状态',
      dataIndex: 'enabled',
      width: 100,
      render: (_, record) => (
        <Switch
          checked={record.enabled}
          checkedChildren="启用"
          unCheckedChildren="停用"
          onChange={() => {}}
        />
      ),
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      valueType: 'dateTime',
      width: 180,
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
              } catch (error) {
                message.error('执行失败：参数契约不匹配');
              }
            }}
            style={{ color: '#999', cursor: 'not-allowed' }}
          >
            执行（未实现）
          </a>
          <a
            onClick={async () => {
              try {
                await deleteAutomation(record.id);
                message.success('规则已删除');
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
    <>
      <Alert
        message="自动化执行功能未完成"
        description="前端发送 { automation_id } 但后端期望 { event, ticket_ids, dry_run }。
        参数契约不匹配，执行按钮无法正常工作。需要修复前端服务或后端 API 以统一参数格式。"
        type="warning"
        showIcon
        style={{ marginBottom: 16 }}
      />
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
          const result = await listAutomations({
            page: params.current,
            page_size: params.pageSize,
          });
          return {
            data: result?.data || [],
            total: result?.total || 0,
            success: true,
          };
        } catch (error) {
          console.error('获取自动化规则失败:', error);
          return { data: [], total: 0, success: true };
        }
      }}
      pagination={{ defaultPageSize: 20 }}
    />
    </>
  );
};

export default AutomationPage;
