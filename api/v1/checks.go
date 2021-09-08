package v1

import (
	"fmt"

	"github.com/flanksource/canary-checker/api/external"
	"github.com/flanksource/kommons"
	v1 "k8s.io/api/core/v1"
)

type Check struct {
	Type, Endpoint, Description, Icon string
}

func (c Check) GetType() string {
	return c.Type
}

func (c Check) GetEndpoint() string {
	return c.Endpoint
}

func (c Check) GetDescription() string {
	return c.Description
}

func (c Check) GetIcon() string {
	return c.Icon
}

type HTTPCheck struct {
	Description `yaml:",inline" json:",inline"`
	Templatable `yaml:",inline" json:",inline"`
	// HTTP endpoint to check.  Mutually exclusive with Namespace
	Endpoint string `yaml:"endpoint" json:"endpoint,omitempty" template:"true"`
	// Namespace to crawl for TLS endpoints.  Mutually exclusive with Endpoint
	Namespace string `yaml:"namespace" json:"namespace,omitempty" template:"true"`
	// Maximum duration in milliseconds for the HTTP request. It will fail the check if it takes longer.
	ThresholdMillis int `yaml:"thresholdMillis" json:"thresholdMillis,omitempty"`
	// Expected response codes for the HTTP Request.
	ResponseCodes []int `yaml:"responseCodes" json:"responseCodes,omitempty"`
	// Exact response content expected to be returned by the endpoint.
	ResponseContent string `yaml:"responseContent" json:"responseContent,omitempty"`
	// Path and value to of expect JSON response by the endpoint
	ResponseJSONContent JSONCheck `yaml:"responseJSONContent,omitempty" json:"responseJSONContent,omitempty"`
	// Maximum number of days until the SSL Certificate expires.
	MaxSSLExpiry int `yaml:"maxSSLExpiry" json:"maxSSLExpiry,omitempty"`
	// HTTP method to call - defaults to GET
	Method string `yaml:"method,omitempty" json:"method,omitempty"`
	//NTLM when set to true will do authentication using NTLM v1 protocol
	NTLM bool `yaml:"ntlm,omitempty" json:"ntlm,omitempty"`
	//NTLM when set to true will do authentication using NTLM v2 protocol
	NTLMv2 bool `yaml:"ntlmv2,omitempty" json:"ntlmv2,omitempty"`
	// HTTP request body contents
	Body string `yaml:"body,omitempty" json:"body,omitempty" template:"true"`
	// HTTP Header fields to be used in the query
	Headers []kommons.EnvVar `yaml:"headers,omitempty" json:"headers,omitempty"`
	// Credentials for authentication headers:
	Authentication *Authentication `yaml:"authentication,omitempty" json:"authentication,omitempty"`
}

func (c HTTPCheck) GetEndpoint() string {
	return c.Endpoint
}

func (c HTTPCheck) GetType() string {
	return "http"
}

func (c HTTPCheck) GetMethod() string {
	if c.Method != "" {
		return c.Method
	}
	return "GET"
}

type TCPCheck struct {
	Description     `yaml:",inline" json:",inline"`
	Endpoint        string `yaml:"endpoint" json:"endpoint,omitempty"`
	ThresholdMillis int64  `yaml:"thresholdMillis" json:"thresholdMillis,omitempty"`
}

func (t TCPCheck) GetEndpoint() string {
	return t.Endpoint
}

func (t TCPCheck) GetType() string {
	return "tcp"
}

type ICMPCheck struct {
	Description         `yaml:",inline" json:",inline"`
	Endpoint            string `yaml:"endpoint" json:"endpoint,omitempty"`
	ThresholdMillis     int64  `yaml:"thresholdMillis" json:"thresholdMillis,omitempty"`
	PacketLossThreshold int64  `yaml:"packetLossThreshold" json:"packetLossThreshold,omitempty"`
	PacketCount         int    `yaml:"packetCount" json:"packetCount,omitempty"`
}

func (c ICMPCheck) GetEndpoint() string {
	return c.Endpoint
}

func (c ICMPCheck) GetType() string {
	return "icmp"
}

type Bucket struct {
	Name     string `yaml:"name" json:"name,omitempty"`
	Region   string `yaml:"region" json:"region,omitempty"`
	Endpoint string `yaml:"endpoint" json:"endpoint,omitempty"`
}

