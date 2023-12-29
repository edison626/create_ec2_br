package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type EC2Config struct {
	ImageId         string
	InstanceType    string
	KeyName         string
	SecurityGroupID string
	SubnetID        string
	TagValue        string
	VolumeSize      int64
}

type Asset struct {
	ID           string   `json:"id"`
	Hostname     string   `json:"hostname"`
	IP           string   `json:"ip"`
	Platform     string   `json:"platform"`
	Protocols    []string `json:"protocols"`
	Protocol     string   `json:"protocol"`
	Port         int      `json:"port"`
	IsActive     bool     `json:"is_active"`
	PublicIP     string   `json:"public_ip"`
	Number       string   `json:"number"`
	Comment      string   `json:"comment"`
	Vendor       string   `json:"vendor"`
	Model        string   `json:"model"`
	SN           string   `json:"sn"`
	CPUModel     string   `json:"cpu_model"`
	CPUCount     int      `json:"cpu_count"`
	CPUCores     int      `json:"cpu_cores"`
	CPUVcpus     int      `json:"cpu_vcpus"`
	Memory       string   `json:"memory"`
	DiskTotal    string   `json:"disk_total"`
	DiskInfo     string   `json:"disk_info"`
	OS           string   `json:"os"`
	OSVersion    string   `json:"os_version"`
	OSArch       string   `json:"os_arch"`
	HostnameRaw  string   `json:"hostname_raw"`
	Domain       string   `json:"domain"`
	AdminUser    string   `json:"admin_user"`
	Nodes        []string `json:"nodes"`
	NodesDisplay []string `json:"nodes_display"`
	Labels       []string `json:"labels"`
}

