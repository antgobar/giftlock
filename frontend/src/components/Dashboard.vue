<template>
  <section class="section">
    <div class="container">
      <div class="columns is-centered">
        <div class="column is-8">
          <div class="box">
            <h1 class="title is-2 has-text-centered">Welcome to Gift Lock!</h1>

            <!-- Create Group Form -->
            <div class="field">
              <form @submit.prevent="createGroup">
                <div class="field has-addons">
                  <div class="control is-expanded">
                    <input 
                      class="input" 
                      type="text" 
                      placeholder="Group name"
                      v-model="newGroup.name"
                      maxlength="255"
                      required 
                    />
                  </div>
                  <div class="control">
                    <button 
                      class="button is-primary" 
                      type="submit"
                      :class="{ 'is-loading': isCreating }"
                      :disabled="isCreating || !newGroup.name.trim()"
                    >
                      Create Group
                    </button>
                  </div>
                </div>
                <div class="field">
                  <div class="control">
                    <input 
                      class="input" 
                      type="text" 
                      placeholder="Description (optional)"
                      v-model="newGroup.description"
                      maxlength="1000"
                    />
                  </div>
                  <p class="help" v-if="newGroup.description.length > 0">
                    {{ newGroup.description.length }}/1000 characters
                  </p>
                </div>
              </form>
            </div>

            <!-- Error Message -->
            <div v-if="error" class="notification is-danger">
              <button class="delete" @click="error = ''"></button>
              {{ error }}
            </div>

            <!-- Success Message -->
            <div v-if="successMessage" class="notification is-success">
              <button class="delete" @click="successMessage = ''"></button>
              {{ successMessage }}
            </div>

            <div class="buttons is-centered">
              <button class="button is-info">
                My Gifts
              </button>
            </div>
          </div>

          <!-- Groups List -->
          <div class="box" v-if="groups.length > 0">
            <h2 class="title is-4">My Groups</h2>
            <div class="columns is-multiline">
              <div class="column is-half" v-for="group in groups" :key="group.id">
                <div class="card">
                  <div class="card-content">
                    <div class="media">
                      <div class="media-content">
                        <p class="title is-5">{{ group.name }}</p>
                        <p class="subtitle is-6" v-if="group.description">{{ group.description }}</p>
                        <p class="subtitle is-6">Created: {{ formatDate(group.createdAt) }}</p>
                      </div>
                    </div>
                    <div class="content">
                      <div class="buttons">
                        <button class="button is-small is-primary">
                          View Gifts
                        </button>
                        <button class="button is-small is-info">
                          Manage
                        </button>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- Empty State -->
          <div class="box has-text-centered" v-else-if="!isLoading">
            <p class="title is-5">No groups yet</p>
            <p class="subtitle">Create your first gift group above!</p>
          </div>

          <!-- Loading State -->
          <div class="box has-text-centered" v-if="isLoading">
            <p class="title is-5">Loading your groups...</p>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup>
import { onMounted, ref } from 'vue';
import { useRouter } from 'vue-router';

const router = useRouter();

// Reactive data
const groups = ref([]);
const newGroup = ref({
  name: '',
  description: ''
});
const isLoading = ref(false);
const isCreating = ref(false);
const error = ref('');
const successMessage = ref('');

// Fetch groups from API
const fetchGroups = async () => {
  isLoading.value = true;
  error.value = '';
  
  try {
    const response = await fetch('/api/groups', {
      method: 'GET',
      credentials: 'include'
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    const data = await response.json();
    groups.value = data || [];
  } catch (err) {
    console.error('Error fetching groups:', err);
    error.value = 'Failed to load groups. Please try again.';
  } finally {
    isLoading.value = false;
  }
};

// Create a new group
const createGroup = async () => {
  const trimmedName = newGroup.value.name.trim();
  const trimmedDescription = newGroup.value.description.trim();

  if (!trimmedName) {
    error.value = 'Group name is required';
    return;
  }

  if (trimmedName.length > 255) {
    error.value = 'Group name must be 255 characters or less';
    return;
  }

  if (trimmedDescription.length > 1000) {
    error.value = 'Group description must be 1000 characters or less';
    return;
  }

  isCreating.value = true;
  error.value = '';
  successMessage.value = '';

  try {
    const requestBody = {
      name: trimmedName,
      description: trimmedDescription
    };

    const response = await fetch('/api/groups', {
      method: 'POST',
      credentials: 'include',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(requestBody)
    });

    if (!response.ok) {
      const errorText = await response.text();
      throw new Error(errorText || `HTTP error! status: ${response.status}`);
    }

    const createdGroup = await response.json();
    groups.value.unshift(createdGroup); // Add to beginning of list
    
    // Reset form
    newGroup.value = {
      name: '',
      description: ''
    };
    
    successMessage.value = `Group "${createdGroup.name}" created successfully!`;
  } catch (err) {
    console.error('Error creating group:', err);
    error.value = err.message || 'Failed to create group. Please try again.';
  } finally {
    isCreating.value = false;
  }
};

// Format date for display
const formatDate = (dateString) => {
  const date = new Date(dateString);
  return date.toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric'
  });
};

// Load groups when component mounts
onMounted(() => {
  fetchGroups();
});
</script>

