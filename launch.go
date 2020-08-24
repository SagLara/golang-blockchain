package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	BChain "github.com/tensor-programming/golang-blockchain/blockchain"

	//"github.com/tensor-programming/golang-blockchain/blockchain"

	//BlockChainPrivada "github.com/SagLara/golang-blockchain/blockchain"
	"github.com/gin-gonic/gin"
)

type User struct {
	email    string `form:"email" json:"email" binding:"required"`
	nombre   string `form:"nombre" json:"email" `
	password string `form:"password" json:"password" binding:"required"`
}

func HomePage(c *gin.Context) {
	body := c.Request.Body
	value, err := ioutil.ReadAll(body)
	if err != nil {
		fmt.Println(err.Error())
	}
	c.JSON(200, gin.H{
		"message": string(value),
	})
}
func GetMensaje(c *gin.Context) {
	fmt.Println("hola mundo")
	c.JSON(200, gin.H{
		"message": "Un mensaje",
	})
}

func PostHomePage(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Post Home Page",
	})
}

func PostNuevoUsuario(c *gin.Context) {

	body := c.Request.Body
	value, err := ioutil.ReadAll(body)

	fmt.Println("Resgistrando ando")
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(404, gin.H{
			"error": "Fallo al crear",
		})
	} else {
		/* var user User */
		var dat map[string]interface{}
		if err := json.Unmarshal(value, &dat); err != nil {
			panic(err)
		}
		fmt.Println("Email:", dat["email"])
		fmt.Println("Nombre :", dat["nombre"])
		fmt.Println("Pasword :", dat["password"])
		c.JSON(200, gin.H{
			"message": dat,
			"id":      id,
		})
		id = id + 1
	}
}

func PostLogin(c *gin.Context) {
	fmt.Println("Inicio servicio postLogin")
	fmt.Println(c.Request.Body)
	var json User
	if c.BindJSON(&json) == nil {
		c.JSON(200, gin.H{
			"nombre": json.email,
		})

	} else {
		c.JSON(404, gin.H{
			"mensssage": "error",
		})
	}
}

func QueryStrings(c *gin.Context) {
	name := c.Query("name")
	age := c.Query("age")

	c.JSON(200, gin.H{
		"name": name,
		"age":  age,
	})
}
func PathParameters(c *gin.Context) {
	name := c.Param("name")
	age := c.Param("age")

	c.JSON(200, gin.H{
		"name": name,
		"age":  age,
	})
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200/")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Origin, Authorization, Accept, Client-Security-Token, Accept-Encoding, x-access-token")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		c.Next()
	}
}

//cd go\projects\Pb-server-basic
var id int
var chain *BChain.BlockChain

//var chain = BChain.InitBlockChain()
//var chain BChain.InitBlockChain()

func main() {

	os.OpenFile("./tmp/blocks/LOCK", os.O_EXCL, 0)
	//defer chain.Database.Close()
	//BlockChainPrivada.Init()
	//BlockChainPrivada.AddBlockMain("asdasdas")
	defer os.Exit(0)

	chain = BChain.InitBlockChain()

	BChain.PrintChain(chain)
	//chain := blockchain.InitBlockChain()
	defer chain.Database.Close()

	id = 0
	fmt.Println("Hola mundo")

	//runtime.Goexit()

	r := gin.Default()
	r.Use(Cors())

	v1 := r.Group("api")
	{
		v1.GET("/", HomePage)
		v1.GET("/mensaje", GetMensaje)
		v1.POST("/", PostHomePage)
		v1.POST("/newUser", PostNuevoUsuario)
		v1.POST("/login", PostLogin)
		v1.GET("/query", QueryStrings)             //http://localhost:8080/query?name=david&age=21
		v1.GET("/path/:name/:age", PathParameters) //http://localhost:8080/path/david/21/
		v1.OPTIONS("/", PostHomePage)
	}
	//r.Run(":8080")

	chain.AddBlock("Prueba global")
	BChain.PrintChain(chain)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("listen: %s\n", err)
		}

	}()

	//Espera que se interrumpa el servicio para darle un apagado elegante y correcto
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")

}
