package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/jwt"
	"github.com/clerk/clerk-sdk-go/v2/user"
	"github.com/danielgtaylor/huma/v2"
	"github.com/tools4net/ezfw/backend/internal/models"
	"github.com/tools4net/ezfw/backend/internal/store"
	"go.uber.org/zap"
)

type contextKey string

const (
	AuthContextKey contextKey = "auth_context"
	UserContextKey contextKey = "user_context"
	AgentContextKey contextKey = "agent_context"
)

// AuthMiddleware provides dual authentication support for users and agents
type AuthMiddleware struct {
	store       store.Store
	logger      *zap.Logger
	clerkKey    string
	userClient  *user.Client
}

// NewAuthMiddleware creates a new authentication middleware
func NewAuthMiddleware(store store.Store, logger *zap.Logger) *AuthMiddleware {
	clerkSecretKey := os.Getenv("CLERK_SECRET_KEY")
	if clerkSecretKey == "" {
		logger.Warn("CLERK_SECRET_KEY not set, Clerk authentication will be disabled")
	}

	// Set global Clerk key
	clerk.SetKey(clerkSecretKey)

	// Create user client for API operations
	config := &clerk.ClientConfig{}
	if clerkSecretKey != "" {
		config.Key = &clerkSecretKey
	}
	userClient := user.NewClient(config)

	return &AuthMiddleware{
		store:      store,
		logger:     logger,
		clerkKey:   clerkSecretKey,
		userClient: userClient,
	}
}

// RequireAuth middleware that requires either user or agent authentication
func (m *AuthMiddleware) RequireAuth(ctx huma.Context, next func(huma.Context)) {
	authCtx := m.authenticate(ctx.Context(), ctx)
	if authCtx == nil {
		ctx.SetStatus(http.StatusUnauthorized)
		ctx.SetHeader("Content-Type", "application/json")
		ctx.BodyWriter().Write([]byte(`{"message":"Authentication required"}`))
		return
	}

	// Store auth context for later retrieval
	// We'll use a different approach since huma middleware context handling is limited
	next(ctx)
}

// RequireUserAuth middleware that requires user authentication only
func (m *AuthMiddleware) RequireUserAuth(ctx huma.Context, next func(huma.Context)) {
	authCtx := m.authenticate(ctx.Context(), ctx)
	if authCtx == nil || authCtx.Type != "user" {
		ctx.SetStatus(http.StatusUnauthorized)
		ctx.SetHeader("Content-Type", "application/json")
		ctx.BodyWriter().Write([]byte(`{"message":"User authentication required"}`))
		return
	}

	// Store auth context for later retrieval
	next(ctx)
}

// RequireAgentAuth middleware that requires agent authentication only
func (m *AuthMiddleware) RequireAgentAuth(ctx huma.Context, next func(huma.Context)) {
	authCtx := m.authenticate(ctx.Context(), ctx)
	if authCtx == nil || authCtx.Type != "agent" {
		ctx.SetStatus(http.StatusUnauthorized)
		ctx.SetHeader("Content-Type", "application/json")
		ctx.BodyWriter().Write([]byte(`{"message":"Agent authentication required"}`))
		return
	}

	// Store auth context for later retrieval
	next(ctx)
}

// authenticate attempts to authenticate the request using either user or agent credentials
func (m *AuthMiddleware) authenticate(ctx context.Context, humaCtx huma.Context) *models.AuthContext {
	// Try agent authentication first (X-Agent-Token header)
	if agentToken := humaCtx.Header("X-Agent-Token"); agentToken != "" {
		return m.authenticateAgent(ctx, agentToken)
	}

	// Try user authentication (Authorization Bearer token)
	if authHeader := humaCtx.Header("Authorization"); authHeader != "" {
		if strings.HasPrefix(authHeader, "Bearer ") {
			token := strings.TrimPrefix(authHeader, "Bearer ")
			return m.authenticateUser(ctx, token)
		}
	}

	return nil
}