type S3Check struct {
	Description `yaml:",inline" json:",inline"`
	Bucket      Bucket `yaml:"bucket" json:"bucket,omitempty"`
	AccessKey   string `yaml:"accessKey" json:"accessKey,omitempty"`
	SecretKey   string `yaml:"secretKey" json:"secretKey,omitempty"`
	ObjectPath  string `yaml:"objectPath" json:"objectPath,omitempty"`
	// Skip TLS verify when connecting to s3
	SkipTLSVerify bool `yaml:"skipTLSVerify" json:"skipTLSVerify,omitempty"`
}

func (c S3Check) GetEndpoint() string {
	return fmt.Sprintf("%s/%s", c.Bucket.Endpoint, c.Bucket.Name)
}

func (c S3Check) GetType() string {
	return "s3"
}

type AWSConnection struct {
	AccessKey kommons.EnvVar `yaml:"accessKey" json:"accessKey"`
	SecretKey kommons.EnvVar `yaml:"secretKey" json:"secretKey"`
	Region    string         `yaml:"region" json:"region"`
	Endpoint  string         `yaml:"endpoint" json:"endpoint,omitempty"`
	// Skip TLS verify when connecting to aws
	SkipTLSVerify bool `yaml:"skipTLSVerify" json:"skipTLSVerify,omitempty"`
}

type S3BucketCheck struct {
	Description   `yaml:",inline" json:",inline"`
	Templatable   `yaml:",inline" json:",inline"`
	AWSConnection `yaml:",inline" json:",inline"`
	FolderTest    `yaml:",inline" json:",inline"`
	Filter        FolderFilter `yaml:"filter,omitempty" json:"filter,omitempty"`
	Bucket        string       `yaml:"bucket" json:"bucket"`
	// Use path style path: http://s3.amazonaws.com/BUCKET/KEY instead of http://BUCKET.s3.amazonaws.com/KEY
	UsePathStyle bool `yaml:"usePathStyle" json:"usePathStyle,omitempty"`
}

func (c S3BucketCheck) GetEndpoint() string {
	if c.AWSConnection.Endpoint != "" {
		return fmt.Sprintf("%s/%s", c.AWSConnection.Endpoint, c.Bucket)
	} else {
		return c.Bucket
	}
}

func (c S3BucketCheck) GetType() string {
	return "s3Bucket"
}

type CloudWatchCheck struct {
	Description   `yaml:",inline" json:",inline"`
	AWSConnection `yaml:",inline" json:",inline"`
	Filter        CloudWatchFilter `yaml:"filter,omitempty" json:"filter,omitempty"`
}

type CloudWatchFilter struct {
	ActionPrefix *string  `yaml:"actionPrefix,omitempty" json:"actionPrefix,omitempty"`
	AlarmPrefix  *string  `yaml:"alarmPrefix,omitempty" json:"alarmPrefix,omitempty"`
	Alarms       []string `yaml:"alarms,omitempty" json:"alarms,omitempty"`
	State        string   `yaml:"state,omitempty" json:"state,omitempty"`
}

func (c CloudWatchCheck) GetEndpoint() string {
	endpoint := c.Region
	if c.Filter.ActionPrefix != nil {
		endpoint += "-" + *c.Filter.ActionPrefix
	}
	if c.Filter.AlarmPrefix != nil {
		endpoint += "-" + *c.Filter.AlarmPrefix
	}
	return endpoint
}

func (c CloudWatchCheck) GetType() string {
	return "cloudwatch"
}

type GCPConnection struct {
	Endpoint    string          `yaml:"endpoint" json:"endpoint,omitempty"`
	Credentials *kommons.EnvVar `yaml:"credentials" json:"credentials,omitempty"`
}

type GCSBucketCheck struct {
	Description   `yaml:",inline" json:",inline"`
	Templatable   `yaml:",inline" json:",inline"`
	FolderTest    `yaml:",inline" json:",inline"`
	GCPConnection `yaml:",inline" json:",inline"`
	Filter        FolderFilter `yaml:"filter,omitempty" json:"filter,omitempty"`
	Bucket        string       `yaml:"bucket" json:"bucket"`
}

func (c GCSBucketCheck) GetEndpoint() string {
	return c.Bucket
}

func (c GCSBucketCheck) GetType() string {
	return "gcsBucket"
}

