package users

import (
	"sync"
	"trovo-wallet-api/internal/cache"
	usersDB "trovo-wallet-api/internal/components/users/db"
	userModels "trovo-wallet-api/internal/components/users/models"

	tErrors "trovo-wallet-api/internal/errors"

	"gorm.io/gorm"
)

//GetUserInfo gets the user Information
func GetUserInfo(identifier string, dynamicLinkServiceUrlChan chan string, db *gorm.DB, redisCache *cache.RedisCache) (userInfo *userModels.UserInfo, err error) {

	//get user from DB
	user, err := usersDB.GetUser(identifier, db)
	if err != nil {
		return nil, err
	}
	if user.Suspended == 1 {
		return userInfo, &tErrors.ErrorUsernameIsSuspended{
			Username: user.Username,
		}
	}
	//set userInfo
	userInfo.UserData = user.ToJSON()

	//Get user wallet balances
	GetUserWalletAssetBalances(&user, dynamicLinkServiceUrlChan, db, redisCache)

	//Get ThirdParty Wallet Access
	userInfo.ThirdPartyWalletAccess = user.Fetch3rdPartyWallets(db)

	return
}

func GetUserWalletAssetBalances(user *userModels.User, dynamicLinkServiceUrlChan chan string, db *gorm.DB, redisCache *cache.RedisCache) (userWalletBalances map[string]userModels.AssetBalances, err error) {

	userWalletBalances = make(map[string]userModels.AssetBalances)
	var rangeWG sync.WaitGroup
	var rangeMT sync.Mutex
	for _, wallet := range user.UserWallets {
		rangeWG.Add(1)
		go func(v userModels.UserWallet, m *sync.Mutex) {
			var wg sync.WaitGroup
			// var m sync.Mutex
			//use go routine to fetch

			var assetBalances userModels.AssetBalances
			assetBalances.Unclaimed = make(map[string]userModels.Balance)
			assetBalances.Claimed = make(map[string]userModels.Balance)

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
			}(v, &wg, m)
			wg.Add(1)
			go func(vg2 userModels.UserWallet, w *sync.WaitGroup, ml *sync.Mutex) {
				defer w.Done()
				claimedWalletBalance, errR1 := vg2.GetBalance(db, false, dynamicLinkServiceUrlChan, redisCache)

				if errR1 == nil {
					//Claimed Assets
					ml.Lock()
					assetBalances.Claimed = claimedWalletBalance
					ml.Unlock()
				}

			}(v, &wg, m)
			wg.Wait()

			m.Lock()
			userWalletBalances[v.ID] = assetBalances
			m.Unlock()
		}(wallet, &rangeMT)

	}
	rangeWG.Wait()

	return
}
