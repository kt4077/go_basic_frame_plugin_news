# 新闻资讯插件

面向 Go Basic Frame `0.0.2` 的官方插件包示例，提供新闻分类、文章草稿、发布/下线、封面上传以及用户端新闻查询能力。

## 功能

- 新闻分类分页查询、新增、修改和安全删除；
- 新闻文章分页查询、新增、修改、软删除、发布和下线；
- 封面仅保存相对路径，接口根据宿主当前存储配置返回完整地址；
- 用户端提供启用分类、已发布文章列表及文章详情接口；
- 详情浏览量使用数据库原子自增，避免并发覆盖；
- 管理端路由接入宿主鉴权、菜单权限和操作日志体系。

## 目录

```text
server_api/internal/plugins/news/   后端插件源码
admin_client/src/plugins/news/      管理端插件页面
database/migrations/                版本化数据库迁移
docs/API.md                         管理端与用户端接口文档
plugin.json                         插件清单、迁移与菜单声明
```

## 环境要求

- Go Basic Frame 核心版本：`0.0.2`
- Go：以宿主项目 `go.mod` 为准
- Node.js / pnpm：以宿主管理端项目为准
- MySQL：8.0 或兼容版本

## 制作发行包

发行 ZIP 的根目录必须直接包含 `plugin.json`，不能额外嵌套仓库目录：

```bash
zip -r go_basic_frame_plugin_news-v1.0.0.zip \
  plugin.json server_api admin_client database README.md CHANGELOG.md
```

在核心后端目录验证并安装：

```bash
go run . plugin validate ../go_basic_frame_plugin_news/go_basic_frame_plugin_news-v1.0.0.zip
go run . plugin install ../go_basic_frame_plugin_news/go_basic_frame_plugin_news-v1.0.0.zip \
  --server-root . --admin-root ../admin_client
go run . plugin generate --server-root .
```

安装命令会执行迁移、按稳定业务键创建菜单，并把前后端源码安装到宿主目录。安装后需要重新构建并重启管理端 API、用户端 API 和管理端前端，然后在“插件管理”中启用插件。

## 接口

管理端接口均位于 `/admin/plugin/news`，由宿主权限中间件保护。公开接口如下：

完整的参数、响应、鉴权、错误场景与调用示例参见 [API 接口文档](./docs/API.md)。

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
- Model 不定义关联，预加载关系只在 Resp 中声明；
- Param、Resp、Model 字段必须保留 Tag 与注释；
- 管理端继续使用 TypeScript、ES6 和箭头函数，并复用宿主 UI 变量与组件。
- 管理端列表页与新增/修改表单组件分离，表单回填、校验和保存逻辑由独立 `*Form.vue` 维护。

## 相关项目

| 官网地址 | 管理端地址 | 接口端地址 |
| --- | --- | --- |
| [项目官网](https://www.tutudati.com/) | [管理端仓库](https://gitee.com/open-source-project-open/go_basic_frame_admin) | [接口端仓库](https://gitee.com/open-source-project-open/go_basic_frame_api) |

## License

请在正式发布前根据项目授权策略补充许可证文件。
