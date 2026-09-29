<script setup lang="ts">
import type { NewsArticle, NewsCategory } from '../../types/news'
import AppEmptyState from '@/composables/components/AppEmptyState.vue'
import AppPageSkeleton from '@/composables/components/AppPageSkeleton.vue'
import AppSubPageHeader from '@/composables/components/AppSubPageHeader.vue'
import { appFeedback } from '@/composables/useAppFeedback'
import { usePageShare } from '@/composables/usePageShare'
import { getNewsAdvertisements, getNewsArticles, getNewsCategories } from '../../api/news'
import NewsArticleCard from '../../components/NewsArticleCard.vue'
import NewsFeedAd from '../../components/NewsFeedAd.vue'
import NewsPopupAd from '../../components/NewsPopupAd.vue'
import type { NewsAdvertisement } from '../../types/news'

interface NewsPageOptions {
  keyword?: string
}

definePage({
  name: 'plugin-news-index',
  style: {
    navigationBarTitleText: '资讯中心',
    navigationStyle: 'custom',
  },
})

usePageShare()

const pageSize = 10
const categories = ref<NewsCategory[]>([])
const articles = ref<NewsArticle[]>([])
const selectedCategory = ref('')
const keyword = ref('')
const page = ref(1)
const total = ref(0)
const loading = ref(true)
const loadingMore = ref(false)
const errorMessage = ref('')
const advertisements = ref<NewsAdvertisement[]>([])
const feedAdvertisement = computed(() => advertisements.value.find(item => item.format === 1))
const popupAdvertisement = computed(() => advertisements.value.find(item => item.format === 3))

const loadAdvertisements = async () => {
  try { return await getNewsAdvertisements(1) }
  catch { return [] as NewsAdvertisement[] }
}

const hasMore = computed(() => articles.value.length < total.value)
const categoryScrollTarget = computed(() => selectedCategory.value ? `category-${selectedCategory.value}` : 'category-all')

const loadArticles = async (reset = false) => {
  if (loadingMore.value) {
    return
  }
  if (reset) {
    page.value = 1
    errorMessage.value = ''
  }
  loadingMore.value = true
  try {
    const result = await getNewsArticles({
      category_slug: selectedCategory.value || undefined,
      keyword: keyword.value.trim() || undefined,
      page: page.value,
      page_size: pageSize,
    })
    articles.value = reset ? result.list : [...articles.value, ...result.list]
    total.value = result.total
  }
  catch (error) {
    errorMessage.value = error instanceof Error ? error.message : '资讯加载失败'
    if (!reset) {
      appFeedback.error(errorMessage.value)
    }
  }
  finally {
    loadingMore.value = false
  }
}

const loadPage = async () => {
  loading.value = true
  errorMessage.value = ''
  try {
    const [categoryList, advertisementList] = await Promise.all([
      getNewsCategories(),
      loadAdvertisements(),
      loadArticles(true),
    ])
    categories.value = categoryList
    advertisements.value = advertisementList
  }
  catch (error) {
    errorMessage.value = error instanceof Error ? error.message : '资讯加载失败'
  }
  finally {
    loading.value = false
  }
}

const selectCategory = async (slug: string) => {
  if (selectedCategory.value === slug) {
    return
  }
  selectedCategory.value = slug
  await loadArticles(true)
}

const search = async () => {
  await loadArticles(true)
}

const openArticle = (article: NewsArticle) => {
  uni.navigateTo({
    url: `/plugins/news/pages/detail/index?slug=${encodeURIComponent(article.slug)}`,
  })
}

onLoad((options?: NewsPageOptions) => {
  keyword.value = options?.keyword || ''
  loadPage()
})

onPullDownRefresh(async () => {
  await loadPage()
  uni.stopPullDownRefresh()
})

onReachBottom(async () => {
  if (!hasMore.value || loadingMore.value) {
    return
  }
  page.value += 1
  await loadArticles()
})
</script>