type ResticCheck struct {
	Description `yaml:",inline" json:",inline"`
	// Repository The restic repository path eg: rest:https://user:pass@host:8000/ or rest:https://host:8000/ or s3:s3.amazonaws.com/bucket_name
	Repository string `yaml:"repository" json:"repository"`
	// Password for the restic repository
	Password *kommons.EnvVar `yaml:"password" json:"password"`
	// MaxAge for backup freshness
	MaxAge string `yaml:"maxAge" json:"maxAge"`
	// CheckIntegrity when enabled will check the Integrity and consistency of the restic reposiotry
	CheckIntegrity bool `yaml:"checkIntegrity,omitempty" json:"checkIntegrity,omitempty"`
	// AccessKey access key id for connection with aws s3, minio, wasabi, alibaba oss
	AccessKey *kommons.EnvVar `yaml:"accessKey,omitempty" json:"accessKey,omitempty"`
	// SecretKey secret access key for connection with aws s3, minio, wasabi, alibaba oss
	SecretKey *kommons.EnvVar `yaml:"secretKey,omitempty" json:"secretKey,omitempty"`
	// CaCert path to the root cert. In case of self-signed certificates
	CaCert string `yaml:"caCert,omitempty" json:"caCert,omitempty"`
}

func (c ResticCheck) GetEndpoint() string {
	return c.Repository
}

func (c ResticCheck) GetType() string {
	return "restic"
}

type JmeterCheck struct {
	Description `yaml:",inline" json:",inline"`
	// Jmx defines tge ConfigMap or Secret reference to get the JMX test plan
	Jmx kommons.EnvVar `yaml:"jmx" json:"jmx"`
	// Host is the server against which test plan needs to be executed
	Host string `yaml:"host,omitempty" json:"host,omitempty"`
	// Port on which the server is running
	Port int32 `yaml:"port,omitempty" json:"port,omitempty"`
	// Properties defines the local Jmeter properties
	Properties []string `yaml:"properties,omitempty" json:"properties,omitempty"`
	// SystemProperties defines the java system property
	SystemProperties []string `yaml:"systemProperties,omitempty" json:"systemProperties,omitempty"`
	// ResponseDuration under which the all the test should pass
	ResponseDuration string `yaml:"responseDuration,omitempty" json:"responseDuration,omitempty"`
}

func (c JmeterCheck) GetEndpoint() string {
	return fmt.Sprintf(c.Host + ":" + string(c.Port))
}

func (c JmeterCheck) GetType() string {
	return "jmeter"
}

type DockerPullCheck struct {
	Description    `yaml:",inline" json:",inline"`
	Image          string          `yaml:"image" json:"image,omitempty"`
	Auth           *Authentication `yaml:"auth,omitempty" json:"auth,omitempty"`
	ExpectedDigest string          `yaml:"expectedDigest" json:"expectedDigest,omitempty"`
	ExpectedSize   int64           `yaml:"expectedSize" json:"expectedSize,omitempty"`
}

func (c DockerPullCheck) GetEndpoint() string {
	return c.Image
}

func (c DockerPullCheck) GetType() string {
	return "dockerPull"
}

type DockerPushCheck struct {
	Description `yaml:",inline" json:",inline"`
	Image       string          `yaml:"image" json:"image,omitempty"`
	Auth        *Authentication `yaml:"auth" json:"auth"`
}

func (c DockerPushCheck) GetEndpoint() string {
	return c.Image
}

func (c DockerPushCheck) GetType() string {
	return "dockerPush"
}

type ContainerdPullCheck struct {
	Description    `yaml:",inline" json:",inline"`
	Image          string         `yaml:"image" json:"image,omitempty"`
	Auth           Authentication `yaml:"auth,omitempty" json:"auth,omitempty"`
	ExpectedDigest string         `yaml:"expectedDigest" json:"expectedDigest,omitempty"`
	ExpectedSize   int64          `yaml:"expectedSize" json:"expectedSize,omitempty"`
}

func (c ContainerdPullCheck) GetEndpoint() string {
	return c.Image
}

func (c ContainerdPullCheck) GetType() string {
	return "containerdPull"
}

type ContainerdPushCheck struct {
	Description `yaml:",inline" json:",inline"`
	Image       string `yaml:"image" json:"image,omitempty"`
	Username    string `yaml:"username" json:"username,omitempty"`
	Password    string `yaml:"password" json:"password,omitempty"`
}

func (c ContainerdPushCheck) GetEndpoint() string {
	return c.Image
}

func (c ContainerdPushCheck) GetType() string {
	return "containerdPush"
}

type RedisCheck struct {
	Description `yaml:",inline" json:",inline"`
	Addr        string          `yaml:"addr" json:"addr" template:"true"`
	Auth        *Authentication `yaml:"auth,omitempty" json:"auth,omitempty"`
	DB          int             `yaml:"db" json:"db"`
}

