package users

import (
	"sync"
	"trovo-wallet-api/internal/cache"
	usersDB "trovo-wallet-api/internal/components/users/db"
	userModels "trovo-wallet-api/internal/components/users/models"

	tErrors "trovo-wallet-api/internal/errors"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

//GetUserInfo gets the user Information
func GetUserInfo(identifier string, dynamicLinkServiceUrlChan chan string, db *gorm.DB, redisCache *cache.RedisCache) (userInfo userModels.UserInfo, err error) {

	//get user from DB
	user, err := usersDB.GetUser(identifier, db)
	if err != nil {
		return userModels.UserInfo{}, err
	}
	if user.Suspended == 1 {
		return userInfo, &tErrors.ErrorUsernameIsSuspended{
			Username: user.Username,
		}
	}
	//set userInfo
	// log.Printf("[GetUserInfo] retrieved User record:[%+v]\n", user)
	userData := user.ToJSON()
	// log.Printf("[GetUserInfo] userData:[%+v]\n", userData)
	userInfo.UserData = userData

	//Get user wallet balances
	userInfo.AssetBalances = make(map[string]userModels.AssetBalances)
	assetBalances, err := GetUserWalletAssetBalances(&user, dynamicLinkServiceUrlChan, db, redisCache)
	if err == nil {
		userInfo.AssetBalances = assetBalances
	}
	// log.Println("[GetUserWalletAssetBalances] finished user wallets json")

	userInfo.ThirdPartyWalletAccess = make([]userModels.ThirdPartyWalletAccess, 0)
	//Get ThirdParty Wallet Access

	userInfo.ThirdPartyWalletAccess = user.Fetch3rdPartyWallets(db, redisCache)

	return
}

func GetUserWalletAssetBalances(user *userModels.User, dynamicLinkServiceUrlChan chan string, db *gorm.DB, redisCache *cache.RedisCache) (userWalletBalances map[string]userModels.AssetBalances, err error) {

	// userWalletBalances = make(map[string]userModels.AssetBalances)
	userWalletBalances = make(map[string]userModels.AssetBalances)
	for _, wallet := range user.UserWallets {

		var wg sync.WaitGroup
		var m sync.Mutex
		//use go routine to fetch

		var assetBalances userModels.AssetBalances
		assetBalances.Unclaimed = make(map[string]userModels.Balance)
		// assetBalances.Unclaimed = make([]userModels.Balance, 0)
		assetBalances.Claimed = make(map[string]userModels.Balance)
		// assetBalances.Claimed = make([]userModels.Balance, 0)

		wg.Add(1)
		go func(vg1 userModels.UserWallet, w *sync.WaitGroup, ml *sync.Mutex) {
			defer w.Done()
			unclaimedBalance, errR1 := vg1.GetBalance(db, true, dynamicLinkServiceUrlChan, redisCache)

			if errR1 == nil {
				//Unclaimed Assets
				ml.Lock()
				assetBalances.Unclaimed = unclaimedBalance
				ml.Unlock()

			}
		}(wallet, &wg, &m)
		wg.Add(1)
		go func(vg2 userModels.UserWallet, w *sync.WaitGroup, ml *sync.Mutex) {
			defer w.Done()
			claimedWalletBalance, errR1 := vg2.GetBalance(db, false, dynamicLinkServiceUrlChan, redisCache)

			if errR1 == nil {
				//Claimed Assets
				ml.Lock()
				assetBalances.Claimed = claimedWalletBalance
				ml.Unlock()
			} else {
				//get default xbn balance
				// log.Println("returning zero balance for ", vg2.ID)
				ml.Lock()
				//set Default XBN balance
				assetBalances.Claimed[":"] = userModels.Balance{
					AssetIssuer: "",
					AssetCode:   "",
					Amount:      decimal.Zero,
				}
				ml.Unlock()
			}

		}(wallet, &wg, &m)
		wg.Wait()
		// log.Printf("[GetUserWalletAssetBalances] finished user wallet balance:[%+v]\n", assetBalances)
		m.Lock()
		userWalletBalances[wallet.ID] = assetBalances
		m.Unlock()

		// log.Println("exited inner wait")

	}
	// log.Println("done...")

	return
}
