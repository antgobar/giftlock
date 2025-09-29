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
            :deleting-group-id="isDeletingGroup ? groupToDelete?.id : null"
            @view-gifts="handleViewGifts"
            @manage-group="handleManageGroup"
            @delete-group="handleDeleteGroup"
          />

          <!-- Delete Confirmation Modal -->
          <ConfirmModal
            :is-active="showDeleteModal"
            :title="'Delete Group'"
            :message="`Are you sure you want to delete the group '${groupToDelete?.name}'? This action cannot be undone.`"
            :confirm-text="'Delete'"
            :loading-text="'Deleting...'"
            :is-loading="isDeletingGroup"
            @confirm="confirmDeleteGroup"
            @cancel="cancelDeleteGroup"
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
import ConfirmModal from './ConfirmModal.vue';
import CreateGroupForm from './CreateGroupForm.vue';
import GroupsList from './GroupsList.vue';

const router = useRouter();

// Use groups composable for data and API calls
const { groups, isLoading, error, fetchGroups, createGroup, deleteGroup } = useGroups();

// Form-specific state
const isCreating = ref(false);
const formError = ref('');
const successMessage = ref('');
const createGroupFormRef = ref(null);

// Delete modal state
const showDeleteModal = ref(false);
const groupToDelete = ref(null);
const isDeletingGroup = ref(false);

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

// Handle group deletion with confirmation
const handleDeleteGroup = async (group) => {
  groupToDelete.value = group;
  showDeleteModal.value = true;
};

const confirmDeleteGroup = async () => {
  if (!groupToDelete.value) return;
  
  isDeletingGroup.value = true;
  
  try {
    await deleteGroup(groupToDelete.value.id);
    successMessage.value = `Group "${groupToDelete.value.name}" deleted successfully!`;
    showDeleteModal.value = false;
    groupToDelete.value = null;
  } catch (err) {
    console.error('Error deleting group:', err);
    formError.value = err.message || 'Failed to delete group. Please try again.';
  } finally {
    isDeletingGroup.value = false;
  }
};

const cancelDeleteGroup = () => {
  showDeleteModal.value = false;
  groupToDelete.value = null;
};

// Load groups when component mounts
onMounted(() => {
  fetchGroups();
});
</script>

