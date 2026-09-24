import { request } from '@/api/http'
import type { PageResult } from '@/types/common'
import type { ArticleItem, ArticleListReq, ArticleSaveReq, CategoryItem, CategoryListReq, CategorySaveReq } from '../types/news'

export const getCategoryList = (params: CategoryListReq) => request<PageResult<CategoryItem>>({ url:'/admin/plugin/news/category/list', method:'get', params })
export const saveCategory = (data: CategorySaveReq) => request<CategoryItem>({ url:'/admin/plugin/news/category/save', method:'post', data })
export const deleteCategory = (id: number) => request<null>({ url:'/admin/plugin/news/category/delete', method:'post', data:{ id } })
export const getArticleList = (params: ArticleListReq) => request<PageResult<ArticleItem>>({ url:'/admin/plugin/news/article/list', method:'get', params })
export const getArticleDetail = (id: number) => request<ArticleItem>({ url:'/admin/plugin/news/article/detail', method:'get', params:{ id } })
export const saveArticle = (data: ArticleSaveReq) => request<ArticleItem>({ url:'/admin/plugin/news/article/save', method:'post', data })
export const updateArticleStatus = (id: number, status: number) => request<null>({ url:'/admin/plugin/news/article/status', method:'post', data:{ id,status } })
export const deleteArticle = (id: number) => request<null>({ url:'/admin/plugin/news/article/delete', method:'post', data:{ id } })
