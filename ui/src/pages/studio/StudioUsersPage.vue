<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, reactive, shallowRef, watch } from 'vue'
import { Eye, KeyRound, Pencil, Trash2 } from 'lucide-vue-next'

import { useAuthStore } from '@/features/auth/stores/authStore'
import { studioApi } from '@/features/studio'
import type { StudioUpdateUserProfileRequest, StudioUser } from '@/features/studio'
import AvatarUploader from '@/shared/components/AvatarUploader.vue'
import BaseButton from '@/shared/components/BaseButton.vue'
import BaseInput from '@/shared/components/BaseInput.vue'
import ChangePasswordDialog from '@/shared/components/ChangePasswordDialog.vue'
import FormField from '@/shared/components/FormField.vue'
import ModalDialog from '@/shared/components/ModalDialog.vue'
import PageHeader from '@/shared/components/PageHeader.vue'
import PageLoadingState from '@/shared/components/PageLoadingState.vue'
import PasswordInput from '@/shared/components/PasswordInput.vue'
import ReadonlyField from '@/shared/components/ReadonlyField.vue'
import StatusAlert from '@/shared/components/StatusAlert.vue'
import StatusBadge from '@/shared/components/StatusBadge.vue'
import { useConfirm, useToast } from '@/shared/composables'

type UserMode = 'view' | 'edit'
type EditableRole = 'member' | 'admin'
type EditableStatus = 'active' | 'disabled' | 'locked'
type UserFormStatus = EditableStatus | 'pending' | 'deleted'

const authStore = useAuthStore()
const { confirm } = useConfirm()
const { pushToast } = useToast()

const users = shallowRef<StudioUser[]>([])
const total = shallowRef(0)
const isLoading = shallowRef(true)
const isMutating = shallowRef(false)
const errorMessage = shallowRef('')
const selectedUser = shallowRef<StudioUser | null>(null)
const passwordTarget = shallowRef<StudioUser | null>(null)
const userMode = shallowRef<UserMode>('view')
const resetPassword = shallowRef('')
const isSelfPasswordDialogOpen = shallowRef(false)
let filterTimer: number | undefined

const filters = reactive({
  keyword: '',
  role: '',
  status: '',
})

const editForm = reactive({
  role: 'member' as EditableRole,
  status: 'active' as UserFormStatus,
})

const profileForm = reactive({
  username: '',
  email: '',
  nickname: '',
  avatar_url: '',
})

const displayUsers = computed(() =>
  users.value.map((user) => ({
    ...user,
    displayName: user.nickname || user.username,
    lastLogin: formatUnixTime(user.last_login_at),
    isSelf: user.user_id === authStore.currentUser?.user_id,
  })),
)

const dialogTitle = computed(() => (userMode.value === 'edit' ? 'Edit user' : 'User details'))
const selectedUserIsSelf = computed(() => selectedUser.value?.user_id === authStore.currentUser?.user_id)

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
  editForm.status = normalizeFormStatus(user.status)
  profileForm.username = user.username
  profileForm.email = user.email
  profileForm.nickname = user.nickname ?? ''
  profileForm.avatar_url = user.avatar_url ?? ''
  resetPassword.value = ''
}

function openPasswordReset(user: StudioUser): void {
  if (user.user_id === authStore.currentUser?.user_id) {
    isSelfPasswordDialogOpen.value = true
    return
  }
  passwordTarget.value = user
  resetPassword.value = ''
}

function closeDialog(): void {
  selectedUser.value = null
  resetPassword.value = ''
}

function closePasswordReset(): void {
  passwordTarget.value = null
  resetPassword.value = ''
}

async function saveUserEdits(): Promise<void> {
  if (!selectedUser.value) {
    return
  }
  const target = selectedUser.value
  await runUserMutation(async () => {
    let nextUser = target
    const profilePatch = buildProfilePatch(target)
    if (Object.keys(profilePatch).length > 0) {
      nextUser = (await studioApi.updateUserProfile(target.user_id, profilePatch, { accessToken: authStore.accessToken })).user
    }
    if (!selectedUserIsSelf.value && editForm.role !== nextUser.role) {
      nextUser = (await studioApi.updateUserRole(target.user_id, { role: editForm.role }, { accessToken: authStore.accessToken })).user
    }
    if (!selectedUserIsSelf.value && isEditableStatus(editForm.status) && editForm.status !== nextUser.status) {
      nextUser = (await studioApi.updateUserStatus(target.user_id, { status: editForm.status }, { accessToken: authStore.accessToken })).user
    }
    replaceUser(nextUser)
    selectedUser.value = nextUser
    updateCurrentUserSnapshot(nextUser)
    pushToast({ tone: 'success', title: 'User updated', message: `${nextUser.email} has been updated.` })
  })
}

