<template>
  <nav class="navbar is-primary has-shadow" role="navigation" aria-label="main navigation">
    <div class="navbar-brand">
      <router-link to="/" class="navbar-item brand-link">
        <span class="icon-text">
          <span class="icon">
            <span class="gift-icon">🎁</span>
          </span>
          <span class="brand-text">Gift Lock</span>
        </span>
      </router-link>

      <a 
        role="button" 
        class="navbar-burger" 
        :class="{ 'is-active': isMenuOpen }" 
        aria-label="menu" 
        :aria-expanded="isMenuOpen"
        @click="toggleMenu"
      >
        <span aria-hidden="true"></span>
        <span aria-hidden="true"></span>
        <span aria-hidden="true"></span>
      </a>
    </div>

    <div class="navbar-menu" :class="{ 'is-active': isMenuOpen }">
      <div class="navbar-start">
        <router-link to="/" class="navbar-item" @click="closeMenu">
          <span class="icon-text">
            <span class="icon">
              <i class="fas fa-home"></i>
            </span>
            <span>Home</span>
          </span>
        </router-link>
        
        <router-link 
          to="/dashboard" 
          class="navbar-item" 
          v-if="isAuthenticated" 
          @click="closeMenu"
        >
          <span class="icon-text">
            <span class="icon">
              <i class="fas fa-tachometer-alt"></i>
            </span>
            <span>Dashboard</span>
          </span>
        </router-link>
      </div>

      <div class="navbar-end">
        <!-- Desktop Auth Buttons -->
        <div class="navbar-item auth-buttons is-hidden-touch" v-if="!isAuthenticated">
          <div class="buttons">
            <router-link to="/register" class="button is-success is-outlined" @click="closeMenu">
              <span class="icon">
                <i class="fas fa-user-plus"></i>
              </span>
              <span>Sign up</span>
            </router-link>
            <router-link to="/login" class="button is-light" @click="closeMenu">
              <span class="icon">
                <i class="fas fa-sign-in-alt"></i>
              </span>
              <span>Log in</span>
            </router-link>
          </div>
        </div>
        
        <!-- Mobile Auth Links (shown in burger menu) -->
        <template v-if="!isAuthenticated">
          <router-link to="/login" class="navbar-item is-hidden-desktop" @click="closeMenu">
            <span class="icon-text">
              <span class="icon">
                <i class="fas fa-sign-in-alt"></i>
              </span>
              <span>Log in</span>
            </span>
          </router-link>
          <router-link to="/register" class="navbar-item is-hidden-desktop" @click="closeMenu">
            <span class="icon-text">
              <span class="icon">
                <i class="fas fa-user-plus"></i>
              </span>
              <span>Sign up</span>
            </span>
          </router-link>
        </template>
        
        <!-- Desktop User Dropdown -->
        <div class="navbar-item has-dropdown is-hoverable is-hidden-touch" v-if="isAuthenticated">
          <a class="navbar-link user-link">
            <span class="icon-text">
              <span class="icon">
                <i class="fas fa-user-circle"></i>
              </span>
              <span>{{ username || 'User' }}</span>
            </span>
          </a>
          <div class="navbar-dropdown">
            <router-link to="/dashboard" class="navbar-item" @click="closeMenu">
              <span class="icon-text">
                <span class="icon">
                  <i class="fas fa-tachometer-alt"></i>
                </span>
                <span>Dashboard</span>
              </span>
            </router-link>
            <hr class="navbar-divider">
            <a class="navbar-item" @click="handleLogout">
              <span class="icon-text">
                <span class="icon">
                  <i class="fas fa-sign-out-alt"></i>
                </span>
                <span>Logout</span>
              </span>
            </a>
          </div>
        </div>

        <!-- Mobile User Menu (shown in burger menu) -->
        <template v-if="isAuthenticated">
          <div class="navbar-item is-hidden-desktop user-info">
            <span class="icon-text">
              <span class="icon">
                <i class="fas fa-user-circle"></i>
              </span>
              <span>{{ username || 'User' }}</span>
            </span>
          </div>
          <hr class="navbar-divider is-hidden-desktop">
          <a class="navbar-item is-hidden-desktop" @click="handleLogout">
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
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%) !important;
  border-bottom: 3px solid rgba(255, 255, 255, 0.1);
}

.brand-link {
  font-weight: bold;
  font-size: 1.25rem;
}

.brand-text {
  color: white !important;
  font-weight: 700;
  letter-spacing: 0.5px;
}

.gift-icon {
  font-size: 1.5em;
  margin-right: 0.25rem;
}

.navbar-item {
  color: white !important;
  transition: background-color 0.3s ease;
}

.navbar-item:hover {
  background-color: rgba(255, 255, 255, 0.1) !important;
  color: white !important;
}

.navbar-link {
  color: white !important;
}

.navbar-link:hover {
  background-color: rgba(255, 255, 255, 0.1) !important;
  color: white !important;
}

.navbar-burger {
  color: white !important;
}

.navbar-burger:hover {
  background-color: rgba(255, 255, 255, 0.1) !important;
}

.auth-buttons .button {
  margin-left: 0.5rem;
}

.button.is-success.is-outlined {
  border-color: white;
  color: white;
}

.button.is-success.is-outlined:hover {
  background-color: white;
  color: #667eea;
  border-color: white;
}

.button.is-light {
  background-color: rgba(255, 255, 255, 0.9);
  color: #667eea;
  border: none;
}

.button.is-light:hover {
  background-color: white;
  color: #667eea;
}

.user-info {
  font-weight: 500;
}

.navbar-dropdown {
  background-color: white;
  box-shadow: 0 8px 16px rgba(0, 0, 0, 0.1);
  border-radius: 6px;
  margin-top: 0.5rem;
}

.navbar-dropdown .navbar-item {
  color: #363636 !important;
}

.navbar-dropdown .navbar-item:hover {
  background-color: #f5f5f5 !important;
  color: #667eea !important;
}

.navbar-divider {
  background-color: #dbdbdb;
}

/* Mobile specific styles */
@media screen and (max-width: 1023px) {
  .navbar-menu {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    box-shadow: 0 8px 16px rgba(0, 0, 0, 0.1);
  }
  
  .navbar-menu.is-active {
    display: block;
  }
  
  .navbar-divider.is-hidden-desktop {
    background-color: rgba(255, 255, 255, 0.2);
    margin: 0.5rem 0;
  }
}
</style>