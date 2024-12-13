package httpimplement

import (
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/labstack/echo"
	"github.com/moazedy/todo/internal/adapter/driven/db/repoimplement"
	"github.com/moazedy/todo/internal/domain/model"
	"github.com/moazedy/todo/internal/domain/srvimplement"
	"github.com/moazedy/todo/internal/port/driven/db/repository"
	"github.com/moazedy/todo/internal/port/driver/http"
	"github.com/moazedy/todo/internal/port/driver/service"
	"github.com/moazedy/todo/pkg/infra/config"
	"github.com/moazedy/todo/pkg/infra/queue"
	"github.com/moazedy/todo/pkg/infra/storage"
	"github.com/moazedy/todo/pkg/infra/tx"
	"gorm.io/gorm"
)

type serverItems struct {
	// basics
	dbConnection        *gorm.DB
	txFactory           tx.TXFactory
	storageAgentFactory storage.StorageAgentFactory
	awsClient           *s3.S3
	sqsClient           *sqs.Client
	sqsClientFactory    queue.SQSClientFactory

	// repo layer
	todoItemRepoFactory repository.GenericRepoFactory[model.TodoItem]
	fileRepo            repository.File
	queueRepo           repository.Queue

	// service layer
	todoService service.TodoItem
	fileService service.File

	// presentation layer
	todoController http.TodoItem
	fileController http.File
}

var items serverItems

func register(app *echo.Echo, cfg config.Config) {
	items.dbConnection = tx.GetDB(cfg.Postgres)
	items.txFactory = tx.NewTXFactory(cfg.Postgres.IsMock, items.dbConnection)
	items.awsClient = storage.CreateAWSS3Client(cfg.Storage.Endpoint, cfg.Storage.AccessKey, cfg.Storage.SecretKey, cfg.Storage.Bucket)
	items.storageAgentFactory = storage.NewStorageAgentFactory(cfg.Storage.IsMock, items.awsClient, cfg.Storage.Bucket)
	items.sqsClient = queue.NewSQSClient(cfg.Queue)
	items.sqsClientFactory = queue.NewSQSClientFactory(cfg.Queue.IsMock, items.sqsClient)

	items.todoItemRepoFactory = repoimplement.NewGenericRepoFactory[model.TodoItem](cfg.Postgres.IsMock)
	items.fileRepo = repoimplement.NewFile(items.storageAgentFactory)
	items.queueRepo = repoimplement.NewQueue(items.sqsClientFactory, cfg.Queue.QueueUrl)

	items.todoService = srvimplement.NewTodoItem(items.txFactory, items.todoItemRepoFactory, items.queueRepo)
	items.fileService = srvimplement.NewFile(items.fileRepo, cfg.Storage.MaxFileSize)

	items.todoController = NewTodoItem(items.todoService)

	items.fileController = NewFile(items.fileService)

	app.POST("/file", items.fileController.Upload)
	app.GET("/file/:file_name", items.fileController.Download)

	app.POST("/todo/item", items.todoController.Create)
	app.PUT("/todo/item", items.todoController.Update)
	app.DELETE("/todo/item/:id", items.todoController.Delete)
	app.GET("/todo/item/:id", items.todoController.GetByID)
}
