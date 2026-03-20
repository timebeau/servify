import { defaultTheme } from '@vuepress/theme-default';
import { defineUserConfig } from 'vuepress';
import { viteBundler } from '@vuepress/bundler-vite';
import { docsNavbar, docsSidebar } from './site-structure';

export default defineUserConfig({
  lang: 'zh-CN',
  title: 'Servify Docs',
  description: 'Servify architecture, implementation backlogs, and integration guides',
  bundler: viteBundler(),
  head: [
    ['meta', { name: 'theme-color', content: '#0f172a' }],
  ],
  theme: defaultTheme({
    logo: '/icon.png',
    navbar: docsNavbar,
    sidebar: docsSidebar,
    editLink: false,
    contributors: false,
    lastUpdated: false,
  }),
});
