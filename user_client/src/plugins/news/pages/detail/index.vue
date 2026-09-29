<script setup lang="ts">
import type { NewsAdvertisement, NewsArticle } from '../../types/news'
import AppEmptyState from '@/composables/components/AppEmptyState.vue'
import AppPageSkeleton from '@/composables/components/AppPageSkeleton.vue'
import AppSubPageHeader from '@/composables/components/AppSubPageHeader.vue'
import { appFeedback } from '@/composables/useAppFeedback'
import { usePageShare } from '@/composables/usePageShare'
import { formatDate } from '@/utils/date'
import { createNewsRequestUID, getNewsAdvertisements, getNewsArticle, getNewsArticles, getNewsInteractionState, recordNewsShare, toggleNewsArticleCollection, toggleNewsArticleLike } from '../../api/news'
import { tokenStorage } from '@/utils/storage'
import NewsArticleFooter from '../../components/NewsArticleFooter.vue'
import NewsCommentSection from '../../components/NewsCommentSection.vue'
import NewsFeedAd from '../../components/NewsFeedAd.vue'
import NewsShareSheet from '../../components/NewsShareSheet.vue'
import NewsPopupAd from '../../components/NewsPopupAd.vue'

interface NewsDetailOptions {
  slug?: string
}

definePage({
  name: 'plugin-news-detail',
  style: {
    navigationBarTitleText: '资讯详情',
    navigationStyle: 'custom',
  },
})

const article = ref<NewsArticle | null>(null)
const loading = ref(true)
const errorMessage = ref('')
const slug = ref('')
const previousArticle = ref<NewsArticle | null>(null)
const nextArticle = ref<NewsArticle | null>(null)
const commentCount = ref(0)
const liked = ref(false)
const likeSubmitting = ref(false)
const collected = ref(false)
const collectionSubmitting = ref(false)
const shareVisible = ref(false)
const commentSectionRef = ref<InstanceType<typeof NewsCommentSection>>()
const advertisements = ref<NewsAdvertisement[]>([])
const feedAdvertisement = computed(() => advertisements.value.find(item => item.format === 1))
const popupAdvertisement = computed(() => advertisements.value.find(item => item.format === 3))

const loadAdvertisements = async () => {
  try { return await getNewsAdvertisements(2) }
  catch { return [] as NewsAdvertisement[] }
}

usePageShare(() => ({
  title: article.value?.title || '资讯详情',
  path: `/plugins/news/pages/detail/index?slug=${encodeURIComponent(slug.value)}`,
  imageUrl: article.value?.cover_url || undefined,
}))

const loadAdjacentArticles = async (currentArticle: NewsArticle) => {
  try {
    const related = await getNewsArticles({
      category_slug: currentArticle.category?.slug || undefined,
      page: 1,
      page_size: 100,
    })
    const currentIndex = related.list.findIndex(item => item.slug === currentArticle.slug)
    previousArticle.value = currentIndex >= 0 ? related.list[currentIndex + 1] || null : null
    nextArticle.value = currentIndex > 0 ? related.list[currentIndex - 1] || null : null
  }
  catch {
    previousArticle.value = null
    nextArticle.value = null
  }
}

const loadArticle = async () => {
  if (!slug.value) {
    errorMessage.value = '缺少文章标识'
    loading.value = false
    return
  }
  loading.value = true
  errorMessage.value = ''
  try {
    const [result, advertisementList] = await Promise.all([getNewsArticle(slug.value), loadAdvertisements()])
    advertisements.value = advertisementList
    article.value = result
    commentCount.value = result.comment_count
    if (tokenStorage.get()) {
      const state = await getNewsInteractionState(result.uid || result.slug)
      liked.value = state.liked
      collected.value = state.collected
    }
    void loadAdjacentArticles(result)
  }
  catch (error) {
    errorMessage.value = error instanceof Error ? error.message : '资讯加载失败'
  }
  finally {
    loading.value = false
  }
}

const openSource = () => {
  const sourceURL = article.value?.source_url?.trim()
  if (!sourceURL || !/^https?:\/\//i.test(sourceURL)) {
    return
  }
  uni.navigateTo({
    url: `/subpackages/webview/index?url=${encodeURIComponent(sourceURL)}`,
  })
}

const handlePendingAction = (name: string) => {
  appFeedback.info(`${name}功能即将开放`)
}

const openArticle = (target: NewsArticle) => {
  uni.redirectTo({
    url: `/plugins/news/pages/detail/index?slug=${encodeURIComponent(target.slug)}`,
  })
}

