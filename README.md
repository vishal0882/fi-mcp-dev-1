# fi-mcp-lite

A minimal, hackathon-ready version of the Fi MCP server. This project provides a lightweight mock server for use in hackathons, demos, and development, simulating the core features of the production Fi MCP server with dummy data and simplified authentication.

## Purpose

- **fi-mcp-lite** is designed for hackathon participants and developers who want to experiment with the Fi MCP API without accessing real user data or production systems.
- It serves dummy financial data and uses a dummy authentication flow, making it safe and easy to use in non-production environments.

## Features

- **Simulates Fi MCP API**: Implements endpoints for net worth, credit report, EPF details, mutual fund transactions, and bank transactions.
- **Dummy Data**: All responses are served from static JSON files in `test_data_dir/`, representing various user scenarios.
- **Dummy Authentication**: Simple login flow using allowed phone numbers (directory names in `test_data_dir/`). No real OTP or user verification.
- **Hackathon-Ready**: No real integrations, no sensitive data, and easy to reset or extend.

## Directory Structure

- `main.go` — Entrypoint, sets up the server and endpoints.
- `middlewares/auth.go` — Implements dummy authentication and session management.
- `test_data_dir/` — Contains directories named after allowed phone numbers. Each directory holds JSON files for different API responses (e.g., `fetch_net_worth.json`).
- `static/` — HTML files for the login and login-successful pages.

## Dummy Data Scenarios

The dummy data covers a variety of user states. Example scenarios:

- **All assets connected**: Banks, EPF, Indian stocks, US stocks, credit report, large or small mutual fund portfolios.
- **All assets except bank account**: No bank account, but other assets present.
- **Multiple banks and UANs**: Multiple bank accounts and EPF UANs, partial transaction coverage.
- **No assets connected**: Only a savings account balance is present.
- **No credit report**: All assets except credit report.

## Test Data Scenarios

| Phone Number | Description                                                                                                                                                                                                                                        |
|-------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| 1111111111  | No assets connected. Only saving account balance present                                                                                                                                                                                           |
| 2222222222  | All assets connected (Banks account, EPF, Indian stocks, US stocks, Credit report). Large mutual fund portfolio with 9 funds                                                                                                                       |
| 3333333333  | All assets connected (Banks account, EPF, Indian stocks, US stocks, Credit report). Small mutual fund portfolio with only 1 fund                                                                                                                   |
| 4444444444  | All assets connected (Banks account, EPF, Indian stocks, US stocks). Small mutual fund portfolio with only 1 fund. With 2 UAN account connected . With 3 different bank with multiple account in them . Only have transactions for 2 bank accounts |
| 5555555555  | All assets connected except credit score (Banks account, EPF, Indian stocks, US stocks). Small mutual fund portfolio with only 1 fund. With 3 different bank with multiple account in them. Only have transactions for 2 bank accounts             |
| 6666666666  | All assets connected except bank account (EPF, Indian stocks, US stocks). Large mutual fund portfolio with 9 funds. No bank account connected                                                                                                      |

## Example: Dummy Data File

A sample `fetch_net_worth.json` (truncated for brevity):

```json
{
  "netWorthResponse": {
    "assetValues": [
      {"netWorthAttribute": "ASSET_TYPE_MUTUAL_FUND", "value": {"currencyCode": "INR", "units": "84642"}},
      {"netWorthAttribute": "ASSET_TYPE_EPF", "value": {"currencyCode": "INR", "units": "211111"}}
    ],
    "liabilityValues": [
      {"netWorthAttribute": "LIABILITY_TYPE_VEHICLE_LOAN", "value": {"currencyCode": "INR", "units": "5000"}}
    ],
    "totalNetWorthValue": {"currencyCode": "INR", "units": "658305"}
  }
}
```

## Authentication Flow

- When a tool/API is called, the server checks for a valid session.
- If not authenticated, the user is prompted to log in via a web page (`/mockWebPage?sessionId=...`).
- Enter any allowed phone number (see directories in `test_data_dir/`). OTP is not validated.
- On successful login, the session is stored in memory for the duration of the server run.

## Running the Server

### Prerequisites
- Go 1.23 or later ([installation instructions](https://go.dev/doc/install))

### Install dependencies
```sh
go mod tidy
```

### Start the server
```sh
FI_MCP_PORT=8080 go run .
```

The server will start on [http://localhost:8080](http://localhost:8080).

## Usage
- Follow instructions in this [guide](https://fi.money/features/getting-started-with-fi-mcp) to setup client
- Replace url with locally running server, for example: `http://localhost:8080/mcp/stream`
- When prompted for login, use one of the above phone numbers
- Otp/Passcode can be anything on the webpage

