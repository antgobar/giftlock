<template>
  <div class="box">
    <h1 class="title is-2 has-text-centered">Welcome to Gift Lock!</h1>

    <div class="field">
      <form @submit.prevent="handleSubmit">
        <div class="field has-addons">
          <div class="control is-expanded">
            <input 
              class="input" 
              type="text" 
              placeholder="Group name"
              v-model="form.name"
              maxlength="255"
              required 
            />
          </div>
          <div class="control">
            <button 
              class="button is-primary" 
              type="submit"
              :class="{ 'is-loading': isCreating }"
              :disabled="isCreating || !form.name.trim()"
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
              v-model="form.description"
              maxlength="1000"
            />
          </div>
          <p class="help" v-if="form.description.length > 0">
            {{ form.description.length }}/1000 characters
          </p>
        </div>
      </form>
    </div>

    <NotificationMessage 
      v-if="error" 
      type="danger" 
      :message="error" 
      @dismiss="$emit('clear-error')"
    />

    <NotificationMessage 
      v-if="successMessage" 
      type="success" 
      :message="successMessage" 
      @dismiss="$emit('clear-success')"
    />
  </div>
</template>

<script setup>
import { ref } from 'vue';
import NotificationMessage from './NotificationMessage.vue';

const props = defineProps({
  isCreating: {
    type: Boolean,
    default: false
  },
  error: {
    type: String,
    default: ''
  },
  successMessage: {
    type: String,
    default: ''
  }
});

const emit = defineEmits(['create-group', 'clear-error', 'clear-success']);

const form = ref({
  name: '',
  description: ''
});

const handleSubmit = () => {
  const trimmedName = form.value.name.trim();
  const trimmedDescription = form.value.description.trim();

  if (!trimmedName) {
    emit('clear-success');
    return;
  }

  emit('create-group', {
    name: trimmedName,
    description: trimmedDescription
  });
};

// Reset form (exposed for parent component)
const resetForm = () => {
  form.value = {
    name: '',
    description: ''
  };
};

defineExpose({
  resetForm
});
</script>