func (c RedisCheck) GetType() string {
	return "redis"
}

func (c RedisCheck) GetEndpoint() string {
	return c.Addr
}

type SQLCheck struct {
	Description `yaml:",inline" json:",inline"`
	Templatable `yaml:",inline" json:",inline"`
	Connection  `yaml:",inline" json:",inline"`
	driver      string `yaml:"-" json:"-"`
	Query       string `yaml:"query" json:"query,omitempty" template:"true"`
	// Number rows to check for
	Result int `yaml:"results" json:"results,omitempty"`
}

func (c *SQLCheck) GetQuery() string {
	if c.Query == "" {
		return "SELECT 1"
	}
	return c.Query
}

func (c SQLCheck) GetDriver() string {
	return c.driver
}

func (c *SQLCheck) SetDriver(driver string) {
	c.driver = driver
}

func (c SQLCheck) GetType() string {
	return c.GetDriver()
}

type PostgresCheck struct {
	SQLCheck `yaml:",inline" json:",inline"`
}

type MssqlCheck struct {
	SQLCheck `yaml:",inline" json:",inline"`
}

type PodCheck struct {
	Description          `yaml:",inline" json:",inline"`
	Namespace            string `yaml:"namespace" json:"namespace,omitempty" template:"true"`
	Spec                 string `yaml:"spec" json:"spec,omitempty"`
	ScheduleTimeout      int64  `yaml:"scheduleTimeout" json:"scheduleTimeout,omitempty"`
	ReadyTimeout         int64  `yaml:"readyTimeout" json:"readyTimeout,omitempty"`
	HTTPTimeout          int64  `yaml:"httpTimeout" json:"httpTimeout,omitempty"`
	DeleteTimeout        int64  `yaml:"deleteTimeout" json:"deleteTimeout,omitempty"`
	IngressTimeout       int64  `yaml:"ingressTimeout" json:"ingressTimeout,omitempty"`
	HTTPRetryInterval    int64  `yaml:"httpRetryInterval" json:"httpRetryInterval,omitempty"`
	Deadline             int64  `yaml:"deadline" json:"deadline,omitempty"`
	Port                 int64  `yaml:"port" json:"port,omitempty"`
	Path                 string `yaml:"path" json:"path,omitempty" template:"true"`
	IngressName          string `yaml:"ingressName" json:"ingressName,omitempty" template:"true" `
	IngressHost          string `yaml:"ingressHost" json:"ingressHost,omitempty" template:"true"`
	ExpectedContent      string `yaml:"expectedContent" json:"expectedContent,omitempty" template:"true"`
	ExpectedHTTPStatuses []int  `yaml:"expectedHttpStatuses" json:"expectedHttpStatuses,omitempty"`
	PriorityClass        string `yaml:"priorityClass" json:"priorityClass,omitempty"`
}

func (c PodCheck) GetEndpoint() string {
	return c.Name
}

func (c PodCheck) String() string {
	return "pod/" + c.Name
}

func (c PodCheck) GetType() string {
	return "pod"
}

type LDAPCheck struct {
	Description   `yaml:",inline" json:",inline"`
	Host          string          `yaml:"host" json:"host,omitempty" template:"true"`
	Auth          *Authentication `yaml:"auth" json:"auth,omitempty"`
	BindDN        string          `yaml:"bindDN" json:"bindDN,omitempty"`
	UserSearch    string          `yaml:"userSearch" json:"userSearch,omitempty"`
	SkipTLSVerify bool            `yaml:"skipTLSVerify" json:"skipTLSVerify,omitempty"`
}

func (c LDAPCheck) GetEndpoint() string {
	return c.Host
}

func (c LDAPCheck) GetType() string {
	return "ldap"
}

