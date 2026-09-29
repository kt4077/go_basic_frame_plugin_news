<script setup lang="ts">
import { computed, ref } from 'vue'
import { ArrowDown, ChatDotRound, Hide, UserFilled, View } from '@element-plus/icons-vue'
import { ElImageViewer } from 'element-plus'
import type { ArticleCommentItem } from '../types/news'

const props = withDefaults(defineProps<{
  item: ArticleCommentItem
  depth?: number
}>(), {
  depth: 0,
})

const emit = defineEmits<{
  statusChange: [item: ArticleCommentItem, status: number]
  loadReplies: [item: ArticleCommentItem]
}>()

const imageViewerVisible = ref(false)

const remainingReplyCount = computed(() => Math.max(0, (props.item.reply_total || props.item.reply_count || 0) - (props.item.children?.length || 0)))

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

const statusType = (status: number) => {
  if (status === CommentStatus.Visible) return 'success'
  if (status === CommentStatus.Hidden) return 'info'
  return 'warning'
}

const forwardStatusChange = (item: ArticleCommentItem, status: number) => {
  emit('statusChange', item, status)
}

const forwardLoadReplies = (item: ArticleCommentItem) => {
  emit('loadReplies', item)
}

</script>

<template>
  <div class="comment-tree-node" :class="{ 'is-reply': depth > 0 }">
    <article class="comment-card">
      <el-avatar :size="depth > 0 ? 36 : 42" :src="item.avatar_url" class="comment-card__avatar">
        <el-icon :size="depth > 0 ? 17 : 20"><UserFilled /></el-icon>
      </el-avatar>
      <div class="comment-card__body">
        <div class="comment-card__header">
          <div class="comment-card__user">
            <strong>{{ item.nickname }}</strong>
            <span>{{ item.account }}</span>
            <el-tag :type="statusType(item.status)" size="small" effect="light">{{ CommentStatusLabels[item.status] }}</el-tag>
          </div>
          <div class="comment-card__right">
            <time>{{ item.created_at }}</time>
            <div class="comment-card__operations">
              <el-tooltip v-if="item.status !== CommentStatus.Visible" content="显示评论">
                <el-icon class="op-icon is-edit" @click="emit('statusChange', item, CommentStatus.Visible)"><View /></el-icon>
              </el-tooltip>
              <el-tooltip v-if="item.status !== CommentStatus.Hidden" content="关闭评论">
                <el-icon class="op-icon is-danger" @click="emit('statusChange', item, CommentStatus.Hidden)"><Hide /></el-icon>
              </el-tooltip>
            </div>
          </div>
        </div>
        <div v-if="item.reply_nickname" class="comment-card__reply">回复 {{ item.reply_nickname }}</div>
        <div class="comment-card__content-row">
          <span class="comment-card__content">{{ item.content }}</span>
          <button v-if="depth > 0 && item.images?.length" type="button" class="comment-card__image-link" @click="imageViewerVisible = true">
            查看图片（{{ item.images.length }}）
          </button>
        </div>
        <div v-if="depth === 0 && item.images?.length" class="comment-card__images">
          <el-image
            v-for="(image, index) in item.images"
            :key="image"
            :src="image"
            :preview-src-list="item.images"
            :initial-index="index"
            fit="cover"
            preview-teleported
            hide-on-click-modal
            class="comment-card__image"
          />
        </div>
        <el-image-viewer v-if="imageViewerVisible" :url-list="item.images" teleported @close="imageViewerVisible = false" />
        <div class="comment-card__meta">
          <span>IP属地：{{ item.ip_location }}</span>
          <span><el-icon><View /></el-icon> 点赞 {{ item.like_count }}</span>
          <span v-if="item.reply_count"><el-icon><ChatDotRound /></el-icon> 回复 {{ item.reply_count }}</span>
        </div>
      </div>
    </article>

    <div v-if="depth === 0 && (item.children?.length || item.reply_count)" class="comment-tree-children">
      <ArticleCommentTreeItem
        v-for="child in item.children || []"
        :key="child.uid"
        :item="child"
        :depth="depth + 1"
        @status-change="forwardStatusChange"
        @load-replies="forwardLoadReplies"
      />
      <div v-if="remainingReplyCount > 0" class="comment-tree-more">
        <el-button type="primary" link :icon="ArrowDown" :loading="item.reply_loading" @click="emit('loadReplies', item)">
          查看更多回复（剩余 {{ remainingReplyCount }} 条）
        </el-button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.comment-tree-node { position: relative; }
.comment-tree-node + .comment-tree-node { border-top: 1px solid var(--el-border-color-lighter); }
.comment-tree-children { position: relative; padding: 0 0 4px 30px; margin: 0 0 8px 20px; }
.comment-tree-children .comment-tree-node + .comment-tree-node { border-top: 1px dashed var(--el-border-color-lighter); }
.comment-tree-more { display: flex; align-items: center; min-height: 38px; padding: 2px 0 8px 4px; }
.comment-tree-more .el-button { font-size: 12px; }
.comment-card, .comment-card__user, .comment-card__right, .comment-card__operations, .comment-card__meta { display: flex; }
.comment-card { gap: 13px; padding: 18px 4px; }
.comment-tree-node.is-reply > .comment-card { padding-top: 13px; padding-bottom: 13px; }
.comment-card__avatar { flex-shrink: 0; background: var(--el-color-primary-light-9); color: var(--el-color-primary); }
.comment-card__body { min-width: 0; flex: 1; }
.comment-card__header { display: flex; align-items: flex-start; justify-content: space-between; gap: 12px; }
.comment-card__user { min-width: 0; align-items: center; gap: 8px; flex-wrap: wrap; }
.comment-card__user strong { color: var(--el-text-color-primary); font-size: 14px; }
.comment-card__user span, .comment-card__header time { color: var(--el-text-color-secondary); font-size: 12px; }
.comment-card__right { flex-shrink: 0; align-items: center; gap: 10px; }
.comment-card__operations { align-items: center; gap: 6px; }
.comment-card__reply { margin-top: 10px; color: var(--el-color-primary); font-size: 12px; }
.comment-card__content-row { min-width: 0; margin: 9px 0 12px; color: var(--el-text-color-primary); font-size: 14px; line-height: 1.7; overflow-wrap: anywhere; }
.comment-card__content { white-space: pre-wrap; }
.comment-card__image-link { display: inline; padding: 0; margin-left: 10px; border: 0; background: transparent; color: var(--el-color-primary); cursor: pointer; font: inherit; font-size: 12px; font-weight: 600; line-height: inherit; vertical-align: baseline; }
.comment-card__image-link:hover { color: var(--el-color-primary-dark-2); text-decoration: underline; }
.comment-card__images { display: flex; align-items: center; gap: 8px; margin: 2px 0 12px; flex-wrap: wrap; }
.comment-card__image { width: 76px; height: 76px; overflow: hidden; border: 1px solid var(--el-border-color-lighter); border-radius: 10px; background: var(--el-fill-color-light); cursor: zoom-in; }
.comment-card__meta { align-items: center; gap: 18px; color: var(--el-text-color-secondary); font-size: 12px; flex-wrap: wrap; }
.comment-card__meta span { display: inline-flex; align-items: center; gap: 4px; }
.comment-card__meta code { margin-left: auto; color: var(--el-text-color-placeholder); }

@media (max-width: 768px) {
  .comment-tree-children { padding-left: 18px; margin-left: 12px; }
  .comment-card__header, .comment-card__right { align-items: flex-start; flex-direction: column; }
  .comment-card__operations { align-self: flex-end; }
}
</style>
