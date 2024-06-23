<script setup lang="ts">
import { DropdownMenuRadioGroup, useForwardPropsEmits } from 'radix-vue'
import type {
  DropdownMenuRadioGroupEmits,
  DropdownMenuRadioGroupProps
} from 'radix-vue'

const props = defineProps<DropdownMenuRadioGroupProps>()
const emits = defineEmits<DropdownMenuRadioGroupEmits>()

const forwarded = useForwardPropsEmits(props, emits)
</script>

<template>
  <DropdownMenuRadioGroup v-bind="forwarded">
    <slot />
  </DropdownMenuRadioGroup>
</template>
