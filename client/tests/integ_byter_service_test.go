package tests

// var fileHdlr = utils.NewFileHdlr()

// func TestByterService(t *testing.T) {
// 	serviceAPI := InitServiceAPI(cfg, logger)

// 	authentication(t, serviceAPI)
// 	defer auth.Identifier.SaveCurrentUser(testerUserCLI)
// 	testByterService(t, serviceAPI)

// }

// func authentication(t *testing.T, srv *api.RemoteService) {
// 	// Registration.
// 	if testerUserCLI.User.Token == "" || testerUserCLI.User.Email == "" {
// 		token, err := srv.Registration(*testerUserCLI.User)
// 		assert.NoError(t, err)
// 		testerUserCLI.User.Token = token
// 	}
// 	// Authentication.
// 	_, err := srv.Authentication(*testerUserCLI.User)
// 	assert.NoError(t, err)
// }

// func testByterService(t *testing.T, srv *api.RemoteService) {
// 	testUpload(t, srv)

// }

// func testUpload(t *testing.T, srv *api.RemoteService) {
// 	// TODO.Работать будет только локально. Пока так.
// 	modelBytes, err := model.NewBytesFromFile(fileHdlr, "/home/ali/dev/GophKeeper/client/tests/store/test.txt")
// 	assert.NoError(t, err)

// 	obj, err := srv.Upload(testerUserCLI, *modelBytes)

// 	res, ok := obj.(*rpc.UploadBytesResponse)
// 	if !ok {
// 		assert.Equal(t, true, ok)
// 	}

// 	fmt.Println(">>>", res.Status, res.UpdatedAt)
// }
