package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

/* -------------------------------------------*/
// Context는 gin의 가장 중요한 부분이다.
// 미들웨어 간에 변수를 전달하고 흐름을 관리하며 요청의 JSON을 검증하고
// JSON 응답을 렌더링 할 수 있다.
func v1EndpointHandler(c *gin.Context) {
	// String 메소드는 주어진 문자열을 response Body에 적는다.
	c.String(http.StatusOK, "v1: %s %s", c.Request.Method, c.Request.URL.Path)
}

func v2EndpointHandler(c *gin.Context) {
	c.String(http.StatusOK, "v2: %s %s", c.Request.Method, c.Request.URL.Path)
}

/* -------------------------------------------*/

/* -------------------------------------------*/

type AddParams struct {
	X float64 `json:"x"` // 이 경우에서 json package는 `json:"x"`를 사용하여
	Y float64 `json:"y"` // X의 값을 해당 json객체의 x값으로 인코딩한다.
}

func add(c *gin.Context) {
	var ap AddParams

	// shouldBindJSON 메소드는 전달받은 구조체 포인터를 JSON형식으로 바인딩한다.
	// 바인딩이란? -> 프로그램의 어떤 기본단위가 가지는 무엇을 결정짓는 것
	// 여기서는 구조체 포인터 타입인 &ap를 JSON타입으로 결정한 것임
	if err := c.ShouldBindJSON(&ap); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Calculator error"})
		return
	}

	// JSON 메소드는 주어진 구조체(인터페이스)를 JSON형식으로 변환한다.
	c.JSON(http.StatusOK, gin.H{"answer": ap.X + ap.Y})
}

/* -------------------------------------------*/

/* -------------------------------------------*/
type Product struct {
	Id int `json:"id" xml:"Id" yaml:"id"`
	// Id 값을 해당 json 객체의 id 값으로 인코딩한다.
	// Id 값을 해당 xml 객체의 Id 값으로 인코딩한다.
	// Id 값을 해당 yaml 객체의 id 값으로 인코딩한다.
	Name string `json:"name" xml:"Name" yaml:"name"`
}

/* -------------------------------------------*/

/* -------------------------------------------*/
type PrintJob struct {
	// gin은 struct-tag를 기반으로 한 validation 기능을 제공한다.
	// binding: ... 바인딩 시에 그다음 validator(...)를 위반하면 오류를 반환한다.
	// required는 null값이 들어와선 안된다고 표현하는 것이며
	// gte는 greater than OR equald을 말하는 것이라 10000이상을 적어야한다고 표현하는 것이며
	// lte less than OR equal을 말하는 것이라 100이하를 적어야 한다는 뜻이다.
	JobId int `json:"jobId" binding:"required,gte=10000"`
	Pages int `json:"pages" binding:"required,gte=1,lte=100"`
}

/* -------------------------------------------*/

/* -------------------------------------------*/
// FindUserAgent는 사용자 정의 미들웨어다
// 미들웨어란 양쪽을 연결하여 데이터를 주고받을 수 있도록
// 중간에서 매개 역할을 하는 소프트 웨어다.
func FindUserAgent() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println(c.GetHeader("User-Agent"))
		c.Next()
	}
}

/* -------------------------------------------*/

/* -------------------------------------------*/
type PrintJobs struct {
	Format    string `json:"format" binding:"required"`
	InvoiceId int    `json:"invoiceid" binding:"required,gte=0"`
	JobId     int    `json:"jobid" binding:"gte=0"`
}

type Invoice struct {
	InvoiceId   int    `json:"invoiceid"`
	CustomerId  int    `json:"customerid" binding:"required,gte=0"`
	Price       int    `json:"price" binding:"required,gte=0"`
	Description string `json:"description" binding:"required"`
}

func createPrintJob(invoiceId int) {
	client := resty.New()
	var p PrintJobs
	_, err := client.R().
		SetBody(PrintJobs{Format: "A4", InvoiceId: invoiceId}).
		SetResult(&p).
		Post("http://localhost:5000/print-jobs")

	if err != nil {
		log.Println("InvoiceGenerator: unable to connect PrinterService")
		return
	}
	log.Printf("InvoiceGenerator: created print job #%v via PrinterService", p.JobId)
}

/* -------------------------------------------*/

func main() {
	// Defualt gin router를 사용
	r := gin.Default()

	/* -------------------------------------------*/
	r.POST("/print-jobs", func(c *gin.Context) {
		var p PrintJobs
		if err := c.ShouldBindJSON(&p); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input!"})
			return
		}

		log.Printf("PrintService: creating new print job from invoice #%v...", p.InvoiceId)
		rand.Seed(time.Now().UnixNano())
		p.JobId = rand.Intn(1000)
		log.Printf("PrintService: created print job #%v", p.JobId)
		c.JSON(http.StatusOK, p)
	})

	r.POST("/invoices", func(c *gin.Context) {
		var iv Invoice
		if err := c.ShouldBindJSON(&iv); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input!"})
			return
		}

		log.Println("InvoiceGenerator: creating new invoice...")
		rand.Seed(time.Now().UnixNano())
		iv.InvoiceId = rand.Intn(1000)
		log.Printf("InvoiceGenerator: created invoice #%v", iv.InvoiceId)

		createPrintJob(iv.InvoiceId)
		c.JSON(http.StatusOK, iv)
	})
	/* -------------------------------------------*/

	/* -------------------------------------------*/
	r.Use(FindUserAgent())
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Middleware works!"})
	})
	/* -------------------------------------------*/

	/* -------------------------------------------*/
	// Use 메소드는 첫번째 인자로 미들웨어를 받는데 이것을 라우터에 연결한다.
	// 이 메소드를 통해 연결된 미들웨어는 모든 요청에 대해 핸들러 체인에 포함된다.
	r.Use(cors.Default())
	r.GET("/cors", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "CORS works!"})
	})
	/* -------------------------------------------*/

	/* -------------------------------------------*/
	r.POST("/print", func(c *gin.Context) {
		var p PrintJob
		if err := c.ShouldBindJSON(&p); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input!"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("PrintJob #%v started!", p.JobId)})
	})
	/* -------------------------------------------*/

	/* -------------------------------------------*/
	r.GET("/productJSON", func(c *gin.Context) {
		product := Product{1, "Apple"}
		c.JSON(http.StatusOK, product)
	})

	r.GET("/productXML", func(c *gin.Context) {
		product := Product{2, "Banana"}
		c.XML(http.StatusOK, product)
	})
	r.GET("/productYAML", func(c *gin.Context) {
		product := Product{3, "Mango"}
		c.JSON(http.StatusOK, product)
	})
	/* -------------------------------------------*/

	/* -------------------------------------------*/
	r.POST("/add", add)
	/* -------------------------------------------*/

	/* -------------------------------------------*/
	// 동일 경로prefix가있는 모든 경로를 그룹화 한다.
	v1 := r.Group("v1")

	v1.GET("/products", v1EndpointHandler)
	v1.GET("/poducts/:productId", v1EndpointHandler)
	v1.POST("/products", v1EndpointHandler)
	v1.PUT("/products/:productId", v1EndpointHandler)
	v1.DELETE("/products/:productId", v1EndpointHandler)

	v2 := r.Group("v2")

	v2.GET("/products", v2EndpointHandler)
	v2.GET("/poducts/:productId", v2EndpointHandler)
	v2.POST("/products", v2EndpointHandler)
	v2.PUT("/products/:productId", v2EndpointHandler)
	v2.DELETE("/products/:productId", v2EndpointHandler)
	/* -------------------------------------------*/

	r.Run(":5000")
}
