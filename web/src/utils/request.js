import axios from 'axios'
import { ElLoading, ElMessage } from 'element-plus'

// 创建axios实例
const service = axios.create({
    // 生产环境未加载 .env.production 时 VITE_BASE_API 为空，会错误请求 /user/login；与后端 /api 对齐
    baseURL: import.meta.env.VITE_BASE_API || '/api',
    timeout: 10000,
    headers: {
        'Content-Type': 'application/json'
    }
})

// 加载状态管理（解决多请求冲突问题）
let loadingCount = 0
let loadingInstance = null

// 显示加载
const showLoading = (config) => {
    if (config.donNotShowLoading) return

    loadingCount++
    if (loadingCount === 1) {
        loadingInstance = ElLoading.service({
            fullscreen: true,
            text: '加载中...',
            background: 'rgba(0, 0, 0, 0.1)'
        })
    }
}

// 隐藏加载
const hideLoading = (config) => {
    if (config.donNotShowLoading) return

    loadingCount--
    if (loadingCount <= 0) {
        loadingInstance?.close()
        loadingInstance = null
        loadingCount = 0
    }
}

// 请求拦截器
service.interceptors.request.use(
    (config) => {
        showLoading(config)
        // 仅管理端请求携带 admin token（单一固定 token）
        const adminToken = localStorage.getItem('admin_token') || ''
        const url = config.url || ''
        if (adminToken && typeof url === 'string' && url.startsWith('/admin/')) {
            config.headers['X-Admin-Token'] = adminToken
        }
        return config
    },
    (error) => {
        hideLoading(error.config)
        ElMessage.error(`请求准备失败: ${error.message || '未知错误'}`)
        return Promise.reject(error)
    }
)

// 响应拦截器
service.interceptors.response.use(
    (response) => {
        hideLoading(response.config)
        const res = response.data

        // 处理业务状态码
        if (typeof res.code !== 'undefined') {
            // 成功状态码（根据实际业务调整）
            if (res.code === 0) {
                return res
            } else {
                ElMessage.error(res.msg || `操作失败（${res.code}）`)
                return Promise.reject(res)
            }
        }

        return res
    },
    (error) => {
        hideLoading(error.config)

        // 网络错误处理
        if (!error.response) {
            ElMessage.error('网络连接异常，请检查网络')
            return Promise.reject(error)
        }

        // 401 未授权处理
        if (error.response.status === 401) {
            ElMessage.error('管理端校验失败（401）')
        } else {
            // 其他HTTP错误
            const status = error.response.status
            const msg = error.response.data?.msg || `请求错误（${status}）`
            ElMessage.error(msg)
        }

        return Promise.reject(error)
    }
)

export default service