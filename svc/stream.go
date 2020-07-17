package svc

import (
	"github.com/deepch/vdk/format/rtsp"
	"log"
	"time"
)

// 根据RTSP地址建立一个rtsp连接
func DailRTSP(name, url string) {
	for {
		log.Println(name, "connect", url)
		rtsp.DebugRtsp = true
		session, err := rtsp.Dial(url)
		if err != nil {
			log.Println(name, err)
			time.Sleep(5 * time.Second)
			continue
		}
		session.RtpKeepAliveTimeout = 10 * time.Second
		if err != nil {
			log.Println(name, err)
			time.Sleep(5 * time.Second)
			continue
		}
		codec, err := session.Streams()
		if err != nil {
			log.Println(name, err)
			time.Sleep(5 * time.Second)
			continue
		}
		RTSPConfig.coAd(name, codec)
		for {
			pkt, err := session.ReadPacket()
			if err != nil {
				log.Println(name, err)
				break
			}
			RTSPConfig.cast(name, pkt)
		}
		err = session.Close()
		if err != nil {
			log.Println("session Close error", err)
		}
		log.Println(name, "reconnect wait 5s")
		time.Sleep(5 * time.Second)
	}
}
