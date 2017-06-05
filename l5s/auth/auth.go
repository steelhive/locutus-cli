package auth

import (
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
)

func GetSession(profile string, region string) *session.Session {
    if profile == "" {
        conf := aws.NewConfig().WithRegion(region)
        return session.Must(session.NewSession(conf))
    }
    opts := session.Options{
        Profile: profile,
        SharedConfigState: session.SharedConfigEnable,
    }
    return session.Must(session.NewSessionWithOptions(opts))
}
