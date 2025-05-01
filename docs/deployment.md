# AWS Deployment Guide

## Prerequisites

Prior to deployment, ensure the following components are prepared:

- AWS CLI configured with appropriate IAM credentials and permissions
- Docker runtime environment for container image building and testing
- Registered domain name with appropriate DNS configuration (for public-facing services)

## Architectural Overview

Our deployment architecture implements a multi-layered security model following the defense-in-depth principle:

1. Service workloads isolated in private subnets without direct internet connectivity
2. Controlled and auditable internet egress via NAT gateways
3. Service-to-service communication secured through IAM policies and security groups
4. Private AWS service access via VPC endpoints implementing PrivateLink
5. External traffic ingress controlled exclusively through API Gateway
6. Service discovery implemented via AWS Cloud Map
7. Container orchestration provided by Amazon ECS

## Virtual Private Cloud (VPC)

### Foundational Concepts

A VPC provides a logically isolated network segment within AWS where cloud resources are provisioned. It forms the network foundation for the entire application infrastructure.

Proper VPC configuration is essential for ensuring both security and functionality. Resources are deployed within private subnets (lacking direct ingress/egress internet connectivity) with controlled internet access facilitated through Network Address Translation (NAT) gateways when required.

### Subnet Architecture

An enterprise-grade VPC implementation should include:

1. **Public subnets**: Host infrastructure components requiring direct internet connectivity

   - Require an explicit route table entry to an Internet Gateway
   - Resources receive public IP address assignment (either static or dynamic)

2. **Private subnets**: Host application workloads, databases, and internal services
   - No direct internet routes defined
   - Enhanced security through network isolation
   - Distributed across multiple Availability Zones for high availability and fault tolerance

For production deployments, implement the following configuration:

1. **Subnet Distribution**:

   - Minimum of one public and one private subnet per Availability Zone
   - Deploy across at least two Availability Zones for fault tolerance
   - Standard configuration: 2 AZs = 2 public + 2 private subnets (minimum)

2. **Resource Placement Strategy**:
   - Public subnets: Reserved exclusively for infrastructure requiring direct internet connectivity (load balancers, NAT gateways, bastion hosts)
   - Private subnets: Application servers, databases, and internal service components

Security best practice: Deploy the majority of resources in private subnets, placing components in public subnets only when direct internet connectivity is an absolute requirement.

### Route Tables

Route tables define the network traffic flow between subnets and external networks, providing critical network segmentation and traffic control.

A route table specifies the outbound network traffic paths from a subnet or gateway. Key routing targets include:

- **Internet Gateways**: Enable direct internet connectivity for resources with public IP addresses
- **NAT Gateways**: Allow outbound-only internet connectivity for resources in private subnets while blocking inbound connection attempts

**Critical configuration note**: Each subnet requires explicit route table association. Without specified association, subnets default to the main route table, potentially creating unintended routing paths.

### Internet Access for Private Resources

The traffic flow for internet access from a private subnet follows this path:

1. Private subnet resource initiates an outbound connection request to an internet destination
2. Private subnet route table evaluates the destination IP against defined routes
3. The default route (0.0.0.0/0) matches as a catch-all for external destinations
4. Traffic is forwarded to the designated NAT gateway in a public subnet
5. NAT gateway performs source network address translation and forwards the request using its public IP
6. Public subnet route table routes traffic via its internet gateway route (0.0.0.0/0)
7. Internet gateway performs the necessary translation between VPC and internet routing
8. Response traffic traverses the reverse path back to the originating resource

Our architecture leverages this pattern to enable OpenTelemetry telemetry export to Grafana Cloud for observability.

### VPC Endpoints

VPC endpoints provide direct, private connectivity to AWS services without traversing the public internet, enhancing both security and performance.

Two distinct endpoint types are available:

1. **Gateway endpoints**: Implemented for S3 and DynamoDB services

   - No additional cost
   - Implemented as route table entries
   - Recommended for CloudWatch logs storage in S3

