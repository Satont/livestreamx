<script setup lang="ts">
import { AlertCircle } from 'lucide-vue-next'
import { ref, watch } from 'vue'

import { useProfile, useProfileUpdate } from '@/api/profile.js'
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger
} from '@/components/ui/dialog'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'

const { data: profile } = useProfile()
const updater = useProfileUpdate()

const form = ref({
  name: '',
  displayName: ''
})
const error = ref('')

watch(
  profile,
  (v) => {
    if (!v) return
    form.value.name = v.userProfile.name
    form.value.displayName = v.userProfile.displayName
  },
  { immediate: true }
)

async function saveChanges() {
  const { error: mutationError } = await updater.executeMutation({
    input: {
      name: form.value.name,
      displayName: form.value.displayName
    }
  })

  if (mutationError) {
    error.value = mutationError.message
  } else {
    error.value = ''
  }
}
</script>

<template>
  <Dialog>
    <DialogTrigger as-child>
      <Button
        size="sm"
        variant="ghost"
        class="w-full flex justify-between"
        :disabled="!profile"
      >
        Change name
      </Button>
    </DialogTrigger>
    <DialogContent class="sm:max-w-[425px]">
      <DialogHeader>
        <DialogTitle>Edit profile</DialogTitle>
      </DialogHeader>
      <div class="grid gap-4 py-4">
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
            Username
          </Label>
          <Input
            id="username"
            v-model="form.displayName"
            class="col-span-3"
            :maxlength="25"
          />
        </div>

        <Alert
          v-if="error"
          variant="destructive"
        >
          <AlertCircle class="w-4 h-4" />
          <AlertTitle>Error</AlertTitle>
          <AlertDescription>
            {{ error }}
          </AlertDescription>
        </Alert>
      </div>
      <DialogFooter>
        <Button
          type="submit"
          @click="saveChanges"
        >
          Save changes
        </Button>
      </DialogFooter>
    </DialogContent>
  </Dialog>
</template>
