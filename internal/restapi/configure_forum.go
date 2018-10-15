// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"net/http"
	"strings"

	assetfs "github.com/elazarl/go-bindata-assetfs"
	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	middleware "github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"

	"github.com/SinimaWath/tp-db/internal/modules/assets/assets_ui"
	"github.com/SinimaWath/tp-db/internal/restapi/operations"
)

//go:generate swagger generate server --target .. --name Forum --spec ../../api/swagger.yml
//go:generate go-bindata -pkg assets_ui -o ../modules/assets/assets_ui/assets_ui.go -prefix ../../api/swagger-ui/ ../../api/swagger-ui/...
//go:generate go-bindata -pkg assets_db -o ../modules/assets/assets_db/assets_db.go -prefix ../../assets/ ../../assets/...

type DatabaseFlags struct {
	Database string `long:"database" description:"database connection parameters" default:"sqlite3:tech-db-hello.db"`
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

	api.ClearHandler = operations.ClearHandlerFunc(func(params operations.ClearParams) middleware.Responder {
		return middleware.NotImplemented("operation .Clear has not yet been implemented")
	})
	api.ForumCreateHandler = operations.ForumCreateHandlerFunc(func(params operations.ForumCreateParams) middleware.Responder {
		return middleware.NotImplemented("operation .ForumCreate has not yet been implemented")
	})
	api.ForumGetOneHandler = operations.ForumGetOneHandlerFunc(func(params operations.ForumGetOneParams) middleware.Responder {
		return middleware.NotImplemented("operation .ForumGetOne has not yet been implemented")
	})
	api.ForumGetThreadsHandler = operations.ForumGetThreadsHandlerFunc(func(params operations.ForumGetThreadsParams) middleware.Responder {
		return middleware.NotImplemented("operation .ForumGetThreads has not yet been implemented")
	})
	api.ForumGetUsersHandler = operations.ForumGetUsersHandlerFunc(func(params operations.ForumGetUsersParams) middleware.Responder {
		return middleware.NotImplemented("operation .ForumGetUsers has not yet been implemented")
	})
	api.PostGetOneHandler = operations.PostGetOneHandlerFunc(func(params operations.PostGetOneParams) middleware.Responder {
		return middleware.NotImplemented("operation .PostGetOne has not yet been implemented")
	})
	api.PostUpdateHandler = operations.PostUpdateHandlerFunc(func(params operations.PostUpdateParams) middleware.Responder {
		return middleware.NotImplemented("operation .PostUpdate has not yet been implemented")
	})
	api.PostsCreateHandler = operations.PostsCreateHandlerFunc(func(params operations.PostsCreateParams) middleware.Responder {
		return middleware.NotImplemented("operation .PostsCreate has not yet been implemented")
	})
	api.StatusHandler = operations.StatusHandlerFunc(func(params operations.StatusParams) middleware.Responder {
		return middleware.NotImplemented("operation .Status has not yet been implemented")
	})
	api.ThreadCreateHandler = operations.ThreadCreateHandlerFunc(func(params operations.ThreadCreateParams) middleware.Responder {
		return middleware.NotImplemented("operation .ThreadCreate has not yet been implemented")
	})
	api.ThreadGetOneHandler = operations.ThreadGetOneHandlerFunc(func(params operations.ThreadGetOneParams) middleware.Responder {
		return middleware.NotImplemented("operation .ThreadGetOne has not yet been implemented")
	})
	api.ThreadGetPostsHandler = operations.ThreadGetPostsHandlerFunc(func(params operations.ThreadGetPostsParams) middleware.Responder {
		return middleware.NotImplemented("operation .ThreadGetPosts has not yet been implemented")
	})
	api.ThreadUpdateHandler = operations.ThreadUpdateHandlerFunc(func(params operations.ThreadUpdateParams) middleware.Responder {
		return middleware.NotImplemented("operation .ThreadUpdate has not yet been implemented")
	})
	api.ThreadVoteHandler = operations.ThreadVoteHandlerFunc(func(params operations.ThreadVoteParams) middleware.Responder {
		return middleware.NotImplemented("operation .ThreadVote has not yet been implemented")
	})
	api.UserCreateHandler = operations.UserCreateHandlerFunc(func(params operations.UserCreateParams) middleware.Responder {
		return middleware.NotImplemented("operation .UserCreate has not yet been implemented")
	})
	api.UserGetOneHandler = operations.UserGetOneHandlerFunc(func(params operations.UserGetOneParams) middleware.Responder {
		return middleware.NotImplemented("operation .UserGetOne has not yet been implemented")
	})
	api.UserUpdateHandler = operations.UserUpdateHandlerFunc(func(params operations.UserUpdateParams) middleware.Responder {
		return middleware.NotImplemented("operation .UserUpdate has not yet been implemented")
	})

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