2. **Interface endpoints (AWS PrivateLink)**: Required for all other AWS services
   - Incurs endpoint provisioning and data transfer costs
   - Implemented as Elastic Network Interfaces (ENIs) with private IP addresses
   - Required for [ECR connectivity](https://docs.aws.amazon.com/AmazonECR/latest/userguide/vpc-endpoints.html#ecr-setting-up-vpc-create) from private subnets

**Critical implementation note**: ECR implementation requires multiple endpoint configurations:

- com.amazonaws.[region].ecr.api - API endpoint for ECR operations
- com.amazonaws.[region].ecr.dkr - Docker registry endpoint for image operations
- com.amazonaws.[region].s3 - Gateway endpoint for image layer storage

### VPC Implementation Checklist

- [ ] Create VPC with appropriate CIDR allocation (e.g., 10.0.0.0/16)
- [ ] Provision public and private subnets across multiple Availability Zones
- [ ] Configure and attach Internet Gateway to the VPC
- [ ] Deploy NAT Gateways in each public subnet for availability
- [ ] Define and associate appropriate route tables for all subnets
- [ ] Implement VPC endpoints for required AWS services

## Relational Database Service (RDS)

### Security Implementation

When provisioning RDS instances:

1. Set [Connectivity] > [Public access] to "No"
   - Restrict database access to resources within the VPC boundary
2. Deploy across multiple Availability Zones using Multi-AZ configuration
3. Create dedicated subnet groups spanning multiple AZs for database instance placement

## Elastic Container Registry (ECR)

ECR provides a managed container image registry for storing, versioning, and distributing container images within the AWS ecosystem.

## Elastic Container Service (ECS)

ECS delivers container orchestration capabilities for deploying, managing, and scaling containerized applications.

### Cluster Architecture

A cluster represents a logical isolation boundary for container workloads, serving as the foundation for service deployments.

- Implement separate clusters for each environment (development, staging, production)
- Select appropriate compute platform:
  - EC2: Provides granular control and cost optimization for predictable workloads
  - Fargate: Delivers serverless container management for variable workloads and reduced operational overhead

### Service Definition

A service manages the deployment lifecycle and availability of task instances, handling scheduling, placement, and recovery.

Key service configuration parameters:

- Task definition reference
- Desired task count
- Deployment strategy configuration (rolling deployments, blue/green with CodeDeploy)
- Network configuration (VPC, subnet selection)
- Service discovery integration with Cloud Map

**Implementation consideration**: Configure appropriate health check parameters and deployment circuit breakers to prevent simultaneous replacement of all tasks during updates.

### Task Definition

Task definitions specify the execution parameters for your containerized application. Tasks are the actual running instances of your application (similar to containers launched via docker compose) and represent the deployable unit that receives resources like ephemeral IP addresses:

- Container image references (repository URI, tag/digest)
- Resource allocations (CPU units, memory limits)
- Port mappings and network configuration
- Environment variables and parameter references
- Volume mounts and persistent storage configuration
- IAM role assignments

**Task networking characteristics**:

- Container-to-container communication within a task occurs via localhost interface
- Tasks receive dynamically allocated private IP addresses within the VPC subnet
- Each task maintains an isolated network namespace for container workloads

**Security consideration**: Tasks deployed in private subnets require either NAT gateway or VPC endpoint connectivity to pull container images from ECR.

### Task IAM Roles

Implement least privilege security model when defining IAM roles for task execution:

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "ecr:GetAuthorizationToken",
        "ecr:BatchCheckLayerAvailability",
        "ecr:GetDownloadUrlForLayer",
        "ecr:BatchGetImage"
      ],
      "Resource": "*"
    }
  ]
}
```

**Security distinction**: The task execution role is distinct from the task role. The execution role authorizes ECS container agent operations (image pulling, log delivery), while the task role provides permissions for application code execution.

## Cloud Map

Cloud Map provides service discovery functionality to address the ephemeral nature of container IP addressing, enabling API Gateway to locate and route traffic to the appropriate service instances.

### Service Discovery Configuration

1. Create a namespace (typically a private DNS namespace for internal service resolution)
2. Define service records within the namespace
3. Configure ECS services to register instances automatically with Cloud Map

## API Gateway

API Gateway serves as the unified entry point for external clients accessing internal services.

### Implementation Procedure

1. Create VPC link

   - Required for services deployed in private subnets
   - Establishes secure connectivity between API Gateway and VPC resources

2. Configure HTTP API

   - HTTP APIs offer cost advantages over REST APIs for most implementation scenarios
   - Route configuration:
     - Path parameters with proxy integration (`{proxy+}`): [AWS Documentation](https://docs.aws.amazon.com/apigateway/latest/developerguide/api-gateway-set-up-simple-proxy.html)
     - Default routes: Use `$default` or `/path/{proxy+}` for catch-all routing
   - Integration configuration:
     - Select "Private resource" integration type for VPC resource access
     - Specify Cloud Map namespace and service for target discovery
     - Associate with the appropriate VPC link

3. Configure custom domain with ACM certificate
   - Select "Public" domain visibility
   - For Cloudflare-managed domains with proxying enabled (orange cloud icon), ensure SSL/TLS encryption mode is at minimum "Full (strict)" through Cloudflare dashboard [SSL/TLS] > [Overview]
   - Define API mappings to the configured HTTP API

**Regional compatibility note**: ACM certificates must be provisioned in the same region as the API Gateway deployment.

### Rate Limiting Configuration

Implement throttling and quotas to protect backend services:

- Define usage plans for different client profiles
- Configure request rate limits and burst capacity
- Implement API key management for client identification

### Troubleshooting Methodology

1. Container image pull failures:

   - Validate ECR IAM permissions for task execution role
   - Verify VPC endpoint or NAT gateway configuration

2. Inter-service communication issues:

   - Examine security group ingress/egress rules
   - Verify service discovery registration and DNS resolution

3. API Gateway connectivity failures:
   - Validate VPC link configuration and status
   - Verify Cloud Map service instance registration
