export type User = {
    username: string,
    firstName: string,
    lastName: string,
    email: string,
    mobileCountryCode: string,
    mobile: string,
    referrer: string,
    isCorporateUser: boolean,
    publicKey: string,
    secretKeys: string[],
}