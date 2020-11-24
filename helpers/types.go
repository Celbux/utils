package helpers

import (
	"cloud.google.com/go/bigquery"
	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	"cloud.google.com/go/datastore"
	"cloud.google.com/go/errorreporting"
	"cloud.google.com/go/logging"
	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/storage"
	"context"
)

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

var isDeclared bool

/*
 * ReadFile(file, method)
 * Default Method: ReadModeSingle (Single is stored in the 0th element)
 * ReadModeSingle returns a single string, including all \n
 * ReadModeSingleCollapsed returns a single string, with all \n stripped away
 * ReadModeMultiline returns an array of string, split on \n (each line)
 */
const (
	ReadModeSingle = iota
	ReadModeSingleCollapsed
	ReadModeMultiline
)

type Response struct {
	Data string `json:"data"`
	Error string `json:"error"`
}
