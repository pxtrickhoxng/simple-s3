package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/spf13/cobra"
)

var (
	bucketName string
	region string
	profile string
)

var rootCmd = &cobra.Command{
	Use:   "simple-s3",
	Short: "A simple CLI for basic S3 operations",
}

var bucketCmd = &cobra.Command{
	Use:   "bucket",
	Short: "Set the active S3 bucket",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		bucketName, _ = cmd.Flags().GetString("name")
		region, _ = cmd.Flags().GetString("region")
		profile, _ = cmd.Flags().GetString("profile")

		if cmd.Name() != "list" {
			if bucketName == "" || region == "" {
				return fmt.Errorf("bucket and region must be set via flags")
			}
			fmt.Println("Using bucket:", bucketName, "in", region)
			if profile != "" {
				fmt.Println("Using AWS profile:", profile)
			}
		}

		return nil
	},
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create the currently selected S3 bucket",
	Run: func(cmd *cobra.Command, args []string) {
		if bucketName == "" || region == "" {
			fmt.Println("Bucket and region must be set first with the 'bucket' command")
			return
		}

		cfgOpts := []func(*config.LoadOptions) error{}
		if profile != "" {
			cfgOpts = append(cfgOpts, config.WithSharedConfigProfile(profile))
		}
		cfgOpts = append(cfgOpts, config.WithRegion(region))

		cfg, err := config.LoadDefaultConfig(context.TODO(), cfgOpts...)
		if err != nil {
			fmt.Println("Error loading AWS config:", err)
			return
		}

		input := &s3.CreateBucketInput{
			Bucket: &bucketName,
		}
		
		// By default, S3 API sets region as us-east-1; This handles cases where users specify another region
		if region != "us-east-1" {
			input.CreateBucketConfiguration = &types.CreateBucketConfiguration{
				LocationConstraint: types.BucketLocationConstraint(region),
			}
		}

		client := s3.NewFromConfig(cfg)

		_, err = client.CreateBucket(context.TODO(), input)

		if err != nil {
			fmt.Println("Error creating bucket:", err)
			return
		}

		fmt.Println("Bucket created successfully:", bucketName)
	},
}

var deleteCmd = &cobra.Command{
    Use:   "delete",
    Short: "Delete the currently selected S3 bucket",
    Run: func(cmd *cobra.Command, args []string) {
        if bucketName == "" || region == "" {
            fmt.Println("Bucket and region must be set first with the 'bucket' command")
            return
        }

        cfgOpts := []func(*config.LoadOptions) error{}
        if profile != "" {
            cfgOpts = append(cfgOpts, config.WithSharedConfigProfile(profile))
        }
        cfgOpts = append(cfgOpts, config.WithRegion(region))

        cfg, err := config.LoadDefaultConfig(context.TODO(), cfgOpts...)
        if err != nil {
            fmt.Println("Error loading AWS config:", err)
            return
        }

        client := s3.NewFromConfig(cfg)

        input := &s3.DeleteBucketInput{
            Bucket: &bucketName,
        }

        _, err = client.DeleteBucket(context.TODO(), input)
        if err != nil {
            fmt.Println("Error deleting bucket:", err)
            return
        }

        fmt.Println("Bucket deleted successfully:", bucketName)
    },
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all S3 buckets in the account",
	Run: func(cmd *cobra.Command, args []string) {
		cfgOpts := []func(*config.LoadOptions) error{}
		if profile != "" {
			cfgOpts = append(cfgOpts, config.WithSharedConfigProfile(profile))
		}

		cfg, err := config.LoadDefaultConfig(context.TODO(), cfgOpts...)
		if err != nil {
			fmt.Println("Error loading AWS config:", err)
			return
		}

		client := s3.NewFromConfig(cfg)

		output, err := client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
		if err != nil {
			fmt.Println("Error listing buckets:", err)
			return
		}

		fmt.Println("Buckets:")
		for _, b := range output.Buckets {
			fmt.Printf(" - %s (created: %v)\n", *b.Name, *b.CreationDate)
		}
	},
}

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show information about the currently selected bucket",
	Run: func(cmd *cobra.Command, args []string) {
		if bucketName == "" || region == "" {
			fmt.Println("Bucket and region must be set via flags")
			return
		}

		cfgOpts := []func(*config.LoadOptions) error{}
		if profile != "" {
			cfgOpts = append(cfgOpts, config.WithSharedConfigProfile(profile))
		}
		cfgOpts = append(cfgOpts, config.WithRegion(region))

		cfg, err := config.LoadDefaultConfig(context.TODO(), cfgOpts...)
		if err != nil {
			fmt.Println("Error loading AWS config:", err)
			return
		}

		client := s3.NewFromConfig(cfg)

		loc, err := client.GetBucketLocation(context.TODO(), &s3.GetBucketLocationInput{
			Bucket: &bucketName,
		})
		if err != nil {
			fmt.Println("Error getting bucket location:", err)
			return
		}

		regionStr := string(loc.LocationConstraint)
		if regionStr == "" {
			regionStr = "us-east-1" // AWS returns empty string for us-east-1
		}

		fmt.Printf("Bucket: %s\nRegion: %s\n", bucketName, regionStr)
	},
}

