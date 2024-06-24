package mappers

import (
	"github.com/satont/stream/apps/api/internal/gql/gqlmodel"
	"github.com/satont/stream/apps/api/internal/repositories/role"
)

func (c *Mapper) DbRoleToGql(r role.Role) gqlmodel.Role {
	features := make([]gqlmodel.RoleFeature, 0, len(r.Features))
	for _, f := range r.Features {
		feature := gqlmodel.RoleFeature(f)
		if feature.IsValid() {
			features = append(features, gqlmodel.RoleFeature(f))
		}
	}

	return gqlmodel.Role{
		ID:        r.ID,
		Name:      r.Name,
		ImageURL:  r.ImageUrl,
		Features:  features,
		ChannelID: r.ChannelID,
	}
}
