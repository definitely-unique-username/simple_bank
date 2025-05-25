package gapi

import (
	"context"
	"database/sql"

	db "github.com/definitely-unique-username/simple_bank/db/sqlc"
	"github.com/definitely-unique-username/simple_bank/pb"
	"github.com/definitely-unique-username/simple_bank/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	user, err := s.store.GetUser(ctx, req.GetUsername())

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "user not found: %s", err)
		}

		return nil, status.Errorf(codes.Internal, "filaed to find user: %s", err)
	}

	err = util.VerifyPassword(req.GetPassword(), user.Hash)

	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "email or password is invalid: %s", err)
	}

	accessToken, accessTokenPayload, err := s.tokenMaker.CreateToken(user.ID, s.config.AccessTokenDuration)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "faield to create access token: %s", err)
	}

	refreshToken, refreshTokenPayload, err := s.tokenMaker.CreateToken(user.ID, s.config.RefreshTokenDuratuin)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "faield to create refresh token: %s", err)
	}

	metadata := s.extractMetadata(ctx)

	session, err := s.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshTokenPayload.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    metadata.UserAgent,
		ClientIp:     metadata.ClintIP,
		IsBlocked:    false,
		ExpiresAt:    refreshTokenPayload.ExpiredAt,
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal, "faield to create session: %s", err)
	}

	return &pb.LoginUserResponse{
		SessioId:              session.ID.String(),
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  timestamppb.New(accessTokenPayload.ExpiredAt),
		User:                  convertUser(user),
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: timestamppb.New(refreshTokenPayload.ExpiredAt),
	}, nil
}
