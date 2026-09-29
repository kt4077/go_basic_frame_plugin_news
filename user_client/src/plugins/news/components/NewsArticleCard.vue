<script setup lang="ts">
import type { NewsArticle } from '../types/news'
import { formatDate } from '@/utils/date'

defineProps<{
  article: NewsArticle
}>()

const emit = defineEmits<{
  select: [article: NewsArticle]
}>()
</script>

<template>
  <view class="article-card" @click="emit('select', article)">
    <image
      v-if="article.cover_url"
      class="article-card__cover"
      :src="article.cover_url"
      mode="aspectFill"
      lazy-load
    />
    <view v-else class="article-card__cover article-card__cover--empty">
      <wd-icon name="image" color="var(--app-primary)" size="28px" />
    </view>

    <view class="article-card__content">
      <text class="article-card__title">
        {{ article.title }}
      </text>
      <text class="article-card__summary">
        {{ article.summary || '点击查看资讯详情' }}
      </text>
      <view class="article-card__meta">
        <view class="article-card__meta-item">
          <wd-icon name="tag" color="var(--app-primary)" size="12px" />
          <text>{{ article.category?.name || '资讯' }}</text>
        </view>
        <view class="article-card__meta-item">
          <wd-icon name="time-line" color="var(--app-icon-muted)" size="12px" />
          <text>{{ formatDate(article.published_at, 'MM-DD') }}</text>
        </view>
        <view class="article-card__meta-item">
          <wd-icon name="eye" color="var(--app-icon-muted)" size="12px" />
          <text>{{ article.total_view_count }}</text>
        </view>
      </view>
    </view>
  </view>
</template>

<style lang="scss" scoped>
.article-card {
  display: flex;
  gap: 22rpx;
  padding: 28rpx 0;
  border-bottom: 1px solid var(--app-border);
}

.article-card__cover {
  width: 220rpx;
  height: 196rpx;
  flex-shrink: 0;
  border-radius: 18rpx;
}

.article-card__cover--empty {
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--app-primary-soft);
}

.article-card__content {
  display: flex;
  min-width: 0;
  min-height: 196rpx;
  flex: 1;
  flex-direction: column;
}

.article-card__meta {
  display: flex;
  align-items: center;
}

.article-card__title {
  display: -webkit-box;
  overflow: hidden;
  color: var(--app-text);
  font-size: 29rpx;
  font-weight: 700;
  line-height: 1.42;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 2;
}

.article-card__summary {
  display: -webkit-box;
  overflow: hidden;
  margin-top: 14rpx;
  color: var(--app-text-secondary);
  font-size: 22rpx;
  line-height: 1.5;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 2;
}

.article-card__meta {
  gap: 18rpx;
  margin-top: 18rpx;
  color: var(--app-icon-muted);
  font-size: 20rpx;
}

.article-card__meta-item {
  display: flex;
  min-width: 0;
  align-items: center;
  gap: 6rpx;
  white-space: nowrap;
}

.article-card__meta-item:first-child {
  overflow: hidden;
  color: var(--app-primary);
}

.article-card__meta-item:first-child text {
  overflow: hidden;
  max-width: 116rpx;
  text-overflow: ellipsis;
}
</style>
