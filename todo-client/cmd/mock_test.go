package cmd

import (
	"net/http"
	"net/http/httptest"
)

var testResp = map[string]struct {
	Status int
	Body   string
}{
	"resultsMany": {
		Status: http.StatusOK,
		Body: `
		{
			"data": [
			  {
				"Task": "Task 1",
				"Done": false,
				"CreatedAt": "2019-10-28T08:23:38.310097076-04:00",
				"CompletedAt": "0001-01-01T00:00:00Z"
			  },
			  {
				"Task": "Task 2",
				"Done": false,
				"CreatedAt": "2019-10-28T08:23:38.310097076-04:00",
				"CompletedAt": "0001-01-01T00:00:00Z"
			  }
			],
			"date": 12425245,
			"total_results": 2
		}`,
	},
	"resultOne": {
		Status: http.StatusOK,
		Body: `
		{
			"data":[
				{
					"Task":"Task 1",
					"Done":false,
					"CreatedAt":"2019-10-28T08:23:38.310097076-04:00",
					"CompletedAt":"0001-01-01T00:00:00Z"
				}
			],
			"date":1234566567,
			"total_results":1
		}`,
	},
	"noResult": {
		Status: http.StatusOK,
		Body: `
		{
			"data":[],
			"date":12345667,
			"total_results":0
		}`,
	},
	"root": {
		Status: http.StatusOK,
		Body:   "Hello from root route",
	},
	"notFound": {
		Status: http.StatusNotFound,
		Body:   "404 - not found",
	},
}

func mockServer(handler http.HandlerFunc) (string, func()) {
	server := httptest.NewServer(handler)

	return server.URL, func() {
		server.Close()
	}
}
