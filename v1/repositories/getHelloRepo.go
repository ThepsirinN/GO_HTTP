package repositoriesV1

import "context"

func (r *repoV1) GetHelloRepo(ctx context.Context, str string) string {
	return str
}
