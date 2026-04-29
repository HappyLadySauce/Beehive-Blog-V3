<script setup lang="ts">
import { computed, reactive, shallowRef, watch } from 'vue'
import { useI18n } from 'vue-i18n'

import { authApi } from '@/features/auth/api/authApi'
import SsoProviderButtons from '@/features/auth/components/SsoProviderButtons.vue'
import { getPendingSsoEmail } from '@/features/auth/composables/useSsoFlow'
import { useAuthStore } from '@/features/auth/stores/authStore'
import type { AuthProvider } from '@/features/auth/types'
import { studioApi } from '@/features/studio'
import AvatarUploader from '@/shared/components/AvatarUploader.vue'
import BaseButton from '@/shared/components/BaseButton.vue'
import BaseInput from '@/shared/components/BaseInput.vue'
import FormField from '@/shared/components/FormField.vue'
import PageHeader from '@/shared/components/PageHeader.vue'
import PasswordInput from '@/shared/components/PasswordInput.vue'
import StatusAlert from '@/shared/components/StatusAlert.vue'
import { useToast } from '@/shared/composables'

interface SsoEmailMessage {
  type: 'beehive:sso-email'
  provider: AuthProvider
  code: string
  state: string
  redirectURI: string
}

const authStore = useAuthStore()
const { t } = useI18n()
const { pushToast } = useToast()
const form = reactive({
  nickname: '',
  avatarUrl: '',
})
const emailForm = reactive({
  email: '',
  currentPassword: '',
})
const isSavingProfile = shallowRef(false)
const isSavingEmail = shallowRef(false)
const successMessage = shallowRef('')
const errorMessage = shallowRef('')

const displayName = computed(() => form.nickname || authStore.currentUser?.username || t('profile.accountFallbackName'))
const currentEmail = computed(() => authStore.currentUser?.email ?? '')
const canSubmitEmailPassword = computed(() =>
  emailForm.email.trim().includes('@') && emailForm.currentPassword.trim().length > 0,
)

watch(
  () => authStore.currentUser,
  (user) => {
    form.nickname = user?.nickname ?? user?.username ?? ''
    form.avatarUrl = user?.avatar_url ?? ''
    emailForm.email = user?.email ?? ''
    emailForm.currentPassword = ''
  },
  { immediate: true },
)

async function saveProfile(): Promise<void> {
  successMessage.value = ''
  errorMessage.value = ''

  if (form.nickname.trim().length === 0) {
    errorMessage.value = t('profile.validation.displayNameRequired')
    return
  }

  isSavingProfile.value = true
  try {
    const response = await studioApi.updateProfile(
      {
        nickname: form.nickname.trim(),
        avatar_url: form.avatarUrl.trim(),
      },
      { accessToken: authStore.accessToken },
    )
    authStore.setCurrentUser(response.user)
    successMessage.value = t('profile.status.profileSaved')
    pushToast({ tone: 'success', title: t('profile.toast.profileSavedTitle'), message: t('profile.toast.profileSavedMessage') })
  }
  catch (error) {
    errorMessage.value = error instanceof Error ? error.message : t('profile.fallback.saveProfileFailed')
  }
  finally {
    isSavingProfile.value = false
  }
}

async function updateEmailWithPassword(): Promise<void> {
  successMessage.value = ''
  errorMessage.value = ''
  if (!canSubmitEmailPassword.value) {
    errorMessage.value = t('profile.validation.emailPasswordRequired')
    return
  }

  isSavingEmail.value = true
  try {
    const response = await authApi.updateEmail(
      {
        email: emailForm.email.trim(),
        verification_method: 'password',
        current_password: emailForm.currentPassword,
      },
      { accessToken: authStore.accessToken },
    )
    authStore.setCurrentUser(response.user)
    emailForm.currentPassword = ''
    successMessage.value = t('profile.status.emailUpdated')
    pushToast({ tone: 'success', title: t('profile.toast.emailUpdatedTitle'), message: response.user.email })
  }
  catch (error) {
    errorMessage.value = error instanceof Error ? error.message : t('profile.fallback.updateEmailFailed')
  }
  finally {
    isSavingEmail.value = false
  }
}

