package main

import (
    "fmt"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/awserr"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/ec2"
    "github.com/aws/aws-sdk-go/service/elb"
    "github.com/aws/aws-sdk-go/service/autoscaling"
    "encoding/json"
)

func main() {
    sess := session.Must(session.NewSessionWithOptions(session.Options{
        Profile: "xxx",
        Config: aws.Config{
            Region: aws.String("xxx"),
        }}))

    svc := ec2.New(sess)
    svc2 := elb.New(sess)
    svc3 := autoscaling.New(sess)

    fmt.Println(`Enter any of the following number: 
    1 > To fing the untagged images
    2 > To find the untagged EC2 instances
    3 > To find the untagged snapshots
    4 > To find the untagged volumes
    5 > To find the untagged loadbalancers
    6 > To find the untagged EIPs
    7 > To find the untagged networkinterfaces
    8 > To find the untagged ASGs`)

    var nu int
    fmt.Scanln(&nu)
    fmt.Println("Fetching the results...! Wait !")
    switch {
    case nu == 1:
        fmt.Println(imag(svc))
    case nu == 2:
        fmt.Println(ec(svc))
    case nu == 3:
        fmt.Println(snap(svc))
    case nu == 4:
        fmt.Println(vo(svc))
    case nu == 5:
        fmt.Println(loadb(svc2))
    case nu == 6:
        fmt.Println(eipfunc(svc))
    case nu == 7:
        fmt.Println(net(svc))
    case nu == 8:
        fmt.Println(asg(svc3))
    default:
        fmt.Println("Invalid")
    }

}

// The imag function is used to find all the untagged images.

func imag(svc *ec2.EC2) interface{} {

    input := &ec2.DescribeImagesInput{
        Owners: []*string{
            aws.String("self"),
        },
    }

    result, err := svc.DescribeImages(input)
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
        //   return
    }
    var tag_ima []string
    var ima []string
    for _, im := range result.Images {
        ima = append(ima, *im.ImageId)

    }
    for _, ima_out := range ima {
        input := &ec2.DescribeTagsInput{
            Filters: []*ec2.Filter{
                {
                    Name: aws.String("resource-id"),
                    Values: []*string{
                        aws.String(ima_out),
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
            //   return
        }
        flag := 0
        for _, ta := range result.Tags {
            if *ta.Key == "PROJECT" {
                flag = 1
            }
        }
        if flag == 0 {
            tag_ima = append(tag_ima, ima_out)
        }
    }
    if len(tag_ima) > 0 {
        fmt.Println("Tag is missing for the below images :")
        imaJson, _ := json.Marshal(tag_ima)
        return string(imaJson)

    } else {
        out, _ := json.Marshal("All images are tagged")
        return string(out)
    }
}

// The ec function is used to find all the untagged ec2 instances.

func ec(svc *ec2.EC2) interface{} {
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
        //   return
    }
    var tagmiss []string
    var inst []string
    for _, out := range result.Reservations {
        inst = append(inst, *out.Instances[0].InstanceId)
    }

    // all the ec2 instance ids will be appended to the inst slice

    for _, out := range inst {
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
            //   return
        }
        flag := 0
        for _, ta := range result.Tags {
            if *ta.Key == "PROJECT" {
                flag = 1
            }
            // fmt.Println(ta)
            // fmt.Println("---")
        }
        if flag == 0 {
            tagmiss = append(tagmiss, out)
        }
    }
    if len(tagmiss) > 0 {
        fmt.Println("Tag is missing for the below instances :")
        ecJson, _ := json.Marshal(tagmiss)
        return string(ecJson)
    } else {
        out, _ := json.Marshal("All instances are tagged")
        return string(out)
    }
}


// The snap function is used to find all the untagged snapshots.

