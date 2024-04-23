export type Permission = {
    createdAt?: Date,
    updatedAt?: Date,
    walletPublicKey: string,
    targetUsername: string,
    fullName: string,
    permission: string,
    permissionState: PermissionState,    
}
  
enum PermissionState {Revoked, Modified, Added}