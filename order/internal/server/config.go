package server

const (
	DEBUG   = "debug"
	RELEASE = "release"
)

type Config struct {
	Port string
	Mode string
}
