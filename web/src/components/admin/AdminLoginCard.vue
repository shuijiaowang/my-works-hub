<script setup>
import { computed, ref } from 'vue'
import { setAdminLoggedIn } from '@/adminAuth'
import service from '@/utils/request'

const emit = defineEmits(['logged-in'])

const tokenInUrl = new URLSearchParams(location.search).get('token') || ''
const token = ref(localStorage.getItem('admin_token') || tokenInUrl)
const password = ref('')
const status = ref('')
const saving = ref(false)

const hasToken = computed(() => !!token.value)

const saveToken = () => {
  localStorage.setItem('admin_token', token.value || '')
}

const login = async () => {
  if (!token.value) {
    status.value = '缺少 token（请用 /admin?token=xxx 打开）'
    return
  }
  saving.value = true
  status.value = ''
  try {
    const res = await service.post(`/admin/login?token=${encodeURIComponent(token.value)}`, {
      password: password.value,
    })
    status.value = res?.msg || '登录成功'
    saveToken()
    setAdminLoggedIn(true)
    emit('logged-in')
  } catch (e) {
    status.value = e?.msg || e?.message || '登录失败'
  } finally {
    saving.value = false
  }
}

</script>

<template>
  <section class="card">
    <p class="desc">
      管理端请求会校验固定 token。首次可通过 <code>/admin?token=xxx</code> 进入并登录。
    </p>
    <el-form label-width="120px">
      <el-form-item label="Admin Token">
        <el-input v-model="token" placeholder="token（固定一个）" />
      </el-form-item>
      <el-form-item label="登录密码">
        <el-input v-model="password" type="password" show-password placeholder="password" />
      </el-form-item>
      <el-form-item>
        <el-button type="primary" :disabled="!hasToken" :loading="saving" @click="login">登录</el-button>
      </el-form-item>
    </el-form>

    <div v-if="status" class="status">{{ status }}</div>
  </section>
</template>

<style scoped>
.desc {
  margin: 0 0 12px;
  color: rgba(0, 0, 0, 0.65);
}
.card {
  padding: 16px;
  border: 1px solid var(--color-border);
  border-radius: 12px;
  background: #fff;
}
.status {
  margin-top: 8px;
  color: rgba(0, 0, 0, 0.75);
  word-break: break-all;
}
</style>
