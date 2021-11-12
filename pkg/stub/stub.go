package stub

import (
	"fmt"
	"time"
)

// IntSlowOperation
// Slow operator imitation
func IntSlowOperation(timeDuration time.Duration) int{
	time.Sleep(timeDuration)
	return 111
}

// StrSlowOperation
// Slow operator imitation
func StrSlowOperation(timeDuration time.Duration, strVars...string) string{
	time.Sleep(timeDuration)
	return fmt.Sprint(strVars)
}