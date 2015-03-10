# hiradio
[![Build Status](https://travis-ci.org/parkghost/hiradio.png)](https://travis-ci.org/parkghost/hiradio) 
[![Coverage Status](https://coveralls.io/repos/parkghost/hiradio/badge.svg)](https://coveralls.io/r/parkghost/hiradio)
[![GoDoc](https://godoc.org/github.com/parkghost/hiradio?status.svg)](https://godoc.org/github.com/parkghost/hiradio)

*hiradio helps to play radio via Hichannel*

## Usage
```
hiradio helps to play radio via Hichannel
Usage:

        hiradio [options] command [arg...]

The commands are:

    list                     List radio stations
    info                     Display radio information and program list
    play                     Play radio on player

Use "hiradio command -h" for more information about a command.
```

## Installation
```
go get github.com/parkghost/hiradio/...
```

## Commands

#### list
```text
$ hiradio list
編號  類型      排行  頻道                            現在播放節目
 222  音樂         1  HitFm聯播網 Taipei 北部         週日 HIT DJ
 156  音樂         2  KISS RADIO 大眾廣播電台         
 206  音樂         4  中廣音樂網i radio               i 溜達
 308  音樂         6  KISS RADIO 網路音樂台           音樂NON STOP
 205  音樂         7  中廣流行網 i like               蔣公廚房
 228  音樂         8  Classical Taiwan愛樂電台        歐洲音樂
  88  音樂         9  HitFm聯播網 中部                週日 HIT DJ
 212  音樂        10  BestRadio 台北好事              好事輕鬆點
 213  音樂        13  BestRadio 高雄港都              港都GoGo Sunday(小妍)
 248  音樂        14  AppleLine 蘋果線上              音樂
 ...
```

#### info [ChannelID]
```text
$ hiradio info 222
編號: 222
頻道: HitFm聯播網 Taipei 北部
類型: 音樂
地點: 北區(基北桃竹苗)
簡介: 熱情Play 只想聽音樂，為全方位的音樂電台。節目內容豐富多元，熱情、專業的DJ全天候放送流行音樂，讓您零時差的迅速抓住最新、最流行的中西方娛樂消息。
節目表:
   00:00 ~ 02:00  LOVE DJ
   02:00 ~ 09:00  只想聽音樂
   09:00 ~ 12:00  賴床 DJ
   12:00 ~ 13:00  HITO 日亞排行榜
   13:00 ~ 14:00  HIT週末!
   14:00 ~ 17:00  嗑音樂
>> 17:00 ~ 18:00  週日 HIT DJ
   18:00 ~ 20:00  HITO唱片行
   20:00 ~ 22:00  ROCK DJ
   22:00 ~ 24:00  HITO LATE NIGHT SHOW
```

#### play [options] [ChannelID]
```text
$ hiradio play -player /usr/bin/vlc 222
Press ctrl-c to exit
```

## License
This project is licensed under the MIT license