func snap(svc *ec2.EC2) interface{} {
    input := &ec2.DescribeSnapshotsInput{OwnerIds: []*string{
        aws.String("self"),
    }}
    result, err := svc.DescribeSnapshots(input)
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
        //   return
    }
    var tag_snap []string
    var snap []string
    for _, sn := range result.Snapshots {
        snap = append(snap, *sn.SnapshotId)
    }
    for _, snap_out := range snap {
        input := &ec2.DescribeTagsInput{
            Filters: []*ec2.Filter{
                {
                    Name: aws.String("resource-id"),
                    Values: []*string{
                        aws.String(snap_out),
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
            //  return
        }
        flag := 0
        for _, ta := range result.Tags {
            if *ta.Key == "PROJECT" {
                flag = 1
            }
        }
        if flag == 0 {
            tag_snap = append(tag_snap, snap_out)
        }
    }
    if len(tag_snap) > 0 {
        fmt.Println("Tag is missing for the below snapshots :")
        snapJson, _ := json.Marshal(tag_snap)
        return string(snapJson)
    } else {
        out, _ := json.Marshal("All snapshots are tagged")
        return string(out)
    }
}

// The vo function is used to find all the untagged volumes

func vo(svc *ec2.EC2) interface{} {
    input := &ec2.DescribeVolumesInput{}

    result, err := svc.DescribeVolumes(input)
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
        //  return
    }
    var vol []string
    var vol_tag []string
    for _, vol_out := range result.Volumes {
        if vol_out.Attachments != nil {
            vol = append(vol, *vol_out.Attachments[0].VolumeId)
        }
    }

    for _, vol_out := range vol {
        input := &ec2.DescribeTagsInput{
            Filters: []*ec2.Filter{
                {
                    Name: aws.String("resource-id"),
                    Values: []*string{
                        aws.String(vol_out),
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
            // return
        }

        flag := 0
        for _, ta := range result.Tags {
            if *ta.Key == "PROJECT" {
                flag = 1
            }
        }

        if flag == 0 {
            vol_tag = append(vol_tag, vol_out)
        }
    }
    if len(vol_tag) > 0 {
        fmt.Println("Tag is missing for the below volumes :")
        volJson, _ := json.Marshal(vol_tag)
        return string(volJson)
    } else {
        out, _ := json.Marshal("All volumes are tagged")
        return string(out)
    }
}

// The loadb function is used to find all the untagged loadbalancers
func loadb(svc *elb.ELB) interface{} {
    input := &elb.DescribeLoadBalancersInput{}

    result, err := svc.DescribeLoadBalancers(input)
    if err != nil {
        if aerr, ok := err.(awserr.Error); ok {
            switch aerr.Code() {
            case elb.ErrCodeAccessPointNotFoundException:
                fmt.Println(elb.ErrCodeAccessPointNotFoundException, aerr.Error())
            default:
                fmt.Println(aerr.Error())
            }
        } else {
            // Print the error, cast err to awserr.Error to get the Code and
            // Message from an error.
            fmt.Println(err.Error())
        }
    }

    // saving all the load balancer names to slice lb

    var lb []string
    var tagmiss []string
    for _, final := range result.LoadBalancerDescriptions {
        lb = append(lb, *final.LoadBalancerName)
    }

    // Iterating through slice lb to get the tag details of each lb. Complete tag details will be saved in variable a for individual lb

    for _, a := range lb {
        input := &elb.DescribeTagsInput{
            LoadBalancerNames: []*string{
                aws.String(a),
            },
        }
        result, err := svc.DescribeTags(input)
        if err != nil {
            if aerr, ok := err.(awserr.Error); ok {
                switch aerr.Code() {
                case elb.ErrCodeAccessPointNotFoundException:
                    fmt.Println(elb.ErrCodeAccessPointNotFoundException, aerr.Error())
                default:
                    fmt.Println(aerr.Error())
                }
            } else {
                // Print the error, cast err to awserr.Error to get the Code and
                // Message from an error.
                fmt.Println(err.Error())
            }
        }

        // checking project tag is available for the elb, if available, flag will be set to 1
        flag := 0
        for _, ta := range result.TagDescriptions[0].Tags {
            if *ta.Key == "PROJECT" {
                flag = 1
            }
        }

        // flag zero means no project tag and lb name will be saved to tagmiss slice
        if flag == 0 {

            tagmiss = append(tagmiss, a)
        }
    }

    // checking the length of slice tagmiss. If size greater than zero means tags missing for some lb. Iterate through tagmiss slice and print the name of each lb
    if len(tagmiss) > 0 {
        fmt.Println("Tag is missing for the below elbs :")
        lbJson, _ := json.Marshal(tagmiss)
        return string(lbJson)
    } else {
        out, _ := json.Marshal("All elbs are tagged")
        return string(out)
    }
}

// The eipfunc function is used to find all the untagged eips

func eipfunc(svc *ec2.EC2) interface{}  {
input := &ec2.DescribeAddressesInput{}

result, err := svc.DescribeAddresses(input)
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
//    return
}
var tag_eip []string
var ip []string
for _,eip := range result.Addresses {
	ip = append(ip,*eip.AllocationId)}

for _,out := range ip {
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
  //  return
}
flag := 0
for _,ta := range result.Tags {
	if *ta.Key == "PROJECT" {
		flag = 1 }
	// fmt.Println(ta)
	// fmt.Println("---")
}
if flag == 0 {
tag_eip = append(tag_eip, out)
} 
}
//fmt.Println("Tag is missing for the below instances :")
if len(tag_eip) > 0 {
fmt.Println("Tag is missing for the below eips :")
eipJson, _ := json.Marshal(tag_eip)
        return string(eipJson)
 } else {
out, _ := json.Marshal("All eips are tagged")
        return string(out)
}
}

// The net function is used to find all the untagged networkinterfaces

func net(svc *ec2.EC2) interface{}  {
input := &ec2.DescribeNetworkInterfacesInput{}

result, err := svc.DescribeNetworkInterfaces(input)
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
    // return
}
var netw []string
var tag_nw []string
for _, nw := range result.NetworkInterfaces {
	netw = append(netw,*nw.NetworkInterfaceId)
}
for _,out := range netw {
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
    //return
}
flag := 0
for _,ta := range result.Tags {
	if *ta.Key == "PROJECT" {
		flag = 1 }
	// fmt.Println(ta)
	// fmt.Println("---")
}
if flag == 0 {
tag_nw = append(tag_nw, out)
} 
}
//fmt.Println("Tag is missing for the below instances :")
if len(tag_nw) > 0 {
fmt.Println("Tag is missing for the below instances :")
        nwJson, _ := json.Marshal(tag_nw)
	    return string(nwJson) } else {
            return string("All instances are tagged")
}
}

