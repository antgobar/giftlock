<template>
  <div class="modal" :class="{ 'is-active': isActive }">
    <div class="modal-background" @click="$emit('cancel')"></div>
    <div class="modal-card">
      <header class="modal-card-head">
        <p class="modal-card-title">{{ title }}</p>
        <button class="delete" aria-label="close" @click="$emit('cancel')"></button>
      </header>
      <section class="modal-card-body">
        <p>{{ message }}</p>
      </section>
      <footer class="modal-card-foot">
        <button 
          class="button is-danger" 
          @click="$emit('confirm')"
          :disabled="isLoading"
        >
          <span v-if="isLoading">{{ loadingText }}</span>
          <span v-else>{{ confirmText }}</span>
        </button>
        <button 
          class="button" 
          @click="$emit('cancel')"
          :disabled="isLoading"
        >
          Cancel
        </button>
      </footer>
    </div>
  </div>
</template>

<script setup>
const props = defineProps({
  isActive: {
    type: Boolean,
    default: false
  },
  title: {
    type: String,
    default: 'Confirm Action'
  },
  message: {
    type: String,
    required: true
  },
  confirmText: {
    type: String,
    default: 'Confirm'
  },
  loadingText: {
    type: String,
    default: 'Loading...'
  },
  isLoading: {
    type: Boolean,
    default: false
  }
});

const emit = defineEmits(['confirm', 'cancel']);
</script>