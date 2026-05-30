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
    <div class="header">
      <h1 class="title">作者：睡觉王</h1>
      <div class="actions">
        <el-button type="primary" :loading="loading" @click="load">刷新</el-button>
      </div>
    </div>
    <el-alert v-if="errMsg" type="error" :title="errMsg" show-icon :closable="false" />

    <el-empty v-if="!loading && !projects.length && !errMsg" description="暂无数据" />

    <section v-else class="cards">
      <div
        v-for="it in projects"
        :key="getFolderName(it)"
        class="card"
        role="button"
        tabindex="0"
        @click="router.push(`/projects/${getFolderName(it)}`)"
        @keydown.enter.prevent="router.push(`/projects/${getFolderName(it)}`)"
      >
        <div class="cover">
          <img :src="getCover(it)" alt="" />
        </div>
        <div class="body">
          <div class="name" :title="getName(it)">{{ getName(it) }}</div>
          <div class="intro" :title="getIntro(it)">{{ getIntro(it) || '-' }}</div>
          <div v-if="getTags(it).length" class="tags">
            <el-tag v-for="t in getTags(it)" :key="t" size="small" class="tag">{{ t }}</el-tag>
          </div>
        </div>
      </div>
    </section>
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
  align-items: flex-end;
  justify-content: space-between;
  gap: 16px;
}
.title {
  font-size: 28px;
  margin: 0;
}
.desc {
  margin: 8px 0 16px;
  color: rgba(0, 0, 0, 0.65);
}
.cards {
  margin-top: 12px;
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 14px;
}
@media (max-width: 980px) {
  .cards {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}
@media (max-width: 640px) {
  .cards {
    grid-template-columns: 1fr;
  }
}
.card {
  border: 1px solid var(--color-border);
  border-radius: 14px;
  background: #fff;
  overflow: hidden;
  cursor: pointer;
  transition: transform 0.08s ease, box-shadow 0.08s ease;
}
.card:hover {
  transform: translateY(-1px);
  box-shadow: 0 8px 18px rgba(0, 0, 0, 0.06);
}
.card:focus-visible {
  outline: 2px solid rgba(64, 158, 255, 0.6);
  outline-offset: 2px;
}
.cover {
  height: 160px;
  background: rgba(0, 0, 0, 0.03);
}
.cover img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
}
.body {
  padding: 12px;
}
.name {
  font-weight: 800;
  color: rgba(0, 0, 0, 0.82);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  margin-bottom: 6px;
}
.intro {
  color: rgba(0, 0, 0, 0.7);
  font-size: 13px;
  line-height: 1.5;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  min-height: 39px;
}
.tags {
  margin-top: 10px;
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}
.tag {
  max-width: 100%;
}
</style>

