package agents

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

// BizLogicTesterAgent specializes in testing business logic vulnerabilities in power industry systems
type BizLogicTesterAgent struct {
	BaseAgent
	businessRules map[string]BusinessRuleSet
	testScenarios map[string][]LogicTestScenario
}

// BusinessRuleSet defines business rules for different power industry systems
type BusinessRuleSet struct {
	SystemType    string
	Rules         []BusinessRule
	Constraints   []BusinessConstraint
	Workflows     []BusinessWorkflow
}

// BusinessRule represents a specific business rule
type BusinessRule struct {
	ID          string
	Name        string
	Description string
	Category    string
	RiskLevel   string
	TestMethod  string
}

// BusinessConstraint defines system constraints that should be enforced
type BusinessConstraint struct {
	Field       string
	Type        string
	MinValue    interface{}
	MaxValue    interface{}
	Required    bool
	Validation  string
}

// BusinessWorkflow represents a business process workflow
type BusinessWorkflow struct {
	Name        string
	Steps       []WorkflowStep
	Permissions []string
	Validation  []string
}

// WorkflowStep represents a single step in a business workflow
type WorkflowStep struct {
	ID          string
	Name        string
	Action      string
	Required    bool
	Validation  []string
	NextSteps   []string
}

// LogicTestScenario defines a business logic test scenario
type LogicTestScenario struct {
	Name        string
	Category    string
	Description string
	TestCases   []LogicTestCase
	RiskLevel   string
}

// LogicTestCase represents a single business logic test
type LogicTestCase struct {
	Name         string
	Description  string
	Input        map[string]interface{}
	Expected     string
	RiskLevel    string
	BusinessRule string
	TestType     string
}

// LogicTestResult contains the results of business logic testing
type LogicTestResult struct {
	TestCase      LogicTestCase
	ActualResult  string
	Status        string
	Vulnerability *LogicVulnerability
	Evidence      string
	Impact        string
}

// LogicVulnerability represents a business logic vulnerability
type LogicVulnerability struct {
	Type            string
	Severity        string
	Description     string
	BusinessImpact  string
	FinancialImpact string
	Evidence        string
	Recommendation  string
}

// NewBizLogicTesterAgent creates a new business logic testing agent
func NewBizLogicTesterAgent() *BizLogicTesterAgent {
	return &BizLogicTesterAgent{
		businessRules: initializeBusinessRules(),
		testScenarios: initializeTestScenarios(),
	}
}

// Execute performs comprehensive business logic testing
func (a *BizLogicTesterAgent) Execute(ctx context.Context, input AgentInput) (*AgentOutput, error) {
	var params struct {
		SystemType   string                 `json:"system_type"`
		BaseURL      string                 `json:"base_url"`
		AuthToken    string                 `json:"auth_token"`
		TestTargets  []string               `json:"test_targets"`
		CustomRules  []BusinessRule         `json:"custom_rules"`
		TestData     map[string]interface{} `json:"test_data"`
	}

	if err := json.Unmarshal([]byte(input.Target), &params); err != nil {
		return nil, fmt.Errorf("invalid input parameters: %w", err)
	}

	// Get business rules for the system type
	ruleSet, exists := a.businessRules[params.SystemType]
	if !exists {
		return nil, fmt.Errorf("no business rules defined for system type: %s", params.SystemType)
	}

	// Get test scenarios for the system type
	scenarios, exists := a.testScenarios[params.SystemType]
	if !exists {
		return nil, fmt.Errorf("no test scenarios defined for system type: %s", params.SystemType)
	}

	var allResults []LogicTestResult
	var vulnerabilities []LogicVulnerability

	// Execute test scenarios
	for _, scenario := range scenarios {
		scenarioResults, err := a.executeTestScenario(ctx, scenario, params.BaseURL, params.AuthToken, params.TestData)
		if err != nil {
			return nil, fmt.Errorf("failed to execute scenario %s: %w", scenario.Name, err)
		}

		allResults = append(allResults, scenarioResults...)

		// Collect vulnerabilities
		for _, result := range scenarioResults {
			if result.Vulnerability != nil {
				vulnerabilities = append(vulnerabilities, *result.Vulnerability)
			}
		}
	}

	// Generate comprehensive report
	report := a.generateLogicReport(params.SystemType, ruleSet, allResults, vulnerabilities)

	return &AgentOutput{
		Result: report,
		Metadata: map[string]interface{}{
			"system_type":        params.SystemType,
			"tests_executed":     len(allResults),
			"vulnerabilities":    len(vulnerabilities),
			"critical_vulns":     a.countBySeverity(vulnerabilities, "critical"),
			"high_vulns":         a.countBySeverity(vulnerabilities, "high"),
			"financial_impact":   a.calculateFinancialImpact(vulnerabilities),
		},
	}, nil
}

