package docker

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	//"net/http"
	"time"
)

type Port struct {
	PrivatePort int64
	PublicPort  int64
	Type        string
}

// APIContainers: ListContainers
type APIContainers struct {
	Id         string `json:"Id"`
	Image      string `json:"Image,omitempty"`
	Command    string `json:"Command,omitempty"`
	Created    int64  `json:"Created,omitempty"`
	Status     string `json:"Status,omitempty"`
	Ports      []Port `json:"Ports,omitempty"`
	SizeRw     int64  `json:"SizeRw,omitempty"`
	SizeRootFS int64  `json:"SizeRootFS,omitempty"`
}

// State represents the state of a container.
type State struct {
	Running    bool      `json:"Running,omitempty" yaml:"Running,omitempty"`
	Paused     bool      `json:"Paused,omitempty" yaml:"Paused,omitempty"`
	Restarting bool      `json:"Restarting,omitempty" yaml:"Restarting,omitempty"`
	OOMKilled  bool      `json:"OOMKilled,omitempty" yaml:"OOMKilled,omitempty"`
	Pid        int       `json:"Pid,omitempty" yaml:"Pid,omitempty"`
	ExitCode   int       `json:"ExitCode,omitempty" yaml:"ExitCode,omitempty"`
	Error      string    `json:"Error,omitempty" yaml:"Error,omitempty"`
	StartedAt  time.Time `json:"StartedAt,omitempty" yaml:"StartedAt,omitempty"`
	FinishedAt time.Time `json:"FinishedAt,omitempty" yaml:"FinishedAt,omitempty"`
}

// SwarmNode containers information about which Swarm node the container is on
type SwarmNode struct {
	ID     string            `json:"ID,omitempty" yaml:"ID,omitempty"`
	IP     string            `json:"IP,omitempty" yaml:"IP,omitempty"`
	Addr   string            `json:"Addr,omitempty" yaml:"Addr,omitempty"`
	Name   string            `json:"Name,omitempty" yaml:"Name,omitempty"`
	CPUs   int64             `json:"CPUs,omitempty" yaml:"CPUs,omitempty"`
	Memory int64             `json:"Memory,omitempty" yaml:"Memory,omitempty"`
	Labels map[string]string `json:"Labels,omitempty" yaml:"Labels,omitempty"`
}

// PortBinding represents the host/container port mapping as returned in the
// `docker inspect` json
type PortBinding struct {
	HostIP   string `json:"HostIP,omitempty" yaml:"HostIP,omitempty"`
	HostPort string `json:"HostPort,omitempty" yaml:"HostPort,omitempty"`
}

// PortMapping represents a deprecated field in the `docker inspect` output,
// and its value as found in NetworkSettings should always be nil
type PortMapping map[string]string

// NetworkSettings contains network-related information about a container
type NetworkSettings struct {
	IPAddress   string                 `json:"IPAddress,omitempty" yaml:"IPAddress,omitempty"`
	IPPrefixLen int                    `json:"IPPrefixLen,omitempty" yaml:"IPPrefixLen,omitempty"`
	Gateway     string                 `json:"Gateway,omitempty" yaml:"Gateway,omitempty"`
	Bridge      string                 `json:"Bridge,omitempty" yaml:"Bridge,omitempty"`
	PortMapping map[string]PortMapping `json:"PortMapping,omitempty" yaml:"PortMapping,omitempty"`
	Ports       map[Port][]PortBinding `json:"Ports,omitempty" yaml:"Ports,omitempty"`
}

