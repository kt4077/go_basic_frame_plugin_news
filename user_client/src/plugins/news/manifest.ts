import type { UserPluginManifest } from '../types'

const manifest: UserPluginManifest = {
  id: 'news',
  name: '资讯中心',
  version: '0.1.4',
  subpackageRoot: 'plugins/news',
  pages: ['pages/index/index', 'pages/detail/index'],
  entryPage: '/plugins/news/pages/index/index',
  minCoreVersion: '0.0.1',
  enabled: true,
}

export default manifest
