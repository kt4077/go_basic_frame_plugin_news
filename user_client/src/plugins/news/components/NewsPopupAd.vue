<script setup lang="ts">
declare const wx: { createInterstitialAd?: (options: { adUnitId: string }) => { show: () => Promise<unknown>; destroy?: () => void } }
const props = defineProps<{ unitId: string }>()
let advertisement: { show: () => Promise<unknown>; destroy?: () => void; onError?: (handler: () => void) => void } | undefined

onMounted(() => {
  const unitId = props.unitId.trim()
  if (!unitId) return
  // #ifdef MP-WEIXIN
  advertisement = wx.createInterstitialAd?.({ adUnitId: unitId })
  advertisement?.show().catch(() => undefined)
  // #endif
})

onUnmounted(() => advertisement?.destroy?.())
</script>

<template><view /></template>
