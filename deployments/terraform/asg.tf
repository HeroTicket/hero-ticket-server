resource "aws_autoscaling_group" "hero_ticket_issuer_node_asg" {
  name             = "hero-ticket-issuer-node-asg"
  max_size         = 1
  min_size         = 1
  desired_capacity = 1

  launch_template {
    id      = aws_launch_template.hero_ticket_issuer_node_template.id
    version = "$Latest"
  }

  health_check_type         = "ELB"
  health_check_grace_period = 300

  vpc_zone_identifier = [
    aws_subnet.hero_ticket_public_subnet[0].id,
    aws_subnet.hero_ticket_public_subnet[1].id,
  ]

  target_group_arns = [

  ]

  tag {
    key                 = "Name"
    value               = "Hero Ticket Issuer Node ASG"
    propagate_at_launch = true
  }
}

# resource "aws_autoscaling_group" "hero_ticket_server_asg" {
#   name             = "hero-ticket-server-asg"
#   max_size         = 1
#   min_size         = 1
#   desired_capacity = 1

#   launch_template {
#     id      = aws_launch_template.hero_ticket_server_template.id
#     version = "$Latest"
#   }

#   health_check_type         = "ELB"
#   health_check_grace_period = 300

#   vpc_zone_identifier = [
#     aws_subnet.hero_ticket_public_subnet[0].id,
#     aws_subnet.hero_ticket_public_subnet[1].id,
#   ]

#   target_group_arns = [

#   ]

#   tag {
#     key                 = "Name"
#     value               = "Hero Ticket Server ASG"
#     propagate_at_launch = true
#   }
# }
