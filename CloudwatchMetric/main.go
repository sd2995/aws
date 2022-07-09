package main

import (
	"math/rand"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"

	"fmt"
)

func main() {
	// Initialize a session that the SDK uses to load
	// credentials from the shared credentials file ~/.aws/credentials
	// and configuration from the shared configuration file ~/.aws/config.
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Config: aws.Config{
			// ...
			Region: aws.String("AWS_REGION"),
			Credentials: credentials.NewStaticCredentials(
				"AWS_ACCESS_KEY",
				"AWS_SECRET_KEY",
				"",
			),
		},
	}))

	// Create new cloudwatch client.
	svc := cloudwatch.New(sess)

	_, err := svc.PutMetricData(&cloudwatch.PutMetricDataInput{
		Namespace: aws.String("MyAppTelemetry"),
		MetricData: []*cloudwatch.MetricDatum{
			&cloudwatch.MetricDatum{
				MetricName: aws.String("CustomMetric"),
				Unit:       aws.String("Count"),
				Value:      aws.Float64(rand.Float64() * 10),
				Dimensions: []*cloudwatch.Dimension{
					&cloudwatch.Dimension{
						Name:  aws.String("mymobileapp"),
						Value: aws.String("abc.xyz"),
					},
				},
			},
		},
	})
	if err != nil {
		fmt.Println("Error adding metrics:", err.Error())
		return
	}

	// Get information about metrics
	result, err := svc.ListMetrics(&cloudwatch.ListMetricsInput{
		Namespace: aws.String("MyAppTelemetry"),
	})
	if err != nil {
		fmt.Println("Error getting metrics:", err.Error())
		return
	}

	for _, metric := range result.Metrics {
		fmt.Println(*metric.MetricName)

		for _, dim := range metric.Dimensions {
			fmt.Println(*dim.Name+":", *dim.Value)
			fmt.Println()
		}
	}
}