// 创建一个EC2Config切片包含不同的实例配置
func configEC2Instances(batch string, wukong []string) []EC2Config {
	return []EC2Config{
		//----------- 测试aws 类型服务器 --------------
		{
			ImageId:      "ami-01c1b05f1b398fb4b", //Amazon Linux
			InstanceType: "t3.small",
			TagValue:     batch + "br-testing-" + wukong[0],
			VolumeSize:   100,
		},
		{
			ImageId:      "ami-0070a028fc72c8f48", // rocky
			InstanceType: "t3.small",
			TagValue:     batch + "br-prod-web-proxy02-testing-" + wukong[1],
			VolumeSize:   100,
		},
		//-----------End 测试aws 类型服务器 --------------

		// {
		// 	ImageId:      "ami-01da42fa32830f2d0",
		// 	InstanceType: "t3.medium",
		// 	TagValue:     batch + "br-prod-web-proxy01",
		// 	VolumeSize:   100,
		// },
		// {
		// 	ImageId:      "ami-0070a028fc72c8f48", // rocky
		// 	InstanceType: "t3.small",
		// 	TagValue:     batch + "br-prod-web-proxy02",
		// 	VolumeSize:   100,
		// },
		// {
		// 	ImageId:      "ami-01da42fa32830f2d0",
		// 	InstanceType: "t3.small",
		// 	TagValue:     batch + "vn-prod-callback",
		// 	VolumeSize:   100,
		// },
		// {
		// 	ImageId:      "ami-01da42fa32830f2d0",
		// 	InstanceType: "t3.small",
		// 	TagValue:     batch + "vn-prod-houtai",
		// 	VolumeSize:   100,
		// },
		// {
		// 	ImageId:      "ami-0e8849aa060c28662",
		// 	InstanceType: "c5.xlarge",
		// 	TagValue:     batch + "vn-prod-mongodb-01",
		// 	VolumeSize:   500,
		// },
		// {
		// 	ImageId:      "ami-0e8849aa060c28662",
		// 	InstanceType: "c5.xlarge",
		// 	TagValue:     batch + "vn-prod-mongodb-02",
		// 	VolumeSize:   500,
		// },
		// {
		// 	ImageId:      "ami-0e8849aa060c28662",
		// 	InstanceType: "c5.xlarge",
		// 	TagValue:     batch + "vn-prod-mongodb-03",
		// 	VolumeSize:   500,
		// },
		// {
		// 	ImageId:      "ami-0e8849aa060c28662",
		// 	InstanceType: "t3.xlarge",
		// 	TagValue:     batch + "vn-prod-cgcron-clinet-task01",
		// 	VolumeSize:   300,
		// },
		// {
		// 	ImageId:      "ami-0e8849aa060c28662",
		// 	InstanceType: "c5.xlarge",
		// 	TagValue:     batch + "vn-prod-cgcron-clinet-task02",
		// 	VolumeSize:   300,
		// },
		// {
		// 	ImageId:      "ami-0e8849aa060c28662",
		// 	InstanceType: "c5.xlarge",
		// 	TagValue:     batch + "prod-redis-etcd01",
		// 	VolumeSize:   150,
		// },
		// {
		// 	ImageId:      "ami-0e8849aa060c28662",
		// 	InstanceType: "c5.xlarge",
		// 	TagValue:     batch + "prod-redis-etcd02",
		// 	VolumeSize:   150,
		// },
		// {
		// 	ImageId:      "ami-0e8849aa060c28662",
		// 	InstanceType: "c5.xlarge",
		// 	TagValue:     batch + "prod-redis-etcd03",
		// 	VolumeSize:   150,
		// },
		// {
		// 	ImageId:      "ami-0e8849aa060c28662",
		// 	InstanceType: "c5.2xlarge",
		// 	TagValue:     batch + "starrocks-be-01",
		// 	VolumeSize:   300,
		// },
		// {
		// 	ImageId:      "ami-0e8849aa060c28662",
		// 	InstanceType: "c5.2xlarge",
		// 	TagValue:     batch + "starrocks-be-02",
		// 	VolumeSize:   300,
		// },
		// {
		// 	ImageId:      "ami-0e8849aa060c28662",
		// 	InstanceType: "c5.2xlarge",
		// 	TagValue:     batch + "starrocks-be-03",
		// 	VolumeSize:   300,
		// },
		// {
		// 	ImageId:      "ami-0e8849aa060c28662",
		// 	InstanceType: "c5.xlarge",
		// 	TagValue:     batch + "starrocks-fe-01",
		// 	VolumeSize:   100,
		// },
		// {
		// 	ImageId:      "ami-0e8849aa060c28662",
		// 	InstanceType: "c5.xlarge",
		// 	TagValue:     batch + "starrocks-fe-02",
		// 	VolumeSize:   100,
		// },
		// {
		// 	ImageId:      "ami-0e8849aa060c28662",
		// 	InstanceType: "c5.xlarge",
		// 	TagValue:     batch + "starrocks-fe-03",
		// 	VolumeSize:   100,
		// },
		// {
		// 	ImageId:      "ami-0e8849aa060c28662",
		// 	InstanceType: "c5.xlarge",
		// 	TagValue:     batch + "wukong-yq-a06",
		// 	VolumeSize:   200,
		// },
		// {
		// 	ImageId:      "ami-0e8849aa060c28662",
		// 	InstanceType: "c5.xlarge",
		// 	TagValue:     batch + "wukong-yq-a07",
		// 	VolumeSize:   200,
		// },
		// {
		// 	ImageId:      "ami-0e8849aa060c28662",
		// 	InstanceType: "c5.xlarge",
		// 	TagValue:     batch + "wukong-yq-a08",
		// 	VolumeSize:   200,
		// },
		// {
		// 	ImageId:      "ami-0e8849aa060c28662",
		// 	InstanceType: "c5.2xlarge",
		// 	TagValue:     batch + "yq-prod-cloudcanal-01",
		// 	VolumeSize:   500,
		// },
		// {
		// 	ImageId:      "ami-0e8849aa060c28662",
		// 	InstanceType: "c5.2xlarge",
		// 	TagValue:     batch + "yq-prod-cloudcanal-02",
		// 	VolumeSize:   500,
		// },
		// {
		// 	ImageId:      "ami-0e8849aa060c28662",
		// 	InstanceType: "c5.2xlarge",
		// 	TagValue:     batch + "yq-prod-cloudcanal-03",
		// 	VolumeSize:   500,
		// },
		// {
		// 	ImageId:      "ami-0e8849aa060c28662",
		// 	InstanceType: "c5.xlarge",
		// 	TagValue:     batch + "yq-rocketmq-prod1",
		// 	VolumeSize:   100,
		// },
		// {
		// 	ImageId:      "ami-0e8849aa060c28662",
		// 	InstanceType: "c5.xlarge",
		// 	TagValue:     batch + "yq-rocketmq-prod2",
		// 	VolumeSize:   100,
		// },
		// {
		// 	ImageId:      "ami-0e8849aa060c28662",
		// 	InstanceType: "c5.xlarge",
		// 	TagValue:     batch + "vn-prod-zinc-beanstalkd",
		// 	VolumeSize:   100,
		// },
		// {
		// 	ImageId:      "ami-0e8849aa060c28662",
		// 	InstanceType: "c5.xlarge",
		// 	TagValue:     batch + "vn-prod-rmqtt01",
		// 	VolumeSize:   100,
		// },
		// {
		// 	ImageId:      "ami-0e8849aa060c28662",
		// 	InstanceType: "c5.xlarge",
		// 	TagValue:     batch + "vn-prod-rmqtt02",
		// 	VolumeSize:   100,
		// },
	}
}

