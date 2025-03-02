package doh

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/miekg/dns"
	"github.com/semihalev/log"
	"github.com/stretchr/testify/assert"
)

func handleTest(w http.ResponseWriter, r *http.Request) {
	handle := func(req *dns.Msg) *dns.Msg {
		msg, _ := dns.Exchange(req, "8.8.8.8:53")

		return msg
	}

	var handleFn func(http.ResponseWriter, *http.Request)
	if r.Method == http.MethodGet && r.URL.Query().Get("dns") == "" {
		handleFn = HandleJSON(handle)
	} else {
		handleFn = HandleWireFormat(handle)
	}

	handleFn(w, r)
}

func Test_disQuery(t *testing.T) {
	// t.Parallel()

	w := httptest.NewRecorder()

	// 测试数据地址查询
	request, err := http.NewRequest("GET", "/dis-query/data/address?data_identifier=3e542bdf-b330-489e-aaa1-609b6d3ef704.data.fuxi.", nil)
	assert.NoError(t, err)

	request.RemoteAddr = "127.0.0.1:0"

	handleDISTest(w, request)

	assert.Equal(t, w.Code, http.StatusOK)

	data, err := ioutil.ReadAll(w.Body)
	assert.NoError(t, err)

	log.Info("data", string(data))

	var m ReturnMsg
	err = json.Unmarshal(data, &m)
	assert.NoError(t, err)

	log.Info("data_address", m.Data["data_address"])
	assert.NotNil(t, m.Data["data_address"])

	// 测试身份公钥查询
	w = httptest.NewRecorder()

	request, err = http.NewRequest("GET", "/dis-query/users/public-key?identity_identifier=fyftest.user.fuxi.", nil)
	assert.NoError(t, err)

	request.RemoteAddr = "127.0.0.1:0"

	handleDISTest(w, request)

	assert.Equal(t, w.Code, http.StatusOK)

	data, err = ioutil.ReadAll(w.Body)
	assert.NoError(t, err)

	log.Info("data", string(data))

	var m2 ReturnMsg
	err = json.Unmarshal(data, &m2)
	assert.NoError(t, err)

	log.Info("public_key", m2.Data["public_key"])
	assert.NotNil(t, m2.Data["public_key"])

	// 测试POD地址查询
	w = httptest.NewRecorder()

	request, err = http.NewRequest("GET", "/dis-query/users/pod?identity_identifier=fuyufantest.user.fuxi", nil)
	assert.NoError(t, err)

	request.RemoteAddr = "127.0.0.1:0"

	handleDISTest(w, request)

	assert.Equal(t, w.Code, http.StatusOK)

	data, err = ioutil.ReadAll(w.Body)
	assert.NoError(t, err)

	log.Info("data", string(data))

	var m3 ReturnMsg
	err = json.Unmarshal(data, &m3)
	assert.NoError(t, err)

	log.Info("PodAddress", m3.Data["pod_address"])
	assert.NotNil(t, m3.Data["pod_address"])

	// 测试所有者标识（RP）
	w = httptest.NewRecorder()

	request, err = http.NewRequest("GET", "/dis-query/data/owner?data_identifier=3e542bdf-b330-489e-aaa1-609b6d3ef704.data.fuxi.", nil)
	assert.NoError(t, err)

	request.RemoteAddr = "127.0.0.1:0"

	handleDISTest(w, request)

	assert.Equal(t, w.Code, http.StatusOK)

	data, err = ioutil.ReadAll(w.Body)
	assert.NoError(t, err)

	log.Info("data", string(data))

	var m4 ReturnMsg
	err = json.Unmarshal(data, &m4)
	assert.NoError(t, err)

	log.Info("OwnerID", m4.Data["owner_identity_identifier"])
	assert.NotNil(t, m4.Data["owner_identity_identifier"])

	// 测试授权信息（TXT）查询
	w = httptest.NewRecorder()

	request, err = http.NewRequest("GET", "/dis-query/authorization/info?data_identifier=3e542bdf-b330-489e-aaa1-609b6d3ef704.data.fuxi.&creator_identity_identifier=LWJWUCI3YCQFZKUUPSSV7SQG6ZAPLLOLCVQYKOO76FUBVPAC36LA====", nil)
	assert.NoError(t, err)

	request.RemoteAddr = "127.0.0.1:0"

	handleDISTest(w, request)

	assert.Equal(t, w.Code, http.StatusOK)

	data, err = ioutil.ReadAll(w.Body)
	assert.NoError(t, err)

	log.Info("data", string(data))

	var m5 ReturnMsg
	err = json.Unmarshal(data, &m5)
	assert.NoError(t, err)

	log.Info("authorization_info", m5.Data["authorization_info"])
	assert.NotNil(t, m5.Data["authorization_info"])

	// 测试数据摘要（TXT）查询
	w = httptest.NewRecorder()

	request, err = http.NewRequest("GET", "/dis-query/data/digest?data_identifier=3e542bdf-b330-489e-aaa1-609b6d3ef704.data.fuxi.", nil)
	assert.NoError(t, err)

	request.RemoteAddr = "127.0.0.1:0"

	handleDISTest(w, request)

	assert.Equal(t, w.Code, http.StatusOK)

	data, err = ioutil.ReadAll(w.Body)
	assert.NoError(t, err)

	log.Info("data", string(data))

	var m6 ReturnMsg
	err = json.Unmarshal(data, &m6)
	assert.NoError(t, err)

	log.Info("data_digest", m6.Data["data_digest"])
	assert.NotNil(t, m6.Data["data_digest"])

	// 测试hub地址（SRV）查询
	w = httptest.NewRecorder()

	request, err = http.NewRequest("GET", "/dis-query/hub/address?domain=user.fuxi.", nil)
	assert.NoError(t, err)

	request.RemoteAddr = "127.0.0.1:0"

	handleDISTest(w, request)

	assert.Equal(t, w.Code, http.StatusOK)

	data, err = ioutil.ReadAll(w.Body)
	assert.NoError(t, err)

	log.Info("data", string(data))

	var m7 ReturnMsg
	err = json.Unmarshal(data, &m7)
	assert.NoError(t, err)

	log.Info("hub_address", m7.Data["hub_address"])
	assert.NotNil(t, m7.Data["hub_address"])

}

