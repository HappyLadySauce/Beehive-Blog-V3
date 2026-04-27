<script setup lang="ts">
import { RouterView } from 'vue-router'

import ConfirmDialogProvider from '@/shared/components/ConfirmDialogProvider.vue'
import ToastProvider from '@/shared/components/ToastProvider.vue'
import { useTheme } from '@/shared/composables'

useTheme()
</script>

<template>
  <RouterView />
  <ToastProvider />
  <ConfirmDialogProvider />
</template>
