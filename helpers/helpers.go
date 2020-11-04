package helpers

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"strings"
	"time"

	"cloud.google.com/go/bigquery"
	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	"cloud.google.com/go/datastore"
	"cloud.google.com/go/errorreporting"
	"cloud.google.com/go/logging"
	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/storage"
	"github.com/golang/gddo/httputil/header"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	taskspb "google.golang.org/genproto/googleapis/cloud/tasks/v2"
	ltype "google.golang.org/genproto/googleapis/logging/type"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

/*
 * GCP Helpers
 * Celbux helpers for easy & neat integration with GCP in main app engine code
 * Requirement: Call InitialiseClients(projectID) in main app start up
 */

// GCP Clients
var ErrorClient *errorreporting.Client
var DatastoreClient *datastore.Client
var StorageClient *storage.Client
var LoggingClient *logging.Client
var TasksClient *cloudtasks.Client
var BigQueryClient *bigquery.Client
var PubSubClient *pubsub.Client
var Ctx context.Context

var KindSuffix = GetTimeString()

func GetProjectID() (string, error) {
	// Get Project ID
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if projectID == "" {
		return "", status.Error(codes.NotFound, "env var GOOGLE_CLOUD_PROJECT must be set")
	}

	return projectID, nil
}

//InitialiseClients provides all required GCP clients for use in main app engine code
//It also takes in an optional serviceAccount
//This is the location of the *.json file containing your GCP API key
//serviceAccount is the relative path location to the file
func InitialiseClients(projectID string, serviceAccount... string) error {
	Ctx = context.Background()
	// Initialise error to prevent shadowing
	var err error
	if len(serviceAccount) > 0 {
		setGCPKey(serviceAccount[0])
	}

	// Creates error client
	if ErrorClient == nil {
		ErrorClient, err = errorreporting.NewClient(Ctx, projectID, errorreporting.Config{
			ServiceName: projectID + "-service",
			OnError: func(err error) {
				log.Printf("Could not log error: %v", err)
			},
		})
		if err != nil {
			log.Fatal(err)
		}
	}

	// Creates datastore client
	if DatastoreClient == nil {
		DatastoreClient, err = datastore.NewClient(Ctx, projectID)
		if err != nil {
			return LogError(err)
		}
	}

	// Creates logging client
	if LoggingClient == nil {
		LoggingClient, err = logging.NewClient(Ctx, projectID)
		if err != nil {
			return LogError(err)
		}
	}

	// Creates storage client
	if StorageClient == nil {
		StorageClient, err = storage.NewClient(Ctx)
		if err != nil {
			return LogError(err)
		}
	}

	// Creates storage client
	if TasksClient == nil {
		TasksClient, err = cloudtasks.NewClient(Ctx)
		if err != nil {
			return LogError(err)
		}
	}

	// Creates BigQuery client
	if BigQueryClient == nil {
		BigQueryClient, err = bigquery.NewClient(Ctx, projectID)
		if err != nil {
			return LogError(err)
		}
	}

	// Creates PubSub client
	if PubSubClient == nil {
		PubSubClient, err = pubsub.NewClient(Ctx, projectID)
		if err != nil {
			return LogError(err)
		}
	}

	return nil
}

func RunBigQuery(query string) error {
	q := BigQueryClient.Query(query)
	q.Location = "EU"
	job, err := q.Run(Ctx)
	if err != nil {
		return err
	}

	_, err = job.Wait(Ctx)
	if err != nil {
		return err
	}

	return nil
}

