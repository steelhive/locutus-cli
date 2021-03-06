package query

import (
    "os"
    "fmt"
    "github.com/aws/aws-sdk-go/service/ec2"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/aws/ec2metadata"
)


func New(session *session.Session) *Query {
    return &Query{
        Client: *ec2.New(session),
        Metadata: *ec2metadata.New(session),
    }
}


type Query struct {
    Client ec2.EC2
    Metadata ec2metadata.EC2Metadata
}


func (q *Query) GetSelf() ec2metadata.EC2InstanceIdentityDocument {
    res, err := q.Metadata.GetInstanceIdentityDocument()
    if err != nil {
        msg := fmt.Sprintf("Query Error: %s", err)
        fmt.Println(msg)
        os.Exit(2)
    }
    return res
}

func (q *Query) GetRole(key *string) *string {
    metadata := q.GetSelf()
    if *key == "" {
        *key = "locutus:role"
    }
    name := "instance-id"
    values := make([]*string, 1)
    values[0] = &metadata.InstanceID

    filters := []*ec2.Filter{
        &ec2.Filter{
            Name: &name,
            Values: values,
        },
    }
    input := &ec2.DescribeInstancesInput{
        Filters: filters,
    }

    unknown := ""
    instances := q.GetInstances(input)
    for _, instance := range instances {
        for _, tag := range instance.Tags {
            if *tag.Key == *key {
                return tag.Value
            }
        }
    }
    return &unknown
}

func (q *Query) GetInstances(input *ec2.DescribeInstancesInput) []ec2.Instance {
    res, err := q.Client.DescribeInstances(input)
    if err != nil {
        msg := fmt.Sprintf("Query Error: %s", err)
        fmt.Println(msg)
        os.Exit(2)
    }
    instances := make([]ec2.Instance, len(res.Reservations))
    for i, reservation := range res.Reservations {
        for _, instance := range reservation.Instances {
            instances[i] = *instance
        }
    }
    return instances
}


func (q *Query) GetPrivateIPs(key *string, values *[]string) []string {
    input := &ec2.DescribeInstancesInput{}

    if *key != "" && len(*values) > 0 {
        name := fmt.Sprintf("tag:%s", *key)

        tagValues := make([]*string, len(*values))

        for i, tag := range *values {
            tagValues[i] = &tag
        }

        filters := []*ec2.Filter{
            &ec2.Filter{
                Name: &name,
                Values: tagValues,
            },
        }

        input = &ec2.DescribeInstancesInput{
            Filters: filters,
        }
    }

    instances := q.GetInstances(input)
    ips := make([]string, len(instances))
    for i, instance := range instances {
        ips[i] = *instance.PrivateIpAddress
    }
    return ips
}
