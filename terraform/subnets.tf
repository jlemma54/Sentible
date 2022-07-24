resource "aws_subnet" "sentible-prod-public-1" {
    vpc_id = "${aws_vpc.sentible-prod-vpc.id}"
    cidr_block = "10.0.0.0/24"
    map_public_ip_on_launch = true
    availability_zone = "us-east-1a"


    tags = {
        Name = "sentible-prod-public-1"
    }
}

resource "aws_subnet" "sentible-prod-private-1" {
    vpc_id = "${aws_vpc.sentible-prod-vpc.id}"
    cidr_block = "10.0.1.0/24"
    map_public_ip_on_launch = false
    availability_zone = "us-east-1a"


    tags = {
        Name = "sentible-prod-private-1"
    }

}


resource "aws_subnet" "sentible-prod-private-2" {
    vpc_id = "${aws_vpc.sentible-prod-vpc.id}"
    cidr_block = "10.0.2.0/24"
    map_public_ip_on_launch = false
    availability_zone = "us-east-1a"
}

resource "aws_subnet" "spare-subnet" {
    vpc_id = "${aws_vpc.sentible-prod-vpc.id}"
    cidr_block = "10.0.3.0/24"
    map_public_ip_on_launch = false
    availability_zone = "us-east-1b" 
}

resource "aws_subnet" "spare-public" {
    vpc_id = "${aws_vpc.sentible-prod-vpc.id}"
    cidr_block = "10.0.4.0/24"
    map_public_ip_on_launch = true
    availability_zone = "us-east-1b"  

    tags = {
        Name = "Spare public subnet group"
    }
}

resource "aws_db_subnet_group" "rds-private-subnet" {
  name       = "rds-private-subnet"
  subnet_ids = ["${aws_subnet.sentible-prod-public-1.id}", "${aws_subnet.spare-public.id}"]

  tags = {
    Name = "My DB subnet group"
  }
}

