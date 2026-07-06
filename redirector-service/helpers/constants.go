package constants

import "time"

const (
	URL_EXCHANGE              = "url_exchange"
	URL_CREATED_QUEUE         = "url.created"
	URL_ANALYTICS_EVENT_QUEUE = "url.analytics"

	// TTL CONSTANT
	UrlTTlConstant = 365 * 5 * 25 * time.Hour
)
