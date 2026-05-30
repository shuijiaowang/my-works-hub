<script setup>
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { fetchAllProjects } from '@/api/projects'

const router = useRouter()

const loading = ref(false)
const errMsg = ref('')
const data = ref(null)

const projects = computed(() => data.value?.data?.projects || [])

const DEFAULT_COVER = '/api/resources/未上传.png'

const getCover = (it) => {
  return it?.coverUrl || DEFAULT_COVER
}

const getFolderName = (it) => it?.project?.FolderName ?? it?.project?.folderName ?? ''

const getName = (it) => it?.project?.Name ?? it?.project?.name ?? ''

const getIntro = (it) => it?.project?.Intro ?? it?.project?.intro ?? ''

const getTags = (it) => {
  const raw = it?.project?.Tags ?? it?.project?.tags ?? ''
  return String(raw)
    .split(',')
    .map((s) => s.trim())
    .filter(Boolean)
}

const load = async () => {
  loading.value = true
  errMsg.value = ''
  try {
    data.value = await fetchAllProjects()
  } catch (e) {
    errMsg.value = e?.msg || e?.message || '请求失败'
  } finally {
    loading.value = false
  }
}

onMounted(load)
</script>

<template>
  <main class="page">
    <header class="hero">
      <div class="hero-content">
        <p class="hero-label">Portfolio</p>
        <h1 class="hero-title">睡觉王</h1>
        <p class="hero-desc">程序开发作品集 · 记录与分享个人项目</p>
      </div>
      <button class="refresh-btn" :disabled="loading" @click="load">
        <span class="refresh-icon" :class="{ spinning: loading }">↻</span>
        刷新
      </button>
    </header>

    <el-alert v-if="errMsg" type="error" :title="errMsg" show-icon :closable="false" class="alert" />

    <div v-if="loading && !projects.length" class="loading-state">
      <div v-for="n in 6" :key="n" class="skeleton-card">
        <div class="skeleton-cover" />
        <div class="skeleton-body">
          <div class="skeleton-line wide" />
          <div class="skeleton-line" />
          <div class="skeleton-line short" />
        </div>
      </div>
    </div>

    <el-empty v-else-if="!loading && !projects.length && !errMsg" description="暂无项目" class="empty" />

    <section v-else class="cards">
      <article
        v-for="it in projects"
        :key="getFolderName(it)"
        class="card"
        role="button"
        tabindex="0"
        @click="router.push(`/projects/${getFolderName(it)}`)"
        @keydown.enter.prevent="router.push(`/projects/${getFolderName(it)}`)"
      >
        <div class="cover">
          <img :src="getCover(it)" :alt="getName(it)" loading="lazy" />
          <div class="cover-overlay" />
        </div>
        <div class="body">
          <h2 class="name" :title="getName(it)">{{ getName(it) }}</h2>
          <p class="intro" :title="getIntro(it)">{{ getIntro(it) || '暂无简介' }}</p>
          <div v-if="getTags(it).length" class="tags">
            <span v-for="t in getTags(it)" :key="t" class="tag">{{ t }}</span>
          </div>
        </div>
      </article>
    </section>
  </main>
</template>

<style scoped>
.page {
  max-width: 1120px;
  margin: 0 auto;
  padding: 32px 24px 64px;
}

.hero {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  gap: 24px;
  padding-bottom: 32px;
  border-bottom: 1px solid var(--portfolio-border);
  margin-bottom: 32px;
}

.hero-label {
  font-family: var(--portfolio-font-mono);
  font-size: 12px;
  font-weight: 500;
  letter-spacing: 0.12em;
  text-transform: uppercase;
  color: var(--portfolio-accent);
  margin-bottom: 8px;
}

.hero-title {
  font-size: clamp(28px, 4vw, 36px);
  font-weight: 700;
  letter-spacing: -0.03em;
  color: var(--portfolio-text);
  margin: 0;
  line-height: 1.2;
}

.hero-desc {
  margin-top: 10px;
  font-size: 14px;
  color: var(--portfolio-text-muted);
  max-width: 420px;
}