var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload a file to the currently selected S3 bucket",
	Run: func(cmd *cobra.Command, args []string) {
		if bucketName == "" || region == "" {
			fmt.Println("Bucket and region must be set via flags")
			return
		}

		filePath, _ := cmd.Flags().GetString("file")
		key, _ := cmd.Flags().GetString("key")

		if filePath == "" || key == "" {
			fmt.Println("You must provide both --file and --key")
			return
		}

		cfgOpts := []func(*config.LoadOptions) error{}
		if profile != "" {
			cfgOpts = append(cfgOpts, config.WithSharedConfigProfile(profile))
		}
		cfgOpts = append(cfgOpts, config.WithRegion(region))

		cfg, err := config.LoadDefaultConfig(context.TODO(), cfgOpts...)
		if err != nil {
			fmt.Println("Error loading AWS config:", err)
			return
		}

		client := s3.NewFromConfig(cfg)

		file, err := os.Open(filePath)
		if err != nil {
			fmt.Println("Error opening file:", err)
			return
		}
		defer file.Close()

		_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
			Bucket: &bucketName,
			Key:    &key,
			Body:   file,
		})
		if err != nil {
			fmt.Println("Error uploading file:", err)
			return
		}

		fmt.Println("File uploaded successfully:", key)
	},
}

var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download a file from the currently selected S3 bucket",
	Run: func(cmd *cobra.Command, args []string) {
		if bucketName == "" || region == "" {
			fmt.Println("Bucket and region must be set via flags")
			return
		}

		key, _ := cmd.Flags().GetString("key")
		output, _ := cmd.Flags().GetString("output")

		if key == "" || output == "" {
			fmt.Println("You must provide both --key and --output")
			return
		}

		cfgOpts := []func(*config.LoadOptions) error{}
		if profile != "" {
			cfgOpts = append(cfgOpts, config.WithSharedConfigProfile(profile))
		}
		cfgOpts = append(cfgOpts, config.WithRegion(region))

		cfg, err := config.LoadDefaultConfig(context.TODO(), cfgOpts...)
		if err != nil {
			fmt.Println("Error loading AWS config:", err)
			return
		}

		client := s3.NewFromConfig(cfg)

		outputObj, err := client.GetObject(context.TODO(), &s3.GetObjectInput{
			Bucket: &bucketName,
			Key:    &key,
		})
		if err != nil {
			fmt.Println("Error downloading file:", err)
			return
		}
		defer outputObj.Body.Close()

		outFile, err := os.Create(output)
		if err != nil {
			fmt.Println("Error creating output file:", err)
			return
		}
		defer outFile.Close()

		_, err = io.Copy(outFile, outputObj.Body)
		if err != nil {
			fmt.Println("Error saving file:", err)
			return
		}

		fmt.Println("File downloaded successfully:", output)
	},
}

