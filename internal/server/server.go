package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sake-kasu/sake-hack-backend/api/generated"
	"github.com/sake-kasu/sake-hack-backend/internal/config"
	"github.com/sake-kasu/sake-hack-backend/internal/database"
	sakeUsecase "github.com/sake-kasu/sake-hack-backend/internal/features/sake/application/usecase"
	sakeInfraQuery "github.com/sake-kasu/sake-hack-backend/internal/features/sake/infrastructure/query"
	sakeInfraRepo "github.com/sake-kasu/sake-hack-backend/internal/features/sake/infrastructure/repository"
	sakePresentation "github.com/sake-kasu/sake-hack-backend/internal/features/sake/presentation"
	storageExternal "github.com/sake-kasu/sake-hack-backend/internal/features/storage/infrastructure/external"
	"github.com/sake-kasu/sake-hack-backend/internal/logger"
	"github.com/sake-kasu/sake-hack-backend/internal/middleware"
	"github.com/valkey-io/valkey-go"
	"go.uber.org/zap"
)

// Server はHTTPサーバー
type Server struct {
	router       *gin.Engine
	httpServer   *http.Server
	postgresPool *pgxpool.Pool
	valkeyClient valkey.Client
	config       *config.Config
}

// New は新しいサーバーを作成する
func New(
	cfg *config.Config,
	postgresPool *pgxpool.Pool,
	valkeyClient valkey.Client,
) *Server {
	// Ginモード設定
	gin.SetMode(cfg.Server.Mode)

	// ルーター作成
	router := gin.New()

	// グローバルミドルウェア設定
	router.Use(gin.Recovery())
	router.Use(middleware.RequestID())
	router.Use(middleware.CORS(middleware.CORSConfig{
		AllowedOrigins:   cfg.CORS.AllowedOrigins,
		AllowedMethods:   cfg.CORS.AllowedMethods,
		AllowedHeaders:   cfg.CORS.AllowedHeaders,
		ExposedHeaders:   cfg.CORS.ExposedHeaders,
		AllowCredentials: cfg.CORS.AllowCredentials,
		MaxAge:           cfg.CORS.MaxAge,
	}))

	// サーバーインスタンス作成
	server := &Server{
		router:       router,
		postgresPool: postgresPool,
		valkeyClient: valkeyClient,
		config:       cfg,
	}

	// ルート設定
	server.setupRoutes()

	return server
}

// Start はサーバーを起動する
func (s *Server) Start() error {
	addr := fmt.Sprintf(":%d", s.config.Server.Port)

	s.httpServer = &http.Server{
		Addr:    addr,
		Handler: s.router,
	}

	logger.Get().Info("サーバーを起動中です...",
		zap.String("addr", addr),
		zap.String("mode", s.config.Server.Mode),
	)

	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("サーバーの起動に失敗しました: %w", err)
	}

	return nil
}

// Shutdown はサーバーをグレースフルシャットダウンする
func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.config.Server.GracefulShutdownTimeout)
	defer cancel()

	logger.Get().Info("サーバーをシャットダウンしています...")

	// HTTPサーバーシャットダウン
	if err := s.httpServer.Shutdown(ctx); err != nil {
		logger.Get().Error("サーバーのシャットダウン時にエラーが発生しました", zap.Error(err))
		return err
	}

	// データベース接続クローズ
	if s.postgresPool != nil {
		s.postgresPool.Close()
	}

	// Valkeyクライアントクローズ
	if s.valkeyClient != nil {
		s.valkeyClient.Close()
	}

	logger.Get().Info("サーバーを正常にシャットダウンしました")
	return nil
}

// HealthCheckResponse はヘルスチェックのレスポンス
type HealthCheckResponse struct {
	Status    string                       `json:"status"`
	Timestamp string                       `json:"timestamp"`
	Database  *HealthCheckDatabaseResponse `json:"database,omitempty"`
}

// HealthCheckDatabaseResponse はデータベースのヘルスチェック結果
type HealthCheckDatabaseResponse struct {
	Postgres string `json:"postgres,omitempty"`
	Valkey   string `json:"valkey,omitempty"`
}

// setupRoutes はルートを設定する
func (s *Server) setupRoutes() {
	// ヘルスチェックエンドポイント
	s.router.GET("/health", s.handleHealth)

	// S3クライアント初期化
	ctx := context.Background()
	s3Client, err := storageExternal.NewS3Client(ctx, s.config.Storage)
	if err != nil {
		logger.Get().Fatal("S3クライアントの初期化に失敗しました", zap.Error(err))
	}

	// Sake Feature
	sakeQuery := sakeInfraQuery.NewSakeQuery(s.postgresPool)
	sakeRepo := sakeInfraRepo.NewSakeRepository(s.postgresPool)

	listSakesUC := sakeUsecase.NewListSakesUsecase(sakeQuery)
	getSakeDetailUC := sakeUsecase.NewGetSakeDetailUsecase(sakeQuery)
	createStockUC := sakeUsecase.NewCreateStockUsecase(sakeRepo)
	updateStockUC := sakeUsecase.NewUpdateStockUsecase(sakeRepo)
	deleteStockUC := sakeUsecase.NewDeleteStockUsecase(sakeRepo)
	patchStockUC := sakeUsecase.NewPatchStockUsecase(sakeRepo)
	createUploadUrlUC := sakeUsecase.NewCreateUploadUrlUsecase(s3Client)
	listKindsUC := sakeUsecase.NewListKindsUsecase(sakeQuery)
	listBreweriesUC := sakeUsecase.NewListBreweriesUsecase(sakeQuery)
	listDrinkStylesUC := sakeUsecase.NewListDrinkStylesUsecase(sakeQuery)

	sakeServer := sakePresentation.NewSakeServerImpl(
		listSakesUC,
		getSakeDetailUC,
		createStockUC,
		updateStockUC,
		deleteStockUC,
		patchStockUC,
		createUploadUrlUC,
		listKindsUC,
		listBreweriesUC,
		listDrinkStylesUC,
		s3Client,
	)

	// APIグループを作成(/apiプレフィックス)
	apiGroup := s.router.Group("/api")

	// OpenAPI ServerInterfaceをGinに登録
	generated.RegisterHandlers(apiGroup, sakeServer)
}

// handleHealth はヘルスチェックハンドラ
func (s *Server) handleHealth(c *gin.Context) {
	ctx := c.Request.Context()

	response := HealthCheckResponse{
		Status:    "ok",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Database:  &HealthCheckDatabaseResponse{},
	}

	// PostgreSQLヘルスチェック
	if s.postgresPool != nil {
		if err := database.HealthCheckPostgresPool(ctx, s.postgresPool); err != nil {
			response.Database.Postgres = "error"
			response.Status = "degraded"
		} else {
			response.Database.Postgres = "ok"
		}
	}

	// Valkeyヘルスチェック
	if s.valkeyClient != nil {
		if err := database.HealthCheckValkey(ctx, s.valkeyClient); err != nil {
			response.Database.Valkey = "error"
			response.Status = "degraded"
		} else {
			response.Database.Valkey = "ok"
		}
	}

	statusCode := http.StatusOK
	if response.Status != "ok" {
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, response)
}
