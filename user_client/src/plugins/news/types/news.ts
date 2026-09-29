export interface NewsCategory {
  id: number
  name: string
  slug: string
}

export interface NewsArticle {
  uid: string
  category: NewsCategory
  title: string
  slug: string
  summary: string
  author: string
  source_name: string
  source_url: string
  cover_url: string
  content?: string
  view_count: number
  virtual_view_count: number
  total_view_count: number
  like_count: number
  share_count: number
  collection_count: number
  comment_count: number
  published_at: string | null
}

export interface NewsComment {
  uid: string
  nickname: string
  avatar_url: string
  reply_nickname?: string
  content: string
  images: string[]
  like_count: number
  reply_count: number
  is_liked: boolean
  is_mine: boolean
  created_at: string
  children?: NewsComment[]
  reply_page?: number
  reply_total?: number
}

export interface NewsCommentListResult {
  list: NewsComment[]
  total: number
}

export interface NewsArticleListResult {
  list: NewsArticle[]
  total: number
}

export interface NewsArticleListParams {
  category_slug?: string
  keyword?: string
  page: number
  page_size: number
}
