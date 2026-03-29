import { getToken, clearToken, parseJwtPayload, setUserInfo } from './utils/auth';

const loginPath = '/login';

/** 初始状态：获取当前用户信息 */
export async function getInitialState(): Promise<{
  currentUser?: API.CurrentUser;
}> {
  const token = getToken();
  if (!token) return {};

  try {
    const user = parseJwtPayload(token);
    if (user) {
      setUserInfo(user);
      return { currentUser: user };
    }
  } catch {
    // ignore
  }

  return {};
}

/** ProLayout 配置 */
export const layout = () => ({
  logo: 'https://gw.alipayobjects.com/mdn/rms_b5fcc5/afts/img/A*1NHAQYduQiQAAAAAAAAAAAAAARQnAQ',
  menu: { locale: false },
  logout: () => {
    clearToken();
    window.location.replace(loginPath);
  },
});

/** 请求配置：JWT 拦截器 */
export const request = {
  timeout: 30000,
  responseInterceptors: [
    (response: any) => {
      if (response.status === 401) {
        clearToken();
        window.location.replace(loginPath);
      }
      return response;
    },
  ],
};
