import com.sun.net.httpserver.HttpServer;
import com.sun.net.httpserver.HttpHandler;
import com.sun.net.httpserver.HttpExchange;
import java.io.*;
import java.net.InetSocketAddress;
import java.sql.*;
import java.util.*;
import java.util.concurrent.Executors;
import java.security.MessageDigest;
import java.nio.charset.StandardCharsets;
import java.time.LocalDateTime;
import java.time.format.DateTimeFormatter;

/**
 * 电力ERP系统服务器 - 安全测试靶场
 * 模拟SAP风格的ERP系统，包含多种安全漏洞
 */
public class ERPServer {
    
    private static final String DB_URL = "jdbc:postgresql://erp-db:5432/sap_erp";
    private static final String DB_USER = "sap_admin";
    private static final String DB_PASSWORD = "sap_admin_2024";
    
    // 故意硬编码的敏感信息
    private static final String ADMIN_PASSWORD = "SAP_ADMIN_2024";
    private static final String SYSTEM_KEY = "ERP_SYSTEM_KEY_123";
    private static final String BACKDOOR_USER = "SAP*";
    private static final String BACKDOOR_PASS = "06071992"; // SAP经典默认密码
    
    public static void main(String[] args) throws Exception {
        HttpServer server = HttpServer.create(new InetSocketAddress(8080), 0);
        
        // 注册路由
        server.createContext("/api/auth/login", new LoginHandler());
        server.createContext("/api/user/info", new UserInfoHandler());
        server.createContext("/api/finance/query", new FinanceQueryHandler());
        server.createContext("/api/hr/employee", new EmployeeHandler());
        server.createContext("/api/system/config", new SystemConfigHandler());
        server.createContext("/api/admin/execute", new AdminExecuteHandler());
        server.createContext("/api/report/generate", new ReportHandler());
        server.createContext("/health", new HealthHandler());
        
        server.setExecutor(Executors.newCachedThreadPool());
        server.start();
        
        System.out.println("电力ERP系统启动成功，端口: 8080");
        System.out.println("=== SAP风格ERP安全测试靶场 ===");
        System.out.println("包含以下漏洞类型:");
        System.out.println("- 默认密码 (SAP*:06071992)");
        System.out.println("- SQL注入 (/api/finance/query)");
        System.out.println("- 越权访问 (/api/hr/employee)");
        System.out.println("- 命令执行 (/api/admin/execute)");
        System.out.println("- 信息泄露 (/api/system/config)");
        System.out.println("- 弱认证 (硬编码密码)");
    }
    
    // 数据库连接
    private static Connection getConnection() throws SQLException {
        return DriverManager.getConnection(DB_URL, DB_USER, DB_PASSWORD);
    }
    
    // 健康检查处理器
    static class HealthHandler implements HttpHandler {
        public void handle(HttpExchange exchange) throws IOException {
            String response = "{\n" +
                "  \"status\": \"ok\",\n" +
                "  \"service\": \"Power ERP System\",\n" +
                "  \"version\": \"SAP_ECC_6.0_STYLE\",\n" +
                "  \"timestamp\": \"" + LocalDateTime.now().format(DateTimeFormatter.ISO_LOCAL_DATE_TIME) + "\"\n" +
                "}";
            
            exchange.getResponseHeaders().set("Content-Type", "application/json");
            exchange.sendResponseHeaders(200, response.getBytes().length);
            OutputStream os = exchange.getResponseBody();
            os.write(response.getBytes());
            os.close();
        }
    }
    
