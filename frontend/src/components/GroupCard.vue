<template>
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
          <button 
            class="button is-small is-primary"
            @click="$emit('view-gifts', group)"
          >
            View Gifts
          </button>
          <button 
            class="button is-small is-info"
            @click="$emit('manage-group', group)"
          >
            Manage
          </button>
          <button 
            class="button is-small is-danger"
            @click="$emit('delete-group', group)"
            :disabled="isDeleting"
          >
            <span v-if="isDeleting">Deleting...</span>
            <span v-else>Delete</span>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
const props = defineProps({
  group: {
    type: Object,
    required: true
  },
  isDeleting: {
    type: Boolean,
    default: false
  }
});

const emit = defineEmits(['view-gifts', 'manage-group', 'delete-group']);

// Format date for display
const formatDate = (dateString) => {
  const date = new Date(dateString);
  return date.toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric'
  });
};
</script>