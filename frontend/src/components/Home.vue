<template>
  <section class="section">
    <div class="container">
      <div class="has-text-centered">
        <h1 class="title">Welcome to Gift Lock</h1>
        <p class="subtitle">Secure gift exchanges made simple</p>
        
        <div class="buttons is-centered mt-5" v-if="!isAuthenticated">
          <router-link to="/register" class="button is-primary">
            Sign Up
          </router-link>
          <router-link to="/login" class="button">
            Log In
          </router-link>
        </div>
        
        <div class="buttons is-centered mt-5" v-else>
          <router-link to="/dashboard" class="button is-primary">
            Go to Dashboard
          </router-link>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup>
import axios from 'axios'
import { onMounted, ref } from 'vue'

const isAuthenticated = ref(false)

const checkAuthStatus = async () => {
  try {
    const response = await axios.get('/api/me')
    isAuthenticated.value = response.status === 200
  } catch (error) {
    isAuthenticated.value = false
  }
}

onMounted(() => {
  checkAuthStatus()
})
</script>