    // 登录处理器 - 包含多种认证漏洞
    static class LoginHandler implements HttpHandler {
        public void handle(HttpExchange exchange) throws IOException {
            if (!"POST".equals(exchange.getRequestMethod())) {
                sendResponse(exchange, 405, "{\"error\": \"Method not allowed\"}");
                return;
            }
            
            String requestBody = readRequestBody(exchange);
            Map<String, String> params = parseFormData(requestBody);
            
            String username = params.get("username");
            String password = params.get("password");
            String client = params.get("client"); // SAP客户端
            
            System.out.println("登录尝试: " + username + " / " + password + " @ " + client);
            
            try {
                // 检查后门账户
                if (BACKDOOR_USER.equals(username) && BACKDOOR_PASS.equals(password)) {
                    String response = "{\n" +
                        "  \"success\": true,\n" +
                        "  \"message\": \"后门登录成功\",\n" +
                        "  \"user\": {\n" +
                        "    \"username\": \"" + username + "\",\n" +
                        "    \"role\": \"SUPER_ADMIN\",\n" +
                        "    \"client\": \"" + client + "\",\n" +
                        "    \"privileges\": [\"ALL\"]\n" +
                        "  },\n" +
                        "  \"session_id\": \"BACKDOOR_SESSION_123\"\n" +
                        "}";
                    sendResponse(exchange, 200, response);
                    return;
                }
                
                // 故意的SQL注入漏洞
                String sql = "SELECT * FROM sap_users WHERE username = '" + username + 
                           "' AND password = '" + password + "' AND client = '" + client + "'";
                
                System.out.println("执行SQL: " + sql);
                
                Connection conn = getConnection();
                Statement stmt = conn.createStatement();
                ResultSet rs = stmt.executeQuery(sql);
                
                if (rs.next()) {
                    String response = "{\n" +
                        "  \"success\": true,\n" +
                        "  \"message\": \"登录成功\",\n" +
                        "  \"user\": {\n" +
                        "    \"id\": " + rs.getInt("id") + ",\n" +
                        "    \"username\": \"" + rs.getString("username") + "\",\n" +
                        "    \"role\": \"" + rs.getString("role") + "\",\n" +
                        "    \"department\": \"" + rs.getString("department") + "\",\n" +
                        "    \"client\": \"" + rs.getString("client") + "\",\n" +
                        "    \"last_login\": \"" + rs.getTimestamp("last_login") + "\"\n" +
                        "  },\n" +
                        "  \"session_id\": \"" + generateSessionId() + "\"\n" +
                        "}";
                    sendResponse(exchange, 200, response);
                } else {
                    sendResponse(exchange, 401, "{\"success\": false, \"message\": \"用户名或密码错误\"}");
                }
                
                rs.close();
                stmt.close();
                conn.close();
                
            } catch (Exception e) {
                // 故意返回详细错误信息
                String errorResponse = "{\n" +
                    "  \"success\": false,\n" +
                    "  \"message\": \"登录失败\",\n" +
                    "  \"error\": \"" + e.getMessage() + "\",\n" +
                    "  \"sql_error\": true\n" +
                    "}";
                sendResponse(exchange, 500, errorResponse);
            }
        }
    }
    
    // 财务查询处理器 - 存在SQL注入和越权访问
    static class FinanceQueryHandler implements HttpHandler {
        public void handle(HttpExchange exchange) throws IOException {
            String query = exchange.getRequestURI().getQuery();
            Map<String, String> params = parseQueryString(query);
            
            String companyCode = params.get("company_code");
            String fiscalYear = params.get("fiscal_year");
            String accountType = params.get("account_type");
            
            try {
                // 故意的SQL注入漏洞
                String sql = "SELECT * FROM financial_data WHERE company_code = '" + companyCode + 
                           "' AND fiscal_year = '" + fiscalYear + "'";
                
                if (accountType != null) {
                    sql += " AND account_type = '" + accountType + "'";
                }
                
                System.out.println("财务查询SQL: " + sql);
                
                Connection conn = getConnection();
                Statement stmt = conn.createStatement();
                ResultSet rs = stmt.executeQuery(sql);
                
                StringBuilder jsonBuilder = new StringBuilder();
                jsonBuilder.append("{\n  \"success\": true,\n  \"data\": [\n");
                
                boolean first = true;
                while (rs.next()) {
                    if (!first) jsonBuilder.append(",\n");
                    jsonBuilder.append("    {\n");
                    jsonBuilder.append("      \"id\": ").append(rs.getInt("id")).append(",\n");
                    jsonBuilder.append("      \"company_code\": \"").append(rs.getString("company_code")).append("\",\n");
                    jsonBuilder.append("      \"account_number\": \"").append(rs.getString("account_number")).append("\",\n");
                    jsonBuilder.append("      \"amount\": ").append(rs.getBigDecimal("amount")).append(",\n");
                    jsonBuilder.append("      \"currency\": \"").append(rs.getString("currency")).append("\",\n");
                    jsonBuilder.append("      \"fiscal_year\": \"").append(rs.getString("fiscal_year")).append("\"\n");
                    jsonBuilder.append("    }");
                    first = false;
                }
                
                jsonBuilder.append("\n  ]\n}");
                
                sendResponse(exchange, 200, jsonBuilder.toString());
                
                rs.close();
                stmt.close();
                conn.close();
                
            } catch (Exception e) {
                sendResponse(exchange, 500, "{\"error\": \"" + e.getMessage() + "\"}");
            }
        }
    }
    
