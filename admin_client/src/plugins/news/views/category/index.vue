<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Delete, Edit, Plus, Refresh, Search } from '@element-plus/icons-vue'
import AppPagination from '@/components/AppPagination.vue'
import { formatDateTimeCell } from '@/utils/datetime'
import CategoryForm from '../../components/CategoryForm.vue'
import { deleteCategory, getCategoryList } from '../../api/news'
import { CategoryStatus, CategoryStatusLabels } from '../../enums/news'
import type { CategoryItem } from '../../types/news'

const loading = ref(false)
const list = ref<CategoryItem[]>([])
const total = ref(0)
const formRef = ref<InstanceType<typeof CategoryForm>>()
const query = reactive({ keyword: '', status: undefined as number | undefined, page: 1, page_size: 20 })

const load = async () => {
  loading.value = true
  try {
    const result = await getCategoryList(query)
    list.value = result.list
    total.value = result.total
  } finally {
    loading.value = false
  }
}

const reset = () => {
  query.keyword = ''
  query.status = undefined
  query.page = 1
  load()
}

const remove = async (row: CategoryItem) => {
  await ElMessageBox.confirm(`确认删除分类「${row.name}」吗？`, '提示', { type: 'warning' })
  await deleteCategory(row.id)
  ElMessage.success('删除成功')
  await load()
}

onMounted(load)
</script>

<template>
  <div class="page-card">
    <div class="toolbar">
      <div class="toolbar-left">
        <h4>新闻分类</h4>
        <el-input v-model="query.keyword" clearable placeholder="名称/标识" style="width: 200px" @keyup.enter="query.page = 1; load()" />
        <el-select v-model="query.status" clearable placeholder="状态" style="width: 110px" @change="query.page = 1; load()">
          <el-option v-for="(label, value) in CategoryStatusLabels" :key="value" :label="label" :value="Number(value)" />
        </el-select>
        <el-button :icon="Search" type="primary" @click="query.page = 1; load()">搜索</el-button>
        <el-button :icon="Refresh" @click="reset">重置</el-button>
      </div>
      <el-button v-perm="'POST:/admin/plugin/news/category/save'" :icon="Plus" type="primary" @click="formRef?.openCreate()">新增分类</el-button>
    </div>
    <el-table v-loading="loading" :data="list" stripe>
      <el-table-column prop="name" label="分类名称" min-width="160" />
      <el-table-column prop="slug" label="分类标识" min-width="150" />
      <el-table-column prop="sort" label="排序" width="90" />
      <el-table-column label="状态" width="100">
        <template #default="{ row }"><el-tag :type="row.status === CategoryStatus.Enabled ? 'success' : 'info'">{{ CategoryStatusLabels[row.status as CategoryStatus] }}</el-tag></template>
      </el-table-column>
      <el-table-column prop="remark" label="备注" min-width="180" show-overflow-tooltip />
      <el-table-column prop="updated_at" label="更新时间" width="170" :formatter="formatDateTimeCell" />
      <el-table-column label="操作" width="110" fixed="right">
        <template #default="{ row }">
          <div class="table-operations">
            <el-tooltip content="修改"><el-icon v-perm="'POST:/admin/plugin/news/category/save'" class="op-icon is-edit" @click="formRef?.openEdit(row)"><Edit /></el-icon></el-tooltip>
            <el-tooltip content="删除"><el-icon v-perm="'POST:/admin/plugin/news/category/delete'" class="op-icon is-danger" @click="remove(row)"><Delete /></el-icon></el-tooltip>
          </div>
        </template>
      </el-table-column>
    </el-table>
    <AppPagination v-model:page="query.page" v-model:page-size="query.page_size" :total="total" @change="load" />
    <CategoryForm ref="formRef" @success="load" />
  </div>
</template>

<style scoped>
.toolbar-left { display: flex; align-items: center; gap: 12px; flex-wrap: wrap; }
</style>
