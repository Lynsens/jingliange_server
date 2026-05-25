package model

type OpsAccessLogItem struct {
	Time      string `json:"time"`
	IP        string `json:"ip"`
	Method    string `json:"method"`
	Path      string `json:"path"`
	Protocol  string `json:"protocol"`
	Status    int    `json:"status"`
	Bytes     int64  `json:"bytes"`
	Referer   string `json:"referer"`
	UserAgent string `json:"user_agent"`
	Raw       string `json:"raw"`
}

type OpsTextLogItem struct {
	Time  string `json:"time,omitempty"`
	Level string `json:"level,omitempty"`
	Raw   string `json:"raw"`
}

type OpsPathStat struct {
	Path  string `json:"path"`
	Count int    `json:"count"`
}

type OpsStatusStat struct {
	Status int `json:"status"`
	Count  int `json:"count"`
}

type OpsSummary struct {
	Date          string          `json:"date"`
	TotalRequests int             `json:"total_requests"`
	UniqueIPs     int             `json:"unique_ips"`
	Status4xx     int             `json:"status_4xx"`
	Status5xx     int             `json:"status_5xx"`
	TotalBytes    int64           `json:"total_bytes"`
	TopPaths      []OpsPathStat   `json:"top_paths"`
	StatusCounts  []OpsStatusStat `json:"status_counts"`
	SourceExists  bool            `json:"source_exists"`
}

type OpsAccessLogResponse struct {
	Date         string             `json:"date"`
	Items        []OpsAccessLogItem `json:"items"`
	SourceExists bool               `json:"source_exists"`
}

type OpsTextLogResponse struct {
	Date         string           `json:"date"`
	Items        []OpsTextLogItem `json:"items"`
	SourceExists bool             `json:"source_exists"`
}
