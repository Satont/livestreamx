<script setup lang="ts">
import { useProfile } from '~/api/profile'
import Profile from '~/components/Profile.vue'

const { data: profile, executeQuery: refetchProfile } =
  await useProfile().useData()
const logout = useProfile().useLogout()

async function doLogout() {
  await logout.executeMutation({})
  await refetchProfile()
}
</script>

<template>
  <UiDialog>
    <UiDialogTrigger v-if="!profile?.userProfile">
      <UiButton size="sm"> Login </UiButton>
    </UiDialogTrigger>

    <UiDropdownMenu v-else>
      <UiDropdownMenuTrigger as-child>
        <UiButton
          class="text-md flex items-center gap-2"
          size="sm"
          variant="ghost"
        >
          {{ profile.userProfile.displayName }}
          <NuxtImg
            :src="profile.userProfile.avatarUrl"
            alt="avatar"
            class="size-7 rounded-full"
          />
        </UiButton>
      </UiDropdownMenuTrigger>
      <UiDropdownMenuContent class="min-w-48 mr-4">
        <UiDropdownMenuItem
          as-child
          class="w-full"
        >
          <UiDialogTrigger>
            <Icon
              name="lucide:user"
              class="size-4"
            />
            Profile
          </UiDialogTrigger>
        </UiDropdownMenuItem>
        <UiDropdownMenuItem as-child>
          <NuxtLink to="/dashboard/settings/stream">
            <Icon
              name="lucide:video"
              class="size-4"
            />
            Stream
          </NuxtLink>
        </UiDropdownMenuItem>
        <UiDropdownMenuSeparator />
        <UiDropdownMenuItem @click="doLogout">
          <Icon
            name="lucide:log-out"
            class="size-4"
          />
          Logout
        </UiDropdownMenuItem>
      </UiDropdownMenuContent>
    </UiDropdownMenu>
    <Profile />
  </UiDialog>
</template>