type NamespaceCheck struct {
	Description          `yaml:",inline" json:",inline"`
	CheckName            string            `yaml:"checkName" json:"checkName,omitempty" template:"true"`
	NamespaceNamePrefix  string            `yaml:"namespaceNamePrefix" json:"namespaceNamePrefix,omitempty"`
	NamespaceLabels      map[string]string `yaml:"namespaceLabels" json:"namespaceLabels,omitempty"`
	NamespaceAnnotations map[string]string `yaml:"namespaceAnnotations" json:"namespaceAnnotations,omitempty"`
	PodSpec              string            `yaml:"podSpec" json:"podSpec,omitempty"`
	ScheduleTimeout      int64             `yaml:"scheduleTimeout" json:"schedule_timeout,omitempty"`
	ReadyTimeout         int64             `yaml:"readyTimeout" json:"readyTimeout,omitempty"`
	HTTPTimeout          int64             `yaml:"httpTimeout" json:"httpTimeout,omitempty"`
	DeleteTimeout        int64             `yaml:"deleteTimeout" json:"deleteTimeout,omitempty"`
	IngressTimeout       int64             `yaml:"ingressTimeout" json:"ingressTimeout,omitempty"`
	HTTPRetryInterval    int64             `yaml:"httpRetryInterval" json:"httpRetryInterval,omitempty"`
	Deadline             int64             `yaml:"deadline" json:"deadline,omitempty"`
	Port                 int64             `yaml:"port" json:"port,omitempty"`
	Path                 string            `yaml:"path" json:"path,omitempty"`
	IngressName          string            `yaml:"ingressName" json:"ingressName,omitempty" template:"true"`
	IngressHost          string            `yaml:"ingressHost" json:"ingressHost,omitempty" template:"true"`
	ExpectedContent      string            `yaml:"expectedContent" json:"expectedContent,omitempty" template:"true"`
	ExpectedHTTPStatuses []int64           `yaml:"expectedHttpStatuses" json:"expectedHttpStatuses,omitempty"`
	PriorityClass        string            `yaml:"priorityClass" json:"priorityClass,omitempty"`
}

func (c NamespaceCheck) GetEndpoint() string {
	return c.CheckName
}

func (c NamespaceCheck) String() string {
	return "namespace/" + c.CheckName
}

func (c NamespaceCheck) GetType() string {
	return "namespace"
}

type DNSCheck struct {
	Description     `yaml:",inline" json:",inline"`
	Server          string   `yaml:"server" json:"server,omitempty"`
	Port            int      `yaml:"port" json:"port,omitempty"`
	Query           string   `yaml:"query,omitempty" json:"query,omitempty"`
	QueryType       string   `yaml:"querytype" json:"querytype,omitempty"`
	MinRecords      int      `yaml:"minrecords,omitempty" json:"minrecords,omitempty"`
	ExactReply      []string `yaml:"exactreply,omitempty" json:"exactreply,omitempty"`
	Timeout         int      `yaml:"timeout" json:"timeout,omitempty"`
	ThresholdMillis int      `yaml:"thresholdMillis" json:"thresholdMillis,omitempty"`
	// SrvReply    SrvReply `yaml:"srvReply,omitempty" json:"srvReply,omitempty"`
}

func (c DNSCheck) GetEndpoint() string {
	s := fmt.Sprintf("%s/%s", c.QueryType, c.Query)
	if c.Server != "" {
		s += "@" + c.Server
		if c.Port != 0 {
			s += fmt.Sprintf(":%d", c.Port)
		}
	}
	return s
}

func (c DNSCheck) GetType() string {
	return "dns"
}

type HelmCheck struct {
	Description `yaml:",inline" json:",inline"`
	Chartmuseum string          `yaml:"chartmuseum" json:"chartmuseum,omitempty"`
	Project     string          `yaml:"project,omitempty" json:"project,omitempty"`
	Auth        *Authentication `yaml:"auth,omitempty" json:"auth,omitempty"`
	CaFile      string          `yaml:"cafile,omitempty" json:"cafile,omitempty"`
}

func (c HelmCheck) GetEndpoint() string {
	return fmt.Sprintf("%s/%s", c.Chartmuseum, c.Project)
}

func (c HelmCheck) GetType() string {
	return "helm"
}

type JunitCheck struct {
	Description `yaml:",inline" json:",inline"`
	TestResults string `yaml:"testResults" json:"testResults"`
	Templatable `yaml:",inline" json:",inline"`
	// Timeout in minutes to wait for specified container to finish its job. Defaults to 5 minutes
	Timeout int        `yaml:"timeout,omitempty" json:"timeout,omitempty"`
	Spec    v1.PodSpec `yaml:"spec" json:"spec"`
}

func (c JunitCheck) GetEndpoint() string {
	if c.Description.String() != "" {
		return c.Description.String()
	}
	if len(c.Spec.Containers) > 0 {
		if c.Spec.Containers[0].Name != "" {
			return c.Spec.Containers[0].Name
		}
		if c.Spec.Containers[0].Image != "" {
			return c.Spec.Containers[0].Image
		}
	}
	return c.TestResults
}

func (c JunitCheck) GetTimeout() int {
	if c.Timeout != 0 {
		return c.Timeout
	}
	return 5
}

func (c JunitCheck) GetType() string {
	return "junit"
}

