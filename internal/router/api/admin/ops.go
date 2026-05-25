package admin

import (
	"bufio"
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lynsens/jingliange_server/internal/model"
	"github.com/lynsens/jingliange_server/pkg/app"
	"github.com/lynsens/jingliange_server/pkg/e"
	"github.com/lynsens/jingliange_server/pkg/setting"
)

const (
	defaultOpsLogLimit = 100
	maxOpsLogLimit     = 500
)

var nginxAccessPattern = regexp.MustCompile(`^(\S+) \S+ \S+ \[([^\]]+)\] "([^"]*)" (\d{3}) (\S+) "([^"]*)" "([^"]*)"`)

type parsedAccessLog struct {
	item model.OpsAccessLogItem
	date string
	ok   bool
}

// @Summary 管理员获取运维流量概览
// @Description 只读读取配置指定的 Nginx access log，统计指定日期的访问量。
// @Tags AdminOps
// @Produce json
// @Param date query string false "日期，格式 YYYY-MM-DD"
// @Success 200 {object} app.Response{data=model.OpsSummary}
// @Router /api/admin/ops/summary [get]
func GetOpsSummary(c *gin.Context) {
	appG := app.Gin{C: c}
	date := normalizeOpsDate(c.Query("date"))

	records, exists, err := readAccessLogs(setting.OpsSetting.NginxAccessLogPath, date, 0, "")
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	summary := summarizeAccessLogs(date, records, exists)
	appG.Response(http.StatusOK, e.SUCCESS, summary)
}

// @Summary 管理员获取 Nginx 访问日志
// @Description 只读读取配置指定的 Nginx access log。
// @Tags AdminOps
// @Produce json
// @Param date query string false "日期，格式 YYYY-MM-DD"
// @Param limit query int false "最大返回条数，最多 500"
// @Param keyword query string false "关键词"
// @Success 200 {object} app.Response{data=model.OpsAccessLogResponse}
// @Router /api/admin/ops/access-logs [get]
func GetOpsAccessLogs(c *gin.Context) {
	appG := app.Gin{C: c}
	date := normalizeOpsDate(c.Query("date"))
	limit := normalizeOpsLimit(c.Query("limit"))

	records, exists, err := readAccessLogs(setting.OpsSetting.NginxAccessLogPath, date, limit, c.Query("keyword"))
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, model.OpsAccessLogResponse{
		Date:         date,
		Items:        records,
		SourceExists: exists,
	})
}

// @Summary 管理员获取 Nginx 错误日志
// @Description 只读读取配置指定的 Nginx error log。
// @Tags AdminOps
// @Produce json
// @Param date query string false "日期，格式 YYYY-MM-DD"
// @Param limit query int false "最大返回条数，最多 500"
// @Param keyword query string false "关键词"
// @Success 200 {object} app.Response{data=model.OpsTextLogResponse}
// @Router /api/admin/ops/error-logs [get]
func GetOpsErrorLogs(c *gin.Context) {
	appG := app.Gin{C: c}
	date := normalizeOpsDate(c.Query("date"))
	limit := normalizeOpsLimit(c.Query("limit"))

	items, exists, err := readTextLogs(setting.OpsSetting.NginxErrorLogPath, date, limit, "", c.Query("keyword"))
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, model.OpsTextLogResponse{
		Date:         date,
		Items:        items,
		SourceExists: exists,
	})
}

// @Summary 管理员获取后端应用日志
// @Description 只读读取配置指定的后端应用日志目录。
// @Tags AdminOps
// @Produce json
// @Param date query string false "日期，格式 YYYY-MM-DD"
// @Param limit query int false "最大返回条数，最多 500"
// @Param level query string false "日志级别，例如 INFO、WARN、ERROR"
// @Param keyword query string false "关键词"
// @Success 200 {object} app.Response{data=model.OpsTextLogResponse}
// @Router /api/admin/ops/app-logs [get]
func GetOpsAppLogs(c *gin.Context) {
	appG := app.Gin{C: c}
	date := normalizeOpsDate(c.Query("date"))
	limit := normalizeOpsLimit(c.Query("limit"))
	path := getAppLogPath(date)

	items, exists, err := readTextLogs(path, date, limit, c.Query("level"), c.Query("keyword"))
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, model.OpsTextLogResponse{
		Date:         date,
		Items:        items,
		SourceExists: exists,
	})
}

