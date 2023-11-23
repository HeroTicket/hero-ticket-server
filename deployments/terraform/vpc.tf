resource "aws_vpc" "hero_ticket_vpc" {
  cidr_block = var.vpc_cidr_block
  tags = merge(
    var.common_tags,
    {
      "Name" = "Hero Ticket VPC"
    }
  )
}

resource "aws_subnet" "hero_ticket_public_subnets" {
  count = length(var.public_subnet_cidr_blocks)

  vpc_id            = aws_vpc.hero_ticket_vpc.id
  cidr_block        = var.public_subnet_cidr_blocks[count.index]
  availability_zone = var.availability_zones[count.index]

  tags = merge(
    var.common_tags,
    {
      "Name" = "Hero Ticket Public Subnet ${count.index + 1}"
    }
  )
}

resource "aws_subnet" "hero_ticket_private_subnets" {
  count = length(var.private_subnet_cidr_blocks)

  vpc_id            = aws_vpc.hero_ticket_vpc.id
  cidr_block        = var.private_subnet_cidr_blocks[count.index]
  availability_zone = var.availability_zones[count.index]

  tags = merge(
    var.common_tags,
    {
      "Name" = "Hero Ticket Private Subnet ${count.index + 1}"
    }
  )
}

resource "aws_internet_gateway" "hero_ticket_igw" {
  vpc_id = aws_vpc.hero_ticket_vpc.id

  tags = merge(
    var.common_tags,
    {
      "Name" = "Hero Ticket Internet Gateway"
    }
  )
}

resource "aws_route_table" "hero_ticket_public_route_table" {
  vpc_id = aws_vpc.hero_ticket_vpc.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.hero_ticket_igw.id
  }

  tags = merge(
    var.common_tags,
    {
      "Name" = "Hero Ticket Public Route Table"
    }
  )
}

resource "aws_route_table" "hero_ticket_private_route_table" {
  vpc_id = aws_vpc.hero_ticket_vpc.id

  tags = merge(
    var.common_tags,
    {
      "Name" = "Hero Ticket Private Route Table"
    }
  )
}

resource "aws_route_table_association" "hero_ticket_public_subnet_association" {
  count          = length(var.public_subnet_cidr_blocks)
  subnet_id      = aws_subnet.hero_ticket_public_subnets[count.index].id
  route_table_id = aws_route_table.hero_ticket_public_route_table.id
}

resource "aws_route_table_association" "hero_ticket_private_subnet_association" {
  count          = length(var.private_subnet_cidr_blocks)
  subnet_id      = aws_subnet.hero_ticket_private_subnets[count.index].id
  route_table_id = aws_route_table.hero_ticket_private_route_table.id
}
