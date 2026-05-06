import { defineStore } from 'pinia'
import { useStorage } from '@vueuse/core'
import { login } from '@/api/user'
import { ElLoading, ElMessage } from 'element-plus'
import router from '@/router'

export const useExampleStore = defineStore('example', () => {
    const token = useStorage('token', '')
    return {

    }
})