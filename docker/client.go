package docker

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	//"github.com/docker/docker/pkg/stdcopy"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	httpclient  *http.Client
	endpoint    string
	endpointUrl *url.URL
}

func NewClient(endpoint string) *Client {
	endpointurl, err := ParseURL(endpoint)
	if err != nil {
		log.Fatal("Error new a client")
		return nil
	}
	return &Client{
		httpclient:  &http.Client{}, //http.DefaultClient,
		endpoint:    endpoint,
		endpointUrl: endpointurl,
	}
} //Does not support tls for now

func ParseURL(endpoint string) (*url.URL, error) {
	url, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	if url.Scheme == "unix" || url.Scheme == "http" {
		return url, nil
	} else if url.Scheme == "tcp" {
		url.Scheme = "http"
		return url, nil
	} else {
		return nil, errors.New("Invalid endpoint")
	}

}

type DoOption struct {
	data interface{}
}

func fakeDial(protocol, addr string) (net.Conn, error) {
	return net.Dial(protocol, addr)
}

func (c *Client) do(method string, url string, opt DoOption) ([]byte, int, error) {
	var param io.Reader

	if opt.data != nil {
		buf, err := json.Marshal(opt.data)
		if err != nil {
			return nil, -1, err
		}
		param = bytes.NewBuffer(buf)
	}

	req, err := http.NewRequest(method, url, param)
	if opt.data != nil {
		req.Header.Set("Content-Type", "application/json") //Set sets the header entries associated with key to the single element value. It replaces any existing values associated with key.
	} else if method == "POST" {
		req.Header.Set("Content-Type", "plain/text")
	}

	var resp *http.Response
	if c.endpointUrl.Scheme == "unix" {
		/*//dial, err := net.Dial("unix", c.endpointUrl.Path)
		tr := &http.Transport{
			Dial: func(proto,add )}
		//tr.Dial("unix", c.endpointUrl.Path)
		client := &http.Client{Transport: tr}

		resp, err = client.Do(req) */ /*Things are totally wrong here.*/

		// 设置 TimeOut
		DefaultClient := http.Client{
			Transport: &http.Transport{
				Dial: func(netw, addr string) (net.Conn, error) {
					deadline := time.Now().Add(30 * time.Second)
					c, err := net.DialTimeout(netw, addr, time.Second*30)
					if err != nil {
						return nil, err
					}
					c.SetDeadline(deadline)
					return c, nil
				},
			},
		}
		// -------------------------------------------
		// 执行
		resp, err = DefaultClient.Do(req)

	} else {
		resp, err = c.httpclient.Do(req)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, -1, err
	}
	fmt.Println(string(body))
	return body, resp.StatusCode, nil

}

type StreamOption struct {
	header         map[string]string
	in             io.Reader
	stdout         io.Writer
	stderr         io.Writer
	rawJsonStream  bool
	setRawTerminal bool
}

type jsonMessage struct {
	Status   string `json:"status,omitempty"`
	Progress string `json:"progress,omitempty"`
	Error    string `json:"error,omitempty"`
	Stream   string `json:"stream,omitempty"`
}

