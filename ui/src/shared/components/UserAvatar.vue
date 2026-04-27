<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  name: string
  src?: string
  size?: 'sm' | 'md' | 'lg'
}>()

const avatarClass = computed(() => ['user-avatar', `user-avatar--${props.size ?? 'md'}`])
const initials = computed(() => {
  const words = props.name.trim().split(/\s+/).filter(Boolean)
  if (words.length === 0) {
    return 'U'
  }
  return words.slice(0, 2).map((word) => word.charAt(0).toUpperCase()).join('')
})
</script>

<template>
  <span :class="avatarClass" aria-hidden="true">
    <img v-if="src" :src="src" :alt="name" />
    <span v-else>{{ initials }}</span>
  </span>
</template>

<style scoped>
.user-avatar {
  flex: 0 0 auto;
  display: inline-grid;
  place-items: center;
  overflow: hidden;
  border: 1px solid rgb(255 255 255 / 72%);
  border-radius: 999px;
  color: #ffffff;
  background: var(--bb-color-primary);
  font-weight: 800;
  line-height: 1;
}

.user-avatar--sm {
  width: 28px;
  height: 28px;
  font-size: 0.72rem;
}

.user-avatar--md {
  width: 36px;
  height: 36px;
  font-size: 0.82rem;
}

.user-avatar--lg {
  width: 56px;
  height: 56px;
  font-size: 1rem;
}

.user-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}
</style>
