resource "aws_launch_template" "hero_ticket_issuer_node_template" {
  name          = "hero-ticket-issuer-node-template"
  description   = "Hero Ticket Issuer Node Template"
  image_id      = "ami-0c9c942bd7bf113a2"
  instance_type = "t2.micro"
  key_name      = aws_key_pair.hero_ticket_keypair.key_name

  block_device_mappings {
    device_name = "/dev/sda1"
    ebs {
      volume_size = 30
      volume_type = "gp2"
    }
  }

  network_interfaces {
    associate_public_ip_address = true
    security_groups             = [aws_security_group.hero_ticket_issuer_node_sg.id]
  }

  user_data = filebase64("./scripts/setup_issuer_node.sh")

  tags = merge(
    var.common_tags,
    {
      "Name" = "Hero Ticket Issuer Node Template"
    }
  )
}

resource "aws_launch_template" "hero_ticket_server_template" {
  name          = "hero-ticket-server-template"
  description   = "Hero Ticket Server Template"
  image_id      = "ami-0c9c942bd7bf113a2"
  instance_type = "t2.micro"
  key_name      = aws_key_pair.hero_ticket_keypair.key_name

  block_device_mappings {
    device_name = "/dev/sda1"
    ebs {
      volume_size = 15
      volume_type = "gp2"
    }
  }

  network_interfaces {
    associate_public_ip_address = true
    security_groups             = [aws_security_group.hero_ticket_server_sg.id]
  }

  iam_instance_profile {
    name = aws_iam_instance_profile.hero_ticket_ec2_to_ecr_instance_profile.name
  }

  user_data = filebase64("./scripts/setup_server.sh")

  tags = merge(
    var.common_tags,
    {
      "Name" = "Hero Ticket Server Template"
    }
  )
}

data "aws_instances" "hero_ticket_issuer_node_instances" {
  instance_tags = {
    Name = "Hero Ticket Issuer Node ASG"
  }
  instance_state_names = ["running"]
}

data "aws_instances" "hero_ticket_server_instances" {
  instance_tags = {
    Name = "Hero Ticket Server ASG"
  }
  instance_state_names = ["running"]
}

output "hero_ticket_issuer_node_public_ips" {
  value = data.aws_instances.hero_ticket_issuer_node_instances.public_ips
}

output "hero_ticket_server_public_ips" {
  value = data.aws_instances.hero_ticket_server_instances.public_ips
}
