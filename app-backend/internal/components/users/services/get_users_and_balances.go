package users

import (
	"log"
	"sync"
	usersDB "trovo-wallet-api/internal/components/users/db"
	userModels "trovo-wallet-api/internal/components/users/models"
	"trovo-wallet-api/internal/sharedconfig"

	tErrors "trovo-wallet-api/internal/errors"
)

// GetUserInfo gets the user Information
func GetUserInfo(identifier string, signerAddress string, gc *sharedconfig.GlobalConfig) (userInfo userModels.UserInfo, err error) {
	var owner bool
	//get user from DB
	primarySigner, err := usersDB.GetUserFromPrimarySigner(signerAddress, gc.DB, gc)
	if err != nil {
		return userModels.UserInfo{}, &tErrors.CustomError{Param: "primarySigner",
			Err:        "error invalid primary signer",
			ErrMessage: "Your request signer is not valid",
		}
	}

	//get user from DB
	user, err := usersDB.GetUser(identifier, gc.DB, gc)
	if err != nil {
		return userModels.UserInfo{}, err
	}

	if user.Suspended == 1 {
		return userInfo, &tErrors.ErrorUsernameIsSuspended{
			Username: user.Username,
		}
	}
	if primarySigner.Username == user.Username {
		owner = true
	}

	//set userInfo
	// log.Printf("[GetUserInfo] retrieved User record:[%+v]\n", user)
	userData := user.ToJSON(gc)
	// log.Printf("[GetUserInfo] userData:[%+v]\n", userData)

	userInfo.UserData = userData

	if !owner {
		userInfo.UserData.UserWallets = nil
		userInfo.AssetBalances = nil
		userInfo.WalletsSharedWithUser = nil
		userInfo.DefaultAssets = nil
		// userInfo.MySharedAccess = nil
	}

	//Get user wallet balances
	if owner {
		userInfo.AssetBalances = make(map[string]userModels.AssetBalances)

		// assetBalances, err := GetUserWalletAssetBalances(&user, gc)
		assetBalances, err := user.GetUserWalletAssetBalances(gc)
		if err == nil {
			userInfo.AssetBalances = assetBalances
		}
		// log.Println("[GetUserWalletAssetBalances] finished user wallets json")

		userInfo.NFTs = make(map[string][]userModels.NFT)

		userNFTs, err := GetUserNFTs(&user, gc)
		if err == nil {
			userInfo.NFTs = userNFTs
		}

		userInfo.WalletsSharedWithUser = make([]userModels.WalletsSharedWithUser, 0)
		//Get ThirdParty Wallet Access

		userInfo.WalletsSharedWithUser = user.FetchWalletsPermissionsSharedWithUser(gc)

		userInfo.DefaultAssets = user.GetDefaultAssets(gc)
	}

	return
}

func GetUserNFTs(user *userModels.User, gc *sharedconfig.GlobalConfig) (userNFTs map[string][]userModels.NFT, err error) {

	userNFTs = make(map[string][]userModels.NFT)
	for _, wallet := range user.UserWallets {

		var wg sync.WaitGroup
		var m sync.Mutex
		//use go routine to fetch

		wg.Add(1)
		go func(vg2 userModels.UserWallet, w *sync.WaitGroup, ml *sync.Mutex) {
			defer w.Done()
			//get only the NFTs in the main wallet
			nfts, errR1 := vg2.GetNFTs(false, gc)

			if errR1 != nil {
				//log server error
				log.Printf("[GetUserNFTs] error getting NFT asset for user:[%s] wallet:[%s] error:[%+v]\n", user.Username, vg2.ID, errR1)

			}
			//Claimed Assets
			ml.Lock()
			userNFTs[vg2.ID] = nfts
			ml.Unlock()

		}(wallet, &wg, &m)
		wg.Wait()

		// log.Println("exited inner wait")

	}
	// log.Println("done...")

	return
}

func GetUserWalletAssetBalances(user *userModels.User, gc *sharedconfig.GlobalConfig) (userWalletBalances map[string]userModels.AssetBalances, err error) {

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
			unclaimedBalance, errR1 := vg1.GetSortedUserBalance(true, gc)

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
			claimedWalletBalance, errR1 := vg2.GetSortedUserBalance(false, gc)

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

func GetWalletAssetBalances(wallet *userModels.UserWallet, gc *sharedconfig.GlobalConfig) (assetBalances userModels.AssetBalances, err error) {

	var wg sync.WaitGroup
	var m sync.Mutex
	//use go routine to fetch

	assetBalances.Unclaimed = make([]userModels.Balance, 0)
	assetBalances.Claimed = make([]userModels.Balance, 0)

	wg.Add(1)
	go func(vg1 *userModels.UserWallet, w *sync.WaitGroup, ml *sync.Mutex) {
		defer w.Done()
		unclaimedBalance, errR1 := vg1.GetSortedUserBalance(true, gc)

		if errR1 == nil {
			//Unclaimed Assets
			ml.Lock()
			assetBalances.Unclaimed = unclaimedBalance
			ml.Unlock()

		}
	}(wallet, &wg, &m)
	wg.Add(1)
	go func(vg2 *userModels.UserWallet, w *sync.WaitGroup, ml *sync.Mutex) {
		defer w.Done()
		claimedWalletBalance, errR1 := vg2.GetSortedUserBalance(false, gc)

		if errR1 != nil {
			//log server error
			log.Printf("[GetWalletAssetBalances] error getting claimed wallet balance for wallet:[%s] error:[%+v]\n", vg2.ID, errR1)

		}
		//Claimed Assets
		ml.Lock()
		assetBalances.Claimed = claimedWalletBalance
		ml.Unlock()

	}(wallet, &wg, &m)
	wg.Wait()

	return
}
