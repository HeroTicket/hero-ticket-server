resource "tls_private_key" "private_key" {
  algorithm = "RSA"
  rsa_bits  = 4096
}

resource "local_file" "private_key_file" {
  content         = tls_private_key.private_key.private_key_pem
  filename        = "${var.keypair_name}.pem"
  file_permission = "0400"
}

resource "aws_key_pair" "hero_ticket_keypair" {
  key_name   = var.keypair_name
  public_key = tls_private_key.private_key.public_key_openssh

  tags = merge(
    var.common_tags,
    {
      Name = "Hero Ticket Key Pair"
    }
  )
}
