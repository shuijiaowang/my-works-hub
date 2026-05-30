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
const activeMediaId = ref('')

const editMode = ref(false)
const canEdit = ref(false)
const form = ref({
  name: '',
  intro: '',
  links: [],
  guide: '',
  tags: '',
  isPublic: true,
})

const item = computed(() => data.value?.data?.project || null)
const project = computed(() => item.value?.project || null)
const folderName = computed(() => route.params.folderName)

const projectName = computed(() => project.value?.Name ?? project.value?.name ?? '')
const projectIntro = computed(() => project.value?.Intro ?? project.value?.intro ?? '')
const normalizeLinks = (raw) => {
  const list = raw?.Links ?? raw?.links ?? []
  if (!Array.isArray(list)) return []
  return list.map((l) => ({
    name: l?.Name ?? l?.name ?? '',
    url: l?.URL ?? l?.url ?? '',
  }))
}

const projectLinks = computed(() => normalizeLinks(project.value))
const projectGuide = computed(() => project.value?.Guide ?? project.value?.guide ?? '')
const projectPublicText = computed(() => String(project.value?.IsPublic ?? project.value?.isPublic ?? ''))
const tagsList = computed(() => {
  const raw = project.value?.Tags ?? project.value?.tags ?? ''
  return String(raw)
    .split(',')
    .map((s) => s.trim())
    .filter(Boolean)
})

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
    links: normalizeLinks(p),
    guide: p.Guide ?? p.guide ?? '',
    tags: p.Tags ?? p.tags ?? '',
    isPublic: Boolean(p.IsPublic ?? p.isPublic ?? true),
  }
}

const loadMedia = async () => {
  if (!folderName.value) return
  mediaLoading.value = true
  try {
    const res = canEdit.value ? await fetchAdminProjectMedia(folderName.value) : await fetchProjectMedia(folderName.value)
    mediaItems.value = res?.data?.items || []

    const items = mediaItems.value || []
    if (!items.length) {
      activeMediaId.value = ''
    } else {
      const exists = activeMediaId.value && items.some((m) => m.id === activeMediaId.value)
      if (!exists) activeMediaId.value = items[0].id
    }
  } catch (e) {
    // ignore: media is optional
    mediaItems.value = []
    activeMediaId.value = ''
  } finally {
    mediaLoading.value = false
  }
}