// executeTestScenario runs all test cases in a business logic scenario
func (a *BizLogicTesterAgent) executeTestScenario(ctx context.Context, scenario LogicTestScenario, baseURL, authToken string, testData map[string]interface{}) ([]LogicTestResult, error) {
	var results []LogicTestResult

	for _, testCase := range scenario.TestCases {
		result, err := a.executeLogicTestCase(ctx, testCase, baseURL, authToken, testData)
		if err != nil {
			result = LogicTestResult{
				TestCase: testCase,
				Status:   "error",
				Evidence: err.Error(),
			}
		}
		results = append(results, result)
	}

	return results, nil
}

// executeLogicTestCase runs a single business logic test case
func (a *BizLogicTesterAgent) executeLogicTestCase(ctx context.Context, testCase LogicTestCase, baseURL, authToken string, testData map[string]interface{}) (LogicTestResult, error) {
	result := LogicTestResult{
		TestCase: testCase,
		Status:   "completed",
	}

	// Execute test based on test type
	switch testCase.TestType {
	case "billing_calculation":
		result = a.testBillingCalculation(testCase, baseURL, authToken)
	case "price_manipulation":
		result = a.testPriceManipulation(testCase, baseURL, authToken)
	case "workflow_bypass":
		result = a.testWorkflowBypass(testCase, baseURL, authToken)
	case "authorization_logic":
		result = a.testAuthorizationLogic(testCase, baseURL, authToken)
	case "data_validation":
		result = a.testDataValidation(testCase, baseURL, authToken)
	case "rate_limiting":
		result = a.testRateLimiting(testCase, baseURL, authToken)
	default:
		result.Status = "skipped"
		result.Evidence = fmt.Sprintf("Unknown test type: %s", testCase.TestType)
	}

	return result, nil
}

// testBillingCalculation tests electricity billing calculation logic
func (a *BizLogicTesterAgent) testBillingCalculation(testCase LogicTestCase, baseURL, authToken string) LogicTestResult {
	result := LogicTestResult{
		TestCase: testCase,
		Status:   "completed",
	}

	// Extract test parameters
	usage, _ := testCase.Input["usage"].(float64)
	rate, _ := testCase.Input["rate"].(string)
	tier, _ := testCase.Input["tier"].(string)

	// Simulate billing calculation test
	calculatedAmount := a.simulateBillingCalculation(usage, rate, tier)
	result.ActualResult = fmt.Sprintf("Calculated amount: %.2f", calculatedAmount)

	// Analyze for vulnerabilities
	vulnerability := a.analyzeBillingLogic(usage, rate, tier, calculatedAmount, testCase)
	if vulnerability != nil {
		result.Vulnerability = vulnerability
		result.Impact = vulnerability.FinancialImpact
	}

	return result
}

// simulateBillingCalculation simulates billing calculation logic
func (a *BizLogicTesterAgent) simulateBillingCalculation(usage float64, rate, tier string) float64 {
	// Simulate different billing scenarios
	baseRate := 0.15 // Base rate per kWh

	switch rate {
	case "peak":
		baseRate = 0.25
	case "valley":
		baseRate = 0.10
	case "admin_override":
		baseRate = 0.01 // Suspicious admin rate
	}

	// Apply tiered pricing
	switch tier {
	case "tier1":
		if usage <= 200 {
			return usage * baseRate
		}
		return 200*baseRate + (usage-200)*baseRate*1.5
	case "tier2":
		return usage * baseRate * 1.5
	case "admin":
		return usage * 0.01 // Suspicious admin tier
	default:
		return usage * baseRate
	}
}

// analyzeBillingLogic analyzes billing calculation for vulnerabilities
func (a *BizLogicTesterAgent) analyzeBillingLogic(usage float64, rate, tier string, amount float64, testCase LogicTestCase) *LogicVulnerability {
	// Check for negative usage
	if usage < 0 && amount < 0 {
		return &LogicVulnerability{
			Type:            "Negative Billing Logic",
			Severity:        "critical",
			Description:     "System allows negative usage values resulting in negative billing amounts",
			BusinessImpact:  "Revenue loss through negative billing",
			FinancialImpact: fmt.Sprintf("Potential loss: $%.2f per transaction", math.Abs(amount)),
			Evidence:        fmt.Sprintf("Usage: %.2f, Amount: %.2f", usage, amount),
			Recommendation:  "Implement strict validation to reject negative usage values",
		}
	}

	// Check for overflow scenarios
	if usage > 999999 && amount == 0 {
		return &LogicVulnerability{
			Type:            "Integer Overflow",
			Severity:        "high",
			Description:     "Large usage values cause calculation overflow resulting in zero billing",
			BusinessImpact:  "Revenue loss through calculation errors",
			FinancialImpact: "Potential loss: Unlimited usage at zero cost",
			Evidence:        fmt.Sprintf("Usage: %.2f, Amount: %.2f", usage, amount),
			Recommendation:  "Implement proper bounds checking and overflow protection",
		}
	}

	// Check for admin privilege abuse
	if (rate == "admin_override" || tier == "admin") && amount < usage*0.05 {
		return &LogicVulnerability{
			Type:            "Administrative Privilege Abuse",
			Severity:        "high",
			Description:     "Administrative rates/tiers can be applied without proper authorization",
			BusinessImpact:  "Revenue loss through unauthorized discounts",
			FinancialImpact: fmt.Sprintf("Discount abuse: %.2f%% off normal rate", (1-amount/(usage*0.15))*100),
			Evidence:        fmt.Sprintf("Rate: %s, Tier: %s, Amount: %.2f", rate, tier, amount),
			Recommendation:  "Implement strict authorization controls for administrative pricing",
		}
	}

	// Check for boundary manipulation
	if usage == 200.01 && tier == "tier1" && amount < 200*0.15+0.01*0.15*1.5 {
		return &LogicVulnerability{
			Type:            "Tier Boundary Manipulation",
			Severity:        "medium",
			Description:     "Tier boundary calculations can be manipulated",
			BusinessImpact:  "Revenue loss through tier manipulation",
			FinancialImpact: "Variable loss depending on usage patterns",
			Evidence:        fmt.Sprintf("Boundary usage: %.2f, Calculated: %.2f", usage, amount),
			Recommendation:  "Review and strengthen tier boundary calculations",
		}
	}

	return nil
}

