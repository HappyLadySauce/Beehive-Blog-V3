<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, reactive, shallowRef, watch } from 'vue'

import { useAuthStore } from '@/features/auth/stores/authStore'
import { studioApi } from '@/features/studio'
import type { StudioUser } from '@/features/studio'
import ActionTagButton from '@/shared/components/ActionTagButton.vue'
import BaseButton from '@/shared/components/BaseButton.vue'
import BaseInput from '@/shared/components/BaseInput.vue'
import FormField from '@/shared/components/FormField.vue'
import PageHeader from '@/shared/components/PageHeader.vue'
import PageLoadingState from '@/shared/components/PageLoadingState.vue'
import PasswordInput from '@/shared/components/PasswordInput.vue'
import ReadonlyField from '@/shared/components/ReadonlyField.vue'
import SideDrawer from '@/shared/components/SideDrawer.vue'
import StatusAlert from '@/shared/components/StatusAlert.vue'
import StatusBadge from '@/shared/components/StatusBadge.vue'
import { useConfirm, useToast } from '@/shared/composables'

type UserMode = 'view' | 'edit'

const authStore = useAuthStore()
const { confirm } = useConfirm()
const { pushToast } = useToast()

const users = shallowRef<StudioUser[]>([])
const total = shallowRef(0)
const isLoading = shallowRef(true)
const isMutating = shallowRef(false)
const errorMessage = shallowRef('')
const selectedUser = shallowRef<StudioUser | null>(null)
const userMode = shallowRef<UserMode>('view')
const resetPassword = shallowRef('')
let filterTimer: number | undefined

const filters = reactive({
  keyword: '',
  role: '',
  status: '',
})

const editForm = reactive({
  role: 'member' as 'member' | 'admin',
  status: 'active' as 'active' | 'disabled' | 'locked',
})

const displayUsers = computed(() =>
  users.value.map((user) => ({
    ...user,
    displayName: user.nickname || user.username,
    lastLogin: formatUnixTime(user.last_login_at),
    isSelf: user.user_id === authStore.currentUser?.user_id,
  })),
)

const drawerTitle = computed(() => (userMode.value === 'edit' ? 'Edit user' : 'User details'))

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
  } catch (error) {
    errorMessage.value = error instanceof Error ? error.message : 'Unable to load users.'
    pushToast({ tone: 'danger', title: 'Users unavailable', message: errorMessage.value })
  } finally {
    isLoading.value = false
  }
}

function scheduleLoadUsers(): void {
  window.clearTimeout(filterTimer)
  filterTimer = window.setTimeout(() => {
    void loadUsers()
  }, 300)
}

function openUser(user: StudioUser, mode: UserMode): void {
  selectedUser.value = user
  userMode.value = mode
  editForm.role = normalizeRole(user.role)
  editForm.status = normalizeStatus(user.status)
  resetPassword.value = ''
}

function closeDrawer(): void {
  selectedUser.value = null
  resetPassword.value = ''
}

async function saveUserEdits(): Promise<void> {
  if (!selectedUser.value) {
    return
  }
  const target = selectedUser.value
  await runUserMutation(async () => {
    let nextUser = target
    if (editForm.role !== target.role) {
      nextUser = (await studioApi.updateUserRole(target.user_id, { role: editForm.role }, { accessToken: authStore.accessToken })).user
    }
    if (editForm.status !== nextUser.status) {
      nextUser = (await studioApi.updateUserStatus(target.user_id, { status: editForm.status }, { accessToken: authStore.accessToken })).user
    }
    replaceUser(nextUser)
    selectedUser.value = nextUser
    pushToast({ tone: 'success', title: 'User updated', message: `${nextUser.email} has been updated.` })
  })
}

async function submitPasswordReset(): Promise<void> {
  if (!selectedUser.value || resetPassword.value.trim() === '') {
    return
  }
  const approved = await confirm({
    title: 'Reset user password?',
    message: `This will replace the current password for ${selectedUser.value.email} and revoke existing sessions.`,
    confirmText: 'Reset password',
    tone: 'danger',
  })
  if (!approved) {
    return
  }
  await runUserMutation(async () => {
    await studioApi.resetUserPassword(
      selectedUser.value!.user_id,
      { new_password: resetPassword.value },
      { accessToken: authStore.accessToken },
    )
    pushToast({ tone: 'success', title: 'Password reset', message: `${selectedUser.value!.email} can now use the new password.` })
    resetPassword.value = ''
  })
}

async function deleteUser(user: StudioUser): Promise<void> {
  const approved = await confirm({
    title: 'Delete user?',
    message: `${user.email} will be soft deleted, hidden from the default list, and all active sessions will be revoked.`,
    confirmText: 'Delete user',
    tone: 'danger',
  })
  if (!approved) {
    return
  }
  await runUserMutation(async () => {
    await studioApi.deleteUser(user.user_id, { accessToken: authStore.accessToken })
    users.value = users.value.filter((item) => item.user_id !== user.user_id)
    total.value = Math.max(0, total.value - 1)
    if (selectedUser.value?.user_id === user.user_id) {
      closeDrawer()
    }
    pushToast({ tone: 'success', title: 'User deleted', message: `${user.email} was removed from the active list.` })
  })
}

async function runUserMutation(action: () => Promise<void>): Promise<void> {
  errorMessage.value = ''
  isMutating.value = true
  try {
    await action()
  } catch (error) {
    errorMessage.value = error instanceof Error ? error.message : 'Unable to update user.'
    pushToast({ tone: 'danger', title: 'User update failed', message: errorMessage.value })
  } finally {
    isMutating.value = false
  }
}

