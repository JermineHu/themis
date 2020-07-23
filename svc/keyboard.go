package svc

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/JermineHu/themis/models"
	keyboard "github.com/JermineHu/themis/svc/gen/keyboard"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"goa.design/goa/v3/security"
	"io"
	"log"
	"strings"
)

// keyboard service example implementation.
// The example methods log the requests and return zero values.
type keyboardsrvc struct {
	logger *log.Logger
}

// 设置参数
var (
	hostMap = map[uint64]map[string]keyboard.BrokerServerStream{}
)

// NewKeyboard returns the keyboard service implementation.
func NewKeyboard(logger *log.Logger) keyboard.Service {
	return &keyboardsrvc{logger}
}

// JWTAuth implements the authorization logic for service "keyboard" for the
// "jwt" security scheme.
func (s *keyboardsrvc) JWTAuth(ctx context.Context, token string, scheme *security.JWTScheme) (context.Context, error) {
	ctx, err := JWTCheck(ctx, token, scheme)
	if err != nil {
		return ctx, keyboard.MakeUnauthorized(err)
	}
	return ctx, err
}

// 键盘日志分页列表；
func (s *keyboardsrvc) List(ctx context.Context, p *keyboard.ListPayload) (res *keyboard.KeyboardList, err error) {
	res = &keyboard.KeyboardList{}
	if p == nil {
		lpl := keyboard.ListPayload{}
		lpl.OffsetTail = 20
		p = &lpl
	}
	list, count, err := models.GetKeyboardList(p)
	if err != nil {
		return
	}
	res.Count = &count
	ls := []*keyboard.KeyboardResult{}

	for i := range list {
		item := keyboard.KeyboardResult{}
		err = copier.Copy(&item, &list[i])
		if err != nil {
			return
		}
		ls = append(ls, &item)

	}
	res.PageData = ls
	return
}

// 键盘日志分页列表；
func (s *keyboardsrvc) ListByHostID(ctx context.Context, p *keyboard.ListByHostIDPayload) (res keyboard.KeyboardResultCollection, err error) {
	res = keyboard.KeyboardResultCollection{}

	pr := keyboard.ListPayload{}
	pr.Where = &keyboard.Keyboard{
		HostID: p.HostID,
	}
	pr.OrderBy = "id"
	pr.IsDesc = true
	list, _, err := models.GetKeyboardListByHostID(&pr)
	if err != nil {
		return
	}
	ls := []*keyboard.KeyboardResult{}

	for i := range list {
		item := keyboard.KeyboardResult{}
		err = copier.Copy(&item, &list[i])
		if err != nil {
			return
		}
		ls = append(ls, &item)
	}

	res = ls
	return
}

// 创建日志数据
func (s *keyboardsrvc) Log(ctx context.Context, p *keyboard.Keyboard) (res *keyboard.KeyboardResult, err error) {

	res = &keyboard.KeyboardResult{}
	cp := models.Keyboard{}
	v, err := json.Marshal(p.Keys)
	if err != nil {
		err = keyboard.MakeBadRequest(err)
		return
	}
	cp.Keys = v
	cp.HostID = p.HostID

	err = models.CreateKeyboard(&cp)
	if err != nil {
		err = keyboard.MakeBadRequest(err)
		return
	}
	err = copier.Copy(&res, &cp)
	if err != nil {
		return
	}
	return
}

// 根据主机ID删除，清空日志数据键盘
func (s *keyboardsrvc) Clear(ctx context.Context, p *keyboard.ClearPayload) (res bool, err error) {
	count, err := models.DeleteKeyboardByHostID(p.HostID)
	res = count > 0
	return
}

// 用于建立广播消息的服务
func (s *keyboardsrvc) Broker(ctx context.Context, p *keyboard.BrokerPayload, stream keyboard.BrokerServerStream) (err error) {

	sessionID := strings.ReplaceAll(uuid.New().String(), "-", "")
	hID := p.HostID
	if p.Token != nil {
		host_id, err := GetHostIDByJWT(*p.Token)
		if err != nil {
			return err
		}
		if host_id != nil {
			if *host_id != p.HostID {
				err = errors.New("token中的主机id与要发送事件的主机不一致！")
				return err
			}
			hID = *host_id
		}
	}

	keybCh := make(chan *keyboard.KeyboardEvent) // 定义接收数据的管道
	errCh := make(chan error)                    // 定义接收错误的管道
	go func() {
		for {
			if stream != nil {
				str, err := stream.Recv()
				if err != nil {
					if err != io.EOF {
						errCh <- err
					}
					close(keybCh)
					close(errCh)
					return
				}
				keybCh <- str
			}

		}
	}()

	if _, ok := hostMap[hID]; !ok {
		x := map[string]keyboard.BrokerServerStream{}
		hostMap[hID] = x
		x[sessionID] = stream
	} else {
		hostMap[hID][sessionID] = stream
	}
	//ticker  := time.NewTicker(time.Second * 5)

	// Listen for context cancellation and stream input simultaneously.
	for done := false; !done; {
		select {
		//case <-ticker.C: // 每隔1秒发送一个心跳
		//	hm := keyboard.KeyboardEvent{
		//		Type: "heartbeat",
		//	}
		//	if err = stream.Send(&hm); err != nil { // 发送心跳数据
		//		log.Println("write:", err)
		//		return err
		//	}
		case keyb := <-keybCh:
			if strings.EqualFold("keyboard", keyb.Type) && keyb.KeyboardInfo != nil {
				h := keyboard.Keyboard{
					HostID: &hID,
					Keys:   keyb.KeyboardInfo.Keys,
				}
				go func() {
					_, err := s.Log(ctx, &h) // 存储消息操作
					if err != nil {
						log.Fatal("键盘数据记录失败：", err)
					}
				}()

				if err = s.SenCast(hID, keyb); err != nil { // 将收到的消息再广播回去
					return err
				}

				//ks := []string{}
				//for k := range keyb.KeyboardInfo.Keys {
				//	ks = append(ks, *keyb.KeyboardInfo.Keys[k].KeyCode)
				//}
				//kstr := strings.Join(ks, "-")
				//
				//if _, ok := kbMap[kstr]; ok {
				//	go func() {
				//		_, err := s.Log(ctx, &h) // 存储消息操作
				//		if err != nil {
				//			log.Fatal("键盘数据记录失败：", err)
				//		}
				//	}()
				//	if err = s.SenCast(hID, keyb); err != nil { // 将收到的消息再广播回去
				//		return err
				//	}
				//}
			}
		case err := <-errCh: // 错误处理
			if err != nil {
				return err
			}
			done = true
			s.DelMapByKey(hID, sessionID)
		case <-ctx.Done():
			done = true
			s.DelMapByKey(hID, sessionID)
		}
	}
	s.DelMapByKey(hID, sessionID)
	return stream.Close()
}

// 发送广播，让对应会话中的人都收到消息
func (s *keyboardsrvc) SenCast(hostID uint64, event *keyboard.KeyboardEvent) (err error) {
	if v, ok := hostMap[hostID]; ok {
		for k := range v {
			if x, yes := v[k]; yes {
				if err = x.Send(event); err != nil { // 将收到的消息再广播回去
					return err
				}
			}
		}
	}
	return nil
}

// 发送广播，让对应会话中的人都收到消息
func (s *keyboardsrvc) DelMapByKey(hostID uint64, sessionID string) {
	if v, ok := hostMap[hostID]; ok {
		delete(hostMap[hostID], sessionID)
		if v, ok = hostMap[hostID]; ok {
			if len(v) == 0 {
				delete(hostMap, hostID)
			}
		}
	}
}
