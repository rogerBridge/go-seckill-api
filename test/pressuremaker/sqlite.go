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

// refresh tokens get from users
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

// 创建ConcurrentNum个用户的token, 并写入sqlite, 便于本地使用
func CreateToken() {
	tokenList, err := GetTokenListSingle()
	if err != nil {
		logger.Fatalf("当获取token时出错: %s", err)
	}
	tokens := make([]Token, 0, ConcurrentNum)
	for i := range tokenList {
		tokens = append(tokens, Token{
			Username: tokenList[i].Username,
			Token:    tokenList[i].Token,
		})
	}
	logger.Debugf("%v", tokens[0])

	//for i := range tokens {
	//	db.Select("Username", "Token").Create(&tokens[i])
	//}

	db.CreateInBatches(tokens, 1000)
	//db.Create(&tokens)
	//db.Model(&Token{}).Create(&tokens)
}

func QueryToken() []Token {
	tokens := make([]Token, ConcurrentNum)
	db.Find(&tokens)
	return tokens
}

func GetTokenListFromSqlite() []string {
	tokenList := make([]string, 0, ConcurrentNum)
	tokens := QueryToken()
	for i := range tokens {
		tokenList = append(tokenList, tokens[i].Token)
	}
	return tokenList
}