func normalizeOpsDate(value string) string {
	if _, err := time.Parse("2006-01-02", value); err == nil {
		return value
	}
	return time.Now().Format("2006-01-02")
}

func normalizeOpsLimit(value string) int {
	limit, err := strconv.Atoi(value)
	if err != nil || limit <= 0 {
		return defaultOpsLogLimit
	}
	if limit > maxOpsLogLimit {
		return maxOpsLogLimit
	}
	return limit
}

func readAccessLogs(path, date string, limit int, keyword string) ([]model.OpsAccessLogItem, bool, error) {
	file, err := os.Open(filepath.Clean(path))
	if errors.Is(err, os.ErrNotExist) {
		return []model.OpsAccessLogItem{}, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	defer file.Close()

	keyword = strings.ToLower(strings.TrimSpace(keyword))
	items := make([]model.OpsAccessLogItem, 0)

	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for scanner.Scan() {
		line := scanner.Text()
		if keyword != "" && !strings.Contains(strings.ToLower(line), keyword) {
			continue
		}

		parsed := parseNginxAccessLog(line)
		if !parsed.ok || parsed.date != date {
			continue
		}

		items = appendWithLimit(items, parsed.item, limit)
	}

	if err := scanner.Err(); err != nil {
		return nil, true, err
	}

	return items, true, nil
}

func parseNginxAccessLog(line string) parsedAccessLog {
	matches := nginxAccessPattern.FindStringSubmatch(line)
	if len(matches) != 8 {
		return parsedAccessLog{}
	}

	accessTime, err := time.Parse("02/Jan/2006:15:04:05 -0700", matches[2])
	if err != nil {
		return parsedAccessLog{}
	}

	method, requestPath, protocol := parseNginxRequest(matches[3])
	status, _ := strconv.Atoi(matches[4])
	bytes, _ := strconv.ParseInt(strings.ReplaceAll(matches[5], "-", "0"), 10, 64)

	return parsedAccessLog{
		item: model.OpsAccessLogItem{
			Time:      accessTime.Format(time.RFC3339),
			IP:        matches[1],
			Method:    method,
			Path:      requestPath,
			Protocol:  protocol,
			Status:    status,
			Bytes:     bytes,
			Referer:   matches[6],
			UserAgent: matches[7],
			Raw:       line,
		},
		date: accessTime.Format("2006-01-02"),
		ok:   true,
	}
}

func parseNginxRequest(value string) (string, string, string) {
	parts := strings.Fields(value)
	if len(parts) >= 3 {
		return parts[0], parts[1], parts[2]
	}
	if len(parts) == 2 {
		return parts[0], parts[1], ""
	}
	return "", value, ""
}

func readTextLogs(path, date string, limit int, level string, keyword string) ([]model.OpsTextLogItem, bool, error) {
	file, err := os.Open(filepath.Clean(path))
	if errors.Is(err, os.ErrNotExist) {
		return []model.OpsTextLogItem{}, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	defer file.Close()

	level = strings.ToUpper(strings.TrimSpace(level))
	keyword = strings.ToLower(strings.TrimSpace(keyword))
	items := make([]model.OpsTextLogItem, 0)

	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for scanner.Scan() {
		line := scanner.Text()
		if !textLogMatchesDate(line, date) {
			continue
		}
		if level != "" && !strings.Contains(strings.ToUpper(line), "["+level+"]") && !strings.Contains(strings.ToUpper(line), level) {
			continue
		}
		if keyword != "" && !strings.Contains(strings.ToLower(line), keyword) {
			continue
		}

		items = appendWithLimit(items, model.OpsTextLogItem{
			Time:  extractTextLogTime(line),
			Level: extractTextLogLevel(line),
			Raw:   line,
		}, limit)
	}

	if err := scanner.Err(); err != nil {
		return nil, true, err
	}

	return items, true, nil
}

func textLogMatchesDate(line, date string) bool {
	if strings.Contains(line, date) {
		return true
	}
	return strings.Contains(line, strings.ReplaceAll(date, "-", "/"))
}

func extractTextLogTime(line string) string {
	candidates := []string{
		"2006/01/02 15:04:05",
		"2006-01-02 15:04:05",
	}
	for _, layout := range candidates {
		if len(line) < len(layout) {
			continue
		}
		raw := line[:len(layout)]
		if t, err := time.ParseInLocation(layout, raw, time.Local); err == nil {
			return t.Format(time.RFC3339)
		}
	}
	return ""
}

func extractTextLogLevel(line string) string {
	for _, level := range []string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"} {
		if strings.Contains(strings.ToUpper(line), "["+level+"]") {
			return level
		}
	}
	return ""
}

func summarizeAccessLogs(date string, items []model.OpsAccessLogItem, exists bool) model.OpsSummary {
	uniqueIPs := make(map[string]struct{})
	pathCounts := make(map[string]int)
	statusCounts := make(map[int]int)
	var totalBytes int64
	status4xx := 0
	status5xx := 0

	for _, item := range items {
		uniqueIPs[item.IP] = struct{}{}
		pathCounts[item.Path]++
		statusCounts[item.Status]++
		totalBytes += item.Bytes
		if item.Status >= 400 && item.Status < 500 {
			status4xx++
		}
		if item.Status >= 500 && item.Status < 600 {
			status5xx++
		}
	}

	return model.OpsSummary{
		Date:          date,
		TotalRequests: len(items),
		UniqueIPs:     len(uniqueIPs),
		Status4xx:     status4xx,
		Status5xx:     status5xx,
		TotalBytes:    totalBytes,
		TopPaths:      topPathStats(pathCounts, 8),
		StatusCounts:  statusStats(statusCounts),
		SourceExists:  exists,
	}
}

func topPathStats(counts map[string]int, limit int) []model.OpsPathStat {
	stats := make([]model.OpsPathStat, 0, len(counts))
	for path, count := range counts {
		stats = append(stats, model.OpsPathStat{Path: path, Count: count})
	}
	sort.Slice(stats, func(i, j int) bool {
		if stats[i].Count == stats[j].Count {
			return stats[i].Path < stats[j].Path
		}
		return stats[i].Count > stats[j].Count
	})
	if len(stats) > limit {
		return stats[:limit]
	}
	return stats
}

func statusStats(counts map[int]int) []model.OpsStatusStat {
	stats := make([]model.OpsStatusStat, 0, len(counts))
	for status, count := range counts {
		stats = append(stats, model.OpsStatusStat{Status: status, Count: count})
	}
	sort.Slice(stats, func(i, j int) bool {
		return stats[i].Status < stats[j].Status
	})
	return stats
}

func appendWithLimit[T any](items []T, item T, limit int) []T {
	items = append(items, item)
	if limit > 0 && len(items) > limit {
		copy(items[0:], items[1:])
		items = items[:limit]
	}
	return items
}

func getAppLogPath(date string) string {
	appLogDir := setting.OpsSetting.AppLogDir
	if appLogDir == "" {
		appLogDir = filepath.Join(setting.AppSetting.RuntimeRootPath, setting.AppSetting.LogSavePath)
	}

	if !filepath.IsAbs(appLogDir) && setting.OpsSetting.AppLogDir == "" {
		appLogDir = filepath.Clean(appLogDir)
	}

	logTime, err := time.Parse("2006-01-02", date)
	if err != nil {
		logTime = time.Now()
	}

	fileName := setting.AppSetting.LogSaveName + logTime.Format(setting.AppSetting.TimeFormat) + "." + setting.AppSetting.LogFileExt
	return filepath.Join(appLogDir, fileName)
}