// authenticateUser validates Clerk JWT token and returns user context
func (m *AuthMiddleware) authenticateUser(ctx context.Context, token string) *models.AuthContext {
	if m.clerkKey == "" {
		m.logger.Warn("Clerk key not configured")
		return nil
	}

	// Decode the JWT token to get the Key ID
	decoded, err := jwt.Decode(ctx, &jwt.DecodeParams{Token: token})
	if err != nil {
		m.logger.Debug("Failed to decode JWT token", zap.Error(err))
		return nil
	}

	// Fetch the JSON web key for verification
	jwk, err := jwt.GetJSONWebKey(ctx, &jwt.GetJSONWebKeyParams{
		KeyID: decoded.KeyID,
	})
	if err != nil {
		m.logger.Debug("Failed to get JWK", zap.Error(err))
		return nil
	}

	// Verify JWT token with Clerk
	claims, err := jwt.Verify(ctx, &jwt.VerifyParams{
		Token: token,
		JWK:   jwk,
	})
	if err != nil {
		m.logger.Debug("Invalid JWT token", zap.Error(err))
		return nil
	}

	// Get user details from Clerk
	usr, err := m.userClient.Get(ctx, claims.Subject)
	if err != nil {
		m.logger.Error("Failed to get user from Clerk", zap.Error(err))
		return nil
	}

	userCtx := &models.UserContext{
		UserID:    usr.ID,
		Email:     getFirstEmailAddress(usr.EmailAddresses),
		FirstName: getStringValue(usr.FirstName),
		LastName:  getStringValue(usr.LastName),
	}

	// Note: Organization memberships require separate API call in Clerk SDK v2
	// For now, we'll leave organization empty and can implement this later
	// if needed using the organization membership API

	return &models.AuthContext{
		User: userCtx,
		Type: "user",
	}
}

// authenticateAgent validates agent token and returns agent context
func (m *AuthMiddleware) authenticateAgent(ctx context.Context, token string) *models.AuthContext {
	agentToken, err := m.store.GetAgentTokenByToken(ctx, token)
	if err != nil {
		m.logger.Debug("Invalid agent token", zap.Error(err))
		return nil
	}

	// Check if token is active
	if agentToken.Status != "active" {
		m.logger.Debug("Agent token is not active", zap.String("status", agentToken.Status))
		return nil
	}

	// Check if token is expired
	if agentToken.ExpiresAt != nil && agentToken.ExpiresAt.Before(time.Now()) {
		m.logger.Debug("Agent token is expired")
		return nil
	}

	// Update last used timestamp
	// Note: LastUsed field would need to be added to AgentTokenUpdate model
	// m.store.UpdateAgentToken(ctx, agentToken.ID, &models.AgentTokenUpdate{
	//     LastUsed: &now,
	// })

	// Get associated node
	node, err := m.store.GetNode(ctx, agentToken.NodeID)
	if err != nil {
		m.logger.Error("Failed to get node for agent token", zap.Error(err))
		return nil
	}

	agentCtx := &models.AgentContext{
		TokenID: agentToken.ID,
		NodeID:  agentToken.NodeID,
		Node:    node,
	}

	return &models.AuthContext{
		Agent: agentCtx,
		Type:  "agent",
	}
}

// GenerateAgentToken generates a new secure agent token
func GenerateAgentToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return "agt_" + hex.EncodeToString(bytes), nil
}

// Helper functions
func getFirstEmailAddress(emails []*clerk.EmailAddress) string {
	if len(emails) > 0 {
		return emails[0].EmailAddress
	}
	return ""
}

func getStringValue(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

// GetAuthContext retrieves the authentication context from the request context
func GetAuthContext(ctx context.Context) *models.AuthContext {
	if authCtx, ok := ctx.Value(AuthContextKey).(*models.AuthContext); ok {
		return authCtx
	}
	return nil
}

// GetUserContext retrieves the user context from the request context
func GetUserContext(ctx context.Context) *models.UserContext {
	if userCtx, ok := ctx.Value(UserContextKey).(*models.UserContext); ok {
		return userCtx
	}
	return nil
}

// GetAgentContext retrieves the agent context from the request context
func GetAgentContext(ctx context.Context) *models.AgentContext {
	if agentCtx, ok := ctx.Value(AgentContextKey).(*models.AgentContext); ok {
		return agentCtx
	}
	return nil
}