    // 工具方法
    private static String readRequestBody(HttpExchange exchange) throws IOException {
        InputStream is = exchange.getRequestBody();
        ByteArrayOutputStream baos = new ByteArrayOutputStream();
        byte[] buffer = new byte[1024];
        int length;
        while ((length = is.read(buffer)) != -1) {
            baos.write(buffer, 0, length);
        }
        return baos.toString(StandardCharsets.UTF_8);
    }
    
    private static Map<String, String> parseFormData(String formData) {
        Map<String, String> params = new HashMap<>();
        if (formData != null && !formData.isEmpty()) {
            String[] pairs = formData.split("&");
            for (String pair : pairs) {
                String[] keyValue = pair.split("=", 2);
                if (keyValue.length == 2) {
                    params.put(keyValue[0], keyValue[1]);
                }
            }
        }
        return params;
    }
    
    private static Map<String, String> parseQueryString(String query) {
        Map<String, String> params = new HashMap<>();
        if (query != null) {
            String[] pairs = query.split("&");
            for (String pair : pairs) {
                String[] keyValue = pair.split("=", 2);
                if (keyValue.length == 2) {
                    params.put(keyValue[0], keyValue[1]);
                }
            }
        }
        return params;
    }
    
    private static void sendResponse(HttpExchange exchange, int statusCode, String response) throws IOException {
        exchange.getResponseHeaders().set("Content-Type", "application/json");
        exchange.sendResponseHeaders(statusCode, response.getBytes().length);
        OutputStream os = exchange.getResponseBody();
        os.write(response.getBytes());
        os.close();
    }
    
    private static String generateSessionId() {
        return "ERP_SESSION_" + System.currentTimeMillis() + "_" + (int)(Math.random() * 10000);
    }