func (c *Client) stream(method string, url string, opt StreamOption) error {

	req, err := http.NewRequest(method, url, opt.in)
	if err != nil {
		return err
	}
	log.Println("step 1: New a http req")
	for key, val := range opt.header {
		req.Header.Set(key, val)
	}

	if opt.stdout == nil {
		opt.stdout = ioutil.Discard
	}
	if opt.stderr == nil {
		opt.stderr = ioutil.Discard
	}

	var resp *http.Response
	resp, err = c.httpclient.Do(req)
	if err != nil {
		return err
	}
	log.Println("Step 2: Do req")
	defer resp.Body.Close()
	/*_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	log.Println("Step 3: Read the resp body")
	fmt.Println(string(body))
	log.Print("Step 4: Copy to stdout")*/
	if resp.Header.Get("Content-Type") == "application/json" {
		log.Println("Transport json directly")
		if opt.rawJsonStream {
			_, err = io.Copy(opt.stdout, resp.Body)
			if err != nil {
				return err
			}
		}

		dec := json.NewDecoder(resp.Body)
		for {
			var m jsonMessage
			err = dec.Decode(&m)
			if err == io.EOF {
				break
			} else if err != nil {
				return err
			}

			if m.Stream != "" {
				fmt.Fprint(opt.stdout, m.Stream)
			} else if m.Progress != "" {
				fmt.Fprintf(opt.stdout, "%s/%s\r", m.Status, m.Progress)
			} else if m.Error != "" {
				return errors.New(m.Error)
			}
			if m.Status != "" {
				fmt.Fprintln(opt.stdout, m.Status)
			}
		}
	} else {
		if opt.setRawTerminal {
			log.Println("io.Copy")
			_, err = io.Copy(opt.stdout, resp.Body)
		} else {
			log.Println("stdCopy")
			//_, err = stdcopy.StdCopy(opt.stdout, opt.stderr, resp.Body) // awkard
			bytes, err := io.Copy(opt.stdout, resp.Body)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("bytes", bytes)
		}
		return err

	}
	return nil

}

/*type Kind uint

const (
	Invalid Kind = iota
	Bool
	Int
	Int8
	Int16
	Int32
	Int64
	Uint
	Uint8
	Uint16
	Uint32
	Uint64
	Uintptr
	Float32
	Float64
	Complex64
	Complex128
	Array
	Chan
	Func
	Interface
	Map
	Ptr
	Slice
	String
	Struct
	UnsafePointer
)
		//http://golang.org/src/reflect/type.go
		// Methods applicable only to some types, depending on Kind.
   		// The methods allowed for each kind are:
   		//
   		//	Int*, Uint*, Float*, Complex*: Bits
		//	Array: Elem, Len
		//	Chan: ChanDir, Elem
   		//	Func: In, NumIn, Out, NumOut, IsVariadic.
   		//	Map: Key, Elem
   		//	Ptr: Elem
   		//	Slice: Elem
   		//	Struct: Field, FieldByIndex, FieldByName, FieldByNameFunc, NumField
*/

func queryString(opts interface{}) string {
	if opts == nil {
		return ""
	}
	value := reflect.ValueOf(opts)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	if value.Kind() != reflect.Struct {
		return ""
	}
	items := url.Values(map[string][]string{})
	for i := 0; i < value.NumField(); i++ {
		field := value.Type().Field(i)
		if field.PkgPath != "" {
			continue
		}
		key := field.Tag.Get("qs")
		if key == "" {
			key = strings.ToLower(field.Name)
		} else if key == "-" {
			continue
		}
		addQueryStringValue(items, key, value.Field(i))
	}
	return items.Encode()
}

func addQueryStringValue(items url.Values, key string, v reflect.Value) {
	switch v.Kind() {
	case reflect.Bool:
		if v.Bool() {
			items.Add(key, "1")
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if v.Int() > 0 {
			items.Add(key, strconv.FormatInt(v.Int(), 10))
		}
	case reflect.Float32, reflect.Float64:
		if v.Float() > 0 {
			items.Add(key, strconv.FormatFloat(v.Float(), 'f', -1, 64))
		}
	case reflect.String:
		if v.String() != "" {
			items.Add(key, v.String())
		}
	case reflect.Ptr:
		if !v.IsNil() {
			if b, err := json.Marshal(v.Interface()); err == nil {
				items.Add(key, string(b))
			}
		}
	case reflect.Map:
		if len(v.MapKeys()) > 0 {
			if b, err := json.Marshal(v.Interface()); err == nil {
				items.Add(key, string(b))
			}
		}
	case reflect.Array, reflect.Slice:
		vLen := v.Len()
		if vLen > 0 {
			for i := 0; i < vLen; i++ {
				addQueryStringValue(items, key, v.Index(i))
			}
		}
	}
}
