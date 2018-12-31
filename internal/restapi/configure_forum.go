// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"log"
	"net/http"
	"strings"

	"gopkg.in/jackc/pgx.v2"

	assetfs "github.com/elazarl/go-bindata-assetfs"
	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	"github.com/go-openapi/swag"

	"github.com/SinimaWath/tp-db/internal/modules/assets/assets_ui"
	"github.com/SinimaWath/tp-db/internal/modules/service"
	"github.com/SinimaWath/tp-db/internal/restapi/operations"
)

//go:generate swagger generate server --target .. --name Forum --spec ../../api/swagger.yml

type DatabaseFlags struct {
	Database string `long:"database" description:"database connection parameters" default:""`
}

var dbFlags DatabaseFlags

func configureFlags(api *operations.ForumAPI) {
	api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{
		{"database", "database connection parameters", &dbFlags},
	}
}

func configureAPI(api *operations.ForumAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()

	api.BinConsumer = runtime.ByteStreamConsumer()

	api.JSONProducer = runtime.JSONProducer()

	log.Printf("configure Api")
	log.Println(dbFlags.Database)
	config, err := pgx.ParseConnectionString(dbFlags.Database)
	if err != nil {
		log.Fatal(err)
	}

	var handler service.ForumHandler = service.NewForumPgsql(&config)
	api.ClearHandler = operations.ClearHandlerFunc(handler.Clear)
	api.ForumCreateHandler = operations.ForumCreateHandlerFunc(handler.ForumCreate)
	api.UserCreateHandler = operations.UserCreateHandlerFunc(handler.UserCreate)
	api.UserGetOneHandler = operations.UserGetOneHandlerFunc(handler.UserGetOne)
	api.UserUpdateHandler = operations.UserUpdateHandlerFunc(handler.UserUpdate)

	api.ForumGetOneHandler = operations.ForumGetOneHandlerFunc(handler.ForumGetOne)
	api.ForumGetThreadsHandler = operations.ForumGetThreadsHandlerFunc(handler.ForumGetThreads)
	api.ForumGetUsersHandler = operations.ForumGetUsersHandlerFunc(handler.ForumGetUsers)

	api.ThreadCreateHandler = operations.ThreadCreateHandlerFunc(handler.ThreadCreate)
	api.ThreadGetOneHandler = operations.ThreadGetOneHandlerFunc(handler.ThreadGetOne)
	api.ThreadUpdateHandler = operations.ThreadUpdateHandlerFunc(handler.ThreadUpdate)
	api.ThreadVoteHandler = operations.ThreadVoteHandlerFunc(handler.ThreadVote)

	api.PostsCreateHandler = operations.PostsCreateHandlerFunc(handler.PostsCreate)
	api.PostGetOneHandler = operations.PostGetOneHandlerFunc(handler.PostGetOne)
	api.PostUpdateHandler = operations.PostUpdateHandlerFunc(handler.PostUpdate)
	api.ThreadGetPostsHandler = operations.ThreadGetPostsHandlerFunc(handler.ThreadGetPosts)
	api.StatusHandler = operations.StatusHandlerFunc(handler.Status)
	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix"
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return uiMiddleware(handler)
}

func uiMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/swagger.json" {
			handler.ServeHTTP(w, r)
			return
		}
		// Serving Swagger UI
		if r.URL.Path == "/api/" {
			r.URL.Path = "/api"
		}
		if r.URL.Path != "/api" && !strings.HasPrefix(r.URL.Path, "/api/") {
			http.FileServer(&assetfs.AssetFS{
				Asset:     assets_ui.Asset,
				AssetDir:  assets_ui.AssetDir,
				AssetInfo: assets_ui.AssetInfo,
			}).ServeHTTP(w, r)
			return
		}
		handler.ServeHTTP(w, r)
	})
}
