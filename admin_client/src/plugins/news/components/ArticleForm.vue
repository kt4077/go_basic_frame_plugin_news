<script setup lang="ts">
import { reactive, ref } from 'vue'
import { ElMessage } from 'element-plus'
import AvatarUpload from '@/components/AvatarUpload.vue'
import { getArticleDetail, saveArticle } from '../api/news'
import type { ArticleItem, ArticleSaveReq, CategoryItem } from '../types/news'

defineProps<{ categories: CategoryItem[] }>()
const emit = defineEmits<{ (event: 'success'): void }>()
const visible = ref(false)
const saving = ref(false)
const loading = ref(false)
const createForm = (): ArticleSaveReq => ({ id: 0, category_id: 0, title: '', slug: '', summary: '', cover: '', content: '', sort: 0 })
const form = reactive<ArticleSaveReq>(createForm())

const openCreate = () => {
  Object.assign(form, createForm())
  visible.value = true
}

const openEdit = async (row: ArticleItem) => {
  visible.value = true
  loading.value = true
  try {
    const detail = await getArticleDetail(row.id)
    Object.assign(form, {
      id: detail.id,
      category_id: detail.category_id,
      title: detail.title,
      slug: detail.slug,
      summary: detail.summary,
      cover: detail.cover_url || detail.cover,
      content: detail.content,
      sort: detail.sort,
    })
  } catch (error) {
    visible.value = false
    throw error
  } finally {
    loading.value = false
  }
}

const submit = async () => {
  if (!form.category_id || !form.title.trim() || !form.slug.trim() || !form.content.trim()) {
    ElMessage.warning('请完整填写分类、标题、标识和正文')
    return
  }
  saving.value = true
  try {
    await saveArticle({ ...form })
    ElMessage.success('保存成功')
    visible.value = false
    emit('success')
  } finally {
    saving.value = false
  }
}

defineExpose({ openCreate, openEdit })
</script>

<template>
  <el-dialog v-model="visible" :title="form.id ? '修改文章' : '新增文章'" width="760px" destroy-on-close>
    <el-form v-loading="loading" label-width="88px">
      <div class="form-grid">
        <el-form-item label="所属分类" required>
          <el-select v-model="form.category_id" style="width: 100%">
            <el-option v-for="item in categories" :key="item.id" :label="item.name" :value="item.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="排序"><el-input-number v-model="form.sort" :min="0" :max="999999" /></el-form-item>
      </div>
      <el-form-item label="文章标题" required><el-input v-model="form.title" maxlength="200" show-word-limit /></el-form-item>
      <el-form-item label="文章标识" required><el-input v-model="form.slug" maxlength="128" placeholder="仅字母和数字，发布后不建议修改" /></el-form-item>
      <el-form-item label="文章摘要"><el-input v-model="form.summary" type="textarea" :rows="2" maxlength="500" show-word-limit /></el-form-item>
      <el-form-item label="文章封面"><AvatarUpload v-model="form.cover" :size="120" shape="square" /></el-form-item>
      <el-form-item label="文章正文" required><el-input v-model="form.content" type="textarea" :rows="12" placeholder="支持输入HTML或纯文本内容" /></el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="visible = false">取消</el-button>
      <el-button type="primary" :loading="saving" :disabled="loading" @click="submit">保存草稿</el-button>
    </template>
  </el-dialog>
</template>

<style scoped>
.form-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 16px; }
</style>
