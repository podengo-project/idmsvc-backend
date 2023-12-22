// Package smoke contains all the smoke test for the service
//
// By smoke test we refers to application requests that get a success
// response for a normal execution of the service, so we check the thinks
// work as expected in normal situations.
//
// If you are thinking a failure and error situations to be covered
// then you will be looking at `internal/test/integration` package
// where you will define the remaining tests to evoque error situations
// and validate they are correctly covered and managed by the application.
package smoke
