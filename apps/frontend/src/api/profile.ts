import { useQuery } from "@urql/vue";
import { graphql } from "@/gql";
import { createGlobalState } from "@vueuse/core";

export const useProfile = createGlobalState(() => {
	return useQuery({
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
				}
			}
		`),
		variables: {}
	});
})