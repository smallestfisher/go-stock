package main

import (
	"context"
	"encoding/json"
	"fmt"
	"go-stock/backend/data"
	"go-stock/backend/db"
	log "go-stock/backend/logger"
	"go-stock/backend/models"
	"go-stock/backend/util"
	"net/http"
	"os"
	"reflect"
	"runtime/debug"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var Version = "1.0.0-web"
var VersionCommit = "web-build"
var OFFICIAL_STATEMENT = "Web Mode"
var BuildKey = "cc1e0d684e32f176c56ff1fcf384dcd9"

func checkDir(dir string) {
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		os.Mkdir(dir, os.ModePerm)
		log.SugaredLogger.Info("create dir: " + dir)
	}
}

func PanicHandler() {
	if r := recover(); r != nil {
		fmt.Printf("Recovered from panic: %v\n", r)
		debug.PrintStack()
	}
}

func AutoMigrate() {
	db.Dao.AutoMigrate(&data.StockInfo{})
	db.Dao.AutoMigrate(&data.StockBasic{})
	db.Dao.AutoMigrate(&data.FollowedStock{})
	db.Dao.AutoMigrate(&data.IndexBasic{})
	db.Dao.AutoMigrate(&data.Settings{})
	db.Dao.AutoMigrate(&models.AIResponseResult{})
	db.Dao.AutoMigrate(&models.StockInfoHK{})
	db.Dao.AutoMigrate(&models.StockInfoUS{})
	db.Dao.AutoMigrate(&data.FollowedFund{})
	db.Dao.AutoMigrate(&data.FundBasic{})
	db.Dao.AutoMigrate(&models.PromptTemplate{})
	db.Dao.AutoMigrate(&data.Group{})
	db.Dao.AutoMigrate(&data.GroupStock{})
	db.Dao.AutoMigrate(&models.Tags{})
	db.Dao.AutoMigrate(&models.Telegraph{})
	db.Dao.AutoMigrate(&models.TelegraphTags{})
	db.Dao.AutoMigrate(&models.LongTigerRankData{})
	db.Dao.AutoMigrate(&data.AIConfig{})
	db.Dao.AutoMigrate(&models.BKDict{})
	db.Dao.AutoMigrate(&models.WordAnalyze{})
	db.Dao.AutoMigrate(&models.SentimentResultAnalyze{})
	db.Dao.AutoMigrate(&models.AiRecommendStocks{})
	db.Dao.AutoMigrate(&models.AllStockInfo{})
	db.Dao.AutoMigrate(&models.CronTask{})
	db.Dao.AutoMigrate(&models.AiAssistantSession{})
}

// Use the same assets as main.go

