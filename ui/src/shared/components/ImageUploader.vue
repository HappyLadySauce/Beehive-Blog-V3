<script setup lang="ts">
import { ImagePlus, X } from 'lucide-vue-next'
import { computed, useTemplateRef } from 'vue'
import { useI18n } from 'vue-i18n'

import { useAuthStore } from '@/features/auth/stores/authStore'
import { useAvatarUpload } from '@/features/uploads/useAvatarUpload'
import type { FileUploadScope } from '@/features/uploads/types'

import BaseButton from './BaseButton.vue'

const props = withDefaults(
  defineProps<{
    modelValue?: string
    scope?: FileUploadScope
    label?: string
    hint?: string
  }>(),
  {
    modelValue: '',
    scope: 'content_cover',
    label: 'Upload image',
    hint: 'PNG, JPEG, WebP, or AVIF. Max 5MB.',
  },
)

const emit = defineEmits<{
  'update:modelValue': [value: string]
}>()

const { t } = useI18n()
const authStore = useAuthStore()
const fileInput = useTemplateRef<HTMLInputElement>('fileInput')
const { isUploading, errorMessage, uploadImage } = useAvatarUpload()
const hasPreview = computed(() => props.modelValue.trim() !== '')
const resolvedLabel = computed(() => props.label || t('uploads.imageUpload'))
const resolvedHint = computed(() => props.hint || t('uploads.imageHint'))

function openFilePicker(): void {
  fileInput.value?.click()
}

function clearImage(): void {
  emit('update:modelValue', '')
}

async function handleFileChange(event: Event): Promise<void> {
  const input = event.target
  if (!(input instanceof HTMLInputElement) || !input.files?.[0]) {
    return
  }
  try {
    const publicURL = await uploadImage(input.files[0], authStore.accessToken, props.scope)
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
  <div class="image-uploader">
    <div class="image-uploader__preview" :class="{ 'image-uploader__preview--empty': !hasPreview }">
      <img v-if="hasPreview" :src="modelValue" alt="" />
      <ImagePlus v-else :size="24" aria-hidden="true" />
    </div>
    <div class="image-uploader__body">
      <div class="image-uploader__actions">
        <BaseButton variant="secondary" type="button" :busy="isUploading" @click="openFilePicker">
          <ImagePlus :size="16" aria-hidden="true" />
          {{ resolvedLabel }}
        </BaseButton>
        <BaseButton v-if="hasPreview" variant="ghost" type="button" @click="clearImage">
          <X :size="16" aria-hidden="true" />
          {{ t('uploads.removeImage') }}
        </BaseButton>
      </div>
      <p v-if="errorMessage" class="image-uploader__error">{{ errorMessage }}</p>
      <p v-else class="image-uploader__hint">{{ resolvedHint }}</p>
    </div>
    <input
      ref="fileInput"
      class="image-uploader__input"
      type="file"
      accept="image/png,image/jpeg,image/webp,image/avif"
      @change="handleFileChange"
    >
  </div>
</template>

<style scoped>
.image-uploader {
  display: grid;
  gap: 10px;
}

.image-uploader__preview {
  width: 100%;
  aspect-ratio: 16 / 9;
  overflow: hidden;
  border: 1px solid var(--bb-color-line);
  border-radius: 8px;
  background: var(--bb-color-surface-elevated);
}

.image-uploader__preview img {
  width: 100%;
  height: 100%;
  display: block;
  object-fit: cover;
}

.image-uploader__preview--empty {
  display: grid;
  place-items: center;
  color: var(--bb-color-muted);
  background: var(--bb-color-subtle);
}

.image-uploader__body,
.image-uploader__actions {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
  align-items: center;
}

.image-uploader__body {
  justify-content: space-between;
}

.image-uploader__hint,
.image-uploader__error {
  margin: 0;
  font-size: 0.86rem;
}

.image-uploader__hint {
  color: var(--bb-color-muted);
}

.image-uploader__error {
  color: var(--bb-color-danger);
}

.image-uploader__input {
  display: none;
}
</style>
