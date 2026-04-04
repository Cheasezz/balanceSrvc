package testsuite

import (
	"testing"
	"time"

	"github.com/Cheasezz/balanceSrvc/config"
	"github.com/Cheasezz/balanceSrvc/internal/app"
	"github.com/Cheasezz/balanceSrvc/pkg/logger"
	blnc "github.com/Cheasezz/balanceSrvc/protos/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	cfgPath = "../config/integrationTests.yml"
)

type TestSuite struct {
	*testing.T
	CtxTimeout    time.Duration
	BalanceClient blnc.BalanceClient
}

func New(t *testing.T) *TestSuite {
	t.Helper()

	cfg := config.MustLoadByPath(cfgPath)
	log := logger.New(cfg.Env)
	app := app.New(log, cfg)

	bufDialer := app.GRPCApp.RunBufConn()

	t.Cleanup(func() {
		t.Helper()
		app.Close()
	})

	cc, err := grpc.NewClient(
		"passthrough:///bufnet",
		grpc.WithContextDialer(bufDialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("grpc server connection failed: %v", err)
	}

	return &TestSuite{
		T:             t,
		CtxTimeout:    cfg.GRPC.Timeout,
		BalanceClient: blnc.NewBalanceClient(cc),
	}
}
