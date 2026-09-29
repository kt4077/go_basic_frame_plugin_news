<script setup lang="ts">
import type { NewsComment } from '../types/news'
import { formatDate } from '@/utils/date'

const props = withDefaults(defineProps<{
  item: NewsComment
  depth?: number
  defaultAvatar: string
}>(), {
  depth: 0,
})

const emit = defineEmits<{
  reply: [item: NewsComment]
  like: [item: NewsComment]
  preview: [current: string, images: string[]]
}>()
</script>

<template>
  <view class="comment-node" :class="{ 'is-reply': depth > 0 }">
    <image class="comment-node__avatar" :class="{ 'is-small': depth > 0 }" :src="item.avatar_url || defaultAvatar" mode="aspectFill" />
    <view class="comment-node__body">
      <view class="comment-node__identity">
        <text class="comment-node__author">{{ item.nickname }}</text>
        <text v-if="depth > 0 && item.reply_nickname" class="comment-node__reply-to">回复 {{ item.reply_nickname }}</text>
      </view>
      <text class="comment-node__content">{{ item.content }}</text>
      <view v-if="item.images?.length" class="comment-node__images">
        <image
          v-for="image in item.images"
          :key="image"
          class="comment-node__image"
          :class="{ 'is-reply-image': depth > 0 }"
          :src="image"
          :mode="depth > 0 ? 'aspectFit' : 'aspectFill'"
          @click="emit('preview', image, item.images || [])"
        />
      </view>
      <view class="comment-node__meta">
        <text>{{ formatDate(item.created_at, 'YYYY-MM-DD') }}</text>
        <view class="comment-node__action" @tap.stop="emit('like', item)">
          <wd-icon :name="item.is_liked ? 'thumb-up-fill' : 'thumb-up'" :color="item.is_liked ? 'var(--app-primary)' : 'var(--app-icon-muted)'" :size="depth > 0 ? '13px' : '14px'" />
          <text>点赞{{ item.like_count ? ` ${item.like_count}` : '' }}</text>
        </view>
        <view class="comment-node__action" @click="emit('reply', item)">
          <wd-icon name="message" color="var(--app-icon-muted)" :size="depth > 0 ? '13px' : '14px'" />
          <text>回复{{ item.reply_total ? ` ${item.reply_total}` : '' }}</text>
        </view>
      </view>
    </view>
  </view>
</template>

<style lang="scss" scoped>
.comment-node { display: flex; gap: 18rpx; padding: 22rpx 0; border-bottom: 1px solid var(--app-border); }
.comment-node.is-reply { gap: 14rpx; padding: 18rpx 0 0; border-bottom: 0; }
.comment-node__avatar { display: block; width: 58rpx; height: 58rpx; flex-shrink: 0; box-sizing: border-box; border: 2rpx solid var(--app-border); border-radius: 50%; background: var(--app-surface-soft); }
.comment-node__avatar.is-small { width: 44rpx; height: 44rpx; }
.comment-node__body { min-width: 0; flex: 1; }
.comment-node__identity, .comment-node__meta, .comment-node__action { display: flex; align-items: center; }
.comment-node__identity { min-width: 0; gap: 10rpx; }
.comment-node__author { overflow: hidden; max-width: 220rpx; color: var(--app-text-secondary); font-size: 23rpx; text-overflow: ellipsis; white-space: nowrap; }
.comment-node__reply-to { color: var(--app-primary); font-size: 20rpx; }
.comment-node__content { display: block; margin-top: 8rpx; color: var(--app-text); font-size: 26rpx; line-height: 1.65; }
.is-reply .comment-node__content { margin-top: 6rpx; font-size: 24rpx; line-height: 1.6; }
.comment-node__images { display: flex; gap: 10rpx; flex-wrap: wrap; margin-top: 14rpx; }
.comment-node__image { width: 150rpx; height: 150rpx; border-radius: 16rpx; background: var(--app-surface-soft); }
.comment-node__image.is-reply-image { width: 190rpx; height: 190rpx; border: 1px solid var(--app-border); border-radius: 10rpx; }
.comment-node__meta { gap: 22rpx; margin-top: 12rpx; color: var(--app-icon-muted); font-size: 20rpx; }
.comment-node__action { gap: 6rpx; }
</style>
