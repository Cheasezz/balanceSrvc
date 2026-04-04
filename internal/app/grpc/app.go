package grpcapp

// type Config struct {
// 	Port    int           `yaml:"port" env-required:"true"`
// 	Timeout time.Duration `yaml:"timeout"`
// }

// type App struct {
// 	log    logger.Logger
// 	Server *grpc.Server
// 	port   int
// }

// func New(l logger.Logger, cfg Config, s *service.Service) *App {
// 	loggingOpts := []logging.Option{
// 		logging.WithLogOnEvents(
// 			// logging.StartCall, logging.FinishCall,
// 			logging.PayloadReceived, logging.PayloadSent,
// 		),
// 	}

// 	recoveryOpts := []recovery.Option{
// 		recovery.WithRecoveryHandler(func(p any) (err error) {
// 			l.Error("Recovered from panic", slog.Any("panic", p))

// 			return status.Errorf(codes.Internal, "internal error")
// 		}),
// 	}

// 	gRPCServer := grpc.NewServer(grpc.ChainUnaryInterceptor(
// 		recovery.UnaryServerInterceptor(recoveryOpts...),
// 		logging.UnaryServerInterceptor(InterceptorLogger(l), loggingOpts...),
// 	))

// 	grpcSrv.Register(gRPCServer, l, s, cfg.Env)

// 	return &App{
// 		log:    l,
// 		Server: gRPCServer,
// 		port:   cfg.Port,
// 	}
// }

// func InterceptorLogger(l logger.Logger) logging.Logger {
// 	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
// 		l.Log(ctx, int(lvl), msg, fields...)
// 	})
// }

// func (a *App) MustRun() {
// 	if err := a.Run(); err != nil {
// 		panic(err)
// 	}
// }

// func (a *App) Run() error {
// 	const op = "grpcapp.Run"
// 	log := a.log.With("op", op)

// 	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
// 	if err != nil {
// 		return fmt.Errorf("%s: %w", op, err)
// 	}

// 	log.Info("gRPC server is running", "addr", l.Addr().String())

// 	if err = a.Server.Serve(l); err != nil {
// 		return fmt.Errorf("%s: %w", op, err)
// 	}

// 	return nil
// }

// func (a *App) RunBufConn() func(context.Context, string) (net.Conn, error) {
// 	lis := bufconn.Listen(1024 * 1024)

// 	go func() {
// 		if err := a.Server.Serve(lis); err != nil {
// 			panic(err)
// 		}
// 	}()

// 	bufDialer := func(context.Context, string) (net.Conn, error) {
// 		return lis.Dial()
// 	}

// 	return bufDialer
// }

// func (a *App) Close() {
// 	a.Server.GracefulStop()
// }
