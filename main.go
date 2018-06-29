package main

import (
    "github.com/aws/aws-lambda-go/lambda"
    "net/http"
    "io/ioutil"
    "log"
    "context"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/sns"
    "github.com/aws/aws-sdk-go/aws"
    "bytes"
    "os"
)

//noinspection GoUnusedFunction
func main() {
    lambda.Start(lichterfestNotify)
}

func lichterfestUnchanged() (bool) {
    resp, err := http.Get(os.Getenv("Url"))
    if err != nil {
        log.Fatal(err)
    }
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Fatal(err)
    }
    return bytes.Contains(body, []byte(os.Getenv("UrlContains")))
}

//noinspection GoUnusedParameter
func lichterfestNotify(ctx context.Context, name string) (message string, err error) {
    if lichterfestUnchanged() {
        return "not_changed", nil
    }

    svc := sns.New(session.New())
    params := &sns.PublishInput{
        Message: aws.String(os.Getenv("NotificationMessage")), // This is the message itself (can be XML / JSON / Text - anything you want)
        TopicArn: aws.String(os.Getenv("AwsSnsApplicationArn")),  //Get this from the Topic in the AWS console.
    }
    _, err = svc.Publish(params)

    if err != nil {                    //Check for errors
        return "", err
    }

    // Pretty-print the response data.
    return "changed", nil
}