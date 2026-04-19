import React from 'react';
import { ProTable } from '@ant-design/pro-components';
import type { ProColumns } from '@ant-design/pro-components';
import { Alert, Tag, Button } from 'antd';
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
        <Tag color={record.status === 'active' ? 'green' : 'default'}>
          {record.status === 'active' ? '已连接' : '未连接'}
        </Tag>
      ),
    },
    {
      title: '开始时间',
      dataIndex: 'started_at',
      valueType: 'dateTime',
      width: 180,
    },
    {
      title: '时长(秒)',
      dataIndex: 'duration',
      width: 100,
    },
    {
      title: '操作',
      valueType: 'option',
      width: 80,
      render: () => <a>配置</a>,
    },
  ];

  const transcriptColumns: ProColumns<any>[] = [
    { title: 'ID', dataIndex: 'id', width: 80 },
    { title: '协议ID', dataIndex: 'protocol_id', width: 120 },
    {
      title: '内容',
      dataIndex: 'content',
      ellipsis: true,
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      valueType: 'dateTime',
      width: 180,
    },
  ];

  return (
    <div>
      <Alert
        message="语音管理页面功能未完成"
        description="当前语音管理页面存在前后端契约不匹配问题：
        1. 协议列表：后端返回结构不支持分页，前端期待分页数据
        2. 转写记录：后端要求 call_id 参数，但前端未发送
        Voice provider 当前仅支持 disabled，无真实录音/转写实现。
        此页面仅用于 UI 演示，不具备实际功能。"
        type="warning"
        showIcon
        style={{ marginBottom: 16 }}
      />
      <ProTable<API.VoiceProtocol>
        headerTitle="语音协议"
        rowKey="id"
        columns={protocolColumns}
        toolBarRender={() => [
          <Button key="add" type="primary" disabled>
            添加协议
          </Button>,
        ]}
        request={async (params) => {
          try {
            const result = await listProtocols({
              page: params.current,
              page_size: params.pageSize,
            });
            // Backend returns { success, data } without pagination total
            // This is a contract mismatch - returning static data for UI demo
            return {
              data: result?.data || [],
              total: (result?.data || []).length, // Fake total for demo
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
            // Backend requires call_id but frontend doesn't send it
            // This will return 400 - contract mismatch
            const result = await listTranscripts({
              page: params.current,
              page_size: params.pageSize,
              // call_id: '' // TODO: Fix backend API to support listing without call_id
            });
            return {
              data: result?.data || [],
              total: result?.total || 0,
              success: true,
            };
          } catch (error) {
            console.error('获取转写记录失败:', error);
            // Return empty instead of showing error
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