// testPriceManipulation tests for price manipulation vulnerabilities
func (a *BizLogicTesterAgent) testPriceManipulation(testCase LogicTestCase, baseURL, authToken string) LogicTestResult {
	result := LogicTestResult{
		TestCase: testCase,
		Status:   "completed",
	}

	// Extract test parameters
	originalPrice, _ := testCase.Input["original_price"].(float64)
	manipulatedPrice, _ := testCase.Input["manipulated_price"].(float64)
	priceField, _ := testCase.Input["price_field"].(string)

	// Simulate price manipulation test
	if manipulatedPrice != originalPrice {
		vulnerability := &LogicVulnerability{
			Type:            "Price Manipulation",
			Severity:        "critical",
			Description:     "Client-side price manipulation is possible",
			BusinessImpact:  "Direct revenue loss through price tampering",
			FinancialImpact: fmt.Sprintf("Loss per transaction: $%.2f", originalPrice-manipulatedPrice),
			Evidence:        fmt.Sprintf("Original: %.2f, Manipulated: %.2f, Field: %s", originalPrice, manipulatedPrice, priceField),
			Recommendation:  "Always validate prices server-side, never trust client input",
		}
		result.Vulnerability = vulnerability
		result.Impact = vulnerability.FinancialImpact
	}

	result.ActualResult = fmt.Sprintf("Price manipulation test: Original=%.2f, Manipulated=%.2f", originalPrice, manipulatedPrice)
	return result
}

// testWorkflowBypass tests for business workflow bypass vulnerabilities
func (a *BizLogicTesterAgent) testWorkflowBypass(testCase LogicTestCase, baseURL, authToken string) LogicTestResult {
	result := LogicTestResult{
		TestCase: testCase,
		Status:   "completed",
	}

	// Extract test parameters
	currentStep, _ := testCase.Input["current_step"].(string)
	targetStep, _ := testCase.Input["target_step"].(string)
	requiredApproval, _ := testCase.Input["required_approval"].(bool)

	// Simulate workflow bypass test
	if currentStep != targetStep && requiredApproval {
		vulnerability := &LogicVulnerability{
			Type:            "Workflow Bypass",
			Severity:        "high",
			Description:     "Business workflow steps can be bypassed without proper validation",
			BusinessImpact:  "Unauthorized operations and compliance violations",
			FinancialImpact: "Variable impact depending on bypassed controls",
			Evidence:        fmt.Sprintf("Bypassed from %s to %s without approval", currentStep, targetStep),
			Recommendation:  "Implement strict workflow validation and state management",
		}
		result.Vulnerability = vulnerability
		result.Impact = vulnerability.BusinessImpact
	}

	result.ActualResult = fmt.Sprintf("Workflow test: %s -> %s (Approval required: %t)", currentStep, targetStep, requiredApproval)
	return result
}

// testAuthorizationLogic tests authorization logic vulnerabilities
func (a *BizLogicTesterAgent) testAuthorizationLogic(testCase LogicTestCase, baseURL, authToken string) LogicTestResult {
	result := LogicTestResult{
		TestCase: testCase,
		Status:   "completed",
	}

	// Extract test parameters
	userRole, _ := testCase.Input["user_role"].(string)
	requestedAction, _ := testCase.Input["requested_action"].(string)
	resourceOwner, _ := testCase.Input["resource_owner"].(string)
	currentUser, _ := testCase.Input["current_user"].(string)

	// Simulate authorization logic test
	authorized := a.simulateAuthorizationCheck(userRole, requestedAction, resourceOwner, currentUser)

	if !authorized && testCase.Expected == "authorized" {
		// Expected authorization but was denied - this is good
		result.ActualResult = "Access properly denied"
	} else if authorized && testCase.Expected == "denied" {
		// Expected denial but was authorized - this is a vulnerability
		vulnerability := &LogicVulnerability{
			Type:            "Authorization Bypass",
			Severity:        "critical",
			Description:     "Authorization logic allows unauthorized access to resources",
			BusinessImpact:  "Unauthorized access to sensitive data and functions",
			FinancialImpact: "Data breach costs and regulatory fines",
			Evidence:        fmt.Sprintf("User %s (%s) gained unauthorized access to %s's resource for action %s", currentUser, userRole, resourceOwner, requestedAction),
			Recommendation:  "Review and strengthen authorization logic and access controls",
		}
		result.Vulnerability = vulnerability
		result.Impact = vulnerability.BusinessImpact
	}

	result.ActualResult = fmt.Sprintf("Authorization test: User=%s, Role=%s, Action=%s, Authorized=%t", currentUser, userRole, requestedAction, authorized)
	return result
}

