import type {
  NewsArticle,
  NewsArticleListParams,
  NewsArticleListResult,
  NewsCategory,
  NewsAdvertisement,
  NewsCommentListResult,
} from '../types/news'
import { request } from '@/utils/request'
import { resolvePlatformSource } from '@/composables/constants/platform'
import { tokenStorage } from '@/utils/storage'

export const getNewsCategories = () =>
  request<NewsCategory[]>({
    url: '/api/plugin/news/categories',
    method: 'GET',
    auth: false,
    loading: 'none',
  })

export const getNewsAdvertisements = (position: 1 | 2) =>
  request<NewsAdvertisement[]>({
    url: '/api/plugin/news/advertisements',
    method: 'GET',
    data: { position },
    auth: false,
    loading: 'none',
  })

const normalizeArticleListParams = (params: NewsArticleListParams): NewsArticleListParams => ({
  ...(params.category_slug ? { category_slug: params.category_slug } : {}),
  ...(params.keyword ? { keyword: params.keyword } : {}),
  page: params.page,
  page_size: params.page_size,
})

export const getNewsArticles = (params: NewsArticleListParams) =>
  request<NewsArticleListResult, NewsArticleListParams>({
    url: '/api/plugin/news/articles',
    method: 'GET',
    data: normalizeArticleListParams(params),
    auth: false,
    loading: 'none',
  })

export const getNewsArticle = (slug: string) =>
  request<NewsArticle, { slug: string }>({
    url: '/api/plugin/news/article/detail',
    method: 'GET',
    data: { slug },
    auth: false,
    loading: 'none',
  })

export const getNewsComments = (articleUID: string, page = 1, pageSize = 10) =>
  request<NewsCommentListResult>({ url: '/api/plugin/news/article/comments', method: 'GET', data: { article_uid: articleUID, page, page_size: pageSize }, auth: false, loading: 'none' })

export const getNewsCommentReplies = (commentUID: string, page = 1, pageSize = 5) =>
  request<NewsCommentListResult>({ url: '/api/plugin/news/comment/replies', method: 'GET', data: { comment_uid: commentUID, page, page_size: pageSize }, auth: false, loading: 'none' })

export const createNewsComment = (articleUID: string, data: { content: string; images: string[]; request_uid: string }) =>
  request<{ uid: string; status: number }>({ url: '/api/auth/plugin/news/article/comment', method: 'POST', data: { article_uid: articleUID, ...data } })

export const replyNewsComment = (commentUID: string, data: { content: string; images: string[]; request_uid: string }) =>
  request<{ uid: string; status: number }>({ url: '/api/auth/plugin/news/comment/reply', method: 'POST', data: { comment_uid: commentUID, ...data } })

export interface NewsInteractionToggleResult {
  active: boolean
  count: number
}

export const toggleNewsArticleLike = (articleUID: string) =>
  request<NewsInteractionToggleResult>({ url: '/api/auth/plugin/news/article/like', method: 'POST', data: { article_uid: articleUID } })

export const toggleNewsArticleCollection = (articleUID: string) =>
  request<NewsInteractionToggleResult>({ url: '/api/auth/plugin/news/article/collection', method: 'POST', data: { article_uid: articleUID } })

export const toggleNewsCommentLike = (commentUID: string) =>
  request<NewsInteractionToggleResult>({ url: '/api/auth/plugin/news/comment/like', method: 'POST', data: { comment_uid: commentUID } })

export const getNewsInteractionState = (articleUID: string) =>
  request<{ liked: boolean; collected: boolean }>({ url: '/api/auth/plugin/news/article/interaction', method: 'GET', data: { article_uid: articleUID }, loading: 'none' })

export const recordNewsShare = (articleUID: string, shareChannel: number, requestUID: string) =>
  request<null>({ url: '/api/auth/plugin/news/article/share', method: 'POST', data: { article_uid: articleUID, request_uid: requestUID, share_channel: shareChannel }, loading: 'none' })

export const createNewsRequestUID = () => `news_${Date.now()}_${Math.random().toString(36).slice(2, 12)}`

export const uploadNewsCommentImage = (filePath: string) => new Promise<string>((resolve, reject) => {
  const baseURL = import.meta.env.VITE_API_BASE_URL
  const appVersion = import.meta.env.VITE_APP_VERSION
  const token = tokenStorage.get()
  uni.uploadFile({
    url: `${baseURL}/api/auth/upload/file`,
    filePath,
    name: 'file',
    header: {
      'X-Platform-Source': String(resolvePlatformSource()),
      'X-App-Version': appVersion,
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
    },
    success: (response) => {
      try {
        const payload = JSON.parse(response.data) as { code: number; msg: string; data?: { relative_path: string } }
        if (response.statusCode >= 200 && response.statusCode < 300 && payload.code === 0 && payload.data?.relative_path) {
          resolve(payload.data.relative_path)
          return
        }
        reject(new Error(payload.msg || '图片上传失败'))
      } catch {
        reject(new Error('图片上传响应异常'))
      }
    },
    fail: error => reject(new Error(error.errMsg || '图片上传失败')),
  })
})