async function submitPasswordReset(): Promise<void> {
  if (!passwordTarget.value || resetPassword.value.trim() === '') {
    return
  }
  const approved = await confirm({
    title: 'Reset user password?',
    message: `This will replace the current password for ${passwordTarget.value.email} and revoke existing sessions.`,
    confirmText: 'Reset password',
    tone: 'danger',
  })
  if (!approved) {
    return
  }
  await runUserMutation(async () => {
    await studioApi.resetUserPassword(
      passwordTarget.value!.user_id,
      { new_password: resetPassword.value },
      { accessToken: authStore.accessToken },
    )
    pushToast({ tone: 'success', title: 'Password reset', message: `${passwordTarget.value!.email} can now use the new password.` })
    closePasswordReset()
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
      closeDialog()
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

function buildProfilePatch(target: StudioUser): StudioUpdateUserProfileRequest {
  const patch: StudioUpdateUserProfileRequest = {}
  const username = profileForm.username.trim()
  const email = profileForm.email.trim()
  const nickname = profileForm.nickname.trim()
  const avatarURL = profileForm.avatar_url.trim()

  if (username !== target.username) {
    patch.username = username
  }
  if (email !== target.email) {
    patch.email = email
  }
  if (nickname !== (target.nickname ?? '')) {
    patch.nickname = nickname
  }
  if (avatarURL !== (target.avatar_url ?? '')) {
    patch.avatar_url = avatarURL
  }

  return patch
}

function updateCurrentUserSnapshot(user: StudioUser): void {
  if (user.user_id !== authStore.currentUser?.user_id) {
    return
  }
  authStore.setCurrentUser({
    user_id: user.user_id,
    username: user.username,
    email: user.email,
    nickname: user.nickname ?? '',
    avatar_url: user.avatar_url ?? '',
    role: user.role,
    status: user.status,
  })
}

function normalizeRole(role: string): EditableRole {
  return role === 'admin' ? 'admin' : 'member'
}

function normalizeFormStatus(status: string): UserFormStatus {
  if (status === 'pending' || status === 'disabled' || status === 'locked' || status === 'deleted') {
    return status
  }
  return 'active'
}

function isEditableStatus(status: UserFormStatus): status is EditableStatus {
  return status === 'active' || status === 'disabled' || status === 'locked'
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
      description="Search accounts, review access, and manage user changes from row actions."
    />

    <div class="users-page__filters">
      <FormField class="users-page__search" label="Search" for-id="user-search">
        <BaseInput id="user-search" v-model="filters.keyword" placeholder="Search username, email, or nickname..." />
      </FormField>
      <label class="users-page__select">
        <span>Role</span>
        <select v-model="filters.role" class="users-page__select-control">
          <option value="">All roles</option>
          <option value="member">Member</option>
          <option value="admin">Admin</option>
        </select>
      </label>
      <label class="users-page__select">
        <span>Status</span>
        <select v-model="filters.status" class="users-page__select-control">
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
      <table class="users-page__grid">
        <thead class="users-page__head">
          <tr class="users-page__row">
            <th class="users-page__cell users-page__cell--user" scope="col">User</th>
            <th class="users-page__cell" scope="col">Role</th>
            <th class="users-page__cell" scope="col">Status</th>
            <th class="users-page__cell" scope="col">Last login</th>
            <th class="users-page__cell users-page__cell--actions" scope="col">Actions</th>
          </tr>
        </thead>
        <tbody class="users-page__body">
          <tr v-if="displayUsers.length === 0" class="users-page__row">
            <td class="users-page__cell users-page__empty" colspan="5">No users found.</td>
          </tr>
          <tr v-for="user in displayUsers" v-else :key="user.user_id" class="users-page__row">
            <td class="users-page__cell users-page__cell--user">
              <strong>{{ user.displayName }}</strong>
              <span>{{ user.email }}</span>
            </td>
            <td class="users-page__cell"><StatusBadge :value="user.role" /></td>
            <td class="users-page__cell"><StatusBadge :value="user.status" /></td>
            <td class="users-page__cell">{{ user.lastLogin }}</td>
            <td class="users-page__cell users-page__cell--actions">
              <div class="users-page__actions">
                <button
                  class="users-page__icon-action"
                  type="button"
                  :aria-label="`View ${user.email}`"
                  :title="`View ${user.email}`"
                  @click="openUser(user, 'view')"
                >
                  <Eye :size="17" aria-hidden="true" />
                </button>
                <button
                  class="users-page__icon-action users-page__icon-action--primary"
                  type="button"
                  :disabled="isMutating"
                  :aria-label="`Edit ${user.email}`"
                  :title="`Edit ${user.email}`"
                  @click="openUser(user, 'edit')"
                >
                  <Pencil :size="17" aria-hidden="true" />
                </button>
                <button
                  class="users-page__icon-action"
                  type="button"
                  :disabled="isMutating"
                  :aria-label="user.isSelf ? `Change password for ${user.email}` : `Reset password for ${user.email}`"
                  :title="user.isSelf ? `Change password for ${user.email}` : `Reset password for ${user.email}`"
                  @click="openPasswordReset(user)"
                >
                  <KeyRound :size="17" aria-hidden="true" />
                </button>
                <button
                  class="users-page__icon-action users-page__icon-action--danger"
                  type="button"
                  :disabled="isMutating"
                  :aria-label="`Delete ${user.email}`"
                  :title="`Delete ${user.email}`"
                  @click="deleteUser(user)"
                >
                  <Trash2 :size="17" aria-hidden="true" />
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <p v-if="!isLoading" class="users-page__count">{{ total }} total users</p>

    <ModalDialog :open="selectedUser !== null" :title="dialogTitle" :description="selectedUser?.email" size="lg" @close="closeDialog">
      <div v-if="selectedUser" class="users-page__modal">
        <div v-if="userMode === 'view'" class="users-page__detail-grid">
          <ReadonlyField label="Username" :value="selectedUser.username" />
          <ReadonlyField label="Nickname" :value="selectedUser.nickname" />
          <ReadonlyField label="Email" :value="selectedUser.email" />
          <ReadonlyField label="User ID" :value="selectedUser.user_id" />
          <ReadonlyField label="Created" :value="formatUnixTime(selectedUser.created_at)" />
          <ReadonlyField label="Updated" :value="formatUnixTime(selectedUser.updated_at)" />
        </div>

        <template v-if="userMode === 'edit'">
          <div class="users-page__edit-grid">
            <FormField label="Username" for-id="edit-username">
              <BaseInput id="edit-username" v-model="profileForm.username" autocomplete="username" />
            </FormField>
            <FormField label="Email" for-id="edit-email">
              <BaseInput id="edit-email" v-model="profileForm.email" autocomplete="email" />
            </FormField>
            <FormField label="Nickname" for-id="edit-nickname">
              <BaseInput id="edit-nickname" v-model="profileForm.nickname" autocomplete="nickname" />
            </FormField>
            <AvatarUploader class="users-page__avatar-upload" v-model="profileForm.avatar_url" :name="profileForm.nickname || profileForm.username" />
          </div>
          <div v-if="!selectedUserIsSelf" class="users-page__edit-grid">
            <label class="users-page__select">
              <span>Role</span>
              <select v-model="editForm.role" class="users-page__select-control" :disabled="isMutating">
                <option value="member">Member</option>
                <option value="admin">Admin</option>
              </select>
            </label>
            <label class="users-page__select">
              <span>Status</span>
              <select v-model="editForm.status" class="users-page__select-control" :disabled="isMutating">
                <option v-if="editForm.status === 'pending'" value="pending" disabled>Pending</option>
                <option value="active">Active</option>
                <option value="disabled">Disabled</option>
                <option value="locked">Locked</option>
                <option v-if="editForm.status === 'deleted'" value="deleted" disabled>Deleted</option>
              </select>
            </label>
          </div>
        </template>
      </div>
      <template #footer>
        <BaseButton v-if="userMode === 'edit'" :busy="isMutating" @click="saveUserEdits">Save changes</BaseButton>
        <BaseButton variant="ghost" @click="closeDialog">Close</BaseButton>
      </template>
    </ModalDialog>

    <ModalDialog
      :open="passwordTarget !== null"
      title="Reset password"
      :description="passwordTarget?.email"
      @close="closePasswordReset"
    >
      <form class="users-page__reset" novalidate @submit.prevent="submitPasswordReset">
        <FormField label="New password" for-id="reset-password">
          <PasswordInput id="reset-password" v-model="resetPassword" autocomplete="new-password" />
        </FormField>
      </form>
      <template #footer>
        <BaseButton :busy="isMutating" :disabled="resetPassword.trim() === ''" @click="submitPasswordReset">
          Reset password
        </BaseButton>
        <BaseButton variant="ghost" @click="closePasswordReset">Close</BaseButton>
      </template>
    </ModalDialog>

    <ChangePasswordDialog :open="isSelfPasswordDialogOpen" @close="isSelfPasswordDialogOpen = false" />
  </section>
</template>

<style scoped>
.users-page {
  display: grid;
  gap: 22px;
}

.users-page__filters {
  display: grid;
  grid-template-columns: minmax(220px, 1fr) repeat(2, minmax(150px, 180px));
  align-items: end;
  gap: 14px;
}

.users-page__search {
  min-width: 0;
}

.users-page__select {
  display: grid;
  gap: 6px;
  color: var(--bb-color-muted);
  font-size: 0.9rem;
  font-weight: 700;
}

.users-page__select-control {
  min-height: 44px;
  border: 1px solid var(--bb-color-line);
  border-radius: 8px;
  padding: 0 12px;
  color: var(--bb-color-text);
  background: var(--bb-color-surface);
  transition: border-color 160ms ease, box-shadow 160ms ease, background-color 160ms ease;
}

.users-page__select-control:hover:not(:disabled) {
  border-color: var(--bb-color-primary);
}

.users-page__select-control:disabled {
  color: var(--bb-color-muted);
  background: var(--bb-color-subtle);
  cursor: not-allowed;
}

.users-page__select-control:focus-visible,
.users-page__table:focus-visible {
  outline: none;
  box-shadow: 0 0 0 3px var(--bb-color-focus);
}

.users-page__table {
  overflow-x: auto;
  border: 1px solid var(--bb-color-line);
  border-radius: 10px;
  background: var(--bb-color-surface);
  box-shadow: var(--bb-shadow-soft);
}

.users-page__grid {
  width: 100%;
  min-width: 900px;
  border-collapse: collapse;
}

.users-page__cell {
  border-bottom: 1px solid var(--bb-color-line);
  padding: 13px 12px;
  text-align: left;
  vertical-align: middle;
}

.users-page__head .users-page__cell {
  color: var(--bb-color-muted);
  font-size: 0.78rem;
  text-transform: uppercase;
  background: var(--bb-color-surface);
}

.users-page__body .users-page__row:nth-child(even) {
  background: var(--bb-color-subtle);
}

.users-page__body .users-page__row:hover {
  background: var(--bb-color-primary-soft);
}

.users-page__body .users-page__row:last-child .users-page__cell {
  border-bottom: 0;
}

.users-page__cell--user {
  width: 34%;
}

.users-page__body .users-page__cell--user {
  display: grid;
  gap: 4px;
}

.users-page__cell--user strong {
  color: var(--bb-color-text-strong);
}

.users-page__cell--user span,
.users-page__count,
.users-page__empty {
  color: var(--bb-color-muted);
}

.users-page__cell--actions {
  width: 220px;
  text-align: right;
}

.users-page__actions {
  display: inline-flex;
  flex-wrap: wrap;
  justify-content: flex-end;
  gap: 8px;
  align-items: center;
}

.users-page__icon-action {
  width: 32px;
  height: 32px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border: 1px solid transparent;
  border-radius: 8px;
  color: var(--bb-color-muted);
  background: transparent;
  text-decoration: none;
  transition: color 160ms ease, border-color 160ms ease, background-color 160ms ease, box-shadow 160ms ease;
}

.users-page__icon-action:hover,
.users-page__icon-action:focus-visible {
  outline: none;
  color: var(--bb-color-text-strong);
  border-color: var(--bb-color-line);
  background: var(--bb-color-surface-elevated);
  box-shadow: 0 0 0 3px var(--bb-color-focus);
}

.users-page__icon-action--primary {
  color: var(--bb-color-text);
}

.users-page__icon-action--danger {
  color: var(--bb-color-danger);
}

.users-page__icon-action:disabled {
  opacity: 0.45;
  cursor: not-allowed;
}

.users-page__modal,
.users-page__reset {
  display: grid;
  gap: 18px;
}

.users-page__detail-grid,
.users-page__edit-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
}

.users-page__reset {
  border-top: 1px solid var(--bb-color-line);
  padding-top: 18px;
}

.users-page__avatar-upload {
  grid-column: 1 / -1;
}

@media (max-width: 780px) {
  .users-page__filters,
  .users-page__detail-grid,
  .users-page__edit-grid {
    grid-template-columns: 1fr;
  }
}
</style>
