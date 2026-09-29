<script setup lang="ts">
import { computed, reactive, ref, watch } from 'vue'
import { ChatLineRound, CircleCheck, Clock, Refresh, Search } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { getArticleCommentList, getArticleCommentReplies, updateArticleCommentStatus } from '../api/news'
import type { ArticleCommentItem, ArticleItem } from '../types/news'
import ArticleCommentTreeItem from './ArticleCommentTreeItem.vue'

const props = defineProps<{
  modelValue: boolean
  article: ArticleItem | null
}>()

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
}>()

const CommentStatus = {
  Pending: 1,
  Visible: 2,
  Hidden: 3,
} as const

const CommentStatusLabels: Record<number, string> = {
  [CommentStatus.Pending]: '待审核',
  [CommentStatus.Visible]: '已展示',
  [CommentStatus.Hidden]: '已关闭',
}

const comments = ref<ArticleCommentItem[]>([])
const loading = ref(false)
const total = ref(0)
const query = reactive({ keyword: '', status: undefined as number | undefined, page: 1, page_size: 4 })
const visibleCount = computed(() => comments.value.filter(item => item.status === CommentStatus.Visible).length)
const pendingCount = computed(() => comments.value.filter(item => item.status === CommentStatus.Pending).length)

const loadComments = async () => {
  if (!props.article?.uid) return
  loading.value = true
  try {
    const result = await getArticleCommentList({ ...query, article_uid: props.article.uid })
    comments.value = result.list.map(item => ({ ...item, children: [], reply_page: 0, reply_total: item.reply_count }))
    total.value = result.total
    await Promise.all(comments.value.filter(item => item.reply_count > 0).map(loadReplies))
  } finally {
    loading.value = false
  }
}

const loadReplies = async (item: ArticleCommentItem) => {
  if (item.reply_loading) return
  item.reply_loading = true
  try {
    const nextPage = (item.reply_page || 0) + 1
    const result = await getArticleCommentReplies({ root_uid: item.uid, page: nextPage, page_size: 2 })
    item.children = [...(item.children || []), ...result.list]
    item.reply_page = nextPage
    item.reply_total = result.total
  } finally {
    item.reply_loading = false
  }
}

const reset = () => {
  query.keyword = ''
  query.status = undefined
  query.page = 1
  loadComments()
}

const updateCommentStatus = (items: ArticleCommentItem[], uid: string, status: number): boolean => {
  for (const item of items) {
    if (item.uid === uid) {
      item.status = status
      return true
    }
    if (item.children?.length && updateCommentStatus(item.children, uid, status)) {
      return true
    }
  }
  return false
}

const changeStatus = async (item: ArticleCommentItem, status: number) => {
  await updateArticleCommentStatus(item.uid, status)
  updateCommentStatus(comments.value, item.uid, status)
  ElMessage.success(status === CommentStatus.Visible ? '评论已显示' : '评论已关闭')
}

watch(
  () => props.modelValue,
  (visible) => {
    if (visible) {
      query.keyword = ''
      query.status = undefined
      query.page = 1
      loadComments()
    }
  },
)

watch(
  () => [query.page, query.page_size],
  loadComments,
)

let searchTimer: ReturnType<typeof setTimeout> | undefined
watch(
  () => [query.keyword, query.status],
  () => {
    query.page = 1
    if (searchTimer) clearTimeout(searchTimer)
    searchTimer = setTimeout(loadComments, 300)
  },
)
</script>

