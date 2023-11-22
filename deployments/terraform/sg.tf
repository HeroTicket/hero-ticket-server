resource "aws_security_group" "hero_ticket_alb_sg" {
  name        = "hero-ticket-alb-sg"
  description = "Hero Ticket ALB Security Group"
  vpc_id      = aws_vpc.hero_ticket_vpc.id

  ingress {
    description = "Allow HTTP inbound traffic"
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    description = "Allow HTTPS inbound traffic"
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    description = "Allow all outbound traffic"
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = merge(
    var.common_tags,
    {
      "Name" = "Hero Ticket ALB Security Group"
    }
  )
}

resource "aws_security_group" "hero_ticket_issuer_node_sg" {
  name        = "hero-ticket-issuer-node-sg"
  description = "Hero Ticket Issuer Node Security Group"
  vpc_id      = aws_vpc.hero_ticket_vpc.id

  ingress {
    description = "Allow SSH inbound traffic from Manager"
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = [var.manager_ip]
  }

  ingress {
    description     = "Allow 3001 inbound traffic from ALB"
    from_port       = 3001
    to_port         = 3001
    protocol        = "tcp"
    security_groups = [aws_security_group.hero_ticket_alb_sg.id]
  }

  egress {
    description = "Allow all outbound traffic"
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_security_group" "hero_ticket_server_sg" {
  name        = "hero-ticket-server-sg"
  description = "Hero Ticket Server Security Group"
  vpc_id      = aws_vpc.hero_ticket_vpc.id

  ingress {
    description = "Allow SSH inbound traffic from Manager"
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = [var.manager_ip]
  }

  ingress {
    description     = "Allow 8080 inbound traffic from ALB"
    from_port       = 8080
    to_port         = 8080
    protocol        = "tcp"
    security_groups = [aws_security_group.hero_ticket_alb_sg.id]
  }

  ingress {
    description     = "Allow 1323 inbound traffic from Issuer Node"
    from_port       = 1323
    to_port         = 1323
    protocol        = "tcp"
    security_groups = [aws_security_group.hero_ticket_issuer_node_sg.id]
  }

  egress {
    description = "Allow all outbound traffic"
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = merge(
    var.common_tags,
    {
      "Name" = "Hero Ticket Server Security Group"
    }
  )
}
