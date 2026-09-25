export interface ILogin {
  loginType: string;
  username: string;
}

export interface ILoginRes {
  dynamicLink: string;
  qrCode: string;
  loginId: string;
}

export interface IVerifyLogin {
  targetUser: string;
  loginID: string;
}

export interface IWalletConnectResponse {
  accessToken: string;
  refreshToken: string;
  userInfo: {
    CountryCode: string;
    adminLevel: 0 | 1;
    contactPhone: string;
    email: string;
    first_name: string;
    // imageThumbnail: string;
    kycLevel: 0 | 1;
    last_name: string;
    maxAssetPerOffer: number;
    maxAssetPerOrder: number;
    mobile: string;
    offline: 0 | 1;
    suspended: 0 | 1;
    takerFee: string;
    takerReputation: {
      cancelations: string;
      completed: string;
      pending: string;
      rank: string;
      reportsAgainst: string;
      trades: string;
    };
    tradeStats: {
      cancelations: string;
      completed: string;
      pending: string;
      rank: string;
      reportsAgainst: string;
      trades: string;
    };
    username: string;
  };
}
