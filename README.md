# vault-secret-receiver

integration with https://github.com/Boostport/kubernetes-vault project 

Why? 
* i need use secrets in go program in main k8s container but that repo provides only file 
i.e doesnt export those variables into ENV or something else

Logic:
* Unmarshal vault-token file into VaultUrl and VaultToken 
* Use that token to recive secrets from vault and use it in another go module 
* *Posible TODO*: remove vault-token file from disk (i believe we dont need it anymore)


## Quick start

```
vaultSettings := getVaultSettings()
client := initClient(vaultSettings)
client.readSecret("secretName","secretPath")
```