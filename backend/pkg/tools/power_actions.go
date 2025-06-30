package tools

// PowerPentesterAction represents the action for power industry penetration testing
type PowerPentesterAction struct {
	SystemType   string                 `json:"system_type" jsonschema:"enum=marketing_2.0,enum=iguowang,enum=sap,description=Type of power industry system to test (marketing_2.0 for Power Marketing System 2.0, iguowang for i国网APP, sap for SAP ERP)"`
	Target       string                 `json:"target" jsonschema:"description=Target system URL or IP address"`
	TestScope    string                 `json:"test_scope,omitempty" jsonschema:"enum=full,enum=quick,enum=focused,description=Scope of penetration testing (full for comprehensive testing, quick for basic checks, focused for specific vulnerabilities)"`
	AuthToken    string                 `json:"auth_token,omitempty" jsonschema:"description=Authentication token if available"`
	CustomTests  []string               `json:"custom_tests,omitempty" jsonschema:"description=List of specific test types to execute"`
	TestData     map[string]interface{} `json:"test_data,omitempty" jsonschema:"description=Custom test data for specific scenarios"`
}

// APITesterAction represents the action for API security testing
type APITesterAction struct {
	SystemType   string            `json:"system_type" jsonschema:"enum=marketing_2.0,enum=iguowang,enum=sap,description=Type of power industry system"`
	BaseURL      string            `json:"base_url" jsonschema:"description=Base URL of the API to test"`
	AuthToken    string            `json:"auth_token,omitempty" jsonschema:"description=Authentication token for API access"`
	Headers      map[string]string `json:"headers,omitempty" jsonschema:"description=Additional HTTP headers for requests"`
	TestSuites   []string          `json:"test_suites,omitempty" jsonschema:"description=Specific test suites to execute (authentication, authorization, billing_api, etc.)"`
	Endpoints    []string          `json:"endpoints,omitempty" jsonschema:"description=Specific API endpoints to test"`
}

// BizLogicTesterAction represents the action for business logic testing
type BizLogicTesterAction struct {
	SystemType   string                 `json:"system_type" jsonschema:"enum=marketing_2.0,enum=iguowang,enum=sap,description=Type of power industry system"`
	BaseURL      string                 `json:"base_url" jsonschema:"description=Base URL of the system to test"`
	AuthToken    string                 `json:"auth_token,omitempty" jsonschema:"description=Authentication token for system access"`
	TestTargets  []string               `json:"test_targets,omitempty" jsonschema:"description=Specific business logic areas to test (billing, workflow, authorization, etc.)"`
	CustomRules  []string               `json:"custom_rules,omitempty" jsonschema:"description=Custom business rules to validate"`
	TestData     map[string]interface{} `json:"test_data,omitempty" jsonschema:"description=Test data for business logic scenarios"`
}

// ComplianceAgentAction represents the action for compliance testing
type ComplianceAgentAction struct {
	SystemType   string   `json:"system_type" jsonschema:"enum=marketing_2.0,enum=iguowang,enum=sap,description=Type of power industry system"`
	BaseURL      string   `json:"base_url" jsonschema:"description=Base URL of the system to assess"`
	Frameworks   []string `json:"frameworks,omitempty" jsonschema:"description=Compliance frameworks to assess against (iso27001, nist_csf, gdpr, etc.)"`
	Standards    []string `json:"standards,omitempty" jsonschema:"description=Regulatory standards to assess (nerc_cip, ferc_standards, sox_compliance, etc.)"`
	TestScope    string   `json:"test_scope,omitempty" jsonschema:"enum=full,enum=data_protection,enum=access_control,enum=audit_logging,description=Scope of compliance assessment"`
	AuthToken    string   `json:"auth_token,omitempty" jsonschema:"description=Authentication token for system access"`
}

// BillingLogicTestAction represents the action for billing logic testing
type BillingLogicTestAction struct {
	Endpoint     string                 `json:"endpoint" jsonschema:"description=Billing API endpoint to test"`
	UserID       string                 `json:"user_id,omitempty" jsonschema:"description=User ID for testing context"`
	Token        string                 `json:"token,omitempty" jsonschema:"description=Authentication token"`
	TestScenarios []string              `json:"test_scenarios,omitempty" jsonschema:"description=Specific billing test scenarios (boundary_values, tiered_pricing, time_of_use, etc.)"`
	CustomData   map[string]interface{} `json:"custom_data,omitempty" jsonschema:"description=Custom test data for billing scenarios"`
}

