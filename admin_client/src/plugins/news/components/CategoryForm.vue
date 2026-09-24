<script setup lang="ts">
import { reactive, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { saveCategory } from '../api/news'
import { CategoryStatus } from '../enums/news'
import type { CategoryItem, CategorySaveReq } from '../types/news'
import { createRandomSlug } from '../utils/slug'

const emit = defineEmits<{ (event: 'success'): void }>()
const visible = ref(false)
const saving = ref(false)
const createForm = (): CategorySaveReq => ({ id: 0, name: '', slug: createRandomSlug('category'), sort: 0, status: CategoryStatus.Enabled, remark: '' })
const form = reactive<CategorySaveReq>(createForm())

const openCreate = () => {
  Object.assign(form, createForm())
  visible.value = true
}

const openEdit = (row: CategoryItem) => {
  Object.assign(form, { id: row.id, name: row.name, slug: row.slug, sort: row.sort, status: row.status, remark: row.remark })
  visible.value = true
}

const submit = async () => {
  if (!form.name.trim() || !form.slug.trim()) {
    ElMessage.warning('请填写分类名称和标识')
    return
  }
  saving.value = true
  try {
    await saveCategory({ ...form })
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
  <el-dialog v-model="visible" :title="form.id ? '修改分类' : '新增分类'" width="520px" destroy-on-close>
    <el-form label-width="88px">
      <el-form-item label="分类名称" required><el-input v-model="form.name" maxlength="64" /></el-form-item>
      <el-form-item label="分类标识" required>
        <el-input v-model="form.slug" :disabled="form.id > 0" maxlength="64" placeholder="仅支持字母和数字" />
        <div class="field-tip">{{ form.id ? '标识用于稳定访问，创建后不可修改' : '已自动生成，也可以在保存前修改' }}</div>
      </el-form-item>
      <el-form-item label="排序"><el-input-number v-model="form.sort" :min="0" :max="999999" /></el-form-item>
      <el-form-item label="状态">
        <el-radio-group v-model="form.status">
          <el-radio-button :value="CategoryStatus.Enabled">启用</el-radio-button>
          <el-radio-button :value="CategoryStatus.Disabled">停用</el-radio-button>
        </el-radio-group>
      </el-form-item>
      <el-form-item label="备注"><el-input v-model="form.remark" type="textarea" :rows="3" maxlength="255" show-word-limit /></el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="visible = false">取消</el-button>
      <el-button type="primary" :loading="saving" @click="submit">保存</el-button>
    </template>
  </el-dialog>
</template>

<style scoped>
.field-tip { margin-top: 6px; color: var(--el-text-color-secondary); font-size: 12px; line-height: 1.5; }
</style>
