<script setup lang="ts">
import { ref, watch } from 'vue'

import { useProfile } from '@/api/profile.js'
import { AuthedUserProviderType } from '@/gql/graphql.js'
import ProfileProvider from '~/components/profile/profile-provider.vue'

const { useData, useUpdateMutation, useDeleteAccount } = useProfile()
const updater = useUpdateMutation()
const deleter = useDeleteAccount()

const { data: profile, executeQuery: refetchProfile } = useData()

const form = ref<{
  name: string
  displayName: string
}>({
  name: '',
  displayName: ''
})
const formError = ref('')

watch(
  profile,
  (v) => {
    if (!v) return
    form.value.name = v.userProfile.name
    form.value.displayName = v.userProfile.displayName
  },
  { immediate: true }
)

async function deleteAccount() {
  await deleter.executeMutation({})
  refetchProfile({ requestPolicy: 'network-only' })
}

async function saveChanges() {
  const { error: mutationError } = await updater.executeMutation({
    input: {
      name: form.value.name,
      displayName: form.value.displayName
    }
  })

  if (mutationError) {
    formError.value = mutationError.message
  } else {
    formError.value = ''
    useToast().toast({
      title: 'Profile updated',
      description: 'Your profile has been updated',
      variant: 'success'
    })
  }
}

const deleteConfirmationOpened = ref(false)
</script>

<template>
  <UiDialogOrSheet>
    <UiDialogHeader>
      <UiDialogTitle>
        {{ profile ? 'Profile' : 'Login' }}
      </UiDialogTitle>
    </UiDialogHeader>

    <template v-if="profile">
      <div class="flex flex-col gap-4">
        <div class="grid grid-cols-4 items-center gap-4">
          <UiLabel
            for="name"
            class="text-right"
          >
            Name
          </UiLabel>
          <UiInput
            id="name"
            v-model="form.name"
            class="col-span-3"
            :maxlength="25"
          />
        </div>
        <div class="grid grid-cols-4 items-center gap-4">
          <UiLabel
            for="username"
            class="text-right"
          >
            Display name
          </UiLabel>
          <UiInput
            id="username"
            v-model="form.displayName"
            class="col-span-3"
            :maxlength="25"
          />
        </div>

        <UiAlert
          v-if="formError"
          variant="destructive"
        >
          <UiAlertTitle>Error</UiAlertTitle>
          <UiAlertDescription>
            {{ formError }}
          </UiAlertDescription>
        </UiAlert>
      </div>

      <UiSeparator class="mt-4" />
    </template>

    <div class="flex gap-2">
      <ProfileProvider
        v-for="p of AuthedUserProviderType"
        :key="p"
        :provider="p"
      />
    </div>

    <template v-if="profile">
      <UiSeparator />

      <UiDialogFooter
        v-if="!!profile"
        class="flex justify-between flex-col sm:flex-col md:flex-row gap-2"
      >
        <UiPopover v-model:open="deleteConfirmationOpened">
          <UiPopoverTrigger as-child>
            <UiButton
              variant="destructive"
              class="place-self-start w-full md:w-auto"
            >
              Delete account
            </UiButton>
          </UiPopoverTrigger>
          <UiPopoverContent class="flex flex-col gap-0.5">
            <span>Are you sure you wanna delete your account?</span>
            <div class="flex gap-2 justify-end">
              <UiButton
                @click="deleteConfirmationOpened = false"
                variant="secondary"
              >
                Cancel
              </UiButton>
              <UiButton @click="deleteAccount">Confirm</UiButton>
            </div>
          </UiPopoverContent>
        </UiPopover>
        <UiButton
          type="submit"
          @click="saveChanges"
        >
          Save changes
        </UiButton>
      </UiDialogFooter>
    </template>
  </UiDialogOrSheet>
</template>
