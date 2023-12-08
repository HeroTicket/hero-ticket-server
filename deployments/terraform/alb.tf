resource "aws_alb" "hero_ticket_alb" {
  name               = "hero-ticket-alb"
  internal           = false
  load_balancer_type = "application"
  security_groups = [
    aws_security_group.hero_ticket_alb_sg.id,
  ]
  subnets                          = aws_subnet.hero_ticket_public_subnets[*].id
  enable_cross_zone_load_balancing = true
  idle_timeout                     = 300

  tags = merge(
    var.common_tags,
    {
      Name = "Hero Ticket ALB"
    }
  )
}

resource "aws_alb_target_group" "hero_ticket_issuer_node_target_group" {
  name     = "alb-target-issue-node"
  port     = 3001
  protocol = "HTTP"
  vpc_id   = aws_vpc.hero_ticket_vpc.id

  health_check {
    path                = "/status"
    protocol            = "HTTP"
    interval            = 300
    timeout             = 120
    healthy_threshold   = 5
    unhealthy_threshold = 10
    matcher             = "200-399"
  }

  stickiness {
    enabled = true
    type    = "lb_cookie"
  }

  tags = merge(
    var.common_tags,
    {
      Name = "Hero Ticket Issuer Node Target Group"
    }
  )
}

resource "aws_alb_target_group" "hero_ticket_server_target_group" {
  name     = "alb-target-server"
  port     = 8080
  protocol = "HTTP"
  vpc_id   = aws_vpc.hero_ticket_vpc.id

  health_check {
    path                = "/status"
    protocol            = "HTTP"
    interval            = 300
    timeout             = 120
    healthy_threshold   = 5
    unhealthy_threshold = 10
    matcher             = "200-399"
  }

  stickiness {
    enabled = true
    type    = "lb_cookie"
  }

  tags = merge(
    var.common_tags,
    {
      Name = "Hero Ticket Server Target Group"
    }
  )
}

resource "aws_lb_listener" "hero_ticket_http_listener" {
  load_balancer_arn = aws_alb.hero_ticket_alb.arn
  port              = 80
  protocol          = "HTTP"

  default_action {
    type = "redirect"
    redirect {
      port        = "443"
      protocol    = "HTTPS"
      status_code = "HTTP_301"
    }
  }

  tags = merge(
    var.common_tags,
    {
      Name = "Hero Ticket HTTP Listener"
    }
  )
}

resource "aws_lb_listener" "hero_ticket_https_listener" {
  load_balancer_arn = aws_alb.hero_ticket_alb.arn
  port              = 443
  protocol          = "HTTPS"
  ssl_policy        = "ELBSecurityPolicy-2016-08"
  certificate_arn   = aws_acm_certificate.hero_ticket_cert.arn

  default_action {
    type = "redirect"
    redirect {
      port        = "443"
      protocol    = "HTTPS"
      status_code = "HTTP_301"
      host        = aws_route53_zone.hero_ticket_zone.name
    }
  }

  tags = merge(
    var.common_tags,
    {
      Name = "Hero Ticket HTTPS Listener"
    }
  )
}

resource "aws_lb_listener_rule" "hero_ticket_alb_listener_rule_https_issuer_node" {
  listener_arn = aws_lb_listener.hero_ticket_https_listener.arn
  priority     = 2

  action {
    type             = "forward"
    target_group_arn = aws_alb_target_group.hero_ticket_issuer_node_target_group.arn
  }

  condition {
    host_header {
      values = [
        aws_route53_record.hero_ticket_issuer_node_record.name,
      ]
    }
  }

  tags = merge(
    var.common_tags,
    {
      Name = "Hero Ticket ALB Listener Rule HTTPS Issuer Node"
    }
  )
}

resource "aws_lb_listener_rule" "hero_ticket_alb_listener_rule_https_server" {
  listener_arn = aws_lb_listener.hero_ticket_https_listener.arn
  priority     = 1

  action {
    type             = "forward"
    target_group_arn = aws_alb_target_group.hero_ticket_server_target_group.arn
  }

  condition {
    host_header {
      values = [
        aws_route53_record.hero_ticket_server_record.name,
      ]
    }
  }

  tags = merge(
    var.common_tags,
    {
      Name = "Hero Ticket ALB Listener Rule HTTPS Server"
    }
  )
}
