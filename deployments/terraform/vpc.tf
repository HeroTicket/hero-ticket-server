resource "aws_vpc" "hero-ticket-vpc" {
  cidr_block = var.vpc_cidr_block
  tags = merge(
    var.common_tags,
    {
      "Name" = "Hero Ticket VPC"
    }
  )
}

resource "aws_subnet" "hero-ticket-public-subnet" {
  count = length(var.public_subnet_cidr_blocks)

  vpc_id            = aws_vpc.hero-ticket-vpc.id
  cidr_block        = var.public_subnet_cidr_blocks[count.index]
  availability_zone = var.availability_zones[count.index]

  tags = merge(
    var.common_tags,
    {
      "Name" = "Hero Ticket Public Subnet ${count.index + 1}"
    }
  )
}

resource "aws_subnet" "hero-ticket-private-subnet" {
  count = length(var.private_subnet_cidr_blocks)

  vpc_id            = aws_vpc.hero-ticket-vpc.id
  cidr_block        = var.private_subnet_cidr_blocks[count.index]
  availability_zone = var.availability_zones[count.index]

  tags = merge(
    var.common_tags,
    {
      "Name" = "Hero Ticket Private Subnet ${count.index + 1}"
    }
  )
}

resource "aws_internet_gateway" "hero-ticket-igw" {
  vpc_id = aws_vpc.hero-ticket-vpc.id

  tags = merge(
    var.common_tags,
    {
      "Name" = "Hero Ticket Internet Gateway"
    }
  )
}

resource "aws_route_table" "hero-ticket-public-route-table" {
  vpc_id = aws_vpc.hero-ticket-vpc.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.hero-ticket-igw.id
  }

  tags = merge(
    var.common_tags,
    {
      "Name" = "Hero Ticket Public Route Table"
    }
  )
}

resource "aws_route_table" "hero-ticket-private-route-table" {
  vpc_id = aws_vpc.hero-ticket-vpc.id

  tags = merge(
    var.common_tags,
    {
      "Name" = "Hero Ticket Private Route Table"
    }
  )
}

resource "aws_route_table_association" "hero-ticket-public-subnet-association" {
  count          = length(var.public_subnet_cidr_blocks)
  subnet_id      = aws_subnet.hero-ticket-public-subnet[count.index].id
  route_table_id = aws_route_table.hero-ticket-public-route-table.id
}

resource "aws_route_table_association" "hero-ticket-private-subnet-association" {
  count          = length(var.private_subnet_cidr_blocks)
  subnet_id      = aws_subnet.hero-ticket-private-subnet[count.index].id
  route_table_id = aws_route_table.hero-ticket-private-route-table.id
}
