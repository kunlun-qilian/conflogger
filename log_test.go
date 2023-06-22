package conflogger

import (
	"context"
	"testing"

	"github.com/go-courier/metax"
	"github.com/sirupsen/logrus"
)

var logger = Log{
	Name:   "test",
	Level:  "Debug",
	Format: "json",
}

func init() {
	logger.SetDefaults()
	logger.Init()
}

func TestLog(t *testing.T) {
	ctx := metax.ContextWithMeta(context.Background(), metax.Meta{"_id": {"from context"}})
	ctx = metax.ContextWithMeta(ctx, metax.Meta{"_id": {"from context"}, "operator": {"test@0.0.0"}})

	logrus.WithContext(ctx).Info("Info")
	logrus.WithContext(ctx).Warning("Warn")
	logrus.WithContext(ctx).Error("Error")
	logrus.WithContext(ctx).WithField("test2", 2).Info("test")
}
