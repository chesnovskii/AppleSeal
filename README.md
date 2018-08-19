# AppleSeal

## Envs   

**`SEAL_TG_TOKEN`** - Telegram bot token   

_example:_   
```
export SEAL_TG_TOKEN="some_very_secret_token"

echo $SEAL_TG_TOKEN
some_very_secret_token
```

**`APPLE_SEAL_DEBUG`**  - logging debug switch
```
export APPLE_SEAL_DEBUG="true"

echo $APPLE_SEAL_DEBUG
true
```

## Building   

```
go get git@github.com:newrushbolt/AppleSeal.git
cd $GOPATH/src/github.com/newrushbolt/AppleSeal
make all
```

