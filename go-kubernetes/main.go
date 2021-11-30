package main

import (
  "fmt"
  "net/http"
  "io/ioutil"
     "log"
     "path/filepath"
     "os"
     "time"
    "github.com/aws/aws-sdk-go/aws/credentials"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/s3/s3manager"
    "github.com/aws/aws-sdk-go/service/s3"
    "flag"
     
)



// Definition to upload the file to S3
func upload_to_s3(filename string , accessKey string, accessSecret string ,myBucket string ,Region string ){
    var awsConfig *aws.Config
    if accessKey == "" || accessSecret == "" {
        //load default credentials
        awsConfig = &aws.Config{
            Region: aws.String(Region),
        }
    } else {
        awsConfig = &aws.Config{
            Region:      aws.String(Region),
            Credentials: credentials.NewStaticCredentials(accessKey, accessSecret, ""),
          //  Endpoint:   aws.String(Endpoint),

        }
    }

    // The session the S3 Uploader will use
    ses , err := session.NewSession(awsConfig)
    if err != nil {
        log.Println(err)
      }
    sess := session.Must(ses , err)
    svc := s3.New(sess)
  	input := &s3.ListObjectsInput{
  			 Bucket: aws.String(myBucket),
  		}

  	 resultf, _ := svc.ListObjects(input)
  	 //w.Header().Set("Content-Type", "text/html")
		 for _, item := range resultf.Contents {
         nfile := "ipfix-data-objects/" + filename
         if nfile == *item.Key  {
           //log.Println(filename + " : already exists in given s3 bucket" )
           return
          }
		  	}
    // Create an uploader with the session and default options
    //uploader := s3manager.NewUploader(sess)
    // Create an uploader with the session and custom options
    uploader := s3manager.NewUploader(sess, func(u *s3manager.Uploader) {
        u.PartSize = 5 * 1024 * 1024 // The minimum/default allowed part size is 5MB
        u.Concurrency = 2            // default is 5
    })

    //open the file
    f, err := os.Open(filename)
    if err != nil {
        log.Printf("failed to open file %q, %v", filename, err)
        return
    }
    //defer f.Close()
    // Upload the file to S3.
    result, err := uploader.Upload(&s3manager.UploadInput{
        Bucket: aws.String(myBucket),
        Key:    aws.String("ipfix-data-objects/" + filename),
        Body:   f,
    })

    //in case it fails to upload
    if err != nil {
        log.Printf("failed to upload file, %v", err)
        return
    }
    log.Printf("file uploaded to, %s\n", result.Location)
}

func send2s3(accessKey string, accessSecret string , myBucket string ,Region string ) {

    for {
      files, err := ioutil.ReadDir("./")
      if err != nil {
          log.Fatal(err)
      }

      for _, f := range files {
      	    fileExtension := filepath.Ext(f.Name())
            fileName := f.Name()
      	    if fileExtension != ".gz" {
                  json_file_name := fileName[:len(fileName)-len(filepath.Ext(fileName))]
                  _, status := os.Stat(json_file_name)
                  if os.IsNotExist(status) {
                    upload_to_s3(f.Name(), accessKey, accessSecret, myBucket, Region)
                  }
              	}
              }
      time.Sleep( 5 * time.Minute)  // loop waits for 5 mintutes to push new files if existed
    }
}


func homePage(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "My Awesome Go App")
}

func setupRoutes() {
  http.HandleFunc("/", homePage)
}

func main() {
  fmt.Println("Go Web App Started on Port 3000")
  accessKey := flag.String("accessKey", "", "User dictionary file")
  accessSecret := flag.String("accessSecret", "", "Log IPFIX message statistics")
  myBucket := flag.String("myBucket", "", "Log traffic rates (Procera)")
  Region := flag.String("Region", "", "Display received flow records in JSON format")
  push2s3 := flag.Bool("push2s3", false , "push")
  flag.Parse()
  if (*push2s3){  send2s3(*accessKey,*accessSecret , *myBucket , *Region)}
  
  setupRoutes()
  http.ListenAndServe(":3000", nil)
}