<template>
  <el-drawer
    :model-value="modelValue"
    direction="rtl"
    size="760px"
    destroy-on-close
    class="article-comment-drawer"
    @update:model-value="emit('update:modelValue', $event)"
  >
    <template #header>
      <div class="drawer-heading">
        <div class="drawer-heading__icon"><el-icon><ChatLineRound /></el-icon></div>
        <div class="drawer-heading__copy"><h4>文章评论</h4><p>{{ article?.title || '未选择文章' }}</p></div>
      </div>
    </template>

    <div class="comment-overview">
      <div class="overview-item"><el-icon><ChatDotRound /></el-icon><div><strong>{{ total }}</strong><span>全部评论</span></div></div>
      <div class="overview-item is-success"><el-icon><CircleCheck /></el-icon><div><strong>{{ visibleCount }}</strong><span>正常展示</span></div></div>
      <div class="overview-item is-warning"><el-icon><Clock /></el-icon><div><strong>{{ pendingCount }}</strong><span>待审核</span></div></div>
    </div>

    <div class="comment-toolbar">
      <el-input v-model="query.keyword" :prefix-icon="Search" clearable placeholder="昵称、账号或评论内容" />
      <el-select v-model="query.status" clearable placeholder="评论状态" style="width: 140px">
        <el-option v-for="(label, value) in CommentStatusLabels" :key="value" :label="label" :value="Number(value)" />
      </el-select>
      <el-button :icon="Refresh" @click="reset">重置</el-button>
    </div>

    <div class="comment-result"><span>评论数据</span><small>共 {{ total }} 条，回复按评论层级展开</small></div>
    <el-empty v-if="!loading && comments.length === 0" description="暂无符合条件的评论" />
    <div v-else v-loading="loading" class="comment-list">
      <ArticleCommentTreeItem v-for="item in comments" :key="item.uid" :item="item" @status-change="changeStatus" @load-replies="loadReplies" />
    </div>
    <div v-if="total" class="comment-pagination">
      <el-pagination
        v-model:current-page="query.page"
        v-model:page-size="query.page_size"
        :page-sizes="[4, 6, 8]"
        :total="total"
        layout="total, sizes, prev, pager, next"
        background
      />
    </div>

    <template #footer>
      <div class="drawer-footer"><span>评论显示状态修改后立即生效</span><el-button @click="emit('update:modelValue', false)">关闭</el-button></div>
    </template>
  </el-drawer>
</template>

<style scoped>
.drawer-heading, .drawer-heading__copy, .overview-item, .overview-item > div, .drawer-footer { display: flex; }
.drawer-heading { min-width: 0; align-items: center; gap: 12px; }
.drawer-heading__icon { display: flex; width: 42px; height: 42px; flex-shrink: 0; align-items: center; justify-content: center; overflow: visible; border-radius: 12px; background: var(--el-color-primary-light-9); color: var(--el-color-primary); }
.drawer-heading__icon .el-icon { width: 22px; height: 22px; overflow: visible; font-size: 22px; }
.drawer-heading__icon :deep(svg) { overflow: visible; }
.drawer-heading__copy { min-width: 0; flex-direction: column; gap: 4px; }
.drawer-heading h4, .drawer-heading p { overflow: hidden; margin: 0; text-overflow: ellipsis; white-space: nowrap; }
.drawer-heading h4 { color: var(--el-text-color-primary); font-size: 17px; }
.drawer-heading p { max-width: 580px; color: var(--el-text-color-secondary); font-size: 13px; font-weight: normal; }
.comment-overview { display: flex; gap: 12px; margin-bottom: 18px; }
.overview-item { min-width: 0; flex: 1; align-items: center; gap: 12px; padding: 15px; border: 1px solid var(--el-border-color-lighter); border-radius: 12px; background: var(--el-fill-color-light); color: var(--el-color-primary); }
.overview-item.is-success { color: var(--el-color-success); }
.overview-item.is-warning { color: var(--el-color-warning); }
.overview-item > .el-icon { font-size: 25px; }
.overview-item > div { min-width: 0; flex-direction: column; gap: 2px; }
.overview-item strong { color: var(--el-text-color-primary); font-size: 21px; }
.overview-item span { color: var(--el-text-color-secondary); font-size: 12px; }
.comment-toolbar { display: flex; align-items: center; gap: 10px; padding-bottom: 18px; border-bottom: 1px solid var(--el-border-color-lighter); }
.comment-toolbar .el-input { min-width: 220px; flex: 1; }
.comment-result { display: flex; align-items: center; justify-content: space-between; padding: 18px 2px 10px; }
.comment-result span { color: var(--el-text-color-primary); font-weight: 600; }
.comment-result small { color: var(--el-text-color-secondary); }
.comment-list { display: flex; flex-direction: column; }
.comment-pagination { display: flex; justify-content: flex-end; padding: 18px 0 4px; border-top: 1px solid var(--el-border-color-lighter); }
.drawer-footer { align-items: center; justify-content: space-between; }
.drawer-footer span { color: var(--el-text-color-secondary); font-size: 12px; }

@media (max-width: 768px) {
  .comment-overview, .comment-toolbar { flex-wrap: wrap; }
  .overview-item { flex-basis: calc(50% - 6px); }
  .comment-toolbar .el-input, .comment-toolbar .el-select { width: 100% !important; flex-basis: 100%; }
  .comment-pagination { overflow-x: auto; justify-content: flex-start; }
}
</style>
