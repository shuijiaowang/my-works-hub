import { ref } from 'vue'

export const adminLoggedIn = ref(sessionStorage.getItem('admin_logged_in') === '1')

export function setAdminLoggedIn(value) {
  adminLoggedIn.value = value
  if (value) {
    sessionStorage.setItem('admin_logged_in', '1')
  } else {
    sessionStorage.removeItem('admin_logged_in')
  }
}