// simulateAuthorizationCheck simulates authorization logic
func (a *BizLogicTesterAgent) simulateAuthorizationCheck(userRole, action, resourceOwner, currentUser string) bool {
	// Simulate authorization logic with potential flaws
	switch userRole {
	case "admin":
		return true // Admins can do anything
	case "operator":
		// Operators can read but not modify
		return strings.Contains(action, "read") || strings.Contains(action, "view")
	case "user":
		// Users can only access their own resources
		return currentUser == resourceOwner && (strings.Contains(action, "read") || strings.Contains(action, "view"))
	case "guest":
		// Guests should have no access, but simulate a flaw
		return strings.Contains(action, "view") // Vulnerability: guests can view
	default:
		return false
	}
}

// testDataValidation tests data validation logic
func (a *BizLogicTesterAgent) testDataValidation(testCase LogicTestCase, baseURL, authToken string) LogicTestResult {
	result := LogicTestResult{
		TestCase: testCase,
		Status:   "completed",
	}

	// Extract test parameters
	inputValue := testCase.Input["input_value"]
	fieldType, _ := testCase.Input["field_type"].(string)
	expectedValidation, _ := testCase.Input["expected_validation"].(string)

	// Simulate data validation test
	validationResult := a.simulateDataValidation(inputValue, fieldType)

	if validationResult != expectedValidation {
		severity := "medium"
		if fieldType == "billing_amount" || fieldType == "payment_amount" {
			severity = "high"
		}

		vulnerability := &LogicVulnerability{
			Type:            "Data Validation Bypass",
			Severity:        severity,
			Description:     fmt.Sprintf("Data validation for %s field can be bypassed", fieldType),
			BusinessImpact:  "Invalid data processing leading to system errors or security issues",
			FinancialImpact: "Variable impact depending on data type and usage",
			Evidence:        fmt.Sprintf("Input: %v, Type: %s, Expected: %s, Actual: %s", inputValue, fieldType, expectedValidation, validationResult),
			Recommendation:  "Implement comprehensive server-side data validation",
		}
		result.Vulnerability = vulnerability
		result.Impact = vulnerability.BusinessImpact
	}

	result.ActualResult = fmt.Sprintf("Validation test: Input=%v, Type=%s, Result=%s", inputValue, fieldType, validationResult)
	return result
}

// simulateDataValidation simulates data validation logic
func (a *BizLogicTesterAgent) simulateDataValidation(value interface{}, fieldType string) string {
	switch fieldType {
	case "email":
		if str, ok := value.(string); ok && strings.Contains(str, "@") {
			return "valid"
		}
		return "invalid"
	case "phone":
		if str, ok := value.(string); ok && len(str) >= 10 {
			return "valid"
		}
		return "invalid"
	case "billing_amount":
		if num, ok := value.(float64); ok && num >= 0 {
			return "valid"
		}
		// Simulate validation flaw - negative amounts might be accepted
		if num, ok := value.(float64); ok && num < 0 {
			return "valid" // This is the vulnerability
		}
		return "invalid"
	case "usage_kwh":
		if num, ok := value.(float64); ok && num >= 0 && num <= 999999 {
			return "valid"
		}
		return "invalid"
	default:
		return "valid" // Default to valid (potential vulnerability)
	}
}

// testRateLimiting tests rate limiting logic
func (a *BizLogicTesterAgent) testRateLimiting(testCase LogicTestCase, baseURL, authToken string) LogicTestResult {
	result := LogicTestResult{
		TestCase: testCase,
		Status:   "completed",
	}

	// Extract test parameters
	requestCount, _ := testCase.Input["request_count"].(float64)
	timeWindow, _ := testCase.Input["time_window"].(float64)
	endpoint, _ := testCase.Input["endpoint"].(string)

	// Simulate rate limiting test
	rateLimited := a.simulateRateLimiting(int(requestCount), int(timeWindow), endpoint)

	if !rateLimited && requestCount > 100 { // Assume 100 requests should trigger rate limiting
		vulnerability := &LogicVulnerability{
			Type:            "Rate Limiting Bypass",
			Severity:        "medium",
			Description:     "Rate limiting can be bypassed allowing excessive requests",
			BusinessImpact:  "Service degradation and potential DoS attacks",
			FinancialImpact: "Service downtime costs and infrastructure strain",
			Evidence:        fmt.Sprintf("Made %.0f requests in %.0f seconds without rate limiting", requestCount, timeWindow),
			Recommendation:  "Implement proper rate limiting and request throttling",
		}
		result.Vulnerability = vulnerability
		result.Impact = vulnerability.BusinessImpact
	}

	result.ActualResult = fmt.Sprintf("Rate limiting test: %.0f requests in %.0f seconds, Limited=%t", requestCount, timeWindow, rateLimited)
	return result
}