// JumpServer 创建新资产的请求
func CreateNewAsset(jmsurl, token string, assetClietToken string, assetHostName string,
	assetIP string, varAssetNote string, varAssetNodeDisplay string) {
	// 创建资产数据
	newAsset := Asset{
		ID:           assetClietToken, //确认是否是UUID 是 ClientToken
		Hostname:     assetHostName,
		IP:           assetIP,
		Platform:     "Linux",
		Protocols:    []string{"ssh/22"},
		Protocol:     "ssh",
		Port:         22,
		IsActive:     true,
		PublicIP:     assetIP,
		AdminUser:    "463fb17d-1257-40ea-8dbd-ddae4ddae199",
		Nodes:        []string{varAssetNote},        // 修改目录 UUID
		NodesDisplay: []string{varAssetNodeDisplay}, // 修改目录 UUID
		Labels:       []string{},
		// 填写其他字段...
	}

	// 将资产数据转换为 JSON
	jsonData, err := json.Marshal(newAsset)
	if err != nil {
		log.Fatal(err)
	}

	// 构造 POST 请求
	url := jmsurl + "/api/v1/assets/assets/"
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
	}

	// 添加必要的头部
	req.Header.Add("Authorization", "Token "+token)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-JMS-ORG", "00000000-0000-0000-0000-000000000002")

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(body))
}

