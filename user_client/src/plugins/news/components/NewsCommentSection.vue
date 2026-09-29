<script setup lang="ts">
import { appFeedback } from '@/composables/useAppFeedback'
import { createNewsComment, createNewsRequestUID, getNewsCommentReplies, getNewsComments, replyNewsComment, toggleNewsCommentLike, uploadNewsCommentImage } from '../api/news'
import type { NewsComment } from '../types/news'
import NewsCommentTreeItem from './NewsCommentTreeItem.vue'

const props = defineProps<{ articleUid: string; totalCount?: number }>()

const defaultAvatar = '/static/images/menu-placeholder.png'
const content = ref('')
const focused = ref(false)
const replyTo = ref('')
const replyUID = ref('')
const selectedImages = ref<string[]>([])
const showEmojiPanel = ref(false)
const likingCommentUIDs = ref<string[]>([])
const emojis = [
  '😀',
  '😃',
  '😄',
  '😁',
  '😆',
  '😂',
  '🤣',
  '😊',
  '🙂',
  '🙃',
  '😉',
  '😍',
  '🥰',
  '😘',
  '😋',
  '😎',
  '🤩',
  '🥳',
  '🤔',
  '🤭',
  '🤗',
  '😴',
  '🥺',
  '😭',
  '😤',
  '😡',
  '😱',
  '🤯',
  '👍',
  '👎',
  '👏',
  '🙌',
  '🤝',
  '💪',
  '🙏',
  '👌',
  '✌️',
  '🤟',
  '❤️',
  '💛',
  '💚',
  '💙',
  '💜',
  '🔥',
  '✨',
  '🌟',
  '🎉',
  '🎁',
  '💯',
  '🌹',
  '🍀',
  '☀️',
  '🌈',
  '🚀',
  '💡',
  '✅',
]
const comments = ref<NewsComment[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = 10

interface CommentDisplayRow {
  key: string
  type: 'comment' | 'more'
  item: NewsComment
  depth: number
  remaining: number
}

const commentRows = computed(() => {
  const rows: CommentDisplayRow[] = []
  const appendRows = (items: NewsComment[], depth: number) => {
    items.forEach((item) => {
      rows.push({ key: `comment-${item.uid}`, type: 'comment', item, depth, remaining: 0 })
      if (item.children?.length) {
        appendRows(item.children, depth + 1)
      }
      const remaining = Math.max(0, (item.reply_total || item.reply_count || 0) - (item.children?.length || 0))
      if (remaining > 0) {
        rows.push({ key: `more-${item.uid}`, type: 'more', item, depth: depth + 1, remaining })
      }
    })
  }
  appendRows(comments.value, 0)
  return rows
})

const loadComments = async (reset = false) => {
  if (!props.articleUid) return
  if (reset) page.value = 1
  const result = await getNewsComments(props.articleUid, page.value, pageSize)
  const rows = result.list.map(item => ({ ...item, children: [], reply_page: 0, reply_total: item.reply_count }))
  comments.value = reset ? rows : [...comments.value, ...rows]
  total.value = result.total
  await Promise.all(rows.filter(item => item.reply_count > 0).map(loadReplies))
}

const loadReplies = async (comment: NewsComment) => {
  const nextPage = (comment.reply_page || 0) + 1
  const result = await getNewsCommentReplies(comment.uid, nextPage, 5)
  const rows = result.list.map(item => ({ ...item, children: [], reply_page: 0, reply_total: item.reply_count }))
  comment.children = [...(comment.children || []), ...rows]
  comment.reply_page = nextPage
  comment.reply_total = result.total
}

const focusComposer = (target = '', uid = '') => {
  replyTo.value = target
  replyUID.value = uid
  focused.value = true
}

const appendEmoji = (emoji: string) => {
  content.value += emoji
  focused.value = true
}

const toggleEmojiPanel = () => {
  showEmojiPanel.value = !showEmojiPanel.value
}

const chooseImages = () => {
  uni.chooseImage({
    count: Math.max(1, 3 - selectedImages.value.length),
    sizeType: ['compressed'],
    success: (result) => {
      selectedImages.value = [...selectedImages.value, ...result.tempFilePaths].slice(0, 3)
    },
  })
}

const removeImage = (index: number) => {
  selectedImages.value.splice(index, 1)
}

const previewImages = (current: string, images: string[]) => {
  uni.previewImage({ current, urls: images })
}

const submit = async () => {
  const value = content.value.trim()
  if (!value && selectedImages.value.length === 0) {
    appFeedback.warning('请输入留言内容')
    return
  }
  try {
    const uploadedImages = await Promise.all(selectedImages.value.map(uploadNewsCommentImage))
    const payload = { content: value, images: uploadedImages, request_uid: createNewsRequestUID() }
    if (replyUID.value) {
      await replyNewsComment(replyUID.value, payload)
    } else {
      await createNewsComment(props.articleUid, payload)
    }
  } catch (error) {
    appFeedback.error(error instanceof Error ? error.message : '留言发布失败')
    return
  }
  content.value = ''
  replyTo.value = ''
  replyUID.value = ''
  selectedImages.value = []
  showEmojiPanel.value = false
  focused.value = false
  appFeedback.success('留言已提交，审核通过后展示')
}

const toggleLike = async (comment: NewsComment) => {
  if (!comment.uid || likingCommentUIDs.value.includes(comment.uid)) {
    return
  }

  likingCommentUIDs.value = [...likingCommentUIDs.value, comment.uid]
  try {
    const result = await toggleNewsCommentLike(comment.uid)
    comment.is_liked = result.active
    comment.like_count = result.count
  } catch (error) {
    appFeedback.error(error instanceof Error ? error.message : '点赞操作失败')
  } finally {
    likingCommentUIDs.value = likingCommentUIDs.value.filter(uid => uid !== comment.uid)
  }
}

defineExpose({ focusComposer })
watch(() => props.articleUid, () => loadComments(true), { immediate: true })
</script>

<template>
  <view id="news-comments" class="comment-section">
    <view class="comment-section__heading">
      <text>留言</text>
      <text class="comment-section__count">
        {{ props.totalCount ?? total }}
      </text>
    </view>

    <view class="comment-composer" :class="{ 'is-focused': focused }">
      <view v-if="replyTo" class="comment-composer__reply">
        回复 {{ replyTo }}
        <wd-icon name="close" color="var(--app-icon-muted)" size="14px" @click="replyTo = ''; replyUID = ''" />
      </view>
      <wd-textarea
        v-model="content"
        :focus="focused"
        :maxlength="300"
        auto-height
        clearable
        compact
        custom-style="--wot-textarea-padding: 0; --wot-textarea-inner-min-height: 128rpx;"
        custom-textarea-class="comment-composer__textarea"
        no-border
        disable-default-padding
        placeholder="写留言"
        confirm-type="send"
        @focus="focused = true"
        @confirm="submit"
      />
      <view v-if="selectedImages.length" class="comment-composer__images">
        <view v-for="(image, index) in selectedImages" :key="image" class="comment-composer__image-wrap">
          <image class="comment-composer__image" :src="image" mode="aspectFill" @click="previewImages(image, selectedImages)" />
          <view class="comment-composer__image-remove" @click="removeImage(index)">
            <wd-icon name="close" color="#ffffff" size="10px" />
          </view>
        </view>
      </view>
      <view v-if="showEmojiPanel" class="comment-composer__emoji-panel">
        <view
          v-for="emoji in emojis"
          :key="emoji"
          class="comment-composer__emoji"
          @click="appendEmoji(emoji)"
        >
          {{ emoji }}
        </view>
      </view>
      <view class="comment-composer__toolbar">
        <view class="comment-composer__tools">
          <view
            class="comment-composer__tool"
            :class="{ 'is-active': showEmojiPanel }"
            @click="toggleEmojiPanel"
          >
            <wd-icon
              name="face-smile-fill"
              color="var(--app-primary)"
              size="18px"
            />
          </view>
          <view class="comment-composer__tool is-image" @click="chooseImages">
            <wd-icon name="image" color="var(--app-primary)" size="18px" />
          </view>
        </view>
        <wd-button size="small" :disabled="!content.trim() && selectedImages.length === 0" @click="submit">
          发布
        </wd-button>
      </view>
    </view>

    <view class="comment-list">
      <template v-for="row in commentRows" :key="row.key">
        <view v-if="row.type === 'comment'" class="comment-list__row" :style="{ paddingLeft: `${Math.min(row.depth, 4) * 38}rpx` }">
          <NewsCommentTreeItem
            :item="row.item"
            :depth="row.depth"
            :default-avatar="defaultAvatar"
            @reply="item => focusComposer(item.nickname, item.uid)"
            @like="toggleLike"
            @preview="previewImages"
          />
        </view>
        <view v-else class="comment-list__more" :style="{ paddingLeft: `${Math.min(row.depth, 4) * 38 + 58}rpx` }" @click="loadReplies(row.item)">
          <text>查看全部 {{ row.remaining }} 条回复</text>
          <wd-icon name="arrow-down" color="var(--app-icon-muted)" size="12px" />
        </view>
      </template>
    </view>
  </view>
</template>

<style lang="scss" scoped>
.comment-section { padding-top: 52rpx; margin-top: 48rpx; border-top: 1px solid var(--app-border); }
.comment-section__heading { display: flex; align-items: center; gap: 10rpx; color: var(--app-text); font-size: 32rpx; font-weight: 700; }
.comment-section__count { color: var(--app-primary); font-size: 22rpx; font-weight: 600; }
.comment-composer { width: 100%; box-sizing: border-box; margin-top: 24rpx; padding: 10rpx 18rpx 12rpx; border: 1px solid var(--app-border); border-radius: 30rpx; background: var(--app-surface-soft); transition: border-color 0.2s ease; }
.comment-composer.is-focused { border-color: var(--app-primary); }
.comment-composer :deep(.wd-textarea) { min-height: 128rpx; padding: 0 4rpx; border-radius: 22rpx; background: transparent !important; }
.comment-composer :deep(.wd-textarea__body) { min-height: 128rpx; }
.comment-composer :deep(.wd-textarea__inner),
.comment-composer :deep(.comment-composer__textarea) { min-height: 128rpx !important; background: transparent !important; color: var(--app-text) !important; font-size: 25rpx; line-height: 48rpx; }
.comment-composer :deep(.wd-textarea__placeholder) { color: var(--app-icon-muted); }
.comment-composer__reply { display: flex; align-items: center; justify-content: space-between; padding: 12rpx 4rpx 4rpx; color: var(--app-primary); font-size: 21rpx; }
.comment-composer__images { display: flex; gap: 12rpx; flex-wrap: wrap; padding: 10rpx 0; }
.comment-composer__image-wrap { position: relative; }
.comment-composer__image { width: 112rpx; height: 112rpx; border-radius: 14rpx; }
.comment-composer__image-remove { position: absolute; top: -8rpx; right: -8rpx; display: flex; width: 30rpx; height: 30rpx; align-items: center; justify-content: center; border-radius: 50%; background: var(--app-feedback-error); }
.comment-composer__emoji-panel { display: flex; gap: 8rpx; flex-wrap: wrap; padding: 18rpx 0 12rpx; margin-top: 8rpx; border-top: 1px solid var(--app-border); }
.comment-composer__emoji { display: flex; width: 62rpx; height: 58rpx; align-items: center; justify-content: center; border-radius: 14rpx; background: var(--app-surface); color: var(--app-text); font-size: 30rpx; }
.comment-composer__toolbar, .comment-composer__tools { display: flex; align-items: center; }
.comment-composer__toolbar { min-height: 48rpx; justify-content: space-between; padding-top: 4rpx; }
.comment-composer__tools { gap: 12rpx; }
.comment-composer__tool { display: flex; width: 52rpx; height: 52rpx; align-items: center; justify-content: center; box-sizing: border-box; border: 1px solid transparent; border-radius: 16rpx; background: var(--app-primary-soft); transition: background-color 0.2s ease, border-color 0.2s ease, transform 0.2s ease; }
.comment-composer__tool.is-active { border-color: var(--app-primary); background: var(--app-primary-soft); }
.comment-composer__tool:active { transform: scale(0.94); }
.comment-list { display: flex; flex-direction: column; margin-top: 30rpx; }
.comment-list__row { box-sizing: border-box; }
.comment-list__more { display: flex; align-items: center; gap: 6rpx; padding-top: 14rpx; padding-bottom: 8rpx; color: var(--app-primary); font-size: 22rpx; }
.comment-item { display: flex; gap: 18rpx; padding: 22rpx 0; border-bottom: 1px solid var(--app-border); }
.comment-item.is-reply { margin-left: 68rpx; }
.comment-item__avatar { display: block; width: 58rpx; height: 58rpx; flex-shrink: 0; box-sizing: border-box; border: 2rpx solid var(--app-border); border-radius: 50%; background: var(--app-surface-soft); }
.comment-item__body { min-width: 0; flex: 1; }
.comment-item__topline, .comment-item__identity, .comment-item__meta, .comment-item__action { display: flex; align-items: center; }
.comment-item__topline { justify-content: space-between; gap: 16rpx; }
.comment-item__identity { min-width: 0; gap: 10rpx; }
.comment-item__author { overflow: hidden; color: var(--app-text-secondary); font-size: 23rpx; text-overflow: ellipsis; white-space: nowrap; }
.comment-item__author-tag { color: var(--app-primary); font-size: 20rpx; }
.comment-item__date { flex-shrink: 0; color: var(--app-icon-muted); font-size: 20rpx; }
.comment-item__meta { gap: 22rpx; margin-top: 12rpx; color: var(--app-icon-muted); font-size: 20rpx; }
.comment-item__action { gap: 6rpx; }
.comment-item__reply-to { color: var(--app-primary); font-size: 20rpx; }
.comment-item__content { display: block; margin-top: 8rpx; color: var(--app-text); font-size: 26rpx; line-height: 1.65; }
.comment-item__images { display: flex; gap: 10rpx; flex-wrap: wrap; margin-top: 14rpx; }
.comment-item__image { width: 150rpx; height: 150rpx; border-radius: 16rpx; background: var(--app-surface-soft); }
.comment-replies { display: flex; flex-direction: column; gap: 22rpx; margin-top: 22rpx; }
.comment-reply { display: flex; gap: 14rpx; padding-left: 10rpx; }
.comment-reply .comment-item__avatar { width: 44rpx; height: 44rpx; }
.comment-reply .comment-item__author { max-width: 190rpx; }
.comment-reply .comment-item__content { margin-top: 6rpx; font-size: 24rpx; line-height: 1.6; }
.comment-item__images.is-reply-images { gap: 14rpx; margin-top: 12rpx; }
.comment-item__images.is-reply-images .comment-item__image { width: 190rpx; height: 190rpx; border: 1px solid var(--app-border); border-radius: 10rpx; background: var(--app-surface-soft); }
.comment-replies__more { display: flex; align-items: center; gap: 6rpx; padding: 18rpx 0 0 58rpx; color: var(--app-primary); font-size: 22rpx; }
</style>
