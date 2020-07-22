package vcache
//http客户端使用proto.Unmarshal进行获得的数据的解码

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"io/ioutil"
	"net/http"
	"net/url"
	"vcache/vcachepb"
)

type httpGetter struct {
	baseURL string
}

func (h *httpGetter) Get(in *vcachepb.Request, out *vcachepb.Response) error {
	u := fmt.Sprintf("%v%v/%v", h.baseURL, url.QueryEscape(in.GetGroup()),
		url.QueryEscape(in.GetKey()))
	//构造好完整的url后，就对url进行请求，并得到回复
	res, err := http.Get(u)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned: %v", res.Status)
	}

	//如果正确回复了200，那么就对其res.body进行处理
	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("reading response body:%v", err)
	}

	//获得到butes要对bytes进行protobuf解码
	if err := proto.Unmarshal(bytes, out); err != nil {
		return fmt.Errorf("decoding response body: %v", err)
	}
	return nil
}
