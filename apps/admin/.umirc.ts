import { defineConfig } from '@umijs/max';
import routes from './config/routes';
import proxy from './config/proxy';
import defaultSettings from './config/defaultSettings';

export default defineConfig({
  antd: {},
  access: {},
  model: {},
  initialState: {},
  request: { dataField: '' },
  layout: {
    title: 'Servify 管理后台',
    locale: false,
    ...defaultSettings,
  },
  routes,
  proxy,
  history: { type: 'browser' },
  jsMinifier: 'terser',
  hash: true,
});
