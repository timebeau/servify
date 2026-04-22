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

  const transcriptColumns: ProColumns<API.VoiceTranscript>[] = [
    { title: '协议ID', dataIndex: 'call_id', width: 120, search: true },
    {
      title: '内容',
      dataIndex: 'content',
      ellipsis: true,
      search: false,
    },
    {
      title: '语言',
      dataIndex: 'language',
      width: 120,
      search: false,
    },
    {
      title: '状态',
      dataIndex: 'finalized',
      width: 100,
      search: false,
      render: (_, record) => <Tag color={record.finalized ? 'green' : 'default'}>{record.finalized ? '已定稿' : '处理中'}</Tag>,
    },
    {
      title: '创建时间',
      dataIndex: 'appended_at',
      valueType: 'dateTime',
      width: 180,
      search: false,
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
              data: result.data,
              total: result.total,
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

      <ProTable<API.VoiceTranscript>
        headerTitle="语音转写记录"
        rowKey={(record) => `${record.call_id}-${record.appended_at}-${record.content}`}
        columns={transcriptColumns}
        style={{ marginTop: 16 }}
        request={async (params) => {
          try {
            const result = await listTranscripts({
              call_id: typeof params.call_id === 'string' ? params.call_id : undefined,
              page: params.current,
              page_size: params.pageSize,
            });
            return {
              data: result.data,
              total: result.total,
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