<template>
  <view class="news-page">
    <AppSubPageHeader title="资讯中心" />

    <AppPageSkeleton :loading="loading" :rows="5">
      <view class="news-hero">
        <view class="news-hero__copy">
          <text class="news-hero__eyebrow">
            DISCOVER NEWS
          </text>
          <text class="news-hero__title">
            发现值得阅读的新鲜内容
          </text>
          <text class="news-hero__desc">
            精选资讯，轻松掌握生活与服务动态
          </text>
        </view>
        <view class="news-hero__icon">
          <wd-icon name="book" color="#ffffff" size="34px" />
        </view>
      </view>

      <view class="news-search">
        <wd-icon name="search-line" color="var(--app-icon-muted)" size="18px" />
        <input
          v-model="keyword"
          class="news-search__input"
          placeholder="搜索资讯标题"
          placeholder-class="news-search__placeholder"
          confirm-type="search"
          @confirm="search"
        >
        <view class="news-search__button" @click="search">
          搜索
        </view>
      </view>

      <view class="category-nav">
        <scroll-view
          class="category-scroll"
          scroll-x
          scroll-with-animation
          :scroll-into-view="categoryScrollTarget"
          :show-scrollbar="false"
        >
          <view class="category-list">
            <view
              id="category-all"
              class="category-item"
              :class="{ 'is-active': selectedCategory === '' }"
              @click="selectCategory('')"
            >
              全部
            </view>
            <view
              v-for="category in categories"
              :id="`category-${category.slug}`"
              :key="category.id"
              class="category-item"
              :class="{ 'is-active': selectedCategory === category.slug }"
              @click="selectCategory(category.slug)"
            >
              {{ category.name }}
            </view>
          </view>
        </scroll-view>
        <view class="category-nav__fade" />
      </view>

      <AppEmptyState
        v-if="errorMessage && articles.length === 0"
        icon="close-circle"
        title="资讯暂时无法加载"
        :description="errorMessage"
      />
      <AppEmptyState
        v-else-if="articles.length === 0"
        title="暂无资讯内容"
        description="换个分类或关键词再看看"
      />
      <view v-else class="article-list">
        <template v-for="(article,index) in articles" :key="article.uid">
          <NewsArticleCard :article="article" @select="openArticle" />
          <NewsFeedAd v-if="index===0 && feedAdvertisement" :unit-id="feedAdvertisement.ad_id" :title="feedAdvertisement.name" :description="feedAdvertisement.description" />
        </template>
        <view class="article-list__footer">
          <wd-loading v-if="loadingMore" color="var(--app-primary)" size="18px" />
          <text v-else>
            {{ hasMore ? '上拉加载更多' : '已经到底啦' }}
          </text>
        </view>
      </view>
    </AppPageSkeleton>
    <NewsPopupAd v-if="popupAdvertisement" :unit-id="popupAdvertisement.ad_id" />
  </view>
</template>

<style lang="scss" scoped>
.news-page {
  min-height: 100vh;
  padding-bottom: calc(54rpx + env(safe-area-inset-bottom));
  background: var(--app-bg);
}

.news-hero {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin: 26rpx 30rpx 0;
  padding: 34rpx;
  border-radius: 34rpx;
  background: linear-gradient(135deg, var(--app-primary), #8178ff 62%, var(--app-accent));
  box-shadow: 0 18rpx 46rpx var(--app-shadow-strong);
  color: #ffffff;
}

.news-hero__copy {
  display: flex;
  min-width: 0;
  flex: 1;
  flex-direction: column;
}

.news-hero__eyebrow { font-size: 18rpx; letter-spacing: 4rpx; opacity: 0.76; }
.news-hero__title { margin-top: 13rpx; font-size: 34rpx; font-weight: 750; line-height: 1.35; }
.news-hero__desc { margin-top: 10rpx; font-size: 22rpx; opacity: 0.78; }

.news-hero__icon {
  display: flex;
  width: 92rpx;
  height: 92rpx;
  flex-shrink: 0;
  align-items: center;
  justify-content: center;
  margin-left: 22rpx;
  border: 1px solid rgba(255, 255, 255, 0.28);
  border-radius: 30rpx;
  background: rgba(255, 255, 255, 0.15);
}

.news-search {
  display: flex;
  height: 78rpx;
  align-items: center;
  margin: 24rpx 30rpx 0;
  padding: 0 10rpx 0 24rpx;
  border: 1px solid var(--app-border);
  border-radius: 26rpx;
  background: var(--app-surface);
}

.news-search__input { min-width: 0; height: 100%; flex: 1; padding: 0 18rpx; color: var(--app-text); font-size: 25rpx; }
.news-search__placeholder { color: var(--app-icon-muted); }
.news-search__button { padding: 14rpx 23rpx; border-radius: 20rpx; background: var(--app-primary); color: #ffffff; font-size: 23rpx; font-weight: 600; }

.category-nav { position: relative; display: flex; align-items: center; margin-top: 22rpx; }
.category-scroll { min-width: 0; flex: 1; white-space: nowrap; }
.category-list { display: inline-flex; gap: 12rpx; padding: 0 52rpx 5rpx 30rpx; }
.category-item { overflow: hidden; max-width: 190rpx; padding: 12rpx 23rpx; border: 1px solid var(--app-border); border-radius: 999rpx; background: var(--app-surface); color: var(--app-text-secondary); font-size: 22rpx; text-overflow: ellipsis; white-space: nowrap; }
.category-item.is-active { border-color: transparent; background: var(--app-primary); color: #ffffff; box-shadow: 0 8rpx 20rpx var(--app-shadow-strong); }
.category-nav__fade { position: absolute; right: 0; width: 42rpx; height: 58rpx; pointer-events: none; background: linear-gradient(90deg, transparent, var(--app-bg)); }

.article-list { display: flex; flex-direction: column; padding: 12rpx 30rpx 0; }
.article-list__footer { display: flex; min-height: 72rpx; align-items: center; justify-content: center; color: var(--app-text-secondary); font-size: 22rpx; }
</style>
