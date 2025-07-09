package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/epifi/fi-mcp-lite/middlewares"
	"github.com/epifi/fi-mcp-lite/pkg"
)

var authMiddleware *middlewares.AuthMiddleware

func main() {
	authMiddleware = middlewares.NewAuthMiddleware()
	s := server.NewMCPServer(
		"Hackathon MCP",
		"0.1.0",
		// Notifies clients when new tools gets added or any changes in tools
		server.WithInstructions("A financial portfolio management MCP server that provides secure access to users' financial data through Fi Money, a financial hub for all things money. This MCP server enables users to:\n- Access comprehensive net worth analysis with asset/liability breakdowns\n- Retrieve detailed transaction histories for mutual funds and Employee Provident Fund accounts\n- View credit reports with scores, loan details, and account histories, this also contains user's date of birth that can be used for calculating their age\n\nIf the person asks, you can tell about Fi Money that it is money management platform that offers below services in partnership with regulated entities:\n\nAVAILABLE SERVICES:\n- Digital savings account with zero Forex cards\n- Invest in Indian Mutual funds, US Stocks (partnership with licensed brokers), Smart and Fixed Deposits.\n- Instant Personal Loans \n- Faster UPI and Bank Transfers payments\n- Credit score monitoring and reports\n\nIMPORTANT LIMITATIONS:\n- This MCP server retrieves only actual user data via Net worth tracker and based on consent provided by the user  and does not generate hypothetical or estimated financial information\n- In this version of the MCP server, user's historical bank transactions, historical stocks transaction data, salary (unless categorically declared) is not present. Don't assume these data points for any kind of analysis.\n\nCRITICAL INSTRUCTIONS FOR FINANCIAL DATA:\n\n1. DATA BOUNDARIES: Only provide information that exists in the user's Fi Money Net worth tracker. Never estimate, extrapolate, or generate hypothetical financial data.\n\n2. SPENDING ANALYSIS: If user asks about spending patterns, categories, or analysis tell the user we currently don't offer that data through the MCP:\n   - For detailed spending insights, direct them to: \"For comprehensive spending analysis and categorization, please use the Fi Money mobile app which provides detailed spending insights and budgeting tools.\"\n\n3. MISSING DATA HANDLING: If requested data is not available:\n   - Clearly state what data is missing\n   - Explain how user can connect additional accounts in Fi Money app\n   - Never fill gaps with estimated or generic information\n"),
		server.WithToolCapabilities(true),
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
		server.WithToolHandlerMiddleware(authMiddleware.AuthMiddleware),
	)

	s.AddTool(mcp.NewTool("fetch_net_worth",
		mcp.WithDescription("Calculate comprehensive net worth using ONLY actual data from accounts users connected on Fi Money including:\n- Bank account balances\n- Mutual fund investment holdings\n- Indian Stocks investment holdings\n- Total US Stocks investment (If investing through Fi Money app)\n- EPF account balances\n- Credit card debt and loan balances (if credit report connected)\n- Any other assets/liabilities linked to Fi Money platform\n\nUse this tool when the user asks for:\n- Their current net worth along with their holdings\n- Analysis of their mutual funds, stocks holdings\n- Use their assets information to give more personalized information \n\nIMPORTANT: This tool returns ONLY declared account data. It does not estimate or calculate values for unconnected accounts. If user has assets/liabilities not connected to Fi Money, they will not be included in calculations.\n\nERROR HANDLING: \n- If mutual funds are not detected, please ask the user whether they hold any mutual funds and prompt them to connect their mutual fund accounts via the Fi app.\n- If no financial accounts are connected that the user specifically asked for, returns empty result with message to connect those accounts in Fi Money app."),
	), dummyHandler)

	s.AddTool(mcp.NewTool("fetch_credit_report",
		mcp.WithDescription("Retrieve comprehensive credit report including scores, active loans, credit card utilization, payment history, date of birth and recent inquiries from connected credit bureaus.\n\nUse this tool when the user asks for:\n- Their credit score related information\n- This tool can provide their age related information through date of birth\n- Which loans to close first based on the loan that has the highest interest rate\n\nIMPORTANT LIMITATIONS:\n- Only returns credit data if user has successfully connected their credit profile in Fi Money app\n- Cannot access credit data from bureaus not integrated with Fi Money platform\n\nERROR HANDLING: If no credit score data is found, respond with: \"No credit score data available. Please connect your credit score in the Fi Money app and try again.\""),
	), dummyHandler)

	s.AddTool(mcp.NewTool("fetch_epf_details",
		mcp.WithDescription("Retrieve detailed EPF (Employee Provident Fund) account information including:\n- Account balance and contributions\n- Employer and employee contribution history\n- Interest earned and credited amounts\n\nIMPORTANT LIMITATIONS:\n- Only returns EPF data if user has successfully linked their EPF account through Fi Money app\n- Data accuracy depends on EPFO integration and user's UAN verification status\n- Transaction history and account statements and passbooks are not available\n\nERROR HANDLING: If EPF account not connected, direct user to link EPF account in Fi Money app using their UAN and other required details.\n"),
	), dummyHandler)

	s.AddTool(mcp.NewTool("fetch_mf_transactions",
		mcp.WithDescription("Retrieve detailed transaction history from accounts connected to Fi Money platform including:\n- Mutual fund transactions\n\nUse this tool when the user asks for:\n- Their portfolio level XIRR (For now, it is only applicable to Mutual funds)\n\nIMPORTANT LIMITATIONS:\n- Current version does NOT include historical bank transactions, historical stocks transactions, NPS, deposits.\n\nERROR HANDLING: Returns only available transaction data with clear indication of data source limitations.\n"),
	), dummyHandler)

	s.AddTool(mcp.NewTool("fetch_bank_transactions",
		mcp.WithDescription("Retrieve detailed bank transactions for each bank account connected to Fi money platform"),
	), dummyHandler)

	// Configure streamable HTTP server with proper endpoints
	httpMux := http.NewServeMux()
	streamableServer := server.NewStreamableHTTPServer(s,
		server.WithEndpointPath("/stream"),
	)
	httpMux.Handle("/mcp/", streamableServer)
	httpMux.HandleFunc("/mockWebPage", webPageHandler)
	httpMux.HandleFunc("/login", loginHandler)
	port := pkg.GetPort()
	log.Println("starting server on port:", port)
	if servErr := http.ListenAndServe(fmt.Sprintf(":%s", port), httpMux); servErr != nil {
		log.Fatalln("error starting server", servErr)
	}
}

func dummyHandler(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return mcp.NewToolResultText("dummy handler"), nil
}

func webPageHandler(w http.ResponseWriter, r *http.Request) {
	sessionId := r.URL.Query().Get("sessionId")
	if sessionId == "" {
		http.Error(w, "sessionId is required", http.StatusBadRequest)
		return
	}

	tmpl, err := template.ParseFiles("static/login.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		SessionId            string
		AllowedMobileNumbers []string
	}{
		SessionId:            sessionId,
		AllowedMobileNumbers: pkg.GetAllowedMobileNumbers(),
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	sessionId := r.FormValue("sessionId")
	phoneNumber := r.FormValue("phoneNumber")

	if sessionId == "" || phoneNumber == "" {
		http.Error(w, "sessionId and phoneNumber are required", http.StatusBadRequest)
		return
	}

	authMiddleware.AddSession(sessionId, phoneNumber)

	tmpl, err := template.ParseFiles("static/login_successful.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
