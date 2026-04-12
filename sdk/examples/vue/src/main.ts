import { createApp } from 'vue';
import App from './App.vue';
import { ServifyPlugin } from '@servify/vue';
import './style.css';

const app = createApp(App);

app.use(ServifyPlugin, {
  config: {
    apiUrl: 'http://localhost:8080',
    customerName: 'Vue 用户',
    customerEmail: 'vue@example.com',
    debug: true,
  },
  onInitialized: () => {
    console.log('Servify Vue SDK 初始化成功');
  },
  onError: (error: Error) => {
    console.error('Servify Vue SDK 错误:', error);
    alert(`SDK 错误: ${error.message}`);
  },
});

app.mount('#app');
