import React from 'react';
import { Result } from 'antd';

const SettingsPage: React.FC = () => {
  return (
    <Result
      status="info"
      title="系统设置"
      subTitle="当前配置通过 config.yml 文件管理。如需修改配置，请编辑 config.yml 文件并重启服务。"
    />
  );
};

export default SettingsPage;
