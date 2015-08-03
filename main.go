package main

import (
    "fmt"
    "sort"
    "regexp"
    "log"
    "os"
    "os/user"
    "strings"
    "path/filepath"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/credentials"
    "github.com/aws/aws-sdk-go/service/s3"
    "gopkg.in/alecthomas/kingpin.v2"
)

type Manifests []string

func (m Manifests) Len() int {
    return len(m)
}

func (m Manifests) Swap(i, j int) {
    m[i], m[j] = m[j], m[i]
}

func (m Manifests) Less(i, j int) bool {
    return m[i] < m[j]
}

func Exists(name string) bool {
    _, err := os.Stat(name)
    return !os.IsNotExist(err)
}

func NotExists(name string) bool {
    return !Exists(name)
}

func AbsPath(fname string) (f string, err error) {
    var fpath string
    matched, _ := regexp.Match("^~/", []byte(fname))
    if matched {
        usr, _ := user.Current()
        fpath = strings.Replace(fname, "~", usr.HomeDir, 1)
    } else {
        fpath, err = filepath.Abs(fname)
    }

    return fpath, err
}

func GetAccessKey() (accessKey string) {
   accessKey = os.Getenv("AWS_ACCESS_KEY")
   if accessKey == "" {
       accessKey = os.Getenv("AWS_ACCESS_KEY_ID")
   }

   return accessKey
}

func GetSecretKey() (secretKey string) {
   secretKey = os.Getenv("AWS_SECRET_KEY")
   if secretKey == "" {
       secretKey = os.Getenv("AWS_SECRET_ACCESS_KEY")
   }

   return secretKey
}

func CredentialFilePath(file string) (credentialFile string, err error) {
    if file == "" {
        credentialFile, err = AbsPath(credentials1)
    } else {
        credentialFile, err = AbsPath(file)
    }

    if file == "" && NotExists(credentialFile) {
        credentialFile, err = AbsPath(credentials2)
    }

    return credentialFile, err
}

func AWSConfig(creds *credentials.Credentials, region string) (config aws.Config) {
    config = aws.Config{
        Credentials: creds,
    }

    if region != "" {
        config.Region = aws.String(region)
    }

    return config
}

func S3ListObjectOutput(cli *s3.S3, bucket string) (listObjectOutput *s3.ListObjectsOutput, err error) {
    listObjectOutput, err = cli.ListObjects(&s3.ListObjectsInput{
        Bucket: aws.String(bucket),
    })

    return listObjectOutput, err
}

func GetManifests(listObjectOutput *s3.ListObjectsOutput) (manifests Manifests) {
    re := regexp.MustCompile(".*.ya?ml$")
    for _, c := range listObjectOutput.Contents {
        if re.Match([]byte(*c.Key)) {
            manifests = append(manifests, *c.Key)
        }
    }

    return manifests
}

func OutputAll(manifests Manifests) {
    for _, manifest := range manifests {
        fmt.Println(manifest)
    }
    os.Exit(0) 
}

func GetManifest(manifests Manifests, num int) (manifest string) {
    if num > manifests.Len() || num <= 0 {
        log.Fatalf("Out of range (%d items)", manifests.Len())
    }

    return manifests[num]
}

func Nth(manifests Manifests, num int) {
    fmt.Println(GetManifest(manifests, num - 1))
    os.Exit(0)
}

func Latest(manifests Manifests) {
    fmt.Println(GetManifest(manifests, 0))
    os.Exit(0)
}

func Oldest(manifests Manifests) {
    fmt.Println(GetManifest(manifests, manifests.Len() - 1))
    os.Exit(0)
}

var (
    all = kingpin.Flag("all", "All manifests").Short('a').Bool()
    num     = kingpin.Flag("num", "N-th Manifest").Short('n').Default("1").Int()
    bucket  = kingpin.Flag("bucket", "Bucket").Short('b').Required().String()
    region  = kingpin.Flag("region", "Region").Default("ap-northeast-1").String()
    file    = kingpin.Flag("file", "Credentials file(Default ~/.aws/credentials, ~/.aws/config)").Short('f').String()
    profile = kingpin.Flag("profile", "Profile").Default("default").String()
    oldest = kingpin.Flag("oldest", "The oldest manifest").Bool()
    credentials1  = "~/.aws/credentials"
    credentials2  = "~/.aws/config"
)

func main() {
    kingpin.Version("0.0.1")
    kingpin.Parse()

    var manifests Manifests
    var creds *credentials.Credentials
    var config aws.Config
    var credentialFile string
    var err error

    credentialFile, err = CredentialFilePath(*file)
    if err != nil {
        log.Fatal(err)
    }

    if Exists(credentialFile) {
        creds = credentials.NewSharedCredentials(credentialFile, *profile)
    } else {
        accessKey := GetAccessKey()
        secretKey := GetSecretKey()
        creds = credentials.NewStaticCredentials(accessKey, secretKey, "")
    }
    
    config = AWSConfig(creds, *region)
    cli := s3.New(&config)

    listObjectOutput, err := S3ListObjectOutput(cli, *bucket)
    if err != nil {
        log.Fatal(err)
    }

    manifests = GetManifests(listObjectOutput)
    
    sort.Sort(sort.Reverse(manifests))
    
    if *all {
        OutputAll(manifests)
    }

    if *oldest {
        Oldest(manifests)
    }

    Nth(manifests, *num)
}
