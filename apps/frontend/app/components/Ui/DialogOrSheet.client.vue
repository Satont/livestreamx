<script setup lang="ts">
import { breakpointsTailwind, useBreakpoints } from '@vueuse/core'

const breakPoints = useBreakpoints(breakpointsTailwind)
const md = breakPoints.greaterOrEqual('md')
</script>

<template>
  <UiDialogContent
    v-if="md"
    class="max-w-3xl overflow-hidden rounded-2xl outline-none"
  >
    <slot />
  </UiDialogContent>
  <UiSheetContent
    v-else
    side="bottom"
    class="flex flex-col"
  >
    <slot />
  </UiSheetContent>
</template>
