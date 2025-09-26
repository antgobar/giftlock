<template>
  <div class="dashboard-container">
    <div class="dashboard">
      <h1>Welcome!</h1>
      <p>You have successfully logged in.</p>
      <button @click="handleLogout" class="logout-button">
        Logout
      </button>
    </div>
  </div>
</template>

<script setup>
import axios from 'axios'
import { useRouter } from 'vue-router'

const router = useRouter()

const handleLogout = async () => {
  try {
    await axios.post('/api/logout')
    router.push('/login')
  } catch (error) {
    console.error('Logout error:', error)
    // Even if logout fails, redirect to login
    router.push('/login')
  }
}
</script>

<style scoped>
.dashboard-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background-color: #f5f5f5;
}

.dashboard {
  background: white;
  padding: 2rem;
  border-radius: 8px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  text-align: center;
  max-width: 500px;
  width: 100%;
}

.dashboard h1 {
  color: #333;
  margin-bottom: 1rem;
}

.dashboard p {
  color: #666;
  margin-bottom: 2rem;
}

.logout-button {
  padding: 0.75rem 1.5rem;
  background-color: #dc3545;
  color: white;
  border: none;
  border-radius: 4px;
  font-size: 1rem;
  cursor: pointer;
  transition: background-color 0.2s;
}

.logout-button:hover {
  background-color: #c82333;
}
</style>