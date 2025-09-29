import { ref } from 'vue';

export function useGroups() {
    const groups = ref([]);
    const isLoading = ref(false);
    const error = ref('');

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
    const createGroup = async (groupData) => {
        const { name, description } = groupData;

        if (!name) {
            throw new Error('Group name is required');
        }

        if (name.length > 255) {
            throw new Error('Group name must be 255 characters or less');
        }

        if (description.length > 1000) {
            throw new Error('Group description must be 1000 characters or less');
        }

        try {
            const requestBody = {
                name,
                description
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

            return createdGroup;
        } catch (err) {
            console.error('Error creating group:', err);
            throw err;
        }
    };

    return {
        groups,
        isLoading,
        error,
        fetchGroups,
        createGroup
    };
}