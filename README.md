# Nyarrent site

# Inti steps
- Install `go`
- Install [mongoDB](https://www.mongodb.com/docs/manual/tutorial/install-mongodb-on-ubuntu/)

## Prerequireties
```sh
$ sudo apt install transmission-daemon zip ffmpeg
```

## Mandatory Environmentals
```bash
export NYARRENT_URI="mongodb://localhost:27017"
export NYARRENT_DATABASE_NAME="nyarrent"
export AS_API_KEY="xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
```

## Confugure ports in Powershell for WSL
```
netsh interface portproxy add v4tov4 listenport=80 listenaddress=0.0.0.0 connectport=3000 connectaddress=$($(wsl hostname -I).Trim());
```
