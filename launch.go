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
	"strings"
	"time"
	"math/rand"
	"strconv"

	/* BChain "github.com/tensor-programming/golang-blockchain/blockchain" */

	//"github.com/tensor-programming/golang-blockchain/blockchain"

	BChain "github.com/SagLara/golang-blockchain/blockchain"
	"github.com/gin-gonic/gin"
)

type User struct {
	email    string `form:"email" json:"email" binding:"required"`
	nombre   string `form:"nombre" json:"email" `
	password string `form:"password" json:"password" binding:"required"`
}

func UserToJson(user User) string {
	jsonString := "{"
	jsonString += "email:" + user.email + ","
	jsonString += "nombre:" + user.nombre + ","
	jsonString += "password:" + user.password + ","
	jsonString += "}"
	return jsonString
}

func GetMensaje(c *gin.Context) {
	bloques := "<h1>BLOQUES</h1><hr>"
	iter := chain.Iterator()
	for {
		block := iter.Next()
		bloques += "<br><ul>"
		
		bloques += "<li>Prev. Hash:"+fmt.Sprintf("%x\n", block.PrevHash)+"</li>"
		bloques += "<li>Data:"+string(block.Data)+"</li>"
		bloques += "<li>Hash:"+fmt.Sprintf("%x\n",block.Hash)+"</li>"
		pow := BChain.NewProof(block)
		bloques += "<li>PoW:"+ strconv.FormatBool(pow.Validate())+"</li>"
		if len(block.PrevHash) == 0 {
			break
		}
		bloques += "</ul>"
	}
	b := []byte(bloques)
	c.Data(http.StatusOK, "text/html; charset=utf-8", b)
}

func ExistEmail(email string) bool {
	/* emailJSON := "email:" + email + "," */
	iter := chain.Iterator()
	for {
		block := iter.Next()
		if len(block.PrevHash) == 0 {
			break
		}
		split := strings.Split(string(block.Data)[1:], ",")
		splitDato := strings.Split(split[0], ":")
		emailBlock := strings.ReplaceAll(splitDato[1]," ","")
		
		if (emailBlock == email) {
			return true
		}
		
	}
	return false
}

func PostNuevoUsuario(c *gin.Context) {

	body := c.Request.Body
	value, err := ioutil.ReadAll(body)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(404, gin.H{
			"error": "Fallo al recibir los datos",
		})
	} else {
		/* var user User */
		var dat map[string]string
		if err := json.Unmarshal(value, &dat); err != nil {
			panic(err)
		}
		var newUser User
		newUser.email = dat["email"]
		newUser.nombre = dat["nombre"]
		newUser.password = dat["password"]
		if (!ExistEmail(newUser.email)){
			err = chain.AddBlock(UserToJson(newUser))
			if err == nil {
				c.JSON(200, gin.H{
					"message": dat,
					"id":      id,
				})
				id = id + 1
			} else {
				c.JSON(404, gin.H{
					"Error": err,
				})
			}
		}else{
			c.JSON(404, gin.H{
				"mensaje": "El email ya tiene una cuenta asociada",
			})
		}
	}
}
func LoginUser(email string, password string) User {
	var user User
	user.nombre=""
	
	iter := chain.Iterator()
	for {
		block := iter.Next()
		if len(block.PrevHash) == 0 {
			break
		}
		split := strings.Split(string(block.Data)[1:], ",")
		splitDatoEmail := strings.Split(split[0], ":")
		emailBlock := strings.ReplaceAll(splitDatoEmail[1]," ","")
		
		if (emailBlock == email) {
			splitDatoPass := strings.Split(split[2], ":")
			passBlock := strings.ReplaceAll(splitDatoPass[1]," ","")
		
			if(passBlock==password){
		
				/* var user User */
				user.email=emailBlock
				user.password=passBlock
				user.nombre= strings.Split(split[1], ":")[1]
				return user
			}
		}
		
	}
	return user
}
func randToken() string {
	b := make([]byte, 8)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func PostLogin(c *gin.Context) {
	body := c.Request.Body
	value, err := ioutil.ReadAll(body)

	if err != nil {
		fmt.Println(err.Error())
		c.AbortWithStatusJSON(400, gin.H{
			"error": "Fallo con lo datos",
		})
	} else {
		/* var user User */
		var dat map[string]string
		if err := json.Unmarshal(value, &dat); err != nil {
			panic(err)
		}

		userReturn := LoginUser(dat["email"],dat["password"])
		
		if (len(userReturn.nombre)>1){
		
			if err == nil {
				c.JSON(200, gin.H{
					"nombre": userReturn.nombre,
					"idToken": randToken(),
				})
			} else {
				c.AbortWithStatusJSON(500, gin.H{
					"mensaje": err,
				})
			}
		}else{
			c.AbortWithStatusJSON(400, gin.H{
				"mensaje": "Correo o contraseña incorrecta",
			  })

			/* c.Error() */
			
			/* c.JSON(401, gin.H{
				"Error": "El usuario no pudo autenticar",
			}) */
		}
	}
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

var id int
var chain *BChain.BlockChain

func main() {

	defer os.Exit(0)

	chain = BChain.InitBlockChain()
	/* BChain.PrintChain(chain) */
	defer chain.Database.Close()

	id = 0

	r := gin.Default()
	r.Use(Cors())

	v1 := r.Group("api")
	{
		v1.GET("/mensaje", GetMensaje)
		v1.POST("/newUser", PostNuevoUsuario)
		v1.POST("/login", PostLogin)
	}

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
