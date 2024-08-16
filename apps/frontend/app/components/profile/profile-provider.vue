<script setup lang="ts">
import { useProfile } from '@/api/profile.js'
import { AuthedUserProviderType } from '@/gql/graphql.js'

const props = defineProps<{ provider: AuthedUserProviderType }>()

const { data: profile } = useProfile().useData()

const authLink = computed(() => {
  return `/api/auth/${props.provider.toLowerCase()}?redirectUri=${window.location.origin}`
})

const providerData = computed(() => {
  return profile.value?.userProfile?.providers?.find(
    (account) => account.provider === props.provider
  )
})

const iconName = computed(() => {
  switch (props.provider) {
    case AuthedUserProviderType.Github:
      return 'lucide:github'
    case AuthedUserProviderType.Twitch:
      return 'lucide:twitch'
  }
})

async function handleProviderAction() {
  window.location.replace(authLink.value)
}
</script>

<template>
  <div
    class="w-full flex items-center gap-2 justify-between bg-accent py-2 px-2 rounded"
    variant="secondary"
    :as="providerData ? 'div' : 'a'"
    :href="authLink"
  >
    <Icon
      :name="iconName"
      class="size-5"
    />
    <div
      v-if="providerData"
      class="flex items-center gap-2"
    >
      <span class="text-md">{{ providerData.displayName }}</span>
      <NuxtImg
        :src="providerData.avatarUrl"
        alt="avatar"
        class="size-7 rounded-full"
      />
    </div>
    <UiButton
      v-else
      size="xs"
      @click="handleProviderAction"
    >
      {{ profile?.userProfile ? 'Connect' : 'Login' }}
    </UiButton>
  </div>
</template>
