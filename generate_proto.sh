protoc -I proto/ proto/account/v1/user.proto --go_out=paths=source_relative:.
protoc -I proto/ proto/buyer/v1/buyer.proto --go_out=paths=source_relative:.
protoc -I proto/ proto/seller/v1/seller.proto --go_out=paths=source_relative:.
protoc -I proto/ proto/transporter/v1/transporter.proto --go_out=paths=source_relative:.