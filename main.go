package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// Monitoring, logging, -> apakah proses berhasil, response seperti apa
// analisis, apakah prosesnya cepat, kalau lama itu di codingan mana, dst

// analitic, google analitic

// Auth

// rate limiter, user ini batasi max request 10 req per detik

// logic bisnis

// membuat logic lebih efisien dengan middleware
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		timeNow := time.Now()

		c.Set("example", "123456")

		// forward / melanjutkan ke handler, middleware selanjutnya
		c.Next()
		log.Println("Time elapsed: ", time.Since(timeNow))
	}
}

func Auth(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// ngirim query url dengan key berupa token, value
		// request
		token := c.GetHeader("Authorization")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
			return
		}

		// validasi token
		claims, err := ValidateTokenJWT(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
			return
		}

		if role == "" || claims.Role == "admin" {
			c.Next()
			return
		}

		if claims.Role != role {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "role not valid"})
			return
		} else {
			c.Next()
		}
		// response
	}
}

// key value
// init middleware
// lempar data dari middleware ke handler
// lempar middleware ke middleware lain

func main() {
	r := gin.Default()

	// set middleware global
	r.Use(Logger())

	// set middleware di group
	v1 := r.Group("/v1", Logger())
	{

		v1.POST("/login", func(c *gin.Context) {
			// ngirim query / body req username, password

			// username = admin, password = password

			username := c.Query("username")
			password := c.Query("password")

			if username != "admin" || password != "password" {
				c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
				return
			}

			// generate token
			token, err := GenerateTokenJWT(username, "pengelola")
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "error generate token JWT"})
				return
			}

			c.JSON(http.StatusOK, gin.H{"token": token})
		})

		v1.GET("/example", func(c *gin.Context) {
			example := c.MustGet("example").(string)

			c.JSON(http.StatusOK, gin.H{"example": example})
		})

		v1.GET("/hello", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "hello"})
		})

		v1.GET("/index", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{"message": "index"})
		})

		v1.GET("/admin", Auth("admin"), func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{"message": "hello admin"})
		})

		v1.GET("/user", Auth("user"), func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{"message": "hello user"})
		})

		v1.GET("/pengelola", Auth("pengelola"), func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "hello pengelola"})
		})
	}

	v1Private := r.Group("/v1/private", Logger(), Auth(""))
	{
		v1Private.GET("/home", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "home"})
		})
	}

	//

	r.Run(":8080")
}

// helper function

type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

var (
	jwtKey = []byte("iniadalahrahasiapenting")
)

// generate token JWT
func GenerateTokenJWT(username string, role string) (string, error) {
	expirationTime := time.Now().Add(5 * time.Minute) // 5 menit

	// Buat claims berisi data username dan role yang akan kita embed ke JWT
	claims := &Claims{
		Username: username,
		Role:     role,
		StandardClaims: jwt.StandardClaims{
			// expiry time menggunakan time millisecond
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Buat token menggunakan encoded claim dengan salah satu algoritma yang dipakai
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Buat JWT string dari token yang sudah dibuat menggunakan JWT key yang telah dideklarasikan (proses encoding JWT)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		// return internal error ketika ada kesalahan saat pembuatan JWT string
		return "", err
	}

	return tokenString, nil
}

// validasi token JWT
func ValidateTokenJWT(tknStr string) (*Claims, error) {
	claims := &Claims{}

	// parse JWT token ke dalam claims
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return nil, err
	}

	if !tkn.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
