import { createGlobalState } from "@vueuse/core";
import { ref } from "vue";
import { useMutation, useQuery, useSubscription } from "@urql/vue";
import { graphql, useFragment } from "@/gql";
import { watch } from "vue";
import { MessageFragmentFragment } from "@/gql/graphql.ts";

export const ChatMessage_Fragment = graphql(`
    fragment MessageFragment on ChatMessage {
        id
        segments {
            type
            content
            ...on MessageSegmentMention {
                user {
                    id
                    color
                    displayName
                }
            }
        }
        sender {
            id
            avatarUrl
            color
            createdAt
            name
            displayName
        }
        createdAt
    }
`)

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

		const fragment = useFragment(ChatMessage_Fragment, data.chatMessages)
		messages.value = [...messages.value, fragment]
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

		const fragments = useFragment(ChatMessage_Fragment, data.chatMessagesLatest)
		messages.value = fragments
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