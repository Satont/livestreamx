<script setup lang="ts">
import {
	DropdownMenu,
	DropdownMenuContent,
	DropdownMenuCheckboxItem,
	DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { Button } from '@/components/ui/button'
import { Settings } from 'lucide-vue-next'

import { showAvatars } from "@/composables/show-avatars.js";
import { showTimestamps } from "@/composables/show-timestamps.js";
</script>

<template>
	<DropdownMenu>
		<DropdownMenuTrigger asChild>
			<Button size="sm" variant="ghost"><Settings /></Button>
		</DropdownMenuTrigger>
		<DropdownMenuContent>
			<DropdownMenuCheckboxItem
				v-model:checked="showAvatars"
			>
				Show Avatars
			</DropdownMenuCheckboxItem>
			<DropdownMenuCheckboxItem
				v-model:checked="showTimestamps"
			>
				Show Time
			</DropdownMenuCheckboxItem>
		</DropdownMenuContent>
	</DropdownMenu>
</template>
