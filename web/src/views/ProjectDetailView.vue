<script setup>
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { deleteProjectMedia, deleteProjectZip, downloadProjectZip, fetchAdminProjectMedia, fetchProjectDetail, fetchProjectMedia, moveProjectMedia, updateProject, uploadProjectMedia, uploadProjectZip } from '@/api/projects'

const route = useRoute()
const router = useRouter()

const loading = ref(false)
const saving = ref(false)
const errMsg = ref('')
const data = ref(null)
const mediaLoading = ref(false)
const mediaItems = ref([])
const zipBusy = ref(false)

const editMode = ref(false)
const canEdit = ref(false)
const form = ref({
  name: '',
  intro: '',
  gitRepo: '',
  guide: '',
  tags: '',
  isPublic: true,
})

const item = computed(() => data.value?.data?.project || null)
const project = computed(() => item.value?.project || null)
const id = computed(() => route.params.id)

const baseName = (p) => {
  const s = String(p || '')
  return s.split(/[/\\\\]/).pop()
}

const zipFiles = computed(() => {
  const files = item.value?.zipFiles || []
  return files.map((p) => baseName(p)).filter(Boolean)
})

const fillFormFromProject = () => {
  const p = project.value || {}
  form.value = {
    name: p.Name ?? p.name ?? '',
    intro: p.Intro ?? p.intro ?? '',
    gitRepo: p.GitRepo ?? p.gitRepo ?? '',
    guide: p.Guide ?? p.guide ?? '',
    tags: p.Tags ?? p.tags ?? '',
    isPublic: Boolean(p.IsPublic ?? p.isPublic ?? true),
  }
}

const loadMedia = async () => {
  if (!id.value) return
  mediaLoading.value = true
  try {
    const res = canEdit.value ? await fetchAdminProjectMedia(id.value) : await fetchProjectMedia(id.value)
    mediaItems.value = res?.data?.items || []
  } catch (e) {
    // ignore: media is optional
    mediaItems.value = []
  } finally {
    mediaLoading.value = false
  }
}

const load = async () => {
  loading.value = true
  errMsg.value = ''
  try {
    data.value = await fetchProjectDetail(id.value)
    fillFormFromProject()
    await loadMedia()
  } catch (e) {
    errMsg.value = e?.msg || e?.message || '请求失败'
  } finally {
    loading.value = false
  }
}

const enterEdit = () => {
  if (!canEdit.value) {
    ElMessage.warning('修改需要管理员登录（请先到 /admin 登录并设置 admin token）')
    return
  }
  editMode.value = true
}

const onSave = async () => {
  saving.value = true
  errMsg.value = ''
  try {
    await updateProject(id.value, {
      name: form.value.name,
      intro: form.value.intro,
      gitRepo: form.value.gitRepo,
      guide: form.value.guide,
      tags: form.value.tags,
      isPublic: form.value.isPublic,
    })
    editMode.value = false
    await load()
  } catch (e) {
    errMsg.value = e?.msg || e?.message || '保存失败'
  } finally {
    saving.value = false
  }
}

const onUpload = async (opt) => {
  const file = opt?.file
  if (!file) return
  try {
    await uploadProjectMedia(id.value, file)
    ElMessage.success('上传成功')
    await loadMedia()
    opt?.onSuccess?.()
  } catch (e) {
    opt?.onError?.(e)
  }
}

const onDeleteMedia = async (mid) => {
  try {
    await deleteProjectMedia(id.value, mid)
    ElMessage.success('删除成功')
    await loadMedia()
  } catch (e) {
    // handled by interceptor
  }
}

const onMove = async (mid, dir) => {
  try {
    await moveProjectMedia(id.value, mid, dir)
    await loadMedia()
  } catch (e) {}
}

const onUploadZip = async (opt) => {
  const file = opt?.file
  if (!file) return
  zipBusy.value = true
  try {
    await uploadProjectZip(id.value, file)
    ElMessage.success('上传成功（同名会覆盖更新）')
    await load()
    opt?.onSuccess?.()
  } catch (e) {
    opt?.onError?.(e)
  } finally {
    zipBusy.value = false
  }
}

const onDeleteZip = async (fileName) => {
  zipBusy.value = true
  try {
    await deleteProjectZip(id.value, fileName)
    ElMessage.success('删除成功')
    await load()
  } finally {
    zipBusy.value = false
  }
}

const onDownloadZip = async (fileName) => {
  zipBusy.value = true
  try {
    const blob = await downloadProjectZip(id.value, fileName)
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = fileName
    document.body.appendChild(a)
    a.click()
    a.remove()
    URL.revokeObjectURL(url)
  } finally {
    zipBusy.value = false
  }
}

