#######################################
# EKS クラスター用 IAM ロールとポリシー
#######################################

# EKS クラスターが操作するための IAM ロールを定義
resource "aws_iam_role" "eks_cluster_role" {
  name = "eks-cluster-role"

  # EKS がこのロールを引き受けるための信頼ポリシー
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action    = "sts:AssumeRole"
      Effect    = "Allow"
      Principal = {
        Service = "eks.amazonaws.com" # EKS サービスに対してのみ許可
      }
    }]
  })
}

# EKS クラスターに必要な基本ポリシーをアタッチ
resource "aws_iam_role_policy_attachment" "eks_cluster_policy" {
  role       = aws_iam_role.eks_cluster_role.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKSClusterPolicy"
}

#######################################
# EKS ノードグループ用 IAM ロールとポリシー
#######################################

# EC2（EKSノード）が使用する IAM ロールを定義
resource "aws_iam_role" "eks_node_role" {
  name = "eks-node-group-role"

  # EC2 インスタンスにこのロールを引き受ける許可を付与
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action    = "sts:AssumeRole"
      Effect    = "Allow"
      Principal = {
        Service = "ec2.amazonaws.com" # EC2 サービス（EKSノード）に対して許可
      }
    }]
  })
}

# EKSノード（EC2）がクラスターに参加するためのポリシー
resource "aws_iam_role_policy_attachment" "node_AmazonEKSWorkerNodePolicy" {
  role       = aws_iam_role.eks_node_role.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKSWorkerNodePolicy"
}

# コンテナイメージの取得（ECR読み取り）を許可
resource "aws_iam_role_policy_attachment" "node_AmazonEC2ContainerRegistryReadOnly" {
  role       = aws_iam_role.eks_node_role.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly"
}

# CNI（Container Network Interface）プラグインの実行権限を付与
resource "aws_iam_role_policy_attachment" "node_AmazonEKS_CNI_Policy" {
  role       = aws_iam_role.eks_node_role.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKS_CNI_Policy"
}
