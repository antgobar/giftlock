<template>
  <div class="login-container">
    <div class="login-form">
      <h2 class="title is-3">Sign In</h2>
      
      <form @submit.prevent="handleLogin">
        <div class="field">
          <label class="label" for="username">Username</label>
          <div class="control">
            <input
              class="input"
              type="text"
              id="username"
              v-model="credentials.username"
              required
              :disabled="isLoading"
              placeholder="Enter your username"
            />
          </div>
        </div>
        
        <div class="field">
          <label class="label" for="password">Password</label>
          <div class="control">
            <input
              class="input"
              type="password"
              id="password"
              v-model="credentials.password"
              required
              :disabled="isLoading"
              placeholder="Enter your password"
            />
          </div>
        </div>
        
        <div class="field">
          <div class="control">
            <button class="button is-primary is-fullwidth" type="submit" :disabled="isLoading">
              <span v-if="isLoading">
                <i class="fas fa-spinner fa-spin"></i>
                Signing in...
              </span>
              <span v-else>
                Sign In
              </span>
            </button>
          </div>
        </div>
      </form>
      
      <div class="notification is-danger" v-if="errorMessage">
        {{ errorMessage }}
      </div>
      
      <div class="notification is-success" v-if="successMessage">
        {{ successMessage }}
      </div>

      <div class="has-text-centered mt-4">
        <p>Don't have an account? 
          <router-link to="/register" class="has-text-primary">Sign up here</router-link>
        </p>
      </div>
    </div>
  </div>
</template>

<script setup>
import axios from 'axios'
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

const credentials = reactive({
  username: '',
  password: ''
})

const isLoading = ref(false)
const errorMessage = ref('')
const successMessage = ref('')

const handleLogin = async () => {
  // Clear previous messages
  errorMessage.value = ''
  successMessage.value = ''
  isLoading.value = true

  try {
    const response = await axios.post('/api/login', {
      username: credentials.username,
      password: credentials.password
    })

    if (response.status === 200) {
      successMessage.value = 'Login successful!'
      // Clear form
      credentials.username = ''
      credentials.password = ''
      
      // Redirect to dashboard
      setTimeout(() => {
        router.push('/dashboard')
      }, 1000) // Small delay to show success message
    }
  } catch (error) {
    console.error('Login error:', error)
    
    if (error.response) {
      // Server responded with error status
      if (error.response.status === 401) {
        errorMessage.value = 'Invalid username or password'
      } else if (error.response.status === 400) {
        errorMessage.value = 'Invalid request. Please check your input.'
      } else {
        errorMessage.value = 'Login failed. Please try again.'
      }
    } else if (error.request) {
      // Request was made but no response received
      errorMessage.value = 'Network error. Please check your connection.'
    } else {
      // Something else happened
      errorMessage.value = 'An unexpected error occurred.'
    }
  } finally {
    isLoading.value = false
  }
}
</script>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: calc(100vh - 52px); /* Account for navbar height */
  background-color: #f5f5f5;
  padding: 1rem;
}

.login-form {
  background: white;
  padding: 2rem;
  border-radius: 8px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  width: 100%;
  max-width: 400px;
}

.title {
  text-align: center;
  margin-bottom: 1.5rem !important;
  color: #333;
}

.field:not(:last-child) {
  margin-bottom: 1rem;
}
</style>