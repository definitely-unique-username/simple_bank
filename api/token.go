package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type renewAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type renewAccessTokenResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

func (s *Server) renewAccessToken(ctx *gin.Context) {
	var req renewAccessTokenRequest

	if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
		fmt.Printf("Error binding request: %v\n", err)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	fmt.Printf("Received refresh token: %s\n", req.RefreshToken)

	refreshTokenPayload, err := s.tokenMaker.VerifyToken(req.RefreshToken)
	if err != nil {
		fmt.Printf("Error verifying token: %v\n", err)
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	fmt.Printf("Token payload ID: %v\n", refreshTokenPayload.ID)

	session, err := s.store.GetSession(ctx, refreshTokenPayload.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Printf("Session not found for ID: %v\n", refreshTokenPayload.ID)
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		fmt.Printf("Error getting session: %v\n", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	fmt.Printf("Found session: %+v\n", session)

	if session.IsBlocked {
		fmt.Printf("Session is blocked: %v\n", session.ID)
		ctx.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("blocked session")))
		return
	}

	// TODO: check if sessions user is the same as in refresh token payload

	if session.RefreshToken != req.RefreshToken {
		fmt.Printf("Token mismatch - Session token: %s, Request token: %s\n",
			session.RefreshToken, req.RefreshToken)
		ctx.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("mismatch session token")))
		return
	}

	if time.Now().After(session.ExpiresAt) {
		fmt.Printf("Session expired at: %v\n", session.ExpiresAt)
		ctx.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("expired session")))
		return
	}

	accessToken, accessTokenPayload, err := s.tokenMaker.CreateToken(refreshTokenPayload.UserID, s.config.AccessTokenDuration)
	if err != nil {
		fmt.Printf("Error creating access token: %v\n", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	fmt.Printf("Successfully created new access token, expires at: %v\n", accessTokenPayload.ExpiredAt)

	ctx.JSON(http.StatusOK, &renewAccessTokenResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessTokenPayload.ExpiredAt,
	})
}
