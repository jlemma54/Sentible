resource "aws_internet_gateway" "sentible-prod-igw" {
    vpc_id = "${aws_vpc.sentible-prod-vpc.id}"
    tags = {
        Name = "sentible-prod-igw"
    }
}

resource "aws_route_table" "sentible-prod-public-crt" {
    vpc_id = "${aws_vpc.sentible-prod-vpc.id}"
    
    route {

        cidr_block = "0.0.0.0/0" 
        gateway_id = "${aws_internet_gateway.sentible-prod-igw.id}" 
    }
    
    tags = {
        Name = "sentible-prod-public-crt"
    }
}

resource "aws_route_table_association" "sentible-prod-crta-public-subnet-1"{
    subnet_id = "${aws_subnet.sentible-prod-public-1.id}"
    route_table_id = "${aws_route_table.sentible-prod-public-crt.id}"
}

resource "aws_route_table_association" "sentible-spare-association" {
    subnet_id = "${aws_subnet.spare-public.id}"
    route_table_id = "${aws_route_table.sentible-prod-public-crt.id}"
}
resource "aws_security_group" "ssh-allowed" {
    name = "ssh-allowed"
    description = "standard ec2 security group"
    vpc_id = "${aws_vpc.sentible-prod-vpc.id}"
    
    egress {
        description = "allow all outbound traffic"
        from_port = 0
        to_port = 0
        protocol = -1
        cidr_blocks = ["0.0.0.0/0"]
    }
    ingress {
        from_port = 22
        to_port = 22
        protocol = "tcp"
        cidr_blocks = ["144.92.38.224/27"]
    }

    ingress {
        from_port = 80
        to_port = 80
        protocol = "tcp"
        cidr_blocks = ["144.92.38.224/27"]
    }

    ingress {
        from_port = 8080
        to_port = 8080
        protocol = "tcp"
        cidr_blocks = ["144.92.38.224/27"]
    }


   ingress {
        from_port = 443
        to_port = 443
        protocol = "tcp"
        cidr_blocks = ["144.92.38.224/27"]
    } 

    


    tags = {
        Name = "ssh-allowed"
    }
}

resource "aws_security_group" "gitlab-sg" {
    name = "gitlab-sg"
    description = "security group for gitlab runner"
    vpc_id = "${aws_vpc.sentible-prod-vpc.id}"
    
    egress {
        description = "allow all outbound traffic"
        from_port = 0
        to_port = 0
        protocol = -1
        cidr_blocks = ["0.0.0.0/0"]
    }
    ingress {
        from_port = 22
        to_port = 22
        protocol = "tcp"
        cidr_blocks = ["144.92.38.224/27"]
    }

    ingress {
        from_port = 80
        to_port = 80
        protocol = "tcp"
        cidr_blocks = ["144.92.38.224/27"]
    }

    ingress {
        from_port = 8080
        to_port = 8080
        protocol = "tcp"
        cidr_blocks = ["144.92.38.224/27"]
    }

    ingress {
        from_port = 443
        to_port = 443
        protocol = "tcp"
        cidr_blocks = ["144.92.38.224/27"]
    }

    tags = {
        Name = "ssh-allowed"
    }
}


resource "aws_security_group" "load-balancer" {
   vpc_id = "${aws_vpc.sentible-prod-vpc.id}"

   egress {
      from_port = 80
      to_port = 80
      protocol = "tcp"
      cidr_blocks = [aws_vpc.sentible-prod-vpc.cidr_block] 
   } 

   ingress {
      from_port = 80
      to_port = 80
      protocol = "tcp"
      cidr_blocks = ["0.0.0.0/0"] 
   }

   ingress {
       from_port = 80
       to_port = 8080
       protocol = "tcp"
       cidr_blocks = ["0.0.0.0/0"]
   }

}

resource "aws_security_group" "db-sg" {
    vpc_id = "${aws_vpc.sentible-prod-vpc.id}"

    # ingress {
    #     from_port = 0
    #     to_port = 0
    #     protocol = -1
    #     cidr_blocks = ["0.0.0.0/0"]
    # }

    # egress {
    #     from_port = 0
    #     to_port = 0
    #     protocol = -1
    #     cidr_blocks = ["0.0.0.0/0"]
    # }

    egress {
        from_port = 3306
        to_port = 3306
        protocol = "tcp"
        cidr_blocks = [aws_vpc.sentible-prod-vpc.cidr_block]
    }

    ingress {
        from_port = 3306
        to_port = 3306
        protocol = "tcp"
        cidr_blocks = [aws_vpc.sentible-prod-vpc.cidr_block]
    }

    ingress {
        from_port = 3306
        to_port = 3306
        protocol = "tcp"
        cidr_blocks = ["144.92.38.224/27"]
    }

    egress {
        from_port = 3306
        to_port = 3306
        protocol = "tcp"
        cidr_blocks = ["144.92.38.224/27"]
    }
}



