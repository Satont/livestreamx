import { useSubscription } from '@urql/vue'
import { createGlobalState } from '@vueuse/core'

import { useChat } from '@/api/chat.ts'
import { graphql } from '@/gql'

export const useStream = createGlobalState(() => {
  const { channelData } = useChat()

  const useStreamState = () =>
    useSubscription({
      query: graphql(`
        subscription StreamState($channelID: UUID!) {
          streamInfo(channelID: $channelID) {
            chatters {
              user {
                id
                displayName
                color
                createdAt
                name
                avatarUrl
              }
            }
            viewers
            startedAt
          }
        }
      `),
      get variables() {
        return {
          channelID: channelData.value?.fetchUserByName.id
        }
      },
      pause() {
        return !channelData.value?.fetchUserByName.id
      }
    })

  return {
    useStreamState
  }
})
