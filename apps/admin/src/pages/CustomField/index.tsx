import React from 'react';
import { ProTable } from '@ant-design/pro-components';
import type { ProColumns } from '@ant-design/pro-components';
import { Button, Space, Tag, message } from 'antd';
import { PlusOutlined } from '@ant-design/icons';
import { listCustomFields, deleteCustomField } from '@/services/customField';
import { getErrorMessage } from '@/utils/error';

const CustomFieldPage: React.FC = () => {
  const columns: ProColumns<API.CustomField>[] = [
    {
      title: 'ID',
      dataIndex: 'id',
      width: 80,
    },
    {
      title: '字段名',
      dataIndex: 'name',
      search: true,
    },
    {
      title: '字段标识',
      dataIndex: 'key',
    },
    {
      title: '字段类型',
      dataIndex: 'field_type',
      width: 120,
      render: (_, record) => <Tag>{record.field_type}</Tag>,
    },
    {
      title: '适用资源',
      dataIndex: 'entity_type',
      width: 120,
      render: (_, record) => <Tag color="blue">{record.entity_type}</Tag>,
    },
    {
      title: '必填',
      dataIndex: 'required',
      width: 80,
      render: (_, record) => (
        <Tag color={record.required ? 'red' : 'default'}>
          {record.required ? '是' : '否'}
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
                await deleteCustomField(record.id);
                message.success('字段已删除');
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
    <ProTable<API.CustomField>
      headerTitle="自定义字段"
      rowKey="id"
      columns={columns}
      toolBarRender={() => [
        <Button key="create" type="primary" icon={<PlusOutlined />}>
          新建字段
        </Button>,
      ]}
      request={async (params) => {
        try {
          const result = await listCustomFields({
            page: params.current,
            page_size: params.pageSize,
            entity_type: params.entity_type,
          });
          return {
            data: result.data,
            total: result.total,
            success: true,
          };
        } catch (error: unknown) {
          console.error('获取自定义字段失败:', error);
          return { data: [], total: 0, success: true };
        }
      }}
      pagination={{ defaultPageSize: 20 }}
    />
  );
};

export default CustomFieldPage;