const toggleLike = async () => {
  const articleIdentifier = article.value?.uid || article.value?.slug
  if (!articleIdentifier || likeSubmitting.value) {
    return
  }

  likeSubmitting.value = true
  try {
    const result = await toggleNewsArticleLike(articleIdentifier)
    liked.value = result.active
    article.value!.like_count = result.count
  } catch (error) {
    appFeedback.error(error instanceof Error ? error.message : '点赞操作失败')
  } finally {
    likeSubmitting.value = false
  }
}

const handleShare = async (channel: number) => {
  if (!article.value || !tokenStorage.get()) return
  await recordNewsShare(article.value.uid, channel, createNewsRequestUID())
  article.value.share_count += 1
}

const toggleCollection = async () => {
  const articleIdentifier = article.value?.uid || article.value?.slug
  if (!articleIdentifier || collectionSubmitting.value) {
    return
  }

  collectionSubmitting.value = true
  try {
    const result = await toggleNewsArticleCollection(articleIdentifier)
    collected.value = result.active
    article.value!.collection_count = result.count
  } catch (error) {
    appFeedback.error(error instanceof Error ? error.message : '收藏操作失败')
  } finally {
    collectionSubmitting.value = false
  }
}

const openComments = async () => {
  uni.pageScrollTo({ selector: '#news-comments', duration: 260 })
  await nextTick()
  commentSectionRef.value?.focusComposer()
}

onLoad((options?: NewsDetailOptions) => {
  slug.value = options?.slug || ''
  loadArticle()
})
</script>

<template>
  <view class="detail-page">
    <AppSubPageHeader title="资讯详情" />
    <AppPageSkeleton :loading="loading" :rows="6">
      <AppEmptyState
        v-if="errorMessage || !article"
        icon="no-content"
        title="暂无资讯内容"
        :description="errorMessage"
        full-page
      />
      <article v-else class="detail-article">
        <view class="detail-header">
          <text class="detail-header__title">
            {{ article.title }}
          </text>
          <scroll-view class="detail-meta-scroll" scroll-x :show-scrollbar="false">
            <view class="detail-header__meta">
              <view class="detail-header__meta-item is-category">
                <wd-icon name="tag" color="var(--app-primary)" size="13px" />
                <text>{{ article.category?.name || '资讯' }}</text>
              </view>
              <view v-if="article.source_name" class="detail-header__meta-item is-source" @click="openSource">
                <wd-icon name="link" color="var(--app-primary)" size="13px" />
                <text>{{ article.source_name }}</text>
              </view>
              <view class="detail-header__meta-item">
                <wd-icon name="time-line" color="var(--app-icon-muted)" size="13px" />
                <text>{{ formatDate(article.published_at, 'YYYY-MM-DD') }}</text>
              </view>
              <view class="detail-header__meta-item">
                <wd-icon name="eye" color="var(--app-icon-muted)" size="13px" />
                <text>{{ article.total_view_count }}</text>
              </view>
            </view>
          </scroll-view>
          <text v-if="article.summary" class="detail-header__summary">
            {{ article.summary }}
          </text>
        </view>
        <view class="detail-content">
          <mp-html :content="article.content || ''" selectable />
        </view>
        <NewsArticleFooter
          :article="article"
          :previous-article="previousArticle"
          :next-article="nextArticle"
          @select="openArticle"
          @open-source="openSource"
        />
        <NewsFeedAd v-if="feedAdvertisement" :unit-id="feedAdvertisement.ad_id" :title="feedAdvertisement.name" :description="feedAdvertisement.description" />
        <NewsCommentSection
          ref="commentSectionRef"
          :article-uid="article.uid"
          :total-count="commentCount"
        />
      </article>

      <view v-if="article" class="detail-actions">
        <view class="detail-actions__source" @click="openSource">
          <view class="detail-actions__avatar">
            <wd-icon name="book" color="#ffffff" size="18px" />
          </view>
          <view class="detail-actions__source-copy">
            <text class="detail-actions__source-name">
              {{ article.source_name || article.category?.name || '资讯中心' }}
            </text>
            <text class="detail-actions__author">
              {{ article.author || '官方发布' }}
            </text>
          </view>
          <view class="detail-actions__follow" @click.stop="handlePendingAction('关注')">
            + 关注
          </view>
        </view>
        <view class="detail-actions__buttons">
          <view class="detail-actions__button" :class="{ 'is-disabled': likeSubmitting }" @tap.stop="toggleLike">
            <wd-icon :name="liked ? 'thumb-up-fill' : 'thumb-up'" :color="liked ? 'var(--app-primary)' : 'var(--app-icon)'" size="20px" />
            <text>{{ article.like_count }}</text>
          </view>
          <view class="detail-actions__button" @click="shareVisible = true">
            <wd-icon name="share-internal" color="var(--app-icon)" size="20px" />
            <text>分享</text>
          </view>
          <view class="detail-actions__button" :class="{ 'is-disabled': collectionSubmitting }" @tap.stop="toggleCollection">
            <wd-icon :name="collected ? 'heart-fill' : 'heart'" :color="collected ? 'var(--app-primary)' : 'var(--app-icon)'" size="20px" />
            <text>{{ article.collection_count || 0 }}</text>
          </view>
          <view class="detail-actions__button" @click="openComments">
            <wd-icon name="message" color="var(--app-icon)" size="20px" />
            <text>{{ commentCount }}</text>
          </view>
        </view>
      </view>
      <NewsShareSheet
        v-if="article"
        v-model="shareVisible"
        :title="article.title"
        :path="`/plugins/news/pages/detail/index?slug=${encodeURIComponent(slug)}`"
        @share="handleShare"
      />
      <NewsPopupAd v-if="popupAdvertisement" :unit-id="popupAdvertisement.ad_id" />
    </AppPageSkeleton>
  </view>
