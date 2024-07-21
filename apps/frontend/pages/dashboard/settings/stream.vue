<script setup lang="ts">
import { ref, watch } from 'vue'

import { useProfile } from '@/api/profile.js'

const { useData, useUpdateMutation } = useProfile()
const { data: profile } = await useData()
const streamAddress = import.meta.client
  ? `rtsp://${window.location.hostname}`
  : ''
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

const showStreamKey = ref(false)
function copyText(text: string) {
  navigator.clipboard.writeText(text)
  useToast().toast({
    title: 'Copied',
    variant: 'success'
  })
}

async function saveChanges() {
  const { error } = await updater.executeMutation({
    input: {
      sevenTvEmoteSetId: form.value.sevenTvEmoteSetId
    }
  })
  if (!error) {
    useToast().toast({
      title: 'Settings saved',
      description: 'Your settings have been saved successfully',
      variant: 'success'
    })
  }
}

const streamKey = computed(() => {
  const data = profile.value?.userProfile
  if (!data) {
    return ''
  }
  return `${data.name}?key=${data.streamKey}`
})
</script>

<template>
  <div
    class="container flex flex-col gap-4 mt-4"
    v-if="profile"
  >
    <div class="flex flex-col gap-4">
      <UiLabel for="streamServer">Stream server</UiLabel>
      <div class="w-full relative">
        <UiInput
          id="streamServer"
          v-model="streamAddress"
          @input.prevent
          class="pr-8"
          readonly
        />
        <UiButton
          size="xs"
          variant="ghost"
          @click="copyText(streamAddress)"
          class="absolute right-2 bottom-0 top-1.5"
        >
          <Icon
            name="lucide:copy"
            class="size-4"
          />
        </UiButton>
      </div>
    </div>

    <div class="flex flex-col gap-4">
      <UiLabel for="streamKey">Stream key</UiLabel>
      <div class="w-full relative">
        <UiInput
          id="streamKey"
          v-model="streamKey"
          :type="showStreamKey ? 'text' : 'password'"
          class="pr-8"
          @input.prevent
          readonly
        />
        <div class="absolute right-2 bottom-0 top-1.5 flex gap-0.5">
          <UiButton
            size="xs"
            variant="ghost"
            @click="copyText(streamKey)"
          >
            <Icon
              name="lucide:copy"
              class="size-4"
            />
          </UiButton>
          <UiButton
            size="xs"
            variant="ghost"
            @click="showStreamKey = !showStreamKey"
          >
            <Icon
              :name="showStreamKey ? 'lucide:eye-off' : 'lucide:eye'"
              class="size-4"
            />
          </UiButton>
        </div>
      </div>
    </div>

    <div class="flex flex-col gap-4">
      <UiLabel for="seventv"> SevenTV Emote Set ID </UiLabel>
      <UiInput
        id="seventv"
        v-model="form.sevenTvEmoteSetId"
        class="col-span-3"
        :maxlength="25"
      />
    </div>

    <UiButton
      @click="saveChanges"
      class="place-self-end w-max"
    >
      Save changes
    </UiButton>
  </div>
</template>
