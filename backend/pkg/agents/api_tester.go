package agents

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// APITesterAgent specializes in testing power industry API endpoints
type APITesterAgent struct {
	BaseAgent
	httpClient *http.Client
	testSuites map[string]APITestSuite
}

// APITestSuite contains test cases for specific API categories
type APITestSuite struct {
	Name        string
	Category    string
	TestCases   []APITestCase
	RiskLevel   string
	Description string
}

// APITestCase represents a single API security test
type APITestCase struct {
	Name         string
	Method       string
	Endpoint     string
	Headers      map[string]string
	Payload      string
	ExpectedCode int
	TestType     string
	RiskLevel    string
	Description  string
}

// APITestResult contains the results of API testing
type APITestResult struct {
	TestCase       APITestCase
	ActualCode     int
	ResponseBody   string
	ResponseTime   time.Duration
	Vulnerability  *APIVulnerability
	Status         string
	Error          string
}

// APIVulnerability represents a discovered API vulnerability
type APIVulnerability struct {
	Type        string
	Severity    string
	Description string
	Impact      string
	Evidence    string
	CVSS        float64
	CWE         string
}

// NewAPITesterAgent creates a new API testing agent
func NewAPITesterAgent() *APITesterAgent {
	return &APITesterAgent{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		testSuites: initializeAPITestSuites(),
	}
}

// Execute performs comprehensive API security testing
func (a *APITesterAgent) Execute(ctx context.Context, input AgentInput) (*AgentOutput, error) {
	var params struct {
		BaseURL     string            `json:"base_url"`
		SystemType  string            `json:"system_type"`
		AuthToken   string            `json:"auth_token"`
		Headers     map[string]string `json:"headers"`
		TestSuites  []string          `json:"test_suites"`
	}

	if err := json.Unmarshal([]byte(input.Target), &params); err != nil {
		return nil, fmt.Errorf("invalid input parameters: %w", err)
	}

	// Select appropriate test suites based on system type
	selectedSuites := a.selectTestSuites(params.SystemType, params.TestSuites)
	
	var allResults []APITestResult
	var vulnerabilities []APIVulnerability

	// Execute each test suite
	for _, suite := range selectedSuites {
		suiteResults, err := a.executeTestSuite(ctx, suite, params.BaseURL, params.AuthToken, params.Headers)
		if err != nil {
			return nil, fmt.Errorf("failed to execute test suite %s: %w", suite.Name, err)
		}

		allResults = append(allResults, suiteResults...)

		// Collect vulnerabilities
		for _, result := range suiteResults {
			if result.Vulnerability != nil {
				vulnerabilities = append(vulnerabilities, *result.Vulnerability)
			}
		}
	}

	// Generate comprehensive report
	report := a.generateAPIReport(params.SystemType, allResults, vulnerabilities)

	return &AgentOutput{
		Result: report,
		Metadata: map[string]interface{}{
			"system_type":        params.SystemType,
			"tests_executed":     len(allResults),
			"vulnerabilities":    len(vulnerabilities),
			"critical_vulns":     a.countBySeverity(vulnerabilities, "critical"),
			"high_vulns":         a.countBySeverity(vulnerabilities, "high"),
		},
	}, nil
}

// selectTestSuites chooses appropriate test suites based on system type
func (a *APITesterAgent) selectTestSuites(systemType string, requestedSuites []string) []APITestSuite {
	var selected []APITestSuite

	// If specific suites requested, use those
	if len(requestedSuites) > 0 {
		for _, suiteName := range requestedSuites {
			if suite, exists := a.testSuites[suiteName]; exists {
				selected = append(selected, suite)
			}
		}
		return selected
	}

	// Otherwise, select based on system type
	switch systemType {
	case "marketing_2.0":
		selected = append(selected, 
			a.testSuites["authentication"],
			a.testSuites["authorization"], 
			a.testSuites["billing_api"],
			a.testSuites["customer_data"],
		)
	case "iguowang":
		selected = append(selected,
			a.testSuites["mobile_auth"],
			a.testSuites["sms_verification"],
			a.testSuites["mobile_api"],
			a.testSuites["payment_api"],
		)
	case "sap":
		selected = append(selected,
			a.testSuites["sap_auth"],
			a.testSuites["enterprise_api"],
			a.testSuites["admin_functions"],
		)
	default:
		// Use all available test suites
		for _, suite := range a.testSuites {
			selected = append(selected, suite)
		}
	}

	return selected
}

