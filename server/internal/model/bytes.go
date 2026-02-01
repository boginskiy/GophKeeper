package model

import (
	"context"
	"os"
	"time"

	"github.com/boginskiy/GophKeeper/server/internal/manager"
	"github.com/boginskiy/GophKeeper/server/internal/utils"
)

type Bytes struct {
	Name         string
	Path         string
	Descr        *os.File
	SentSize     string
	ReceivedSize string
	Type         string
	Owner        string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// insertValuesFromCtx insert FileSize, FileName, FileOwner in *model.Bytes
func (b *Bytes) InsertValuesFromCtx(ctx context.Context) error {
	size, errSize := manager.TakeClientValueFromCtx(ctx, "total_size", 0)
	name, errName := manager.TakeClientValueFromCtx(ctx, "file_name", 0)
	owner, errOwner := manager.TakeServerValueFromCtx(ctx, manager.EmailCtx)

	if errSize != nil || errName != nil || errOwner != nil {
		return utils.DefinErr(errSize, errName, errOwner)
	}
	b.SentSize = size
	b.Name = name
	b.Owner = owner
	return nil
}

func (b *Bytes) GetOwner() string {
	return b.Owner
}

func (b *Bytes) GetFileName() string {
	return b.Name
}
