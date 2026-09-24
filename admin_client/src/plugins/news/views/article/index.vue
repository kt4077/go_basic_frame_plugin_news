<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Delete, Download, Edit, Plus, Search, Upload } from '@element-plus/icons-vue'
import AppPagination from '@/components/AppPagination.vue'
import { formatDateTimeCell } from '@/utils/datetime'
import ArticleForm from '../../components/ArticleForm.vue'
import { deleteArticle, getArticleList, getCategoryList, updateArticleStatus } from '../../api/news'
import { ArticleEnableStatus, ArticleEnableStatusLabels, ArticleStatus, ArticleStatusLabels, CategoryStatus } from '../../enums/news'
import type { ArticleItem, CategoryItem } from '../../types/news'

const loading = ref(false)
const route = useRoute()
const router = useRouter()
const list = ref<ArticleItem[]>([])
const categories = ref<CategoryItem[]>([])
const total = ref(0)
const query = reactive({ keyword: '', category_id: undefined as number | undefined, status: undefined as number | undefined, enable_status: undefined as number | undefined, page: 1, page_size: 20 })
const formVisible = computed(() => route.query.action === 'create' || route.query.action === 'edit')
const articleId = computed(() => route.query.action === 'edit' ? Number(route.query.id) || undefined : undefined)

const load = async () => {
  loading.value = true
  try {
    const result = await getArticleList(query)
    list.value = result.list
    total.value = result.total
  } finally {
    loading.value = false
  }
}

const loadCategories = async () => {
  const result = await getCategoryList({ page: 1, page_size: 100, status: CategoryStatus.Enabled })
  categories.value = result.list
}

const changeStatus = async (row: ArticleItem, status: ArticleStatus) => {
  const action = status === ArticleStatus.Published ? '发布' : '下线'
  await ElMessageBox.confirm(`确认${action}「${row.title}」吗？`, '提示', { type: 'warning' })
  await updateArticleStatus(row.id, status)
  ElMessage.success(`${action}成功`)
  await load()
}

const remove = async (row: ArticleItem) => {
  await ElMessageBox.confirm(`确认删除「${row.title}」吗？`, '提示', { type: 'warning' })
  await deleteArticle(row.id)
  ElMessage.success('删除成功')
  await load()
}

const statusType = (status: number) => status === ArticleStatus.Published ? 'success' : status === ArticleStatus.Offline ? 'info' : 'warning'
const openCreate = () => router.push({ path: route.path, query: { action: 'create' } })
const openEdit = (row: ArticleItem) => router.push({ path: route.path, query: { action: 'edit', id: String(row.id) } })
const closeForm = () => router.replace({ path: route.path })
const saved = async () => { await closeForm(); await load() }
onMounted(async () => { await Promise.all([load(), loadCategories()]) })
</script>

