<script setup lang="ts">
import { Upload } from 'lucide-vue-next'
import { computed, useTemplateRef } from 'vue'
import { useI18n } from 'vue-i18n'

import { useAuthStore } from '@/features/auth/stores/authStore'
import { useAvatarUpload } from '@/features/file-manager/useAvatarUpload'

import BaseButton from './BaseButton.vue'
import UserAvatar from './UserAvatar.vue'

const props = withDefaults(
  defineProps<{
    modelValue?: string
    name?: string
  }>(),
  {
    modelValue: '',
    name: 'User',
  },
)

const emit = defineEmits<{
  'update:modelValue': [value: string]
}>()

const { t } = useI18n()
const authStore = useAuthStore()
const fileInput = useTemplateRef<HTMLInputElement>('fileInput')
const { isUploading, errorMessage, uploadAvatar } = useAvatarUpload()
const previewName = computed(() => props.name || t('uploads.avatarFallbackName'))

function openFilePicker(): void {
  fileInput.value?.click()
}

async function handleFileChange(event: Event): Promise<void> {
  const input = event.target
  if (!(input instanceof HTMLInputElement) || !input.files?.[0]) {
    return
  }
  try {
    const publicURL = await uploadAvatar(input.files[0], authStore.accessToken)
    emit('update:modelValue', publicURL)
  }
  catch {
    // Error state is owned by useAvatarUpload for inline rendering.
    // 错误状态由 useAvatarUpload 管理并在组件内展示。
  }
  finally {
    input.value = ''
  }
}
</script>

<template>
  <div class="avatar-uploader">
    <UserAvatar :name="previewName" :src="modelValue" size="lg" />
    <div class="avatar-uploader__body">
      <BaseButton variant="secondary" type="button" :busy="isUploading" @click="openFilePicker">
        <Upload :size="16" aria-hidden="true" />
        {{ t('uploads.avatarUpload') }}
      </BaseButton>
      <p v-if="errorMessage" class="avatar-uploader__error">{{ errorMessage }}</p>
      <p v-else class="avatar-uploader__hint">{{ t('uploads.avatarHint') }}</p>
    </div>
    <input
      ref="fileInput"
      class="avatar-uploader__input"
      type="file"
      accept="image/png,image/jpeg,image/webp,image/avif"
      @change="handleFileChange"
    >
  </div>
</template>

<style scoped>
.avatar-uploader {
  display: flex;
  align-items: center;
  gap: 14px;
}

.avatar-uploader__body {
  min-width: 0;
  display: grid;
  gap: 6px;
}

.avatar-uploader__hint,
.avatar-uploader__error {
  margin: 0;
  font-size: 0.86rem;
}

.avatar-uploader__hint {
  color: var(--bb-color-muted);
}

.avatar-uploader__error {
  color: var(--bb-color-danger);
}

.avatar-uploader__input {
  display: none;
}
</style>
