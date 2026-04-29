<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, reactive, shallowRef, watch } from 'vue'
import { Eye, KeyRound, Pencil, Trash2 } from 'lucide-vue-next'
import { useI18n } from 'vue-i18n'

import { useAuthStore } from '@/features/auth/stores/authStore'
import { studioApi } from '@/features/studio'
import type { StudioUpdateUserProfileRequest, StudioUser } from '@/features/studio'
import AvatarUploader from '@/shared/components/AvatarUploader.vue'
import BaseButton from '@/shared/components/BaseButton.vue'
import BaseInput from '@/shared/components/BaseInput.vue'
import BaseSelect, { type BaseSelectOption } from '@/shared/components/BaseSelect.vue'
import ChangePasswordDialog from '@/shared/components/ChangePasswordDialog.vue'
import FormField from '@/shared/components/FormField.vue'
import IconActionButton from '@/shared/components/IconActionButton.vue'
import ModalDialog from '@/shared/components/ModalDialog.vue'
import PageHeader from '@/shared/components/PageHeader.vue'
import PageLoadingState from '@/shared/components/PageLoadingState.vue'
import PasswordInput from '@/shared/components/PasswordInput.vue'
import ReadonlyField from '@/shared/components/ReadonlyField.vue'
import StatusAlert from '@/shared/components/StatusAlert.vue'
import StatusBadge from '@/shared/components/StatusBadge.vue'
import { useConfirm, useToast } from '@/shared/composables'
import { useLocale } from '@/shared/i18n'

type UserMode = 'view' | 'edit'
type EditableRole = 'member' | 'admin'
type EditableStatus = 'active' | 'disabled' | 'locked'
type UserFormStatus = EditableStatus | 'pending' | 'deleted'

const authStore = useAuthStore()
const { t } = useI18n()
const { locale } = useLocale()
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
  includeDeleted: false,
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
    deletedAt: formatUnixTime(user.deleted_at),
    isSelf: user.user_id === authStore.currentUser?.user_id,
    isDeleted: user.status === 'deleted',
  })),
)

const roleOptions = computed<BaseSelectOption[]>(() => [
  { value: '', label: t('roles.all') },
  { value: 'member', label: t('roles.member') },
  { value: 'admin', label: t('roles.admin') },
])
const editableRoleOptions = computed<BaseSelectOption[]>(() => roleOptions.value.filter((option) => option.value !== ''))
const statusOptions = computed<BaseSelectOption[]>(() => [
  { value: '', label: t('userStatus.all') },
  { value: 'pending', label: t('userStatus.pending') },
  { value: 'active', label: t('userStatus.active') },
  { value: 'disabled', label: t('userStatus.disabled') },
  { value: 'locked', label: t('userStatus.locked') },
  { value: 'deleted', label: t('userStatus.deleted'), disabled: !filters.includeDeleted },
])
const editableStatusOptions = computed<BaseSelectOption[]>(() => [
  ...(editForm.status === 'pending' ? [{ value: 'pending', label: t('userStatus.pending'), disabled: true }] : []),
  { value: 'active', label: t('userStatus.active') },
  { value: 'disabled', label: t('userStatus.disabled') },
  { value: 'locked', label: t('userStatus.locked') },
  ...(editForm.status === 'deleted' ? [{ value: 'deleted', label: t('userStatus.deleted'), disabled: true }] : []),
])
const dialogTitle = computed(() => (userMode.value === 'edit' ? t('users.editDialog.title') : t('users.viewDialog.title')))
const selectedUserIsSelf = computed(() => selectedUser.value?.user_id === authStore.currentUser?.user_id)
const selectedUserIsDeleted = computed(() => selectedUser.value?.status === 'deleted')

function formatUnixTime(value?: number): string {
  if (!value) {
    return t('common.never')
  }
  return new Intl.DateTimeFormat(locale.value, {
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
        include_deleted: filters.includeDeleted,
        page: 1,
        page_size: 50,
      },
      { accessToken: authStore.accessToken },
    )
    users.value = response.items
    total.value = response.total
  }
  catch (error) {
    errorMessage.value = error instanceof Error ? error.message : t('users.unavailableTitle')
    pushToast({ tone: 'danger', title: t('users.unavailableTitle'), message: errorMessage.value })
  }
  finally {
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
  userMode.value = user.status === 'deleted' ? 'view' : mode
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
  if (!selectedUser.value || selectedUser.value.status === 'deleted') {
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
    pushToast({ tone: 'success', title: t('users.toast.savedTitle'), message: t('users.toast.savedMessage') })
  })
}

