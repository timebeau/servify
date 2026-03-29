import React from 'react';
import { LoginForm, ProFormText } from '@ant-design/pro-components';
import { LockOutlined, UserOutlined } from '@ant-design/icons';
import { message } from 'antd';
import { navigateTo } from '@/lib/navigation';
import { setToken, parseJwtPayload, setUserInfo } from '@/utils/auth';

const LoginPage: React.FC = () => {
  const handleSubmit = async (values: { username: string; password: string }) => {
    try {
      const resp = await fetch('/api/v1/auth/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(values),
      });

      if (!resp.ok) {
        const err = await resp.json().catch(() => ({ message: '登录失败' }));
        message.error(err.message || '登录失败');
        return;
      }

      const data = await resp.json();
      const token = data.token || data.data?.token;
      if (!token) {
        message.error('服务端未返回有效 Token');
        return;
      }

      setToken(token);
      const user = parseJwtPayload(token);
      if (user) setUserInfo(user);

      message.success('登录成功');
      navigateTo('/dashboard');
    } catch {
      message.error('网络错误，请检查后端服务是否启动');
    }
  };

  return (
    <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100vh' }}>
      <LoginForm
        title="Servify"
        subTitle="智能客服管理后台"
        onFinish={handleSubmit}
      >
        <ProFormText
          name="username"
          fieldProps={{ size: 'large', prefix: <UserOutlined /> }}
          placeholder="用户名"
          rules={[{ required: true, message: '请输入用户名' }]}
        />
        <ProFormText.Password
          name="password"
          fieldProps={{ size: 'large', prefix: <LockOutlined /> }}
          placeholder="密码"
          rules={[{ required: true, message: '请输入密码' }]}
        />
      </LoginForm>
    </div>
  );
};

export default LoginPage;
