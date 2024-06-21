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

const fontSize = computed({
	get() {
		return [chatFontSize.value]
	},
	set(values: number[]) {
		chatFontSize.value = values.at(0)
	}
})
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
					<h1>Font size</h1>
					<Slider
						v-model="fontSize"
						:min="10"
						:max="50"
						:step="1"
						class="w-full"
					/>
					<div class="flex justify-between">
						<span>
							10px
						</span>
						<span>
							{{ fontSize.at(0) }}px
						</span>
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
