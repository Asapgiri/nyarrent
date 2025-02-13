function createRandomString(length) {
    const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";
    let result = "";
    for (let i = 0; i < length; i++) {
        result += chars.charAt(Math.floor(Math.random() * chars.length));
    }
    return result;
}

function getCookie(cname) {
    let name = cname + "=";
    let decodedCookie = decodeURIComponent(document.cookie);
    let ca = decodedCookie.split(';');
    for(let i = 0; i <ca.length; i++) {
        let c = ca[i];
        while (c.charAt(0) == ' ') {
            c = c.substring(1);
        }
        if (c.indexOf(name) == 0) {
            return c.substring(name.length, c.length);
        }
    }
    return "";
}

if ("" == getCookie("formathash")) {
    document.cookie = "formathash="+createRandomString(32);
}

function refreshTorrent(hash) {
    fetch('/api/gettorrent/' + hash)
        .then(resp => resp.json())
        .then(json => {
            if ('success' == json.result) {
                torrent = json.arguments.torrents[0];
                percent = (torrent.haveValid / torrent.sizeWhenDone) * 100;

                pbars = document.getElementsByClassName('progress-'+hash)
                for (let i = 0; i < pbars.length; i++) {
                    pbars[i].style.width = percent+'%';
                    pbars[i].innerText = percent+'%';
                }

                if (percent < 100) {
                    setTimeout(() => {refreshTorrent(hash)}, 1000)
                }
            }
        })
}