async function submitPasswordReset(): Promise<void> {
  if (!passwordTarget.value || passwordTarget.value.status === 'deleted' || resetPassword.value.trim() === '') {
    return
  }
  const approved = await confirm({
    title: t('users.passwordDialog.resetTitle'),
    message: t('users.actions.resetPassword', { email: passwordTarget.value.email }),
    confirmText: t('common.reset'),
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
    pushToast({ tone: 'success', title: t('users.toast.passwordTitle'), message: t('users.toast.passwordMessage') })
    closePasswordReset()
  })
}

async function deleteUser(user: StudioUser): Promise<void> {
  if (user.status === 'deleted') {
    return
  }
  const approved = await confirm({
    title: t('users.confirm.deleteTitle'),
    message: t('users.confirm.deleteMessage', { email: user.email }),
    confirmText: t('users.confirm.deleteAction'),
    tone: 'danger',
  })
  if (!approved) {
    return
  }
  await runUserMutation(async () => {
    await studioApi.deleteUser(user.user_id, { accessToken: authStore.accessToken })
    if (selectedUser.value?.user_id === user.user_id) {
      closeDialog()
    }
    if (passwordTarget.value?.user_id === user.user_id) {
      closePasswordReset()
    }
    await loadUsers()
    pushToast({ tone: 'success', title: t('users.toast.deletedTitle'), message: t('users.toast.deletedMessage', { email: user.email }) })
  })
}

