<template>
  <div class="login-container">
    <div class="login-form">
      <h2>Login</h2>
      
      <form @submit.prevent="handleLogin">
        <div class="form-group">
          <label for="username">Username:</label>
          <input
            type="text"
            id="username"
            v-model="credentials.username"
            required
            :disabled="isLoading"
          />
        </div>
        
        <div class="form-group">
          <label for="password">Password:</label>
          <input
            type="password"
            id="password"
            v-model="credentials.password"
            required
            :disabled="isLoading"
          />
        </div>
        
        <div class="form-group">
          <button type="submit" :disabled="isLoading">
            {{ isLoading ? 'Logging in...' : 'Login' }}
          </button>
        </div>
      </form>
      
      <div v-if="errorMessage" class="error-message">
        {{ errorMessage }}
      </div>
      
      <div v-if="successMessage" class="success-message">
        {{ successMessage }}
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
  min-height: 100vh;
  background-color: #f5f5f5;
}

.login-form {
  background: white;
  padding: 2rem;
  border-radius: 8px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  width: 100%;
  max-width: 400px;
}

.login-form h2 {
  text-align: center;
  margin-bottom: 1.5rem;
  color: #333;
}

.form-group {
  margin-bottom: 1rem;
}

.form-group label {
  display: block;
  margin-bottom: 0.5rem;
  font-weight: bold;
  color: #555;
}

.form-group input {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 1rem;
  box-sizing: border-box;
}

.form-group input:focus {
  outline: none;
  border-color: #007bff;
  box-shadow: 0 0 0 2px rgba(0, 123, 255, 0.25);
}

.form-group input:disabled {
  background-color: #f8f9fa;
  cursor: not-allowed;
}

.form-group button {
  width: 100%;
  padding: 0.75rem;
  background-color: #007bff;
  color: white;
  border: none;
  border-radius: 4px;
  font-size: 1rem;
  cursor: pointer;
  transition: background-color 0.2s;
}

.form-group button:hover:not(:disabled) {
  background-color: #0056b3;
}

.form-group button:disabled {
  background-color: #6c757d;
  cursor: not-allowed;
}

.error-message {
  margin-top: 1rem;
  padding: 0.75rem;
  background-color: #f8d7da;
  color: #721c24;
  border: 1px solid #f5c6cb;
  border-radius: 4px;
  text-align: center;
}

.success-message {
  margin-top: 1rem;
  padding: 0.75rem;
  background-color: #d4edda;
  color: #155724;
  border: 1px solid #c3e6cb;
  border-radius: 4px;
  text-align: center;
}
</style>