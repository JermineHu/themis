package svc

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/JermineHu/themis/models"
	rtsp "github.com/JermineHu/themis/svc/gen/rtsp"
	"github.com/deepch/vdk/av"
	"github.com/deepch/vdk/codec/h264parser"
	"github.com/jinzhu/copier"
	"github.com/pion/webrtc/v2"
	"github.com/pion/webrtc/v2/pkg/media"
	"goa.design/goa/v3/security"
	"log"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// rtsp service example implementation.
// The example methods log the requests and return zero values.
type rtspsrvc struct {
	logger *log.Logger
}

func (s *rtspsrvc) Receive(ctx context.Context, r *rtsp.ReceivePayload) (res string, err error) {
	data := r.Data
	suuid := r.RtspID
	log.Println("Request", suuid)
	if RTSPConfig.ext(suuid) {
		/*

			Get Codecs INFO

		*/
		codecs := RTSPConfig.coGe(suuid)
		if codecs == nil {
			log.Println("Codec error")
			return res, rtsp.MakeBadRequest(fmt.Errorf("获取解码信息错误"))
		}
		sps := codecs[0].(h264parser.CodecData).SPS()
		pps := codecs[0].(h264parser.CodecData).PPS()
		/*

			Recive Remote SDP as Base64

		*/
		sd, err := base64.StdEncoding.DecodeString(data)
		if err != nil {
			log.Println("DecodeString error", err)
			return res, rtsp.MakeBadRequest(fmt.Errorf("DecodeString error:%v", err))

		}
		/*

			Create Media MediaEngine

		*/

		mediaEngine := webrtc.MediaEngine{}
		offer := webrtc.SessionDescription{
			Type: webrtc.SDPTypeOffer,
			SDP:  string(sd),
		}
		err = mediaEngine.PopulateFromSDP(offer)
		if err != nil {
			log.Println("PopulateFromSDP error", err)
			return res, rtsp.MakeBadRequest(fmt.Errorf("PopulateFromSDP error"))
		}

		var payloadType uint8
		for _, videoCodec := range mediaEngine.GetCodecsByKind(webrtc.RTPCodecTypeVideo) {
			if videoCodec.Name == "H264" && strings.Contains(videoCodec.SDPFmtpLine, "packetization-mode=1") {
				payloadType = videoCodec.PayloadType
				break
			}
		}
		if payloadType == 0 {
			log.Println("Remote peer does not support H264")
			return res, rtsp.MakeBadRequest(fmt.Errorf("Remote peer does not support H264"))
		}
		if payloadType != 126 {
			log.Println("Video might not work with codec", payloadType)
		}
		log.Println("Work payloadType", payloadType)
		api := webrtc.NewAPI(webrtc.WithMediaEngine(mediaEngine))

		peerConnection, err := api.NewPeerConnection(webrtc.Configuration{
			ICEServers: []webrtc.ICEServer{
				{
					URLs:       []string{"turn:numb.viagenie.ca"},
					Credential: "muazkh",
					Username:   "webrtc@live.com",
				},
				{
					URLs: []string{"stun:stun.l.google.com:19302"},
				},
			},
		})
		if err != nil {
			log.Println("NewPeerConnection error", err)
			return res, rtsp.MakeBadRequest(fmt.Errorf("NewPeerConnection error：%v", err))
		}
		/*

			ADD KeepAlive Timer

		*/
		timer1 := time.NewTimer(time.Second * 2)
		peerConnection.OnDataChannel(func(d *webrtc.DataChannel) {
			// Register text message handling
			d.OnMessage(func(msg webrtc.DataChannelMessage) {
				//fmt.Printf("Message from DataChannel '%s': '%s'\n", d.Label(), string(msg.Data))
				timer1.Reset(2 * time.Second)
			})
		})
		/*

			ADD Video Track

		*/
		videoTrack, err := peerConnection.NewTrack(payloadType, rand.Uint32(), "video", suuid+"_pion")
		if err != nil {
			log.Fatalln("NewTrack", err)
		}
		_, err = peerConnection.AddTransceiverFromTrack(videoTrack,
			webrtc.RtpTransceiverInit{
				Direction: webrtc.RTPTransceiverDirectionSendonly,
			},
		)
		if err != nil {
			log.Println("AddTransceiverFromTrack error", err)
			return res, rtsp.MakeBadRequest(fmt.Errorf("AddTransceiverFromTrack error：%v", err))
		}
		_, err = peerConnection.AddTrack(videoTrack)
		if err != nil {
			log.Println("AddTrack error", err)
			return res, rtsp.MakeBadRequest(fmt.Errorf("AddTrack error：%v", err))
		}
		/*

			ADD Audio Track

		*/
		var audioTrack *webrtc.Track
		if len(codecs) > 1 && (codecs[1].Type() == av.PCM_ALAW || codecs[1].Type() == av.PCM_MULAW) {
			switch codecs[1].Type() {
			case av.PCM_ALAW:
				audioTrack, err = peerConnection.NewTrack(webrtc.DefaultPayloadTypePCMA, rand.Uint32(), "audio", suuid+"audio")
			case av.PCM_MULAW:
				audioTrack, err = peerConnection.NewTrack(webrtc.DefaultPayloadTypePCMU, rand.Uint32(), "audio", suuid+"audio")
			}
			if err != nil {
				log.Println(err)
				return res, rtsp.MakeBadRequest(fmt.Errorf(" error：%v", err))
			}
			_, err = peerConnection.AddTransceiverFromTrack(audioTrack,
				webrtc.RtpTransceiverInit{
					Direction: webrtc.RTPTransceiverDirectionSendonly,
				},
			)
			if err != nil {
				log.Println("AddTransceiverFromTrack error", err)
				return res, rtsp.MakeBadRequest(fmt.Errorf("AddTransceiverFromTrack error：%v", err))
			}
			_, err = peerConnection.AddTrack(audioTrack)
			if err != nil {
				log.Println(err)
				return res, rtsp.MakeBadRequest(fmt.Errorf("error：%v", err))
			}
		}
		if err := peerConnection.SetRemoteDescription(offer); err != nil {
			log.Println("SetRemoteDescription error", err, offer.SDP)
			return res, rtsp.MakeBadRequest(fmt.Errorf("SetRemoteDescription error：%v", err))
		}
		answer, err := peerConnection.CreateAnswer(nil)
		if err != nil {
			log.Println("CreateAnswer error", err)
			return res, rtsp.MakeBadRequest(fmt.Errorf("CreateAnswer error：%v", err))
		}

		if err = peerConnection.SetLocalDescription(answer); err != nil {
			log.Println("SetLocalDescription error", err)
			return res, rtsp.MakeBadRequest(fmt.Errorf("SetLocalDescription error：%v", err))
		}
		res = base64.StdEncoding.EncodeToString([]byte(answer.SDP))
		if err != nil {
			log.Println("Writer SDP error", err)
			return res, rtsp.MakeBadRequest(fmt.Errorf("Writer SDP error：%v", err))
		}
		control := make(chan bool, 10)
		var first int = 1
		peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
			if first > 0 {
				log.Printf("Connection State has changed %s \n", connectionState.String())
				//	if connectionState != webrtc.ICEConnectionStateConnected {
				//		log.Println("Client Close Exit")
				//		err := peerConnection.Close()
				//		if err != nil {
				//			log.Println("peerConnection Close error", err)
				//		}
				//		control <- true
				//		return
				//	}
				if connectionState == webrtc.ICEConnectionStateConnected {
					go func() {
						cuuid, ch := RTSPConfig.clAd(suuid)
						log.Println("start stream", suuid, "client", cuuid)
						defer func() {
							log.Println("stop stream", suuid, "client", cuuid)
							defer RTSPConfig.clDe(suuid, cuuid)
						}()
						var Vpre time.Duration
						var start bool
						timer1.Reset(5 * time.Second)
						for {
							select {
							case <-timer1.C:
								log.Println("Client Close Keep-Alive Timer")
								peerConnection.Close()
							case <-control:
								return
							case pck := <-ch:
								//timer1.Reset(2 * time.Second)
								if pck.IsKeyFrame {
									start = true
								}
								if !start {
									continue
								}
								if pck.IsKeyFrame {
									pck.Data = append([]byte("\000\000\001"+string(sps)+"\000\000\001"+string(pps)+"\000\000\001"), pck.Data[4:]...)

								} else {
									pck.Data = pck.Data[4:]
								}
								var Vts time.Duration
								if pck.Idx == 0 && videoTrack != nil {
									if Vpre != 0 {
										Vts = pck.Time - Vpre
									}
									samples := uint32(90000 / 1000 * Vts.Milliseconds())
									err := videoTrack.WriteSample(media.Sample{Data: pck.Data, Samples: samples})
									if err != nil {
										return
									}
									Vpre = pck.Time
								} else if pck.Idx == 1 && audioTrack != nil {
									err := audioTrack.WriteSample(media.Sample{Data: pck.Data, Samples: uint32(len(pck.Data))})
									if err != nil {
										return
									}
								}
							}
						}

					}()
					first = 0
				}
			}
		})
	}
	return
}

