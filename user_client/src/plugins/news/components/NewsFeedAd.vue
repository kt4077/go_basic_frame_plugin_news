<script setup lang="ts">
interface AdErrorEvent {
  detail?: {
    errCode?: number
    errMsg?: string
  }
}

const props = defineProps<{
  /** 后端广告配置接口返回的广告位 ID。 */
  unitId: string
  title?: string
  description?: string
}>()

const failedUnitId = ref('')
const resolvedUnitId = computed(() => props.unitId.trim())
const visible = computed(() => Boolean(resolvedUnitId.value) && failedUnitId.value !== resolvedUnitId.value)

const handleError = (_event: AdErrorEvent) => {
  failedUnitId.value = resolvedUnitId.value
}
</script>

<template>
  <!-- #ifdef MP-WEIXIN -->
  <view v-if="visible" class="article-card news-feed-ad">
    <view class="article-card__cover"><text>广告</text></view>
    <view class="article-card__content">
      <view class="article-card__heading"><text class="article-card__title">{{ title || '精选推荐' }}</text><text class="article-card__badge">广告</text></view>
      <text class="article-card__summary">{{ description || '更多精彩内容，点击了解详情' }}</text>
      <ad-custom class="article-card__ad" :unit-id="resolvedUnitId" @error="handleError" />
    </view>
  </view>
  <!-- #endif -->
</template>

<style lang="scss" scoped>
.article-card{display:flex;gap:22rpx;padding:28rpx 0;border-bottom:1px solid var(--app-border)}
.article-card__cover{display:flex;width:220rpx;height:196rpx;flex-shrink:0;align-items:center;justify-content:center;border-radius:18rpx;background:var(--app-primary-soft);color:var(--app-primary);font-size:24rpx;font-weight:700}
.article-card__content{display:flex;min-width:0;min-height:196rpx;flex:1;flex-direction:column}.article-card__heading{display:flex;align-items:flex-start;gap:10rpx}.article-card__title{display:-webkit-box;overflow:hidden;flex:1;color:var(--app-text);font-size:29rpx;font-weight:700;line-height:1.42;-webkit-box-orient:vertical;-webkit-line-clamp:2}.article-card__badge{flex-shrink:0;padding:2rpx 8rpx;border:1px solid var(--app-border);border-radius:6rpx;color:var(--app-icon-muted);font-size:16rpx}.article-card__summary{display:-webkit-box;overflow:hidden;margin-top:14rpx;color:var(--app-text-secondary);font-size:22rpx;line-height:1.5;-webkit-box-orient:vertical;-webkit-line-clamp:2}.article-card__ad{overflow:hidden;max-height:74rpx;margin-top:auto}
</style>
