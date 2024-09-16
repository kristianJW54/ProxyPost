package internal

import "time"

type ClientConn struct {
	RemoteAddr   string
	UserAgent    string
	Events       []ConnEvents
	RequestCount int
}

type ConnEvents struct {
	Method      string
	URL         string
	Referer     string
	RequestTime time.Duration
	Headers     HeaderInfo
}

type HeaderInfo struct {
	Status         string
	Accept         string
	AcceptEncoding string
	Connection     string
	ContentType    string
	Cookie         string
	Referer        string
	Platform       string
	Other          map[string]string
}