// SSE event structure
type SSEEvent struct {
	Name string      `json:"name"`
	Data interface{} `json:"data"`
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.SugaredLogger.Error("panic: ", r)
			log.SugaredLogger.Error("stack: ", string(debug.Stack()))
		}
	}()

	// Initialize backend
	checkDir("data")
	db.Init("")
	
	// MUST set IsWebMode before any util.Emit calls to avoid Wails context errors
	util.IsWebMode = true

	log.SugaredLogger.Info("Initializing database schema...")
	AutoMigrate()
	
	// Sentiment analysis can be slow and may panic, run it in background
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.SugaredLogger.Errorf("Panic during InitAnalyzeSentiment: %v", r)
			}
			util.Emit(context.Background(), "loadingMsg", "done")
		}()
		log.SugaredLogger.Info("Initializing sentiment analysis data in background...")
		data.InitAnalyzeSentiment()
	}()

	log.SugaredLogger.Info("Starting Web Server...")

	app := NewApp()
	app.ctx = context.Background()

	// Ensure we are ready as soon as possible for the UI
	// (The actual heavy data loading is happening in the background)
	util.IsReady = true
	util.Emit(app.ctx, "loadingMsg", "done")

	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Basic Auth
	user := os.Getenv("APP_USER")
	pass := os.Getenv("APP_PASSWORD")
	if user != "" && pass != "" {
		e.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
			if username == user && password == pass {
				return true, nil
			}
			return false, nil
		}))
		log.SugaredLogger.Infof("Basic Auth enabled for user: %s", user)
	} else {
		log.SugaredLogger.Warn("Basic Auth NOT enabled. Set APP_USER and APP_PASSWORD to enable.")
	}

	// SSE Endpoint for Wails Events
	e.GET("/events", func(c echo.Context) error {
		w := c.Response().Writer
		c.Response().Header().Set(echo.HeaderContentType, "text/event-stream")
		c.Response().Header().Set(echo.HeaderCacheControl, "no-cache")
		c.Response().Header().Set(echo.HeaderConnection, "keep-alive")

		// Register this client
		eventChan := util.RegisterListener()
		defer util.UnregisterListener(eventChan)

		// If already ready, send "done" immediately to unblock frontend
		if util.IsReady {
			dataBytes, _ := json.Marshal(util.WebEvent{Name: "loadingMsg", Data: "done"})
			fmt.Fprintf(w, "data: %s\n\n", dataBytes)
			w.(http.Flusher).Flush()
		}

		// Stream events from channel to SSE
		for {
			select {
			case event, ok := <-eventChan:
				if !ok {
					return nil
				}
				dataBytes, _ := json.Marshal(event)
				fmt.Fprintf(w, "data: %s\n\n", dataBytes)
				w.(http.Flusher).Flush()
			case <-c.Request().Context().Done():
				return nil
			}
		}
	})

	// Generic API handler to call App methods via reflection
	e.POST("/api/:method", func(c echo.Context) error {
		methodName := c.Param("method")
		// Capitalize first letter to match Go exported methods
		methodName = strings.ToUpper(methodName[:1]) + methodName[1:]

		method := reflect.ValueOf(app).MethodByName(methodName)
		if !method.IsValid() {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Method not found"})
		}

		// Parse arguments from body
		var args []interface{}
		if c.Request().ContentLength > 0 {
			if err := json.NewDecoder(c.Request().Body).Decode(&args); err != nil {
				return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid arguments: " + err.Error()})
			}
		}

		// Prepare call arguments
		methodType := method.Type()
		if len(args) != methodType.NumIn() {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": fmt.Sprintf("Wrong number of arguments. Expected %d, got %d", methodType.NumIn(), len(args)),
			})
		}

		log.SugaredLogger.Infof("Web API Call: %s with args: %v", methodName, args)

		in := make([]reflect.Value, len(args))
		for i := 0; i < len(args); i++ {
			expectedType := methodType.In(i)
			val := reflect.ValueOf(args[i])
			
			// Handle type conversion (especially for numbers from JSON)
			if val.IsValid() && val.Type().ConvertibleTo(expectedType) {
				in[i] = val.Convert(expectedType)
			} else {
				// Try complex conversion via JSON
				jsonBytes, _ := json.Marshal(args[i])
				newVal := reflect.New(expectedType).Interface()
				if err := json.Unmarshal(jsonBytes, newVal); err == nil {
					in[i] = reflect.ValueOf(newVal).Elem()
				} else {
					return c.JSON(http.StatusBadRequest, map[string]string{
						"error": fmt.Sprintf("Argument %d type mismatch: expected %v", i, expectedType),
					})
				}
			}
		}

		// Call method
		results := method.Call(in)

		log.SugaredLogger.Infof("Web API Method %s returned %d values", methodName, len(results))

		// Check for error in return values (standard Go pattern is last return value)
		if len(results) > 0 {
			lastResult := results[len(results)-1]
			if lastResult.Type().Implements(reflect.TypeOf((*error)(nil)).Elem()) {
				if !lastResult.IsNil() {
					err := lastResult.Interface().(error)
					return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
				}
			}
		}

		// Convert results
		var out []interface{}
		for _, r := range results {
			// Don't include the error in the regular result list if we already checked it
			if r.Type().Implements(reflect.TypeOf((*error)(nil)).Elem()) {
				continue
			}
			out = append(out, r.Interface())
		}

		if len(out) == 0 {
			return c.NoContent(http.StatusNoContent)
		}

		// If it's a single return, send it directly (unless it's a nil error)
		if len(out) == 1 {
			if out[0] == nil {
				return c.NoContent(http.StatusOK)
			}
			return c.JSON(http.StatusOK, out[0])
		}

		return c.JSON(http.StatusOK, out)
	})

	// Static files
	e.StaticFS("/", echo.MustSubFS(webAssets, "frontend/dist"))

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	e.Logger.Fatal(e.Start(":" + port))
}
