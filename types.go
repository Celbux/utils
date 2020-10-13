package celbuxhelpers

import (
	"cloud.google.com/go/bigquery"
	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	"cloud.google.com/go/datastore"
	"cloud.google.com/go/errorreporting"
	"cloud.google.com/go/logging"
	"cloud.google.com/go/storage"
)

// KindSuffix TODO annotate
var KindSuffix = GetTimeString()

/* GCP Helpers
 * Celbux helpers for easy & neat integration with GCP in main app engine code
 * Requirement: Call initialiseClients(projectID) in main app start up
 */

// ErrorClient TODO annotate
var ErrorClient *errorreporting.Client

// DatastoreClient TODO annotate
var DatastoreClient *datastore.Client

// StorageClient TODO annotate
var StorageClient *storage.Client

// LoggingClient TODO annotate
var LoggingClient *logging.Client

// TasksClient TODO annotate
var TasksClient *cloudtasks.Client

// BigQueryClient TODO annotate
var BigQueryClient *bigquery.Client
