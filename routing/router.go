package routing

import (
	"awsx-api/handlers"
	"awsx-api/handlers/getLandingZoneDetails"
	"github.com/gorilla/mux"
	"net/http"
)

// Route describes a single route
type Route struct {
	Name          string
	Method        string
	Pattern       string
	HandlerFunc   http.HandlerFunc
	Authenticated bool
}

// Routes holds an array of Route. A note on swagger documentation. The path variables and query parameters
// are defined in ../doc.go.  YOu need to manually associate params and routes.
type Routes struct {
	Routes []Route
}

// NewRoutes creates and returns all the API routes
func NewRoutes() (r *Routes) {
	r = new(Routes)

	r.Routes = []Route{

		// swagger:route GET /awsx/appconfig
		// ---
		// Produces aws appconfig resource summary
		//
		// 		Produces:
		//		- application/json
		//
		//		Schemes: http, https
		//
		// responses:
		//    400: badRequestError
		//    503: serviceUnavailableError
		//		500: internalError
		//		200: metricsStatsResponse
		//{
		//	"AwsxAppconfig",
		//	"GET",
		//	"/awsx/appconfig",
		//	handlers.GetAppconfig,
		//	true,
		//},
		//{
		//	"AwsxLambda",
		//	"GET",
		//	"/awsx/lambda",
		//	handlers.GetLambdas,
		//	true,
		//},
		//{
		//	"AwsxLambdaWithTags",
		//	"GET",
		//	"/awsx/lambda/function-with-tag",
		//	handlers.GetLambdaWithTags,
		//	true,
		//},
		//{
		//	"AwsxTotalLambda",
		//	"GET",
		//	"/awsx/lambda/total-functions",
		//	handlers.GetNumberOfLambdas,
		//	true,
		//},
		//{
		//	"AwsxEks",
		//	"GET",
		//	"/awsx/eks",
		//	handlers.GetEks,
		//	true,
		//},
		//{
		//	"AwsxEcs",
		//	"GET",
		//	"/awsx/ecs",
		//	handlers.GetEcs,
		//	true,
		//},
		//{
		//	"AwsxVpc",
		//	"GET",
		//	"/awsx/vpc",
		//	handlers.GetVpc,
		//	true,
		//},
		//{
		//	"AwsxRds",
		//	"GET",
		//	"/awsx/rds",
		//	handlers.GetRds,
		//	true,
		//},
		//{
		//	"AwsxEc2",
		//	"GET",
		//	"/awsx/ec2",
		//	handlers.GetEc2,
		//	true,
		//},
		{
			"readyz",
			"GET",
			"/app-health/awsx-api/readyz",
			handlers.Readiness,
			false,
		},
		{
			"livez",
			"GET",
			"/app-health/awsx-api/livez",
			handlers.Liveness,
			false,
		},
		{
			"AwsxCloudWatchQueryApi",
			"GET",
			"/awsx-api/getQueryOutput",
			handlers.ExecuteQuery,
			true,
		},
		{
			"AwsxCloudWatchQueryApi",
			"GET",
			"/awsx-api/execute-query",
			getLandingZoneDetails.ExecuteLandingzoneQueries,
			true,
		},
		{
			"AwsxCloudTrailEventsApi",
			"GET",
			"/awsx-api/getEvents",
			handlers.GetAwsEvents,
			true,
		},
		// {
		// 	"AwsxEc2",
		// 	"GET",
		// 	"/awsx-api/getQueryOutput",
		// 	EC2.ExecuteQuery,
		// 	true,
		// },
		// {
		// 	"AwsxEc2",
		// 	"GET",
		// 	"/awsx-api/ec2/getQueryOutput",
		// 	EC2.ExecuteNetworkQuery,
		// 	true,
		// },
		// {
		// 	"AwsxEks",
		// 	"GET",
		// 	"/awsx-api/eksCpu/getQueryOutput",
		// 	EKS.GetEKScpuUtilizationPanel,
		// 	true,
		// },
		// {
		// 	"AwsxEks",
		// 	"GET",
		// 	"/awsx-api/eksMemory/getQueryOutput",
		// 	EKS.GetEKSMemoryUtilizationPanel,
		// 	true,
		// },
		// {
		// 	"AwsxEks",
		// 	"GET",
		// 	"/awsx-api/eksNetwork/getQueryOutput",
		// 	EKS.GetEKSNetworkUtilizationPanel,
		// 	true,
		// },
		// {
		// 	"AwsxEks",
		// 	"GET",
		// 	"/awsx-api/eksStorage/getQueryOutput",
		// 	EKS.GetEKSStorageUtilizationPanel,
		// 	true,
		// },
		// {
		// 	"AwsxEks",
		// 	"GET",
		// 	"/awsx-api/eksCpuRequests/getQueryOutput",
		// 	EKS.GetCpuRequestsPanel,
		// 	true,
		// },
		// {
		// 	"AwsxEks",
		// 	"GET",
		// 	"/awsx-api/eksCpuLimits/getQueryOutput",
		// 	EKS.GetCpuLimitPanel,
		// 	true,
		// },
		// {
		// 	"AwsxEks",
		// 	"GET",
		// 	"/awsx-api/eksAllocatableCpu/getQueryOutput",
		// 	EKS.GetAllocatableCpuPanel,
		// 	true,
		// },
		// {
		// 	"AwsxEks",
		// 	"GET",
		// 	"/awsx-api/eksCpuGraph/getQueryOutput",
		// 	EKS.GetCpuGraphPanel,
		// 	true,
		// },
		// {
		// 	"AwsxECS",
		// 	"GET",
		// 	"/awsx-api/ecsCpuRequests/getQueryOutput",
		// 	ECS.GetContainerPanel,
		// 	true,
		// },
		// {
		// 	"AwsxECS",
		// 	"GET",
		// 	"/awsx-api/ecsMemoryRequests/getQueryOutput",
		// 	ECS.GetECSMemoryUtilizationPanel,
		// 	true,
		// },
		// {
		// 	"AwsxECS",
		// 	"GET",
		// 	"/awsx-api/ecsStorageRequests/getQueryOutput",
		// 	ECS.GetStorageUtilizationPanel,
		// 	true,
		// },
		// {
		// 	"AwsxECS",
		// 	"GET",
		// 	"/awsx-api/ecsNetworkRequests/getQueryOutput",
		// 	ECS.GetECSNetworkUtilizationPanel,
		// 	true,
		// },

		// {
		// 	"AwsxEks",
		// 	"GET",
		// 	"/awsx-api/getQueryOutput",
		// 	EKS.ExecuteQuery
		// 	true,
		// },

		//{
		//	"AwsxS3",
		//	"GET",
		//	"/awsx/s3",
		//	handlers.GetS3,
		//	true,
		//},
		//{
		//	"AwsxS3WithTags",
		//	"GET",
		//	"/awsx/s3/bucket-with-tag",
		//	handlers.GetS3WithTags,
		//	true,
		//},
		//{
		//	"AwsxCdn",
		//	"GET",
		//	"/awsx/cdn",
		//	handlers.GetCdn,
		//	true,
		//},
		//{
		//	"AwsxKinesys",
		//	"GET",
		//	"/awsx/kinesys",
		//	handlers.GetKinesys,
		//	true,
		//},
		//{
		//	"AwsxDynamodb",
		//	"GET",
		//	"/awsx/dynamodb",
		//	handlers.GetDynamodb,
		//	true,
		//},
		//{
		//	"AwsxWaf",
		//	"GET",
		//	"/awsx/waf",
		//	handlers.GetWaf,
		//	true,
		//},
		//{
		//	"AwsxGlue",
		//	"GET",
		//	"/awsx/glue",
		//	handlers.GetGlue,
		//	true,
		//},
		//{
		//	"AwsxKms",
		//	"GET",
		//	"/awsx/kms",
		//	handlers.GetKms,
		//	true,
		//},
		//{
		//	"AwsxAppmesh",
		//	"GET",
		//	"/awsx/appmesh",
		//	handlers.GetAppmesh,
		//	true,
		//},
	}

	return
}

