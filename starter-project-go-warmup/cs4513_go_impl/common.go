// Package cs4513_go_impl performs top-k word count and parallel sum
// by convention, the package name should match the last segment of its import path (see ./cs4513_go_test/q1_test.go)
package cs4513_go_impl

import "log"

// Propagate error if it exists
func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
