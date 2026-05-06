<script setup>
import { reactive, ref } from 'vue'
import service from '@/utils/request'

const props = defineProps({
  modelValue: { type: Boolean, default: false },
})
const emit = defineEmits(['update:modelValue', 'created'])

const saving = ref(false)

const form = reactive({
  name: '',
  folderName: '',
  isPublic: true,
  intro: '',
  codeStartAt: null,
})

const close = () => emit('update:modelValue', false)

const submit = async () => {
  if (!form.name.trim()) return
  saving.value = true
  try {
    const res = await service.post('/admin/projects', {
      name: form.name,
      folderName: form.folderName,
      isPublic: form.isPublic,
      intro: form.intro,
      codeStartAt: form.codeStartAt,
    })
    emit('created', res?.data)
    close()
  } finally {
    saving.value = false
  }
}
</script>

<template>
  <el-dialog :model-value="modelValue" title="新建项目" width="560px" @close="close">
    <el-form label-width="110px">
      <el-form-item label="项目名" required>
        <el-input v-model="form.name" placeholder="必填" />
      </el-form-item>
      <el-form-item label="文件夹名">
        <el-input v-model="form.folderName" placeholder="可不填（自动生成）" />
      </el-form-item>
      <el-form-item label="是否公开">
        <el-switch v-model="form.isPublic" active-text="公开" inactive-text="不公开" />
      </el-form-item>
      <el-form-item label="简介">
        <el-input v-model="form.intro" type="textarea" :rows="3" placeholder="可不填" />
      </el-form-item>
      <el-form-item label="开始时间">
        <el-date-picker
          v-model="form.codeStartAt"
          type="date"
          placeholder="可不填（默认今天）"
          style="width: 100%"
          format="YYYY-MM-DD"
        />
      </el-form-item>
    </el-form>

    <template #footer>
      <el-button @click="close">取消</el-button>
      <el-button type="primary" :loading="saving" :disabled="!form.name.trim()" @click="submit">
        创建
      </el-button>
    </template>
  </el-dialog>
</template>
