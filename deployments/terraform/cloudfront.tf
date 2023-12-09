resource "aws_s3_bucket" "hero_ticket_client_bucket" {
  bucket = "heroticket.xyz"

  tags = merge(
    var.common_tags,
    {
      Name = "Hero Ticket Client Bucket"
    }
  )
}

resource "aws_s3_bucket_website_configuration" "hero_ticket_client_bucket_website_configuration" {
  bucket = aws_s3_bucket.hero_ticket_client_bucket.id

  index_document {
    suffix = "index.html"
  }

  error_document {
    key = "index.html"
  }
}

resource "aws_s3_bucket_ownership_controls" "hero_ticket_client_bucket_ownership_controls" {
  bucket = aws_s3_bucket.hero_ticket_client_bucket.id

  rule {
    object_ownership = "BucketOwnerPreferred"
  }
}

resource "aws_s3_bucket_public_access_block" "hero_ticket_client_bucket_public_access_block" {
  bucket = aws_s3_bucket.hero_ticket_client_bucket.id

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

resource "aws_s3_bucket_policy" "allow_cloudfront_access_to_s3_bucket" {
  bucket = aws_s3_bucket.hero_ticket_client_bucket.id
  policy = data.aws_iam_policy_document.allow_cloudfront_access_to_s3_bucket.json
}

data "aws_iam_policy_document" "allow_cloudfront_access_to_s3_bucket" {
  statement {
    actions = ["s3:GetObject"]
    resources = [
      aws_s3_bucket.hero_ticket_client_bucket.arn,
      "${aws_s3_bucket.hero_ticket_client_bucket.arn}/*",
    ]

    principals {
      type        = "AWS"
      identifiers = [aws_cloudfront_origin_access_identity.s3_bucket_origin_access_identity.iam_arn]
    }
  }
}

resource "aws_cloudfront_origin_access_identity" "s3_bucket_origin_access_identity" {
  comment = "Access identity for S3 bucket"
}

resource "aws_cloudfront_distribution" "hero_ticket_client_distribution" {
  aliases             = ["heroticket.xyz"]
  enabled             = true
  is_ipv6_enabled     = true
  default_root_object = "index.html"
  price_class         = "PriceClass_200"

  origin {
    domain_name = aws_s3_bucket.hero_ticket_client_bucket.bucket_regional_domain_name
    origin_id   = aws_s3_bucket.hero_ticket_client_bucket.id

    custom_origin_config {
      http_port              = "80"
      https_port             = "443"
      origin_protocol_policy = "http-only"
      origin_ssl_protocols   = ["TLSv1", "TLSv1.1", "TLSv1.2"]
    }
  }

  default_cache_behavior {
    allowed_methods  = ["GET", "HEAD"]
    cached_methods   = ["GET", "HEAD"]
    target_origin_id = aws_s3_bucket.hero_ticket_client_bucket.id

    forwarded_values {
      query_string = false

      cookies {
        forward = "none"
      }
    }

    compress               = true
    viewer_protocol_policy = "redirect-to-https"
    min_ttl                = 0
    default_ttl            = 3600
    max_ttl                = 86400
  }

  restrictions {
    geo_restriction {
      restriction_type = "none"
    }
  }

  viewer_certificate {
    acm_certificate_arn = aws_acm_certificate.hero_ticket_cert2.arn
    ssl_support_method  = "sni-only"
  }

  tags = merge(
    var.common_tags,
    {
      Name = "Hero Ticket Client Distribution"
    }
  )
}
