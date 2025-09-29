<template>
  <div v-if="groups.length > 0" class="box">
    <h2 class="title is-4">My Groups</h2>
    <div class="columns is-multiline">
      <div class="column is-half" v-for="group in groups" :key="group.id">
        <GroupCard 
          :group="group"
          :is-deleting="deletingGroupId === group.id"
          @view-gifts="$emit('view-gifts', $event)"
          @manage-group="$emit('manage-group', $event)"
          @delete-group="$emit('delete-group', $event)"
        />
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
</template>

<script setup>
import GroupCard from './GroupCard.vue';

const props = defineProps({
  groups: {
    type: Array,
    required: true
  },
  isLoading: {
    type: Boolean,
    default: false
  },
  deletingGroupId: {
    type: String,
    default: null
  }
});

const emit = defineEmits(['view-gifts', 'manage-group', 'delete-group']);
</script>