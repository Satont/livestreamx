import { createGlobalState } from "@vueuse/core";
import { ref } from "vue";
import { useMutation, useQuery, useSubscription } from "@urql/vue";
import { graphql } from "@/gql";
import { watch } from "vue";
import { MessageFragmentFragment } from "@/gql/graphql.ts";

export const useChat = createGlobalState(() => {
	const messages = ref<MessageFragmentFragment[]>([])

	const sub = useSubscription({
		query: graphql(`
			subscription NewChatMessages {
				chatMessages {
					...MessageFragment
				}
			}
		`),
		variables: {},
	})

	watch(sub.data, (data) => {
		if (!data) return

		messages.value = [...messages.value, data.chatMessages]
	})

	const initialMessages = useQuery({
		query: graphql(`
			query ChatMessages {
				chatMessagesLatest {
					...MessageFragment
				}
			}
		`),
		variables: {},
	})

	watch(initialMessages.data, (data) => {
		if (!data) return

		messages.value = data.chatMessagesLatest
	})

	const useSendMessage = () => useMutation(graphql(`
		mutation SendMessage($opts: SendMessageInput!) {
			sendMessage(input: $opts)
		}
	`))

	return {
		messages,
		useSendMessage,
	}
})