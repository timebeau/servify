import { createApp } from 'vue';
import App from './App.vue';
import { ServifyPlugin } from '@servify/vue';
import './style.css';

const app = createApp(App);

app.use(ServifyPlugin, {
  config: {
    apiUrl: 'http://localhost:8080',
    wsUrl: 'ws://localhost:8080/api/v1/ws',
    customerName: 'Vue User',
    customerEmail: 'vue@example.com',
    debug: true,
  },
  onInitialized: () => {
    console.log('Servify Vue SDK initialized');
  },
  onError: (error: Error) => {
    console.error('Servify Vue SDK error:', error);
    window.alert(`SDK error: ${error.message}`);
  },
});

app.mount('#app');
