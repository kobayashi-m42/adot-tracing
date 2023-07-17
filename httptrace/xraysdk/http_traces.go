package xraysdk

import (
	"log"
	"net/http"

	"github.com/aws/aws-xray-sdk-go/xray"
)

func AwsSdkCall(w http.ResponseWriter, r *http.Request, s3 *S3Client) {
	if _, err := s3.Client.ListBuckets(r.Context(), nil); err != nil {
		log.Println(err)
		return
	}
	_, _ = w.Write([]byte("Tracing aws sdk call"))
}

func OutgoingHttpCall(w http.ResponseWriter, r *http.Request) {
	req, _ := http.NewRequestWithContext(r.Context(), "GET", "https://aws.amazon.com/", nil)

	response, err := xray.Client(nil).Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer response.Body.Close()
	_, _ = w.Write([]byte("Tracing http call"))
}
