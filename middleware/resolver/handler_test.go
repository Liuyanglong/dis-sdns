package resolver

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/miekg/dns"
	"github.com/semihalev/log"
	"github.com/semihalev/sdns/config"
	"github.com/semihalev/sdns/dnsutil"
	"github.com/semihalev/sdns/middleware"
	_ "github.com/semihalev/sdns/middleware/edns"
	"github.com/semihalev/sdns/mock"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	cfg := makeTestConfig()

	middleware.Setup(cfg)

	m.Run()

	os.Exit(0)
}

func makeTestConfig() *config.Config {
	log.Root().SetHandler(log.LvlFilterHandler(0, log.StdoutHandler))

	cfg := new(config.Config)
	cfg.RootServers = []string{"192.5.5.241:53"}
	cfg.RootKeys = []string{
		".			172800	IN	DNSKEY	257 3 8 AwEAAaz/tAm8yTn4Mfeh5eyI96WSVexTBAvkMgJzkKTOiW1vkIbzxeF3+/4RgWOq7HrxRixHlFlExOLAJr5emLvN7SWXgnLh4+B5xQlNVz8Og8kvArMtNROxVQuCaSnIDdD5LKyWbRd2n9WGe2R8PzgCmr3EgVLrjyBxWezF0jLHwVN8efS3rCj/EWgvIWgb9tarpVUDK/b58Da+sqqls3eNbuv7pr+eoZG+SrDK6nWeL3c6H5Apxz7LjVc1uTIdsIXxuOLYA4/ilBmSVIzuDWfdRUfhHdY6+cn8HFRm+2hM8AnXGXws9555KrUB5qihylGa8subX2Nn6UwNR1AkUTV74bU=",
	}
	cfg.Maxdepth = 30
	cfg.Expire = 600
	cfg.CacheSize = 1024
	cfg.Timeout.Duration = 2 * time.Second

	return cfg
}

func Test_handler(t *testing.T) {
	ctx := context.Background()

	handler := middleware.Get("resolver").(*DNSHandler)

	time.Sleep(2 * time.Second)

	fmt.Println(handler.resolver.rootkeys)
	fmt.Println(handler.cfg.BindDOH)

	assert.Equal(t, "resolver", handler.Name())

	m := new(dns.Msg)
	m.RecursionDesired = true
	m.SetQuestion("www.apple.com.", dns.TypeA)
	r := handler.handle(ctx, m)
	assert.Equal(t, len(r.Answer) > 0, true)

	m = new(dns.Msg)
	m.RecursionDesired = true
	// test again for caches
	m.SetQuestion("www.apple.com.", dns.TypeA)
	r = handler.handle(ctx, m)
	assert.Equal(t, len(r.Answer) > 0, true)

	m = new(dns.Msg)
	m.RecursionDesired = true
	m.SetEdns0(dnsutil.DefaultMsgSize, true)
	m.SetQuestion("dnssec-failed.org.", dns.TypeA)
	r = handler.handle(ctx, m)
	assert.Equal(t, len(r.Answer) == 0, true)

	m = new(dns.Msg)
	m.RecursionDesired = true
	m.SetQuestion("example.com.", dns.TypeA)
	r = handler.handle(ctx, m)
	assert.Equal(t, len(r.Answer) > 0, true)

	m = new(dns.Msg)
	m.RecursionDesired = true
	m.SetQuestion(".", dns.TypeANY)
	r = handler.handle(ctx, m)
	assert.Equal(t, r.Rcode, dns.RcodeNotImplemented)

	m = new(dns.Msg)
	m.RecursionDesired = true
	m.SetQuestion(".", dns.TypeNS)
	m.RecursionDesired = false
	r = handler.handle(ctx, m)
	assert.NotEqual(t, r.Rcode, dns.RcodeServerFailure)
}

func Test_HandlerHINFO(t *testing.T) {
	ctx := context.Background()
	cfg := makeTestConfig()
	handler := New(cfg)

	m := new(dns.Msg)
	m.SetQuestion(".", dns.TypeHINFO)

	debugns = true
	resp := handler.handle(ctx, m)

	assert.Equal(t, true, len(resp.Ns) > 0)
}

func Test_HandlerServe(t *testing.T) {
	cfg := makeTestConfig()
	h := New(cfg)

	ch := middleware.NewChain([]middleware.Handler{})
	mw := mock.NewWriter("tcp", "127.0.0.1:0")

	req := new(dns.Msg)
	req.SetQuestion(".", dns.TypeNS)

	ch.Reset(mw, req)

	h.ServeDNS(context.Background(), ch)
	assert.Equal(t, true, ch.Writer.Written())
}