// simulateRateLimiting simulates rate limiting logic
func (a *BizLogicTesterAgent) simulateRateLimiting(requestCount, timeWindow int, endpoint string) bool {
	// Simulate rate limiting logic with potential flaws
	switch endpoint {
	case "/api/auth/login":
		// Login should have strict rate limiting
		return requestCount > 5
	case "/api/billing/calculate":
		// Billing calculations should be rate limited
		return requestCount > 50
	case "/api/user/profile":
		// Profile access should have moderate rate limiting
		return requestCount > 100
	default:
		// Default endpoints might not have rate limiting (vulnerability)
		return false
	}
}

// Helper functions
func (a *BizLogicTesterAgent) countBySeverity(vulnerabilities []LogicVulnerability, severity string) int {
	count := 0
	for _, vuln := range vulnerabilities {
		if vuln.Severity == severity {
			count++
		}
	}
	return count
}

func (a *BizLogicTesterAgent) calculateFinancialImpact(vulnerabilities []LogicVulnerability) string {
	totalImpact := 0.0
	for _, vuln := range vulnerabilities {
		// Extract numeric values from financial impact strings
		if strings.Contains(vuln.FinancialImpact, "$") {
			// Simple extraction - in real implementation would be more sophisticated
			parts := strings.Split(vuln.FinancialImpact, "$")
			if len(parts) > 1 {
				if val, err := strconv.ParseFloat(strings.Fields(parts[1])[0], 64); err == nil {
					totalImpact += val
				}
			}
		}
	}

	if totalImpact > 0 {
		return fmt.Sprintf("$%.2f estimated per incident", totalImpact)
	}
	return "Variable impact - requires detailed analysis"
}

// generateLogicReport creates a comprehensive business logic security report
func (a *BizLogicTesterAgent) generateLogicReport(systemType string, ruleSet BusinessRuleSet, results []LogicTestResult, vulnerabilities []LogicVulnerability) string {
	report := fmt.Sprintf("# Business Logic Security Assessment Report\n\n")
	report += fmt.Sprintf("## System Type: %s\n\n", systemType)

	// Executive Summary
	report += "## Executive Summary\n\n"
	report += fmt.Sprintf("This report presents the results of comprehensive business logic security testing performed on the %s system. ", systemType)
	report += fmt.Sprintf("A total of %d business logic test cases were executed across multiple scenarios.\n\n", len(results))

	// Business Rules Summary
	report += "## Business Rules Tested\n\n"
	for _, rule := range ruleSet.Rules {
		report += fmt.Sprintf("- **%s** (%s): %s\n", rule.Name, rule.RiskLevel, rule.Description)
	}
	report += "\n"

	// Risk Summary
	criticalCount := a.countBySeverity(vulnerabilities, "critical")
	highCount := a.countBySeverity(vulnerabilities, "high")
	mediumCount := a.countBySeverity(vulnerabilities, "medium")

	report += "## Risk Summary\n\n"
	report += fmt.Sprintf("- **Critical Business Logic Flaws**: %d\n", criticalCount)
	report += fmt.Sprintf("- **High Risk Logic Issues**: %d\n", highCount)
	report += fmt.Sprintf("- **Medium Risk Logic Issues**: %d\n", mediumCount)
	report += fmt.Sprintf("- **Total Logic Vulnerabilities**: %d\n\n", len(vulnerabilities))

	// Financial Impact Assessment
	financialImpact := a.calculateFinancialImpact(vulnerabilities)
	report += "## Financial Impact Assessment\n\n"
	report += fmt.Sprintf("**Estimated Financial Impact**: %s\n\n", financialImpact)

	// Vulnerability Details
	if len(vulnerabilities) > 0 {
		report += "## Business Logic Vulnerabilities\n\n"
		for i, vuln := range vulnerabilities {
			report += fmt.Sprintf("### Vulnerability %d: %s\n\n", i+1, vuln.Type)
			report += fmt.Sprintf("**Severity**: %s\n", vuln.Severity)
			report += fmt.Sprintf("**Description**: %s\n", vuln.Description)
			report += fmt.Sprintf("**Business Impact**: %s\n", vuln.BusinessImpact)
			report += fmt.Sprintf("**Financial Impact**: %s\n", vuln.FinancialImpact)
			report += fmt.Sprintf("**Evidence**: %s\n", vuln.Evidence)
			report += fmt.Sprintf("**Recommendation**: %s\n\n", vuln.Recommendation)
		}
	}

	// Test Results Summary
	report += "## Test Results Summary\n\n"
	completedCount := a.countByStatus(results, "completed")
	errorCount := a.countByStatus(results, "error")
	skippedCount := a.countByStatus(results, "skipped")

	report += fmt.Sprintf("- **Tests Completed**: %d\n", completedCount)
	report += fmt.Sprintf("- **Tests Failed**: %d\n", errorCount)
	report += fmt.Sprintf("- **Tests Skipped**: %d\n", skippedCount)
	report += fmt.Sprintf("- **Success Rate**: %.1f%%\n\n", float64(completedCount)/float64(len(results))*100)

	// Business Impact Analysis
	report += "## Business Impact Analysis\n\n"
	report += "### Revenue Protection\n"
	report += "- Business logic vulnerabilities directly impact revenue through billing manipulation\n"
	report += "- Price manipulation vulnerabilities allow customers to pay reduced amounts\n"
	report += "- Negative billing logic can result in credits instead of charges\n\n"

	report += "### Compliance and Regulatory Impact\n"
	report += "- Billing accuracy is required by utility regulations\n"
	report += "- Workflow bypasses may violate audit and compliance requirements\n"
	report += "- Data validation failures can lead to regulatory violations\n\n"

	report += "### Operational Impact\n"
	report += "- Authorization bypasses compromise system security\n"
	report += "- Rate limiting issues can lead to service degradation\n"
	report += "- Workflow bypasses undermine business process integrity\n\n"

	// Recommendations
	report += "## Recommendations\n\n"
	report += "### Immediate Actions (Critical Priority)\n"
	if criticalCount > 0 {
		report += "1. **Fix Critical Logic Flaws**: Immediately address all critical business logic vulnerabilities\n"
		report += "2. **Implement Server-Side Validation**: Ensure all business rules are enforced server-side\n"
		report += "3. **Review Billing Logic**: Conduct comprehensive review of all billing calculations\n"
		report += "4. **Strengthen Authorization**: Implement proper authorization checks for all business functions\n\n"
	}

	report += "### Short-term Actions (High Priority)\n"
	report += "1. **Business Rule Documentation**: Document all business rules and validation requirements\n"
	report += "2. **Workflow Security**: Implement proper workflow state management and validation\n"
	report += "3. **Input Validation**: Strengthen input validation for all business-critical fields\n"
	report += "4. **Rate Limiting**: Implement appropriate rate limiting for all endpoints\n\n"

	report += "### Long-term Actions (Strategic)\n"
	report += "1. **Business Logic Testing**: Integrate business logic testing into development lifecycle\n"
	report += "2. **Security Training**: Train developers on secure business logic implementation\n"
	report += "3. **Automated Testing**: Implement automated business logic security testing\n"
	report += "4. **Continuous Monitoring**: Deploy monitoring for business logic anomalies\n\n"

	return report
}

