<script setup lang="ts">
import { LogOut, User, Video } from 'lucide-vue-next'
import { ref } from 'vue'

import { useProfile } from '@/api/profile.js'
import ProfileSettings from '@/components/profile-settings.vue'
import { Button } from '@/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger
} from '@/components/ui/dropdown-menu'

const { useData, useLogout } = useProfile()
const { data: profile, executeQuery: refetchProfile } = useData()
const logout = useLogout()

async function doLogout() {
  await logout.executeMutation({})
  await refetchProfile({ requestPolicy: 'network-only' })
}

const profileOpened = ref(false)
</script>

<template>
  <Button
    v-if="!profile"
    class="cursor-pointer"
    as="a"
    size="xs"
    @click="profileOpened = true"
  >
    Login
  </Button>

  <DropdownMenu v-else>
    <DropdownMenuTrigger>
      <Button
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
    </DropdownMenuTrigger>
    <DropdownMenuContent class="min-w-48">
      <DropdownMenuItem @select="profileOpened = true">
        <User class="dropdown-icon" />
        Profile
      </DropdownMenuItem>
      <DropdownMenuItem as-child>
        <router-link to="/dashboard/settings/stream">
          <Video class="dropdown-icon" /> Stream
        </router-link>
      </DropdownMenuItem>
      <DropdownMenuSeparator />
      <DropdownMenuItem @select="doLogout">
        <LogOut class="dropdown-icon" />
        Logout
      </DropdownMenuItem>
    </DropdownMenuContent>
  </DropdownMenu>

  <ProfileSettings v-model:modal-opened="profileOpened" />
</template>

<style scoped>
.dropdown-icon {
  @apply mr-2 size-4;
}
</style>