func main() {

	// 在jenkins 输入值
	varJmsServerURL := os.Getenv("JmsServerURL")
	varJMSToken := os.Getenv("JMSToken")
	varBatch := os.Getenv("Batch")
	varAssetNote := os.Getenv("AssetNote")
	varAssetNodeDisplay := os.Getenv("AssetNodeDisplay")
	varWukong := os.Getenv("WukongPlatform")
	fmt.Printf("JmsServerURL : %s\n", varJmsServerURL)
	fmt.Printf("JMSToken : %s\n", varJMSToken)
	fmt.Printf("Batch : %s\n", varBatch)
	fmt.Printf("AssetNote : %s\n", varAssetNote)
	fmt.Printf("AssetNodeDisplay : %s\n", varAssetNodeDisplay)
	fmt.Printf("WukongPlatform : %s\n", varWukong)

	if varJmsServerURL == "" || varJMSToken == "" || varBatch == "" || varAssetNote == "" || varAssetNodeDisplay == "" || varWukong == "" {
		log.Fatalf("值不能为空")
	}

	//创建AWS会话和EC2服务客户端代码
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("sa-east-1"), // 替换为您的AWS区域
	})
	if err != nil {
		fmt.Println("创建会话失败:", err)
		return
	}

	svc := ec2.New(sess)

	//实例内容 - 函数
	arrayWukong := strings.Split(varWukong, ",")
	configs := configEC2Instances(varBatch, arrayWukong)

	// ... 创建实例的代码
	for _, config := range configs {
		runResult, err := svc.RunInstances(&ec2.RunInstancesInput{
			ImageId:      aws.String(config.ImageId),
			InstanceType: aws.String(config.InstanceType),
			KeyName:      aws.String("ec2-br"),
			MinCount:     aws.Int64(1),
			MaxCount:     aws.Int64(1), // 只创建一台实例
			SecurityGroupIds: []*string{
				aws.String("sg-04ef918475fdf4f6c"),
			},
			SubnetId: aws.String("subnet-07a05d60e8bfb62be"), // 替换为您的子网ID
			BlockDeviceMappings: []*ec2.BlockDeviceMapping{
				{
					DeviceName: aws.String("/dev/sdh"),
					Ebs: &ec2.EbsBlockDevice{
						VolumeSize: aws.Int64(config.VolumeSize), // 存储
						VolumeType: aws.String("gp2"),
					},
				},
			},
			TagSpecifications: []*ec2.TagSpecification{
				{
					ResourceType: aws.String("instance"),
					Tags: []*ec2.Tag{
						{
							Key:   aws.String("Name"),
							Value: aws.String(config.TagValue), // 指定实例名称
						},
					},
				},
			},
		})

		if err != nil {
			fmt.Println("无法创建实例:", err)
			return
		}

		fmt.Println("已成功创建实例:", runResult.Instances)

		instanceId := runResult.Instances[0].InstanceId

		// 等待实例变为running状态
		fmt.Println("等待实例启动...")
		for {
			descInstances, err := svc.DescribeInstances(&ec2.DescribeInstancesInput{
				InstanceIds: []*string{instanceId},
			})
			if err != nil {
				fmt.Println("无法获取实例状态:", err)
				return
			}

			state := descInstances.Reservations[0].Instances[0].State.Name
			if *state == "running" {
				break
			}

			time.Sleep(10 * time.Second)
		}
		fmt.Println("实例已启动,正在分配弹性IP...")

		// 申请弹性IP
		allocRes, err := svc.AllocateAddress(&ec2.AllocateAddressInput{
			Domain: aws.String("vpc-0cadb665c480c21d1"), // VPC网络
		})
		if err != nil {
			fmt.Println("无法分配弹性IP:", err)
			return
		}

		// 关联弹性IP到实例
		_, err = svc.AssociateAddress(&ec2.AssociateAddressInput{
			InstanceId:   instanceId,
			AllocationId: allocRes.AllocationId,
		})
		if err != nil {
			fmt.Println("无法关联弹性IP:", err)
			return
		}
		fmt.Println("弹性IP已成功关联到实例:", *instanceId)

		// 获取弹性IP的详细信息
		describeAddressesOutput, err := svc.DescribeAddresses(&ec2.DescribeAddressesInput{
			AllocationIds: []*string{allocRes.AllocationId},
		})
		if err != nil {
			fmt.Println("无法获取弹性IP的详细信息:", err)
			return
		}

		// 检查是否有返回的地址
		if len(describeAddressesOutput.Addresses) > 0 {
			eip := describeAddressesOutput.Addresses[0].PublicIp
			fmt.Println("关联的弹性IP地址是:", *eip)
		} else {
			fmt.Println("未找到弹性IP的详细信息")
		}

		fmt.Println("~~~~~~~值并配置Jumpserver API~~~~~~~~~~~~")
		var assetInstanceName string
		for _, tag := range runResult.Instances[0].Tags {
			if *tag.Key == "Name" {
				assetInstanceName = *tag.Value
				break
			}
		}

		assetsClientToken := runResult.Instances[0].ClientToken
		assetIP := describeAddressesOutput.Addresses[0].PublicIp

		fmt.Println("ClientToken:", assetsClientToken)
		fmt.Println("Host Name:", assetInstanceName)
		fmt.Println("Host IP:", assetIP)

		CreateNewAsset(varJmsServerURL, varJMSToken, *assetsClientToken, assetInstanceName, *assetIP, varAssetNote, varAssetNodeDisplay)

	}
}