async function handleEmailAuthorized(message: SsoEmailMessage): Promise<void> {
  const pendingEmail = getPendingSsoEmail() || emailForm.email.trim()
  if (!pendingEmail) {
    errorMessage.value = t('profile.validation.pendingEmailMissing')
    return
  }

  isSavingEmail.value = true
  try {
    const response = await authApi.updateEmail(
      {
        email: pendingEmail,
        verification_method: 'sso',
        provider: message.provider,
        code: message.code,
        state: message.state,
        redirect_uri: message.redirectURI,
      },
      { accessToken: authStore.accessToken },
    )
    authStore.setCurrentUser(response.user)
    emailForm.email = response.user.email
    successMessage.value = t('profile.status.emailUpdated')
    pushToast({ tone: 'success', title: t('profile.toast.emailUpdatedTitle'), message: response.user.email })
  }
  catch (error) {
    errorMessage.value = error instanceof Error ? error.message : t('profile.fallback.updateEmailFailed')
  }
  finally {
    isSavingEmail.value = false
  }
}
</script>

<template>
  <section class="profile-page">
    <PageHeader
      :eyebrow="t('profile.eyebrow')"
      :title="t('profile.title')"
      :description="t('profile.description')"
    />

    <StatusAlert v-if="successMessage" tone="success" :title="t('profile.successTitle')">
      {{ successMessage }}
    </StatusAlert>
    <StatusAlert v-if="errorMessage" tone="danger" :title="t('profile.errorTitle')">
      {{ errorMessage }}
    </StatusAlert>

    <form class="profile-page__panel" novalidate @submit.prevent="saveProfile">
      <AvatarUploader v-model="form.avatarUrl" :name="displayName" />

      <FormField :label="t('common.displayName')" for-id="profile-nickname">
        <BaseInput id="profile-nickname" v-model="form.nickname" autocomplete="name" required />
      </FormField>

      <BaseButton type="submit" :busy="isSavingProfile">{{ t('profile.saveProfile') }}</BaseButton>
    </form>

    <section class="profile-page__panel" aria-labelledby="email-title">
      <div class="profile-page__section-heading">
        <h2 id="email-title">{{ t('profile.emailSectionTitle') }}</h2>
        <p>{{ t('profile.currentEmail', { email: currentEmail || t('common.none') }) }}</p>
      </div>

      <FormField :label="t('common.newEmail')" for-id="profile-email">
        <BaseInput id="profile-email" v-model="emailForm.email" type="email" autocomplete="email" inputmode="email" required />
      </FormField>

      <form class="profile-page__email-password" novalidate @submit.prevent="updateEmailWithPassword">
        <FormField :label="t('common.currentPassword')" for-id="profile-email-password">
          <PasswordInput id="profile-email-password" v-model="emailForm.currentPassword" autocomplete="current-password" />
        </FormField>
        <BaseButton type="submit" :busy="isSavingEmail" :disabled="!canSubmitEmailPassword">
          {{ t('profile.verifyWithPassword') }}
        </BaseButton>
      </form>

      <SsoProviderButtons
        surface="email"
        :email="emailForm.email.trim()"
        :access-token="authStore.accessToken"
        return-to="/account/profile"
        @email-authorized="handleEmailAuthorized"
      />
    </section>
  </section>
</template>

<style scoped>
.profile-page {
  display: grid;
  gap: 24px;
}

.profile-page__panel {
  width: min(640px, 100%);
  display: grid;
  gap: 16px;
  border: 1px solid var(--bb-color-line);
  border-radius: 10px;
  padding: 18px;
  background: var(--bb-color-surface);
  box-shadow: var(--bb-shadow-soft);
}

.profile-page__section-heading {
  display: grid;
  gap: 4px;
}

.profile-page__section-heading h2,
.profile-page__section-heading p {
  margin: 0;
}

.profile-page__section-heading h2 {
  color: var(--bb-color-text-strong);
  font-size: 1.05rem;
}

.profile-page__section-heading p {
  color: var(--bb-color-muted);
}

.profile-page__email-password {
  display: grid;
  grid-template-columns: minmax(0, 1fr) auto;
  align-items: end;
  gap: 12px;
}

@media (max-width: 640px) {
  .profile-page__email-password {
    grid-template-columns: 1fr;
  }
}
</style>