func Test_disAuth(t *testing.T) {
	t.Parallel()

	// 授权验证
	w := httptest.NewRecorder()

	// userid := "LWJWUCI3YCQFZKUUPSSV7SQG6ZAPLLOLCVQYKOO76FUBVPAC36LA===="
	// dataid := "3e542bdf-b330-489e-aaa1-609b6d3ef704.data.fuxi."

	// cred, err := loadUserCredentials("../../test/userb.yaml")
	// assert.NoError(t, err)

	// accPrivKey, _, err := fetchKeyPair(cred)
	// assert.NoError(t, err)

	// cred, err := loadUserCredentials("../../test/fuyufantest.user.fuxi.yaml")
	// assert.NoError(t, err)

	// podPrivKey, _, err := fetchKeyPair(cred)
	// assert.NoError(t, err)

	// podSignature, err := sign(podPrivKey, hash([]byte(dataid+userid)))
	// assert.NoError(t, err)

	// accSignature, err := sign(accPrivKey, hash([]byte(dataid+userid)))
	// assert.NoError(t, err)

	// log.Info("accSignature", base64.StdEncoding.EncodeToString(accSignature))
	// log.Info("podSignature", base64.StdEncoding.EncodeToString(podSignature))

	request, err := http.NewRequest("GET", "/dis-query/authorization/authentication?identity_identifier=uTahrtp2OlS78hVxdM3GcPd1Z/AI2ELJLhe04JAZn8w=&data_identifier=fcb70ba0-6b9f-11ed-aef6-63114a05a8a0.data.fuxi.", nil)
	assert.NoError(t, err)

	request.RemoteAddr = "127.0.0.1:0"
	// request.Header.Set("Authorization", "Bearer "+base64.StdEncoding.EncodeToString(podSignature))

	handleDISTest(w, request)

	assert.Equal(t, w.Code, http.StatusOK)

	data, err := ioutil.ReadAll(w.Body)
	assert.NoError(t, err)

	log.Info("data", string(data))

	var m ReturnMsg
	err = json.Unmarshal(data, &m)
	assert.NoError(t, err)

	log.Info("Authorization Info", m.Data["authorization_info"])
	assert.NotNil(t, m.Data["authorization_info"])

	// 完整性验证
	w = httptest.NewRecorder()

	request, err = http.NewRequest("GET", "/dis-query/data/authentication?data_identifier=3e542bdf-b330-489e-aaa1-609b6d3ef704.data.fuxi.&data_digest=digest234-updated", nil)
	assert.NoError(t, err)

	request.RemoteAddr = "127.0.0.1:0"

	handleDISTest(w, request)

	assert.Equal(t, w.Code, http.StatusOK)

	data, err = ioutil.ReadAll(w.Body)
	assert.NoError(t, err)

	log.Info("data", string(data))

	var m2 ReturnMsg
	err = json.Unmarshal(data, &m2)
	assert.NoError(t, err)

	log.Info("Data Authentication", m2.Data["pass"])
	assert.NotNil(t, m2.Data["pass"])

}