func (s *rtspsrvc) Codec(ctx context.Context, payload *rtsp.CodecPayload) (res interface{}, err error) {
	rtsp_id := payload.RtspID
	if RTSPConfig.ext(rtsp_id) {
		codecs := RTSPConfig.coGe(rtsp_id)
		if codecs == nil {
			return
		}
		res = codecs
	}
	return
}

// NewRtsp returns the rtsp service implementation.
func NewRtsp(logger *log.Logger) rtsp.Service {
	return &rtspsrvc{logger}
}

// JWTAuth implements the authorization logic for service "rtsp" for the "jwt"
// security scheme.
func (s *rtspsrvc) JWTAuth(ctx context.Context, token string, scheme *security.JWTScheme) (context.Context, error) {
	ctx, err := JWTCheck(ctx, token, scheme)
	if err != nil {
		return ctx, rtsp.MakeUnauthorized(err)
	}
	return ctx, err
}

// 流的数据列表；
func (s *rtspsrvc) List(ctx context.Context, p *rtsp.ListPayload) (res *rtsp.RtspList, err error) {
	res = &rtsp.RtspList{}
	if p == nil {
		lpl := rtsp.ListPayload{}
		lpl.OffsetTail = 20
		p = &lpl
	}
	list, count, err := models.GetRtspList(p)
	if err != nil {
		return
	}
	res.Count = &count
	ls := []*rtsp.RtspResult{}

	for i := range list {
		item := rtsp.RtspResult{}
		err = copier.Copy(&item, &list[i])
		if err != nil {
			return
		}
		ct := list[i].CreatedAt.Format(time.RFC3339)
		ut := list[i].UpdatedAt.Format(time.RFC3339)
		item.CreatedAt = &ct
		item.UpdatedAt = &ut
		ls = append(ls, &item)

	}
	res.PageData = ls
	return
}

