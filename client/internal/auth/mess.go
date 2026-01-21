package auth

// type Mess struct {
// 	Client *client.ClientCLI
// 	User *user.UserCLI
// }

// const ATTEMPTS = 3

// func GetPassword(client *client.ClientCLI, user *user.UserCLI) (string, error) {
// 	client.SendMess("Enter password...")
// 	return user.ReceiveMess()
// }

// func GetEmail(client *client.ClientCLI, user *user.UserCLI) (string, error) {
// 	client.SendMess("Enter email...")
// 	return user.ReceiveMess()
// }

// type GetterFn func(*client.ClientCLI, *user.UserCLI) (string, error)
// type CheckerFn func(*user.UserCLI, string) bool

// func TryToGetSeveralTimes(get GetterFn, check CheckerFn) GetterFn {
// 	return func(client *client.ClientCLI, user *user.UserCLI) (string, error) {

// 		for repeat := 0; repeat < ATTEMPTS; repeat++ {
// 			result, err := get(client, user)

// 			if err == nil && check(user, result) {
// 				return result, nil
// 			}
// 			client.SendMess("Uncorrected credentials. Try again...")
// 		}
// 		return "", errs.ErrUncorrectCredentials
// 	}
// }
