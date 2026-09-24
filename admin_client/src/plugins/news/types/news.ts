import type { PageQuery } from '@/types/common'

export interface CategoryItem { id:number; name:string; slug:string; sort:number; status:number; remark:string; created_at:string; updated_at:string }
export interface CategoryListReq extends PageQuery { keyword?:string; status?:number }
export interface CategorySaveReq { id:number; name:string; slug:string; sort:number; status:number; remark:string }

export interface ArticleCategoryItem { id:number; name:string; slug:string }
export interface ArticleItem { id:number; category_id:number; category:ArticleCategoryItem; title:string; slug:string; summary:string; cover:string; cover_url:string; content:string; status:number; enable_status:number; sort:number; view_count:number; virtual_view_count:number; total_view_count:number; published_at?:string; created_at:string; updated_at:string }
export interface ArticleListReq extends PageQuery { keyword?:string; category_id?:number; status?:number; enable_status?:number }
export interface ArticleSaveReq { id:number; category_id:number; title:string; slug:string; summary:string; cover:string; content:string; enable_status:number; virtual_view_count:number; sort:number }
