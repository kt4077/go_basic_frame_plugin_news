<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { Delete, Edit, Plus } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { deleteAdvertisementConfig, getAdvertisementConfigs, getAdvertisementOptions, saveAdvertisementConfig } from '../../api/news'
import type { AdvertisementConfig, AdvertisementOption } from '../../types/news'

const formatLabels: Record<number,string> = { 1:'信息流', 2:'激励视频', 3:'弹窗广告' }
const positionLabels: Record<number,string> = { 1:'新闻列表', 2:'新闻详情' }
const loading = ref(false)
const visible = ref(false)
const rows = ref<AdvertisementConfig[]>([])
const options = ref<AdvertisementOption[]>([])
const form = reactive({ id:0, advertisement_id:0, position:1, status:1 })
const load = async () => { loading.value = true; try { [rows.value, options.value] = await Promise.all([getAdvertisementConfigs(), getAdvertisementOptions()]) } finally { loading.value = false } }
const openCreate = () => { Object.assign(form,{ id:0, advertisement_id:0, position:1, status:1 }); visible.value = true }
const openEdit = (row: AdvertisementConfig) => { Object.assign(form,{ id:row.id, advertisement_id:row.advertisement_id, position:row.position, status:row.config_status }); visible.value = true }
const submit = async () => { if (!form.advertisement_id) return ElMessage.warning('请选择广告'); await saveAdvertisementConfig({ ...form }); ElMessage.success('保存成功'); visible.value=false; await load() }
const remove = async (row: AdvertisementConfig) => { await ElMessageBox.confirm(`确认删除“${row.name}”的新闻广告配置吗？`,'删除确认',{type:'warning'}); await deleteAdvertisementConfig(row.id); ElMessage.success('删除成功'); await load() }
onMounted(load)
</script>

<template>
  <div class="page-container"><el-card shadow="never">
    <div class="toolbar"><el-alert title="候选广告仅显示核心广告配置中已绑定新闻插件的广告。" type="info" :closable="false" /><el-button v-perm="'POST:/admin/plugin/news/advertisement/save'" type="primary" :icon="Plus" @click="openCreate">新增配置</el-button></div>
    <el-table v-loading="loading" :data="rows" row-key="id">
      <el-table-column prop="name" label="广告名称" min-width="150" /><el-table-column prop="ad_id" label="广告ID" min-width="180" show-overflow-tooltip />
      <el-table-column label="广告形式" width="110"><template #default="{row}">{{ formatLabels[row.format] || '-' }}</template></el-table-column>
      <el-table-column label="展示位置" width="120"><template #default="{row}"><el-tag effect="plain">{{ positionLabels[row.position] }}</el-tag></template></el-table-column>
      <el-table-column prop="description" label="广告描述" min-width="220" show-overflow-tooltip />
      <el-table-column label="核心状态" width="100"><template #default="{row}"><el-tag :type="row.status===1?'success':'info'">{{ row.status===1?'启用':'停用' }}</el-tag></template></el-table-column>
      <el-table-column label="投放状态" width="100"><template #default="{row}"><el-tag :type="row.config_status===1?'success':'info'">{{ row.config_status===1?'启用':'停用' }}</el-tag></template></el-table-column>
      <el-table-column label="操作" width="110" fixed="right"><template #default="{row}"><el-icon v-perm="'POST:/admin/plugin/news/advertisement/save'" class="op-icon is-edit" @click="openEdit(row)"><Edit /></el-icon><el-icon v-perm="'POST:/admin/plugin/news/advertisement/delete'" class="op-icon is-danger" @click="remove(row)"><Delete /></el-icon></template></el-table-column>
    </el-table>
  </el-card>
  <el-dialog v-model="visible" :title="form.id?'修改广告配置':'新增广告配置'" width="560px"><el-form label-width="100px"><el-form-item label="广告" required><el-select v-model="form.advertisement_id" filterable style="width:100%"><el-option v-for="item in options" :key="item.id" :label="`${item.name}（${formatLabels[item.format]}）${item.status===2?' - 已停用':''}`" :value="item.id" /></el-select></el-form-item><el-form-item label="展示位置" required><el-radio-group v-model="form.position"><el-radio-button :value="1">新闻列表</el-radio-button><el-radio-button :value="2">新闻详情</el-radio-button></el-radio-group></el-form-item><el-form-item label="投放状态" required><el-radio-group v-model="form.status"><el-radio-button :value="1">启用</el-radio-button><el-radio-button :value="2">停用</el-radio-button></el-radio-group></el-form-item></el-form><template #footer><el-button @click="visible=false">取消</el-button><el-button type="primary" @click="submit">保存</el-button></template></el-dialog>
  </div>
</template>
<style scoped>.toolbar{display:flex;align-items:center;justify-content:space-between;gap:18px;margin-bottom:16px}.toolbar .el-alert{flex:1}.op-icon+.op-icon{margin-left:14px}</style>
