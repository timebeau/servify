import React from 'react';
import { ProTable } from '@ant-design/pro-components';
import type { ProColumns, ActionType } from '@ant-design/pro-components';
import { Button, Space, Tag, message } from 'antd';
import { PlusOutlined } from '@ant-design/icons';
import { useRef } from 'react';
import { listMacros, deleteMacro } from '@/services/macro';

const MacroPage: React.FC = () => {
  const actionRef = useRef<ActionType>();
  const columns: ProColumns<API.Macro>[] = [
    {
      title: 'ID',
      dataIndex: 'id',
      width: 80,
    },
    {
      title: '名称',
      dataIndex: 'name',
      search: true,
    },
    {
      title: '描述',
      dataIndex: 'description',
      ellipsis: true,
      search: false,
    },
    {
      title: '语言',
      dataIndex: 'language',
      width: 100,
      search: false,
    },
    {
      title: '状态',
      dataIndex: 'active',
      width: 100,
      search: false,
      render: (_, record) => (
        <Tag color={record.active ? 'green' : 'default'}>
          {record.active ? '启用' : '停用'}
        </Tag>
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
          <a>查看</a>
          <a>编辑</a>
          <a
            onClick={async () => {
              try {
                await deleteMacro(record.id);
                message.success('宏已删除');
                actionRef.current?.reload();
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
    <ProTable<API.Macro>
      headerTitle="宏模板"
      rowKey="id"
      actionRef={actionRef}
      columns={columns}
      toolBarRender={() => [
        <Button key="create" type="primary" icon={<PlusOutlined />}>
          新建宏
        </Button>,
      ]}
      request={async (params) => {
        try {
          const result = await listMacros();
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
        } catch (error) {
          console.error('获取宏列表失败:', error);
          return { data: [], total: 0, success: true };
        }
      }}
      pagination={{ defaultPageSize: 20 }}
    />
  );
};

export default MacroPage;
