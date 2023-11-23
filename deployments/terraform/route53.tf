resource "aws_route53_zone" "hero_ticket_zone" {
  name = "heroticket.xyz"

  tags = merge(
    var.common_tags,
    {
      "Name" = "Hero Ticket Hosted Zone"
  })
}

resource "aws_route53_record" "hero_ticket_issuer_node_record" {
  zone_id = aws_route53_zone.hero_ticket_zone.zone_id
  name    = "issuer.${aws_route53_zone.hero_ticket_zone.name}"
  type    = "A"

  alias {
    name                   = aws_alb.hero_ticket_alb.dns_name
    zone_id                = aws_alb.hero_ticket_alb.zone_id
    evaluate_target_health = false
  }
}
