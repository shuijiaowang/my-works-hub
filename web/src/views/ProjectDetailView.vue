<script setup>
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { adminLoggedIn } from '@/adminAuth'
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
  canEdit.value = adminLoggedIn.value && !!localStorage.getItem('admin_token')
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
        <div v-if="tagsList.length" class="header-tags">
          <span v-for="t in tagsList" :key="t" class="header-tag">{{ t }}</span>
        </div>
      </div>

      <div class="actions">
        <button class="back-btn" @click="router.push('/projects')">← 返回</button>
        <el-button v-if="canEdit" type="primary" size="small" :disabled="loading || !project" @click="enterEdit">编辑</el-button>
      </div>
    </div>

    <el-alert v-if="errMsg" type="error" :title="errMsg" show-icon :closable="false" style="margin-top: 12px" />

    <el-skeleton v-if="loading" :rows="8" animated style="margin-top: 12px" />

    <section v-else-if="project" class="layout">
      <div class="hero-card">
        <div class="hero-media">
          <div v-if="mediaLoading" class="hero-empty">加载中…</div>
          <div v-else-if="!mediaItems?.length" class="hero-empty">暂无预览</div>
          <div v-else class="hero-viewer">
            <div class="hero-stage">
              <video v-if="activeMedia?.kind === 'video'" :src="activeMedia?.url" controls playsinline />
              <img v-else :src="activeMedia?.url" :alt="projectName" loading="lazy" />
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
          <p v-if="projectIntro" class="intro">{{ projectIntro }}</p>
          <p v-else class="intro empty-text">暂无简介</p>

          <div v-if="projectLinks.length" class="links">
            <a
              v-for="(link, i) in projectLinks"
              :key="i"
              class="link-btn"
              :href="link.url"
              target="_blank"
              rel="noreferrer"
            >
              <span class="link-icon">↗</span>
              {{ link.name || link.url }}
            </a>
          </div>
        </div>
      </div>

      <div v-if="zipFiles?.length" class="section">
        <h2 class="section-title">
          <span class="section-icon">⬇</span>
          资源下载
        </h2>
        <div class="files">
          <div v-for="f in zipFiles" :key="f" class="file-row">
            <span class="file-name" :title="f">{{ f }}</span>
            <button class="download-btn" :disabled="zipBusy" @click="onDownloadZip(f)">下载</button>
          </div>
        </div>
      </div>

      <div v-if="projectGuide" class="section">
        <h2 class="section-title">
          <span class="section-icon">📖</span>
          教程
        </h2>
        <div class="guide">{{ projectGuide }}</div>
      </div>
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
  max-width: 880px;
  margin: 0 auto;
  padding: 28px 24px 64px;
}

.header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 16px;
  margin-bottom: 28px;
}

.crumb {
  display: flex;
  align-items: center;
  gap: 6px;
  font-family: var(--portfolio-font-mono);
  font-size: 12px;
  color: var(--portfolio-text-muted);
  margin-bottom: 10px;
}
.crumb-link {
  cursor: pointer;
  user-select: none;
  transition: color 0.15s;
}
.crumb-link:hover {
  color: var(--portfolio-accent);
}
.crumb-link:focus-visible {
  outline: 2px solid var(--portfolio-accent);
  outline-offset: 2px;
  border-radius: 4px;
}
.crumb-sep {
  opacity: 0.5;
}
.crumb-current {
  color: var(--portfolio-text-secondary);
}

.title {
  font-size: clamp(22px, 3.5vw, 28px);
  font-weight: 700;
  letter-spacing: -0.03em;
  color: var(--portfolio-text);
  margin: 0;
  line-height: 1.25;
}

.header-tags {
  margin-top: 10px;
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}
.header-tag {
  font-family: var(--portfolio-font-mono);
  font-size: 11px;
  font-weight: 500;
  padding: 3px 8px;
  border-radius: 6px;
  background: var(--portfolio-accent-soft);
  color: var(--portfolio-accent);
  border: 1px solid rgba(37, 99, 235, 0.15);
}

.actions {
  display: flex;
  gap: 8px;
  align-items: center;
  flex-shrink: 0;
}
.back-btn {
  padding: 7px 12px;
  border: 1px solid var(--portfolio-border);
  border-radius: var(--portfolio-radius);
  background: var(--portfolio-surface);
  color: var(--portfolio-text-secondary);
  font-size: 13px;
  cursor: pointer;
  transition: border-color 0.15s, box-shadow 0.15s;
}
.back-btn:hover {
  border-color: var(--color-border-hover);
  box-shadow: var(--portfolio-shadow-sm);
}

.layout {
  display: grid;
  gap: 20px;
}

.hero-card {
  border: 1px solid var(--portfolio-border);
  border-radius: var(--portfolio-radius-lg);
  background: var(--portfolio-surface);
  overflow: hidden;
  box-shadow: var(--portfolio-shadow-sm);
  display: grid;
  grid-template-columns: minmax(0, 1.1fr) minmax(0, 0.9fr);
}
@media (max-width: 720px) {
  .hero-card {
    grid-template-columns: 1fr;
  }
}

.hero-media {
  background: var(--portfolio-border-light);
  border-right: 1px solid var(--portfolio-border);
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 200px;
}
@media (max-width: 720px) {
  .hero-media {
    border-right: none;
    border-bottom: 1px solid var(--portfolio-border);
  }
}

.hero-empty {
  padding: 48px 20px;
  text-align: center;
  color: var(--portfolio-text-muted);
  font-size: 13px;
  width: 100%;
}

.hero-viewer {
  width: 100%;
  padding: 16px;
  display: grid;
  gap: 10px;
}

