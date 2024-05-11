package EC2

import (
	"awsx-api/handlers/getElementDetails/EC2"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

//	func TestGetQueryOutput(t *testing.T) {
//		// Define test cases
//		tests := []struct {
//			elementId      string
//			elementType    string
//			query          string
//			responseType   string
//			expected       int                    // Expected HTTP status code
//			expectedParams map[string]interface{} // Expected response parameters
//		}{
//			// Define your test cases here
//			{
//				elementId:    "14923",
//				elementType:  "EC2",
//				query:        "cpu_usage_idle_panel",
//				responseType: "frame",
//				expected:     http.StatusOK,
//				expectedParams: map[string]interface{}{
//					"elementId":    "14923",
//					"elementType":  "EC2",
//					"query":        "cpu_usage_idle_panel",
//					"responseType": "frame",
//				},
//			},
//			// Add more test cases as needed
//		}
//
//		// Iterate over test cases
//		for _, tc := range tests {
//			req, err := http.NewRequest("GET", "/awsx-api/getQueryOutput?elementId="+tc.elementId+"&elementType="+tc.elementType+"&query="+tc.query+"&responseType="+tc.responseType, nil)
//			if err != nil {
//				t.Fatalf("could not create request: %v", err)
//			}
//			rr := httptest.NewRecorder()
//
//			EC2.GetCPUUsageIdlePanel(rr, req)
//
//			// Check the status code
//			//if status := rr.Code; status != tc.expected {
//			//	t.Errorf("handler returned wrong status code: got %v want %v",
//			//		status, tc.expected)
//			//
//			//}
//			//
//			//// Check the Content-Type header
//			//expectedContentType := "application/json"
//			//if contentType := rr.Header().Get("Content-Type"); contentType != expectedContentType {
//			//	t.Errorf("handler returned wrong content type: got %v want %v",
//			//		contentType, expectedContentType)
//			//}
//			//
//			//// Decode the response body
//			//var responseBody map[string]interface{}
//			//if err := json.NewDecoder(rr.Body).Decode(&responseBody); err != nil {
//			//	t.Errorf("error decoding response body: %v", err)
//			//}
//			//
//			//// Check the response parameters
//			//for key, value := range tc.expectedParams {
//			//	if val, ok := responseBody[key]; !ok || val != value {
//			//		t.Errorf("unexpected value for parameter %s: got %v want %v", key, val, value)
//			//	}
//			//}
//		}
//		//expected := `{"CPU_Idle":[{"Timestamp":"your-timestamp","Value":your-value}]}`
//		//if rr.Body.String() != expected {
//		//	t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
//		//}
//	}
//
//	func TestGetQueryOutput(t *testing.T) {
//		// Define test cases
//		tests := []struct {
//			elementId      string
//			elementType    string
//			query          string
//			responseType   string
//			expected       int                    // Expected HTTP status code
//			expectedParams map[string]interface{} // Expected response parameters
//		}{
//			// Define your test cases here
//			{
//				elementId:    "14923",
//				elementType:  "EC2",
//				query:        "cpu_usage_idle_panel",
//				responseType: "frame",
//				expected:     http.StatusOK,
//				expectedParams: map[string]interface{}{
//					"elementId":    "14923",
//					"elementType":  "EC2",
//					"query":        "cpu_usage_idle_panel",
//					"responseType": "frame",
//				},
//			},
//			// Add more test cases as needed
//		}
//
//		// Iterate over test cases
//		for _, tc := range tests {
//			req, err := http.NewRequest("GET", "/awsx-api/getQueryOutput?elementId="+tc.elementId+"&elementType="+tc.elementType+"&query="+tc.query+"&responseType="+tc.responseType, nil)
//			if err != nil {
//				t.Fatalf("could not create request: %v", err)
//			}
//			rr := httptest.NewRecorder()
//
//			EC2.GetCPUUsageIdlePanel(rr, req)
//
//			// Check the status code
//			//if status := rr.Code; status != tc.expected {
//			//	t.Errorf("handler returned wrong status code: got %v want %v",
//			//		status, tc.expected)
//			//
//			//}
//			//
//			//// Check the Content-Type header
//			//expectedContentType := "application/json"
//			//if contentType := rr.Header().Get("Content-Type"); contentType != expectedContentType {
//			//	t.Errorf("handler returned wrong content type: got %v want %v",
//			//		contentType, expectedContentType)
//			//}
//			//
//			//// Decode the response body
//			//var responseBody map[string]interface{}
//			//if err := json.NewDecoder(rr.Body).Decode(&responseBody); err != nil {
//			//	t.Errorf("error decoding response body: %v", err)
//			//}
//			//
//			//// Check the response parameters
//			//for key, value := range tc.expectedParams {
//			//	if val, ok := responseBody[key]; !ok || val != value {
//			//		t.Errorf("unexpected value for parameter %s: got %v want %v", key, val, value)
//			//	}
//			//}
//		}
//		//expected := `{"CPU_Idle":[{"Timestamp":"your-timestamp","Value":your-value}]}`
//		//if rr.Body.String() != expected {
//		//	t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
//		//}
//	}
//
//func TestGetQueryOutput(t *testing.T) {
//	rr := httptest.NewRecorder()
//	request, err := http.NewRequest(http.MethodGet, "/awsx-api/getQueryOutput?elementId=14923&elementType=EC2&query=cpu_usage_idle_panel&responseType=frame&startTime=2024-03-28T00:00:00Z&endTime=2024-03-28T23:59:59Z", nil)
//	if err != nil {
//		return
//	}
//	EC2.GetCPUUsageIdlePanel(rr, request)
//	if rr.Result().StatusCode != http.StatusOK {
//		t.Errorf("expected 200 but got %d", rr.Result().StatusCode)
//	}
//	var responseBody map[string]interface{}
//	if err := json.NewDecoder(rr.Result().Body).Decode(&responseBody); err != nil {
//		t.Errorf("error decoding response body: %v", err)
//	}
//
//	// Define expected JSON string
//	expected := `{
//    "CPU_Idle": {
//        "Messages": null,
//        "MetricDataResults": [
//            {
//                "Id": "m1",
//                "Label": "cpu_usage_idle",
//                "Messages": null,
//                "StatusCode": "Complete",
//                "Timestamps": [
//                    "2024-03-28T23:55:00Z",
//                    "2024-03-28T23:50:00Z",
//                    "2024-03-28T23:45:00Z"
//                ],
//                "Values": [
//                    99.32285959009052,
//                    99.33969600692487,
//                    99.31958768160477
//                ]
//            }
//        ],
//        "NextToken": null
//    }
//}`
//
//	// Compare response body with expected JSON string
//	expectedBytes := []byte(expected)
//	responseBytes, err := json.Marshal(responseBody)
//	if err != nil {
//		t.Errorf("error encoding response body: %v", err)
//	}
//	if !bytes.Equal(expectedBytes, responseBytes) {
//		t.Errorf("expected response body does not match actual:\nExpected: %s\nActual: %s", expected, string(responseBytes))
//	}
//}

//	func TestGetCPUUsageIdlePanel(t *testing.T) {
//		server := httptest.NewServer(http.HandlerFunc(EC2.GetCPUUsageIdlePanel))
//		defer server.Close()
//		reqs, err := http.Get(server.URL)
//
//		if err != nil {
//			t.Error(err)
//		}
//		if reqs.StatusCode != http.StatusOK {
//			t.Errorf("expected 200 but got %d", reqs.StatusCode)
//		}
//		var responseBody map[string]interface{}
//		if err := json.NewDecoder(reqs.Body).Decode(&responseBody); err != nil {
//			t.Errorf("error decoding response body: %v", err)
//		}
//
//		// Define expected JSON string
//		expected := `{
//	   "CPU_Idle": {
//	       "Messages": null,
//	       "MetricDataResults": [
//	           {
//	               "Id": "m1",
//	               "Label": "cpu_usage_idle",
//	               "Messages": null,
//	               "StatusCode": "Complete",
//	               "Timestamps": [
//	                   "2024-03-28T23:55:00Z",
//	                   "2024-03-28T23:50:00Z",
//	                   "2024-03-28T23:45:00Z",
//	                   "2024-03-28T23:40:00Z",
//	                   "2024-03-28T23:35:00Z",
//	                   "2024-03-28T23:30:00Z",
//	                   "2024-03-28T23:25:00Z",
//	                   "2024-03-28T23:20:00Z",
//	                   "2024-03-28T23:15:00Z",
//	                   "2024-03-28T23:10:00Z",
//	                   "2024-03-28T23:05:00Z",
//	                   "2024-03-28T23:00:00Z",
//	                   "2024-03-28T22:55:00Z",
//	                   "2024-03-28T22:50:00Z",
//	                   "2024-03-28T22:45:00Z",
//	                   "2024-03-28T22:40:00Z",
//	                   "2024-03-28T22:35:00Z",
//	                   "2024-03-28T22:30:00Z",
//	                   "2024-03-28T22:25:00Z",
//	                   "2024-03-28T22:20:00Z",
//	                   "2024-03-28T22:15:00Z",
//	                   "2024-03-28T22:10:00Z",
//	                   "2024-03-28T22:05:00Z",
//	                   "2024-03-28T22:00:00Z",
//	                   "2024-03-28T21:55:00Z",
//	                   "2024-03-28T21:50:00Z",
//	                   "2024-03-28T21:45:00Z",
//	                   "2024-03-28T21:40:00Z",
//	                   "2024-03-28T21:35:00Z",
//	                   "2024-03-28T21:30:00Z",
//	                   "2024-03-28T21:25:00Z",
//	                   "2024-03-28T21:20:00Z",
//	                   "2024-03-28T21:15:00Z",
//	                   "2024-03-28T21:10:00Z",
//	                   "2024-03-28T21:05:00Z",
//	                   "2024-03-28T21:00:00Z",
//	                   "2024-03-28T20:55:00Z",
//	                   "2024-03-28T20:50:00Z",
//	                   "2024-03-28T20:45:00Z",
//	                   "2024-03-28T20:40:00Z",
//	                   "2024-03-28T20:35:00Z",
//	                   "2024-03-28T20:30:00Z",
//	                   "2024-03-28T20:25:00Z",
//	                   "2024-03-28T20:20:00Z",
//	                   "2024-03-28T20:15:00Z",
//	                   "2024-03-28T20:10:00Z",
//	                   "2024-03-28T20:05:00Z",
//	                   "2024-03-28T20:00:00Z",
//	                   "2024-03-28T19:55:00Z",
//	                   "2024-03-28T19:50:00Z",
//	                   "2024-03-28T19:45:00Z",
//	                   "2024-03-28T19:40:00Z",
//	                   "2024-03-28T19:35:00Z",
//	                   "2024-03-28T19:30:00Z",
//	                   "2024-03-28T19:25:00Z",
//	                   "2024-03-28T19:20:00Z",
//	                   "2024-03-28T19:15:00Z",
//	                   "2024-03-28T19:10:00Z",
//	                   "2024-03-28T19:05:00Z",
//	                   "2024-03-28T19:00:00Z",
//	                   "2024-03-28T18:55:00Z",
//	                   "2024-03-28T18:50:00Z",
//	                   "2024-03-28T18:45:00Z",
//	                   "2024-03-28T18:40:00Z",
//	                   "2024-03-28T18:35:00Z",
//	                   "2024-03-28T18:30:00Z",
//	                   "2024-03-28T18:25:00Z",
//	                   "2024-03-28T18:20:00Z",
//	                   "2024-03-28T18:15:00Z",
//	                   "2024-03-28T18:10:00Z",
//	                   "2024-03-28T18:05:00Z",
//	                   "2024-03-28T18:00:00Z",
//	                   "2024-03-28T17:55:00Z",
//	                   "2024-03-28T17:50:00Z",
//	                   "2024-03-28T17:45:00Z",
//	                   "2024-03-28T17:40:00Z",
//	                   "2024-03-28T17:35:00Z",
//	                   "2024-03-28T17:30:00Z",
//	                   "2024-03-28T17:25:00Z",
//	                   "2024-03-28T17:20:00Z",
//	                   "2024-03-28T17:15:00Z",
//	                   "2024-03-28T17:10:00Z",
//	                   "2024-03-28T17:05:00Z",
//	                   "2024-03-28T17:00:00Z",
//	                   "2024-03-28T16:55:00Z",
//	                   "2024-03-28T16:50:00Z",
//	                   "2024-03-28T16:45:00Z",
//	                   "2024-03-28T16:40:00Z",
//	                   "2024-03-28T16:35:00Z",
//	                   "2024-03-28T16:30:00Z",
//	                   "2024-03-28T16:25:00Z",
//	                   "2024-03-28T16:20:00Z",
//	                   "2024-03-28T16:15:00Z",
//	                   "2024-03-28T16:10:00Z",
//	                   "2024-03-28T16:05:00Z",
//	                   "2024-03-28T16:00:00Z",
//	                   "2024-03-28T15:55:00Z",
//	                   "2024-03-28T15:50:00Z",
//	                   "2024-03-28T15:45:00Z",
//	                   "2024-03-28T15:40:00Z",
//	                   "2024-03-28T15:35:00Z",
//	                   "2024-03-28T15:30:00Z",
//	                   "2024-03-28T15:25:00Z",
//	                   "2024-03-28T15:20:00Z",
//	                   "2024-03-28T15:15:00Z",
//	                   "2024-03-28T15:10:00Z",
//	                   "2024-03-28T15:05:00Z",
//	                   "2024-03-28T15:00:00Z",
//	                   "2024-03-28T14:55:00Z",
//	                   "2024-03-28T14:50:00Z",
//	                   "2024-03-28T14:45:00Z",
//	                   "2024-03-28T14:40:00Z",
//	                   "2024-03-28T14:35:00Z",
//	                   "2024-03-28T14:30:00Z",
//	                   "2024-03-28T14:25:00Z",
//	                   "2024-03-28T14:20:00Z",
//	                   "2024-03-28T14:15:00Z",
//	                   "2024-03-28T14:10:00Z",
//	                   "2024-03-28T14:05:00Z",
//	                   "2024-03-28T14:00:00Z",
//	                   "2024-03-28T13:55:00Z",
//	                   "2024-03-28T13:50:00Z",
//	                   "2024-03-28T13:45:00Z",
//	                   "2024-03-28T13:40:00Z",
//	                   "2024-03-28T13:35:00Z",
//	                   "2024-03-28T13:30:00Z",
//	                   "2024-03-28T13:25:00Z",
//	                   "2024-03-28T13:20:00Z",
//	                   "2024-03-28T13:15:00Z",
//	                   "2024-03-28T13:10:00Z",
//	                   "2024-03-28T13:05:00Z",
//	                   "2024-03-28T13:00:00Z",
//	                   "2024-03-28T12:55:00Z",
//	                   "2024-03-28T12:50:00Z",
//	                   "2024-03-28T12:45:00Z",
//	                   "2024-03-28T12:40:00Z",
//	                   "2024-03-28T12:35:00Z",
//	                   "2024-03-28T12:30:00Z",
//	                   "2024-03-28T12:25:00Z",
//	                   "2024-03-28T12:20:00Z",
//	                   "2024-03-28T12:15:00Z",
//	                   "2024-03-28T12:10:00Z",
//	                   "2024-03-28T12:05:00Z",
//	                   "2024-03-28T12:00:00Z",
//	                   "2024-03-28T11:55:00Z",
//	                   "2024-03-28T11:50:00Z",
//	                   "2024-03-28T11:45:00Z",
//	                   "2024-03-28T11:40:00Z",
//	                   "2024-03-28T11:35:00Z",
//	                   "2024-03-28T11:30:00Z",
//	                   "2024-03-28T11:25:00Z",
//	                   "2024-03-28T11:20:00Z",
//	                   "2024-03-28T11:15:00Z",
//	                   "2024-03-28T11:10:00Z",
//	                   "2024-03-28T11:05:00Z",
//	                   "2024-03-28T11:00:00Z",
//	                   "2024-03-28T10:55:00Z",
//	                   "2024-03-28T10:50:00Z",
//	                   "2024-03-28T10:45:00Z",
//	                   "2024-03-28T10:40:00Z",
//	                   "2024-03-28T10:35:00Z",
//	                   "2024-03-28T10:30:00Z",
//	                   "2024-03-28T10:25:00Z",
//	                   "2024-03-28T10:20:00Z",
//	                   "2024-03-28T10:15:00Z",
//	                   "2024-03-28T10:10:00Z",
//	                   "2024-03-28T10:05:00Z",
//	                   "2024-03-28T10:00:00Z",
//	                   "2024-03-28T09:55:00Z",
//	                   "2024-03-28T09:50:00Z",
//	                   "2024-03-28T09:45:00Z",
//	                   "2024-03-28T09:40:00Z",
//	                   "2024-03-28T09:35:00Z",
//	                   "2024-03-28T09:30:00Z",
//	                   "2024-03-28T09:25:00Z",
//	                   "2024-03-28T09:20:00Z",
//	                   "2024-03-28T09:15:00Z",
//	                   "2024-03-28T09:10:00Z",
//	                   "2024-03-28T09:05:00Z",
//	                   "2024-03-28T09:00:00Z",
//	                   "2024-03-28T08:55:00Z",
//	                   "2024-03-28T08:50:00Z",
//	                   "2024-03-28T08:45:00Z",
//	                   "2024-03-28T08:40:00Z",
//	                   "2024-03-28T08:35:00Z",
//	                   "2024-03-28T08:30:00Z",
//	                   "2024-03-28T08:25:00Z",
//	                   "2024-03-28T08:20:00Z",
//	                   "2024-03-28T08:15:00Z",
//	                   "2024-03-28T08:10:00Z",
//	                   "2024-03-28T08:05:00Z",
//	                   "2024-03-28T08:00:00Z",
//	                   "2024-03-28T07:55:00Z",
//	                   "2024-03-28T07:50:00Z",
//	                   "2024-03-28T07:45:00Z",
//	                   "2024-03-28T07:40:00Z",
//	                   "2024-03-28T07:35:00Z",
//	                   "2024-03-28T07:30:00Z",
//	                   "2024-03-28T07:25:00Z",
//	                   "2024-03-28T07:20:00Z",
//	                   "2024-03-28T07:15:00Z",
//	                   "2024-03-28T07:10:00Z",
//	                   "2024-03-28T07:05:00Z",
//	                   "2024-03-28T07:00:00Z",
//	                   "2024-03-28T06:55:00Z",
//	                   "2024-03-28T06:50:00Z",
//	                   "2024-03-28T06:45:00Z",
//	                   "2024-03-28T06:40:00Z",
//	                   "2024-03-28T06:35:00Z",
//	                   "2024-03-28T06:30:00Z",
//	                   "2024-03-28T06:25:00Z",
//	                   "2024-03-28T06:20:00Z",
//	                   "2024-03-28T06:15:00Z",
//	                   "2024-03-28T06:10:00Z",
//	                   "2024-03-28T06:05:00Z",
//	                   "2024-03-28T06:00:00Z",
//	                   "2024-03-28T05:55:00Z",
//	                   "2024-03-28T05:50:00Z",
//	                   "2024-03-28T05:45:00Z",
//	                   "2024-03-28T05:40:00Z",
//	                   "2024-03-28T05:35:00Z",
//	                   "2024-03-28T05:30:00Z",
//	                   "2024-03-28T05:25:00Z",
//	                   "2024-03-28T05:20:00Z",
//	                   "2024-03-28T05:15:00Z",
//	                   "2024-03-28T05:10:00Z",
//	                   "2024-03-28T05:05:00Z",
//	                   "2024-03-28T05:00:00Z",
//	                   "2024-03-28T04:55:00Z",
//	                   "2024-03-28T04:50:00Z",
//	                   "2024-03-28T04:45:00Z",
//	                   "2024-03-28T04:40:00Z",
//	                   "2024-03-28T04:35:00Z",
//	                   "2024-03-28T04:30:00Z",
//	                   "2024-03-28T04:25:00Z",
//	                   "2024-03-28T04:20:00Z",
//	                   "2024-03-28T04:15:00Z",
//	                   "2024-03-28T04:10:00Z",
//	                   "2024-03-28T04:05:00Z",
//	                   "2024-03-28T04:00:00Z",
//	                   "2024-03-28T03:55:00Z",
//	                   "2024-03-28T03:50:00Z",
//	                   "2024-03-28T03:45:00Z",
//	                   "2024-03-28T03:40:00Z",
//	                   "2024-03-28T03:35:00Z",
//	                   "2024-03-28T03:30:00Z",
//	                   "2024-03-28T03:25:00Z",
//	                   "2024-03-28T03:20:00Z",
//	                   "2024-03-28T03:15:00Z",
//	                   "2024-03-28T03:10:00Z",
//	                   "2024-03-28T03:05:00Z",
//	                   "2024-03-28T03:00:00Z",
//	                   "2024-03-28T02:55:00Z",
//	                   "2024-03-28T02:50:00Z",
//	                   "2024-03-28T02:45:00Z",
//	                   "2024-03-28T02:40:00Z",
//	                   "2024-03-28T02:35:00Z",
//	                   "2024-03-28T02:30:00Z",
//	                   "2024-03-28T02:25:00Z",
//	                   "2024-03-28T02:20:00Z",
//	                   "2024-03-28T02:15:00Z",
//	                   "2024-03-28T02:10:00Z",
//	                   "2024-03-28T02:05:00Z",
//	                   "2024-03-28T02:00:00Z",
//	                   "2024-03-28T01:55:00Z",
//	                   "2024-03-28T01:50:00Z",
//	                   "2024-03-28T01:45:00Z",
//	                   "2024-03-28T01:40:00Z",
//	                   "2024-03-28T01:35:00Z",
//	                   "2024-03-28T01:30:00Z",
//	                   "2024-03-28T01:25:00Z",
//	                   "2024-03-28T01:20:00Z",
//	                   "2024-03-28T01:15:00Z",
//	                   "2024-03-28T01:10:00Z",
//	                   "2024-03-28T01:05:00Z",
//	                   "2024-03-28T01:00:00Z",
//	                   "2024-03-28T00:55:00Z",
//	                   "2024-03-28T00:50:00Z",
//	                   "2024-03-28T00:45:00Z",
//	                   "2024-03-28T00:40:00Z",
//	                   "2024-03-28T00:35:00Z",
//	                   "2024-03-28T00:30:00Z",
//	                   "2024-03-28T00:25:00Z",
//	                   "2024-03-28T00:20:00Z",
//	                   "2024-03-28T00:15:00Z",
//	                   "2024-03-28T00:10:00Z",
//	                   "2024-03-28T00:05:00Z",
//	                   "2024-03-28T00:00:00Z"
//	               ],
//	               "Values": [
//	                   99.32285959009052,
//	                   99.33969600692487,
//	                   99.31958768160477,
//	                   99.05487418004839,
//	                   99.33300196852683,
//	                   99.36337912807161,
//	                   99.31307553339792,
//	                   99.36338921997366,
//	                   99.3229691925602,
//	                   99.05810026742208,
//	                   99.33314808907282,
//	                   99.3665641054654,
//	                   99.35308846751818,
//	                   99.39337940441942,
//	                   99.33303852862545,
//	                   99.09153260552254,
//	                   99.3097199504307,
//	                   99.37652950547347,
//	                   99.31966517115323,
//	                   99.3731974386433,
//	                   99.32642650650007,
//	                   99.10166152001133,
//	                   99.32957666617764,
//	                   99.34990849629098,
//	                   99.31294875062508,
//	                   99.39340014533872,
//	                   99.3062025219165,
//	                   99.05132603604932,
//	                   99.33959827428882,
//	                   99.36638600588212,
//	                   99.31624434073242,
//	                   99.3898499339474,
//	                   99.28940156721971,
//	                   99.01793020307495,
//	                   99.31608533144322,
//	                   99.3731996674282,
//	                   99.34302809331798,
//	                   99.35666988059306,
//	                   99.31331749866001,
//	                   99.06874279872842,
//	                   99.3468873339135,
//	                   99.34014451954,
//	                   99.33335798871983,
//	                   99.38016546101058,
//	                   99.34000483943456,
//	                   99.04522383942505,
//	                   99.31323102280888,
//	                   99.3803035073058,
//	                   99.35344112377884,
//	                   99.3467756913802,
//	                   99.32005300840075,
//	                   99.02538588760339,
//	                   99.29650197448748,
//	                   99.35350502006379,
//	                   99.34994733726508,
//	                   99.38360477646624,
//	                   99.29650201191032,
//	                   99.05524490561761,
//	                   99.34680599253076,
//	                   99.10926224436753,
//	                   99.28329466057663,
//	                   99.39346580540261,
//	                   98.96421120483578,
//	                   99.06533839398081,
//	                   99.32000756787612,
//	                   99.39353267941789,
//	                   99.27662806929527,
//	                   99.3400974080099,
//	                   99.29651375767692,
//	                   99.04209595556605,
//	                   99.32000917348654,
//	                   99.35675015565556,
//	                   99.30997022137586,
//	                   99.350083092166,
//	                   99.28985513714026,
//	                   98.7499443516879,
//	                   99.31333323352588,
//	                   99.330122302802,
//	                   99.28991180466848,
//	                   99.33331527663377,
//	                   99.27994848470583,
//	                   99.02163609656398,
//	                   99.27988622121724,
//	                   99.34011981579748,
//	                   99.29307334020372,
//	                   99.34346150130972,
//	                   99.30319094149351,
//	                   99.02476832213517,
//	                   99.30993656420469,
//	                   99.32329142635194,
//	                   99.09904342208077,
//	                   99.30661500481308,
//	                   99.27969591177926,
//	                   99.02827307687566,
//	                   99.25639407749513,
//	                   99.30674354187626,
//	                   99.26637422775454,
//	                   99.1693245146565,
//	                   99.30998877276704,
//	                   99.34995735990148,
//	                   98.96501427003507,
//	                   99.3333967087585,
//	                   99.08260132611062,
//	                   93.81975552163665,
//	                   99.31988915954928,
//	                   99.34003791826179,
//	                   98.93814625827608,
//	                   99.29661700518616,
//	                   99.32324990101696,
//	                   99.30674408885108,
//	                   99.31317669125738,
//	                   99.35341866806999,
//	                   98.95833365706604,
//	                   99.34005701175194,
//	                   99.29984239132304,
//	                   99.28994314660068,
//	                   99.29637234404566,
//	                   99.32004402356135,
//	                   99.01835750819825,
//	                   99.30675246567658,
//	                   99.31323391534497,
//	                   99.33336744587612,
//	                   99.31998571607554,
//	                   99.35676914283297,
//	                   98.97175085450804,
//	                   99.29674095511551,
//	                   99.29324898594534,
//	                   99.2965271808176,
//	                   99.32665776934206,
//	                   99.33333437688151,
//	                   98.98831840238456,
//	                   99.33669847321087,
//	                   99.3268222006556,
//	                   99.34343571177138,
//	                   99.28661436036292,
//	                   99.31668333569453,
//	                   99.00179172525047,
//	                   99.26288610238626,
//	                   99.33009591009777,
//	                   99.32640183721415,
//	                   99.32974858095562,
//	                   99.3331570003754,
//	                   98.98803443519049,
//	                   99.35653406737127,
//	                   99.33628184029763,
//	                   99.34996806577068,
//	                   99.31981564896707,
//	                   99.36662461589928,
//	                   99.00134323990008,
//	                   99.34314881867138,
//	                   99.34644393578336,
//	                   99.3565570757034,
//	                   99.31953089473129,
//	                   99.3700478145669,
//	                   98.97127322040912,
//	                   99.34654280220484,
//	                   99.29982963862462,
//	                   99.33309745324689,
//	                   99.32281021828656,
//	                   99.36315364342134,
//	                   99.03137151269922,
//	                   99.34653714947171,
//	                   99.28972944915263,
//	                   99.34991247993847,
//	                   99.28622888310625,
//	                   99.36333027549303,
//	                   98.98458268593212,
//	                   99.33976456415226,
//	                   99.35643685211797,
//	                   99.31977576033805,
//	                   99.33959999174381,
//	                   99.34973277903289,
//	                   98.99076649943638,
//	                   99.36636973198792,
//	                   99.32287595873018,
//	                   99.35623292009628,
//	                   99.29602531177225,
//	                   99.16856317620638,
//	                   99.06118045935527,
//	                   99.31972914727646,
//	                   99.3294963111261,
//	                   99.3329301117731,
//	                   99.31623249486867,
//	                   99.33648296836114,
//	                   98.9942768483626,
//	                   99.37304060597312,
//	                   99.27932770668711,
//	                   99.28933525932213,
//	                   99.30289335611887,
//	                   99.30284117048743,
//	                   98.9908812299888,
//	                   99.30620625376051,
//	                   99.30948510526396,
//	                   99.33286548006646,
//	                   99.29588312820812,
//	                   99.33977800905608,
//	                   98.97784831189216,
//	                   99.31310704901809,
//	                   99.31971233999808,
//	                   99.3129244760004,
//	                   99.29301326155219,
//	                   99.32611253260083,
//	                   98.97084197851794,
//	                   99.34294216474959,
//	                   99.32265284620303,
//	                   99.34291179133852,
//	                   99.3093491814004,
//	                   99.30284613958158,
//	                   98.13643772909056,
//	                   99.11170724965363,
//	                   99.13824324723363,
//	                   99.23572319515895,
//	                   99.15192954889848,
//	                   99.04769712980806,
//	                   73.80206801823189,
//	                   99.35006620014569,
//	                   99.31973258082239,
//	                   99.35337432052435,
//	                   99.31650762735501,
//	                   99.36995061796154,
//	                   99.0213516606282,
//	                   99.33331198169284,
//	                   99.3431725015075,
//	                   99.36667867513063,
//	                   99.24285182983014,
//	                   99.38665956936761,
//	                   99.3230731143407,
//	                   99.03834350670225,
//	                   99.32982884342951,
//	                   99.34659274266414,
//	                   99.35327660764182,
//	                   99.3498496381281,
//	                   99.32307311961155,
//	                   99.0584401077439,
//	                   99.28634278092082,
//	                   99.36332294918728,
//	                   99.31977523203138,
//	                   99.36654553406403,
//	                   99.33306381469902,
//	                   99.04163517766655,
//	                   99.33323561710704,
//	                   99.38671675706738,
//	                   99.34653382638915,
//	                   99.39013880428732,
//	                   99.33645380740556,
//	                   99.098349621923,
//	                   99.2963492874093,
//	                   99.34318874079551,
//	                   99.33642853491995,
//	                   99.35658856799748,
//	                   99.35322498478263,
//	                   97.3171625688009,
//	                   99.36659554488605,
//	                   99.39030829847972,
//	                   99.34005476598114,
//	                   99.43035606196558,
//	                   99.35321263075016,
//	                   99.0750490188673,
//	                   99.36666178671366,
//	                   99.3834324000375,
//	                   99.3465556992215,
//	                   99.38683081463846,
//	                   99.34324994063203,
//	                   99.1118757976784,
//	                   99.34995403197493,
//	                   99.39337933047648,
//	                   99.33658680056621,
//	                   99.37674689400586,
//	                   99.33326928459704,
//	                   99.08493751054588,
//	                   99.3666388043836,
//	                   99.37668791716833,
//	                   99.31316373819713,
//	                   99.40347400416314,
//	                   99.3297704539645,
//	                   99.07492883786846,
//	                   99.31622068552176,
//	                   99.36985074478712,
//	                   99.35653350915004,
//	                   99.37332605511423,
//	                   99.32973395259478,
//	                   99.03447373592526,
//	                   99.3530969316227,
//	                   99.38000601396428,
//	                   99.33639141809263,
//	                   99.41345519211787,
//	                   99.34656241198316,
//	                   98.97771059729699
//	               ]
//	           }
//	       ],
//	       "NextToken": null
//	   }
//	}`
//
//		// Compare response body with expected JSON string
//		expectedBytes := []byte(expected)
//		responseBytes, err := json.Marshal(responseBody)
//		if err != nil {
//			t.Errorf("error encoding response body: %v", err)
//		}
//		if !bytes.Equal(expectedBytes, responseBytes) {
//			t.Errorf("expected response body does not match actual:\nExpected: %s\nActual: %s", expected, string(responseBytes))
//		}
//	}
//
// func TestGetQueryOutput(t *testing.T) {
//
//	// Mock API server
//	mockServer := mockAPIServer()
//	defer mockServer.Close()
//
//	// Replace API base URL with the mock server URL
//	APIBaseURL := "/awsx-api/getQueryOutput?elementId=14923&elementType=EC2&query=cpu_usage_idle_panel&responseType=frame&startTime=2024-03-28T00:00:00Z&endTime=2024-03-28T23:59:59Z"
//	originalAPIBaseURL := APIBaseURL
//	APIBaseURL = mockServer.URL
//	defer func() { APIBaseURL = originalAPIBaseURL }()
//
//	// Now, you can proceed with your test as before
//	rr := httptest.NewRecorder()
//	request, err := http.NewRequest(http.MethodGet, "/awsx-api/getQueryOutput?elementId=14923&elementType=EC2&query=cpu_usage_idle_panel&responseType=frame&startTime=2024-03-28T00:00:00Z&endTime=2024-03-28T23:59:59Z", nil)
//	if err != nil {
//		t.Fatal("Failed to create request:", err)
//	}
//	EC2.GetCPUUsageIdlePanel(rr, request)
//
//	// Check status code
//	if rr.Result().StatusCode != http.StatusOK {
//		t.Errorf("expected 200 but got %d", rr.Result().StatusCode)
//	}
//
//	// Decode response body
//
// }
//
//	func mockAPIServer() *httptest.Server {
//		// Create a new mock server
//		return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//			// Define the response body
//			responseBody := `{
//	   "CPU_Idle": {
//	       "Messages": null,
//	       "MetricDataResults": [
//	           {
//	               "Id": "m1",
//	               "Label": "cpu_usage_idle",
//	               "Messages": null,
//	               "StatusCode": "Complete",
//	               "Timestamps": [
//	                   "2024-03-28T23:55:00Z",
//	                   "2024-03-28T23:50:00Z"
//	               ],
//	               "Values": [
//	                   99.32285959009052,
//	                   99.33969600692487
//	               ]
//	           }
//	       ],
//	       "NextToken": null
//	   }
//	}`
//
//			// Set response headers
//			w.Header().Set("Content-Type", "application/json")
//			w.WriteHeader(http.StatusOK)
//
//			// Write the response body
//			_, _ = w.Write([]byte(responseBody))
//		}))
//	}

func TestRouting(t *testing.T) {
	mockServer := mockAPIServer()
	defer mockServer.Close()

	request, err := http.NewRequest(http.MethodGet, mockServer.URL+"/awsx-api/getQueryOutput?elementId=14923&elementType=EC2&query=cpu_usage_idle_panel&responseType=frame&startTime=2024-03-28T00:00:00Z&endTime=2024-03-28T23:59:59Z", nil)
	if err != nil {
		t.Fatal("Failed to create request:", err)
	}

	rr := httptest.NewRecorder()

	EC2.GetCPUUsageIdlePanel(rr, request)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status code 200 but got %d", rr.Code)
	}

	// Decode response body and check the content
	var responseBody map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &responseBody); err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	// Check if the CPU_Idle key exists
	cpuIdleData, ok := responseBody["CPU_Idle"]
	if !ok {
		t.Error("Expected 'CPU_Idle' key in response body, but not found")
	}

	// Check if there are MetricDataResults
	metricDataResults, ok := cpuIdleData.(map[string]interface{})["MetricDataResults"].([]interface{})
	if !ok || len(metricDataResults) == 0 {
		t.Error("Expected 'MetricDataResults' in 'CPU_Idle' data, but not found or empty")
	}

	for _, metricData := range metricDataResults {
		metricDataMap, ok := metricData.(map[string]interface{})
		if !ok {
			t.Error("Invalid format for MetricDataResult")
			continue
		}

		// Check the presence of required fields
		requiredFields := []string{"Id", "Label", "StatusCode", "Timestamps", "Values"}
		for _, field := range requiredFields {
			if _, found := metricDataMap[field]; !found {
				t.Errorf("Expected field '%s' in MetricDataResult, but not found", field)
			}
		}

		timestamps, ok := metricDataMap["Timestamps"].([]interface{})
		if !ok {
			t.Error("Invalid format for Timestamps")
			continue
		}
		values, ok := metricDataMap["Values"].([]interface{})
		if !ok {
			t.Error("Invalid format for Values")
			continue
		}

		if len(timestamps) != len(values) {
			t.Errorf("Length of Timestamps (%d) does not match length of Values (%d)", len(timestamps), len(values))
		}

	}

}
func TestGetCPUUsageIdlePanel(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(EC2.GetCPUUsageIdlePanel))
	defer server.Close()
	reqs, err := http.Get(server.URL)

	if err != nil {
		t.Error(err)
	}
	if reqs.StatusCode != http.StatusOK {
		t.Errorf("expected status code 200 but got %d", reqs.StatusCode)
	}

	// Decode response body and check the content
	var responseBody map[string]interface{}
	if err := json.NewDecoder(reqs.Body).Decode(&responseBody); err != nil {
		t.Errorf("error decoding response body: %v", err)
	}

	// Check if the CPU_Idle key exists
	cpuIdleData, ok := responseBody["CPU_Idle"]
	if !ok {
		t.Error("Expected 'CPU_Idle' key in response body, but not found")
	}

	// Check if there are MetricDataResults
	metricDataResults, ok := cpuIdleData.(map[string]interface{})["MetricDataResults"].([]interface{})
	if !ok || len(metricDataResults) == 0 {
		t.Error("Expected 'MetricDataResults' in 'CPU_Idle' data, but not found or empty")
	}

	for _, metricData := range metricDataResults {
		metricDataMap, ok := metricData.(map[string]interface{})
		if !ok {
			t.Error("Invalid format for MetricDataResult")
			continue
		}

		// Check the presence of required fields
		requiredFields := []string{"Id", "Label", "StatusCode", "Timestamps", "Values"}
		for _, field := range requiredFields {
			if _, found := metricDataMap[field]; !found {
				t.Errorf("Expected field '%s' in MetricDataResult, but not found", field)
			}
		}

		timestamps, ok := metricDataMap["Timestamps"].([]interface{})
		if !ok {
			t.Error("Invalid format for Timestamps")
			continue
		}
		values, ok := metricDataMap["Values"].([]interface{})
		if !ok {
			t.Error("Invalid format for Values")
			continue
		}

		if len(timestamps) != len(values) {
			t.Errorf("Length of Timestamps (%d) does not match length of Values (%d)", len(timestamps), len(values))
		}

	}

}

func mockAPIServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	}))
}