// APIFuzzerAction represents the action for API fuzzing
type APIFuzzerAction struct {
	BaseURL      string            `json:"base_url" jsonschema:"description=Base URL of the API to fuzz"`
	Endpoints    []string          `json:"endpoints" jsonschema:"description=List of API endpoints to fuzz"`
	Headers      map[string]string `json:"headers,omitempty" jsonschema:"description=HTTP headers to include in requests"`
	AuthToken    string            `json:"auth_token,omitempty" jsonschema:"description=Authentication token"`
	PayloadTypes []string          `json:"payload_types,omitempty" jsonschema:"description=Types of payloads to use (auth_bypass, injection, business_logic, etc.)"`
	Intensity    string            `json:"intensity,omitempty" jsonschema:"enum=light,enum=medium,enum=aggressive,description=Fuzzing intensity level"`
}

// SAPScannerAction represents the action for SAP security scanning
type SAPScannerAction struct {
	Target       string   `json:"target" jsonschema:"description=SAP system target (IP address or hostname)"`
	Port         int      `json:"port,omitempty" jsonschema:"description=SAP system port (default 8000)"`
	Client       string   `json:"client,omitempty" jsonschema:"description=SAP client number (e.g., 000, 100)"`
	ScanTypes    []string `json:"scan_types,omitempty" jsonschema:"description=Types of scans to perform (default_creds, config_issues, known_vulns, etc.)"`
	Credentials  []string `json:"credentials,omitempty" jsonschema:"description=Credentials to test (format: username:password:client)"`
}

// MobileSecurityAction represents the action for mobile security testing
type MobileSecurityAction struct {
	AppType      string   `json:"app_type" jsonschema:"enum=iguowang,enum=power_mobile,description=Type of mobile application"`
	BaseURL      string   `json:"base_url" jsonschema:"description=Base URL of the mobile app backend"`
	TestAreas    []string `json:"test_areas,omitempty" jsonschema:"description=Areas to test (sms_auth, api_security, data_storage, etc.)"`
	AuthToken    string   `json:"auth_token,omitempty" jsonschema:"description=Authentication token"`
	DeviceInfo   string   `json:"device_info,omitempty" jsonschema:"description=Device information for testing context"`
}

// PowerDataAnalyzerAction represents the action for power data analysis
type PowerDataAnalyzerAction struct {
	DataSource   string                 `json:"data_source" jsonschema:"description=Data source to analyze (api_endpoint, log_file, database, etc.)"`
	DataType     string                 `json:"data_type" jsonschema:"enum=usage_patterns,enum=billing_data,enum=customer_info,enum=system_logs,description=Type of power industry data to analyze"`
	AnalysisType string                 `json:"analysis_type" jsonschema:"enum=security_anomalies,enum=privacy_risks,enum=compliance_issues,description=Type of analysis to perform"`
	Parameters   map[string]interface{} `json:"parameters,omitempty" jsonschema:"description=Analysis parameters and configuration"`
	TimeRange    string                 `json:"time_range,omitempty" jsonschema:"description=Time range for data analysis (e.g., last_24h, last_week)"`
}

// PowerPentesterResult represents the result of power industry penetration testing
type PowerPentesterResult struct {
	SystemType       string                 `json:"system_type"`
	TestsExecuted    int                    `json:"tests_executed"`
	Vulnerabilities  int                    `json:"vulnerabilities"`
	CriticalVulns    int                    `json:"critical_vulns"`
	HighVulns        int                    `json:"high_vulns"`
	Report           string                 `json:"report"`
	Recommendations  []string               `json:"recommendations"`
	ComplianceImpact map[string]interface{} `json:"compliance_impact"`
}