// executeTestSuite runs all test cases in a test suite
func (a *APITesterAgent) executeTestSuite(ctx context.Context, suite APITestSuite, baseURL, authToken string, headers map[string]string) ([]APITestResult, error) {
	var results []APITestResult

	for _, testCase := range suite.TestCases {
		result, err := a.executeTestCase(ctx, testCase, baseURL, authToken, headers)
		if err != nil {
			result = APITestResult{
				TestCase: testCase,
				Status:   "error",
				Error:    err.Error(),
			}
		}
		results = append(results, result)
	}

	return results, nil
}

// executeTestCase runs a single API test case
func (a *APITesterAgent) executeTestCase(ctx context.Context, testCase APITestCase, baseURL, authToken string, headers map[string]string) (APITestResult, error) {
	// Construct full URL
	url := strings.TrimSuffix(baseURL, "/") + "/" + strings.TrimPrefix(testCase.Endpoint, "/")

	// Prepare request
	req, err := http.NewRequestWithContext(ctx, testCase.Method, url, strings.NewReader(testCase.Payload))
	if err != nil {
		return APITestResult{}, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	for key, value := range testCase.Headers {
		req.Header.Set(key, value)
	}

	// Set authentication if provided
	if authToken != "" {
		req.Header.Set("Authorization", "Bearer "+authToken)
	}

	// Set content type for POST/PUT requests
	if testCase.Method == "POST" || testCase.Method == "PUT" {
		req.Header.Set("Content-Type", "application/json")
	}

	// Execute request
	startTime := time.Now()
	resp, err := a.httpClient.Do(req)
	responseTime := time.Since(startTime)

	if err != nil {
		return APITestResult{}, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body := make([]byte, 4096) // Limit response size
	n, _ := resp.Body.Read(body)
	responseBody := string(body[:n])

	// Create result
	result := APITestResult{
		TestCase:     testCase,
		ActualCode:   resp.StatusCode,
		ResponseBody: responseBody,
		ResponseTime: responseTime,
		Status:       "completed",
	}

	// Analyze for vulnerabilities
	vulnerability := a.analyzeAPIResponse(testCase, resp.StatusCode, responseBody)
	if vulnerability != nil {
		result.Vulnerability = vulnerability
	}

	return result, nil
}

// analyzeAPIResponse analyzes API response for security vulnerabilities
func (a *APITesterAgent) analyzeAPIResponse(testCase APITestCase, statusCode int, responseBody string) *APIVulnerability {
	switch testCase.TestType {
	case "sql_injection":
		if strings.Contains(responseBody, "SQL") || strings.Contains(responseBody, "syntax error") || 
		   strings.Contains(responseBody, "mysql") || strings.Contains(responseBody, "postgres") {
			return &APIVulnerability{
				Type:        "SQL Injection",
				Severity:    "critical",
				Description: "SQL injection vulnerability detected in API endpoint",
				Impact:      "Complete database compromise, data theft, data manipulation",
				Evidence:    responseBody,
				CVSS:        9.8,
				CWE:         "CWE-89",
			}
		}

	case "auth_bypass":
		if statusCode == 200 && strings.Contains(testCase.Payload, "OR 1=1") {
			return &APIVulnerability{
				Type:        "Authentication Bypass",
				Severity:    "critical",
				Description: "Authentication can be bypassed using SQL injection",
				Impact:      "Unauthorized access to user accounts and sensitive data",
				Evidence:    fmt.Sprintf("Status: %d, Payload: %s", statusCode, testCase.Payload),
				CVSS:        9.1,
				CWE:         "CWE-287",
			}
		}

	case "idor":
		if statusCode == 200 && strings.Contains(responseBody, "id") {
			return &APIVulnerability{
				Type:        "Insecure Direct Object Reference",
				Severity:    "high",
				Description: "Direct object references allow unauthorized data access",
				Impact:      "Access to other users' sensitive information",
				Evidence:    fmt.Sprintf("Endpoint: %s, Response contains user data", testCase.Endpoint),
				CVSS:        7.5,
				CWE:         "CWE-639",
			}
		}

	case "privilege_escalation":
		if statusCode == 200 && (strings.Contains(responseBody, "admin") || strings.Contains(responseBody, "elevated")) {
			return &APIVulnerability{
				Type:        "Privilege Escalation",
				Severity:    "critical",
				Description: "API allows unauthorized privilege escalation",
				Impact:      "Administrative access to system functions",
				Evidence:    responseBody,
				CVSS:        8.8,
				CWE:         "CWE-269",
			}
		}

	case "information_disclosure":
		if statusCode == 200 && (strings.Contains(responseBody, "password") || 
		   strings.Contains(responseBody, "secret") || strings.Contains(responseBody, "key")) {
			return &APIVulnerability{
				Type:        "Information Disclosure",
				Severity:    "high",
				Description: "Sensitive information exposed in API response",
				Impact:      "Exposure of credentials and sensitive configuration",
				Evidence:    responseBody,
				CVSS:        7.5,
				CWE:         "CWE-200",
			}
		}

	case "business_logic":
		if statusCode == 200 && strings.Contains(testCase.Payload, "-") && strings.Contains(responseBody, "amount") {
			return &APIVulnerability{
				Type:        "Business Logic Flaw",
				Severity:    "high",
				Description: "Business logic allows invalid operations",
				Impact:      "Financial loss through business logic manipulation",
				Evidence:    fmt.Sprintf("Negative amount processed: %s", responseBody),
				CVSS:        7.2,
				CWE:         "CWE-840",
			}
		}
	}

	return nil
}

// generateAPIReport creates a comprehensive API security report
func (a *APITesterAgent) generateAPIReport(systemType string, results []APITestResult, vulnerabilities []APIVulnerability) string {
	report := fmt.Sprintf("# API Security Assessment Report\n\n")
	report += fmt.Sprintf("## System Type: %s\n\n", systemType)

	// Executive Summary
	report += "## Executive Summary\n\n"
	report += fmt.Sprintf("This report presents the results of comprehensive API security testing performed on the %s system. ", systemType)
	report += fmt.Sprintf("A total of %d API endpoints were tested using %d test cases.\n\n", a.countUniqueEndpoints(results), len(results))

	// Risk Summary
	criticalCount := a.countBySeverity(vulnerabilities, "critical")
	highCount := a.countBySeverity(vulnerabilities, "high")
	mediumCount := a.countBySeverity(vulnerabilities, "medium")

	report += "## Risk Summary\n\n"
	report += fmt.Sprintf("- **Critical Vulnerabilities**: %d\n", criticalCount)
	report += fmt.Sprintf("- **High Risk Vulnerabilities**: %d\n", highCount)
	report += fmt.Sprintf("- **Medium Risk Vulnerabilities**: %d\n", mediumCount)
	report += fmt.Sprintf("- **Total Vulnerabilities**: %d\n\n", len(vulnerabilities))

	// Vulnerability Details
	if len(vulnerabilities) > 0 {
		report += "## Vulnerability Details\n\n"
		for i, vuln := range vulnerabilities {
			report += fmt.Sprintf("### Vulnerability %d: %s\n\n", i+1, vuln.Type)
			report += fmt.Sprintf("**Severity**: %s (CVSS: %.1f)\n", vuln.Severity, vuln.CVSS)
			report += fmt.Sprintf("**CWE**: %s\n", vuln.CWE)
			report += fmt.Sprintf("**Description**: %s\n", vuln.Description)
			report += fmt.Sprintf("**Impact**: %s\n", vuln.Impact)
			report += fmt.Sprintf("**Evidence**: %s\n\n", vuln.Evidence)
		}
	}

	// Test Results Summary
	report += "## Test Results Summary\n\n"
	successCount := a.countByStatus(results, "completed")
	errorCount := a.countByStatus(results, "error")

	report += fmt.Sprintf("- **Tests Completed**: %d\n", successCount)
	report += fmt.Sprintf("- **Tests Failed**: %d\n", errorCount)
	report += fmt.Sprintf("- **Success Rate**: %.1f%%\n\n", float64(successCount)/float64(len(results))*100)

	// Recommendations
	report += "## Recommendations\n\n"
	report += "### Immediate Actions (Critical Priority)\n"
	if criticalCount > 0 {
		report += "1. **Address Critical Vulnerabilities**: Immediately patch all critical security flaws\n"
		report += "2. **Disable Vulnerable Endpoints**: Temporarily disable endpoints with critical vulnerabilities\n"
		report += "3. **Implement Emergency Monitoring**: Deploy additional monitoring for affected systems\n\n"
	}

	report += "### Short-term Actions (High Priority)\n"
	report += "1. **Input Validation**: Implement comprehensive input validation for all API endpoints\n"
	report += "2. **Authentication Strengthening**: Enhance authentication mechanisms and session management\n"
	report += "3. **Authorization Review**: Conduct thorough review of authorization controls\n"
	report += "4. **Security Testing**: Integrate automated security testing into CI/CD pipeline\n\n"

	report += "### Long-term Actions (Strategic)\n"
	report += "1. **Security Architecture Review**: Conduct comprehensive security architecture assessment\n"
	report += "2. **Developer Training**: Implement secure coding training for development teams\n"
	report += "3. **Security Standards**: Establish and enforce API security standards\n"
	report += "4. **Continuous Monitoring**: Deploy comprehensive security monitoring and alerting\n\n"

	return report
}

// Helper functions
func (a *APITesterAgent) countUniqueEndpoints(results []APITestResult) int {
	endpoints := make(map[string]bool)
	for _, result := range results {
		endpoints[result.TestCase.Endpoint] = true
	}
	return len(endpoints)
}

func (a *APITesterAgent) countBySeverity(vulnerabilities []APIVulnerability, severity string) int {
	count := 0
	for _, vuln := range vulnerabilities {
		if vuln.Severity == severity {
			count++
		}
	}
	return count
}

func (a *APITesterAgent) countByStatus(results []APITestResult, status string) int {
	count := 0
	for _, result := range results {
		if result.Status == status {
			count++
		}
	}
	return count
}

// initializeAPITestSuites creates predefined test suites for different categories
func initializeAPITestSuites() map[string]APITestSuite {
	return map[string]APITestSuite{
		"authentication": {
			Name:        "Authentication Security",
			Category:    "auth",
			RiskLevel:   "critical",
			Description: "Tests authentication mechanisms for bypass and weaknesses",
			TestCases: []APITestCase{
				{
					Name:         "SQL Injection in Login",
					Method:       "POST",
					Endpoint:     "/api/auth/login",
					Payload:      `{"username": "admin' OR 1=1--", "password": "any"}`,
					ExpectedCode: 401,
					TestType:     "sql_injection",
					RiskLevel:    "critical",
					Description:  "Tests for SQL injection in authentication",
				},
				{
					Name:         "Authentication Bypass",
					Method:       "POST", 
					Endpoint:     "/api/auth/login",
					Payload:      `{"username": "admin", "password": "' OR '1'='1"}`,
					ExpectedCode: 401,
					TestType:     "auth_bypass",
					RiskLevel:    "critical",
					Description:  "Tests for authentication bypass vulnerabilities",
				},
			},
		},
		"authorization": {
			Name:        "Authorization Controls",
			Category:    "authz",
			RiskLevel:   "high",
			Description: "Tests authorization and access control mechanisms",
			TestCases: []APITestCase{
				{
					Name:         "IDOR User Access",
					Method:       "GET",
					Endpoint:     "/api/users/1",
					ExpectedCode: 403,
					TestType:     "idor",
					RiskLevel:    "high",
					Description:  "Tests for insecure direct object references",
				},
				{
					Name:         "Privilege Escalation",
					Method:       "PUT",
					Endpoint:     "/api/user/profile",
					Payload:      `{"user_id": "admin", "role": "administrator"}`,
					ExpectedCode: 403,
					TestType:     "privilege_escalation",
					RiskLevel:    "critical",
					Description:  "Tests for privilege escalation vulnerabilities",
				},
			},
		},
		"billing_api": {
			Name:        "Billing API Security",
			Category:    "business_logic",
			RiskLevel:   "critical",
			Description: "Tests billing and payment API endpoints",
			TestCases: []APITestCase{
				{
					Name:         "Negative Amount Processing",
					Method:       "POST",
					Endpoint:     "/api/billing/calculate",
					Payload:      `{"usage": -999999, "rate": "standard"}`,
					ExpectedCode: 400,
					TestType:     "business_logic",
					RiskLevel:    "high",
					Description:  "Tests business logic for negative billing amounts",
				},
			},
		},
		"customer_data": {
			Name:        "Customer Data Protection",
			Category:    "data_protection",
			RiskLevel:   "high",
			Description: "Tests customer data access and protection",
			TestCases: []APITestCase{
				{
					Name:         "Customer Data Enumeration",
					Method:       "GET",
					Endpoint:     "/api/customers",
					ExpectedCode: 401,
					TestType:     "information_disclosure",
					RiskLevel:    "high",
					Description:  "Tests for unauthorized customer data access",
				},
			},
		},
		"mobile_auth": {
			Name:        "Mobile Authentication",
			Category:    "mobile",
			RiskLevel:   "high",
			Description: "Tests mobile-specific authentication mechanisms",
			TestCases: []APITestCase{
				{
					Name:         "SMS Verification Bypass",
					Method:       "POST",
					Endpoint:     "/api/auth/verify-sms",
					Payload:      `{"phone": "13800138000", "code": "000000"}`,
					ExpectedCode: 401,
					TestType:     "auth_bypass",
					RiskLevel:    "high",
					Description:  "Tests SMS verification bypass",
				},
			},
		},
		"sms_verification": {
			Name:        "SMS Verification Security",
			Category:    "mobile",
			RiskLevel:   "medium",
			Description: "Tests SMS verification implementation",
			TestCases: []APITestCase{
				{
					Name:         "SMS Code Enumeration",
					Method:       "POST",
					Endpoint:     "/api/auth/send-sms",
					Payload:      `{"phone": "13800138000"}`,
					ExpectedCode: 200,
					TestType:     "information_disclosure",
					RiskLevel:    "medium",
					Description:  "Tests for SMS code disclosure",
				},
			},
		},
		"mobile_api": {
			Name:        "Mobile API Security",
			Category:    "mobile",
			RiskLevel:   "high",
			Description: "Tests mobile application API endpoints",
			TestCases: []APITestCase{
				{
					Name:         "User List Exposure",
					Method:       "GET",
					Endpoint:     "/api/user/list",
					ExpectedCode: 401,
					TestType:     "information_disclosure",
					RiskLevel:    "high",
					Description:  "Tests for user information exposure",
				},
			},
		},
		"payment_api": {
			Name:        "Payment API Security",
			Category:    "financial",
			RiskLevel:   "critical",
			Description: "Tests payment processing API security",
			TestCases: []APITestCase{
				{
					Name:         "Payment Amount Manipulation",
					Method:       "POST",
					Endpoint:     "/api/payment/create",
					Payload:      `{"amount": 0.01, "bill_id": "12345"}`,
					ExpectedCode: 400,
					TestType:     "business_logic",
					RiskLevel:    "critical",
					Description:  "Tests payment amount manipulation",
				},
			},
		},
		"sap_auth": {
			Name:        "SAP Authentication",
			Category:    "enterprise",
			RiskLevel:   "critical",
			Description: "Tests SAP system authentication",
			TestCases: []APITestCase{
				{
					Name:         "SAP Default Credentials",
					Method:       "POST",
					Endpoint:     "/api/auth/login",
					Payload:      `{"username": "SAP*", "password": "06071992", "client": "000"}`,
					ExpectedCode: 401,
					TestType:     "auth_bypass",
					RiskLevel:    "critical",
					Description:  "Tests for SAP default credentials",
				},
			},
		},
		"enterprise_api": {
			Name:        "Enterprise API Security",
			Category:    "enterprise",
			RiskLevel:   "high",
			Description: "Tests enterprise system API endpoints",
			TestCases: []APITestCase{
				{
					Name:         "System Configuration Exposure",
					Method:       "GET",
					Endpoint:     "/api/system/config",
					ExpectedCode: 401,
					TestType:     "information_disclosure",
					RiskLevel:    "high",
					Description:  "Tests for system configuration exposure",
				},
			},
		},
		"admin_functions": {
			Name:        "Administrative Functions",
			Category:    "admin",
			RiskLevel:   "critical",
			Description: "Tests administrative function security",
			TestCases: []APITestCase{
				{
					Name:         "Command Execution",
					Method:       "POST",
					Endpoint:     "/api/admin/execute",
					Payload:      `{"command": "whoami"}`,
					ExpectedCode: 401,
					TestType:     "privilege_escalation",
					RiskLevel:    "critical",
					Description:  "Tests for command execution vulnerabilities",
				},
			},
		},
	}
}
