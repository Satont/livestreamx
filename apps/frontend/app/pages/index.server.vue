<script setup lang="ts">
import { formatDuration, intervalToDuration } from 'date-fns'

import { useStreamsList } from '~/api/stream'

const { data: streams } = await useStreamsList()

const now = new Date()

function computeStreamUptime(startedAt: Date) {
  const duration = intervalToDuration({
    start: startedAt,
    end: now
  })

  return formatDuration(duration, {
    format: ['hours', 'minutes']
  })
}
</script>

<template>
  <div class="w-full h-full justify-center mt-4 container">
    <h1 class="text-2xl">Browse channels</h1>

    <h1
      class="text-xl text-accent-foreground text-center"
      v-if="!streams?.streams.length"
    >
      No one online
    </h1>

    <div
      v-else
      class="flex pt-4 gap-4"
    >
      <NuxtLink
        v-for="stream of streams?.streams"
        :key="stream.channel.id"
        class="flex flex-col border border-border rounded-md bg-accent w-80"
        :to="stream.channel.name"
      >
        <div class="w-full h-full relative">
          <NuxtImg
            :src="stream.thumbnailUrl"
            class="w-auto h-44"
          />

          <span
            class="flex gap-2 items-center absolute py-0.5 px-2 bottom-1 right-1 bg-accent-foreground rounded-md text-sm text-accent"
          >
            <Icon
              name="lucide:eye"
              class="size-4"
            />
            {{ stream.viewers }}
          </span>

          <span
            class="flex gap-2 items-center absolute top-1 px-2 py-0.5 right-1 bg-accent-foreground rounded-md text-sm text-accent"
          >
            <Icon
              name="lucide:clock"
              class="size-4"
            />
            {{ computeStreamUptime(stream.startedAt) }}
          </span>
        </div>

        <div class="p-4 flex gap-2 h-auto">
          <NuxtImg
            :src="stream.channel.avatarUrl"
            class="size-6 rounded-full"
          />
          {{ stream.channel.displayName }}
        </div>
      </NuxtLink>
    </div>
  </div>
</template>
