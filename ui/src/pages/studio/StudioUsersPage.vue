<script setup lang="ts">
import { computed, onMounted, reactive, shallowRef } from 'vue'

import { useAuthStore } from '@/features/auth/stores/authStore'
import { studioApi } from '@/features/studio'
import type { StudioUser } from '@/features/studio'
import BaseButton from '@/shared/components/BaseButton.vue'
import BaseInput from '@/shared/components/BaseInput.vue'
import FormField from '@/shared/components/FormField.vue'
import PageHeader from '@/shared/components/PageHeader.vue'
import PasswordInput from '@/shared/components/PasswordInput.vue'
import StatusAlert from '@/shared/components/StatusAlert.vue'

const authStore = useAuthStore()
const users = shallowRef<StudioUser[]>([])
const total = shallowRef(0)
const isLoading = shallowRef(true)
const isMutating = shallowRef(false)
const errorMessage = shallowRef('')
const successMessage = shallowRef('')
const resetTarget = shallowRef<StudioUser | null>(null)
const resetPassword = shallowRef('')

const filters = reactive({
  keyword: '',
  role: '',
  status: '',
})

const editableUsers = computed(() =>
  users.value.map((user) => ({
    ...user,
    displayName: user.nickname || user.username,
    lastLogin: formatUnixTime(user.last_login_at),
    isSelf: user.user_id === authStore.currentUser?.user_id,
  })),
)

function formatUnixTime(value?: number): string {
  if (!value) {
    return 'Never'
  }
  return new Intl.DateTimeFormat('en', {
    dateStyle: 'medium',
    timeStyle: 'short',
  }).format(new Date(value * 1000))
}

async function loadUsers(): Promise<void> {
  isLoading.value = true
  errorMessage.value = ''
  try {
    const response = await studioApi.listUsers(
      {
        keyword: filters.keyword.trim(),
        role: filters.role,
        status: filters.status,
        page: 1,
        page_size: 50,
      },
      { accessToken: authStore.accessToken },
    )
    users.value = response.items
    total.value = response.total
  }
  catch (error) {
    errorMessage.value = error instanceof Error ? error.message : 'Unable to load users.'
  }
  finally {
    isLoading.value = false
  }
}

async function updateRole(user: StudioUser, role: 'member' | 'admin'): Promise<void> {
  await runUserMutation(async () => {
    const response = await studioApi.updateUserRole(user.user_id, { role }, { accessToken: authStore.accessToken })
    replaceUser(response.user)
    successMessage.value = `${response.user.email} role updated.`
  })
}

async function updateStatus(user: StudioUser, status: 'active' | 'disabled' | 'locked'): Promise<void> {
  await runUserMutation(async () => {
    const response = await studioApi.updateUserStatus(user.user_id, { status }, { accessToken: authStore.accessToken })
    replaceUser(response.user)
    successMessage.value = `${response.user.email} status updated.`
  })
}

async function submitPasswordReset(): Promise<void> {
  if (!resetTarget.value) {
    return
  }
  await runUserMutation(async () => {
    await studioApi.resetUserPassword(
      resetTarget.value!.user_id,
      { new_password: resetPassword.value },
      { accessToken: authStore.accessToken },
    )
    successMessage.value = `${resetTarget.value!.email} password reset.`
    resetTarget.value = null
    resetPassword.value = ''
  })
}

async function runUserMutation(action: () => Promise<void>): Promise<void> {
  errorMessage.value = ''
  successMessage.value = ''
  isMutating.value = true
  try {
    await action()
  }
  catch (error) {
    errorMessage.value = error instanceof Error ? error.message : 'Unable to update user.'
  }
  finally {
    isMutating.value = false
  }
}

function replaceUser(nextUser: StudioUser): void {
  users.value = users.value.map((user) => (user.user_id === nextUser.user_id ? nextUser : user))
}

function openPasswordReset(user: StudioUser): void {
  resetTarget.value = user
  resetPassword.value = ''
}

onMounted(loadUsers)
</script>

