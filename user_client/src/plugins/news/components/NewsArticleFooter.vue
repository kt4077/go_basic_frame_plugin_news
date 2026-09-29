<script setup lang="ts">
import type { NewsArticle } from '../types/news'

defineProps<{
  article: NewsArticle
  previousArticle: NewsArticle | null
  nextArticle: NewsArticle | null
}>()

const emit = defineEmits<{
  select: [article: NewsArticle]
  openSource: []
}>()
</script>

<template>
  <view class="article-footer">
    <view class="article-footer__navigation">
      <view class="article-footer__links">
        <view
          class="article-footer__link"
          :class="{ 'is-disabled': !previousArticle }"
          @click="previousArticle && emit('select', previousArticle)"
        >
          <text class="article-footer__link-label">
            ‹ 上一篇
          </text>
          <text class="article-footer__link-title">
            {{ previousArticle?.title || '已经是第一篇' }}
          </text>
        </view>
        <view class="article-footer__divider" />
        <view
          class="article-footer__link article-footer__link--right"
          :class="{ 'is-disabled': !nextArticle }"
          @click="nextArticle && emit('select', nextArticle)"
        >
          <text class="article-footer__link-label">
            下一篇 ›
          </text>
          <text class="article-footer__link-title">
            {{ nextArticle?.title || '已经是最后一篇' }}
          </text>
        </view>
      </view>
    </view>

    <view class="article-footer__notice">
      <text>作者提示：个人观点，仅供参考</text>
      <view class="article-footer__notice-meta">
        <text
          class="article-footer__original"
          :class="{ 'is-disabled': !article.source_url }"
          @click="article.source_url && emit('openSource')"
        >
          阅读原文
        </text>
        <text>阅读 {{ article.total_view_count }}</text>
      </view>
    </view>
  </view>
</template>

<style lang="scss" scoped>
.article-footer { margin-top: 54rpx; }
.article-footer__navigation { padding: 26rpx; border: 1px solid var(--app-border); border-radius: 24rpx; background: var(--app-surface-soft); }
.article-footer__links { display: flex; align-items: stretch; }
.article-footer__link { display: flex; min-width: 0; flex: 1; flex-direction: column; padding-right: 22rpx; }
.article-footer__link--right { align-items: flex-end; padding-right: 0; padding-left: 22rpx; text-align: right; }
.article-footer__link-label { color: var(--app-text); font-size: 24rpx; font-weight: 650; }
.article-footer__link-title { display: -webkit-box; overflow: hidden; margin-top: 10rpx; color: var(--app-text-secondary); font-size: 22rpx; line-height: 1.45; -webkit-box-orient: vertical; -webkit-line-clamp: 2; }
.article-footer__link.is-disabled { opacity: 0.45; }
.article-footer__divider { width: 1px; flex-shrink: 0; background: var(--app-border); }
.article-footer__notice { display: flex; flex-direction: column; gap: 10rpx; margin-top: 30rpx; color: var(--app-icon-muted); font-size: 22rpx; line-height: 1.5; }
.article-footer__notice-meta { display: flex; gap: 18rpx; align-items: center; }
.article-footer__original { color: var(--app-primary); }
.article-footer__original.is-disabled { color: var(--app-icon-muted); }
</style>
