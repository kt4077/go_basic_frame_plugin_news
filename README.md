# 新闻资讯插件

面向 Go Basic Frame `0.0.1` 的官方插件包示例，当前插件版本 `0.1.4`，提供新闻分类、文章发布、广告配置、评论、点赞、收藏、分享以及用户端新闻阅读能力。

## 功能

- 新闻分类分页查询、新增、修改和安全删除；
- 新闻文章分页查询、新增、修改、软删除、发布和下线；
- 新增与修改使用独立页面，正文使用 WangEditor 富文本编辑器；
- 管理端支持复用已有资讯，并以新的业务标识创建独立资讯；
- 支持文章启用状态、虚拟阅读数量，并分别保留真实与展示阅读数量；
- 支持文章作者、来源平台名称和来源外链，用户端通过受控 WebView 打开外链；
- 封面仅保存相对路径，接口根据宿主当前存储配置返回完整地址；
- 用户端提供启用分类、已发布文章列表及文章详情接口；
- 用户端提供独立 uni-app 子包，包含分类筛选、标题搜索、分页列表和富文本详情；
- 详情浏览量使用数据库原子自增，避免并发覆盖；
- 管理端路由接入宿主鉴权、菜单权限和操作日志体系。
- 支持评论层级分页、显示/关闭审核、文章与评论点赞、取消点赞和分享记录；
- 对外业务关联使用唯一 `uid`，用户关联使用 `sys_member.sn`，不暴露数据库主键和会员内部标识；
- 分类和文章在新增时自动生成可编辑标识，创建后锁定标识以保持访问地址稳定。

### 复用资讯

管理端文章列表的“复用资讯”图标会打开独立新增页面，不会直接修改原文章。页面预填分类、标题、摘要、作者、来源、封面、正文、启用状态和排序，并在标题后增加“（复用）”提示。

复用时不会继承数据库主键、文章 UID、发布状态、真实阅读量、虚拟阅读量、点赞、收藏、分享、评论等业务数据。文章标识会重新生成约 11 位的短值，格式为 `n + 6 位时间片 + 4 位随机值`，例如 `n4k8x2a7b3c`；保存前允许编辑，保存后按照现有规则锁定。

## 目录

```text
server_api/internal/plugins/news/   后端插件源码
admin_client/src/plugins/news/      管理端插件页面
user_client/src/plugins/news/       用户端 uni-app 插件子包
database/migrations/                版本化数据库迁移
docs/api/v0.1.4.md                  当前完整接口文档
docs/database/v0.1.4.md             当前数据库与引用说明
docs/updates/v0.1.4.md              当前安装、兼容和更新说明
plugin.json                         插件清单、迁移与菜单声明
```

## 环境要求

- Go Basic Frame 核心版本：`0.0.1`
- Go：以宿主项目 `go.mod` 为准
- Node.js / pnpm：以宿主管理端项目为准
- 用户端：uni-app、Vue 3、TypeScript、Wot Design Uni，以宿主依赖为准
- 管理端依赖：`@wangeditor/editor`、`@wangeditor/editor-for-vue`
- MySQL：8.0 或兼容版本

## 制作发行包

发行 ZIP 的根目录必须直接包含 `plugin.json`，不能额外嵌套仓库目录：

```bash
zip -r go_basic_frame_plugin_news-v0.1.4.zip \
  plugin.json server_api admin_client user_client database README.md CHANGELOG.md docs
```

在核心后端目录验证并安装：

```bash
go run . plugin validate ../go_basic_frame_plugin_news/go_basic_frame_plugin_news-v0.1.4.zip
go run . plugin install ../go_basic_frame_plugin_news/go_basic_frame_plugin_news-v0.1.4.zip \
  --server-root . --admin-root ../admin_client --user-root ../user_client \
  --config ./config.yaml --apply-database
```

管理端需要提供 `@wangeditor/editor` 和 `@wangeditor/editor-for-vue` 依赖。安装器会同时安装后端、管理端和用户端子包；`user_client` 在构建时自动发现插件清单和子包目录。安装成功后需重新构建用户端与管理端，并重启管理端 API、用户端 API。只有后续手工修复后端插件源码或恢复安装前代码时，才需要单独执行 `plugin generate`。