func NewRouter() *mux.Router {

	// conf := config.Get()
	// webRoot := "" //conf.Server.WebRoot
	// webRootWithSlash := webRoot + "/"

	appRouter := mux.NewRouter().StrictSlash(true)
	// appRouter := rootRouter

	// staticFileServer := http.FileServer(http.Dir(conf.Server.StaticContentRootDirectory))s

	// if webRoot != "/" {
	// 	// help the user out - if a request comes in for "/", redirect to our true webroot
	// 	rootRouter.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 		http.Redirect(w, r, webRootWithSlash, http.StatusFound)
	// 	})

	// 	appRouter = rootRouter.PathPrefix(conf.Server.WebRoot).Subrouter()
	// 	staticFileServer = http.StripPrefix(webRootWithSlash, staticFileServer)

	// 	// Because of OIDC, when we receive a request for the webroot without
	// 	// the trailing slash, we can not redirect the user to the correct
	// 	// webroot as the hash params are lost (and they are not sent to the
	// 	// server).
	// 	//
	// 	// See https://github.com/kiali/kiali/issues/3103
	// 	rootRouter.HandleFunc(webRoot, func(w http.ResponseWriter, r *http.Request) {
	// 		r.URL.Path = webRootWithSlash
	// 		rootRouter.ServeHTTP(w, r)
	// 	})
	// } else {
	// webRootWithSlash = "/"
	// }

	// fileServerHandler := func(w http.ResponseWriter, r *http.Request) {
	// 	urlPath := r.RequestURI
	// 	if r.URL != nil {
	// 		urlPath = r.URL.Path
	// 	}

	// 	if urlPath == webRootWithSlash || urlPath == webRoot || urlPath == webRootWithSlash+"index.html" {
	// 		serveIndexFile(w)
	// 	} else if urlPath == webRootWithSlash+"env.js" {
	// 		serveEnvJsFile(w)
	// 	} else {
	// 		staticFileServer.ServeHTTP(w, r)
	// 	}
	// }

	// appRouter = appRouter.StrictSlash(true)

	// Build our API server routes and install them.
	apiRoutes := NewRoutes()
	// authenticationHandler, _ := handlers.NewAuthenticationHandler()
	for _, route := range apiRoutes.Routes {
		// handlerFunction := metricHandler(route.HandlerFunc, route)
		// if route.Authenticated {
		// 	handlerFunction = authenticationHandler.Handle(handlerFunction)
		// } else {
		// 	handlerFunction = authenticationHandler.HandleUnauthenticated(handlerFunction)
		// }
		appRouter.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}

	// if authController := authentication.GetAuthController(); authController != nil {
	// 	if ac, ok := authController.(*authentication.OpenIdAuthController); ok {
	// 		ac.PostRoutes(appRouter)
	// 	}
	// }

	// All client-side routes are prefixed with /console.
	// They are forwarded to index.html and will be handled by react-router.
	// appRouter.PathPrefix("/console").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	serveIndexFile(w)
	// })

	// if authController := authentication.GetAuthController(); authController != nil {
	// 	if ac, ok := authController.(*authentication.OpenIdAuthController); ok {
	// 		authCallback := ac.GetAuthCallbackHandler(http.HandlerFunc(fileServerHandler))
	// 		rootRouter.Methods("GET").Path(webRootWithSlash).Handler(authCallback)
	// 	}
	// }

	// rootRouter.PathPrefix(webRootWithSlash).HandlerFunc(fileServerHandler)

	return appRouter
}

