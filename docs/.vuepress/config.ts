import { defaultTheme } from '@vuepress/theme-default';
import { defineUserConfig } from 'vuepress';
import { viteBundler } from '@vuepress/bundler-vite';

export default defineUserConfig({
  lang: 'zh-CN',
  title: 'Servify Docs',
  description: 'Servify architecture, implementation backlogs, and integration guides',
  bundler: viteBundler(),
  theme: defaultTheme({
    logo: '/icon.png',
    navbar: [
      { text: '文档首页', link: '/' },
      { text: '实施计划', link: '/implementation/' },
      { text: 'WeKnora', link: '/WEKNORA_INTEGRATION.html' },
      { text: 'CI', link: '/CI_SELF_HOSTED.html' },
    ],
    sidebar: {
      '/': [
        '/',
        '/WEKNORA_INTEGRATION',
        '/CI_SELF_HOSTED',
        '/implementation/',
      ],
      '/implementation/': [
        '/implementation/',
        '/implementation/01-platform-and-runtime',
        '/implementation/02-ai-and-knowledge',
        '/implementation/03-business-modules',
        '/implementation/04-sdk-and-channel-adapters',
        '/implementation/05-engineering-hardening',
        '/implementation/06-voice-and-protocol-expansion',
        '/implementation/07-sdk-multi-surface',
        '/implementation/08-ai-provider-expansion',
      ],
    },
  }),
});
