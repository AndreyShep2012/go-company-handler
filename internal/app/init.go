package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/golang-jwt/jwt"
	slogfiber "github.com/samber/slog-fiber"
	"gopkg.in/mgo.v2/bson"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func initLogger(level string) {
	logLevel := slog.LevelInfo
	switch level {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	}

	opts := &slog.HandlerOptions{Level: logLevel}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, opts))
	slog.SetDefault(logger)
}

func initFiberServer(apiRoot, jwtKey string) (*fiber.App, fiber.Router) {
	fiberServer := fiber.New(fiber.Config{
		CaseSensitive: false,
	})

	fiberServer.Use(panicMiddleware())

	apiRouter := fiberServer.Group(apiRoot)
	apiRouter.Use(requestid.New())
	apiRouter.Use(jwtMiddleware([]byte(jwtKey)))

	extendedLogs := false
	if slog.Default().Enabled(context.Background(), slog.LevelDebug) {
		extendedLogs = true
	}

	config := slogfiber.Config{WithRequestBody: extendedLogs, WithResponseBody: extendedLogs, WithRequestHeader: extendedLogs, WithResponseHeader: extendedLogs}
	apiRouter.Use(slogfiber.NewWithConfig(slog.Default(), config))

	return fiberServer, apiRouter
}

func initMongo(ctx context.Context, mongoUri, databaseName, collectionName string, connectTimeoutSec int) *mongo.Collection {
	slog.Info("connecting to mongo: ", "uri", mongoUri, "timeout", connectTimeoutSec)

	clientOptions := options.Client().ApplyURI(mongoUri)
	clientOptions.SetConnectTimeout(time.Duration(connectTimeoutSec) * time.Second)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		panic("failed to connect to mongo: " + err.Error())
	}

	if err := client.Ping(context.Background(), nil); err != nil {
		panic("failed to ping mongo: " + err.Error())
	}

	collection := client.Database(databaseName).Collection(collectionName)

	// Create a unique index on the name field
	indexModel := mongo.IndexModel{
		Keys:    bson.M{"name": 1},
		Options: options.Index().SetUnique(true),
	}
	_, err = collection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		panic("failed to create unique index on name field: " + err.Error())
	}

	slog.Info("connected to mongo")
	return collection
}

func panicMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		defer func() {
			if r := recover(); r != nil {
				slog.Error("panic happened, err: ", r, "stack", string(debug.Stack()))

				var ok bool
				if err, ok = r.(error); !ok {
					err = fmt.Errorf("panic: %v", r)
				}
			}
		}()

		return c.Next()
	}
}

func jwtMiddleware(secretKey []byte) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		if c.Method() == http.MethodGet {
			return c.Next()
		}

		tokenString := c.Get("Authorization")
		if tokenString == "" {
			slog.Debug("get empty Authorization header")
			return fiber.ErrUnauthorized
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		slog.Debug("get token", "token", tokenString)

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return secretKey, nil
		})

		if err != nil {
			slog.Debug("parse token error", "error", err.Error())
			return fiber.ErrUnauthorized
		}

		if !token.Valid {
			slog.Debug("token is not valid")
			return fiber.ErrUnauthorized
		}

		slog.Debug("token parsed successfully")
		return c.Next()
	}
}
