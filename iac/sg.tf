#######################################
# Security Group
#######################################

# ALB
resource "aws_security_group" "alb_sg" {
  name = "share-basket-alb-sg"
  description = "Allow inbound HTTP/HTTPS traffic from internet"
  vpc_id = aws_vpc.main.id

  ingress {
    description = "Allow HTTP"
    from_port = 80
    to_port = 80
    protocol = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    description = "Allow HTTPS"
    from_port = 443
    to_port = 443
    protocol = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port = 0
    to_port = 0
    protocol = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "share-basket-alb-sg"
  }
}

# EKS
resource "aws_security_group" "eks_node_sg" {
  name = "share-basket-eks-node-sg"
  description = "Allow traffic from ALB to EKS nodes"
  vpc_id = aws_vpc.main.id

  ingress {
    description = "Allow HTTP from ALB"
    from_port = 80
    to_port = 80
    protocol = "tcp"
    security_groups = [aws_security_group.alb_sg.id]
  }

  ingress {
    description = "Node to Node communication"
    from_port = 0
    to_port = 65535
    protocol = "tcp"
    self = true
  }

  # 一旦広く許可しておいて、後から閉じていく想定
  egress {
    from_port = 0
    to_port =  0
    protocol = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "share-basket-eks-node-sg"
  }
}
