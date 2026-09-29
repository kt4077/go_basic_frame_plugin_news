import { request } from '@/api/http'
import type { PageResult } from '@/types/common'
import type { AdvertisementConfig, AdvertisementOption, AdvertisementSaveReq, ArticleCommentItem, ArticleCommentListReq, ArticleItem, ArticleListReq, ArticleSaveReq, CategoryItem, CategoryListReq, CategorySaveReq } from '../types/news'

export const getCategoryList = (params: CategoryListReq) => request<PageResult<CategoryItem>>({ url:'/admin/plugin/news/category/list', method:'get', params })
export const saveCategory = (data: CategorySaveReq) => request<CategoryItem>({ url:'/admin/plugin/news/category/save', method:'post', data })
export const deleteCategory = (id: number) => request<null>({ url:'/admin/plugin/news/category/delete', method:'post', data:{ id } })
export const getArticleList = (params: ArticleListReq) => request<PageResult<ArticleItem>>({ url:'/admin/plugin/news/article/list', method:'get', params })
export const getArticleDetail = (id: number) => request<ArticleItem>({ url:'/admin/plugin/news/article/detail', method:'get', params:{ id } })
export const saveArticle = (data: ArticleSaveReq) => request<ArticleItem>({ url:'/admin/plugin/news/article/save', method:'post', data })
export const updateArticleStatus = (id: number, status: number) => request<null>({ url:'/admin/plugin/news/article/status', method:'post', data:{ id,status } })
export const deleteArticle = (id: number) => request<null>({ url:'/admin/plugin/news/article/delete', method:'post', data:{ id } })
export const getArticleCommentList = (params: ArticleCommentListReq) => request<PageResult<ArticleCommentItem>>({ url:'/admin/plugin/news/article/comment/list', method:'get', params })
export const getArticleCommentReplies = (params: ArticleCommentListReq) => request<PageResult<ArticleCommentItem>>({ url:'/admin/plugin/news/article/comment/replies', method:'get', params })
export const updateArticleCommentStatus = (comment_uid: string, status: number) => request<null>({ url:'/admin/plugin/news/article/comment/status', method:'post', data:{ comment_uid, status } })
export const getAdvertisementConfigs = () => request<AdvertisementConfig[]>({ url:'/admin/plugin/news/advertisement/list', method:'get' })
export const getAdvertisementOptions = () => request<AdvertisementOption[]>({ url:'/admin/plugin/news/advertisement/options', method:'get' })
export const saveAdvertisementConfig = (data: AdvertisementSaveReq) => request<null>({ url:'/admin/plugin/news/advertisement/save', method:'post', data })
export const deleteAdvertisementConfig = (id: number) => request<null>({ url:'/admin/plugin/news/advertisement/delete', method:'post', data:{ id } })
