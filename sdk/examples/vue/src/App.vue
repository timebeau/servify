<script setup lang="ts">
import { computed, ref, watch } from 'vue';
import {
  useServifyReady,
  useChat,
  useAI,
  useSatisfaction,
  useRemoteAssist,
} from '@servify/vue';

const { isConnected } = useServifyReady();
const {
  session,
  messages,
  agent,
  isLoading,
  error,
  isAgentTyping,
  startChat,
  sendMessage,
  endChat,
  uploadFile,
} = useChat();
const { askAI } = useAI();
const { submitRating } = useSatisfaction();
const {
  state: remoteAssistState,
  isActive: remoteAssistActive,
  error: remoteAssistError,
  remoteStream,
  startRemoteAssist,
  endRemoteAssist,
} = useRemoteAssist();

const messageText = ref('');
const remoteVideo = ref<HTMLVideoElement | null>(null);

watch(remoteStream, (stream) => {
  if (remoteVideo.value) {
    remoteVideo.value.srcObject = stream;
  }
});

const remoteMediaLabel = computed(() => (remoteStream.value ? '已接入' : '未接入'));

async function handleStartChat() {
  try {
    await startChat({
      priority: 'normal',
      message: '您好，我需要帮助',
    });
  } catch (err) {
    alert('开始聊天失败: ' + (err as Error).message);
  }
}

async function handleSendMessage() {
  if (!session.value || !messageText.value.trim()) {
    return;
  }

  try {
    await sendMessage(messageText.value);
    messageText.value = '';
  } catch (err) {
    alert('发送消息失败: ' + (err as Error).message);
  }
}

async function handleEndChat() {
  try {
    await endRemoteAssist();
    await endChat();
  } catch (err) {
    alert('结束聊天失败: ' + (err as Error).message);
  }
}

async function handleAskAI() {
  const question = window.prompt('请输入您的问题：');
  if (!question) {
    return;
  }

  try {
    const response = await askAI(question);
    window.alert(`AI 回答(${(response.confidence * 100).toFixed(1)}%): ${response.answer}`);
  } catch (err) {
    alert('AI 问答失败: ' + (err as Error).message);
  }
}

async function handleRating() {
  const raw = window.prompt('请为服务评分 (1-5)：');
  const rating = Number(raw);
  if (!Number.isInteger(rating) || rating < 1 || rating > 5) {
    alert('请输入有效评分');
    return;
  }

  try {
    await submitRating({ rating });
    alert('感谢您的评价');
  } catch (err) {
    alert('提交评价失败: ' + (err as Error).message);
  }
}

async function handleStartRemoteAssist() {
  if (!session.value) {
    alert('请先开始聊天，再发起远程协助');
    return;
  }

  try {
    await startRemoteAssist({ captureScreen: true, audio: false });
  } catch (err) {
    alert('发起远程协助失败: ' + (err as Error).message);
  }
}

async function handleStopRemoteAssist() {
  try {
    await endRemoteAssist();
  } catch (err) {
    alert('结束远程协助失败: ' + (err as Error).message);
  }
}

async function handleFileChange(event: Event) {
  const input = event.target as HTMLInputElement;
  const file = input.files?.[0];
  if (!file || !session.value) {
    return;
  }

  try {
    await uploadFile(file);
    alert(`文件已上传: ${file.name}`);
  } catch (err) {
    alert('文件上传失败: ' + (err as Error).message);
  } finally {
    input.value = '';
  }
}
</script>

<template>
  <div class="container">
    <header class="header">
      <h1>Servify 客服系统</h1>
      <p>Vue SDK 示例</p>
    </header>

    <section class="status-bar">
      <div class="status-left">
        <span class="status-dot" :class="isConnected ? 'connected' : 'disconnected'"></span>
        <span>{{ isConnected ? '已连接' : '已断开连接' }}</span>
      </div>
      <div v-if="agent">客服代理：<strong>{{ agent.name }}</strong></div>
    </section>

    <section class="chat-area">
      <div class="message system">欢迎使用 Servify 客服系统！请开始聊天。</div>

      <div v-for="message in messages" :key="message.id" class="message" :class="message.sender_type">
        {{ message.content }}
      </div>

      <div v-if="error" class="message system">错误: {{ error.message }}</div>
      <div class="message system">远程协助状态: {{ remoteAssistState }}</div>
      <div v-if="remoteAssistError" class="message system">
        远程协助错误: {{ remoteAssistError.message }}
      </div>
      <div class="message system">远端媒体: {{ remoteMediaLabel }}</div>
    </section>

    <section class="remote-preview">
      <div class="preview-title">远程协助媒体预览</div>
      <video
        v-show="remoteStream"
        ref="remoteVideo"
        autoplay
        playsinline
        muted
        class="preview-video"
      />
      <div v-if="!remoteStream" class="preview-empty">尚未收到远端媒体流</div>
    </section>

    <div v-if="isAgentTyping" class="typing-indicator">客服正在输入...</div>

    <section class="controls">
      <div class="input-group">
        <input
          v-model="messageText"
          type="text"
          maxlength="1000"
          placeholder="输入您的消息..."
          :disabled="!session"
          @keydown.enter.prevent="handleSendMessage"
        />
        <label class="btn secondary">
          📎
          <input type="file" class="hidden-input" @change="handleFileChange" />
        </label>
        <button class="btn primary" :disabled="!session || !messageText.trim() || isLoading" @click="handleSendMessage">
          发送
        </button>
      </div>

      <div class="button-row">
        <button class="btn success" :disabled="!!session || isLoading" @click="handleStartChat">开始聊天</button>
        <button class="btn danger" :disabled="!session || isLoading" @click="handleEndChat">结束聊天</button>
        <button class="btn secondary" :disabled="isLoading" @click="handleAskAI">AI 助手</button>
        <button class="btn secondary" :disabled="!session" @click="handleRating">评价服务</button>
        <button class="btn secondary" :disabled="!session || remoteAssistActive" @click="handleStartRemoteAssist">
          开始屏幕协助
        </button>
        <button class="btn danger" :disabled="!remoteAssistActive" @click="handleStopRemoteAssist">
          结束屏幕协助
        </button>
      </div>
    </section>
  </div>
</template>
