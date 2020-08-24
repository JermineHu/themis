package svc

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/JermineHu/themis/common"
	"github.com/JermineHu/themis/models"
	keyboard "github.com/JermineHu/themis/svc/gen/keyboard"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"goa.design/goa/v3/security"
	"io"
	"log"
	"strings"
	"time"
)

// keyboard service example implementation.
// The example methods log the requests and return zero values.
type keyboardsrvc struct {
	logger *log.Logger
}

// 设置参数
var (
	hostMap  = map[uint64]map[string]keyboard.BrokerServerStream{}
	hostsMap = map[uint64]map[string]keyboard.BrokerForHostsServerStream{} // 列表
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
	_, err = JWTCheckForHost(ctx, *p.Token)
	if err != nil {
		return nil, keyboard.MakeUnauthorized(err)
	}
	res = &keyboard.KeyboardResult{}
	cp := models.Keyboard{}
	v, err := json.Marshal(p.Keys)
	if err != nil {
		err = keyboard.MakeBadRequest(err)
		return
	}
	cp.Keys = v
	cp.HostID = p.HostID
	ks := []string{}
	for k := range p.Keys {
		ks = append(ks, *p.Keys[k].KeyCode)
	}
	kstr := strings.Join(ks, "-")
	cp.EventCode = &kstr
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
	ticker := time.NewTicker(time.Second * 5)

	// Listen for context cancellation and stream input simultaneously.
	for done := false; !done; {
		select {
		case <-ticker.C: // 每隔1秒发送一个心跳
			hm := keyboard.KeyboardEvent{
				Type: "heartbeat",
			}
			if err = stream.Send(&hm); err != nil { // 发送心跳数据
				log.Println("write:", err)
				return err
			}
		case keyb := <-keybCh:
			if strings.EqualFold("keyboard", keyb.Type) && keyb.KeyboardInfo != nil {
				h := keyboard.Keyboard{
					HostID: &hID,
					Keys:   keyb.KeyboardInfo.Keys,
					Token:  p.Token,
				}
				ks := []string{}
				for k := range keyb.KeyboardInfo.Keys {
					ks = append(ks, *keyb.KeyboardInfo.Keys[k].KeyCode)
				}
				kstr := strings.Join(ks, "-")
				if _, ok := kbMap[kstr]; ok {
					go func() {
						result, err := s.Log(ctx, &h) // 存储消息操作
						if err != nil {
							log.Fatal("键盘数据记录失败：", err)
						}
						if strings.EqualFold(kbMap[kstr].EventType, common.EVENT_TYPE_DELETE_PRVE_DATA) {
							models.DeletePrevKeyboardByIDAndHostID(*result.ID, *result.HostID) // 删除当前主机的上一个
						}
					}()
					cl := kbMap[kstr].Color
					en := kbMap[kstr].Desc
					et := kbMap[kstr].EventType
					keyb.KeyboardInfo.Color = &cl
					keyb.KeyboardInfo.EName = &en
					keyb.KeyboardInfo.EType = &et
					if strings.EqualFold(kbMap[kstr].EventType, common.EVENT_TYPE_CLEAN) {
						models.DeleteKeyboardByHostID(p.HostID) // 删除该主机下的所有数据
					}
					if err = s.SenCast(hID, keyb); err != nil { // 将收到的消息再广播回去
						log.Fatal("发送失败：", err)
					}
					if len(hostsMap) > 0 {
						keybChList <- keyb // 向列表同步消息
					}
				}
			}
		case err := <-errCh: // 错误处理
			if err != nil {
				s.DelMapByKey(hID, sessionID)
				return err
			}
			done = true

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
func (s *keyboardsrvc) SenCastForListPage(hostID uint64, event *keyboard.KeyboardEvent) (err error) {
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

// 发送广播，让对应会话中的人都收到消息
func (s *keyboardsrvc) DelMapByKeyForListPage(hostIDs []uint64, sessionID string) {
	for i := range hostIDs {
		if v, ok := hostsMap[hostIDs[i]]; ok {
			delete(hostsMap[hostIDs[i]], sessionID)
			if v, ok = hostsMap[hostIDs[i]]; ok {
				if len(v) == 0 {
					delete(hostsMap, hostIDs[i])
				}
			}
		}
	}
}

// 根据主机ID获取统计数据
func (s *keyboardsrvc) Statistics(ctx context.Context, p *keyboard.StatisticsPayload) (res map[string]*keyboard.Keyboard, err error) {
	res = make(map[string]*keyboard.Keyboard)
	ks, err := models.StatisticsKeyboardEventByHostID(p.HostID)
	if err != nil {
		return nil, keyboard.MakeBadRequest(err)
	}
	for i := range ks {
		k := &keyboard.Keyboard{}
		k.HostID = ks[i].HostID
		if ks[i].EventCode != nil {
			if _, ok := kbMap[*ks[i].EventCode]; ok {
				et := kbMap[*ks[i].EventCode].EventType
				cl := kbMap[*ks[i].EventCode].Color
				en := kbMap[*ks[i].EventCode].Desc
				k.EType = &et
				k.Color = &cl
				k.EName = &en
			}
			k.Count = &ks[i].Count
			res[*ks[i].EventCode] = k
		}

	}
	return
}

var keybChList = make(chan *keyboard.KeyboardEvent, 100) // 定义接收数据的管道
// 用于建立广播消息的服务
func (s *keyboardsrvc) BrokerForHosts(ctx context.Context, p *keyboard.BrokerForHostsPayload, stream keyboard.BrokerForHostsServerStream) (err error) {

	sessionID := strings.ReplaceAll(uuid.New().String(), "-", "")

	ids, err := models.GetAllHostID()
	if err != nil {
		return err
	}

	//批量设置所监听的主机
	for i := range ids {
		if _, ok := hostsMap[ids[i]]; !ok {
			x := map[string]keyboard.BrokerForHostsServerStream{}
			hostsMap[ids[i]] = x
			x[sessionID] = stream
		} else {
			hostsMap[ids[i]][sessionID] = stream
		}
	}

	////批量设置所监听的主机
	//for i := range p.HostIds {
	//	if _, ok := hostsMap[p.HostIds[i]]; !ok {
	//		x := map[string]keyboard.BrokerForHostsServerStream{}
	//		hostsMap[p.HostIds[i]] = x
	//		x[sessionID] = stream
	//	} else {
	//		hostsMap[p.HostIds[i]][sessionID] = stream
	//	}
	//}

	ticker := time.NewTicker(time.Second * 5)

	// Listen for context cancellation and stream input simultaneously.
	for done := false; !done; {
		select {
		case <-ticker.C: // 每隔1秒发送一个心跳
			hm := keyboard.KeyboardEvent{
				Type: "heartbeat",
			}
			if err = stream.Send(&hm); err != nil { // 发送心跳数据
				log.Println("write:", err)
				return err
			}
		case keyb := <-keybChList:
			if strings.EqualFold("keyboard", keyb.Type) && keyb.KeyboardInfo != nil {
				ks := []string{}
				for k := range keyb.KeyboardInfo.Keys {
					ks = append(ks, *keyb.KeyboardInfo.Keys[k].KeyCode)
				}
				kstr := strings.Join(ks, "-")
				if _, ok := kbMap[kstr]; ok {
					cl := kbMap[kstr].Color
					en := kbMap[kstr].Desc
					keyb.KeyboardInfo.Color = &cl
					keyb.KeyboardInfo.EName = &en
					if err = s.SenCastForListPage(*keyb.KeyboardInfo.HostID, keyb); err != nil { // 将收到的消息再广播回去
						log.Fatal("发送失败：", err)
					}
				}
			}
		case <-ctx.Done():
			done = true
			s.DelMapByKeyForListPage(p.HostIds, sessionID)
		}
	}
	s.DelMapByKeyForListPage(p.HostIds, sessionID)
	return stream.Close()
}
