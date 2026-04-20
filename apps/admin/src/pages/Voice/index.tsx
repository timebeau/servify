import React from 'react';
import { ProTable } from '@ant-design/pro-components';
import type { ProColumns } from '@ant-design/pro-components';
import { Tag } from 'antd';
import { listProtocols, listTranscripts } from '@/services/voice';

const VoicePage: React.FC = () => {
  const protocolColumns: ProColumns<API.VoiceProtocol>[] = [
    { title: 'ID', dataIndex: 'id', width: 80 },
    { title: '类型', dataIndex: 'type', width: 120 },
    {
      title: '状态',
      dataIndex: 'status',
      width: 100,
      render: (_, record) => (
        <Tag color={record.status === 'active' || record.status === 'configured' ? 'green' : 'default'}>
          {record.status === 'active' ? '已连接' : record.status === 'configured' ? '已配置' : '未连接'}
        </Tag>
      ),
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      valueType: 'dateTime',
      width: 180,
    },
  ];

  const transcriptColumns: ProColumns<any>[] = [
    { title: 'ID', dataIndex: 'id', width: 80 },
    { title: '协议ID', dataIndex: 'call_id', width: 120 },
    {
      title: '内容',
      dataIndex: 'content',
      ellipsis: true,
    },
    {
      title: '创建时间',
      dataIndex: 'appended_at',
      valueType: 'dateTime',
      width: 180,
    },
  ];

  return (
    <div>
      <ProTable<API.VoiceProtocol>
        headerTitle="语音协议"
        rowKey="id"
        columns={protocolColumns}
        request={async (params) => {
          try {
            const result = await listProtocols({
              page: params.current,
              page_size: params.pageSize,
            });
            return {
              data: result?.data?.data || result?.data || [],
              total: result?.data?.total || (result?.data?.data || result?.data || []).length,
              success: true,
            };
          } catch (error) {
            console.error('获取语音协议失败:', error);
            return { data: [], total: 0, success: true };
          }
        }}
        search={false}
        pagination={{ defaultPageSize: 10 }}
      />

      <ProTable<any>
        headerTitle="语音转写记录"
        rowKey="id"
        columns={transcriptColumns}
        style={{ marginTop: 16 }}
        request={async (params) => {
          try {
            const result = await listTranscripts({
              page: params.current,
              page_size: params.pageSize,
            });
            return {
              data: result?.data || [],
              total: result?.total || result?.data?.length || 0,
              success: true,
            };
          } catch (error) {
            console.error('获取转写记录失败:', error);
            return { data: [], total: 0, success: true };
          }
        }}
        search={{ filterType: 'light' }}
        pagination={{ defaultPageSize: 10 }}
      />
    </div>
  );
};

export default VoicePage;
