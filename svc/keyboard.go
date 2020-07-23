package svc

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/JermineHu/themis/models"
	"github.com/jinzhu/copier"
	"io"
	"log"

	keyboard "github.com/JermineHu/themis/svc/gen/keyboard"
	"goa.design/goa/v3/security"
)

// keyboard service example implementation.
// The example methods log the requests and return zero values.
type keyboardsrvc struct {
	logger *log.Logger
}

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

	host_id, err := GetHostIDByJWT(p.Token)
	if err != nil {
		return
	}
	if host_id != nil {
		if host_id != p.HostID {
			err = errors.New("token中的主机id与要发送事件的主机不一致！")
			return err
		}
	}
	keybCh := make(chan *keyboard.Keyboard) // 定义接收数据的管道
	errCh := make(chan error)               // 定义接收错误的管道
	go func() {
		for {
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
	}()

	// Listen for context cancellation and stream input simultaneously.
	for done := false; !done; {
		select {
		case keyb := <-keybCh:
			//s.storeMessage(keyb)  // 存储消息操作
			go func() {
				h := keyboard.Keyboard{
					HostID: host_id,
					Keys:   keyb.Keys,
				}
				_, err := s.Log(ctx, &h)
				if err != nil {
					log.Fatal("键盘数据记录失败：", err)
				}
			}()

			if err = stream.Send(keyb); err != nil { // 将收到的消息再广播回去
				return err
			}
		case err := <-errCh: // 错误处理
			if err != nil {
				return err
			}
			done = true
		case <-ctx.Done():
			done = true
		}
	}
	return stream.Close()
}
