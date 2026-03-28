import React from 'react';
import { ProTable } from '@ant-design/pro-components';
import type { ProColumns } from '@ant-design/pro-components';
import { Button, Space, message } from 'antd';
import { PlusOutlined } from '@ant-design/icons';
import { listMacros, deleteMacro } from '@/services/macro';

const MacroPage: React.FC = () => {
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
      title: '分类',
      dataIndex: 'category',
      width: 120,
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
          <a>查看</a>
          <a>编辑</a>
          <a
            onClick={async () => {
              try {
                await deleteMacro(record.id);
                message.success('宏已删除');
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
      columns={columns}
      toolBarRender={() => [
        <Button key="create" type="primary" icon={<PlusOutlined />}>
          新建宏
        </Button>,
      ]}
      request={async (params) => {
        try {
          const result = await listMacros({
            page: params.current,
            page_size: params.pageSize,
            category: params.category,
          });
          return {
            data: result?.data || [],
            total: result?.total || 0,
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
