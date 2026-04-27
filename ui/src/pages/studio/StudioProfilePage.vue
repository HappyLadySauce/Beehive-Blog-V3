<script setup lang="ts">
import { computed, reactive, shallowRef, watch } from 'vue'

import { useAuthStore } from '@/features/auth/stores/authStore'
import { studioApi } from '@/features/studio'
import BaseButton from '@/shared/components/BaseButton.vue'
import BaseInput from '@/shared/components/BaseInput.vue'
import FormField from '@/shared/components/FormField.vue'
import PageHeader from '@/shared/components/PageHeader.vue'
import StatusAlert from '@/shared/components/StatusAlert.vue'
import UserAvatar from '@/shared/components/UserAvatar.vue'

const authStore = useAuthStore()
const form = reactive({
  nickname: '',
  avatarUrl: '',
})
const isSaving = shallowRef(false)
const successMessage = shallowRef('')
const errorMessage = shallowRef('')

const displayName = computed(() => form.nickname || authStore.currentUser?.username || 'Admin')

watch(
  () => authStore.currentUser,
  (user) => {
    form.nickname = user?.nickname ?? user?.username ?? ''
    form.avatarUrl = user?.avatar_url ?? ''
  },
  { immediate: true },
)

async function handleSubmit() {
  successMessage.value = ''
  errorMessage.value = ''

  if (form.nickname.trim().length === 0) {
    errorMessage.value = 'Enter a display name.'
    return
  }

  isSaving.value = true
  try {
    const response = await studioApi.updateProfile(
      {
        nickname: form.nickname.trim(),
        avatar_url: form.avatarUrl.trim(),
      },
      { accessToken: authStore.accessToken },
    )
    authStore.setCurrentUser(response.user)
    successMessage.value = 'Profile saved.'
  }
  catch (error) {
    errorMessage.value = error instanceof Error ? error.message : 'Unable to save profile.'
  }
  finally {
    isSaving.value = false
  }
}
</script>

<template>
  <section class="profile-page">
    <PageHeader
      eyebrow="Account"
      title="Profile"
      description="Update the nickname and avatar shown across account surfaces."
    />

    <StatusAlert v-if="successMessage" tone="success" title="Profile updated">
      {{ successMessage }}
    </StatusAlert>
    <StatusAlert v-if="errorMessage" tone="danger" title="Profile update failed">
      {{ errorMessage }}
    </StatusAlert>

    <form class="profile-page__form" novalidate @submit.prevent="handleSubmit">
      <div class="profile-page__avatar">
        <UserAvatar :name="displayName" :src="form.avatarUrl" size="lg" />
        <div>
          <strong>{{ displayName }}</strong>
          <span>{{ authStore.currentUser?.role ?? 'admin' }}</span>
        </div>
      </div>

      <FormField label="Display name" for-id="profile-nickname">
        <BaseInput id="profile-nickname" v-model="form.nickname" autocomplete="name" required />
      </FormField>

      <FormField label="Avatar URL" for-id="profile-avatar">
        <BaseInput id="profile-avatar" v-model="form.avatarUrl" type="url" autocomplete="url" />
      </FormField>

      <BaseButton type="submit" :busy="isSaving">Save profile</BaseButton>
    </form>
  </section>
</template>

<style scoped>
.profile-page {
  display: grid;
  gap: 24px;
}

.profile-page__form {
  width: min(560px, 100%);
  display: grid;
  gap: 16px;
}

.profile-page__avatar {
  display: flex;
  align-items: center;
  gap: 14px;
  border: 1px solid var(--bb-color-line);
  border-radius: 8px;
  padding: 16px;
  background: var(--bb-color-surface);
}

.profile-page__avatar div {
  min-width: 0;
  display: grid;
}

.profile-page__avatar strong,
.profile-page__avatar span {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.profile-page__avatar span {
  color: var(--bb-color-muted);
}
</style>
