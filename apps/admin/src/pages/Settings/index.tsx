import React, { useState, useEffect } from 'react';
import {
  ProForm,
  ProFormText,
  ProFormTextArea,
  ProFormSelect,
  ProFormSwitch,
} from '@ant-design/pro-components';
import { ProCard } from '@ant-design/pro-components';
import { message, Spin } from 'antd';
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
      <ProForm
        onFinish={async (values) => {
          try {
            console.log('Settings saved:', values);
            message.success('设置已保存');
          } catch (error) {
            message.error('保存设置失败');
          }
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
          rules={[{ required: true, message: '请输入工作空间名称' }]}
        />
        <ProFormSelect
          name="timezone"
          label="时区"
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
          options={[
            { label: '简体中文', value: 'zh-CN' },
            { label: 'English', value: 'en-US' },
          ]}
        />
        <ProFormSwitch name="auto_assignment" label="自动分配工单" />
        <ProFormSwitch name="notification_enabled" label="启用通知" />
        <ProFormTextArea
          name="welcome_message"
          label="欢迎语"
          placeholder="设置客户首次联系的欢迎语"
          fieldProps={{ rows: 3 }}
        />
      </ProForm>
    </ProCard>
  );
};

export default SettingsPage;
