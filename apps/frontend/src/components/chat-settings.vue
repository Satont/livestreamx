<script setup lang="ts">
import { Button } from '@/components/ui/button'
import { Settings } from 'lucide-vue-next'

import { showAvatars } from "@/composables/show-avatars.js";
import { showTimestamps } from "@/composables/show-timestamps.js";
import { Slider } from "@/components/ui/slider";
import { chatFontSize } from "@/composables/chat-font-size.ts";
import { computed } from "vue";
import {
	Popover,
	PopoverContent,
	PopoverTrigger,
} from '@/components/ui/popover'
import { Separator } from "@/components/ui/separator";
import { Switch } from "@/components/ui/switch";
import {
	NumberField,
	NumberFieldContent,
	NumberFieldDecrement,
	NumberFieldIncrement,
	NumberFieldInput,
} from '@/components/ui/number-field'

const fontSize = computed({
	get() {
		return [chatFontSize.value]
	},
	set(values: number[]) {
		chatFontSize.value = values.at(0)
	}
})

const min = 10
const max = 50
</script>

<template>
	<Popover>
		<PopoverTrigger asChild>
			<Button size="sm" variant="ghost"><Settings /></Button>
		</PopoverTrigger>
		<PopoverContent class="p-2 w-80">
			<div class="w-full flex flex-col">
				<Button
					@click="showAvatars = !showAvatars"
					size="sm"
					class="flex gap-2 justify-between"
					variant="ghost"
				>
					<span>
						Show avatars
					</span>
					<Switch :checked="showAvatars" class="data-[state=unchecked]:bg-zinc-600"  />
				</Button>
				<Button
					@click="showTimestamps = !showTimestamps"
					size="sm"
					class="flex gap-2 justify-between"
					variant="ghost"
				>
					<span>
						Show time
					</span>
					<Switch :checked="showTimestamps" class="data-[state=unchecked]:bg-zinc-600"  />
				</Button>
				<Separator class="my-4" />
				<div class="flex flex-col gap-2 px-3">
					<div class="flex justify-between items-center">
						<h1 class="text-lg text-slate-900 font-medium dark:text-slate-200">Font size</h1>
						<NumberField :min :max v-model="chatFontSize" class="w-36">
							<NumberFieldContent>
								<NumberFieldDecrement />
								<NumberFieldInput />
								<NumberFieldIncrement />
							</NumberFieldContent>
						</NumberField>
					</div>
					<Slider
						v-model="fontSize"
						:min
						:max
						:step="1"
						class="w-full"
					/>
					<div class="flex justify-between">
						<span>{{ min }}px</span>
						<span>{{ max }}px</span>
					</div>
				</div>
			</div>
		</PopoverContent>
	</Popover>

<!--	<DropdownMenu>-->
<!--		<DropdownMenuTrigger asChild>-->
<!--			<Button size="sm" variant="ghost"><Settings /></Button>-->
<!--		</DropdownMenuTrigger>-->
<!--		<DropdownMenuContent>-->
<!--			<DropdownMenuCheckboxItem-->
<!--				v-model:checked="showAvatars"-->
<!--			>-->
<!--				Show Avatars-->
<!--			</DropdownMenuCheckboxItem>-->
<!--			<DropdownMenuCheckboxItem-->
<!--				v-model:checked="showTimestamps"-->
<!--			>-->
<!--				Show Time-->
<!--			</DropdownMenuCheckboxItem>-->
<!--			<DropdownMenuItem>-->
<!--				<Slider-->
<!--					v-model="fontSize"-->
<!--					:max="50"-->
<!--					:step="1"-->
<!--					:class="cn('w-3/5', $attrs.class ?? '')"-->
<!--				/>-->
<!--			</DropdownMenuItem>-->
<!--		</DropdownMenuContent>-->
<!--	</DropdownMenu>-->
</template>
