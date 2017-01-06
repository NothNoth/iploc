# iploc - IP localization

iploc is a unix shell command witten in go __detecting IP__ strings in text and __geolocating them__.

## "Hey, how do I use this?" - Bob

Easy:

    echo "116.31.116.44" | iploc
    [116.31.116.44:China Telecom Guangdong/CN/Xinxi]


## "Mmhh ok, but that will not change my world" - John

What about this use case ?

    cat ~/auth.log | grep Failed | ./iploc
	Jan  6 06:55:25 pogo sshd[9124]: Failed password for root from [116.31.116.44:China Telecom Guangdong/CN/Xinxi] port 47017 ssh2
	Jan  6 06:55:41 pogo sshd[9129]: Failed password for root from [116.31.116.44:China Telecom Guangdong/CN/Xinxi] port 42884 ssh2
	Jan  6 06:55:59 pogo sshd[9133]: Failed password for root from [116.31.116.44:China Telecom Guangdong/CN/Xinxi] port 48093 ssh2
	Jan  6 06:57:03 pogo sshd[9147]: Failed password for root from [116.31.116.44:China Telecom Guangdong/CN/Xinxi] port 41735 ssh2
	Jan  6 06:57:05 pogo sshd[9147]: Failed password for root from [116.31.116.44:China Telecom Guangdong/CN/Xinxi] port 41735 ssh2
	Jan  6 06:57:06 pogo sshd[9152]: Failed password for root from [89.132.5.72:UPC Magyarorszag Kft./HU/Gyomro] port 4938 ssh2

## "That's so cool, I need that thing!" -- Cassandra

Pretty easy:

    git clone git@github.com:NothNoth/iploc.git
    cd iploc/iploc/
    go build

## "How does this thing work?" -- Bruce

iploc basically performs regexps to extract ip addresses from the input text stream then uses and external online service for IP identification (http://ip-api.com/)

## "Can I use this and create a billion $ business?"

iploc is distributed under BSD licence.