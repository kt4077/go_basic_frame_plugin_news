<script setup lang="ts">
import { appFeedback } from '@/composables/useAppFeedback'

interface Props {
  modelValue: boolean
  title: string
  path: string
}

const props = defineProps<Props>()

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  share: [channel: number]
}>()

const visible = computed({
  get: () => props.modelValue,
  set: value => emit('update:modelValue', value),
})

const resolveShareLink = () => {
  // #ifdef H5
  return window.location.href
  // #endif

  // #ifndef H5
  return props.path
  // #endif
}

const copyLink = () => {
  uni.setClipboardData({
    data: resolveShareLink(),
    success: () => {
      emit('share', 4)
      visible.value = false
      appFeedback.success('链接已复制')
    },
    fail: () => appFeedback.error('复制失败，请稍后重试'),
  })
}

const shareToTimeline = () => {
  emit('share', 2)
  visible.value = false
  nextTick(() => {
    appFeedback.info('请点击右上角菜单，选择“分享到朋友圈”')
  })
}

const shareBySystem = async () => {
  // #ifdef H5
  if (typeof navigator !== 'undefined' && navigator.share) {
    try {
      await navigator.share({ title: props.title, url: resolveShareLink() })
      emit('share', 3)
      visible.value = false
      return
    }
    catch (error) {
      if (error instanceof Error && error.name === 'AbortError') {
        return
      }
    }
  }
  // #endif

  copyLink()
}
</script>

<template>
  <wd-action-sheet
    v-model="visible"
    title="分享资讯"
    cancel-text="取消"
    root-portal
    custom-class="news-share-sheet"
  >
    <view class="share-options">
      <!-- #ifdef MP-WEIXIN -->
      <button class="share-option" open-type="share" @click="emit('share', 1); visible = false">
        <view class="share-option__icon is-wechat">
          <wd-icon name="message" color="#ffffff" size="24px" />
        </view>
        <text>微信好友</text>
      </button>
      <view class="share-option" @click="shareToTimeline">
        <view class="share-option__icon is-timeline">
          <wd-icon name="user-group" color="#ffffff" size="24px" />
        </view>
        <text>朋友圈</text>
      </view>
      <!-- #endif -->

      <!-- #ifdef H5 -->
      <view class="share-option" @click="shareBySystem">
        <view class="share-option__icon is-system">
          <wd-icon name="share-internal" color="#ffffff" size="24px" />
        </view>
        <text>系统分享</text>
      </view>
      <!-- #endif -->

      <view class="share-option" @click="copyLink">
        <view class="share-option__icon is-copy">
          <wd-icon name="link" color="#ffffff" size="24px" />
        </view>
        <text>复制链接</text>
      </view>
    </view>
  </wd-action-sheet>
</template>

<style lang="scss" scoped>
:global(.news-share-sheet) { background: var(--app-surface) !important; color: var(--app-text); }
.share-options { display: flex; align-items: flex-start; gap: 46rpx; padding: 34rpx 38rpx 44rpx; }
.share-option { display: flex; width: 112rpx; flex-direction: column; align-items: center; gap: 14rpx; padding: 0; margin: 0; border: 0; background: transparent; color: var(--app-text-secondary); font-size: 22rpx; line-height: 1.3; }
.share-option::after { display: none; }
.share-option__icon { display: flex; width: 84rpx; height: 84rpx; align-items: center; justify-content: center; border-radius: 24rpx; box-shadow: 0 10rpx 24rpx var(--app-shadow-soft); }
.share-option__icon.is-wechat { background: linear-gradient(145deg, #37c86f, #18a957); }
.share-option__icon.is-timeline { background: linear-gradient(145deg, #31c896, #16a97d); }
.share-option__icon.is-system { background: linear-gradient(145deg, var(--app-primary), var(--app-accent)); }
.share-option__icon.is-copy { background: linear-gradient(145deg, #5d79eb, #7b66e8); }
</style>