.hero-stage {
  width: 100%;
  max-height: 280px;
  border-radius: var(--portfolio-radius);
  overflow: hidden;
  background: #0f172a;
  display: flex;
  align-items: center;
  justify-content: center;
}
.hero-stage img,
.hero-stage video {
  max-width: 100%;
  max-height: 280px;
  width: auto;
  height: auto;
  object-fit: contain;
  display: block;
}

.hero-thumbs {
  display: flex;
  gap: 8px;
  overflow-x: auto;
  padding-bottom: 2px;
}
.hero-thumb {
  width: 72px;
  height: 48px;
  flex: 0 0 auto;
  border: 2px solid transparent;
  border-radius: 8px;
  background: var(--portfolio-surface);
  padding: 0;
  cursor: pointer;
  overflow: hidden;
  transition: border-color 0.15s;
}
.hero-thumb.active {
  border-color: var(--portfolio-accent);
}
.hero-thumb:focus-visible {
  outline: 2px solid var(--portfolio-accent);
  outline-offset: 2px;
}
.hero-thumb-inner {
  width: 100%;
  height: 100%;
  overflow: hidden;
  background: var(--portfolio-border-light);
}
.hero-thumb-inner img,
.hero-thumb-inner video {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
}

.hero-meta {
  padding: 24px 24px 28px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  gap: 16px;
}
.intro {
  color: var(--portfolio-text-secondary);
  line-height: 1.75;
  font-size: 14px;
  margin: 0;
}
.empty-text {
  color: var(--portfolio-text-muted);
  font-style: italic;
}

.links {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}
.link-btn {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 8px 14px;
  border: 1px solid var(--portfolio-border);
  border-radius: var(--portfolio-radius);
  background: var(--portfolio-surface);
  color: var(--portfolio-accent);
  font-size: 13px;
  font-weight: 500;
  text-decoration: none;
  transition: border-color 0.15s, background 0.15s, box-shadow 0.15s;
}
.link-btn:hover {
  border-color: var(--portfolio-accent);
  background: var(--portfolio-accent-soft);
  box-shadow: var(--portfolio-shadow-sm);
}
.link-icon {
  font-size: 12px;
  opacity: 0.7;
}

.section {
  border: 1px solid var(--portfolio-border);
  border-radius: var(--portfolio-radius-lg);
  background: var(--portfolio-surface);
  padding: 20px 24px;
  box-shadow: var(--portfolio-shadow-sm);
}
.section-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  font-weight: 600;
  color: var(--portfolio-text);
  margin: 0 0 14px;
  letter-spacing: -0.01em;
}
.section-icon {
  font-size: 14px;
  opacity: 0.7;
}

.guide {
  line-height: 1.8;
  font-size: 14px;
  color: var(--portfolio-text-secondary);
  white-space: pre-wrap;
  font-family: var(--portfolio-font-mono);
  font-size: 13px;
  background: var(--portfolio-border-light);
  padding: 16px;
  border-radius: var(--portfolio-radius);
  border: 1px solid var(--portfolio-border);
}

.files {
  display: grid;
  gap: 8px;
}
.file-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 10px 14px;
  border-radius: var(--portfolio-radius);
  background: var(--portfolio-border-light);
  border: 1px solid var(--portfolio-border);
}
.file-name {
  font-family: var(--portfolio-font-mono);
  font-size: 12px;
  color: var(--portfolio-text-secondary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.download-btn {
  flex-shrink: 0;
  padding: 5px 12px;
  border: 1px solid var(--portfolio-accent);
  border-radius: 6px;
  background: var(--portfolio-accent-soft);
  color: var(--portfolio-accent);
  font-size: 12px;
  font-weight: 500;
  cursor: pointer;
  transition: background 0.15s;
}
.download-btn:hover:not(:disabled) {
  background: rgba(37, 99, 235, 0.15);
}
.download-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.muted {
  color: var(--portfolio-text-muted);
}

.form {
  max-width: 100%;
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
  max-width: 100%;
}
.media-card {
  width: 220px;
  flex: 0 0 auto;
  border: 1px solid var(--portfolio-border);
  border-radius: var(--portfolio-radius);
  background: var(--portfolio-surface);
  overflow: hidden;
}
.thumb {
  height: 140px;
  background: var(--portfolio-border-light);
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
  color: var(--portfolio-text-muted);
  word-break: break-all;
  margin-bottom: 8px;
}
.media-ops {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.panel {
  border: 1px solid var(--portfolio-border);
  border-radius: var(--portfolio-radius);
  background: var(--portfolio-surface);
  padding: 12px;
}
.panel-title {
  font-weight: 600;
  margin-bottom: 10px;
  color: var(--portfolio-text);
  font-size: 13px;
}
.panel-body {
  min-width: 0;
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
  font-family: var(--portfolio-font-mono);
  font-size: 12px;
}
.zip-ops {
  display: flex;
  gap: 8px;
  flex: 0 0 auto;
}
.file {
  padding: 6px 8px;
  border-radius: 8px;
  background: var(--portfolio-border-light);
  font-family: var(--portfolio-font-mono);
  font-size: 12px;
}
.files {
  display: grid;
  gap: 6px;
}

.drawer-head {
  position: sticky;
  top: 0;
  z-index: 1;
  background: var(--portfolio-surface);
  border-bottom: 1px solid var(--portfolio-border);
  padding: 12px 12px 10px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}
.drawer-title {
  font-weight: 700;
  color: var(--portfolio-text);
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
  min-width: 0;
}
.admin-panel {
  padding: 12px;
  min-width: 0;
  overflow: hidden;
}
:deep(.admin-drawer .el-drawer__body) {
  overflow-x: hidden;
}
</style>

