package embedded

import (
	"embed"
)

//go:embed Dockerfile config.json seeds migrations
var EmbeddedFiles embed.FS
