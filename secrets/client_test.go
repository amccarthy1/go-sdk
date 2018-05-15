package secrets

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestClientGetVersion(t *testing.T) {
	assert := assert.New(t)

	client, err := New()
	assert.Nil(err)

	mountMetaJSON := `{"request_id":"e114c628-6493-28ed-0975-418a75c7976f","lease_id":"","renewable":false,"lease_duration":0,"data":{"accessor":"kv_45f6a162","config":{"default_lease_ttl":0,"force_no_cache":false,"max_lease_ttl":0,"plugin_name":""},"description":"key/value secret storage","local":false,"options":{"version":"2"},"path":"secret/","seal_wrap":false,"type":"kv"},"wrap_info":null,"warnings":null,"auth":null}`

	m := NewMockHTTPClient().WithString("GET", URL("%s/v1/sys/internal/ui/mounts/secret/", client.Remote().String()), mountMetaJSON)
	client.WithHTTPClient(m)

	version, err := client.getVersion()
	assert.Nil(err)
	assert.Equal(Version2, version)
}

func TestClientGetMountMeta(t *testing.T) {
	assert := assert.New(t)

	client, err := New()
	assert.Nil(err)

	mountMetaJSON := `{"request_id":"e114c628-6493-28ed-0975-418a75c7976f","lease_id":"","renewable":false,"lease_duration":0,"data":{"accessor":"kv_45f6a162","config":{"default_lease_ttl":0,"force_no_cache":false,"max_lease_ttl":0,"plugin_name":""},"description":"key/value secret storage","local":false,"options":{"version":"2"},"path":"secret/","seal_wrap":false,"type":"kv"},"wrap_info":null,"warnings":null,"auth":null}`

	m := NewMockHTTPClient().WithString("GET", URL("%s/v1/sys/internal/ui/mounts/secret/", client.Remote().String()), mountMetaJSON)
	client.WithHTTPClient(m)

	mountMeta, err := client.getMountMeta()
	assert.Nil(err)
	assert.NotNil(mountMeta)
	assert.Equal(Version2, mountMeta.Data.Options["version"])
}

func TestClientJSONBody(t *testing.T) {
	assert := assert.New(t)

	client, err := New()
	assert.Nil(err)

	output, err := client.jsonBody(map[string]interface{}{
		"foo": "bar",
	})
	assert.Nil(err)
	defer output.Close()

	contents, err := ioutil.ReadAll(output)
	assert.Nil(err)
	assert.Equal("{\"foo\":\"bar\"}\n", string(contents))
}

func TestClientReadJSON(t *testing.T) {
	assert := assert.New(t)

	client, err := New()
	assert.Nil(err)

	jsonBody := bytes.NewBuffer([]byte(`{"foo":"bar"}`))

	output := map[string]interface{}{}
	assert.Nil(client.readJSON(jsonBody, &output))
	assert.Equal("bar", output["foo"])
}

func TestClientCopyRemote(t *testing.T) {
	assert := assert.New(t)

	client, err := New()
	assert.Nil(err)

	copy := client.copyRemote()
	copy.Host = "not_" + copy.Host

	anotherCopy := client.copyRemote()
	assert.NotEqual(anotherCopy.Host, copy.Host)
}

func TestClientDiscard(t *testing.T) {
	assert := assert.New(t)

	client, err := New()
	assert.Nil(err)

	assert.NotNil(client.discard(nil, fmt.Errorf("this is only a test")))

	assert.Nil(client.discard(client.jsonBody(map[string]interface{}{
		"foo": "bar",
	})))
}