<template>
  <div :class="notificationClass">
    <button class="delete" @click="$emit('dismiss')"></button>
    {{ message }}
  </div>
</template>

<script setup>
import { computed } from 'vue';

const props = defineProps({
  type: {
    type: String,
    required: true,
    validator: (value) => ['success', 'danger', 'warning', 'info'].includes(value)
  },
  message: {
    type: String,
    required: true
  }
});

const emit = defineEmits(['dismiss']);

const notificationClass = computed(() => {
  return `notification is-${props.type}`;
});
</script>