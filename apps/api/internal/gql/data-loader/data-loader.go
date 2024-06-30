package data_loader

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/satont/stream/apps/api/internal/gql/gqlmodel"
	"github.com/satont/stream/apps/api/internal/gql/mappers"
	"github.com/satont/stream/apps/api/internal/repositories/user"
	"github.com/vikstrous/dataloadgen"
)

type ctxKey string

const (
	LoadersKey = ctxKey("dataloaders")
)

type DataLoader struct {
	mapper   *mappers.Mapper
	userRepo user.Repository

	baseUserLoader *dataloadgen.Loader[uuid.UUID, *gqlmodel.BaseUser]
	chatUserLoader *dataloadgen.Loader[uuid.UUID, *gqlmodel.ChatUser]
}

type Opts struct {
	UserRepo user.Repository
	Mapper   *mappers.Mapper
}

func New(opts Opts) *DataLoader {
	loader := &DataLoader{
		mapper:   opts.Mapper,
		userRepo: opts.UserRepo,
	}

	loader.baseUserLoader = dataloadgen.NewLoader(
		loader.getBaseUsersByIds,
		dataloadgen.WithWait(time.Millisecond),
	)

	loader.chatUserLoader = dataloadgen.NewLoader(
		loader.getChatUsersByIds,
		dataloadgen.WithWait(time.Millisecond),
	)

	return loader
}

// For returns the dataloader for a given context
func GetLoaderForRequest(ctx context.Context) *DataLoader {
	return ctx.Value(LoadersKey).(*DataLoader)
}
