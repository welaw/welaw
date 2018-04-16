package database_test

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/welaw/welaw/backend/database"
	"github.com/welaw/welaw/proto"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func createTestUser(t *testing.T, db database.Database) *proto.User {
	seq := randSeq(6)
	u := &proto.User{
		Username:     "test-" + seq,
		FullName:     "test full name",
		Email:        "test-" + seq + "@welaw.org",
		EmailPrivate: false,
		Biography:    "test biography",
		PictureUrl:   "/assets/test.png",
		ProviderId:   "123456",
		Upstream:     "",
	}
	u, err := db.CreateUser(u)
	require.NoError(t, err)
	return u
}

//func TestGetUserByProviderID(t *testing.T) {
//err := godotenv.Load(testEnvPath)
//require.NoError(t, err)

//logger := newTestLogger()
//db := backend.NewTestDatabase(logger)

//u := createTestUser(t, db)

//err = db.DeleteUser(u.Username)
//require.NoError(t, err)

//u, err = db.GetUserByProviderId(u.ProviderId)
//require.NoError(t, err)

//fmt.Printf("test user: %v\n", u)
//}
