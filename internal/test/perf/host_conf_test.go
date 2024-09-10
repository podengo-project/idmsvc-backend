package perf

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/podengo-project/idmsvc-backend/internal/api/header"
	"github.com/podengo-project/idmsvc-backend/internal/api/public"
	"github.com/podengo-project/idmsvc-backend/internal/domain/model"
	builder_api "github.com/podengo-project/idmsvc-backend/internal/test/builder/api"
	"github.com/shirou/gopsutil/v4/cpu"
)

type requestError struct {
	StatusCode int
	Err        error
}

type Timing struct {
	Start time.Time
	End   time.Time
}

func (t Timing) Duration() time.Duration {
	return t.End.Sub(t.Start)
}

type HostConfTest struct {
	requestCount       int
	concurrentRequests int
	requestTimeout     time.Duration
}

func TestHostConf(t *testing.T) {
	orgCount := 500

	testCases := []HostConfTest{
		// request count, concurrent requests, request timeout
		{10, 1, 1 * time.Second},
		{100, 1, 1 * time.Second},
		{1000, 1, 1 * time.Second},
		{10, 10, 1 * time.Second},
		{100, 100, 1 * time.Second},
		{1000, 100, 1 * time.Second},
		{2000, 100, 1 * time.Second},
		{4000, 100, 1 * time.Second},
		{8000, 100, 1 * time.Second},
		{16000, 100, 20 * time.Second},
		{16000, 200, 20 * time.Second},
		{16000, 400, 20 * time.Second},
		{16000, 800, 20 * time.Second},
		{16000, 1600, 20 * time.Second},
	}

	cpuInfo, err := cpu.Info()
	if err != nil {
		t.Fatalf("Failed to get CPU info: %v", err)
	}
	for i, info := range cpuInfo {
		t.Logf("CPU %d: %s", i, info.ModelName)
		t.Logf("Cores %d: %d", i, info.Cores)
	}

	domains := prepDomains(t, orgCount)
	processMonitor, err := NewProcessMonitor("bin/service", 500*time.Millisecond, 10*time.Minute)
	if err != nil {
		t.Fatalf("Failed to create process monitor: %v", err)
	}

	for _, tc := range testCases {
		testName := "Count:" + strconv.Itoa(tc.requestCount) + " Concurrent:" + strconv.Itoa(tc.concurrentRequests)
		t.Run(testName, func(t *testing.T) {
			requests := createHostConfRequests(t, tc.requestCount, domains)
			client := &http.Client{
				Timeout: tc.requestTimeout,
			}
			processMonitor.ClearHistory()
			processMonitor.Monitor() // Collect the last state
			processMonitor.Start()
			errors, timing := doRequestsConcurrently(client, requests, tc.concurrentRequests)
			processMonitor.Stop()
			processMonitor.Monitor() // Collect the last state

			t.Logf("Errors:                %d, %f%%", len(errors), float64(len(errors)*100)/float64(tc.requestCount))
			t.Logf("Start:                 %v", timing.Start)
			t.Logf("End:                   %v", timing.End)
			t.Logf("Duration:              %v", timing.Duration())
			t.Logf("Requests per second:   %f", float64(tc.requestCount)/timing.Duration().Seconds())
			t.Logf("Average response time: %v", timing.Duration()/time.Duration(tc.requestCount))
			t.Logf("Max CPU:               %f%%", processMonitor.MaxCPUPercent())
			t.Logf("Max memory:            %s", HumanBytes(processMonitor.MaxMemRSS()))
		})
	}
}

// prepDomains creates domains for the test and returns them.
// The domains are created directly in the database.
func prepDomains(t *testing.T, orgCount int) []*model.Domain {
	prepData := NewPrepData()
	t.Log("Creating domains: ", orgCount)

	if orgCount == 0 {
		orgCount = 500
	}
	domains, err := CreateDomains(*prepData, orgCount, 10000000)
	t.Cleanup(func() {
		delErr := DeleteDomains(*prepData, domains)
		if delErr != nil {
			t.Fatalf("Failed to delete domains: %v", delErr)
		}
	})
	if err != nil {
		t.Fatalf("Failed to create domains: %v", err)
	}
	return domains
}

func getBaseURL() string {
	return "http://localhost:8000/api/idmsvc/v1"
}

func createHostConfRequest(
	baseURL, inventoryID, fqdn string, domain *model.Domain,
) (*http.Request, error) {
	url := baseURL + "/host-conf/" + inventoryID + "/" + fqdn
	identity := builder_api.NewSystemXRHID().WithOrgID(domain.OrgId).Build()
	domainType := public.RhelIdm
	hostConf := builder_api.NewHostConf().
		WithDomainName(domain.DomainName).
		WithDomainType(&domainType).
		Build()
	reqBody, err := json.Marshal(hostConf)
	if err != nil {
		return nil, err
	}
	reader := bytes.NewReader(reqBody)

	req, err := http.NewRequest(http.MethodPost, url, reader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(header.HeaderXRHID, header.EncodeXRHID(&identity))

	return req, nil
}

func createHostConfRequests(t *testing.T, requestCount int, domains []*model.Domain) []*http.Request {
	t.Log("Request count: ", requestCount)
	baseURL := getBaseURL()
	orgCount := len(domains)
	requests := make([]*http.Request, requestCount)
	var err error

	for i := 0; i < requestCount; i++ {
		inventoryID := uuid.NewString()
		domain := domains[i%orgCount]
		fqdn := "test" + strconv.Itoa(i) + "." + *domain.DomainName
		requests[i], err = createHostConfRequest(baseURL, inventoryID, fqdn, domain)
		if err != nil {
			t.Fatalf("Failed to create host-conf request: %v", err)
		}
	}
	return requests
}

func doRequestsConcurrently(
	client *http.Client, requests []*http.Request, concurrentRequests int,
) (errors []*requestError, timing Timing) {
	n := len(requests)
	concurentSem := make(chan int, concurrentRequests)
	errChan := make(chan *requestError, n)
	timing.Start = time.Now()
	for i := 0; i < n; i++ {
		concurentSem <- 1
		go func(i int) {
			resp, reqErr := client.Do(requests[i])
			errChan <- parseRequestResult(resp, reqErr)
			<-concurentSem
		}(i)
	}
	timing.End = time.Now()

	for i := 0; i < n; i++ {
		err := <-errChan
		if err != nil {
			errors = append(errors, err)
		}
	}
	return
}

func parseRequestResult(resp *http.Response, respErr error) *requestError {
	if respErr != nil {
		return &requestError{Err: respErr}
	}
	closeErr := resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return &requestError{StatusCode: resp.StatusCode}
	}
	if closeErr != nil {
		return &requestError{Err: closeErr}
	}
	return nil
}