func EncodeStruct(w http.ResponseWriter, obj interface{}) error {
	// Writes the encoded marshalled json into the http writer mainly for the purpose of a response
	(w).Header().Set("Content-Type", "application/json; charset=utf-8")
	err := json.NewEncoder(w).Encode(obj)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

func DecodeStruct(w http.ResponseWriter, r *http.Request, obj interface{}) error {
	// Decode request into provided struct pointer
	if r.Header.Get("Content-Type") != "" {
		value, _ := header.ParseValueAndParams(r.Header, "Content-Type")
		if value != "application/json" {
			return status.Error(codes.FailedPrecondition, "content type must be application/json")
		}
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1048576)
	err := json.NewDecoder(r.Body).Decode(&obj)

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	return nil
}

func GLog(name string, text string, severity *ltype.LogSeverity) {
	//severity is nillable. Debug by default
	// Sets log name to unix nano second
	logger := LoggingClient.Logger(name)

	// Set severity based on params. Default Severity: DEBUG
	var logSeverity logging.Severity
	if severity == nil {
		logSeverity = logging.Severity(ltype.LogSeverity_DEBUG)
	} else {
		logSeverity = logging.Severity(*severity)
	}

	// Adds an entry to the log buffer.
	logger.Log(logging.Entry{
		Payload:  text,
		Severity: logSeverity,
	})

	// Print for local environment. Displayed as default severity in GCP
	if IsDev() {
		fmt.Printf("LOCAL [%v] %v\n", logSeverity.String(), text)
	}
}

func LogError(err error) error {
	// Log for Logs Viewer
	ErrorClient.Report(errorreporting.Entry{
		Error: err,
	})

	// Log for Local
	fmt.Printf("Error: %v", err.Error())

	// Optional for quick-hand returns in other func
	return err
}

func DownloadObject(bucket string, object string) ([]byte, error) {
	//DownloadObject downloads an object from Cloud Storage
	rc, err := StorageClient.Bucket(bucket).Object(object).NewReader(context.Background())
	if err != nil {
		return nil, fmt.Errorf("Object(%q).NewReader: %v", object, err)
	}
	defer rc.Close()

	data, err := ioutil.ReadAll(rc)
	if err != nil {
		return nil, fmt.Errorf("ioutil.ReadAll: %v", err)
	}

	return data, nil
}

func QueueHTTPRequest(ctx context.Context, queuePath string, request *taskspb.HttpRequest) (*taskspb.Task, error) {
	// Build the Task payload.
	req := &taskspb.CreateTaskRequest{
		Parent: queuePath,
		Task: &taskspb.Task{
			MessageType: &taskspb.Task_HttpRequest{
				HttpRequest: request,
			},
		},
	}

	ctxDeadline, cancel := context.WithDeadline(ctx, time.Now().Add(time.Second * 30))
	defer cancel()

	createdTask, err := TasksClient.CreateTask(ctxDeadline, req)
	if err != nil {
		return nil, LogError(err)
	}

	return createdTask, nil
}

func PrintHTTPBody(resp *http.Response) (string, error) {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", LogError(err)
	}
	return string(body), nil
}

func Encrypt(data string) string {
	return b64.URLEncoding.EncodeToString([]byte(data))
}

func Decrypt(data string) (string, error) {
	s, err := b64.URLEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}
	return string(s), nil
}

func GetTestName() string {
	// Gets the current running method by reflection.
	// this is useful for linking tests to functions for logging.

	fpcs := make([]uintptr, 1)
	runtime.Callers(2, fpcs)
	caller := runtime.FuncForPC(fpcs[0] - 1)
	r := strings.Replace(caller.Name(), "github.com/MSpaceDev/JiraOnTheGO/src/service", "", -1)
	return strings.Replace(r, ".", "", -1)
}

func WriteFile(data string, name string) error {
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write([]byte(data))
	if err != nil {
		return err
	}
	return nil
}

/*
 * Default: ReadModeSingle (Single is stored in the 0th element)
 * ReadModeSingle returns a single string, including all \n
 * ReadModeSingleCollapsed returns a single string, with all \n stripped away
 * ReadModeMultiline returns an array of string, split on \n (each line)
 */
const (
	ReadModeSingle = iota
	ReadModeSingleCollapsed
	ReadModeMultiline
)
func ReadFile(file string, method... int) ([]string, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
	    return nil, err
	}

	if method == nil {
		method = []int { ReadModeSingle }
	}

	switch method[0] {
	case ReadModeSingle:
		return []string { string(data) }, nil
	case ReadModeSingleCollapsed:
		return []string { strings.Replace(string(data), "\n", "", -1) }, nil
	case ReadModeMultiline:
		return strings.Split(string(data), "\n"), nil
	}

	return []string { string(data) }, nil
}

func GetTimeString() string {
	loc, _ := time.LoadLocation("Africa/Johannesburg")
	startTime := time.Now().In(loc).String()
	return startTime[:len(startTime)-18]
}

func GetKind(kind string) string {
	if IsDev() {
		return kind + KindSuffix
	}
	return kind
}

func SetKind(val string) {
	if IsDev() {
		KindSuffix = GetTimeString()
		KindSuffix += val
		fmt.Printf("KindSuffix :%v\n", KindSuffix)
	}
}

// IsDev returns true when this app is NOT deployed, and is run locally
func IsDev() bool {
	return !appengine.IsDevAppServer()
}

func IsNil(i interface{}) bool {
	if i == nil {
		return true
	}
	switch reflect.TypeOf(i).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		return reflect.ValueOf(i).IsNil()
	}
	return false
}

func setGCPKey(key string) {
	absPath, err := filepath.Abs(key)
	if err != nil {
		fmt.Printf("could not find key at location: %v", key)
	}


	err = os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", absPath)
	if err != nil {
		fmt.Printf("could not find key at location: %v", absPath)
	}
}

func Match(data string, regex string) ([][]string, error) {
	r, err := regexp.Compile(regex)
	if err != nil {
	    return nil, err
	}


	return r.FindAllStringSubmatch(data, -1), nil
}