// The asg function is used to find all the untagged AutoScalingGroups

func asg(svc *autoscaling.AutoScaling) interface{} {
input := &autoscaling.DescribeAutoScalingGroupsInput{}
result, err := svc.DescribeAutoScalingGroups(input)
if err != nil {
    if aerr, ok := err.(awserr.Error); ok {
        switch aerr.Code() {
        case autoscaling.ErrCodeInvalidNextToken:
            fmt.Println(autoscaling.ErrCodeInvalidNextToken, aerr.Error())
        case autoscaling.ErrCodeResourceContentionFault:
            fmt.Println(autoscaling.ErrCodeResourceContentionFault, aerr.Error())
        default:
            fmt.Println(aerr.Error())
        }
    } else {
        // Print the error, cast err to awserr.Error to get the Code and
        // Message from an error.
        fmt.Println(err.Error())
    }
    //return
}
var tagmiss []string
var asg []string

for _,out := range result.AutoScalingGroups {
//fmt.Println(*result.AutoScalingGroups[0].AutoScalingGroupName)
asg = append(asg,*out.AutoScalingGroupName)
}
for _,asg_out := range asg {
input := &autoscaling.DescribeTagsInput{
    Filters: []*autoscaling.Filter{
        {
            Name: aws.String("auto-scaling-group"),
            Values: []*string{
                aws.String(asg_out),
            },
        },
    },
}
result, err := svc.DescribeTags(input)
if err != nil {
    if aerr, ok := err.(awserr.Error); ok {
        switch aerr.Code() {
        case autoscaling.ErrCodeInvalidNextToken:
            fmt.Println(autoscaling.ErrCodeInvalidNextToken, aerr.Error())
        case autoscaling.ErrCodeResourceContentionFault:
            fmt.Println(autoscaling.ErrCodeResourceContentionFault, aerr.Error())
        default:
            fmt.Println(aerr.Error())
        }
    } else {
        // Print the error, cast err to awserr.Error to get the Code and
        // Message from an error.
        fmt.Println(err.Error())
    }
   // return
}
flag := 0
for _,ta := range result.Tags {
	if *ta.Key == "PROJECT" {
		flag = 1 }
}

if flag == 0 {

tagmiss = append(tagmiss, asg_out)
} 
}

if len(tagmiss) > 0 {
fmt.Println("Tag is missing for the below ASGs :")
	asgJson, _ := json.Marshal(tagmiss)
	return string(asgJson)
} else {
out, _ := json.Marshal("All ASGs are tagged")
 return string(out)
}
}

