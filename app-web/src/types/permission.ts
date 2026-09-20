export type Permission = {
    createdAt?: Date,
    updatedAt?: Date,
    walletAddress: string,
    targetUsername: string,
    fullName: string,
    permission: string,
    permissionState?: PermissionState | null,    
}
  
export enum PermissionState {
    Revoked = 'Revoked',
    Modified = 'Modified',
    Added = 'Added'
}