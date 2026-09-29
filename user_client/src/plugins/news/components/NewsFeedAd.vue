<script setup lang="ts">
interface AdErrorEvent {
  detail?: {
    errCode?: number
    errMsg?: string
  }
}

const props = defineProps<{
  /** 接口返回的微信信息流广告位 ID；未传入时使用环境变量兜底。 */
  unitId?: string
}>()

const fallbackUnitId = import.meta.env.VITE_WECHAT_NEWS_FEED_AD_UNIT_ID?.trim() || ''
const failedUnitId = ref('')
const resolvedUnitId = computed(() => props.unitId?.trim() || fallbackUnitId)
const visible = computed(() => Boolean(resolvedUnitId.value) && failedUnitId.value !== resolvedUnitId.value)

const handleError = (_event: AdErrorEvent) => {
  failedUnitId.value = resolvedUnitId.value
}
</script>

<template>
  <!-- #ifdef MP-WEIXIN -->
  <view v-if="visible" class="news-feed-ad">
    <ad-custom
      :unit-id="resolvedUnitId"
      @error="handleError"
    />
  </view>
  <!-- #endif -->
</template>

<style lang="scss" scoped>
.news-feed-ad { overflow: hidden; margin-top: 42rpx; border-radius: 22rpx; background: var(--app-surface-soft); }
</style>
