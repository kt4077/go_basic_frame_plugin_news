import type { PageQuery } from '@/types/common'

export interface CategoryItem { id:number; name:string; slug:string; sort:number; status:number; remark:string; created_at:string; updated_at:string }
export interface CategoryListReq extends PageQuery { keyword?:string; status?:number }
export interface CategorySaveReq { id:number; name:string; slug:string; sort:number; status:number; remark:string }

export interface ArticleCategoryItem { id:number; name:string; slug:string }
export interface ArticleItem { id:number; uid:string; category_id:number; category:ArticleCategoryItem; title:string; slug:string; summary:string; author:string; source_name:string; source_url:string; cover:string; cover_url:string; content:string; status:number; enable_status:number; sort:number; view_count:number; virtual_view_count:number; total_view_count:number; like_count:number; share_count:number; collection_count:number; comment_count:number; published_at?:string; created_at:string; updated_at:string }
export interface ArticleListReq extends PageQuery { keyword?:string; category_id?:number; status?:number; enable_status?:number; order_by?:string; order?:'asc'|'desc' }
export interface ArticleSaveReq { id:number; category_id:number; title:string; slug:string; summary:string; author:string; source_name:string; source_url:string; cover:string; content:string; enable_status:number; virtual_view_count:number; sort:number }

export interface ArticleCommentItem {
  uid:string
  article_uid:string
  parent_uid:string
  root_uid:string
  nickname:string
  account:string
  avatar_url:string
  content:string
  images:string[]
  created_at:string
  ip_location:string
  like_count:number
  reply_count:number
  reply_nickname?:string
  status:number
  children?:ArticleCommentItem[]
  reply_page?:number
  reply_total?:number
  reply_loading?:boolean
}

export interface ArticleCommentListReq extends PageQuery { article_uid?:string; root_uid?:string; keyword?:string; status?:number }