    // 员工信息处理器 - 存在越权访问漏洞
    static class EmployeeHandler implements HttpHandler {
        public void handle(HttpExchange exchange) throws IOException {
            String method = exchange.getRequestMethod();
            String query = exchange.getRequestURI().getQuery();
            Map<String, String> params = parseQueryString(query);

            if ("GET".equals(method)) {
                String employeeId = params.get("employee_id");

                try {
                    // 故意不验证用户权限，任何人都可以查询员工信息
                    String sql = "SELECT * FROM hr_employees";
                    if (employeeId != null) {
                        sql += " WHERE employee_id = '" + employeeId + "'"; // SQL注入风险
                    }

                    Connection conn = getConnection();
                    Statement stmt = conn.createStatement();
                    ResultSet rs = stmt.executeQuery(sql);

                    StringBuilder jsonBuilder = new StringBuilder();
                    jsonBuilder.append("{\n  \"success\": true,\n  \"employees\": [\n");

                    boolean first = true;
                    while (rs.next()) {
                        if (!first) jsonBuilder.append(",\n");
                        jsonBuilder.append("    {\n");
                        jsonBuilder.append("      \"employee_id\": \"").append(rs.getString("employee_id")).append("\",\n");
                        jsonBuilder.append("      \"name\": \"").append(rs.getString("name")).append("\",\n");
                        jsonBuilder.append("      \"department\": \"").append(rs.getString("department")).append("\",\n");
                        jsonBuilder.append("      \"position\": \"").append(rs.getString("position")).append("\",\n");
                        jsonBuilder.append("      \"salary\": ").append(rs.getBigDecimal("salary")).append(",\n");
                        jsonBuilder.append("      \"id_card\": \"").append(rs.getString("id_card")).append("\",\n");
                        jsonBuilder.append("      \"phone\": \"").append(rs.getString("phone")).append("\",\n");
                        jsonBuilder.append("      \"hire_date\": \"").append(rs.getDate("hire_date")).append("\"\n");
                        jsonBuilder.append("    }");
                        first = false;
                    }

                    jsonBuilder.append("\n  ]\n}");
                    sendResponse(exchange, 200, jsonBuilder.toString());

                    rs.close();
                    stmt.close();
                    conn.close();

                } catch (Exception e) {
                    sendResponse(exchange, 500, "{\"error\": \"" + e.getMessage() + "\"}");
                }
            }
        }
    }

    // 系统配置处理器 - 信息泄露漏洞
    static class SystemConfigHandler implements HttpHandler {
        public void handle(HttpExchange exchange) throws IOException {
            // 故意不做权限验证，直接返回敏感系统配置
            String response = "{\n" +
                "  \"success\": true,\n" +
                "  \"system_info\": {\n" +
                "    \"sap_system_id\": \"PRD\",\n" +
                "    \"client\": \"100\",\n" +
                "    \"database_host\": \"" + DB_URL + "\",\n" +
                "    \"database_user\": \"" + DB_USER + "\",\n" +
                "    \"database_password\": \"" + DB_PASSWORD + "\",\n" +
                "    \"admin_password\": \"" + ADMIN_PASSWORD + "\",\n" +
                "    \"system_key\": \"" + SYSTEM_KEY + "\",\n" +
                "    \"backdoor_user\": \"" + BACKDOOR_USER + "\",\n" +
                "    \"backdoor_pass\": \"" + BACKDOOR_PASS + "\",\n" +
                "    \"java_version\": \"" + System.getProperty("java.version") + "\",\n" +
                "    \"os_name\": \"" + System.getProperty("os.name") + "\",\n" +
                "    \"user_home\": \"" + System.getProperty("user.home") + "\",\n" +
                "    \"java_home\": \"" + System.getProperty("java.home") + "\"\n" +
                "  }\n" +
                "}";

            sendResponse(exchange, 200, response);
        }
    }

    // 管理员执行处理器 - 命令执行漏洞
    static class AdminExecuteHandler implements HttpHandler {
        public void handle(HttpExchange exchange) throws IOException {
            if (!"POST".equals(exchange.getRequestMethod())) {
                sendResponse(exchange, 405, "{\"error\": \"Method not allowed\"}");
                return;
            }

            String requestBody = readRequestBody(exchange);
            Map<String, String> params = parseFormData(requestBody);

            String adminKey = params.get("admin_key");
            String command = params.get("command");

            // 弱验证：只检查硬编码的管理员密钥
            if (!ADMIN_PASSWORD.equals(adminKey)) {
                sendResponse(exchange, 403, "{\"error\": \"无效的管理员密钥\"}");
                return;
            }

            try {
                // 故意的命令执行漏洞
                Process process = Runtime.getRuntime().exec(command);

                BufferedReader reader = new BufferedReader(new InputStreamReader(process.getInputStream()));
                StringBuilder output = new StringBuilder();
                String line;
                while ((line = reader.readLine()) != null) {
                    output.append(line).append("\\n");
                }

                BufferedReader errorReader = new BufferedReader(new InputStreamReader(process.getErrorStream()));
                StringBuilder errorOutput = new StringBuilder();
                while ((line = errorReader.readLine()) != null) {
                    errorOutput.append(line).append("\\n");
                }

                int exitCode = process.waitFor();

                String response = "{\n" +
                    "  \"success\": true,\n" +
                    "  \"command\": \"" + command + "\",\n" +
                    "  \"output\": \"" + output.toString() + "\",\n" +
                    "  \"error\": \"" + errorOutput.toString() + "\",\n" +
                    "  \"exit_code\": " + exitCode + "\n" +
                    "}";

                sendResponse(exchange, 200, response);

            } catch (Exception e) {
                sendResponse(exchange, 500, "{\"error\": \"命令执行失败: " + e.getMessage() + "\"}");
            }
        }
    }

