package models

import "github.com/sujamess/fastgql/integration/remote_api"

type Viewer struct {
	User *remote_api.User
}