func Test_dohJSON(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()

	request, err := http.NewRequest("GET", "/dns-query?name=www.google.com&type=a&do=true&cd=true", nil)
	assert.NoError(t, err)

	request.RemoteAddr = "127.0.0.1:0"

	handleTest(w, request)

	assert.Equal(t, w.Code, http.StatusOK)

	data, err := ioutil.ReadAll(w.Body)
	assert.NoError(t, err)

	var dm Msg
	err = json.Unmarshal(data, &dm)
	assert.NoError(t, err)

	assert.Equal(t, len(dm.Answer) > 0, true)

}

func Test_dohJSONerror(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()

	request, err := http.NewRequest("GET", "/dns-query?name=", nil)
	assert.NoError(t, err)

	request.RemoteAddr = "127.0.0.1:0"

	handleTest(w, request)

	assert.Equal(t, w.Code, http.StatusBadRequest)
}

func Test_dohJSONaccepthtml(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()

	request, err := http.NewRequest("GET", "/dns-query?name=www.google.com", nil)
	assert.NoError(t, err)

	request.RemoteAddr = "127.0.0.1:0"

	request.Header.Add("Accept", "text/html")
	handleTest(w, request)

	assert.Equal(t, w.Code, http.StatusOK)
	assert.Equal(t, w.Header().Get("Content-Type"), "application/x-javascript")
}

func Test_dohWireGET(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()

	req := new(dns.Msg)
	req.SetQuestion("www.google.com.", dns.TypeA)
	req.RecursionDesired = true

	data, err := req.Pack()
	assert.NoError(t, err)

	dq := base64.RawURLEncoding.EncodeToString(data)

	request, err := http.NewRequest("GET", fmt.Sprintf("/dns-query?dns=%s", dq), nil)
	assert.NoError(t, err)

	request.RemoteAddr = "127.0.0.1:0"

	handleTest(w, request)

	assert.Equal(t, w.Code, http.StatusOK)

	data, err = ioutil.ReadAll(w.Body)
	assert.NoError(t, err)

	msg := new(dns.Msg)
	err = msg.Unpack(data)
	assert.NoError(t, err)

	assert.Equal(t, msg.Rcode, dns.RcodeSuccess)

	assert.Equal(t, len(msg.Answer) > 0, true)
}

func Test_dohWireGETerror(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()

	request, err := http.NewRequest("GET", "/dns-query?dns=", nil)
	assert.NoError(t, err)

	request.RemoteAddr = "127.0.0.1:0"

	handleTest(w, request)

	assert.Equal(t, w.Code, http.StatusBadRequest)
}

func Test_dohWireGETbadquery(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()

	request, err := http.NewRequest("GET", "/dns-query?dns=Df4", nil)
	assert.NoError(t, err)

	request.RemoteAddr = "127.0.0.1:0"

	handleTest(w, request)

	assert.Equal(t, w.Code, http.StatusBadRequest)
}

func Test_dohWireHEAD(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()

	request, err := http.NewRequest("HEAD", "/dns-query?dns=", nil)
	assert.NoError(t, err)

	request.RemoteAddr = "127.0.0.1:0"

	handleTest(w, request)

	assert.Equal(t, w.Code, http.StatusMethodNotAllowed)
}

func Test_dohWirePOST(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()

	req := new(dns.Msg)
	req.SetQuestion("www.google.com.", dns.TypeA)
	req.RecursionDesired = true

	data, err := req.Pack()
	assert.NoError(t, err)

	request, err := http.NewRequest("POST", "/dns-query", bytes.NewReader(data))
	assert.NoError(t, err)

	request.RemoteAddr = "127.0.0.1:0"
	request.Header.Add("Content-Type", "application/dns-message")

	handleTest(w, request)

	assert.Equal(t, w.Code, http.StatusOK)

	data, err = ioutil.ReadAll(w.Body)
	assert.NoError(t, err)

	msg := new(dns.Msg)
	err = msg.Unpack(data)
	assert.NoError(t, err)

	assert.Equal(t, msg.Rcode, dns.RcodeSuccess)

	assert.Equal(t, len(msg.Answer) > 0, true)
}

func Test_dohWirePOSTError(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()

	request, err := http.NewRequest("POST", "/dns-query", bytes.NewReader([]byte{}))
	assert.NoError(t, err)

	request.RemoteAddr = "127.0.0.1:0"
	request.Header.Add("Content-Type", "text/html")

	handleTest(w, request)

	assert.Equal(t, w.Code, http.StatusUnsupportedMediaType)
}
