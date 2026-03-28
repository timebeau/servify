import React from 'react';
import { ProTable } from '@ant-design/pro-components';
import type { ProColumns } from '@ant-design/pro-components';
import { Tag, Button, Space, message } from 'antd';
import { PlusOutlined } from '@ant-design/icons';
import { history } from '@umijs/max';
import { listDocs, deleteDoc } from '@/services/knowledge';

const KnowledgeListPage: React.FC = () => {
  const columns: ProColumns<API.KnowledgeDoc>[] = [
    {
      title: 'ID',
      dataIndex: 'id',
      width: 80,
    },
    {
      title: '标题',
      dataIndex: 'title',
      ellipsis: true,
      search: true,
    },
    {
      title: '分类',
      dataIndex: 'category',
      width: 120,
    },
    {
      title: '状态',
      dataIndex: 'status',
      width: 100,
      render: (_, record) => <Tag>{record.status}</Tag>,
    },
    {
      title: '更新时间',
      dataIndex: 'updated_at',
      valueType: 'dateTime',
      width: 180,
      sorter: true,
    },
    {
      title: '操作',
      valueType: 'option',
      width: 160,
      render: (_, record) => (
        <Space>
          <a onClick={() => history.push(`/knowledge/detail/${record.id}`)}>查看</a>
          <a>编辑</a>
          <a
            onClick={async () => {
              try {
                await deleteDoc(record.id);
                message.success('文档已删除');
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
    <ProTable<API.KnowledgeDoc>
      headerTitle="知识库文档"
      rowKey="id"
      columns={columns}
      toolBarRender={() => [
        <Button key="create" type="primary" icon={<PlusOutlined />}>
          新建文档
        </Button>,
      ]}
      request={async (params) => {
        try {
          const result = await listDocs({
            page: params.current,
            page_size: params.pageSize,
            search: params.title,
            category: params.category,
            status: params.status,
          });
          return {
            data: result?.data || [],
            total: result?.total || 0,
            success: true,
          };
        } catch (error) {
          console.error('获取知识库列表失败:', error);
          return { data: [], total: 0, success: true };
        }
      }}
      pagination={{ defaultPageSize: 20 }}
    />
  );
};

export default KnowledgeListPage;
