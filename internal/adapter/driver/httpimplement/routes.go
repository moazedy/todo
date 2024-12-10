package httpimplement

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/labstack/echo"
	"github.com/moazedy/todo/internal/adapter/driven/db/repoimplement"
	"github.com/moazedy/todo/internal/domain/model"
	"github.com/moazedy/todo/internal/domain/srvimplement"
	"github.com/moazedy/todo/internal/port/driven/db/repository"
	"github.com/moazedy/todo/internal/port/driver/http"
	"github.com/moazedy/todo/internal/port/driver/service"
	"github.com/moazedy/todo/pkg/infra/config"
	"github.com/moazedy/todo/pkg/infra/storage"
	"github.com/moazedy/todo/pkg/infra/tx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type serverItems struct {
	// basics
	dbConnection *gorm.DB
	txFactory    tx.TXFactory
	storageAgent storage.StorageAgent
	awsClient    *s3.Client

	// repo layer
	todoItemRepoFactory repository.GenericRepoFactory[model.TodoItem]
	fileRepo            repository.File

	// service layer
	todoService service.TodoItem
	fileService service.File

	// presentation layer
	todoController http.TodoItem
	fileController http.File
}

var items serverItems

func register(app *echo.Echo, cfg config.Config) {
	items.dbConnection = initDB(cfg.Postgres)
	items.txFactory = tx.NewTXFactory(items.dbConnection)
	items.awsClient = createAWSS3Client(cfg.Storage.Endpoint, cfg.Storage.AccessKey, cfg.Storage.SecretKey)
	items.storageAgent = storage.NewStorageAgent(items.awsClient, cfg.Storage.Bucket)

	items.todoItemRepoFactory = repoimplement.NewGenericRepoFactory[model.TodoItem]()
	items.fileRepo = repoimplement.NewFile(items.storageAgent)

	items.todoService = srvimplement.NewTodoItem(items.txFactory, items.todoItemRepoFactory)
	items.fileService = srvimplement.NewFile(items.fileRepo)

	items.todoController = NewTodoItem(items.todoService)

	items.fileController = NewFile(items.fileService)

	app.POST("/file", items.fileController.Upload)
	app.GET("/file/:file_id", items.fileController.Download)

	app.POST("/todo/item", items.todoController.Create)
	app.PUT("/todo/item", items.todoController.Update)
	app.DELETE("/todo/item/:id", items.todoController.Delete)
	app.GET("/todo/item/:id", items.todoController.GetByID)
}

// opening connection with database
func initDB(cfg config.PostgresConfig) *gorm.DB {
	db, err := gorm.Open(postgres.Open(cfg.ToString()), &gorm.Config{})
	if err != nil {
		println(err)
		panic("failed to connect with db")
	}

	// Check if the desired database exists
	err = db.Exec(
		fmt.Sprintf("CREATE DATABASE %s",
			cfg.Name,
		)).Error
	if err != nil {
		println(err.Error())
	}

	// Reconnect to the newly created database
	db, err = gorm.Open(postgres.Open(cfg.ToStringWithDbName()), &gorm.Config{})
	if err != nil {
		println(err)
		panic("failed to connect with db")
	}

	err = db.Exec(
		`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`,
	).Error
	if err != nil {
		panic(err.Error())
	}

	// NOTE : register all entities in here
	err = db.AutoMigrate(
		&model.TodoItem{},
	)
	if err != nil {
		println(err.Error())
		panic("failed to migrate")
	}

	return db
}

func createAWSS3Client(endpoint, accessKey, secretKey string) *s3.Client {
	// Load the Shared AWS Configuration (~/.aws/config)
	cfg, err := awsconfig.LoadDefaultConfig(context.TODO(), awsconfig.WithRegion("us-west-2"))
	if err != nil {
		log.Fatal(err)
	}

	cfg.Credentials = aws.CredentialsProviderFunc(func(ctx context.Context) (aws.Credentials, error) {
		return aws.Credentials{
			AccessKeyID:     accessKey,
			SecretAccessKey: secretKey,
		}, nil
	})
	cfg.BaseEndpoint = aws.String(endpoint)
	fmt.Printf("AWS S3 options -> endpoint: %s, access key: %s, secret key: %s\n", endpoint, accessKey, secretKey)

	// Create an Amazon S3 service client
	client := s3.NewFromConfig(cfg)

	return client
}
