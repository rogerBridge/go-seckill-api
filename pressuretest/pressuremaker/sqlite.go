package pressuremaker

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Token struct {
	gorm.Model
	Username string `json:"username"`
	Token    string `json:"token"`
}

//func connectToSqlite() (*gorm.DB, error) {
//	var db, err = gorm.Open(sqlite.Open("./pressuremaker/token.db"), &gorm.Config{})
//	if err != nil {
//		return db, err
//	}
//	return db, nil
//}

var db, errConnectToSqlite = gorm.Open(sqlite.Open("./pressuremaker/token.db"), &gorm.Config{})

func InitSqlite() {
	db.Exec("DELETE FROM tokens")

	if errConnectToSqlite != nil {
		logger.Fatalf("While connect to sqlite, error: %s", errConnectToSqlite)
	}
	err := db.AutoMigrate(&Token{})
	if err != nil {
		logger.Fatalf("While Migrate sqlite, error: %s", err)
	}
}

func CreateToken() {
	tokenList, err := GetTokenListSingle()
	if err != nil {
		logger.Fatalf("当获取token时出错: %s", err)
	}
	tokens := make([]Token, 0, 10000)
	for i := range tokenList {
		tokens = append(tokens, Token{
			Username: tokenList[i].Username,
			Token:    tokenList[i].Token,
		})
	}
	logger.Infoln("%+v", tokens[0])
	for i := range tokens {
		db.Create(&tokens[i])
	}
	//db.Create(&tokens)
	//db.Model(&Token{}).Create(&tokens)
}

func QueryToken() []Token {
	tokens := make([]Token, 10000)
	db.Find(&tokens)
	return tokens
}

func GetTokenListFromSqlite() []string {
	tokenList := make([]string, 0, 10000)
	tokens := QueryToken()
	for i := range tokens {
		tokenList = append(tokenList, tokens[i].Token)
	}
	return tokenList
}
