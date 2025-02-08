# Nyarrent site

# Inti steps
- Install `go`
- Install [mongoDB](https://www.mongodb.com/docs/manual/tutorial/install-mongodb-on-ubuntu/)

## Prerequireties
```sh
$ sudo apt install transmission-daemon zip ffmpeg
```

## Configs
The `.config.json` file will be created after the first run with the following defaults:
```json
{
  "Http": {
    "Url": "",
    "Port": "3000"
  },
  "Dbase": {
    "Url": "mongodb://localhost:27017",
    "Name": "nyarrent"
  },
  "AnimeAPI": {
    "ApiUrl": "https://animeschedule.net/api/v3",
    "ApiKey": ""
  }
}
```

## Confugure ports in Powershell for WSL
```
netsh interface portproxy add v4tov4 listenport=80 listenaddress=0.0.0.0 connectport=3000 connectaddress=$($(wsl hostname -I).Trim());
```