function replaceUser(nextUser: StudioUser): void {
  users.value = users.value.map((user) => (user.user_id === nextUser.user_id ? nextUser : user))
}

function normalizeRole(role: string): 'member' | 'admin' {
  return role === 'admin' ? 'admin' : 'member'
}

function normalizeStatus(status: string): 'active' | 'disabled' | 'locked' {
  if (status === 'disabled' || status === 'locked') {
    return status
  }
  return 'active'
}

watch(() => [filters.keyword, filters.role, filters.status], scheduleLoadUsers)

onMounted(loadUsers)
onBeforeUnmount(() => window.clearTimeout(filterTimer))
</script>

<template>
  <section class="users-page">
    <PageHeader
      eyebrow="Studio"
      title="Users"
      description="Search accounts and manage role, status, password recovery, and soft deletion from row actions."
    />

    <div class="users-page__filters">
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
    </div>

    <StatusAlert v-if="errorMessage" tone="danger" title="Users unavailable">{{ errorMessage }}</StatusAlert>
    <PageLoadingState v-else-if="isLoading" title="Loading users" :rows="5" />

    <div v-else class="users-page__table" role="region" aria-label="Studio users" tabindex="0">
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
          <tr v-if="displayUsers.length === 0">
            <td colspan="5">No users found.</td>
          </tr>
          <tr v-for="user in displayUsers" v-else :key="user.user_id">
            <td>
              <strong>{{ user.displayName }}</strong>
              <span>{{ user.email }}</span>
            </td>
            <td><StatusBadge :value="user.role" /></td>
            <td><StatusBadge :value="user.status" /></td>
            <td>{{ user.lastLogin }}</td>
            <td>
              <div class="users-page__actions">
                <ActionTagButton @click="openUser(user, 'view')">View</ActionTagButton>
                <ActionTagButton tone="primary" :disabled="user.isSelf || isMutating" @click="openUser(user, 'edit')">Edit</ActionTagButton>
                <ActionTagButton tone="danger" :disabled="user.isSelf || isMutating" @click="deleteUser(user)">Delete</ActionTagButton>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <p v-if="!isLoading" class="users-page__count">{{ total }} total users</p>

    <SideDrawer :open="selectedUser !== null" :title="drawerTitle" :description="selectedUser?.email" @close="closeDrawer">
      <div v-if="selectedUser" class="users-page__drawer">
        <div class="users-page__detail-grid">
          <ReadonlyField label="Username" :value="selectedUser.username" />
          <ReadonlyField label="Nickname" :value="selectedUser.nickname" />
          <ReadonlyField label="Email" :value="selectedUser.email" />
          <ReadonlyField label="User ID" :value="selectedUser.user_id" />
          <ReadonlyField label="Created" :value="formatUnixTime(selectedUser.created_at)" />
          <ReadonlyField label="Updated" :value="formatUnixTime(selectedUser.updated_at)" />
        </div>

        <template v-if="userMode === 'edit'">
          <div class="users-page__edit-grid">
            <label class="users-page__select">
              <span>Role</span>
              <select v-model="editForm.role" :disabled="isMutating">
                <option value="member">Member</option>
                <option value="admin">Admin</option>
              </select>
            </label>
            <label class="users-page__select">
              <span>Status</span>
              <select v-model="editForm.status" :disabled="isMutating">
                <option value="active">Active</option>
                <option value="disabled">Disabled</option>
                <option value="locked">Locked</option>
              </select>
            </label>
          </div>
          <form class="users-page__reset" @submit.prevent="submitPasswordReset">
            <FormField label="New password" for-id="reset-password">
              <PasswordInput id="reset-password" v-model="resetPassword" autocomplete="new-password" />
            </FormField>
            <BaseButton type="submit" variant="secondary" :busy="isMutating" :disabled="resetPassword.trim() === ''">
              Reset password
            </BaseButton>
          </form>
        </template>
      </div>
      <template #footer>
        <BaseButton v-if="userMode === 'edit'" :busy="isMutating" @click="saveUserEdits">Save changes</BaseButton>
        <BaseButton variant="ghost" @click="closeDrawer">Close</BaseButton>
      </template>
    </SideDrawer>
  </section>
</template>

<style scoped>
.users-page {
  display: grid;
  gap: 24px;
}

.users-page__filters {
  display: grid;
  grid-template-columns: minmax(180px, 1fr) repeat(2, minmax(140px, 180px));
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
  transition: border-color 160ms ease, box-shadow 160ms ease, background-color 160ms ease;
}

select:hover:not(:disabled) {
  border-color: var(--bb-color-primary);
}

select:disabled {
  color: var(--bb-color-muted);
  background: var(--bb-color-subtle);
  cursor: not-allowed;
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
  box-shadow: var(--bb-shadow-soft);
}

table {
  width: 100%;
  min-width: 820px;
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
  background: var(--bb-color-subtle);
}

tbody tr:hover {
  background: var(--bb-color-subtle);
}

tr:last-child td {
  border-bottom: 0;
}

td:first-child {
  display: grid;
  gap: 3px;
}

td:first-child span,
.users-page__count {
  color: var(--bb-color-muted);
}

.users-page__actions {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.users-page__drawer,
.users-page__reset {
  display: grid;
  gap: 16px;
}

.users-page__detail-grid,
.users-page__edit-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
}

.users-page__reset {
  border-top: 1px solid var(--bb-color-line);
  padding-top: 16px;
}

@media (max-width: 760px) {
  .users-page__filters,
  .users-page__detail-grid,
  .users-page__edit-grid {
    grid-template-columns: 1fr;
  }
}
</style>