type SmbCheck struct {
	Description `yaml:",inline" json:",inline"`
	Templatable `yaml:",inline" json:",inline"`
	Filter      FolderFilter `yaml:"filter,omitempty" json:"filter,omitempty"`
	FolderTest  `yaml:",inline" json:",inline"`
	//Server location of smb server. Can be hostname/ip or in '\\server\e$\a\b\c' syntax
	//Where server is the hostname e$ is the sharename and a/b/c is the searchPath location
	Server string `yaml:"server" json:"server"`
	//Port on which smb server is running. Defaults to 445
	Port int             `yaml:"port,omitempty" json:"port,omitempty"`
	Auth *Authentication `yaml:"auth" json:"auth"`
	//Domain...
	Domain string `yaml:"domain,omitempty" json:"domain,omitempty"`
	// Workstation...
	Workstation string `yaml:"workstation,omitempty" json:"workstation,omitempty"`
	//Sharename to mount from the samba server
	Sharename string `yaml:"sharename,omitempty" json:"sharename,omitempty"`
	//SearchPath sub-path inside the mount location
	SearchPath string `yaml:"searchPath,omitempty" json:"searchPath,omitempty" `
}

func (c SmbCheck) GetEndpoint() string {
	return fmt.Sprintf("%s:%d/%s-%s", c.Server, c.GetPort(), c.Sharename, c.Description)
}

func (c SmbCheck) GetType() string {
	return "smb"
}

func (c SmbCheck) GetPort() int {
	if c.Port != 0 {
		return c.Port
	}
	return 445
}

type PrometheusCheck struct {
	Description `yaml:",inline" json:",inline"`
	Templatable `yaml:",inline" json:",inline"`
	// Address of the prometheus server
	Host string `yaml:"host" json:"host" template:"true" `
	// PromQL query
	Query string `yaml:"query" json:"query" template:"true"`
}

func (c PrometheusCheck) GetType() string {
	return "prometheus"
}

func (c PrometheusCheck) GetEndpoint() string {
	return fmt.Sprintf("%v/%v", c.Host, c.Description)
}

type MongoDBCheck struct {
	Description `yaml:",inline" json:",inline"`
	// Monogodb connection string, e.g.  mongodb://:27017/?authSource=admin, See https://docs.mongodb.com/manual/reference/connection-string/
	Connection `yaml:",inline" json:",inline"`
}

func (c MongoDBCheck) GetType() string {
	return "mongodb"
}

/*

```yaml
http:
  - endpoints:
      - https://httpstat.us/200
      - https://httpstat.us/301
    thresholdMillis: 3000
    responseCodes: [201,200,301]
    responseContent: ""
    maxSSLExpiry: 60
  - endpoints:
      - https://httpstat.us/500
    thresholdMillis: 3000
    responseCodes: [500]
    responseContent: ""
    maxSSLExpiry: 60
  - endpoints:
      - https://httpstat.us/500
    thresholdMillis: 3000
    responseCodes: [302]
    responseContent: ""
    maxSSLExpiry: 60
```
*/
type HTTP struct {
	HTTPCheck `yaml:",inline" json:"inline"`
}

/*

```yaml
dns:
  - server: 8.8.8.8
    port: 53
    query: "flanksource.com"
    querytype: "A"
    minrecords: 1
    exactreply: ["34.65.228.161"]
    timeout: 10
```
*/
type DNS struct {
	DNSCheck `yaml:",inline" json:"inline"`
}

/*
DockerPull check will try to pull a Docker image from specified registry, verify it's checksum and size.

```yaml

docker:
  - image: docker.io/library/busybox:1.31.1
    auth:
		username:
			value: some-user
		password:
			value: some-password
    expectedDigest: 6915be4043561d64e0ab0f8f098dc2ac48e077fe23f488ac24b665166898115a
    expectedSize: 1219782
```

*/
type DockerPull struct {
	DockerPullCheck `yaml:",inline" json:"inline"`
}

/*
DockerPush check will try to push a Docker image to specified registry.

```yaml

dockerPush:
  - image: ttl.sh/flanksource-busybox:1.30
    auth:
      username:
        value: $DOCKER_USERNAME
      password:
        value: $DOCKER_PASSWORD
```

*/
type DockerPush struct {
	DockerPushCheck `yaml:",inline" json:"inline"`
}