onMounted(() => {
  canEdit.value = !!localStorage.getItem('admin_token')
  load()
})
</script>

<template>
  <main class="page">
    <div class="header">
      <div>
        <h1 class="title">项目详情</h1>
        <div v-if="project" class="sub">
          <span class="name">{{ project?.Name || project?.name }}</span>
          <span class="sep">·</span>
          <span class="muted">ID: {{ project?.ID || project?.id }}</span>
          <span class="sep">·</span>
          <span class="muted">folder: {{ project?.FolderName || project?.folderName }}</span>
        </div>
      </div>

      <div class="actions">
        <el-button @click="router.push('/projects')">返回列表</el-button>
        <el-button v-if="!editMode && canEdit" type="primary" :disabled="loading || !project" @click="enterEdit">修改</el-button>
        <el-button v-if="!editMode && !canEdit" :disabled="loading || !project" @click="enterEdit">修改（需管理员）</el-button>
        <template v-else>
          <el-button :disabled="saving" @click="editMode = false; fillFormFromProject()">取消</el-button>
          <el-button type="primary" :loading="saving" @click="onSave">保存</el-button>
        </template>
      </div>
    </div>

    <el-alert v-if="errMsg" type="error" :title="errMsg" show-icon :closable="false" style="margin-top: 12px" />

    <el-skeleton v-if="loading" :rows="8" animated style="margin-top: 12px" />

    <section v-else-if="project" class="card">
      <el-form v-if="editMode" label-width="90px" class="form">
        <el-form-item label="名称">
          <el-input v-model="form.name" placeholder="项目名" />
        </el-form-item>
        <el-form-item label="简介">
          <el-input v-model="form.intro" type="textarea" :rows="3" placeholder="简介" />
        </el-form-item>
        <el-form-item label="Git">
          <el-input v-model="form.gitRepo" placeholder="Git Repo URL" />
        </el-form-item>
        <el-form-item label="标签">
          <el-input v-model="form.tags" placeholder="逗号分隔" />
        </el-form-item>
        <el-form-item label="公开">
          <el-switch v-model="form.isPublic" />
        </el-form-item>
        <el-form-item label="教程">
          <el-input v-model="form.guide" type="textarea" :rows="8" placeholder="文字教程" />
        </el-form-item>
      </el-form>

      <div v-else class="grid">
        <div class="block">
          <div class="block-title">简介</div>
          <div class="block-body">{{ project?.Intro || project?.intro || '-' }}</div>
        </div>

        <div class="block">
          <div class="block-title">Git</div>
          <div class="block-body">
            <a v-if="project?.GitRepo || project?.gitRepo" :href="project?.GitRepo || project?.gitRepo" target="_blank" rel="noreferrer">
              {{ project?.GitRepo || project?.gitRepo }}
            </a>
            <span v-else>-</span>
          </div>
        </div>

        <div class="block">
          <div class="block-title">Tags</div>
          <div class="block-body">{{ project?.Tags || project?.tags || '-' }}</div>
        </div>

        <div class="block">
          <div class="block-title">公开</div>
          <div class="block-body">{{ String(project?.IsPublic ?? project?.isPublic) }}</div>
        </div>

        <div class="block" style="grid-column: 1 / -1">
          <div class="block-title">教程</div>
          <div class="block-body pre">{{ project?.Guide || project?.guide || '-' }}</div>
        </div>

        <div class="block" style="grid-column: 1 / -1">
          <div class="block-title">媒体（图片/视频）</div>
          <div class="block-body">
            <div class="media-actions">
              <el-upload v-if="canEdit" :show-file-list="false" :http-request="onUpload" accept="image/*,video/*">
                <el-button size="small" type="primary" :loading="mediaLoading">上传</el-button>
              </el-upload>
              <el-button size="small" :loading="mediaLoading" @click="loadMedia">刷新</el-button>
            </div>

            <div v-if="mediaLoading" class="muted" style="margin-top: 8px">加载中...</div>
            <div v-else-if="!mediaItems?.length" class="muted" style="margin-top: 8px">（空）</div>

            <div v-else class="media-strip">
              <div v-for="m in mediaItems" :key="m.id" class="media-card">
                <div class="thumb">
                  <video v-if="m.kind === 'video'" :src="m.url" controls muted />
                  <img v-else :src="m.url" alt="" />
                </div>
                <div class="media-meta">
                  <div class="media-name">{{ m.originalName || m.fileName }}</div>
                  <div v-if="canEdit" class="media-ops">
                    <el-button size="small" @click="onMove(m.id, 'left')">←</el-button>
                    <el-button size="small" @click="onMove(m.id, 'right')">→</el-button>
                    <el-button size="small" type="danger" @click="onDeleteMedia(m.id)">删除</el-button>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div class="block">
          <div class="block-title">media</div>
          <div class="block-body">
            <div class="path">{{ item?.mediaDir }}</div>
            <div v-if="item?.mediaFiles?.length" class="files">
              <div v-for="f in item.mediaFiles" :key="f" class="file">{{ f }}</div>
            </div>
            <div v-else class="muted">（空）</div>
          </div>
        </div>

        <div class="block">
          <div class="block-title">zip</div>
          <div class="block-body">
            <div class="path">{{ item?.zipDir }}</div>
            <div class="media-actions" style="margin-bottom: 8px">
              <el-upload v-if="canEdit" :show-file-list="false" :http-request="onUploadZip" accept=".zip,application/zip">
                <el-button size="small" type="primary" :loading="zipBusy">上传 zip</el-button>
              </el-upload>
              <el-button size="small" :loading="zipBusy" @click="load">刷新</el-button>
            </div>

            <div v-if="zipFiles?.length" class="files">
              <div v-for="f in zipFiles" :key="f" class="file zip-row">
                <span class="zip-name">{{ f }}</span>
                <span v-if="canEdit" class="zip-ops">
                  <el-button size="small" :loading="zipBusy" @click="onDownloadZip(f)">下载</el-button>
                  <el-button size="small" type="danger" :loading="zipBusy" @click="onDeleteZip(f)">删除</el-button>
                </span>
              </div>
            </div>
            <div v-else class="muted">（空）</div>
          </div>
        </div>
      </div>
    </section>

    <el-empty v-else description="未找到项目" style="margin-top: 16px" />
  </main>
