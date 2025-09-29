<template>
  <section class="section">
    <div class="container">
      <div class="columns is-centered">
        <div class="column is-8">
          <div class="box">
            <h1 class="title is-2 has-text-centered">Welcome to Gift Lock!</h1>
            <p class="subtitle is-4 has-text-centered">You have successfully logged in.</p>
            
            <div class="content">
              <p>This is your dashboard where you can:</p>
              <ul>
                <li>Create and manage gift exchange groups</li>
                <li>Add gifts to your wishlist</li>
                <li>View upcoming gift exchange events</li>
                <li>Manage your account settings</li>
              </ul>
            </div>

            <div class="buttons is-centered">
              <button class="button is-primary is-large">
                <span class="icon">
                  <i class="fas fa-users"></i>
                </span>
                <span>Create Group</span>
              </button>
              
              <button class="button is-info is-large">
                <span class="icon">
                  <i class="fas fa-list"></i>
                </span>
                <span>My Gifts</span>
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup>
import { useRouter } from 'vue-router';

const router = useRouter()

// Logout is now handled by the navbar component
</script>