// Container:CreateContainers
type Container struct {
	ID string `json:"Id" yaml:"Id"`

	Created time.Time `json:"Created,omitempty" yaml:"Created,omitempty"`

	Path string   `json:"Path,omitempty" yaml:"Path,omitempty"`
	Args []string `json:"Args,omitempty" yaml:"Args,omitempty"`

	Config *Config `json:"Config,omitempty" yaml:"Config,omitempty"`
	State  State   `json:"State,omitempty" yaml:"State,omitempty"`
	Image  string  `json:"Image,omitempty" yaml:"Image,omitempty"`

	Node *SwarmNode `json:"Node,omitempty" yaml:"Node,omitempty"`

	NetworkSettings *NetworkSettings `json:"NetworkSettings,omitempty" yaml:"NetworkSettings,omitempty"`

	SysInitPath    string `json:"SysInitPath,omitempty" yaml:"SysInitPath,omitempty"`
	ResolvConfPath string `json:"ResolvConfPath,omitempty" yaml:"ResolvConfPath,omitempty"`
	HostnamePath   string `json:"HostnamePath,omitempty" yaml:"HostnamePath,omitempty"`
	HostsPath      string `json:"HostsPath,omitempty" yaml:"HostsPath,omitempty"`
	Name           string `json:"Name,omitempty" yaml:"Name,omitempty"`
	Driver         string `json:"Driver,omitempty" yaml:"Driver,omitempty"`

	Volumes    map[string]string `json:"Volumes,omitempty" yaml:"Volumes,omitempty"`
	VolumesRW  map[string]bool   `json:"VolumesRW,omitempty" yaml:"VolumesRW,omitempty"`
	HostConfig *HostConfig       `json:"HostConfig,omitempty" yaml:"HostConfig,omitempty"`
	ExecIDs    []string          `json:"ExecIDs,omitempty" yaml:"ExecIDs,omitempty"`

	AppArmorProfile string `json:"AppArmorProfile,omitempty" yaml:"AppArmorProfile,omitempty"`
}

// KeyValuePair is a type for generic key/value pairs as used in the Lxc
// configuration
type KeyValuePair struct {
	Key   string `json:"Key,omitempty" yaml:"Key,omitempty"`
	Value string `json:"Value,omitempty" yaml:"Value,omitempty"`
}

// RestartPolicy represents the policy for automatically restarting a container.
//
// Possible values are:
//
//   - always: the docker daemon will always restart the container
//   - on-failure: the docker daemon will restart the container on failures, at
//                 most MaximumRetryCount times
//   - no: the docker daemon will not restart the container automatically
type RestartPolicy struct {
	Name              string `json:"Name,omitempty" yaml:"Name,omitempty"`
	MaximumRetryCount int    `json:"MaximumRetryCount,omitempty" yaml:"MaximumRetryCount,omitempty"`
}

// Device represents a device mapping between the Docker host and the
// container.
type Device struct {
	PathOnHost        string `json:"PathOnHost,omitempty" yaml:"PathOnHost,omitempty"`
	PathInContainer   string `json:"PathInContainer,omitempty" yaml:"PathInContainer,omitempty"`
	CgroupPermissions string `json:"CgroupPermissions,omitempty" yaml:"CgroupPermissions,omitempty"`
}

// LogConfig defines the log driver type and the configuration for it.
type LogConfig struct {
	Type   string            `json:"Type,omitempty" yaml:"Type,omitempty"`
	Config map[string]string `json:"Config,omitempty" yaml:"Config,omitempty"`
}

