resource "aws_route53_zone" "hero_ticket_zone" {
  name = "heroticket.xyz"

  tags = merge(
    var.common_tags,
    {
      "Name" = "Hero Ticket Hosted Zone"
  })
}

# resource "aws_route53_record" "www" {
#   zone_id = aws_route53_zone.hero_ticket_zone.zone_id
#   name    = "www.${aws_route53_zone.hero_ticket_zone.name}"
#   type    = "A"
# }
