package svc

import (
	"crypto/rand"
	"fmt"
	"github.com/JermineHu/themis/models"
	"github.com/JermineHu/themis/svc/gen/rtsp"
	"github.com/deepch/vdk/av"
	"log"
	"strconv"
)

//RTSPConfig global
var RTSPConfig *ConfigST

//ConfigST struct
type ConfigST struct {
	//	Server  ServerST            `json:"server"`
	Streams map[string]StreamST `json:"streams"`
}

////ServerST struct
//type ServerST struct {
//	HTTPPort string `json:"http_port"`
//}

//StreamST struct
type StreamST struct {
	URL    string `json:"url"`
	Status bool   `json:"status"`
	Codecs []av.CodecData
	Cl     map[string]viwer
}
type viwer struct {
	c chan av.Packet
}

func loadConfig() *ConfigST {
	tmp := ConfigST{}
	lpl := rtsp.ListPayload{}
	lpl.OffsetTail = 500
	list, _, err := models.GetRtspList(&lpl)
	if err != nil {
		log.Fatal("rtsp数据初始化失败！")
	}
	for i := range list {
		st := StreamST{}
		st.Cl = make(map[string]viwer)
		st.URL = *list[i].RtspURL
		id := strconv.FormatUint(list[i].ID, 10)
		tmp.Streams[id] = st
	}
	go func() { // 初始化流，建立连接
		for k, v := range tmp.Streams {
			go DailRTSP(k, v.URL)
		}
	}()
	return &tmp
}

func (element *ConfigST) cast(uuid string, pck av.Packet) {
	for _, v := range element.Streams[uuid].Cl {
		if len(v.c) < cap(v.c) {
			v.c <- pck
		}
	}
}

func (element *ConfigST) ext(suuid string) bool {
	_, ok := element.Streams[suuid]
	return ok
}

func (element *ConfigST) coAd(suuid string, codecs []av.CodecData) {
	t := element.Streams[suuid]
	t.Codecs = codecs
	element.Streams[suuid] = t
}

func (element *ConfigST) coGe(suuid string) []av.CodecData {
	return element.Streams[suuid].Codecs
}

func (element *ConfigST) clAd(suuid string) (string, chan av.Packet) {
	cuuid := pseudoUUID()
	ch := make(chan av.Packet, 100)
	element.Streams[suuid].Cl[cuuid] = viwer{c: ch}
	return cuuid, ch
}

func (element *ConfigST) list() (string, []string) {
	var res []string
	var fist string
	for k := range element.Streams {
		if fist == "" {
			fist = k
		}
		res = append(res, k)
	}
	return fist, res
}
func (element *ConfigST) clDe(suuid, cuuid string) {
	delete(element.Streams[suuid].Cl, cuuid)
}

func pseudoUUID() (uuid string) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	uuid = fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return
}
