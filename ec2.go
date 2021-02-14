package main
import (
    "fmt"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
     "github.com/aws/aws-sdk-go/service/ec2"
"github.com/aws/aws-sdk-go/aws/awserr"
// "github.com/aws/aws-sdk-go/service/elbv2"
)

func main(){
sess := session.Must(session.NewSessionWithOptions(session.Options{
    Profile: "xxx",
    Config: aws.Config{
        Region: aws.String("xxx"),
    },}))
svc := ec2.New(sess)
input := &ec2.DescribeInstancesInput{}

result, err := svc.DescribeInstances(input)
if err != nil {
    if aerr, ok := err.(awserr.Error); ok {
        switch aerr.Code() {
        default:
            fmt.Println(aerr.Error())
        }
    } else {
        // Print the error, cast err to awserr.Error to get the Code and
        // Message from an error.
        fmt.Println(err.Error())
    }
    return
}
var tagmiss[] string
var inst[] string
for _,out := range result.Reservations {
	inst = append(inst,*out.Instances[0].InstanceId)
}

// all the ec2 instance ids will be appended to the inst slice

for _,out := range inst {
input := &ec2.DescribeTagsInput{
    Filters: []*ec2.Filter{
        {
            Name: aws.String("resource-id"),
            Values: []*string{
                aws.String(out),
            },
        },
    },
}
result, err := svc.DescribeTags(input)
if err != nil {
    if aerr, ok := err.(awserr.Error); ok {
        switch aerr.Code() {
        default:
            fmt.Println(aerr.Error())
        }
    } else {
        // Print the error, cast err to awserr.Error to get the Code and
        // Message from an error.
        fmt.Println(err.Error())
    }
    return
}
flag := 0
for _,ta := range result.Tags {
	if *ta.Key == "PROJECT" {
		flag = 1 }
	// fmt.Println(ta)
	// fmt.Println("---")
}
if flag == 0 {
tagmiss = append(tagmiss, out)
} 
}
//fmt.Println("Tag is missing for the below instances :")
if len(tagmiss) > 0 {
fmt.Println("Tag is missing for the below instances :")

	for _, i := range tagmiss {
		fmt.Println(i)
}} else {
fmt.Println("All instances are tagged")
}
}
