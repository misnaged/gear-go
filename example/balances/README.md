## transfer_keep_alive example

by default sender is `Alice` and receiver is `Bob`

---
### Receiver
To change the receiver you need to change args, as shown in code below:
```go
	args := []any{map[string]interface{}{"Id": "PUT_RECEIVER_ID_HERE"}, 12345} 
```
----
### Sender

Sender could be changed by setting keyring values in application config.
These values could be changed manually, by editting config/init.go file.

only two categories are available to be set: **Sr25519** and **Ed25519** 
example is below:
```go
    // 
    // 
	viper.SetDefault("keyring.category", "Ed25519") 

	viper.SetDefault("keyring.seed", "0xe5be9a5092b81bca64be81d212e7f2f9eba183bb7a90954f7b76361f6edb5c0a") // Alice seed 
```

Another way to change these values is setting them by env var as shown below:

```shell
export KEYRING_CATEGORY="Ed25519"
export KEYRING_SEED="0xe5be9a5092b81bca64be81d212e7f2f9eba183bb7a90954f7b76361f6edb5c0a"
```

## To check account status see *example/storage/get_account_value*


