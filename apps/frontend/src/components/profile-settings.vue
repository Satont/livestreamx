<script setup lang="ts">
import { AlertCircle } from 'lucide-vue-next'
import { ref, watch } from 'vue'

import { useProfile } from '@/api/profile.ts'
import NavbarProfileProviderButton from '@/components/navbar-profile-provider-button.vue'
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogFooter,
  DialogHeader,
  DialogTitle
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
import { AuthedUserProviderType } from '@/gql/graphql.js'

const modalOpened = defineModel('modalOpened', {
  required: true,
  default: false
})

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
  }
}

const deleteConfirmationOpened = ref(false)
function closeModal() {
  modalOpened.value = false
}
</script>

<template>
  <Dialog v-model:open="modalOpened">
    <DialogOrSheet
      @closeAutoFocus="closeModal"
      @escapeKeyDown="closeModal"
      @interactOutside="closeModal"
    >
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

        <Separator class="my-4" />
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
