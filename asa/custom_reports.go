package asa

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gocarina/gocsv"
)

type ImpressionShareReportService service

type ImpressionShareReportRequest struct {
	Field     string `json:"field,omitempty"`
	Limit     int    `json:"limit,omitempty"`
	Offset    int    `json:"offset,omitempty"`
	SortOrder string `json:"sortOrder,omitempty"`
}

type ImpressionShareReport struct {
	ID               int64    `json:"id"`
	Name             string   `json:"name"`
	StartTime        string   `json:"startTime"`
	EndTime          string   `json:"endTime"`
	Granularity      string   `json:"granularity"`
	DownloadUri      string   `json:"downloadUri"`
	Dimensions       []string `json:"dimensions"`
	Metrics          []string `json:"metrics"`
	State            string   `json:"state"`
	CreationTime     string   `json:"creationTime"`
	ModificationTime string   `json:"modificationTime"`
	DateRange        string   `json:"dateRange"`
}

// The Impression Share report request body.
// https://developer.apple.com/documentation/apple_search_ads/customreportrequest
type CustomReportRequest struct {
	// The date range of the report request. A date range is required only when using WEEKLY granularity in Impression Share Report.
	// Default: LAST_WEEK
	// Possible Values: LAST_WEEK, LAST_2_WEEKS, LAST_4_WEEKS
	DateRange CustomReportDateRange `json:"dateRange,omitempty"`

	// The end time of the report. The format is YYYY-MM-DD, such as 2022-06-30.
	EndTime Date `json:"endTime,omitempty"`

	// The report data organized by day or week.
	// Impression Share reports with a WEEKLY granularity value can’t have custom startTime and endTime in the request payload.
	// Default: DAILY
	// Possible Values: DAILY, WEEKLY
	Granularity CustomReportGranularity `json:"granularity,omitempty"`

	// (Required) A free-text field. The maximum length is 50 characters.
	Name string `json:"name"`

	// Selector is an optional parameter to filter API results using the countryOrRegion and adamId fields. For countryOrRegion, use an alpha-2 country code value. The IN operator is available to use with Impression Share reports.
	// See SovCondition for selector descriptions and see Selector for structural guidance with selectors.
	Selector Selector `json:"selector,omitempty"`

	// The start time of the report. The format is YYYY-MM-DD, such as 2022-06-01.
	StartTime Date `json:"startTime,omitempty"`
}

type CustomReportDateRange string

const (
	CustomReportDateRangeLastWeek   CustomReportDateRange = "LAST_WEEK"
	CustomReportDateRangeLast2Weeks CustomReportDateRange = "LAST_2_WEEKS"
	CustomReportDateRangeLast4Weeks CustomReportDateRange = "LAST_4_WEEKS"
)

type CustomReportGranularity string

const (
	CustomReportGranularityDaily  CustomReportGranularity = "DAILY"
	CustomReportGranularityWeekly CustomReportGranularity = "WEEKLY"
)

type ImpressionShareReportsResponse struct {
	Data       []ImpressionShareReport `json:"data"`
	Pagination *PageDetail             `json:"pagination"`
	Error      *ErrorResponseBody      `json:"error"`
}

type ImpressionShareReportResponse struct {
	Data       *ImpressionShareReport `json:"data"`
	Pagination *PageDetail            `json:"pagination"`
	Error      *ErrorResponseBody     `json:"error"`
}

type DailyImpressionShareReport struct {
	Records []DailyImpressionShareReportRecord `json:"records"`
}

type DailyImpressionShareReportRecord struct {
	Date                Date    `csv:"date"`
	AppName             string  `csv:"appName"`
	AdamId              int64   `csv:"adamId"`
	CountryOrRegion     string  `csv:"countryOrRegion"`
	SearchTerm          string  `csv:"searchTerm"`
	LowImpressionShare  float64 `csv:"lowImpressionShare"`
	HighImpressionShare float64 `csv:"highImpressionShare"`
	Rank                string  `csv:"rank"`
	SearchPopularity    int     `csv:"searchPopularity"`
}

// Use this endpoint to obtain a reportId to use in a Get a Single Impression Share Report request. This endopoint supports selectors. See CustomReportRequest for selector structure.
// You can generate up to 10 reports within 24 hours.
// You can create reports for a range of up to 30 days for any time period after 2020-04-12.
// You can’t edit or remove report fields.
// Impression Share reports with a WEEKLY granularity value can’t have custom startTime and endTime in the request payload. Use dateRange instead. See CustomReportRequest.
// https://developer.apple.com/documentation/apple_search_ads/impression_share_report
func (s *ImpressionShareReportService) CreateImpressionShareReport(ctx context.Context, request CustomReportRequest) (*ImpressionShareReportResponse, *Response, error) {
	url := "custom-reports"
	res := new(ImpressionShareReportResponse)
	resp, err := s.client.post(ctx, url, request, res)

	return res, resp, err
}

// Fetches all Impression Share reports containing metrics and metadata.
// Use this endpoint to return all Impression Share reports containing metrics and metadata. Use query parameters as needed.
// The rate limit for this endpoint is 150 reports within 15 minutes.
// https://developer.apple.com/documentation/apple_search_ads/get_all_impression_share_reports
func (s *ImpressionShareReportService) GetAllImpressionShareReports(ctx context.Context, params *ImpressionShareReportRequest) (*ImpressionShareReportsResponse, *Response, error) {
	url := "custom-reports"
	res := new(ImpressionShareReportsResponse)
	resp, err := s.client.get(ctx, url, params, res)

	return res, resp, err
}

// Fetches a single Impression Share report containing metrics and metadata.
// Use this endpoint to return a single Impression Share report containing metrics and metadata. Use a reportId as a resource in the URI.
// The rate limit for this endpoint is 30 reports within 15 minutes.
// https://developer.apple.com/documentation/apple_search_ads/get_a_single_impression_share_report
func (s *ImpressionShareReportService) GetSingleImpressionShareReport(ctx context.Context, reportID int64) (*ImpressionShareReportResponse, *Response, error) {
	url := fmt.Sprintf("custom-reports/%d", reportID)
	res := new(ImpressionShareReportResponse)
	resp, err := s.client.get(ctx, url, nil, res)

	return res, resp, err
}

// Downloads the report from the downloadUri and parses the CSV data.
func (s *ImpressionShareReportService) DownloadReport(ctx context.Context, report ImpressionShareReport) (*DailyImpressionShareReport, error) {
	// Create a new HTTP request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, report.DownloadUri, nil)
	if err != nil {
		return nil, err
	}

	// Send the HTTP request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(resp.Body)

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %s", err)
	}

	// Convert the response body to a string
	csvData := string(body)

	// Create a DailyImpressionShareReport to hold our parsed data
	dailyReport := &DailyImpressionShareReport{}

	// Parse the CSV data
	if err := gocsv.UnmarshalString(csvData, &dailyReport.Records); err != nil {
		return nil, fmt.Errorf("failed to parse CSV data: %s", err)
	}

	return dailyReport, nil
}
