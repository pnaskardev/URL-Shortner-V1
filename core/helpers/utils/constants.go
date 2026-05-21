package utils

type SERVICES string

const (
	SHORTENER_SERVICE SERVICES = "shortener_service"
)

var GlobalConstants map[SERVICES]string = map[SERVICES]string{
	SHORTENER_SERVICE: "http://localhost:8001/api/shorten",
}
