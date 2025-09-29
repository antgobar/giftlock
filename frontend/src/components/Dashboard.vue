<template>
  <section class="section">
    <div class="container">
      <div class="columns is-centered">
        <div class="column is-8">
          <CreateGroupForm 
            ref="createGroupFormRef"
            :is-creating="isCreating"
            :error="formError"
            :success-message="successMessage"
            @create-group="handleCreateGroup"
            @clear-error="formError = ''"
            @clear-success="successMessage = ''"
          />

          <GroupsList 
            :groups="groups"
            :is-loading="isLoading"
            @view-gifts="handleViewGifts"
            @manage-group="handleManageGroup"
          />
        </div>
      </div>
    </div>
  </section>
</template>

<script setup>
import { onMounted, ref } from 'vue';
import { useRouter } from 'vue-router';
import { useGroups } from '../composables/useGroups.js';
import CreateGroupForm from './CreateGroupForm.vue';
import GroupsList from './GroupsList.vue';

const router = useRouter();

// Use groups composable for data and API calls
const { groups, isLoading, error, fetchGroups, createGroup } = useGroups();

// Form-specific state
const isCreating = ref(false);
const formError = ref('');
const successMessage = ref('');
const createGroupFormRef = ref(null);

// Handle group creation
const handleCreateGroup = async (groupData) => {
  isCreating.value = true;
  formError.value = '';
  successMessage.value = '';

  try {
    const createdGroup = await createGroup(groupData);
    
    // Reset form
    createGroupFormRef.value?.resetForm();
    
    successMessage.value = `Group "${createdGroup.name}" created successfully!`;
  } catch (err) {
    console.error('Error creating group:', err);
    formError.value = err.message || 'Failed to create group. Please try again.';
  } finally {
    isCreating.value = false;
  }
};

// Handle group actions
const handleViewGifts = (group) => {
  // TODO: Navigate to gifts view for this group
  console.log('View gifts for group:', group);
};

const handleManageGroup = (group) => {
  // TODO: Navigate to group management page
  console.log('Manage group:', group);
};

// Load groups when component mounts
onMounted(() => {
  fetchGroups();
});
</script>

