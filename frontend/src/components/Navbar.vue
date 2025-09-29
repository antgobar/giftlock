<template>
  <nav class="navbar is-primary has-shadow" role="navigation" aria-label="main navigation">
    <div class="navbar-brand">
      <router-link to="/" class="navbar-item">
        <span class="icon-text">
          <span class="icon">
            <i class="fas fa-gift"></i>
          </span>
          <span class="title is-5 has-text-white">Gift Lock</span>
        </span>
      </router-link>
    </div>

    <div class="navbar-menu is-active">
      <div class="navbar-end">
        <!-- Unauthenticated user buttons -->
        <template v-if="!isAuthenticated">
          <router-link to="/login" class="navbar-item">
            <span class="icon-text">
              <span class="icon">
                <i class="fas fa-sign-in-alt"></i>
              </span>
              <span>Login</span>
            </span>
          </router-link>
          <router-link to="/register" class="navbar-item">
            <span class="icon-text">
              <span class="icon">
                <i class="fas fa-user-plus"></i>
              </span>
              <span>Register</span>
            </span>
          </router-link>
        </template>
        
        <!-- Authenticated user items -->
        <template v-if="isAuthenticated">
          <router-link to="/dashboard" class="navbar-item">
            <span class="icon-text">
              <span class="icon">
                <i class="fas fa-tachometer-alt"></i>
              </span>
              <span>Dashboard</span>
            </span>
          </router-link>
          <a class="navbar-item" @click="handleLogout">
            <span class="icon-text">
              <span class="icon">
                <i class="fas fa-sign-out-alt"></i>
              </span>
              <span>Logout</span>
            </span>
          </a>
        </template>
      </div>
    </div>
  </nav>
</template>

<script setup>
import axios from 'axios'
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const isAuthenticated = ref(false)
const username = ref('')

const checkAuthStatus = async () => {
  try {
    const response = await axios.get('/api/me')
    if (response.status === 200) {
      isAuthenticated.value = true
      username.value = response.data.username
    }
  } catch (error) {
    isAuthenticated.value = false
    username.value = ''
  }
}

const handleLogout = async () => {
  try {
    await axios.post('/api/logout')
  } catch (error) {
    console.error('Logout error:', error)
  } finally {
    isAuthenticated.value = false
    username.value = ''
    router.push('/')
  }
}

// Check auth status on component mount
onMounted(() => {
  checkAuthStatus()
})

// Listen for route changes to update auth status
router.afterEach(() => {
  checkAuthStatus()
})

// Export functions for parent components to use
defineExpose({
  updateAuthStatus: checkAuthStatus
})
</script>

