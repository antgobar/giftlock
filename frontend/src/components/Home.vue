<template>
  <div class="home-container">
    <!-- Hero Section -->
    <section class="hero is-primary is-large">
      <div class="hero-body">
        <div class="container has-text-centered">
          <h1 class="title is-1">
            <span class="icon is-large">
              <i class="fas fa-gift"></i>
            </span>
            Gift Lock
          </h1>
          <h2 class="subtitle is-3">
            Secure your gift exchanges with friends and family
          </h2>
          <p class="is-size-5 mb-5">
            Organize secret Santa, birthday surprises, and holiday gift exchanges 
            without spoiling the surprise. Keep gift lists private until the big reveal!
          </p>
          <div class="buttons is-centered" v-if="!isAuthenticated">
            <router-link to="/register" class="button is-white is-large">
              <strong>Get Started</strong>
            </router-link>
            <router-link to="/login" class="button is-primary is-inverted is-outlined is-large">
              Sign In
            </router-link>
          </div>
          <div class="buttons is-centered" v-else>
            <router-link to="/dashboard" class="button is-white is-large">
              <strong>Go to Dashboard</strong>
            </router-link>
          </div>
        </div>
      </div>
    </section>

    <!-- Features Section -->
    <section class="section">
      <div class="container">
        <div class="columns is-multiline">
          <div class="column is-4">
            <div class="card">
              <div class="card-content has-text-centered">
                <div class="icon is-large has-text-primary mb-4">
                  <i class="fas fa-users fa-3x"></i>
                </div>
                <h3 class="title is-4">Group Management</h3>
                <p>Create and manage gift exchange groups with friends, family, or colleagues. Invite members and organize your events seamlessly.</p>
              </div>
            </div>
          </div>
          
          <div class="column is-4">
            <div class="card">
              <div class="card-content has-text-centered">
                <div class="icon is-large has-text-primary mb-4">
                  <i class="fas fa-lock fa-3x"></i>
                </div>
                <h3 class="title is-4">Secure & Private</h3>
                <p>Your gift lists remain completely private until the reveal date. No peeking, no spoilers - just pure surprise and joy!</p>
              </div>
            </div>
          </div>
          
          <div class="column is-4">
            <div class="card">
              <div class="card-content has-text-centered">
                <div class="icon is-large has-text-primary mb-4">
                  <i class="fas fa-calendar-alt fa-3x"></i>
                </div>
                <h3 class="title is-4">Scheduled Reveals</h3>
                <p>Set reveal dates for your gift exchanges. Perfect for Christmas morning, birthdays, or any special occasion you're celebrating.</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- How It Works Section -->
    <section class="section has-background-light">
      <div class="container">
        <h2 class="title is-2 has-text-centered mb-5">How It Works</h2>
        <div class="columns is-variable is-8">
          <div class="column is-3 has-text-centered">
            <div class="icon is-large has-text-primary mb-3">
              <i class="fas fa-user-plus fa-2x"></i>
            </div>
            <h4 class="title is-5">1. Sign Up</h4>
            <p>Create your free Gift Lock account in seconds</p>
          </div>
          
          <div class="column is-3 has-text-centered">
            <div class="icon is-large has-text-primary mb-3">
              <i class="fas fa-users-cog fa-2x"></i>
            </div>
            <h4 class="title is-5">2. Create Group</h4>
            <p>Set up your gift exchange group and invite participants</p>
          </div>
          
          <div class="column is-3 has-text-centered">
            <div class="icon is-large has-text-primary mb-3">
              <i class="fas fa-list fa-2x"></i>
            </div>
            <h4 class="title is-5">3. Add Gifts</h4>
            <p>Add your gift ideas securely to your private list</p>
          </div>
          
          <div class="column is-3 has-text-centered">
            <div class="icon is-large has-text-primary mb-3">
              <i class="fas fa-surprise fa-2x"></i>
            </div>
            <h4 class="title is-5">4. Reveal Day</h4>
            <p>Enjoy the surprise when gifts are revealed on the special day!</p>
          </div>
        </div>
      </div>
    </section>

    <!-- CTA Section -->
    <section class="section is-medium" v-if="!isAuthenticated">
      <div class="container has-text-centered">
        <h2 class="title is-3 mb-4">Ready to Lock in the Perfect Gift Exchange?</h2>
        <p class="subtitle is-5 mb-5">Join thousands of families and friends making gift giving more exciting!</p>
        <router-link to="/register" class="button is-primary is-large">
          <strong>Start Your Free Account</strong>
        </router-link>
      </div>
    </section>
  </div>
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

<style scoped>
.home-container {
  min-height: calc(100vh - 52px); /* Account for navbar height */
}

.hero {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.card {
  height: 100%;
  transition: transform 0.2s ease-in-out;
}

.card:hover {
  transform: translateY(-2px);
}

.icon.is-large i {
  font-size: 3rem !important;
}

@media screen and (max-width: 768px) {
  .hero .title {
    font-size: 2.5rem !important;
  }
  
  .hero .subtitle {
    font-size: 1.25rem !important;
  }
  
  .hero .is-size-5 {
    font-size: 1rem !important;
  }
}
</style>