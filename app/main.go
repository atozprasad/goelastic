package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic"
	"github.com/teris-io/shortid"
)

const (
	elasticIndexName = "products"
	elasticTypeName  = "product"
)

type Document struct {
	ID        string    `json:"id"`
	Name     string    `json:"name"`
}

type productcategory struct {
Id      string  `json:"_id"`
Name    string  `json:"name"`
}
type DocumentRequest struct {
      Sku            int     `json:"sku"`
      Name           string  `json:"name"`
      Ptype          string  `json:"type"`
      Price          float32 `json:"price"`
      Upc            string     `json:"upc"`
    Category       []productcategory
      Shipping       float32 `json:"shipping"`
      Description    string  `json:"description"`
      Manufacturer   string  `json:"manufacturer"`
      Model           string  `json:"model"`
      Url             string    `json:"url"`
      Image           string    `json:"image"`
  }


type DocumentResponse struct {
	Name     string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type SearchResponse struct {
	Time      string             `json:"time"`
	Hits      string             `json:"hits"`
	Documents []DocumentResponse `json:"documents"`
}

var (
	elasticClient *elastic.Client
)

func main() {
	var err error
	elasticURL:="http://localhost:9200"
	// Create Elastic client and wait for Elasticsearch to be ready
	for {
		elasticClient, err = elastic.NewClient(
			elastic.SetURL(elasticURL),
			elastic.SetSniff(false),
		)
		if err != nil {
			log.Println(err)
			// Retry every 3 seconds
			time.Sleep(3 * time.Second)
		} else {
			break
		}
	}
	// Start HTTP server
	r := gin.Default()
	r.POST("/bulkupdate", createDocumentsEndpoint)
	r.GET("/search", searchEndpoint)
	if err = r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}

func createDocumentsEndpoint(c *gin.Context) {
	// Parse request
	var docs []DocumentRequest
	if err := c.BindJSON(&docs); err != nil {
		errorResponse(c, http.StatusBadRequest, "Malformed request body")
		return
	}
	// Insert documents in bulk
	bulk := elasticClient.
		Bulk().
		Index(elasticIndexName).
		Type(elasticTypeName)
	for _, d := range docs {
		doc := Document{
			ID:        shortid.MustGenerate(),
			Name:     d.Name,
		}
		bulk.Add(elastic.NewBulkIndexRequest().Id(doc.ID).Doc(doc))
	}
	if _, err := bulk.Do(c.Request.Context()); err != nil {
		log.Println(err)
		errorResponse(c, http.StatusInternalServerError, "Failed to create documents")
		return
	}
	c.Status(http.StatusOK)
}

func searchEndpoint(c *gin.Context) {
	// Parse request
	query := c.Query("query")
	if query == "" {
		errorResponse(c, http.StatusBadRequest, "Query not specified")
		return
	}
	skip := 0
	take := 10
	if i, err := strconv.Atoi(c.Query("skip")); err == nil {
		skip = i
	}
	if i, err := strconv.Atoi(c.Query("take")); err == nil {
		take = i
	}
	// Perform search
	esQuery := elastic.NewMultiMatchQuery(query, "name").
		Fuzziness("2").
		MinimumShouldMatch("2")
	result, err := elasticClient.Search().
		Index(elasticIndexName).
		Query(esQuery).
		From(skip).Size(take).
		Do(c.Request.Context())
	if err != nil {
		log.Println(err)
		errorResponse(c, http.StatusInternalServerError, "Something went wrong")
		return
	}
	res := SearchResponse{
		Time: fmt.Sprintf("%d", result.TookInMillis),
		Hits: fmt.Sprintf("%d", result.Hits.TotalHits),
	}
	// Transform search results before returning them
	docs := make([]DocumentResponse, 0)
	for _, hit := range result.Hits.Hits {
		var doc DocumentResponse
		json.Unmarshal(*hit.Source, &doc)
		docs = append(docs, doc)
	}
	res.Documents = docs
	c.JSON(http.StatusOK, res)
}

func errorResponse(c *gin.Context, code int, err string) {
	c.JSON(code, gin.H{
		"error": err,
	})
}
