package lichterfest_lambda

import (
    "github.com/aws/aws-lambda-go/lambda"
    "net/http"
    "io/ioutil"
    "log"
    "github.com/Iwark/pushnotification"
    "github.com/aws/aws-sdk-go/aws"
    "os"
    "context"
)

//noinspection GoUnusedFunction
func main() {
    lambda.Start(lichterfestNotify)
}

func lichterfestChanged() (bool) {
    resp, err := http.Get("https://www.tickets-lichterfest.de/")
    if err != nil {
        log.Fatal(err)
    }
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Fatal(err)
    }
    sbody := string(body[:len(expectedBody)])
    return sbody != expectedBody
}

//noinspection GoUnusedParameter
func lichterfestNotify(ctx context.Context, name string) bool {
    if !lichterfestChanged() {
        return false
    }

    push := pushnotification.Service{
        AWSAccessKey:         os.Getenv("AWS_ACCESS_KEY"),
        AWSAccessSecret:      os.Getenv("AWS_ACCESS_SECRET"),
        AWSSNSApplicationARN: os.Getenv("AWS_SNS_APPLICATION_ARN"),
        AWSRegion:            os.Getenv("AWS_REGION"),
    }
    err := push.Send(os.Getenv("DEVICE_TOKEN"), &pushnotification.Data{
        Alert: aws.String("test message"),
        Sound: aws.String("default"),
        Badge: aws.Int(1),
    })
    if err != nil {
        log.Fatal(err)
    }
    return true
}

const expectedBody = `<html>
<head>
    <link href="https://assets.jimstatic.com/under-construction.css.fb4d2e84f40565710387a60d91649d42.css" rel="stylesheet" type="text/css" media="all"/>
    <meta content="noindex" name="robots">
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <title>https://www.tickets-lichterfest.de/</title>
</head>
<body>
<div class="background"></div>

<section class="pp-content-wrapper">
    <div class="pp-content-wrapper-inner">
        <div class="pp-content">
            <h1 class="pp-headline">
                Hallo!                    </h1>
            <p class="pp-description">
                Diese Webseite befindet sich gerade im Aufbau und ist bald online.                    </p>

        </div>
    </div>
</section>

<div id="contact-overlay">
    <div class="pp-contactform-align">
        <div id="contact-form" class="pp-contactform-wrapper">

            <a id="contact-form-close" title="Close this element" href="#">
                Schließen                    </a>

            <form action="" method="POST" id="contact" class="pp-contactform">
                <fieldset class="pp-contactform-group">
                    <legend class="pp-contactform-title">
                        Kontakt                            </legend>
                    <div class="pp-contactform-item">
                        <label class="pp-contactform-label">Name</label>
                        <input class="pp-contactform-input" type="text" id="contact-name" name="contact-name" placeholder="Name" />
                    </div>
                    <div class="pp-contactform-item">
                        <label class="pp-contactform-label">E-Mail-Adresse</label>
                        <input class="pp-contactform-input" type="email" id="contact-mail" name="contact-mail" placeholder="E-Mail" />
                    </div>
                    <div class="pp-contactform-item">
                        <label class="pp-contactform-label">Nachricht</label>
                        <textarea class="pp-contactform-textarea" id="contact-text" name="contact-text"></textarea>
                    </div>

                    <div class="pp-contactform-item-error pp-contactform-item-server-error">
                        <div class="jui-message-box jui-message-box--danger pp-contactform-error pp-contactform-item-server-error">
                            <span class="jui-message-box__icon icon-warning-sign"></span>
                            <div class="jui-message-box__content">
                                Ups, da ist etwas schief gelaufen!
                            </div>
                        </div>
                    </div>

                    <div class="pp-contactform-item-error pp-contactform-item-captcha-error">
                        <div class="jui-message-box jui-message-box--danger ">
                            <span class="jui-message-box__icon icon-warning-sign"></span>
                            <div class="jui-message-box__content">
                                Captcha ist nicht korrekt.                                    </div>
                        </div>
                    </div>
                    <div class="pp-contactform-item-success">
                        <div class="jui-success">
                                  <span class="jui-success__icon icon-ok">

                                  </span>
                        </div>
                        <span class="pp-contactform-thanks">Danke!</span>
                    </div>

                    <div class="pp-contactform-item pp-contactform-item-submit">
                        <button class='jui-button jui-button--primary jui-button--big' href='#' name="submit">
                                    <span class="jui-button__text">
                                        Senden                                    </span>
                        </button>
                    </div>
                </fieldset>
            </form>
        </div>

    </div>
</div>

<footer class="pp-footer">
</footer>

<script src="https://assets.jimstatic.com/js/init.js.2d641d170fd679ccb2fe.js"></script>
<script src="https://assets.jimstatic.com/under-construction.js.ca19598918a2cf7724d3.js"></script>
</body>
</html>

`