resource "aws_acm_certificate" "hero_ticket_cert" {
  domain_name               = aws_route53_zone.hero_ticket_zone.name
  subject_alternative_names = ["*.${aws_route53_zone.hero_ticket_zone.name}"]
  validation_method         = "DNS"

  lifecycle {
    create_before_destroy = true
  }

  tags = merge(
    var.common_tags,
    {
      "Name" = "Hero Ticket Certificate"
    }
  )
}

resource "aws_route53_record" "hero_ticket_cert_validation_record" {
  zone_id = aws_route53_zone.hero_ticket_zone.zone_id
  name    = tolist(aws_acm_certificate.hero_ticket_cert.domain_validation_options)[0].resource_record_name
  type    = tolist(aws_acm_certificate.hero_ticket_cert.domain_validation_options)[0].resource_record_type
  records = [tolist(aws_acm_certificate.hero_ticket_cert.domain_validation_options)[0].resource_record_value]
  ttl     = 60
}

resource "aws_acm_certificate_validation" "hero_ticket_cert_validation" {
  certificate_arn         = aws_acm_certificate.hero_ticket_cert.arn
  validation_record_fqdns = [aws_route53_record.hero_ticket_cert_validation_record.fqdn]
}