// HostConfig contains the container options related to starting a container on
// a given host
type HostConfig struct {
	Binds           []string               `json:"Binds,omitempty" yaml:"Binds,omitempty"`
	CapAdd          []string               `json:"CapAdd,omitempty" yaml:"CapAdd,omitempty"`
	CapDrop         []string               `json:"CapDrop,omitempty" yaml:"CapDrop,omitempty"`
	ContainerIDFile string                 `json:"ContainerIDFile,omitempty" yaml:"ContainerIDFile,omitempty"`
	LxcConf         []KeyValuePair         `json:"LxcConf,omitempty" yaml:"LxcConf,omitempty"`
	Privileged      bool                   `json:"Privileged,omitempty" yaml:"Privileged,omitempty"`
	PortBindings    map[Port][]PortBinding `json:"PortBindings,omitempty" yaml:"PortBindings,omitempty"`
	Links           []string               `json:"Links,omitempty" yaml:"Links,omitempty"`
	PublishAllPorts bool                   `json:"PublishAllPorts,omitempty" yaml:"PublishAllPorts,omitempty"`
	DNS             []string               `json:"Dns,omitempty" yaml:"Dns,omitempty"` // For Docker API v1.10 and above only
	DNSSearch       []string               `json:"DnsSearch,omitempty" yaml:"DnsSearch,omitempty"`
	ExtraHosts      []string               `json:"ExtraHosts,omitempty" yaml:"ExtraHosts,omitempty"`
	VolumesFrom     []string               `json:"VolumesFrom,omitempty" yaml:"VolumesFrom,omitempty"`
	NetworkMode     string                 `json:"NetworkMode,omitempty" yaml:"NetworkMode,omitempty"`
	IpcMode         string                 `json:"IpcMode,omitempty" yaml:"IpcMode,omitempty"`
	PidMode         string                 `json:"PidMode,omitempty" yaml:"PidMode,omitempty"`
	RestartPolicy   RestartPolicy          `json:"RestartPolicy,omitempty" yaml:"RestartPolicy,omitempty"`
	Devices         []Device               `json:"Devices,omitempty" yaml:"Devices,omitempty"`
	LogConfig       LogConfig              `json:"LogConfig,omitempty" yaml:"LogConfig,omitempty"`
	ReadonlyRootfs  bool                   `json:"ReadonlyRootfs,omitempty" yaml:"ReadonlyRootfs,omitempty"`
	SecurityOpt     []string               `json:"SecurityOpt,omitempty" yaml:"SecurityOpt,omitempty"`
	CgroupParent    string                 `json:"CgroupParent,omitempty" yaml:"CgroupParent,omitempty"`
	Memory          int64                  `json:"Memory,omitempty" yaml:"Memory,omitempty"`
	MemorySwap      int64                  `json:"MemorySwap,omitempty" yaml:"MemorySwap,omitempty"`
	CPUShares       int64                  `json:"CpuShares,omitempty" yaml:"CpuShares,omitempty"`
	CPUSet          string                 `json:"Cpuset,omitempty" yaml:"Cpuset,omitempty"`
	CPUQuota        int64                  `json:"CpuQuota,omitempty" yaml:"CpuQuota,omitempty"`
	CPUPeriod       int64                  `json:"CpuPeriod,omitempty" yaml:"CpuPeriod,omitempty"`
}

// Config is the list of configuration options used when creating a container.
// Config does not contain the options that are specific to starting a container on a
// given host.  Those are contained in HostConfig
type Config struct {
	Hostname        string              `json:"Hostname,omitempty" yaml:"Hostname,omitempty"`
	Domainname      string              `json:"Domainname,omitempty" yaml:"Domainname,omitempty"`
	User            string              `json:"User,omitempty" yaml:"User,omitempty"`
	Memory          int64               `json:"Memory,omitempty" yaml:"Memory,omitempty"`
	MemorySwap      int64               `json:"MemorySwap,omitempty" yaml:"MemorySwap,omitempty"`
	CPUShares       int64               `json:"CpuShares,omitempty" yaml:"CpuShares,omitempty"`
	CPUSet          string              `json:"Cpuset,omitempty" yaml:"Cpuset,omitempty"`
	AttachStdin     bool                `json:"AttachStdin,omitempty" yaml:"AttachStdin,omitempty"`
	AttachStdout    bool                `json:"AttachStdout,omitempty" yaml:"AttachStdout,omitempty"`
	AttachStderr    bool                `json:"AttachStderr,omitempty" yaml:"AttachStderr,omitempty"`
	PortSpecs       []string            `json:"PortSpecs,omitempty" yaml:"PortSpecs,omitempty"`
	ExposedPorts    map[Port]struct{}   `json:"ExposedPorts,omitempty" yaml:"ExposedPorts,omitempty"`
	Tty             bool                `json:"Tty,omitempty" yaml:"Tty,omitempty"`
	OpenStdin       bool                `json:"OpenStdin,omitempty" yaml:"OpenStdin,omitempty"`
	StdinOnce       bool                `json:"StdinOnce,omitempty" yaml:"StdinOnce,omitempty"`
	Env             []string            `json:"Env,omitempty" yaml:"Env,omitempty"`
	Cmd             []string            `json:"Cmd,omitempty" yaml:"Cmd,omitempty"`
	DNS             []string            `json:"Dns,omitempty" yaml:"Dns,omitempty"` // For Docker API v1.9 and below only
	Image           string              `json:"Image,omitempty" yaml:"Image,omitempty"`
	Volumes         map[string]struct{} `json:"Volumes,omitempty" yaml:"Volumes,omitempty"`
	VolumesFrom     string              `json:"VolumesFrom,omitempty" yaml:"VolumesFrom,omitempty"`
	WorkingDir      string              `json:"WorkingDir,omitempty" yaml:"WorkingDir,omitempty"`
	MacAddress      string              `json:"MacAddress,omitempty" yaml:"MacAddress,omitempty"`
	Entrypoint      []string            `json:"Entrypoint,omitempty" yaml:"Entrypoint,omitempty"`
	NetworkDisabled bool                `json:"NetworkDisabled,omitempty" yaml:"NetworkDisabled,omitempty"`
	SecurityOpts    []string            `json:"SecurityOpts,omitempty" yaml:"SecurityOpts,omitempty"`
	OnBuild         []string            `json:"OnBuild,omitempty" yaml:"OnBuild,omitempty"`
	Labels          map[string]string   `json:"Labels,omitempty" yaml:"Labels,omitempty"`
}