var listObjectsCmd = &cobra.Command{
	Use:   "list-objects",
	Short: "List objects in the currently selected S3 bucket",
	Run: func(cmd *cobra.Command, args []string) {
		if bucketName == "" || region == "" {
			fmt.Println("Bucket and region must be set via flags")
			return
		}

		prefix, _ := cmd.Flags().GetString("prefix")

		cfgOpts := []func(*config.LoadOptions) error{}
		if profile != "" {
			cfgOpts = append(cfgOpts, config.WithSharedConfigProfile(profile))
		}
		cfgOpts = append(cfgOpts, config.WithRegion(region))

		cfg, err := config.LoadDefaultConfig(context.TODO(), cfgOpts...)
		if err != nil {
			fmt.Println("Error loading AWS config:", err)
			return
		}

		client := s3.NewFromConfig(cfg)

		output, err := client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
			Bucket: &bucketName,
			Prefix: &prefix,
		})
		if err != nil {
			fmt.Println("Error listing objects:", err)
			return
		}

		if len(output.Contents) == 0 {
			fmt.Println("No objects found in bucket")
			return
		}

		fmt.Println("Objects:")
		for _, obj := range output.Contents {
			fmt.Printf(" - %s (size: %d)\n", *obj.Key, obj.Size)
		}
	},
}

var deleteObjectCmd = &cobra.Command{
	Use:   "delete-object",
	Short: "Delete an object from the currently selected S3 bucket",
	Run: func(cmd *cobra.Command, args []string) {
		if bucketName == "" || region == "" {
			fmt.Println("Bucket and region must be set via flags")
			return
		}

		key, _ := cmd.Flags().GetString("key")
		if key == "" {
			fmt.Println("You must provide --key to delete an object")
			return
		}

		cfgOpts := []func(*config.LoadOptions) error{}
		if profile != "" {
			cfgOpts = append(cfgOpts, config.WithSharedConfigProfile(profile))
		}
		cfgOpts = append(cfgOpts, config.WithRegion(region))

		cfg, err := config.LoadDefaultConfig(context.TODO(), cfgOpts...)
		if err != nil {
			fmt.Println("Error loading AWS config:", err)
			return
		}

		client := s3.NewFromConfig(cfg)

		_, err = client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
			Bucket: &bucketName,
			Key:    &key,
		})
		if err != nil {
			fmt.Println("Error deleting object:", err)
			return
		}

		fmt.Println("Object deleted successfully:", key)
	},
}

func init() {
	bucketCmd.PersistentFlags().StringVar(&bucketName, "name", "", "S3 bucket name (required)")
	bucketCmd.PersistentFlags().StringVar(&region, "region", "", "S3 region (required)")
	bucketCmd.PersistentFlags().StringVar(&profile, "profile", "", "AWS profile (optional)")

	uploadCmd.Flags().String("file", "", "Local file path to upload (required)")
	uploadCmd.Flags().String("key", "", "S3 object key / filename (required)")

	downloadCmd.Flags().String("key", "", "S3 object key / filename to download (required)")
	downloadCmd.Flags().String("output", "", "Local path to save the file (required)")

	listObjectsCmd.Flags().String("prefix", "", "Filter objects by prefix (optional)")

	deleteObjectCmd.Flags().String("key", "", "S3 object key / filename to delete (required)")
}

func main() {
	rootCmd.AddCommand(bucketCmd)

	bucketCmd.AddCommand(createCmd)
	bucketCmd.AddCommand(deleteCmd)
	bucketCmd.AddCommand(listCmd)
	bucketCmd.AddCommand(infoCmd)
	bucketCmd.AddCommand(uploadCmd)
	bucketCmd.AddCommand(downloadCmd)
	bucketCmd.AddCommand(listObjectsCmd)
	bucketCmd.AddCommand(deleteObjectCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
