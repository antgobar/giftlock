<template>
  <nav class="navbar is-dark" role="navigation" aria-label="main navigation">
    <div class="navbar-brand">
      <router-link to="/" class="navbar-item">
        <strong>🎁 Gift Lock</strong>
      </router-link>

      <a role="button" class="navbar-burger" :class="{ 'is-active': isMenuOpen }" aria-label="menu" aria-expanded="false" @click="toggleMenu">
        <span aria-hidden="true"></span>
        <span aria-hidden="true"></span>
        <span aria-hidden="true"></span>
      </a>
    </div>

    <div class="navbar-menu" :class="{ 'is-active': isMenuOpen }">
      <div class="navbar-start">
        <router-link to="/" class="navbar-item" @click="closeMenu">
          Home
        </router-link>
        <router-link to="/dashboard" class="navbar-item" v-if="isAuthenticated" @click="closeMenu">
          Dashboard
        </router-link>
      </div>

      <div class="navbar-end">
        <div class="navbar-item" v-if="!isAuthenticated">
          <div class="buttons">
            <router-link to="/register" class="button is-primary" @click="closeMenu">
              <strong>Sign up</strong>
            </router-link>
            <router-link to="/login" class="button is-light" @click="closeMenu">
              Log in
            </router-link>
          </div>
        </div>
        
        <div class="navbar-item has-dropdown is-hoverable" v-if="isAuthenticated">
          <a class="navbar-link">
            {{ username || 'User' }}
          </a>
          <div class="navbar-dropdown">
            <router-link to="/dashboard" class="navbar-item" @click="closeMenu">
              Dashboard
            </router-link>
            <hr class="navbar-divider">
            <a class="navbar-item" @click="handleLogout">
              Logout
            </a>
          </div>
        </div>
      </div>
    </div>
  </nav>
</template>

<script setup>
import axios from 'axios'
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const isMenuOpen = ref(false)
const isAuthenticated = ref(false)
const username = ref('')

const toggleMenu = () => {
  isMenuOpen.value = !isMenuOpen.value
}

const closeMenu = () => {
  isMenuOpen.value = false
}

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
    closeMenu()
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

<style scoped>
.navbar {
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.navbar-brand .navbar-item {
  font-size: 1.25rem;
}

.navbar-item.router-link-active:not(.button) {
  background-color: rgba(255, 255, 255, 0.1);
}

@media screen and (max-width: 1023px) {
  .navbar-menu {
    box-shadow: 0 8px 16px rgba(0, 0, 0, 0.1);
  }
}
</style>