package main

import (
	"log"
	"net/http"
	"strings"

	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	docs "github.com/sajid-khan-js/k8s-adventures/go-app/docs"
	"github.com/sajid-khan-js/k8s-adventures/go-app/modules/k8sclient"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

func setupRouter() *gin.Engine {

	router := gin.Default()

	router.Use(DummyMiddleware)

	router.GET("/namespaces", getNamespaces)
	router.GET("/namespaces/:name", getNamespaceByName)
	router.POST("/namespaces", postNamespace)

	router.GET("/healthz", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// Mock readiness e.g. app might need to connect to db, load data, cache warm etc.
	// https://blog.gopheracademy.com/advent-2017/kubernetes-ready-service/
	isReady := &atomic.Value{}
	isReady.Store(false)
	// go routine, not blocking
	go func() {
		log.Printf("Readyz probe is negative by default...")
		time.Sleep(15 * time.Second)
		isReady.Store(true)
		log.Printf("Readyz probe is positive.")
	}()
	router.GET("/readyz", func(c *gin.Context) {
		if isReady == nil || !isReady.Load().(bool) {
			c.Status(http.StatusServiceUnavailable)
		} else {
			c.Status(http.StatusOK)
		}

	})

	// https://prometheus.io/docs/guides/go-application/ and https://gabrieltanner.org/blog/collecting-prometheus-metrics-in-golang/
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// localhost:8080/swagger/index.html
	docs.SwaggerInfo.BasePath = "/"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}

// getNamespaces godoc
// @Summary List namespaces in a K8s cluster.
// @Description Get all namespaces and all Pods in a cluster.
// @Produce json
// @Success 200 {object} []Namespace
// @Router /namespaces/ [get]
// curl -v -L localhost:8080/namespaces
func getNamespaces(c *gin.Context) {

	clientSet, err := k8sclient.InitClient()
	if err != nil {
		log.Print(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong on our side"})
		return
	}

	rawNamespaces, err := k8sclient.GetNamespaces(*clientSet)
	if err != nil {
		log.Print(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong on our side"})
		return
	}

	var namespaces []Namespace

	for _, ns := range rawNamespaces {
		var namespace Namespace
		namespace.Name = ns

		// TODO share this code with getNamespaceByName
		podsInNamespace, err := k8sclient.GetPods(*clientSet, ns)
		if err != nil {
			log.Print(err)
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong on our side"})
			return
		}

		for i, n := range podsInNamespace {
			var p Pod
			p.Name = i
			p.Status = n
			namespace.Pods = append(namespace.Pods, p)
		}

		namespaces = append(namespaces, namespace)
	}

	// Return serialized namespaces
	c.IndentedJSON(http.StatusOK, namespaces)
}

// getNamespaceByName godoc
// @Summary Get a K8s namespace and it's Pods.
// @Description Get Pods in a namespace.
// @Param        name   path      string  true  "Namespace name"
// @Produce json
// @Success 200 {object} Namespace
// @Router /namespaces/{name} [get]
// curl -v -L localhost:8080/namespaces/default
func getNamespaceByName(c *gin.Context) {

	name := c.Param("name")

	clientSet, err := k8sclient.InitClient()
	if err != nil {
		log.Print(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong on our side"})
		return
	}

	podsInNamespace, err := k8sclient.GetPods(*clientSet, name)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "not found"):
			log.Print(err)
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Namespace '" + name + "' not found"})
			return
		default:
			log.Print(err)
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong on our side"})
			return
		}
	}

	if len(podsInNamespace) > 0 {
		var ns Namespace
		ns.Name = name
		for i, n := range podsInNamespace {
			var p Pod
			p.Name = i
			p.Status = n
			ns.Pods = append(ns.Pods, p)
		}

		c.IndentedJSON(http.StatusOK, ns)

	} else {
		// Probably should just return the empty slice of Pods to be consistent, but leaving this here for demo purposes
		c.IndentedJSON(http.StatusOK, gin.H{"message": "Namespace '" + name + "' has no Pods"})
	}

}

// postNamespace godoc
// @Summary Create a new namespace.
// @Description Create a new namespace in the cluster.
// @Param        name   body      Namespace  true  "Namespace name"
// @Accept  json
// @Produce json
// @Success 200 {object} Namespace
// @Router /namespaces [post]
/*
curl -v -L 'localhost:8080/namespaces/' \
--header 'Content-Type: application/json' \
--data-raw '{
    "name": "mynamespace"
}'
*/
func postNamespace(c *gin.Context) {

	var newNamespace Namespace

	// Validate based on JSON struct tags in Namespace struct
	if err := c.ShouldBindJSON(&newNamespace); err != nil {
		// TODO differentiate between validation errors e.g. missing required key (name) and incorrect value (namespace name not RFC 1123)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Namespace name must comply with RFC 1123"})
		return
	}

	name := newNamespace.Name

	clientSet, err := k8sclient.InitClient()
	if err != nil {
		log.Print(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong on our side"})
		return
	}

	err = k8sclient.CreateNamespace(*clientSet, name)

	if err != nil {
		switch {
		case strings.Contains(err.Error(), "already exists"):
			log.Print(err)
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Namespace '" + name + "' already exists"})
			return
		default:
			log.Print(err)
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong on our side"})
			return
		}
	}

	c.IndentedJSON(http.StatusCreated, newNamespace)
}
