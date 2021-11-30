# kuberenets-secrets

>go mod init venupammi26/go-docker
>
>go mod download
>
>go build
>
>go get github.com/aws/aws-sdk-go/aws
>
>go get github.com/aws/aws-sdk-go/aws/credentials
>
>go get github.com/aws/aws-sdk-go/aws/session
>
>go get github.com/aws/aws-sdk-go/service/s3
>
>go get github.com/aws/aws-sdk-go/service/s3/s3manager
>
>go run main.go




>docker login
>
>docker build . -t venupammi/s3-go-kube
>
>docker run -it -p 3000:3000 venupammi/s3-go-kube
>
>sudo docker run exec venupammi/s3-go-kube bash
>
>sudo docker push venupammi/s3-go-kube
>



>kubectl kustomize | kubectl apply -n default -f -



>./main --accessKey=$s3AccessKeyId --accessSecret=$s3secretKey --myBucket=$myBucket --Region=$Region --push2s3=True
