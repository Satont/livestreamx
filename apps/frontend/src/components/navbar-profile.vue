<script setup lang="ts">
import { AlertCircle, Copy, Eye, EyeOff } from 'lucide-vue-next'
import { ref, watch } from 'vue'

import { useProfile } from '@/api/profile.js'
import NavbarProfileProviderButton from '@/components/navbar-profile-provider-button.vue'
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger
} from '@/components/ui/dialog'
import DialogOrSheet from '@/components/ui/dialog-or-sheet.vue'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import {
  Popover,
  PopoverContent,
  PopoverTrigger
} from '@/components/ui/popover'
import { Separator } from '@/components/ui/separator'
import { useProfileModalState } from '@/composables/use-profile-modal-state.js'
import { AuthedUserProviderType } from '@/gql/graphql.js'

const { useLogout, useData, useUpdateMutation, useDeleteAccount } = useProfile()
const { data: profile, executeQuery: refetchProfile } = useData()
const logout = useLogout()
const updater = useUpdateMutation()
const deleter = useDeleteAccount()
const { opened } = useProfileModalState()

const form = ref<{
  name: string
  displayName: string
  sevenTvEmoteSetId?: string
}>({
  name: '',
  displayName: '',
  sevenTvEmoteSetId: undefined
})
const formError = ref('')

watch(
  profile,
  (v) => {
    if (!v) return
    form.value.name = v.userProfile.name
    form.value.displayName = v.userProfile.displayName
    form.value.sevenTvEmoteSetId = v.userProfile.sevenTvEmoteSetId ?? undefined
  },
  { immediate: true }
)

async function logoutUser() {
  await logout.executeMutation({})
  refetchProfile({ requestPolicy: 'network-only' })
}

async function deleteAccount() {
  await deleter.executeMutation({})
  refetchProfile({ requestPolicy: 'network-only' })
}

async function saveChanges() {
  const { error: mutationError } = await updater.executeMutation({
    input: {
      name: form.value.name,
      displayName: form.value.displayName,
      sevenTvEmoteSetId: form.value.sevenTvEmoteSetId
    }
  })

  if (mutationError) {
    formError.value = mutationError.message
  } else {
    formError.value = ''
  }
}

const deleteConfirmationOpened = ref(false)

const streamAddress = `rtsp://${window.location.host}`

const showStreamKey = ref(false)
function copyText(text: string) {
  navigator.clipboard.writeText(text)
}
</script>

<template>
  <Dialog v-model:open="opened">
    <DialogTrigger>
      <Button
        v-if="!profile"
        class="cursor-pointer"
        as="a"
        size="xs"
      >
        Login
      </Button>
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
    </DialogTrigger>
    <DialogOrSheet>
      <DialogHeader>
        <DialogTitle>
          {{ profile ? 'Profile' : 'Login' }}
        </DialogTitle>
      </DialogHeader>

      <template v-if="profile">
        <div class="flex flex-col gap-4">
          <div class="grid grid-cols-4 items-center gap-4">
            <Label
              for="name"
              class="text-right"
            >
              Name
            </Label>
            <Input
              id="name"
              v-model="form.name"
              class="col-span-3"
              :maxlength="25"
            />
          </div>
          <div class="grid grid-cols-4 items-center gap-4">
            <Label
              for="username"
              class="text-right"
            >
              Display name
            </Label>
            <Input
              id="username"
              v-model="form.displayName"
              class="col-span-3"
              :maxlength="25"
            />
          </div>

          <Alert
            v-if="formError"
            variant="destructive"
          >
            <AlertCircle class="w-4 h-4" />
            <AlertTitle>Error</AlertTitle>
            <AlertDescription>
              {{ formError }}
            </AlertDescription>
          </Alert>
        </div>

        <Separator />
      </template>

      <div class="flex flex-col gap-2">
        <NavbarProfileProviderButton
          v-for="p of AuthedUserProviderType"
          :key="p"
          :provider="p"
        />
      </div>

      <template v-if="profile">
        <Separator />

        <div class="flex flex-col gap-2">
          <div class="flex flex-col gap-4">
            <Label for="streamServer">Stream server</Label>
            <div class="w-full relative">
              <Input
                id="streamServer"
                disabled
                :default-value="streamAddress"
                @click="copyText(streamAddress)"
                class="pr-8"
              />
              <Button
                size="xs"
                variant="ghost"
                @click="copyText(streamAddress)"
                class="absolute right-2 bottom-0 top-1.5"
              >
                <Copy class="size-4" />
              </Button>
            </div>
          </div>

          <div class="flex flex-col gap-4">
            <Label for="streamKey">Stream key</Label>
            <div class="w-full relative">
              <Input
                id="streamKey"
                :default-value="profile.userProfile.streamKey"
                :type="showStreamKey ? 'text' : 'password'"
                disabled
                class="pr-8"
              />
              <div class="absolute right-2 bottom-0 top-1.5 flex gap-0.5">
                <Button
                  size="xs"
                  variant="ghost"
                  @click="copyText(profile?.userProfile.streamKey)"
                >
                  <Copy class="size-4" />
                </Button>
                <Button
                  size="xs"
                  variant="ghost"
                  @click="showStreamKey = !showStreamKey"
                >
                  <component
                    :is="showStreamKey ? EyeOff : Eye"
                    class="size-4"
                  />
                </Button>
              </div>
            </div>
          </div>
        </div>
      </template>

      <template v-if="profile">
        <div class="grid grid-cols-4 items-center gap-4">
          <Label
            for="name"
            class="text-right"
          >
            SevenTV Emote Set
          </Label>
          <Input
            id="name"
            v-model="form.sevenTvEmoteSetId"
            class="col-span-3"
            :maxlength="25"
          />
        </div>
      </template>

      <template v-if="profile">
        <Separator />

        <DialogFooter
          v-if="!!profile"
          class="flex justify-between flex-col sm:flex-col md:flex-row gap-2"
        >
          <Popover v-model:open="deleteConfirmationOpened">
            <PopoverTrigger as-child>
              <Button
                variant="destructive"
                class="place-self-start w-full md:w-auto"
              >
                Delete account
              </Button>
            </PopoverTrigger>
            <PopoverContent class="flex flex-col gap-0.5">
              <span>Are you sure you wanna delete your account?</span>
              <div class="flex gap-2 justify-end">
                <Button
                  @click="deleteConfirmationOpened = false"
                  variant="secondary"
                >
                  Cancel
                </Button>
                <Button @click="deleteAccount">Confirm</Button>
              </div>
            </PopoverContent>
          </Popover>
          <Button
            @click="logoutUser"
            variant="destructive"
          >
            Logout
          </Button>
          <Button
            type="submit"
            @click="saveChanges"
          >
            Save changes
          </Button>
        </DialogFooter>
      </template>
    </DialogOrSheet>
  </Dialog>
</template>