/*
S3 check will:

* list objects in the bucket to check for Read permissions
* PUT an object into the bucket for Write permissions
* download previous uploaded object to check for Get permissions

```yaml

s3:
  - buckets:
      - name: "test-bucket"
        region: "us-east-1"
        endpoint: "https://test-bucket.s3.us-east-1.amazonaws.com"
    secretKey: "<access-key>"
    accessKey: "<secret-key>"
    objectPath: "path/to/object"
```
*/
type S3 struct {
	S3Check `yaml:",inline" json:"inline"`
}

/*
This check will

- search objects matching the provided object path pattern
- check that latest object is no older than provided MaxAge value in seconds
- check that latest object size is not smaller than provided MinSize value in bytes.

```yaml
s3Bucket:
  - bucket: foo
    accessKey: "<access-key>"
    secretKey: "<secret-key>"
    region: "us-east-2"
    endpoint: "https://s3.us-east-2.amazonaws.com"
    objectPath: "(.*)archive.zip$"
    readWrite: true
    maxAge: 5000000
    minSize: 50000
```
*/
type S3Bucket struct {
	S3BucketCheck `yaml:",inline" json:"inline"`
}

type TCP struct {
	TCPCheck `yaml:",inline" json:"inline"`
}

/*
```yaml
pod:
  - name: golang
    namespace: default
    spec: |
      apiVersion: v1
      kind: Pod
      metadata:
        name: hello-world-golang
        namespace: default
        labels:
          app: hello-world-golang
      spec:
        containers:
          - name: hello
            image: quay.io/toni0/hello-webserver-golang:latest
    port: 8080
    path: /foo/bar
    ingressName: hello-world-golang
    ingressHost: "hello-world-golang.127.0.0.1.nip.io"
    scheduleTimeout: 2000
    readyTimeout: 5000
    httpTimeout: 2000
    deleteTimeout: 12000
    ingressTimeout: 5000
    deadline: 29000
    httpRetryInterval: 200
    expectedContent: bar
    expectedHttpStatuses: [200, 201, 202]
```
*/
type Pod struct {
	PodCheck `yaml:",inline" json:"inline"`
}

/*

The LDAP check will:

* bind using provided user/password to the ldap host. Supports ldap/ldaps protocols.
* search an object type in the provided bind DN.s

```yaml

ldap:
  - host: ldap://127.0.0.1:10389
    auth:
      username:
        value: uid=admin,ou=system
      password:
        value: secret
    bindDN: ou=users,dc=example,dc=com
    userSearch: "(&(objectClass=organizationalPerson))"
  - host: ldap://127.0.0.1:10389
    auth:
      username:
        value: uid=admin,ou=system
      password:
        value: secret
    bindDN: ou=groups,dc=example,dc=com
    userSearch: "(&(objectClass=groupOfNames))"
```
*/
type LDAP struct {
	LDAPCheck `yaml:",inline" json:"inline"`
}

/*

The Namespace check will:

* create a new namespace using the labels/annotations provided

```yaml

namespace:
  - namePrefix: "test-name-prefix-"
		labels:
			team: test
		annotations:
			"foo.baz.com/foo": "bar"
```
*/
type Namespace struct {
	NamespaceCheck `yaml:",inline" json:"inline"`
}

/*
This test will check ICMP packet loss and duration.

```yaml

icmp:
  - endpoints:
      - https://google.com
      - https://yahoo.com
    thresholdMillis: 400
    packetLossThreshold: 0.5
    packetCount: 2
```
*/
type ICMP struct {
	ICMPCheck `yaml:",inline" json:"inline"`
}

/*
This check will try to connect to a specified Postgresql database, run a query against it and verify the results.

```yaml

postgres:
  - connection: "user=postgres password=mysecretpassword host=192.168.0.103 port=15432 dbname=postgres sslmode=disable"
    query:  "SELECT 1"
		results: 1
```
*/
type Postgres struct {
	PostgresCheck `yaml:",inline" json:"inline"`
}

/*
This check will try to connect to a specified MsSQL database, run a query against it and verify the results.

```yaml

mssql:
  - connection: 'server=localhost;user id=sa;password=Some_S3cure_p@sswd;port=1433;database=test'
    query: "SELECT 1"
	results: 1
```
*/
type MsSQL struct {
	MssqlCheck `yaml:",inline" json:"inline"`
}

type Helm struct {
	HelmCheck `yaml:",inline" json:"inline"`
}

type SrvReply struct {
	Target   string `yaml:"target,omitempty"`
	Port     int    `yaml:"port,omitempty"`
	Priority int    `yaml:"priority,omitempty"`
	Weight   int    `yaml:"wight,omitempty"`
}