微信小程序新闻详情页支持原生模板信息流广告。请在微信公众平台创建模板广告位，并在宿主 `user_client` 对应环境文件中配置广告位 ID；留空时广告组件不会渲染：

```dotenv
VITE_WECHAT_NEWS_FEED_AD_UNIT_ID=adunit-xxxxxxxxxxxxxxxx
```

该广告仅编译到微信小程序端，加载失败时自动隐藏，不会在 H5、App 或其他小程序平台留下空白占位。

如果数据库已经安装原 `1.1.0` 至 `1.1.4`，必须先备份数据库、停用插件并重启服务，再使用 `plugin upgrade` 和显式的 `--allow-version-rebase` 参数。安装器会检查原 `1.0.0`、`1.1.0` 迁移均已成功，不重复执行 DDL，并保留业务数据。完整命令见 [v0.0.1 更新说明](./docs/updates/v0.0.1.md)。

## 接口

管理端接口均位于 `/admin/plugin/news`，由宿主权限中间件保护。公开接口如下：

完整的参数、响应、鉴权、安全限制与调用示例参见 [API 接口文档](./docs/api/v0.1.4.md)。

业务表、内外部引用、文件归属及 Remove/Purge 影响参见 [数据库说明](./docs/database/v0.1.4.md)。

| 方法 | 地址 | 说明 |
| --- | --- | --- |
| GET | `/api/plugin/news/categories` | 启用分类列表 |
| GET | `/api/plugin/news/articles` | 已发布文章分页列表 |
| GET | `/api/plugin/news/article/detail?slug=xxx` | 已发布文章详情 |
| GET | `/api/plugin/news/article/comments?article_uid=xxx` | 一级评论分页 |
| GET | `/api/plugin/news/comment/replies?comment_uid=xxx` | 评论回复分页 |

评论、回复、点赞、取消点赞、分享与互动状态接口需要用户登录，详见当前 API 文档。

列表参数支持 `category_slug`、`keyword`、`page` 和 `page_size`。

## 开发约束

- 不直接改写宿主核心表，菜单和插件记录由安装器维护；
- 不使用固定菜单 ID 或预置自增业务 ID；
- 新版本只新增迁移文件，禁止修改已发布迁移；
- 每次发布必须按 SemVer 提升版本号：不兼容变更提升主版本，新功能提升次版本，兼容修复或UI优化提升修订版本；禁止覆盖同版本发行包；
- Model 不定义关联，预加载关系只在 Resp 中声明；
- Param、Resp、Model 字段必须保留 Tag 与注释；
- 管理端继续使用 TypeScript、ES6 和箭头函数，并复用宿主 UI 变量与组件。
- 管理端列表页与新增/修改表单组件分离，表单回填、校验和保存逻辑由独立 `*Form.vue` 维护。
- 管理端新增依赖必须在 README 中声明，宿主构建前必须完成依赖安装。
- 用户端插件必须位于 `user_client/src/plugins/{plugin_id}` 独立子包；入口只允许由首页菜单等业务入口跳转，不得加入或修改自定义 TabBar 页面。
- 用户端页面复用宿主请求、导航、主题、交互和空状态能力，使用 Flexbox，兼容亮色、暗色及声明支持的平台。
- 每次发布必须新增 `docs/updates/v{version}.md`，说明变化、迁移、升级步骤和回滚方式。
- 每次发布必须新增完整的 `docs/api/v{version}.md` 和 `docs/database/v{version}.md`；即使接口或数据库没有变化，也要明确记录无变化，禁止覆盖历史版本文档。

## 相关项目

| 官网地址 | 管理端地址 | 接口端地址 |
| --- | --- | --- |
| [项目官网](https://www.tutudati.com/) | [管理端仓库](https://gitee.com/open-source-project-open/go_basic_frame_admin) | [接口端仓库](https://gitee.com/open-source-project-open/go_basic_frame_api) |

## License

请在正式发布前根据项目授权策略补充许可证文件。
