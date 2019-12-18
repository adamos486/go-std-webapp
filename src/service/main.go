package main

import (
	"database/sql"
	"fmt"
	osLog "log"
	"net/http"
	"os"
	"service/auth"
	"service/auth/basic"
	"service/auth/token/jwt"
	"service/database"
	"service/handlers/index"
	"service/handlers/recovery"
	"service/handlers/request"
	"service/identity"
	"service/log"

	"github.com/go-chi/chi"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"

	"go.uber.org/zap"
)

func checkForProd() bool {
	return os.Getenv("ENV") == os.Getenv("PROD")
}

func dotEnv() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
}

func main() {
	//dotenv first
	dotEnv()

	//Determine if server is running prod.
	isProd := checkForProd()

	//Initialize db client
	db := setupDBClient()
	defer func() {
		closeErr := db.Close()
		if closeErr != nil {
			panic(closeErr)
		}
	}()

	//Initialize log client
	logger := setupLogClient(isProd)

	//Initialize auth client
	authClient := setupAuthClient()

	//Initialize route handlers
	indexRoute := index.New(logger, db)
	identityRoute := setupIdentity(logger, db, authClient)

	//Configure chi router
	router := setupChiRouter(authClient, logger)

	//Configure routes
	router.Get("/", indexRoute.Handler)
	router.Get("/identity", identityRoute.Handler)
	router.Post("/identity", identityRoute.CreateIdentity)
	router.Post("/auth", identityRoute.AuthIdentity)

	//Serve
	fmt.Println("Starting up server @ localhost:9000/")
	err := http.ListenAndServe(":9000", router)
	if err != nil {
		osLog.Fatal(err)
	}
}

func setupAuthClient() *auth.Client {
	basicAuth := basic.NewAuth()
	jwtService := jwt.NewService()
	return auth.NewClient(basicAuth, jwtService)
}

func setupChiRouter(auth *auth.Client, log log.ProdInterface) *chi.Mux {
	router := chi.NewRouter()

	router.Use(request.GenerateRequestIDMiddle)

	request.SetupLogger(log)
	router.Use(request.Logger)

	basic.SetupAuthMiddleware(auth, log)
	router.Use(basic.AuthMiddleware)

	recovery.SetupRecover(log)
	router.Use(recovery.Recover)

	return router
}

func setupDBClient() *sql.DB {
	db, err := sql.Open("postgres", os.Getenv("DB_ADDR"))
	if err != nil {
		panic(err)
	}
	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(25)
	return db
}

func setupIdentity(logger log.ProdInterface, db database.DBInterface,
	auth auth.Interface) *identity.HandlerObject {
	identityService := identity.NewServiceObject(logger, db)
	return identity.NewHandlerObject(logger, identityService, auth)
}

func setupLogClient(prod bool) *zap.Logger {
	var logger *zap.Logger
	var zapErr error
	if !prod {
		logger, zapErr = zap.NewDevelopment()
		if zapErr != nil {
			panic(zapErr)
		}
	} else {
		logger, zapErr = zap.NewProduction()
		if zapErr != nil {
			panic(zapErr)
		}
	}
	return logger
}