/*
This check will try to connect to a specified Redis instance, run a ping against it and verify the pong response.

```yaml

redis:
  - addr: "redis-service.default:6379"
	db: 0
	description: "The redis test"
```
*/

type Redis struct {
	RedisCheck `yaml:",inline" json:"inline"`
}

/*

This check will connect to a restic repository and perform Integrity and backupFreshness Tests
```yaml
restic:
	- repository: s3:http://minio.infra/restic-repo
      password:
        value: S0M3p@sswd
      maxAge: 5h30m
      checkIntegrity: true
      accessKey:
        value: some-access-key
      secretKey:
        value: some-secret-key
      description: The restic test
```
*/

type Restic struct {
	ResticCheck `yaml:",inline" json:"inline"`
}

/*
Jmeter check will run jmeter cli against the supplied host
```yaml
jmeter:
    - jmx:
        name: jmx-test-plan
        valueFrom:
          configMapKeyRef:
             key: jmeter-test.xml
             name: jmeter
      host: "some-host"
      port: 8080
      properties:
        - remote_hosts=127.0.0.1
      systemProperties:
        - user.dir=/home/mstover/jmeter_stuff
      description: The Jmeter test
```
*/
type Jmeter struct {
	JmeterCheck `yaml:",inline" json:",inline"`
}

/*
Junit check will wait for the given pod to be completed than parses all the xml files present in the defined testResults directory
```yaml
junit:
  - testResults: "/tmp/junit-results/"
	description: "junit demo test"
    spec:
      containers:
        - name: jes
          image: docker.io/tarun18/junit-test-pass
          command: ["/start.sh"]
```
*/
type Junit struct {
	JunitCheck `yaml:",inline" json:",inline"`
}

/*
Smb check will connect to the given samba server with given credentials
find the age of the latest updated file and compare it with minAge
count the number of file present and compare with minCount if defined
```yaml
smb:
   - server: 192.168.1.9
	 auth:
       username:
          value: samba
       password:
           value: password
     sharename: "Some Public Folder"
     minAge: 10h
	 maxAge: 20h
	 searchPath: a/b/c
     description: "Success SMB server"
```

User can define server in `\\server\e$\a\b\c` format where `server` is the host
`e$` is the sharename and `a/b/c` represent the sub-dir inside mount location where the test will run to verify
```yaml
smb:
   - server: '\\192.168.1.5\Some Public Folder\somedir'
     auth:
		 username:
		   value: samba
		 password:
		   valueFrom:
			 secretKeyRef:
			   key: smb-password
			   name: smb
     sharename: "Tarun Khandelwal’s Public Folder"
     minAge: 10h
     maxAge: 100h
     description: "Success SMB server"
```
*/
type Smb struct {
	SmbCheck `yaml:",inline" json:",inline"`
}

type EC2Check struct {
	Description   `yaml:",inline" json:",inline"`
	AWSConnection `yaml:",inline" json:",inline"`
	AMI           string                    `yaml:"ami,omitempty" json:"ami,omitempty"`
	UserData      string                    `yaml:"userData,omitempty" json:"userData,omitempty"`
	SecurityGroup string                    `yaml:"securityGroup,omitempty" json:"securityGroup,omitempty"`
	KeepAlive     bool                      `yaml:"keepAlive,omitempty" json:"keepAlive,omitempty"`
	WaitTime      int                       `yaml:"waitTime,omitempty" json:"waitTime,omitempty"`
	TimeOut       int                       `yaml:"timeOut,omitempty" json:"timeOut,omitempty"`
	CanaryRef     []v1.LocalObjectReference `yaml:"canaryRef,omitempty" json:"canaryRef,omitempty"`
}

func (c EC2Check) GetEndpoint() string {
	return c.Region
}

func (c EC2Check) GetType() string {
	return "ec2"
}

var AllChecks = []external.Check{
	HTTPCheck{},
	TCPCheck{},
	ICMPCheck{},
	S3Check{},
	S3BucketCheck{},
	DockerPullCheck{},
	DockerPushCheck{},
	ContainerdPullCheck{},
	ContainerdPushCheck{},
	PostgresCheck{},
	MssqlCheck{},
	RedisCheck{},
	PodCheck{},
	LDAPCheck{},
	ResticCheck{},
	NamespaceCheck{},
	DNSCheck{},
	HelmCheck{},
	JmeterCheck{},
	JunitCheck{},
	SmbCheck{},
	EC2Check{},
	PrometheusCheck{},
	GCSBucketCheck{},
	MongoDBCheck{},
	CloudWatchCheck{},
}
