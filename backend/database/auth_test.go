package database_test

//func TestHasPermission(t *testing.T) {
//err := godotenv.Load(testEnvPath)
//require.NoError(t, err)

//logger := newTestLogger()
//db := backend.NewTestDatabase(logger)

//u := createTestUser(t, db)

//permission, err := db.HasPermission(u.Uid, "delete_user")
//require.NoError(t, err)
//require.False(t, permission)

//err = db.CreateUserRoles(u.Username, []string{"admin"})
//require.NoError(t, err)

//permission, err = db.HasPermission(u.Uid, "delete_user")
//require.NoError(t, err)
//require.True(t, permission)

//err = db.DeleteUser(u.Username)
//require.NoError(t, err)
//}