<template>
  <ArticleForm v-if="formVisible" :article-id="articleId" :categories="categories" @cancel="closeForm" @success="saved" />
  <div v-else class="page-card">
    <div class="toolbar">
      <div class="toolbar-left">
        <h4>新闻文章</h4>
        <el-input v-model="query.keyword" clearable placeholder="标题/标识" style="width: 200px" @keyup.enter="query.page = 1; load()" />
        <el-select v-model="query.category_id" clearable placeholder="分类" style="width: 140px" @change="query.page = 1; load()"><el-option v-for="item in categories" :key="item.id" :label="item.name" :value="item.id" /></el-select>
        <el-select v-model="query.status" clearable placeholder="状态" style="width: 120px" @change="query.page = 1; load()"><el-option v-for="(label, value) in ArticleStatusLabels" :key="value" :label="label" :value="Number(value)" /></el-select>
        <el-select v-model="query.enable_status" clearable placeholder="启用状态" style="width: 120px" @change="query.page = 1; load()"><el-option v-for="(label, value) in ArticleEnableStatusLabels" :key="value" :label="label" :value="Number(value)" /></el-select>
        <el-button :icon="Search" type="primary" @click="query.page = 1; load()">搜索</el-button>
      </div>
      <el-button v-perm="'POST:/admin/plugin/news/article/save'" :icon="Plus" type="primary" @click="openCreate">新增文章</el-button>
    </div>
    <el-table v-loading="loading" :data="list" stripe>
      <el-table-column label="封面" width="110" align="center">
        <template #default="{ row }"><el-image v-if="row.cover_url" :src="row.cover_url" :preview-src-list="[row.cover_url]" :initial-index="0" fit="contain" preview-teleported hide-on-click-modal class="cover" /><div v-else class="cover cover-empty">暂无封面</div></template>
      </el-table-column>
      <el-table-column prop="title" label="标题" min-width="220" show-overflow-tooltip />
      <el-table-column prop="slug" label="标识" min-width="180" show-overflow-tooltip />
      <el-table-column label="摘要" min-width="260" show-overflow-tooltip><template #default="{ row }"><span class="summary-text">{{ row.summary || '—' }}</span></template></el-table-column>
      <el-table-column prop="category.name" label="分类" width="130" />
      <el-table-column label="状态" width="100"><template #default="{ row }"><el-tag :type="statusType(row.status)">{{ ArticleStatusLabels[row.status as ArticleStatus] }}</el-tag></template></el-table-column>
      <el-table-column label="启用状态" width="100"><template #default="{ row }"><el-tag :type="row.enable_status === ArticleEnableStatus.Enabled ? 'success' : 'info'">{{ ArticleEnableStatusLabels[row.enable_status as ArticleEnableStatus] }}</el-tag></template></el-table-column>
      <el-table-column prop="view_count" label="真实阅读" width="100" align="center" />
      <el-table-column prop="virtual_view_count" label="虚拟阅读" width="100" align="center" />
      <el-table-column prop="total_view_count" label="总阅读" width="100" align="center" />
      <el-table-column prop="sort" label="排序" width="80" />
      <el-table-column prop="published_at" label="发布时间" width="170" :formatter="formatDateTimeCell" />
      <el-table-column label="操作" width="170" fixed="right">
        <template #default="{ row }"><div class="table-operations">
          <el-tooltip content="修改"><el-icon v-perm="'POST:/admin/plugin/news/article/save'" class="op-icon is-edit" @click="openEdit(row)"><Edit /></el-icon></el-tooltip>
          <el-tooltip v-if="row.status !== ArticleStatus.Published" content="发布"><el-icon v-perm="'POST:/admin/plugin/news/article/status'" class="op-icon" @click="changeStatus(row, ArticleStatus.Published)"><Upload /></el-icon></el-tooltip>
          <el-tooltip v-else content="下线"><el-icon v-perm="'POST:/admin/plugin/news/article/status'" class="op-icon is-warning" @click="changeStatus(row, ArticleStatus.Offline)"><Download /></el-icon></el-tooltip>
          <el-tooltip content="删除"><el-icon v-perm="'POST:/admin/plugin/news/article/delete'" class="op-icon is-danger" @click="remove(row)"><Delete /></el-icon></el-tooltip>
        </div></template>
      </el-table-column>
    </el-table>
    <AppPagination v-model:page="query.page" v-model:page-size="query.page_size" :total="total" @change="load" />
  </div>
</template>

<style scoped>
.toolbar-left { display: flex; align-items: center; gap: 12px; flex-wrap: wrap; }
.cover { width: 72px; height: 48px; margin: 0 auto; }
.cover :deep(.el-image__inner) { cursor: zoom-in; }
.cover-empty { display: flex; align-items: center; justify-content: center; border: 1px dashed var(--el-border-color); border-radius: 6px; color: var(--el-text-color-placeholder); font-size: 11px; }
.summary-text { color: var(--el-text-color-regular); line-height: 1.6; }
</style>
