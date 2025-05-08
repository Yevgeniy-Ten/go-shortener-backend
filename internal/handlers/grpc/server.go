package grpc

import (
	"context"
	"errors"
	"net"
	"shorter/internal/cookies"
	"shorter/internal/domain"
	"shorter/internal/handlers"
	"shorter/internal/logger"
	"shorter/internal/urlstorage/database/urls"
	"shorter/pkg"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type ShortenerGRPCServer struct {
	UnimplementedShortenerServer
	Storage domain.Storage
	Log     *logger.ZapLogger
	Config  *handlers.Config
}

func (s *ShortenerGRPCServer) ShortenURL(ctx context.Context, req *ShortenURLRequest) (*ShortenURLResponse, error) {
	userID, err := cookies.GetUserIDFromMetadata(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "unauthorized: "+err.Error())
	}
	if !pkg.ValidateURL(req.OriginalUrl) {
		return nil, status.Error(codes.InvalidArgument, "Некорректный ShortURL.")
	}
	urlID, err := s.Storage.URLS.Save(req.OriginalUrl, userID)
	if err != nil {
		var duplicateError *urls.DuplicateError
		if errors.As(err, &duplicateError) {
			return &ShortenURLResponse{ShortUrl: s.Config.ServerAddr + "/" + duplicateError.ShortURL}, status.Error(codes.AlreadyExists, "duplicate")
		}
		s.Log.Log.Error("gRPC: Error when save", zap.Error(err))
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &ShortenURLResponse{ShortUrl: s.Config.ServerAddr + "/" + urlID}, nil
}

func (s *ShortenerGRPCServer) ShortenURLsBatch(ctx context.Context, req *ShortenURLsBatchRequest) (*ShortenURLsBatchResponse, error) {
	userID, err := cookies.GetUserIDFromMetadata(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "unauthorized: "+err.Error())
	}
	if len(req.OriginalUrls) == 0 {
		return nil, status.Error(codes.InvalidArgument, "empty batch")
	}
	urlsForSave := make([]domain.URLS, 0, len(req.OriginalUrls))
	for i, url := range req.OriginalUrls {
		if !pkg.ValidateURL(url) {
			return nil, status.Errorf(codes.InvalidArgument, "Некорректный ShortURL в позиции %d", i)
		}
		urlsForSave = append(urlsForSave, domain.URLS{CorrelationID: "", URL: url})
	}
	err = s.Storage.URLS.SaveBatch(urlsForSave, userID)
	if err != nil {
		s.Log.Log.Error("gRPC: Error when save batch", zap.Error(err))
		return nil, status.Error(codes.Internal, "internal error")
	}
	shortUrls := make([]string, 0, len(urlsForSave))
	for _, u := range urlsForSave {
		shortUrls = append(shortUrls, s.Config.ServerAddr+"/"+u.CorrelationID)
	}
	return &ShortenURLsBatchResponse{ShortUrls: shortUrls}, nil
}

func (s *ShortenerGRPCServer) GetOriginalURL(ctx context.Context, req *GetOriginalURLRequest) (*GetOriginalURLResponse, error) {
	url, err := s.Storage.URLS.GetURL(req.ShortUrl)
	if err != nil {
		var urlIsDeletedError *urls.URLIsDeletedError
		if errors.As(err, &urlIsDeletedError) {
			return nil, status.Error(codes.NotFound, "url is deleted")
		}
		s.Log.Log.Error("gRPC: Error when get", zap.Error(err))
		return nil, status.Error(codes.NotFound, "not found")
	}
	if url == "" {
		return nil, status.Error(codes.InvalidArgument, "empty url")
	}
	return &GetOriginalURLResponse{OriginalUrl: url}, nil
}

func (s *ShortenerGRPCServer) GetUserURLs(ctx context.Context, req *GetUserURLsRequest) (*GetUserURLsResponse, error) {
	userID, err := cookies.GetUserIDFromMetadata(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "unauthorized: "+err.Error())
	}
	urls, err := s.Storage.URLS.GetUserURLs(userID, s.Config.ServerAddr)
	if err != nil {
		s.Log.Log.Error("gRPC: Error when get user urls", zap.Error(err))
		return nil, status.Error(codes.Internal, "internal error")
	}
	if len(urls) == 0 {
		return nil, status.Error(codes.NotFound, "no urls for user")
	}
	resp := &GetUserURLsResponse{}
	for _, u := range urls {
		resp.Urls = append(resp.Urls, &UserURL{
			ShortUrl:    u.ShortURL,
			OriginalUrl: u.OriginalURL,
		})
	}
	return resp, nil
}

func (s *ShortenerGRPCServer) DeleteUserURLs(ctx context.Context, req *DeleteUserURLsRequest) (*DeleteUserURLsResponse, error) {
	userID, err := cookies.GetUserIDFromMetadata(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "unauthorized: "+err.Error())
	}
	if len(req.ShortUrls) == 0 {
		return nil, status.Error(codes.InvalidArgument, "no urls to delete")
	}
	err = s.Storage.URLS.DeleteURLs(req.ShortUrls, userID)
	if err != nil {
		s.Log.Log.Error("gRPC: Error when delete urls", zap.Error(err))
		return &DeleteUserURLsResponse{Success: false}, status.Error(codes.Internal, "internal error")
	}
	return &DeleteUserURLsResponse{Success: true}, nil
}

func (s *ShortenerGRPCServer) GetInternalStats(ctx context.Context, req *GetInternalStatsRequest) (*GetInternalStatsResponse, error) {
	if s.Config.TrustedSubnet == "" {
		return nil, status.Error(codes.PermissionDenied, "no trusted subnet configured")
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.PermissionDenied, "no metadata")
	}
	ips := md["x-real-ip"]
	if len(ips) == 0 {
		return nil, status.Error(codes.PermissionDenied, "no X-Real-IP header")
	}
	clientIP := ips[0]

	_, ipnet, err := net.ParseCIDR(s.Config.TrustedSubnet)
	if err != nil {
		return nil, status.Error(codes.Internal, "invalid trusted subnet")
	}
	if !ipnet.Contains(net.ParseIP(clientIP)) {
		return nil, status.Error(codes.PermissionDenied, "forbidden")
	}

	stats, err := s.Storage.URLS.GetStats()
	if err != nil {
		s.Log.Log.Error("gRPC: Error when get stats", zap.Error(err))
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &GetInternalStatsResponse{
		Urls:  int64(stats.URLs),
		Users: int64(stats.Users),
	}, nil
}

func RunGRPCServer(addr string, storage domain.Storage, log *logger.ZapLogger, config *handlers.Config) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	s := grpc.NewServer()
	RegisterShortenerServer(s, &ShortenerGRPCServer{
		Storage: storage,
		Log:     log,
		Config:  config,
	})
	log.InfoCtx(context.Background(), "gRPC server listening", zap.String("addr", addr))
	return s.Serve(lis)
}

func (s *ShortenerGRPCServer) CreateUser(ctx context.Context, req *CreateUserRequest) (*CreateUserResponse, error) {
	if s.Storage.User == nil {
		return nil, status.Error(codes.Unimplemented, "user storage not available")
	}
	userID, err := s.Storage.User.Create()
	if err != nil {
		s.Log.Log.Error("gRPC: Error creating user", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to create user")
	}
	return &CreateUserResponse{UserId: int32(userID)}, nil
}
