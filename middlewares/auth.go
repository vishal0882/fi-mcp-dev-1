package middlewares

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type AuthMiddleware struct {
	sessionStore map[string]string
}

func NewAuthMiddleware() *AuthMiddleware {
	return &AuthMiddleware{
		sessionStore: make(map[string]string),
	}
}

func (m *AuthMiddleware) AuthMiddleware(next server.ToolHandlerFunc) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// fetch sessionId from context
		// this gets populated for every tool call
		sessionId := server.ClientSessionFromContext(ctx).SessionID()
		phoneNumber, ok := m.sessionStore[sessionId]
		if !ok {
			loginUrl, getLoginUrlErr := m.getLoginUrl(ctx, sessionId)
			if getLoginUrlErr != nil {
				return mcp.NewToolResultText("something went wrong"), nil
			}
			res := fmt.Sprintf("Please login by clicking this link: [Login](%s)\n\nIf your client supports clickable links, you can render and present it and ask them to click the link above. Otherwise, display the URL and ask them to copy and paste it into their browser: %s\n\nAfter completing the login in your browser, let me know and I'll continue with your request.", loginUrl, loginUrl)
			return mcp.NewToolResultText(res), nil
		}
		ctx = context.WithValue(ctx, "phone_number", phoneNumber)
		return next(ctx, req)
	}
}

// GetLoginUrl fetches dynamic login url for given sessionId
func (m *AuthMiddleware) getLoginUrl(ctx context.Context, sessionId string) (string, error) {
	return "http://localhost:8080/mockWebPage?sessionId=" + sessionId, nil
}

func (m *AuthMiddleware) AddSession(sessionId, phoneNumber string) {
	m.sessionStore[sessionId] = phoneNumber
}
