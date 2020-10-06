package buyer

import (
	"time"

	accountpb "github.com/ihtkas/farm/account/v1"
)

// Profile has buyer information about the User
type Profile struct {
	*accountpb.User
}

// Order has request for single product
type Order struct {
	ProductID  uint64
	Quote      uint // TODO valid currency
	Quantity   uint
	RequiredOn time.Time
}
