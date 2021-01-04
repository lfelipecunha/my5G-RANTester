package control_test_engine

import "fmt"

// generated a IMSI from integer.
func ImsiGenerator(i int) string {
    return fmt.Sprintf("imsi-20893%08d", i)
}