// APITesterResult represents the result of API security testing
type APITesterResult struct {
	SystemType        string                 `json:"system_type"`
	EndpointsTested   int                    `json:"endpoints_tested"`
	TestsExecuted     int                    `json:"tests_executed"`
	Vulnerabilities   int                    `json:"vulnerabilities"`
	CriticalVulns     int                    `json:"critical_vulns"`
	HighVulns         int                    `json:"high_vulns"`
	Report            string                 `json:"report"`
	VulnerabilityList []string               `json:"vulnerability_list"`
	Recommendations   []string               `json:"recommendations"`
}

// BizLogicTesterResult represents the result of business logic testing
type BizLogicTesterResult struct {
	SystemType        string                 `json:"system_type"`
	TestsExecuted     int                    `json:"tests_executed"`
	Vulnerabilities   int                    `json:"vulnerabilities"`
	CriticalVulns     int                    `json:"critical_vulns"`
	HighVulns         int                    `json:"high_vulns"`
	FinancialImpact   string                 `json:"financial_impact"`
	Report            string                 `json:"report"`
	BusinessRules     []string               `json:"business_rules_tested"`
	Recommendations   []string               `json:"recommendations"`
}

// ComplianceAgentResult represents the result of compliance testing
type ComplianceAgentResult struct {
	SystemType        string                 `json:"system_type"`
	ComplianceScore   float64                `json:"compliance_score"`
	FrameworksTested  []string               `json:"frameworks_tested"`
	StandardsTested   []string               `json:"standards_tested"`
	RequirementsTested int                   `json:"requirements_tested"`
	Findings          int                    `json:"findings"`
	CriticalFindings  int                    `json:"critical_findings"`
	Report            string                 `json:"report"`
	RegulatoryImpact  map[string]interface{} `json:"regulatory_impact"`
	Recommendations   []string               `json:"recommendations"`
}

// BillingLogicTestResult represents the result of billing logic testing
type BillingLogicTestResult struct {
	Endpoint             string   `json:"endpoint"`
	TotalTests           int      `json:"total_tests"`
	CompletedTests       int      `json:"completed_tests"`
	VulnerabilitiesFound int      `json:"vulnerabilities_found"`
	Report               string   `json:"report"`
	BusinessImpact       string   `json:"business_impact"`
	Recommendations      []string `json:"recommendations"`
}

// APIFuzzerResult represents the result of API fuzzing
type APIFuzzerResult struct {
	BaseURL              string   `json:"base_url"`
	EndpointsTested      int      `json:"endpoints_tested"`
	TotalRequests        int      `json:"total_requests"`
	VulnerabilitiesFound int      `json:"vulnerabilities_found"`
	CriticalVulns        int      `json:"critical_vulns"`
	HighVulns            int      `json:"high_vulns"`
	Report               string   `json:"report"`
	Recommendations      []string `json:"recommendations"`
}

// SAPScannerResult represents the result of SAP security scanning
type SAPScannerResult struct {
	Target              string   `json:"target"`
	VulnerabilitiesFound int     `json:"vulnerabilities_found"`
	DefaultCredsFound   bool     `json:"default_creds_found"`
	ConfigIssues        int      `json:"config_issues"`
	Report              string   `json:"report"`
	Recommendations     []string `json:"recommendations"`
}

// MobileSecurityResult represents the result of mobile security testing
type MobileSecurityResult struct {
	AppType             string   `json:"app_type"`
	TestAreasCompleted  []string `json:"test_areas_completed"`
	VulnerabilitiesFound int     `json:"vulnerabilities_found"`
	CriticalVulns       int      `json:"critical_vulns"`
	Report              string   `json:"report"`
	Recommendations     []string `json:"recommendations"`
}

// PowerDataAnalyzerResult represents the result of power data analysis
type PowerDataAnalyzerResult struct {
	DataSource       string                 `json:"data_source"`
	DataType         string                 `json:"data_type"`
	AnalysisType     string                 `json:"analysis_type"`
	AnomaliesFound   int                    `json:"anomalies_found"`
	PrivacyRisks     int                    `json:"privacy_risks"`
	ComplianceIssues int                    `json:"compliance_issues"`
	Report           string                 `json:"report"`
	Insights         map[string]interface{} `json:"insights"`
	Recommendations  []string               `json:"recommendations"`
}