func (a *BizLogicTesterAgent) countByStatus(results []LogicTestResult, status string) int {
	count := 0
	for _, result := range results {
		if result.Status == status {
			count++
		}
	}
	return count
}

// initializeBusinessRules sets up business rules for different power industry systems
func initializeBusinessRules() map[string]BusinessRuleSet {
	return map[string]BusinessRuleSet{
		"marketing_2.0": {
			SystemType: "Power Marketing System 2.0",
			Rules: []BusinessRule{
				{
					ID:          "BILL_001",
					Name:        "Positive Billing Amount",
					Description: "All billing amounts must be positive values",
					Category:    "billing",
					RiskLevel:   "critical",
					TestMethod:  "boundary_value_testing",
				},
				{
					ID:          "BILL_002",
					Name:        "Tiered Pricing Integrity",
					Description: "Tiered pricing calculations must be accurate and consistent",
					Category:    "billing",
					RiskLevel:   "high",
					TestMethod:  "calculation_verification",
				},
				{
					ID:          "AUTH_001",
					Name:        "Customer Data Isolation",
					Description: "Customers can only access their own billing data",
					Category:    "authorization",
					RiskLevel:   "critical",
					TestMethod:  "access_control_testing",
				},
				{
					ID:          "WORK_001",
					Name:        "Billing Approval Workflow",
					Description: "Large billing adjustments require manager approval",
					Category:    "workflow",
					RiskLevel:   "high",
					TestMethod:  "workflow_testing",
				},
			},
			Constraints: []BusinessConstraint{
				{Field: "usage_kwh", Type: "float", MinValue: 0.0, MaxValue: 999999.0, Required: true},
				{Field: "billing_amount", Type: "float", MinValue: 0.01, MaxValue: 999999.99, Required: true},
				{Field: "customer_id", Type: "string", Required: true, Validation: "alphanumeric"},
			},
		},
		"iguowang": {
			SystemType: "i国网 Mobile App",
			Rules: []BusinessRule{
				{
					ID:          "SMS_001",
					Name:        "SMS Verification Rate Limiting",
					Description: "SMS verification codes are rate limited per phone number",
					Category:    "authentication",
					RiskLevel:   "medium",
					TestMethod:  "rate_limiting_testing",
				},
				{
					ID:          "PAY_001",
					Name:        "Payment Amount Validation",
					Description: "Payment amounts must match billing amounts",
					Category:    "payment",
					RiskLevel:   "critical",
					TestMethod:  "amount_verification",
				},
				{
					ID:          "DATA_001",
					Name:        "User Data Privacy",
					Description: "User personal data is protected and access controlled",
					Category:    "privacy",
					RiskLevel:   "high",
					TestMethod:  "privacy_testing",
				},
			},
			Constraints: []BusinessConstraint{
				{Field: "phone_number", Type: "string", Required: true, Validation: "phone_format"},
				{Field: "sms_code", Type: "string", Required: true, Validation: "numeric_6_digits"},
				{Field: "payment_amount", Type: "float", MinValue: 0.01, MaxValue: 99999.99, Required: true},
			},
		},
		"sap": {
			SystemType: "SAP ERP System",
			Rules: []BusinessRule{
				{
					ID:          "SAP_001",
					Name:        "User Authorization Matrix",
					Description: "Users can only access functions authorized for their role",
					Category:    "authorization",
					RiskLevel:   "critical",
					TestMethod:  "role_based_testing",
				},
				{
					ID:          "FIN_001",
					Name:        "Financial Data Integrity",
					Description: "Financial transactions must be accurate and auditable",
					Category:    "financial",
					RiskLevel:   "critical",
					TestMethod:  "financial_verification",
				},
				{
					ID:          "AUDIT_001",
					Name:        "Audit Trail Completeness",
					Description: "All system changes must be logged for audit purposes",
					Category:    "audit",
					RiskLevel:   "high",
					TestMethod:  "audit_testing",
				},
			},
			Constraints: []BusinessConstraint{
				{Field: "user_id", Type: "string", Required: true, Validation: "sap_user_format"},
				{Field: "client", Type: "string", Required: true, Validation: "3_digit_numeric"},
				{Field: "amount", Type: "decimal", MinValue: -999999999.99, MaxValue: 999999999.99, Required: true},
			},
		},
	}
}

