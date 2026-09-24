# 新闻资讯插件

面向 Go Basic Frame `0.0.2` 的官方插件包示例，提供新闻分类、文章草稿、发布/下线、封面上传以及用户端新闻查询能力。

## 功能

- 新闻分类分页查询、新增、修改和安全删除；
- 新闻文章分页查询、新增、修改、软删除、发布和下线；
- 新增与修改使用独立页面，正文使用 WangEditor 富文本编辑器；
- 支持文章启用状态、虚拟阅读数量，并分别保留真实与展示阅读数量；
- 封面仅保存相对路径，接口根据宿主当前存储配置返回完整地址；
- 用户端提供启用分类、已发布文章列表及文章详情接口；
- 详情浏览量使用数据库原子自增，避免并发覆盖；
- 管理端路由接入宿主鉴权、菜单权限和操作日志体系。
- 分类和文章在新增时自动生成可编辑标识，创建后锁定标识以保持访问地址稳定。

## 目录

```text
server_api/internal/plugins/news/   后端插件源码
admin_client/src/plugins/news/      管理端插件页面
database/migrations/                版本化数据库迁移
docs/API.md                         当前版本接口文档入口
docs/DATABASE.md                    当前版本数据库说明入口
docs/api/                           按版本保存的完整接口文档
docs/database/                      按版本保存的数据库说明
docs/updates/                       各版本独立更新说明
plugin.json                         插件清单、迁移与菜单声明
```

## 环境要求

- Go Basic Frame 核心版本：`0.0.2`
- Go：以宿主项目 `go.mod` 为准
- Node.js / pnpm：以宿主管理端项目为准
- 管理端依赖：`@wangeditor/editor`、`@wangeditor/editor-for-vue`
- MySQL：8.0 或兼容版本

## 制作发行包

发行 ZIP 的根目录必须直接包含 `plugin.json`，不能额外嵌套仓库目录：

```bash
zip -r go_basic_frame_plugin_news-v1.1.3.zip \
  plugin.json server_api admin_client database README.md CHANGELOG.md docs
```

在核心后端目录验证并安装：

```bash
go run . plugin validate ../go_basic_frame_plugin_news/go_basic_frame_plugin_news-v1.1.3.zip
go run . plugin upgrade ../go_basic_frame_plugin_news/go_basic_frame_plugin_news-v1.1.3.zip \
  --server-root . --admin-root ../admin_client \
  --config ./config.yaml --apply-database
```

升级前先在宿主管理端执行 `pnpm add @wangeditor/editor @wangeditor/editor-for-vue@next`。升级命令需要按宿主说明增加 `--config ./config.yaml --apply-database` 才会执行 v1.1.0 数据库迁移。安装器会自动重新生成插件注册文件；升级成功后直接构建并重启管理端 API、用户端 API 和管理端前端。只有后续手工修复插件源码或恢复安装前代码时，才需要单独执行 `plugin generate`。

## 接口

管理端接口均位于 `/admin/plugin/news`，由宿主权限中间件保护。公开接口如下：

完整的参数、响应、鉴权、错误场景与调用示例参见 [API 接口文档](./docs/API.md)。

业务表、内外部引用、文件归属及 Remove/Purge 影响参见 [数据库说明](./docs/DATABASE.md)。

| 方法 | 地址 | 说明 |
| --- | --- | --- |
| GET | `/api/plugin/news/categories` | 启用分类列表 |
| GET | `/api/plugin/news/articles` | 已发布文章分页列表 |
| GET | `/api/plugin/news/articles/:slug` | 已发布文章详情 |

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
- 每次发布必须新增 `docs/updates/v{version}.md`，说明变化、迁移、升级步骤和回滚方式。
- 每次发布必须新增完整的 `docs/api/v{version}.md` 和 `docs/database/v{version}.md`；即使接口或数据库没有变化，也要明确记录无变化，禁止覆盖历史版本文档。

## 相关项目

| 官网地址 | 管理端地址 | 接口端地址 |
| --- | --- | --- |
| [项目官网](https://www.tutudati.com/) | [管理端仓库](https://gitee.com/open-source-project-open/go_basic_frame_admin) | [接口端仓库](https://gitee.com/open-source-project-open/go_basic_frame_api) |

## License

请在正式发布前根据项目授权策略补充许可证文件。