.refresh-btn {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 8px 14px;
  border: 1px solid var(--portfolio-border);
  border-radius: var(--portfolio-radius);
  background: var(--portfolio-surface);
  color: var(--portfolio-text-secondary);
  font-size: 13px;
  cursor: pointer;
  transition: border-color 0.15s, box-shadow 0.15s;
  flex-shrink: 0;
}
.refresh-btn:hover:not(:disabled) {
  border-color: var(--color-border-hover);
  box-shadow: var(--portfolio-shadow-sm);
}
.refresh-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}
.refresh-icon {
  display: inline-block;
  font-size: 16px;
  line-height: 1;
}
.refresh-icon.spinning {
  animation: spin 0.8s linear infinite;
}
@keyframes spin {
  to { transform: rotate(360deg); }
}

.alert {
  margin-bottom: 24px;
}

.loading-state {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 20px;
}
@media (max-width: 980px) {
  .loading-state { grid-template-columns: repeat(2, minmax(0, 1fr)); }
}
@media (max-width: 640px) {
  .loading-state { grid-template-columns: 1fr; }
}

.skeleton-card {
  border: 1px solid var(--portfolio-border);
  border-radius: var(--portfolio-radius-lg);
  overflow: hidden;
  background: var(--portfolio-surface);
}
.skeleton-cover {
  height: 168px;
  background: linear-gradient(90deg, #eef1f5 25%, #f8f9fb 50%, #eef1f5 75%);
  background-size: 200% 100%;
  animation: shimmer 1.4s infinite;
}
.skeleton-body {
  padding: 16px;
  display: grid;
  gap: 8px;
}
.skeleton-line {
  height: 12px;
  border-radius: 4px;
  background: linear-gradient(90deg, #eef1f5 25%, #f8f9fb 50%, #eef1f5 75%);
  background-size: 200% 100%;
  animation: shimmer 1.4s infinite;
}
.skeleton-line.wide { width: 70%; height: 16px; }
.skeleton-line.short { width: 45%; }
@keyframes shimmer {
  0% { background-position: 200% 0; }
  100% { background-position: -200% 0; }
}

.empty {
  padding: 48px 0;
}

.cards {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 20px;
}
@media (max-width: 980px) {
  .cards { grid-template-columns: repeat(2, minmax(0, 1fr)); }
}
@media (max-width: 640px) {
  .cards { grid-template-columns: 1fr; }
  .hero { flex-direction: column; align-items: flex-start; }
}

.card {
  border: 1px solid var(--portfolio-border);
  border-radius: var(--portfolio-radius-lg);
  background: var(--portfolio-surface);
  overflow: hidden;
  cursor: pointer;
  transition: transform 0.2s ease, box-shadow 0.2s ease, border-color 0.2s ease;
}
.card:hover {
  transform: translateY(-3px);
  box-shadow: var(--portfolio-shadow-hover);
  border-color: var(--color-border-hover);
}
.card:focus-visible {
  outline: 2px solid var(--portfolio-accent);
  outline-offset: 2px;
}

.cover {
  position: relative;
  height: 168px;
  background: var(--portfolio-border-light);
  overflow: hidden;
}
.cover img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
  transition: transform 0.3s ease;
}
.card:hover .cover img {
  transform: scale(1.03);
}
.cover-overlay {
  position: absolute;
  inset: 0;
  background: linear-gradient(to top, rgba(15, 23, 42, 0.04), transparent 40%);
  pointer-events: none;
}

.body {
  padding: 16px 16px 18px;
}
.name {
  font-size: 15px;
  font-weight: 600;
  color: var(--portfolio-text);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  margin: 0 0 6px;
  letter-spacing: -0.01em;
}
.intro {
  color: var(--portfolio-text-muted);
  font-size: 13px;
  line-height: 1.55;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  min-height: 40px;
  margin: 0;
}
.tags {
  margin-top: 12px;
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}
.tag {
  font-family: var(--portfolio-font-mono);
  font-size: 11px;
  font-weight: 500;
  padding: 3px 8px;
  border-radius: 6px;
  background: var(--portfolio-accent-soft);
  color: var(--portfolio-accent);
  border: 1px solid rgba(37, 99, 235, 0.15);
}
</style>