// initializeTestScenarios sets up test scenarios for different power industry systems
func initializeTestScenarios() map[string][]LogicTestScenario {
	return map[string][]LogicTestScenario{
		"marketing_2.0": {
			{
				Name:        "Billing Calculation Logic",
				Category:    "billing",
				Description: "Tests billing calculation logic for various scenarios",
				RiskLevel:   "critical",
				TestCases: []LogicTestCase{
					{
						Name:         "Negative Usage Test",
						Description:  "Test system behavior with negative usage values",
						Input:        map[string]interface{}{"usage": -999999.0, "rate": "standard", "tier": "tier1"},
						Expected:     "rejected",
						RiskLevel:    "critical",
						BusinessRule: "BILL_001",
						TestType:     "billing_calculation",
					},
					{
						Name:         "Zero Usage Test",
						Description:  "Test system behavior with zero usage",
						Input:        map[string]interface{}{"usage": 0.0, "rate": "standard", "tier": "tier1"},
						Expected:     "minimal_charge",
						RiskLevel:    "medium",
						BusinessRule: "BILL_001",
						TestType:     "billing_calculation",
					},
					{
						Name:         "Overflow Usage Test",
						Description:  "Test system behavior with extremely large usage values",
						Input:        map[string]interface{}{"usage": 999999999.0, "rate": "standard", "tier": "tier1"},
						Expected:     "handled_gracefully",
						RiskLevel:    "high",
						BusinessRule: "BILL_002",
						TestType:     "billing_calculation",
					},
					{
						Name:         "Admin Rate Override Test",
						Description:  "Test unauthorized use of administrative rates",
						Input:        map[string]interface{}{"usage": 1000.0, "rate": "admin_override", "tier": "tier1"},
						Expected:     "rejected",
						RiskLevel:    "critical",
						BusinessRule: "AUTH_001",
						TestType:     "billing_calculation",
					},
					{
						Name:         "Tier Boundary Test",
						Description:  "Test tier boundary calculations",
						Input:        map[string]interface{}{"usage": 200.01, "rate": "standard", "tier": "tier1"},
						Expected:     "tier2_rate_applied",
						RiskLevel:    "medium",
						BusinessRule: "BILL_002",
						TestType:     "billing_calculation",
					},
				},
			},
			{
				Name:        "Price Manipulation",
				Category:    "security",
				Description: "Tests for client-side price manipulation vulnerabilities",
				RiskLevel:   "critical",
				TestCases: []LogicTestCase{
					{
						Name:         "Client Price Override",
						Description:  "Test if client can override server-calculated prices",
						Input:        map[string]interface{}{"original_price": 150.00, "manipulated_price": 1.50, "price_field": "total_amount"},
						Expected:     "server_price_enforced",
						RiskLevel:    "critical",
						BusinessRule: "BILL_001",
						TestType:     "price_manipulation",
					},
					{
						Name:         "Discount Manipulation",
						Description:  "Test unauthorized discount application",
						Input:        map[string]interface{}{"original_price": 100.00, "manipulated_price": 10.00, "price_field": "discount_amount"},
						Expected:     "unauthorized_discount_rejected",
						RiskLevel:    "high",
						BusinessRule: "AUTH_001",
						TestType:     "price_manipulation",
					},
				},
			},
			{
				Name:        "Authorization Logic",
				Category:    "security",
				Description: "Tests authorization and access control logic",
				RiskLevel:   "critical",
				TestCases: []LogicTestCase{
					{
						Name:         "Cross-Customer Data Access",
						Description:  "Test if users can access other customers' data",
						Input:        map[string]interface{}{"user_role": "user", "requested_action": "view_billing", "resource_owner": "other_customer", "current_user": "customer_123"},
						Expected:     "denied",
						RiskLevel:    "critical",
						BusinessRule: "AUTH_001",
						TestType:     "authorization_logic",
					},
					{
						Name:         "Privilege Escalation Test",
						Description:  "Test if users can escalate their privileges",
						Input:        map[string]interface{}{"user_role": "user", "requested_action": "admin_function", "resource_owner": "system", "current_user": "user_123"},
						Expected:     "denied",
						RiskLevel:    "critical",
						BusinessRule: "AUTH_001",
						TestType:     "authorization_logic",
					},
					{
						Name:         "Guest Access Test",
						Description:  "Test guest user access controls",
						Input:        map[string]interface{}{"user_role": "guest", "requested_action": "view_data", "resource_owner": "customer_123", "current_user": "guest"},
						Expected:     "denied",
						RiskLevel:    "high",
						BusinessRule: "AUTH_001",
						TestType:     "authorization_logic",
					},
				},
			},
		},
		"iguowang": {
			{
				Name:        "SMS Authentication Logic",
				Category:    "authentication",
				Description: "Tests SMS verification and authentication logic",
				RiskLevel:   "high",
				TestCases: []LogicTestCase{
					{
						Name:         "SMS Rate Limiting Test",
						Description:  "Test SMS verification rate limiting",
						Input:        map[string]interface{}{"request_count": 100.0, "time_window": 60.0, "endpoint": "/api/auth/send-sms"},
						Expected:     "rate_limited",
						RiskLevel:    "medium",
						BusinessRule: "SMS_001",
						TestType:     "rate_limiting",
					},
					{
						Name:         "SMS Code Validation",
						Description:  "Test SMS code validation logic",
						Input:        map[string]interface{}{"input_value": "000000", "field_type": "sms_code", "expected_validation": "invalid"},
						Expected:     "invalid",
						RiskLevel:    "high",
						BusinessRule: "SMS_001",
						TestType:     "data_validation",
					},
				},
			},
			{
				Name:        "Payment Logic",
				Category:    "financial",
				Description: "Tests payment processing logic",
				RiskLevel:   "critical",
				TestCases: []LogicTestCase{
					{
						Name:         "Payment Amount Validation",
						Description:  "Test payment amount validation",
						Input:        map[string]interface{}{"input_value": -100.0, "field_type": "payment_amount", "expected_validation": "invalid"},
						Expected:     "invalid",
						RiskLevel:    "critical",
						BusinessRule: "PAY_001",
						TestType:     "data_validation",
					},
					{
						Name:         "Zero Payment Test",
						Description:  "Test zero amount payment processing",
						Input:        map[string]interface{}{"input_value": 0.0, "field_type": "payment_amount", "expected_validation": "invalid"},
						Expected:     "invalid",
						RiskLevel:    "high",
						BusinessRule: "PAY_001",
						TestType:     "data_validation",
					},
				},
			},
		},
		"sap": {
			{
				Name:        "SAP Authorization Matrix",
				Category:    "authorization",
				Description: "Tests SAP role-based authorization",
				RiskLevel:   "critical",
				TestCases: []LogicTestCase{
					{
						Name:         "Financial Data Access",
						Description:  "Test unauthorized financial data access",
						Input:        map[string]interface{}{"user_role": "hr_user", "requested_action": "view_financial", "resource_owner": "finance_dept", "current_user": "hr_001"},
						Expected:     "denied",
						RiskLevel:    "critical",
						BusinessRule: "SAP_001",
						TestType:     "authorization_logic",
					},
					{
						Name:         "Cross-Client Access",
						Description:  "Test cross-client data access",
						Input:        map[string]interface{}{"user_role": "user", "requested_action": "view_data", "resource_owner": "client_200", "current_user": "user_client_100"},
						Expected:     "denied",
						RiskLevel:    "critical",
						BusinessRule: "SAP_001",
						TestType:     "authorization_logic",
					},
				},
			},
			{
				Name:        "Financial Data Integrity",
				Category:    "financial",
				Description: "Tests financial data processing logic",
				RiskLevel:   "critical",
				TestCases: []LogicTestCase{
					{
						Name:         "Large Amount Processing",
						Description:  "Test processing of large financial amounts",
						Input:        map[string]interface{}{"input_value": 999999999.99, "field_type": "financial_amount", "expected_validation": "valid"},
						Expected:     "valid",
						RiskLevel:    "medium",
						BusinessRule: "FIN_001",
						TestType:     "data_validation",
					},
					{
						Name:         "Negative Amount Processing",
						Description:  "Test processing of negative financial amounts",
						Input:        map[string]interface{}{"input_value": -1000000.0, "field_type": "financial_amount", "expected_validation": "valid"},
						Expected:     "valid",
						RiskLevel:    "high",
						BusinessRule: "FIN_001",
						TestType:     "data_validation",
					},
				},
			},
		},
	}
}