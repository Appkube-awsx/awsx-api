package handlers

import (
	"awsx-api/helpers"
	"awsx-api/models"
	"awsx-api/util"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/grafana-tools/sdk"
	"github.com/sirupsen/logrus"
)

// GrafanaQueryRangeHandler is used for handling Grafana Range queries
func GrafanaQueryRangeHandler(w http.ResponseWriter, req *http.Request) {
	// if req.Method != http.MethodGet {
	// 	w.WriteHeader(http.StatusNotFound)
	// 	return
	// }

	reqQuery := req.URL.Query()
	client := util.NewGrafanaClient()
	data, err := GrafanaQueryRange(client, req.Context(), reqQuery.Get("url"), reqQuery.Get("api-key"), &reqQuery)
	if err != nil {
		// h.log.Error(ErrGrafanaQuery(err))
		// http.Error(w, ErrGrafanaQuery(err).Error(), http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(data)
}

func GrafanaQueryRange(g *models.GrafanaClient, ctx context.Context, BaseURL, APIKey string, queryData *url.Values) ([]byte, error) {
	if queryData == nil {
		return nil, errors.New("query data passed is nil")
	}

	c, err := sdk.NewClient(BaseURL, APIKey, g.HttpClient)
	if err != nil {
		return nil, util.CommonError(err)
	}

	ds, err := c.GetDatasourceByName(ctx, queryData.Get("ds"))
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	var reqURL string
	if g.PromMode {
		reqURL = fmt.Sprintf("%s/api/v1/query_range", BaseURL)
	} else {
		reqURL = fmt.Sprintf("%s/api/datasources/proxy/%d/api/v1/query_range", BaseURL, ds.ID)
	}

	newURL, _ := url.Parse(reqURL)
	q := url.Values{}
	q.Set("query", queryData.Get("query"))
	q.Set("start", queryData.Get("start"))
	q.Set("end", queryData.Get("end"))
	q.Set("step", queryData.Get("step"))
	newURL.RawQuery = q.Encode()
	queryURL := newURL.String()
	data, err := g.MakeRequest(ctx, queryURL, APIKey)
	if err != nil {
		return nil, util.CommonError(err)
	}
	return data, nil
}

func GetDs(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Starting getDs")
	grafanaUrl := r.URL.Query().Get("grafanaUrl")
	apiKey := r.URL.Query().Get("apiKey")
	if grafanaUrl == "" {
		fmt.Println("Grafana url not provided")
		return
	} else if apiKey == "" {
		fmt.Println("Grafana api key (userId:password) not provided")
		return
	}
	pref := &models.Preference{
		Grafana: &models.Grafana{
			// GrafanaURL:    "http://grafana.synectiks.net",
			// GrafanaAPIKey: "admin:password",
			GrafanaURL:    grafanaUrl,
			GrafanaAPIKey: apiKey,
		},
	}
	user := &models.User{
		UserID:    "admin",
		FirstName: "admin",
	}

	b := GrafanaBoardsHandler(pref, user)
	json.NewEncoder(w).Encode(b)
	fmt.Println("getDs completed")
}

func GrafanaBoardsHandler(prefObj *models.Preference, user *models.User) []*models.GrafanaBoard {
	// if req.Method != http.MethodGet && req.Method != http.MethodPost {
	// 	w.WriteHeader(http.StatusNotFound)
	// 	return
	// }

	// No POST for now. Commented
	// if req.Method == http.MethodPost {
	// 	h.SaveSelectedGrafanaBoardsHandler(w, req, prefObj, user, p)
	// 	return
	// }
	fmt.Println("Starting GrafanaBoardsHandler")
	client := util.NewGrafanaClient()
	if prefObj.Grafana == nil || prefObj.Grafana.GrafanaURL == "" {
		// h.log.Error(ErrGrafanaConfig)
		// http.Error(w, "Invalid grafana endpoint", http.StatusBadRequest)
		fmt.Println("Grafana URL null. Exiting")
		return nil
	}
	req := &http.Request{
		Method: "GET",
	}

	if err := helpers.Validate(client, req.Context(), prefObj.Grafana.GrafanaURL, prefObj.Grafana.GrafanaAPIKey); err != nil {
		// h.log.Error(ErrGrafanaScan(err))
		// http.Error(w, "Unable to connect to grafana", http.StatusInternalServerError)
		fmt.Println("Unable to connect to grafana. Exiting")
		return nil
	}

	var dashboardSearch = "" //req.URL.Query().Get("dashboardSearch")
	boards, err := GetGrafanaBoards(client, req.Context(), prefObj.Grafana.GrafanaURL, prefObj.Grafana.GrafanaAPIKey, dashboardSearch)
	if err != nil {
		// h.log.Error(ErrGrafanaBoards(err))
		// http.Error(w, "unable to get grafana boards", http.StatusInternalServerError)
		return nil
	}
	// fmt.Println(boards)
	w := new(bytes.Buffer)
	err = json.NewEncoder(w).Encode(boards)

	if err != nil {
		// obj := "boards payload"
		// h.log.Error(ErrMarshal(err, obj))
		// http.Error(w, "Unable to marshal the boards payload", http.StatusInternalServerError)
		fmt.Println("Unable to marshal the boards payload")
		return nil
	}
	fmt.Println("GrafanaBoardsHandler completed: ", w.String())
	return boards
}

func GetGrafanaBoards(g *models.GrafanaClient, ctx context.Context, BaseURL, APIKey, dashboardSearch string) ([]*models.GrafanaBoard, error) {
	fmt.Println("Starting GetGrafanaBoards")
	if strings.HasSuffix(BaseURL, "/") {
		BaseURL = strings.Trim(BaseURL, "/")
	}
	c, err := sdk.NewClient(BaseURL, APIKey, g.HttpClient)
	if err != nil {
		return nil, util.CommonError(err)
	}

	boardLinks, err := c.SearchDashboards(ctx, dashboardSearch, false)
	if err != nil {
		return nil, util.CommonError(err)
	}
	boards := []*models.GrafanaBoard{}
	for _, link := range boardLinks {
		if link.Type != "dash-db" {
			continue
		}
		// TODO Need to do the unitest for Grafana helper
		board, _, err := c.GetDashboardByUID(ctx, link.UID)
		// fmt.Println("DashBoard...... ", board)
		if err != nil {
			fmt.Println("ERROR in GetDashboardByUID")
			return nil, util.DashboardError(err, link.UID)
		}
		// b, _ := json.Marshal(board)
		// logrus.Debugf("Board before foramating: %s", b)
		// fmt.Println("")
		// fmt.Printf("Board before foramating %s", b)
		grafBoard, err := helpers.ProcessBoard(g, ctx, c, &board, &link)
		if err != nil {
			return nil, err
		}
		fmt.Println()
		fmt.Println()

		// fmt.Println("grafana dsboard")
		// fmt.Println(grafBoard)
		// b, _ := json.Marshal(grafBoard)
		// logrus.Debugf("Board after foramating: %s", b)
		// fmt.Println("")
		// fmt.Printf("Board after foramating %s", b)
		boards = append(boards, grafBoard)
	}
	// fmt.Println("Board after foramating ", boards)
	return boards, nil
}

func GrafanaQueryHandler(w http.ResponseWriter, r *http.Request) {

	grafanaUrl := r.URL.Query().Get("grafanaUrl")
	apiKey := r.URL.Query().Get("apiKey")
	if grafanaUrl == "" {
		fmt.Println("Grafana url not provided")
		return
	} else if apiKey == "" {
		fmt.Println("Grafana api key (userId:password) not provided")
		return
	}

	prefObj := &models.Preference{
		Grafana: &models.Grafana{
			// GrafanaURL:    "http://grafana.synectiks.net",
			// GrafanaAPIKey: "admin:password",
			GrafanaURL:    grafanaUrl,
			GrafanaAPIKey: apiKey,
		},
	}
	// user := &User{
	// 	UserID:    "admin",
	// 	FirstName: "admin",
	// }

	reqQuery := r.URL.Query()

	// if prefObj.Grafana == nil || prefObj.Grafana.GrafanaURL == "" {
	// 	err := ErrGrafanaConfig
	// 	h.log.Error(err)
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }
	client := util.NewGrafanaClient()
	data, err := GrafanaQuery(client, r.Context(), prefObj.Grafana.GrafanaURL, prefObj.Grafana.GrafanaAPIKey, &reqQuery)
	if err != nil {
		// h.log.Error(ErrGrafanaQuery(err))
		// http.Error(w, ErrGrafanaQuery(err).Error(), http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(data)
}

func GrafanaQuery(g *models.GrafanaClient, ctx context.Context, BaseURL, APIKey string, queryData *url.Values) ([]byte, error) {
	if queryData == nil {
		return nil, errors.New("query data passed is nil")
	}
	query := strings.TrimSpace(queryData.Get("query"))
	dsID := queryData.Get("dsid")
	var queryURL string
	switch {
	case strings.HasPrefix(query, "label_values("):
		val := strings.Replace(query, "label_values(", "", 1)
		val = strings.TrimSpace(strings.TrimSuffix(val, ")"))
		if strings.Contains(val, ",") {
			start := queryData.Get("start")
			end := queryData.Get("end")
			comInd := strings.LastIndex(val, ", ")
			if comInd > -1 {
				val = val[:comInd]
			}
			for key := range *queryData {
				if key != "query" && key != "dsid" && key != "start" && key != "end" {
					val1 := queryData.Get(key)
					val = strings.Replace(val, "$"+key, val1, -1)
				}
			}
			var reqURL string
			if g.PromMode {
				reqURL = fmt.Sprintf("%s/api/v1/series", BaseURL)
			} else {
				reqURL = fmt.Sprintf("%s/api/datasources/proxy/%s/api/v1/series", BaseURL, dsID)
			}
			queryURLInst, _ := url.Parse(reqURL)
			qParams := queryURLInst.Query()
			qParams.Set("match[]", val)
			if start != "" && end != "" {
				qParams.Set("start", start)
				qParams.Set("end", end)
			}
			queryURLInst.RawQuery = qParams.Encode()
			queryURL = queryURLInst.String()
		} else {
			if g.PromMode {
				queryURL = fmt.Sprintf("%s/api/v1/label/%s/values", BaseURL, val)
			} else {
				queryURL = fmt.Sprintf("%s/api/datasources/proxy/%s/api/v1/label/%s/values", BaseURL, dsID, val)
			}
		}
	case strings.HasPrefix(query, "query_result("):
		val := strings.Replace(query, "query_result(", "", 1)
		val = strings.TrimSpace(strings.TrimSuffix(val, ")"))
		for key := range *queryData {
			if key != "query" && key != "dsid" {
				val1 := queryData.Get(key)
				val = strings.Replace(val, "$"+key, val1, -1)
			}
		}
		var reqURL string
		if g.PromMode {
			reqURL = fmt.Sprintf("%s/api/v1/query", BaseURL)
		} else {
			reqURL = fmt.Sprintf("%s/api/datasources/proxy/%s/api/v1/query", BaseURL, dsID)
		}
		newURL, _ := url.Parse(reqURL)
		q := url.Values{}
		q.Set("query", val)
		newURL.RawQuery = q.Encode()
		queryURL = newURL.String()
	default:
		return json.Marshal(map[string]interface{}{
			"status": "success",
			"data":   []string{query},
		})
	}
	logrus.Debugf("derived query url: %s", queryURL)

	data, err := g.MakeRequest(ctx, queryURL, APIKey)
	if err != nil {
		return nil, errors.New("error getting data from Grafana API")
	}
	return data, nil
}
