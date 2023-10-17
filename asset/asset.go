package asset

import "github.com/goalm/lib/utils"

type asset interface {
	RedempT(start utils.Date) int
}
