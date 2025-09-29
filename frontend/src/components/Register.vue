<template>
  <section class="section">
    <div class="container">
      <div class="columns is-centered">
        <div class="column is-4">
          <div class="box">
            <h2 class="title is-3 has-text-centered">Create Account</h2>
            
            <form @submit.prevent="handleRegister">
              <div class="field">
                <label class="label" for="username">Username</label>
                <div class="control">
                  <input
                    class="input"
                    :class="{ 'is-danger': errors.username }"
                    type="text"
                    id="username"
                    v-model="credentials.username"
                    :disabled="isLoading"
                    placeholder="Enter your username"
                  />
                </div>
                <p class="help is-danger" v-if="errors.username">{{ errors.username }}</p>
              </div>
              
              <div class="field">
                <label class="label" for="password">Password</label>
                <div class="control">
                  <input
                    class="input"
                    :class="{ 'is-danger': errors.password }"
                    type="password"
                    id="password"
                    v-model="credentials.password"
                    :disabled="isLoading"
                    placeholder="Enter your password"
                  />
                </div>
                <p class="help is-danger" v-if="errors.password">{{ errors.password }}</p>
              </div>

              <div class="field">
                <label class="label" for="confirmPassword">Confirm Password</label>
                <div class="control">
                  <input
                    class="input"
                    :class="{ 'is-danger': errors.confirmPassword }"
                    type="password"
                    id="confirmPassword"
                    v-model="credentials.confirmPassword"
                    :disabled="isLoading"
                    placeholder="Confirm your password"
                  />
                </div>
                <p class="help is-danger" v-if="errors.confirmPassword">{{ errors.confirmPassword }}</p>
              </div>
              
              <div class="field">
                <div class="control">
                  <button class="button is-primary is-fullwidth" type="submit" :class="{ 'is-loading': isLoading }" :disabled="isLoading">
                    Create Account
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
              <p>Already have an account? 
                <router-link to="/login" class="has-text-primary">Sign in here</router-link>
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup>
import axios from 'axios'
import { computed, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

const credentials = reactive({
  username: '',
  password: '',
  confirmPassword: ''
})

const errors = reactive({
  username: '',
  password: '',
  confirmPassword: ''
})

const isLoading = ref(false)
const errorMessage = ref('')
const successMessage = ref('')

const isFormValid = computed(() => {
  return credentials.username.trim().length >= 3 &&
         credentials.password.trim().length >= 6 &&
         credentials.confirmPassword.trim() !== '' &&
         credentials.password === credentials.confirmPassword
})

const validateForm = () => {
  // Clear previous errors
  errors.username = ''
  errors.password = ''
  errors.confirmPassword = ''

  let isValid = true

  // Username validation
  if (!credentials.username.trim()) {
    errors.username = 'Username is required'
    isValid = false
  } else if (credentials.username.length < 3) {
    errors.username = 'Username must be at least 3 characters long'
    isValid = false
  }

  // Password validation
  if (!credentials.password.trim()) {
    errors.password = 'Password is required'
    isValid = false
  } else if (credentials.password.length < 6) {
    errors.password = 'Password must be at least 6 characters long'
    isValid = false
  }

  // Confirm password validation
  if (!credentials.confirmPassword.trim()) {
    errors.confirmPassword = 'Please confirm your password'
    isValid = false
  } else if (credentials.password !== credentials.confirmPassword) {
    errors.confirmPassword = 'Passwords do not match'
    isValid = false
  }

  return isValid
}

const handleRegister = async () => {
  // Clear previous messages
  errorMessage.value = ''
  successMessage.value = ''

  if (!validateForm()) {
    return
  }

  isLoading.value = true

  try {
    const response = await axios.post('/api/register', {
      username: credentials.username,
      password: credentials.password
    })

    if (response.status === 201) {
      successMessage.value = 'Account created successfully! Redirecting to login...'
      
      // Clear form
      credentials.username = ''
      credentials.password = ''
      credentials.confirmPassword = ''
      
      // Redirect to login after success
      setTimeout(() => {
        router.push('/login')
      }, 2000)
    }
  } catch (error) {
    console.error('Registration error:', error)
    
    if (error.response) {
      // Server responded with error status
      if (error.response.status === 409) {
        errorMessage.value = 'Username is already taken. Please choose a different username.'
        errors.username = 'This username is already taken'
      } else if (error.response.status === 400) {
        errorMessage.value = 'Invalid input. Please check your information and try again.'
      } else {
        errorMessage.value = 'Registration failed. Please try again.'
      }
    } else if (error.request) {
      // Request was made but no response received
      errorMessage.value = 'Network error. Please check your connection and try again.'
    } else {
      // Something else happened
      errorMessage.value = 'An unexpected error occurred. Please try again.'
    }
  } finally {
    isLoading.value = false
  }
}
</script>

