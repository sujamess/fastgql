package models

import "github.com/arsmn/fastgql/integration/remote_api"

type Viewer struct {
	User *remote_api.User
}
