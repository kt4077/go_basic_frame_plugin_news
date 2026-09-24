<script setup lang="ts">
import { nextTick, onBeforeUnmount, onMounted, reactive, ref, shallowRef, watch } from 'vue'
import { ArrowLeft, DocumentChecked } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { Editor, Toolbar } from '@wangeditor/editor-for-vue'
import type { IDomEditor, IEditorConfig, IToolbarConfig } from '@wangeditor/editor'
import '@wangeditor/editor/dist/css/style.css'
import AvatarUpload from '@/components/AvatarUpload.vue'
import { uploadFile } from '@/api/upload'
import { getArticleDetail, saveArticle } from '../api/news'
import { ArticleEnableStatus } from '../enums/news'
import type { ArticleSaveReq, CategoryItem } from '../types/news'
import { createRandomSlug } from '../utils/slug'

const props = defineProps<{ articleId?: number; categories: CategoryItem[] }>()
const emit = defineEmits<{ (event: 'success'): void; (event: 'cancel'): void }>()
const formRef = ref<FormInstance>()
const editorRef = shallowRef<IDomEditor>()
const saving = ref(false)
const loading = ref(false)
const createForm = (): ArticleSaveReq => ({
  id: 0,
  category_id: 0,
  title: '',
  slug: createRandomSlug('article'),
  summary: '',
  cover: '',
  content: '',
  enable_status: ArticleEnableStatus.Enabled,
  virtual_view_count: 0,
  sort: 0,
})
const form = reactive<ArticleSaveReq>(createForm())
const detailStats = reactive({ viewCount: 0, totalViewCount: 0 })
const rules: FormRules<ArticleSaveReq> = {
  category_id: [{ required: true, message: '请选择所属分类', trigger: 'change' }],
  title: [{ required: true, message: '请输入文章标题', trigger: 'blur' }],
  slug: [
    { required: true, message: '请输入文章标识', trigger: 'blur' },
    { pattern: /^[a-zA-Z0-9][a-zA-Z0-9-_]*$/, message: '仅支持字母、数字、中划线和下划线', trigger: 'blur' },
  ],
  content: [{ required: true, message: '请输入文章正文', trigger: 'change' }],
}
const toolbarConfig: Partial<IToolbarConfig> = {}
const editorConfig: Partial<IEditorConfig> = {
  placeholder: '请输入文章正文，可插入图片、链接、引用和代码块',
  MENU_CONF: {
    uploadImage: {
      async customUpload(file: File, insertFn: (url: string, alt?: string, href?: string) => void) {
        const result = await uploadFile(file)
        insertFn(result.url, result.file_name, result.url)
      },
    },
  },
}

const loadDetail = async () => {
  Object.assign(form, createForm())
  detailStats.viewCount = 0
  detailStats.totalViewCount = 0
  formRef.value?.clearValidate()
  if (!props.articleId) return
  loading.value = true
  try {
    const detail = await getArticleDetail(props.articleId)
    Object.assign(form, {
      id: detail.id,
      category_id: detail.category_id,
      title: detail.title,
      slug: detail.slug,
      summary: detail.summary,
      cover: detail.cover_url || detail.cover,
      content: detail.content,
      enable_status: detail.enable_status,
      virtual_view_count: detail.virtual_view_count,
      sort: detail.sort,
    })
    detailStats.viewCount = detail.view_count
    detailStats.totalViewCount = detail.total_view_count
  } catch (error) {
    emit('cancel')
    throw error
  } finally {
    loading.value = false
  }
}

const handleCreated = (editor: IDomEditor) => { editorRef.value = editor }
const submit = async () => {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return
  saving.value = true
  try {
    await saveArticle({ ...form })
    ElMessage.success(form.id ? '文章修改成功' : '文章新增成功')
    emit('success')
  } finally {
    saving.value = false
  }
}

watch(() => props.articleId, async () => { await nextTick(); await loadDetail() })
onMounted(loadDetail)
onBeforeUnmount(() => editorRef.value?.destroy())
</script>