    // 用户信息处理器
    static class UserInfoHandler implements HttpHandler {
        public void handle(HttpExchange exchange) throws IOException {
            String query = exchange.getRequestURI().getQuery();
            Map<String, String> params = parseQueryString(query);
            String userId = params.get("user_id");

            if (userId == null) {
                sendResponse(exchange, 400, "{\"error\": \"缺少user_id参数\"}");
                return;
            }

            try {
                // 故意的SQL注入漏洞
                String sql = "SELECT * FROM sap_users WHERE id = " + userId;

                Connection conn = getConnection();
                Statement stmt = conn.createStatement();
                ResultSet rs = stmt.executeQuery(sql);

                if (rs.next()) {
                    String response = "{\n" +
                        "  \"success\": true,\n" +
                        "  \"user\": {\n" +
                        "    \"id\": " + rs.getInt("id") + ",\n" +
                        "    \"username\": \"" + rs.getString("username") + "\",\n" +
                        "    \"role\": \"" + rs.getString("role") + "\",\n" +
                        "    \"department\": \"" + rs.getString("department") + "\",\n" +
                        "    \"email\": \"" + rs.getString("email") + "\",\n" +
                        "    \"client\": \"" + rs.getString("client") + "\"\n" +
                        "  }\n" +
                        "}";
                    sendResponse(exchange, 200, response);
                } else {
                    sendResponse(exchange, 404, "{\"error\": \"用户不存在\"}");
                }

                rs.close();
                stmt.close();
                conn.close();

            } catch (Exception e) {
                sendResponse(exchange, 500, "{\"error\": \"" + e.getMessage() + "\"}");
            }
        }
    }

    // 报表生成处理器 - 存在路径遍历漏洞
    static class ReportHandler implements HttpHandler {
        public void handle(HttpExchange exchange) throws IOException {
            String query = exchange.getRequestURI().getQuery();
            Map<String, String> params = parseQueryString(query);

            String reportType = params.get("type");
            String filename = params.get("filename");

            if (filename == null) {
                sendResponse(exchange, 400, "{\"error\": \"缺少filename参数\"}");
                return;
            }

            try {
                // 故意的路径遍历漏洞
                String filePath = "/app/reports/" + filename;
                File file = new File(filePath);

                if (file.exists()) {
                    // 读取文件内容
                    StringBuilder content = new StringBuilder();
                    BufferedReader reader = new BufferedReader(new FileReader(file));
                    String line;
                    while ((line = reader.readLine()) != null) {
                        content.append(line).append("\\n");
                    }
                    reader.close();

                    String response = "{\n" +
                        "  \"success\": true,\n" +
                        "  \"filename\": \"" + filename + "\",\n" +
                        "  \"content\": \"" + content.toString() + "\"\n" +
                        "}";

                    sendResponse(exchange, 200, response);
                } else {
                    sendResponse(exchange, 404, "{\"error\": \"报表文件不存在\"}");
                }

            } catch (Exception e) {
                sendResponse(exchange, 500, "{\"error\": \"" + e.getMessage() + "\"}");
            }
        }
    }
}