const load = async () => {
  loading.value = true
  errMsg.value = ''
  try {
    data.value = await fetchProjectDetail(folderName.value)
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
    await updateProject(folderName.value, {
      name: form.value.name,
      intro: form.value.intro,
      links: form.value.links,
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
    await uploadProjectMedia(folderName.value, file)
    ElMessage.success('上传成功')
    await loadMedia()
    opt?.onSuccess?.()
  } catch (e) {
    opt?.onError?.(e)
  }
}

const onDeleteMedia = async (mid) => {
  try {
    await deleteProjectMedia(folderName.value, mid)
    ElMessage.success('删除成功')
    await loadMedia()
  } catch (e) {
    // handled by interceptor
  }
}

const onMove = async (mid, dir) => {
  try {
    await moveProjectMedia(folderName.value, mid, dir)
    await loadMedia()
  } catch (e) {}
}

const onUploadZip = async (opt) => {
  const file = opt?.file
  if (!file) return
  zipBusy.value = true
  try {
    await uploadProjectZip(folderName.value, file)
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
    await deleteProjectZip(folderName.value, fileName)
    ElMessage.success('删除成功')
    await load()
  } finally {
    zipBusy.value = false
  }
}

const onDownloadZip = async (fileName) => {
  zipBusy.value = true
  try {
    const blob = await downloadProjectZip(folderName.value, fileName)
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

const addLink = () => {
  form.value.links.push({ name: '', url: '' })
}

const removeLink = (index) => {
  form.value.links.splice(index, 1)
}

onMounted(() => {
  canEdit.value = !!localStorage.getItem('admin_token')
  load()
})

const activeMedia = computed(() => {
  const items = mediaItems.value || []
  if (!items.length) return null
  const found = items.find((m) => m.id === activeMediaId.value)
  return found || items[0] || null
})

const setActiveMedia = (m) => {
  if (!m?.id) return
  activeMediaId.value = m.id
}
</script>

<template>
  <main class="page">
    <div class="header">
      <div>
        <div class="crumb">
          <span class="crumb-link" role="button" tabindex="0" @click="router.push('/projects')" @keydown.enter.prevent="router.push('/projects')">
            我的作品
          </span>
          <span class="crumb-sep">/</span>
          <span class="crumb-current">详情</span>
        </div>
        <h1 class="title">{{ projectName || '项目详情' }}</h1>
        <div v-if="project" class="sub">
          <span class="muted">ID: {{ project?.ID || project?.id }}</span>
          <span class="sep">·</span>
          <span class="muted">folder: {{ project?.FolderName || project?.folderName }}</span>
        </div>
      </div>

      <div class="actions">
        <el-button @click="router.push('/projects')">返回列表</el-button>
        <el-button v-if="canEdit" type="primary" :disabled="loading || !project" @click="enterEdit">管理/修改</el-button>
        <el-button v-else :disabled="loading || !project" @click="enterEdit">管理/修改（需管理员）</el-button>
      </div>
    </div>

    <el-alert v-if="errMsg" type="error" :title="errMsg" show-icon :closable="false" style="margin-top: 12px" />

    <el-skeleton v-if="loading" :rows="8" animated style="margin-top: 12px" />

    <section v-else-if="project" class="layout">
      <div class="content">
        <div class="hero">
          <div class="hero-media">
            <div v-if="mediaLoading" class="hero-empty muted">媒体加载中...</div>
            <div v-else-if="!mediaItems?.length" class="hero-empty muted">暂无媒体</div>
            <div v-else class="hero-viewer">
              <div class="hero-stage">
                <video v-if="activeMedia?.kind === 'video'" :src="activeMedia?.url" controls playsinline />
                <img v-else :src="activeMedia?.url" alt="" loading="lazy" />
              </div>

              <div v-if="mediaItems.length > 1" class="hero-thumbs" role="tablist" aria-label="媒体缩略图">
                <button
                  v-for="m in mediaItems"
                  :key="m.id"
                  class="hero-thumb"
                  :class="{ active: m.id === activeMediaId }"
                  type="button"
                  role="tab"
                  :aria-selected="m.id === activeMediaId"
                  @click="setActiveMedia(m)"
                >
                  <div class="hero-thumb-inner">
                    <video v-if="m.kind === 'video'" :src="m.url" muted playsinline />
                    <img v-else :src="m.url" alt="" loading="lazy" />
                  </div>
                </button>
              </div>
            </div>
          </div>

          <div class="hero-meta">
            <div class="intro" v-if="projectIntro">{{ projectIntro }}</div>
            <div v-else class="muted">暂无简介</div>

            <div v-if="tagsList.length" class="tags">
              <el-tag v-for="t in tagsList" :key="t" class="tag" size="small" effect="plain">{{ t }}</el-tag>
            </div>

            <div v-if="projectLinks.length" class="links">
              <div v-for="(link, i) in projectLinks" :key="i" class="link-item">
                <a class="git" :href="link.url" target="_blank" rel="noreferrer">{{ link.name || link.url }}</a>
                <span v-if="link.name" class="muted link-raw">{{ link.url }}</span>
              </div>
            </div>
          </div>
        </div>

        <div class="section">
          <div class="section-title">资源下载</div>
          <div class="panel-body">
            <div class="res-title">Zip</div>
            <div v-if="zipFiles?.length" class="files">
              <div v-for="f in zipFiles" :key="f" class="file zip-row">
                <span class="zip-name" :title="f">{{ f }}</span>
                <span class="zip-ops">
                  <el-button size="small" :loading="zipBusy" @click="onDownloadZip(f)">下载</el-button>
                </span>
              </div>
            </div>
            <div v-else class="muted">（空）</div>
          </div>
        </div>

        <div class="section">
          <div class="section-title">教程</div>
          <div class="guide pre">{{ projectGuide || '-' }}</div>
        </div>
      </div>

      <aside class="side">
        <div class="panel">
          <div class="panel-title">信息</div>
          <div class="kv">
            <div class="k">公开</div>
            <div class="v">{{ projectPublicText }}</div>
          </div>
        </div>
      </aside>
    </section>

    <el-empty v-else description="未找到项目" style="margin-top: 16px" />

    <el-drawer
      v-model="editMode"
      :size="520"
      direction="rtl"
      :with-header="false"
      :close-on-click-modal="false"
      :close-on-press-escape="false"
      class="admin-drawer"
    >
      <div class="drawer-head">
        <div class="drawer-title">管理 / 修改</div>
        <div class="drawer-actions">
          <el-button :disabled="saving" @click="editMode = false; fillFormFromProject()">取消</el-button>
          <el-button type="primary" :loading="saving" @click="onSave">保存</el-button>
        </div>
      </div>

      <div class="drawer-body">
        <el-form label-width="90px" class="form">
          <el-form-item label="名称">
            <el-input v-model="form.name" placeholder="项目名" />
          </el-form-item>
          <el-form-item label="简介">
            <el-input v-model="form.intro" type="textarea" :rows="3" placeholder="简介" />
          </el-form-item>
          <el-form-item label="链接">
            <div class="link-editor">
              <div v-for="(link, i) in form.links" :key="i" class="link-row">
                <el-input v-model="link.name" placeholder="名称，如 GitHub" />
                <el-input v-model="link.url" placeholder="链接 URL" />
                <el-button @click="removeLink(i)">删除</el-button>
              </div>
              <el-button type="primary" plain @click="addLink">添加链接</el-button>
            </div>
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

        <div class="panel admin-panel">
          <div class="panel-title">媒体（图片/视频）</div>
          <div class="panel-body">
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

        <div class="panel admin-panel">
          <div class="panel-title">zip</div>
          <div class="panel-body">
            <div class="media-actions" style="margin-bottom: 8px">
              <el-upload v-if="canEdit" :show-file-list="false" :http-request="onUploadZip" accept=".zip,application/zip">
                <el-button size="small" type="primary" :loading="zipBusy">上传 zip</el-button>
              </el-upload>
              <el-button size="small" :loading="zipBusy" @click="load">刷新</el-button>
            </div>

            <div v-if="zipFiles?.length" class="files">
              <div v-for="f in zipFiles" :key="f" class="file zip-row">
                <span class="zip-name" :title="f">{{ f }}</span>
                <span class="zip-ops">
                  <el-button size="small" :loading="zipBusy" @click="onDownloadZip(f)">下载</el-button>
                  <el-button size="small" type="danger" :loading="zipBusy" @click="onDeleteZip(f)">删除</el-button>
                </span>
              </div>
            </div>
            <div v-else class="muted">（空）</div>
          </div>
        </div>
      </div>
    </el-drawer>
  </main>
</template>

<style scoped>
.page {
  max-width: 1440px;
  margin: 0 auto;
  padding: clamp(16px, 2.4vw, 28px);
}
.header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 16px;
}
.crumb {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 12px;
  color: rgba(0, 0, 0, 0.6);
  margin-bottom: 6px;
}
.crumb-link {
  cursor: pointer;
  user-select: none;
}
.crumb-link:hover {
  color: rgba(64, 158, 255, 0.95);
}
.crumb-link:focus-visible {
  outline: 2px solid rgba(64, 158, 255, 0.45);
  outline-offset: 2px;
  border-radius: 6px;
}
.crumb-sep {
  opacity: 0.6;
}
.crumb-current {
  color: rgba(0, 0, 0, 0.75);
}
.title {
  font-size: 30px;
  margin: 0;
  letter-spacing: -0.2px;
}
.sub {
  margin-top: 6px;
  color: rgba(0, 0, 0, 0.6);
  font-size: 12px;
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

.layout {
  margin-top: 14px;
  display: grid;
  grid-template-columns: minmax(0, 1fr);
  gap: 16px;
  align-items: start;
}
@media (max-width: 980px) {
  .layout {
    grid-template-columns: 1fr;
  }
}
.content {
  min-width: 0;
}
.side {
  display: grid;
  gap: 12px;
}
@media (max-width: 980px) {
  .side {
    position: static;
  }
}
.hero {
  border: 1px solid var(--color-border);
  border-radius: 16px;
  background: #fff;
  overflow: hidden;
}
.hero-media {
  background: rgba(0, 0, 0, 0.03);
}
.hero-empty {
  padding: 28px 16px;
  text-align: center;
}
.hero-viewer {
  display: grid;
  gap: 10px;
  padding: 12px;
}
.hero-stage {
  width: 100%;
  aspect-ratio: 16 / 9;
  border-radius: 14px;
  overflow: hidden;
  background: rgba(0, 0, 0, 0.05);
}
.hero-stage img,
.hero-stage video {
  width: 100%;
  height: 100%;
  object-fit: contain;
  display: block;
  background: #000;
}
.hero-thumbs {
  display: flex;
  gap: 10px;
  overflow-x: auto;
  padding-bottom: 6px;
}
.hero-thumb {
  width: 120px;
  height: 72px;
  flex: 0 0 auto;
  border: 1px solid rgba(0, 0, 0, 0.15);
  border-radius: 10px;
  background: rgba(255, 255, 255, 0.8);
  padding: 0;
  cursor: pointer;
}
.hero-thumb.active {
  border-color: rgba(64, 158, 255, 0.9);
  box-shadow: 0 0 0 2px rgba(64, 158, 255, 0.25);
}
.hero-thumb:focus-visible {
  outline: 2px solid rgba(64, 158, 255, 0.55);
  outline-offset: 2px;
}
.hero-thumb-inner {
  width: 100%;
  height: 100%;
  border-radius: 10px;
  overflow: hidden;
  background: rgba(0, 0, 0, 0.06);
}
.hero-thumb-inner img,
.hero-thumb-inner video {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
}
.hero-meta {
  padding: 14px 16px 16px;
}
.intro {
  color: rgba(0, 0, 0, 0.78);
  line-height: 1.7;
  font-size: 14px;
}
.tags {
  margin-top: 12px;
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}
.tag {
  max-width: 100%;
}
.link-item {
  display: grid;
  gap: 6px;
}
.link-editor {
  width: 100%;
  display: grid;
  gap: 8px;
}
.link-row {
  display: grid;
  grid-template-columns: minmax(0, 1fr) minmax(0, 2fr) auto;
  gap: 8px;
  align-items: center;
}
.links {
  margin-top: 12px;
  display: grid;
  gap: 6px;
}
.git {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border: 1px solid rgba(64, 158, 255, 0.35);
  color: rgba(64, 158, 255, 0.95);
  background: rgba(64, 158, 255, 0.06);
  padding: 8px 10px;
  border-radius: 12px;
  text-decoration: none;
  font-weight: 600;
}
.git:hover {
  background: rgba(64, 158, 255, 0.09);
}
.link-raw {
  font-size: 12px;
  word-break: break-all;
}

.section {
  margin-top: 12px;
  border: 1px solid var(--color-border);
  border-radius: 16px;
  background: #fff;
  padding: 14px 16px 16px;
}
.section-title {
  font-weight: 800;
  color: rgba(0, 0, 0, 0.82);
  margin-bottom: 10px;
}
.guide {
  line-height: 1.75;
}

.panel {
  border: 1px solid var(--color-border);
  border-radius: 16px;
  background: #fff;
  padding: 12px;
}
.panel-title {
  font-weight: 800;
  margin-bottom: 10px;
  color: rgba(0, 0, 0, 0.82);
}
.panel-body {
  min-width: 0;
}
.res-title {
  font-weight: 700;
  margin-bottom: 8px;
}
.kv {
  display: grid;
  grid-template-columns: 86px minmax(0, 1fr);
  gap: 8px 10px;
  align-items: start;
}
.k {
  color: rgba(0, 0, 0, 0.58);
  font-size: 12px;
  padding-top: 2px;
}
.v {
  color: rgba(0, 0, 0, 0.82);
  word-break: break-word;
  font-size: 13px;
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
  max-width: 100%;
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

.drawer-head {
  position: sticky;
  top: 0;
  z-index: 1;
  background: #fff;
  border-bottom: 1px solid var(--color-border);
  padding: 12px 12px 10px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}
.drawer-title {
  font-weight: 900;
  color: rgba(0, 0, 0, 0.82);
}
.drawer-actions {
  display: flex;
  gap: 8px;
  flex: 0 0 auto;
}
.drawer-body {
  padding: 12px;
  display: grid;
  gap: 12px;
}
.admin-panel {
  padding: 12px;
}
</style>

