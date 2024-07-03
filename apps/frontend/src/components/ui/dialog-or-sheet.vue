<script setup lang="ts">
import { breakpointsTailwind, useBreakpoints } from '@vueuse/core'

import { DialogContent } from '@/components/ui/dialog'
import { SheetContent } from '@/components/ui/sheet'

const breakPoints = useBreakpoints(breakpointsTailwind)
const md = breakPoints.greaterOrEqual('md')
</script>

<template>
  <DialogContent
    v-if="md"
    class="max-w-3xl overflow-hidden rounded-2xl outline-none"
  >
    <slot />
  </DialogContent>
  <SheetContent
    v-else
    side="bottom"
    class="flex flex-col"
  >
    <slot />
  </SheetContent>
</template>
