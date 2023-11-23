resource "aws_iam_policy" "hero_ticket_ec2_to_ecr_policy" {
  name        = "hero_ticket_ec2_to_ecr_policy"
  path        = "/"
  description = "Policy for EC2 to ECR"

  policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "ecr:GetDownloadUrlForLayer",
                "ecr:BatchGetImage",
                "ecr:DescribeImages",
                "ecr:GetAuthorizationToken",
                "ecr:ListImages"
            ],
            "Resource": "*"
        }
    ]
}
EOF

  tags = merge(
    var.common_tags,
    {
      "Name" = "Hero Ticket EC2 to ECR Policy"
    }
  )
}

resource "aws_iam_role" "hero_ticket_ec2_to_ecr_role" {
  name = "hero_ticket_ec2_to_ecr_role"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Effect": "Allow",
      "Principal": {
        "Service": ["ec2.amazonaws.com"]
        }
      }
    ]
  }
EOF

  tags = merge(
    var.common_tags,
    {
      "Name" = "Hero Ticket EC2 to ECR Role"
    }
  )
}

resource "aws_iam_role_policy_attachment" "hero_ticket_ec2_to_ecr_role_policy_attachment" {
  role       = aws_iam_role.hero_ticket_ec2_to_ecr_role.name
  policy_arn = aws_iam_policy.hero_ticket_ec2_to_ecr_policy.arn
}

resource "aws_iam_instance_profile" "hero_ticket_ec2_to_ecr_instance_profile" {
  name = "hero_ticket_ec2_to_ecr_instance_profile"
  role = aws_iam_role.hero_ticket_ec2_to_ecr_role.name
  tags = merge(
    var.common_tags,
    {
      "Name" = "Hero Ticket EC2 to ECR Instance Profile"
    }
  )
}
