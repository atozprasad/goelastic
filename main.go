package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
	"io/ioutil"
	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic"
	"github.com/teris-io/shortid"
)

const (
	elasticIndexName = "products"
	elasticTypeName  = "product"
	elasticURL1 ="http://35.202.99.46:9200"
	elasticURL2 ="http://35.192.32.150:9200"
	elasticURL3 ="http://35.224.21.162:9200"
	elasticUser ="elastic"
	elasticPwd ="hKVd9xXQ"
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
      Description    string  `json:"description"`
      Manufacturer   string  `json:"manufacturer"`
      Model           string  `json:"model"`
      Url             string    `json:"url"`
      Image           string    `json:"image"`
			Shipping       float32 `json:"shipping"`
  }


type DocumentResponse struct {
	Name     string    `json:"name"`
//	CreatedAt time.Time `json:"created_at"`
}
type SearchResponse struct {
	Time      string             `json:"time"`
	Hits      string             `json:"hits"`
	Documents []DocumentResponse `json:"documents"`
}

//####################################Autocomplete#######################################

	type AutocompleteDocument struct {
			Score		 float32		`json:"score"`
			Text		 string			`json:"text"`
			Name		 string			`json:"name"`
	    AutocompletePayload struct  {
				FirstName	string `json:"first_name"`
		 		ID				int 	 `json:"id"`
		 		LastName	string `json:"last_name"`
			}  `json:"payload"`
	}
 	type AutocompleteResponse struct {
 		Objects []AutocompleteDocument
 }
//####################################Autocomplete#######################################
var (
	elasticClient *elastic.Client
	)
	var skip = 0
	var take = 10

var dataFilePath string="./products.json"

func main() {
	var err error

	// Create Elastic client and wait for Elasticsearch to be ready
	for {

		// client, err := elastic.NewClient(
		//   elastic.SetURL("http://127.0.0.1:9200", "http://127.0.0.1:9201"),
		//   elastic.SetBasicAuth("user", "secret"))


		elasticClient, err = elastic.NewClient(
			elastic.SetURL(elasticURL1,elasticURL2,elasticURL3),
			elastic.SetBasicAuth(elasticUser, elasticPwd),
			elastic.SetSniff(false))

		if err != nil {
			log.Println(err)
			// Retry every 3 seconds
			if elastic.IsConnErr(err) !=true {
				log.Println("Yes it is connection Error")
			}

			time.Sleep(6 * time.Second)
		} else {
			break
		}
	}
	// Start HTTP server
	r := gin.Default()
	r.POST("/bulkupdload", bulkUploadEndpoint)
	r.POST("/fileupload", createFromFileEndpoint)
	r.GET("/search", searchEndpoint)
	r.GET("/autocomplete", autocompleteEndpoint)


	if err = r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}


func rootService(c *gin.Context) {
	c.Param("article_id")
	errorResponse(c, http.StatusBadRequest, c.Request.URL.Path)
	return

}


func bulkUploadEndpoint(c *gin.Context) {
	// Parse request
		var docs []DocumentRequest
		if err := c.BindJSON(&docs); err != nil {
			errorResponse(c, http.StatusBadRequest, "Malformed request body")
			return
		}
		err:=AddProductsToIndex(docs,c,elasticClient)
		 if err != nil {
			errorResponse(c, http.StatusBadRequest, "Failed to post data in Elastic")
 			return
		 }
		 c.Status(http.StatusOK)
	 }

func createFromFileEndpoint(c *gin.Context) {
		raw, err := ioutil.ReadFile(dataFilePath)
	 if err != nil {
		    errorResponse(c, http.StatusBadRequest, "Failed to read data from file. [Reason] " + err.Error())
			 return
	 }
	 var docs []DocumentRequest
	 err = json.Unmarshal(raw, &docs)
	 if err != nil {
		   errorResponse(c, http.StatusBadRequest, "Failed to convert data to json. [Reason] " + err.Error())
			 return
	 }

		if err != nil {
			 errorResponse(c, http.StatusBadRequest, "Failed to read data from file. [Reason] " + err.Error())
			 return
		}
		err=AddProductsToIndex(docs,c,elasticClient)
		 if err != nil {
			errorResponse(c, http.StatusBadRequest, "Failed to post data in Elastic. [Reason] "+err.Error())
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
func autocompleteEndpoint(c *gin.Context) {
	// Parse request
		query := c.Query("query")
		if query == "" {
			errorResponse(c, http.StatusBadRequest, "Query not specified")
			return
		}
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
		res :=[]AutocompleteDocument{}
		// Transform search results before returning them
		docs := make([]AutocompleteDocument, 0)
		i:=0
		for _, hit := range result.Hits.Hits {
			i++
			var doc AutocompleteDocument
			json.Unmarshal(*hit.Source, &doc)
			doc.Text=doc.Name
			doc.Score=1.0
	    doc.AutocompletePayload.FirstName=""
			doc.AutocompletePayload.ID=i
			doc.AutocompletePayload.LastName=""
			docs = append(docs, doc)
		}
		res= docs
		c.JSON(http.StatusOK, res)
}
func errorResponse(c *gin.Context, code int, err string) {
		c.JSON(code, gin.H{
			"error": err,
		})
}


func AddProductsToIndex(docs []DocumentRequest,c *gin.Context, elasticClient *elastic.Client ) error {
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
			errorResponse(c, http.StatusInternalServerError, "Failed to create documents  [Reason] " + err.Error())
			return err
		}
		return nil
}