</template>

<style scoped>
.page {
  max-width: 1100px;
  margin: 0 auto;
  padding: 24px;
}
.header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 16px;
}
.title {
  font-size: 28px;
  margin: 0;
}
.sub {
  margin-top: 6px;
  color: rgba(0, 0, 0, 0.6);
  font-size: 12px;
}
.name {
  font-weight: 700;
  color: rgba(0, 0, 0, 0.82);
}
.sep {
  margin: 0 6px;
}
.muted {
  color: rgba(0, 0, 0, 0.6);
}
.actions {
  display: flex;
  gap: 10px;
  align-items: center;
  flex-wrap: wrap;
}
.card {
  margin-top: 12px;
  padding: 12px;
  border: 1px solid var(--color-border);
  border-radius: 12px;
  background: #fff;
}
.grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
}
@media (max-width: 900px) {
  .grid {
    grid-template-columns: 1fr;
  }
}
.block {
  padding: 12px;
  border: 1px solid var(--color-border);
  border-radius: 12px;
  background: #fff;
}
.block-title {
  font-weight: 600;
  margin-bottom: 6px;
}
.block-body {
  color: rgba(0, 0, 0, 0.82);
  word-break: break-word;
}
.pre {
  white-space: pre-wrap;
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, 'Liberation Mono', 'Courier New', monospace;
  font-size: 12px;
  background: rgba(0, 0, 0, 0.03);
  padding: 10px;
  border-radius: 10px;
}
.path {
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, 'Liberation Mono', 'Courier New', monospace;
  font-size: 12px;
  color: rgba(0, 0, 0, 0.7);
  margin-bottom: 8px;
}
.files {
  display: grid;
  gap: 6px;
}
.file {
  padding: 6px 8px;
  border-radius: 8px;
  background: rgba(0, 0, 0, 0.04);
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, 'Liberation Mono', 'Courier New', monospace;
  font-size: 12px;
}
.form {
  max-width: 760px;
}

.media-actions {
  display: flex;
  gap: 10px;
  align-items: center;
  flex-wrap: wrap;
}

.media-strip {
  margin-top: 10px;
  display: flex;
  gap: 12px;
  overflow-x: auto;
  padding-bottom: 6px;
}

.media-card {
  width: 220px;
  flex: 0 0 auto;
  border: 1px solid var(--color-border);
  border-radius: 12px;
  background: #fff;
  overflow: hidden;
}

.thumb {
  height: 140px;
  background: rgba(0, 0, 0, 0.03);
  display: flex;
  align-items: center;
  justify-content: center;
}

.thumb img,
.thumb video {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
}

.media-meta {
  padding: 10px;
}

.media-name {
  font-size: 12px;
  color: rgba(0, 0, 0, 0.7);
  word-break: break-all;
  margin-bottom: 8px;
}

.media-ops {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.zip-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
}

.zip-name {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.zip-ops {
  display: flex;
  gap: 8px;
  flex: 0 0 auto;
}
</style>

