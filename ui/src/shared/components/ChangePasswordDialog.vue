<script setup lang="ts">
import { reactive, shallowRef, watch } from 'vue'
import { useI18n } from 'vue-i18n'

import { useAuthStore } from '@/features/auth/stores/authStore'
import { studioApi } from '@/features/studio'
import { useToast } from '@/shared/composables'

import BaseButton from './BaseButton.vue'
import FormField from './FormField.vue'
import ModalDialog from './ModalDialog.vue'
import PasswordInput from './PasswordInput.vue'
import StatusAlert from './StatusAlert.vue'

defineProps<{
  open: boolean
}>()

const emit = defineEmits<{
  close: []
  saved: []
}>()

const authStore = useAuthStore()
const { t } = useI18n()
const { pushToast } = useToast()
const isSaving = shallowRef(false)
const errorMessage = shallowRef('')
const form = reactive({
  oldPassword: '',
  newPassword: '',
})

watch(
  () => form.newPassword,
  () => {
    errorMessage.value = ''
  },
)

function resetForm(): void {
  form.oldPassword = ''
  form.newPassword = ''
  errorMessage.value = ''
}

async function submit(): Promise<void> {
  if (form.oldPassword.trim() === '' || form.newPassword.trim() === '') {
    errorMessage.value = t('password.validation.currentAndNewRequired')
    return
  }
  isSaving.value = true
  try {
    await studioApi.changePassword(
      {
        old_password: form.oldPassword,
        new_password: form.newPassword,
      },
      { accessToken: authStore.accessToken },
    )
    pushToast({ tone: 'success', title: t('password.toast.changedTitle'), message: t('password.toast.changedMessage') })
    emit('saved')
    emit('close')
    resetForm()
  }
  catch (error) {
    errorMessage.value = error instanceof Error ? error.message : t('password.fallback.changeFailed')
  }
  finally {
    isSaving.value = false
  }
}

function close(): void {
  emit('close')
  resetForm()
}
</script>

<template>
  <ModalDialog :open="open" :title="t('password.dialog.title')" :description="t('password.dialog.description')" @close="close">
    <form class="change-password-dialog" novalidate @submit.prevent="submit">
      <StatusAlert v-if="errorMessage" tone="danger" :title="t('password.dialog.errorTitle')">{{ errorMessage }}</StatusAlert>
      <FormField :label="t('common.currentPassword')" for-id="dialog-old-password">
        <PasswordInput id="dialog-old-password" v-model="form.oldPassword" autocomplete="current-password" required />
      </FormField>
      <FormField :label="t('common.newPassword')" for-id="dialog-new-password">
        <PasswordInput id="dialog-new-password" v-model="form.newPassword" autocomplete="new-password" required />
      </FormField>
    </form>
    <template #footer>
      <BaseButton :busy="isSaving" @click="submit">{{ t('password.dialog.submit') }}</BaseButton>
      <BaseButton variant="ghost" @click="close">{{ t('common.close') }}</BaseButton>
    </template>
  </ModalDialog>
</template>

<style scoped>
.change-password-dialog {
  display: grid;
  gap: 16px;
}
</style>
