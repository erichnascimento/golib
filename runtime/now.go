package runtime

import "time"

var Now func() time.Time = time.Now
