package admin

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/lynsens/jingliange_server/pkg/app"
	"github.com/lynsens/jingliange_server/pkg/logging"
	"github.com/lynsens/jingliange_server/pkg/setting"
	"github.com/lynsens/jingliange_server/pkg/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupOpsTest(t *testing.T) (string, string) {
	t.Helper()

	gin.SetMode(gin.TestMode)
	tmpDir := t.TempDir()
	wd, err := os.Getwd()
	require.NoError(t, err)
	runtimeRoot, err := filepath.Rel(wd, tmpDir)
	require.NoError(t, err)

	originalApp := *setting.AppSetting
	originalOps := *setting.OpsSetting

	setting.AppSetting.JwtSecret = "test_secret_key_for_ops"
	setting.AppSetting.JwtExpire = 720
	setting.AppSetting.RuntimeRootPath = runtimeRoot + string(os.PathSeparator)
	setting.AppSetting.LogSavePath = ""
	setting.AppSetting.LogSaveName = "app"
	setting.AppSetting.LogFileExt = "log"
	setting.AppSetting.TimeFormat = "20060102"
	setting.OpsSetting.NginxAccessLogPath = filepath.Join(tmpDir, "access.log")
	setting.OpsSetting.NginxErrorLogPath = filepath.Join(tmpDir, "error.log")
	setting.OpsSetting.AppLogDir = tmpDir
	util.InitJWT()
	logging.Setup()

	t.Cleanup(func() {
		logging.Close()
		*setting.AppSetting = originalApp
		*setting.OpsSetting = originalOps
		util.InitJWT()
	})

	adminToken, err := util.GenerateAdminToken("admin")
	require.NoError(t, err)
	userToken, err := util.GenerateToken("user")
	require.NoError(t, err)

	return adminToken, userToken
}

func setupOpsRouter() *gin.Engine {
	router := gin.New()
	adminGroup := router.Group("/api/admin")
	adminGroup.Use(util.AdminJWT())
	{
		adminGroup.GET("/ops/summary", GetOpsSummary)
		adminGroup.GET("/ops/access-logs", GetOpsAccessLogs)
		adminGroup.GET("/ops/error-logs", GetOpsErrorLogs)
		adminGroup.GET("/ops/app-logs", GetOpsAppLogs)
	}
	return router
}

func performOpsRequest(t *testing.T, router *gin.Engine, path string, token string) (int, app.Response) {
	t.Helper()

	req, err := http.NewRequest("GET", path, nil)
	require.NoError(t, err)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var response app.Response
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	return w.Code, response
}

func TestOpsAccessLogParsingAndSummary(t *testing.T) {
	adminToken, _ := setupOpsTest(t)
	router := setupOpsRouter()

	accessLog := `1.1.1.1 - - [25/May/2026:01:00:00 +0800] "GET / HTTP/2.0" 200 120 "-" "Mozilla/5.0"
2.2.2.2 - - [25/May/2026:01:01:00 +0800] "POST /api/v1/menu/getMenu HTTP/2.0" 200 220 "-" "MiniProgram"
1.1.1.1 - - [25/May/2026:01:02:00 +0800] "GET /missing HTTP/2.0" 404 20 "-" "Mozilla/5.0"
3.3.3.3 - - [25/May/2026:01:03:00 +0800] "GET /boom HTTP/2.0" 500 10 "-" "Mozilla/5.0"
4.4.4.4 - - [24/May/2026:23:59:00 +0800] "GET /old HTTP/2.0" 200 5 "-" "Mozilla/5.0"
`
	require.NoError(t, os.WriteFile(setting.OpsSetting.NginxAccessLogPath, []byte(accessLog), 0644))

	statusCode, response := performOpsRequest(t, router, "/api/admin/ops/summary?date=2026-05-25", adminToken)
	require.Equal(t, http.StatusOK, statusCode)
	data := response.Data.(map[string]interface{})
	assert.Equal(t, float64(4), data["total_requests"])
	assert.Equal(t, float64(3), data["unique_ips"])
	assert.Equal(t, float64(1), data["status_4xx"])
	assert.Equal(t, float64(1), data["status_5xx"])
	assert.Equal(t, float64(370), data["total_bytes"])
	assert.Equal(t, true, data["source_exists"])

	statusCode, response = performOpsRequest(t, router, "/api/admin/ops/access-logs?date=2026-05-25&limit=2&keyword=mozilla", adminToken)
	require.Equal(t, http.StatusOK, statusCode)
	data = response.Data.(map[string]interface{})
	items := data["items"].([]interface{})
	require.Len(t, items, 2)
	assert.Equal(t, "/missing", items[0].(map[string]interface{})["path"])
	assert.Equal(t, "/boom", items[1].(map[string]interface{})["path"])
}

func TestOpsTextLogsAndMissingFiles(t *testing.T) {
	adminToken, _ := setupOpsTest(t)
	router := setupOpsRouter()

	appLogPath := filepath.Join(setting.OpsSetting.AppLogDir, "app20260525.log")
	appLog := `2026/05/25 01:00:00 [INFO][menu.go:33] ok
2026/05/25 01:01:00 [ERROR][menu.go:76] failed
2026/05/24 23:59:00 [ERROR][old.go:1] old
`
	require.NoError(t, os.WriteFile(appLogPath, []byte(appLog), 0644))

	statusCode, response := performOpsRequest(t, router, "/api/admin/ops/app-logs?date=2026-05-25&level=ERROR&limit=10", adminToken)
	require.Equal(t, http.StatusOK, statusCode)
	data := response.Data.(map[string]interface{})
	items := data["items"].([]interface{})
	require.Len(t, items, 1)
	assert.Contains(t, items[0].(map[string]interface{})["raw"], "failed")
	assert.Equal(t, true, data["source_exists"])

	statusCode, response = performOpsRequest(t, router, "/api/admin/ops/error-logs?date=2026-05-25", adminToken)
	require.Equal(t, http.StatusOK, statusCode)
	data = response.Data.(map[string]interface{})
	assert.Equal(t, false, data["source_exists"])
	assert.Len(t, data["items"], 0)
}

func TestOpsRoutesRequireAdminJWT(t *testing.T) {
	adminToken, userToken := setupOpsTest(t)
	router := setupOpsRouter()
	require.NoError(t, os.WriteFile(setting.OpsSetting.NginxAccessLogPath, []byte(""), 0644))

	statusCode, _ := performOpsRequest(t, router, "/api/admin/ops/summary?date=2026-05-25", "")
	assert.Equal(t, http.StatusUnauthorized, statusCode)

	statusCode, _ = performOpsRequest(t, router, "/api/admin/ops/summary?date=2026-05-25", "invalid")
	assert.Equal(t, http.StatusUnauthorized, statusCode)

	statusCode, _ = performOpsRequest(t, router, "/api/admin/ops/summary?date=2026-05-25", userToken)
	assert.Equal(t, http.StatusUnauthorized, statusCode)

	statusCode, _ = performOpsRequest(t, router, "/api/admin/ops/summary?date=2026-05-25", adminToken)
	assert.Equal(t, http.StatusOK, statusCode)
}

func TestNormalizeOpsLimit(t *testing.T) {
	assert.Equal(t, defaultOpsLogLimit, normalizeOpsLimit(""))
	assert.Equal(t, 20, normalizeOpsLimit("20"))
	assert.Equal(t, maxOpsLogLimit, normalizeOpsLimit("9999"))
	assert.Equal(t, defaultOpsLogLimit, normalizeOpsLimit("-1"))
}
