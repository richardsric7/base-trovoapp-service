package users

import (
	"log"
	"sync"
	"trovo-wallet-api/internal/cache"
	usersDB "trovo-wallet-api/internal/components/users/db"
	userModels "trovo-wallet-api/internal/components/users/models"
	"trovo-wallet-api/internal/middleware"

	tErrors "trovo-wallet-api/internal/errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

//GetUserInfo gets the user Information
func GetUserInfo(identifier string, dynamicLinkServiceUrlChan chan string, db *gorm.DB, redisCache *cache.RedisCache, c *gin.Context) (userInfo userModels.UserInfo, err error) {

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
	var owner bool
	for _, wallet := range user.UserWallets {
		if wallet.Tag == nil {
			//primary wallet
			owner = wallet.Signer == middleware.ExtractSigner(c)
		}
	}

	//set userInfo
	// log.Printf("[GetUserInfo] retrieved User record:[%+v]\n", user)
	userData := user.ToJSON()
	// log.Printf("[GetUserInfo] userData:[%+v]\n", userData)
	userInfo.UserData = userData

	if !owner {
		userInfo.UserData.UserWallets = nil
		userInfo.AssetBalances = nil
		userInfo.ThirdPartyWalletAccess = nil
		userInfo.DefaultAssets = nil
	}

	//Get user wallet balances
	if owner {
		userInfo.AssetBalances = make(map[string]userModels.AssetBalances)
		assetBalances, err := GetUserWalletAssetBalances(&user, dynamicLinkServiceUrlChan, db, redisCache)
		if err == nil {
			userInfo.AssetBalances = assetBalances
		}
		// log.Println("[GetUserWalletAssetBalances] finished user wallets json")

		userInfo.ThirdPartyWalletAccess = make([]userModels.ThirdPartyWalletAccess, 0)
		//Get ThirdParty Wallet Access

		userInfo.ThirdPartyWalletAccess = user.Fetch3rdPartyWallets(db, redisCache)

		userInfo.DefaultAssets = user.GetDefaultAssets(db, redisCache)
	}

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
		// assetBalances.Unclaimed = make(map[string]userModels.Balance)
		assetBalances.Unclaimed = make([]userModels.Balance, 0)
		// assetBalances.Claimed = make(map[string]userModels.Balance)
		assetBalances.Claimed = make([]userModels.Balance, 0)

		wg.Add(1)
		go func(vg1 userModels.UserWallet, w *sync.WaitGroup, ml *sync.Mutex) {
			defer w.Done()
			unclaimedBalance, errR1 := vg1.GetSortedUserBalance(db, true, dynamicLinkServiceUrlChan, redisCache)

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
			claimedWalletBalance, errR1 := vg2.GetSortedUserBalance(db, false, dynamicLinkServiceUrlChan, redisCache)

			if errR1 != nil {
				//log server error
				log.Printf("[GetUserWalletAssetBalances] error getting claimed wallet balance for user:[%s] wallet:[%s] error:[%+v]\n", user.Username, vg2.ID, errR1)

			}
			//Claimed Assets
			ml.Lock()
			assetBalances.Claimed = claimedWalletBalance
			ml.Unlock()

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
