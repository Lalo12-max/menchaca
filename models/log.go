package models

import "time"

type Log struct {
	IDLog        int       `json:"id_log"`
	Method       string    `json:"method"`
	Path         string    `json:"path"`
	Protocol     string    `json:"protocol"`
	StatusCode   int       `json:"status_code"`
	ResponseTime int       `json:"response_time"` 
	UserAgent    *string   `json:"user_agent"`
	IP           string    `json:"ip"`
	Hostname     string    `json:"hostname"`
	Body         *string   `json:"body"`
	Params       *string   `json:"params"`
	Query        *string   `json:"query"`
	Email        *string   `json:"email"`        
	Username     *string   `json:"username"`     
	Role         *string   `json:"role"`         
	LogLevel     string    `json:"log_level"`    
	Environment  string    `json:"environment"`  
	NodeVersion  string    `json:"node_version"` 
	PID          int       `json:"pid"`          
	Timestamp    time.Time `json:"timestamp"`
	URL          string    `json:"url"` 
	CreatedAt    time.Time `json:"created_at"`
}
