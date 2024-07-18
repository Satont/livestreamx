<script setup lang="ts">
import { Copy, Eye, EyeOff } from 'lucide-vue-next'
import { computed, ref, watch } from 'vue'
import { toast } from 'vue-sonner'

import { useProfile } from '@/api/profile.ts'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'

const { useData, useUpdateMutation } = useProfile()
const { data: profile } = useData()
const streamAddress = `rtsp://${window.location.host}`
const updater = useUpdateMutation()

const form = ref<{
  sevenTvEmoteSetId?: string
}>({
  sevenTvEmoteSetId: undefined
})

watch(
  profile,
  (v) => {
    if (!v) return
    form.value.sevenTvEmoteSetId = v.userProfile.sevenTvEmoteSetId ?? undefined
  },
  { immediate: true }
)

const streamKey = computed(() => {
  const data = profile.value?.userProfile
  if (!data) {
    return ''
  }

  return `${data.name}?key=${data.streamKey}`
})

const showStreamKey = ref(false)
function copyText(text: string) {
  navigator.clipboard.writeText(text)
}

async function saveChanges() {
  const { error } = await updater.executeMutation({
    input: {
      sevenTvEmoteSetId: form.value.sevenTvEmoteSetId
    }
  })
  if (!error) {
    toast.success('Settings saved', {
      description: 'Your settings have been saved successfully',
      dismissible: true
    })
  }
}
</script>

<template>
  <div
    class="container flex flex-col gap-4 mt-4"
    v-if="profile"
  >
    <div class="flex flex-col gap-4">
      <Label for="streamServer">Stream server</Label>
      <div class="w-full relative">
        <Input
          id="streamServer"
          :value="streamAddress"
          :default-value="streamAddress"
          @input.prevent
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
          :default-value="streamKey"
          :value="streamKey"
          :type="showStreamKey ? 'text' : 'password'"
          class="pr-8"
          @input.prevent
        />
        <div class="absolute right-2 bottom-0 top-1.5 flex gap-0.5">
          <Button
            size="xs"
            variant="ghost"
            @click="copyText(streamKey)"
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

    <div class="flex flex-col gap-4">
      <Label for="seventv"> SevenTV Emote Set ID </Label>
      <Input
        id="seventv"
        v-model="form.sevenTvEmoteSetId"
        class="col-span-3"
        :maxlength="25"
      />
    </div>

    <Button
      @click="saveChanges"
      class="place-self-end w-max"
    >
      Save changes
    </Button>
  </div>
</template>
