<script setup lang="ts">
import { computed } from 'vue'

import { useProfile } from '@/api/profile.ts'
import { Button } from '@/components/ui/button'
import { AuthedUserProviderType } from '@/gql/graphql.ts'

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
</script>

<template>
  <Button
    class="w-full flex items-center gap-2 justify-between"
    variant="secondary"
    :as="providerData ? 'div' : 'a'"
    :href="authLink"
    :disabled="providerData"
  >
    <p class="text-md font-bold">
      {{ provider.slice(0, 1) + provider.toLowerCase().slice(1) }}
    </p>
    <div
      v-if="providerData"
      class="flex items-center gap-2"
    >
      <span class="text-md">{{ providerData.displayName }}</span>
      <img
        :src="providerData.avatarUrl"
        alt="avatar"
        class="size-7 rounded-full"
      />
    </div>
  </Button>
</template>