func (c *Client) ListContainers() ([]APIContainers, error) {
	method := "GET"
	url := c.endpoint + "/containers/json"
	//multi param: GET /containers/json?a=1&size=1

	body, _, err := c.do(method, url, DoOption{})
	if err != nil {
		return nil, err
	}

	var contanier []APIContainers
	err = json.Unmarshal(body, &contanier)
	if err != nil {
		return nil, err
	}
	return contanier, nil
}

type CreateContainerOption struct {
	Name       string
	Config     *Config
	HostConfig *HostConfig
}

func (c *Client) CreateContainers(opt CreateContainerOption) (*Container, error) {
	method := "POST"
	url := c.endpoint + "/containers/create"

	body, _, err := c.do(
		method, url, DoOption{
			data: struct {
				*Config
				*HostConfig
			}{
				opt.Config,
				opt.HostConfig,
			},
		})
	//fmt.Println(string(body))
	var container *Container
	err = json.Unmarshal(body, &container)
	if err != nil {
		return nil, err
	}
	return container, nil
}

type GetContainerLogOption struct {
	Follow     bool
	Stdout     bool
	Stderr     bool
	Timestamps bool
	Tail       string
	Container  string    `qs:"-"`
	OutStream  io.Writer `qs:"-"`
	ErrStream  io.Writer `qs:"-"`

	SetRawTerminal bool
}

func (c *Client) GetContainerLogs(opt GetContainerLogOption) error {

	if opt.Container == "" {
		return errors.New("Not specified containers.")
	}
	if opt.Tail == "" {
		opt.Tail = "all"
	}
	method := "GET"
	url := c.endpoint + "/containers/" + opt.Container + "/logs?" + queryString(opt)

	return c.stream(method, url, StreamOption{
		stdout:         opt.OutStream,
		stderr:         opt.ErrStream,
		setRawTerminal: opt.SetRawTerminal,
	})
}

type StopContainerOption struct {
	Time      int64  `qs:"t"`
	Container string `qs:"-"`
}

func (c *Client) StopContainer(opt StopContainerOption) error {

	method := "POST"
	//	url := c.endpoint + "/containers/" + opt.container + "/stop?" + queryString(opt)
	url := fmt.Sprintf("%s/containers/%s/stop?%s", c.endpoint, opt.Container, queryString(opt))

	_, code, err := c.do(method, url, DoOption{})
	if code == 404 { // code == http.StatusNotFound
		return errors.New("No such id")
	} else if code == 500 {
		return errors.New("Server error")
	} else if err != nil {
		return err
	} else {
		return nil
	}

}