</template>

<style lang="scss" scoped>
.detail-page { min-height: 100vh; background: var(--app-surface); }
.detail-article { max-width: 860rpx; min-height: 100vh; box-sizing: border-box; padding: 38rpx 34rpx calc(190rpx + env(safe-area-inset-bottom)); margin: 0 auto; background: var(--app-surface); color: var(--app-text); }
.detail-header { display: flex; flex-direction: column; padding-bottom: 30rpx; }
.detail-header__title { color: var(--app-text); font-size: 44rpx; font-weight: 750; line-height: 1.42; letter-spacing: 0.5rpx; }
.detail-header__summary { margin-top: 24rpx; padding: 22rpx 24rpx; border-radius: 18rpx; background: var(--app-surface-soft); color: var(--app-text-secondary); font-size: 25rpx; line-height: 1.75; }
.detail-meta-scroll { width: 100%; margin-top: 24rpx; white-space: nowrap; }
.detail-header__meta, .detail-header__meta-item { display: inline-flex; align-items: center; }
.detail-header__meta { gap: 24rpx; padding-right: 24rpx; color: var(--app-icon-muted); font-size: 22rpx; }
.detail-header__meta-item { gap: 7rpx; }
.detail-header__meta-item.is-category, .detail-header__meta-item.is-source { color: var(--app-primary); }
.detail-content { padding-top: 34rpx; border-top: 1px solid var(--app-border); color: var(--app-text); font-size: 30rpx; line-height: 1.9; }
.detail-content :deep(p) { margin: 0 0 28rpx; }
.detail-content :deep(img) { max-width: 100%; height: auto; }

.detail-actions { position: fixed; z-index: 30; right: 0; bottom: 0; left: 0; display: flex; min-height: 104rpx; align-items: center; justify-content: space-between; gap: 22rpx; box-sizing: border-box; padding: 14rpx 26rpx calc(14rpx + env(safe-area-inset-bottom)); border-top: 1px solid var(--app-border); background: var(--app-floating-surface); box-shadow: 0 -12rpx 38rpx var(--app-shadow-soft); backdrop-filter: blur(24rpx); }
.detail-actions__source { display: flex; min-width: 0; flex: 1; align-items: center; }
.detail-actions__avatar { display: flex; width: 58rpx; height: 58rpx; flex-shrink: 0; align-items: center; justify-content: center; border-radius: 50%; background: linear-gradient(145deg, var(--app-primary), var(--app-accent)); }
.detail-actions__source-copy { display: flex; min-width: 0; flex-direction: column; margin-left: 13rpx; }
.detail-actions__source-name { overflow: hidden; color: var(--app-text); font-size: 23rpx; font-weight: 650; text-overflow: ellipsis; white-space: nowrap; }
.detail-actions__author { overflow: hidden; margin-top: 3rpx; color: var(--app-icon-muted); font-size: 19rpx; text-overflow: ellipsis; white-space: nowrap; }
.detail-actions__follow { flex-shrink: 0; padding: 9rpx 15rpx; margin-left: 12rpx; border-radius: 12rpx; background: var(--app-primary); color: #ffffff; font-size: 20rpx; font-weight: 650; }
.detail-actions__buttons { display: flex; flex-shrink: 0; align-items: center; gap: 20rpx; }
.detail-actions__button { display: flex; min-width: 46rpx; align-items: center; justify-content: center; flex-direction: column; gap: 3rpx; padding: 0; margin: 0; border: 0; background: transparent; color: var(--app-icon-muted); font-size: 18rpx; line-height: 1.2; }
.detail-actions__button.is-disabled { opacity: 0.55; }

@media (max-width: 390px) {
  .detail-actions { gap: 12rpx; padding-right: 18rpx; padding-left: 18rpx; }
  .detail-actions__source-copy { max-width: 120rpx; }
  .detail-actions__buttons { gap: 13rpx; }
}
</style>