// statusResponseWriter contains a ResponseWriter and a StatusCode to read in the metrics middleware
type statusResponseWriter struct {
	http.ResponseWriter
	StatusCode int
}

// WriteHeader will be called by any function that needs to set an status code, in this function the StatusCode is also set
func (srw *statusResponseWriter) WriteHeader(code int) {
	srw.ResponseWriter.WriteHeader(code)
	srw.StatusCode = code
}

// updateMetric evaluates the StatusCode, if there is an error, increase the API failure counter, otherwise save the duration
// func updateMetric(route string, srw *statusResponseWriter, timer *prometheus.Timer) {
// 	// Always measure the duration even if the API call ended in an error
// 	timer.ObserveDuration()
// 	// Increase the error counter on 500 and 503 errors
// 	if srw.StatusCode == http.StatusInternalServerError || srw.StatusCode == http.StatusServiceUnavailable {
// 		internalmetrics.GetAPIFailureMetric(route).Inc()
// 	}
// }

// func metricHandler(next http.Handler, route Route) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		// By default, if there is no call to WriteHeader, an 200 will be
// 		srw := &statusResponseWriter{
// 			ResponseWriter: w,
// 			StatusCode:     http.StatusOK,
// 		}
// 		promtimer := internalmetrics.GetAPIProcessingTimePrometheusTimer(route.Name)
// 		defer updateMetric(route.Name, srw, promtimer)
// 		next.ServeHTTP(srw, r)
// 	})
// }

// serveEnvJsFile generates the env.js file needed by the UI from Kiali configs. The
// generated file is sent to the HTTP response.
// func serveEnvJsFile(w http.ResponseWriter) {
// 	conf := config.Get()
// 	var body string
// 	if len(conf.Server.WebHistoryMode) > 0 {
// 		body += fmt.Sprintf("window.HISTORY_MODE='%s';", conf.Server.WebHistoryMode)
// 	}

// 	body += "window.WEB_ROOT = document.getElementsByTagName('base')[0].getAttribute('href').replace(/^https?:\\/\\/[^#?\\/]+/g, '').replace(/\\/+$/g, '')"

// 	w.Header().Set("content-type", "text/javascript")
// 	_, err := io.WriteString(w, body)
// 	if err != nil {
// 		log.Errorf("HTTP I/O error [%v]", err.Error())
// 	}
// }

// serveIndexFile takes UI's index.html as a template to generate a modified index file that takes
// into account the web_root path configured in the Kiali CR. The result is sent to the HTTP response.
// func serveIndexFile(w http.ResponseWriter) {
// 	webRootPath := config.Get().Server.WebRoot
// 	webRootPath = strings.TrimSuffix(webRootPath, "/")

// 	path, _ := filepath.Abs("./console/index.html")
// 	b, err := ioutil.ReadFile(path)
// 	if err != nil {
// 		log.Errorf("File I/O error [%v]", err.Error())
// 		handlers.RespondWithDetailedError(w, http.StatusInternalServerError, "Unable to read index.html template file", err.Error())
// 		return
// 	}

// 	html := string(b)
// 	newHTML := html

// 	if len(webRootPath) != 0 {
// 		searchStr := `<base href="/"`
// 		newStr := `<base href="` + webRootPath + `/"`
// 		newHTML = strings.Replace(html, searchStr, newStr, -1)
// 	}

// 	w.Header().Set("content-type", "text/html")
// 	_, err = io.WriteString(w, newHTML)
// 	if err != nil {
// 		log.Errorf("HTTP I/O error [%v]", err.Error())
// 	}
// }
