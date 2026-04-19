import React, { useState, useEffect } from 'react';
import {
  ProForm,
  ProFormText,
  ProFormTextArea,
  ProFormSelect,
  ProFormSwitch,
} from '@ant-design/pro-components';
import { ProCard } from '@ant-design/pro-components';
import { Alert, Spin } from 'antd';
import { getPortalConfig } from '@/services/portalConfig';
import { getWorkspaceOverview } from '@/services/workspace';

const SettingsPage: React.FC = () => {
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchConfig = async () => {
      setLoading(true);
      try {
        const [portalResult, workspaceResult] = await Promise.allSettled([
          getPortalConfig(),
          getWorkspaceOverview(),
        ]);
        // Config loaded - ProForm will use initialValues
        // In a real implementation, you would set form values here
        if (portalResult.status === 'fulfilled' && portalResult.value) {
          // portalConfig available
        }
        if (workspaceResult.status === 'fulfilled' && workspaceResult.value) {
          // workspace config available
        }
      } catch (error) {
        console.error('获取配置失败:', error);
      } finally {
        setLoading(false);
      }
    };
    fetchConfig();
  }, []);

  if (loading) {
    return (
      <div style={{ textAlign: 'center', padding: 80 }}>
        <Spin size="large" tip="加载中..." />
      </div>
    );
  }

  return (
    <ProCard title="系统设置">
      <Alert
        message="设置页面功能未完成"
        description="当前设置页面仅用于演示 UI。实际的配置保存/加载 API 尚未实现，修改设置不会持久化。如需修改配置，请直接编辑 config.yml 文件并重启服务。"
        type="warning"
        showIcon
        style={{ marginBottom: 24 }}
      />
      <ProForm
        onFinish={async (values) => {
          // TODO: Implement actual settings persistence API
          console.log('Settings saved (NOT IMPLEMENTED):', values);
          // This is a placeholder - the real implementation requires:
          // 1. Backend API endpoint(s) for updating workspace/portal config
          // 2. Config write-back with validation and hot-reload
          // 3. Proper error handling and user feedback
        }}
        disabled
        submitter={{
          render: () => null, // Hide submit button since it's not implemented
        }}
        initialValues={{
          workspace_name: '',
          timezone: 'Asia/Shanghai',
          language: 'zh-CN',
          auto_assignment: true,
          notification_enabled: true,
        }}
        style={{ maxWidth: 640 }}
      >
        <ProFormText
          name="workspace_name"
          label="工作空间名称"
          placeholder="请输入工作空间名称"
          disabled
        />
        <ProFormSelect
          name="timezone"
          label="时区"
          disabled
          options={[
            { label: 'UTC+8 北京时间', value: 'Asia/Shanghai' },
            { label: 'UTC+9 东京时间', value: 'Asia/Tokyo' },
            { label: 'UTC+0 GMT', value: 'GMT' },
            { label: 'UTC-5 东部时间', value: 'US/Eastern' },
          ]}
        />
        <ProFormSelect
          name="language"
          label="语言"
          disabled
          options={[
            { label: '简体中文', value: 'zh-CN' },
            { label: 'English', value: 'en-US' },
          ]}
        />
        <ProFormSwitch name="auto_assignment" label="自动分配工单" disabled />
        <ProFormSwitch name="notification_enabled" label="启用通知" disabled />
        <ProFormTextArea
          name="welcome_message"
          label="欢迎语"
          placeholder="设置客户首次联系的欢迎语"
          fieldProps={{ rows: 3 }}
          disabled
        />
      </ProForm>
    </ProCard>
  );
};

export default SettingsPage;
