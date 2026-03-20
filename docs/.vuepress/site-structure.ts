const implementationPages = [
  '/implementation/01-platform-and-runtime',
  '/implementation/02-ai-and-knowledge',
  '/implementation/03-business-modules',
  '/implementation/04-sdk-and-channel-adapters',
  '/implementation/05-engineering-hardening',
  '/implementation/06-voice-and-protocol-expansion',
  '/implementation/07-sdk-multi-surface',
  '/implementation/08-ai-provider-expansion',
];

const guidePages = [
  '/ARCHITECTURE',
  '/WEKNORA_INTEGRATION',
  '/CI_SELF_HOSTED',
  '/release-versioning',
  '/testing-pyramid',
  '/MERMAID_COMPATIBILITY',
];

export const docsNavbar = [
  { text: '首页', link: '/' },
  { text: '架构', link: '/ARCHITECTURE' },
  { text: '实施计划', link: '/implementation/' },
  {
    text: '专题',
    children: [
      { text: 'WeKnora 集成', link: '/WEKNORA_INTEGRATION' },
      { text: 'CI / Runner', link: '/CI_SELF_HOSTED' },
      { text: '版本发布', link: '/release-versioning' },
      { text: '测试金字塔', link: '/testing-pyramid' },
      { text: 'Mermaid 兼容性', link: '/MERMAID_COMPATIBILITY' },
    ],
  },
];

export const docsSidebar = {
  '/implementation/': [
    {
      text: '实施计划',
      children: ['/implementation/', ...implementationPages],
    },
  ],
  '/': [
    {
      text: '站点导航',
      children: ['/', ...guidePages],
    },
    {
      text: '实施计划',
      children: ['/implementation/', ...implementationPages],
    },
  ],
};