async function runUserMutation(action: () => Promise<void>): Promise<void> {
  errorMessage.value = ''
  isMutating.value = true
  try {
    await action()
  }
  catch (error) {
    errorMessage.value = error instanceof Error ? error.message : t('users.unavailableTitle')
    pushToast({ tone: 'danger', title: t('users.unavailableTitle'), message: errorMessage.value })
  }
  finally {
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

watch(() => filters.includeDeleted, (includeDeleted) => {
  if (!includeDeleted && filters.status === 'deleted') {
    filters.status = ''
  }
})

watch(() => [filters.keyword, filters.role, filters.status, filters.includeDeleted], scheduleLoadUsers)

onMounted(loadUsers)
onBeforeUnmount(() => window.clearTimeout(filterTimer))
</script>

<template>
  <section class="users-page">
    <PageHeader
      :eyebrow="t('users.eyebrow')"
      :title="t('users.title')"
      :description="t('users.description')"
    />

    <div class="users-page__filters">
      <FormField class="users-page__search" :label="t('common.search')" for-id="user-search">
        <BaseInput id="user-search" v-model="filters.keyword" :placeholder="t('users.searchPlaceholder')" />
      </FormField>
      <FormField :label="t('users.columns.role')" for-id="user-role-filter">
        <BaseSelect id="user-role-filter" v-model="filters.role" :options="roleOptions" :aria-label="t('users.columns.role')" />
      </FormField>
      <FormField :label="t('users.columns.status')" for-id="user-status-filter">
        <BaseSelect id="user-status-filter" v-model="filters.status" :options="statusOptions" :aria-label="t('users.columns.status')" />
      </FormField>
      <label class="users-page__deleted-toggle">
        <input v-model="filters.includeDeleted" class="users-page__deleted-checkbox" type="checkbox">
        <span>{{ t('users.filters.includeDeleted') }}</span>
      </label>
    </div>
    <p class="users-page__filter-hint">{{ t('users.filters.includeDeletedHint') }}</p>

    <StatusAlert v-if="errorMessage" tone="danger" :title="t('users.unavailableTitle')">{{ errorMessage }}</StatusAlert>
    <PageLoadingState v-else-if="isLoading" :title="t('users.loadingTitle')" :rows="5" />

    <div v-else class="users-page__table" role="region" :aria-label="t('users.regionLabel')" tabindex="0">
      <table class="users-page__grid">
        <thead class="users-page__head">
          <tr class="users-page__row">
            <th class="users-page__cell users-page__cell--user" scope="col">{{ t('users.columns.user') }}</th>
            <th class="users-page__cell" scope="col">{{ t('users.columns.role') }}</th>
            <th class="users-page__cell" scope="col">{{ t('users.columns.status') }}</th>
            <th class="users-page__cell" scope="col">{{ t('users.columns.lastLogin') }}</th>
            <th class="users-page__cell users-page__cell--actions" scope="col">{{ t('common.actions') }}</th>
          </tr>
        </thead>
        <tbody class="users-page__body">
          <tr v-if="displayUsers.length === 0" class="users-page__row">
            <td class="users-page__cell users-page__empty" colspan="5">{{ t('users.empty') }}</td>
          </tr>
          <tr v-for="user in displayUsers" v-else :key="user.user_id" class="users-page__row">
            <td class="users-page__cell users-page__cell--user">
              <strong>{{ user.displayName }}</strong>
              <span>{{ user.email }}</span>
              <span v-if="user.isDeleted">{{ t('users.deletedAt', { value: user.deletedAt }) }}</span>
            </td>
            <td class="users-page__cell"><StatusBadge :value="user.role" /></td>
            <td class="users-page__cell"><StatusBadge :value="user.status" /></td>
            <td class="users-page__cell">{{ user.lastLogin }}</td>
            <td class="users-page__cell users-page__cell--actions">
              <div class="users-page__actions">
                <IconActionButton :aria-label="t('users.actions.viewUser', { email: user.email })" :title="t('users.actions.viewUser', { email: user.email })" @click="openUser(user, 'view')">
                  <Eye :size="17" aria-hidden="true" />
                </IconActionButton>
                <IconActionButton
                  tone="primary"
                  :disabled="isMutating || user.isDeleted"
                  :aria-label="user.isDeleted ? t('users.actions.editDeletedUser', { email: user.email }) : t('users.actions.editUser', { email: user.email })"
                  :title="user.isDeleted ? t('users.actions.editDeletedUser', { email: user.email }) : t('users.actions.editUser', { email: user.email })"
                  @click="openUser(user, 'edit')"
                >
                  <Pencil :size="17" aria-hidden="true" />
                </IconActionButton>
                <IconActionButton
                  :disabled="isMutating || user.isDeleted"
                  :aria-label="user.isDeleted ? t('users.actions.resetDeletedUserPassword', { email: user.email }) : user.isSelf ? t('users.actions.changePassword', { email: user.email }) : t('users.actions.resetPassword', { email: user.email })"
                  :title="user.isDeleted ? t('users.actions.resetDeletedUserPassword', { email: user.email }) : user.isSelf ? t('users.actions.changePassword', { email: user.email }) : t('users.actions.resetPassword', { email: user.email })"
                  @click="openPasswordReset(user)"
                >
                  <KeyRound :size="17" aria-hidden="true" />
                </IconActionButton>
                <IconActionButton
                  tone="danger"
                  :disabled="isMutating || user.isDeleted"
                  :aria-label="user.isDeleted ? t('users.actions.deleteDeletedUser', { email: user.email }) : t('users.actions.deleteUser', { email: user.email })"
                  :title="user.isDeleted ? t('users.actions.deleteDeletedUser', { email: user.email }) : t('users.actions.deleteUser', { email: user.email })"
                  @click="deleteUser(user)"
                >
                  <Trash2 :size="17" aria-hidden="true" />
                </IconActionButton>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <p v-if="!isLoading" class="users-page__count">{{ t('users.count', { count: total }) }}</p>

    <ModalDialog :open="selectedUser !== null" :title="dialogTitle" :description="selectedUser?.email" size="lg" @close="closeDialog">
      <div v-if="selectedUser" class="users-page__modal">
        <StatusAlert
          v-if="selectedUserIsDeleted"
          tone="warning"
          :title="t('users.deletedBanner.title')"
        >
          {{ t('users.deletedBanner.message') }}
        </StatusAlert>
        <div v-if="userMode === 'view'" class="users-page__detail-grid">
          <ReadonlyField :label="t('users.editDialog.username')" :value="selectedUser.username" />
          <ReadonlyField :label="t('users.editDialog.nickname')" :value="selectedUser.nickname" />
          <ReadonlyField :label="t('users.editDialog.email')" :value="selectedUser.email" />
          <ReadonlyField :label="t('users.columns.role')" :value="t(`roles.${normalizeRole(selectedUser.role)}`)" />
          <ReadonlyField :label="t('users.columns.status')" :value="t(`userStatus.${normalizeFormStatus(selectedUser.status)}`)" />
          <ReadonlyField :label="t('users.viewDialog.userId')" :value="selectedUser.user_id" />
          <ReadonlyField :label="t('users.viewDialog.created')" :value="formatUnixTime(selectedUser.created_at)" />
          <ReadonlyField :label="t('users.viewDialog.updated')" :value="formatUnixTime(selectedUser.updated_at)" />
          <ReadonlyField v-if="selectedUser.deleted_at" :label="t('users.columns.deletedAt')" :value="formatUnixTime(selectedUser.deleted_at)" />
        </div>

        <template v-if="userMode === 'edit'">
          <div class="users-page__avatar-section">
            <AvatarUploader class="users-page__avatar-upload" v-model="profileForm.avatar_url" :name="profileForm.nickname || profileForm.username" />
          </div>
          <div class="users-page__edit-grid">
            <FormField :label="t('users.editDialog.username')" for-id="edit-username">
              <BaseInput id="edit-username" v-model="profileForm.username" autocomplete="username" />
            </FormField>
            <FormField :label="t('users.editDialog.email')" for-id="edit-email">
              <BaseInput id="edit-email" v-model="profileForm.email" autocomplete="email" />
            </FormField>
            <FormField :label="t('users.editDialog.nickname')" for-id="edit-nickname">
              <BaseInput id="edit-nickname" v-model="profileForm.nickname" autocomplete="nickname" />
            </FormField>
          </div>
          <div v-if="!selectedUserIsSelf" class="users-page__edit-grid">
            <FormField :label="t('users.columns.role')" for-id="edit-role">
              <BaseSelect id="edit-role" v-model="editForm.role" :options="editableRoleOptions" :disabled="isMutating" :aria-label="t('users.columns.role')" />
            </FormField>
            <FormField :label="t('users.columns.status')" for-id="edit-status">
              <BaseSelect id="edit-status" v-model="editForm.status" :options="editableStatusOptions" :disabled="isMutating" :aria-label="t('users.columns.status')" />
            </FormField>
          </div>
        </template>
      </div>
      <template #footer>
        <BaseButton v-if="userMode === 'edit' && !selectedUserIsDeleted" :busy="isMutating" @click="saveUserEdits">{{ t('common.saveChanges') }}</BaseButton>
        <BaseButton variant="ghost" @click="closeDialog">{{ t('common.close') }}</BaseButton>
      </template>
    </ModalDialog>

    <ModalDialog
      :open="passwordTarget !== null"
      :title="t('users.passwordDialog.resetTitle')"
      :description="passwordTarget?.email"
      @close="closePasswordReset"
    >
      <form class="users-page__reset" novalidate @submit.prevent="submitPasswordReset">
        <FormField :label="t('users.passwordDialog.newPassword')" for-id="reset-password">
          <PasswordInput id="reset-password" v-model="resetPassword" autocomplete="new-password" />
        </FormField>
      </form>
      <template #footer>
        <BaseButton :busy="isMutating" :disabled="resetPassword.trim() === ''" @click="submitPasswordReset">
          {{ t('common.reset') }}
        </BaseButton>
        <BaseButton variant="ghost" @click="closePasswordReset">{{ t('common.close') }}</BaseButton>
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
  grid-template-columns: minmax(220px, 1fr) repeat(2, minmax(150px, 180px)) minmax(180px, 220px);
  align-items: end;
  gap: 14px;
}

.users-page__search {
  min-width: 0;
}

.users-page__deleted-toggle {
  min-height: 44px;
  display: inline-flex;
  align-items: center;
  gap: 10px;
  color: var(--bb-color-text);
  cursor: pointer;
}

.users-page__deleted-checkbox {
  width: 16px;
  height: 16px;
  accent-color: var(--bb-color-primary);
}

.users-page__filter-hint {
  margin-top: -8px;
  color: var(--bb-color-muted);
  font-size: 0.92rem;
}

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

.users-page__avatar-section {
  border: 1px solid var(--bb-color-line);
  border-radius: 10px;
  padding: 14px;
  background: var(--bb-color-subtle);
}

.users-page__avatar-upload {
  width: 100%;
}

@media (max-width: 780px) {
  .users-page__filters,
  .users-page__detail-grid,
  .users-page__edit-grid {
    grid-template-columns: 1fr;
  }
}
</style>
