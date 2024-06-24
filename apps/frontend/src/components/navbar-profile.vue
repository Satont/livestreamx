<script setup lang="ts">
import { useProfile } from '@/api/profile.ts'
import { Button } from '@/components/ui/button'

const { data: profile, error } = useProfile()
const loginUri = `/api/auth/twitch?redirectUri=${window.location.origin}`
</script>

<template>
  <div class="flex items-center">
    <Button
      class="cursor-pointer"
      v-if="!profile || error"
      as="a"
      size="xs"
      :href="loginUri"
      >Login</Button
    >
    <Button
      v-else
      class="flex items-center gap-2 text-md"
      size="sm"
      variant="ghost"
    >
      {{ profile.userProfile.displayName }}
      <img
        :src="profile.userProfile.avatarUrl"
        alt="avatar"
        class="size-7 rounded-full"
      />
    </Button>
  </div>
</template>
