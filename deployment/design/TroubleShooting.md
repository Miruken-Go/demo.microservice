# Trouble Shooting

## Orphaned System Permissions
    
    [{"code":"RoleAssignmentUpdateNotPermitted","message":"Tenant ID, application ID, principal ID, and scope are not allowed to be updated."}]

This is a bug azure.  Deleting resources does not delete role assignments for system generated ids.

To fix look at resources with role assignments created in bicep.  Such as KeyVaults that need "Key Vault Secrets User" permissions
* Navigate to the keyvaults > IAM > Role Assignments
* Delete any assinments where the user name is "Identity not found."  These are orphaned permissions