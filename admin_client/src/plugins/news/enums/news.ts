export enum CategoryStatus { Enabled = 1, Disabled = 2 }
export const CategoryStatusLabels: Record<CategoryStatus, string> = { [CategoryStatus.Enabled]: '启用', [CategoryStatus.Disabled]: '停用' }

export enum ArticleStatus { Draft = 1, Published = 2, Offline = 3 }
export const ArticleStatusLabels: Record<ArticleStatus, string> = { [ArticleStatus.Draft]: '草稿', [ArticleStatus.Published]: '已发布', [ArticleStatus.Offline]: '已下线' }

export enum ArticleEnableStatus { Enabled = 1, Disabled = 2 }
export const ArticleEnableStatusLabels: Record<ArticleEnableStatus, string> = { [ArticleEnableStatus.Enabled]: '启用', [ArticleEnableStatus.Disabled]: '禁用' }