<template>
  <section class="users-page">
    <PageHeader
      eyebrow="Studio"
      title="Users"
      description="Search accounts, adjust roles and statuses, and reset credentials for recovery."
    />

    <form class="users-page__filters" @submit.prevent="loadUsers">
      <FormField label="Search" for-id="user-search">
        <BaseInput id="user-search" v-model="filters.keyword" placeholder="Username, email, or nickname" />
      </FormField>
      <label class="users-page__select">
        <span>Role</span>
        <select v-model="filters.role">
          <option value="">All roles</option>
          <option value="member">Member</option>
          <option value="admin">Admin</option>
        </select>
      </label>
      <label class="users-page__select">
        <span>Status</span>
        <select v-model="filters.status">
          <option value="">All statuses</option>
          <option value="pending">Pending</option>
          <option value="active">Active</option>
          <option value="disabled">Disabled</option>
          <option value="locked">Locked</option>
        </select>
      </label>
      <BaseButton type="submit" :busy="isLoading">Apply</BaseButton>
    </form>

    <StatusAlert v-if="successMessage" tone="success" title="User updated">
      {{ successMessage }}
    </StatusAlert>
    <StatusAlert v-if="errorMessage" tone="danger" title="Users unavailable">
      {{ errorMessage }}
    </StatusAlert>
    <StatusAlert v-else-if="isLoading" tone="info" title="Loading users">
      User records are being loaded from gateway.
    </StatusAlert>

    <div class="users-page__table" role="region" aria-label="Studio users" tabindex="0">
      <table>
        <thead>
          <tr>
            <th scope="col">User</th>
            <th scope="col">Role</th>
            <th scope="col">Status</th>
            <th scope="col">Last login</th>
            <th scope="col">Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-if="editableUsers.length === 0">
            <td colspan="5">No users found.</td>
          </tr>
          <tr v-for="user in editableUsers" v-else :key="user.user_id">
            <td>
              <strong>{{ user.displayName }}</strong>
              <span>{{ user.email }}</span>
            </td>
            <td>
              <select
                :value="user.role"
                :disabled="user.isSelf || isMutating"
                :aria-label="`Change ${user.email} role`"
                @change="updateRole(user, (($event.target as HTMLSelectElement).value as 'member' | 'admin'))"
              >
                <option value="member">Member</option>
                <option value="admin">Admin</option>
              </select>
            </td>
            <td>
              <select
                :value="user.status"
                :disabled="user.isSelf || isMutating"
                :aria-label="`Change ${user.email} status`"
                @change="updateStatus(user, (($event.target as HTMLSelectElement).value as 'active' | 'disabled' | 'locked'))"
              >
                <option value="active">Active</option>
                <option value="disabled">Disabled</option>
                <option value="locked">Locked</option>
              </select>
            </td>
            <td>{{ user.lastLogin }}</td>
            <td>
              <BaseButton variant="secondary" :disabled="user.isSelf || isMutating" @click="openPasswordReset(user)">
                Reset password
              </BaseButton>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <p class="users-page__count">{{ total }} total users</p>

    <form v-if="resetTarget" class="users-page__reset" @submit.prevent="submitPasswordReset">
      <div>
        <strong>Set new password</strong>
        <span>{{ resetTarget.email }}</span>
      </div>
      <FormField label="New password" for-id="reset-password">
        <PasswordInput id="reset-password" v-model="resetPassword" autocomplete="new-password" />
      </FormField>
      <div class="users-page__reset-actions">
        <BaseButton type="submit" :busy="isMutating">Save password</BaseButton>
        <BaseButton variant="ghost" @click="resetTarget = null">Cancel</BaseButton>
      </div>
    </form>
  </section>
</template>

<style scoped>
.users-page {
  display: grid;
  gap: 24px;
}

.users-page__filters {
  display: grid;
  grid-template-columns: minmax(180px, 1fr) repeat(2, minmax(140px, 180px)) auto;
  align-items: end;
  gap: 12px;
}

.users-page__select {
  display: grid;
  gap: 6px;
  color: var(--bb-color-muted);
  font-size: 0.92rem;
  font-weight: 650;
}

select {
  min-height: 44px;
  border: 1px solid var(--bb-color-line);
  border-radius: 8px;
  padding: 0 10px;
  color: var(--bb-color-text);
  background: var(--bb-color-surface);
}

select:focus-visible,
.users-page__table:focus-visible {
  outline: none;
  box-shadow: 0 0 0 3px var(--bb-color-focus);
}

.users-page__table {
  overflow-x: auto;
  border: 1px solid var(--bb-color-line);
  border-radius: 8px;
  background: var(--bb-color-surface);
}

table {
  width: 100%;
  min-width: 760px;
  border-collapse: collapse;
}

th,
td {
  border-bottom: 1px solid var(--bb-color-line);
  padding: 12px;
  text-align: left;
  vertical-align: middle;
}

th {
  color: var(--bb-color-muted);
  font-size: 0.8rem;
  text-transform: uppercase;
}

tr:last-child td {
  border-bottom: 0;
}

td:first-child {
  display: grid;
  gap: 3px;
}

td:first-child span,
.users-page__reset span,
.users-page__count {
  color: var(--bb-color-muted);
}

.users-page__reset {
  width: min(520px, 100%);
  display: grid;
  gap: 14px;
  border: 1px solid var(--bb-color-line);
  border-radius: 8px;
  padding: 16px;
  background: var(--bb-color-surface);
}

.users-page__reset div:first-child {
  display: grid;
  gap: 2px;
}

.users-page__reset-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

@media (max-width: 760px) {
  .users-page__filters {
    grid-template-columns: 1fr;
  }
}
</style>
