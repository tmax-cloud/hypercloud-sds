package kubectl

import (
	"strings"
)

var (
	errMsgWhenDeletingNonExistingCRD = "no matches for kind"
	errMsgWhenGettingNonExistingCRD  = "the server doesn't have a resource type"
)

// CRDAlreadyExists returns false if the errMsg means that the client operate on cr, even if its crd does not exist,
// otherwise return true
// crd 가 없는데 cr 을 get, delete 하는 경우는 에러로 판단하지 않기 위함
func CRDAlreadyExists(errMsg string) bool {
	return !(strings.Contains(errMsg, errMsgWhenDeletingNonExistingCRD) || // 없는 crd 를 delete 할 때 나오는 에러
		strings.Contains(errMsg, errMsgWhenGettingNonExistingCRD)) // 없는 crd 를 get 할 때 나오는 에러
}
