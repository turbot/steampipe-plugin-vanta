package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/turbot/steampipe-plugin-vanta/restapi"
	"github.com/turbot/steampipe-plugin-vanta/restapi/model"
)

func main() {
	// Get API token from environment variable
	token := os.Getenv("VANTA_API_TOKEN")
	if token == "" {
		log.Fatal("VANTA_API_TOKEN environment variable is required")
	}

	// Create a new Vanta client
	ctx := context.Background()
	client, err := restapi.New(ctx, restapi.WithToken(token))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Example 1: List all tests with pagination
	fmt.Println("=== Example 1: List All Tests ===")
	options := &model.ListTestsOptions{
		PageSize: 10, // Limit to 10 tests per page
	}

	result, err := client.ListTests(ctx, options)
	if err != nil {
		log.Fatalf("Failed to list tests: %v", err)
	}

	fmt.Printf("Found %d tests\n", len(result.Results.Data))
	for _, test := range result.Results.Data {
		fmt.Printf("- %s: %s (Status: %s, Category: %s)\n",
			test.ID, test.Name, test.Status, test.Category)
	}

	// Example 2: Get a specific test by ID
	fmt.Println("\n=== Example 2: Get Specific Test ===")
	if len(result.Results.Data) > 0 {
		testID := result.Results.Data[0].ID
		test, err := client.GetTestByID(ctx, testID)
		if err != nil {
			log.Printf("Failed to get test %s: %v", testID, err)
		} else {
			fmt.Printf("Test: %s\n", test.Name)
			fmt.Printf("Description: %s\n", test.Description)
			fmt.Printf("Last Test Run: %v\n", test.LastTestRunDate)
			if test.Owner != nil {
				fmt.Printf("Owner: %s (%s)\n", test.Owner.DisplayName, test.Owner.EmailAddress)
			}
		}
	}

	// Example 3: List tests by category using filter
	fmt.Println("\n=== Example 3: List Tests by Category ===")
	infraOptions := &model.ListTestsOptions{
		PageSize:       5,
		CategoryFilter: "Infrastructure",
	}
	infraTests, err := client.ListTests(ctx, infraOptions)
	if err != nil {
		log.Printf("Failed to list infrastructure tests: %v", err)
	} else {
		fmt.Printf("Found %d infrastructure tests\n", len(infraTests.Results.Data))
		for _, test := range infraTests.Results.Data {
			fmt.Printf("- %s: %s\n", test.ID, test.Name)
		}
	}

	// Example 4: List failing tests using status filter
	fmt.Println("\n=== Example 4: List Failing Tests ===")
	failingOptions := &model.ListTestsOptions{
		PageSize:     5,
		StatusFilter: "NEEDS_ATTENTION",
	}
	failingTests, err := client.ListTests(ctx, failingOptions)
	if err != nil {
		log.Printf("Failed to list failing tests: %v", err)
	} else {
		fmt.Printf("Found %d failing tests\n", len(failingTests.Results.Data))
		for _, test := range failingTests.Results.Data {
			fmt.Printf("- %s: %s (Status: %s)\n", test.ID, test.Name, test.Status)
			if test.RemediationStatusInfo != nil {
				fmt.Printf("  Remediation: %s (Items: %d)\n",
					test.RemediationStatusInfo.Status, test.RemediationStatusInfo.ItemCount)
			}
		}
	}

	// Example 5: List AWS tests using integration filter
	fmt.Println("\n=== Example 5: List AWS Tests ===")
	awsOptions := &model.ListTestsOptions{
		PageSize:          5,
		IntegrationFilter: "aws",
	}
	awsTests, err := client.ListTests(ctx, awsOptions)
	if err != nil {
		log.Printf("Failed to list AWS tests: %v", err)
	} else {
		fmt.Printf("Found %d AWS tests\n", len(awsTests.Results.Data))
		for _, test := range awsTests.Results.Data {
			fmt.Printf("- %s: %s\n", test.ID, test.Name)
			fmt.Printf("  Integrations: %v\n", test.Integrations)
		}
	}

	// Example 6: List deactivated tests using status filter
	fmt.Println("\n=== Example 6: List Deactivated Tests ===")
	deactivatedOptions := &model.ListTestsOptions{
		PageSize:     3,
		StatusFilter: "DEACTIVATED",
	}
	deactivatedTests, err := client.ListTests(ctx, deactivatedOptions)
	if err != nil {
		log.Printf("Failed to list deactivated tests: %v", err)
	} else {
		fmt.Printf("Found %d deactivated tests\n", len(deactivatedTests.Results.Data))
		for _, test := range deactivatedTests.Results.Data {
			fmt.Printf("- %s: %s\n", test.ID, test.Name)
			if test.DeactivatedStatusInfo != nil && test.DeactivatedStatusInfo.DeactivatedReason != nil {
				fmt.Printf("  Reason: %s\n", *test.DeactivatedStatusInfo.DeactivatedReason)
			}
		}
	}

	// Example 7: List tests by owner
	fmt.Println("\n=== Example 7: List Tests by Owner ===")
	if len(result.Results.Data) > 0 && result.Results.Data[0].Owner != nil {
		ownerID := result.Results.Data[0].Owner.ID
		ownerOptions := &model.ListTestsOptions{
			PageSize:    5,
			OwnerFilter: ownerID,
		}
		ownerTests, err := client.ListTests(ctx, ownerOptions)
		if err != nil {
			log.Printf("Failed to list tests by owner: %v", err)
		} else {
			fmt.Printf("Found %d tests owned by %s\n", len(ownerTests.Results.Data), ownerID)
			for _, test := range ownerTests.Results.Data {
				fmt.Printf("- %s: %s\n", test.ID, test.Name)
			}
		}
	}

	// Example 8: Pagination through all tests
	fmt.Println("\n=== Example 8: Pagination Example ===")
	allTestCount := 0
	cursor := ""

	for {
		pageOptions := &model.ListTestsOptions{
			PageSize:   50,
			PageCursor: cursor,
		}

		pageResult, err := client.ListTests(ctx, pageOptions)
		if err != nil {
			log.Printf("Failed to get page: %v", err)
			break
		}

		allTestCount += len(pageResult.Results.Data)
		fmt.Printf("Page with %d tests (Total so far: %d)\n",
			len(pageResult.Results.Data), allTestCount)

		// Check if there are more pages
		if !pageResult.Results.PageInfo.HasNextPage {
			break
		}

		cursor = pageResult.Results.PageInfo.EndCursor
	}

	fmt.Printf("Total tests found: %d\n", allTestCount)

	// Example 9: Advanced filtering - passing tests in Data storage category
	fmt.Println("\n=== Example 9: Advanced Filtering ===")
	advancedOptions := &model.ListTestsOptions{
		PageSize:       10,
		StatusFilter:   "OK",
		CategoryFilter: "Data storage",
	}
	advancedTests, err := client.ListTests(ctx, advancedOptions)
	if err != nil {
		log.Printf("Failed to list tests with advanced filtering: %v", err)
	} else {
		fmt.Printf("Found %d passing data storage tests\n", len(advancedTests.Results.Data))
		for _, test := range advancedTests.Results.Data {
			fmt.Printf("- %s: %s (Status: %s, Category: %s)\n",
				test.ID, test.Name, test.Status, test.Category)
		}
	}

	// Example 10: List tests in rollout
	fmt.Println("\n=== Example 10: List Tests in Rollout ===")
	inRollout := true
	rolloutOptions := &model.ListTestsOptions{
		PageSize:    5,
		IsInRollout: &inRollout,
	}
	rolloutTests, err := client.ListTests(ctx, rolloutOptions)
	if err != nil {
		log.Printf("Failed to list rollout tests: %v", err)
	} else {
		fmt.Printf("Found %d tests in rollout\n", len(rolloutTests.Results.Data))
		for _, test := range rolloutTests.Results.Data {
			fmt.Printf("- %s: %s\n", test.ID, test.Name)
		}
	}
}
