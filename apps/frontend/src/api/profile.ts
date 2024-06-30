import { useMutation, useQuery } from '@urql/vue'
import { createGlobalState } from '@vueuse/core'

import { graphql } from '@/gql'

export const useProfile = createGlobalState(() => {
  const useData = () =>
    useQuery({
      query: graphql(`
        query UserProfile {
          userProfile {
            id
            name
            displayName
            isBanned
            createdAt
            color
            avatarUrl
            isAdmin
            providers {
              provider
              userId
              name
              displayName
              avatarUrl
            }
            __typename
            streamKey
          }
        }
      `),
      variables: {}
    })

  const useUpdateMutation = () =>
    useMutation(
      graphql(`
        mutation UpdateProfile($input: UpdateUserProfileInput!) {
          updateUserProfile(input: $input) {
            __typename
          }
        }
      `)
    )

  const useLogout = () =>
    useMutation(
      graphql(`
        mutation UserLogout {
          logout
        }
      `)
    )

  const useDeleteAccount = () =>
    useMutation(
      graphql(`
        mutation DeleteAccount {
          deleteAccount
        }
      `)
    )

  return {
    useData,
    useUpdateMutation,
    useLogout,
    useDeleteAccount
  }
})