// 创建RTSP数据
func (s *rtspsrvc) Create(ctx context.Context, p *rtsp.Rtsp) (res *rtsp.RtspResult, view string, err error) {
	res = &rtsp.RtspResult{}
	view = "default"
	cp := models.Rtsp{}
	err = copier.Copy(&cp, p)
	if err != nil {
		return
	}

	if ok, _ := regexp.MatchString("\\d+", p.HostID); !ok {
		return nil, "", rtsp.MakeBadRequest(errors.New("主机为必填！如果没有主机请先通过agent进行注册！"))
	}

	err = models.CreateRtsp(&cp)
	if err != nil {
		err = rtsp.MakeBadRequest(err)
		return
	}
	st := StreamST{}
	st.Cl = make(map[string]viwer)
	st.URL = *p.RtspURL
	id := strconv.FormatUint(cp.ID, 10)
	RTSPConfig.Streams[id] = st
	go DailRTSP(id, st.URL) // 建立连接
	err = copier.Copy(&res, &cp)
	if err != nil {
		return
	}
	return
}

// 根据id修改数据
func (s *rtspsrvc) Update(ctx context.Context, p *rtsp.Rtsp) (res *rtsp.RtspResult, view string, err error) {
	res = &rtsp.RtspResult{}
	view = "default"
	cp := models.Rtsp{}
	err = copier.Copy(&cp, p)
	if err != nil {
		return
	}
	err = models.UpdateRtspByID(*p.ID, &cp)
	if err != nil {
		err = rtsp.MakeBadRequest(err)
		return
	}
	st := StreamST{}
	st.Cl = make(map[string]viwer)
	st.URL = *p.RtspURL
	id := strconv.FormatUint(cp.ID, 10)
	RTSPConfig.Streams[id] = st
	go DailRTSP(id, st.URL) // 建立连接
	err = copier.Copy(&res, &cp)
	if err != nil {
		return
	}
	return
}

// 根据id删除
func (s *rtspsrvc) Delete(ctx context.Context, p *rtsp.DeletePayload) (res bool, err error) {
	count, err := models.DeleteRtspByID(p.ID)
	res = count > 0
	if res {
		delete(RTSPConfig.Streams, strconv.FormatUint(p.ID, 10)) // 删除map中的数据
	}
	return
}

// 根据id获取信息
func (s *rtspsrvc) Show(ctx context.Context, p *rtsp.ShowPayload) (res *rtsp.RtspResult, view string, err error) {
	res = &rtsp.RtspResult{}
	view = "default"
	cp, err := models.GetRtspById(p.ID)
	if err != nil {
		err = rtsp.MakeBadRequest(err)
		return
	}
	err = copier.Copy(&res, &cp)
	if err != nil {
		return
	}
	ct := cp.CreatedAt.Format(time.RFC3339)
	ut := cp.UpdatedAt.Format(time.RFC3339)
	res.CreatedAt = &ct
	res.UpdatedAt = &ut
	return
}