<template>
  <div v-loading="loading" class="article-editor-page">
    <div class="page-card editor-header">
      <div class="header-main">
        <el-button :icon="ArrowLeft" circle plain aria-label="返回文章列表" @click="emit('cancel')" />
        <div>
          <h4>{{ form.id ? '修改文章' : '新增文章' }}</h4>
          <p>完善文章内容和展示设置，保存后可在文章列表中发布。</p>
        </div>
      </div>
      <div class="header-actions">
        <el-button @click="emit('cancel')">取消</el-button>
        <el-button v-perm="'POST:/admin/plugin/news/article/save'" type="primary" :icon="DocumentChecked" :loading="saving" @click="submit">保存文章</el-button>
      </div>
    </div>

    <el-form ref="formRef" :model="form" :rules="rules" label-position="top">
      <div class="editor-layout">
        <section class="page-card content-card">
          <div class="section-title"><h4>文章内容</h4><span>标题、摘要与正文内容</span></div>
          <el-form-item label="文章标题" prop="title">
            <el-input v-model="form.title" maxlength="200" show-word-limit placeholder="请输入清晰准确的文章标题" />
          </el-form-item>
          <el-form-item label="文章摘要" prop="summary">
            <el-input v-model="form.summary" type="textarea" :rows="3" maxlength="500" show-word-limit placeholder="用于列表和分享场景的简短介绍" />
          </el-form-item>
          <el-form-item label="文章正文" prop="content">
            <div class="wang-editor">
              <Toolbar :editor="editorRef" :default-config="toolbarConfig" mode="default" class="editor-toolbar" />
              <Editor v-model="form.content" :default-config="editorConfig" mode="default" class="editor-content" @on-created="handleCreated" />
            </div>
          </el-form-item>
        </section>

        <aside class="editor-sidebar">
          <section class="page-card setting-card">
            <div class="section-title"><h4>基础设置</h4><span>分类和访问标识</span></div>
            <el-form-item label="所属分类" prop="category_id">
              <el-select v-model="form.category_id" placeholder="请选择分类" style="width: 100%">
                <el-option v-for="item in categories" :key="item.id" :label="item.name" :value="item.id" />
              </el-select>
            </el-form-item>
            <el-form-item label="文章标识" prop="slug">
              <el-input v-model="form.slug" :disabled="form.id > 0" maxlength="128" placeholder="仅支持字母、数字、中划线和下划线" />
              <div class="field-tip">{{ form.id ? '标识用于稳定访问，创建后不可修改' : '已自动生成，也可以在保存前修改' }}</div>
            </el-form-item>
            <div class="compact-grid">
              <el-form-item label="启用状态" prop="enable_status">
                <el-radio-group v-model="form.enable_status">
                  <el-radio-button :value="ArticleEnableStatus.Enabled">启用</el-radio-button>
                  <el-radio-button :value="ArticleEnableStatus.Disabled">禁用</el-radio-button>
                </el-radio-group>
              </el-form-item>
              <el-form-item label="排序" prop="sort">
                <el-input-number v-model="form.sort" :min="0" :max="999999" controls-position="right" />
              </el-form-item>
            </div>
          </section>

          <section class="page-card setting-card">
            <div class="section-title"><h4>展示设置</h4><span>封面与浏览数量</span></div>
            <el-form-item label="文章封面">
              <AvatarUpload v-model="form.cover" :size="132" shape="square" />
            </el-form-item>
            <el-form-item label="虚拟阅读数量" prop="virtual_view_count">
              <el-input-number v-model="form.virtual_view_count" :min="0" :max="999999999" controls-position="right" />
              <div class="field-tip">前台展示数量 = 真实阅读数量 + 虚拟阅读数量</div>
            </el-form-item>
            <div v-if="form.id" class="view-stats">
              <div><span>真实阅读</span><strong>{{ detailStats.viewCount }}</strong></div>
              <div><span>当前展示</span><strong>{{ detailStats.totalViewCount }}</strong></div>
            </div>
          </section>
        </aside>
      </div>
    </el-form>
  </div>
</template>

<style scoped>
.article-editor-page { display: flex; flex-direction: column; gap: 16px; }
.editor-header { display: flex; align-items: center; justify-content: space-between; gap: 20px; padding: 18px 22px; }
.header-main, .header-actions { display: flex; align-items: center; gap: 14px; }
.header-main h4, .section-title h4 { margin: 0; font-size: 16px; }
.header-main p, .section-title span { margin: 5px 0 0; color: var(--el-text-color-secondary); font-size: 12px; }
.editor-layout { display: grid; grid-template-columns: 340px minmax(0, 1fr); gap: 16px; align-items: start; }
.content-card, .setting-card { padding: 22px; }
.content-card { grid-column: 2; grid-row: 1; }
.editor-sidebar { display: flex; grid-column: 1; grid-row: 1; flex-direction: column; gap: 16px; }
.section-title { margin-bottom: 20px; padding-bottom: 14px; border-bottom: 1px solid var(--el-border-color-lighter); }
.compact-grid { display: grid; grid-template-columns: 1fr 110px; gap: 14px; }
.compact-grid :deep(.el-input-number), .setting-card :deep(.el-input-number) { width: 100%; }
.wang-editor { width: 100%; overflow: hidden; border: 1px solid var(--el-border-color); border-radius: 8px; background: var(--el-bg-color); }
.editor-toolbar { border-bottom: 1px solid var(--el-border-color); }
.editor-content { min-height: 460px; overflow-y: hidden; }
.field-tip { margin-top: 7px; color: var(--el-text-color-secondary); font-size: 12px; line-height: 1.5; }
.view-stats { display: grid; grid-template-columns: 1fr 1fr; gap: 10px; margin-top: 4px; }
.view-stats div { padding: 12px; border-radius: 8px; background: var(--el-fill-color-light); }
.view-stats span, .view-stats strong { display: block; }
.view-stats span { color: var(--el-text-color-secondary); font-size: 12px; }
.view-stats strong { margin-top: 5px; font-size: 18px; }
@media (max-width: 1100px) { .editor-layout { grid-template-columns: 1fr; } .content-card, .editor-sidebar { grid-column: 1; } .content-card { grid-row: 2; } .editor-sidebar { display: grid; grid-row: 1; grid-template-columns: 1fr 1fr; } }
@media (max-width: 720px) { .editor-header { align-items: flex-start; flex-direction: column; } .header-actions { width: 100%; justify-content: flex-end; } .editor-sidebar { display: flex; } .compact-grid { grid-template-columns: 1fr; } }
</style>
