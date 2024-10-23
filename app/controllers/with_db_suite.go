package controllers

import (
	"app/db"
	"app/repositories"
	"app/services"
	"bytes"
	"database/sql"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/DATA-DOG/go-txdb"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type WithDBSuite struct {
	suite.Suite
}

var (
	DBCon *gorm.DB
	token string
)

// func (s *WithDBSuite) SetupSuite()                           {} // テストスイート実施前の処理
// func (s *WithDBSuite) TearDownSuite()                        {} // テストスイート終了後の処理
// func (s *WithDBSuite) SetupTest()                            {} // テストケース実施前の処理
// func (s *WithDBSuite) TearDownTest()                         {} // テストケース終了後の処理
// func (s *WithDBSuite) BeforeTest(suiteName, testName string) {} // テストケース実施前の処理
// func (s *WithDBSuite) AfterTest(suiteName, testName string)  {} // テストケース終了後の処理

func init() {
	txdb.Register("txdb-controller", "mysql", db.GetDsn())
}

func (s *WithDBSuite) SetDBCon() {
	db, err := sql.Open("txdb-controller", "connect")
	if err != nil {
		log.Fatalln(err)
	}

	DBCon, err = gorm.Open(mysql.New(mysql.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		s.T().Fatalf("failed to initialize GORM DB: %v", err)
	}
}

func (s *WithDBSuite) CloseDB() {
	db, _ := DBCon.DB()
	db.Close()
}

func (s *WithDBSuite) SignIn() {
	userRepository := repositories.NewUserRepository(DBCon)
	authService := services.NewAuthService(userRepository)
	authController := NewAuthController(authService)

	// gin contextの生成
	authRecorder := httptest.NewRecorder()
	ginContext, _ := gin.CreateTestContext(authRecorder)

	// NOTE: リクエストの生成
	body := bytes.NewBufferString("{\"email\":\"test@example.com\",\"password\":\"password\"}")
	req, _ := http.NewRequest("POST", "/auth/sign_in", body)
	req.Header.Set("Content-Type", "application/json")
	ginContext.Request = req

	// NOTE: ログインし、tokenに認証情報を格納
	authController.SignIn(ginContext)
	token = authRecorder.Result().Cookies()[0].Value
}