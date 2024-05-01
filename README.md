# ECS Drainer (Lambda)

This project is an extension of the ideas presented in the [AWS Blog post](https://aws.amazon.com/ru/blogs/compute/how-to-automate-container-instance-draining-in-amazon-ecs/) and the [GitHub sample](https://github.com/aws-samples/ecs-cid-sample), with key differences:

- Receives AutoScaling Hooks events through CloudWatch rules, supporting multiple ECS Clusters with a single function.
- Utilizes the [Serverless Framework](https://github.com/serverless/serverless) for deployment.
- Implemented in Golang.
- Supports draining Spot Instances via [Spot Instance Interruption Notice](https://docs.aws.amazon.com/en_us/AWSEC2/latest/UserGuide/spot-interruptions.html#spot-instance-termination-notices).

## Motivation

During ECS instance AMI updates, the Auto Scaling Group (ASG) may replace instances without draining them, potentially causing brief container downtimes. This function automates the draining process for ECS cluster instances, enhancing availability.

## How It Works

The *ecs-drain-lambda* function:

1. Receives a CloudWatch event:
    - An **ANY** AutoScaling Lifecycle Terminate event configured via [EC2 Auto Scaling Lifecycle Hooks](https://docs.aws.amazon.com/autoscaling/ec2/userguide/lifecycle-hooks.html) for the `autoscaling:EC2_INSTANCE_TERMINATING` event.
    - An **ANY** Spot Instance Interruption Notice (Note: AWS does not guarantee that instances will be drained in time; instances could be terminated before notice arrival).
2. Retrieves the ID of the terminating instance.
3. Extracts the ECS Cluster name from the instance's UserData (`ECS_CLUSTER=xxxxxxxxx` format).
4. Initiates the draining process if ECS Tasks are running on the instance.
5. Waits for all ECS Tasks to shutdown.
6. Completes the Lifecycle Hook, allowing the ASG to proceed with instance termination.

## Requirements

- [Serverless Framework](https://github.com/serverless/serverless)
- [Golang](https://golang.org/doc/install)
- GNU Make
- Configured EC2 Auto Scaling Lifecycle Hooks for the `autoscaling:EC2_INSTANCE_TERMINATING` event.

### CloudFormation Example for Lifecycle Hook

```yaml
ASGTerminateHook:
  Type: "AWS::AutoScaling::LifecycleHook"
  Properties:
    AutoScalingGroupName: !Ref ECSAutoScalingGroup
    DefaultResult: "ABANDON"
    HeartbeatTimeout: "900"
    LifecycleTransition: "autoscaling:EC2_INSTANCE_TERMINATING"
```

## Usage

```

git clone github.com/moabukar/ecs-drainer.git
cd ecs-drainer
make deploy
# To specify a different AWS region, use:
sls deploy -v --region 
```

## Deploy with Terraform
For Terraform deployment, refer to the ecs-drain-lambda Terraform Module.

## Limitations

- The function waits up to 15 minutes to complete the Drain process. If it exceeds this time, it times out.

- Failure of the function triggers the default lifecycle hook action (ABANDON or CONTINUE), both of which will allow the instance to terminate. ABANDON will halt any remaining actions, such as other lifecycle hooks, while CONTINUE will allow them to complete.
