package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"golang.org/x/net/context/ctxhttp"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	_ "github.com/aws/aws-xray-sdk-go/plugins/ec2"
	_ "github.com/aws/aws-xray-sdk-go/plugins/ecs"
	"github.com/aws/aws-xray-sdk-go/xray"

	"github.com/pottava/dockerized-aws-x-ray/sample/src/lib"
)

const thisApplicationsName = "myApp"

func main() {
	xray.Configure(xray.Config{LogLevel: env("AWS_XRAY_LOG_LEVEL", "info")})

	http.Handle("/", wrap(index))
	http.Handle("/http/", wrap(httpRequests))
	http.Handle("/db/", wrap(database))
	http.Handle("/s3/", wrap(s3ListFiles))
	http.Handle("/delay/", wrap(delay))
	http.Handle("/mixed/", wrap(mixed))

	log.Printf("[service] listening on port %s", env("PORT", "80"))
	log.Fatal(http.ListenAndServe(":"+env("PORT", "80"), nil))
}

func env(key, def string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return def
	}
	return value
}

func wrap(handler http.HandlerFunc) http.Handler {
	return xray.Handler(xray.NewFixedSegmentNamer(thisApplicationsName), http.HandlerFunc(handler))
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, `<!DOCTYPE html><html lang="en">
<head>
  <meta charset="utf-8">
  <title>Samples | AWS X-Ray SDK for Go</title>
  <link rel="icon" href="data:;base64,=">
</head><body>
<h3>HTTP Requests</h3>
<ul>
  <li><a href="/http/200">External HTTP Requests (200)</a></li>
  <li><a href="/http/400">Internal HTTP Requests (400: Bad Request)</a></li>
  <li><a href="/http/401">Internal HTTP Requests (401: Unauthorized)</a></li>
  <li><a href="/http/403">Internal HTTP Requests (403: Forbidden)</a></li>
  <li><a href="/http/500">Internal HTTP Requests (500: Internal Server Error)</a></li>
  <li>..</li>
</ul>
<h3>Accessing databases</h3>
<ul>
  <li><a href="/db/">Internal MySQL</a></li>
</ul>
<h3>AWS services</h3>
<ul>
  <li><a href="/s3/">List s3 buckets</a></li>
  <li><a href="/s3/your-bucket-name">List s3 objects</a></li>
</ul>
<h3>Delay</h3>
<ul>
  <li><a href="/delay/500">Sleep 500ms</a></li>
  <li><a href="/delay/1000">Sleep 1s</a></li>
  <li>..</li>
</ul>
<h3>Mixed</h3>
<ul>
  <li><a href="/mixed/">Accessing</a></li>
</ul>
</body></html>`)
}

func httpRequests(w http.ResponseWriter, r *http.Request) {
	var result string
	switch path := r.URL.Path[len("/http/"):]; path {
	case "200":
		result = httpGet(r, "http://ip-api.com/json")
	default:
		result = httpGet(r, "http://err/errors/"+path)
	}
	lib.RenderJSON(w, result, nil)
}

func database(w http.ResponseWriter, r *http.Request) {
	result := ""
	err := dbConn().QueryRow(r.Context(), "SELECT 1").Scan(&result)
	lib.RenderJSON(w, result, err)
}

func s3ListFiles(w http.ResponseWriter, r *http.Request) {
	bucket := r.URL.Path[len("/s3/"):]
	if len(bucket) == 0 {
		result, err := s3Client().ListBucketsWithContext(r.Context(), &s3.ListBucketsInput{})
		lib.RenderJSON(w, result, err)
		return
	}
	result, err := s3Client().ListObjectsWithContext(r.Context(), &s3.ListObjectsInput{
		Bucket: aws.String(bucket),
	})
	lib.RenderJSON(w, result, err)
}

func delay(w http.ResponseWriter, r *http.Request) {
	duration := r.URL.Path[len("/delay/"):]
	if len(duration) == 0 {
		duration = "100"
	}
	lib.RenderJSON(w, httpGet(r, "http://err/sleep/"+duration), nil)
}

func mixed(w http.ResponseWriter, r *http.Request) {
	// HTTP Requests
	httpResult := httpGet(r, "http://err/errors/429")

	// Access database
	mysqlVersion := ""
	row := dbConn().QueryRow(r.Context(), "SELECT version()")
	row.Scan(&mysqlVersion)

	// List s3 buckets
	buckets, err := s3Client().ListBucketsWithContext(r.Context(), &s3.ListBucketsInput{})
	lib.RenderJSON(w, fmt.Sprintf("%s, %s, %v", httpResult, mysqlVersion, buckets), err)
}

func httpGet(r *http.Request, uri string) string {
	resp, _ := ctxhttp.Get(r.Context(), xray.Client(nil), uri)
	if json, err := ioutil.ReadAll(resp.Body); err == nil {
		return string(json)
	}
	return ""
}

func s3Client() *s3.S3 {
	client := s3.New(session.Must(session.NewSession(&aws.Config{
		DisableSSL: aws.Bool(true),
	})))
	xray.AWS(client.Client)
	return client
}

func dbConn() *xray.DB {
	db, _ := xray.SQL("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		env("MYSQL_USER", "user"), env("MYSQL_PASSWORD", "pass"),
		env("MYSQL_HOST", "db"), env("MYSQL_PORT", "3306"),
		env("MYSQL_DATABASE", "")))
	return db
}
